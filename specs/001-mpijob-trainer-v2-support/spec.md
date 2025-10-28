# Feature Specification: MPIJobs Support in OpenShift AI via Trainer V2

**Feature Branch**: `001-mpijob-trainer-v2-support`
**Created**: 2025-10-28
**Status**: Draft
**Input**: User description: "Support for MPIJob v2 in OpenShift AI using KubeFlow Trainer v2"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Data Scientist Creates First MPIJob (Priority: P1)

A data scientist with existing experience using PyTorchJob in OpenShift AI needs to run a distributed training workload using MPI (Horovod) with the same familiar tools and interfaces.

**Why this priority**: This is the foundational user journey that enables all other scenarios. If users cannot create and run MPIJobs successfully, the feature provides no value. Supports Goal 1 (Unified Training Experience).

**Independent Test**: Can be fully tested by creating an MPIJob using the Python SDK, verifying it runs to completion, and confirming the same workflow works identically to PyTorchJob creation. Delivers immediate value by enabling MPI workloads.

**Acceptance Scenarios**:

1. **Given** a data scientist has access to an OpenShift AI project with Trainer V2 enabled and MPI runtime configured, **When** they use the Python SDK to create an MPIJob specifying worker count, resources per worker, and container image, **Then** the MPIJob is created successfully and appears in the unified Training Jobs dashboard with an MPI framework indicator.

2. **Given** an MPIJob has been submitted, **When** the data scientist checks the job status via the Dashboard, **Then** they see status transitions from "Pending" to "Gang Scheduling" to "Running" with real-time updates and a topology view showing launcher plus worker pod health.

3. **Given** an MPIJob is running, **When** the data scientist accesses logs through the Dashboard, **Then** they can view aggregated logs from launcher and worker pods with filtering and search capabilities, similar to PyTorchJob log access.

4. **Given** an MPIJob completes successfully, **When** the data scientist reviews the job summary, **Then** they see completion status, duration, resource utilization metrics, and can access final model artifacts from configured storage.

---

### User Story 2 - MLOps Engineer Troubleshoots Failed MPIJob (Priority: P1)

An MLOps engineer receives an alert that a production MPIJob has failed and needs to quickly identify the root cause and remediate the issue using unified observability tools.

**Why this priority**: Production failures directly impact business operations. Without effective troubleshooting capabilities, teams lose confidence in the platform and revert to manual workarounds. Supports Goal 2 (Enterprise-Grade Observability).

**Independent Test**: Can be tested by deliberately causing common MPI failure scenarios (network policy blocking, resource constraints, SSH failures) and verifying the engineer can diagnose and resolve each issue using only the Dashboard and standard CLI tools within 15 minutes.

**Acceptance Scenarios**:

1. **Given** an MPIJob has failed, **When** the MLOps engineer accesses the job in the Dashboard, **Then** they see a clear failure reason (e.g., "Launcher pod failed", "Worker-2 crashed"), topology view highlighting which specific pods failed, and access to relevant logs and events.

2. **Given** the engineer is investigating a gang scheduling timeout, **When** they view the MPIJob status, **Then** they see diagnostic information showing "Waiting for 4 worker pods. Only 2 nodes available" with actionable recommendations.

3. **Given** an MPIJob worker pod failed with an OOM error, **When** the engineer examines worker-specific logs, **Then** they can identify the exact error ("CUDA out of memory"), understand which worker failed, and have sufficient context to adjust resource allocations and resubmit the job.

4. **Given** an MPIJob failed due to network policy issues, **When** the engineer reviews launcher logs showing "SSH connection refused", **Then** they can verify network policy configuration via CLI, identify the missing pod-to-pod communication rule, and apply the required policy.

---

### User Story 3 - Administrator Configures MPI Runtime Templates (Priority: P2)

An OpenShift AI administrator needs to define and publish standardized MPI runtime templates that enforce organizational policies, resource constraints, and best practices for all users.

**Why this priority**: Runtime templates enable centralized governance and reduce user configuration errors. While essential for enterprise deployments, the feature can initially function with basic runtime configurations while template capabilities are refined.

