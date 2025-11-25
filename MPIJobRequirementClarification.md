# MPIJob Requirement Clarification

## 1. Overall Scenario and Requirements

### What is your scenario for using MPIJob?

We use MPIJob to run distributed AI training workloads that require tightly coupled communication between processes, such as large-scale deep learning model training and scientific simulations.

### What makes MPIJob the right solution for your workloads compared to other job types in OpenShift AI KubeFlow Trainer?

MPIJob provides native support for MPI communication patterns, enabling efficient inter-process communication and synchronization across nodes. Unlike generic batch jobs or single-node training jobs, MPIJob supports multi-node parallelism with optimized communication, which is critical for scaling HPC and AI workloads.

### What are your specific technical requirements for MPIJob?

- Support for multi-node MPI execution with process affinity and topology awareness
- Integration with GPU resources and CUDA-aware MPI for optimized GPU communication
- Support for RDMA-capable networking (InfiniBand or RoCE) to reduce latency
- Ability to schedule jobs with gang scheduling semantics to ensure all pods start simultaneously
- Compatibility with shared filesystems for checkpointing and data sharing

## 2. Current Implementation

### How are you using MPIJobs in your environment today?

We currently deploy MPIJobs using the upstream MPI Operator from Kubeflow, customized to our environment. Jobs are submitted via OpenShift AI pipelines for LLM/SLM model continuous pretraining and fine-tuning.

### What modifications have you made to your OpenShift AI cluster to support MPIJobs?

- Integrated RDMA-capable network plugins and enabled InfiniBand support on nodes
- Configured shared NFS volumes for persistent storage
- Tuned kubelet and container runtime settings for low-latency networking and GPU isolation
- Implemented passwordless SSH within MPI pods for MPI runtime communication

### What challenges are you facing with your current setup?

- Scheduling MPIJobs efficiently during peak demand due to lack of gang scheduling support
- Managing network configuration complexity for RDMA and InfiniBand across heterogeneous nodes
- Ensuring library compatibility between MPI versions and CUDA drivers
- Debugging MPI job failures due to network timeouts or resource contention

## 3. Infrastructure and Integration

### Our focus is on integrating MPIJob with Trainer V2. What are your current workflow patterns that need to be supported?

We run multi-step training workflows that include data preprocessing, distributed training via MPIJob, and post-processing. Integration with Trainer V2 should support seamless job submission, monitoring, and artifact management within the pipeline.

### Do you have specific network requirements (RDMA, InfiniBand)?

Yes, we require RDMA support over InfiniBand or RoCE to achieve low-latency, high-throughput communication for MPI workloads.

Additionally, we require support for newest FT/RL algorithms such as PPO/GRPO/GAPO etc. However, some of these algorithms require NCCL_IB and NCCL_P2P to be disabled.

### We are planning Kueue integration for MPIJob scheduling. What scheduling policies and priorities do you need?

- Gang scheduling to ensure all MPI pods start simultaneously
- Priority classes to prioritize critical training jobs over less urgent workloads
- Preemption support to free resources for high-priority MPIJobs
- Fair-share scheduling to balance resource usage across teams

## 4. Compute and Accelerator Requirements

### Do you require GPU support for your MPI workloads?

Yes, GPU acceleration is critical for our deep learning training workloads.

### Do you need CUDA-aware MPI or GPU-optimized communication?

Yes, CUDA-aware MPI and GPUDirect RDMA are required to minimize data transfer overhead between GPUs across nodes.

### What are your typical resource requirements (CPU, memory, GPU) per job?

- **CPU**: 8-40 cores per pod depending on workload
- **Memory**: 1942-4096 GB per pod
- **GPU**: 8 NVIDIA H100 or H200 GPUs per pod

### What scaling requirements do you have (number of pods/processes)?

We scale from 2 to 6 pods per MPIJob, consisting of 1 launcher pod and 1 to 5 GPU worker nodes. With 5 nodes and 40 GPUs, this requires up to 40 parallel processes across all GPUs.

## 5. Software and Filesystem Requirements

### What shared filesystem do you use?

We use MinIO for shared storage accessible by all MPI pods.

### Do you need passwordless SSH or other communication channels?

Yes, passwordless SSH is required within MPI pods to enable MPI runtime communication and process management.

### What library compatibility requirements do you have?

- MPI implementations compatible with OpenMPI 4.x or MVAPICH2
- CUDA Toolkit versions aligned with GPU drivers (e.g., CUDA 11.x)
- NCCL for optimized collective communication

### What software dependencies are critical for your workloads?

- Python 3.11+ with TensorFlow or PyTorch
- MPI libraries (OpenMPI)
- CUDA and cuDNN libraries
- Monitoring and logging agents compatible with OpenShift

## 6. Performance and Scale

### What are your typical job sizes (number of processes/nodes)?

Jobs typically run on 2 to 5 nodes, with 8-16 processes per node.

### What performance benchmarks do you need to meet?

- Achieve near-linear scaling efficiency up to 128 GPUs
- Maintain network latency below 2 microseconds for MPI messages
- Complete training jobs within SLA-defined time windows (e.g., 8-12 hours)

### How do you handle peak demand and resource allocation?

- Use priority-based scheduling with preemption to allocate resources dynamically
- Employ Kueue for queue management and gang scheduling to optimize resource utilization
- Implement autoscaling policies for GPU nodes based on workload demand

## 7. Additional Requirements Summary

### Network Configuration

RDMA support over InfiniBand or RoCE is required to achieve the necessary low-latency, high-throughput communication for MPI workloads.

### Resource Management and Scheduling

We intend to use Kueue for orchestrating MPIJob scheduling. Gang scheduling is essential to ensure all MPI pods launch simultaneously, preventing partial job starts and resource wastage. Priority classes and preemption support are also desired to manage workload priorities effectively.

### GPU and Accelerator Integration

CUDA-aware MPI and GPUDirect RDMA are critical for optimizing GPU-to-GPU communication across nodes, reducing overhead and improving training throughput.

### Shared Filesystem

We use NFS as shared filesystems accessible by all MPI pods for data sharing, checkpointing, and logging.

### Additional Software Requirements

- Passwordless SSH within MPI pods for MPI runtime communication
- Compatibility with MPI libraries such as OpenMPI 4.x or MVAPICH2
- Alignment of CUDA Toolkit and driver versions to ensure compatibility with GPU-accelerated MPI
