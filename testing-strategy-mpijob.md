# MPIJob Testing Strategy and Quality Assurance Requirements

**Feature**: MPIJobs Support in Red Hat OpenShift AI using KubeFlow Trainer V2
**Document Version**: 1.0
**Created**: 2025-10-28
**Author**: Neil (QA Engineer & Test Automation Architect)

---

## Executive Summary

This document outlines comprehensive testing strategies for implementing MPIJob support in OpenShift AI through KubeFlow Trainer V2. The testing approach covers unit tests, integration tests, contract tests, and end-to-end validation scenarios based on industry best practices from KubeFlow Training Operator, Kubernetes operator testing patterns, and OpenShift AI's existing test infrastructure.

**Key Testing Objectives:**
- Validate MPIJob creation, execution, and monitoring through CLI, SDK, and Dashboard
- Ensure multi-tenancy isolation and RBAC enforcement
- Verify gang scheduling behavior and distributed training orchestration
- Validate observability, logging, and failure diagnostics
- Ensure contract compatibility between components (controller, dashboard, SDK)

---

## Table of Contents

1. [Testing Framework Requirements](#1-testing-framework-requirements)
2. [Unit Testing Strategy](#2-unit-testing-strategy)
3. [Integration Testing Strategy](#3-integration-testing-strategy)
4. [Contract Testing Strategy](#4-contract-testing-strategy)
5. [End-to-End Testing Strategy](#5-end-to-end-testing-strategy)
6. [Test Automation and CI/CD](#6-test-automation-and-cicd)
7. [Test Data and Environment Management](#7-test-data-and-environment-management)
8. [Performance and Scale Testing](#8-performance-and-scale-testing)
9. [Security and Compliance Testing](#9-security-and-compliance-testing)
10. [Test Execution Schedule and Coverage](#10-test-execution-schedule-and-coverage)

---

## 1. Testing Framework Requirements

### 1.1 Go-Based Testing (Operator/Controller)

**Test Category**: Unit & Integration Tests
**Testing Approach**: Use Ginkgo BDD framework with Gomega matchers and EnvTest for controller logic

**Required Tooling**:
- **Ginkgo v2**: BDD testing framework for Go
- **Gomega**: Matcher library for expressive assertions
- **EnvTest**: Kubernetes API server and etcd for integration testing
- **controller-runtime/fake**: Mock Kubernetes clients for unit tests

**Success Criteria**:
- All controller reconciliation logic tested with >80% code coverage
- CRD validation and defaulting verified
- Webhook validation tested for invalid inputs
- Tests run in <5 minutes locally without real cluster

**Edge Cases**:
- Launcher pod creation failure scenarios
- Worker pod partial failures during gang scheduling
- SSH key generation and distribution errors
- Runtime template validation failures

**Implementation Details from Research**:
```
Based on Kubeflow Training Operator patterns:
- Use suite_test.go with BeforeSuite/AfterSuite for EnvTest lifecycle
- Test reconcile loops with Eventually() for async operations
- Mock Kubernetes API server responses for edge cases
- Use fake clients for testing error handling without cluster
```

---

### 1.2 Python-Based Testing (SDK, Dashboard, E2E)

**Test Category**: Integration, Contract, E2E Tests
**Testing Approach**: Use pytest framework with fixtures for test data and environment management

**Required Tooling**:
- **pytest**: Test framework for Python SDK and E2E tests
- **pytest-bdd**: Behavior-driven testing for user scenarios
- **pytest-xdist**: Parallel test execution
- **requests**: HTTP client for Dashboard API testing
- **kubernetes Python client**: For E2E cluster interactions

**Success Criteria**:
- Python SDK methods tested with unit and integration tests
- Dashboard API endpoints validated with contract tests
- E2E scenarios executable against live OpenShift clusters
- Test results published to CI/CD pipeline

**Edge Cases**:
- Network timeouts during API calls
- Cluster resource constraints during E2E tests
- Authentication/authorization failures
- Concurrent job submission conflicts

**Implementation Details from Research**:
```
Based on Red Hat OpenShift AI testing patterns (red-hat-data-services/openshift-ai-tests):
- Use Poetry for dependency management
- Organize tests by component (sdk/, dashboard/, e2e/)
- Use --log-cli-level for detailed logging
- Implement fixtures for cluster setup and teardown
```

---

### 1.3 Testing Kubernetes Operators Best Practices

**Key Patterns from Research**:

1. **EnvTest Integration**:
   - Spins up local Kubernetes API server and etcd
   - No kubelet, controller-manager, or scheduler needed
   - Tests run quickly without full cluster overhead
   - Validates CRD schemas, webhooks, and controller logic

2. **Ginkgo BDD Structure**:
   - `Describe` blocks for grouping related tests
   - `Context` blocks for different scenarios
   - `It` blocks for individual test cases
   - `BeforeEach/AfterEach` for test setup/cleanup

3. **Gomega Assertions**:
   - Use `Eventually()` for async operations (polling with timeout)
   - Use `Consistently()` for stable state verification
   - Expressive matchers: `HaveOccurred()`, `BeNil()`, `Equal()`

4. **Controller Testing Strategy**:
   - Test reconcile function independently with mock clients
   - Verify correct resource creation (launcher pod, worker pods, services)
   - Test error handling and retry logic
   - Validate status updates and conditions

---

## 2. Unit Testing Strategy

### 2.1 Controller Unit Tests (Go)

**Test Category**: Unit Test
**Testing Approach**: Test individual functions and methods in isolation using fake Kubernetes clients

**Required Tooling**:
- Go testing framework (`testing` package)
- `controller-runtime/pkg/client/fake`
- Gomega matchers for assertions

**Test Coverage Areas**:

#### 2.1.1 TrainJob Controller Reconciliation
```
Test Cases:
- TC-UNIT-001: Reconcile creates launcher pod with correct configuration
- TC-UNIT-002: Reconcile creates worker pods with correct replica count
- TC-UNIT-003: Reconcile generates SSH keys and creates secret
- TC-UNIT-004: Reconcile updates TrainJob status based on pod status
- TC-UNIT-005: Reconcile handles launcher pod failure gracefully
- TC-UNIT-006: Reconcile handles partial worker pod failures
- TC-UNIT-007: Reconcile respects resource quotas and rejects over-quota jobs
- TC-UNIT-008: Reconcile validates runtime reference exists
```

**Success Criteria**:
- All reconciliation paths covered (happy path + error paths)
- Status updates correctly reflect pod states
- Proper error messages for failure scenarios
- Tests execute in <100ms per test case

**Edge Cases**:
- Runtime template not found
- Invalid worker count (0, negative)
- Resource requests exceed node capacity
- Namespace being deleted during reconciliation

---

#### 2.1.2 CRD Validation and Defaulting
```
Test Cases:
- TC-UNIT-009: Webhook rejects invalid worker count (<1)
- TC-UNIT-010: Webhook rejects invalid slots per worker (<1)
- TC-UNIT-011: Webhook applies default values for optional fields
- TC-UNIT-012: Webhook validates runtime reference format
- TC-UNIT-013: Webhook rejects conflicting resource specifications
- TC-UNIT-014: Webhook validates container image format
```

**Success Criteria**:
- All validation rules tested
- Default values applied correctly
- Clear error messages for validation failures

---

#### 2.1.3 Gang Scheduling Logic
```
Test Cases:
- TC-UNIT-015: Gang scheduler waits for all pods before starting
- TC-UNIT-016: Gang scheduler times out if pods don't start
- TC-UNIT-017: Gang scheduler handles partial pod scheduling
- TC-UNIT-018: Gang scheduler respects priority classes
```

**Success Criteria**:
- Atomic scheduling verified (all pods or none)
- Timeout behavior matches configuration
- Status updates show gang scheduling progress

---

### 2.2 SDK Unit Tests (Python)

**Test Category**: Unit Test
**Testing Approach**: Test SDK client methods using mocked Kubernetes API responses

**Required Tooling**:
- pytest
- unittest.mock for mocking Kubernetes API
- pytest fixtures for test data

**Test Coverage Areas**:

#### 2.2.1 MPIJob Creation via SDK
```
Test Cases:
- TC-SDK-001: create_train_job() generates correct TrainJob manifest
- TC-SDK-002: create_train_job() validates required parameters
- TC-SDK-003: create_train_job() applies default values
- TC-SDK-004: create_train_job() handles API errors gracefully
- TC-SDK-005: create_train_job() supports all MPI-specific parameters
```

**Success Criteria**:
- Correct API calls made to Kubernetes
- Manifest structure matches TrainJob CRD schema
- Error handling provides actionable messages

---

#### 2.2.2 MPIJob Status and Monitoring
```
Test Cases:
- TC-SDK-006: get_train_job_status() returns correct status
- TC-SDK-007: get_train_job_status() handles non-existent jobs
- TC-SDK-008: list_train_jobs() filters by runtime correctly
- TC-SDK-009: get_train_job_logs() retrieves launcher logs
- TC-SDK-010: get_train_job_logs() retrieves worker logs
```

**Success Criteria**:
- Status accurately reflects cluster state
- Log retrieval handles multiple pods correctly
- Filtering and pagination work as expected

---

## 3. Integration Testing Strategy

### 3.1 Controller Integration Tests (Go + EnvTest)

**Test Category**: Integration Test
**Testing Approach**: Run controller against EnvTest Kubernetes API server to validate end-to-end controller behavior

**Required Tooling**:
- EnvTest (controller-runtime/envtest)
- Ginkgo/Gomega
- Test CRD manifests

**Test Coverage Areas**:

#### 3.1.1 Full MPIJob Lifecycle
```
Test Cases:
- TC-INT-001: Create TrainJob → Launcher pod created → Workers created → Job succeeds
- TC-INT-002: Create TrainJob → Gang scheduling → All pods ready → MPI initialization
- TC-INT-003: Create TrainJob → Launcher fails → Status shows launcher failure
- TC-INT-004: Create TrainJob → Worker fails → Status shows worker failure
- TC-INT-005: Delete TrainJob → All pods deleted → Resources cleaned up
```

**Success Criteria**:
- Full reconciliation loop executes correctly
- Status updates propagate through lifecycle
- Resource cleanup completes without leaks

**Edge Cases**:
- Rapid job creation and deletion
- Multiple jobs in same namespace
- Controller restart during job execution
- API server unavailability

---

#### 3.1.2 Gang Scheduling Integration
```
Test Cases:
- TC-INT-006: Gang scheduler waits for resources → All pods scheduled atomically
- TC-INT-007: Gang scheduler timeout → Job fails with clear message
- TC-INT-008: Gang scheduler handles node failure during scheduling
- TC-INT-009: Gang scheduler respects priority classes and preemption
```

**Success Criteria**:
- All pods scheduled together or none scheduled
- Timeout behavior correct with configurable values
- Events emitted for gang scheduling progress

---

#### 3.1.3 RBAC and Multi-Tenancy
```
Test Cases:
- TC-INT-010: User without permissions cannot create TrainJob
- TC-INT-011: User cannot access TrainJob in different namespace
- TC-INT-012: Service account has minimal required permissions
- TC-INT-013: Resource quotas enforced for launcher + workers aggregate
```

**Success Criteria**:
- RBAC policies correctly enforced
- Namespace isolation maintained
- Resource quotas prevent over-allocation

**Edge Cases**:
- RBAC changes during job execution
- Quota changes after job created but before scheduled
- Cross-namespace resource references rejected

---

### 3.2 SDK Integration Tests (Python + Live Cluster)

**Test Category**: Integration Test
**Testing Approach**: Test SDK against real OpenShift cluster to validate API compatibility

**Required Tooling**:
- pytest
- kubernetes Python client
- OpenShift test cluster (dev/staging)

**Test Coverage Areas**:

#### 3.2.1 SDK-to-Controller Integration
```
Test Cases:
- TC-INT-014: SDK creates TrainJob → Controller reconciles → Job runs
- TC-INT-015: SDK retrieves status → Matches controller status updates
- TC-INT-016: SDK deletes job → Controller cleans up resources
- TC-INT-017: SDK streams logs → Logs from launcher and workers retrieved
```

**Success Criteria**:
- SDK operations translate to correct Kubernetes API calls
- Status retrieved matches actual cluster state
- Log streaming works for all pod types

---

#### 3.2.2 Dashboard-to-Controller Integration
```
Test Cases:
- TC-INT-018: Dashboard API creates TrainJob → Job appears in list
- TC-INT-019: Dashboard API retrieves status → Real-time updates visible
- TC-INT-020: Dashboard API retrieves logs → Launcher and worker logs accessible
- TC-INT-021: Dashboard API deletes job → Resources removed from cluster
```

**Success Criteria**:
- Dashboard API correctly interacts with Kubernetes
- Real-time status updates work (WebSocket/SSE)
- UI reflects actual cluster state accurately

---

### 3.3 Network and Communication Testing

**Test Category**: Integration Test
**Testing Approach**: Validate launcher-to-worker SSH communication and network policies

**Required Tooling**:
- pytest
- Network policy test utilities
- SSH connectivity test scripts

**Test Coverage Areas**:

#### 3.3.1 Launcher-Worker SSH Communication
```
Test Cases:
- TC-INT-022: Launcher connects to all workers via SSH
- TC-INT-023: SSH key distribution works for all workers
- TC-INT-024: MPI hostfile generated with correct worker addresses
- TC-INT-025: mpirun successfully spawns processes on all workers
```

**Success Criteria**:
- SSH connections established within timeout
- MPI initialization completes successfully
- All workers responsive to launcher commands

**Edge Cases**:
- Worker pod IP changes during job (pod restart)
- SSH connection timeout/retry behavior
- Incomplete SSH key distribution
- Network partition between launcher and workers

---

#### 3.3.2 Network Policy Enforcement
```
Test Cases:
- TC-INT-026: Network policy allows intra-namespace pod-to-pod SSH
- TC-INT-027: Network policy blocks cross-namespace SSH
- TC-INT-028: Network policy allows MPI communication ports
- TC-INT-029: Network policy works with Istio disabled (annotation)
```

**Success Criteria**:
- Required traffic allowed between launcher and workers
- Cross-tenant traffic blocked
- Istio sidecar injection disabled for MPI pods

**Edge Cases**:
- Default-deny network policies in namespace
- Multiple network policies applying to same pods
- Network policy changes during job execution

---

## 4. Contract Testing Strategy

### 4.1 API Contract Between Components

**Test Category**: Contract Test
**Testing Approach**: Validate API contracts between Trainer V2 controller, Dashboard backend, and SDK client

**Required Tooling**:
- Pact (consumer-driven contract testing)
- pytest for contract test execution
- OpenAPI schema validators

**Test Coverage Areas**:

#### 4.1.1 TrainJob CRD Schema Contract
```
Test Cases:
- TC-CONTRACT-001: TrainJob CRD schema version compatibility (v2alpha1)
- TC-CONTRACT-002: Required fields present in all SDK/Dashboard/CLI operations
- TC-CONTRACT-003: Optional fields handled correctly when absent
- TC-CONTRACT-004: Status subresource schema matches expectations
- TC-CONTRACT-005: MPI-specific fields (slotsPerWorker, etc.) present
```

**Success Criteria**:
- CRD schema changes don't break existing clients
- All clients produce valid manifests
- Status structure consistent across versions

**Edge Cases**:
- Schema version migration (alpha to beta)
- Deprecated fields still supported
- Unknown fields handled gracefully

---

#### 4.1.2 Dashboard API Contract
```
Test Cases:
- TC-CONTRACT-006: Dashboard GET /api/trainjobs returns correct structure
- TC-CONTRACT-007: Dashboard POST /api/trainjobs accepts MPI parameters
- TC-CONTRACT-008: Dashboard GET /api/trainjobs/{id}/logs streams logs correctly
- TC-CONTRACT-009: Dashboard DELETE /api/trainjobs/{id} removes job
- TC-CONTRACT-010: Dashboard WebSocket provides real-time status updates
```

**Success Criteria**:
- API responses match OpenAPI specification
- Breaking changes detected by contract tests
- Backward compatibility maintained across releases

---

#### 4.1.3 SDK Client Contract
```
Test Cases:
- TC-CONTRACT-011: SDK TrainingClient.create_train_job() signature stable
- TC-CONTRACT-012: SDK returns expected data structures (dicts/objects)
- TC-CONTRACT-013: SDK error handling consistent with API errors
- TC-CONTRACT-014: SDK supports all TrainJob CRD fields
```

**Success Criteria**:
- SDK method signatures backward compatible
- Return types match documentation
- Error types consistent and catchable

---

### 4.2 Storage System Contracts

**Test Category**: Contract Test
**Testing Approach**: Validate contracts with storage systems for training data and model artifacts

**Required Tooling**:
- pytest
- boto3 (S3 testing)
- OpenShift OCS test utilities

**Test Coverage Areas**:

#### 4.2.1 S3-Compatible Storage
```
Test Cases:
- TC-CONTRACT-015: MPIJob mounts S3 bucket credentials correctly
- TC-CONTRACT-016: Workers can read training data from S3
- TC-CONTRACT-017: Launcher can write model artifacts to S3
- TC-CONTRACT-018: S3 authentication failures handled gracefully
```

**Success Criteria**:
- S3 access works from all pods (launcher + workers)
- Credentials securely passed via secrets
- Storage errors reported in job status

---

#### 4.2.2 OpenShift Container Storage (OCS)
```
Test Cases:
- TC-CONTRACT-019: MPIJob mounts OCS PVCs correctly
- TC-CONTRACT-020: Workers can read/write to shared PVC
- TC-CONTRACT-021: PVC cleanup after job completion
```

**Success Criteria**:
- PVC mounting consistent across pods
- Read/write performance acceptable
- No orphaned PVCs after job deletion

---

## 5. End-to-End Testing Strategy

### 5.1 User Scenario E2E Tests

**Test Category**: End-to-End Test
**Testing Approach**: Test complete user workflows from job creation to completion using reference workloads

**Required Tooling**:
- pytest with pytest-bdd for scenario testing
- OpenShift test cluster with GPU nodes
- Reference training workloads (Horovod, Intel MPI)

**Test Coverage Areas**:

#### 5.1.1 Data Scientist Workflow (User Story 1)
```
Test Cases:
- TC-E2E-001: Create MPIJob via SDK → Monitor via Dashboard → Retrieve results
  - Setup: Dev cluster with 4 GPU nodes
  - Workload: Horovod MNIST training
  - Success: Job completes, accuracy >95%, logs accessible

- TC-E2E-002: Create MPIJob via Dashboard → Monitor status → Download artifacts
  - Setup: Staging cluster with 8 GPU nodes
  - Workload: Horovod ResNet-50 training
  - Success: Job completes, model saved to S3, Dashboard shows metrics

- TC-E2E-003: Create MPIJob via CLI → Stream logs → Cancel mid-execution
  - Setup: Dev cluster with 2 CPU nodes
  - Workload: Intel MPI parameter server
  - Success: Job starts, logs visible, graceful cancellation
```

**Success Criteria**:
- Time-to-first-job <15 minutes for new users
- Job status updates visible within 10 seconds
- Logs accessible with <2 clicks in Dashboard
- 100% success rate for reference workloads

**Edge Cases**:
- Cluster resource contention during job submission
- Network interruption during log streaming
- Dashboard connection lost during job execution

---

#### 5.1.2 MLOps Troubleshooting Workflow (User Story 2)
```
Test Cases:
- TC-E2E-004: Job fails (OOM) → Engineer identifies failure → Adjusts resources → Reruns
  - Setup: Cluster with limited GPU memory
  - Failure: Worker OOM error
  - Success: Error visible in Dashboard within 30s, clear diagnostic message

- TC-E2E-005: Gang scheduling timeout → Engineer checks resources → Reduces workers → Success
  - Setup: Cluster with limited node capacity
  - Failure: Gang scheduling timeout (5 min)
  - Success: Clear message "Waiting for 8 workers, only 4 nodes available"

- TC-E2E-006: Network policy blocks SSH → Engineer reviews logs → Adds policy → Reruns
  - Setup: Namespace with restrictive network policy
  - Failure: Launcher cannot connect to workers
  - Success: Error message points to network policy issue
```

**Success Criteria**:
- Failure reason visible within 1 minute of failure
- Root cause identifiable from logs/events in <10 minutes
- Remediation steps clear from error messages
- 90% of failures self-serviceable without support

---

#### 5.1.3 Administrator Runtime Configuration (User Story 3)
```
Test Cases:
- TC-E2E-007: Admin creates ClusterTrainingRuntime → Test runtime → Publish → Users create jobs
  - Setup: Cluster admin access
  - Runtime: Horovod with A100 GPU support
  - Success: Test job passes, users can reference runtime

- TC-E2E-008: Admin updates runtime configuration → Existing jobs unaffected → New jobs use update
  - Setup: Existing running MPIJob
  - Update: Change default slots per worker
  - Success: Running job continues, new job uses new defaults
```

**Success Criteria**:
- Runtime testing validates configuration before publication
- Runtime updates don't disrupt running jobs
- Users can discover and use published runtimes

---

#### 5.1.4 Multi-Tenancy Isolation (User Story 4)
```
Test Cases:
- TC-E2E-009: Team A creates job → Team B cannot access → Namespace isolation verified
  - Setup: Two namespaces with different RBAC
  - Test: Team B attempts to access Team A's job via CLI, SDK, Dashboard
  - Success: All access attempts denied with RBAC error

- TC-E2E-010: User attempts to exceed quota → Job rejected → Clear error message
  - Setup: Namespace with 4 GPU quota
  - Test: User submits job requiring 8 GPUs
  - Success: Job rejected immediately with quota error
```

**Success Criteria**:
- Zero cross-tenant access violations
- Resource quotas enforced before scheduling
- RBAC errors clear and actionable

---

### 5.2 Reference Workload Validation

**Test Category**: End-to-End Test
**Testing Approach**: Run standard distributed training workloads to validate functional correctness and performance

**Required Tooling**:
- Horovod training examples
- Intel MPI benchmarks
- TensorFlow/PyTorch with MPI
- Performance profiling tools

**Reference Workloads**:

#### 5.2.1 Horovod PyTorch Example
```
Workload: TC-E2E-011
Description: Horovod distributed PyTorch training on MNIST dataset
Configuration:
  - 4 worker nodes
  - 2 GPUs per worker (8 total GPUs)
  - 2 slots per worker
  - Batch size: 64 per GPU
Expected Results:
  - Training completes in <10 minutes
  - Final accuracy: >98%
  - Scaling efficiency: >85% (vs single GPU)
  - All 8 GPUs utilized during training
```

**Success Criteria**:
- Workload completes successfully
- Accuracy matches expected baseline
- Resource utilization >80% for all GPUs
- MPI communication overhead <15%

---

#### 5.2.2 Intel MPI TensorFlow Benchmark
```
Workload: TC-E2E-012
Description: TensorFlow ResNet-50 training with Intel MPI
Configuration:
  - 8 worker nodes (CPU only)
  - 4 slots per worker (32 MPI ranks)
  - Batch size: 32 per rank
  - Synthetic data for reproducibility
Expected Results:
  - Training throughput: >500 images/sec aggregate
  - Scaling efficiency: >70% (vs single node)
  - MPI collectives: <200ms per iteration
  - No worker failures or hangs
```

**Success Criteria**:
- Performance meets baseline benchmarks
- All MPI ranks participate correctly
- No communication bottlenecks detected
- Reproducible results across runs

---

#### 5.2.3 Horovod TensorFlow BERT Fine-tuning
```
Workload: TC-E2E-013
Description: Distributed BERT fine-tuning with Horovod
Configuration:
  - 4 worker nodes
  - 4 GPUs per worker (16 total GPUs)
  - 4 slots per worker
  - GLUE dataset (SQuAD)
Expected Results:
  - Fine-tuning completes in <2 hours
  - Final F1 score: >88%
  - GPU memory utilization: 70-90%
  - Model artifacts saved to S3
```

**Success Criteria**:
- Large-scale training completes successfully
- Model quality matches expected metrics
- Storage integration works (data in, model out)
- No out-of-memory errors

---

### 5.3 Failure Scenario Testing

**Test Category**: End-to-End Test
**Testing Approach**: Deliberately inject failures to validate error handling and recovery

**Test Coverage Areas**:

#### 5.3.1 Resource Constraint Failures
```
Test Cases:
- TC-E2E-014: Insufficient GPU resources → Gang scheduling timeout → Clear error
  - Inject: Request 16 GPUs on 8 GPU cluster
  - Expected: Timeout after 5 min, error message shows available resources

- TC-E2E-015: Memory limit too low → Worker OOM → Job fails with diagnostic
  - Inject: Set memory limit to 4GB for 10GB model
  - Expected: Worker OOM, Dashboard highlights which worker failed

- TC-E2E-016: Node failure during training → Worker pod evicted → Job fails gracefully
  - Inject: Drain node during training
  - Expected: Job fails, status shows node drain event
```

**Success Criteria**:
- All failures detected and reported within 2 minutes
- Error messages actionable and specific
- System doesn't enter unrecoverable state

---

#### 5.3.2 Network Failure Scenarios
```
Test Cases:
- TC-E2E-017: Network partition between launcher and worker → SSH failure → Job fails
  - Inject: Block TCP port 22 with network policy
  - Expected: Launcher logs show "Connection refused", job fails

- TC-E2E-018: Intermittent network latency → MPI communication timeout → Retry succeeds
  - Inject: Add 500ms latency between pods
  - Expected: MPI initialization retries, eventually succeeds or fails with timeout

- TC-E2E-019: DNS resolution failure → Worker address lookup fails → Clear diagnostic
  - Inject: Break DNS in namespace
  - Expected: Error message indicates DNS issue
```

**Success Criteria**:
- Network failures don't cause silent hangs
- Timeout values configurable and respected
- Diagnostic messages identify network as root cause

---

#### 5.3.3 Configuration and Setup Failures
```
Test Cases:
- TC-E2E-020: Invalid container image → ImagePullBackOff → Error visible in Dashboard
  - Inject: Use non-existent image tag
  - Expected: Dashboard shows "ImagePullBackOff" for launcher/workers

- TC-E2E-021: Missing runtime reference → Job validation fails → Clear error before submission
  - Inject: Reference non-existent runtime
  - Expected: Webhook rejects job with "Runtime 'xyz' not found"

- TC-E2E-022: Istio sidecar injection conflicts → MPI initialization hangs → Diagnostic message
  - Inject: Enable Istio without disabling sidecar
  - Expected: Timeout with message suggesting Istio annotation
```

**Success Criteria**:
- Configuration errors caught at submission time (webhook)
- Runtime errors visible in Dashboard within 1 minute
- Error messages reference documentation for resolution

---

### 5.4 Observability and Logging Validation

**Test Category**: End-to-End Test
**Testing Approach**: Validate metrics, logs, and events are correctly captured and accessible

**Test Coverage Areas**:

#### 5.4.1 Metrics Collection
```
Test Cases:
- TC-E2E-023: MPIJob emits job duration metric → Visible in Prometheus
- TC-E2E-024: MPIJob emits success/failure rate → Dashboard shows aggregate stats
- TC-E2E-025: MPIJob emits gang scheduling duration → Monitoring shows scheduling overhead
- TC-E2E-026: MPIJob emits resource utilization → GPU/CPU/memory metrics visible
```

**Success Criteria**:
- All metrics available in Prometheus within 30 seconds
- Metrics match actual job behavior (durations, statuses)
- Metrics queryable and usable for alerting

---

#### 5.4.2 Log Aggregation
```
Test Cases:
- TC-E2E-027: Launcher logs aggregated → Accessible via Dashboard and CLI
- TC-E2E-028: Worker logs aggregated → Individual and combined views available
- TC-E2E-029: MPI initialization logs captured → SSH setup and hostfile generation visible
- TC-E2E-030: Application logs captured → Training progress and errors visible
```

**Success Criteria**:
- All pod logs accessible through unified interface
- Logs retained for minimum 7 days
- Log search and filtering functional
- Timestamps and pod labels correct

---

#### 5.4.3 Event Tracking
```
Test Cases:
- TC-E2E-031: Job lifecycle events captured → Created, Scheduling, Running, Completed
- TC-E2E-032: Pod events captured → Pending, Running, Failed with reasons
- TC-E2E-033: Error events captured → OOM, ImagePull, CrashLoopBackOff
- TC-E2E-034: Custom MPI events captured → Gang scheduling, SSH setup, MPI initialization
```

**Success Criteria**:
- Events visible in Dashboard timeline
- Events include timestamps and context
- Events retained for 7 days minimum

---

## 6. Test Automation and CI/CD

### 6.1 Continuous Integration Testing

**Test Category**: CI/CD Automation
**Testing Approach**: Automated test execution on every code commit and pull request

**Required Tooling**:
- GitHub Actions or OpenShift Pipelines
- Container registry for test images
- Ephemeral test clusters (Kind, Minikube, or OpenShift Dev Sandbox)

**CI/CD Pipeline Structure**:

```
Stage 1: Pre-Commit Validation (Local Developer)
  - Linting: golangci-lint, pylint
  - Unit tests: Go unit tests, Python unit tests
  - Duration: <2 minutes

Stage 2: Pull Request Validation (CI)
  - Unit tests: Full test suite (Go + Python)
  - Integration tests: EnvTest-based controller tests
  - Contract tests: API schema validation
  - Duration: <10 minutes

Stage 3: Merge to Main (CI)
  - Unit + Integration + Contract tests
  - E2E tests: Subset of critical paths on dev cluster
  - Duration: <30 minutes

Stage 4: Release Candidate (CI)
  - Full E2E test suite on staging cluster
  - Performance benchmarks
  - Security scans
  - Duration: <2 hours

Stage 5: Production Release (Manual + CI)
  - Smoke tests on production cluster
  - Rollback validation
  - Duration: <30 minutes
```

**Success Criteria**:
- All tests automated and run on every PR
- Test failures block merges
- Flaky tests identified and fixed within 1 week
- CI pipeline completes in <15 minutes for fast feedback

---

### 6.2 Test Environment Management

**Test Category**: Infrastructure
**Testing Approach**: Manage multiple test environments for different testing stages

**Test Environments**:

#### 6.2.1 Local Development Environment
```
Purpose: Developer unit and integration testing
Infrastructure:
  - EnvTest for controller tests (no real cluster)
  - Docker Desktop or Podman for local SDK tests
  - Mock Kubernetes API for fast iteration
Resources: Developer laptop
Duration: Tests run in <5 minutes
```

#### 6.2.2 Ephemeral CI Environment
```
Purpose: Pull request validation
Infrastructure:
  - Kind cluster (Kubernetes in Docker)
  - Lightweight test workloads
  - No GPU requirements
Resources: CI runners (GitHub Actions, OpenShift Pipelines)
Duration: Cluster created in <2 minutes, tests run in <10 minutes
Lifecycle: Created per PR, destroyed after tests
```

#### 6.2.3 Persistent Dev/Staging Cluster
```
Purpose: E2E testing with realistic workloads
Infrastructure:
  - OpenShift cluster with 4-8 GPU nodes
  - Multiple namespaces for team testing
  - Persistent storage (S3, OCS)
Resources: Dedicated OpenShift cluster
Duration: Long-running (months)
Lifecycle: Persistent, reset between test runs
```

#### 6.2.4 Pre-Production Cluster
```
Purpose: Release candidate validation
Infrastructure:
  - Mirror of production configuration
  - Same scale and hardware as production
  - Production-like network policies and RBAC
Resources: Dedicated OpenShift cluster
Duration: Long-running (months)
Lifecycle: Persistent, production configuration
```

**Success Criteria**:
- Environments available 99%+ of time
- Environment setup automated (Infrastructure as Code)
- Test data and cleanup automated
- Resource costs optimized (ephemeral where possible)

---

### 6.3 Test Data Management

**Test Category**: Test Data Strategy
**Testing Approach**: Manage training datasets, models, and test fixtures for reproducible testing

**Test Data Sources**:

#### 6.3.1 Synthetic Data
```
Use Cases: Performance benchmarks, reproducibility tests
Examples:
  - Random tensors for throughput testing
  - Generated MNIST-like images for functional testing
  - Synthetic parameter server data
Advantages: Fast generation, no external dependencies, consistent
Storage: Generated on-demand in tests
```

#### 6.3.2 Standard ML Datasets
```
Use Cases: Functional validation, accuracy testing
Examples:
  - MNIST (classification baseline)
  - CIFAR-10 (image classification)
  - GLUE/SQuAD (NLP tasks)
Advantages: Well-known baselines, community benchmarks
Storage: S3 bucket or OCS PVC, cached locally
```

#### 6.3.3 Test Fixtures
```
Use Cases: Unit and integration tests
Examples:
  - TrainJob YAML manifests
  - Runtime template definitions
  - Mock API responses
  - Test logs and events
Storage: Git repository (version controlled)
```

**Test Data Management Strategy**:
- Small datasets (<100MB) stored in Git LFS
- Large datasets (>1GB) stored in S3 with versioning
- Synthetic data generated programmatically
- Test data cleanup automated after test runs
- Dataset versions pinned for reproducibility

---

## 7. Test Data and Environment Management

### 7.1 Test Cluster Configuration

**Required Cluster Features**:
```
Cluster Requirements:
  - OpenShift 4.12+ or Kubernetes 1.25+
  - GPU nodes: 4-8 nodes with NVIDIA A100 or V100 GPUs
  - CPU nodes: 8+ nodes with 16 CPU cores each
  - Storage: S3-compatible (MinIO or AWS S3) + OCS/Rook-Ceph
  - Networking: Calico or OpenShift SDN with network policy support
  - Operators: KubeFlow Trainer V2, MPI Operator, Kueue (optional)
  - Monitoring: Prometheus, Grafana, OpenShift Logging
```

**Test Namespaces**:
```
- mpijob-e2e-tests: End-to-end test execution
- mpijob-sdk-tests: SDK integration testing
- mpijob-dashboard-tests: Dashboard integration testing
- team-a-test: Multi-tenancy testing (team A)
- team-b-test: Multi-tenancy testing (team B)
- mpijob-perf-tests: Performance and scale testing
```

---

### 7.2 Test Data Lifecycle

**Test Data Preparation**:
```
Before Test Execution:
  1. Download/generate required datasets
  2. Upload to test S3 bucket or PVC
  3. Create test secrets (S3 credentials, SSH keys)
  4. Deploy runtime templates to cluster
  5. Verify cluster resources available
```

**Test Data Cleanup**:
```
After Test Execution:
  1. Delete test TrainJobs (and associated pods)
  2. Remove test secrets
  3. Clean up S3 test data (or retain for debugging)
  4. Delete test runtime templates
  5. Verify no resource leaks (orphaned pods, PVCs)
```

**Cleanup Automation**:
```python
# Example pytest fixture for test cleanup
@pytest.fixture(scope="function")
def test_namespace(kube_client):
    """Create test namespace and clean up after test."""
    namespace = f"mpijob-test-{uuid.uuid4().hex[:8]}"
    kube_client.create_namespace(namespace)
    yield namespace
    kube_client.delete_namespace(namespace, wait=True)
```

---

## 8. Performance and Scale Testing

### 8.1 Performance Benchmarking

**Test Category**: Performance Test
**Testing Approach**: Measure and validate MPIJob performance against baselines

**Test Coverage Areas**:

#### 8.1.1 Job Submission Latency
```
Test Cases:
- TC-PERF-001: Measure time from job submission to launcher pod creation
  - Baseline: <10 seconds (90th percentile)
  - Test: Submit 100 jobs sequentially, measure latency

- TC-PERF-002: Measure time from launcher ready to all workers ready
  - Baseline: <30 seconds for 8 workers
  - Test: Gang scheduling overhead measurement

- TC-PERF-003: Measure time from all pods ready to MPI initialization complete
  - Baseline: <60 seconds for 8 workers (SSH + MPI setup)
  - Test: Measure MPI framework overhead
```

**Success Criteria**:
- 90th percentile latency within baseline
- No regressions from previous releases
- Latency scales linearly with worker count

---

#### 8.1.2 Training Throughput
```
Test Cases:
- TC-PERF-004: Measure training throughput (images/sec) for Horovod ResNet-50
  - Baseline: >90% scaling efficiency vs single GPU
  - Test: 1 GPU vs 8 GPU distributed training

- TC-PERF-005: Measure MPI communication overhead vs training compute time
  - Baseline: MPI overhead <15% of total iteration time
  - Test: Profile MPI collectives vs forward/backward pass

- TC-PERF-006: Measure scaling efficiency from 4 to 16 workers
  - Baseline: >85% scaling efficiency
  - Test: Fixed batch size per worker, measure throughput
```

**Success Criteria**:
- Scaling efficiency meets industry benchmarks
- MPI overhead acceptable for distributed training
- Performance comparable to bare-metal MPI

---

#### 8.1.3 Resource Utilization
```
Test Cases:
- TC-PERF-007: Measure GPU utilization during training
  - Baseline: >80% average GPU utilization
  - Test: Monitor GPU usage throughout job

- TC-PERF-008: Measure network bandwidth utilization
  - Baseline: Network not saturated (<80% bandwidth)
  - Test: Monitor network I/O during MPI collectives

- TC-PERF-009: Measure CPU overhead for launcher pod
  - Baseline: Launcher CPU <10% of single worker CPU
  - Test: Validate launcher resource allocation
```

**Success Criteria**:
- GPUs fully utilized during training
- Network not bottleneck
- Launcher overhead minimal

---

### 8.2 Scale Testing

**Test Category**: Scale Test
**Testing Approach**: Test MPIJob behavior at various scales to validate limits

**Test Coverage Areas**:

#### 8.2.1 Worker Count Scaling
```
Test Cases:
- TC-SCALE-001: Create MPIJob with 2, 4, 8, 16, 32 workers
  - Success: All jobs complete successfully
  - Measure: Time to schedule, time to complete, failure rate

- TC-SCALE-002: Maximum worker count supported
  - Test: Determine upper limit (cluster size, gang scheduling)
  - Document: Recommended max workers (e.g., 64 workers)
```

**Success Criteria**:
- Jobs succeed up to documented maximum
- Clear error messages when exceeding limits
- Performance degradation acceptable up to max

---

#### 8.2.2 Concurrent Job Scaling
```
Test Cases:
- TC-SCALE-003: Submit 10, 20, 50 MPIJobs concurrently
  - Success: All jobs schedule fairly (queue-based)
  - Measure: Queue wait time, scheduling fairness

- TC-SCALE-004: Multiple teams submitting concurrent jobs
  - Test: Team A submits 5 jobs, Team B submits 5 jobs simultaneously
  - Success: Fair scheduling, namespace isolation maintained
```

**Success Criteria**:
- Concurrent jobs don't interfere
- Gang scheduling handles contention gracefully
- Queue positions clear to users

---

#### 8.2.3 Cluster Load Testing
```
Test Cases:
- TC-SCALE-005: Saturate cluster resources with MPIJobs
  - Test: Submit jobs until all GPU/CPU resources allocated
  - Success: Cluster remains stable, no scheduler failures

- TC-SCALE-006: Long-running MPIJobs (24+ hours)
  - Test: Training job runs for extended duration
  - Success: No resource leaks, stable performance throughout
```

**Success Criteria**:
- Cluster stable under full load
- No memory leaks in controller
- Long-running jobs complete successfully

---

## 9. Security and Compliance Testing

### 9.1 Security Testing

**Test Category**: Security Test
**Testing Approach**: Validate security controls and identify vulnerabilities

**Test Coverage Areas**:

#### 9.1.1 Multi-Tenancy Security
```
Test Cases:
- TC-SEC-001: Tenant A cannot access Tenant B's MPIJob via CLI
  - Test: Use Tenant A credentials to access Tenant B namespace
  - Expected: Access denied (RBAC error)

- TC-SEC-002: Tenant A cannot access Tenant B's logs or pod exec
  - Test: Attempt to exec into Tenant B's worker pod
  - Expected: Permission denied

- TC-SEC-003: Tenant A's launcher cannot connect to Tenant B's workers
  - Test: Network policy isolation testing
  - Expected: Connection refused (network policy block)

- TC-SEC-004: Resource quota prevents Tenant A from exhausting cluster
  - Test: Submit oversized jobs to exhaust resources
  - Expected: Quota enforcement blocks excessive allocation
```

**Success Criteria**:
- Zero cross-tenant access violations
- All access attempts logged (audit trail)
- Resource quotas enforced consistently

---

#### 9.1.2 RBAC and Permission Testing
```
Test Cases:
- TC-SEC-005: User with read-only role cannot create MPIJobs
  - Test: Attempt job creation with read-only ServiceAccount
  - Expected: Forbidden error

- TC-SEC-006: User without namespace access cannot list jobs
  - Test: Query jobs in unauthorized namespace
  - Expected: Empty list or error (depending on RBAC config)

- TC-SEC-007: Launcher pod ServiceAccount has minimal permissions
  - Test: Audit launcher ServiceAccount RBAC
  - Expected: Only required permissions (pod exec, secret read)
```

**Success Criteria**:
- Least privilege principle enforced
- RBAC policies documented and auditable
- No privilege escalation vulnerabilities

---

#### 9.1.3 SSH Key Security
```
Test Cases:
- TC-SEC-008: SSH keys ephemeral (deleted with job)
  - Test: Create job, verify SSH secret created, delete job, verify secret deleted
  - Expected: No orphaned SSH keys

- TC-SEC-009: SSH keys not exposed in logs
  - Test: Review launcher and worker logs for leaked credentials
  - Expected: No sensitive data in logs

- TC-SEC-010: SSH keys stored in secrets with proper RBAC
  - Test: Verify secret permissions restrict access to job namespace
  - Expected: Only authorized ServiceAccounts can read secrets
```

**Success Criteria**:
- SSH keys never exposed in logs or events
- Keys deleted immediately after job completion
- Secret access controlled by RBAC

---

#### 9.1.4 Network Security
```
Test Cases:
- TC-SEC-011: Encrypted communication option for sensitive workloads
  - Test: Enable TLS for launcher-worker communication
  - Expected: TLS handshake succeeds, data encrypted

- TC-SEC-012: Network policies prevent unauthorized pod-to-pod traffic
  - Test: Attempt connection from non-MPI pod to worker
  - Expected: Connection blocked by network policy

- TC-SEC-013: Istio sidecar disabled to prevent MPI conflicts
  - Test: Verify annotation disables sidecar injection
  - Expected: No Istio sidecar in launcher or worker pods
```

**Success Criteria**:
- Optional encryption available for compliance
- Network policies enforce least privilege
- Service mesh conflicts avoided

---

### 9.2 Compliance Testing

**Test Category**: Compliance Test
**Testing Approach**: Validate compliance with industry standards and regulations

**Test Coverage Areas**:

#### 9.2.1 Audit Logging
```
Test Cases:
- TC-COMP-001: MPIJob creation event logged with user identity
  - Test: Create job, verify audit log entry
  - Expected: Log includes timestamp, user, namespace, operation

- TC-COMP-002: MPIJob modification event logged
  - Test: Update job, verify audit log entry
  - Expected: Log includes old/new values, user, timestamp

- TC-COMP-003: MPIJob deletion event logged
  - Test: Delete job, verify audit log entry
  - Expected: Log includes user, timestamp, namespace

- TC-COMP-004: Audit logs retained for compliance period (e.g., 90 days)
  - Test: Verify log retention policy
  - Expected: Logs available for minimum retention period
```

**Success Criteria**:
- All lifecycle events audited
- Audit logs immutable and tamper-proof
- Logs retained per compliance requirements

---

#### 9.2.2 FIPS Compliance
```
Test Cases:
- TC-COMP-005: MPI runtime uses FIPS-validated crypto modules
  - Test: Verify OpenMPI/Intel MPI built with FIPS mode
  - Expected: Cryptographic operations FIPS-compliant

- TC-COMP-006: SSH key generation uses FIPS algorithms
  - Test: Inspect SSH key algorithm (RSA 2048+)
  - Expected: FIPS-approved algorithms only
```

**Success Criteria**:
- All cryptographic operations FIPS-compliant
- FIPS mode testable in CI/CD
- Documentation covers FIPS configuration

---

#### 9.2.3 Data Privacy (GDPR, HIPAA)
```
Test Cases:
- TC-COMP-007: Log sanitization prevents PII exposure
  - Test: Submit job with sensitive data, review logs
  - Expected: No PII visible in logs

- TC-COMP-008: Data deletion on request (right to be forgotten)
  - Test: Delete job and associated artifacts
  - Expected: All data removed from cluster and storage
```

**Success Criteria**:
- PII never logged or exposed
- Data deletion complete and verifiable
- Compliance validated by security team

---

## 10. Test Execution Schedule and Coverage

### 10.1 Test Execution Frequency

```
Daily (Developer Commits):
  - Unit tests: All (Go + Python)
  - Duration: <5 minutes
  - Trigger: Local pre-commit hook

Per Pull Request:
  - Unit tests: All
  - Integration tests: Controller + SDK
  - Contract tests: API schemas
  - Duration: <10 minutes
  - Trigger: GitHub Actions on PR

Per Merge to Main:
  - Unit + Integration + Contract tests
  - E2E tests: Smoke suite (critical paths)
  - Duration: <30 minutes
  - Trigger: GitHub Actions on merge

Nightly (Dev Cluster):
  - Full E2E test suite
  - Reference workload validation
  - Performance benchmarks
  - Duration: <2 hours
  - Trigger: Scheduled (2 AM)

Weekly (Staging Cluster):
  - Full E2E + Performance + Scale tests
  - Security scans
  - Compliance checks
  - Duration: <4 hours
  - Trigger: Scheduled (Saturday)

Release Candidate:
  - Full test suite on staging
  - Manual exploratory testing
  - Customer acceptance testing
  - Duration: 1-2 days
  - Trigger: Release branch creation

Production Release:
  - Smoke tests on production
  - Rollback validation
  - Duration: <1 hour
  - Trigger: Deployment to production
```

---

### 10.2 Test Coverage Goals

```
Code Coverage Targets:
  - Controller (Go): >80% line coverage
  - SDK (Python): >85% line coverage
  - Dashboard API (Python): >80% line coverage

Functional Coverage:
  - All functional requirements: 100% coverage
  - User stories (P1): 100% automated E2E tests
  - User stories (P2): 80% automated E2E tests
  - User stories (P3): Manual testing acceptable
  - Edge cases: 100% documented, >80% automated

Test Pyramid:
  - Unit tests: 60% of total tests
  - Integration tests: 30% of total tests
  - E2E tests: 10% of total tests

Test Stability:
  - Flaky test rate: <2%
  - Test failure investigation: <24 hours
  - Test maintenance: <10% of engineering time
```

---

### 10.3 Test Metrics and Reporting

**Key Metrics**:

```
Test Execution Metrics:
  - Total test count: ~500 tests (unit + integration + E2E)
  - Test execution time: <30 minutes (CI), <2 hours (nightly)
  - Test pass rate: >95% (target)
  - Flaky test rate: <2%

Code Coverage Metrics:
  - Controller code coverage: >80%
  - SDK code coverage: >85%
  - Overall coverage: >80%

Quality Metrics:
  - Defect escape rate: <5% (defects found in production)
  - Time to fix test failures: <24 hours
  - Test automation rate: >90% of functional requirements

Performance Metrics:
  - Job submission latency: <10 seconds (90th percentile)
  - Training throughput: >90% scaling efficiency
  - Resource utilization: >80% GPU utilization
```

**Test Reporting**:

```
Daily Reports:
  - Test pass/fail summary
  - New test failures
  - Flaky test identification

Weekly Reports:
  - Test coverage trends
  - Test execution time trends
  - Defect metrics

Release Reports:
  - Test coverage summary
  - Known issues and workarounds
  - Performance benchmark results
  - Security and compliance validation
```

---

## 11. Why Do We Need to Do This?

As a seasoned QA engineer, I need to ask: **Why is comprehensive testing for MPIJob support critical?**

**Business Impact**:
- **Market Differentiation**: OpenShift AI will be first enterprise platform with unified Trainer V2 interface for all frameworks including MPI. Poor quality would undermine this competitive advantage.
- **Enterprise Customers**: 92% of customers require multi-tenancy. Security vulnerabilities or isolation failures would block adoption in regulated industries (finance, healthcare, government).
- **Customer Trust**: MLOps engineers will abandon the platform if MPIJobs fail unpredictably or lack proper diagnostics. Support ticket load would overwhelm support teams.

**Technical Complexity**:
- **Distributed Systems**: MPI involves launcher-worker coordination, gang scheduling, SSH communication, and network policies. Each integration point is a potential failure mode.
- **Multi-Component**: MPIJobs span controller, SDK, Dashboard, storage, networking. Contract failures between components would cause cryptic errors.
- **Kubernetes Orchestration**: Gang scheduling, resource quotas, RBAC, network policies—all must work together correctly.

**Risk Without Testing**:
- **Production Failures**: Incomplete testing leads to customer production incidents, emergency patches, and loss of confidence.
- **Support Burden**: Poor diagnostics and error handling result in 3x higher support ticket volume (industry benchmark).
- **Security Incidents**: Multi-tenancy violations or credential leaks could cause compliance failures and customer churn.

---

## 12. How Am I Going to Test This?

**Testing Approach Summary**:

1. **Local Development**: EnvTest + unit tests for fast iteration (<5 min feedback)
2. **CI/CD Automation**: PR validation with integration + contract tests (<10 min)
3. **Dev Cluster E2E**: Nightly comprehensive tests with reference workloads (<2 hours)
4. **Staging Validation**: Weekly full test suite including scale and performance
5. **Production Smoke Tests**: Post-deployment validation of critical paths

**Key Testing Tools**:
- **Go**: Ginkgo + Gomega + EnvTest (controller testing)
- **Python**: pytest + pytest-bdd (SDK, Dashboard, E2E)
- **Kubernetes**: kubectl, oc CLI, Kubernetes Python client
- **Observability**: Prometheus queries, log analysis, event inspection

**Challenges and Mitigation**:
- **Challenge**: GPU resources expensive for testing
  - **Mitigation**: Use CPU-based tests for functional validation, GPU tests for critical paths only
- **Challenge**: E2E tests take too long
  - **Mitigation**: Parallel execution (pytest-xdist), test prioritization, smoke suite for fast feedback
- **Challenge**: Flaky tests in distributed systems
  - **Mitigation**: Proper timeouts, Eventually() patterns, retry logic, test isolation

---

## 13. Can I Test This Locally?

**Yes, with limitations:**

**What Can Be Tested Locally**:
- **Unit Tests**: All controller and SDK unit tests run without cluster (EnvTest, mocks)
- **Integration Tests**: Controller integration tests with EnvTest (simulated Kubernetes API)
- **Contract Tests**: API schema validation, SDK method contracts

**What Requires Cluster**:
- **E2E Tests**: Full MPIJob lifecycle requires real OpenShift/Kubernetes cluster
- **Network Tests**: SSH communication, network policies require actual pod networking
- **Performance Tests**: GPU resources, multi-node distributed training

**Local Testing Setup**:

```bash
# Install dependencies
go install github.com/onsi/ginkgo/v2/ginkgo@latest
go install github.com/onsi/gomega@latest
pip install poetry
poetry install

# Run unit tests (local, <2 min)
make test-unit

# Run integration tests (local with EnvTest, <5 min)
make test-integration

# Run E2E tests (requires cluster, ~30 min)
export KUBECONFIG=~/.kube/config
make test-e2e

# Run specific test suite
ginkgo -v ./pkg/controller/...  # Go tests
pytest tests/unit/  # Python unit tests
pytest tests/e2e/test_basic_mpijob.py  # Specific E2E test
```

**Developer Workflow**:
1. Make code changes
2. Run unit tests locally (<2 min)
3. Run integration tests locally (<5 min)
4. Push to PR, CI runs full suite (<10 min)
5. Nightly E2E validates on real cluster

---

## 14. Can You Provide Me Details About...

### 14.1 Test Infrastructure Requirements

**I need the following to automate MPIJob testing:**

**Cluster Access**:
- **Dev Cluster**: OpenShift 4.12+ with 4-8 GPU nodes (NVIDIA A100 or V100)
- **CI Cluster**: Ephemeral Kind/Minikube for PR validation (CPU-only acceptable)
- **Staging Cluster**: Production-mirror for release validation

**Credentials and Permissions**:
- **Cluster Admin**: For creating ClusterTrainingRuntime, RBAC policies
- **Namespace Admin**: For creating test namespaces, resource quotas
- **Service Account**: For CI/CD automation (create jobs, access logs)

**Storage**:
- **S3 Bucket**: For test datasets and model artifacts (MinIO or AWS S3)
- **OCS/PVC**: For persistent storage testing
- **Container Registry**: For test images (Quay, Docker Hub, or internal registry)

**Monitoring and Observability**:
- **Prometheus**: For metrics collection and validation
- **Grafana**: For performance dashboards
- **OpenShift Logging**: For log aggregation testing

---

### 14.2 Test Data Requirements

**Training Datasets**:
- **MNIST**: Small dataset for quick validation (~10MB)
- **CIFAR-10**: Medium dataset for functional testing (~200MB)
- **ImageNet subset**: Large dataset for performance testing (~10GB)
- **Synthetic data**: Generated programmatically for scale tests

**Container Images**:
- **Horovod**: `horovod/horovod:latest` with PyTorch/TensorFlow
- **Intel MPI**: Custom image with Intel MPI + training code
- **Test images**: Custom images for failure scenario testing

**Test Fixtures**:
- **TrainJob YAMLs**: Pre-defined manifests for common scenarios
- **Runtime templates**: Horovod, Intel MPI, OpenMPI configurations
- **Network policies**: Test policies for security validation

---

### 14.3 CI/CD Integration Requirements

**GitHub Actions Workflow**:

```yaml
name: MPIJob Test Suite

on:
  pull_request:
    branches: [main]
  push:
    branches: [main]
  schedule:
    - cron: '0 2 * * *'  # Nightly at 2 AM

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run Go unit tests
        run: make test-unit
      - name: Run Python unit tests
        run: poetry run pytest tests/unit/

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup EnvTest
        run: make setup-envtest
      - name: Run integration tests
        run: make test-integration

  e2e-tests:
    runs-on: self-hosted  # Requires cluster access
    if: github.event_name == 'push' || github.event_name == 'schedule'
    steps:
      - uses: actions/checkout@v3
      - name: Run E2E tests
        run: poetry run pytest tests/e2e/ --cluster=${{ secrets.TEST_CLUSTER_URL }}
```

---

### 14.4 Test Reporting and Dashboards

**Test Result Dashboards**:
- **CI Dashboard**: GitHub Actions or OpenShift Pipelines UI
- **Test Coverage**: Codecov or SonarQube for code coverage trends
- **Performance Metrics**: Grafana dashboards for benchmark tracking
- **Flaky Test Detection**: Test analytics platform (e.g., BuildPulse, Flaky Test Tracker)

**Key Metrics to Track**:
- Test pass rate over time
- Test execution duration trends
- Code coverage trends
- Flaky test identification
- Performance benchmark regression detection

---

## 15. Summary and Recommendations

### 15.1 Testing Priorities

**Phase 1 (MVP - Weeks 1-4)**:
1. Unit tests for controller reconciliation logic
2. Integration tests with EnvTest (gang scheduling, pod creation)
3. Basic E2E tests (create job, monitor, delete)
4. Contract tests for SDK and Dashboard APIs

**Phase 2 (Post-MVP - Weeks 5-8)**:
1. Comprehensive E2E test suite (all user stories)
2. Reference workload validation (Horovod, Intel MPI)
3. Failure scenario testing (resource constraints, network issues)
4. Multi-tenancy and security testing

**Phase 3 (Hardening - Weeks 9-12)**:
1. Performance and scale testing
2. Long-running job stability tests
3. Observability and logging validation
4. Compliance testing (FIPS, audit logs)

---

### 15.2 Key Success Metrics

**Test Coverage Goals**:
- >80% code coverage for controller and SDK
- 100% functional requirement coverage
- >90% test automation rate

**Quality Goals**:
- <2% flaky test rate
- <5% defect escape rate
- <24 hours to fix test failures

**Performance Goals**:
- Job submission latency <10 seconds (90th percentile)
- Training throughput >90% scaling efficiency
- Test execution <30 minutes (CI), <2 hours (nightly)

---

### 15.3 Final Recommendations

1. **Invest in Test Infrastructure Early**: Set up EnvTest, CI/CD, and test clusters in Week 1 to enable parallel development and testing.

2. **Test Pyramid Discipline**: Prioritize unit tests (fast, reliable) over E2E tests (slow, flaky). Aim for 60% unit, 30% integration, 10% E2E.

3. **Reference Workloads**: Use industry-standard Horovod and Intel MPI examples for functional validation and performance benchmarking.

4. **Failure Scenario Focus**: Distributed systems fail in complex ways. Dedicate 30% of testing effort to failure scenarios and error handling.

5. **Multi-Tenancy Non-Negotiable**: Security testing cannot be deferred. Any multi-tenancy violation blocks enterprise adoption.

6. **Observability First-Class**: Logs, metrics, and events are as important as functionality. Test observability thoroughly.

7. **Automate Everything**: Manual testing doesn't scale. Automate unit, integration, contract, and E2E tests from Day 1.

8. **Continuous Validation**: Run tests on every PR (fast feedback) and nightly (comprehensive validation). Don't wait for release cycle to find issues.

---

## Appendix A: Test Case Summary

### Total Test Cases by Category

```
Unit Tests: ~150 test cases
  - Controller: 50 tests
  - SDK: 40 tests
  - Webhooks: 20 tests
  - Utilities: 40 tests

Integration Tests: ~80 test cases
  - Controller + EnvTest: 30 tests
  - SDK + Cluster: 20 tests
  - Dashboard + API: 15 tests
  - Network + Security: 15 tests

Contract Tests: ~30 test cases
  - CRD Schema: 10 tests
  - Dashboard API: 10 tests
  - SDK Client: 10 tests

End-to-End Tests: ~50 test cases
  - User workflows: 20 tests
  - Reference workloads: 5 tests
  - Failure scenarios: 15 tests
  - Observability: 10 tests

Performance Tests: ~15 test cases
  - Latency: 5 tests
  - Throughput: 5 tests
  - Scale: 5 tests

Security Tests: ~25 test cases
  - Multi-tenancy: 10 tests
  - RBAC: 8 tests
  - SSH keys: 4 tests
  - Network: 3 tests

Total: ~350 test cases
```

---

## Appendix B: Testing Tools and Frameworks

### Go Testing Stack
- **Ginkgo v2**: BDD test framework
- **Gomega**: Matcher/assertion library
- **EnvTest**: Kubernetes API simulation
- **controller-runtime**: Kubernetes controller libraries
- **golangci-lint**: Code quality and linting

### Python Testing Stack
- **pytest**: Test framework
- **pytest-bdd**: Behavior-driven development
- **pytest-xdist**: Parallel test execution
- **poetry**: Dependency management
- **kubernetes**: Kubernetes Python client
- **requests**: HTTP client for API testing
- **pylint**: Code quality

### Infrastructure and CI/CD
- **GitHub Actions**: CI/CD automation
- **Kind**: Kubernetes in Docker (local testing)
- **OpenShift Pipelines**: CI/CD for OpenShift
- **Docker/Podman**: Container runtime
- **MinIO**: S3-compatible storage for testing

### Observability and Monitoring
- **Prometheus**: Metrics collection
- **Grafana**: Visualization and dashboards
- **OpenShift Logging**: Log aggregation
- **kubectl/oc**: Kubernetes CLI tools

---

## Appendix C: Reference Workloads

### Horovod PyTorch MNIST
```python
# Example reference workload for testing
# Source: https://github.com/horovod/horovod/blob/master/examples/pytorch/pytorch_mnist.py

import torch
import horovod.torch as hvd

hvd.init()
torch.cuda.set_device(hvd.local_rank())

# Training loop...
# Expected: >98% accuracy, completes in <10 minutes
```

### Intel MPI TensorFlow Benchmark
```bash
# Example Intel MPI benchmark command
# Expected: >500 images/sec aggregate throughput

mpirun -np 32 -ppn 4 \
  python train.py \
  --model=resnet50 \
  --batch_size=32 \
  --num_batches=100
```

### Horovod TensorFlow BERT Fine-tuning
```python
# Example BERT fine-tuning workload
# Expected: F1 score >88%, completes in <2 hours

import horovod.tensorflow as hvd
from transformers import TFBertForSequenceClassification

hvd.init()
# Fine-tuning loop...
```

---

**Document End**

This testing strategy provides comprehensive coverage for MPIJob implementation in OpenShift AI. For questions or clarifications, please contact the QA team.
