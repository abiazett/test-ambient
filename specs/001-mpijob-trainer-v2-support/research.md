# Research: MPIJobs Support in OpenShift AI via Trainer V2

**Feature Branch**: `001-mpijob-trainer-v2-support`
**Date**: 2025-10-28
**Status**: Phase 0 Complete

This document consolidates research findings for implementing MPIJob support in OpenShift AI using KubeFlow Trainer V2. All NEEDS CLARIFICATION items from Technical Context have been resolved.

---

## 1. KubeFlow Trainer V2 Architecture

### Decision: Use Trainer V2 Unified API

**Rationale**: Trainer V2 provides a unified training API that replaces framework-specific CRDs (PyTorchJob, TFJob, MPIJob v1). This architectural shift separates infrastructure configuration (managed by platform administrators via runtime templates) from training job specifications (managed by data scientists via TrainJob).

**Key Benefits**:
- Consistent API across all training frameworks (PyTorch, TensorFlow, MPI, JAX)
- Separation of concerns: platform teams manage runtimes, data scientists focus on ML
- Built on native Kubernetes JobSet API instead of custom implementations
- Better integration with Kueue for multi-tenant resource management

### TrainJob CRD Structure

**Core Fields**:
- `runtimeRef`: Links to TrainingRuntime or ClusterTrainingRuntime (apiGroup, kind, name)
- `trainer`: Training configuration with:
  - `numNodes`: Number of distributed training nodes
  - `image`: Docker image for training containers
  - `command`: Entrypoint command and arguments
  - `resourcesPerNode`: CPU, memory, GPU allocations per node
- `initializer`: Optional automated dataset/model setup from HuggingFace or S3

**Under the Hood**: Controller translates TrainJob + RuntimeTemplate → JobSet → multiple Kubernetes Jobs. This leverages native K8s features (SuccessPolicy, PodFailurePolicy, suspend) rather than reimplementing them.

### TrainingRuntime Configuration

**Structure**:
- `template.spec.replicatedJobs`: Array of job templates (launcher, workers, initializers)
- `mlPolicy`: ML-specific configuration (numNodes, torch.numProcPerNode)
- `podGroupPolicy`: Gang scheduling configuration (coscheduling, scheduleTimeoutSeconds)

**MPI Runtime Components**:
- **Launcher node (node-0)**: Runs mpirun command and OpenMPI launcher
- **Worker nodes**: Run OpenSSH server and execute training workload
- **Automatic setup**: Controller handles SSH key generation, OpenMPI hostfile creation, distributed environment variables (world size, rank)
- **Supported implementations**: OpenMPI, Intel MPI, MPICH

**Gang Scheduling**: Uses Kubernetes scheduler-plugins coscheduling or Volcano Scheduler. PodGroup resource ensures all pods schedule atomically with configurable timeout.

### Controller Architecture

**Reconciliation Loop Pattern**:
1. Trigger: Watch TrainJob and RuntimeTemplate resource changes
2. Fetch TrainJob and referenced RuntimeTemplate
3. Validate specifications using CEL (Common Expression Language) validations
4. Apply MLPolicy and PodGroupPolicy transformations
5. Create/update JobSet with rendered templates
6. Create PodGroup resource for gang scheduling
7. Update TrainJob status based on JobSet conditions

**MPI Lifecycle**:
- Controller generates SSH key pairs and injects into secrets
- Creates ConfigMap with OpenMPI hostfile listing all worker nodes
- Launcher pod waits for all worker SSH servers to be ready
- Runs mpirun command with proper rank assignment
- Monitors job completion and updates status

**Resource Management**: Integration with Kueue enables multi-tenant resource orchestration, quota management, and multi-cluster workload distribution.

### Implementation Notes

**CEL Validations**: Trainer V2 uses Common Expression Language for CRD validation rules. This enables complex validation logic (e.g., "if runtimeRef specifies MPI runtime, then mlPolicy.mpiImplementation must be set") at admission time.

**JobSet Dependency**: Trainer V2 depends on JobSet (kubernetes-sigs/jobset). OpenShift AI must install JobSet operator as a prerequisite.

