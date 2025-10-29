# Data Model: MPIJob Support for OpenShift AI

**Date**: 2025-10-29
**Status**: Complete
**Source**: Extracted from spec.md Key Entities section and ARCHITECTURE.md

## Overview

This document defines the data entities and their relationships for the MPIJob feature. The data model follows Kubernetes Custom Resource patterns and integrates with KubeFlow Trainer V2 API specifications.

---

## Core Entities

### 1. MPIJob (Custom Resource)

**Description**: Represents a distributed MPI training job with launcher and worker pods.

**API Version**: `kubeflow.org/v2beta1`

**Schema**:

```yaml
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: string                      # Job identifier (DNS-1123 subdomain)
  namespace: string                 # Kubernetes namespace for multi-tenancy
  labels:
    training.kubeflow.org/job-name: string
    training.kubeflow.org/job-type: "mpijob"
  annotations:
    sidecar.istio.io/inject: "false"  # CRITICAL: Prevent Istio injection
spec:
  mpiReplicaSpecs:
    Launcher:
      replicas: 1                   # Always 1 launcher
      template:
        spec:
          containers:
          - name: string
            image: string           # Container image with MPI runtime
            command: []string       # Entrypoint command
            args: []string          # Command arguments
            env: []EnvVar
            volumeMounts: []VolumeMount
            resources:
              requests:
                cpu: string
                memory: string
              limits:
                cpu: string
                memory: string
    Worker:
      replicas: int32               # Number of worker pods (2-100)
      template:
        spec:
          containers:
          - name: string
            image: string           # Same or different from launcher
            resources:
              requests:
                cpu: string
                memory: string
                nvidia.com/gpu: string
              limits:
                cpu: string
                memory: string
                nvidia.com/gpu: string
            volumeMounts: []VolumeMount
          nodeSelector: map[string]string
          tolerations: []Toleration
          affinity: Affinity
  cleanPodPolicy: string            # "Running" | "All" | "None"
  slotsPerWorker: int32             # MPI slots per worker (default: 1)
  runPolicy:
    activeDeadlineSeconds: int64    # Max job duration
    backoffLimit: int32             # Retry count on failure
    ttlSecondsAfterFinished: int32  # Cleanup delay
status:
  conditions:
  - type: string                    # "Pending" | "Running" | "Succeeded" | "Failed"
    status: string                  # "True" | "False" | "Unknown"
    reason: string                  # Reason for status
    message: string                 # Human-readable details
    lastTransitionTime: timestamp
  replicaStatuses:
    Launcher:
      active: int32
      succeeded: int32
      failed: int32
    Worker:
      active: int32
      succeeded: int32
      failed: int32
  startTime: timestamp
  completionTime: timestamp
```

**Attributes**:

| Field | Type | Required | Description | Validation |
|-------|------|----------|-------------|------------|
| `metadata.name` | string | Yes | Job identifier | DNS-1123 subdomain (max 63 chars) |
| `metadata.namespace` | string | Yes | Namespace for multi-tenancy | Must exist, user must have RBAC access |
| `spec.mpiReplicaSpecs.Worker.replicas` | int32 | Yes | Number of worker pods | 2-100 (enforced by admission webhook) |
| `spec.mpiReplicaSpecs.Worker.template.spec.containers[0].image` | string | Yes | Container image | Must be pullable, contain MPI runtime |
| `spec.mpiReplicaSpecs.Worker.template.spec.containers[0].resources.limits["nvidia.com/gpu"]` | string | No | GPUs per worker | 0-16 (configurable via LimitRange) |
| `spec.slotsPerWorker` | int32 | No | MPI slots per worker | Default: 1, Range: 1-8 |
| `spec.cleanPodPolicy` | string | No | Pod cleanup behavior | Default: "Running", Options: "All", "None" |
| `status.conditions[].type` | string | No | Current job phase | "Pending", "Running", "Succeeded", "Failed" |

**Relationships**:
- 1 MPIJob → 1 Launcher Pod
- 1 MPIJob → N Worker Pods (N = `spec.mpiReplicaSpecs.Worker.replicas`)
- 1 MPIJob → 1 ConfigMap (MPI hostfile)
- 1 MPIJob → 1 Secret (SSH keys)
- 1 MPIJob → 1 NetworkPolicy (optional, for network isolation)

**State Transitions**:

```
Created → Pending → Running → Succeeded
                           └→ Failed
```

