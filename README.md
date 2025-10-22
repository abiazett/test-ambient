# Test-ambient repo

## MPIJob Support for OpenShift AI

This project implements MPIJob support for OpenShift AI, enabling distributed training using the Message Passing Interface (MPI) protocol. The implementation follows the plan outlined in `001-support-for-mpijob/plan.md` and the tasks in `001-support-for-mpijob/tasks.md`.

## Project Structure

The project is organized into the following directories:

```
├── kubernetes/           # Kubernetes manifests and configurations
│   ├── crds/             # Custom Resource Definitions
│   ├── operator/         # KubeFlow Training Operator configuration
│   ├── rbac/             # RBAC roles and permissions
│   ├── network/          # NetworkPolicy templates
│   └── scheduling/       # Volcano scheduler configuration
├── cli/                  # Command Line Interface (Go)
│   ├── cmd/              # CLI commands
│   │   ├── training/     # Training-related commands
│   └── test/             # CLI tests
├── sdk/                  # Python SDK for OpenShift AI
│   ├── odh/              # SDK package
│   │   └── training/     # Training-specific modules
│   ├── examples/         # Example usage
│   └── tests/            # SDK tests
├── frontend/             # Web UI components (React)
│   ├── src/              # Source code
│   │   ├── components/   # React components
│   │   └── services/     # API services
├── api/                  # API definitions
│   └── training/         # Training APIs
├── service/              # Backend services
│   └── training/         # Training services
├── operator/             # Custom operators
│   └── controllers/      # Controller implementations
├── test/                 # Test resources
│   ├── manifests/        # Test manifests
│   ├── integration/      # Integration tests
│   └── e2e/              # End-to-end tests
└── scripts/              # Development and deployment scripts
    └── dev-setup/        # Development environment setup
```

## Implementation Progress

### Completed Tasks

- [x] Created project directory structure
- [x] Initialized Kubernetes components
  - [x] MPIJob CRD definition
  - [x] RBAC roles and permissions
  - [x] NetworkPolicy templates
  - [x] Volcano scheduler configuration
  - [x] KubeFlow Training Operator V2 deployment manifests
- [x] Set up CI/CD pipeline configuration
  - [x] CI workflow for build and test
  - [x] CD workflow for releases
  - [x] PR template
- [x] Created development environment setup scripts
  - [x] Environment setup
  - [x] Test job creation
  - [x] Environment management
- [x] Initialized Go module for CLI
  - [x] Basic CLI structure
  - [x] Training command group
  - [x] MPIJob create command
  - [x] MPIJob delete command
  - [x] MPIJob list command
  - [x] MPIJob describe command
  - [x] MPIJob logs command
- [x] Initialized Python package for SDK
  - [x] SDK package structure
  - [x] MPIJob models
  - [x] MPIJobClient class
  - [x] MPIJob class
  - [x] ResourceSpec class
  - [x] Example usage
- [x] Configured React project for UI components
  - [x] Application structure
  - [x] MPIJob form
  - [x] Job list view
  - [x] Job detail view
  - [x] Worker topology visualization

### Next Steps

The following tasks should be completed next according to the implementation plan:

1. **Core Implementations**
   - [ ] Implement MPIJob CRD validation in `kubernetes/crds/mpijob_validation.yaml`
   - [ ] Create MPIJob controller reconciliation logic in `operator/controllers/mpijob_controller.go`
   - [ ] Implement MPIJob status tracking in `operator/controllers/mpijob_status.go`
   - [ ] Build Training Service API gateway in `api/training/mpijob.go`
   - [ ] Implement Training Service backend with gRPC in `service/training/mpijob.go`

2. **Test Implementations**
   - [ ] Create MPIJob CRD test manifests in `test/manifests/mpijob/`
   - [ ] Set up integration test framework in `test/integration/`
   - [ ] Implement test utilities for job validation in `test/utils/job_validator.go`
   - [ ] Create E2E test suite structure in `test/e2e/`

3. **Integration Tasks**
   - [ ] Implement RBAC integration in `kubernetes/rbac/mpijob_roles.yaml`
   - [ ] Create NetworkPolicy templates in `kubernetes/network/mpijob_network_policies.yaml`
   - [ ] Implement Volcano scheduler integration in `kubernetes/scheduling/volcano_config.yaml`

## Getting Started

### Development Environment Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/openshift-ai/mpijob.git
   cd mpijob
   ```

2. Set up the development environment:
   ```bash
   cd scripts/dev-setup
   ./setup.sh
   ```

3. Start the development environment:
   ```bash
   ./manage_dev_env.sh start
   ```

4. Create test jobs:
   ```bash
   ./create_test_jobs.sh
   ```

### Building Components

#### CLI

```bash
cd cli
go build -o bin/odh .
```

#### SDK

```bash
cd sdk
pip install -e .
```

#### Frontend

```bash
cd frontend
npm install
npm start
```

## Documentation

- CLI Documentation: [cli/README.md](cli/README.md)
- SDK Documentation: [sdk/README.md](sdk/README.md)
- Frontend Documentation: [frontend/README.md](frontend/README.md)
- Development Environment: [scripts/dev-setup/README.md](scripts/dev-setup/README.md)

## License

This project is licensed under the Apache License 2.0.
