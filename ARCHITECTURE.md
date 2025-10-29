# MPIJobs Architecture for RedHat OpenShift AI
## Implementation Architecture using KubeFlow Trainer V2

**Version:** 1.0
**Date:** 2025-10-29
**Status:** Draft
**Target Release:** RHOAI 2.17 (Q3 2026)

---

## Executive Summary

This document provides implementation-ready architectural guidance for integrating MPIJobs support into RedHat OpenShift AI through KubeFlow Trainer V2. The architecture aligns with RHOAI's north star vision of providing enterprise-grade ML operations while maintaining hybrid cloud flexibility and open source foundations.

**Key Architectural Principles:**
- **Upstream-First**: Leverage KubeFlow Trainer V2 and MPI Operator v2 without forking
- **Enterprise-Grade**: Security, multi-tenancy, and observability as first-class concerns
- **Hybrid Cloud Native**: Identical experience across on-premises and cloud deployments
- **Unified Experience**: Seamless integration with existing RHOAI Dashboard, CLI, and SDK patterns

**Architecture at a Glance:**
- MPI Operator v2 deployed as managed component within RHOAI operator
- Three interface layers: Dashboard UI, Python SDK, CLI (all speaking to Kubernetes API)
- Standard Kubernetes primitives (Jobs, Pods, Services, Secrets) for MPI orchestration
- Prometheus/Grafana for observability, OpenShift RBAC for security

---

## 1. System Architecture

### 1.1 Component Overview

```
┌─────────────────────────────────────────────────────────────────────┐
│                    User Interaction Layer                            │
├──────────────────┬──────────────────┬──────────────────────────────┤
│  ODH Dashboard   │   Python SDK     │       CLI (odh)              │
│  (React/TypeScript│  (openshift-ai)  │    (kubectl plugin)          │
└────────┬─────────┴────────┬─────────┴───────────┬──────────────────┘
         │                  │                     │
         └──────────────────┼─────────────────────┘
                            │
                     ┌──────▼────────┐
                     │  Kubernetes   │
                     │   API Server  │
                     └──────┬────────┘
                            │
         ┌──────────────────┼─────────────────────┐
         │                  │                     │
    ┌────▼─────┐    ┌──────▼────────┐    ┌──────▼─────┐
    │ Training │    │  MPI Operator │    │   Kueue    │
    │ Operator │    │      v2       │    │ (optional) │
    │   v2     │    │  (Controller) │    └────────────┘
    └────┬─────┘    └──────┬────────┘
         │                 │
         └────────┬────────┘
                  │
         ┌────────▼─────────────────────────────────┐
         │         MPIJob Custom Resource           │
         │  (kubeflow.org/v2beta1/MPIJob)          │
         └────────┬─────────────────────────────────┘
                  │
         ┌────────▼─────────────────────────────────┐
         │        Kubernetes Resources              │
         │  ┌──────────┐  ┌──────────┐             │
         │  │ Launcher │  │ Workers  │             │
         │  │   Job    │  │ StatefulSet│            │
         │  └────┬─────┘  └────┬─────┘             │
         │       │             │                    │
         │  ┌────▼─────┐  ┌────▼─────┐             │
         │  │  Pods    │  │  Pods    │             │
         │  │ (1 pod)  │  │ (N pods) │             │
         │  └──────────┘  └──────────┘             │
         └──────────────────────────────────────────┘
                  │
         ┌────────▼─────────────────────────────────┐
         │        Infrastructure Layer              │
         │  ┌───────────┐  ┌──────────┐            │
         │  │ GPU Nodes │  │ Storage  │            │
         │  │ (NVIDIA)  │  │ (PVC/S3) │            │
         │  └───────────┘  └──────────┘            │
         └──────────────────────────────────────────┘
```

### 1.2 Integration Architecture

**Deployment Model: Sub-Controller Pattern**

The MPI Operator v2 will be deployed as a sub-controller managed by the RHOAI operator, following the established pattern used for other KubeFlow Training Operators:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: training-operator
  namespace: redhat-ods-applications
spec:
  replicas: 1
  template:
    spec:
      containers:
      - name: training-operator
        image: quay.io/opendatahub/training-operator:v2.0.0
        # Single operator manages all job types: MPIJob, PyTorchJob, TFJob
```

**Key Architectural Decisions:**

**ADR-001: Single Training Operator vs. Separate MPI Operator**
- **Decision**: Deploy KubeFlow Training Operator v2 as single controller managing all job types
- **Rationale**:
  - Upstream KubeFlow Trainer V2 unified TFJob, PyTorchJob, MPIJob, XGBoostJob under single operator
  - Reduces operational complexity (one operator to upgrade/monitor vs. multiple)
  - Enables code sharing for common functionality (job lifecycle, status reporting)
- **Trade-offs**:
  - Single point of failure (mitigated by HA deployment with leader election)
  - Larger blast radius for bugs (mitigated by comprehensive testing)
- **Alternatives Considered**:
  - Separate MPI Operator deployment: More isolation but higher operational overhead
  - Embedded within RHOAI operator: Tighter coupling, harder to upstream contributions

**ADR-002: CRD Version and API Compatibility**
- **Decision**: Use `kubeflow.org/v2beta1` API group for MPIJob
- **Rationale**:
  - Aligns with upstream KubeFlow Trainer V2 (currently v2beta1, GA v2 planned for 2026)
  - Maintains compatibility with upstream examples and documentation
  - Clear upgrade path to v2 GA without breaking changes
- **Migration Strategy**:
  - Conversion webhooks for v1alpha1 → v2beta1 (if customers have existing MPIJobs)
  - API versioning documented in compatibility matrix

### 1.3 Data Flow Architecture

**Job Submission Flow:**

```
1. User submits MPIJob via CLI/SDK/Dashboard
   ↓
2. Request to Kubernetes API Server (RBAC validation)
   ↓
3. MPIJob CR created in namespace
   ↓
4. Training Operator watches MPIJob CR
   ↓
5. Operator creates:
   - ConfigMap (MPI hostfile with worker addresses)
   - Secret (SSH keys for MPI communication)
   - Job (launcher pod)
   - StatefulSet (worker pods with predictable DNS names)
   ↓
6. Kueue (if enabled) queues job based on resource availability
   ↓
7. Kubernetes scheduler places pods on nodes (GPU affinity)
   ↓
8. Launcher pod waits for all workers to be Ready
   ↓
9. Launcher executes mpirun command
   ↓
10. MPI processes communicate across workers via SSH
   ↓
11. Training completes, launcher reports status to MPIJob CR
   ↓
12. User views results in Dashboard/CLI/SDK
```

**MPI Communication Flow:**

```
┌──────────────────────────────────────────────────────────────┐
│                   MPI Communication Mesh                      │
│                                                               │
│  ┌─────────────┐         SSH/MPI Protocol                    │
│  │  Launcher   │◄──────────────────────────────────┐         │
│  │   Pod       │                                   │         │
│  └──────┬──────┘                                   │         │
│         │ mpirun -np 4 -hostfile /etc/mpi/hostfile│         │
│         │ python train.py                          │         │
│         ▼                                          │         │
│  ┌──────────────────────────────────────────────┐ │         │
│  │         MPI Allreduce/Broadcast              │ │         │
│  └──────┬───────────────────┬───────────────┬───┘ │         │
│         │                   │               │     │         │
│  ┌──────▼─────┐     ┌──────▼─────┐  ┌──────▼─────┐ │       │
│  │  Worker-0  │◄────►│  Worker-1  │◄►│  Worker-2  │◄┘       │
│  │   Pod      │     │   Pod      │  │   Pod      │         │
│  │            │     │            │  │            │         │
│  │  Rank 0    │     │  Rank 1    │  │  Rank 2    │         │
│  └────────────┘     └────────────┘  └────────────┘         │
│                                                               │
│  Network: OpenShift SDN/OVN-Kubernetes (default)             │
│  Protocol: SSH (port 2222) for MPI process spawning          │
│  Data: TCP sockets for MPI collectives (NCCL for GPUs)       │
└──────────────────────────────────────────────────────────────┘
```

---

## 2. Technical Stack & Implementation

### 2.1 Core Components

| Component | Technology | Version | Purpose |
|-----------|-----------|---------|---------|
| **Training Operator** | Go | v2.0+ | Reconcile MPIJob CRs, manage lifecycle |
| **MPI Implementation** | Open MPI | 4.1+ | Inter-process communication runtime |
| **Container Runtime** | CRI-O / containerd | OpenShift default | Execute training pods |
| **GPU Driver** | NVIDIA GPU Operator | 23.9+ | GPU resource management |
| **Networking** | OpenShift SDN / OVN-Kubernetes | OpenShift default | Pod-to-pod communication |
| **Storage** | ODF / NFS / Ceph | Customer choice | Training data and checkpoints |

### 2.2 CLI Implementation

**Technology:** Go (kubectl plugin), distributed as `oc-odh` plugin

**Architecture Pattern:** Thin client calling Kubernetes API

```go
// pkg/cmd/mpijob/create.go
type CreateOptions struct {
    Name         string
    Namespace    string
    Workers      int
    Image        string
    Command      []string
    GPU          int
    Memory       string
    CPU          int
    client       kubernetes.Interface
}

