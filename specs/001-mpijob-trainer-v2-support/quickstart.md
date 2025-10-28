# Quickstart: MPIJobs in OpenShift AI

**Goal**: Run your first MPIJob in OpenShift AI using Trainer V2 in 15 minutes.

**Prerequisites**:
- OpenShift AI cluster with Trainer V2 installed
- Access to a namespace with GPU resources (or CPU for testing)
- `oc` CLI or Python SDK installed
- Basic familiarity with distributed training concepts

---

## Option 1: Python SDK (Recommended)

### Step 1: Install SDK

```bash
pip install kubeflow-training-sdk>=2.0.0
```

### Step 2: Create Your First MPIJob

Save this as `mnist_mpi_job.py`:

```python
from kubeflow.training import TrainingClient

# Initialize client
client = TrainingClient()

# Create MPIJob for Horovod MNIST training
job = client.create_mpi_job(
    name="mnist-horovod-quickstart",
    runtime="mpi-horovod-gpu",  # Use pre-configured runtime
    num_nodes=2,  # 2 workers for simple distributed training
    image="horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-py3.10-gpu",
    command=[
        "horovodrun",
        "-np", "2",
        "--verbose",
        "python", "/examples/pytorch/pytorch_mnist.py"
    ],
    resources_per_node={
        "requests": {
            "cpu": "4",
            "memory": "16Gi",
            "nvidia.com/gpu": "1"
        },
        "limits": {
            "nvidia.com/gpu": "1"
        }
    },
    namespace="your-namespace",
    env=[
        {"name": "NCCL_DEBUG", "value": "INFO"}
    ]
)

print(f"✓ Created MPIJob: {job.metadata.name}")
print(f"Status: {job.status.phase}")
```

### Step 3: Run and Monitor

```bash
python mnist_mpi_job.py
```

**Output**:
```
✓ Created MPIJob: mnist-horovod-quickstart
Status: Pending
```

### Step 4: Check Status

```python
# Check job status
job = client.get_job("mnist-horovod-quickstart", namespace="your-namespace")

print(f"Phase: {job.status.phase}")
print(f"Launcher: {job.status.launcher_pod_status.phase}")
print(f"Workers: {len([w for w in job.status.worker_pod_statuses if w.phase == 'Running'])}/{len(job.status.worker_pod_statuses)} ready")

# Wait for completion (blocks until done)
final_job = client.wait_for_job_completion(
    "mnist-horovod-quickstart",
    namespace="your-namespace",
    timeout=600  # 10 minutes
)

if final_job.status.phase == "Succeeded":
    print("✓ Training completed successfully!")
else:
    print(f"✗ Training failed: {final_job.status.failure_reason}")
```

### Step 5: View Logs

```python
# Get launcher logs
launcher_logs = client.get_job_logs(
    "mnist-horovod-quickstart",
    namespace="your-namespace",
    pod_type="launcher"
)
print(launcher_logs)

# Get worker 0 logs
worker_logs = client.get_job_logs(
    "mnist-horovod-quickstart",
    namespace="your-namespace",
    pod_type="worker",
    pod_index=0
)
print(worker_logs)
```

### Step 6: Clean Up

```python
client.delete_job("mnist-horovod-quickstart", namespace="your-namespace")
print("✓ Job deleted")
```

---

## Option 2: CLI (kubectl/oc)

### Step 1: Verify Runtime is Available

```bash
oc get clustertrainingruntimes

# Expected output:
# NAME                MPI IMPLEMENTATION   TIMEOUT   AGE
# mpi-horovod-gpu     OpenMPI              300       5d
# mpi-intel-cpu       IntelMPI             300       5d
```

### Step 2: Create TrainJob YAML

Save this as `mnist-mpijob.yaml`:

```yaml
apiVersion: kubeflow.org/v2alpha1
kind: TrainJob
metadata:
  name: mnist-horovod-quickstart
  namespace: your-namespace
spec:
  runtimeRef:
    kind: ClusterTrainingRuntime
    name: mpi-horovod-gpu
  trainer:
    numNodes: 2
    image: horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-py3.10-gpu
    command:
      - horovodrun
      - -np
      - "2"
      - --verbose
      - python
      - /examples/pytorch/pytorch_mnist.py
    resourcesPerNode:
      requests:
        cpu: "4"
        memory: "16Gi"
        nvidia.com/gpu: "1"
      limits:
        nvidia.com/gpu: "1"
    env:
      - name: NCCL_DEBUG
        value: INFO
```

