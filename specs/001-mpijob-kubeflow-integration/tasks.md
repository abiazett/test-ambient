# Tasks: MPIJob Support for OpenShift AI

**Input**: Design documents from `/specs/001-mpijob-kubeflow-integration/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: Tests are NOT explicitly requested in the feature specification, so test tasks are not included. Implementation tasks include validation and testing as part of the implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story. This feature spans multiple repositories (Training Operator, CLI, SDK, Dashboard) as described in plan.md.

## Format: `- [ ] [ID] [P?] [Story?] Description`
- **[P]**: Can run in parallel (different files/repos, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions
This is a distributed system spanning multiple repositories:
- **training-operator/**: Go operator (upstream KubeFlow)
- **cli/**: Go CLI tool (kubectl plugin)
- **sdk/**: Python SDK package
- **dashboard/**: React/TypeScript UI + Flask backend
- **docs/**: Documentation
- **monitoring/**: Prometheus/Grafana

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Repository initialization and basic project structure

- [ ] T001 Fork KubeFlow Training Operator v2 repository and create RHOAI packaging branch
- [ ] T002 [P] Create CLI repository structure with cmd/mpijob/ and pkg/client/ directories
- [ ] T003 [P] Initialize Python SDK project with openshift_ai_sdk/training/ module structure
- [ ] T004 [P] Create documentation structure in docs/ with quickstart/, reference/, troubleshooting/ directories
- [ ] T005 [P] Setup monitoring manifests in monitoring/prometheus/ and monitoring/grafana/dashboards/

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Training Operator Foundation

- [ ] T006 Implement MPIJob CRD schema in training-operator/manifests/crds/ based on contracts/mpijob-crd.yaml
- [ ] T007 [P] Create MPIJob controller reconciliation loop in training-operator/pkg/controller/mpijob/reconciler.go
- [ ] T008 [P] Implement SSH key generation and Secret/ConfigMap creation in training-operator/pkg/controller/mpijob/ssh.go
- [ ] T009 [P] Implement hostfile generation in training-operator/pkg/controller/mpijob/hostfile.go
- [ ] T010 Implement gang scheduling integration with Kueue in training-operator/pkg/controller/mpijob/podgroup.go
- [ ] T011 [P] Create admission webhook for quota validation in training-operator/pkg/webhook/mpijob_validator.go
- [ ] T012 [P] Implement MPIJob status updates and condition transitions in training-operator/pkg/controller/mpijob/status.go

### Shared Data Models and Utilities

- [ ] T013 [P] Implement ResourceRequirements validation in training-operator/pkg/apis/kubeflow.org/v2beta1/validation.go
- [ ] T014 [P] Create NetworkPolicy template for MPIJob isolation in training-operator/pkg/controller/mpijob/networkpolicy.go
- [ ] T015 [P] Implement Istio sidecar injection disabling logic in training-operator/pkg/controller/mpijob/annotations.go

### Prometheus Metrics Foundation

- [ ] T016 [P] Implement MPIJob metrics exporters in training-operator/pkg/metrics/mpijob_metrics.go
- [ ] T017 [P] Create Prometheus recording rules in monitoring/prometheus/training-operator-rules.yaml

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Data Scientist Single-Job Execution (Priority: P1) 🎯 MVP

**Goal**: Enable data scientists to submit and run a 4-worker MPIJob from Jupyter notebook using Python SDK with Horovod training script, achieving 30% time-to-train reduction vs single-node

**Independent Test**: Submit a simple 4-worker MPIJob via Python SDK from Jupyter notebook with Horovod training script, verify job completion with trained model artifacts saved to persistent storage. Verify real-time job status updates and log viewing from Dashboard.

### Python SDK Implementation for User Story 1

- [ ] T018 [P] [US1] Create MPIJob class in sdk/openshift_ai_sdk/training/mpijob.py with create(), get_status(), get_logs() methods
- [ ] T019 [P] [US1] Implement ResourceRequirements class in sdk/openshift_ai_sdk/training/resources.py
- [ ] T020 [US1] Implement Kubernetes API client wrapper in sdk/openshift_ai_sdk/training/client.py (depends on T018, T019)
- [ ] T021 [US1] Add MPIJob manifest generation in sdk/openshift_ai_sdk/training/mpijob.py (_build_manifest method)
- [ ] T022 [US1] Implement wait_for_completion() method with timeout handling in sdk/openshift_ai_sdk/training/mpijob.py
- [ ] T023 [US1] Add log streaming and aggregation in sdk/openshift_ai_sdk/training/mpijob.py (get_logs with launcher/worker filtering)

### Dashboard UI Implementation for User Story 1

- [ ] T024 [P] [US1] Create MPIJobList component in dashboard/frontend/src/pages/TrainingJobs/MPIJobList.tsx
- [ ] T025 [P] [US1] Create MPIJobDetails component in dashboard/frontend/src/pages/TrainingJobs/MPIJobDetails.tsx with Overview tab
- [ ] T026 [P] [US1] Implement real-time status polling (5s interval) in dashboard/frontend/src/api/training.ts
- [ ] T027 [P] [US1] Create MPIJobLogs component in dashboard/frontend/src/pages/TrainingJobs/MPIJobLogs.tsx with launcher/worker filtering
- [ ] T028 [US1] Implement Dashboard backend REST API for MPIJob listing in dashboard/backend/app/routes/training.py (GET /api/training/mpijobs)
- [ ] T029 [US1] Implement Dashboard backend REST API for MPIJob details in dashboard/backend/app/routes/training.py (GET /api/training/mpijobs/{name})
- [ ] T030 [US1] Implement Dashboard backend REST API for MPIJob logs in dashboard/backend/app/routes/training.py (GET /api/training/mpijobs/{name}/logs)

### Error Handling and User Feedback for User Story 1

- [ ] T031 [US1] Implement clear error messages for common failures in training-operator/pkg/controller/mpijob/error_messages.go (OOMKilled, CUDA errors, resource constraints)
- [ ] T032 [P] [US1] Add error display with remediation suggestions in dashboard/frontend/src/pages/TrainingJobs/MPIJobDetails.tsx
- [ ] T033 [P] [US1] Implement SDK exception handling with actionable messages in sdk/openshift_ai_sdk/training/exceptions.py

**Checkpoint**: At this point, User Story 1 should be fully functional - data scientists can submit MPIJobs from Jupyter, monitor via Dashboard, and view logs

---

## Phase 4: User Story 3 - Administrator Multi-Tenant Resource Management (Priority: P1)

**Goal**: Enable platform administrators to manage shared GPU resources across multiple teams with fair allocation, quotas, and complete observability

**Independent Test**: Configure namespace-level resource quotas, have multiple teams submit concurrent MPIJobs, verify quotas are enforced, jobs are fairly scheduled via Kueue, and administrators can view all jobs across namespaces. Verify RBAC isolation prevents cross-namespace access.

**Note**: US3 is P1 (same as US1) because multi-tenancy and resource governance are critical for enterprise deployments

### RBAC and Multi-Tenancy for User Story 3

- [ ] T034 [P] [US3] Create RBAC manifests in training-operator/manifests/rbac/ with ClusterRole for mpijob-user and mpijob-admin
- [ ] T035 [P] [US3] Implement namespace-scoped RBAC validation in training-operator/pkg/webhook/mpijob_validator.go (verify user can create in namespace)
- [ ] T036 [US3] Implement quota pre-validation in admission webhook in training-operator/pkg/webhook/mpijob_validator.go (check ResourceQuota before admission)

### Kueue Gang Scheduling Integration for User Story 3

- [ ] T037 [US3] Create ClusterQueue CRD examples in training-operator/manifests/kueue/clusterqueue.yaml with GPU flavors
- [ ] T038 [P] [US3] Create LocalQueue CRD examples in training-operator/manifests/kueue/localqueue.yaml for team namespaces
- [ ] T039 [US3] Implement automatic Kueue label injection in training-operator/pkg/controller/mpijob/reconciler.go (kueue.x-k8s.io/queue-name)
- [ ] T040 [US3] Implement PodGroup creation for gang scheduling in training-operator/pkg/controller/mpijob/podgroup.go with minAvailable semantics

### Administrator Dashboard for User Story 3

- [ ] T041 [P] [US3] Create admin view for cluster-wide MPIJob listing in dashboard/frontend/src/pages/TrainingJobs/AdminMPIJobList.tsx
- [ ] T042 [P] [US3] Implement GPU utilization metrics dashboard in dashboard/frontend/src/pages/TrainingJobs/ClusterGPUMetrics.tsx
- [ ] T043 [P] [US3] Add namespace filtering and quota display in dashboard/frontend/src/pages/TrainingJobs/AdminMPIJobList.tsx
- [ ] T044 [US3] Implement backend API for cluster-wide job listing in dashboard/backend/app/routes/training.py (GET /api/training/mpijobs?all-namespaces=true) with RBAC check
- [ ] T045 [US3] Implement backend API for quota status in dashboard/backend/app/routes/training.py (GET /api/training/namespaces/{ns}/quota)

### Quota Error Handling for User Story 3

- [ ] T046 [US3] Implement quota exceeded error messages with suggestions in training-operator/pkg/webhook/mpijob_validator.go (show current usage, available, requested)
- [ ] T047 [P] [US3] Add quota error display in SDK in sdk/openshift_ai_sdk/training/exceptions.py (QuotaExceededError)
- [ ] T048 [P] [US3] Add pre-submission quota validation in Dashboard in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx

**Checkpoint**: At this point, User Story 3 should be fully functional - administrators can manage multi-tenant GPU resources with quotas, fair scheduling, and observability

---

## Phase 5: User Story 2 - MLOps Engineer Programmatic Workflow Integration (Priority: P2)

**Goal**: Enable MLOps engineers to integrate distributed training into CI/CD pipelines using CLI and SDK for automated model retraining workflows

**Independent Test**: Create a CI/CD pipeline script that uses CLI to submit MPIJob from YAML file, polls job status programmatically, retrieves logs on failure, and exits with appropriate status codes. Verify SDK can query job status and retrieve logs.

### CLI Implementation for User Story 2

- [ ] T049 [P] [US2] Implement 'oc odh mpijob create' command in cli/cmd/mpijob/create.go with YAML file support
- [ ] T050 [P] [US2] Implement 'oc odh mpijob get' command in cli/cmd/mpijob/get.go with JSON/YAML output formatting
- [ ] T051 [P] [US2] Implement 'oc odh mpijob describe' command in cli/cmd/mpijob/describe.go with detailed status output
- [ ] T052 [P] [US2] Implement 'oc odh mpijob logs' command in cli/cmd/mpijob/logs.go with --launcher, --worker, --follow flags
- [ ] T053 [P] [US2] Implement 'oc odh mpijob delete' command in cli/cmd/mpijob/delete.go
- [ ] T054 [US2] Implement CLI output printer utilities in cli/pkg/printer/mpijob.go for table/JSON/YAML formats
- [ ] T055 [US2] Implement Kubernetes API client in cli/pkg/client/mpijob_client.go for CLI commands

### SDK Enhancements for User Story 2

- [ ] T056 [P] [US2] Add delete() method to MPIJob class in sdk/openshift_ai_sdk/training/mpijob.py
- [ ] T057 [P] [US2] Implement get_worker_statuses() method in sdk/openshift_ai_sdk/training/mpijob.py returning list of worker pod states
- [ ] T058 [P] [US2] Add structured JSON output for status in sdk/openshift_ai_sdk/training/mpijob.py (to_dict method)
- [ ] T059 [US2] Implement retry logic with exponential backoff in sdk/openshift_ai_sdk/training/client.py for transient API errors

### CI/CD Integration Examples for User Story 2

- [ ] T060 [P] [US2] Create example CI/CD pipeline script in docs/examples/cicd-pipeline.sh demonstrating CLI usage
- [ ] T061 [P] [US2] Create example SDK-based pipeline in docs/examples/sdk-pipeline.py for programmatic job management
- [ ] T062 [US2] Document exit codes and error handling patterns in docs/reference/cli-reference.md

**Checkpoint**: At this point, User Story 2 should be fully functional - MLOps engineers can integrate MPIJobs into CI/CD pipelines using CLI and SDK

---

## Phase 6: User Story 4 - Data Scientist UI-Driven Job Creation (Priority: P2)

**Goal**: Enable data scientists with basic Kubernetes knowledge to create and monitor MPIJobs through an intuitive web interface without writing YAML or code

**Independent Test**: Navigate to ODH Dashboard, use job creation wizard to configure MPIJob with guided forms (no YAML required), submit job, and monitor progress through UI tabs (Overview, Workers, Logs, Metrics). Verify user can create job without Kubernetes/YAML expertise.

### Dashboard Job Creation Wizard for User Story 4

- [ ] T063 [P] [US4] Create MPIJobCreateWizard component in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx with multi-step form
- [ ] T064 [P] [US4] Implement form step 1: Basic Information in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx (name, namespace)
- [ ] T065 [P] [US4] Implement form step 2: Container Configuration in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx (image, command, args)
- [ ] T066 [P] [US4] Implement form step 3: Resources in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx (workers, CPU, memory, GPU)
- [ ] T067 [P] [US4] Implement form step 4: Advanced Options in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx (volumes, node selector, env vars)
- [ ] T068 [P] [US4] Implement form step 5: Review in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx with generated YAML preview
- [ ] T069 [US4] Add form validation and field-level help tooltips in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx
- [ ] T070 [US4] Implement Dashboard backend API for MPIJob creation in dashboard/backend/app/routes/training.py (POST /api/training/mpijobs)

### Enhanced Dashboard Details Views for User Story 4

- [ ] T071 [P] [US4] Add Configuration tab to MPIJobDetails in dashboard/frontend/src/pages/TrainingJobs/MPIJobDetails.tsx showing YAML
- [ ] T072 [P] [US4] Create MPIJobWorkers component in dashboard/frontend/src/pages/TrainingJobs/MPIJobWorkers.tsx with worker pod list and status
- [ ] T073 [P] [US4] Create MPIJobMetrics component in dashboard/frontend/src/pages/TrainingJobs/MPIJobMetrics.tsx with embedded Grafana panels
- [ ] T074 [P] [US4] Create MPIJobEvents component in dashboard/frontend/src/pages/TrainingJobs/MPIJobEvents.tsx showing Kubernetes events
- [ ] T075 [US4] Implement tab navigation in dashboard/frontend/src/pages/TrainingJobs/MPIJobDetails.tsx (Overview, Configuration, Workers, Logs, Metrics, Events)

### Unified Training Jobs List for User Story 4

- [ ] T076 [US4] Extend job list to include TFJob and PyTorchJob in dashboard/frontend/src/pages/TrainingJobs/TrainingJobsList.tsx
- [ ] T077 [P] [US4] Add filtering by job type, status, date range in dashboard/frontend/src/pages/TrainingJobs/TrainingJobsList.tsx
- [ ] T078 [US4] Implement backend API for unified job listing in dashboard/backend/app/routes/training.py (GET /api/training/jobs with type filter)

### In-Context Help and Documentation for User Story 4

- [ ] T079 [P] [US4] Add field-level tooltips for all form fields in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx
- [ ] T080 [P] [US4] Add documentation links to quickstart guide in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx
- [ ] T081 [US4] Implement inline validation messages in dashboard/frontend/src/pages/TrainingJobs/MPIJobCreateWizard.tsx with suggested fixes

**Checkpoint**: At this point, User Story 4 should be fully functional - non-technical data scientists can create and monitor MPIJobs through intuitive UI without YAML/code

---

## Phase 7: User Story 5 - Data Scientist LLM Fine-Tuning at Scale (Priority: P3)

**Goal**: Enable data scientists to fine-tune 7B parameter LLMs across 8 nodes with 32 total GPUs using DeepSpeed with MPI backend, achieving ≥85% scaling efficiency

**Independent Test**: Submit MPIJob with 8 workers and 4 GPUs per worker (32 total GPUs), using DeepSpeed ZeRO optimization, verify training achieves ≥85% scaling efficiency and successful model checkpoint saved. This validates large-scale LLM capability.

### Large-Scale Job Configuration for User Story 5

- [ ] T082 [P] [US5] Create reference architecture YAML for 8-node LLM training in docs/reference/architectures/8-node-llm-training.yaml
- [ ] T083 [P] [US5] Document DeepSpeed configuration in docs/examples/deepspeed-config.json with ZeRO stage 2/3 examples
- [ ] T084 [P] [US5] Create example LLM training script in docs/examples/train-gpt-7b.py with DeepSpeed integration

### Performance Optimization for User Story 5

- [ ] T085 [P] [US5] Implement NCCL environment variable injection in training-operator/pkg/controller/mpijob/reconciler.go (NCCL_DEBUG, NCCL_IB_DISABLE, NCCL_SOCKET_IFNAME)
- [ ] T086 [P] [US5] Add pod anti-affinity for worker spreading in training-operator/pkg/controller/mpijob/reconciler.go (prevent workers on same node)
- [ ] T087 [P] [US5] Implement GPU topology-aware scheduling hints in training-operator/pkg/controller/mpijob/affinity.go

### Scaling Efficiency Metrics for User Story 5

- [ ] T088 [P] [US5] Create Grafana dashboard for scaling efficiency in monitoring/grafana/dashboards/mpijob-scaling-efficiency.json
- [ ] T089 [P] [US5] Implement scaling efficiency calculation in monitoring/prometheus/training-operator-rules.yaml (actual vs expected throughput)
- [ ] T090 [US5] Add scaling efficiency metrics to Dashboard in dashboard/frontend/src/pages/TrainingJobs/MPIJobMetrics.tsx

### SDK Enhancements for Large-Scale Jobs for User Story 5

- [ ] T091 [P] [US5] Add affinity and tolerations configuration to MPIJob class in sdk/openshift_ai_sdk/training/mpijob.py
- [ ] T092 [P] [US5] Implement environment variable configuration in sdk/openshift_ai_sdk/training/mpijob.py
- [ ] T093 [US5] Add volume mount configuration for checkpoint storage in sdk/openshift_ai_sdk/training/mpijob.py

**Checkpoint**: At this point, User Story 5 should be fully functional - data scientists can fine-tune large LLMs at scale with ≥85% scaling efficiency

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories and production readiness

### Documentation

- [ ] T094 [P] Validate and update quickstart guide in docs/quickstart/mpijobs-getting-started.md based on implemented features
- [ ] T095 [P] Create CLI reference documentation in docs/reference/cli-reference.md with all commands and flags
- [ ] T096 [P] Generate SDK API documentation in docs/reference/sdk-api.md from Python docstrings
- [ ] T097 [P] Create MPIJob CRD spec reference in docs/reference/mpijob-crd-spec.md
- [ ] T098 [P] Create comprehensive troubleshooting guide in docs/troubleshooting/mpijobs-common-issues.md covering all pitfalls from research.md
- [ ] T099 Document Istio sidecar pitfall prominently in docs/quickstart/mpijobs-getting-started.md and CLAUDE.md

### Observability

- [ ] T100 [P] Create comprehensive Grafana dashboard in monitoring/grafana/dashboards/mpijob-overview.json with GPU utilization, job status, queue metrics
- [ ] T101 [P] Implement distributed tracing integration in training-operator/pkg/controller/mpijob/reconciler.go using OpenTelemetry
- [ ] T102 [P] Add structured logging throughout operator in training-operator/pkg/controller/mpijob/ with consistent log levels

### Testing and Quality

- [ ] T103 [P] Create integration test suite in training-operator/tests/integration/ covering gang scheduling, SSH communication, status transitions
- [ ] T104 [P] Create E2E test for full job lifecycle in training-operator/tests/e2e/ using 4-worker ResNet-50 example
- [ ] T105 [P] Implement SDK unit tests in sdk/tests/unit/ achieving 80% coverage
- [ ] T106 [P] Create CLI integration tests in cli/tests/e2e/ covering all commands
- [ ] T107 [P] Implement Dashboard component tests in dashboard/frontend/src/pages/TrainingJobs/__tests__/ using Jest and React Testing Library
- [ ] T108 Run quickstart validation to ensure 30-minute time-to-first-job target

### Performance and Benchmarking

- [ ] T109 [P] Create performance benchmarking suite in docs/benchmarks/scaling-efficiency.sh with ResNet-50, BERT, GPT-7B tests
- [ ] T110 [P] Validate 85% scaling efficiency target for 8-node MPIJobs with benchmarks
- [ ] T111 Optimize Dashboard page load time to <2 seconds for 100 jobs

### Security Hardening

- [ ] T112 [P] Implement FIPS 140-2 compliant SSH key generation in training-operator/pkg/controller/mpijob/ssh.go using RSA-4096
- [ ] T113 [P] Add NetworkPolicy automatic creation in training-operator/pkg/controller/mpijob/networkpolicy.go for job isolation
- [ ] T114 [P] Implement audit logging for all MPIJob lifecycle events in training-operator/pkg/controller/mpijob/audit.go
- [ ] T115 Add security scanning for container images in CI/CD pipeline

### Production Readiness

- [ ] T116 [P] Create RHOAI operator integration manifests in training-operator/manifests/rhoai/
- [ ] T117 [P] Implement operator upgrade path from upstream KubeFlow in training-operator/manifests/upgrade/
- [ ] T118 [P] Create disaster recovery documentation in docs/operations/backup-restore.md
- [ ] T119 Implement admission webhook to enforce Istio sidecar annotation in training-operator/pkg/webhook/mpijob_validator.go
- [ ] T120 Update CLAUDE.md with all implementation patterns and common pitfalls

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-7)**: All depend on Foundational phase completion
  - US1 (P1): Data Scientist Single-Job Execution - No dependencies on other stories
  - US3 (P1): Administrator Multi-Tenant Resource Management - No dependencies on other stories (can run parallel with US1)
  - US2 (P2): MLOps Engineer Programmatic Workflow - Depends on US1 for SDK foundation, adds CLI
  - US4 (P2): Data Scientist UI-Driven Job Creation - Depends on US1 for backend APIs, adds UI wizard
  - US5 (P3): Data Scientist LLM Fine-Tuning at Scale - Depends on US1 for core functionality, adds scale optimizations
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: FOUNDATIONAL - SDK + Dashboard viewing - Must complete first
- **User Story 3 (P1)**: Can start after Foundational, parallel with US1 - Independent multi-tenancy features
- **User Story 2 (P2)**: Builds on US1 SDK, adds CLI - Requires US1 SDK complete
- **User Story 4 (P2)**: Builds on US1 Dashboard backend, adds creation wizard - Requires US1 Dashboard APIs complete
- **User Story 5 (P3)**: Builds on US1 core functionality, adds scale optimizations - Requires US1 complete

### Recommended Implementation Order

1. **MVP (Phases 1+2+3)**: Setup → Foundational → User Story 1
   - Delivers: Data scientists can submit and monitor MPIJobs from Jupyter
   - Time estimate: 6-8 weeks

2. **Enterprise (Add US3)**: Multi-tenant resource management
   - Delivers: Production-ready multi-tenancy and quotas
   - Time estimate: +2-3 weeks

3. **MLOps (Add US2)**: CLI for automation
   - Delivers: CI/CD integration capability
   - Time estimate: +2 weeks

4. **Self-Service (Add US4)**: Dashboard creation wizard
   - Delivers: UI-driven job creation for non-technical users
   - Time estimate: +2-3 weeks

5. **Scale (Add US5)**: LLM fine-tuning optimizations
   - Delivers: Large-scale training at ≥85% efficiency
   - Time estimate: +2 weeks

### Parallel Opportunities

**Within Foundational Phase (Phase 2):**
- T007 (controller), T008 (SSH), T009 (hostfile), T011 (webhook), T012 (status), T013 (validation), T014 (NetworkPolicy), T015 (annotations), T016 (metrics), T017 (Prometheus rules) can all run in parallel

**Within User Story 1 (Phase 3):**
- SDK tasks: T018, T019 (MPIJob class, ResourceRequirements) can run in parallel
- Dashboard tasks: T024, T025, T026, T027 (UI components) can run in parallel
- Error handling: T032, T033 (Dashboard error display, SDK exceptions) can run in parallel

**Within User Story 3 (Phase 4):**
- RBAC: T034, T035 (manifests, validation) can run in parallel
- Kueue: T037, T038 (ClusterQueue, LocalQueue) can run in parallel
- Dashboard: T041, T042, T043 (admin views) can run in parallel
- Errors: T047, T048 (SDK quota errors, Dashboard validation) can run in parallel

**Across User Stories After Foundational:**
- US1 and US3 can be developed in parallel (independent teams)
- Once US1 complete: US2, US4, US5 can be started in parallel

---

## Parallel Example: Foundational Phase

```bash
# All these tasks can be launched in parallel after Phase 1:
Task: "Implement MPIJob controller reconciliation loop in training-operator/pkg/controller/mpijob/reconciler.go"
Task: "Implement SSH key generation in training-operator/pkg/controller/mpijob/ssh.go"
Task: "Implement hostfile generation in training-operator/pkg/controller/mpijob/hostfile.go"
Task: "Create admission webhook in training-operator/pkg/webhook/mpijob_validator.go"
Task: "Implement MPIJob status updates in training-operator/pkg/controller/mpijob/status.go"
Task: "Create NetworkPolicy template in training-operator/pkg/controller/mpijob/networkpolicy.go"
Task: "Implement MPIJob metrics exporters in training-operator/pkg/metrics/mpijob_metrics.go"
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 3 Only)