func (o *CreateOptions) Run() error {
    mpijob := &kubeflowv2beta1.MPIJob{
        ObjectMeta: metav1.ObjectMeta{
            Name:      o.Name,
            Namespace: o.Namespace,
        },
        Spec: kubeflowv2beta1.MPIJobSpec{
            MPIReplicaSpecs: map[kubeflowv2beta1.MPIReplicaType]*kubeflowv2beta1.ReplicaSpec{
                kubeflowv2beta1.MPIReplicaTypeLauncher: {
                    Replicas: ptr.To(int32(1)),
                    Template: corev1.PodTemplateSpec{
                        Spec: corev1.PodSpec{
                            Containers: []corev1.Container{{
                                Name:    "launcher",
                                Image:   o.Image,
                                Command: o.Command,
                            }},
                        },
                    },
                },
                kubeflowv2beta1.MPIReplicaTypeWorker: {
                    Replicas: ptr.To(int32(o.Workers)),
                    Template: corev1.PodTemplateSpec{
                        Spec: corev1.PodSpec{
                            Containers: []corev1.Container{{
                                Name:  "worker",
                                Image: o.Image,
                                Resources: corev1.ResourceRequirements{
                                    Limits: corev1.ResourceList{
                                        "nvidia.com/gpu": resource.MustParse(strconv.Itoa(o.GPU)),
                                    },
                                },
                            }},
                        },
                    },
                },
            },
        },
    }

    return o.client.Create(context.Background(), mpijob)
}
```

**Command Structure:**

```bash
# Imperative commands
oc odh mpijob create <name> --workers 4 --gpu 1 --image <img> --command "python train.py"
oc odh mpijob get [name]
oc odh mpijob describe <name>
oc odh mpijob logs <name> [--launcher|--worker N] [--follow]
oc odh mpijob delete <name>

# Declarative approach
oc odh mpijob create -f mpijob.yaml
oc apply -f mpijob.yaml  # Standard kubectl also works
```

### 2.3 Python SDK Implementation

**Technology:** Python 3.9+, leverages kubernetes-client

**Distribution:** PyPI package `openshift-ai-sdk` (extends existing package)

**Architecture Pattern:** Object-oriented wrapper over Kubernetes API

```python
# openshift_ai_sdk/training/mpijob.py
from dataclasses import dataclass
from typing import List, Dict, Optional
from kubernetes import client, config

@dataclass
class ResourceRequirements:
    cpu: int = 4
    memory: str = "16Gi"
    gpu: int = 0

class MPIJob:
    """MPIJob resource for distributed MPI-based training."""

    def __init__(
        self,
        name: str,
        namespace: str = "default",
        workers: int = 2,
        image: str = None,
        command: List[str] = None,
        resources: ResourceRequirements = None,
        volume_mounts: Optional[List[Dict]] = None,
        env: Optional[Dict[str, str]] = None,
        node_selector: Optional[Dict[str, str]] = None,
    ):
        self.name = name
        self.namespace = namespace
        self.workers = workers
        self.image = image
        self.command = command or []
        self.resources = resources or ResourceRequirements()
        self.volume_mounts = volume_mounts or []
        self.env = env or {}
        self.node_selector = node_selector or {}

        config.load_incluster_config()
        self.api = client.CustomObjectsApi()
        self.core_api = client.CoreV1Api()

    def create(self) -> "MPIJob":
        """Create MPIJob in cluster."""
        mpijob_manifest = self._build_manifest()

        self.api.create_namespaced_custom_object(
            group="kubeflow.org",
            version="v2beta1",
            namespace=self.namespace,
            plural="mpijobs",
            body=mpijob_manifest,
        )
        return self

    def wait_for_completion(self, timeout: int = 3600) -> str:
        """Block until job completes or times out."""
        import time
        start_time = time.time()

        while time.time() - start_time < timeout:
            status = self.get_status()
            if status in ["Succeeded", "Failed"]:
                return status
            time.sleep(5)

        raise TimeoutError(f"Job {self.name} did not complete within {timeout}s")

    def get_status(self) -> str:
        """Get current job status."""
        mpijob = self.api.get_namespaced_custom_object(
            group="kubeflow.org",
            version="v2beta1",
            namespace=self.namespace,
            plural="mpijobs",
            name=self.name,
        )

        conditions = mpijob.get("status", {}).get("conditions", [])
        for condition in conditions:
            if condition["type"] in ["Succeeded", "Failed"]:
                return condition["type"]
        return "Running"

    def get_logs(self, worker_index: Optional[int] = None) -> str:
        """Retrieve logs from launcher or specific worker."""
        if worker_index is None:
            pod_name = f"{self.name}-launcher"
        else:
            pod_name = f"{self.name}-worker-{worker_index}"

        return self.core_api.read_namespaced_pod_log(
            name=pod_name,
            namespace=self.namespace,
        )

    def delete(self):
        """Delete MPIJob."""
        self.api.delete_namespaced_custom_object(
            group="kubeflow.org",
            version="v2beta1",
            namespace=self.namespace,
            plural="mpijobs",
            name=self.name,
        )

    @classmethod
    def get(cls, name: str, namespace: str = "default") -> "MPIJob":
        """Retrieve existing MPIJob."""
        config.load_incluster_config()
        api = client.CustomObjectsApi()

        mpijob_obj = api.get_namespaced_custom_object(
            group="kubeflow.org",
            version="v2beta1",
            namespace=namespace,
            plural="mpijobs",
            name=name,
        )

        # Reconstruct MPIJob object from cluster state
        return cls._from_manifest(mpijob_obj)

    def _build_manifest(self) -> Dict:
        """Build Kubernetes manifest for MPIJob."""
        return {
            "apiVersion": "kubeflow.org/v2beta1",
            "kind": "MPIJob",
            "metadata": {
                "name": self.name,
                "namespace": self.namespace,
            },
            "spec": {
                "mpiReplicaSpecs": {
                    "Launcher": {
                        "replicas": 1,
                        "template": {
                            "spec": {
                                "containers": [{
                                    "name": "launcher",
                                    "image": self.image,
                                    "command": self.command,
                                    "env": [{"name": k, "value": v} for k, v in self.env.items()],
                                }]
                            }
                        }
                    },
                    "Worker": {
                        "replicas": self.workers,
                        "template": {
                            "spec": {
                                "containers": [{
                                    "name": "worker",
                                    "image": self.image,
                                    "resources": {
                                        "limits": {
                                            "cpu": str(self.resources.cpu),
                                            "memory": self.resources.memory,
                                            "nvidia.com/gpu": str(self.resources.gpu),
                                        }
                                    },
                                    "volumeMounts": self.volume_mounts,
                                }],
                                "nodeSelector": self.node_selector,
                            }
                        }
                    }
                }
            }
        }
```

**Usage Example:**

```python
from openshift_ai_sdk.training import MPIJob, ResourceRequirements

# Create and submit job
job = MPIJob(
    name="pytorch-distributed",
    namespace="my-project",
    workers=4,
    image="quay.io/my-org/pytorch-horovod:latest",
    command=["mpirun", "python", "train.py", "--epochs", "10"],
    resources=ResourceRequirements(cpu=8, memory="64Gi", gpu=1),
    volume_mounts=[{
        "name": "training-data",
        "mountPath": "/data"
    }],
)

job.create()

# Monitor asynchronously
status = job.get_status()
print(f"Job status: {status}")

# Or block until completion
final_status = job.wait_for_completion(timeout=7200)
print(f"Job completed with status: {final_status}")

