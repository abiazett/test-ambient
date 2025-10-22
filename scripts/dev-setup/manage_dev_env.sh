#!/bin/bash
# Manage development environment for MPIJob support

# Constants
KIND_CLUSTER_NAME="mpijob-dev"

# Command-line argument parsing
function print_usage() {
    echo "Usage: $0 <command>"
    echo "Commands:"
    echo "  start     - Start the development environment"
    echo "  stop      - Stop the development environment"
    echo "  status    - Check status of the development environment"
    echo "  logs      - View logs from MPIJob controller"
    echo "  deploy    - Deploy latest code changes to the development environment"
    echo "  clean     - Remove all MPIJobs from the test namespace"
    echo "  reset     - Reset the development environment (delete and recreate)"
}

# Start development environment
function start_env() {
    echo "Starting development environment..."

    # Check if cluster exists
    if ! kind get clusters | grep -q "^${KIND_CLUSTER_NAME}$"; then
        echo "Cluster ${KIND_CLUSTER_NAME} not found. Please run setup.sh first."
        exit 1
    fi

    # Start cluster if not running
    if ! kubectl cluster-info --context kind-${KIND_CLUSTER_NAME} &> /dev/null; then
        kind export kubeconfig --name ${KIND_CLUSTER_NAME}
        echo "Cluster started."
    else
        echo "Cluster is already running."
    fi

    # Set correct kubectl context
    kubectl config use-context kind-${KIND_CLUSTER_NAME}

    echo "Development environment is ready."
    echo "Use 'kubectl get pods -A' to verify all components are running."
}

# Stop development environment
function stop_env() {
    echo "Stopping development environment..."

    # Check if cluster exists
    if ! kind get clusters | grep -q "^${KIND_CLUSTER_NAME}$"; then
        echo "Cluster ${KIND_CLUSTER_NAME} not found."
        exit 1
    fi

    # Stop the cluster
    kind delete cluster --name ${KIND_CLUSTER_NAME}

    echo "Development environment stopped."
}

# Check status of development environment
function check_status() {
    echo "Checking development environment status..."

    # Check if cluster exists
    if ! kind get clusters | grep -q "^${KIND_CLUSTER_NAME}$"; then
        echo "Cluster ${KIND_CLUSTER_NAME} not found. Please run setup.sh first."
        exit 1
    fi

    # Check if cluster is running
    if ! kubectl cluster-info --context kind-${KIND_CLUSTER_NAME} &> /dev/null; then
        echo "Cluster is not running. Use 'start' command to start it."
        exit 1
    fi

    # Set correct kubectl context
    kubectl config use-context kind-${KIND_CLUSTER_NAME}

    # Check component status
    echo "=== Kubernetes Components ==="
    kubectl get pods -n kube-system

    echo -e "\n=== KubeFlow Training Operator ==="
    kubectl get pods -n kubeflow

    echo -e "\n=== Volcano Scheduler ==="
    kubectl get pods -n volcano-system

    echo -e "\n=== MPIJobs ==="
    kubectl get mpijobs -A

    echo -e "\n=== Nodes ==="
    kubectl get nodes

    echo -e "\nDevelopment environment is ready."
}

# View logs from MPIJob controller
function view_logs() {
    echo "Viewing logs from MPIJob controller..."

    # Set correct kubectl context
    kubectl config use-context kind-${KIND_CLUSTER_NAME}

    # Get logs from training operator
    kubectl logs -f deployment/training-operator -n kubeflow
}

# Deploy latest code changes
function deploy_changes() {
    echo "Deploying latest code changes to development environment..."

    # Set correct kubectl context
    kubectl config use-context kind-${KIND_CLUSTER_NAME}

    # Update CRDs
    kubectl apply -f ../../kubernetes/crds/mpijob_crd.yaml
    kubectl apply -f ../../kubernetes/crds/mpijob_validation.yaml

    # Update RBAC
    kubectl apply -f ../../kubernetes/rbac/mpijob_roles.yaml

    # Update NetworkPolicy templates
    kubectl apply -f ../../kubernetes/network/mpijob_network_policies.yaml

    # Restart training operator to pick up changes
    kubectl rollout restart deployment/training-operator -n kubeflow

    echo "Changes deployed. Waiting for training operator to restart..."
    kubectl rollout status deployment/training-operator -n kubeflow

    echo "Deployment complete."
}

# Clean up MPIJobs
function clean_jobs() {
    echo "Cleaning up MPIJobs from test namespace..."

    # Set correct kubectl context
    kubectl config use-context kind-${KIND_CLUSTER_NAME}

    # Delete all MPIJobs in the test namespace
    kubectl delete mpijobs --all -n mpijob-test

    echo "Cleanup complete."
}

# Reset development environment
function reset_env() {
    echo "Resetting development environment..."

    # Stop the environment
    stop_env

    # Run setup script
    ./setup.sh

    echo "Development environment reset complete."
}

# Main
if [ $# -eq 0 ]; then
    print_usage
    exit 1
fi

case "$1" in
    start)
        start_env
        ;;
    stop)
        stop_env
        ;;
    status)
        check_status
        ;;
    logs)
        view_logs
        ;;
    deploy)
        deploy_changes
        ;;
    clean)
        clean_jobs
        ;;
    reset)
        reset_env
        ;;
    *)
        print_usage
        exit 1
        ;;
esac