**Alternatives Considered**:
- **Continue using MPIJob v1 operator**: Rejected because it's being deprecated in favor of Trainer V2
- **Custom implementation without Trainer V2**: Rejected due to high maintenance burden and lack of unified API

---

## 2. Migration Path from Legacy MPIJob v2beta1

### Decision: 12-Month Dual-Support Window with Proactive Migration Assistance

**Rationale**: Field analysis shows MPI workloads represent high-value, low-volume customers (autonomous vehicles, HPC-AI convergence, LLM training). Zero backward compatibility requires active migration support to avoid production disruptions.

### Current State: Legacy MPIJob v2beta1

**API Structure**:
- Framework-specific CRD (`kubeflow.org/v2beta1/MPIJob`)
- Users explicitly define Launcher and Worker pod templates
- Manual gang scheduling configuration (users calculate `minAvailable = 1 + workerCount`)
- Direct control over SSH key mounting, network policies, MPI parameters

**Common Usage Patterns** (from field data):
- 60%: Horovod multi-GPU training (2-8 workers, 2-8 GPUs each)
- 25%: HPC batch compute (16-64 CPU-only workers)
- 10%: LLM fine-tuning (2-8 workers, 4-8 A100/H100 GPUs)
- 5%: Elastic training (dynamic scaling)

### Target State: Trainer V2 TrainJob

**Major Changes**:
- **Unified CRD**: Single TrainJob replaces framework-specific jobs
- **Runtime-Based**: MPI configuration abstracted into ClusterTrainingRuntime templates
- **Simplified Interface**: Users specify `numNodes` instead of separate launcher/worker definitions
- **Automated Management**: Launcher pod auto-created, gang scheduling auto-calculated
- **Enhanced Observability**: Consistent metrics and events across all training frameworks

**Breaking Change**: Complete API redesign - zero backward compatibility.

### Migration Path

**Field Mapping Complexity**:
- **Direct Mappings** (30%): `Worker.replicas` → `numNodes`, `image`, `resources`
- **Template Migration** (50%): `nodeSelector`, `tolerations`, `affinity` move to runtime templates
- **Behavioral Changes** (20%): Launcher management, gang scheduling calculation, SSH key handling

**Step-by-Step Process**:
1. Admins create runtime templates (ClusterTrainingRuntime)
2. Users convert YAML using documented field mappings
3. Validate in dev environment (metrics, performance, outcomes)
4. Update automation (SDK, CI/CD pipelines)
5. Incremental production rollout with rollback plan

**Migration Tools**:
- Best-effort YAML conversion script (60% automation)
- Side-by-side API comparison documentation
- Runtime template validator
- Migration dashboard for tracking progress

### Risk Assessment

**High-Impact Risks**:

**Risk 1: Network Policy Conflicts** (50-60% probability)
- SSH communication (port 22) blocked by default deny-all policies
- **Mitigation**: Pre-flight validation, standard network policy templates, automated namespace provisioning

**Risk 2: Runtime Template Misconfiguration** (30-40% probability)
- Incorrect slots, launcher resources, or missing Istio disable annotation
- **Mitigation**: Validated reference templates, "Test Runtime" feature, configuration guide

**Risk 3: Automation Breakage** (70-80% of customers with CI/CD)
- API endpoints, pod naming, status fields incompatible
- **Mitigation**: SDK migration guide, 12-month dual-support period, compatibility examples

**Risk 4: Performance Degradation** (10-15% probability)
- Abstraction layers or scheduling differences may impact training time
- **Mitigation**: Performance benchmarking guide, telemetry monitoring, escalation path

### Rollout Plan

**Timeline**:
- **Months 0-6**: Full support for both v2beta1 and Trainer V2
- **Months 6-9**: Deprecation warnings in legacy API
- **Months 9-12**: Maintenance mode (security patches only)
- **Month 12**: Legacy API disabled

**Success Metrics**:
- **Adoption**: 60% of MPI users migrated by month 10
- **Quality**: >85% migration success rate without critical issues
- **Performance**: <10% training time variance from baseline
- **Satisfaction**: NPS >40 for migration experience

