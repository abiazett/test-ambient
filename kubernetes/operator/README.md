# KubeFlow Training Operator Configuration for MPIJob Support

This directory contains the Kubernetes manifests needed to deploy and configure the KubeFlow Training Operator V2 with support for MPIJobs in OpenShift AI.

## Contents

- `deployment.yaml` - Main deployment specification for the Training Operator
- `service.yaml` - Service for exposing metrics and webhook endpoints
- `role.yaml` - RBAC permissions for the Training Operator
- `kustomization.yaml` - Kustomize configuration for deploying the operator
- `patches/` - Directory containing Kustomize patches
  - `enable_mpijob.yaml` - Patch to enable MPIJob controller and webhook support
- `webhook-setup.yaml` - Admission webhook configuration for MPIJob validation

## Installation

### Prerequisites

- OpenShift cluster 4.12 or higher
- Administrative access to the cluster

### Automated Installation

Use the installation script to deploy the Training Operator:

```bash
# From the project root
./scripts/dev-setup/setup.sh
```

### Manual Installation

1. Create the namespace:
   ```bash
   oc new-project openshift-ai
   ```

2. Apply the Kubernetes manifests:
   ```bash
   oc apply -k kubernetes/operator/
   ```

3. Configure the webhooks:
   ```bash
   # Generate CA certificate and webhook certs
   ./scripts/dev-setup/generate-certs.sh

   # Apply webhook configuration
   envsubst < kubernetes/operator/webhook-setup.yaml | oc apply -f -
   ```

## Configuration

The Training Operator is configured with the following features:

- MPIJob support enabled
- TFJob support enabled
- PyTorchJob support enabled
- XGBoostJob support enabled
- Webhook validation for job specs
- Leader election for high availability

## Verification

After installation, verify that the operator is running correctly:

```bash
oc get pods -n openshift-ai -l app.kubernetes.io/name=kubeflow-training-operator
```

You should see the pod in Running state.

Check that the CRDs are installed:

```bash
oc get crd | grep kubeflow.org
```

You should see the MPIJob CRD among others.

## Troubleshooting

If the operator fails to start:

1. Check the logs:
   ```bash
   oc logs -f deployment/kubeflow-training-operator -n openshift-ai
   ```

2. Verify RBAC permissions:
   ```bash
   oc auth can-i create mpijobs --as=system:serviceaccount:openshift-ai:kubeflow-training-operator
   ```

3. Check webhook configuration:
   ```bash
   oc get mutatingwebhookconfigurations,validatingwebhookconfigurations
   ```