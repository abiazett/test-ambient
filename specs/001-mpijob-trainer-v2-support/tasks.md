# Tasks: MPIJobs Support in OpenShift AI via Trainer V2

**Feature Branch**: `001-mpijob-trainer-v2-support`
**Date**: 2025-10-28
**Input**: Design documents from `/specs/001-mpijob-trainer-v2-support/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/, quickstart.md

**Tests**: NOT explicitly requested in spec.md - focusing on implementation tasks only

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`
- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3, etc.)
- Include exact file paths in descriptions

## Path Conventions

This is a multi-repository project. Tasks span across:
- **kubeflow/training-operator/** - Go controller implementation
- **kubeflow/training-sdk/** - Python SDK implementation
- **opendatahub-io/odh-dashboard/** - React Dashboard UI
- **red-hat-data-services/openshift-ai-tests/** - E2E tests

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and development environment setup

- [ ] T001 Set up local development environment with Kind/Minikube cluster (Kubernetes 1.27+)
- [ ] T002 Install JobSet operator as prerequisite in dev cluster
- [ ] T003 [P] Install scheduler-plugins with coscheduling for gang scheduling in dev cluster
- [ ] T004 [P] Fork and clone kubeflow/training-operator repository
- [ ] T005 [P] Fork and clone kubeflow/training-sdk repository
- [ ] T006 [P] Fork and clone opendatahub-io/odh-dashboard repository
- [ ] T007 [P] Fork and clone red-hat-data-services/openshift-ai-tests repository
- [ ] T008 Create reference ClusterTrainingRuntime templates (mpi-horovod-gpu, mpi-openmpi-cpu) in yaml manifests
- [ ] T009 [P] Configure Go 1.21+ development environment with required dependencies
- [ ] T010 [P] Configure Python 3.11+ development environment with pytest and poetry

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core controller infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T011 Implement MPI runtime handler in kubeflow/training-operator/pkg/controller/trainer/mpi_runtime.go
- [ ] T012 [P] Implement SSH key generation logic in kubeflow/training-operator/pkg/controller/trainer/ssh_keys.go
- [ ] T013 [P] Implement hostfile generation logic in kubeflow/training-operator/pkg/controller/trainer/hostfile.go
- [ ] T014 [P] Implement gang scheduling PodGroup creation in kubeflow/training-operator/pkg/controller/trainer/gang_scheduling.go
- [ ] T015 Implement TrainJob reconciliation loop for MPI runtime in kubeflow/training-operator/pkg/controller/trainer/trainjob_controller.go
- [ ] T016 Add CEL validation rules for MPI runtime in kubeflow/training-operator/pkg/apis/kubeflow.org/v2alpha1/trainjob_types.go
- [ ] T017 Implement JobSet creation from TrainJob + RuntimeTemplate in kubeflow/training-operator/pkg/controller/trainer/jobset_builder.go
- [ ] T018 [P] Add status reporting logic for MPI jobs in kubeflow/training-operator/pkg/controller/trainer/status.go
- [ ] T019 [P] Implement error handling and failure reason detection in kubeflow/training-operator/pkg/controller/trainer/errors.go
- [ ] T020 Update ClusterTrainingRuntime and TrainingRuntime CRDs to support mpiImplementation field in kubeflow/training-operator/pkg/apis/kubeflow.org/v2alpha1/runtime_types.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Data Scientist Creates First MPIJob (Priority: P1) 🎯 MVP

**Goal**: Enable data scientists to create and run MPIJobs using Python SDK with same workflow as PyTorchJob

**Independent Test**: Create an MPIJob using Python SDK, verify it runs to completion, confirm workflow matches PyTorchJob creation

### Implementation for User Story 1

- [ ] T021 [P] [US1] Implement create_mpi_job() method in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T022 [P] [US1] Create TrainJob data model class in kubeflow/training-sdk/kubeflow/training/models.py
- [ ] T023 [P] [US1] Create TrainJobSpec data model class in kubeflow/training-sdk/kubeflow/training/models.py
- [ ] T024 [P] [US1] Create TrainJobStatus data model class in kubeflow/training-sdk/kubeflow/training/models.py
- [ ] T025 [US1] Implement get_job() method in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T026 [US1] Implement wait_for_job_completion() method with timeout in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T027 [US1] Add client-side validation for MPI job parameters in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T028 [P] [US1] Implement get_job_logs() method with pod type filtering in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T029 [P] [US1] Implement delete_job() method in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T030 [US1] Add error handling for common SDK exceptions (ValueError, RuntimeError, K8sApiException) in kubeflow/training-sdk/kubeflow/training/exceptions.py
- [ ] T031 [US1] Update SDK documentation with MPI examples in kubeflow/training-sdk/docs/mpi-jobs.md
- [ ] T032 [US1] Validate SDK against reference runtime template (mpi-horovod-gpu) in dev cluster

**Checkpoint**: At this point, User Story 1 should be fully functional - data scientists can create MPIJobs via SDK

---

## Phase 4: User Story 4 - Multi-Tenant Job Isolation and RBAC (Priority: P1)

**Goal**: Enforce namespace isolation, RBAC policies, and resource quotas for MPIJobs across tenants

**Independent Test**: Create MPIJobs in separate namespaces with different RBAC policies, verify users cannot access other tenants' jobs, confirm resource quotas enforced

### Implementation for User Story 4

- [ ] T033 [P] [US4] Implement namespace validation in controller in kubeflow/training-operator/pkg/controller/trainer/validation.go
- [ ] T034 [P] [US4] Implement resource quota pre-flight checks (launcher + workers) in kubeflow/training-operator/pkg/controller/trainer/quota.go
- [ ] T035 [US4] Add RBAC permission checks in controller admission webhook in kubeflow/training-operator/pkg/webhook/trainjob_webhook.go
- [ ] T036 [US4] Implement namespace-scoped pod labeling for network policies in kubeflow/training-operator/pkg/controller/trainer/labels.go
- [ ] T037 [P] [US4] Create standard NetworkPolicy template for MPI pod-to-pod SSH communication in manifests/network-policy.yaml
- [ ] T038 [P] [US4] Document required RBAC roles (trainjob-editor, trainjob-viewer) in docs/rbac.md
- [ ] T039 [US4] Add audit logging for TrainJob lifecycle events in kubeflow/training-operator/pkg/controller/trainer/audit.go
- [ ] T040 [US4] Implement SSH secret RBAC restrictions (namespace-scoped access only) in kubeflow/training-operator/pkg/controller/trainer/ssh_keys.go
- [ ] T041 [US4] Add validation that launcher/workers run in same namespace in kubeflow/training-operator/pkg/controller/trainer/validation.go

**Checkpoint**: At this point, multi-tenancy is enforced - users can only access jobs in their namespace

---

## Phase 5: User Story 2 - MLOps Engineer Troubleshoots Failed MPIJob (Priority: P1)

**Goal**: Provide clear error messages and diagnostics for MPIJob failures enabling self-service troubleshooting

**Independent Test**: Deliberately cause common failure scenarios (network policy blocking, resource constraints, SSH failures), verify engineer can diagnose and resolve using Dashboard and CLI within 15 minutes

### Implementation for User Story 2

- [ ] T042 [P] [US2] Implement gang scheduling timeout detection with clear error messages in kubeflow/training-operator/pkg/controller/trainer/gang_scheduling.go
- [ ] T043 [P] [US2] Implement SSH connection failure detection in kubeflow/training-operator/pkg/controller/trainer/ssh_keys.go
- [ ] T044 [P] [US2] Implement worker pod failure detection with specific failure reasons (OOM, ImagePull, CrashLoop) in kubeflow/training-operator/pkg/controller/trainer/pod_monitor.go
- [ ] T045 [US2] Add diagnostic events for gang scheduling progress (e.g., "Waiting for 4 worker pods. Only 2 nodes available") in kubeflow/training-operator/pkg/controller/trainer/events.go
- [ ] T046 [US2] Implement network policy validation diagnostics in kubeflow/training-operator/pkg/controller/trainer/network.go
- [ ] T047 [US2] Add structured error messages for 95% of failure scenarios per SC-006 in kubeflow/training-operator/pkg/controller/trainer/errors.go
- [ ] T048 [P] [US2] Implement launcher log aggregation in SDK in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T049 [P] [US2] Implement worker log aggregation with filtering by pod index in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T050 [US2] Add event streaming API for real-time status updates in kubeflow/training-sdk/kubeflow/training/client.py
- [ ] T051 [US2] Create troubleshooting guide with common failure scenarios in docs/troubleshooting.md

**Checkpoint**: At this point, engineers can effectively troubleshoot failures using clear diagnostics

---

## Phase 6: User Story 3 - Administrator Configures MPI Runtime Templates (Priority: P2)

**Goal**: Enable administrators to define and publish standardized MPI runtime templates with organizational policies

**Independent Test**: Create a ClusterTrainingRuntime for Horovod with GPU support, validate through test job, confirm users can reference it in MPIJobs

### Implementation for User Story 3

- [ ] T052 [P] [US3] Create Dashboard UI for runtime template management in opendatahub-io/odh-dashboard/frontend/src/pages/runtimes/RuntimeList.tsx
- [ ] T053 [P] [US3] Implement runtime template creation wizard in opendatahub-io/odh-dashboard/frontend/src/pages/runtimes/CreateRuntimeForm.tsx
- [ ] T054 [US3] Implement runtime template editor form with MPI-specific fields in opendatahub-io/odh-dashboard/frontend/src/pages/runtimes/RuntimeEditor.tsx
- [ ] T055 [US3] Add "Test Runtime" feature that creates minimal test job (1 launcher + 2 workers) in opendatahub-io/odh-dashboard/backend/src/routes/runtimes.ts
- [ ] T056 [P] [US3] Implement runtime template validation (image pull, SSH keys, MPI initialization) in kubeflow/training-operator/pkg/webhook/runtime_webhook.go
- [ ] T057 [P] [US3] Create reference runtime templates (mpi-horovod-gpu, mpi-intel-cpu, mpi-openmpi-distributed) in manifests/runtimes/
- [ ] T058 [US3] Add runtime template configuration guide for administrators in docs/admin-guide.md
- [ ] T059 [US3] Implement runtime template versioning and update handling in kubeflow/training-operator/pkg/controller/trainer/runtime_manager.go
- [ ] T060 [US3] Document required Security Context Constraints for SSH server in docs/security.md

**Checkpoint**: At this point, admins can create and manage runtime templates via Dashboard

---

## Phase 7: User Story 5 - Dashboard UI Job Creation and Monitoring (Priority: P2)

**Goal**: Enable data scientists to create and monitor MPIJobs entirely through Dashboard without writing YAML or using CLI

**Independent Test**: Create an MPIJob using only the Dashboard creation wizard, monitor progress through UI, access logs via Dashboard, confirm feature parity with SDK workflows

### Implementation for User Story 5

- [ ] T061 [P] [US5] Add MPI job type to unified Training Jobs list in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/JobList.tsx
- [ ] T062 [P] [US5] Create MPIJob creation wizard with MPI-specific fields in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/CreateMPIJobWizard.tsx
- [ ] T063 [US5] Implement form validation with contextual help tooltips in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/CreateMPIJobWizard.tsx
- [ ] T064 [P] [US5] Implement job details page with topology visualization (launcher + workers) in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/JobDetails.tsx
- [ ] T065 [P] [US5] Create topology view component showing pod health status in opendatahub-io/odh-dashboard/frontend/src/components/TopologyView.tsx
- [ ] T066 [P] [US5] Implement log viewer with tabbed interface (launcher, worker, aggregated) in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/LogsView.tsx
- [ ] T067 [P] [US5] Add log filtering and search capabilities in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/LogsView.tsx
- [ ] T068 [P] [US5] Implement resource utilization metrics charts (CPU, memory, GPU) in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/MetricsView.tsx
- [ ] T069 [P] [US5] Create events timeline component showing lifecycle events in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/EventsTimeline.tsx
- [ ] T070 [US5] Implement "Clone Job" feature to duplicate configuration in opendatahub-io/odh-dashboard/frontend/src/pages/trainjobs/JobDetails.tsx
- [ ] T071 [US5] Add real-time status updates via WebSocket in opendatahub-io/odh-dashboard/frontend/src/api/websocket.ts
- [ ] T072 [US5] Implement Dashboard REST API endpoints in opendatahub-io/odh-dashboard/backend/src/routes/trainjobs.ts
- [ ] T073 [US5] Add OAuth2 authentication integration in opendatahub-io/odh-dashboard/backend/src/auth/oauth.ts

**Checkpoint**: At this point, users can manage full MPIJob lifecycle through Dashboard UI only

---

## Phase 8: User Story 6 - CLI Job Management for Automation (Priority: P2)

**Goal**: Enable DevOps engineers to create and manage MPIJobs using CLI commands in CI/CD pipelines with same commands as other Trainer V2 jobs

**Independent Test**: Script complete MPIJob lifecycle operations (create, list, describe, delete) using only `oc` CLI commands, verify functionality matches PyTorchJob CLI workflows

### Implementation for User Story 6

- [ ] T074 [P] [US6] Implement `oc get trainjob` filtering by runtime type in kubeflow/training-operator/pkg/apis/kubeflow.org/v2alpha1/trainjob_types.go
- [ ] T075 [P] [US6] Add custom columns for `oc get trainjob` output (runtime, age, completion) in kubeflow/training-operator/config/crd/bases/additionalPrinterColumns
- [ ] T076 [P] [US6] Enhance `oc describe trainjob` output with MPI-specific details in kubeflow/training-operator/pkg/apis/kubeflow.org/v2alpha1/trainjob_status.go
- [ ] T077 [US6] Create CLI usage examples for common operations in docs/cli-guide.md
- [ ] T078 [US6] Add CLI automation examples for CI/CD pipelines in docs/automation.md
- [ ] T079 [US6] Document log access via `oc logs` with pod selector patterns in docs/cli-guide.md

**Checkpoint**: At this point, CLI workflows are fully functional for automation use cases

---

## Phase 9: User Story 7 - Migration from Legacy MPIJob v2beta1 (Priority: P3)

**Goal**: Provide clear migration guidance and tools for users with existing legacy MPIJob v2beta1 specifications

**Independent Test**: Convert a representative legacy MPIJob YAML to TrainJob format using migration documentation, deploy both versions side-by-side in development, validate identical training outcomes

### Implementation for User Story 7

- [ ] T080 [P] [US7] Create side-by-side API comparison table in docs/migration/api-comparison.md
- [ ] T081 [P] [US7] Create field mapping guide (v2beta1 → TrainJob) in docs/migration/field-mapping.md
- [ ] T082 [P] [US7] Document behavioral differences (launcher management, gang scheduling) in docs/migration/behavioral-changes.md
- [ ] T083 [US7] Create example conversions for common scenarios in docs/migration/examples/
- [ ] T084 [US7] Create best-effort YAML conversion script in tools/migrate-mpijob.py
- [ ] T085 [P] [US7] Create runtime template validator tool in tools/validate-runtime.py
- [ ] T086 [P] [US7] Document migration validation process in docs/migration/validation.md
- [ ] T087 [US7] Create migration timeline and deprecation notice in docs/migration/timeline.md
- [ ] T088 [US7] Add migration telemetry to track adoption in kubeflow/training-operator/pkg/controller/trainer/telemetry.go

**Checkpoint**: At this point, users have comprehensive migration guidance and tooling

---

## Phase 10: Observability and Monitoring

**Purpose**: Comprehensive observability for production deployments

- [ ] T089 [P] Implement Prometheus metrics for MPIJobs (duration, success/failure rate, time-to-start) in kubeflow/training-operator/pkg/metrics/metrics.go
- [ ] T090 [P] Implement MPI-specific metrics (gang scheduling duration, MPI initialization time) in kubeflow/training-operator/pkg/metrics/mpi_metrics.go
- [ ] T091 [P] Add structured logging with proper log levels in kubeflow/training-operator/pkg/controller/trainer/logger.go
- [ ] T092 [P] Implement log sanitization to prevent credential exposure in kubeflow/training-operator/pkg/controller/trainer/log_sanitizer.go
- [ ] T093 [P] Create Grafana dashboard template for MPIJob monitoring in manifests/grafana/mpijob-dashboard.json
- [ ] T094 Add OpenShift AI monitoring stack integration in docs/observability.md
- [ ] T095 [P] Implement audit logging for lifecycle events with user attribution in kubeflow/training-operator/pkg/controller/trainer/audit.go

---

## Phase 11: Integration Testing

**Purpose**: Validate end-to-end functionality and integration between components

- [ ] T096 [P] Create E2E test for basic MPIJob lifecycle (create, run, succeed) in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_mpi_job_lifecycle.py
- [ ] T097 [P] Create E2E test for gang scheduling timeout scenarios in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_gang_scheduling.py
- [ ] T098 [P] Create E2E test for multi-tenancy isolation in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_multi_tenancy.py
- [ ] T099 [P] Create E2E test for network policy validation in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_network_policies.py
- [ ] T100 [P] Create E2E test for SSH key generation and distribution in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_ssh_keys.py
- [ ] T101 Create reference workload test: Horovod PyTorch MNIST (2-4 GPUs, >98% accuracy) in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_reference_workloads.py
- [ ] T102 [P] Create reference workload test: Intel MPI TensorFlow benchmark in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_reference_workloads.py
- [ ] T103 [P] Create failure scenario tests (OOM, ImagePull, SSH failure, quota exceeded) in red-hat-data-services/openshift-ai-tests/tests/trainjobs/test_failure_scenarios.py
- [ ] T104 Set up CI/CD pipeline for E2E tests on dev cluster in .github/workflows/e2e-tests.yaml

---

## Phase 12: Performance Testing and Optimization

**Purpose**: Validate performance goals and optimize for production scale

- [ ] T105 [P] Create performance benchmark suite for job submission latency (<10s p90) in red-hat-data-services/openshift-ai-tests/tests/performance/test_submission_latency.py
- [ ] T106 [P] Create performance benchmark for gang scheduling (<30s for 8 workers) in red-hat-data-services/openshift-ai-tests/tests/performance/test_gang_scheduling_latency.py
- [ ] T107 [P] Create scaling efficiency tests (2, 4, 8, 16, 32, 64 workers) in red-hat-data-services/openshift-ai-tests/tests/performance/test_scaling_efficiency.py
- [ ] T108 Test concurrent job submission (10, 20, 50 jobs) in red-hat-data-services/openshift-ai-tests/tests/performance/test_concurrent_jobs.py
- [ ] T109 Test long-running job stability (24+ hours, memory leak detection) in red-hat-data-services/openshift-ai-tests/tests/performance/test_long_running.py
- [ ] T110 Optimize controller reconciliation loop performance in kubeflow/training-operator/pkg/controller/trainer/trainjob_controller.go
- [ ] T111 Profile and optimize SSH key generation performance in kubeflow/training-operator/pkg/controller/trainer/ssh_keys.go

---

## Phase 13: Security and Compliance

**Purpose**: Security hardening and compliance validation for enterprise deployment

- [ ] T112 [P] Create Security Context Constraints template for OpenShift in manifests/scc/mpi-training-scc.yaml
- [ ] T113 [P] Document FIPS compliance for SSH key generation in docs/security/fips.md
- [ ] T114 [P] Implement penetration testing for multi-tenancy isolation per SC-009 in docs/security/pentest-plan.md
- [ ] T115 Validate log sanitization prevents credential exposure in kubeflow/training-operator/pkg/controller/trainer/log_sanitizer.go
- [ ] T116 [P] Document required network policies for MPI communication in docs/security/network-policies.md
- [ ] T117 [P] Implement optional encrypted SSH communication for data-in-transit per FR-SEC-006 in kubeflow/training-operator/pkg/controller/trainer/ssh_keys.go
- [ ] T118 Validate RBAC enforcement across all components in red-hat-data-services/openshift-ai-tests/tests/security/test_rbac.py

---

## Phase 14: Documentation

**Purpose**: Comprehensive user and administrator documentation

- [ ] T119 [P] Create user guide covering all MPIJob features in docs/user-guide.md
- [ ] T120 [P] Create administrator guide for runtime template configuration in docs/admin-guide.md
- [ ] T121 [P] Create troubleshooting guide with common issues and solutions in docs/troubleshooting.md
- [ ] T122 [P] Create API reference documentation in docs/api-reference.md
- [ ] T123 [P] Create architecture documentation explaining how MPIJob works in docs/architecture.md
- [ ] T124 [P] Create performance tuning guide in docs/performance.md
- [ ] T125 [P] Update quickstart.md with validation results from real cluster testing
- [ ] T126 Create video tutorials for Dashboard UI workflows in docs/tutorials/
- [ ] T127 [P] Create FAQ document addressing common questions in docs/faq.md

---

## Phase 15: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements and production readiness

- [ ] T128 [P] Code review and refactoring across all repositories
- [ ] T129 [P] Implement consistent error message formatting across components
- [ ] T130 [P] Add telemetry for adoption tracking (60% adoption by month 10 per SC-002)
- [ ] T131 Performance optimization based on benchmark results
- [ ] T132 [P] Security hardening based on penetration test results
- [ ] T133 Create release notes and changelog in docs/CHANGELOG.md
- [ ] T134 [P] Create demo environment with sample MPIJobs for customer validation
- [ ] T135 Conduct 3-5 pilot customer validations per plan.md
- [ ] T136 [P] Update CLAUDE.md with MPIJob-specific commands and technologies
- [ ] T137 Run complete test suite (unit, integration, E2E, performance, security)
- [ ] T138 Final quickstart.md validation in production-like environment

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-9)**: All depend on Foundational phase completion
  - User stories can proceed in parallel (if staffed) after Phase 2
  - Or sequentially in priority order: US1(P1) → US4(P1) → US2(P1) → US3(P2) → US5(P2) → US6(P2) → US7(P3)
- **Observability (Phase 10)**: Can proceed in parallel with user stories
- **Integration Testing (Phase 11)**: Depends on US1, US2, US4 completion (P1 stories)
- **Performance Testing (Phase 12)**: Depends on Phase 11 completion
- **Security (Phase 13)**: Can proceed in parallel with user stories
- **Documentation (Phase 14)**: Can proceed in parallel with implementation
- **Polish (Phase 15)**: Depends on all desired user stories and testing being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 4 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P1)**: Can start after Foundational (Phase 2) - Benefits from US1 SDK but is independently testable
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 5 (P2)**: Can start after Foundational (Phase 2) - Benefits from US1 SDK but is independently testable
- **User Story 6 (P2)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 7 (P3)**: Can start after Foundational (Phase 2) - Should reference US1-US6 for migration examples

### Within Each User Story

- SDK methods can be implemented in parallel (marked [P])
- Dashboard components can be implemented in parallel (marked [P])
- Documentation can be written in parallel with implementation
- Each story should be independently testable upon completion

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories (US1, US2, US3, US4, US5, US6, US7) can start in parallel if team capacity allows
- Within each story, all tasks marked [P] can run in parallel
- Observability (Phase 10), Security (Phase 13), and Documentation (Phase 14) can proceed in parallel with user story implementation

---

## Parallel Example: User Story 1

```bash
# Launch all SDK methods for User Story 1 together:
Task: "Implement create_mpi_job() method in kubeflow/training-sdk/kubeflow/training/client.py"
Task: "Create TrainJob data model class in kubeflow/training-sdk/kubeflow/training/models.py"
Task: "Create TrainJobSpec data model class in kubeflow/training-sdk/kubeflow/training/models.py"
Task: "Create TrainJobStatus data model class in kubeflow/training-sdk/kubeflow/training/models.py"
Task: "Implement get_job_logs() method with pod type filtering in kubeflow/training-sdk/kubeflow/training/client.py"
Task: "Implement delete_job() method in kubeflow/training-sdk/kubeflow/training/client.py"
```

---

## Implementation Strategy

### MVP First (User Stories 1, 4, 2 Only - All P1)

1. Complete Phase 1: Setup (~1 week)
2. Complete Phase 2: Foundational (~3-4 weeks) - CRITICAL - blocks all stories
3. Complete Phase 3: User Story 1 (~2 weeks)
4. Complete Phase 4: User Story 4 (~1-2 weeks)
5. Complete Phase 5: User Story 2 (~1-2 weeks)
6. Complete Phase 11 (partial): Basic E2E tests for P1 stories
7. **STOP and VALIDATE**: Test P1 user stories independently
8. Deploy/demo if ready

**MVP Scope**: Weeks 1-10 (Setup + Foundation + US1 + US4 + US2 + Basic Tests)

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready (Weeks 1-5)
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!) (Weeks 6-7)
3. Add User Story 4 → Test independently → Deploy/Demo (Weeks 8-9)
4. Add User Story 2 → Test independently → Deploy/Demo (Weeks 10-11)
5. Add User Story 3 → Test independently → Deploy/Demo (Weeks 12-13)
6. Add User Story 5 → Test independently → Deploy/Demo (Weeks 14-16)
7. Add User Story 6 → Test independently → Deploy/Demo (Weeks 17-18)
8. Add User Story 7 → Test independently → Deploy/Demo (Weeks 19-20)
9. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together (Weeks 1-5)
2. Once Foundational is done:
   - Team A: User Story 1 (SDK) + User Story 6 (CLI)
   - Team B: User Story 4 (Multi-tenancy) + User Story 2 (Troubleshooting)
   - Team C: User Story 3 (Runtime Templates) + User Story 5 (Dashboard)
   - Team D: Observability + Security + Documentation
3. Stories complete and integrate independently

**Parallel Team Scope**: Weeks 1-12 (Setup + Foundation + All P1/P2 Stories in Parallel)

---

## Summary

**Total Tasks**: 138 tasks across 15 phases

**Task Breakdown by User Story**:
- Phase 1 (Setup): 10 tasks
- Phase 2 (Foundational): 10 tasks
- User Story 1 (P1): 12 tasks
- User Story 4 (P1): 9 tasks
- User Story 2 (P1): 10 tasks
- User Story 3 (P2): 9 tasks
- User Story 5 (P2): 13 tasks
- User Story 6 (P2): 6 tasks
- User Story 7 (P3): 9 tasks
- Phase 10 (Observability): 7 tasks
- Phase 11 (Integration Testing): 9 tasks
- Phase 12 (Performance): 7 tasks
- Phase 13 (Security): 7 tasks
- Phase 14 (Documentation): 9 tasks
- Phase 15 (Polish): 11 tasks

**Parallel Opportunities Identified**: 87 tasks marked [P] can run in parallel (63% of total)

**Independent Test Criteria**:
- US1: Create MPIJob via SDK, run to completion, match PyTorchJob workflow
- US2: Cause failure scenarios, diagnose and resolve within 15 minutes
- US3: Create ClusterTrainingRuntime, validate via test job
- US4: Create jobs in separate namespaces, verify isolation and quotas
- US5: Create and monitor MPIJob entirely through Dashboard UI
- US6: Script complete lifecycle using only CLI commands
- US7: Convert legacy YAML, deploy side-by-side, validate identical outcomes

**Suggested MVP Scope** (P1 stories only):
- Phase 1: Setup
- Phase 2: Foundational
- Phase 3: User Story 1 (Data Scientist Creates First MPIJob)
- Phase 4: User Story 4 (Multi-Tenant Job Isolation and RBAC)
- Phase 5: User Story 2 (MLOps Engineer Troubleshoots Failed MPIJob)
- Phase 11 (partial): Basic E2E tests

**Estimated Timeline**:
- MVP (P1 only): 10-12 weeks
- MVP + P2 stories: 16-18 weeks
- Full feature (P1 + P2 + P3): 20-24 weeks

**Format Validation**: ✓ ALL tasks follow checklist format (checkbox, ID, labels, file paths)

---

## Notes

- [P] tasks = different files/repos, no dependencies - safe to parallelize
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Multi-repository coordination required - use feature branches across all repos
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Tests are NOT included in main phases as they were not explicitly requested in spec.md
- Focus on enabling the 7 user stories with high-quality implementation and observability
