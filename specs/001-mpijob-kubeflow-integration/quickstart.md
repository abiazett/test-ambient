# MPIJob Quickstart Guide

**Target Audience**: Data scientists and ML engineers new to distributed training with MPIJobs
**Prerequisites**: OpenShift AI 2.17+, Kubernetes cluster with GPU nodes, RBAC access to create MPIJobs
**Time to Complete**: 30 minutes
**Goal**: Train a PyTorch model using Horovod across 4 workers with GPUs

---

## Table of Contents

1. [Prerequisites](#prerequisites)
2. [Verify Cluster Setup](#verify-cluster-setup)
3. [Method 1: Python SDK (Jupyter Notebook)](#method-1-python-sdk-jupyter-notebook)
4. [Method 2: CLI](#method-2-cli)
5. [Method 3: Dashboard UI](#method-3-dashboard-ui)
6. [Monitoring Your Job](#monitoring-your-job)
7. [Troubleshooting](#troubleshooting)
8. [Next Steps](#next-steps)

---

## Prerequisites

### Required Components

✅ **OpenShift AI 2.17+** installed
✅ **Kubernetes 1.27+** cluster
✅ **GPU nodes** with NVIDIA GPU Operator
✅ **Training Operator v2** deployed (managed by RHOAI operator)
✅ **Kueue** (optional, recommended for gang scheduling)

### RBAC Permissions

Verify you have permissions to create MPIJobs:

```bash
oc auth can-i create mpijobs.kubeflow.org -n <your-namespace>
# Expected output: yes
```

If "no", contact your cluster administrator to grant `mpijob-user` role.

### Training Data

For this quickstart, we'll use a sample CIFAR-10 dataset that downloads automatically. For production training, you'd mount a PVC with your dataset.

---

## Verify Cluster Setup

### Check GPU Availability

```bash
# List GPU nodes
oc get nodes -l nvidia.com/gpu.present=true

# Check total GPUs available
oc describe nodes -l nvidia.com/gpu.present=true | grep nvidia.com/gpu
```

Expected output:
```
nvidia.com/gpu:  4
nvidia.com/gpu:  4
```

### Check Training Operator

```bash
oc get deployment training-operator -n redhat-ods-applications
```

Expected output:
```
NAME                 READY   UP-TO-DATE   AVAILABLE   AGE
training-operator    1/1     1            1           5d
```

### Check ResourceQuota (if applicable)

```bash
oc get resourcequota -n <your-namespace>
oc describe resourcequota <quota-name> -n <your-namespace>
```

Ensure you have sufficient GPU quota available.

---

## Method 1: Python SDK (Jupyter Notebook)

### Step 1: Install SDK

In your Jupyter notebook, run:

```python
!pip install openshift-ai-sdk
```

### Step 2: Create Training Script

Create a file `train.py` with this simple PyTorch + Horovod example:

```python
# train.py
import torch
import torch.nn as nn
import torch.optim as optim
import horovod.torch as hvd
from torchvision import datasets, transforms

# Initialize Horovod
hvd.init()

# Pin GPU to be used to process local rank (one GPU per process)
torch.cuda.set_device(hvd.local_rank())

# Define a simple CNN
class SimpleCNN(nn.Module):
    def __init__(self):
        super(SimpleCNN, self).__init__()
        self.conv1 = nn.Conv2d(3, 32, kernel_size=3, padding=1)
        self.conv2 = nn.Conv2d(32, 64, kernel_size=3, padding=1)
        self.fc1 = nn.Linear(64 * 8 * 8, 128)
        self.fc2 = nn.Linear(128, 10)
        self.pool = nn.MaxPool2d(2, 2)
        self.relu = nn.ReLU()

    def forward(self, x):
        x = self.pool(self.relu(self.conv1(x)))
        x = self.pool(self.relu(self.conv2(x)))
        x = x.view(-1, 64 * 8 * 8)
        x = self.relu(self.fc1(x))
        x = self.fc2(x)
        return x

# Load CIFAR-10 dataset
transform = transforms.Compose([
    transforms.ToTensor(),
    transforms.Normalize((0.5, 0.5, 0.5), (0.5, 0.5, 0.5))
])

train_dataset = datasets.CIFAR10(root='/tmp/data', train=True, download=True, transform=transform)

# Partition dataset among workers
train_sampler = torch.utils.data.distributed.DistributedSampler(
    train_dataset, num_replicas=hvd.size(), rank=hvd.rank()
)

train_loader = torch.utils.data.DataLoader(
    train_dataset, batch_size=128, sampler=train_sampler
)

# Build model
model = SimpleCNN().cuda()
optimizer = optim.SGD(model.parameters(), lr=0.01 * hvd.size())

# Wrap optimizer with Horovod DistributedOptimizer
optimizer = hvd.DistributedOptimizer(optimizer, named_parameters=model.named_parameters())

# Broadcast parameters from rank 0 to all other processes
hvd.broadcast_parameters(model.state_dict(), root_rank=0)
hvd.broadcast_optimizer_state(optimizer, root_rank=0)

# Loss function
criterion = nn.CrossEntropyLoss()

# Training loop
model.train()
for epoch in range(5):
    train_sampler.set_epoch(epoch)
    for batch_idx, (data, target) in enumerate(train_loader):
        data, target = data.cuda(), target.cuda()

        optimizer.zero_grad()
        output = model(data)
        loss = criterion(output, target)
        loss.backward()
        optimizer.step()

        if batch_idx % 10 == 0 and hvd.rank() == 0:
            print(f'Epoch {epoch}, Batch {batch_idx}, Loss: {loss.item():.4f}')

if hvd.rank() == 0:
    print('Training complete!')
    torch.save(model.state_dict(), '/tmp/model.pth')
```

### Step 3: Build Container Image

Create a `Dockerfile`:

```dockerfile
FROM nvcr.io/nvidia/pytorch:25.01-py3

# Install Horovod with MPI support
RUN pip install horovod[pytorch]

# Install OpenSSH for MPI communication
RUN apt-get update && apt-get install -y openssh-server
RUN mkdir /var/run/sshd

# Copy training script
COPY train.py /workspace/train.py

WORKDIR /workspace
```

Build and push:

```bash
docker build -t quay.io/<your-org>/pytorch-horovod:quickstart .
docker push quay.io/<your-org>/pytorch-horovod:quickstart
```

### Step 4: Submit MPIJob from Jupyter

In your Jupyter notebook:

```python
from openshift_ai_sdk.training import MPIJob, ResourceRequirements

# Create MPIJob configuration
job = MPIJob(
    name="quickstart-pytorch",
    namespace="<your-namespace>",  # Replace with your namespace
    workers=4,
    image="quay.io/<your-org>/pytorch-horovod:quickstart",  # Replace with your image
    command=["horovodrun", "-np", "4", "--verbose", "python", "/workspace/train.py"],
    resources=ResourceRequirements(
        cpu=8,
        memory="32Gi",
        gpu=1
    ),
)

# Submit job
print("Submitting MPIJob...")
job.create()
print(f"✓ Job '{job.name}' created successfully!")
```

### Step 5: Monitor Progress

```python
# Check status
import time

while True:
    status = job.get_status()
    print(f"Job status: {status}")

    if status in ["Succeeded", "Failed"]:
        break

    time.sleep(10)

if status == "Succeeded":
    print("✓ Training completed successfully!")
else:
    print("✗ Training failed. Check logs for details.")
```

### Step 6: View Logs

```python
# Get launcher logs
logs = job.get_logs()
print("=== Launcher Logs ===")
print(logs)

# Get worker 0 logs
worker_logs = job.get_logs(worker_index=0)
print("=== Worker 0 Logs ===")
print(worker_logs)
```

### Step 7: Cleanup

```python
# Delete job
job.delete()
print("✓ Job deleted successfully!")
```

---

## Method 2: CLI

### Step 1: Create YAML Configuration

Create `mpijob-quickstart.yaml`:

```yaml
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: quickstart-pytorch
  namespace: <your-namespace>
  annotations:
    sidecar.istio.io/inject: "false"  # IMPORTANT: Disable Istio injection
spec:
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      template:
        spec:
          containers:
          - name: launcher
            image: quay.io/<your-org>/pytorch-horovod:quickstart
            command:
            - horovodrun
            - -np
            - "4"
            - --verbose
            - python
            - /workspace/train.py
    Worker:
      replicas: 4
      template:
        spec:
          containers:
          - name: worker
            image: quay.io/<your-org>/pytorch-horovod:quickstart
            resources:
              requests:
                cpu: "8"
                memory: "32Gi"
                nvidia.com/gpu: "1"
              limits:
                cpu: "8"
                memory: "32Gi"
                nvidia.com/gpu: "1"
          nodeSelector:
            nvidia.com/gpu.present: "true"
  cleanPodPolicy: Running
  slotsPerWorker: 1
```

### Step 2: Submit Job

```bash
# Using kubectl
kubectl apply -f mpijob-quickstart.yaml

# Or using oc CLI
oc apply -f mpijob-quickstart.yaml
```

### Step 3: Monitor Progress

```bash
# Check job status
oc get mpijobs quickstart-pytorch -n <your-namespace>

# Watch status changes
oc get mpijobs quickstart-pytorch -n <your-namespace> -w

# Check worker pods
oc get pods -l training.kubeflow.org/job-name=quickstart-pytorch -n <your-namespace>
```

### Step 4: View Logs

```bash
# Launcher logs
oc logs quickstart-pytorch-launcher -n <your-namespace>

# Worker 0 logs
oc logs quickstart-pytorch-worker-0 -n <your-namespace>

# Follow logs in real-time
oc logs -f quickstart-pytorch-launcher -n <your-namespace>
```

### Step 5: Cleanup

```bash
oc delete mpijob quickstart-pytorch -n <your-namespace>
```

---

## Method 3: Dashboard UI

### Step 1: Navigate to Dashboard

1. Open OpenShift AI Dashboard: `https://odh-dashboard.apps.<cluster-domain>`
2. Log in with your OpenShift credentials
3. Select your namespace from the dropdown

### Step 2: Create MPIJob

1. Click **"Training Jobs"** in the left navigation
2. Click **"Create MPIJob"** button
3. Fill in the job creation wizard:

**Basic Information**:
- Name: `quickstart-pytorch`
- Namespace: `<your-namespace>`

**Container Configuration**:
- Image: `quay.io/<your-org>/pytorch-horovod:quickstart`
- Command: `horovodrun -np 4 --verbose python /workspace/train.py`

**Resources**:
- Workers: `4`
- CPU per worker: `8`
- Memory per worker: `32Gi`
- GPU per worker: `1`

**Advanced** (optional):
- Node Selector: `nvidia.com/gpu.present: "true"`
- Clean Pod Policy: `Running`
- Slots Per Worker: `1`

4. Click **"Review"** to see generated YAML
5. Click **"Submit"** to create the job

### Step 3: Monitor Progress

The Dashboard will automatically redirect you to the job details page:

- **Overview Tab**: Shows job status, duration, worker count
- **Workers Tab**: Lists all worker pods with their status and node placement
- **Logs Tab**: View logs from launcher and workers
- **Metrics Tab**: Real-time GPU utilization, memory usage (requires Prometheus)
- **Events Tab**: Kubernetes events related to the job

### Step 4: View Results

Once the job status changes to **"Succeeded"**:
1. Click the **"Logs"** tab
2. View training output and final accuracy
3. Download logs if needed

### Step 5: Cleanup

1. Click **"Actions"** dropdown
2. Select **"Delete Job"**
3. Confirm deletion

---

## Monitoring Your Job

### Using OpenShift Console

1. Navigate to **Workloads → Pods**
2. Filter by label: `training.kubeflow.org/job-name=quickstart-pytorch`
3. View pod details, logs, and events

### Using Grafana Dashboards

If Prometheus/Grafana is configured:

1. Navigate to Grafana: `https://grafana.apps.<cluster-domain>`
2. Open **"MPIJob Overview"** dashboard
3. Select your job from the dropdown
4. View metrics:
   - GPU utilization per worker
   - Memory usage
   - Network bandwidth
   - Training throughput

### Using Python SDK

```python
# Get worker statuses
worker_statuses = job.get_worker_statuses()
for worker in worker_statuses:
    print(f"Worker {worker['index']}: {worker['phase']} on node {worker['node']}")

# Get performance metrics
metrics = job.get_metrics()
print(f"Average GPU utilization: {metrics['aggregated']['avgGpuUtilization']:.2f}%")
```

---

## Troubleshooting

### Job Stuck in "Pending"

**Symptom**: MPIJob status is "Pending" for >5 minutes

**Diagnosis**:
```bash
# Check worker pod status
oc get pods -l training.kubeflow.org/job-name=quickstart-pytorch

# Describe pending pods
oc describe pod quickstart-pytorch-worker-0
```

**Common Causes**:
1. **Insufficient GPUs**: Requested GPUs exceed available
   - **Solution**: Reduce worker count or wait for resources
2. **Image Pull Failure**: Container image not found or auth failed
   - **Solution**: Verify image exists, check registry credentials
3. **Quota Exceeded**: Namespace GPU quota exhausted
   - **Solution**: Delete other jobs or request quota increase
4. **Node Selector Mismatch**: No nodes match selector
   - **Solution**: Check node labels with `oc get nodes --show-labels`

### MPI Communication Failure

**Symptom**: Launcher logs show "ssh: connect to host worker-0 port 2222: Connection refused"

**Diagnosis**:
```bash
# Check if workers are running
oc get pods -l training.kubeflow.org/replica-type=worker

# Check NetworkPolicy
oc get networkpolicies -n <your-namespace>

# Check SSH service in worker pod
oc exec quickstart-pytorch-worker-0 -- ps aux | grep sshd
```

**Common Causes**:
1. **NetworkPolicy Blocking**: SSH port 2222 blocked
   - **Solution**: Apply MPIJob NetworkPolicy template
2. **SSH Server Not Running**: Container image missing SSH server
   - **Solution**: Ensure image includes `openssh-server` and starts sshd
3. **DNS Resolution Failure**: Worker DNS names not resolving
   - **Solution**: Verify StatefulSet created workers (not Deployment)

### Worker Pod OOMKilled

**Symptom**: Worker pod terminated with "OOMKilled" reason

**Diagnosis**:
```bash
oc describe pod quickstart-pytorch-worker-0 | grep -A5 "Last State"
```

**Solution**:
- Increase memory request/limit in job configuration
- Reduce batch size in training script
- Use gradient accumulation to train with effective larger batch size

### Training Accuracy Poor

**Symptom**: Model trains but achieves low accuracy

**Common Causes**:
1. **Learning Rate Too High**: Distributed training requires LR adjustment
   - **Solution**: Scale LR by number of workers (e.g., `lr * hvd.size()`)
2. **Gradient Accumulation Issues**: Gradients not synchronized correctly
   - **Solution**: Verify Horovod `DistributedOptimizer` wraps base optimizer
3. **Data Partitioning**: Same data processed by all workers
   - **Solution**: Use `DistributedSampler` to partition dataset

### Performance Issues

**Symptom**: Training slower than expected or scaling efficiency <70%

**Diagnosis**:
```python
# Check GPU utilization
metrics = job.get_metrics()
for worker in metrics['workers']:
    print(f"Worker {worker['workerIndex']}: {worker['gpuUtilization']:.2f}% GPU")
```

**Common Causes**:
1. **Data Loading Bottleneck**: CPU-bound data loading
   - **Solution**: Increase DataLoader `num_workers`, use prefetching
2. **Small Batch Size**: GPU underutilized
   - **Solution**: Increase batch size per worker
3. **Network Bandwidth**: MPI Allreduce latency
   - **Solution**: Check network with `oc exec quickstart-pytorch-launcher -- iperf3 -c quickstart-pytorch-worker-0`

---

## Next Steps

### Advanced Features

1. **Hyperparameter Tuning**: Integrate with Katib for automated HPO
2. **Model Checkpointing**: Mount PVC to save checkpoints during training
3. **Multi-Node LLM Training**: Scale to 8+ nodes with DeepSpeed
4. **RDMA Networking**: Configure InfiniBand for 16× performance improvement
5. **Fault Tolerance**: Enable checkpoint/resume for long-running jobs

### Example Configurations

**8-Node LLM Fine-Tuning**:
```yaml
spec:
  mpiReplicaSpecs:
    Worker:
      replicas: 8
      template:
        spec:
          containers:
          - name: worker
            resources:
              limits:
                nvidia.com/gpu: "4"  # 32 GPUs total
                memory: "256Gi"
            env:
            - name: NCCL_DEBUG
              value: "INFO"
            - name: NCCL_IB_DISABLE
              value: "1"  # Disable InfiniBand if not available
```

**Checkpoint to S3**:
```python
job = MPIJob(
    # ... other config ...
    env={
        "S3_BUCKET": "my-checkpoints",
        "S3_PREFIX": "quickstart-pytorch",
        "AWS_ACCESS_KEY_ID": "<key-id>",
        "AWS_SECRET_ACCESS_KEY": "<secret-key>",
    },
)
```

### Documentation Resources

- **API Reference**: `/docs/reference/mpijob-crd-spec.md`
- **CLI Commands**: `/docs/reference/cli-reference.md`
- **SDK Documentation**: `/docs/reference/sdk-api.md`
- **Troubleshooting Guide**: `/docs/troubleshooting/mpijobs-common-issues.md`
- **Architecture**: `/docs/architecture/mpijobs-design.md`

### Community & Support

- **Slack**: `#rhoai-mpijobs` channel
- **GitHub Issues**: `github.com/opendatahub-io/training-operator`
- **RedHat Support**: Open support case via Customer Portal

---

## Summary

Congratulations! You've successfully:

✅ Submitted your first MPIJob
✅ Monitored job progress with CLI/SDK/Dashboard
✅ Viewed training logs
✅ Troubleshot common issues

**Time to First Job**: ~30 minutes (as designed!)

**Key Takeaways**:
1. MPIJobs require 2+ workers and MPI-compatible container images
2. Three interfaces available: Python SDK (Jupyter), CLI (kubectl/oc), Dashboard UI
3. Gang scheduling (Kueue) prevents resource deadlocks
4. Always disable Istio sidecar injection with annotation
5. Scale learning rate linearly with worker count

**Next**: Try training your own model with your dataset! 🚀