1. Complete Phase 1: Setup (5 tasks, ~1 week)
2. Complete Phase 2: Foundational (12 tasks, ~4 weeks with parallel execution)
3. Complete Phase 3: User Story 1 (16 tasks, ~3 weeks)
4. Complete Phase 4: User Story 3 (15 tasks, ~2 weeks)
5. **STOP and VALIDATE**: Test both stories independently
6. Deploy MVP to staging cluster
7. Validate against success criteria SC-001 through SC-010

**MVP Delivers:**
- ✅ Data scientists can submit MPIJobs from Jupyter (US1)
- ✅ Real-time job monitoring and logs in Dashboard (US1)
- ✅ Multi-tenant resource management with quotas (US3)
- ✅ Gang scheduling prevents resource deadlock (US3)
- ✅ RBAC isolation across namespaces (US3)

**MVP Success Criteria:**
- SC-001: ≤10 lines of code to submit job from SDK
- SC-002: 30% time reduction vs single-node training
- SC-003: Dashboard loads 100 jobs in <2 seconds
- SC-006: Clear quota error before pod creation
- SC-008: Job submission completes in <5 seconds
- SC-009: Cross-namespace isolation enforced

### Incremental Delivery

1. **Phase 1+2 → Foundation Ready** (5 weeks)
   - Training Operator functional with MPIJob CRD
   - Kueue gang scheduling integrated
   - SSH key management and hostfile generation working