# Retrieve logs
logs = job.get_logs()  # Launcher logs
worker_logs = job.get_logs(worker_index=0)  # Worker 0 logs
```

### 2.4 Dashboard UI Implementation

**Technology:** React 18+, TypeScript, PatternFly 5 (RHOAI standard)

**Architecture Pattern:** REST API backend (Flask/Python) + React frontend

**Component Hierarchy:**

```
src/pages/
├── TrainingJobs/
│   ├── MPIJobList.tsx           # List view with filtering
│   ├── MPIJobDetails.tsx        # Details page with tabs
│   ├── MPIJobCreateWizard.tsx   # Multi-step creation form
│   ├── MPIJobLogs.tsx           # Log viewer with streaming
│   ├── MPIJobMetrics.tsx        # Grafana dashboard embed
│   └── components/
│       ├── JobStatusBadge.tsx
│       ├── WorkerPodTable.tsx
│       └── ResourceChart.tsx
```

**Key Features:**

1. **Unified Job List**: MPIJobs displayed alongside TFJob, PyTorchJob
2. **Job Creation Wizard**:
   - Step 1: Basic info (name, namespace)
   - Step 2: Container image and command
   - Step 3: Resources (workers, CPU, memory, GPU)
   - Step 4: Advanced (volumes, env vars, node selector)
   - Step 5: Review YAML
3. **Details Page Tabs**:
   - Overview: Status, duration, configuration summary
   - Configuration: Full YAML view (editable for resubmission)
   - Workers: Pod list with status, node placement, resource usage
   - Logs: Aggregated logs with worker selection and filtering
   - Metrics: Embedded Grafana panel with GPU utilization
   - Events: Kubernetes events related to job

**API Endpoints (Backend):**

```python
# dashboard/backend/app/routes/training.py
from flask import Blueprint, jsonify, request
from kubernetes import client

training_bp = Blueprint('training', __name__)

@training_bp.route('/api/mpijobs', methods=['GET'])
def list_mpijobs():
    """List all MPIJobs in accessible namespaces."""
    namespace = request.args.get('namespace', 'default')
    api = client.CustomObjectsApi()

    mpijobs = api.list_namespaced_custom_object(
        group="kubeflow.org",
        version="v2beta1",
        namespace=namespace,
        plural="mpijobs",
    )

    return jsonify({
        'items': [_serialize_mpijob(job) for job in mpijobs['items']]
    })

@training_bp.route('/api/mpijobs/<name>', methods=['GET'])
def get_mpijob(name):
    """Get specific MPIJob details."""
    namespace = request.args.get('namespace', 'default')
    api = client.CustomObjectsApi()

    mpijob = api.get_namespaced_custom_object(
        group="kubeflow.org",
        version="v2beta1",
        namespace=namespace,
        plural="mpijobs",
        name=name,
    )

    return jsonify(_serialize_mpijob(mpijob))

@training_bp.route('/api/mpijobs', methods=['POST'])
def create_mpijob():
    """Create new MPIJob."""
    manifest = request.json
    namespace = manifest['metadata']['namespace']

    api = client.CustomObjectsApi()
    mpijob = api.create_namespaced_custom_object(
        group="kubeflow.org",
        version="v2beta1",
        namespace=namespace,
        plural="mpijobs",
        body=manifest,
    )

    return jsonify(_serialize_mpijob(mpijob)), 201

@training_bp.route('/api/mpijobs/<name>/logs', methods=['GET'])
def get_mpijob_logs(name):
    """Stream logs from launcher or worker pods."""
    namespace = request.args.get('namespace', 'default')
    worker_index = request.args.get('worker')

    core_api = client.CoreV1Api()

    if worker_index:
        pod_name = f"{name}-worker-{worker_index}"
    else:
        pod_name = f"{name}-launcher"

    logs = core_api.read_namespaced_pod_log(
        name=pod_name,
        namespace=namespace,
        follow=False,
        tail_lines=1000,
    )

    return jsonify({'logs': logs})
```

**Frontend Component Example:**

```typescript
// src/pages/TrainingJobs/MPIJobDetails.tsx
import React, { useEffect, useState } from 'react';
import {
  Page,
  PageSection,
  Tabs,
  Tab,
  Title,
  Badge,
} from '@patternfly/react-core';
import { MPIJob } from '../../types/training';
import { fetchMPIJob } from '../../api/training';

const MPIJobDetails: React.FC<{ name: string; namespace: string }> = ({ name, namespace }) => {
  const [job, setJob] = useState<MPIJob | null>(null);
  const [activeTabKey, setActiveTabKey] = useState<string>('overview');

  useEffect(() => {
    const loadJob = async () => {
      const data = await fetchMPIJob(name, namespace);
      setJob(data);
    };

    loadJob();
    const interval = setInterval(loadJob, 5000); // Refresh every 5s
    return () => clearInterval(interval);
  }, [name, namespace]);

  if (!job) return <div>Loading...</div>;

  return (
    <Page>
      <PageSection variant="light">
        <Title headingLevel="h1">
          {job.metadata.name}
          <Badge style={{ marginLeft: '1rem' }}>
            {job.status.conditions?.[0]?.type || 'Pending'}
          </Badge>
        </Title>
      </PageSection>

      <PageSection>
        <Tabs activeKey={activeTabKey} onSelect={(_, key) => setActiveTabKey(key as string)}>
          <Tab eventKey="overview" title="Overview">
            <OverviewTab job={job} />
          </Tab>
          <Tab eventKey="configuration" title="Configuration">
            <ConfigurationTab job={job} />
          </Tab>
          <Tab eventKey="workers" title="Workers">
            <WorkersTab job={job} />
          </Tab>
          <Tab eventKey="logs" title="Logs">
            <LogsTab name={name} namespace={namespace} />
          </Tab>
          <Tab eventKey="metrics" title="Metrics">
            <MetricsTab job={job} />
          </Tab>
          <Tab eventKey="events" title="Events">
            <EventsTab name={name} namespace={namespace} />
          </Tab>
        </Tabs>
      </PageSection>
    </Page>
  );
};
```

---

## 3. Networking Architecture

### 3.1 MPI Communication Patterns

**Default Network Stack:**

```
┌──────────────────────────────────────────────────────────┐
│                Application Layer                          │
│  MPI_Send / MPI_Recv / MPI_Allreduce (user code)        │
└───────────────────────┬──────────────────────────────────┘
                        │
┌───────────────────────▼──────────────────────────────────┐
│                  MPI Library                             │
│  Open MPI 4.1+ (process management, collectives)        │
└───────────────────────┬──────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        │               │               │
┌───────▼──────┐ ┌─────▼──────┐ ┌──────▼─────────┐
│     SSH      │ │   TCP/IP   │ │  NCCL (GPU)    │
│  (spawning)  │ │ (data xfer)│ │  (collectives) │
└───────┬──────┘ └─────┬──────┘ └──────┬─────────┘
        │               │               │
        └───────────────┼───────────────┘
                        │
┌───────────────────────▼──────────────────────────────────┐
│           OpenShift SDN / OVN-Kubernetes                 │
│  Pod-to-pod networking (CNI plugin)                      │
└───────────────────────┬──────────────────────────────────┘
                        │
┌───────────────────────▼──────────────────────────────────┐
│              Physical Network                            │
│  1/10/25/100 GbE (typical enterprise)                    │
└──────────────────────────────────────────────────────────┘
```

**Key Network Requirements:**

1. **Pod-to-Pod Connectivity**: All worker pods must reach each other on SSH port (2222)
2. **DNS Resolution**: StatefulSet ensures predictable pod DNS names for hostfile generation
3. **Bandwidth**: Minimum 10 GbE recommended for multi-node training (NCCL saturates 25 GbE for GPU collectives)
4. **Latency**: <1ms intra-rack, <5ms inter-rack for acceptable MPI performance

### 3.2 SSH Key Management

**Architecture Pattern: Ephemeral per-job SSH keys**

```
Job Creation Flow:
1. Training Operator generates SSH keypair
2. Private key → Secret (launcher-ssh-key)
3. Public key → ConfigMap (worker-authorized-keys)
4. Launcher mounts private key at /root/.ssh/id_rsa
5. Workers mount authorized_keys at /root/.ssh/authorized_keys
6. Launcher SSHs to workers to spawn MPI processes
7. Job completion → Operator deletes Secret and ConfigMap
```

**Security Characteristics:**
- Keys are RSA 4096-bit (FIPS 140-2 compliant)
- Keys scoped to single job (no cross-job access)
- Automatic cleanup on job deletion
- No key rotation needed (ephemeral lifecycle)

**ADR-003: SSH vs. PMIx Process Management**
- **Decision**: Use SSH for MVP, evaluate PMIx post-MVP
- **Rationale**:
  - SSH is battle-tested, supported by all MPI implementations
  - PMIx (Process Management Interface Exascale) emerging standard with less tooling maturity
  - Container ecosystem tooling (kubectl exec, logs) works well with SSH-based MPI
- **Future**: PMIx may eliminate SSH dependency, improve startup time by 20-30%

### 3.3 Network Policy Design

**Default NetworkPolicy for MPIJobs:**

```yaml
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mpijob-{job-name}
  namespace: {namespace}