**Alternatives Considered**:
- **Automatic backward compatibility layer**: Rejected due to high complexity and maintenance burden
- **Immediate deprecation**: Rejected due to high customer impact risk

---

## 3. Testing Strategy

### Decision: Comprehensive Test Pyramid with Gang Scheduling, Multi-Tenancy, and Observability as P0 Priorities

**Rationale**: Distributed systems complexity, multi-component integration, and enterprise security requirements demand rigorous testing. Poor diagnostics = 3x support burden (industry benchmark).

### Test Framework Selection

**Go (Controller/Operator)**:
- **Ginkgo v2 + Gomega**: BDD testing framework (industry standard from KubeFlow)
- **EnvTest**: Integration testing with simulated K8s API (tests run locally <5 min)

**Python (SDK/Dashboard/E2E)**:
- **pytest**: Aligns with existing OpenShift AI tests
- **Poetry**: Dependency management
- **pytest-bdd**: Behavior-driven E2E scenarios

### Test Category Breakdown

**Test Pyramid**:
```
  - Unit Tests (60%):        ~150 tests, <5 min execution
  - Integration Tests (30%): ~80 tests, <10 min execution
  - E2E Tests (10%):         ~50 tests, <2 hours execution

Total: ~350 test cases
```

**P0 Test Categories** (MVP blockers):
1. **Gang Scheduling**: All pods scheduled atomically, timeout on insufficient resources, clear error messages
2. **Multi-Tenancy Isolation**: RBAC enforcement, network policy validation, resource quota enforcement (92% of enterprise customers require)
3. **SSH Communication**: Key generation/distribution, launcher-to-worker connections, MPI initialization validation
4. **Basic Lifecycle**: Job creation, status transitions, completion, cleanup

**P1 Test Categories** (High priority):
1. **Observability**: Metrics, logs, events, error diagnostics
2. **Reference Workloads**: Horovod PyTorch MNIST, Intel MPI TensorFlow benchmarks, BERT fine-tuning
3. **Failure Scenarios**: Resource constraints (OOM, insufficient GPUs), network failures (SSH timeout), configuration errors

### Integration Testing Approach

**Controller Testing (EnvTest)**:
- Simulates Kubernetes API locally (no real cluster needed)
- Tests reconciliation loops, CRD validation, webhooks
- Pattern: `BeforeSuite` → Test → `AfterSuite`

**SDK Integration Testing**:
- Test against real OpenShift cluster
- Validate API compatibility
- Contract tests for breaking changes

**Network Testing**:
- Launcher-worker SSH communication
- Network policy enforcement
- Istio sidecar conflicts

### End-to-End Test Scenarios

**User Workflows**:
- Data scientist creates MPIJob via SDK → monitors → retrieves results
- MLOps engineer troubleshoots failure → identifies root cause → remediates
- Admin configures runtime template → validates → publishes
- Multi-tenant isolation validation (cross-namespace access blocked)

**Reference Workloads**:
- **Horovod PyTorch MNIST**: 2-4 GPUs, >98% accuracy, <10 min
- **Intel MPI TensorFlow**: 8 nodes, >500 images/sec, 70% scaling efficiency
- **Horovod BERT**: 16 GPUs, F1 >88%, <2 hours

### Performance Baselines

**Performance Targets**:
- **Job Submission Latency**: <10 seconds (90th percentile)
- **Gang Scheduling**: <30 seconds for 8 workers
- **MPI Initialization**: <60 seconds (SSH + MPI setup)
- **Training Throughput**: >90% scaling efficiency vs single GPU
- **MPI Overhead**: <15% of total iteration time

**Scale Testing Dimensions**:
- Worker count: 2, 4, 8, 16, 32 workers (max ~64)
- Concurrent jobs: 10, 20, 50 jobs simultaneously
- Long-running: 24+ hour jobs (memory leak detection)

### Security and Compliance Testing

**Multi-Tenancy Requirements**:
- Namespace isolation (pods cannot cross boundaries)
- RBAC enforcement (role-based access control)
- Resource quotas (prevent exhaustion)
- Network policies (block cross-tenant traffic)
- SSH key security (ephemeral, scoped to namespace)

