# Getting Started with MPIJobs in OpenShift AI

This guide shows you how to create and run your first MPI-based distributed training job using OpenShift AI's Trainer V2 API.

## Prerequisites

Before you begin, ensure you have:

1. Access to an OpenShift AI cluster (OpenShift 4.14+ with Kubernetes 1.27+)
2. The `oc` or `kubectl` CLI tool installed
3. Python 3.11+ (for SDK option)
4. Appropriate RBAC permissions to create TrainJobs in your namespace

## What is an MPIJob?

An MPIJob in OpenShift AI is a distributed training workload that uses the Message Passing Interface (MPI) protocol to coordinate multiple training nodes. Common use cases include:

- **Horovod** - distributed deep learning with TensorFlow or PyTorch
- **Intel MPI** - optimized CPU-based distributed training
- **OpenMPI** - general-purpose distributed computing

## Three Ways to Create an MPIJob

You can create MPIJobs using any of these interfaces:

1. **Python SDK** (recommended) - programmatic job creation and monitoring
2. **CLI (kubectl/oc)** - YAML-based job definitions
3. **Dashboard UI** - web-based visual interface

## Option 1: Python SDK (Recommended)

### Step 1: Install the SDK

```bash
pip install kubeflow-training
```

### Step 2: Create Your First MPIJob

```python
from kubeflow.training import TrainingClient

# Initialize the client
client = TrainingClient()

# Create an MPIJob for distributed Horovod training
client.create_mpi_job(
    name="mnist-distributed",
    namespace="default",
    runtime="mpi-horovod-gpu",
    num_nodes=4,
    image="horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-mxnet1.9.1-py3.9-gpu",
    command=["horovodrun"],
    args=[
        "-np", "4",
        "--verbose",
        "python",
        "/examples/pytorch/pytorch_mnist.py",
        "--epochs=10",
        "--batch-size=64"
    ],
    resources_per_node={
        "requests": {
            "cpu": "4",
            "memory": "16Gi",
            "nvidia.com/gpu": "1"
        },
        "limits": {
            "cpu": "8",
            "memory": "32Gi",
            "nvidia.com/gpu": "1"
        }
    },
    env=[
        {"name": "PYTHONUNBUFFERED", "value": "1"},
        {"name": "NCCL_DEBUG", "value": "INFO"}
    ]
)

print("MPIJob created successfully!")
```

### Step 3: Monitor Your Job

```python
# Wait for job to complete (with 30-minute timeout)
client.wait_for_job_completion(
    name="mnist-distributed",
    namespace="default",
    timeout=1800
)

# Get job status
job = client.get_job(name="mnist-distributed", namespace="default")
print(f"Job Status: {job.status.phase}")
print(f"Start Time: {job.status.start_time}")
print(f"Completion Time: {job.status.completion_time}")

# Get logs from launcher pod
logs = client.get_job_logs(
    name="mnist-distributed",
    namespace="default",
    pod_type="launcher"
)
print(logs)
```

### Step 4: Clean Up

```python
# Delete the job when finished
client.delete_job(name="mnist-distributed", namespace="default")
```

## Option 2: CLI (kubectl/oc)

### Step 1: Create a TrainJob YAML

Create a file named `mnist-distributed.yaml`:

```yaml
apiVersion: kubeflow.org/v2alpha1
kind: TrainJob
metadata:
  name: mnist-distributed
  namespace: default
spec:
  # Reference to the cluster runtime template
  runtimeRef:
    kind: ClusterTrainingRuntime
    name: mpi-horovod-gpu

  # Training configuration
  trainer:
    # Number of worker nodes
    numNodes: 4

    # Container image with training code
    image: horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-mxnet1.9.1-py3.9-gpu

    # Training command and arguments
    command:
      - horovodrun
    args:
      - -np
      - "4"
      - --verbose
      - python
      - /examples/pytorch/pytorch_mnist.py
      - --epochs=10
      - --batch-size=64

    # Resources per worker node
    resourcesPerNode:
      requests:
        cpu: "4"
        memory: "16Gi"
        nvidia.com/gpu: "1"
      limits:
        cpu: "8"
        memory: "32Gi"
        nvidia.com/gpu: "1"

    # Environment variables
    env:
      - name: PYTHONUNBUFFERED
        value: "1"
      - name: NCCL_DEBUG
        value: INFO
```

### Step 2: Submit the Job

```bash
oc apply -f mnist-distributed.yaml
```

### Step 3: Monitor the Job

```bash
# Check job status
oc get trainjob mnist-distributed

# Get detailed information
oc describe trainjob mnist-distributed

# Watch job progress
oc get trainjob mnist-distributed -w

# Check pod status
oc get pods -l training.kubeflow.org/job-name=mnist-distributed
```

### Step 4: View Logs

```bash
# Get launcher logs
oc logs -l training.kubeflow.org/job-name=mnist-distributed,training.kubeflow.org/job-role=launcher

# Get worker logs (worker 0)
oc logs mnist-distributed-worker-0

# Stream logs in real-time
oc logs -f -l training.kubeflow.org/job-name=mnist-distributed,training.kubeflow.org/job-role=launcher
```

### Step 5: Clean Up

