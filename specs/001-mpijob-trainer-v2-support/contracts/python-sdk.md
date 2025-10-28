# Python SDK Contract: MPIJob Support

**Feature Branch**: `001-mpijob-trainer-v2-support`
**Date**: 2025-10-28
**SDK Version**: >= 2.0.0

This document defines the Python SDK API contract for creating and managing MPIJobs in OpenShift AI.

---

## Overview

The Python SDK extends the KubeFlow Training SDK with MPIJob convenience methods while maintaining compatibility with the generic TrainJob API.

**SDK Repository**: `kubeflow/training-sdk` or `red-hat-data-services/openshift-ai-sdk`

---

## Installation

```bash
pip install kubeflow-training-sdk>=2.0.0
# or
pip install openshift-ai-sdk>=1.5.0
```

---

## API Methods

### 1. Create MPIJob

Creates a new MPIJob using Trainer V2 API.

**Method Signature**:
```python
def create_mpi_job(
    name: str,
    runtime: str,
    num_nodes: int,
    image: str,
    command: Optional[List[str]] = None,
    args: Optional[List[str]] = None,
    resources_per_node: Optional[Dict[str, Dict[str, str]]] = None,
    env: Optional[List[Dict[str, str]]] = None,
    volumes: Optional[List[Dict]] = None,
    volume_mounts: Optional[List[Dict]] = None,
    namespace: str = "default",
    labels: Optional[Dict[str, str]] = None,
    annotations: Optional[Dict[str, str]] = None,
) -> TrainJob
```

**Parameters**:
- `name` (str, required): Unique job name within namespace
- `runtime` (str, required): Name of TrainingRuntime or ClusterTrainingRuntime with MPI support
- `num_nodes` (int, required): Number of worker nodes (1-64)
- `image` (str, required): Container image for training (e.g., "horovod/horovod:latest")
- `command` (List[str], optional): Entrypoint command override
- `args` (List[str], optional): Command arguments
- `resources_per_node` (dict, optional): Resource requests/limits per node
  - Format: `{"requests": {"cpu": "8", "memory": "32Gi", "nvidia.com/gpu": "2"}, "limits": {...}}`
- `env` (List[dict], optional): Environment variables
  - Format: `[{"name": "ENV_VAR", "value": "value"}]`
- `volumes` (List[dict], optional): Kubernetes volume definitions
- `volume_mounts` (List[dict], optional): Volume mount definitions
- `namespace` (str, optional): Target namespace (default: "default")
- `labels` (dict, optional): Job labels
- `annotations` (dict, optional): Job annotations

**Returns**: `TrainJob` object with initial status

**Raises**:
- `ValueError`: Invalid parameters (e.g., num_nodes < 1 or > 64)
- `RuntimeError`: Runtime not found or doesn't support MPI
- `K8sApiException`: Kubernetes API error (e.g., insufficient quota)

**Example**:
```python
from kubeflow.training import TrainingClient

client = TrainingClient()

job = client.create_mpi_job(
    name="mnist-horovod-training",
    runtime="mpi-horovod-gpu",
    num_nodes=4,
    image="horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-py3.10-gpu",
    command=["python", "/workspace/train.py"],
    args=["--epochs", "10", "--batch-size", "128"],
    resources_per_node={
        "requests": {
            "cpu": "8",
            "memory": "32Gi",
            "nvidia.com/gpu": "2"
        },
        "limits": {
            "nvidia.com/gpu": "2"
        }
    },
    env=[
        {"name": "NCCL_DEBUG", "value": "INFO"},
        {"name": "HOROVOD_LOG_LEVEL", "value": "DEBUG"}
    ],
    namespace="ml-project"
)

print(f"Created job: {job.metadata.name}")
print(f"Status: {job.status.phase}")
```

---

### 2. Get Job Status

Retrieves detailed status for a training job.

**Method Signature**:
```python
def get_job(
    name: str,
    namespace: str = "default"
) -> TrainJob
```

**Parameters**:
- `name` (str, required): Job name
- `namespace` (str, optional): Namespace (default: "default")

**Returns**: `TrainJob` object with current status

**Raises**:
- `K8sApiException`: Job not found (404)