### Step 3: Submit Job

```bash
oc create -f mnist-mpijob.yaml

# Expected output:
# trainjob.kubeflow.org/mnist-horovod-quickstart created
```

### Step 4: Monitor Status

```bash
# Check job status
oc get trainjob mnist-horovod-quickstart

# Detailed status
oc describe trainjob mnist-horovod-quickstart

# Watch real-time updates
oc get trainjob mnist-horovod-quickstart --watch
```

**Expected Status Transitions**:
```
NAME                       PHASE     AGE   RUNTIME            NODES
mnist-horovod-quickstart   Pending   5s    mpi-horovod-gpu    2
mnist-horovod-quickstart   Running   30s   mpi-horovod-gpu    2
mnist-horovod-quickstart   Succeeded 5m    mpi-horovod-gpu    2
```

### Step 5: View Logs

```bash
# List pods
oc get pods -l training.kubeflow.org/job-name=mnist-horovod-quickstart

# Expected output:
# NAME                                   READY   STATUS    RESTARTS   AGE
# mnist-horovod-quickstart-launcher-0    1/1     Running   0          2m
# mnist-horovod-quickstart-worker-0      1/1     Running   0          2m
# mnist-horovod-quickstart-worker-1      1/1     Running   0          2m

# View launcher logs
oc logs mnist-horovod-quickstart-launcher-0

# View worker 0 logs
oc logs mnist-horovod-quickstart-worker-0

# Follow logs in real-time
oc logs -f mnist-horovod-quickstart-launcher-0
```

### Step 6: Clean Up

```bash
oc delete trainjob mnist-horovod-quickstart

# Verify deletion
oc get trainjob
```

---

## Option 3: Dashboard UI

### Step 1: Access Dashboard

1. Open OpenShift AI Dashboard: `https://your-cluster/dashboard`
2. Navigate to **Training Jobs** section
3. Click **Create Job**

### Step 2: Fill Job Creation Form

**Basic Settings**:
- **Name**: `mnist-horovod-quickstart`
- **Namespace**: Select your project
- **Runtime Type**: MPI
- **Runtime**: `mpi-horovod-gpu`

**Training Configuration**:
- **Number of Nodes**: `2`
- **Container Image**: `horovod/horovod:0.28.1-tf2.11.0-torch2.0.0-py3.10-gpu`
- **Command**: `horovodrun -np 2 --verbose python /examples/pytorch/pytorch_mnist.py`

**Resources per Node**:
- **CPU Request**: `4`
- **Memory Request**: `16Gi`
- **GPU Request**: `1`
- **GPU Limit**: `1`

**Environment Variables** (optional):
- **Name**: `NCCL_DEBUG`, **Value**: `INFO`

### Step 3: Submit and Monitor

1. Click **Create Job**
2. You'll be redirected to the job details page
3. View real-time status updates:
   - **Topology View**: Shows launcher and 2 worker pods
   - **Status Timeline**: Pending → Running → Succeeded
   - **Logs**: Tabbed interface for launcher and worker logs
   - **Metrics**: CPU, memory, GPU utilization graphs

### Step 4: View Results

- **Logs Tab**: Access launcher and worker logs
- **Events Tab**: View lifecycle events
- **Metrics Tab**: See resource utilization over time

### Step 5: Clean Up

- Click **Actions** → **Delete Job**
- Confirm deletion

---

## Understanding What Happened

### Behind the Scenes

When you created the MPIJob, the following happened automatically:

1. **TrainJob Created**: Your job spec submitted to Kubernetes API
2. **Runtime Applied**: Controller merged your job spec with `mpi-horovod-gpu` runtime template
3. **Gang Scheduling**: PodGroup created ensuring launcher + 2 workers schedule atomically
4. **SSH Keys Generated**: Ephemeral SSH key pair created for launcher-worker communication
5. **Pods Launched**:
   - **Launcher Pod**: Runs `mpirun` command to orchestrate training
   - **Worker Pods** (2x): Run SSH servers and execute training workload