**Independent Test**: Can be tested by creating a ClusterTrainingRuntime for Horovod with GPU support, validating it through a test job, and confirming users can reference this runtime when creating MPIJobs. Delivers value by standardizing MPI configurations across the organization.

**Acceptance Scenarios**:

1. **Given** an administrator has cluster-admin permissions, **When** they navigate to Dashboard → Settings → Training Runtimes and create a new MPI runtime specifying Horovod implementation, default slots per worker, container image, and resource constraints, **Then** the runtime is saved as a ClusterTrainingRuntime and becomes available to all users cluster-wide.

2. **Given** a runtime template has been created, **When** the administrator uses the "Test Runtime" feature, **Then** the system creates a minimal test job (1 launcher + 2 workers), validates image pull success, SSH key generation, and MPI process spawning, and reports test results.

3. **Given** a runtime template has been published, **When** data scientists create new MPIJobs, **Then** they can reference the runtime by name (e.g., "mpi-horovod-gpu"), and the system enforces the configured resource constraints, node selectors, and scheduling policies automatically.

4. **Given** a runtime requires updates, **When** the administrator modifies the ClusterTrainingRuntime configuration, **Then** newly created MPIJobs use the updated configuration while existing running jobs continue with their original runtime specifications.

---

### User Story 4 - Multi-Tenant Job Isolation and RBAC (Priority: P1)

Multiple teams in an enterprise organization need to run MPIJobs in the same OpenShift AI cluster with strong isolation guarantees, respecting namespace boundaries, RBAC policies, and resource quotas.

**Why this priority**: Multi-tenancy is non-negotiable for 92% of enterprise customers. Without proper isolation, security incidents could block adoption entirely. Supports Goal 3 (Multi-Tenancy and Security Consistency).

**Independent Test**: Can be tested by creating MPIJobs in separate namespaces with different RBAC policies and resource quotas, verifying that users cannot access other tenants' jobs, and confirming resource quotas are enforced correctly for launcher plus worker pod aggregates.

**Acceptance Scenarios**:

1. **Given** two teams (Team A and Team B) have separate OpenShift AI projects, **When** Team A creates an MPIJob in their namespace, **Then** Team B cannot view, modify, or delete Team A's MPIJob through any interface (CLI, SDK, Dashboard), and all pod-to-pod MPI communication remains within Team A's namespace.

2. **Given** a namespace has a resource quota of 8 GPUs and 64GB memory, **When** a user attempts to create an MPIJob requiring 4 workers with 3 GPUs and 20GB memory each (12 GPUs total), **Then** the system rejects the job with a clear error message indicating resource quota would be exceeded.

3. **Given** an MPIJob is running in one namespace, **When** network traffic flows between launcher and worker pods, **Then** network policies enforce namespace isolation preventing cross-namespace communication while allowing required intra-namespace MPI traffic.

4. **Given** a user has read-only access to an OpenShift AI project, **When** they attempt to create or delete MPIJobs in that namespace, **Then** the system denies the operation with an RBAC error message indicating insufficient permissions.

---

### User Story 5 - Dashboard UI Job Creation and Monitoring (Priority: P2)

A data scientist who prefers UI-based workflows needs to create and monitor MPIJobs entirely through the Dashboard without writing YAML or using the CLI.

**Why this priority**: 58% of new users prefer UI-based job creation (per user research). While SDK and CLI are critical for automation, Dashboard access lowers the barrier to entry and supports adoption. Can be implemented after core API functionality is stable.

**Independent Test**: Can be tested by creating an MPIJob using only the Dashboard creation wizard, monitoring its progress through the UI, accessing logs via the Dashboard, and confirming feature parity with SDK-based workflows.

**Acceptance Scenarios**:

1. **Given** a data scientist accesses the OpenShift AI Dashboard, **When** they navigate to Training Jobs → Create Job and select MPI as the runtime type, **Then** they see a creation wizard with MPI-specific fields (worker count, slots per worker, MPI runtime selection, container image, resource requests).

2. **Given** the user is filling out the MPIJob creation wizard, **When** they hover over configuration fields, **Then** they see contextual help tooltips explaining each parameter (e.g., "Slots Per Worker: Number of MPI processes per worker node. Typically equal to number of GPUs per worker").

