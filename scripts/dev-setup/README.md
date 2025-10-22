# Development Environment Setup for MPIJob Support

This directory contains scripts to set up and manage a local development environment for the MPIJob support feature in OpenShift AI.

## Prerequisites

The development environment requires the following tools:

- Docker
- kubectl
- Kind (Kubernetes in Docker)
- Go 1.21 or higher
- Python 3.10 or higher
- Node.js and npm
- Git

## Available Scripts

### 1. `setup.sh`

This script sets up the complete development environment:

- Creates a local Kind Kubernetes cluster
- Installs KubeFlow Training Operator
- Installs Volcano Scheduler
- Applies MPIJob CRDs and RBAC configurations
- Sets up Python, Go, and React development environments

Usage:
```bash
./setup.sh
```

### 2. `manage_dev_env.sh`

This script helps manage the development environment:

- Start/stop the Kind cluster
- Check the status of components
- View logs from the MPIJob controller
- Deploy the latest code changes
- Clean up MPIJobs from the test namespace
- Reset the development environment

Usage:
```bash
./manage_dev_env.sh <command>
```

Available commands:
- `start` - Start the development environment
- `stop` - Stop the development environment
- `status` - Check status of the development environment
- `logs` - View logs from MPIJob controller
- `deploy` - Deploy latest code changes
- `clean` - Remove all MPIJobs from the test namespace
- `reset` - Reset the development environment

### 3. `create_test_jobs.sh`

This script creates sample MPIJob instances for testing:

- Basic MPIJob with 4 workers
- MPIJob with GPU resources
- MPIJob with Intel MPI implementation
- MPIJob with network policy

Usage:
```bash
./create_test_jobs.sh
```

## Development Workflow

1. **Initial Setup**

   ```bash
   # Run the setup script
   ./setup.sh
   ```

2. **Daily Development**

   ```bash
   # Start the development environment
   ./manage_dev_env.sh start

   # Check status
   ./manage_dev_env.sh status

   # Make code changes...

   # Deploy changes
   ./manage_dev_env.sh deploy

   # Create test jobs
   ./create_test_jobs.sh

   # View logs
   ./manage_dev_env.sh logs
   ```

3. **Testing**

   ```bash
   # Test the CLI
   cd ../../cli
   go test ./...

   # Test the operator
   cd ../../operator
   go test ./...

   # Test the SDK
   cd ../../sdk
   pytest
   ```

4. **Clean Up**

   ```bash
   # Clean up test jobs
   ./manage_dev_env.sh clean

   # Stop the environment when done
   ./manage_dev_env.sh stop
   ```

## Troubleshooting

- **Kind cluster fails to start**: Ensure Docker is running and you have enough system resources.
- **Pod scheduling issues**: Check node capacity with `kubectl describe nodes`.
- **MPIJob controller issues**: View logs with `./manage_dev_env.sh logs`.
- **Test jobs fail**: Check the error message with `kubectl describe mpijob <job-name> -n mpijob-test`.

## Additional Resources

- [KubeFlow Training Operator Documentation](https://github.com/kubeflow/training-operator)
- [MPI Operator Implementation Guide](https://github.com/kubeflow/mpi-operator)
- [Volcano Scheduler Documentation](https://github.com/volcano-sh/volcano)