# OpenShift AI CLI for MPIJob Support

Command-line interface for OpenShift AI that supports managing MPIJobs.

## Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/openshift-ai/mpijob
cd mpijob/cli

# Build
go build -o odh .

# Install
mv odh /usr/local/bin/
```

### Using Container Image

```bash
# Pull the container image
podman pull quay.io/openshiftai/mpijob-cli:latest

# Create an alias
alias odh='podman run --rm -it -v ~/.kube:/root/.kube:Z quay.io/openshiftai/mpijob-cli:latest'
```

## Usage

The CLI follows a consistent command structure:

```
odh training <action> mpijob [name] [flags]
```

Where:
- `<action>` is one of: create, delete, list, describe, logs
- `mpijob` is the job type
- `[name]` is the name of the job (required for some commands)
- `[flags]` are optional flags for the command

### Examples

#### Creating an MPIJob

```bash
# Create from a YAML file
odh training create mpijob --from-file=mpijob.yaml

# Create using inline parameters
odh training create mpijob my-job --workers=4 --gpu=2 --cpu=4 --memory=16Gi \
  --image=myregistry.com/myimage:latest --command="python" --command="/train.py"
```

#### Listing MPIJobs

```bash
# List all MPIJobs in current namespace
odh training list mpijob

# List MPIJobs in all namespaces
odh training list mpijob --all-namespaces

# Filter by status
odh training list mpijob --status=Running
```

#### Describing an MPIJob

```bash
# Describe an MPIJob
odh training describe mpijob my-job

# Watch status updates in real-time
odh training describe mpijob my-job --watch
```

#### Getting Logs from an MPIJob

```bash
# Get logs from the launcher
odh training logs mpijob my-job

# Get logs from a specific worker
odh training logs mpijob my-job --worker=0

# Follow logs in real-time
odh training logs mpijob my-job --follow
```

#### Deleting an MPIJob

```bash
# Delete an MPIJob
odh training delete mpijob my-job

# Force delete
odh training delete mpijob my-job --force
```

## Development

### Prerequisites

- Go 1.21 or higher
- Access to an OpenShift cluster with MPIJob support

### Building

```bash
# Build the CLI
go build -o odh .

# Run tests
go test -v ./...
```

### Creating a Release

```bash
# Build for multiple platforms
GOOS=linux GOARCH=amd64 go build -o odh-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o odh-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o odh-windows-amd64.exe .
```