3. **Given** an MPIJob has been created via the Dashboard, **When** the user views the job details page, **Then** they see a topology visualization showing launcher and worker pods with individual health status, real-time status updates, and access to logs for each pod.

4. **Given** an MPIJob is displayed in the Dashboard, **When** the user wants to create a similar job with modified parameters, **Then** they can use the "Clone Job" feature to duplicate the configuration and adjust specific parameters before submission.

---

### User Story 6 - CLI Job Management for Automation (Priority: P2)

A DevOps engineer needs to create, monitor, and manage MPIJobs using CLI commands in CI/CD pipelines and automation scripts with the same commands they use for other Trainer V2 jobs.

**Why this priority**: CLI automation is essential for enterprise workflows. However, SDK provides similar capabilities and is often preferred for programmatic use. CLI can be prioritized after SDK functionality is stable.

**Independent Test**: Can be tested by scripting complete MPIJob lifecycle operations (create, list, describe, delete) using only `oc` CLI commands and verifying functionality matches PyTorchJob CLI workflows.

**Acceptance Scenarios**:

1. **Given** a DevOps engineer has `oc` CLI access to an OpenShift AI cluster, **When** they run `oc create -f mpijob.yaml` with a TrainJob manifest specifying MPI runtime, **Then** the MPIJob is created and confirmation message includes job name and namespace.

2. **Given** multiple training jobs exist in a namespace, **When** the engineer runs `oc get trainjob --runtime=mpi`, **Then** they see a filtered list showing only MPIJobs with columns for name, status, runtime, age, and completion time.

3. **Given** an MPIJob is running, **When** the engineer runs `oc describe trainjob <job-name>`, **Then** they see detailed information including status, events, launcher and worker pod status, resource allocations, and any error messages.

4. **Given** an MPIJob needs to be terminated, **When** the engineer runs `oc delete trainjob <job-name>`, **Then** the job and all associated pods (launcher + workers) are deleted gracefully with confirmation message.

---

### User Story 7 - Migration from Legacy MPIJob v2beta1 (Priority: P3)

A data scientist has existing production MPIJobs using the legacy v2beta1 API and needs to migrate to Trainer V2 with clear guidance and minimal disruption.

**Why this priority**: Migration support is important for existing customers but not required for new feature adoption. Documentation and conversion tools can be developed incrementally based on actual migration demand.

**Independent Test**: Can be tested by converting a representative legacy MPIJob YAML to TrainJob format using migration documentation, deploying both versions side-by-side in development, and validating identical training outcomes.

**Acceptance Scenarios**:

1. **Given** a data scientist has a legacy MPIJob v2beta1 YAML specification, **When** they access the migration documentation, **Then** they find a side-by-side API comparison table, field mapping guide, and example conversions for common scenarios.

2. **Given** the user is converting a legacy MPIJob specification, **When** they apply the documented transformation rules, **Then** they can create an equivalent TrainJob manifest that runs with the same runtime behavior (launcher resources, worker count, slot configuration).

3. **Given** a converted MPIJob is submitted to the cluster, **When** the job runs to completion, **Then** the training outcomes (model accuracy, training time, resource utilization) match the legacy MPIJob implementation within acceptable variance.

4. **Given** the organization is planning a full migration, **When** administrators review the deprecation timeline documentation, **Then** they understand the legacy API support window (2 release cycles ≈ 12 months), migration milestones, and access to migration support resources.

---

### Edge Cases

- **Gang Scheduling Timeout**: What happens when an MPIJob requests 8 worker pods but only 4 nodes are available in the cluster? System should provide clear error message indicating insufficient resources, show partial scheduling status in Dashboard, and allow user to adjust worker count without losing job configuration.

- **Partial Worker Failure During Training**: How does the system handle when 1 worker out of 4 crashes mid-training due to OOM or node failure? MPI training typically requires all workers to be healthy, so job should fail with clear indication of which worker failed and why. User should be able to clone and resubmit with adjusted resources.

- **Network Policy Conflicts**: What happens when default namespace network policies block SSH communication between launcher and worker pods? Job creation should succeed but launcher will fail to initialize MPI with connection errors. Dashboard should show diagnostic message indicating network policy issue with link to troubleshooting documentation.

