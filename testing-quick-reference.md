# MPIJob Testing Quick Reference Guide

**Purpose**: Quick reference for developers and QA engineers working on MPIJob testing
**Last Updated**: 2025-10-28

---

## Test Categories at a Glance

```
┌─────────────────────────────────────────────────────────────┐
│                     TEST PYRAMID                            │
│                                                             │
│                       /\                                    │
│                      /E2E\         10% (~50 tests)         │
│                     /------\       User workflows          │
│                    /  INT   \      30% (~80 tests)         │
│                   /----------\     Component integration   │
│                  /   UNIT     \    60% (~150 tests)        │
│                 /--------------\   Fast feedback           │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

---

## Testing Frameworks by Language

### Go (Controller Testing)

```go
// Unit Test Example
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("TrainJob Controller", func() {
    Context("when creating MPIJob", func() {
        It("should create launcher pod", func() {
            // Test implementation
            Eventually(func() error {
                // Check pod exists
            }).Should(Succeed())
        })
    })
})
```

**Tools**:
- Ginkgo v2 (BDD framework)
- Gomega (matchers)
- EnvTest (simulated K8s API)
- controller-runtime/fake (mock clients)

**Run Tests**:
```bash
# Unit tests
make test-unit

# Integration tests (with EnvTest)
make test-integration

# Specific package
ginkgo -v ./pkg/controller/trainjob/
```

---

### Python (SDK, Dashboard, E2E)

```python
# Unit Test Example
import pytest
from kubeflow.training import TrainingClient

def test_create_train_job():
    """Test SDK creates correct TrainJob manifest."""
    client = TrainingClient()
    job = client.create_train_job(
        name="test-mpijob",
        runtime_ref="mpi-horovod",
        num_nodes=4
    )
    assert job.metadata.name == "test-mpijob"
    assert job.spec.trainer.num_nodes == 4
```

**Tools**:
- pytest (test framework)
- pytest-bdd (behavior-driven development)
- pytest-xdist (parallel execution)
- kubernetes Python client

**Run Tests**:
```bash
# Install dependencies
poetry install

# Unit tests
poetry run pytest tests/unit/

# Integration tests (requires cluster)
poetry run pytest tests/integration/ --kubeconfig=$KUBECONFIG

# E2E tests
poetry run pytest tests/e2e/ --cluster=https://api.example.com

# Specific test
poetry run pytest tests/e2e/test_basic_mpijob.py::test_create_via_sdk
```

---

## Test Checklist by Priority

### P0 - MVP Blockers (Must Have)

- [ ] **Gang Scheduling**
  - [ ] All pods scheduled together (TC-INT-006)
  - [ ] Timeout on insufficient resources (TC-INT-007)
  - [ ] Clear error messages on failure (TC-INT-008)

- [ ] **Multi-Tenancy Isolation**
  - [ ] RBAC prevents cross-namespace access (TC-INT-010)
  - [ ] Resource quotas enforced (TC-INT-013)
  - [ ] Network policies block cross-tenant traffic (TC-INT-026)

- [ ] **SSH Communication**
  - [ ] Launcher connects to all workers (TC-INT-022)
  - [ ] SSH keys generated and distributed (TC-INT-023)
  - [ ] MPI initialization successful (TC-INT-025)

- [ ] **Basic Lifecycle**
  - [ ] Create job via SDK/CLI/Dashboard (TC-E2E-001)
  - [ ] Monitor status updates (TC-E2E-001)
  - [ ] Access logs (TC-E2E-001)
  - [ ] Delete job (TC-CLI-005)

### P1 - High Priority (Post-MVP)

- [ ] **Failure Scenarios**
  - [ ] OOM error handling (TC-E2E-015)
  - [ ] Gang scheduling timeout (TC-E2E-014)
  - [ ] Network policy issues (TC-E2E-017)
  - [ ] Image pull failures (TC-E2E-020)

- [ ] **Observability**
  - [ ] Metrics collection (TC-E2E-023 to TC-E2E-026)
  - [ ] Log aggregation (TC-E2E-027 to TC-E2E-030)
  - [ ] Event tracking (TC-E2E-031 to TC-E2E-034)

- [ ] **Reference Workloads**
  - [ ] Horovod PyTorch MNIST (TC-E2E-011)
  - [ ] Intel MPI TensorFlow (TC-E2E-012)
  - [ ] Horovod BERT fine-tuning (TC-E2E-013)

### P2 - Medium Priority

- [ ] **Dashboard UI**
  - [ ] Job creation wizard (TC-UI-002)
  - [ ] Topology view (TC-UI-004)
  - [ ] Log viewer (TC-UI-005)
  - [ ] Clone job feature (TC-E2E-002)

- [ ] **Performance Benchmarks**
  - [ ] Job submission latency <10s (TC-PERF-001)
  - [ ] Gang scheduling <30s (TC-PERF-002)
  - [ ] Scaling efficiency >90% (TC-PERF-004)

### P3 - Lower Priority

- [ ] **Migration**
  - [ ] Legacy MPIJob conversion (TC-E2E-007)
  - [ ] Migration documentation (FR-MIG-001)

- [ ] **Scale Testing**
  - [ ] 32+ worker jobs (TC-SCALE-001)
  - [ ] 50+ concurrent jobs (TC-SCALE-004)
  - [ ] 24+ hour jobs (TC-SCALE-006)

---

## Common Test Patterns

### Pattern 1: Eventually() for Async Operations

```go
// GO - Wait for pod to be ready
Eventually(func() bool {
    pod := &corev1.Pod{}
    err := k8sClient.Get(ctx, podKey, pod)
    if err != nil {
        return false
    }
    return pod.Status.Phase == corev1.PodRunning
}, timeout, interval).Should(BeTrue())
```

```python
# PYTHON - Wait for job to complete
import time
from kubernetes import client

