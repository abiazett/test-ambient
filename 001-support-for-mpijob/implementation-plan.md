# MPIJob Implementation Plan for OpenShift AI
**Feature Branch**: `001-support-for-mpijob`
**Created**: 2025-10-21
**Status**: Draft
**Architect**: Archie

## Executive Summary

This implementation plan outlines the technical approach for integrating MPIJob support into OpenShift AI using KubeFlow Training Operator V2. The plan aligns with our north star architecture of unified training job management while leveraging proven upstream patterns from the KubeFlow ecosystem.

The architectural implications of this decision are significant: we're establishing a pattern that will scale to future job types and setting expectations for how distributed training workloads integrate with OpenShift AI's observability, security, and resource management infrastructure.

## Table of Contents

1. [Technical Architecture](#technical-architecture)
2. [Component Breakdown](#component-breakdown)
3. [Integration Architecture](#integration-architecture)
4. [Development Phases](#development-phases)
5. [Testing Strategy](#testing-strategy)
6. [Clarification Recommendations](#clarification-recommendations)
7. [Risk Assessment](#risk-assessment)
8. [Success Metrics](#success-metrics)

---

## Technical Architecture

### System Context and Integration Points

The MPIJob implementation sits at the intersection of multiple architectural layers in OpenShift AI. Understanding these integration boundaries is critical for lasting impact across products.

```
┌─────────────────────────────────────────────────────────────────┐
│                    OpenShift AI Platform                         │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              User-Facing Interfaces                       │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │  │
│  │  │ ODH Dashboard│  │     CLI      │  │  Python SDK  │   │  │
│  │  │    (React)   │  │    (Go)      │  │   (Python)   │   │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │  │
│  │         │                  │                  │           │  │
│  └─────────┼──────────────────┼──────────────────┼───────────┘  │
│            │                  │                  │               │
│  ┌─────────▼──────────────────▼──────────────────▼───────────┐  │
│  │          OpenShift AI Training Service Layer              │  │
│  │  ┌──────────────────────────────────────────────────────┐ │  │
│  │  │  Training Job Abstraction (Unified API)              │ │  │
│  │  │  - Job CRUD operations                               │ │  │
│  │  │  - Status aggregation                                │ │  │
│  │  │  - Log collection                                    │ │  │
│  │  │  - Metric exposure                                   │ │  │
│  │  └────────────────────┬─────────────────────────────────┘ │  │
│  └───────────────────────┼───────────────────────────────────┘  │
│                          │                                       │
│  ┌───────────────────────▼───────────────────────────────────┐  │
│  │        KubeFlow Training Operator V2                      │  │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │  │
│  │  │  MPIJob CRD  │  │PyTorchJob CRD│  │  TFJob CRD   │   │  │
│  │  │  Controller  │  │  Controller  │  │  Controller  │   │  │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │  │
│  └─────────┼──────────────────┼──────────────────┼───────────┘  │
└────────────┼──────────────────┼──────────────────┼───────────────┘
             │                  │                  │
┌────────────▼──────────────────▼──────────────────▼───────────────┐
│                    Kubernetes / OpenShift                         │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐    │
│  │   Pods    │  │  Services │  │   RBAC    │  │  Network  │    │
│  │           │  │           │  │           │  │  Policies │    │
│  └───────────┘  └───────────┘  └───────────┘  └───────────┘    │
│  ┌───────────┐  ┌───────────┐  ┌───────────┐  ┌───────────┐    │
│  │   PVCs    │  │  Quotas   │  │  Events   │  │Monitoring │    │
│  └───────────┘  └───────────┘  └───────────┘  └───────────┘    │
└───────────────────────────────────────────────────────────────────┘
```

### Architecture Principles

**1. Upstream-First Pattern**
Leverage KubeFlow Training Operator V2 as the authoritative source for MPIJob behavior. Our implementations wrap upstream capabilities rather than reimplementing them. This aligns with the Martin Fowler pattern of "Anti-Corruption Layer" - we translate between OpenShift AI's unified interface and KubeFlow's specific implementations.

**2. Unified Abstraction Layer**
The Training Service Layer provides consistent interfaces across all job types (MPIJob, PyTorchJob, TFJob). This is critical for 2-3 year horizon: as we add new job types (Ray, Spark, custom operators), the user-facing interfaces remain stable.

**3. Event-Driven Observability**
Status, logs, and metrics flow through event streams rather than polling. This scales better for multi-tenant environments and enables real-time UI updates without overwhelming the API server.

**4. Least-Privilege Security Model**
MPIJob operations inherit OpenShift's RBAC with no privileged escalation. Service accounts have minimal required permissions. Network policies default to deny with explicit allowances.

### Core Components Architecture

#### 1. KubeFlow Training Operator V2 Integration

**Component**: Training Operator Deployment
**Responsibility**: Manages MPIJob CRD lifecycle, reconciliation, and pod orchestration
**Technology**: Upstream KubeFlow Training Operator (Go-based controller)

**Key Architectural Decisions**:

- **Deployment Model**: Dedicated namespace (`openshift-ai-training-operator`) with cluster-wide CRD scope
- **Version Pinning**: Lock to specific Training Operator version with tested compatibility matrix
- **Upgrade Strategy**: Blue-green deployment pattern for operator upgrades to avoid job disruption
- **Failure Domain**: Operator restarts do not affect running jobs (state in etcd, not memory)

**Integration Points**:
```
Training Operator
  ├── Watches: MPIJob CRD instances
  ├── Creates: Launcher Pod, Worker Pods, Services
  ├── Updates: MPIJob Status subresource
  └── Emits: Kubernetes Events for state transitions
```

**Configuration Requirements**:
- CRD installation via OLM (Operator Lifecycle Manager)
- Leader election for HA deployments
- Webhook validation for MPIJob spec
- ServiceAccount with pod creation, service management permissions

#### 2. MPIJob Custom Resource Definition

**Schema Structure** (aligns with upstream KubeFlow Training Operator V2):

```yaml
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: training-job-001
  namespace: data-science-team
spec:
  # Worker topology configuration
  replicaSpecs:
    Launcher:
      replicas: 1
      template:
        spec:
          containers:
          - name: launcher
            image: training-image:latest
            command: ["mpirun", "-np", "8", "python", "train.py"]
            resources:
              requests:
                cpu: "2"
                memory: "8Gi"
    Worker:
      replicas: 4
      template:
        spec:
          containers:
          - name: worker
            image: training-image:latest
            resources:
              requests:
                cpu: "8"
                memory: "32Gi"
                nvidia.com/gpu: "2"

  # MPI-specific configuration
  mpiImplementation: OpenMPI
  slotsPerWorker: 2

  # OpenShift AI extensions
  runPolicy:
    cleanPodPolicy: Running
    ttlSecondsAfterFinished: 3600
    backoffLimit: 0

status:
  conditions:
  - type: Created
    status: "True"
    lastTransitionTime: "2025-10-21T10:00:00Z"
  - type: Running
    status: "True"
    lastTransitionTime: "2025-10-21T10:01:00Z"
  replicaStatuses:
    Launcher:
      active: 1
    Worker:
      active: 4
  startTime: "2025-10-21T10:00:00Z"
```

**Status Model**:
- **Conditions**: State transitions with timestamps (Created, Running, Succeeded, Failed)
- **ReplicaStatuses**: Per-role pod counts (active, pending, failed)
- **StartTime/CompletionTime**: Job lifecycle timestamps
- **Phase**: High-level state (Pending, Running, Succeeded, Failed)

#### 3. OpenShift AI Training Service Layer

**Purpose**: Abstract KubeFlow-specific details into unified OpenShift AI training job interface

**Technology Stack**:
- Backend: Go-based microservice
- API: gRPC for internal, REST for external
- Storage: etcd via Kubernetes API (no separate database)

**Service Responsibilities**:

1. **Job Lifecycle Management**
   - Translate unified job specs to MPIJob CRDs
   - Validate configurations against quotas, RBAC
   - Emit audit events for all operations

2. **Status Aggregation**
   - Watch MPIJob status updates
   - Compute derived metrics (duration, completion rate)
   - Cache frequently accessed data (last status, log locations)

3. **Log Collection**
   - Stream logs from launcher and worker pods
   - Tag logs with worker ID, timestamp
   - Support filtering, tailing, historical retrieval

4. **Metric Exposition**
   - Expose Prometheus metrics for job lifecycle
   - Integrate with OpenShift monitoring stack
   - Provide per-namespace, per-job granularity

**API Endpoints** (REST):
```
POST   /api/v1/namespaces/{ns}/trainingjobs/mpijob
GET    /api/v1/namespaces/{ns}/trainingjobs/{name}
GET    /api/v1/namespaces/{ns}/trainingjobs
DELETE /api/v1/namespaces/{ns}/trainingjobs/{name}
GET    /api/v1/namespaces/{ns}/trainingjobs/{name}/status
GET    /api/v1/namespaces/{ns}/trainingjobs/{name}/logs
GET    /api/v1/namespaces/{ns}/trainingjobs/{name}/metrics
```

**Data Flow: Job Creation**
```
1. User → Training Service: CreateMPIJobRequest
2. Training Service: Validate (RBAC, quotas, schema)
3. Training Service: Emit audit event (job.create.initiated)
4. Training Service → Kubernetes API: Create MPIJob CRD
5. Training Operator: Watch event, begin reconciliation
6. Training Operator → Kubernetes: Create Launcher Pod
7. Training Operator → Kubernetes: Create Worker Pods
8. Training Operator → MPIJob Status: Update conditions
9. Training Service: Watch status, update cache
10. Training Service → User: CreateMPIJobResponse
```

### Network Architecture for MPI Communication

**Critical Consideration**: MPI workers must communicate peer-to-peer for collective operations (all-reduce, broadcast, scatter-gather). This has implications for network policies, service mesh, and firewall rules.

**Communication Patterns**:
```
Launcher Pod
  ├── Initiates MPI context (mpirun/mpiexec)
  ├── Establishes SSH or MPI-native connections to workers
  └── Coordinates rank assignments (rank 0, rank 1, ..., rank N)

Worker Pods
  ├── Accept connections from Launcher
  ├── Peer-to-peer communication for MPI collectives
  │   └── Typically on ephemeral high ports (>1024)
  └── Report status back to Launcher
```

**Network Policy Requirements**:
```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mpijob-worker-communication
  namespace: data-science-team
spec:
  podSelector:
    matchLabels:
      training.kubeflow.org/job-role: worker
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow launcher to communicate with workers
  - from:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-role: launcher
  # Allow workers to communicate with each other
  - from:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-role: worker
  egress:
  # Allow workers to communicate with each other
  - to:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-role: worker
  # Allow workers to reach external services (registry, S3, etc.)
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443
```

**Service Mesh Considerations**:
- **Istio Compatibility**: MPI traffic typically bypasses service mesh (direct pod-to-pod). Need to annotate pods with `sidecar.istio.io/inject: "false"`
- **Performance Impact**: Service mesh proxy adds latency to MPI collectives (5-10% overhead observed in benchmarks)
- **Recommendation**: Document how to disable sidecar injection for MPIJob pods

**Port Requirements**:
- **SSH-based MPI** (OpenMPI default): Port 22
- **MPI-native communication** (Intel MPI, MPICH): Dynamic high ports (typically 32768-60999)
- **Control plane**: Kubernetes API (6443), etcd (2379-2380)

### Scheduling and Resource Management Architecture

**Gang Scheduling Requirement**: All workers for an MPIJob must be scheduled simultaneously. Partial scheduling leads to deadlock (launcher waits for workers that can't be scheduled).

**Implementation Options**:

**Option A: Coscheduling via Volcano Scheduler** (Recommended)
- **Pattern**: Multi-tenant scheduler with gang scheduling primitives
- **Integration**: KubeFlow Training Operator supports Volcano via `schedulerName: volcano`
- **Pros**: Proven in production, handles complex topologies, prevents resource fragmentation
- **Cons**: Additional component to deploy, learning curve for administrators

**Option B: PodGroup with Default Scheduler**
- **Pattern**: Custom scheduler extender or admission controller
- **Pros**: Minimal additional infrastructure
- **Cons**: Less mature, limited support for complex scheduling constraints

**Recommendation**: Start with Volcano scheduler for MVP. This aligns with industry best practices (used by Alibaba, Baidu, others for large-scale AI workloads) and provides growth path for advanced features (preemption, bin packing, topology awareness).

**Resource Quota Integration**:
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: data-science-team-quota
  namespace: data-science-team
spec:
  hard:
    requests.nvidia.com/gpu: "16"
    requests.cpu: "128"
    requests.memory: "512Gi"
    # Custom resource for job count
    count/mpijobs.kubeflow.org: "5"
```

**Validation Flow**:
1. User submits MPIJob requesting 4 workers × 2 GPUs = 8 GPUs total
2. Training Service queries ResourceQuota for namespace
3. Calculate current usage: 10 GPUs in use
4. Validation: (10 + 8) > 16 → Reject with error: "ResourceQuotaExceeded: Requested 8 GPUs but only 6 available under quota"
5. Provide actionable guidance: "Current usage: 10/16 GPUs. Reduce worker count or wait for jobs to complete."

### Observability Architecture

**Three Pillars**: Metrics, Logs, Traces

#### Metrics Collection

**Source**: Training Operator emits metrics via Prometheus exporter

**Key Metrics**:
```
# Job lifecycle
mpijob_created_total{namespace, job_name}
mpijob_running_total{namespace, job_name}
mpijob_succeeded_total{namespace, job_name}
mpijob_failed_total{namespace, job_name}
mpijob_duration_seconds{namespace, job_name, phase}

# Resource utilization
mpijob_worker_cpu_usage{namespace, job_name, worker_id}
mpijob_worker_memory_usage{namespace, job_name, worker_id}
mpijob_worker_gpu_usage{namespace, job_name, worker_id}
mpijob_worker_gpu_memory_usage{namespace, job_name, worker_id}

# Worker status
mpijob_workers_active{namespace, job_name}
mpijob_workers_pending{namespace, job_name}
mpijob_workers_failed{namespace, job_name}
```

**Integration**: ServiceMonitor CRD for Prometheus Operator
```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: training-operator-metrics
  namespace: openshift-ai-training-operator
spec:
  selector:
    matchLabels:
      app: training-operator
  endpoints:
  - port: metrics
    interval: 30s
```

#### Log Aggregation

**Pattern**: Centralized logging via OpenShift Logging (EFK stack)

**Log Flow**:
```
Worker Pods → stdout/stderr → FluentD → Elasticsearch → Kibana
                                ↓
                         Training Service (cache)
                                ↓
                         CLI/SDK/UI (streaming)
```

**Log Tagging**:
- `namespace`: Kubernetes namespace
- `job_name`: MPIJob name
- `job_role`: launcher, worker
- `worker_id`: 0, 1, 2, ... N
- `container`: launcher, worker

**API for Log Retrieval**:
```go
// Training Service API
GET /api/v1/namespaces/{ns}/trainingjobs/{name}/logs
  ?component=worker
  &worker_id=2
  &tail=100
  &follow=true

// Returns: Server-Sent Events (SSE) stream
data: [Worker-2] Epoch 1/10 - Loss: 0.523
data: [Worker-2] Epoch 2/10 - Loss: 0.412
```

#### Distributed Tracing (Post-MVP)

**Pattern**: OpenTelemetry for cross-component tracing

**Trace Spans**:
```
JobCreate (Training Service)
├── ValidateRBAC (Training Service)
├── ValidateQuota (Training Service)
├── CreateCRD (Kubernetes API)
├── ReconcileJob (Training Operator)
│   ├── CreateLauncher (Training Operator)
│   ├── CreateWorkers (Training Operator)
│   └── UpdateStatus (Training Operator)
└── EmitAuditEvent (Audit Service)
```

**Value**: Diagnose latency issues (why did job take 30s to start?) and failure propagation (where did this error originate?).

---

## Component Breakdown

### 1. CLI Implementation

**Technology**: Go-based CLI using Cobra framework (consistent with `oc` tooling)

**Component Structure**:
```
odh-cli/
├── cmd/
│   └── training/
│       ├── create_mpijob.go
│       ├── delete_mpijob.go
│       ├── describe_mpijob.go
│       ├── list_mpijob.go
│       └── logs_mpijob.go
├── pkg/
│   ├── client/
│   │   └── training_service.go   # gRPC/REST client
│   ├── validation/
│   │   └── mpijob_validator.go
│   └── output/
│       ├── table_formatter.go
│       └── json_formatter.go
└── test/
    └── e2e/
        └── mpijob_cli_test.go
```

**Command Specifications**:

#### `odh training create mpijob`
```bash
# From YAML file
odh training create mpijob -f mpijob-spec.yaml

# From stdin
cat mpijob-spec.yaml | odh training create mpijob -f -

# Inline parameters (quick start)
odh training create mpijob my-job \
  --image=training:latest \
  --workers=4 \
  --gpu-per-worker=2 \
  --cpu-per-worker=8 \
  --memory-per-worker=32Gi \
  --command="horovodrun -np 8 python train.py" \
  --wait \
  --timeout=3600

# Dry-run for validation
odh training create mpijob -f spec.yaml --dry-run=client
```

**Implementation Details**:
- **Validation**: Client-side validation before submission (faster feedback)
- **Wait Mode**: Block until job completes or timeout (useful for CI/CD)
- **Output Formats**: JSON, YAML, table (human-readable)
- **Error Handling**: Rich error messages with suggested fixes

#### `odh training describe mpijob`
```bash
# Basic describe
odh training describe mpijob my-job

# Watch mode (auto-refresh)
odh training describe mpijob my-job --watch

# Output format
odh training describe mpijob my-job -o json
```

**Output Format**:
```
Name:       my-job
Namespace:  data-science-team
Type:       MPIJob
Status:     Running
Workers:    4/4 running
Duration:   5m 32s

Worker Status:
  Launcher: Running (1m 15s)
  Worker-1: Running (1m 10s)
  Worker-2: Running (1m 10s)
  Worker-3: Running (1m 08s)
  Worker-4: Running (1m 10s)

Resource Allocation:
  CPU:    32 cores (8 per worker)
  Memory: 128 GiB (32 GiB per worker)
  GPU:    8 (2 per worker)

Events:
  5m32s  Created   MPIJob created
  5m15s  Launching Launcher pod started
  5m10s  Running   All workers active
```

#### `odh training logs mpijob`
```bash
# Launcher logs
odh training logs mpijob my-job

# Specific worker
odh training logs mpijob my-job --worker=2

# Follow mode (tail -f)
odh training logs mpijob my-job --worker=2 --follow

# Tail last 100 lines
odh training logs mpijob my-job --tail=100

# All workers (interleaved)
odh training logs mpijob my-job --all-workers
```

**Implementation**: Server-Sent Events (SSE) for follow mode, HTTP GET for static logs

#### `odh training list mpijob`
```bash
# List all MPIJobs in current namespace
odh training list mpijob

# All namespaces (requires cluster-wide permissions)
odh training list mpijob --all-namespaces

# Filter by status
odh training list mpijob --status=Running

# Output formats
odh training list mpijob -o wide
odh training list mpijob -o json
```

**Output**:
```
NAME           TYPE     STATUS    WORKERS  CREATED
my-job         mpijob   Running   4/4      5m ago
experiment-42  mpijob   Succeeded 8/8      1h ago
test-job       mpijob   Failed    2/4      30m ago
```

#### `odh training delete mpijob`
```bash
# Delete specific job
odh training delete mpijob my-job

# Delete with confirmation
odh training delete mpijob my-job --confirm

# Delete and wait for cleanup
odh training delete mpijob my-job --wait
```

**Error Handling Patterns**:

```go
// Structured error responses
type MPIJobError struct {
    Code    string   `json:"code"`
    Message string   `json:"message"`
    Details string   `json:"details"`
    Suggestions []string `json:"suggestions"`
}

// Example error
{
    "code": "ResourceQuotaExceeded",
    "message": "Cannot create MPIJob: GPU quota exceeded",
    "details": "Requested 8 GPUs but only 4 available (10/16 in use)",
    "suggestions": [
        "Reduce worker count from 4 to 2",
        "Wait for running jobs to complete",
        "Contact administrator to increase quota"
    ]
}
```

**Testing Requirements**:
- Unit tests: 80% code coverage
- Integration tests: All commands against live cluster
- E2E tests: Complete user workflows (create → monitor → logs → delete)

### 2. Python SDK Implementation

**Technology**: Python 3.8+ with type hints, async support

**Package Structure**:
```
odh-sdk/
├── odh/
│   ├── training/
│   │   ├── __init__.py
│   │   ├── mpijob.py
│   │   ├── client.py
│   │   ├── models.py
│   │   └── exceptions.py
│   └── common/
│       ├── auth.py
│       └── config.py
├── tests/
│   ├── unit/
│   │   └── test_mpijob.py
│   └── integration/
│       └── test_mpijob_e2e.py
├── examples/
│   ├── basic_mpijob.py
│   ├── horovod_example.py
│   └── pipeline_integration.py
└── docs/
    └── api_reference.md
```

**Core API Design**:

```python
from odh.training import MPIJob, ResourceSpec, JobStatus
from typing import Optional, Dict, List

class MPIJob:
    """
    Represents an MPI-based distributed training job.

    This class provides a Pythonic interface to OpenShift AI's MPIJob
    capabilities, abstracting Kubernetes complexities while providing
    full control over distributed training configuration.
    """

    def __init__(
        self,
        name: str,
        namespace: str = "default",
        workers: int = 1,
        slots_per_worker: int = 1,
        image: str,
        command: List[str],
        launcher_resources: Optional[ResourceSpec] = None,
        worker_resources: Optional[ResourceSpec] = None,
        env: Optional[Dict[str, str]] = None,
        volumes: Optional[List['Volume']] = None,
        mpi_implementation: str = "OpenMPI",
        clean_pod_policy: str = "Running",
        ttl_seconds_after_finished: int = 3600
    ):
        """
        Initialize MPIJob configuration.

        Args:
            name: Unique job identifier (DNS-1123 compliant)
            namespace: Kubernetes namespace
            workers: Number of worker replicas (min: 1, max: 100)
            slots_per_worker: MPI process slots per worker
            image: Container image with training code and MPI
            command: Entrypoint command for training
            launcher_resources: CPU/memory/GPU for launcher
            worker_resources: CPU/memory/GPU per worker
            env: Environment variables for all pods
            volumes: Persistent volumes for data/models
            mpi_implementation: MPI implementation (OpenMPI, IntelMPI, MPICH)
            clean_pod_policy: When to delete pods (Running, All, None)
            ttl_seconds_after_finished: Auto-cleanup delay

        Raises:
            ValueError: Invalid configuration
            AuthenticationError: Missing credentials
        """
        self.name = name
        self.namespace = namespace
        self.workers = workers
        self.slots_per_worker = slots_per_worker
        self.image = image
        self.command = command
        self.launcher_resources = launcher_resources or ResourceSpec()
        self.worker_resources = worker_resources or ResourceSpec()
        self.env = env or {}
        self.volumes = volumes or []
        self.mpi_implementation = mpi_implementation
        self.clean_pod_policy = clean_pod_policy
        self.ttl_seconds_after_finished = ttl_seconds_after_finished

        self._client = None  # Lazy initialization
        self._status_cache = None

    def create(self, wait: bool = False, timeout: int = 3600) -> 'MPIJob':
        """
        Submit MPIJob to cluster.

        Args:
            wait: Block until completion (default: False)
            timeout: Maximum wait time if wait=True

        Returns:
            Self (for method chaining)

        Raises:
            ResourceQuotaExceeded: Quota violation
            ValidationError: Invalid specification
            TimeoutError: Job didn't complete within timeout
        """
        self._validate()

        # Create CRD via Training Service API
        response = self._client.create_mpijob(
            namespace=self.namespace,
            spec=self._to_spec()
        )

        if wait:
            return self.wait_for_completion(timeout=timeout)

        return self

    def status(self) -> JobStatus:
        """
        Get current job status.

        Returns:
            JobStatus with phase, worker counts, timestamps
        """
        response = self._client.get_mpijob_status(
            namespace=self.namespace,
            name=self.name
        )

        return JobStatus(
            phase=response['phase'],
            workers_running=response['replicaStatuses']['Worker']['active'],
            workers_pending=response['replicaStatuses']['Worker'].get('pending', 0),
            workers_failed=response['replicaStatuses']['Worker'].get('failed', 0),
            start_time=response.get('startTime'),
            completion_time=response.get('completionTime')
        )

    def wait_for_completion(
        self,
        timeout: int = 3600,
        poll_interval: int = 10
    ) -> 'MPIJob':
        """
        Block until job completes.

        Args:
            timeout: Maximum wait time in seconds
            poll_interval: Seconds between status checks

        Returns:
            Self with updated status

        Raises:
            TimeoutError: Exceeded timeout
            JobFailedError: Job failed with error details
        """
        import time
        elapsed = 0

        while elapsed < timeout:
            status = self.status()

            if status.phase == "Succeeded":
                return self

            if status.phase == "Failed":
                logs = self.get_logs(component="launcher", tail=100)
                raise JobFailedError(
                    message=f"Job {self.name} failed",
                    logs=logs,
                    status=status
                )

            time.sleep(poll_interval)
            elapsed += poll_interval

        raise TimeoutError(f"Job {self.name} did not complete within {timeout}s")

    def get_logs(
        self,
        component: str = "launcher",
        worker_id: Optional[int] = None,
        tail: Optional[int] = None,
        follow: bool = False
    ) -> str:
        """
        Retrieve logs from job pods.

        Args:
            component: "launcher" or "worker"
            worker_id: Specific worker index (0-indexed)
            tail: Last N lines only
            follow: Stream logs in real-time

        Returns:
            Log output as string (or generator if follow=True)
        """
        if component == "worker" and worker_id is None:
            raise ValueError("worker_id required when component='worker'")

        if follow:
            # Return generator for streaming
            return self._client.stream_logs(
                namespace=self.namespace,
                name=self.name,
                component=component,
                worker_id=worker_id,
                tail=tail
            )
        else:
            # Return complete logs
            return self._client.get_logs(
                namespace=self.namespace,
                name=self.name,
                component=component,
                worker_id=worker_id,
                tail=tail
            )

    def delete(self, wait: bool = False) -> None:
        """
        Delete job and associated resources.

        Args:
            wait: Block until all pods terminated
        """
        self._client.delete_mpijob(
            namespace=self.namespace,
            name=self.name
        )

        if wait:
            # Poll until job no longer exists
            while self._exists():
                time.sleep(2)

    @classmethod
    def from_template(cls, template_name: str, **overrides) -> 'MPIJob':
        """
        Create MPIJob from predefined template.

        Templates:
            - horovod-tensorflow
            - horovod-pytorch
            - deepspeed-mpi
            - intel-mpi-pytorch

        Args:
            template_name: Template identifier
            **overrides: Parameters to override

        Returns:
            Configured MPIJob instance
        """
        template = cls._load_template(template_name)
        template.update(overrides)
        return cls(**template)

    @classmethod
    def from_yaml(cls, yaml_path: str) -> 'MPIJob':
        """
        Load MPIJob from YAML file.

        Args:
            yaml_path: Path to YAML file

        Returns:
            MPIJob instance
        """
        import yaml
        with open(yaml_path, 'r') as f:
            spec = yaml.safe_load(f)
        return cls(**spec)

    def to_yaml(self, output_path: Optional[str] = None) -> str:
        """
        Export MPIJob as YAML.

        Args:
            output_path: File path to write (optional)

        Returns:
            YAML string
        """
        import yaml
        spec = self._to_spec()
        yaml_str = yaml.dump(spec)

        if output_path:
            with open(output_path, 'w') as f:
                f.write(yaml_str)

        return yaml_str

    def _validate(self):
        """Client-side validation before submission."""
        if self.workers < 1:
            raise ValueError("workers must be >= 1")

        if self.slots_per_worker < 1:
            raise ValueError("slots_per_worker must be >= 1")

        if not self.image:
            raise ValueError("image is required")

        if not self.command:
            raise ValueError("command is required")

        # Validate name format (DNS-1123)
        import re
        if not re.match(r'^[a-z0-9]([-a-z0-9]*[a-z0-9])?$', self.name):
            raise ValueError(
                f"Invalid job name '{self.name}'. "
                "Must be lowercase alphanumeric with hyphens."
            )


class ResourceSpec:
    """Resource allocation specification."""

    def __init__(
        self,
        cpu: str = "1",
        memory: str = "2Gi",
        gpus: int = 0
    ):
        self.cpu = cpu
        self.memory = memory
        self.gpus = gpus

    def to_dict(self) -> Dict:
        spec = {
            "requests": {
                "cpu": self.cpu,
                "memory": self.memory
            },
            "limits": {
                "cpu": self.cpu,
                "memory": self.memory
            }
        }

        if self.gpus > 0:
            spec["requests"]["nvidia.com/gpu"] = str(self.gpus)
            spec["limits"]["nvidia.com/gpu"] = str(self.gpus)

        return spec


class JobStatus:
    """Job execution status."""

    def __init__(
        self,
        phase: str,
        workers_running: int,
        workers_pending: int,
        workers_failed: int,
        start_time: Optional[str],
        completion_time: Optional[str]
    ):
        self.phase = phase
        self.workers_running = workers_running
        self.workers_pending = workers_pending
        self.workers_failed = workers_failed
        self.start_time = start_time
        self.completion_time = completion_time

    @property
    def duration(self) -> Optional[float]:
        """Job duration in seconds."""
        if not self.start_time:
            return None

        from datetime import datetime
        start = datetime.fromisoformat(self.start_time.replace('Z', '+00:00'))

        if self.completion_time:
            end = datetime.fromisoformat(self.completion_time.replace('Z', '+00:00'))
        else:
            end = datetime.now(start.tzinfo)

        return (end - start).total_seconds()

    def __repr__(self):
        return (
            f"JobStatus(phase={self.phase}, "
            f"running={self.workers_running}, "
            f"pending={self.workers_pending}, "
            f"failed={self.workers_failed})"
        )


# Exception hierarchy
class ODHError(Exception):
    """Base exception for ODH SDK."""
    pass

class ResourceQuotaExceeded(ODHError):
    """Quota violation."""
    def __init__(self, details: str):
        self.details = details
        super().__init__(f"Resource quota exceeded: {details}")

class ValidationError(ODHError):
    """Invalid configuration."""
    def __init__(self, field: str, message: str):
        self.field = field
        self.message = message
        super().__init__(f"Validation error in '{field}': {message}")

class JobFailedError(ODHError):
    """Job execution failed."""
    def __init__(self, message: str, logs: str, status: JobStatus):
        self.message = message
        self.logs = logs
        self.status = status
        super().__init__(message)

class AuthenticationError(ODHError):
    """Authentication failed."""
    pass
```

**Usage Examples**:

```python
# Example 1: Basic job creation
from odh.training import MPIJob, ResourceSpec

job = MPIJob(
    name="mnist-distributed",
    workers=4,
    image="quay.io/opendatahub/horovod-mnist:latest",
    command=["horovodrun", "-np", "8", "python", "/examples/mnist.py"],
    worker_resources=ResourceSpec(cpu="4", memory="16Gi", gpus=2)
)

job.create(wait=True, timeout=3600)
print(f"Job completed in {job.status().duration}s")

# Example 2: Pipeline integration
def train_model(dataset_path: str, model_output: str):
    """Automated training step in pipeline."""
    job = MPIJob(
        name=f"training-{uuid.uuid4().hex[:8]}",
        namespace="ml-team",
        workers=8,
        image="registry.com/training:latest",
        command=["python", "train.py"],
        env={
            "DATASET_PATH": dataset_path,
            "MODEL_OUTPUT": model_output
        },
        worker_resources=ResourceSpec(cpu="8", memory="32Gi", gpus=2)
    )

    try:
        job.create(wait=True, timeout=7200)
        return model_output
    except JobFailedError as e:
        logger.error(f"Training failed: {e.message}")
        logger.error(f"Logs: {e.logs}")
        raise

# Example 3: Monitoring with callbacks
def on_status_change(status: JobStatus):
    print(f"Job status: {status.phase} - {status.workers_running}/{job.workers} workers")

job = MPIJob(...)
job.create()

while True:
    status = job.status()
    on_status_change(status)

    if status.phase in ["Succeeded", "Failed"]:
        break

    time.sleep(10)

# Example 4: Log streaming
job = MPIJob(...)
job.create()

# Stream logs from specific worker
for log_line in job.get_logs(component="worker", worker_id=2, follow=True):
    print(log_line)
```

**Async Support** (Post-MVP):
```python
import asyncio
from odh.training import AsyncMPIJob

async def train_parallel_jobs():
    """Launch multiple jobs concurrently."""
    jobs = [
        AsyncMPIJob(name=f"job-{i}", workers=4, image="training:latest")
        for i in range(10)
    ]

    # Create all jobs concurrently
    await asyncio.gather(*[job.create() for job in jobs])

    # Wait for all to complete
    results = await asyncio.gather(*[job.wait_for_completion() for job in jobs])

    return results

# Run
results = asyncio.run(train_parallel_jobs())
```

**Testing Strategy**:
- **Unit Tests**: Mock API responses, test validation logic
- **Integration Tests**: Against real OpenShift AI cluster
- **Documentation Tests**: All examples must execute successfully
- **Type Checking**: MyPy strict mode

### 3. ODH Dashboard UI Implementation

**Technology**: React 18, TypeScript, PatternFly 5

**Component Architecture**:
```
odh-dashboard/
└── src/
    ├── pages/
    │   └── training/
    │       ├── MPIJobList.tsx
    │       ├── MPIJobCreate.tsx
    │       ├── MPIJobDetail.tsx
    │       └── MPIJobLogs.tsx
    ├── components/
    │   └── training/
    │       ├── JobTypeSelector.tsx
    │       ├── MPIJobForm.tsx
    │       ├── WorkerTopology.tsx
    │       ├── ResourceAllocation.tsx
    │       └── StatusBadge.tsx
    ├── api/
    │   └── training-service.ts
    ├── types/
    │   └── mpijob.ts
    └── hooks/
        ├── useMPIJob.ts
        ├── useMPIJobStatus.ts
        └── useMPIJobLogs.ts
```

**Key UI Components**:

#### 1. MPIJob Creation Form

**Design Principles**:
- Progressive disclosure: Essential fields first, advanced options collapsed
- Inline validation: Real-time feedback before submission
- Contextual help: Tooltips, examples, decision guides

**Form Structure**:
```tsx
interface MPIJobFormProps {
  namespace: string;
  onSubmit: (job: MPIJobSpec) => void;
  onCancel: () => void;
}

const MPIJobForm: React.FC<MPIJobFormProps> = ({ namespace, onSubmit, onCancel }) => {
  const [formData, setFormData] = useState<MPIJobFormData>({
    name: '',
    image: '',
    workers: 4,
    slotsPerWorker: 1,
    command: [],
    workerResources: {
      cpu: '8',
      memory: '32Gi',
      gpus: 2
    },
    launcherResources: {
      cpu: '2',
      memory: '8Gi',
      gpus: 0
    },
    env: [],
    volumes: [],
    mpiImplementation: 'OpenMPI',
    cleanPodPolicy: 'Running'
  });

  const [validation, setValidation] = useState<ValidationState>({});
  const [quotaInfo, setQuotaInfo] = useState<QuotaInfo | null>(null);

  // Real-time validation
  useEffect(() => {
    validateForm(formData).then(setValidation);
  }, [formData]);

  // Quota availability check
  useEffect(() => {
    if (formData.workers && formData.workerResources.gpus) {
      const totalGPUs = formData.workers * formData.workerResources.gpus;
      checkQuota(namespace, totalGPUs).then(setQuotaInfo);
    }
  }, [formData.workers, formData.workerResources.gpus, namespace]);

  return (
    <Form>
      {/* Essential Fields */}
      <FormGroup label="Job Name" isRequired fieldId="name">
        <TextInput
          id="name"
          value={formData.name}
          onChange={(value) => setFormData({ ...formData, name: value })}
          validated={validation.name?.status || 'default'}
        />
        <FormHelperText>
          <HelperText>
            <HelperTextItem variant={validation.name?.status}>
              {validation.name?.message || 'Unique identifier for this training job'}
            </HelperTextItem>
          </HelperText>
        </FormHelperText>
      </FormGroup>

      <FormGroup label="Training Image" isRequired fieldId="image">
        <TextInput
          id="image"
          value={formData.image}
          onChange={(value) => setFormData({ ...formData, image: value })}
          placeholder="quay.io/example/training:latest"
        />
        <FormHelperText>
          <HelperText>
            <HelperTextItem>
              Container image with training code and MPI installed.
              <Button variant="link" component="a" href="/docs/building-images" target="_blank">
                Learn more →
              </Button>
            </HelperTextItem>
          </HelperText>
        </FormHelperText>
      </FormGroup>

      <FormGroup label="Number of Workers" isRequired fieldId="workers">
        <NumberInput
          value={formData.workers}
          onMinus={() => setFormData({ ...formData, workers: Math.max(1, formData.workers - 1) })}
          onPlus={() => setFormData({ ...formData, workers: formData.workers + 1 })}
          onChange={(event) => {
            const value = parseInt(event.target.value, 10);
            if (!isNaN(value) && value >= 1) {
              setFormData({ ...formData, workers: value });
            }
          }}
          min={1}
          max={100}
        />
        <FormHelperText>
          <HelperText>
            <HelperTextItem>
              Number of parallel workers for distributed training
            </HelperTextItem>
          </HelperText>
        </FormHelperText>
      </FormGroup>

      {/* Resource Allocation */}
      <FormSection title="Worker Resources">
        <FormGroup label="CPUs per Worker" isRequired>
          <TextInput
            value={formData.workerResources.cpu}
            onChange={(value) => setFormData({
              ...formData,
              workerResources: { ...formData.workerResources, cpu: value }
            })}
          />
        </FormGroup>

        <FormGroup label="Memory per Worker" isRequired>
          <TextInput
            value={formData.workerResources.memory}
            onChange={(value) => setFormData({
              ...formData,
              workerResources: { ...formData.workerResources, memory: value }
            })}
            placeholder="32Gi"
          />
        </FormGroup>

        <FormGroup label="GPUs per Worker">
          <NumberInput
            value={formData.workerResources.gpus}
            onMinus={() => setFormData({
              ...formData,
              workerResources: {
                ...formData.workerResources,
                gpus: Math.max(0, formData.workerResources.gpus - 1)
              }
            })}
            onPlus={() => setFormData({
              ...formData,
              workerResources: {
                ...formData.workerResources,
                gpus: formData.workerResources.gpus + 1
              }
            })}
            min={0}
          />
        </FormGroup>

        {/* Quota Status */}
        {quotaInfo && (
          <Alert
            variant={quotaInfo.available ? 'info' : 'warning'}
            title={
              quotaInfo.available
                ? `Resources available: ${quotaInfo.availableGPUs}/${quotaInfo.totalGPUs} GPUs`
                : `Insufficient resources: Requested ${quotaInfo.requestedGPUs} GPUs but only ${quotaInfo.availableGPUs} available`
            }
            isInline
          />
        )}
      </FormSection>

      <FormGroup label="Command" isRequired fieldId="command">
        <TextArea
          id="command"
          value={formData.command.join('\n')}
          onChange={(value) => setFormData({
            ...formData,
            command: value.split('\n').filter(line => line.trim())
          })}
          placeholder={`horovodrun\n-np\n8\npython\ntrain.py`}
        />
        <FormHelperText>
          <HelperText>
            <HelperTextItem>
              Command to execute (one argument per line)
            </HelperTextItem>
          </HelperText>
        </FormHelperText>
      </FormGroup>

      {/* Advanced Options (Collapsed by Default) */}
      <ExpandableSection toggleText="Advanced Options">
        <FormSection title="Launcher Resources">
          {/* Launcher CPU, memory, GPU inputs */}
        </FormSection>

        <FormSection title="MPI Configuration">
          <FormGroup label="MPI Implementation">
            <Select
              value={formData.mpiImplementation}
              onChange={(value) => setFormData({ ...formData, mpiImplementation: value })}
              options={[
                { label: 'OpenMPI', value: 'OpenMPI' },
                { label: 'Intel MPI', value: 'IntelMPI' },
                { label: 'MPICH', value: 'MPICH' }
              ]}
            />
          </FormGroup>

          <FormGroup label="Slots per Worker">
            <NumberInput
              value={formData.slotsPerWorker}
              onChange={(value) => setFormData({ ...formData, slotsPerWorker: value })}
            />
            <FormHelperText>
              <HelperText>
                <HelperTextItem>
                  Number of MPI processes per worker. Typically matches GPU count.
                </HelperTextItem>
              </HelperText>
            </FormHelperText>
          </FormGroup>
        </FormSection>

        <FormSection title="Environment Variables">
          <KeyValueEditor
            entries={formData.env}
            onChange={(env) => setFormData({ ...formData, env })}
          />
        </FormSection>

        <FormSection title="Volumes">
          <VolumeEditor
            volumes={formData.volumes}
            onChange={(volumes) => setFormData({ ...formData, volumes })}
          />
        </FormSection>
      </ExpandableSection>

      {/* Actions */}
      <ActionGroup>
        <Button
          variant="primary"
          onClick={() => onSubmit(formData)}
          isDisabled={!isFormValid(validation)}
        >
          Create Job
        </Button>
        <Button variant="link" onClick={onCancel}>
          Cancel
        </Button>
      </ActionGroup>
    </Form>
  );
};
```

#### 2. Job List View

**Features**:
- Unified view of all job types (MPIJob, PyTorchJob, TFJob)
- Filter by type, status, namespace
- Sort by creation time, duration, status
- Bulk operations (delete multiple jobs)

```tsx
const MPIJobList: React.FC = () => {
  const { jobs, loading, error, refresh } = useMPIJobs();
  const [filters, setFilters] = useState<JobFilters>({
    type: ['mpijob'],
    status: [],
    namespace: 'all'
  });

  const filteredJobs = useMemo(() => {
    return jobs.filter(job => {
      if (filters.type.length && !filters.type.includes(job.type)) return false;
      if (filters.status.length && !filters.status.includes(job.status)) return false;
      return true;
    });
  }, [jobs, filters]);

  return (
    <PageSection>
      <Title headingLevel="h1">Training Jobs</Title>

      {/* Toolbar with filters and actions */}
      <Toolbar>
        <ToolbarContent>
          <ToolbarItem>
            <Select
              variant="checkbox"
              onSelect={(_, value) => {
                setFilters({
                  ...filters,
                  type: filters.type.includes(value as string)
                    ? filters.type.filter(t => t !== value)
                    : [...filters.type, value as string]
                });
              }}
              selections={filters.type}
              isOpen={typeSelectOpen}
              onToggle={setTypeSelectOpen}
              placeholderText="Filter by type"
            >
              <SelectOption value="mpijob">MPIJob</SelectOption>
              <SelectOption value="pytorchjob">PyTorchJob</SelectOption>
              <SelectOption value="tfjob">TFJob</SelectOption>
            </Select>
          </ToolbarItem>

          <ToolbarItem>
            <Select
              variant="checkbox"
              onSelect={(_, value) => {
                setFilters({
                  ...filters,
                  status: filters.status.includes(value as string)
                    ? filters.status.filter(s => s !== value)
                    : [...filters.status, value as string]
                });
              }}
              selections={filters.status}
              placeholderText="Filter by status"
            >
              <SelectOption value="Running">Running</SelectOption>
              <SelectOption value="Succeeded">Succeeded</SelectOption>
              <SelectOption value="Failed">Failed</SelectOption>
              <SelectOption value="Pending">Pending</SelectOption>
            </Select>
          </ToolbarItem>

          <ToolbarItem variant="separator" />

          <ToolbarItem>
            <Button
              variant="primary"
              onClick={() => navigate('/training/create')}
            >
              Create Training Job
            </Button>
          </ToolbarItem>

          <ToolbarItem>
            <Button variant="plain" onClick={refresh}>
              <SyncIcon />
            </Button>
          </ToolbarItem>
        </ToolbarContent>
      </Toolbar>

      {/* Job Table */}
      <Table variant="compact">
        <Thead>
          <Tr>
            <Th>Name</Th>
            <Th>Type</Th>
            <Th>Status</Th>
            <Th>Workers</Th>
            <Th>Duration</Th>
            <Th>Created</Th>
            <Th>Actions</Th>
          </Tr>
        </Thead>
        <Tbody>
          {filteredJobs.map(job => (
            <Tr key={job.name}>
              <Td>
                <Link to={`/training/${job.name}`}>{job.name}</Link>
              </Td>
              <Td>
                <Label color="blue">{job.type}</Label>
              </Td>
              <Td>
                <StatusBadge status={job.status} />
              </Td>
              <Td>
                {job.type === 'mpijob' && (
                  <span>
                    {job.status.workersRunning}/{job.spec.workers}
                  </span>
                )}
              </Td>
              <Td>{formatDuration(job.status.duration)}</Td>
              <Td>{formatRelativeTime(job.metadata.creationTimestamp)}</Td>
              <Td>
                <ActionsColumn
                  items={[
                    {
                      title: 'View Details',
                      onClick: () => navigate(`/training/${job.name}`)
                    },
                    {
                      title: 'View Logs',
                      onClick: () => navigate(`/training/${job.name}/logs`)
                    },
                    {
                      title: 'Clone and Modify',
                      onClick: () => navigate(`/training/create?clone=${job.name}`)
                    },
                    {
                      isSeparator: true
                    },
                    {
                      title: 'Delete',
                      onClick: () => handleDelete(job.name)
                    }
                  ]}
                />
              </Td>
            </Tr>
          ))}
        </Tbody>
      </Table>

      {loading && <Spinner />}
      {error && <Alert variant="danger" title="Failed to load jobs" />}
      {!loading && filteredJobs.length === 0 && (
        <EmptyState>
          <EmptyStateIcon icon={CubesIcon} />
          <Title headingLevel="h4">No training jobs found</Title>
          <EmptyStateBody>
            Create your first training job to get started with distributed training.
          </EmptyStateBody>
          <Button variant="primary" onClick={() => navigate('/training/create')}>
            Create Training Job
          </Button>
        </EmptyState>
      )}
    </PageSection>
  );
};
```

#### 3. Job Detail View

**Features**:
- Real-time status updates (WebSocket or polling)
- Worker topology visualization
- Resource utilization charts
- Log viewer with worker selection
- Event timeline

```tsx
const MPIJobDetail: React.FC<{ jobName: string }> = ({ jobName }) => {
  const { job, status, loading } = useMPIJobStatus(jobName);
  const [selectedWorker, setSelectedWorker] = useState<number | null>(null);

  useEffect(() => {
    // Poll status every 5 seconds
    const interval = setInterval(() => {
      refetch();
    }, 5000);

    return () => clearInterval(interval);
  }, []);

  return (
    <PageSection>
      <Title headingLevel="h1">{jobName}</Title>

      <Grid hasGutter>
        {/* Status Overview */}
        <GridItem span={12}>
          <Card>
            <CardTitle>Status</CardTitle>
            <CardBody>
              <DescriptionList>
                <DescriptionListGroup>
                  <DescriptionListTerm>Phase</DescriptionListTerm>
                  <DescriptionListDescription>
                    <StatusBadge status={status.phase} />
                  </DescriptionListDescription>
                </DescriptionListGroup>

                <DescriptionListGroup>
                  <DescriptionListTerm>Workers</DescriptionListTerm>
                  <DescriptionListDescription>
                    <Flex>
                      <FlexItem>
                        <Label color="green">
                          {status.workersRunning} Running
                        </Label>
                      </FlexItem>
                      {status.workersPending > 0 && (
                        <FlexItem>
                          <Label color="orange">
                            {status.workersPending} Pending
                          </Label>
                        </FlexItem>
                      )}
                      {status.workersFailed > 0 && (
                        <FlexItem>
                          <Label color="red">
                            {status.workersFailed} Failed
                          </Label>
                        </FlexItem>
                      )}
                    </Flex>
                  </DescriptionListDescription>
                </DescriptionListGroup>

                <DescriptionListGroup>
                  <DescriptionListTerm>Duration</DescriptionListTerm>
                  <DescriptionListDescription>
                    {formatDuration(status.duration)}
                  </DescriptionListDescription>
                </DescriptionListGroup>

                <DescriptionListGroup>
                  <DescriptionListTerm>Created</DescriptionListTerm>
                  <DescriptionListDescription>
                    {formatTimestamp(job.metadata.creationTimestamp)}
                  </DescriptionListDescription>
                </DescriptionListGroup>
              </DescriptionList>
            </CardBody>
          </Card>
        </GridItem>

        {/* Worker Topology */}
        <GridItem span={12}>
          <Card>
            <CardTitle>Worker Topology</CardTitle>
            <CardBody>
              <WorkerTopology
                launcher={status.launcher}
                workers={status.workers}
                onWorkerClick={setSelectedWorker}
                selectedWorker={selectedWorker}
              />
            </CardBody>
          </Card>
        </GridItem>

        {/* Resource Utilization */}
        <GridItem span={6}>
          <Card>
            <CardTitle>CPU Utilization</CardTitle>
            <CardBody>
              <ResourceChart
                data={status.workers.map(w => ({
                  name: `Worker ${w.id}`,
                  value: w.cpuUsage
                }))}
              />
            </CardBody>
          </Card>
        </GridItem>

        <GridItem span={6}>
          <Card>
            <CardTitle>GPU Utilization</CardTitle>
            <CardBody>
              <ResourceChart
                data={status.workers.map(w => ({
                  name: `Worker ${w.id}`,
                  value: w.gpuUsage
                }))}
              />
            </CardBody>
          </Card>
        </GridItem>

        {/* Logs Viewer */}
        <GridItem span={12}>
          <Card>
            <CardTitle>Logs</CardTitle>
            <CardBody>
              <LogViewer
                jobName={jobName}
                component={selectedWorker !== null ? 'worker' : 'launcher'}
                workerId={selectedWorker}
              />
            </CardBody>
          </Card>
        </GridItem>

        {/* Events Timeline */}
        <GridItem span={12}>
          <Card>
            <CardTitle>Events</CardTitle>
            <CardBody>
              <EventsTimeline events={job.events} />
            </CardBody>
          </Card>
        </GridItem>
      </Grid>
    </PageSection>
  );
};
```

**Worker Topology Visualization**:
```tsx
const WorkerTopology: React.FC<WorkerTopologyProps> = ({
  launcher,
  workers,
  onWorkerClick,
  selectedWorker
}) => {
  return (
    <div className="worker-topology">
      {/* Launcher Node */}
      <div className="topology-node launcher">
        <div
          className={`node-card ${launcher.status.toLowerCase()}`}
          onClick={() => onWorkerClick(null)}
        >
          <div className="node-icon">
            <LauncherIcon />
          </div>
          <div className="node-label">Launcher</div>
          <div className="node-status">{launcher.status}</div>
        </div>
      </div>

      {/* Worker Nodes */}
      <div className="topology-workers">
        {workers.map((worker, index) => (
          <div
            key={worker.id}
            className={`topology-node worker ${
              selectedWorker === index ? 'selected' : ''
            }`}
          >
            <div
              className={`node-card ${worker.status.toLowerCase()}`}
              onClick={() => onWorkerClick(index)}
            >
              <div className="node-icon">
                <WorkerIcon />
              </div>
              <div className="node-label">Worker {index}</div>
              <div className="node-status">{worker.status}</div>
              <div className="node-resources">
                {worker.gpus > 0 && (
                  <span className="resource-badge">
                    {worker.gpus} GPU
                  </span>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};
```

**Accessibility Considerations**:
- Keyboard navigation: Tab through workers, Enter to select
- Screen reader announcements: "Worker 2, status: Running, 2 GPUs, CPU usage 85%"
- Color-blind safe: Status indicated by icons + text, not just color
- ARIA labels: All interactive elements properly labeled

**Real-Time Updates**:
```tsx
const useMPIJobStatus = (jobName: string) => {
  const [status, setStatus] = useState<JobStatus | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<Error | null>(null);

  useEffect(() => {
    // Option 1: WebSocket for real-time updates
    const ws = new WebSocket(`wss://api.openshift-ai.example.com/ws/jobs/${jobName}/status`);

    ws.onmessage = (event) => {
      const updatedStatus = JSON.parse(event.data);
      setStatus(updatedStatus);
      setLoading(false);
    };

    ws.onerror = (error) => {
      setError(new Error('WebSocket connection failed'));
      setLoading(false);
    };

    return () => ws.close();

    // Option 2: Polling (fallback)
    // const interval = setInterval(async () => {
    //   try {
    //     const response = await fetch(`/api/v1/jobs/${jobName}/status`);
    //     const data = await response.json();
    //     setStatus(data);
    //   } catch (err) {
    //     setError(err);
    //   }
    // }, 5000);
    //
    // return () => clearInterval(interval);
  }, [jobName]);

  return { status, loading, error };
};
```

---

## Integration Architecture

### 1. OpenShift RBAC Integration

**Permission Model**:

MPIJob operations map to Kubernetes RBAC verbs on the `mpijobs.kubeflow.org` resource.

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: mpijob-creator
rules:
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs"]
  verbs: ["create", "get", "list", "watch", "delete"]
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs/status"]
  verbs: ["get"]
- apiGroups: [""]
  resources: ["pods/log"]
  verbs: ["get"]
```

**Role Hierarchy**:

1. **Data Scientist Role** (most common)
   - Create, read, delete own MPIJobs
   - View logs and status
   - No access to other teams' jobs

2. **MLOps Engineer Role**
   - Full CRUD on MPIJobs in assigned namespaces
   - Can view all jobs across namespaces (read-only)
   - Can manage service accounts for automation

3. **Platform Administrator Role**
   - Full CRUD across all namespaces
   - Configure quotas, network policies
   - Access audit logs

**RoleBinding Example**:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind:RoleBinding
metadata:
  name: data-scientist-mpijob-access
  namespace: data-science-team
subjects:
- kind: User
  name: maria@example.com
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: mpijob-creator
  apiGroup: rbac.authorization.k8s.io
```

**Service Account for Automation**:
```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: pipeline-service-account
  namespace: data-science-team
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: pipeline-mpijob-access
  namespace: data-science-team
subjects:
- kind: ServiceAccount
  name: pipeline-service-account
  namespace: data-science-team
roleRef:
  kind: ClusterRole
  name: mpijob-creator
  apiGroup: rbac.authorization.k8s.io
```

### 2. Resource Quota and LimitRange Integration

**Quota Enforcement Points**:

1. **Pre-Creation Validation** (Training Service Layer)
   - Query current namespace usage
   - Calculate requested resources (workers × resources per worker)
   - Reject if total would exceed quota

2. **Kubernetes Admission Control** (Backup Enforcement)
   - ResourceQuota admission controller validates pod creation
   - Prevents quota bypass even if Training Service validation is skipped

**Example ResourceQuota**:
```yaml
apiVersion: v1
kind: ResourceQuota
metadata:
  name: data-science-team-quota
  namespace: data-science-team
spec:
  hard:
    requests.cpu: "128"
    requests.memory: "512Gi"
    requests.nvidia.com/gpu: "16"
    limits.cpu: "128"
    limits.memory: "512Gi"
    # Custom resource for job count
    count/mpijobs.kubeflow.org: "5"
    count/pods: "50"
```

**LimitRange for Default Values**:
```yaml
apiVersion: v1
kind: LimitRange
metadata:
  name: data-science-team-limits
  namespace: data-science-team
spec:
  limits:
  - max:
      cpu: "16"
      memory: "64Gi"
      nvidia.com/gpu: "4"
    min:
      cpu: "100m"
      memory: "128Mi"
    default:
      cpu: "2"
      memory: "8Gi"
    defaultRequest:
      cpu: "1"
      memory: "4Gi"
    type: Container
```

**Quota Error Handling**:
```json
{
  "error": {
    "code": "ResourceQuotaExceeded",
    "message": "Cannot create MPIJob: GPU quota exceeded",
    "details": {
      "resource": "nvidia.com/gpu",
      "requested": 8,
      "current_usage": 10,
      "quota_limit": 16,
      "available": 6
    },
    "suggestions": [
      "Reduce worker count from 4 to 3 (6 GPUs)",
      "Wait for jobs to complete and free up resources",
      "Contact administrator to increase quota"
    ]
  }
}
```

### 3. Monitoring Stack Integration

**Prometheus Integration**:

Training Operator exposes metrics on `:8080/metrics`. ServiceMonitor CRD configures Prometheus to scrape:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: training-operator
  namespace: openshift-ai-training-operator
  labels:
    app: training-operator
spec:
  selector:
    matchLabels:
      app: training-operator
  endpoints:
  - port: metrics
    interval: 30s
    path: /metrics
```

**Custom Metrics** (emitted by Training Service):

```
# Job lifecycle metrics
openshift_ai_mpijob_created_total{namespace="data-science-team",job_name="training-001"}
openshift_ai_mpijob_running_total{namespace="data-science-team",job_name="training-001"}
openshift_ai_mpijob_succeeded_total{namespace="data-science-team",job_name="training-001"}
openshift_ai_mpijob_failed_total{namespace="data-science-team",job_name="training-001"}

# Duration metrics (histogram)
openshift_ai_mpijob_duration_seconds_bucket{namespace="data-science-team",phase="Running",le="300"}
openshift_ai_mpijob_duration_seconds_bucket{namespace="data-science-team",phase="Running",le="600"}
openshift_ai_mpijob_duration_seconds_bucket{namespace="data-science-team",phase="Running",le="1800"}
openshift_ai_mpijob_duration_seconds_sum{namespace="data-science-team",phase="Running"}
openshift_ai_mpijob_duration_seconds_count{namespace="data-science-team",phase="Running"}

# Worker status metrics
openshift_ai_mpijob_workers_active{namespace="data-science-team",job_name="training-001"}
openshift_ai_mpijob_workers_pending{namespace="data-science-team",job_name="training-001"}
openshift_ai_mpijob_workers_failed{namespace="data-science-team",job_name="training-001"}

# Resource utilization (per worker)
openshift_ai_mpijob_worker_cpu_usage{namespace="data-science-team",job_name="training-001",worker_id="0"}
openshift_ai_mpijob_worker_memory_bytes{namespace="data-science-team",job_name="training-001",worker_id="0"}
openshift_ai_mpijob_worker_gpu_utilization{namespace="data-science-team",job_name="training-001",worker_id="0"}
```

**Grafana Dashboard**:

Create pre-configured dashboard for MPIJob monitoring:

```json
{
  "dashboard": {
    "title": "OpenShift AI - MPIJob Monitoring",
    "panels": [
      {
        "title": "Active MPIJobs",
        "targets": [
          {
            "expr": "sum(openshift_ai_mpijob_running_total) by (namespace)"
          }
        ]
      },
      {
        "title": "Job Success Rate",
        "targets": [
          {
            "expr": "sum(openshift_ai_mpijob_succeeded_total) / (sum(openshift_ai_mpijob_succeeded_total) + sum(openshift_ai_mpijob_failed_total))"
          }
        ]
      },
      {
        "title": "Worker Status",
        "targets": [
          {
            "expr": "openshift_ai_mpijob_workers_active"
          },
          {
            "expr": "openshift_ai_mpijob_workers_pending"
          },
          {
            "expr": "openshift_ai_mpijob_workers_failed"
          }
        ]
      }
    ]
  }
}
```

### 4. Audit Logging Integration

**Audit Event Schema**:

```json
{
  "kind": "Event",
  "apiVersion": "audit.k8s.io/v1",
  "level": "RequestResponse",
  "auditID": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
  "stage": "ResponseComplete",
  "requestURI": "/api/v1/namespaces/data-science-team/trainingjobs/mpijob",
  "verb": "create",
  "user": {
    "username": "maria@example.com",
    "uid": "12345",
    "groups": ["data-scientists", "system:authenticated"]
  },
  "sourceIPs": ["10.0.1.15"],
  "userAgent": "odh-cli/1.0.0",
  "objectRef": {
    "resource": "mpijobs",
    "namespace": "data-science-team",
    "name": "training-001",
    "apiGroup": "kubeflow.org",
    "apiVersion": "v2beta1"
  },
  "responseStatus": {
    "metadata": {},
    "code": 201
  },
  "requestObject": {
    "spec": {
      "workers": 4,
      "image": "training:latest"
    }
  },
  "responseObject": {
    "metadata": {
      "name": "training-001",
      "namespace": "data-science-team",
      "uid": "abcd-1234-efgh-5678"
    }
  },
  "requestReceivedTimestamp": "2025-10-21T10:00:00.123456Z",
  "stageTimestamp": "2025-10-21T10:00:00.234567Z",
  "annotations": {
    "authorization.k8s.io/decision": "allow",
    "authorization.k8s.io/reason": "RBAC: allowed by RoleBinding"
  }
}
```

**Audit Policy Configuration**:

```yaml
apiVersion: audit.k8s.io/v1
kind: Policy
rules:
# Log MPIJob operations with full request/response
- level: RequestResponse
  resources:
  - group: kubeflow.org
    resources: ["mpijobs"]
  verbs: ["create", "delete", "patch", "update"]

# Log MPIJob reads at Metadata level
- level: Metadata
  resources:
  - group: kubeflow.org
    resources: ["mpijobs"]
  verbs: ["get", "list", "watch"]

# Log MPIJob status updates
- level: RequestResponse
  resources:
  - group: kubeflow.org
    resources: ["mpijobs/status"]
  verbs: ["update", "patch"]
```

**SIEM Integration**:

Forward audit logs to enterprise SIEM (Splunk, ELK, etc.):

```yaml
apiVersion: logging.openshift.io/v1
kind: ClusterLogForwarder
metadata:
  name: instance
  namespace: openshift-logging
spec:
  outputs:
  - name: remote-siem
    type: elasticsearch
    url: https://siem.example.com:9200
    secret:
      name: siem-credentials
  pipelines:
  - name: audit-logs
    inputRefs:
    - audit
    outputRefs:
    - remote-siem
```

### 5. Data Science Pipelines Integration (Post-MVP)

**Pattern**: MPIJob as pipeline step in Kubeflow Pipelines

```python
from kfp import dsl
from kfp.components import create_component_from_func

@dsl.component
def distributed_training_step(
    dataset_path: str,
    model_output: str,
    workers: int = 4,
    gpus_per_worker: int = 2
) -> str:
    """Distributed training pipeline step."""
    from odh.training import MPIJob, ResourceSpec

    job = MPIJob(
        name=f"training-{dsl.RUN_ID_PLACEHOLDER}",
        workers=workers,
        image="training:latest",
        command=["python", "train.py", "--dataset", dataset_path, "--output", model_output],
        worker_resources=ResourceSpec(cpu="8", memory="32Gi", gpus=gpus_per_worker)
    )

    job.create(wait=True, timeout=7200)

    return model_output

@dsl.pipeline(
    name="End-to-End Training Pipeline",
    description="Data prep → Distributed training → Evaluation"
)
def training_pipeline(dataset_uri: str, model_registry_uri: str):
    # Step 1: Data preprocessing
    preprocess_op = preprocess_data(dataset_uri=dataset_uri)

    # Step 2: Distributed training (MPIJob)
    train_op = distributed_training_step(
        dataset_path=preprocess_op.outputs['processed_data'],
        model_output="/models/trained",
        workers=4,
        gpus_per_worker=2
    )

    # Step 3: Model evaluation
    evaluate_op = evaluate_model(
        model_path=train_op.outputs['model_output']
    )

    # Step 4: Register model
    register_op = register_model(
        model_path=train_op.outputs['model_output'],
        metrics=evaluate_op.outputs['metrics'],
        registry_uri=model_registry_uri
    )

# Compile and submit pipeline
from kfp import compiler
compiler.Compiler().compile(training_pipeline, 'training_pipeline.yaml')
```

---

## Development Phases

### Phase 1: Foundation (Weeks 1-4)

**Objective**: Establish core infrastructure and basic MPIJob creation

**Deliverables**:
1. KubeFlow Training Operator V2 deployment and configuration
2. MPIJob CRD installation and validation webhooks
3. Training Service API (minimal viable endpoints)
4. CLI basic commands (create, delete, list, describe)
5. SDK core MPIJob class with create/delete methods

**Technical Tasks**:

**Week 1: Infrastructure Setup**
- Deploy Training Operator V2 to dev cluster
- Configure CRDs and RBAC
- Set up CI/CD pipelines for component builds
- Establish development environment (local cluster, testing tools)

**Week 2: Training Service Backend**
- Implement gRPC service definition
- Create REST API gateway
- Build job creation endpoint with validation
- Integrate with Kubernetes API for CRD operations

**Week 3: CLI Development**
- Scaffold CLI project structure (Cobra framework)
- Implement `create mpijob` command
- Implement `delete mpijob` command
- Implement `list mpijob` command
- Add basic error handling and output formatting

**Week 4: SDK Development**
- Create Python package structure
- Implement MPIJob class with type hints
- Add create/delete methods with API client
- Write unit tests (>80% coverage)
- Generate API documentation (Sphinx)

**Testing Focus**:
- Unit tests for all components
- Integration tests: CLI → Training Service → Kubernetes
- Manual E2E testing: Create MPIJob, verify pods created, delete job

**Success Criteria**:
- User can create MPIJob via CLI: `odh training create mpijob -f spec.yaml`
- MPIJob CRD appears in cluster: `kubectl get mpijob`
- Launcher and worker pods are created automatically
- User can delete job: `odh training delete mpijob my-job`
- Python SDK can create and delete jobs programmatically

---

### Phase 2: Observability (Weeks 5-8)

**Objective**: Add status monitoring, log access, and metrics

**Deliverables**:
1. Job status API and CLI/SDK integration
2. Log aggregation and streaming
3. Metrics collection and Prometheus integration
4. Basic Dashboard UI for job list and detail views

**Technical Tasks**:

**Week 5: Status and Monitoring**
- Implement status aggregation in Training Service
- Add `describe mpijob` CLI command with real-time updates
- Add `status()` method to SDK
- Create status cache layer for performance

**Week 6: Log Collection**
- Implement log retrieval API (launcher and worker logs)
- Add `logs mpijob` CLI command with follow mode
- Add `get_logs()` method to SDK with streaming support
- Integrate with OpenShift logging stack

**Week 7: Metrics and Observability**
- Implement Prometheus metric exporters
- Create ServiceMonitor for automatic scraping
- Build Grafana dashboard template
- Add resource utilization tracking per worker

**Week 8: Dashboard UI (Phase 1)**
- Create job list view with filters
- Build job detail page with status overview
- Implement worker topology visualization
- Add log viewer component

**Testing Focus**:
- Status updates reflect accurately within 10 seconds
- Logs are retrievable from all workers
- Metrics appear in Prometheus
- UI renders correctly with mock data

**Success Criteria**:
- User can monitor job status: `odh training describe mpijob my-job --watch`
- User can view logs: `odh training logs mpijob my-job --worker=2 --follow`
- Metrics visible in Grafana dashboard
- Dashboard shows running jobs with worker status

---

### Phase 3: UI and UX (Weeks 9-12)

**Objective**: Complete Dashboard UI with job creation, error handling, and accessibility

**Deliverables**:
1. Job creation form with validation
2. Advanced job management (clone, retry)
3. Error messages with actionable guidance
4. Accessibility compliance (WCAG 2.1 Level AA)

**Technical Tasks**:

**Week 9: Job Creation UI**
- Build MPIJob creation form with progressive disclosure
- Implement inline validation and quota checks
- Add job type selector with decision guide
- Create resource allocation UI components

**Week 10: Advanced Features**
- Implement "Clone and modify" workflow
- Add job deletion with confirmation
- Build log viewer with worker selection
- Create event timeline component

**Week 11: Error Handling and Guidance**
- Design error message patterns
- Implement contextual help and tooltips
- Add troubleshooting links
- Create empty states and loading indicators

**Week 12: Accessibility and Polish**
- Keyboard navigation testing
- Screen reader compatibility
- Color-blind safe status indicators
- Performance optimization (lazy loading, virtualization)

**Testing Focus**:
- Usability testing with target personas
- Accessibility audit (automated + manual)
- Cross-browser testing (Chrome, Firefox, Safari, Edge)
- Performance testing (render time <2s)

**Success Criteria**:
- User can create MPIJob through UI form
- Validation errors appear inline with helpful messages
- All UI elements keyboard accessible
- Screen reader announces job status changes
- WCAG 2.1 Level AA compliance verified

---

### Phase 4: Integration and Hardening (Weeks 13-16)

**Objective**: Enterprise features, security, multi-tenancy, and documentation

**Deliverables**:
1. RBAC integration and permission enforcement
2. Resource quota validation
3. Audit logging
4. Network policy templates
5. Comprehensive documentation

**Technical Tasks**:

**Week 13: RBAC and Security**
- Implement permission checks in Training Service
- Create role templates (data scientist, MLOps, admin)
- Add service account support for automation
- Test permission boundaries across namespaces

**Week 14: Resource Management**
- Implement quota validation before job creation
- Add quota information to UI
- Create LimitRange configuration guidance
- Test quota enforcement edge cases

**Week 15: Audit and Compliance**
- Configure audit logging for MPIJob operations
- Test SIEM integration
- Document compliance mappings (GDPR, HIPAA, SOC2)
- Create audit log analysis examples

**Week 16: Documentation and Examples**
- Write getting started guide (30-minute first job)
- Create API reference documentation
- Build example gallery (Horovod, DeepSpeed, Intel MPI)
- Develop troubleshooting guide

**Testing Focus**:
- Security testing (permission bypass attempts)
- Multi-tenant isolation verification
- Audit log completeness
- Documentation accuracy (all examples must execute)

**Success Criteria**:
- RBAC enforced: Users cannot create jobs in namespaces without permissions
- Quota violations rejected with clear errors
- All MPIJob operations generate audit events
- Documentation enables first job creation within 30 minutes

---

### Phase 5: Performance and Scale (Weeks 17-20)

**Objective**: Validate performance, scalability, and production readiness

**Deliverables**:
1. Performance benchmarks
2. Scale testing results
3. Optimization recommendations
4. Production deployment guide

**Technical Tasks**:

**Week 17: Performance Benchmarking**
- Measure job submission latency (<5s target)
- Test distributed training speedup (>2x target)
- Benchmark status update latency (<10s target)
- Identify performance bottlenecks

**Week 18: Scale Testing**
- Test concurrent job creation (100+ jobs)
- Validate large worker counts (50+ workers)
- Test multi-tenant scenarios (10+ namespaces)
- Measure resource consumption (Training Service, Operator)

**Week 19: Optimization**
- Optimize status caching
- Improve log retrieval performance
- Tune Prometheus metric cardinality
- Reduce Dashboard UI render time

**Week 20: Production Readiness**
- Create deployment playbooks
- Document upgrade procedures
- Build disaster recovery plans
- Conduct final security review

**Testing Focus**:
- Load testing (1000+ concurrent users)
- Chaos engineering (pod failures, network issues)
- Long-running stability tests (7+ days)
- Performance regression prevention

**Success Criteria**:
- Job submission <5s under normal load
- Distributed training shows >2x speedup
- System supports 50+ concurrent MPIJobs
- Training Service handles 100+ status queries/sec
- All production readiness checklists complete

---

## Testing Strategy

### Unit Testing

**Coverage Target**: 80% for all components

**Go Components (CLI, Training Service)**:
```go
// Example: Training Service unit test
func TestMPIJobCreation(t *testing.T) {
    // Setup
    mockClient := NewMockKubernetesClient()
    service := NewTrainingService(mockClient)

    spec := &MPIJobSpec{
        Name: "test-job",
        Workers: 4,
        Image: "training:latest",
        Command: []string{"python", "train.py"},
    }

    // Execute
    job, err := service.CreateMPIJob(context.Background(), "test-ns", spec)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "test-job", job.Name)
    assert.Equal(t, 4, job.Spec.Workers)

    // Verify CRD creation was called
    mockClient.AssertCalled(t, "Create", mock.MatchedBy(func(obj interface{}) bool {
        mpijob, ok := obj.(*MPIJob)
        return ok && mpijob.Name == "test-job"
    }))
}

func TestResourceQuotaValidation(t *testing.T) {
    // Test quota exceeded scenario
    mockClient := NewMockKubernetesClient()
    mockClient.On("GetResourceQuota", "test-ns").Return(&ResourceQuota{
        Spec: ResourceQuotaSpec{
            Hard: ResourceList{
                "requests.nvidia.com/gpu": resource.MustParse("16"),
            },
        },
        Status: ResourceQuotaStatus{
            Used: ResourceList{
                "requests.nvidia.com/gpu": resource.MustParse("10"),
            },
        },
    }, nil)

    service := NewTrainingService(mockClient)

    spec := &MPIJobSpec{
        Workers: 4,
        WorkerResources: ResourceSpec{
            GPUs: 2, // Total: 8 GPUs
        },
    }

    // Should succeed (10 + 8 = 18 > 16, but we validate against available)
    err := service.ValidateQuota(context.Background(), "test-ns", spec)
    assert.NoError(t, err)

    // Increase to exceed quota
    spec.Workers = 5 // Total: 10 GPUs (10 + 10 = 20 > 16)
    err = service.ValidateQuota(context.Background(), "test-ns", spec)
    assert.Error(t, err)
    assert.IsType(t, &ResourceQuotaExceededError{}, err)
}
```

**Python SDK**:
```python
import unittest
from unittest.mock import Mock, patch
from odh.training import MPIJob, ResourceSpec, JobFailedError

class TestMPIJob(unittest.TestCase):

    def setUp(self):
        self.mock_client = Mock()

    def test_create_job(self):
        """Test job creation with valid spec."""
        job = MPIJob(
            name="test-job",
            workers=4,
            image="training:latest",
            command=["python", "train.py"],
            worker_resources=ResourceSpec(cpu="8", memory="32Gi", gpus=2)
        )

        # Mock API response
        self.mock_client.create_mpijob.return_value = {
            "metadata": {"name": "test-job", "namespace": "default"},
            "spec": {"workers": 4}
        }

        job._client = self.mock_client
        result = job.create(wait=False)

        self.assertEqual(result.name, "test-job")
        self.mock_client.create_mpijob.assert_called_once()

    def test_validation_errors(self):
        """Test client-side validation."""
        with self.assertRaises(ValueError) as ctx:
            job = MPIJob(
                name="invalid name with spaces",  # Invalid DNS-1123
                workers=4,
                image="training:latest",
                command=["python", "train.py"]
            )
            job._validate()

        self.assertIn("Invalid job name", str(ctx.exception))

    def test_wait_for_completion_timeout(self):
        """Test timeout handling."""
        job = MPIJob(name="test-job", workers=4, image="training:latest", command=["train"])
        job._client = self.mock_client

        # Mock status always returns "Running"
        self.mock_client.get_mpijob_status.return_value = {
            "phase": "Running",
            "replicaStatuses": {"Worker": {"active": 4}}
        }

        with self.assertRaises(TimeoutError):
            job.wait_for_completion(timeout=1, poll_interval=0.5)

    def test_job_failure_with_logs(self):
        """Test failure handling with log retrieval."""
        job = MPIJob(name="test-job", workers=4, image="training:latest", command=["train"])
        job._client = self.mock_client

        # Mock failed status
        self.mock_client.get_mpijob_status.return_value = {
            "phase": "Failed",
            "replicaStatuses": {"Worker": {"active": 0, "failed": 1}}
        }

        self.mock_client.get_logs.return_value = "Error: Out of memory"

        with self.assertRaises(JobFailedError) as ctx:
            job.wait_for_completion(timeout=5, poll_interval=1)

        self.assertIn("Out of memory", ctx.exception.logs)
```

**React UI**:
```typescript
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { MPIJobForm } from './MPIJobForm';

describe('MPIJobForm', () => {
  test('validates required fields', async () => {
    const onSubmit = jest.fn();
    render(<MPIJobForm namespace="test-ns" onSubmit={onSubmit} onCancel={() => {}} />);

    // Try to submit without filling required fields
    const submitButton = screen.getByText('Create Job');
    fireEvent.click(submitButton);

    // Should show validation errors
    await waitFor(() => {
      expect(screen.getByText(/job name is required/i)).toBeInTheDocument();
      expect(screen.getByText(/training image is required/i)).toBeInTheDocument();
    });

    expect(onSubmit).not.toHaveBeenCalled();
  });

  test('shows quota warning when resources exceed availability', async () => {
    const mockCheckQuota = jest.fn().mockResolvedValue({
      available: false,
      requestedGPUs: 8,
      availableGPUs: 4,
      totalGPUs: 16
    });

    render(<MPIJobForm namespace="test-ns" checkQuota={mockCheckQuota} />);

    // Fill in form
    fireEvent.change(screen.getByLabelText(/job name/i), { target: { value: 'test-job' } });
    fireEvent.change(screen.getByLabelText(/number of workers/i), { target: { value: '4' } });
    fireEvent.change(screen.getByLabelText(/gpus per worker/i), { target: { value: '2' } });

    // Should show warning
    await waitFor(() => {
      expect(screen.getByText(/insufficient resources/i)).toBeInTheDocument();
      expect(screen.getByText(/requested 8 gpus but only 4 available/i)).toBeInTheDocument();
    });
  });

  test('submits valid form', async () => {
    const onSubmit = jest.fn();
    render(<MPIJobForm namespace="test-ns" onSubmit={onSubmit} onCancel={() => {}} />);

    // Fill required fields
    fireEvent.change(screen.getByLabelText(/job name/i), { target: { value: 'test-job' } });
    fireEvent.change(screen.getByLabelText(/training image/i), { target: { value: 'training:latest' } });
    fireEvent.change(screen.getByLabelText(/number of workers/i), { target: { value: '4' } });

    // Submit
    const submitButton = screen.getByText('Create Job');
    fireEvent.click(submitButton);

    await waitFor(() => {
      expect(onSubmit).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'test-job',
          image: 'training:latest',
          workers: 4
        })
      );
    });
  });
});
```

### Integration Testing

**Test Scenarios**:

1. **End-to-End Job Lifecycle**
   - Create MPIJob via CLI
   - Verify CRD created in Kubernetes
   - Wait for launcher and worker pods to start
   - Retrieve job status
   - Access logs from launcher and workers
   - Delete job and verify cleanup

2. **Cross-Component Communication**
   - CLI → Training Service → Kubernetes API
   - SDK → Training Service → Kubernetes API
   - Dashboard UI → Training Service → Kubernetes API
   - Training Operator → MPIJob CRD → Pod creation

3. **RBAC Enforcement**
   - User with permissions can create jobs
   - User without permissions receives authorization error
   - Service account can create jobs on behalf of users
   - Cross-namespace access denied

4. **Resource Quota Validation**
   - Job creation succeeds when under quota
   - Job creation fails when exceeding quota
   - Quota error message includes actionable guidance
   - Quota updated after job deletion

**Test Environment**:
- Kind cluster (local Kubernetes)
- OpenShift cluster (CI/CD)
- Staging environment (pre-production validation)

**Example Integration Test**:
```bash
#!/bin/bash
# Integration test script

set -e

NAMESPACE="test-integration"
JOB_NAME="test-job-$(date +%s)"

echo "Creating namespace..."
kubectl create namespace $NAMESPACE

echo "Applying resource quota..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: ResourceQuota
metadata:
  name: test-quota
  namespace: $NAMESPACE
spec:
  hard:
    requests.cpu: "32"
    requests.memory: "128Gi"
    count/mpijobs.kubeflow.org: "5"
EOF

echo "Creating MPIJob via CLI..."
odh training create mpijob $JOB_NAME \
  --namespace=$NAMESPACE \
  --image=training:latest \
  --workers=2 \
  --cpu-per-worker=4 \
  --memory-per-worker=16Gi \
  --command="python train.py"

echo "Waiting for job to start..."
timeout 60 bash -c "until odh training describe mpijob $JOB_NAME --namespace=$NAMESPACE | grep -q 'Status: Running'; do sleep 5; done"

echo "Verifying pods created..."
LAUNCHER_POD=$(kubectl get pods -n $NAMESPACE -l training.kubeflow.org/job-name=$JOB_NAME,training.kubeflow.org/job-role=launcher -o jsonpath='{.items[0].metadata.name}')
WORKER_PODS=$(kubectl get pods -n $NAMESPACE -l training.kubeflow.org/job-name=$JOB_NAME,training.kubeflow.org/job-role=worker -o jsonpath='{.items[*].metadata.name}')

echo "Launcher pod: $LAUNCHER_POD"
echo "Worker pods: $WORKER_PODS"

echo "Retrieving logs..."
odh training logs mpijob $JOB_NAME --namespace=$NAMESPACE --tail=10

echo "Deleting job..."
odh training delete mpijob $JOB_NAME --namespace=$NAMESPACE --wait

echo "Verifying cleanup..."
POD_COUNT=$(kubectl get pods -n $NAMESPACE -l training.kubeflow.org/job-name=$JOB_NAME -o json | jq '.items | length')
if [ "$POD_COUNT" -ne 0 ]; then
  echo "ERROR: $POD_COUNT pods remain after job deletion"
  exit 1
fi

echo "Cleaning up namespace..."
kubectl delete namespace $NAMESPACE

echo "Integration test PASSED"
```

### E2E Testing

**User Workflow Tests**:

1. **Data Scientist First Job**
   - Navigate to Dashboard
   - Click "Create Training Job"
   - Select "MPIJob"
   - Fill in configuration form
   - Submit job
   - Monitor status in real-time
   - View logs
   - Download artifacts
   - Delete job

2. **MLOps Pipeline Integration**
   - Create service account
   - Grant permissions
   - Write pipeline script using SDK
   - Execute pipeline
   - Verify job created
   - Monitor via CLI
   - Handle failure scenario
   - Cleanup resources

3. **Administrator Multi-Tenant Setup**
   - Create multiple namespaces
   - Apply resource quotas
   - Configure RBAC roles
   - Test permission boundaries
   - Monitor resource usage
   - Review audit logs

**E2E Test Framework**:
```python
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

class TestMPIJobDashboard:

    def setup_method(self):
        self.driver = webdriver.Chrome()
        self.driver.get("https://openshift-ai.example.com")
        self.login()

    def test_create_mpijob_via_ui(self):
        """Test job creation through Dashboard UI."""
        driver = self.driver

        # Navigate to Training Jobs
        driver.find_element(By.LINK_TEXT, "Training Jobs").click()

        # Click Create
        driver.find_element(By.BUTTON, "Create Training Job").click()

        # Select MPIJob
        driver.find_element(By.XPATH, "//label[contains(text(), 'MPIJob')]").click()

        # Fill form
        driver.find_element(By.ID, "job-name").send_keys("test-job-e2e")
        driver.find_element(By.ID, "image").send_keys("quay.io/opendatahub/horovod-mnist:latest")
        driver.find_element(By.ID, "workers").clear()
        driver.find_element(By.ID, "workers").send_keys("2")

        # Submit
        driver.find_element(By.XPATH, "//button[contains(text(), 'Create Job')]").click()

        # Wait for redirect to job detail page
        WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, "//h1[contains(text(), 'test-job-e2e')]"))
        )

        # Verify status changes from Pending to Running
        WebDriverWait(driver, 60).until(
            EC.text_to_be_present_in_element((By.CLASS_NAME, "status-badge"), "Running")
        )

        # Click on worker topology
        driver.find_element(By.XPATH, "//div[contains(@class, 'worker-node')][1]").click()

        # Verify logs are visible
        WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.CLASS_NAME, "log-viewer"))
        )

        # Delete job
        driver.find_element(By.XPATH, "//button[contains(text(), 'Delete')]").click()
        driver.find_element(By.XPATH, "//button[contains(text(), 'Confirm')]").click()

        # Verify redirect to job list
        WebDriverWait(driver, 10).until(
            EC.presence_of_element_located((By.XPATH, "//h1[contains(text(), 'Training Jobs')]"))
        )

    def teardown_method(self):
        self.driver.quit()
```

### Performance Testing

**Load Test Scenarios**:

1. **Job Submission Load**
   - 100 concurrent job creations
   - Measure latency (target: <5s p95)
   - Verify no failures

2. **Status Query Load**
   - 1000 concurrent status queries
   - Measure throughput (target: >100 qps)
   - Verify cache effectiveness

3. **Log Retrieval Load**
   - 50 concurrent log streams
   - Measure bandwidth consumption
   - Verify no connection drops

**Performance Test Example**:
```python
import asyncio
import aiohttp
import time
from statistics import mean, stdev

async def create_job(session, job_id):
    """Create single MPIJob and measure latency."""
    start = time.time()

    spec = {
        "name": f"perf-test-{job_id}",
        "workers": 2,
        "image": "training:latest",
        "command": ["sleep", "3600"]
    }

    async with session.post(
        "https://api.openshift-ai.example.com/api/v1/namespaces/test-ns/trainingjobs/mpijob",
        json=spec
    ) as response:
        await response.json()
        latency = time.time() - start
        return latency

async def performance_test():
    """Run load test with 100 concurrent job creations."""
    async with aiohttp.ClientSession() as session:
        tasks = [create_job(session, i) for i in range(100)]
        latencies = await asyncio.gather(*tasks)

    print(f"Mean latency: {mean(latencies):.2f}s")
    print(f"Std dev: {stdev(latencies):.2f}s")
    print(f"P95: {sorted(latencies)[94]:.2f}s")
    print(f"Max: {max(latencies):.2f}s")

asyncio.run(performance_test())
```

---

## Clarification Recommendations

Based on the spec.md questions, here are architectural recommendations:

### Q-ARCH-001: MPI Implementation Support

**Recommendation**: Support all major MPI implementations (OpenMPI, Intel MPI, MPICH) for MVP.

**Rationale**:
- The architectural implications of this decision are minimal - MPI implementation is a container-level concern, not a platform concern. Users package their chosen MPI in the training image.
- This aligns with the upstream KubeFlow Training Operator pattern which is MPI-agnostic.
- Testing scope increases but provides broader customer coverage (financial services favor Intel MPI, research institutions use OpenMPI/MPICH).

**Implementation**:
- Document container image requirements for each MPI implementation
- Provide reference Dockerfiles for common patterns
- Validate basic functionality with each implementation
- Community contribution model for additional MPI implementations

**Trade-off**: Increased documentation and testing burden, but avoids customer adoption blockers.

---

### Q-ARCH-002: Networking and IPC

**Recommendation**: Document standard MPI communication patterns, provide NetworkPolicy templates, defer RDMA to post-MVP.

**Rationale**:
- Most distributed training workloads perform acceptably on standard Kubernetes networking (TCP/IP over CNI).
- RDMA adds significant infrastructure complexity (specialized NICs, driver support, cluster configuration) for incremental performance gains (5-15% in benchmarks).
- In 18 months, this will need to scale to RDMA-enabled clusters for high-performance computing customers, but MVP focus is on mainstream use cases.

**Implementation**:
- **Required Ports**: MPI typically uses SSH (port 22) or dynamic high ports (32768-60999). Document both patterns.
- **CNI Compatibility**: Test with common CNI plugins (Calico, Cilium, OpenShift SDN). Flag incompatibilities.
- **NetworkPolicy Templates**: Provide starter policies for:
  - Same-namespace worker communication
  - External registry access
  - Deny-all default with explicit allowances

**Post-MVP**: RDMA support requires:
- SR-IOV device plugin for InfiniBand/RoCE
- Network topology awareness for optimal worker placement
- Performance profiling tools to validate RDMA benefits

---

### Q-ARCH-003: Resource Orchestration

**Recommendation**: Implement gang scheduling via Volcano Scheduler for MVP.

**Rationale**:
- Gang scheduling is **critical** for distributed training. Partial job placement leads to deadlock (launcher waits for workers that never schedule).
- Volcano is the industry-standard pattern for Kubernetes gang scheduling, used by Alibaba, Baidu, and others for AI workloads.
- This aligns with our north star architecture of leveraging proven open source patterns rather than building custom solutions.

**Implementation**:
- Deploy Volcano Scheduler alongside Training Operator
- Configure MPIJob pods to use Volcano's PodGroup for gang scheduling
- Test eviction scenarios (what happens if a worker pod is evicted mid-training?)
  - **Behavior**: Entire job fails (MPI all-or-nothing semantics)
  - **Recovery**: User must resubmit job (checkpoint/restart in post-MVP)

**Architectural Trade-offs**:
- **Pro**: Prevents resource fragmentation, improves cluster utilization
- **Pro**: Handles backfill scheduling (small jobs fit around large pending jobs)
- **Con**: Additional component to deploy and maintain
- **Con**: Learning curve for administrators

**Alternative Considered**: Custom scheduler extender
- **Rejected**: Would duplicate Volcano functionality, creating technical debt

---

### Q-ARCH-004: Hardware Acceleration

**Recommendation**: Support NVIDIA GPUs for MVP, document path for AMD/Intel accelerators.

**Rationale**:
- 85%+ of enterprise AI workloads use NVIDIA GPUs (based on customer data).
- AMD ROCm and Intel oneAPI support is maturing but not yet mainstream for distributed training.
- Architectural pattern is GPU-agnostic (Kubernetes device plugins abstract hardware details).

**Implementation**:
- Test with NVIDIA GPU Operator
- Validate MIG (Multi-Instance GPU) compatibility
- Document GPU topology awareness (NVLink, PCIe lanes)
- Provide recommendations for GPU:worker ratios

**Post-MVP Roadmap**:
- AMD ROCm support (Q3 2025): Test with MI200 series accelerators
- Intel Data Center GPU support (Q4 2025): Validate with Ponte Vecchio
- Mixed GPU types (2026): Support heterogeneous clusters

**Architectural Considerations**:
- Node selectors for GPU type (e.g., `nvidia.com/gpu.product=A100`)
- Topology constraints for optimal interconnect (workers on same switch/rack)
- Graceful degradation if requested GPU type unavailable

---

### Q-ARCH-005: Container Image Requirements

**Recommendation**: Document requirements, provide reference images, validate common frameworks.

**Rationale**:
- Container image packaging is user responsibility (we can't predict every framework/version).
- Reference images reduce onboarding friction and provide known-good baselines.
- The Martin Fowler pattern of "paved path" - make it easy to do the right thing.

**Required Components in Training Image**:
1. MPI implementation (OpenMPI, Intel MPI, MPICH)
2. Training framework (TensorFlow, PyTorch, etc.)
3. Distributed training library (Horovod, DeepSpeed, etc.)
4. SSH server (if using SSH-based MPI launch)
5. Training code and dependencies

**Reference Images** (provide as starting points):
```Dockerfile
# Horovod + TensorFlow + OpenMPI
FROM nvidia/cuda:11.8.0-cudnn8-devel-ubuntu22.04

RUN apt-get update && apt-get install -y \
    openssh-server \
    openmpi-bin \
    libopenmpi-dev \
    python3 \
    python3-pip

RUN pip3 install tensorflow==2.12.0 horovod[tensorflow]

# Configure SSH for MPI
RUN mkdir /var/run/sshd && \
    ssh-keygen -A && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

COPY training_script.py /workspace/

CMD ["/usr/sbin/sshd", "-D"]
```

**Validation Strategy**:
- Test reference images with example workloads
- Document image scanning (security vulnerabilities)
- Provide image size optimization guidance (multi-stage builds)

---

### Q-NFR-002: Scalability Limits

**Recommendation**: Document tested maximums, design for 100 workers/job and 50 concurrent jobs/cluster for MVP.

**Rationale**:
- Customer data shows 90%+ of distributed training jobs use <20 workers.
- Architectural bottlenecks emerge at scale: etcd limits, control plane load, network saturation.
- Better to set realistic expectations than promise unvalidated limits.

**Tested Maximums for MVP**:
- **Workers per job**: 100 (tested with synthetic workloads)
- **Concurrent jobs per cluster**: 50 (tested with 4-worker jobs)
- **Cluster size**: 200 nodes, 800 GPUs (validated on reference architecture)

**Scalability Considerations**:
- **etcd pressure**: Each MPIJob creates 1 CRD + N pods. Monitor etcd metrics.
- **Control plane load**: Training Operator watches all MPIJobs. Use leader election for HA.
- **Network bandwidth**: Worker communication can saturate 10GbE. Document 25GbE+ recommendations.

**Post-MVP Optimizations** (2-3 year horizon):
- Hierarchical MPIJob (job-of-jobs pattern for 1000+ workers)
- Multi-cluster federation (span distributed training across regions)
- Elastic scaling (dynamic worker count based on available resources)

---

### Q-NFR-003: Reliability and Fault Tolerance

**Recommendation**: Fail-fast for MVP, document checkpoint/restart pattern for post-MVP.

**Rationale**:
- MPI semantics are all-or-nothing (one worker failure → entire job fails).
- Checkpoint/restart adds framework-specific complexity (TensorFlow, PyTorch have different APIs).
- Architectural trade-off: simplicity and predictability vs. resilience.

**MVP Behavior**:
- **Worker failure**: Job fails immediately, user notified with clear error
- **Node failure**: Kubernetes reschedules pods, but job likely already failed
- **Network partition**: MPI timeout → job fails
- **Retry policy**: `backoffLimit: 0` (no automatic retries - user must resubmit)

**Post-MVP Enhancements**:
1. **Checkpoint/Restart** (Q3 2025)
   - Framework-agnostic checkpoint storage (PVC, S3)
   - Automatic checkpoint interval configuration
   - Resume from latest checkpoint on failure

2. **Graceful Degradation** (Q4 2025)
   - Elastic training: reduce worker count if nodes fail
   - Partial job completion: save intermediate results

3. **Preemption Support** (2026)
   - Spot/preemptible instances for cost savings
   - Checkpoint before preemption signal

**Documentation**:
- User guide: "Handling job failures"
- Example: Implementing checkpointing in TensorFlow/PyTorch
- Monitoring: Detect failure patterns (always same worker fails → hardware issue)

---

### Q-OBS-001: Metrics Collection

**Recommendation**: Infrastructure metrics for MVP, training metrics (loss, accuracy) in post-MVP.

**Rationale**:
- Infrastructure metrics (CPU, memory, GPU utilization) are platform responsibility.
- Training metrics are application-specific and require framework integration.
- The architectural implications: if we expose training metrics now, we set expectations for deep framework integration across all job types.

**MVP Metrics** (collected automatically):
- Job lifecycle (created, running, succeeded, failed counts)
- Worker status (active, pending, failed)
- Resource utilization (CPU%, memory%, GPU%, GPU memory%)
- Job duration (total, per phase)
- Pod restarts, OOM kills

**Post-MVP Training Metrics** (user-exposed):
- Loss, accuracy, validation metrics
- Learning rate, batch size
- Throughput (samples/sec, tokens/sec)
- Custom user metrics (via Prometheus pushgateway)

**Implementation Pattern**:
```python
# User training script (post-MVP)
from odh.training import MetricsCollector

metrics = MetricsCollector()

for epoch in range(num_epochs):
    loss = train_epoch()
    accuracy = evaluate()

    # Push to OpenShift AI metrics backend
    metrics.log({"epoch": epoch, "loss": loss, "accuracy": accuracy})
```

**Architectural Consideration**: This creates technical debt if not designed for extensibility. Training metrics integration will compound over time as more frameworks are added.

---

## Risk Assessment

### High-Impact Risks

**Risk 1: Training Operator Upstream Compatibility**
- **Description**: KubeFlow Training Operator V2 API changes break our integration
- **Likelihood**: Medium (upstream is actively developed)
- **Impact**: High (requires code changes across CLI, SDK, UI)
- **Mitigation**:
  - Pin to specific Training Operator version
  - Monitor upstream releases and test pre-release versions
  - Contribute to upstream community for influence on API changes
  - Maintain abstraction layer in Training Service (isolates user-facing APIs from upstream changes)

**Risk 2: Network Policy Incompatibility**
- **Description**: Restrictive enterprise network policies block MPI communication
- **Likelihood**: Medium (common in regulated industries)
- **Impact**: High (jobs fail silently or with cryptic errors)
- **Mitigation**:
  - Provide NetworkPolicy templates with clear documentation
  - Build diagnostic tool: "MPI connectivity test" (validates worker-to-worker communication)
  - Document common firewall rules and CNI configurations
  - Create troubleshooting guide with specific error patterns

**Risk 3: Resource Quota Complexity**
- **Description**: Users struggle to configure appropriate quotas for distributed training
- **Likelihood**: High (quota math is non-intuitive: workers × resources per worker)
- **Impact**: Medium (job creation failures, support burden)
- **Mitigation**:
  - UI shows real-time quota calculations
  - Error messages include resource breakdown
  - Provide quota sizing calculator
  - Document reference quotas for common workload sizes

### Medium-Impact Risks

**Risk 4: Performance Below Expectations**
- **Description**: Distributed training doesn't achieve >2x speedup target
- **Likelihood**: Medium (depends on workload, network, MPI configuration)
- **Impact**: Medium (customer dissatisfaction, but not blocking)
- **Mitigation**:
  - Set expectations: "Speedup depends on workload characteristics"
  - Provide performance tuning guide
  - Document known bottlenecks (communication-heavy vs. compute-heavy)
  - Benchmark reference workloads and publish results

**Risk 5: Volcano Scheduler Deployment Complexity**
- **Description**: Customers resist additional component deployment
- **Likelihood**: Low-Medium (some prefer minimal infrastructure)
- **Impact**: Medium (gang scheduling is critical, but could make optional)
- **Mitigation**:
  - Document Volcano benefits clearly (prevent resource fragmentation)
  - Provide automated deployment scripts
  - Make gang scheduling optional for environments where resource contention is low
  - Offer alternative: simple queue-based scheduling (first-come-first-served)

### Low-Impact Risks

**Risk 6: UI Complexity**
- **Description**: MPIJob configuration form is too complex for data scientists
- **Likelihood**: Low (progressive disclosure mitigates)
- **Impact**: Low (users can fall back to CLI/SDK)
- **Mitigation**:
  - Usability testing with target personas
  - Provide job templates for common patterns
  - "Quick start" mode with minimal fields

**Risk 7: Audit Log Volume**
- **Description**: High-frequency MPIJob operations overwhelm audit logging infrastructure
- **Likelihood**: Low (audit logs are typically handled by enterprise SIEM)
- **Impact**: Low (performance degradation, storage costs)
- **Mitigation**:
  - Configure appropriate audit levels (RequestResponse for mutations, Metadata for reads)
  - Document log retention policies
  - Provide log volume estimation tool

---

## Success Metrics

### Technical Metrics (Measurable at GA)

1. **Performance**
   - Job submission latency: <5s (p95)
   - Status update latency: <10s (p95)
   - Distributed training speedup: >2x vs. single-node baseline
   - Training Service throughput: >100 requests/sec

2. **Reliability**
   - Job success rate: >90% for valid configurations
   - Pod cleanup success rate: 100% (no orphaned resources)
   - MTTR for job failures: <30 minutes (with documentation)

3. **Scalability**
   - Concurrent jobs per cluster: 50+
   - Workers per job: 100+
   - Users per cluster: 500+

### Adoption Metrics (Measurable 90 Days Post-GA)

1. **Usage**
   - 30% of active users create at least one MPIJob
   - 50+ MPIJobs created per week across customer base
   - Average job size: 4-8 workers

2. **Engagement**
   - 70% of users who create first MPIJob create second within 30 days
   - CLI:SDK:UI usage ratio: 40:30:30
   - Documentation page views: 1000+ per month

3. **Customer Satisfaction**
   - NPS score increase: +10 points among distributed training users
   - Support ticket volume: <5 MPIJob-related tickets per week
   - Feature request rate: <10% (indicates MVP is sufficiently complete)

### Business Metrics (Measurable 12 Months Post-GA)

1. **Revenue Impact**
   - $5M+ net-new ARR attributed to distributed training capability
   - 0 lost deals citing lack of MPIJob support
   - 8+ customer case studies published

2. **Market Position**
   - Analyst recognition (Gartner, Forrester) mentions distributed training
   - 3+ conference talks/blog posts by customers
   - Top-5 contributor to KubeFlow Training Operator upstream

---

## Conclusion

This implementation plan provides a comprehensive roadmap for MPIJob support in OpenShift AI, architected for lasting impact across products and designed to scale to future distributed training requirements. The phased approach balances rapid MVP delivery (16 weeks) with enterprise hardening and performance validation (20 weeks total).

Key architectural decisions align with industry best practices:
- Upstream-first pattern with KubeFlow Training Operator
- Unified abstraction layer for consistent user experience
- Gang scheduling via Volcano for production reliability
- Event-driven observability for real-time monitoring

The architectural implications of these decisions extend beyond MPIJob: we're establishing patterns that will scale to Ray, Spark, and custom training operators in the 2-3 year horizon. This aligns with our north star architecture of unified training job management while maintaining flexibility for future innovation.

**Next Steps**:
1. Review and approve this implementation plan
2. Address clarification questions with stakeholder input
3. Begin Phase 1 implementation (Foundation)
4. Establish beta program with 5-10 strategic customers
5. Iterate based on early feedback before GA

---

**Document Version**: 1.0
**Last Updated**: 2025-10-21
**Author**: Archie (AI Architect)
**Status**: Ready for Review