- **Image Pull Failures**: How does the system handle when the specified container image cannot be pulled (authentication, missing tag, registry unavailable)? Launcher and worker pods will fail with ImagePullBackOff status. Dashboard should display clear error message with image name and registry details to aid troubleshooting.

- **Istio Service Mesh Conflicts**: What happens when Istio sidecar injection interferes with MPI SSH communication? MPI initialization may hang or fail. Runtime templates should include `sidecar.istio.io/inject: "false"` annotation as best practice. Documentation should explain this common pitfall.

- **Resource Quota Boundary**: How does the system handle when an MPIJob would exactly consume remaining namespace quota? System should allow job creation if aggregate launcher + worker resources fit within quota. If quota changes after job creation but before scheduling, provide clear error and guidance.

- **Concurrent MPIJob Submissions**: What happens when multiple users submit MPIJobs simultaneously in a resource-constrained cluster? Gang scheduling should queue jobs fairly. Dashboard should show queue position and estimated wait time based on cluster capacity and running jobs.

- **Launcher Pod Termination**: What happens if the launcher pod is manually deleted while workers are still running? Workers should become orphaned and eventually cleaned up by garbage collection. System should mark job as failed with diagnostic message indicating launcher was terminated unexpectedly.

- **SSH Key Generation Failure**: How does the system handle when SSH key generation for MPI communication fails (insufficient entropy, storage issues)? Launcher pod fails to start with clear error message. Dashboard should indicate SSH key generation failure as distinct failure mode with troubleshooting link.

- **Slot Configuration Mismatch**: What happens when user specifies slots_per_worker=4 but workers only have 2 GPUs available? MPI may initialize but training performance will be suboptimal (CPU fallback) or fail with CUDA device errors. Pre-flight validation should warn about slot/GPU mismatch when possible.

## Requirements *(mandatory)*

### Functional Requirements

**CLI Requirements:**

- **FR-CLI-001**: Users MUST be able to create MPIJobs via `oc create` CLI using Trainer V2 TrainJob CRD with MPI runtime specification in YAML format
- **FR-CLI-002**: Users MUST be able to list all MPIJobs using standard `oc get trainjob` command with optional runtime filter to show only MPI jobs
- **FR-CLI-003**: Users MUST be able to view MPIJob status, events, launcher pod status, worker pod array status, and recent log entries via `oc describe trainjob` command
- **FR-CLI-004**: Users MUST be able to access MPIJob logs from launcher and worker pods via `oc logs` command with pod selector support
- **FR-CLI-005**: Users MUST be able to delete MPIJobs via `oc delete trainjob` command, which triggers graceful termination of launcher and all worker pods

**SDK Requirements:**

- **FR-SDK-001**: Python SDK MUST support creating MPIJobs with MPI-specific parameters including worker count, slots per worker, launcher resource configuration, and runtime reference
- **FR-SDK-002**: Python SDK MUST provide MPIJob status checking methods returning detailed status including job phase, launcher readiness, worker readiness count, and failure reasons
- **FR-SDK-003**: Python SDK MUST support MPIJob event streaming for programmatic monitoring of job lifecycle transitions (submitted, scheduling, running, completed, failed)
- **FR-SDK-004**: Python SDK MUST support MPIJob deletion and cancellation operations with optional graceful termination period
- **FR-SDK-005**: Python SDK MUST provide log retrieval methods for accessing launcher and worker pod logs programmatically with filtering by pod and time range

**Dashboard UI Requirements:**

- **FR-UI-001**: Dashboard MUST display MPIJobs in the unified Trainer V2 job list with clear MPI framework indicator icon and runtime name
- **FR-UI-002**: Dashboard MUST provide MPIJob creation wizard with form fields for MPI-specific configuration (worker count, slots per worker, runtime selection, container image, resource requests, command and arguments)
- **FR-UI-003**: Dashboard MUST show MPIJob status with real-time updates including job phase (Pending, Gang Scheduling, Running, Succeeded, Failed), launcher pod status, and worker pod array aggregate status
- **FR-UI-004**: Dashboard MUST display MPIJob topology view showing launcher pod and all worker pods with individual health status indicators (green for healthy, yellow for pending, red for failed)
- **FR-UI-005**: Dashboard MUST provide access to MPIJob logs with tabbed interface for launcher logs, individual worker logs, and aggregated worker logs, including filtering and search capabilities
- **FR-UI-006**: Dashboard MUST show MPIJob resource utilization metrics including CPU, memory, and GPU usage for launcher and worker pods over time
- **FR-UI-007**: Dashboard MUST display MPIJob events timeline showing lifecycle events (job created, gang scheduling initiated, pods starting, training started, job completed) with timestamps