def wait_for_job_completion(job_name, namespace, timeout=300):
    """Wait for TrainJob to complete."""
    api = client.CustomObjectsApi()
    start_time = time.time()

    while time.time() - start_time < timeout:
        job = api.get_namespaced_custom_object(
            group="kubeflow.org",
            version="v2alpha1",
            namespace=namespace,
            plural="trainjobs",
            name=job_name
        )
        status = job.get("status", {}).get("phase")
        if status in ["Succeeded", "Failed"]:
            return status
        time.sleep(5)

    raise TimeoutError(f"Job {job_name} did not complete in {timeout}s")
```

---

### Pattern 2: Test Fixtures for Cleanup

```python
# PYTHON - pytest fixture for namespace cleanup
import pytest
from kubernetes import client, config

@pytest.fixture(scope="function")
def test_namespace():
    """Create test namespace and clean up after test."""
    config.load_kube_config()
    v1 = client.CoreV1Api()

    # Create namespace
    namespace = f"mpijob-test-{uuid.uuid4().hex[:8]}"
    body = client.V1Namespace(
        metadata=client.V1ObjectMeta(name=namespace)
    )
    v1.create_namespace(body=body)

    yield namespace

    # Cleanup
    v1.delete_namespace(name=namespace)
```

---

### Pattern 3: Mock Kubernetes API

```python
# PYTHON - Mock Kubernetes API responses
from unittest.mock import MagicMock, patch

def test_sdk_handles_api_errors():
    """Test SDK gracefully handles API errors."""
    with patch('kubernetes.client.CustomObjectsApi') as mock_api:
        # Mock API error response
        mock_api.return_value.create_namespaced_custom_object.side_effect = \
            client.ApiException(status=400, reason="Bad Request")

        # Test SDK error handling
        client = TrainingClient()
        with pytest.raises(ValueError, match="Invalid job specification"):
            client.create_train_job(name="invalid-job", num_nodes=-1)
```

---

## Test Execution Frequency

```
┌──────────────────┬─────────────────────────┬─────────────┐
│ Trigger          │ Test Suite              │ Duration    │
├──────────────────┼─────────────────────────┼─────────────┤
│ Pre-Commit       │ Unit tests (local)      │ <2 min      │
│ Pull Request     │ Unit + Integration      │ <10 min     │
│ Merge to Main    │ Unit + Int + E2E Smoke  │ <30 min     │
│ Nightly          │ Full E2E + Performance  │ <2 hours    │
│ Weekly           │ Full Suite + Scale      │ <4 hours    │
│ Release          │ Full Suite + Manual     │ 1-2 days    │
└──────────────────┴─────────────────────────┴─────────────┘
```

---

## Test Coverage Targets

```
Code Coverage:
  ✓ Controller (Go):      >80%
  ✓ SDK (Python):         >85%
  ✓ Dashboard API:        >80%

Functional Coverage:
  ✓ Functional Reqs:      100%
  ✓ User Stories (P1):    100% automated
  ✓ User Stories (P2):    80% automated
  ✓ Edge Cases:           80% automated

Test Stability:
  ✓ Flaky Test Rate:      <2%
  ✓ Failure Resolution:   <24 hours
```

---

## Debugging Tips

### Failed Unit Test

```bash
# Run specific test with verbose output
ginkgo -v --focus="should create launcher pod" ./pkg/controller/

# Check test logs
poetry run pytest tests/unit/test_sdk.py -v -s --log-cli-level=DEBUG
```

### Failed Integration Test

```bash
# Check EnvTest logs
export USE_EXISTING_CLUSTER=false
make test-integration 2>&1 | tee test.log

