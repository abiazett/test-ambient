# Implementation Plan: MPIJob Support for OpenShift AI

**Branch**: `001-mpijob-kubeflow-integration` | **Date**: 2025-10-29 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-mpijob-kubeflow-integration/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Integrate KubeFlow MPI Operator v2 into RedHat OpenShift AI to enable distributed MPI-based training for ML models. This feature provides three interfaces (Dashboard UI, Python SDK, CLI) for creating, monitoring, and managing MPIJobs with enterprise-grade security, multi-tenancy, and observability. The implementation follows an upstream-first approach, leveraging KubeFlow Training Operator v2 as a managed component within the RHOAI operator lifecycle.

## Technical Context

**Language/Version**:
- Backend/Operator: Go 1.21+ (Training Operator v2)
- CLI: Go 1.21+ (kubectl plugin pattern)
- SDK: Python 3.9+ (kubernetes-client library)
- Dashboard: React 18+, TypeScript 5.0+, PatternFly 5

**Primary Dependencies**:
- KubeFlow Training Operator v2.0+ (manages MPIJob, PyTorchJob, TFJob)
- Kubernetes API server (custom resource management)
- Open MPI 4.1+ (inter-process communication runtime)
- Kueue (gang scheduling and fair-share resource allocation)
- Prometheus/Grafana (observability)
- NVIDIA GPU Operator 23.9+ (GPU resource management)
- OpenShift SDN/OVN-Kubernetes (pod-to-pod networking)

**Storage**:
- PVC/ODF/NFS/Ceph for training data and model checkpoints
- S3-compatible object storage (optional, for checkpoint persistence)
- Kubernetes Secrets (SSH keys for MPI communication)
- ConfigMaps (MPI hostfiles)

**Testing**:
- Go: go test for operator and CLI
- Python: pytest for SDK
- JavaScript/TypeScript: Jest/React Testing Library for Dashboard
- Integration: Kubernetes test environment with GPU nodes
- Performance: ResNet-50, BERT, GPT-7B scaling efficiency benchmarks

**Target Platform**:
- RedHat OpenShift 4.14+ (Kubernetes 1.27+)
- On-premises and public cloud (AWS, Azure, GCP)
- GPU nodes with NVIDIA GPUs (T4, V100, A100, H100)

**Project Type**:
- Distributed system (operator + multi-interface)
- Web application (Dashboard UI)
- Library/SDK (Python package)
- CLI tool (kubectl plugin)

**Performance Goals**:
- ≥85% scaling efficiency for 8-node MPIJobs (speedup ≥6.8× vs. single node)
- ≥75% scaling efficiency for 16-node MPIJobs
- <5 seconds job submission latency (CLI/SDK)
- <2 seconds Dashboard UI page load for 100 jobs
- <30 seconds worker pod startup (excluding image pull)
- Dashboard real-time updates every 5 seconds

**Constraints**:
- Gang scheduling required (all workers must start atomically)
- SSH-based MPI process spawning (MVP; PMIx evaluation post-MVP)
- Standard TCP/IP networking (MVP; RDMA post-MVP)
- Upstream KubeFlow compatibility (no forking)
- RBAC and namespace isolation for multi-tenancy
- FIPS 140-2 compliance for SSH keys (RSA 4096-bit)

**Scale/Scope**:
- Support 2-100 workers per MPIJob
- Support 50+ concurrent MPIJobs across namespaces
- Platform tested with 512 GPU clusters
- Dashboard UI handles 100+ training jobs without degradation
- Single Training Operator manages all job types (unified controller)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

**Constitution Status**: No project constitution exists yet. Constitution templates are available at `.specify/memory/constitution.md` but are not instantiated for this project.

**Architectural Principles (from ARCHITECTURE.md)**:
1. **Upstream-First**: Leverage KubeFlow Trainer V2 and MPI Operator v2 without forking ✓
2. **Enterprise-Grade**: Security, multi-tenancy, and observability as first-class concerns ✓
3. **Hybrid Cloud Native**: Identical experience across on-premises and cloud deployments ✓
4. **Unified Experience**: Seamless integration with existing RHOAI Dashboard, CLI, and SDK patterns ✓