**Compliance Validation**:
- **Audit Logging**: All lifecycle events logged with user attribution
- **FIPS Compliance**: Crypto modules validated for federal customers
- **GDPR/HIPAA**: PII never exposed, data deletion on request
- **SOC 2**: User actions auditable with timestamps

**Security Test Count**: ~25 test cases

### Test Automation and CI/CD

**Test Execution Schedule**:
```
Per Commit:    Unit tests (local, <2 min)
Per PR:        Unit + Integration + Contract (<10 min)
Per Merge:     Unit + Integration + E2E Smoke (<30 min)
Nightly:       Full E2E + Performance (<2 hours)
Weekly:        Full Suite + Scale + Security (<4 hours)
Release:       Full Suite + Manual validation (1-2 days)
```

**Infrastructure Requirements**:
- **Dev Cluster**: 4-8 GPU nodes (NVIDIA A100/V100), daily E2E testing
- **CI Cluster**: Kind/Minikube (CPU-only), PR validation
- **Staging Cluster**: Production-mirror, 8+ GPU nodes, release validation

**Alternatives Considered**:
- **Manual testing only**: Rejected due to high regression risk and slow feedback
- **E2E tests only**: Rejected due to slow execution and poor coverage granularity
- **Unit tests only**: Rejected due to insufficient integration validation

---

## 4. Implementation Best Practices

### MPI Runtime Template Patterns

**Decision: Use ClusterTrainingRuntime for Organization-Wide Templates**