- **Pending**: Waiting for resources (GPUs, CPU, memory) or image pull
- **Running**: All worker pods are Running, launcher has started `mpirun`
- **Succeeded**: Launcher pod exited with code 0
- **Failed**: Any pod failed (OOMKilled, error exit code, eviction)

**Validation Rules**:
1. `spec.mpiReplicaSpecs.Worker.replicas >= 2` (distributed training requires 2+ workers)
2. `spec.mpiReplicaSpecs.Launcher.replicas == 1` (exactly one launcher)
3. Container image must include:
   - MPI runtime (Open MPI 4.1+ or Intel MPI)
   - SSH server (for process spawning)
   - Training framework (PyTorch, TensorFlow, Horovod)
4. If `nvidia.com/gpu` requested, nodes must have GPU Operator installed
5. Total resources (workers × resources per worker) must not exceed namespace ResourceQuota

---

### 2. Worker Pod

**Description**: Individual compute unit within an MPIJob, executes training workload under MPI coordination.

**Managed By**: Training Operator creates StatefulSet → pods with predictable DNS names

**Schema** (Kubernetes Pod):

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: {mpijob-name}-worker-{index}  # e.g., my-job-worker-0
  namespace: string
  labels:
    training.kubeflow.org/job-name: string
    training.kubeflow.org/replica-type: "worker"
    training.kubeflow.org/replica-index: string
  ownerReferences:
  - apiVersion: kubeflow.org/v2beta1
    kind: MPIJob
    name: string
    uid: string
spec:
  containers:
  - name: worker
    image: string
    resources:
      limits:
        nvidia.com/gpu: string
        memory: string
        cpu: string
    volumeMounts:
    - name: ssh-auth
      mountPath: /root/.ssh
      readOnly: true
  hostname: {mpijob-name}-worker-{index}
  subdomain: {mpijob-name}-worker
  restartPolicy: Never
  volumes:
  - name: ssh-auth
    secret:
      secretName: {mpijob-name}-worker-ssh
      defaultMode: 0600
status:
  phase: string                     # "Pending" | "Running" | "Succeeded" | "Failed"
  podIP: string
  hostIP: string
  conditions:
  - type: "Ready"
    status: string
  containerStatuses:
  - name: "worker"
    state:
      running:
        startedAt: timestamp
      terminated:
        exitCode: int32
        reason: string
```

**Attributes**:

| Field | Type | Description |
|-------|------|-------------|
| `metadata.name` | string | Predictable name for DNS resolution: `{job-name}-worker-{index}` |
| `spec.containers[0].resources.limits["nvidia.com/gpu"]` | string | GPU allocation (0-16 per worker) |
| `status.phase` | string | Pod lifecycle state |
| `status.podIP` | string | Pod IP address (for MPI communication) |
| `status.containerStatuses[0].state` | object | Container execution state (running/terminated) |

**Relationships**:
- N Worker Pods → 1 MPIJob (parent)
- 1 Worker Pod → 1 Node (scheduled by Kubernetes)
- 1 Worker Pod → 1 GPU (or multiple, up to 16)
- Each Worker Pod → 1 Secret (SSH authorized_keys)

**Metrics** (collected via DCGM exporter):
- `DCGM_FI_DEV_GPU_UTIL`: GPU utilization percentage (target: >85%)
- `DCGM_FI_DEV_FB_USED`: GPU memory used (bytes)
- `container_cpu_usage_seconds_total`: CPU usage
- `container_memory_working_set_bytes`: Memory usage
- `container_network_transmit_bytes_total`: Network bandwidth

---

### 3. Job Configuration

**Description**: Declarative specification of MPIJob requirements, stored as Kubernetes Custom Resource.

**Format**: YAML/JSON compatible with KubeFlow Trainer V2 schema

**Example**:

```yaml
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: pytorch-distributed
  namespace: my-project
spec:
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      template:
        spec:
          containers:
          - name: launcher
            image: quay.io/my-org/pytorch-horovod:latest
            command: ["mpirun", "-np", "4", "python", "train.py"]
            env:
            - name: HOROVOD_TIMELINE
              value: "/data/timeline.json"
            volumeMounts:
            - name: training-data
              mountPath: /data
          volumes:
          - name: training-data
            persistentVolumeClaim:
              claimName: imagenet-pvc
    Worker:
      replicas: 4
      template:
        spec:
          containers:
          - name: worker
            image: quay.io/my-org/pytorch-horovod:latest
            resources:
              limits:
                nvidia.com/gpu: "1"
                memory: "64Gi"
                cpu: "16"
            volumeMounts:
            - name: training-data
              mountPath: /data
          volumes:
          - name: training-data
            persistentVolumeClaim:
              claimName: imagenet-pvc
          nodeSelector:
            nvidia.com/gpu.present: "true"
  cleanPodPolicy: Running
  slotsPerWorker: 1