2. **Add US1 → Data Scientist MVP** (3 weeks, 8 weeks total)
   - SDK operational from Jupyter
   - Dashboard viewing functional
   - Logs accessible
   - **Deploy and demo**

3. **Add US3 → Enterprise MVP** (2 weeks, 10 weeks total)
   - Multi-tenancy with RBAC
   - Resource quotas enforced
   - Admin dashboard operational
   - **Deploy and demo**

4. **Add US2 → MLOps Ready** (2 weeks, 12 weeks total)
   - CLI commands functional
   - CI/CD integration examples working
   - **Deploy and demo**

5. **Add US4 → Self-Service UI** (3 weeks, 15 weeks total)
   - Job creation wizard functional
   - Unified training jobs list working
   - **Deploy and demo**

6. **Add US5 → Scale Ready** (2 weeks, 17 weeks total)
   - LLM fine-tuning optimizations
   - 85% scaling efficiency validated
   - **Deploy and demo**

7. **Phase 8 → Production Polish** (2 weeks, 19 weeks total)
   - Documentation complete
   - Security hardening done
   - Benchmarking validated
   - **Production release**

### Parallel Team Strategy

With 3 developers after Foundational phase completes:

**Sprint 1-3 (Weeks 6-8):**
- Developer A: User Story 1 (SDK + Dashboard backend)
- Developer B: User Story 3 (RBAC + Kueue + Admin UI)
- Developer C: Documentation + Monitoring