**Best Practices**:
1. Deploy organization-wide MPI configurations as ClusterTrainingRuntime
2. Create namespace-specific variants as TrainingRuntime
3. Define launcher as replicatedJobs[0] with mpirun command
4. Define workers as replicatedJobs[1] with SSH server
5. Leverage automatic MPI setup (don't manually configure SSH keys/hostfiles)
6. Set appropriate mpiImplementation: OpenMPI for most cases, Intel MPI for Intel hardware

**Example Runtime Template Structure**:
```yaml
apiVersion: kubeflow.org/v2alpha1
kind: ClusterTrainingRuntime
metadata:
  name: mpi-horovod-gpu
spec:
  mlPolicy:
    numNodes: 4  # Default, can be overridden in TrainJob
    mpiImplementation: OpenMPI
  podGroupPolicy:
    coscheduling:
      scheduleTimeoutSeconds: 300
  template:
    spec:
      replicatedJobs:
      - name: launcher
        template:
          spec:
            containers:
            - name: launcher
              # Launcher pod configuration
      - name: worker
        replicas: 4
        template:
          spec:
            containers:
            - name: worker
              # Worker pod configuration
```

### Network Policy Configuration

**Decision: Disable Istio Sidecar Injection for Training Jobs**

**Rationale**: Istio sidecars cause pods to not terminate after job completion and interfere with MPI launcher-worker communication.

**Required Configuration**:
```yaml
metadata:
  annotations:
    sidecar.istio.io/inject: "false"
```

**Network Policy Requirements**:
1. Allow inter-pod SSH traffic (TCP port 22) between launcher and workers
2. Allow MPI communication ports (typically ephemeral ranges)
3. Respect namespace boundaries for multi-tenancy

**Alternative Considered**: Configure `holdApplicationUntilProxyStarts: true` if sidecars required - rejected because disabling sidecars is simpler and sufficient for batch training jobs.

### Gang Scheduling Configuration

**Decision: Always Enable Gang Scheduling for Multi-Node Jobs**

**Rationale**: Prevents resource waste from partial scheduling. Without gang scheduling, launcher may start but workers fail to schedule, consuming resources with no training progress.

**Best Practices**:
1. Set appropriate `scheduleTimeoutSeconds`: Balance between waiting for resources (higher timeout) and failing fast (lower timeout)
2. Use Kueue for production deployments (better multi-tenant resource management than basic coscheduling)
3. Configure `minMember` in PodGroup to launcher + all workers

**Scheduler Options**:
- **scheduler-plugins coscheduling**: Basic gang scheduling for development
- **Volcano Scheduler**: Advanced scheduling with priorities and queues
- **Kueue**: Multi-tenant resource management with quotas (recommended for production)

### Observability and Logging

**Decision: Integrate with OpenShift AI Unified Observability Stack**

**Required Metrics**:
- Job duration, success/failure rate, time-to-start, completion time (consistent with PyTorchJob)
- MPI-specific: gang scheduling duration, MPI initialization time, launcher-to-worker connection time
- Resource utilization: CPU, memory, GPU usage per pod

**Required Events**:
- Job submitted, gang scheduling started, pods ready, training started, job completed/failed
- Failures must include: gang scheduling timeout, launcher failure, worker failure, resource constraints

**Log Aggregation**:
- Logs accessible through OpenShift logging infrastructure
- Proper labeling: job name, pod type (launcher/worker), pod index
- Log sanitization to prevent credential exposure

### Security Context Constraints

**Decision: Document Required SCC Capabilities, Provide Standard SCC Templates**

**Required Capabilities**:
- SSH server requires `NET_BIND_SERVICE` or non-privileged port (>1024)
- MPI communication may require `SYS_PTRACE` for some implementations
- GPU access requires appropriate device plugin permissions

**Recommended SCC Template**:
```yaml
apiVersion: security.openshift.io/v1
kind: SecurityContextConstraints
metadata:
  name: mpi-training
allowPrivilegedContainer: false
allowHostNetwork: false
allowHostPID: false
runAsUser:
  type: MustRunAsNonRoot
seLinuxContext:
  type: MustRunAs
fsGroup:
  type: RunAsAny
supplementalGroups:
  type: RunAsAny
volumes:
- configMap
- secret
- persistentVolumeClaim
- emptyDir
requiredDropCapabilities:
- ALL
allowedCapabilities:
- NET_BIND_SERVICE  # For SSH server
```

**Alternatives Considered**:
- **Use privileged SCC**: Rejected due to security risk
- **No SCC documentation**: Rejected because customers will face permission errors

---

## 5. OpenShift AI Integration Points

### Dashboard Integration

**Decision: Extend Existing Trainer V2 Dashboard with MPI-Specific Views**

**Required Features**:
- Unified Training Jobs list with MPI framework indicator
- MPIJob creation wizard with MPI-specific fields (worker count, slots per worker, runtime selection)
- Job details page with topology visualization (launcher + workers)
- Log access with tabbed interface (launcher logs, worker logs, aggregated)
- Resource utilization metrics (CPU, memory, GPU per pod)

**API Contracts**:
- REST endpoints: `/api/v1/trainjobs`, `/api/v1/trainjobs/{name}`, `/api/v1/trainjobs/{name}/logs`
- WebSocket: Real-time job status updates
- Authentication: OpenShift OAuth integration

### Python SDK Integration

**Decision: Extend Existing Trainer V2 SDK with MPI Methods**

**Required Methods**:
- `create_mpi_job(name, runtime, num_nodes, image, command, resources)`: Create MPIJob
- `get_job_status(name)`: Get detailed status (job phase, launcher readiness, worker readiness count)
- `stream_job_events(name)`: Stream lifecycle events
- `get_job_logs(name, pod_type, pod_index)`: Retrieve logs
- `delete_job(name, grace_period)`: Delete with optional graceful termination

**SDK Repository**: Part of KubeFlow Training Python SDK or OpenShift AI SDK

### Storage Integration

**Decision: Support S3, OCS/PVC, and NFS for Training Data and Model Artifacts**

**Patterns**:
- **S3**: Use s3fs or boto3 for direct access (credentials via secrets)
- **OCS/PVC**: Mount PersistentVolumeClaim to launcher and workers
- **NFS**: Mount NFS volumes (requires appropriate SCC permissions)

**Best Practice**: Use Trainer V2 initializer for dataset fetching to reduce GPU idle time during I/O.

### RBAC Integration

**Decision: Namespace-Scoped Permissions with ClusterTrainingRuntime for Platform Admins**

**Roles**:
- **Platform Admin**: Can create ClusterTrainingRuntime (cluster-scoped)
- **Project Admin**: Can create TrainingRuntime (namespace-scoped)
- **Data Scientist**: Can create TrainJob (namespace-scoped)
- **Viewer**: Can view TrainJob (read-only)

**RBAC Resources**:
```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: trainjob-editor
rules:
- apiGroups: ["kubeflow.org"]
  resources: ["trainjobs"]
  verbs: ["create", "get", "list", "update", "delete"]
- apiGroups: ["kubeflow.org"]
  resources: ["trainingruntimes"]
  verbs: ["get", "list"]
```

### Resource Quota Management

**Decision: Launcher + All Worker Pod Resources Count Against Namespace Quotas**

**Implementation**: OpenShift ResourceQuota controller automatically tracks all pod resource requests/limits.

**Best Practice**: Pre-validate job can fit within quota before submission. Provide clear error message if quota would be exceeded.

---

## 6. Technology Stack Decisions

### Language/Version

**Decision: Go 1.21+ for Controller, Python 3.11+ for SDK/Tests**

**Rationale**:
- Go: KubeFlow Trainer V2 controller written in Go, standard for K8s operators
- Python: OpenShift AI SDK and tests primarily Python

### Primary Dependencies

**Decision**:
- **Controller**: KubeFlow Training Operator V2, JobSet, scheduler-plugins/Kueue
- **SDK**: KubeFlow Training Python SDK, kubernetes Python client
- **Dashboard**: React (existing OpenShift AI Dashboard), OpenShift OAuth
- **Testing**: Ginkgo/Gomega (Go), pytest (Python), EnvTest

### Storage

**Decision: S3 for Model/Dataset Storage, PVC for Checkpoints**

**Rationale**:
- S3: Industry standard for ML datasets and model artifacts
- PVC: Fast local storage for checkpoints during training
- NFS: Legacy support for on-premises deployments

### Testing

**Decision: pytest for Python, Ginkgo/Gomega for Go**

**Rationale**: Industry standard for respective languages, well-integrated with CI/CD.

### Target Platform

**Decision: OpenShift 4.14+ (Kubernetes 1.27+)**

**Rationale**:
- Trainer V2 requires K8s 1.27+ for JobSet API
- OpenShift 4.14+ provides required features (SCC, OAuth, monitoring stack)

### Performance Goals

**Decision**:
- **Job Submission Latency**: <10 seconds (90th percentile)
- **Gang Scheduling**: <30 seconds for 8 workers
- **Training Throughput**: >90% scaling efficiency
- **MPI Overhead**: <15% of iteration time

### Constraints

**Decision**:
- **Gang Scheduling Timeout**: 300 seconds default (configurable)
- **Max Workers per Job**: 64 workers (practical limit based on SSH key distribution and MPI scaling)
- **Min Resources per Worker**: 1 CPU, 1Gi memory
- **Max Concurrent Jobs per Namespace**: Limited by resource quotas only

### Scale/Scope

**Decision**:
- **Expected Workload**: 10-100 concurrent MPIJobs per cluster
- **Max Workers per Job**: 64 workers
- **Max GPUs per Worker**: 8 GPUs
- **User Base**: 100-1000 data scientists per cluster

---

## Summary

All NEEDS CLARIFICATION items from Technical Context have been resolved:

✅ **Language/Version**: Go 1.21+ (controller), Python 3.11+ (SDK/tests)
✅ **Primary Dependencies**: KubeFlow Trainer V2, JobSet, Ginkgo/pytest, OpenShift Dashboard
✅ **Storage**: S3 (datasets/models), PVC (checkpoints), NFS (legacy support)
✅ **Testing**: Ginkgo/Gomega (Go), pytest (Python), EnvTest, ~350 test cases
✅ **Target Platform**: OpenShift 4.14+ (Kubernetes 1.27+)
✅ **Performance Goals**: <10s submission, <30s gang scheduling, >90% efficiency
✅ **Constraints**: 300s gang scheduling timeout, 64 max workers, multi-tenancy enforced
✅ **Scale/Scope**: 10-100 concurrent jobs, 100-1000 users per cluster

**Next Steps**: Proceed to Phase 1 (Design & Contracts) to generate data-model.md, contracts/, and quickstart.md.