```

**Validation Rules**:
1. Schema validation via Kubernetes API server (CRD OpenAPI schema)
2. Admission webhook validates:
   - Resource requests do not exceed namespace quota
   - Container image is valid format
   - GPU request is supported by cluster
3. User-defined validation via `--dry-run=server` before submission

**Relationships**:
- 1 Job Configuration → 1 MPIJob CR (1:1 mapping)
- Job Configuration editable via:
  - CLI: `oc odh mpijob create -f config.yaml`
  - SDK: `MPIJob._build_manifest()`
  - Dashboard: Job creation wizard → generates YAML → submits

---

### 4. Resource Quota

**Description**: Namespace-level limits on GPU/CPU/memory consumption, enforced before MPIJob worker pods are scheduled.

**Managed By**: Kubernetes ResourceQuota controller

**Schema**:

```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: team-cv-quota
  namespace: team-cv
spec:
  hard:
    requests.nvidia.com/gpu: "32"      # Max 32 GPUs across all pods
    requests.cpu: "256"                # Max 256 CPU cores
    requests.memory: "2048Gi"          # Max 2 TB memory
    pods: "100"                        # Max 100 pods
  scopeSelector:
    matchExpressions:
    - operator: In
      scopeName: PriorityClass
      values: ["high-priority"]
status:
  hard:
    requests.nvidia.com/gpu: "32"
  used:
    requests.nvidia.com/gpu: "24"      # Currently allocated GPUs
```

**Attributes**:

| Field | Type | Description |
|-------|------|-------------|
| `spec.hard["requests.nvidia.com/gpu"]` | string | Maximum GPUs for namespace |
| `status.used["requests.nvidia.com/gpu"]` | string | Currently allocated GPUs |

**Relationships**:
- 1 Namespace → 0-N ResourceQuotas
- ResourceQuota enforced before MPIJob worker pods created
- If `(status.used + mpijob.workers × mpijob.gpu_per_worker) > spec.hard`, job rejected

**Error Handling**:

When quota exceeded:

```bash
Error: Quota exceeded
Details: Requested 40 GPUs (20 workers × 2 GPUs), but namespace quota is 32 GPUs.
Current usage: 24 GPUs (3 running jobs)
Available: 8 GPUs

Suggestions:
1. Reduce worker count to 16 (32 GPUs total)
2. Wait for existing jobs to complete
3. Request quota increase from administrator
```

---

### 5. Job Metrics

**Description**: Time-series performance and resource data collected from worker pods, aggregated and displayed in Dashboard.

**Collection**: Prometheus scrapes metrics from:
- Training Operator (job-level metrics)
- kube-state-metrics (pod-level metrics)
- NVIDIA DCGM exporter (GPU metrics)

**Schema** (Prometheus time series):

```
# Job-level metrics (from Training Operator)
training_operator_mpijob_created_total{namespace="team-cv"} 45
training_operator_mpijob_running_total{namespace="team-cv"} 3
training_operator_mpijob_succeeded_total{namespace="team-cv"} 38
training_operator_mpijob_failed_total{namespace="team-cv"} 4
training_operator_mpijob_duration_seconds_bucket{namespace="team-cv",le="3600"} 35

# Pod-level metrics (from kube-state-metrics)
kube_pod_container_resource_requests{pod="my-job-worker-0",resource="nvidia.com/gpu"} 1
kube_pod_container_resource_requests{pod="my-job-worker-0",resource="memory"} 68719476736

# GPU metrics (from DCGM exporter)
DCGM_FI_DEV_GPU_UTIL{pod="my-job-worker-0",gpu="0",namespace="team-cv"} 92.5
DCGM_FI_DEV_FB_USED{pod="my-job-worker-0",gpu="0",namespace="team-cv"} 72000000000
DCGM_FI_PROF_PIPE_TENSOR_ACTIVE{pod="my-job-worker-0",gpu="0"} 98.2

