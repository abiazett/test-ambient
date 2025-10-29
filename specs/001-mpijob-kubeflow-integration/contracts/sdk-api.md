# Python SDK API Contract

**Package**: `openshift-ai-sdk`
**Module**: `openshift_ai_sdk.training`
**Version**: 1.0.0
**Python**: 3.9+

---

## Overview

The OpenShift AI Python SDK provides a Pythonic interface for creating, monitoring, and managing MPIJobs from Jupyter notebooks, scripts, and automation pipelines.

---

## Installation

```bash
pip install openshift-ai-sdk
```

---

## API Classes

### MPIJob

**Description**: Main class for MPIJob lifecycle management.

```python
class MPIJob:
    def __init__(
        self,
        name: str,
        namespace: str = "default",
        workers: int = 2,
        image: str = None,
        command: List[str] = None,
        args: List[str] = None,
        resources: ResourceRequirements = None,
        volume_mounts: Optional[List[Dict]] = None,
        volumes: Optional[List[Dict]] = None,
        env: Optional[Dict[str, str]] = None,
        node_selector: Optional[Dict[str, str]] = None,
        affinity: Optional[Dict] = None,
        tolerations: Optional[List[Dict]] = None,
        slots_per_worker: int = 1,
        clean_pod_policy: str = "Running",
        active_deadline_seconds: Optional[int] = None,
        backoff_limit: int = 6,
        ttl_seconds_after_finished: Optional[int] = None,
    ) -> None:
        """
        Initialize MPIJob configuration.

        Args:
            name: Job identifier (DNS-1123 subdomain, max 63 chars)
            namespace: Kubernetes namespace (default: "default")
            workers: Number of worker pods (2-100)
            image: Container image with MPI runtime
            command: Container entrypoint command
            args: Container command arguments
            resources: CPU, memory, GPU resource requirements
            volume_mounts: List of volume mount configurations
            volumes: List of volume configurations
            env: Environment variables (key-value pairs)
            node_selector: Node selector labels for scheduling
            affinity: Kubernetes affinity rules
            tolerations: Kubernetes tolerations
            slots_per_worker: MPI slots per worker (1-8, default: 1)
            clean_pod_policy: Pod cleanup policy ("Running", "All", "None")
            active_deadline_seconds: Max job duration before killed
            backoff_limit: Retry count on failure (default: 6)
            ttl_seconds_after_finished: Cleanup delay after completion

        Raises:
            ValueError: If name is invalid or workers < 2
        """
```

**Methods**:

```python
def create(self) -> "MPIJob":
    """
    Create MPIJob in cluster.

    Returns:
        Self (for method chaining)

    Raises:
        ApiException: If Kubernetes API call fails
        ValueError: If configuration is invalid
    """

def get_status(self) -> str:
    """
    Get current job status.

    Returns:
        Job status: "Pending", "Running", "Succeeded", "Failed"

    Raises:
        ApiException: If job not found or API call fails
    """

def wait_for_completion(self, timeout: int = 3600, poll_interval: int = 5) -> str:
    """
    Block until job completes or times out.

    Args:
        timeout: Maximum wait time in seconds (default: 3600)
        poll_interval: Seconds between status checks (default: 5)

    Returns:
        Final status: "Succeeded" or "Failed"

    Raises:
        TimeoutError: If job does not complete within timeout
        ApiException: If API call fails
    """

def get_logs(
    self,
    worker_index: Optional[int] = None,
    follow: bool = False,
    tail_lines: int = 1000
) -> str:
    """
    Retrieve logs from launcher or specific worker.

    Args:
        worker_index: Worker pod index (None for launcher logs)
        follow: Stream logs in real-time (default: False)
        tail_lines: Number of lines from end (default: 1000)

    Returns:
        Pod logs as string

    Raises:
        ApiException: If pod not found or not running
    """

def delete(self) -> None:
    """
    Delete MPIJob and associated resources.

    Raises:
        ApiException: If deletion fails
    """

@classmethod
def get(cls, name: str, namespace: str = "default") -> "MPIJob":
    """
    Retrieve existing MPIJob from cluster.

    Args:
        name: Job name
        namespace: Kubernetes namespace (default: "default")

    Returns:
        MPIJob instance reconstructed from cluster state

    Raises:
        ApiException: If job not found
    """

@classmethod
def list(
    cls,
    namespace: str = "default",
    label_selector: Optional[str] = None
) -> List["MPIJob"]:
    """
    List MPIJobs in namespace.

    Args:
        namespace: Kubernetes namespace (default: "default")
        label_selector: Label selector (e.g., "app=training")

    Returns:
        List of MPIJob instances

    Raises:
        ApiException: If API call fails
    """

def get_worker_statuses(self) -> List[Dict[str, Any]]:
    """
    Get status of all worker pods.

    Returns:
        List of worker status dicts with keys:
        - index: Worker index (0-N)
        - name: Pod name
        - phase: Pod phase (Pending/Running/Succeeded/Failed)
        - node: Node name
        - ready: Boolean indicating if pod is ready
        - start_time: Pod start time (datetime)

    Raises:
        ApiException: If API call fails
    """

def get_metrics(self) -> Dict[str, Any]:
    """
    Retrieve performance metrics from Prometheus.

    Returns:
        Dict with keys:
        - workers: List of per-worker metrics
        - aggregated: Cluster-wide aggregates

    Raises:
        ApiException: If Prometheus query fails
    """

def to_yaml(self) -> str:
    """
    Export job configuration as YAML.

    Returns:
        YAML string representation

    Raises:
        ValueError: If configuration is invalid
    """

@classmethod
def from_yaml(cls, yaml_str: str) -> "MPIJob":
    """
    Create MPIJob from YAML configuration.

    Args:
        yaml_str: YAML string

    Returns:
        MPIJob instance

    Raises:
        ValueError: If YAML is invalid
    """
```

