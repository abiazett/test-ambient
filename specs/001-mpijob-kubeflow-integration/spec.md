# Feature Specification: MPIJob Support for OpenShift AI

**Feature Branch**: `001-mpijob-kubeflow-integration`
**Created**: 2025-10-29
**Status**: Draft
**Input**: User description: "Support for MPIJob in OpenShift AI using KubeFlow Trainer v2"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Data Scientist Single-Job Execution (Priority: P1)

A data scientist needs to train a PyTorch model using Horovod across multiple nodes with GPUs to accelerate model development from days to hours.

**Why this priority**: This is the core value proposition - enabling basic multi-node distributed training is the foundational capability that unblocks all other use cases. Without this, no other advanced features matter.

**Independent Test**: Can be fully tested by submitting a simple 4-worker MPIJob via Python SDK from a Jupyter notebook with Horovod training script, and verifying job completion with trained model artifacts saved to persistent storage. Delivers immediate time-to-train reduction for data scientists.

**Acceptance Scenarios**:

1. **Given** a data scientist has a PyTorch training script and container image with Horovod, **When** they submit an MPIJob with 4 workers and 1 GPU per worker using the Python SDK from Jupyter, **Then** the job successfully launches all 4 worker pods, trains the model, and completes within 30% of the time required for single-node training
2. **Given** an MPIJob is running, **When** the data scientist views the job in ODH Dashboard, **Then** they see real-time job status (Pending/Running/Succeeded/Failed), worker pod states, and can access logs from all workers
3. **Given** an MPIJob has failed, **When** the data scientist views the error details, **Then** they see a clear error message indicating root cause (e.g., "Worker 2 failed: CUDA out of memory") with suggested remediation

---

### User Story 2 - MLOps Engineer Programmatic Workflow Integration (Priority: P2)

An MLOps engineer needs to integrate distributed training into CI/CD pipelines for automated model retraining workflows.

**Why this priority**: Once basic execution works (P1), programmatic access enables automation and production workflows. This unlocks the enterprise value of continuous training and MLOps maturity.

**Independent Test**: Can be tested by creating a CI/CD pipeline script that uses CLI commands to submit an MPIJob with declarative YAML, polls job status programmatically, retrieves logs on failure, and exits with appropriate status codes. Delivers automation capability independent of UI.

**Acceptance Scenarios**:

1. **Given** an MLOps engineer has a training job configuration, **When** they submit an MPIJob using CLI with a YAML file in a CI/CD script, **Then** the command completes with exit code 0 and returns structured JSON output containing job name and status
2. **Given** a job is submitted via CLI, **When** the engineer queries job status using the SDK, **Then** they receive detailed status information including worker states, runtime duration, and failure details if applicable
3. **Given** a job has completed, **When** the engineer retrieves logs programmatically, **Then** they can access synchronized logs from launcher and all worker pods with timestamp alignment

---

### User Story 3 - Administrator Multi-Tenant Resource Management (Priority: P1)

A platform administrator needs to manage shared GPU resources across multiple teams with fair allocation, quotas, and complete observability.

**Why this priority**: Multi-tenancy and resource governance are critical for enterprise deployments. Without this, organizations cannot safely share GPU infrastructure, limiting adoption to single-team deployments.

**Independent Test**: Can be tested by configuring namespace-level resource quotas, having multiple teams submit concurrent MPIJobs, and verifying that quotas are enforced, jobs are fairly scheduled, and administrators can view all jobs across namespaces. Delivers enterprise-grade governance.

**Acceptance Scenarios**:

1. **Given** an administrator sets namespace-level GPU quotas (e.g., max 32 GPUs per team), **When** a user attempts to submit an MPIJob that exceeds quota, **Then** they receive a clear error message before any pods are created, explaining the quota limit
2. **Given** multiple teams are running concurrent MPIJobs, **When** the administrator views the admin dashboard, **Then** they see cluster-wide GPU utilization, all active MPIJobs filtered by namespace/user, and queue status
3. **Given** a user submits an MPIJob in their namespace, **When** users in other namespaces attempt to view or modify the job, **Then** they are denied access based on RBAC policies

---

### User Story 4 - Data Scientist UI-Driven Job Creation (Priority: P2)