**Example**:
```python
job = client.get_job("mnist-horovod-training", namespace="ml-project")

print(f"Phase: {job.status.phase}")
print(f"Launcher Status: {job.status.launcher_pod_status.phase}")
print(f"Workers Ready: {sum(1 for w in job.status.worker_pod_statuses if w.phase == 'Running')}/{len(job.status.worker_pod_statuses)}")

if job.status.failure_reason:
    print(f"Failure Reason: {job.status.failure_reason}")
    print(f"Failure Message: {job.status.failure_message}")
```

---

### 3. Wait for Job Completion

Blocks until job completes (Succeeded or Failed) with optional timeout.

**Method Signature**:
```python
def wait_for_job_completion(
    name: str,
    namespace: str = "default",
    timeout: Optional[int] = None,
    polling_interval: int = 10
) -> TrainJob
```

**Parameters**:
- `name` (str, required): Job name
- `namespace` (str, optional): Namespace (default: "default")
- `timeout` (int, optional): Maximum wait time in seconds (None = wait indefinitely)
- `polling_interval` (int, optional): Seconds between status checks (default: 10)

**Returns**: `TrainJob` object in terminal state (Succeeded or Failed)

**Raises**:
- `TimeoutError`: Job did not complete within timeout
- `K8sApiException`: Job not found

**Example**:
```python
job = client.wait_for_job_completion(
    "mnist-horovod-training",
    namespace="ml-project",
    timeout=3600  # 1 hour
)

if job.status.phase == "Succeeded":
    print("Training succeeded!")
    print(f"Duration: {job.status.completion_time - job.status.start_time}")
else:
    print(f"Training failed: {job.status.failure_reason}")
```

---

### 4. Stream Job Events

Streams lifecycle events for real-time monitoring.

**Method Signature**:
```python
def stream_job_events(
    name: str,
    namespace: str = "default"
) -> Iterator[Event]
```

**Parameters**:
- `name` (str, required): Job name
- `namespace` (str, optional): Namespace (default: "default")

**Yields**: `Event` objects as they occur

**Raises**:
- `K8sApiException`: Job not found

**Example**:
```python
for event in client.stream_job_events("mnist-horovod-training", namespace="ml-project"):
    print(f"[{event.timestamp}] {event.type}: {event.reason} - {event.message}")

    if event.reason == "TrainingComplete":
        print("Training finished!")
        break
```

---

### 5. Get Job Logs

Retrieves logs from launcher and/or worker pods.

**Method Signature**:
```python
def get_job_logs(
    name: str,
    namespace: str = "default",
    pod_type: str = "all",
    pod_index: Optional[int] = None,
    tail_lines: Optional[int] = 100,
    since_seconds: Optional[int] = None,
    follow: bool = False
) -> Union[str, Iterator[str]]
```

**Parameters**:
- `name` (str, required): Job name
- `namespace` (str, optional): Namespace (default: "default")
- `pod_type` (str, optional): "launcher", "worker", or "all" (default: "all")
- `pod_index` (int, optional): Worker pod index (0, 1, 2, ...) - only for pod_type="worker"
- `tail_lines` (int, optional): Number of lines from end (default: 100)
- `since_seconds` (int, optional): Logs from last N seconds
- `follow` (bool, optional): Stream logs in real-time (default: False)

**Returns**:
- If `follow=False`: String with all logs
- If `follow=True`: Iterator yielding log lines as they arrive

**Raises**:
- `K8sApiException`: Job or pod not found
- `ValueError`: Invalid pod_type or pod_index

**Example**:
```python
# Get launcher logs
launcher_logs = client.get_job_logs(
    "mnist-horovod-training",
    namespace="ml-project",
    pod_type="launcher",
    tail_lines=200
)
print(launcher_logs)

# Get specific worker logs
worker_logs = client.get_job_logs(
    "mnist-horovod-training",
    namespace="ml-project",
    pod_type="worker",
    pod_index=0
)
print(worker_logs)

# Stream logs in real-time
for log_line in client.get_job_logs(
    "mnist-horovod-training",
    namespace="ml-project",
    pod_type="all",
    follow=True
):
    print(log_line, end='')
```

---

### 6. Delete Job

Deletes a training job and associated resources.

**Method Signature**:
```python
def delete_job(
    name: str,
    namespace: str = "default",
    grace_period_seconds: int = 30
) -> None
```