6. **MPI Initialization**: Launcher connects to workers via SSH, sets up MPI environment
7. **Training Runs**: Horovod distributes MNIST training across 2 GPUs
8. **Completion**: Pods terminate, SSH keys deleted, status updated to Succeeded

### Key Concepts

**Runtime Templates**: Pre-configured infrastructure blueprints managed by platform admins. You just reference them by name.

**Gang Scheduling**: Ensures all pods (launcher + workers) start together. If resources unavailable, job waits or fails cleanly.

**Launcher vs Workers**:
- **Launcher**: Coordinates distributed training (runs `mpirun`)
- **Workers**: Execute actual training workload under MPI control

**Auto-Configuration**: You only specify high-level parameters (num_nodes, image, resources). Controller handles SSH keys, hostfile, environment variables.

---

## Common Issues and Solutions

### Issue 1: Gang Scheduling Timeout

**Error**: `GangSchedulingTimeout: Requested 3 pods, only 2 nodes available`

**Solution**: Reduce `num_nodes` or add more cluster capacity.

### Issue 2: Image Pull Error

**Error**: `ImagePullBackOff: Failed to pull image horovod/horovod:latest`

**Solution**:
- Check image name and tag are correct
- Add image pull secret if using private registry
- Verify network connectivity to registry

### Issue 3: SSH Connection Failed

**Error**: `SSHConnectionFailed: Launcher could not connect to worker-0`

**Solution**:
- Check network policies allow pod-to-pod traffic on port 22
- Verify Istio sidecar injection is disabled (should be by default in runtime template)
- Check worker pod logs for SSH server startup errors

### Issue 4: Resource Quota Exceeded

**Error**: `ResourceQuotaExceeded: Insufficient GPU quota`

**Solution**:
- Reduce GPUs per node or number of nodes
- Request quota increase from cluster admin
- Use CPU-only training for testing

### Issue 5: Worker OOMKilled

**Error**: `WorkerFailed: worker-1 terminated with reason OOMKilled`

**Solution**:
- Increase memory request in `resourcesPerNode`
- Reduce batch size in training script
- Use gradient accumulation to reduce memory

---

## Next Steps

### 1. Customize Your Training

Replace the example command with your own training script:

```python
job = client.create_mpi_job(
    name="my-custom-training",
    runtime="mpi-horovod-gpu",
    num_nodes=4,
    image="my-registry/my-training-image:v1.0",
    command=["python", "/workspace/train.py"],
    args=["--epochs", "50", "--batch-size", "256"],
    resources_per_node={
        "requests": {
            "cpu": "8",
            "memory": "32Gi",
            "nvidia.com/gpu": "2"
        }
    }
)
```

### 2. Add Data Volumes

Mount PersistentVolumeClaim for training data:

```yaml
spec:
  trainer:
    # ... other fields ...
    volumes:
      - name: training-data
        persistentVolumeClaim:
          claimName: mnist-dataset
    volumeMounts:
      - name: training-data
        mountPath: /data
```

### 3. Scale Up

Increase workers for larger training:

```python
job = client.create_mpi_job(
    name="large-scale-training",
    runtime="mpi-horovod-gpu",
    num_nodes=16,  # 16 workers with 8 GPUs each = 128 GPUs total
    resources_per_node={
        "requests": {
            "cpu": "32",
            "memory": "128Gi",
            "nvidia.com/gpu": "8"
        }
    }
)
```

### 4. Use Different MPI Implementation

Switch to Intel MPI for Intel hardware:

```yaml
spec:
  runtimeRef:
    kind: ClusterTrainingRuntime
    name: mpi-intel-cpu  # Intel MPI optimized for Intel CPUs
```

### 5. Create Custom Runtime

Ask your platform admin to create a custom runtime for your specific needs (custom images, node selectors, tolerations).

---

## Performance Tips

### 1. GPU Utilization

Monitor GPU usage to ensure training is compute-bound:

```python
metrics = client.get_job_metrics("my-job", namespace="your-namespace")

for i, worker in enumerate(metrics['workers']):
    gpu_util = worker['gpu'][-1]['value']
    print(f"Worker {i}: {gpu_util}% GPU utilization")

    if gpu_util < 70:
        print(f"⚠ Warning: Worker {i} GPU underutilized. Consider increasing batch size.")
```

**Target**: >80% GPU utilization for compute-intensive training.

### 2. Scaling Efficiency

Test scaling efficiency before committing to large jobs:

```python
# Test with 2, 4, 8 workers
for num_workers in [2, 4, 8]:
    job = client.create_mpi_job(
        name=f"scaling-test-{num_workers}w",
        runtime="mpi-horovod-gpu",
        num_nodes=num_workers,
        # ... other params
    )

    final_job = client.wait_for_job_completion(job.metadata.name)

    training_time = (final_job.status.completion_time - final_job.status.start_time).total_seconds()
    print(f"{num_workers} workers: {training_time}s")
```

**Ideal**: 2x workers = 2x throughput (50% scaling efficiency is typical beyond 8-16 workers).

### 3. Network Optimization

For multi-node jobs, optimize network:

```python
job = client.create_mpi_job(
    # ... other params ...
    env=[
        # Use NCCL for GPU communication (faster than MPI for GPUs)
        {"name": "HOROVOD_GPU_OPERATIONS", "value": "NCCL"},

        # Enable NCCL debug for troubleshooting
        {"name": "NCCL_DEBUG", "value": "INFO"},

        # Use InfiniBand if available
        {"name": "NCCL_IB_DISABLE", "value": "0"}
    ]
)
```

### 4. Batch Size Tuning

Larger batch sizes = better GPU utilization but may affect convergence:

```python
# Start with batch size = 32 per GPU
base_batch_size = 32
num_gpus = num_nodes * gpus_per_node

# Scale batch size linearly with GPUs
global_batch_size = base_batch_size * num_gpus

job = client.create_mpi_job(
    # ... other params ...
    args=["--batch-size", str(global_batch_size)]
)
```

---

## Troubleshooting Commands

```bash
# Check job status
oc get trainjob <job-name>

# Detailed job information
oc describe trainjob <job-name>

# View job events
oc get events --field-selector involvedObject.name=<job-name>

# Check pod status
oc get pods -l training.kubeflow.org/job-name=<job-name>

# View launcher logs
oc logs <job-name>-launcher-0

# View worker logs
oc logs <job-name>-worker-0

# Check gang scheduling
oc get podgroup <job-name>
oc describe podgroup <job-name>

# Check resource quotas
oc describe resourcequota

# Check network policies
oc get networkpolicies

# Check SSH secret (debug only)
oc get secret <job-name>-ssh
```

---

## What's Next?

✓ You've run your first MPIJob in OpenShift AI!

**Continue Learning**:
- [User Guide](../user-guide.md) - Comprehensive feature documentation
- [Architecture](../architecture.md) - How MPIJob works under the hood
- [Performance Tuning](../performance.md) - Optimize for large-scale training
- [Troubleshooting](../troubleshooting.md) - Common issues and solutions
- [Migration Guide](../migration.md) - Migrate from legacy MPIJob v2beta1

**Get Help**:
- [Community Forum](https://discuss.openshift.com)
- [GitHub Issues](https://github.com/opendatahub-io/opendatahub-operator/issues)
- [Red Hat Support](https://access.redhat.com) (for enterprise customers)

---

## Summary

**What You Learned**:
1. Create MPIJob using Python SDK, CLI, or Dashboard
2. Monitor job status and view logs
3. Understand gang scheduling and auto-configuration
4. Troubleshoot common issues
5. Optimize performance for distributed training

**Key Takeaways**:
- MPIJobs in OpenShift AI are simplified through Trainer V2 runtime templates
- Gang scheduling ensures reliable multi-pod coordination
- Auto-configuration handles SSH keys, hostfiles, and MPI environment
- Three interfaces (SDK, CLI, Dashboard) for different workflows
- Comprehensive observability through logs, events, and metrics

**Time to First Job**: ~15 minutes ✓

Ready to scale your distributed training with MPIJobs in OpenShift AI!