A data scientist with basic Kubernetes knowledge needs to create and monitor MPIJobs through an intuitive web interface without writing YAML or code.

**Why this priority**: UI access democratizes distributed training for users without deep technical expertise. This is important for broad adoption but not critical for initial viability (technical users can use SDK/CLI).

**Independent Test**: Can be tested by navigating to ODH Dashboard, using the job creation wizard to configure an MPIJob with guided forms, submitting the job, and monitoring progress through the UI tabs (Overview, Workers, Logs, Metrics). Delivers self-service capability for less technical users.

**Acceptance Scenarios**:

1. **Given** a data scientist navigates to ODH Dashboard, **When** they click "Create MPIJob" and complete the guided form (name, image, command, worker count, resources), **Then** the job is successfully created without requiring YAML knowledge
2. **Given** an MPIJob is running, **When** the user clicks on the job in the list view, **Then** they see a details page with 6 tabs (Overview, Configuration, Workers, Logs, Metrics, Events) with live data
3. **Given** multiple training jobs exist (MPIJobs, TFJobs, PyTorchJobs), **When** the user views the unified training jobs list, **Then** they can filter by job type, status, and date range, with all KubeFlow Trainer job types displayed together

---

### User Story 5 - Data Scientist LLM Fine-Tuning at Scale (Priority: P3)

A data scientist needs to fine-tune a 7B parameter large language model across 8 nodes with 32 total GPUs using DeepSpeed with MPI backend.

**Why this priority**: Large-scale LLM training is a high-value use case that demonstrates platform capability at scale. While important for competitive positioning, it builds on the core P1 capabilities and is not required for initial viability.

**Independent Test**: Can be tested by submitting an MPIJob with 8 workers and 4 GPUs per worker, using DeepSpeed ZeRO optimization, and verifying that training achieves ≥85% scaling efficiency with successful model checkpoint saved. Delivers advanced capability for large-scale workloads.

**Acceptance Scenarios**:

1. **Given** a data scientist has a 7B parameter LLM and DeepSpeed configuration, **When** they submit an MPIJob with 8 workers and 4 GPUs per worker, **Then** all 32 worker pods successfully launch, establish MPI communication mesh, and begin distributed training
2. **Given** the LLM training job is running, **When** measuring training performance, **Then** the job achieves ≥85% scaling efficiency (speedup on 32 GPUs is ≥27× vs. single GPU)
3. **Given** the LLM training completes, **When** the data scientist retrieves the model checkpoint, **Then** the checkpoint is successfully saved to persistent storage (PVC) and accessible for subsequent inference deployment

---

### Edge Cases

- What happens when a worker pod fails mid-training (e.g., OOMKilled or node failure)?
  - The entire MPIJob transitions to "Failed" state within 30 seconds
  - Dashboard UI displays clear failure indication with root cause and resource metrics at time of failure
  - Logs from failed worker are retained and accessible for debugging
  - User can edit configuration (e.g., increase memory) and resubmit

- How does the system handle insufficient GPU resources for requested workers?
  - Job enters "Pending" state with clear message: "Insufficient GPUs: requested 32, available 24"
  - If Kueue is configured, job is queued until resources become available
  - Dashboard shows queue position and estimated wait time
  - User can reduce worker count or delete job to free up queue slot

- What happens when MPI communication fails between workers (e.g., NetworkPolicy blocking)?
  - Worker pods start but MPI initialization fails with timeout
  - Launcher pod logs show SSH connection failures with specific worker IPs
  - Dashboard displays error: "MPI communication failure - check NetworkPolicies"
  - Troubleshooting guide provides remediation steps

- How does the system handle image pull failures or invalid container images?
  - Job transitions to "Failed" before any training begins
  - Error message clearly indicates image pull issue: "Failed to pull image 'invalid-registry.com/model:latest': authentication required"
  - User can update image reference and resubmit
  - No partial worker allocation occurs

- What happens when users submit MPIJobs to namespaces where they lack permissions?
  - API returns immediate 403 Forbidden error with RBAC details
  - CLI displays: "Error: User 'jane' lacks 'create' permission for MPIJobs in namespace 'team-cv'"
  - Dashboard UI prevents job submission with pre-validation
  - No resources are created or scheduled