**Observability Requirements:**

- **FR-OBS-001**: MPIJobs MUST emit metrics compatible with OpenShift AI's existing Trainer V2 monitoring including job duration, success/failure rate, time-to-start, and completion time
- **FR-OBS-002**: MPIJobs MUST emit MPI-specific metrics including gang scheduling duration, MPI initialization time, launcher-to-worker connection time, and slot detection success
- **FR-OBS-003**: MPIJob events MUST be captured and displayed consistently with other Trainer V2 jobs including job submitted, gang scheduling started, pods ready, training started, job completed/failed events
- **FR-OBS-004**: MPIJob logs MUST be aggregated and accessible through OpenShift logging infrastructure with proper labeling for job name, pod type (launcher/worker), and pod index
- **FR-OBS-005**: MPIJob failures MUST generate clear, actionable error messages visible in Dashboard, CLI, and SDK identifying failure reason (gang scheduling timeout, launcher failure, worker failure, resource constraints) and affected components

**Integration Requirements:**

- **FR-INT-001**: MPIJobs MUST respect OpenShift AI's RBAC policies with project-level isolation, preventing users from accessing MPIJobs in namespaces where they lack permissions
- **FR-INT-002**: MPIJobs MUST integrate with OpenShift AI's resource quota management, with launcher plus all worker pod resources counted against namespace quotas
- **FR-INT-003**: MPIJobs MUST support KubeFlow Trainer V2's TrainingRuntime and ClusterTrainingRuntime CRDs for MPI runtime configuration including container images, default parameters, and resource constraints
- **FR-INT-004**: MPIJobs MUST support gang scheduling ensuring launcher and all worker pods are scheduled atomically or not at all, with configurable timeout
- **FR-INT-005**: MPIJobs MUST support pod-to-pod networking requirements for MPI communication while respecting namespace network policies
- **FR-INT-006**: MPIJobs MUST support SSH key generation and distribution for launcher-to-worker communication with ephemeral keys scoped to job lifetime
- **FR-INT-007**: MPIJobs MUST integrate with OpenShift AI's storage systems (S3, OCS, NFS) for training data input and model artifact output with consistent access patterns
- **FR-INT-008**: MPIJobs MUST support Security Context Constraints (SCCs) consistent with OpenShift security policies, with documentation for required SCC capabilities

**Runtime Template Requirements:**

- **FR-RT-001**: Administrators MUST be able to define TrainingRuntime or ClusterTrainingRuntime resources specifying MPI framework variant (Horovod, Intel MPI, OpenMPI), container image, default slots per worker, and default resource allocations
- **FR-RT-002**: Runtime templates MUST support configuration of required annotations (e.g., disable Istio sidecar injection), node selectors, affinity rules, and tolerations
- **FR-RT-003**: Runtime templates MUST support definition of resource constraints including maximum workers per job, maximum GPUs per worker, and required node labels
- **FR-RT-004**: Runtime templates MUST support gang scheduling configuration including scheduling timeout, priority class, and minimum available pods requirement
- **FR-RT-005**: Users MUST be able to reference runtime templates by name when creating MPIJobs, with system validating configuration against template constraints

**Security and Multi-Tenancy Requirements:**

- **FR-SEC-001**: MPIJobs MUST enforce namespace isolation with pod-to-pod communication restricted to same namespace only
- **FR-SEC-002**: MPIJob SSH keys MUST be ephemeral (generated per job) and stored securely in Kubernetes secrets with appropriate RBAC restrictions
- **FR-SEC-003**: MPIJob logs MUST respect log sanitization rules to prevent exposure of credentials, API keys, or sensitive training data
- **FR-SEC-004**: MPIJob lifecycle events (creation, modification, deletion) MUST integrate with OpenShift audit logging with user attribution and timestamps
- **FR-SEC-005**: MPIJobs MUST support network policy requirements for MPI communication with documented standard policies for pod-to-pod SSH access
- **FR-SEC-006**: MPIJobs MUST support optional encrypted communication between launcher and workers for customers requiring data-in-transit encryption