# Verify EnvTest setup
kubebuilder envtest list
```

### Failed E2E Test

```bash
# Check job status
oc get trainjob test-mpijob -o yaml

# Check pod status
oc get pods -l job-name=test-mpijob

# Check logs
oc logs -l job-name=test-mpijob --all-containers

# Check events
oc get events --sort-by='.lastTimestamp'
```

---

## Common Issues and Fixes

### Issue 1: EnvTest Fails to Start

**Symptom**: `Error: unable to start test environment`

**Fix**:
```bash
# Install/update EnvTest binaries
make envtest
export KUBEBUILDER_ASSETS="$(./bin/setup-envtest use -p path)"
```

---

### Issue 2: E2E Test Timeout

**Symptom**: Test hangs or times out waiting for job completion

**Check**:
1. Gang scheduling stuck? `oc describe trainjob <name>`
2. Image pull issues? `oc get pods -l job-name=<name>`
3. Resource constraints? `oc describe node`

**Fix**:
- Increase timeout in test
- Check cluster has sufficient resources
- Verify image accessible in registry

---

### Issue 3: Flaky Network Tests

**Symptom**: SSH connection tests pass sometimes, fail other times

**Fix**:
```python
# Use proper retry logic
def wait_for_ssh_connection(pod_name, timeout=60):
    """Wait for SSH to be ready with retries."""
    for i in range(timeout):
        result = subprocess.run(
            ["oc", "exec", pod_name, "--", "ssh", "worker-0", "echo", "ok"],
            capture_output=True
        )
        if result.returncode == 0:
            return True
        time.sleep(1)
    return False
```

---

## Test Resources

### Test Clusters

```
Development Cluster:
  URL: https://api.dev.openshift.com
  GPU Nodes: 4x A100
  Use: Daily E2E testing

CI Cluster:
  Type: Kind (ephemeral)
  Resources: CPU-only
  Use: PR validation

Staging Cluster:
  URL: https://api.staging.openshift.com
  GPU Nodes: 8x A100
  Use: Release validation
```

### Test Data Locations

```
S3 Bucket: s3://openshift-ai-test-data/mpijob/
  - datasets/mnist/
  - datasets/cifar10/
  - reference-workloads/horovod/
  - reference-workloads/intel-mpi/

Container Registry: quay.io/openshift-ai/test-images/
  - horovod-pytorch:latest
  - intel-mpi-tensorflow:latest
  - test-failure-scenarios:latest
```

---

## Quick Commands

### Run All Tests Locally

```bash
# Unit tests only (fast)
make test-unit

# Integration tests (requires EnvTest)
make test-integration

# E2E tests (requires cluster)
export KUBECONFIG=~/.kube/config
make test-e2e

# Full test suite
make test-all
```

### Run Specific Test Category

```bash
# Go unit tests for controller
ginkgo ./pkg/controller/...

# Python SDK tests
poetry run pytest tests/sdk/

# E2E user workflow tests
poetry run pytest tests/e2e/test_user_workflows.py

# Performance tests
poetry run pytest tests/performance/
```

### Check Test Coverage

```bash
# Go coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Python coverage
poetry run pytest --cov=src --cov-report=html tests/
open htmlcov/index.html
```

---

## Test Naming Conventions

```
Format: test_<component>_<action>_<expected_result>

Examples:
  ✓ test_controller_creates_launcher_pod_successfully
  ✓ test_sdk_handles_invalid_worker_count_with_error
  ✓ test_gang_scheduler_times_out_on_insufficient_resources
  ✓ test_dashboard_displays_job_topology_view

Ginkgo (Go):
  Describe("TrainJob Controller")
    Context("when creating new job")
      It("should create launcher pod")
      It("should create worker pods")
```

---

## CI/CD Integration

### GitHub Actions Workflow

```yaml
name: MPIJob Tests

on: [pull_request, push]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run unit tests
        run: make test-unit

  integration-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run integration tests
        run: make test-integration

  e2e-tests:
    runs-on: self-hosted  # Requires cluster
    if: github.event_name == 'push'
    steps:
      - uses: actions/checkout@v3
      - name: Run E2E tests
        run: make test-e2e
```

---

## Contact and Support

**QA Team Lead**: Neil (Principal QA Engineer)
**Slack Channel**: #openshift-ai-testing
**Test Infrastructure Issues**: File in Jira project RHOAI-TEST
**Test Cluster Access**: Contact DevOps team

---

## Additional Resources

- **Full Testing Strategy**: `testing-strategy-mpijob.md`
- **Research Summary**: `testing-research-summary.md`
- **Feature Spec**: `specs/001-mpijob-trainer-v2-support/spec.md`
- **Functional Requirements**: `specs/001-mpijob-trainer-v2-support/spec.md#requirements`

---

**Last Updated**: 2025-10-28
**Version**: 1.0
