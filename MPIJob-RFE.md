# MPIJob Support for Red Hat OpenShift AI

**Author(s):** [To be filled]
**Tracking Link:** <Jira URL>
**Other references:**
- [KubeFlow Training Operator Documentation](https://www.kubeflow.org/docs/components/training/)
- [MPI Operator GitHub Repository](https://github.com/kubeflow/mpi-operator)

---

# Summary - The "What"

This RFE proposes adding support for MPIJobs (Message Passing Interface Jobs) to Red Hat OpenShift AI through integration with KubeFlow Trainer V2. This enhancement will enable data scientists and MLOps engineers to run distributed AI training workloads that require tightly coupled, multi-node parallelism with optimized inter-process communication. The integration will provide MPIJob capabilities through CLI, SDK, and the OpenShift AI Dashboard UI, offering unified observability and management alongside other KubeFlow Trainer job types.

# Rationale - The "Why"

## Problem Statement

Organizations running large-scale AI training workloads, particularly for large language models (LLMs) and scientific simulations, require efficient multi-node distributed training capabilities. These workloads demand low-latency, high-throughput communication between processes that cannot be adequately served by existing job types in OpenShift AI. Currently, users must implement custom solutions using upstream KubeFlow MPI Operator with manual cluster modifications, leading to:

- **Operational complexity**: Teams must manually integrate and maintain the MPI Operator outside the standard OpenShift AI workflow
- **Inconsistent user experience**: MPIJobs lack the unified management, monitoring, and observability available for other training job types
- **Integration overhead**: Custom implementations require extensive cluster modifications for networking, storage, and scheduling that are not standardized or supported
- **Limited enterprise adoption**: Lack of native support creates barriers for enterprises seeking to leverage OpenShift AI for HPC-grade AI workloads

## Expected Impact

This feature addresses a critical gap for customers running enterprise-scale AI workloads, particularly in sectors such as:
- Financial services (fraud detection, risk modeling)
- Healthcare and life sciences (drug discovery, genomics)
- Research institutions (scientific simulations, climate modeling)
- Technology companies (LLM pre-training and fine-tuning)

The market opportunity is significant as distributed training becomes essential for modern AI applications. Customers are increasingly demanding integrated solutions that combine HPC-grade performance with enterprise Kubernetes management capabilities.

## Current State of MPIJob Support

### Challenges with Current Approaches

Today, organizations requiring MPI-based distributed training must operate outside the standard OpenShift AI experience:

1. **Manual Operator Installation**: Users must deploy and maintain the upstream KubeFlow MPI Operator separately, managing its lifecycle independently from OpenShift AI updates

2. **Complex Cluster Modifications**: Enabling MPIJobs requires significant infrastructure changes:
   - RDMA-capable network plugins and InfiniBand configuration
   - Shared filesystem setup (NFS, MinIO) for data sharing and checkpointing
   - Passwordless SSH configuration within pods for MPI runtime communication
   - Custom kubelet and container runtime tuning for low-latency networking
   - GPU isolation and CUDA-aware MPI library integration

3. **Scheduling Limitations**: Without integrated gang scheduling support, users face:
   - Partial job starts when resources are limited
   - Resource wastage from incomplete pod sets
   - Poor scheduling efficiency during peak demand periods

4. **Fragmented Observability**: MPIJobs exist outside the OpenShift AI Dashboard, requiring separate monitoring tools and creating disconnected user experiences

5. **Library and Compatibility Management**: Users must manually ensure compatibility between:
   - MPI implementations (OpenMPI, MVAPICH2) and CUDA drivers
   - Network layers (RDMA, InfiniBand, RoCE) and MPI libraries
   - NCCL versions and specific algorithm requirements (some newer RL algorithms like PPO/GRPO/GAPO require NCCL_IB and NCCL_P2P to be disabled)

### Why MPIJob is Essential

MPIJob provides capabilities that other job types cannot deliver:

- **Native MPI Communication Patterns**: Enables efficient collective operations (broadcast, reduce, gather) across distributed processes
- **Multi-node Parallelism**: Supports scaling to dozens or hundreds of nodes with synchronized execution
- **CUDA-aware MPI**: Optimizes GPU-to-GPU communication across nodes, critical for multi-GPU training
- **Process Topology Awareness**: Respects network and hardware topology for optimal communication patterns
- **Gang Scheduling Semantics**: Ensures all processes start simultaneously, preventing deadlocks and resource waste

These capabilities are fundamental for workloads such as:
- Large language model pre-training and fine-tuning requiring 40+ GPUs across multiple nodes
- Scientific simulations with tightly coupled computational meshes
- Reinforcement learning algorithms with distributed data collection and training

### Technical Requirements Overview

Based on customer feedback and upstream KubeFlow Trainer V2 capabilities, the integration should support:

- **Multi-node Execution**: 2-5 nodes per job with 8-40 cores and 8 GPUs per node
- **Network Optimization**: RDMA over InfiniBand or RoCE for sub-microsecond latency
- **GPU Acceleration**: CUDA-aware MPI and GPUDirect RDMA for efficient GPU communication
- **Shared Storage**: Integration with shared filesystems (NFS, MinIO) for data and checkpoints
- **Scheduling**: Gang scheduling through Kueue integration with priority classes and preemption
- **User Interfaces**: CLI, SDK (Python), and Dashboard UI access with unified observability

### Prerequisite Dependency

**This RFE depends on KubeFlow Trainer V2 support being implemented in Red Hat OpenShift AI first.** MPIJob integration will build upon the Trainer V2 foundation to ensure consistent architecture, APIs, and user experience across all supported job types.

# Red Hat Alignment

This RFE aligns with Red Hat's corporate strategy in several key areas:

**Product Enhancement**: This feature directly enhances Red Hat OpenShift AI by expanding its distributed training capabilities to support HPC-grade workloads. It positions OpenShift AI as a comprehensive platform for both standard and advanced AI training scenarios.

**Open Source Leadership**: By integrating upstream KubeFlow Trainer V2 and MPI Operator capabilities, Red Hat continues its commitment to contributing to and adopting open source AI/ML ecosystem components while adding enterprise-grade support and integration.

**Customer Success**: This addresses a documented gap for customers requiring distributed training at scale, particularly those in regulated industries who need enterprise support for production AI workloads.

**Competitive Positioning**: Native MPIJob support differentiates OpenShift AI from competing platforms by offering integrated HPC-grade distributed training without requiring custom implementations or third-party tools.

# Potential Customers

Target customer segments include:

1. **Financial Services**: Organizations training large-scale models for fraud detection, risk assessment, and algorithmic trading
2. **Healthcare and Life Sciences**: Research institutions running distributed simulations for drug discovery and genomics analysis
3. **Technology Companies**: Teams performing LLM pre-training and fine-tuning requiring multi-node GPU clusters
4. **Research Institutions**: Universities and national labs running scientific simulations and climate modeling
5. **Manufacturing and Engineering**: Companies using AI for simulation, optimization, and digital twin applications

# Current Customers/Partners & Scope

Current users running MPIJobs outside OpenShift AI report typical workload characteristics:
- Job sizes: 2-5 nodes with 8-16 processes per node
- GPU requirements: 8 NVIDIA H100/H200 GPUs per node
- Resource profiles: 8-40 CPU cores, 192-4096 GB memory per pod
- Scaling targets: Up to 128 GPUs with near-linear scaling efficiency
- Performance requirements: Network latency below 2 microseconds for MPI messages

These users currently implement custom solutions and would benefit from native integration with standardized support.

# Alternatives Considered

## 1. Continue Using Upstream MPI Operator Independently

**Approach**: Users deploy and maintain the upstream KubeFlow MPI Operator alongside OpenShift AI.

**Limitations**:
- Fragmented user experience with separate management interfaces
- No integration with OpenShift AI Dashboard for unified observability
- Increased operational burden for users to maintain compatibility across upgrades
- Lack of enterprise support for the integrated stack

## 2. Use PyTorchJob or TFJob with Custom Distributed Communication

**Approach**: Implement MPI-style communication patterns within PyTorchJob or TFJob frameworks.

**Limitations**:
- These frameworks are optimized for their respective ML libraries, not general HPC workloads
- Performance overhead from non-native MPI implementations
- Limited support for RDMA and specialized networking hardware
- Cannot support non-Python or legacy HPC applications requiring standard MPI

## 3. Custom Job CRDs

**Approach**: Create custom Kubernetes CRDs and controllers specifically for distributed MPI workloads.

**Limitations**:
- Duplicates functionality already available in KubeFlow ecosystem
- Increases maintenance burden for Red Hat
- Fragments the user experience across multiple custom job types
- Lacks community support and upstream contribution opportunities

## Why KubeFlow Trainer V2 Integration is the Right Approach

Integrating MPIJob through KubeFlow Trainer V2 provides:
- **Unified Architecture**: Consistent APIs and user experience across all job types
- **Community Alignment**: Leverages upstream KubeFlow development and best practices
- **Enterprise Support**: Enables Red Hat to provide supported, tested integration
- **Reduced Complexity**: Single integration point rather than multiple disparate systems
- **Future Flexibility**: Foundation for supporting additional job types as KubeFlow ecosystem evolves

---

This RFE document provides the high-level requirements and rationale for MPIJob support. Detailed technical specifications, implementation approaches, and design decisions will be developed during the specification phase following RFE approval.