**Migration Requirements:**

- **FR-MIG-001**: Documentation MUST provide side-by-side API comparison between legacy MPIJob v2beta1 and Trainer V2 TrainJob with field mapping guide
- **FR-MIG-002**: Documentation MUST include example conversions for common MPIJob scenarios (basic training, multi-worker GPU training, CPU-based training)
- **FR-MIG-003**: Documentation MUST explain behavioral differences between legacy and Trainer V2 implementations including launcher management and scheduling

### Key Entities

- **MPIJob**: A distributed training job using Message Passing Interface for multi-node coordination. Represented by Trainer V2 TrainJob CRD with MPI runtime reference. Key attributes: worker count, slots per worker, runtime reference, resource requirements per worker, container image, command/arguments.

- **Launcher Pod**: Orchestration pod responsible for initializing MPI environment, generating SSH keys, establishing connections to worker pods, and executing mpirun command. Automatically managed by runtime controller. Key attributes: resource allocations (typically smaller than workers), SSH key secret reference, hostfile configuration.

- **Worker Pod**: Compute pods executing training workload under MPI coordination. Multiple workers form the distributed training cluster. Key attributes: slot count (MPI processes per worker), resource allocations (CPU, memory, GPU), unique index in worker array, SSH server for launcher connections.

- **MPI Runtime**: Template defining MPI framework configuration including framework variant (Horovod, Intel MPI, OpenMPI), container image, default parameters, and resource constraints. Represented by TrainingRuntime (namespace-scoped) or ClusterTrainingRuntime (cluster-scoped). Key attributes: runtime name, MPI implementation, default slots, resource limits, scheduling constraints.

- **Gang Scheduling Group**: Logical grouping of launcher and all worker pods that must be scheduled atomically. Ensures all required pods start together or job fails with timeout. Key attributes: minimum available pods (launcher + worker count), scheduling timeout, priority class.

- **Training Job Status**: State and health information for MPIJob including job phase (Pending, Running, Succeeded, Failed), launcher readiness, worker readiness count, failure reason, start/completion timestamps. Key attributes: current phase, conditions, pod statuses, events.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Data scientists can create and run their first MPIJob within 15 minutes of accessing documentation, matching the current PyTorchJob benchmark for time-to-first-job

- **SC-002**: 60% of teams using distributed training in OpenShift AI adopt MPIJobs within 6 months of general availability

- **SC-003**: MLOps engineers report 30% or greater time savings on monitoring and troubleshooting distributed training jobs compared to using separate MPI operator tooling

- **SC-004**: Users can complete the full MPIJob lifecycle (create, monitor, retrieve results, delete) using any single interface (CLI, SDK, or Dashboard) without switching tools

- **SC-005**: Support tickets related to distributed training job management decrease by 25% within 3 months of MPIJobs GA

- **SC-006**: 95% of MPIJob failures provide actionable error messages that allow users to self-remediate without contacting support

- **SC-007**: MPIJobs demonstrate functional correctness by successfully completing reference training workloads (Horovod PyTorch example, Intel MPI TensorFlow example) with expected training accuracy outcomes

- **SC-008**: Net Promoter Score (NPS) for unified training interface exceeds 40 among users managing multiple distributed training frameworks

- **SC-009**: Multi-tenancy controls prevent any cross-tenant data access or resource violations in MPIJob workloads, validated through security penetration testing

- **SC-010**: Users can monitor MPIJob progress and access logs with no more than 2 clicks from the main Training Jobs dashboard page

- **SC-011**: MPIJob creation using Dashboard wizard has 90% or higher task completion rate for first-time users without requiring external documentation

- **SC-012**: Gang scheduling timeout for MPIJobs provides clear status updates every 30 seconds showing current pod readiness and cluster resource availability