spec:
  podSelector:
    matchLabels:
      training.kubeflow.org/job-name: {job-name}
  policyTypes:
  - Ingress
  - Egress
  ingress:
  # Allow launcher to SSH into workers
  - from:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-name: {job-name}
          training.kubeflow.org/replica-type: launcher
    ports:
    - protocol: TCP
      port: 2222
  # Allow worker-to-worker MPI communication
  - from:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-name: {job-name}
          training.kubeflow.org/replica-type: worker
    ports:
    - protocol: TCP
      port: 1024-65535  # Dynamic MPI ports
  egress:
  # Allow DNS resolution
  - to:
    - namespaceSelector:
        matchLabels:
          name: openshift-dns
    ports:
    - protocol: UDP
      port: 53
  # Allow access to container registry (for image pulls)
  - to:
    - namespaceSelector: {}
    ports:
    - protocol: TCP
      port: 443
  # Allow worker-to-worker communication
  - to:
    - podSelector:
        matchLabels:
          training.kubeflow.org/job-name: {job-name}
```

**Multi-Tenancy Isolation:**
- NetworkPolicies automatically created per MPIJob
- Pods in different namespaces cannot communicate (namespace isolation)
- Pods in same namespace but different jobs cannot communicate (label-based isolation)

### 3.4 GPU-Optimized Networking (Post-MVP)

**RDMA Support Architecture:**

For high-performance GPU collectives (InfiniBand/RoCE):

```yaml
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: llm-training-rdma
spec:
  mpiReplicaSpecs:
    Worker:
      template:
        metadata:
          annotations:
            k8s.v1.cni.cncf.io/networks: |
              [
                {
                  "name": "ib-network",
                  "namespace": "default",
                  "ips": ["192.168.1.10"]
                }
              ]
        spec:
          containers:
          - name: worker
            resources:
              limits:
                rdma/hca: 1  # Request RDMA device
            volumeMounts:
            - name: rdma-cm
              mountPath: /dev/infiniband
          volumes:
          - name: rdma-cm
            hostPath:
              path: /dev/infiniband
```

**Requirements:**
- Multus CNI for secondary network interfaces
- SR-IOV Network Operator for InfiniBand device passthrough
- RDMA-enabled container images with `libibverbs`

**Performance Expectations:**
- Standard TCP/IP: ~12 GB/s A100 GPU-to-GPU
- RDMA/InfiniBand: ~200 GB/s (16x improvement)

---

## 4. Resource Management

### 4.1 Scheduling Architecture

**Kueue Integration for Fair-Share Scheduling:**

```
┌──────────────────────────────────────────────────────────┐
│                   User Submits MPIJob                     │
└────────────────────────┬─────────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────────┐
│              Training Operator Creates:                   │
│  • WorkloadPriorityClass (maps to Kueue priority)        │
│  • PodGroup (gang scheduling constraint)                 │
└────────────────────────┬─────────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────────┐
│         Kueue Admission Check                            │
│  1. Check ClusterQueue resource availability             │
│  2. Check ResourceQuota limits                           │
│  3. Apply fair-share algorithm (weighted across teams)   │
│  4. Admit or Queue based on priority and availability    │
└────────────────────────┬─────────────────────────────────┘
                         │
         ┌───────────────┴────────────────┐
         │                                │
    ┌────▼──────┐                   ┌────▼──────┐
    │  Admitted │                   │   Queued  │
    │  (Running)│                   │  (Pending)│
    └────┬──────┘                   └────┬──────┘
         │                                │
         │                                │
         │         ┌──────────────────────┘
         │         │ Resource becomes available
         │         │ or higher priority job completes
         │         │
         └─────────▼────────────────────────────────────┐
                   │ Kubernetes Scheduler               │
                   │ Places pods on nodes with GPUs     │
                   └────────────────────────────────────┘
```

**ADR-004: Gang Scheduling Requirement**
- **Decision**: Require Kueue for production MPIJobs (MVP)
- **Rationale**:
  - MPI training requires all workers to start simultaneously
  - Without gang scheduling, partial allocations cause deadlock:
    - Job A gets 20/32 GPUs, waits for 12 more
    - Job B gets 30/32 GPUs, waits for 2 more
    - Cluster has 64 GPUs total → deadlock
  - Kueue provides PodGroup semantics for all-or-nothing scheduling
- **Configuration**:
  ```yaml
  apiVersion: kueue.x-k8s.io/v1beta1
  kind: ResourceFlavor
  metadata:
    name: gpu-flavor
  spec:
    nodeLabels:
      nvidia.com/gpu.present: "true"
  ---
  apiVersion: kueue.x-k8s.io/v1beta1
  kind: ClusterQueue
  metadata:
    name: gpu-cluster-queue
  spec:
    namespaceSelector: {}
    cohort: default-cohort
    preemption:
      withinClusterQueue: LowerPriority
      reclaimWithinCohort: Any
    resourceGroups:
    - coveredResources: ["cpu", "memory", "nvidia.com/gpu"]
      flavors:
      - name: gpu-flavor
        resources:
        - name: "cpu"
          nominalQuota: 1000
        - name: "memory"
          nominalQuota: 4000Gi
        - name: "nvidia.com/gpu"
          nominalQuota: 128
  ```

**Alternative (Post-MVP): Volcano/Coscheduler**
- Lighter-weight gang scheduling option
- Less sophisticated fair-share algorithms than Kueue
- Consider for customers not using Kueue elsewhere

### 4.2 Resource Quota Enforcement

**Three-Level Quota Hierarchy:**

```
1. Cluster-Level: ClusterQueue (Kueue)
   Total GPUs available for MPIJobs across all namespaces

2. Namespace-Level: ResourceQuota (Kubernetes)
   Maximum resources per team/project

3. Job-Level: LimitRange (Kubernetes)
   Min/max resources per individual MPIJob
```

**Example Configuration:**

```yaml
# Namespace quota: team-cv can use max 32 GPUs
apiVersion: v1
kind: ResourceQuota
metadata:
  name: team-cv-quota
  namespace: team-cv
spec:
  hard:
    requests.nvidia.com/gpu: "32"
    requests.cpu: "256"
    requests.memory: "2048Gi"
---
# LimitRange: No single MPIJob can exceed 16 GPUs
apiVersion: v1
kind: LimitRange
metadata:
  name: mpijob-limits
  namespace: team-cv
spec:
  limits:
  - type: "Container"
    max:
      nvidia.com/gpu: "16"
      memory: "256Gi"
    default:
      nvidia.com/gpu: "1"
      memory: "16Gi"
```

**Error Handling:**

When quota exceeded, user receives clear error:

```bash
$ oc odh mpijob create large-training --workers 20 --gpu 2

Error: Quota exceeded
Details: Requested 40 GPUs (20 workers × 2 GPUs), but namespace quota is 32 GPUs.
Current usage: 24 GPUs (3 running jobs)
Available: 8 GPUs

Suggestions:
1. Reduce worker count to 16 (32 GPUs total)
2. Wait for existing jobs to complete
3. Request quota increase from administrator
```

### 4.3 GPU Node Affinity

**Topology-Aware Scheduling:**

For optimal performance, schedule all workers in same rack/zone:

```yaml
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: topology-aware-job
spec:
  mpiReplicaSpecs:
    Worker:
      template:
        spec:
          affinity:
            podAffinity:
              # Prefer to schedule workers on same rack
              preferredDuringSchedulingIgnoredDuringExecution:
              - weight: 100
                podAffinityTerm:
                  labelSelector:
                    matchLabels:
                      training.kubeflow.org/job-name: topology-aware-job
                  topologyKey: topology.kubernetes.io/rack
            nodeAffinity:
              # Require GPU nodes
              requiredDuringSchedulingIgnoredDuringExecution:
                nodeSelectorTerms:
                - matchExpressions:
                  - key: nvidia.com/gpu.present
                    operator: In
                    values:
                    - "true"
```

**GPU Type Selection:**

For heterogeneous GPU clusters (A100 + V100):

```yaml
spec:
  mpiReplicaSpecs:
    Worker:
      template:
        spec:
          nodeSelector:
            nvidia.com/gpu.product: "NVIDIA-A100-SXM4-80GB"
