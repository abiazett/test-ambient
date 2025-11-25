# Feature Refinement - MPIJob Support for Red Hat OpenShift AI

---

## Feature Metadata

| Field | Value |
|-------|-------|
| **Feature Jira Link** | <Add link to the feature (XXXSTRAT-####)> |
| **Status** | Requirements Refinement |
| **Slack Channel / Thread** | <Add slack channel or discussion thread if applicable> |
| **Feature Owner** | <Person> |
| **Delivery Owner** | <Person> |
| **RFE Council Reviewer** | <Person> |
| **Product** | RHOAI (self-managed and managed) |

---

## Feature Details

### Feature Overview

MPIJob support enables data scientists and MLOps engineers to execute distributed AI training workloads requiring tightly coupled, multi-node parallelism using the Message Passing Interface (MPI) protocol. This integration extends Red Hat OpenShift AI's distributed training capabilities through the KubeFlow Trainer V2 framework, providing native support for HPC-grade workloads that demand low-latency inter-process communication.

**Key Technical Capabilities:**
- **Multi-node MPI execution**: Launcher-worker architecture where a launcher pod coordinates MPI runtime across multiple worker pods
- **MPI implementation flexibility**: Support for OpenMPI (default), Intel MPI, and MPICH
- **Gang scheduling integration**: Built-in support for Kueue to ensure all pods start simultaneously
- **GPU-aware communication**: CUDA-aware MPI for optimized GPU-to-GPU data transfer across nodes
- **SSH-based coordination**: Automatic SSH key management for MPI process communication
- **Resource management**: Configurable slots per worker, flexible replica specifications, and cleanup policies

**User Personas and Benefits:**

1. **Data Scientists**
   - Submit MPIJobs through Dashboard UI or Python SDK without understanding Kubernetes internals
   - Monitor training progress through unified observability dashboards
   - Access logs and metrics from launcher and worker pods in a single view
   - Leverage pre-configured container images with MPI libraries and CUDA support

2. **MLOps Engineers**
   - Deploy and manage MPIJobs via CLI and GitOps workflows
   - Configure gang scheduling policies for efficient resource utilization
   - Integrate MPIJobs into automated training pipelines
   - Debug distributed training failures with comprehensive status reporting

3. **Platform Administrators**
   - Manage MPI capabilities lifecycle through OpenShift AI operator
   - Configure cluster-level resources (network policies, storage classes, GPU quotas)
   - Monitor cluster-wide MPIJob resource consumption and performance
   - Enforce security policies for SSH key management and pod security standards

**Current State vs. Future State:**

*Current State:*
- Users must manually deploy upstream KubeFlow MPI Operator outside OpenShift AI
- No integration with OpenShift AI Dashboard - MPIJobs are invisible to the UI
- Manual configuration of networking, storage, and scheduling requirements
- Fragmented observability requiring separate monitoring tools
- No enterprise support for the integrated stack
- Complex cluster modifications required (RDMA setup, SSH configuration, shared storage)

*Future State:*
- MPIJob support delivered as native OpenShift AI capability through Trainer V2
- Unified job submission through Dashboard UI, CLI (`oai-cli`), and Python SDK
- Integrated observability with logs, metrics, and status in Dashboard alongside other job types
- Automated SSH key management and network configuration
- Gang scheduling integration with Kueue for efficient resource allocation
- Enterprise-supported, tested integration with Red Hat support coverage

### The Why

**Technical Gaps Addressed:**

Organizations running large-scale AI training workloads, particularly for large language models (LLMs) and scientific simulations, require efficient multi-node distributed training capabilities that current OpenShift AI job types cannot adequately serve.

**Current Limitations:**
1. **Communication Pattern Gaps**: PyTorchJob and TFJob optimize for data-parallel training with parameter server or all-reduce patterns but lack native support for:
   - Collective MPI operations (MPI_Bcast, MPI_Reduce, MPI_Allgather)
   - Point-to-point communication primitives required by scientific simulations
   - Topology-aware process placement for RDMA and InfiniBand networks
   - CUDA-aware MPI for direct GPU-to-GPU memory transfers

2. **Scheduling Coordination**: MPIJobs require all worker pods ready simultaneously. Without gang scheduling, partial pod sets lead to launcher failures and wasted resources.

3. **Performance Requirements**: Modern LLM training and HPC workloads demand sub-microsecond MPI message latency and near-linear scaling efficiency to 128+ GPUs.

**Why Now:**

1. **Market Demand**
   - LLM pre-training and fine-tuning workloads increasingly require 40+ GPUs across multiple nodes
   - Customers in financial services, healthcare, and research sectors explicitly requesting MPI support
   - Competitive platforms (Azure ML, SageMaker) offer integrated distributed training with MPI

2. **Technical Foundation**
   - KubeFlow Trainer V2 provides standardized architecture for job type integration
   - Upstream MPI Operator has matured with stable v2beta1 API
   - OpenShift infrastructure improvements (Multus CNI, SR-IOV) enable high-performance networking

3. **Ecosystem Alignment**
   - KubeFlow community actively maintains MPI Operator with regular releases
   - Integration aligns with Red Hat's open source strategy
   - Trainer V2 dependency establishes consistent patterns for future job types

**Customer Impact:**

Target customer segments include:
- **Financial Services**: Fraud detection, risk assessment, algorithmic trading requiring large-scale model training
- **Healthcare and Life Sciences**: Drug discovery, genomics analysis with distributed simulations
- **Technology Companies**: LLM pre-training and fine-tuning requiring multi-node GPU clusters
- **Research Institutions**: Scientific simulations and climate modeling
- **Manufacturing**: AI for simulation, optimization, and digital twin applications

### High Level Requirements

**Core Functionality:**

1. **As a Data Scientist**, I want to create MPIJob training workloads through the OpenShift AI Dashboard UI, CLI, and SDK, So that I can run distributed AI training with MPI orchestration using the same interface I use for other job types.

2. **As an MLOps Engineer**, I want to leverage KubeFlow Trainer V2's MPIJob capabilities through OpenShift AI, So that I can orchestrate multi-node distributed training with native MPI communication patterns without manual operator installation.

3. **As a Platform Engineer**, I want to deploy pre-configured ClusterTrainingRuntimes for MPI-based training, So that data scientists can launch MPIJob workloads without understanding complex Kubernetes and MPI configuration details.

4. **As a DevOps Engineer**, I want unified observability for MPIJobs through the OpenShift AI Dashboard, So that I can monitor MPI training jobs alongside PyTorchJobs and other training workloads in a single interface.

**Integration & Scheduling:**

5. **As a Platform Administrator**, I want MPIJobs to integrate with Kueue for gang scheduling, So that multi-node training jobs only start when all required resources are available, preventing partial job starts and resource waste.

6. **As an Infrastructure Team Member**, I want MPIJobs to support RDMA-capable networking (InfiniBand/RoCE), So that training workloads achieve sub-microsecond latency required for efficient multi-node communication.

7. **As a Data Scientist**, I want MPIJobs to automatically handle SSH key generation and hostfile creation, So that I don't need to manually configure passwordless SSH communication between training pods.

**Storage & Data Access:**

8. **As a Data Scientist**, I want MPIJobs to access shared filesystem storage (NFS, S3-compatible), So that all training nodes can access datasets and write checkpoints to shared locations.

9. **As an MLOps Engineer**, I want to initialize datasets and pre-trained models from S3 and HuggingFace sources before MPIJob execution, So that training workloads have access to required data assets.

**GPU & Hardware Acceleration:**

10. **As a Machine Learning Researcher**, I want MPIJobs to support CUDA-aware MPI and GPUDirect RDMA, So that I can achieve efficient GPU-to-GPU communication across nodes for large-scale model training.

11. **As a DevOps Engineer**, I want MPIJobs to work with NVIDIA GPU Operator and other accelerator frameworks, So that training workloads can leverage available GPU resources in the cluster.

**Enterprise Requirements:**

12. **As a Security Administrator**, I want MPIJobs to respect namespace isolation and RBAC policies, So that multi-tenant clusters maintain proper security boundaries.

13. **As a Platform Administrator**, I want MPIJobs to support priority classes and preemption policies, So that critical training workloads can be scheduled appropriately in resource-constrained environments.

### Non-Functional Requirements

**Security:**

1. **Multi-Tenancy Isolation**
   - MPIJobs must respect namespace boundaries and RBAC policies
   - SSH keys generated for MPI communication must be scoped to individual jobs and automatically cleaned up
   - Service accounts and pod security contexts must follow OpenShift AI security standards
   - Secrets containing dataset/model access credentials must be properly scoped to TrainJob namespaces

2. **Network Security**
   - MPI communication endpoints must only be accessible within the job's pod network
   - RDMA/InfiniBand configurations must not bypass network policies
   - Support for NetworkPolicies to restrict inter-pod communication to authorized training nodes

3. **Image Security**
   - Container images for MPI launcher, workers, and initializers must be scanned and validated
   - Support for disconnected/air-gapped environments with internal image registries
   - Image pull secrets properly propagated to all job components

**Performance:**

1. **Network Performance**
   - Support for RDMA over InfiniBand with sub-2 microsecond latency
   - Support for RDMA over Converged Ethernet (RoCE) as alternative to InfiniBand
   - CUDA-aware MPI for direct GPU-to-GPU communication without CPU involvement
   - GPUDirect RDMA for optimal multi-node GPU communication

2. **Scaling Characteristics**
   - Initial release: Support 2-5 nodes with 8-40 cores per node
   - GPU scaling: Up to 128 GPUs (8 GPUs per node) with near-linear scaling efficiency
   - Resource profiles: 192-4096 GB memory per pod
   - Checkpoint/model size: Support for large model files (>100GB)

3. **Storage Performance**
   - Shared storage must support concurrent read/write from multiple training nodes
   - Checkpoint writes should not block training progress
   - Dataset pre-loading to reduce training iteration time

**Scalability:**

1. **Horizontal Scalability**
   - Controller must handle multiple concurrent MPIJobs across namespaces
   - Gang scheduling integration with Kueue for efficient resource allocation
   - PodGroup creation for all-or-nothing scheduling semantics

2. **Resource Management**
   - Proper resource requests/limits propagation from TrainJob to underlying pods
   - Support for node affinity and tolerations for specialized hardware (GPU nodes, RDMA-enabled nodes)
   - Integration with cluster autoscaler for dynamic node provisioning

**Upgrade Considerations:**

1. **Version Compatibility**
   - MPIJob integration depends on KubeFlow Trainer V2 being available in OpenShift AI
   - ClusterTrainingRuntimes are immutable; updates require new runtime versions
   - In-flight MPIJobs should not be disrupted by OpenShift AI upgrades
   - Migration path from self-managed upstream MPI Operator to integrated solution

2. **API Stability**
   - KubeFlow Trainer V2 APIs are currently v2alpha1 (alpha status)
   - OpenShift AI must shield users from upstream API changes during alpha period
   - Clear deprecation policy when APIs graduate from alpha to beta to GA

**High Availability:**

1. **Controller Resilience**
   - Kubeflow Trainer controller must support leader election for HA deployments
   - Controller restarts should not disrupt running MPIJobs
   - Reconciliation loop must handle partial job state during recovery

2. **Job Resilience**
   - Support for MPI job restart policies (maxRestarts configuration)
   - Checkpoint/resume capabilities through shared storage
   - PodFailurePolicy integration for intelligent failure handling

**Observability:**

1. **Metrics**
   - Expose MPIJob-specific metrics (jobs created, succeeded, failed)
   - Integration with Prometheus for metrics collection
   - Resource utilization metrics (GPU, CPU, memory, network)

2. **Logging**
   - Unified log aggregation from launcher and worker pods
   - Structured logging for MPI-specific events (hostfile creation, SSH setup)
   - Integration with OpenShift logging infrastructure

3. **Monitoring**
   - Dashboard UI must display MPIJob status and conditions
   - Visibility into individual pod states (launcher, workers, initializers)
   - Progress tracking for dataset/model initialization phases

**Disconnected Environment Support:**

1. **Air-Gapped Installation**
   - All container images must be available in disconnected registries
   - Support for internal image mirrors
   - No external network dependencies for job execution (after initial setup)

2. **Internal Artifact Repositories**
   - Support for internal S3-compatible storage for datasets/models
   - Support for internal container registries
   - HuggingFace model mirror support for disconnected environments

### Out-of-Scope

**For Initial Release:**

1. **MPI Operator V1 Support**
   - Only MPI Operator V2 functionality through Trainer V2 will be supported
   - No backward compatibility with standalone MPI Operator V1 deployments

2. **Advanced MPI Implementations**
   - Initial release: OpenMPI only
   - Out of scope: Intel MPI, MVAPICH2, MPICH (may be added in future releases)

3. **Framework-Specific Optimizations**
   - Initial focus on generic MPI workloads and basic HuggingFace integration
   - Out of scope: Framework-specific runtimes for TensorFlow, JAX, XGBoost, PaddlePaddle
   - Out of scope: Specialized integrations with Megatron-LM, DeepSpeed (except basic MPI orchestration)

4. **Elastic Training**
   - MPIJobs with dynamic node scaling (elastic training) not supported in initial release
   - Fixed node count required at job submission time
   - Future consideration: Integration with PyTorch Elastic-style auto-scaling

5. **Advanced Scheduling Features**
   - Job preemption and resume with checkpoint recovery (manual only in v1)
   - Multi-cluster training across federated OpenShift AI instances
   - Spot/preemptible instance support with automatic migration

6. **UI Features**
   - Visual job topology/network diagram in Dashboard
   - Interactive MPI configuration wizards
   - Real-time performance profiling in Dashboard
   - Resource recommendation engine for optimal MPI configuration

7. **Storage Options**
   - POSIX-compliant distributed filesystems beyond NFS (e.g., Lustre, GPFS, BeeGFS)
   - Automatic PVC provisioning and lifecycle management
   - Multi-tiered storage (hot/cold data management)

8. **Network Configurations**
   - Specialized network plugins beyond standard CNI
   - Custom network topology awareness and optimization
   - Automatic network configuration detection and tuning

9. **Legacy Support**
   - Migration tooling from MPI Operator V1 to Trainer V2
   - Automatic conversion of existing MPIJob CRs to TrainJob format
   - Side-by-side operation with standalone MPI Operator

### Acceptance Criteria

**Job Lifecycle Management:**

**AC-1: MPIJob Creation and Submission**
- **Given** a user has access to OpenShift AI with MPIJob support enabled
- **When** they submit an MPIJob specification via Dashboard UI, CLI, or SDK
- **Then** the system creates the MPIJob resource, validates the specification, and transitions to Created state
- **And** the user receives confirmation with job ID and submission timestamp

**AC-2: Gang Scheduling Coordination**
- **Given** an MPIJob with 4 worker replicas requiring 2 GPUs each
- **When** the scheduler evaluates resource availability
- **Then** all 4 worker pods are scheduled atomically (all or nothing)
- **And** the launcher pod is created only after all workers are in Ready state (if LauncherCreationPolicy=WaitForWorkersReady)
- **And** the job fails gracefully if resources cannot be allocated within scheduleTimeoutSeconds

**AC-3: MPI Runtime Execution**
- **Given** an MPIJob with launcher and worker pods in Running state
- **When** the launcher initiates mpirun command
- **Then** MPI processes establish communication via SSH on all worker pods
- **And** collective operations (allreduce, broadcast) complete successfully across all ranks
- **And** training logs are accessible from the launcher pod

**AC-4: Job Completion and Cleanup**
- **Given** an MPIJob completes successfully (launcher exits with code 0)
- **When** the completion condition is detected
- **Then** the job status transitions to Succeeded with completion timestamp
- **And** pods are cleaned up according to CleanPodPolicy (Running, All, or None)
- **And** the job is removed after TTLSecondsAfterFinished (if configured)

**Multi-Interface Integration:**

**AC-5: Dashboard UI Submission**
- **Given** a user accesses the OpenShift AI Dashboard
- **When** they navigate to Distributed Training > Create Job > MPIJob
- **Then** a form allows specifying: launcher/worker specs, MPI implementation, slots per worker, scheduling policy
- **And** the form validates resource requests and provides helpful error messages
- **And** successful submission redirects to job detail page

**AC-6: CLI Submission**
- **Given** a user has `oai-cli` installed and configured
- **When** they run `oai-cli training create mpijob -f mpijob.yaml`
- **Then** the CLI validates the YAML specification locally
- **And** submits the job to the cluster and returns job ID
- **And** provides options to follow logs (`--follow`) or wait for completion (`--wait`)

**AC-7: Python SDK Integration**
- **Given** a user has OpenShift AI Python SDK installed
- **When** they programmatically create an MPIJob using SDK objects
- **Then** the SDK provides Pythonic interfaces for MPIJobSpec, ReplicaSpec, and RunPolicy
- **And** the SDK handles authentication and cluster communication transparently
- **And** the SDK supports both synchronous and asynchronous job submission

**Observability and Monitoring:**

**AC-8: Job Status Visibility**
- **Given** an MPIJob is running in the cluster
- **When** a user views the job in Dashboard, CLI, or SDK
- **Then** the status displays: Created, Running, Restarting, Succeeded, Failed, or Suspended
- **And** replica statuses show active/succeeded/failed counts for launcher and workers
- **And** timestamps are shown for startTime, completionTime, lastReconcileTime

**AC-9: Log Aggregation**
- **Given** an MPIJob with 1 launcher and 4 worker pods
- **When** a user requests logs via Dashboard or CLI
- **Then** logs from all pods are accessible with pod name identification
- **And** the launcher logs are highlighted as primary output
- **And** logs support filtering by pod, time range, and log level

**AC-10: Metrics and Events**
- **Given** the MPI Operator exposes Prometheus metrics
- **When** an MPIJob lifecycle event occurs (created, running, succeeded, failed)
- **Then** metrics increment appropriate counters (mpi_operator_jobs_created_total, etc.)
- **And** Kubernetes Events are generated for significant state transitions
- **And** Dashboard displays job-level metrics (runtime duration, resource utilization)

**Resource Management and Scheduling:**

**AC-11: GPU Resource Allocation**
- **Given** an MPIJob requests 8 GPUs per worker across 4 workers
- **When** the job is scheduled
- **Then** each worker pod receives exactly the requested GPU allocation
- **And** GPU device visibility is verified via nvidia-smi in worker pods
- **And** GPU metrics (utilization, memory) are collected and exposed

**AC-12: Gang Scheduling with Kueue**
- **Given** Kueue is configured for the namespace with resource quotas
- **When** an MPIJob specifies schedulingPolicy with queue name
- **Then** Kueue creates a PodGroup with minAvailable set to total replicas
- **And** the PodGroup is admitted only when full resource allocation is available
- **And** lower-priority jobs are preempted if necessary (based on priorityClass)

**AC-13: Network Performance**
- **Given** an MPIJob configured for RDMA over InfiniBand
- **When** MPI collective operations execute
- **Then** message latency is measured and verified < 2 microseconds for small messages
- **And** bandwidth utilization approaches theoretical maximum for the network fabric
- **And** network errors or retransmissions are logged and surfaced

**Error Handling and Recovery:**

**AC-14: Worker Pod Failure Handling**
- **Given** an MPIJob with RestartPolicy=OnFailure
- **When** a worker pod fails (exit code != 0)
- **Then** the pod is restarted according to restartPolicy
- **And** the failure is recorded in replicaStatuses.failed count
- **And** the job transitions to Failed if BackoffLimit is exceeded

**AC-15: Launcher Failure Recovery**
- **Given** an MPIJob where the launcher pod fails
- **When** the failure is detected
- **Then** the job status transitions to Failed with reason and message
- **And** worker pods are cleaned up according to CleanPodPolicy
- **And** Events describe the failure reason (e.g., "LauncherPodFailed: Exit code 137")

**AC-16: Resource Scheduling Timeout**
- **Given** an MPIJob with scheduleTimeoutSeconds=300
- **When** resources are unavailable for 300 seconds
- **Then** the job transitions to Failed with reason "SchedulingTimeout"
- **And** no pods are created or remain in Pending state
- **And** the user is notified via Events and status conditions

**Integration with Trainer V2:**

**AC-17: Trainer V2 API Compatibility**
- **Given** Trainer V2 provides unified APIs for job management
- **When** an MPIJob is created through Trainer V2 interfaces
- **Then** the job appears in the Trainer V2 job list alongside PyTorchJobs and TFJobs
- **And** common operations (list, get, delete, logs) work consistently across job types
- **And** RBAC policies apply uniformly to all Trainer V2 job types

**AC-18: Unified Observability**
- **Given** the OpenShift AI Dashboard displays Trainer V2 jobs
- **When** a user views the Distributed Training page
- **Then** MPIJobs are listed with status, start time, and resource allocation
- **And** clicking an MPIJob shows the same detail view layout as other job types
- **And** log viewing, metrics, and event streams use consistent UI patterns

**Security and Compliance:**

**AC-19: SSH Key Management**
- **Given** MPI requires SSH for inter-pod communication
- **When** an MPIJob is created
- **Then** the operator generates a unique SSH key pair for the job
- **And** the private key is stored in a Secret with restricted permissions
- **And** the public key is mounted to worker pods at sshAuthMountPath
- **And** keys are rotated or removed when the job completes

**AC-20: Pod Security Standards**
- **Given** OpenShift enforces restricted Pod Security Standards
- **When** an MPIJob is submitted
- **Then** launcher and worker pods comply with namespace security policies
- **And** containers run as non-root user (unless explicitly overridden)
- **And** capabilities are dropped to minimal required set
- **And** security context violations are rejected with clear error messages

### Risks & Assumptions

**Technical Risks:**

1. **Upstream API Stability (HIGH RISK)**
   - **Risk**: KubeFlow Trainer V2 is currently alpha (v2alpha1); APIs may change significantly
   - **Impact**: Breaking changes could require rework of OpenShift AI integration
   - **Mitigation**:
     - Maintain abstraction layer between OpenShift AI and upstream APIs
     - Active participation in KubeFlow community to influence API design
     - Prepare for API versioning strategy when moving to beta/GA
   - **Assumption**: Upstream will provide reasonable deprecation windows for alpha → beta → GA transitions

2. **JobSet Dependency (MEDIUM RISK)**
   - **Risk**: Trainer V2 depends on Kubernetes SIG JobSet which is also in active development
   - **Impact**: JobSet API changes could cascade to Trainer V2 and subsequently OpenShift AI
   - **Mitigation**: Monitor JobSet project roadmap; ensure OpenShift includes compatible JobSet versions
   - **Assumption**: JobSet will reach beta/GA status before OpenShift AI MPIJob feature is GA

3. **Complex Infrastructure Prerequisites (HIGH RISK)**
   - **Risk**: MPIJobs require specialized infrastructure (RDMA, InfiniBand, GPU configurations) that may not be available in all customer environments
   - **Impact**: Feature may have limited adoption if infrastructure setup is too complex
   - **Mitigation**:
     - Provide clear infrastructure prerequisite documentation
     - Support degraded mode without RDMA (TCP-based MPI)
     - Offer validated reference architectures for common HPC configurations
   - **Assumption**: Target customers have or can acquire HPC-grade infrastructure for production use

4. **SSH Key Management Security (MEDIUM RISK)**
   - **Risk**: MPI requires SSH communication between pods; key management introduces security considerations
   - **Impact**: Security teams may block adoption due to SSH-based communication model
   - **Mitigation**:
     - Automatic, ephemeral SSH key generation per job
     - Keys scoped to job lifetime and namespace
     - Investigate gRPC-based alternatives in future (upstream discussion)
   - **Assumption**: Security review will accept SSH-based approach with proper key lifecycle management

5. **Gang Scheduling Complexity (MEDIUM RISK)**
   - **Risk**: Gang scheduling through Kueue/Coscheduling adds complexity and potential failure modes
   - **Impact**: Jobs may fail to schedule or experience unexpected delays
   - **Mitigation**:
     - Clear documentation on gang scheduling behavior
     - Monitoring and alerting for scheduling failures
     - Fallback to non-gang scheduling for development/testing scenarios
   - **Assumption**: Kueue integration is mature enough for production use

6. **Multi-Tenancy Isolation (MEDIUM RISK)**
   - **Risk**: MPI jobs with privileged networking (RDMA) may conflict with multi-tenant security requirements
   - **Impact**: Feature may not be usable in multi-tenant shared clusters
   - **Mitigation**:
     - Document single-tenant cluster as preferred deployment for HPC workloads
     - Investigate namespace-scoped RDMA configurations
     - Support dedicated node pools for MPI workloads
   - **Assumption**: Initial customers will use dedicated clusters or node pools for HPC training

**Assumptions:**

1. **Prerequisite Availability**
   - **Assumption**: KubeFlow Trainer V2 support will be implemented and GA in OpenShift AI before MPIJob integration
   - **Validation Needed**: Confirm Trainer V2 timeline aligns with MPIJob feature schedule

2. **Customer Infrastructure**
   - **Assumption**: Target customers have or can provision HPC-grade infrastructure (RDMA networks, high-end GPUs, shared storage)
   - **Validation Needed**: Customer research to confirm infrastructure availability

3. **Operator Compatibility**
   - **Assumption**: NVIDIA GPU Operator and other accelerator operators are compatible with MPI communication patterns
   - **Validation Needed**: Compatibility testing with various GPU operators

4. **Storage Performance**
   - **Assumption**: Available shared storage solutions (NFS, OpenShift Data Foundation) can handle concurrent I/O from multi-node training
   - **Validation Needed**: Performance benchmarks with representative workloads

5. **Network Plugin Compatibility**
   - **Assumption**: OpenShift CNI plugins (OVN-Kubernetes, Multus) can be configured to support RDMA/InfiniBand
   - **Validation Needed**: Network configuration testing and documentation

6. **License Compliance**
   - **Assumption**: OpenMPI and other MPI implementations have compatible licenses (Apache 2.0 preferred)
   - **Validation Needed**: Legal review of all MPI implementation licenses

7. **Upstream Maintenance**
   - **Assumption**: KubeFlow community will continue active maintenance of Trainer V2 and MPI integration
   - **Validation Needed**: Review upstream community health and contributor activity

8. **Resource Availability**
   - **Assumption**: Customers have sufficient cluster resources to run multi-node jobs (typically 40-128 GPUs per job)
   - **Validation Needed**: Sizing guidance and capacity planning documentation

### Supporting Documentation

**Architecture and Design Documents:**
1. MPIJob Integration Architecture (component diagrams, sequence flows, data models)
2. Gang Scheduling Design (Kueue integration, PodGroup patterns, priority policies)
3. Networking and Storage Requirements (RDMA/InfiniBand configs, shared filesystem patterns)

**User Documentation:**
1. Installation and Configuration Guide (enabling MPIJob support, cluster prerequisites)
2. User Guides per Persona (Data Scientist, MLOps Engineer, Administrator)
3. Tutorial and Examples (Hello World MPI, TensorFlow/Horovod, PyTorch MPI backend)

**API and Reference Documentation:**
1. MPIJob API Reference (complete field documentation, validation rules)
2. CLI Reference (`oai-cli training` commands, troubleshooting workflows)
3. Python SDK Documentation (object model, code examples, async patterns)

**Operational Documentation:**
1. Troubleshooting Guide (failure modes, diagnostics, log analysis patterns)
2. Performance Tuning Guide (MPI library selection, network optimization, GPU placement)
3. Migration and Upgrade Guide (from standalone MPI Operator, API versioning)

**Design Artifacts:**
- Wireframes for MPIJob creation form in Dashboard UI
- Job detail page mockups showing launcher/worker status
- Log viewer designs for multi-pod log aggregation
- Error message templates for common validation failures
- Test plans (unit, integration, e2e, performance)

### Additional Clarifying Information

**Engineering Implementation Considerations:**

1. **Dependency on Trainer V2**
   - This feature **requires** KubeFlow Trainer V2 support to be implemented in OpenShift AI first
   - MPIJob integration will leverage Trainer V2 unified APIs, observability infrastructure, RBAC, and Dashboard patterns

2. **MPI Implementation Support Matrix**

   | MPI Implementation | Use Case | Priority |
   |-------------------|----------|----------|
   | OpenMPI (default) | General-purpose, most widely used | P0 - Must Have |
   | Intel MPI | Intel CPU/GPU optimized workloads | P1 - Should Have |
   | MPICH | HPC/scientific computing applications | P1 - Should Have |

3. **Container Image Requirements**
   - Users must provide images with MPI runtime libraries, SSH daemon, CUDA libraries (if GPU), training framework
   - Documentation should provide base image Dockerfiles and multi-stage build examples

4. **Networking Architecture**
   - **Standard Networking (Minimum Viable)**: Pod-to-pod over cluster network, adequate for 2-4 nodes
   - **High-Performance Networking (Recommended for Production)**: RDMA over InfiniBand/RoCE with Multus CNI and SR-IOV

5. **Storage Patterns**
   - **Required**: ReadWriteMany PVCs for shared datasets/checkpoints
   - **Common Configurations**: NFS, MinIO/Ceph, Lustre/GPFS, cloud provider storage

6. **Gang Scheduling Considerations**
   - MPIJob controller creates PodGroup resources when Kueue is available
   - PodGroup.spec.minMember = launcher replicas + worker replicas
   - Fallback to standard Kubernetes scheduling if Kueue not available (with risk of partial pod sets)

7. **Observability Integration Points**
   - MPI Operator exposes Prometheus metrics at /metrics endpoint
   - Dashboard displays total jobs, success rate, average duration, resource utilization
   - Logging integration with OpenShift EFK/Loki stack
   - Event streams displayed in Dashboard job detail page

8. **Multi-Tenancy and Resource Isolation**
   - MPIJobs are namespaced resources with RBAC-controlled access
   - SSH keys isolated per MPIJob
   - Network policies can restrict pod-to-pod communication to same MPIJob
   - PodSecurityStandards enforce container runtime security

9. **Testing and Validation Requirements**
   - Functional testing for each MPI implementation
   - Performance testing (MPI benchmarks, scaling efficiency, GPU collectives)
   - Compatibility testing (Trainer V2 APIs, Kueue, upgrade paths)

10. **Known Limitations and Future Enhancements**
    - **Initial Release Limitations**: No elastic training, no built-in fault tolerance, no automated network config
    - **Future Enhancements**: MLflow integration, automated checkpoint management, dynamic resource scaling

---

## New Feature / Component Prerequisites & Dependencies

### ODH/RHOAI Build Process Onboarding

**Question**: Will this feature require onboarding of a new container Image or component? **YES**

**Required Components:**

1. **KubeFlow Trainer V2 Controller** (`kubeflow-trainer-controller-manager`)
   - Purpose: Core controller for TrainJob CRD and TrainingRuntime management
   - Upstream: https://github.com/kubeflow/trainer
   - License: Apache 2.0 (to be confirmed)

2. **MPI Launcher Image** (`mpi-launcher`)
   - Purpose: Launcher pod that executes mpirun commands
   - Upstream: https://github.com/kubeflow/mpi-operator (V2 implementation)

3. **MPI Worker Base Image** (`mpi-worker-base`)
   - Purpose: Base image with OpenMPI, SSH server, CUDA-aware MPI libraries
   - Dependencies: OpenMPI with CUDA support, NCCL libraries, SSH server

4. **Dataset Initializer Image** (`dataset-initializer`)
   - Purpose: Initialize and preprocess datasets from S3, HuggingFace

5. **Model Initializer Image** (`model-initializer`)
   - Purpose: Download and prepare pre-trained models from HuggingFace, S3

6. **JobSet Controller** (`jobset-controller`)
   - Purpose: Kubernetes SIG JobSet controller (dependency of Trainer V2)
   - Upstream: https://github.com/kubernetes-sigs/jobset

**Action**: Follow [instructions](?tab=t.658ekv1v7j02) to add components to ODH/RHOAI build process (2 week lead time). Contact [Jay Koehler](mailto:jkoehler@redhat.com) with questions.

### License Validation

**Question**: Will this feature require bringing in new upstream projects or sub-projects into the product? **YES**

**Components Requiring License Review:**
1. OpenMPI - Check license (historically BSD-like, confirm version)
2. KubeFlow Trainer - Apache 2.0 (confirm)
3. JobSet - Apache 2.0 (confirm)
4. NCCL - NVIDIA license (may require alternative or customer-provided)
5. HuggingFace libraries - Apache 2.0 (confirm)

**Preferred License**: Apache 2.0

**Action**: Post to `forum-openshift-ai-architecture` with specific licenses if not Apache 2.0, explaining why required.

### Accelerator/Package Support

**Question**: Does this feature require support from the AIPCC team? **YES**

**Package Requests** (per [AIPCC-1](https://issues.redhat.com/browse/AIPCC-1)):

1. **MPI Libraries**
   - OpenMPI with CUDA support
   - CUDA-aware MPI bindings
   - NCCL library integration

2. **Networking Libraries**
   - RDMA libraries (libibverbs, librdmacm)
   - InfiniBand drivers and utilities
   - GPUDirect RDMA support

3. **GPU Packages**
   - CUDA toolkit (version compatible with MPI builds)
   - cuDNN for deep learning workloads
   - NCCL for multi-GPU communication

4. **Python Packages** (for initializer images)
   - HuggingFace transformers
   - HuggingFace datasets
   - Boto3 for S3 access

**Accelerator Support**:
- NVIDIA GPUs (H100, H200, A100, V100 - prioritize latest generations)
- Ensure compatibility with NVIDIA GPU Operator
- Future consideration: AMD ROCm, Intel GPUs (out of scope for initial release)

**Action**: Clone and populate AIPCC ticket, attach to RHAISTRAT feature.

### Architecture Review Check

**Question**: Does the feature have the label "requires_architecture_review"? **YES**
**Question**: Does the related RFE indicate "Requires architecture review: YES"? **YES**

**Architecture Review Topics:**

1. **Integration Architecture**: How MPIJob integrates with existing OpenShift AI architecture
2. **Security Model**: SSH key lifecycle, multi-tenancy with privileged networking, RBAC
3. **Scalability & Performance**: Controller scalability, gang scheduling, storage architecture
4. **Failure Handling**: Job retry semantics, partial failures, PodFailurePolicy integration
5. **Upgrade Strategy**: In-flight jobs during upgrades, TrainingRuntime versioning
6. **Observability Architecture**: Metrics, logs, Dashboard integration

**Forum**: OpenShift AI Architecture Forum
**Timing**: Before committing to specific solution implementation
**Action**: Schedule architecture review session after requirements refinement complete.

### Additional Dependencies

**Upstream Dependencies:**

1. **KubeFlow Upstream Community**
   - Dependency: Trainer V2 API stability and feature completeness
   - Impact: Critical blocker; cannot implement without upstream APIs
   - Contact: KubeFlow Training Working Group

2. **Kubernetes SIG Batch**
   - Dependency: JobSet API stability and features
   - Impact: Critical blocker; Trainer V2 depends on JobSet
   - Contact: Kubernetes SIG Batch Working Group

**OpenShift Team Dependencies:**

3. **OpenShift Networking Team**
   - Dependency: RDMA/InfiniBand support in OpenShift CNI
   - Impact: High; limits performance without RDMA support
   - Contact: OpenShift Networking team

4. **NVIDIA GPU Operator Team**
   - Dependency: GPU Operator compatibility with MPI jobs, GPUDirect RDMA support
   - Impact: Medium; affects GPU communication efficiency
   - Contact: OpenShift AI AIPCC team liaison

5. **OpenShift Storage Team (ODF)**
   - Dependency: Shared storage support (NFS, RWX PVCs) for multi-node training
   - Impact: High; multi-node training requires shared storage
   - Contact: OpenShift Data Foundation team

6. **Kueue Project**
   - Dependency: Gang scheduling integration, PodGroup support
   - Impact: Medium; affects scheduling reliability
   - Contact: Kueue maintainers

**Internal OpenShift AI Team Dependencies:**

7. **OpenShift AI Dashboard Team**
   - Dependency: UI components for MPIJob creation, monitoring, management
   - Impact: Medium; affects user experience
   - Contact: Internal - OpenShift AI Dashboard team

8. **OpenShift AI SDK Team**
   - Dependency: Python SDK support for TrainJob creation and management
   - Impact: Medium; affects programmatic job submission
   - Contact: Internal - OpenShift AI SDK team

9. **Documentation Team**
   - Dependency: Comprehensive documentation for complex infrastructure setup
   - Impact: Medium; critical for customer adoption
   - Contact: Internal - OpenShift AI Documentation team

10. **QE Team**
    - Dependency: Testing infrastructure with RDMA-capable hardware
    - Impact: High; cannot validate full functionality without proper test infrastructure
    - Contact: Internal - OpenShift AI QE team

**Infrastructure Prerequisites (Customer-Facing):**

These must be documented for customers:

1. **Networking Infrastructure**: InfiniBand fabric or RoCE-capable Ethernet, RDMA-capable network adapters, CNI plugin configuration
2. **Compute Infrastructure**: GPU-equipped compute nodes (8 GPUs per node), GPUDirect RDMA support, adequate CPU/memory
3. **Storage Infrastructure**: Shared POSIX filesystem (NFS or equivalent) with RWX access, high-performance storage for checkpoints
4. **OpenShift Configuration**: Compatible OpenShift version, Kueue installed, NVIDIA GPU Operator, sufficient cluster quota

**Additional Requirements:**

- **Security & Compliance Review**: Security review of SSH-based MPI communication (before release)
- **Performance Benchmarking**: Validate performance claims during QE phase
- **Reference Architecture Documentation**: Document validated hardware configurations

---

## High Level Plan

> **Note**: This section will be populated during planning phase after requirements approval and architecture review.

### Team Delivery Plan

| Team(s) Involved in Delivery | Start Date | Work to Deliver (EPIC) | Team Dependencies | T-Shirt Size Estimate | Approval/Comments |
|------------------------------|------------|------------------------|-------------------|---------------------|-------------------|
| [team-ai-core-platform](mailto:team-ai-core-platform@redhat.com) | TBD | Trainer V2 Integration | Upstream KubeFlow, JobSet | TBD | |
| [team-ai-core-platform](mailto:team-ai-core-platform@redhat.com) | TBD | MPIJob Controller Integration | Trainer V2 complete | TBD | |
| Dashboard Team | TBD | Dashboard UI for MPIJob | MPIJob API available | TBD | |
| SDK Team | TBD | Python SDK for MPIJob | MPIJob API available | TBD | |
| CLI Team | TBD | CLI Support for MPIJob | MPIJob API available | TBD | |
| Documentation Team | TBD | User & Admin Documentation | Feature complete | TBD | |
| UXD Team | TBD | UX Design & Validation | Requirements approved | TBD | |
| QE Team | TBD | Test Infrastructure & Validation | Feature complete | TBD | |
| AIPCC Team | TBD | Package Support (MPI, CUDA, NCCL) | Package requests submitted | TBD | |

---

## How to Engage The Documentation and UXD Teams

### Documentation Team

Engage the Documentation team by:

1. Adding the "Documentation" component to the feature
2. Setting the "Product Documentation Required" field to **Yes** in the feature
3. Adding the docs team to the table in the section above
4. Reviewing the [Docs Intake Process](https://docs.google.com/document/d/1G_LKipII0DMX3UxpkxVEpgM9Pk5tHcfZdvnkjn9E1mI/edit?tab=t.0)
5. Making sure to flag any new features and enhancements for Release Notes

**Documentation Scope for MPIJob:**
- Installation and configuration guide (cluster prerequisites, RDMA/GPU setup)
- Multi-persona user guides (Data Scientist, MLOps Engineer, Administrator)
- Tutorial examples (Hello World MPI, LLM fine-tuning, PyTorch distributed)
- API reference documentation
- Troubleshooting and performance tuning guides
- Migration guide from standalone MPI Operator

### UXD Team

Engage the UXD team by:

1. Adding the "UXD" component to the feature
2. Adding the UXD team to the table above
3. Reaching out to [Jenn Giardino](mailto:jgiardin@redhat.com) or [Beau Morley](mailto:bmorley@redhat.com)

**UXD Scope for MPIJob:**
- Dashboard UI design for MPIJob creation form
- Job detail page layouts and status visualization
- Multi-pod log aggregation viewer design
- User research and validation with target personas
- Error message and validation feedback design
- Consistency review with existing Trainer V2 job types