**Parameters**:
- `name` (str, required): Job name
- `namespace` (str, optional): Namespace (default: "default")
- `grace_period_seconds` (int, optional): Grace period for pod termination (default: 30)

**Returns**: None

**Raises**:
- `K8sApiException`: Job not found

**Example**:
```python
client.delete_job("mnist-horovod-training", namespace="ml-project")
print("Job deleted successfully")
```

---

### 7. List Jobs

Lists all training jobs in a namespace with optional filtering.

**Method Signature**:
```python
def list_jobs(
    namespace: str = "default",
    runtime: Optional[str] = None,
    status: Optional[str] = None,
    label_selector: Optional[str] = None
) -> List[TrainJob]
```

**Parameters**:
- `namespace` (str, optional): Namespace (default: "default")
- `runtime` (str, optional): Filter by runtime name
- `status` (str, optional): Filter by status (Pending, Running, Succeeded, Failed)
- `label_selector` (str, optional): Kubernetes label selector (e.g., "team=ml,project=mnist")

**Returns**: List of `TrainJob` objects

**Example**:
```python
# List all jobs
all_jobs = client.list_jobs(namespace="ml-project")

# List running MPI jobs
mpi_jobs = client.list_jobs(
    namespace="ml-project",
    runtime="mpi-horovod-gpu",
    status="Running"
)

for job in mpi_jobs:
    print(f"{job.metadata.name}: {job.status.phase}")
```

---

### 8. Get Job Metrics

Retrieves resource utilization metrics.

**Method Signature**:
```python
def get_job_metrics(
    name: str,
    namespace: str = "default",
    time_range: str = "1h"
) -> Dict[str, Any]
```

**Parameters**:
- `name` (str, required): Job name
- `namespace` (str, optional): Namespace (default: "default")
- `time_range` (str, optional): Time range (1h, 6h, 24h, 7d) (default: "1h")

**Returns**: Dictionary with launcher and worker metrics (CPU, memory, GPU)

**Example**:
```python
metrics = client.get_job_metrics("mnist-horovod-training", namespace="ml-project")

print(f"Launcher CPU: {metrics['launcher']['cpu'][-1]['value']} cores")
print(f"Launcher Memory: {metrics['launcher']['memory'][-1]['value'] / (1024**3):.2f} GiB")

for i, worker_metrics in enumerate(metrics['workers']):
    gpu_util = worker_metrics['gpu'][-1]['value']
    print(f"Worker {i} GPU Utilization: {gpu_util}%")
```

---

### 9. List Available Runtimes

Lists available training runtimes (cluster-wide and namespace-scoped).

**Method Signature**:
```python
def list_runtimes(
    namespace: str = "default",
    framework: Optional[str] = None
) -> List[Dict[str, Any]]
```

**Parameters**:
- `namespace` (str, optional): Namespace (default: "default")
- `framework` (str, optional): Filter by framework (e.g., "mpi", "pytorch")

**Returns**: List of runtime dictionaries

**Example**:
```python
runtimes = client.list_runtimes(namespace="ml-project", framework="mpi")

for runtime in runtimes:
    print(f"{runtime['name']} ({runtime['kind']}): {runtime['mpiImplementation']}")
```

---

## Data Models

### TrainJob

```python
@dataclass
class TrainJob:
    metadata: ObjectMeta
    spec: TrainJobSpec
    status: TrainJobStatus
```

### ObjectMeta

```python
@dataclass
class ObjectMeta:
    name: str
    namespace: str
    labels: Optional[Dict[str, str]] = None
    annotations: Optional[Dict[str, str]] = None
    creation_timestamp: Optional[datetime] = None
```

### TrainJobSpec

```python
@dataclass
class TrainJobSpec:
    runtime_ref: RuntimeRef
    trainer: TrainerSpec
    initializer: Optional[InitializerSpec] = None
```

### RuntimeRef

```python
@dataclass
class RuntimeRef:
    kind: str  # "TrainingRuntime" or "ClusterTrainingRuntime"
    name: str
    api_group: str = "kubeflow.org"
```

### TrainerSpec

```python
@dataclass
class TrainerSpec:
    num_nodes: int
    image: str
    command: Optional[List[str]] = None
    args: Optional[List[str]] = None
    resources_per_node: Optional[ResourceRequirements] = None
    env: Optional[List[EnvVar]] = None
```