**Gate: Complexity Justification**
- **Status**: PASS
- **Rationale**: This feature integrates multiple components (Operator, CLI, SDK, Dashboard) but each is necessary:
  - Operator: Required for Kubernetes resource management
  - CLI: Required for automation and CI/CD integration (FR-006)
  - SDK: Required for Jupyter notebook integration (FR-011)
  - Dashboard: Required for UI-driven job creation (FR-015)
  - All components follow established RHOAI patterns

**Gate: Testing Requirements**
- **Status**: REQUIRES CLARIFICATION → Will be addressed in Phase 0
- **Open Questions**:
  - What is the minimum test coverage threshold for each component?
  - What is the test environment specification (GPU types, cluster size)?
  - Are there specific compliance/security test requirements?

**Gate: Upstream Compatibility**
- **Status**: PASS
- **Rationale**: Using `kubeflow.org/v2beta1` API without modifications, no forking of upstream operators

**Re-evaluation Required**: After Phase 1 design artifacts are complete, re-check compliance with any additional constraints or patterns that emerge.

---

## Post-Design Constitution Re-evaluation

**Date**: 2025-10-29
**Status**: PASS - All gates satisfied

### Gate: Testing Requirements (Previously "NEEDS CLARIFICATION")
- **Status**: RESOLVED → PASS
- **Decision** (from research.md):
  - Kubernetes Operators: 80% unit test coverage
  - Python SDK: 75-80% coverage
  - React Dashboard: 80-90% coverage
  - Multi-layer approach: Unit (80%) → Integration (70%) → E2E (critical paths)
- **Test Environment**: 4 GPU nodes, 8 GPUs total (NVIDIA T4 or better)
- **Compliance**: Standard ISO/IEC 25010 quality model for software testing

### Gate: Architectural Complexity
- **Status**: PASS
- **Justification**: Research validated all architectural decisions
  - Gang scheduling (Kueue): Mandatory to prevent deadlocks (research confirmed)
  - SSH-based MPI: Industry standard, PMIx evaluation post-MVP (research confirmed)
  - Multiple interfaces (Operator, CLI, SDK, Dashboard): Required by functional requirements
  - Distributed repositories: Necessary for upstream contribution and independent lifecycles

### Gate: Performance Targets
- **Status**: PASS
- **Validation** (from research.md):
  - 85% scaling efficiency for 8-node clusters: Achievable with TCP/IP networking
  - ResNet-50 benchmark: 85% efficiency validated in production deployments
  - BERT-Large: 88-90% efficiency with DeepSpeed
  - Network requirements: 10 GbE sufficient for 2-4 nodes, 25-50 GbE for 8-16 nodes

### Gate: Security Compliance
- **Status**: PASS
- **Validation**:
  - FIPS 140-2 compliance: RSA-4096 keys (confirmed standard)
  - RBAC enforcement: Kubernetes-native RBAC patterns (data-model.md)
  - Namespace isolation: NetworkPolicies + ResourceQuotas (data-model.md)
  - Secret management: Ephemeral SSH keys with ownerReferences (contracts/)

### Gate: Upstream Compatibility
- **Status**: PASS
- **Validation**:
  - Using `kubeflow.org/v2beta1` API without modifications (mpijob-crd.yaml)
  - No forking of KubeFlow Training Operator (architecture principle)
  - Conversion webhooks for v1alpha1 → v2beta1 migration (future-proof)

### Gate: Documentation Completeness
- **Status**: PASS
- **Artifacts Generated**:
  - ✓ research.md: Comprehensive technical research with 46 citations
  - ✓ data-model.md: Complete entity definitions with validation rules
  - ✓ contracts/: CRD, Dashboard API (OpenAPI 3.0), SDK API specification
  - ✓ quickstart.md: 30-minute end-to-end tutorial with troubleshooting
  - ✓ plan.md: Implementation plan with technical context and structure

### New Risks Identified
1. **Istio Incompatibility**: Critical pitfall discovered in research
   - **Mitigation**: Documented in quickstart, added to CLAUDE.md, included in CRD examples
2. **Learning Rate Scaling**: Common training accuracy issue
   - **Mitigation**: Documented in quickstart troubleshooting section
3. **Network Bandwidth**: Can be bottleneck for large-scale training
   - **Mitigation**: Performance benchmarking section in research.md, clear recommendations

### Recommendations for Phase 2 (Task Generation)
1. **High Priority**: Implement admission webhook to enforce Istio sidecar annotation
2. **High Priority**: Create integration tests for Kueue gang scheduling
3. **Medium Priority**: Performance benchmarking suite (ResNet-50, BERT)
4. **Medium Priority**: Migration tooling for SageMaker/Azure ML customers
5. **Low Priority**: Distributed tracing integration (OpenTelemetry)

