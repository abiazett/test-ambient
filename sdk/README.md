# OpenShift AI Python SDK

Python SDK for interacting with OpenShift AI services, with support for MPIJob management.

## Installation

### From PyPI

```bash
pip install openshift-ai-sdk
```

### From Source

```bash
# Clone the repository
git clone https://github.com/openshift-ai/mpijob
cd mpijob/sdk

# Install in development mode
pip install -e .
```

## Usage

### Basic Usage

```python
from odh.training import MPIJobClient, MPIJobSpec, WorkerSpec, ResourceSpec

# Initialize client
client = MPIJobClient(namespace="my-namespace")

# Create an MPIJob
worker_spec = {
    "replicas": 4,
    "resources": {
        "requests": {
            "cpu": "4",
            "memory": "16Gi",
            "nvidia.com/gpu": 2
        },
        "limits": {
            "cpu": "4",
            "memory": "16Gi",
            "nvidia.com/gpu": 2
        }
    },
    "image": "example/tf-horovod:latest",
    "command": [
        "python",
        "/train.py",
        "--epochs",
        "10"
    ]
}

# Create job
job = client.create_mpijob(
    name="my-mpijob",
    namespace="my-namespace",
    worker_spec=worker_spec
)

# Wait for job completion
job.wait_for_completion(timeout=3600)

# Get logs from workers
logs = job.get_logs(worker=0)
```

### Creating an MPIJob from YAML

```python
from odh.training import MPIJobClient

# Initialize client
client = MPIJobClient()

# Create from YAML file
job = client.create_mpijob(
    name="my-job",  # Will be overridden by file
    from_file="path/to/mpijob.yaml"
)
```

### Monitoring an MPIJob

```python
from odh.training import MPIJobClient, MPIJobStatus

# Initialize client
client = MPIJobClient()

# Get existing job
job = client.get_mpijob(name="my-job", namespace="my-namespace")

# Define status callback
def on_status_change(status: MPIJobStatus):
    print(f"Job status: {status.conditions[-1].type if status.conditions else 'Unknown'}")
    print(f"Workers running: {job.workers_running}/{job.raw['spec']['mpiReplicaSpecs']['Worker']['replicas']}")

# Monitor job with callback
job.monitor(callback=on_status_change)
```

### Listing MPIJobs

```python
from odh.training import MPIJobClient

# Initialize client
client = MPIJobClient()

# List all jobs in default namespace
jobs = client.list_mpijobs()

# List jobs with label selector
jobs = client.list_mpijobs(
    namespace="my-namespace",
    label_selector="app=my-app"
)

# Print job information
for job in jobs:
    print(f"Job: {job.name}, Status: {job.phase}, Workers: {job.workers_running}")
```

### Deleting an MPIJob

```python
from odh.training import MPIJobClient

# Initialize client
client = MPIJobClient()

# Delete job
client.delete_mpijob(name="my-job", namespace="my-namespace", wait=True)
```

## Examples

The SDK includes several examples in the `examples/` directory:

- `create_mpijob.py`: Create and monitor an MPIJob
- `list_mpijobs.py`: List and filter MPIJobs
- `monitor_mpijob.py`: Monitor an existing MPIJob

To run an example:

```bash
cd examples
python create_mpijob.py --name=test-job --namespace=default --workers=2 --gpus=2
```

## Development

### Prerequisites

- Python 3.8 or higher
- Access to a Kubernetes cluster with KubeFlow Training Operator V2

### Running Tests

```bash
# Install test dependencies
pip install -e ".[test]"

# Run tests
pytest
```

### Code Style

The project uses Black, isort, and mypy for code formatting and type checking:

```bash
# Install development dependencies
pip install -e ".[dev]"

# Format code
black odh
isort odh

# Type checking
mypy odh
```