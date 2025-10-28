# Data Model: MPIJobs Support in OpenShift AI

**Feature Branch**: `001-mpijob-trainer-v2-support`
**Date**: 2025-10-28
**Status**: Phase 1 Design

This document defines the entity model and relationships for MPIJob support in OpenShift AI using KubeFlow Trainer V2.

---

## 1. Entity Definitions

### 1.1 TrainJob (MPIJob Instance)

Represents a distributed training job using Message Passing Interface for multi-node coordination.

**CRD**: `kubeflow.org/v2alpha1/TrainJob`

**Fields**:
- `metadata.name` (string, required): Unique identifier for the job within namespace
- `metadata.namespace` (string, required): Namespace for multi-tenant isolation
- `metadata.labels` (map[string]string, optional): User-defined labels for organization
- `metadata.annotations` (map[string]string, optional): System and user annotations
- `spec.runtimeRef` (RuntimeRef, required): Reference to TrainingRuntime or ClusterTrainingRuntime
  - `apiGroup`: "kubeflow.org"
  - `kind`: "TrainingRuntime" or "ClusterTrainingRuntime"
  - `name`: Runtime template name
- `spec.trainer` (TrainerSpec, required): Training configuration
  - `numNodes` (int, required): Number of distributed training nodes (1-64)
  - `image` (string, required): Docker image for training containers
  - `command` ([]string, optional): Entrypoint command override
  - `args` ([]string, optional): Command arguments
  - `resourcesPerNode` (ResourceRequirements, required): CPU, memory, GPU per node