---

### ResourceRequirements

**Description**: Container resource requests and limits.

```python
@dataclass
class ResourceRequirements:
    cpu: int = 4
    memory: str = "16Gi"
    gpu: int = 0

    def __post_init__(self):
        """
        Validate resource values.

        Raises:
            ValueError: If values are invalid
        """
```

**Attributes**:

| Field | Type | Default | Description | Validation |
|-------|------|---------|-------------|------------|
| `cpu` | int | 4 | CPU cores requested | > 0 |
| `memory` | str | "16Gi" | Memory quantity | Kubernetes quantity format |
| `gpu` | int | 0 | GPU count | 0-16 |

---

## Usage Examples

### Basic Job Submission

```python
from openshift_ai_sdk.training import MPIJob, ResourceRequirements

# Create and submit job
job = MPIJob(
    name="pytorch-distributed",
    namespace="my-project",
    workers=4,
    image="quay.io/my-org/pytorch-horovod:latest",
    command=["mpirun", "-np", "4", "python", "train.py"],
    resources=ResourceRequirements(cpu=8, memory="64Gi", gpu=1),
)

job.create()
print(f"Job created: {job.name}")
```

### Monitoring Job Status

```python
# Check status
status = job.get_status()
print(f"Job status: {status}")

# Wait for completion
try:
    final_status = job.wait_for_completion(timeout=7200)  # 2 hours
    print(f"Job completed with status: {final_status}")
except TimeoutError:
    print("Job did not complete within 2 hours")
```

### Retrieving Logs

```python
# Get launcher logs
launcher_logs = job.get_logs()
print("Launcher logs:")
print(launcher_logs)

# Get worker 0 logs
worker_logs = job.get_logs(worker_index=0)
print("Worker 0 logs:")
print(worker_logs)

# Stream logs in real-time
for line in job.get_logs(follow=True):
    print(line, end='')
```

### Advanced Configuration with Volumes

```python
job = MPIJob(
    name="llm-fine-tuning",
    namespace="ml-team",
    workers=8,
    image="quay.io/my-org/llama-trainer:latest",
    command=["bash", "/scripts/train.sh"],
    resources=ResourceRequirements(cpu=16, memory="128Gi", gpu=4),
    volume_mounts=[
        {
            "name": "training-data",
            "mountPath": "/data"
        },
        {
            "name": "model-checkpoints",
            "mountPath": "/checkpoints"
        }
    ],
    volumes=[
        {
            "name": "training-data",
            "persistentVolumeClaim": {
                "claimName": "training-data-pvc"
            }
        },
        {
            "name": "model-checkpoints",
            "persistentVolumeClaim": {
                "claimName": "checkpoints-pvc"
            }
        }
    ],
    env={
        "HOROVOD_TIMELINE": "/checkpoints/timeline.json",
        "NCCL_DEBUG": "INFO",
    },
    node_selector={
        "nvidia.com/gpu.product": "NVIDIA-A100-SXM4-80GB"
    },
    slots_per_worker=4,
    active_deadline_seconds=28800,  # 8 hours
)

job.create()
```

### Job Management

```python
# List all jobs in namespace
jobs = MPIJob.list(namespace="ml-team")
for job in jobs:
    print(f"{job.name}: {job.get_status()}")

# Get existing job
existing_job = MPIJob.get(name="pytorch-distributed", namespace="my-project")

# Get worker statuses
worker_statuses = existing_job.get_worker_statuses()
for worker in worker_statuses:
    print(f"Worker {worker['index']}: {worker['phase']} on {worker['node']}")

# Get metrics
metrics = existing_job.get_metrics()
print(f"Average GPU utilization: {metrics['aggregated']['avgGpuUtilization']}%")

# Delete job
existing_job.delete()
```

### Export/Import Configuration

```python
# Export to YAML
yaml_config = job.to_yaml()
with open("job-config.yaml", "w") as f:
    f.write(yaml_config)

# Import from YAML
with open("job-config.yaml", "r") as f:
    yaml_str = f.read()
job = MPIJob.from_yaml(yaml_str)
job.create()
```

---

## Error Handling

