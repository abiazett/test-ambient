# MPIJobs Support in OpenShift AI via Trainer V2 - Implementation Workspace

This workspace contains the design documents, reference implementations, and tooling for adding MPIJob support to OpenShift AI using the KubeFlow Trainer V2 unified training API.

## Project Overview

This feature implements support for MPI-based distributed training jobs (Horovod, Intel MPI, OpenMPI) in OpenShift AI. The implementation spans multiple upstream repositories:

- **kubeflow/training-operator** - Go controller implementation
- **kubeflow/training-sdk** - Python SDK implementation
- **opendatahub-io/odh-dashboard** - React Dashboard UI
- **red-hat-data-services/openshift-ai-tests** - E2E tests

## Workspace Structure

```
test-ambient/
├── specs/001-mpijob-trainer-v2-support/  # Feature specification and design
│   ├── spec.md                           # User-facing requirements
│   ├── plan.md                           # Implementation plan
│   ├── tasks.md                          # Detailed task breakdown
│   ├── research.md                       # Technical research findings
│   ├── data-model.md                     # Entity relationships
│   ├── quickstart.md                     # Getting started guide
│   └── contracts/                        # API contracts
│       ├── dashboard-api.yaml            # Dashboard REST API (OpenAPI 3.0)
│       ├── python-sdk.md                 # Python SDK contract
│       └── crd-schemas.yaml              # Kubernetes CRD schemas
│
├── implementations/                      # Reference implementations
│   ├── controller/                       # Go controller code (for kubeflow/training-operator)
│   ├── sdk/                              # Python SDK code (for kubeflow/training-sdk)
│   ├── dashboard/                        # React UI code (for opendatahub-io/odh-dashboard)
│   └── tests/                            # E2E tests (for red-hat-data-services/openshift-ai-tests)
│
├── manifests/                            # Kubernetes manifests and templates
│   ├── runtimes/                         # ClusterTrainingRuntime templates
│   ├── network-policy/                   # NetworkPolicy templates
│   ├── scc/                              # OpenShift Security Context Constraints
│   └── grafana/                          # Grafana dashboard templates
│
├── tools/                                # Migration and validation tools
│   ├── migrate-mpijob.py                 # YAML conversion tool (v2beta1 → TrainJob)
│   └── validate-runtime.py               # Runtime template validator
│
└── docs/                                 # User and administrator documentation
    ├── user-guide/                       # End-user documentation
    ├── admin-guide/                      # Administrator documentation
    ├── migration/                        # Migration guides
    ├── troubleshooting/                  # Troubleshooting guides
    └── api-reference/                    # API reference documentation
```

## Implementation Approach

Since this is a multi-repository project, this workspace serves as:

1. **Design Hub**: All design documents, API contracts, and architectural decisions
2. **Reference Implementation**: Complete implementations ready to be contributed to upstream repos
3. **Tooling**: Migration scripts, validators, and automation tools
4. **Documentation**: Comprehensive docs that will be distributed across repositories

## Getting Started

### For Contributors

1. Review the feature specification: `specs/001-mpijob-trainer-v2-support/spec.md`
2. Understand the implementation plan: `specs/001-mpijob-trainer-v2-support/plan.md`
3. Check the task breakdown: `specs/001-mpijob-trainer-v2-support/tasks.md`
4. Follow the quickstart guide: `specs/001-mpijob-trainer-v2-support/quickstart.md`

### For Users (After Implementation)

See the quickstart guide for step-by-step instructions on creating your first MPIJob:
`specs/001-mpijob-trainer-v2-support/quickstart.md`

## Key Features

- **Unified API**: Same interfaces (SDK, CLI, Dashboard) as PyTorchJob and TensorFlowJob
- **Gang Scheduling**: Atomic scheduling of launcher and worker pods
- **Multi-Tenancy**: Namespace isolation, RBAC, and resource quotas
- **Observability**: Comprehensive metrics, logs, and error diagnostics
- **Migration Support**: Tools and guides for migrating from legacy MPIJob v2beta1

## Technology Stack

- **Controller**: Go 1.21+
- **SDK**: Python 3.11+
- **Dashboard**: React (OpenShift AI Dashboard)
- **Testing**: pytest (Python), Ginkgo (Go), EnvTest (integration)
- **Target Platform**: OpenShift 4.14+ (Kubernetes 1.27+)

## Dependencies

- KubeFlow Training Operator V2
- JobSet (Kubernetes 1.27+)
- Gang Scheduling (scheduler-plugins or Kueue)
- OpenShift 4.14+

## Status

This workspace is actively under development. See `specs/001-mpijob-trainer-v2-support/tasks.md` for the current task status.

## Contributing

This feature spans multiple upstream repositories. Contributions should be made to the appropriate upstream repo:

- **Controller**: https://github.com/kubeflow/training-operator
- **SDK**: https://github.com/kubeflow/training-sdk
- **Dashboard**: https://github.com/opendatahub-io/odh-dashboard
- **Tests**: https://github.com/red-hat-data-services/openshift-ai-tests

## License

Follow the license of each respective upstream repository.

## Contact

For questions or feedback, please refer to the OpenShift AI and KubeFlow community channels.

---

## Original: Test-ambient repo

Create Branches for each RFE test