**Sprint 4-5 (Weeks 9-10):**
- Developer A: User Story 2 (CLI)
- Developer B: User Story 4 (Dashboard wizard)
- Developer C: User Story 5 (scale optimizations)

**Sprint 6 (Week 11):**
- All developers: Testing, benchmarking, documentation polish

---

## Notes

- [P] tasks = different files/repos, no blocking dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently testable
- This is a distributed system - coordinate releases across repos
- Training Operator is upstream contribution - maintain KubeFlow compatibility
- Istio sidecar annotation is CRITICAL - document prominently
- Gang scheduling (Kueue) is mandatory - not optional
- Target 80% test coverage for operator and SDK
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently

## Task Statistics

**Total Tasks**: 120 tasks

**Tasks per User Story:**
- Setup: 5 tasks
- Foundational: 12 tasks (BLOCKS all stories)
- User Story 1 (P1): 16 tasks (SDK + Dashboard viewing)
- User Story 3 (P1): 15 tasks (Multi-tenancy + RBAC)
- User Story 2 (P2): 14 tasks (CLI)
- User Story 4 (P2): 19 tasks (Dashboard creation wizard)
- User Story 5 (P3): 13 tasks (LLM scale optimizations)
- Polish: 26 tasks (cross-cutting)

**Parallel Opportunities:**
- Foundational phase: 11 parallel tasks (T007-T017)
- User Story 1: 10 parallel tasks
- User Story 3: 8 parallel tasks
- User Story 2: 9 parallel tasks
- User Story 4: 13 parallel tasks
- User Story 5: 8 parallel tasks
- Polish: 18 parallel tasks

**MVP Scope (US1 + US3):**
- 48 tasks total (Setup + Foundational + US1 + US3)
- Estimated timeline: 10 weeks with parallel execution
- Delivers core value: distributed training + multi-tenancy

**Independent Test Criteria:**
- ✅ US1: Submit 4-worker job from Jupyter, verify completion and logs
- ✅ US3: Configure quotas, submit concurrent jobs, verify enforcement and isolation
- ✅ US2: Run CLI commands in CI/CD script, verify exit codes and logs
- ✅ US4: Create job via Dashboard wizard without YAML, verify submission
- ✅ US5: Submit 32-GPU LLM job, verify ≥85% scaling efficiency