```bash
oc delete trainjob mnist-distributed
```

## Option 3: Dashboard UI

### Step 1: Access the Dashboard

1. Navigate to OpenShift AI Dashboard in your browser
2. Log in with your credentials
3. Click on "Training Jobs" in the left navigation menu

### Step 2: Create a New Job

1. Click the "Create Training Job" button
2. Select "MPI" as the job type
3. Fill in the form:
   - **Job Name**: `mnist-distributed`
   - **Runtime Template**: `mpi-horovod-gpu`
   - **Number of Nodes**: `4`
   - **Container Image**: `horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-mxnet1.9.1-py3.9-gpu`
   - **Command**: `horovodrun -np 4 --verbose python /examples/pytorch/pytorch_mnist.py --epochs=10`
   - **Resources per Node**:
     - CPU: 4 cores request, 8 cores limit
     - Memory: 16Gi request, 32Gi limit
     - GPU: 1 NVIDIA GPU
4. Click "Create"

### Step 3: Monitor Your Job

The Dashboard will show:
- **Topology View**: Visual representation of launcher and worker pods
- **Status**: Current phase (Pending, Running, Succeeded, Failed)
- **Logs**: Tabbed interface for launcher and worker logs
- **Events**: Timeline of lifecycle events
- **Metrics**: CPU, memory, and GPU utilization charts

### Step 4: View Logs

1. Click on the job name to open the details page
2. Navigate to the "Logs" tab
3. Select "Launcher" or "Worker" to view specific pod logs
4. Use the search and filter tools to find specific log entries

### Step 5: Clean Up

Click the "Delete" button on the job details page, or select the job from the list and click "Delete".

## Understanding What Happened

When you created the MPIJob, the Training Operator performed these actions:

1. **Gang Scheduling**: Created a PodGroup to ensure all pods (1 launcher + 4 workers) schedule atomically
2. **SSH Keys**: Generated ephemeral SSH key pairs for secure communication
3. **Hostfile**: Created an MPI hostfile with worker pod addresses
4. **JobSet**: Created a JobSet with launcher and worker Jobs
5. **Pods**: Scheduled 1 launcher pod and 4 worker pods
6. **Networking**: Applied NetworkPolicy to allow SSH communication on port 2222
7. **Training**: Launcher executed `mpirun` to coordinate distributed training across workers

## Common Issues and Solutions

### Issue: Gang Scheduling Timeout

**Symptoms**: Job stays in "Pending" phase for >5 minutes

**Cause**: Not enough nodes available to schedule all pods simultaneously

**Solution**:
- Reduce `numNodes` to fit available cluster capacity
- Increase cluster node count
- Adjust resource requests to fit available nodes

### Issue: ImagePullBackOff

**Symptoms**: Pods stuck in "ImagePullBackOff" status

**Cause**: Container image not accessible

**Solution**:
- Verify image name and tag are correct
- Check image registry credentials
- Ensure cluster has network access to registry

### Issue: SSH Connection Failed

**Symptoms**: Launcher logs show "ssh: connect to host ... port 2222: Connection refused"

**Cause**: NetworkPolicy blocking SSH communication

**Solution**:
- Apply the MPI NetworkPolicy: `oc apply -f mpi-training-network-policy.yaml`
- Verify NetworkPolicy allows port 2222 between pods

### Issue: Resource Quota Exceeded

**Symptoms**: Job creation fails with "exceeded quota" error

**Cause**: Namespace resource limits exceeded

**Solution**:
- Reduce `numNodes` or `resourcesPerNode`
- Request quota increase from cluster administrator
- Delete unused jobs to free resources

### Issue: OOM (Out of Memory)

**Symptoms**: Worker pods terminated with "OOMKilled" status

**Cause**: Training process exceeded memory limits

**Solution**:
- Increase `resourcesPerNode.limits.memory`
- Reduce batch size in training script
- Optimize model architecture

## Next Steps

Now that you've created your first MPIJob, you can:

1. **Customize Training**: Modify the command and args to run your own training scripts
2. **Add Data Volumes**: Mount PVCs or object storage for training data
3. **Scale Up**: Increase `numNodes` for larger distributed training
4. **Use Different MPI Implementations**: Try `mpi-openmpi-cpu` or `mpi-intel-mpi`
5. **Integrate with CI/CD**: Automate job creation using the SDK in your pipelines

## Additional Resources

- [User Guide](../user-guide/README.md) - Comprehensive user documentation
- [API Reference](../api-reference/README.md) - Detailed API documentation
- [Troubleshooting Guide](../troubleshooting/README.md) - Common issues and solutions
- [Migration Guide](../migration/README.md) - Migrating from legacy MPIJob v2beta1
- [Quickstart Guide](../../specs/001-mpijob-trainer-v2-support/quickstart.md) - 5-minute quickstart

## Getting Help

If you encounter issues not covered in this guide:

1. Check the [Troubleshooting Guide](../troubleshooting/README.md)
2. Review the [FAQ](./faq.md)
3. Search existing issues in the KubeFlow Training Operator repository
4. Ask in the OpenShift AI community forums
5. Contact Red Hat support (if using a supported OpenShift AI distribution)
