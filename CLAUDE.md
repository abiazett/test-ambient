# test-ambient Development Guidelines

Auto-generated from all feature plans. Last updated: 2025-10-29

## Active Technologies

### 001-mpijob-kubeflow-integration: MPIJob Support for OpenShift AI

**Languages & Versions:**
- Backend/Operator: Go 1.21+
- CLI: Go 1.21+ (kubectl plugin)
- SDK: Python 3.9+
- Dashboard: React 18+, TypeScript 5.0+, PatternFly 5

**Primary Dependencies:**
- KubeFlow Training Operator v2.0+
- Kubernetes API (custom resources)
- Open MPI 4.1+
- Kueue (gang scheduling)
- Prometheus/Grafana
- NVIDIA GPU Operator 23.9+

**Target Platform:**
- RedHat OpenShift 4.14+ (Kubernetes 1.27+)
- GPU nodes with NVIDIA GPUs
- On-premises and public cloud

**Testing:**
- Go: go test
- Python: pytest (80% coverage target)
- JavaScript/TypeScript: Jest/React Testing Library
- Integration: Kubernetes test environment with GPUs

**Performance Goals:**
- ≥85% scaling efficiency for 8-node MPIJobs
- <5s job submission latency
- <2s Dashboard UI page load

**Key Constraints:**
- Gang scheduling required (all workers start atomically)
- SSH-based MPI spawning (MVP)
- Standard TCP/IP networking (MVP)
- Upstream KubeFlow compatibility (no forking)
- FIPS 140-2 compliance for SSH keys

## Project Structure

This is a distributed system spanning multiple repositories:

```
# Training Operator (upstream KubeFlow, RHOAI packaging)
training-operator/
├── pkg/
│   ├── apis/kubeflow.org/v2beta1/
│   ├── controller/mpijob/
│   └── common/
├── manifests/
│   ├── crds/
│   └── rbac/
└── tests/

# CLI (odh-cli or oc-odh plugin)
cli/
├── cmd/mpijob/
├── pkg/client/
└── tests/

# Python SDK
sdk/
├── openshift_ai_sdk/training/
├── tests/
└── examples/

# Dashboard UI
dashboard/
├── backend/app/routes/training.py
└── frontend/src/pages/TrainingJobs/

# Documentation
docs/
├── quickstart/
├── reference/
└── troubleshooting/

# Observability
monitoring/
├── prometheus/
└── grafana/dashboards/
```

## Commands

### CLI Commands
```bash
# Create MPIJob
oc odh mpijob create <name> --workers 4 --gpu 1 --image <img> --command "python train.py"
oc odh mpijob create -f mpijob.yaml

# Manage MPIJobs
oc odh mpijob get [name]
oc odh mpijob describe <name>
oc odh mpijob delete <name>

# View logs
oc odh mpijob logs <name> [--launcher|--worker N] [--follow]
```

### Python SDK
```python
from openshift_ai_sdk.training import MPIJob, ResourceRequirements

# Create and submit job
job = MPIJob(
    name="my-job",
    namespace="my-project",
    workers=4,
    image="my-image:latest",
    command=["python", "train.py"],
    resources=ResourceRequirements(cpu=8, memory="64Gi", gpu=1)
)
job.create()

# Monitor
status = job.get_status()
job.wait_for_completion(timeout=3600)
logs = job.get_logs()
job.delete()
```

### Testing
```bash
# Go tests (operator/CLI)
go test ./pkg/...

# Python tests (SDK)
pytest tests/ --cov=openshift_ai_sdk --cov-report=html

# Integration tests
pytest tests/integration/ -m integration

# Performance benchmarks
./tests/performance/scaling_efficiency.sh
```

## Code Style

### Go
- Follow standard Go conventions
- Use `gofmt` for formatting
- Run `go vet` and `golangci-lint`
- 80% test coverage target

### Python
- PEP 8 style guide
- Type hints required
- Docstrings for all public APIs
- 75-80% test coverage target

### TypeScript/React
- PatternFly 5 components
- Follow RHOAI Dashboard patterns
- 80-90% test coverage target

### YAML/Kubernetes
- Use `kubeflow.org/v2beta1` API group
- Always include Istio disable annotation: `sidecar.istio.io/inject: "false"`
- Follow KubeFlow Trainer V2 schema conventions

## Architecture Principles

1. **Upstream-First**: Leverage KubeFlow Trainer V2 without forking
2. **Enterprise-Grade**: Security, multi-tenancy, observability as first-class concerns
3. **Hybrid Cloud Native**: Identical experience across on-prem and cloud
4. **Unified Experience**: Seamless integration with RHOAI Dashboard, CLI, SDK

## Key Architectural Decisions

- **ADR-001**: Single Training Operator manages all job types (MPIJob, PyTorchJob, TFJob)
- **ADR-002**: Use `kubeflow.org/v2beta1` API (GA v2 planned 2026)
- **ADR-003**: SSH for MVP, evaluate PMIx post-MVP
- **ADR-004**: Kueue required for gang scheduling (prevents deadlock)
- **ADR-005**: MPIJob pods have no Kubernetes API access by default

## Common Pitfalls

⚠️ **Critical Issues to Avoid:**
1. **Istio Sidecar**: Must disable with `sidecar.istio.io/inject: "false"` annotation
2. **Pod Cleanup**: Use `cleanPodPolicy: Running` to terminate workers after launcher completes
3. **Worker Count**: Minimum 2 workers required for distributed training
4. **Gang Scheduling**: Without Kueue, resource deadlocks occur with partial allocations
5. **Learning Rate**: Scale LR by number of workers (e.g., `lr * hvd.size()`)

## Recent Changes
- 001-mpijob-kubeflow-integration: Added MPIJob support with KubeFlow Trainer V2 integration

<!-- MANUAL ADDITIONS START -->
<!-- MANUAL ADDITIONS END -->