# Network metrics (from node-exporter)
container_network_transmit_bytes_total{pod="my-job-worker-0"} 123456789012
```

**Attributes**:

| Metric | Type | Labels | Description |
|--------|------|--------|-------------|
| `training_operator_mpijob_created_total` | Counter | namespace | Total MPIJobs created |
| `training_operator_mpijob_duration_seconds` | Histogram | namespace | Job duration distribution |
| `DCGM_FI_DEV_GPU_UTIL` | Gauge | pod, gpu, namespace | GPU utilization % (0-100) |
| `DCGM_FI_DEV_FB_USED` | Gauge | pod, gpu, namespace | GPU memory used (bytes) |

**Relationships**:
- N Metrics → 1 MPIJob (many metrics per job)
- Metrics stored in Prometheus (15-90 day retention)
- Metrics visualized in Grafana dashboards
- Dashboard UI queries Prometheus API for real-time metrics

**Query Examples** (PromQL):

```promql
# Average GPU utilization across all workers
avg(DCGM_FI_DEV_GPU_UTIL{pod=~"my-job-worker-.*"}) by (pod)

# Scaling efficiency: actual throughput vs. expected
(rate(training_samples_processed_total{job="8-node"}[5m]) /
 rate(training_samples_processed_total{job="single-node"}[5m])) / 8

# Network bandwidth per worker
sum(rate(container_network_transmit_bytes_total{pod=~"my-job-worker-.*"}[1m])) by (pod)
```

---

## Supporting Entities

### 6. ConfigMap (MPI Hostfile)

**Description**: Generated by Training Operator, contains list of worker pod DNS names for MPI process spawning.

**Example**:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: my-job-hostfile
  namespace: team-cv
  ownerReferences:
  - apiVersion: kubeflow.org/v2beta1
    kind: MPIJob
    name: my-job
data:
  hostfile: |
    my-job-worker-0 slots=1
    my-job-worker-1 slots=1
    my-job-worker-2 slots=1
    my-job-worker-3 slots=1
```

**Lifecycle**: Created when MPIJob enters Running phase, deleted when MPIJob deleted (via ownerReferences)

---

### 7. Secret (SSH Keys)

**Description**: Ephemeral RSA-4096 SSH keypair for MPI inter-worker communication.

**Example**:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: my-job-ssh-key
  namespace: team-cv
  ownerReferences:
  - apiVersion: kubeflow.org/v2beta1
    kind: MPIJob
    name: my-job
type: Opaque
data:
  id_rsa: <base64-encoded-private-key>      # Mounted in launcher pod
  id_rsa.pub: <base64-encoded-public-key>   # Used to generate authorized_keys
  authorized_keys: <base64-encoded>         # Mounted in worker pods
```

**Security Properties**:
- Keys are RSA-4096 (FIPS 140-2 compliant)
- Keys scoped to single job (no cross-job access)
- Automatic cleanup on job deletion (ownerReferences)
- No key rotation needed (ephemeral lifecycle)

**Lifecycle**: Generated by Training Operator when MPIJob created, deleted when MPIJob deleted

---

### 8. NetworkPolicy (Optional)

**Description**: Kubernetes NetworkPolicy for namespace isolation and MPI communication rules.

**Example**:

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mpijob-my-job
  namespace: team-cv
spec:
  podSelector:
    matchLabels:
      training.kubeflow.org/job-name: my-job
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow launcher to SSH into workers
  - from:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-name: my-job
          training.kubeflow.org/replica-type: launcher
    ports:
    - protocol: TCP
      port: 2222
  # Allow worker-to-worker MPI communication
  - from:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-name: my-job
          training.kubeflow.org/replica-type: worker
    ports:
    - protocol: TCP
      port: 1024-65535
  egress:
  # Allow DNS resolution
  - to:
    - namespaceSelector:
        matchLabels:
          name: openshift-dns
    ports:
    - protocol: UDP
      port: 53
  # Allow worker-to-worker communication
  - to:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-name: my-job
```

**Purpose**: Enforce multi-tenancy isolation (pods in job A cannot connect to pods in job B)

---

## Entity Relationship Diagram