## Requirements *(mandatory)*

### Functional Requirements

**Core MPIJob Orchestration**:

- **FR-001**: System MUST integrate KubeFlow MPI Operator v2 as a managed component within OpenShift AI operator lifecycle
- **FR-002**: System MUST support creation, monitoring, and deletion of MPIJobs with launcher/worker pod architecture
- **FR-003**: System MUST enable configuration of worker replica count (2-100 workers), resource requests/limits (CPU, memory, GPU), and MPI-specific parameters (slots per worker)
- **FR-004**: System MUST provide automatic injection of MPI hostfile and SSH key management for inter-worker communication
- **FR-005**: System MUST delete ephemeral SSH keys when MPIJob is deleted (via owner references)

**CLI Interface**:

- **FR-006**: System MUST extend OpenShift AI CLI with MPIJob commands: create, get, describe, logs, delete
- **FR-007**: System MUST support declarative job definition via YAML files compatible with KubeFlow Trainer V2 MPIJob schema
- **FR-008**: System MUST provide imperative commands for common use cases (e.g., creating job with minimal flags: name, workers, image, command)
- **FR-009**: System MUST enable real-time log streaming from launcher and worker pods with filtering capabilities
- **FR-010**: CLI MUST provide structured output (JSON/YAML) for automation and exit codes indicating success/failure

**Python SDK**:

- **FR-011**: System MUST provide Python SDK with MPIJob class supporting full lifecycle management (create, get, delete, wait_for_completion)
- **FR-012**: SDK MUST enable Pythonic job definition with type hints and validation
- **FR-013**: SDK MUST support synchronous and asynchronous job submission, status checking, and log retrieval
- **FR-014**: SDK MUST be pre-installed and importable in ODH Jupyter notebook environments

**ODH Dashboard UI**:

- **FR-015**: Dashboard MUST add "MPIJobs" section with list view showing all MPIJobs in user's authorized namespace(s)
- **FR-016**: Dashboard MUST implement job creation wizard with guided form (name, image, command, workers, resources) and advanced YAML editor
- **FR-017**: Dashboard MUST provide job details page with 6 tabs: Overview (status, duration), Configuration (YAML), Workers (pod list with status), Logs (launcher + worker logs), Metrics (resource usage), Events (Kubernetes events)
- **FR-018**: Dashboard MUST display unified job list combining MPIJobs with other KubeFlow Trainer job types (TFJob, PyTorchJob) with filtering by type, status, and date range

**Unified Observability**:

- **FR-019**: System MUST integrate MPIJob metrics (job status, pod status, resource usage) into OpenShift AI monitoring stack (Prometheus/Grafana)
- **FR-020**: Dashboard MUST provide real-time job status updates (Pending → Running → Succeeded/Failed) with automatic refresh every 5 seconds
- **FR-021**: System MUST enable aggregated log viewing across launcher and all worker pods with timestamp synchronization
- **FR-022**: System MUST expose key metrics: job duration, worker pod states, GPU utilization per worker, memory usage

**Security and Multi-Tenancy**:

- **FR-023**: System MUST enforce namespace isolation with RBAC policies (users can only create/view MPIJobs in authorized namespaces)
- **FR-024**: MPIJob pods MUST use ServiceAccount-based authentication with minimal required permissions
- **FR-025**: System MUST integrate with OpenShift NetworkPolicies to restrict inter-pod communication to authorized MPIJob workers within same job
- **FR-026**: System MUST log all MPIJob lifecycle events (create, delete, completion, failure) via OpenShift audit logs with user identity and timestamp

**Resource Management**:

- **FR-027**: System MUST support Kubernetes ResourceQuotas and LimitRanges for MPIJob pods (CPU, memory, GPU limits per namespace)
- **FR-028**: System MUST integrate with Kueue for fair-share scheduling and job queuing when resources are insufficient
- **FR-029**: System MUST enable node affinity, taints, and tolerations for GPU node targeting
- **FR-030**: System MUST provide clear, actionable error messages when resource requests exceed quota or no suitable nodes are available

**Documentation**:

- **FR-031**: System MUST provide quickstart guide with end-to-end example (training a PyTorch model with Horovod on 4 workers)
- **FR-032**: Documentation MUST include CLI command reference, SDK API documentation, and Dashboard UI workflows
- **FR-033**: Documentation MUST include troubleshooting guide for common issues (SSH failures, MPI mismatches, GPU allocation problems, OOM errors)
- **FR-034**: Documentation MUST include reference architectures for typical MPIJob configurations (2-node, 8-node, 16-node setups)

### Key Entities

- **MPIJob**: Represents a distributed MPI training job with launcher and worker pods. Key attributes: name, namespace, worker count, container image, command/args, resource requirements (CPU, memory, GPU per worker), volume mounts, environment variables, status (Pending/Running/Succeeded/Failed), start time, completion time, failure reason.

- **Worker Pod**: Individual compute unit within an MPIJob. Key attributes: pod name, worker index, node assignment, GPU allocation, status, logs, resource utilization metrics. Relationship: each MPIJob has 1 launcher pod and N worker pods.

- **Job Configuration**: Declarative specification of MPIJob requirements. Key attributes: YAML/JSON schema compatible with KubeFlow Trainer V2, includes all FR-003 parameters. Relationship: persisted as Kubernetes custom resource, editable via CLI/SDK/UI.

- **Resource Quota**: Namespace-level limits on GPU/CPU/memory consumption. Key attributes: namespace, max GPUs, max CPU, max memory, current usage. Relationship: enforced before MPIJob worker pods are scheduled.

- **Job Metrics**: Time-series performance and resource data. Key attributes: timestamp, GPU utilization per worker, memory usage, network bandwidth, training throughput. Relationship: collected from worker pods, aggregated and displayed in Dashboard, stored in Prometheus.

## Success Criteria *(mandatory)*

### Measurable Outcomes

**Core Functionality**:

- **SC-001**: A data scientist can submit an MPIJob using SDK from Jupyter notebook with ≤10 lines of code and achieve successful training completion
- **SC-002**: MPIJob with 4 workers and 1 GPU per worker completes training in ≤30% of the time required for equivalent single-node training, demonstrating effective distributed scaling
- **SC-003**: Dashboard page displaying 100 training jobs loads within 2 seconds and updates job status every 5 seconds with real-time data

**Observability**:

- **SC-004**: Users can view synchronized logs from launcher and all worker pods within 5 seconds of log generation, with keyword search functional
- **SC-005**: When a worker pod fails, the entire MPIJob transitions to "Failed" state within 30 seconds with clear error message visible in Dashboard

**Resource Management**:

- **SC-006**: When a user attempts to exceed namespace GPU quota, they receive clear error message before any pods are created (within 2 seconds of submission)
- **SC-007**: Platform supports 50+ concurrent MPIJobs across multiple namespaces without API server degradation or scheduling delays
- **SC-008**: Job submission via CLI or SDK completes within 5 seconds from command execution to "Job Created" confirmation

**Security and Multi-Tenancy**:

- **SC-009**: MPIJob worker pods in namespace A cannot communicate with worker pods in namespace B, enforced by NetworkPolicies and validated via network testing
- **SC-010**: SSH keys used for MPI inter-worker communication are unique per job, ephemeral, and automatically deleted when job is deleted

**Usability**:

- **SC-011**: A new user can run their first successful MPIJob within 30 minutes of reading the quickstart guide (measured via user testing with 10+ participants)
- **SC-012**: Dashboard UI includes in-context help tooltips and documentation links for all configuration fields, with ≥80% "helpful" rating from user feedback

**Performance and Scalability**:

- **SC-013**: LLM fine-tuning job with 8 workers and 32 total GPUs achieves ≥85% scaling efficiency (speedup ≥27× vs. single GPU)
- **SC-014**: Worker pods start within 30 seconds of job submission when resources are available (excluding container image pull time)

**Integration**:

- **SC-015**: MPIJobs can mount existing PVCs for training data and model checkpoints, with successful data persistence verified across job restarts
- **SC-016**: MPIJob metrics are visible in existing Prometheus/Grafana dashboards alongside other workload metrics within 10 seconds of data generation