```

---

## 5. Observability Architecture

### 5.1 Metrics Collection

**Prometheus Metrics Exposed by Training Operator:**

```
# Job-level metrics
training_operator_mpijob_created_total{namespace="default"} 45
training_operator_mpijob_running_total{namespace="default"} 3
training_operator_mpijob_succeeded_total{namespace="default"} 38
training_operator_mpijob_failed_total{namespace="default"} 4
training_operator_mpijob_queued_total{namespace="default"} 2

# Job duration histogram
training_operator_mpijob_duration_seconds_bucket{namespace="default",le="600"} 12
training_operator_mpijob_duration_seconds_bucket{namespace="default",le="3600"} 35
training_operator_mpijob_duration_seconds_bucket{namespace="default",le="7200"} 42

# Worker pod metrics (via kube-state-metrics)
kube_pod_container_resource_requests{pod=~".*-worker-.*",resource="nvidia.com/gpu"} 1
kube_pod_container_resource_requests{pod=~".*-worker-.*",resource="memory"} 68719476736

# GPU utilization (via NVIDIA DCGM exporter)
DCGM_FI_DEV_GPU_UTIL{pod="my-job-worker-0",gpu="0"} 92.5
DCGM_FI_DEV_FB_USED{pod="my-job-worker-0",gpu="0"} 72000000000  # 72 GB
```

**Grafana Dashboard Panels:**

```
Dashboard: "MPIJob Overview"
├── Row 1: Cluster-Wide Metrics
│   ├── Panel: Total MPIJobs (Running, Succeeded, Failed, Queued)
│   ├── Panel: GPU Utilization (cluster-wide %)
│   └── Panel: Average Job Duration (last 7 days)
├── Row 2: Per-Job Metrics (select job from dropdown)
│   ├── Panel: Job Status Timeline
│   ├── Panel: Worker Pod Status (table)
│   ├── Panel: GPU Utilization per Worker (graph)
│   └── Panel: Network Throughput (MPI communication)
└── Row 3: Resource Usage
    ├── Panel: CPU Usage per Worker
    ├── Panel: Memory Usage per Worker
    └── Panel: Disk I/O (training data loading)
```

**PromQL Example Queries:**

```promql
# GPU utilization across all workers of a job
avg(DCGM_FI_DEV_GPU_UTIL{
  pod=~"my-job-worker-.*"
}) by (pod)

# Scaling efficiency: actual speedup vs. linear expectation
# (requires custom metric from training script)
(
  rate(training_samples_processed_total{job="single-node"}[5m])
  /
  rate(training_samples_processed_total{job="8-node"}[5m])
) / 8

# Network bandwidth between workers
sum(rate(container_network_transmit_bytes_total{
  pod=~"my-job-worker-.*"
}[1m])) by (pod)
```

### 5.2 Logging Architecture

**Log Aggregation Flow:**

```
┌──────────────────────────────────────────────────────────┐
│         MPIJob Pods (Launcher + Workers)                 │
│  stdout/stderr → Container runtime (CRI-O)               │
└────────────────────────┬─────────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────────┐
│         OpenShift Logging Stack                          │
│  Collector: Fluentd / Vector                             │
│  Storage: Elasticsearch / Loki                           │
│  Query: Kibana / Grafana Explore                         │
└────────────────────────┬─────────────────────────────────┘
                         │
┌────────────────────────▼─────────────────────────────────┐
│         ODH Dashboard / CLI                              │
│  Query logs via Elasticsearch API                        │
│  Display synchronized logs from launcher + workers       │
└──────────────────────────────────────────────────────────┘
```

**Structured Logging Format:**

Training scripts should emit JSON logs for structured parsing:

```json
{
  "timestamp": "2025-10-29T14:32:01.123Z",
  "level": "INFO",
  "component": "trainer",
  "mpi_rank": 0,
  "epoch": 5,
  "batch": 1200,
  "loss": 0.0234,
  "gpu_mem_used_gb": 67.2,
  "message": "Training progress"
}
```

**Dashboard Log Viewer Features:**

1. **Synchronized Timeline**: Logs from all workers aligned by timestamp
2. **Rank Filtering**: Show logs from specific MPI ranks
3. **Severity Filtering**: ERROR, WARN, INFO, DEBUG
4. **Keyword Search**: Full-text search across all logs
5. **Export**: Download logs as JSON/TXT for offline analysis

### 5.3 Distributed Tracing (Post-MVP)

**OpenTelemetry Integration:**

For deep performance analysis of MPI collectives:

```python
# training script with tracing
from opentelemetry import trace
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.exporter.otlp.proto.grpc.trace_exporter import OTLPSpanExporter

tracer = trace.get_tracer(__name__)

with tracer.start_as_current_span("training_epoch"):
    for batch in dataloader:
        with tracer.start_as_current_span("forward_pass"):
            outputs = model(batch)

        with tracer.start_as_current_span("backward_pass"):
            loss.backward()

        with tracer.start_as_current_span("mpi_allreduce"):
            # Horovod/DDP gradient synchronization
            optimizer.synchronize()
```

**Trace Analysis:**
- Identify slow MPI collectives (Allreduce latency spikes)
- Detect imbalanced workloads (rank 0 fast, rank 7 slow)
- Visualize end-to-end training iteration breakdown

---

## 6. Security Architecture

### 6.1 RBAC Model

**Role Hierarchy:**

```
├── mpijob-admin (ClusterRole)
│   ├── Create/delete MPIJobs in any namespace
│   ├── View all MPIJobs cluster-wide
│   └── Manage MPI Operator configuration
│
├── mpijob-user (ClusterRole, namespace-scoped binding)
│   ├── Create/delete MPIJobs in assigned namespaces
│   ├── View own MPIJobs
│   └── Read logs from own job pods
│
└── mpijob-viewer (ClusterRole, namespace-scoped binding)
    ├── View MPIJobs in assigned namespaces
    └── Read logs (no create/delete)
```

**RBAC Definitions:**

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: mpijob-user
rules:
# MPIJob CRUD
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs"]
  verbs: ["create", "delete", "get", "list", "watch"]
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs/status"]
  verbs: ["get", "list", "watch"]

# Pod logs access (for jobs they own)
- apiGroups: [""]
  resources: ["pods/log"]
  verbs: ["get", "list"]

# ConfigMaps/Secrets (for volume mounts)
- apiGroups: [""]
  resources: ["configmaps", "secrets"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: data-scientist-mpijob-user
  namespace: team-cv
subjects:
- kind: Group
  name: data-scientists
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: mpijob-user
  apiGroup: rbac.authorization.k8s.io
```

**ServiceAccount Security:**

MPIJob pods run with restricted ServiceAccounts:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: mpijob-launcher
  namespace: team-cv
automountServiceAccountToken: false  # No K8s API access needed
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: mpijob-worker
  namespace: team-cv
automountServiceAccountToken: false
```

**ADR-005: Least Privilege for Job Pods**
- **Decision**: MPIJob pods have no Kubernetes API access by default
- **Rationale**:
  - Training workloads don't need cluster-wide permissions
  - Reduces blast radius of compromised pod
  - Aligns with zero-trust security model
- **Exception**: If training script needs to interact with K8s API (e.g., push metrics to custom resource), user explicitly grants permissions via RoleBinding

### 6.2 Multi-Tenancy Isolation

**Namespace-Level Isolation:**

```
Namespace: team-cv
├── ResourceQuota: 32 GPUs max
├── NetworkPolicy: Deny cross-namespace traffic
├── RoleBinding: data-scientists group → mpijob-user role
└── LimitRange: Max 16 GPUs per MPIJob

Namespace: team-nlp
├── ResourceQuota: 32 GPUs max
├── NetworkPolicy: Deny cross-namespace traffic
├── RoleBinding: nlp-engineers group → mpijob-user role
└── LimitRange: Max 16 GPUs per MPIJob
```

**Enforcement Mechanisms:**

1. **RBAC**: Users can only create MPIJobs in namespaces where they have RoleBinding
2. **ResourceQuota**: Cluster-level enforcement of GPU limits per namespace
3. **NetworkPolicy**: Pods in team-cv cannot connect to pods in team-nlp
4. **PVC Isolation**: PVCs in team-cv not mountable by pods in team-nlp

**Dashboard Multi-Tenancy:**

- User sees only MPIJobs in namespaces they have RBAC access to
- Admin view shows all namespaces (requires `mpijob-admin` ClusterRole)
- Namespace selector in UI dynamically populated based on user permissions

### 6.3 Secret Management

**Training Script Credentials:**

For accessing external services (S3, databases):

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: training-secrets
  namespace: team-cv
type: Opaque
data:
  S3_ACCESS_KEY: <base64-encoded>
  S3_SECRET_KEY: <base64-encoded>
---
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: s3-training
spec:
  mpiReplicaSpecs:
    Worker:
      template:
        spec:
          containers:
          - name: worker
            envFrom:
            - secretRef:
                name: training-secrets
```

