#!/bin/bash
set -e

# MPIJob Development Environment Setup
# This script sets up a local development environment for the MPIJob support feature

# Check for required tools
check_requirements() {
    echo "Checking system requirements..."

    # Check for kubectl
    if ! command -v kubectl &> /dev/null; then
        echo "kubectl not found. Please install kubectl."
        exit 1
    fi

    # Check for kind
    if ! command -v kind &> /dev/null; then
        echo "kind not found. Would you like to install it? (y/n)"
        read -r install_kind
        if [[ "$install_kind" == "y" ]]; then
            curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.17.0/kind-linux-amd64
            chmod +x ./kind
            sudo mv ./kind /usr/local/bin/kind
        else
            echo "kind is required for local development. Exiting."
            exit 1
        fi
    fi

    # Check for Go
    if ! command -v go &> /dev/null; then
        echo "Go not found. Please install Go 1.21 or higher."
        exit 1
    fi

    # Check for Python
    if ! command -v python3 &> /dev/null; then
        echo "Python 3 not found. Please install Python 3.10 or higher."
        exit 1
    fi

    # Check for npm
    if ! command -v npm &> /dev/null; then
        echo "npm not found. Please install Node.js and npm."
        exit 1
    fi

    # Check for Docker
    if ! command -v docker &> /dev/null; then
        echo "Docker not found. Please install Docker."
        exit 1
    fi

    echo "All system requirements satisfied."
}

# Set up local Kubernetes cluster using kind
setup_kind_cluster() {
    echo "Creating Kind Kubernetes cluster..."

    # Define Kind cluster config with GPU support
    cat > kind-config.yaml << EOF
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: mpijob-dev
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
- role: worker
EOF

    # Create cluster
    kind create cluster --config=kind-config.yaml

    echo "Kind cluster 'mpijob-dev' created successfully."

    # Ensure kubectl is pointing to the right context
    kubectl config use-context kind-mpijob-dev
}

# Install KubeFlow Training Operator
install_training_operator() {
    echo "Installing KubeFlow Training Operator v2..."

    # Create namespace
    kubectl create namespace kubeflow || true

    # Apply Training Operator deployment
    kubectl apply -f https://github.com/kubeflow/training-operator/releases/download/v1.5.0/training-operator.yaml

    # Wait for deployment to be available
    echo "Waiting for Training Operator to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/training-operator -n kubeflow

    echo "KubeFlow Training Operator installed successfully."
}

# Install Volcano Scheduler
install_volcano_scheduler() {
    echo "Installing Volcano Scheduler..."

    # Apply Volcano deployment
    kubectl apply -f https://raw.githubusercontent.com/volcano-sh/volcano/master/installer/volcano-development.yaml

    # Wait for deployment to be available
    echo "Waiting for Volcano Scheduler to be ready..."
    kubectl wait --for=condition=available --timeout=300s deployment/volcano-scheduler -n volcano-system

    echo "Volcano Scheduler installed successfully."
}

# Install project CRDs
install_mpijob_crds() {
    echo "Installing MPIJob CRDs..."

    # Apply MPIJob CRD
    kubectl apply -f ../../kubernetes/crds/mpijob_crd.yaml

    # Apply RBAC roles
    kubectl apply -f ../../kubernetes/rbac/mpijob_roles.yaml

    echo "MPIJob CRDs installed successfully."
}

# Set up Python development environment
setup_python_env() {
    echo "Setting up Python development environment..."

    # Create virtual environment
    python3 -m venv .venv
    source .venv/bin/activate

    # Install dev dependencies
    pip install --upgrade pip
    pip install pytest pytest-mock black flake8 mypy

    # Install SDK in development mode
    cd ../../sdk
    pip install -e .
    cd ../scripts/dev-setup

    echo "Python environment set up successfully."
}

# Set up Go development environment
setup_go_env() {
    echo "Setting up Go development environment..."

    # Initialize Go modules for CLI
    cd ../../cli
    go mod init github.com/openshift-ai/mpijob/cli
    go mod tidy
    cd ../scripts/dev-setup

    # Initialize Go modules for operator
    cd ../../operator
    go mod init github.com/openshift-ai/mpijob/operator
    go mod tidy
    cd ../scripts/dev-setup

    echo "Go environment set up successfully."
}

# Set up frontend development environment
setup_frontend_env() {
    echo "Setting up frontend development environment..."

    cd ../../frontend

    # Initialize npm project
    npm init -y

    # Install React and dependencies
    npm install --save react react-dom react-router-dom @patternfly/react-core @patternfly/react-table @patternfly/react-charts axios formik yup

    # Install dev dependencies
    npm install --save-dev typescript @types/react @types/react-dom jest @testing-library/react @testing-library/jest-dom eslint prettier

    # Create basic tsconfig.json
    cat > tsconfig.json << EOF
{
  "compilerOptions": {
    "target": "es5",
    "lib": [
      "dom",
      "dom.iterable",
      "esnext"
    ],
    "allowJs": true,
    "skipLibCheck": true,
    "esModuleInterop": true,
    "allowSyntheticDefaultImports": true,
    "strict": true,
    "forceConsistentCasingInFileNames": true,
    "noFallthroughCasesInSwitch": true,
    "module": "esnext",
    "moduleResolution": "node",
    "resolveJsonModule": true,
    "isolatedModules": true,
    "noEmit": true,
    "jsx": "react-jsx"
  },
  "include": [
    "src"
  ]
}
EOF

    cd ../scripts/dev-setup

    echo "Frontend environment set up successfully."
}

# Main function
main() {
    echo "=== MPIJob Development Environment Setup ==="

    # Create local directory for dev environment
    mkdir -p mpijob-dev
    cd mpijob-dev

    # Run setup steps
    check_requirements
    setup_kind_cluster
    install_training_operator
    install_volcano_scheduler
    install_mpijob_crds
    setup_python_env
    setup_go_env
    setup_frontend_env

    echo "=== Setup Complete ==="
    echo "You can now start developing MPIJob support features."
    echo "To use the development environment:"
    echo "1. Activate Python venv: source mpijob-dev/.venv/bin/activate"
    echo "2. Switch to Kind context: kubectl config use-context kind-mpijob-dev"
    echo "3. See README.md for more details on development workflows"
}

# Execute main function
main