### Exception Hierarchy

```python
from kubernetes.client.exceptions import ApiException

try:
    job = MPIJob(name="my-job", namespace="my-project", workers=4, image="my-image")
    job.create()
except ValueError as e:
    # Invalid configuration (e.g., workers < 2, invalid name)
    print(f"Configuration error: {e}")
except ApiException as e:
    if e.status == 403:
        # RBAC permission denied
        print("Insufficient permissions to create MPIJob")
    elif e.status == 409:
        # Job with same name already exists
        print("Job already exists")
    elif e.status == 422:
        # Validation error (e.g., quota exceeded)
        print(f"Validation error: {e.body}")
    else:
        # Other API errors
        print(f"API error: {e}")
```

### Common Error Patterns

```python
# Quota exceeded
try:
    job.create()
except ApiException as e:
    if "quota exceeded" in str(e).lower():
        print("Namespace GPU quota exceeded")
        print("Try reducing worker count or deleting existing jobs")

# Image pull failure
try:
    job.wait_for_completion(timeout=300)
except TimeoutError:
    worker_statuses = job.get_worker_statuses()
    for worker in worker_statuses:
        if worker['phase'] == 'Pending':
            # Check for image pull errors
            print(f"Worker {worker['index']} stuck in Pending")

# MPI communication failure
logs = job.get_logs()
if "ssh: connect to host" in logs:
    print("MPI communication failure detected")
    print("Check NetworkPolicy and SSH configuration")
```

---

## Type Hints and IDE Support

The SDK provides complete type hints for IDE autocompletion:

```python
from openshift_ai_sdk.training import MPIJob
from typing import List, Dict, Optional

def submit_training_job(
    name: str,
    workers: int,
    image: str,
    env: Optional[Dict[str, str]] = None
) -> MPIJob:
    job = MPIJob(
        name=name,
        workers=workers,
        image=image,
        env=env or {},
    )
    job.create()
    return job

# IDE will provide autocompletion for all methods
job = submit_training_job("my-job", 4, "my-image")
status: str = job.get_status()  # Type hint: str
```

---

## Async Support (Future Enhancement)

Planned for post-MVP:

```python
import asyncio
from openshift_ai_sdk.training import AsyncMPIJob

async def submit_and_monitor():
    job = AsyncMPIJob(name="my-job", workers=4, image="my-image")
    await job.create()

    status = await job.get_status()
    print(f"Status: {status}")

    final_status = await job.wait_for_completion()
    print(f"Completed: {final_status}")

asyncio.run(submit_and_monitor())
```

---

## SDK Testing

### Unit Tests

```python
import pytest
from unittest.mock import Mock, patch
from openshift_ai_sdk.training import MPIJob

def test_mpijob_create():
    job = MPIJob(name="test-job", workers=2, image="test-image")
    manifest = job._build_manifest()

    assert manifest["metadata"]["name"] == "test-job"
    assert manifest["spec"]["mpiReplicaSpecs"]["Worker"]["replicas"] == 2

def test_mpijob_invalid_workers():
    with pytest.raises(ValueError):
        MPIJob(name="test-job", workers=1, image="test-image")  # < 2 workers

@patch("openshift_ai_sdk.training.mpijob.client.CustomObjectsApi")
def test_mpijob_get_status(mock_api):
    mock_api.return_value.get_namespaced_custom_object.return_value = {
        "status": {
            "conditions": [{"type": "Running", "status": "True"}]
        }
    }

    job = MPIJob(name="test-job", workers=2, image="test-image")
    status = job.get_status()

    assert status == "Running"
```

### Integration Tests

```python
@pytest.mark.integration
def test_mpijob_end_to_end(k8s_cluster):
    job = MPIJob(
        name="integration-test",
        namespace="test-namespace",
        workers=2,
        image="horovod/horovod:latest",
        command=["horovodrun", "-np", "2", "python", "-c", "print('Hello')"],
    )

    job.create()

    try:
        final_status = job.wait_for_completion(timeout=300)
        assert final_status == "Succeeded"

        logs = job.get_logs()
        assert "Hello" in logs
    finally:
        job.delete()
```

---

## Compatibility

| SDK Version | RHOAI Version | Kubernetes API | Python |
|-------------|---------------|----------------|--------|
| 1.0.0       | 2.17+         | v2beta1        | 3.9+   |

---

## Migration from kubectl

```bash
# kubectl approach
kubectl apply -f mpijob.yaml
kubectl get mpijobs
kubectl logs my-job-launcher
kubectl delete mpijob my-job
```

```python
# SDK approach
from openshift_ai_sdk.training import MPIJob

job = MPIJob.from_yaml(open("mpijob.yaml").read())
job.create()

jobs = MPIJob.list()
logs = job.get_logs()
job.delete()
```

**Benefits**:
- Pythonic API (vs. YAML manipulation)
- Type safety and IDE autocompletion
- Error handling with exceptions
- Integrated with Jupyter notebooks
- Synchronous and asynchronous patterns