**Best Practices:**

1. **External Secret Managers**: Integrate with Vault, AWS Secrets Manager
   ```yaml
   apiVersion: external-secrets.io/v1beta1
   kind: ExternalSecret
   metadata:
     name: training-secrets
   spec:
     secretStoreRef:
       name: vault-backend
       kind: ClusterSecretStore
     target:
       name: training-secrets
     data:
     - secretKey: S3_ACCESS_KEY
       remoteRef:
         key: team-cv/s3-credentials
         property: access_key
   ```

2. **Avoid hardcoding in images**: Never bake credentials into container images
3. **Audit logging**: Track secret access via OpenShift audit logs
4. **Rotation**: Implement secret rotation policies (30-90 days)

### 6.4 Compliance & Audit

**HIPAA Compliance Example:**

```yaml
# Encrypted PVC for PHI data
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: encrypted-mri-data
  namespace: medical-research
  annotations:
    volume.beta.kubernetes.io/storage-class: encrypted-cephfs
spec:
  accessModes:
  - ReadWriteMany
  resources:
    requests:
      storage: 2Ti
---
# MPIJob with audit annotations
apiVersion: kubeflow.org/v2beta1
kind: MPIJob
metadata:
  name: hipaa-compliant-training
  namespace: medical-research
  annotations:
    compliance.openshift.io/hipaa: "true"
    compliance.openshift.io/data-classification: "PHI"
    audit.openshift.io/log-to-siem: "true"
spec:
  mpiReplicaSpecs:
    Worker:
      template:
        metadata:
          annotations:
            # Ensure pods scheduled on compliant nodes
            compliance.openshift.io/require-node-labels: "compliance=hipaa"
        spec:
          nodeSelector:
            compliance: hipaa
          volumes:
          - name: encrypted-data
            persistentVolumeClaim:
              claimName: encrypted-mri-data
```

**Audit Log Example:**

```json
{
  "kind": "Event",
  "apiVersion": "audit.k8s.io/v1",
  "level": "RequestResponse",
  "auditID": "abc123",
  "stage": "ResponseComplete",
  "requestURI": "/apis/kubeflow.org/v2beta1/namespaces/medical-research/mpijobs",
  "verb": "create",
  "user": {
    "username": "scientist-jane",
    "groups": ["medical-researchers", "system:authenticated"]
  },
  "objectRef": {
    "resource": "mpijobs",
    "namespace": "medical-research",
    "name": "hipaa-compliant-training"
  },
  "responseStatus": {
    "code": 201
  },
  "requestReceivedTimestamp": "2025-10-29T14:32:01.123Z"
}
```

---

## 7. Key Architectural Decisions (ADRs)

### ADR-001: Single Training Operator vs. Separate MPI Operator
**Status:** Accepted
**Decision:** Deploy KubeFlow Training Operator v2 as single controller managing all job types (MPIJob, PyTorchJob, TFJob).
**Context:** KubeFlow Trainer V2 unified multiple operators. Need to decide on deployment model for RHOAI.
**Consequences:**
- **Pros:** Reduced operational complexity, code sharing, upstream alignment
- **Cons:** Single point of failure, larger blast radius for bugs
- **Mitigation:** HA deployment with leader election, comprehensive testing

### ADR-002: CRD Version and API Compatibility
**Status:** Accepted
**Decision:** Use `kubeflow.org/v2beta1` API group for MPIJob MVP.
**Context:** Upstream KubeFlow Trainer V2 currently at v2beta1, GA v2 planned for 2026.
**Consequences:**
- **Pros:** Upstream compatibility, clear upgrade path, community ecosystem
- **Cons:** Beta API may have breaking changes before GA
- **Mitigation:** Conversion webhooks for API versioning, compatibility matrix documentation

### ADR-003: SSH vs. PMIx Process Management
**Status:** Accepted for MVP, PMIx evaluation post-MVP
**Decision:** Use SSH for MPI process spawning in MVP.
**Context:** MPI processes need to be launched across worker pods. SSH is standard but PMIx is emerging.
**Consequences:**
- **Pros:** Battle-tested, wide tool support, mature ecosystem
- **Cons:** SSH adds ~5-10s startup overhead, requires key management
- **Mitigation:** Ephemeral keys per job, evaluate PMIx post-MVP for 20-30% startup improvement

### ADR-004: Gang Scheduling Requirement
**Status:** Accepted
**Decision:** Require Kueue for production MPIJobs (MVP).
**Context:** MPI training requires all workers to start atomically to avoid deadlock.
**Consequences:**
- **Pros:** Prevents resource deadlock, fair-share scheduling, priority/preemption
- **Cons:** Additional dependency on Kueue operator
- **Mitigation:** Kueue becoming standard in OpenShift AI stack, document as prerequisite

### ADR-005: Least Privilege for Job Pods
**Status:** Accepted
**Decision:** MPIJob pods have no Kubernetes API access by default.
**Context:** Security best practice for zero-trust environments.
**Consequences:**
- **Pros:** Reduced blast radius, compliance-friendly, principle of least privilege
- **Cons:** Users must explicitly grant permissions if training script needs K8s API
- **Mitigation:** Documentation and examples for common use cases requiring API access

### ADR-006: Dashboard UI Framework
**Status:** Accepted
**Decision:** Use PatternFly 5 (React) for Dashboard UI, consistent with existing RHOAI Dashboard.
**Context:** Need to integrate MPIJob UI into existing ODH Dashboard.
**Consequences:**
- **Pros:** Consistent UX, shared component library, accessibility built-in
- **Cons:** Learning curve for contributors not familiar with PatternFly
- **Mitigation:** Component templates and style guide for MPIJob pages

### ADR-007: SDK API Design Pattern
**Status:** Accepted
**Decision:** Object-oriented SDK with `MPIJob` class, methods for lifecycle operations.
**Context:** Need Pythonic API for data scientists, balancing simplicity with flexibility.
**Consequences:**
- **Pros:** Intuitive for Python developers, type hints for IDE support, async patterns
- **Cons:** Abstraction layer may hide underlying Kubernetes complexity
- **Mitigation:** Advanced users can still use kubectl/raw Kubernetes API, SDK is additive

### ADR-008: Networking Approach (MVP)
**Status:** Accepted for MVP, RDMA post-MVP
**Decision:** Use standard OpenShift SDN/OVN-Kubernetes for pod-to-pod communication in MVP.
**Context:** Need to balance simplicity (standard networking) with performance (RDMA).
**Consequences:**
- **Pros:** Works out-of-box on all OpenShift clusters, no special hardware required
- **Cons:** Lower bandwidth than RDMA (10-25 GbE vs. 200 GB/s InfiniBand)
- **Mitigation:** Document RDMA setup guide for post-MVP, sufficient for 80% of use cases

### ADR-009: Log Aggregation Strategy
**Status:** Accepted
**Decision:** Leverage existing OpenShift logging stack (Fluentd/Loki), no custom logging infrastructure.
**Context:** MPIJobs generate logs from launcher + N workers, need centralized view.
**Consequences:**
- **Pros:** Consistent with other OpenShift workloads, no new infrastructure
- **Cons:** Log volume can be high for large MPIJobs (100+ workers)
- **Mitigation:** Log rotation policies, structured logging for filtering, retention limits

### ADR-010: Metrics Collection Architecture
**Status:** Accepted
**Decision:** Use Prometheus for metrics, NVIDIA DCGM exporter for GPU metrics, Grafana for visualization.
**Context:** Need to expose job-level and pod-level metrics for observability.
**Consequences:**
- **Pros:** Industry standard, integrates with existing RHOAI monitoring
- **Cons:** DCGM exporter requires daemonset on all GPU nodes
- **Mitigation:** DCGM exporter packaged as part of NVIDIA GPU Operator (already required)

---

## 8. Implementation Phases

### Phase 1: Foundation (Q4 2025)
**Goal:** KubeFlow Trainer V2 integration into RHOAI

**Deliverables:**
- Training Operator v2 packaged as RHOAI component
- CRDs installed (MPIJob, PyTorchJob, TFJob)
- Operator lifecycle management via RHOAI operator
- Basic CLI commands (`oc odh mpijob create/get/delete`)

**Success Criteria:**
- User can submit MPIJob via kubectl/oc
- Launcher and worker pods created successfully
- MPI communication established (hello world MPI program)

### Phase 2: MVP (Q1 2026)
**Goal:** Complete CLI, SDK, Dashboard UI for MPIJobs