```
┌──────────────────────────────────────────────────────────────────┐
│                         MPIJob (CR)                               │
│  - metadata.name: string                                          │
│  - spec.mpiReplicaSpecs.Worker.replicas: int32                   │
│  - spec.cleanPodPolicy: string                                    │
│  - status.conditions[].type: string                               │
└────────┬─────────────────────────────────┬──────────────┬─────────┘
         │ ownerReferences                 │              │
         │                                 │              │
    ┌────▼─────────┐              ┌────────▼──────┐  ┌───▼──────────┐
    │ ConfigMap    │              │ Secret (SSH)  │  │ NetworkPolicy│
    │ (hostfile)   │              │ (keys)        │  │ (optional)   │
    └──────────────┘              └───────────────┘  └──────────────┘
         │
         │ consumed by
         │
    ┌────▼─────────────────────────────────────────────────────────┐
    │                    Launcher Pod (1)                           │
    │  - mounts: SSH private key                                    │
    │  - command: mpirun -np N --hostfile /etc/mpi/hostfile ...     │
    └──────┬────────────────────────────────────────────────────────┘
           │ SSH spawns processes on
           │
    ┌──────▼────────────────────────────────────────────────────────┐
    │                  Worker Pods (N)                               │
    │  - name: {job-name}-worker-{0..N-1}                           │
    │  - mounts: SSH authorized_keys                                │
    │  - resources.limits["nvidia.com/gpu"]: string                 │
    └──────┬────────────────────────────────────────────────────────┘
           │ scheduled on
           │
    ┌──────▼────────────────────────────────────────────────────────┐
    │                     GPU Nodes                                  │
    │  - nvidia.com/gpu: available GPUs                             │
    └───────────────────────────────────────────────────────────────┘
           │ metrics scraped by
           │
    ┌──────▼────────────────────────────────────────────────────────┐
    │                  Prometheus + DCGM Exporter                    │
    │  - DCGM_FI_DEV_GPU_UTIL                                       │
    │  - DCGM_FI_DEV_FB_USED                                        │
    └──────┬────────────────────────────────────────────────────────┘
           │ visualized in
           │
    ┌──────▼────────────────────────────────────────────────────────┐
    │              Grafana Dashboard / ODH Dashboard UI              │
    └───────────────────────────────────────────────────────────────┘

Constraints:
- ResourceQuota.spec.hard["nvidia.com/gpu"] >= Σ(Worker Pods × GPU per worker)
- Kueue PodGroup ensures Worker Pods admitted atomically
- All Worker Pods must reach Ready state before Launcher starts mpirun
```

---

## Data Flow

### Job Submission Flow

```
1. User submits MPIJob (CLI/SDK/Dashboard)
   ↓
2. Kubernetes API server validates schema (CRD)
   ↓
3. Admission webhook checks ResourceQuota
   ↓
4. MPIJob CR created in etcd
   ↓
5. Training Operator watches MPIJob CR
   ↓
6. Operator creates:
   - ConfigMap (hostfile)
   - Secret (SSH keys)
   - Job (launcher pod)
   - StatefulSet (worker pods)
   ↓
7. Kueue queues job if resources insufficient
   ↓
8. Kubernetes scheduler places pods on nodes
   ↓
9. All worker pods reach Ready state
   ↓
10. Launcher executes mpirun command
   ↓
11. MPI processes communicate via SSH + TCP/IP
   ↓
12. Training completes, launcher exits
   ↓
13. Training Operator updates MPIJob status
   ↓
14. User views results in Dashboard/CLI/SDK
```

### Metrics Collection Flow

```
1. Worker pods run training workload
   ↓
2. DCGM exporter scrapes GPU metrics (every 10s)
   ↓
3. Prometheus scrapes DCGM exporter
   ↓
4. Metrics stored in Prometheus TSDB
   ↓
5. Grafana queries Prometheus API
   ↓
6. Dashboard embeds Grafana panels
   ↓
7. User views real-time metrics
```

---

## Summary

This data model provides:

1. **MPIJob Custom Resource**: Core entity for job definition and status tracking
2. **Worker Pods**: Compute units with GPU allocation and MPI communication
3. **Job Configuration**: YAML specification validated by CRD schema
4. **Resource Quota**: Namespace-level enforcement of GPU limits
5. **Job Metrics**: Time-series data for observability
6. **Supporting Resources**: ConfigMap (hostfile), Secret (SSH keys), NetworkPolicy (isolation)

All entities follow Kubernetes API conventions and integrate with existing RHOAI components. The model supports:
- Multi-tenancy via namespaces and RBAC
- Gang scheduling via Kueue PodGroups
- Observability via Prometheus/Grafana
- Security via ephemeral SSH keys and NetworkPolicies
- Scalability from 2 to 100 workers per job