**Final Assessment**: All constitutional gates PASS. Design is implementation-ready.

## Project Structure

### Documentation (this feature)

```
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

**Note**: This feature integrates with multiple existing repositories in the RHOAI ecosystem. Below is the anticipated structure for new code contributions:

```
# Training Operator (upstream KubeFlow, RHOAI fork/packaging)
training-operator/
├── pkg/
│   ├── apis/
│   │   └── kubeflow.org/
│   │       └── v2beta1/           # MPIJob CRD definitions
│   ├── controller/
│   │   └── mpijob/                # MPIJob reconciliation logic
│   └── common/                    # Shared utilities
├── manifests/
│   ├── crds/                      # MPIJob CustomResourceDefinition
│   └── rbac/                      # ServiceAccount, Role, RoleBinding
└── tests/
    ├── unit/
    └── integration/

# CLI (odh-cli or oc-odh plugin)
cli/
├── cmd/
│   └── mpijob/
│       ├── create.go
│       ├── get.go
│       ├── describe.go
│       ├── logs.go
│       └── delete.go
├── pkg/
│   ├── client/                    # Kubernetes API client
│   └── printer/                   # Output formatting
└── tests/
    └── e2e/

# Python SDK (openshift-ai-sdk package)
sdk/
├── openshift_ai_sdk/
│   └── training/
│       ├── mpijob.py              # MPIJob class
│       ├── resources.py           # ResourceRequirements, etc.
│       └── __init__.py
├── tests/
│   ├── unit/
│   └── integration/
└── examples/
    └── jupyter_notebooks/

# Dashboard UI (odh-dashboard repository)
dashboard/
├── backend/
│   └── app/
│       └── routes/
│           └── training.py        # MPIJob REST API endpoints
└── frontend/
    └── src/
        ├── pages/
        │   └── TrainingJobs/
        │       ├── MPIJobList.tsx
        │       ├── MPIJobDetails.tsx
        │       ├── MPIJobCreateWizard.tsx
        │       ├── MPIJobLogs.tsx
        │       └── MPIJobMetrics.tsx
        ├── api/
        │   └── training.ts        # API client
        └── types/
            └── training.ts        # MPIJob TypeScript types

# Documentation
docs/
├── quickstart/
│   └── mpijobs-getting-started.md
├── reference/
│   ├── cli-reference.md
│   ├── sdk-api.md
│   └── mpijob-crd-spec.md
└── troubleshooting/
    └── mpijobs-common-issues.md

# Observability
monitoring/
├── prometheus/
│   └── training-operator-rules.yaml
└── grafana/
    └── dashboards/
        └── mpijob-overview.json
```

**Structure Decision**: This is a distributed system spanning multiple repositories:
1. **Training Operator**: Go-based Kubernetes operator (upstream KubeFlow contribution + RHOAI packaging)
2. **CLI**: Go-based kubectl plugin following `oc odh` command structure
3. **SDK**: Python library distributed via PyPI
4. **Dashboard**: React/TypeScript frontend + Flask backend integration into existing ODH Dashboard
5. **Documentation**: Markdown docs integrated into RHOAI documentation site

Each component has independent development lifecycle but coordinated releases via RHOAI operator.

## Complexity Tracking

*Fill ONLY if Constitution Check has violations that must be justified*

**Status**: No constitutional violations detected. Complexity is justified by requirements:

| Component | Justification | Alternatives Considered |
|-----------|--------------|------------------------|
| Multiple repositories | Each component serves different user interfaces and has independent release cycles | Single monorepo would create coupling and complicate upstream contributions to KubeFlow |
| Four user interfaces (Operator, CLI, SDK, Dashboard) | FR-006 (CLI), FR-011 (SDK), FR-015 (Dashboard) explicitly require each interface | Web-only interface rejected because MLOps teams require CLI/SDK for automation |
| Gang scheduling (Kueue dependency) | MPI training requires all workers to start atomically (ADR-004) | Without gang scheduling, resource deadlocks occur with partial allocations |
| SSH-based MPI spawning | Battle-tested, supported by all MPI implementations | PMIx is emerging but lacks tooling maturity; evaluate post-MVP |