**Deliverables:**
- Full CLI implementation (all commands from RFE R2)
- Python SDK with `MPIJob` class (RFE R3)
- Dashboard UI: List view, create wizard, details page (RFE R4)
- Observability: Prometheus metrics, Grafana dashboards (RFE R5)
- Security: RBAC roles, NetworkPolicies (RFE R6)
- Documentation: Quickstart, API reference, troubleshooting (RFE R8)

**Success Criteria:**
- All 10 Acceptance Criteria from RFE validated
- Private beta with 5 strategic customers
- End-to-end: Submit MPIJob from Dashboard, monitor in Grafana, view logs, retrieve trained model

### Phase 3: Beta (Q2 2026)
**Goal:** Public beta, production readiness

**Deliverables:**
- Kueue integration for gang scheduling (RFE R7, ADR-004)
- HA deployment for Training Operator
- Performance benchmarks published (ResNet-50, BERT, GPT-7B)
- Customer migration guides (from SageMaker, Azure ML, standalone KubeFlow)
- Field enablement: Sales kit, SA training, partner training

**Success Criteria:**
- 10+ beta customers running production workloads
- Zero P1/P2 bugs for 4 consecutive weeks
- NPS >50 from beta users
- Performance benchmarks meet 85% scaling efficiency target

### Phase 4: GA (Q3 2026)
**Goal:** General availability in RHOAI 2.17

**Deliverables:**
- Full RedHat support coverage
- SLA guarantees (99.9% operator uptime)
- Production hardening (tested with 512 GPU clusters)
- Compliance certifications (FedRAMP, PCI-DSS documentation)

**Success Criteria:**
- Passes RedHat GA readiness checklist
- Zero known P1/P2 bugs
- Documentation completeness >95% (automated checks)
- 30-minute time-to-first-job validated with 20+ users

### Phase 5: Post-GA Enhancements (Q4 2026+)
**Goal:** Advanced features based on customer feedback

**Deliverables:**
- Katib integration for hyperparameter optimization (RFE R10)
- Model Registry integration (RFE R11)
- RDMA/InfiniBand support (RFE R12, ADR-008)
- Checkpoint and fault tolerance (RFE R13)
- Elastic training (dynamic worker scaling)
- PMIx process management (ADR-003)

---

## 9. Testing Strategy

### 9.1 Unit Testing

**Components Under Test:**
- CLI: Command parsing, argument validation, API client logic
- SDK: MPIJob class methods, resource builders, error handling
- Dashboard Backend: API endpoints, RBAC checks, data serialization

**Coverage Target:** 80%+

**Example: SDK Unit Test**

```python
# tests/sdk/test_mpijob.py
import pytest
from unittest.mock import Mock, patch
from openshift_ai_sdk.training import MPIJob, ResourceRequirements

def test_mpijob_create():
    """Test MPIJob creation generates correct manifest."""
    job = MPIJob(
        name="test-job",
        namespace="default",
        workers=4,
        image="pytorch:latest",
        command=["python", "train.py"],
        resources=ResourceRequirements(cpu=8, memory="32Gi", gpu=1),
    )

    manifest = job._build_manifest()

    assert manifest["metadata"]["name"] == "test-job"
    assert manifest["spec"]["mpiReplicaSpecs"]["Worker"]["replicas"] == 4
    assert manifest["spec"]["mpiReplicaSpecs"]["Worker"]["template"]["spec"]["containers"][0]["image"] == "pytorch:latest"
    assert manifest["spec"]["mpiReplicaSpecs"]["Worker"]["template"]["spec"]["containers"][0]["resources"]["limits"]["nvidia.com/gpu"] == "1"

def test_mpijob_wait_for_completion_timeout():
    """Test wait_for_completion raises TimeoutError."""
    job = MPIJob(name="test-job", namespace="default", workers=2, image="test")

    with patch.object(job, 'get_status', return_value='Running'):
        with pytest.raises(TimeoutError):
            job.wait_for_completion(timeout=1)
```

### 9.2 Integration Testing

**Test Environment:**
- OpenShift 4.14+ cluster
- 4 GPU nodes (8 GPUs total, e.g., NVIDIA T4)
- KubeFlow Trainer V2 installed
- Kueue installed

**Test Scenarios:**

```python
# tests/integration/test_mpijob_lifecycle.py
import pytest
from openshift_ai_sdk.training import MPIJob, ResourceRequirements

def test_mpijob_end_to_end(kube_client):
    """Test complete MPIJob lifecycle."""
    job = MPIJob(
        name="integration-test",
        namespace="test-namespace",
        workers=2,
        image="horovod/horovod:latest",
        command=["horovodrun", "-np", "2", "python", "-c", "import horovod.torch as hvd; hvd.init(); print(f'Rank {hvd.rank()}')"],
        resources=ResourceRequirements(cpu=4, memory="8Gi", gpu=1),
    )

    # Create job
    job.create()

    # Wait for completion (5 min timeout)
    status = job.wait_for_completion(timeout=300)
    assert status == "Succeeded"

    # Check logs contain MPI output
    logs = job.get_logs()
    assert "Rank 0" in logs

    # Cleanup
    job.delete()

def test_mpijob_quota_enforcement(kube_client):
    """Test ResourceQuota prevents over-allocation."""
    # Set quota: 2 GPUs max
    create_resource_quota("test-namespace", gpus=2)

    # Try to create job requiring 4 GPUs
    job = MPIJob(
        name="quota-test",
        namespace="test-namespace",
        workers=4,
        image="pytorch:latest",
        resources=ResourceRequirements(gpu=1),
    )

    with pytest.raises(Exception) as exc_info:
        job.create()

    assert "quota exceeded" in str(exc_info.value).lower()

def test_mpijob_multi_tenancy(kube_client):
    """Test namespace isolation."""
    # Create MPIJob in namespace-a
    job_a = MPIJob(name="job-a", namespace="namespace-a", workers=2, image="test")
    job_a.create()

    # User in namespace-b should not see job-a
    with pytest.raises(Exception) as exc_info:
        MPIJob.get("job-a", namespace="namespace-b")

    assert "not found" in str(exc_info.value).lower()
```

### 9.3 Performance Testing

**Benchmarks:**

1. **Scaling Efficiency:**
   - Train ResNet-50 on ImageNet
   - Measure throughput (images/sec) on 1, 2, 4, 8, 16 nodes
   - Target: >85% efficiency up to 8 nodes, >75% up to 16 nodes

2. **Job Submission Latency:**
   - Submit 100 MPIJobs concurrently
   - Measure time from CLI command to job CR created
   - Target: <5 seconds p99

3. **Dashboard UI Responsiveness:**
   - Load job list page with 100 jobs
   - Measure initial page load time
   - Target: <2 seconds p95

**Example: Scaling Efficiency Test**

```bash
#!/bin/bash
# tests/performance/scaling_efficiency.sh

for nodes in 1 2 4 8; do
  echo "Testing $nodes nodes..."

  oc odh mpijob create resnet-$nodes \
    --workers $nodes \
    --gpu 1 \
    --image quay.io/odh/resnet50-benchmark:latest \
    --command "python benchmark.py --batch-size 128"

  # Wait for completion
  oc odh mpijob wait resnet-$nodes --timeout 1h

  # Extract throughput from logs
  throughput=$(oc odh mpijob logs resnet-$nodes | grep "throughput" | awk '{print $2}')
  echo "$nodes nodes: $throughput images/sec"

  # Calculate efficiency
  baseline=1200  # Single-node throughput
  expected=$((baseline * nodes))
  actual=$throughput
  efficiency=$(echo "scale=2; $actual / $expected * 100" | bc)
  echo "Efficiency: $efficiency%"

  # Cleanup
  oc odh mpijob delete resnet-$nodes
done
```

### 9.4 Security Testing

**Test Cases:**

1. **RBAC Enforcement:**
   - User with `mpijob-viewer` role cannot create/delete jobs
   - User can only access jobs in authorized namespaces

2. **NetworkPolicy Isolation:**
   - Pod in job A cannot connect to pod in job B (different namespace)
   - Launcher can SSH to workers, but workers cannot SSH to launcher

3. **Secret Isolation:**
   - Pod in namespace A cannot mount Secret from namespace B

4. **Vulnerability Scanning:**
   - All container images pass Clair/Trivy scans (no critical CVEs)

---

## 10. Migration & Compatibility

### 10.1 From Standalone KubeFlow MPI Operator

**Migration Path:**