- `spec.initializer` (InitializerSpec, optional): Automated dataset/model setup
  - `dataset.storageUri` (string): Dataset location (s3://, hf://)
  - `model.storageUri` (string): Model location
- `status` (TrainJobStatus): Current job state
  - `phase` (string): Pending, Running, Succeeded, Failed
  - `conditions` ([]Condition): Detailed status conditions
  - `startTime` (timestamp): Job start time
  - `completionTime` (timestamp): Job completion time
  - `launcherPodStatus` (PodStatus): Launcher pod state
  - `workerPodStatuses` ([]PodStatus): Worker pod states
  - `failureReason` (string): Human-readable failure reason

**Relationships**:
- References one TrainingRuntime (namespace-scoped) or ClusterTrainingRuntime (cluster-scoped)
- Owns one launcher Pod
- Owns N worker Pods (where N = spec.trainer.numNodes)
- Generates one PodGroup for gang scheduling
- Generates one JobSet to manage launcher and worker Jobs

**Validation Rules**:
- `numNodes` must be >= 1 and <= 64
- `image` must be a valid container image reference
- `resourcesPerNode.requests` must be specified
- `runtimeRef` must point to existing runtime with mpiImplementation set
- `namespace` must have resource quota sufficient for launcher + all workers

**State Transitions**:
```
Pending → Running → Succeeded
             ↓
            Failed
```
- **Pending**: Job created, waiting for gang scheduling
- **Running**: All pods scheduled and training in progress
- **Succeeded**: Training completed successfully
- **Failed**: Job failed (launcher failure, worker failure, gang scheduling timeout)

---

### 1.2 ClusterTrainingRuntime (MPI Runtime Template)

Template defining MPI framework configuration, pod templates, and scheduling policies. Cluster-scoped resource managed by platform administrators.

**CRD**: `kubeflow.org/v2alpha1/ClusterTrainingRuntime`

**Fields**:
- `metadata.name` (string, required): Unique identifier (e.g., "mpi-horovod-gpu")
- `spec.mlPolicy` (MLPolicySpec, required): ML-specific configuration
  - `numNodes` (int, optional): Default number of nodes (can be overridden in TrainJob)
  - `mpiImplementation` (string, required): "OpenMPI", "IntelMPI", or "MPICH"
  - `slotsPerWorker` (int, optional): MPI processes per worker (default: 1)
- `spec.podGroupPolicy` (PodGroupPolicySpec, required): Gang scheduling configuration
  - `coscheduling.scheduleTimeoutSeconds` (int): Timeout for scheduling all pods (default: 300)
  - `priorityClass` (string, optional): Pod priority class
- `spec.template` (RuntimeTemplate, required): Pod template definitions
  - `replicatedJobs[0]` (LauncherJob): Launcher pod template
    - `name`: "launcher"
    - `template.spec.containers[0]`: Launcher container with mpirun command
  - `replicatedJobs[1]` (WorkerJob): Worker pod template
    - `name`: "worker"
    - `replicas`: Default worker count (overridden by TrainJob.trainer.numNodes)
    - `template.spec.containers[0]`: Worker container with SSH server

**Relationships**:
- Referenced by multiple TrainJobs (one-to-many)
- Cluster-scoped (accessible from all namespaces)

**Validation Rules**:
- `mpiImplementation` must be one of: "OpenMPI", "IntelMPI", "MPICH"
- `scheduleTimeoutSeconds` must be >= 60 and <= 3600
- `replicatedJobs` must contain exactly 2 entries (launcher, worker)
- `slotsPerWorker` must be >= 1 and <= 16

**Example Values**:
- `name`: "mpi-horovod-gpu", "mpi-intel-cpu", "mpi-openmpi-distributed"
- `mpiImplementation`: "OpenMPI" (most common), "IntelMPI" (Intel hardware optimization)
- `scheduleTimeoutSeconds`: 300 (5 minutes default), 600 (10 minutes for large clusters)

---

### 1.3 TrainingRuntime (Namespace-Scoped MPI Runtime)

Similar to ClusterTrainingRuntime but namespace-scoped. Allows project-specific runtime customizations.

**CRD**: `kubeflow.org/v2alpha1/TrainingRuntime`

**Fields**: Same as ClusterTrainingRuntime

**Relationships**:
- Referenced by TrainJobs in same namespace only
- Namespace-scoped (project-specific)

**Use Cases**:
- Project-specific container images (e.g., private registry)
- Project-specific resource constraints
- Experimental runtime configurations before promoting to ClusterTrainingRuntime

---

### 1.4 Launcher Pod

Orchestration pod responsible for initializing MPI environment and executing mpirun command. Automatically created by Trainer V2 controller.

**Resource Type**: Kubernetes Pod (generated, not user-defined)

**Fields**:
- `metadata.name` (string): `{trainjob-name}-launcher-{random}`
- `metadata.namespace` (string): Same as TrainJob
- `metadata.labels`:
  - `training.kubeflow.org/job-name`: TrainJob name
  - `training.kubeflow.org/replica-type`: "launcher"
- `spec.containers[0]`:
  - `image`: From TrainJob.spec.trainer.image
  - `command`: `["mpirun", "--allow-run-as-root", "-np", "{totalProcesses}", ...]`
  - `resources`: Typically smaller than workers (CPU-bound coordination)
- `spec.volumes`:
  - SSH keys secret (generated by controller)
  - OpenMPI hostfile ConfigMap (generated by controller)

**Relationships**:
- Owned by one TrainJob
- Connects to N worker Pods via SSH (where N = TrainJob.spec.trainer.numNodes)
- Reads hostfile ConfigMap listing all worker Pod IPs
- Reads SSH private key Secret for authentication

**State Transitions**:
```
Pending → Running → Succeeded
             ↓
          Failed
```

**Lifecycle**:
1. **Pending**: Waiting for gang scheduling to place all pods
2. **Running**: SSH to workers, execute mpirun command
3. **Succeeded**: Training completed, mpirun exit code 0
4. **Failed**: SSH timeout, worker connection failure, training error

---

### 1.5 Worker Pod

Compute pods executing training workload under MPI coordination. Multiple workers form the distributed training cluster.

**Resource Type**: Kubernetes Pod (generated, not user-defined)

**Fields**:
- `metadata.name` (string): `{trainjob-name}-worker-{index}`
- `metadata.namespace` (string): Same as TrainJob
- `metadata.labels`:
  - `training.kubeflow.org/job-name`: TrainJob name
  - `training.kubeflow.org/replica-type`: "worker"
  - `training.kubeflow.org/replica-index`: Worker index (0, 1, 2, ...)
- `spec.containers[0]`:
  - `image`: From TrainJob.spec.trainer.image
  - `command`: SSH server (sshd) for launcher connections
  - `resources`: From TrainJob.spec.trainer.resourcesPerNode
- `spec.volumes`:
  - SSH keys secret (generated by controller)
  - Training data volumes (PVC, S3, NFS from TrainJob)
  - Model checkpoint volumes

**Relationships**:
- Owned by one TrainJob
- Receives connections from one launcher Pod
- Communicates with other worker Pods via MPI (e.g., AllReduce operations)
- Mounts shared storage (PVC, NFS, S3) for datasets and checkpoints

**State Transitions**:
```
Pending → Running → Succeeded
             ↓
          Failed
```

**Failure Modes**:
- **ImagePullBackOff**: Container image not found or authentication failed
- **OOMKilled**: Out of memory (insufficient memory allocation)
- **CrashLoopBackOff**: Training script error, CUDA error, SSH server failure
- **Evicted**: Node resource pressure, preemption by higher priority pods

---

### 1.6 PodGroup (Gang Scheduling Resource)

Logical grouping ensuring launcher and all worker pods schedule atomically. Created by Trainer V2 controller.

**CRD**: `scheduling.sigs.k8s.io/v1alpha1/PodGroup` (or vendor-specific: Volcano, Kueue)

**Fields**:
- `metadata.name` (string): `{trainjob-name}`
- `metadata.namespace` (string): Same as TrainJob
- `spec.minMember` (int): 1 (launcher) + numNodes (workers)
- `spec.scheduleTimeoutSeconds` (int): From ClusterTrainingRuntime.spec.podGroupPolicy
- `spec.priorityClassName` (string, optional): Priority class for scheduling

**Relationships**:
- References one TrainJob
- Controls scheduling of launcher Pod and all worker Pods

**Lifecycle**:
1. **Pending**: Waiting to reserve resources for all pods
2. **Scheduled**: All pods placed on nodes atomically
3. **Timeout**: Could not schedule all pods within scheduleTimeoutSeconds

**Timeout Behavior**:
- If resources unavailable within timeout → TrainJob fails with clear error message
- Error message should indicate: "Gang scheduling timeout: Requested {numNodes+1} pods, only {available} nodes available with sufficient resources"

---

### 1.7 JobSet (Kubernetes Batch Resource)

Manages groups of Kubernetes Jobs for launcher and workers. Created by Trainer V2 controller from TrainJob + RuntimeTemplate.

**CRD**: `jobset.x-k8s.io/v1alpha2/JobSet`

**Fields**:
- `metadata.name` (string): `{trainjob-name}`
- `metadata.namespace` (string): Same as TrainJob
- `spec.replicatedJobs`:
  - Launcher Job (1 replica)
  - Worker Job (numNodes replicas)
- `spec.successPolicy`: Conditions for job success
- `spec.failurePolicy`: Conditions for job failure and restart

**Relationships**:
- Owned by one TrainJob
- Owns launcher Job and worker Job
- Each Job owns respective Pods

**Purpose**: Provides native Kubernetes batch semantics (restart policies, success/failure conditions) without custom controller implementation.

---

### 1.8 SSH Secret (Ephemeral Authentication)

Kubernetes Secret containing SSH key pair for launcher-to-worker authentication. Generated by Trainer V2 controller per TrainJob.

**Resource Type**: Kubernetes Secret

**Fields**:
- `metadata.name` (string): `{trainjob-name}-ssh`
- `metadata.namespace` (string): Same as TrainJob
- `metadata.ownerReferences`: TrainJob (for garbage collection)
- `data.id_rsa` (base64): Private key (launcher)
- `data.id_rsa.pub` (base64): Public key (workers)
- `data.authorized_keys` (base64): Public key in authorized_keys format

**Relationships**:
- Owned by one TrainJob
- Mounted by launcher Pod (private key)
- Mounted by worker Pods (authorized_keys)

**Security Properties**:
- **Ephemeral**: Generated per job, deleted when job completes
- **Namespace-scoped**: Cannot be accessed from other namespaces
- **RBAC-protected**: Only pods in same namespace can mount
- **Non-persistent**: Not stored outside Kubernetes etcd
- **Algorithm**: RSA 4096-bit or Ed25519 (FIPS-compliant for federal customers)

---

### 1.9 Hostfile ConfigMap (MPI Worker Discovery)

Kubernetes ConfigMap containing OpenMPI hostfile listing all worker Pod IPs and slot counts. Generated by Trainer V2 controller.

**Resource Type**: Kubernetes ConfigMap

**Fields**:
- `metadata.name` (string): `{trainjob-name}-hostfile`
- `metadata.namespace` (string): Same as TrainJob
- `metadata.ownerReferences`: TrainJob (for garbage collection)
- `data.hostfile` (string): OpenMPI hostfile format

**Hostfile Format**:
```
worker-0.{trainjob-name}-worker.{namespace}.svc.cluster.local slots={slotsPerWorker}
worker-1.{trainjob-name}-worker.{namespace}.svc.cluster.local slots={slotsPerWorker}
worker-2.{trainjob-name}-worker.{namespace}.svc.cluster.local slots={slotsPerWorker}
...
```

**Relationships**:
- Owned by one TrainJob
- Mounted by launcher Pod at `/etc/mpi/hostfile`

**Purpose**: Tells mpirun command which worker nodes to connect to and how many MPI processes to spawn per worker.

---

### 1.10 Training Job Status (Observability)

State and health information for MPIJob. Part of TrainJob.status field.

**Fields**:
- `phase` (string, required): Current job phase
  - `Pending`: Job created, waiting for gang scheduling
  - `Running`: All pods scheduled, training in progress
  - `Succeeded`: Training completed successfully
  - `Failed`: Job failed (see failureReason)
- `conditions` ([]Condition, required): Detailed status conditions
  - `type`: "GangScheduling", "LauncherReady", "WorkersReady", "TrainingComplete"
  - `status`: "True", "False", "Unknown"
  - `reason`: Machine-readable reason code
  - `message`: Human-readable message
  - `lastTransitionTime`: Timestamp of last condition change
- `launcherPodStatus` (PodStatus, optional):
  - `phase`: "Pending", "Running", "Succeeded", "Failed"
  - `reason`: Failure reason if failed
  - `message`: Detailed error message
- `workerPodStatuses` ([]PodStatus, optional): Array of worker pod states
  - `index`: Worker index (0, 1, 2, ...)
  - `phase`: Pod phase
  - `reason`: Failure reason if failed
  - `message`: Detailed error message
- `startTime` (timestamp, optional): Job start time (when gang scheduling completed)
- `completionTime` (timestamp, optional): Job completion time
- `failureReason` (string, optional): High-level failure category
  - `GangSchedulingTimeout`: Could not schedule all pods within timeout
  - `LauncherFailed`: Launcher pod failed
  - `WorkerFailed`: One or more worker pods failed
  - `SSHConnectionFailed`: Launcher could not connect to workers
  - `MPIInitializationFailed`: MPI initialization error
  - `TrainingError`: User training script error
  - `ResourceQuotaExceeded`: Namespace resource quota insufficient

**Relationships**:
- Part of TrainJob resource
- Updated by Trainer V2 controller based on Pod events

**Update Frequency**: Real-time updates via Kubernetes watch API, reflected in Dashboard within 1 second.

---

## 2. Entity Relationships Diagram

```
┌──────────────────────────────────────────────────────────────┐
│ ClusterTrainingRuntime (Cluster-Scoped)                      │
│ - mpiImplementation: OpenMPI                                  │
│ - scheduleTimeoutSeconds: 300                                 │
│ - launcher pod template                                       │
│ - worker pod template                                         │
└────────────────┬─────────────────────────────────────────────┘
                 │
                 │ referenced by (1:N)
                 │
                 ▼
┌──────────────────────────────────────────────────────────────┐
│ TrainJob (Namespace-Scoped)                                  │
│ - runtimeRef: mpi-horovod-gpu                                 │
│ - numNodes: 4                                                 │
│ - image: horovod/horovod:latest                              │
│ - resourcesPerNode: 8 CPU, 32Gi memory, 2 GPU                │
└────────────────┬─────────────────────────────────────────────┘
                 │
                 │ owns (1:1)
                 │
                 ├─────────────────────────────────────────────┐
                 │                                              │
                 ▼                                              ▼
┌─────────────────────────────────┐    ┌─────────────────────────────────┐
│ PodGroup (Gang Scheduling)      │    │ JobSet (Batch Management)        │
│ - minMember: 5 (1+4)             │    │ - replicatedJobs: [launcher,    │
│ - scheduleTimeoutSeconds: 300    │    │   worker]                        │
└─────────────────────────────────┘    └─────────────┬───────────────────┘
                 │                                    │
                 │ controls scheduling                │ manages
                 │                                    │
                 ▼                                    ▼
┌──────────────────────────────────────────────────────────────┐
│ Launcher Pod                                                  │
│ - command: mpirun -np 8 python train.py                      │
│ - mounts: SSH private key, hostfile                          │
└────────────────┬─────────────────────────────────────────────┘
                 │
                 │ connects via SSH (1:N)
                 │
                 ▼
┌──────────────────────────────────────────────────────────────┐
│ Worker Pods (N=4)                                             │
│ - worker-0: index=0, slots=2, resources=8CPU/32Gi/2GPU       │
│ - worker-1: index=1, slots=2, resources=8CPU/32Gi/2GPU       │
│ - worker-2: index=2, slots=2, resources=8CPU/32Gi/2GPU       │
│ - worker-3: index=3, slots=2, resources=8CPU/32Gi/2GPU       │
│ - mounts: SSH authorized_keys, training data PVC             │
└──────────────────────────────────────────────────────────────┘
                 │
                 │ MPI AllReduce communication (N:N)
                 │
                 ▼
┌──────────────────────────────────────────────────────────────┐
│ Shared Storage (PVC/S3/NFS)                                   │
│ - Training datasets                                           │
│ - Model checkpoints                                           │
│ - Final model artifacts                                       │
└──────────────────────────────────────────────────────────────┘

Supporting Resources (Owned by TrainJob):

┌─────────────────────────────────┐    ┌─────────────────────────────────┐
│ SSH Secret                       │    │ Hostfile ConfigMap               │
│ - id_rsa (private key)           │    │ - worker-0 slots=2              │
│ - authorized_keys (public key)   │    │ - worker-1 slots=2              │
└─────────────────────────────────┘    │ - worker-2 slots=2              │
                                         │ - worker-3 slots=2              │
                                         └─────────────────────────────────┘
```

---

## 3. API Contracts

### 3.1 TrainJob Creation Contract

**Request**: Create TrainJob via kubectl apply or Python SDK

**Validation**:
1. RuntimeRef must point to existing TrainingRuntime or ClusterTrainingRuntime
2. RuntimeRef target must have `mpiImplementation` set (not null)
3. `numNodes` must be >= 1 and <= 64
4. `image` must be valid container image reference
5. `resourcesPerNode.requests` must be specified
6. Namespace must have sufficient resource quota for launcher + (numNodes × resourcesPerNode)

**Controller Actions**:
1. Generate SSH key pair → create Secret
2. Create hostfile ConfigMap with worker DNS names
3. Render runtime template with TrainJob parameters → create JobSet
4. Create PodGroup with minMember = 1 + numNodes
5. Update TrainJob status to Pending

**Expected Response**: TrainJob created with status.phase = Pending

---

### 3.2 Gang Scheduling Contract

**Trigger**: PodGroup created with minMember = N

**Scheduler Actions** (coscheduling plugin or Kueue):
1. Reserve resources for N pods atomically
2. If resources available: Mark PodGroup as Scheduled, allow pods to start
3. If resources unavailable: Wait up to scheduleTimeoutSeconds
4. If timeout exceeded: Fail PodGroup, controller updates TrainJob status to Failed

**Expected Outcomes**:
- **Success**: All pods transition to Running within scheduleTimeoutSeconds
- **Failure**: TrainJob.status.failureReason = "GangSchedulingTimeout", TrainJob.status.conditions includes diagnostic message showing available vs requested resources

---

### 3.3 MPI Initialization Contract

**Prerequisites**: All pods in Running phase

**Launcher Actions**:
1. Read hostfile from /etc/mpi/hostfile
2. Wait for worker SSH servers to be ready (poll SSH port 22)
3. Test SSH connections to all workers using private key
4. Execute mpirun command with hostfile

**Worker Actions**:
1. Start SSH server (sshd) on container startup
2. Mount authorized_keys from Secret
3. Accept SSH connections from launcher
4. Execute training script commands from launcher

**Expected Outcomes**:
- **Success**: MPI initialization completes within 60 seconds, training starts
- **Failure**: Launcher fails with SSH connection error, TrainJob.status.failureReason = "SSHConnectionFailed", diagnostics show which worker(s) unreachable

---

### 3.4 Training Completion Contract

**Training Success**:
- mpirun command exits with code 0
- Launcher pod status: Succeeded
- TrainJob status: Succeeded
- Model artifacts saved to configured storage

**Training Failure**:
- mpirun command exits with non-zero code OR
- Any worker pod crashes (OOM, CUDA error, script exception) OR
- Launcher pod crashes
- TrainJob status: Failed with appropriate failureReason

**Cleanup**:
- SSH Secret deleted automatically (ownerReference)
- Hostfile ConfigMap deleted automatically
- Launcher and worker Pods retained for debugging if cleanPodPolicy allows, otherwise deleted after TTL

---

### 3.5 Dashboard API Contract

**Endpoint**: `GET /api/v1/trainjobs`

**Response**:
```json
{
  "items": [
    {
      "name": "mnist-training",
      "namespace": "ml-project",
      "runtime": "mpi-horovod-gpu",
      "framework": "MPI",
      "status": "Running",
      "startTime": "2025-10-28T10:00:00Z",
      "duration": "5m30s",
      "numNodes": 4,
      "launcherStatus": "Running",
      "workerStatuses": [
        {"index": 0, "status": "Running"},
        {"index": 1, "status": "Running"},
        {"index": 2, "status": "Running"},
        {"index": 3, "status": "Running"}
      ]
    }
  ]
}
```

**Contract**: Response includes all TrainJobs in namespace with MPI runtime, with real-time status updates via WebSocket.

---

### 3.6 SDK API Contract

**Python SDK Method**: `create_mpi_job()`

**Signature**:
```python
def create_mpi_job(
    name: str,
    runtime: str,
    num_nodes: int,
    image: str,
    command: List[str],
    resources_per_node: dict,
    namespace: str = "default"
) -> TrainJob
```

**Contract**: Creates TrainJob resource, validates parameters client-side before submission, returns TrainJob object with initial status.

---

## 4. Key Validation Rules Summary

| Entity | Validation Rule | Error Message |
|--------|----------------|---------------|
| TrainJob | numNodes >= 1 and <= 64 | "numNodes must be between 1 and 64" |
| TrainJob | runtimeRef must exist | "Runtime {name} not found in namespace" |
| TrainJob | Runtime must have mpiImplementation | "Runtime {name} does not support MPI" |
| TrainJob | Namespace quota must be sufficient | "Insufficient namespace quota: requested {X}, available {Y}" |
| ClusterTrainingRuntime | scheduleTimeoutSeconds >= 60 | "scheduleTimeoutSeconds must be at least 60" |
| ClusterTrainingRuntime | replicatedJobs must have launcher and worker | "Runtime must define launcher and worker jobs" |
| PodGroup | minMember must equal launcher + workers | "PodGroup minMember mismatch: expected {N}, got {M}" |
| Worker Pod | Resource requests must be specified | "Worker pod missing resource requests" |

---

## 5. State Machine Summary

### TrainJob Phases

```
[Created]
   ↓
[Pending] ← Gang scheduling in progress
   ↓
   ├─→ [Failed] ← Gang scheduling timeout
   │
   ↓ (All pods scheduled)
[Running] ← Training in progress
   ↓
   ├─→ [Failed] ← Launcher failure, worker failure, training error
   │
   ↓ (Training complete)
[Succeeded]
```

### Key Transitions

- **Pending → Running**: Triggered by gang scheduling success (all pods placed)
- **Running → Succeeded**: Triggered by mpirun exit code 0
- **Running → Failed**: Triggered by pod failure, SSH connection failure, or training error
- **Pending → Failed**: Triggered by gang scheduling timeout

---

## Summary

This data model defines 10 core entities for MPIJob support:
1. **TrainJob**: User-facing training job specification
2. **ClusterTrainingRuntime**: Platform-managed MPI runtime template
3. **TrainingRuntime**: Namespace-scoped runtime variant
4. **Launcher Pod**: MPI orchestration pod
5. **Worker Pod**: Compute pods (N instances)
6. **PodGroup**: Gang scheduling coordinator
7. **JobSet**: Kubernetes batch management
8. **SSH Secret**: Ephemeral authentication
9. **Hostfile ConfigMap**: Worker discovery
10. **Training Job Status**: Observability data

All entities validated with clear error messages, state transitions defined, and API contracts specified for Dashboard, SDK, and CLI integration.

**Next Steps**: Generate API contracts in `/contracts/` directory and create quickstart.md with end-to-end examples.