### TrainJobStatus

```python
@dataclass
class TrainJobStatus:
    phase: str  # "Pending", "Running", "Succeeded", "Failed"
    conditions: List[Condition]
    start_time: Optional[datetime] = None
    completion_time: Optional[datetime] = None
    launcher_pod_status: Optional[PodStatus] = None
    worker_pod_statuses: Optional[List[PodStatus]] = None
    failure_reason: Optional[str] = None
    failure_message: Optional[str] = None
```

### PodStatus

```python
@dataclass
class PodStatus:
    index: Optional[int]  # None for launcher, 0+ for workers
    phase: str  # "Pending", "Running", "Succeeded", "Failed"
    reason: Optional[str] = None
    message: Optional[str] = None
```

### Event

```python
@dataclass
class Event:
    type: str  # "Normal", "Warning"
    reason: str
    message: str
    timestamp: datetime
    source: str
```

---

## Error Handling

### Exception Hierarchy

```python
TrainingSDKException (base)
├── ValueError (parameter validation)
├── RuntimeError (runtime not found, configuration errors)
├── K8sApiException (Kubernetes API errors)
│   ├── NotFoundError (404)
│   ├── ForbiddenError (403)
│   ├── UnauthorizedError (401)
│   └── QuotaExceededError (quota validation)
└── TimeoutError (wait timeout)
```

### Example Error Handling

```python
from kubeflow.training import TrainingClient
from kubeflow.training.exceptions import (
    K8sApiException,
    NotFoundError,
    QuotaExceededError
)

try:
    job = client.create_mpi_job(
        name="mnist-horovod",
        runtime="mpi-horovod-gpu",
        num_nodes=4,
        image="horovod/horovod:latest",
        resources_per_node={"requests": {"cpu": "8", "memory": "32Gi", "nvidia.com/gpu": "2"}},
        namespace="ml-project"
    )
except QuotaExceededError as e:
    print(f"Insufficient quota: {e.message}")
    # Retry with fewer nodes
    job = client.create_mpi_job(..., num_nodes=2)
except K8sApiException as e:
    print(f"Kubernetes error: {e}")
    raise
```

---

## Configuration

### Client Initialization

```python
from kubeflow.training import TrainingClient

# Default: uses kubeconfig from ~/.kube/config
client = TrainingClient()

# Custom kubeconfig path
client = TrainingClient(config_file="/path/to/kubeconfig")

# In-cluster config (when running inside Kubernetes)
client = TrainingClient(in_cluster=True)
```

---

## Backward Compatibility

**Compatibility Notes**:
- SDK version >= 2.0.0 required for Trainer V2 support
- Legacy MPIJob v2beta1 API available via `create_legacy_mpi_job()` for migration period
- All Trainer V2 methods work consistently across MPI, PyTorch, and TensorFlow jobs

**Deprecation Timeline**:
- Legacy API methods deprecated in SDK 2.0.0 (warnings emitted)
- Legacy API methods removed in SDK 3.0.0 (12 months later)

---

## Testing Support

### Mock Client for Unit Tests

```python
from kubeflow.training.testing import MockTrainingClient

# Create mock client for testing
mock_client = MockTrainingClient()

# Register expected job
mock_client.register_job(
    name="test-job",
    namespace="test",
    status_phase="Running"
)

# Use in tests
job = mock_client.get_job("test-job", namespace="test")
assert job.status.phase == "Running"
```

---

## Summary

The Python SDK provides 9 core methods for MPIJob management:
1. `create_mpi_job()` - Create new job
2. `get_job()` - Get job status
3. `wait_for_job_completion()` - Block until completion
4. `stream_job_events()` - Real-time event streaming
5. `get_job_logs()` - Retrieve logs
6. `delete_job()` - Delete job
7. `list_jobs()` - List jobs with filtering
8. `get_job_metrics()` - Resource utilization
9. `list_runtimes()` - Available runtimes

All methods include comprehensive error handling, type hints, and docstrings. SDK is backward compatible with legacy MPIJob v2beta1 API during migration period.

**Next Steps**: Implement SDK methods, write contract tests, and generate quickstart documentation.