```bash
# 1. Export existing MPIJobs from standalone operator
kubectl get mpijobs -o yaml > mpijobs-backup.yaml

# 2. Install RHOAI with Training Operator v2
oc apply -f rhoai-operator.yaml

# 3. Uninstall standalone MPI Operator
kubectl delete deployment mpi-operator -n kubeflow

# 4. Re-create MPIJobs (CRD schema compatible)
oc apply -f mpijobs-backup.yaml
```

**API Compatibility:**
- v1alpha1 → v2beta1 automatic conversion via webhook
- No manual YAML changes required for most jobs

### 10.2 From AWS SageMaker

**Config Converter Tool:**

```python
# tools/migrate_sagemaker.py
import yaml
import json

def convert_sagemaker_to_mpijob(sagemaker_config: dict) -> dict:
    """Convert SageMaker MPI training job config to RHOAI MPIJob."""
    return {
        "apiVersion": "kubeflow.org/v2beta1",
        "kind": "MPIJob",
        "metadata": {
            "name": sagemaker_config["TrainingJobName"],
        },
        "spec": {
            "mpiReplicaSpecs": {
                "Launcher": {
                    "replicas": 1,
                    "template": {
                        "spec": {
                            "containers": [{
                                "name": "launcher",
                                "image": sagemaker_config["AlgorithmSpecification"]["TrainingImage"],
                                "command": sagemaker_config["AlgorithmSpecification"]["ContainerEntrypoint"],
                            }]
                        }
                    }
                },
                "Worker": {
                    "replicas": sagemaker_config["ResourceConfig"]["InstanceCount"],
                    "template": {
                        "spec": {
                            "containers": [{
                                "name": "worker",
                                "image": sagemaker_config["AlgorithmSpecification"]["TrainingImage"],
                                "resources": {
                                    "limits": {
                                        "nvidia.com/gpu": str(_parse_gpu_count(
                                            sagemaker_config["ResourceConfig"]["InstanceType"]
                                        )),
                                    }
                                }
                            }]
                        }
                    }
                }
            }
        }
    }

# Usage:
# python migrate_sagemaker.py --input sagemaker-config.json --output mpijob.yaml
```

---

## 11. Operational Runbooks

### 11.1 Troubleshooting Common Issues

**Issue: MPIJob Stuck in Pending**

**Symptoms:**
```bash
$ oc odh mpijob get my-job
NAME      STATUS    DURATION
my-job    Pending   5m
```

**Diagnosis:**
```bash
# Check worker pod status
oc get pods -l training.kubeflow.org/job-name=my-job

# Common causes:
# 1. Insufficient GPUs
oc describe node | grep nvidia.com/gpu

# 2. Image pull failure
oc describe pod my-job-worker-0 | grep -A5 Events

# 3. Quota exceeded
oc get resourcequota -n <namespace>
```

**Resolution:**
- Insufficient GPUs: Wait for resources or reduce worker count
- Image pull failure: Check registry credentials, image tag exists
- Quota exceeded: Increase quota or delete other jobs

**Issue: MPI Communication Failure**

**Symptoms:**
```
Launcher logs: "ssh: connect to host my-job-worker-0 port 2222: Connection refused"
```

**Diagnosis:**
```bash
# Check NetworkPolicy allows SSH
oc get networkpolicies -n <namespace>

# Check SSH service running in workers
oc exec my-job-worker-0 -- ps aux | grep sshd

# Check DNS resolution
oc exec my-job-launcher -- nslookup my-job-worker-0
```

**Resolution:**
- NetworkPolicy blocking: Apply MPIJob NetworkPolicy template
- SSH not running: Check worker container image includes SSH server
- DNS failure: Verify StatefulSet created workers (not Deployment)

### 11.2 Performance Tuning

**Low GPU Utilization (<70%)**

**Diagnosis:**
```bash
# Check NVIDIA DCGM metrics
oc exec my-job-worker-0 -- nvidia-smi dmon -s um

# Common causes:
# 1. Data loading bottleneck (CPU-bound)
# 2. Small batch size (GPU underutilized)
# 3. Slow network (MPI Allreduce latency)
```

**Resolution:**
- Data loading: Increase workers per GPU, use data prefetching
- Batch size: Increase global batch size, use gradient accumulation
- Network: Check NCCL logs for bandwidth, consider RDMA (post-MVP)

**Poor Scaling Efficiency**

**Diagnosis:**
```python
# Calculate scaling efficiency
baseline_throughput = 1200  # samples/sec on 1 node
actual_throughput = 8500    # samples/sec on 8 nodes
expected_throughput = baseline_throughput * 8  # 9600

efficiency = (actual_throughput / expected_throughput) * 100
print(f"Scaling efficiency: {efficiency}%")  # Should be >85%
```

**Resolution:**
- <70% efficiency: Check network bandwidth, reduce Allreduce frequency
- Imbalanced workers: Check data partitioning, worker CPU allocation
- Stragglers: Use distributed tracing to identify slow workers

---

## 12. Open Questions & Future Work

### 12.1 Open Questions for Engineering Team

1. **Training Operator HA**: What is the leader election strategy for Training Operator? How do we ensure zero downtime during operator upgrades?

2. **Kueue Integration Depth**: Should MPIJob status reflect Kueue queue position, or only Kubernetes pod status?

3. **Dashboard WebSocket vs. Polling**: Use WebSockets for real-time log streaming, or poll every 5s? (Trade-off: complexity vs. latency)

4. **SDK Async Patterns**: Should SDK support `asyncio` for non-blocking operations, or only synchronous? (Target: ML engineers using Jupyter notebooks)

5. **Checkpoint Storage**: Should we provide opinionated guidance on PVC vs. S3 for checkpoints, or leave to user choice?

6. **GPU Metrics Granularity**: Expose per-GPU metrics (8 GPUs per node) or per-pod aggregates? (Trade-off: cardinality vs. detail)

### 12.2 Future Enhancements

**Elastic Training (Post-MVP):**
- Dynamic worker scaling during training (add/remove workers based on loss curve plateau)
- Requires checkpoint/restore support
- Use case: Cost optimization by reducing workers after initial fast learning phase

**Multi-Cluster Federation (Post-MVP):**
- Run single MPIJob across multiple OpenShift clusters
- Use case: Burst to cloud when on-prem GPUs saturated
- Challenges: High inter-cluster network latency, complex scheduling

**AutoML Integration (Post-MVP):**
- Automatic MPIJob configuration tuning (optimal worker count, batch size)
- Use case: Novice users don't know how to configure distributed training
- Approach: Bayesian optimization over MPIJob hyperparameters

**Spot/Preemptible Instance Support (Post-MVP):**
- Use lower-cost spot instances for workers, checkpoint frequently
- Use case: Cost-sensitive training workloads
- Challenges: Checkpoint overhead, handling frequent restarts

---

## 13. Conclusion

This architecture provides a production-ready blueprint for integrating MPIJobs into RedHat OpenShift AI using KubeFlow Trainer V2. The design balances three critical tensions:

1. **Upstream Alignment vs. Enterprise Requirements**: Leverage upstream KubeFlow while adding RHOAI-specific security, multi-tenancy, and observability.

2. **Simplicity vs. Flexibility**: Provide simple interfaces (Dashboard, SDK) for common use cases while allowing power users to access full Kubernetes/MPI capabilities.

3. **Performance vs. Operational Complexity**: Deliver excellent distributed training performance with standard networking (MVP), with clear path to RDMA for advanced users (post-MVP).

**Key Success Factors:**

- **Unified Experience**: MPIJobs feel native to RHOAI, not bolted-on
- **Enterprise-Grade**: Security, multi-tenancy, compliance built-in from day one
- **Hybrid Cloud Ready**: Identical experience on-prem and cloud, avoiding vendor lock-in
- **Observability**: Full visibility into distributed training performance via Grafana/Prometheus
- **Incremental Adoption**: Users start with single-node training, scale to multi-node seamlessly

**Next Steps:**

1. Engineering team reviews architecture, challenges assumptions in ADRs
2. Design reviews with customers (private beta candidates) for UX validation
3. Spike: Training Operator v2 packaging in RHOAI operator (2 weeks)
4. Spike: Dashboard UI mockups in PatternFly (1 week)
5. Kickoff Phase 1 implementation (Q4 2025)

**Questions or Feedback:**
- Architecture discussions: #rhoai-mpijobs Slack channel
- Technical questions: File issues in rhoai-architecture repo
- Customer feedback: Share via product management team

---

**Document Ownership:**
- **Author**: RHOAI Architecture Team
- **Reviewers**: Engineering Leadership, Product Management, Security Team
- **Approval**: VP Engineering, VP Product
- **Last Updated**: 2025-10-29
