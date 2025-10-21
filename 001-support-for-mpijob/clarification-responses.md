# MPIJob Clarification Questions - Technical Recommendations

This document provides specific technical recommendations for the clarification questions identified in the MPIJob support specification.

## Technical Architecture Questions

### Q-ARCH-001: MPI Implementation Support

**Question**: Should MPIJob support all MPI implementations (OpenMPI, Intel MPI, MPICH) or focus on a primary implementation for MVP?

**Recommendation**: Support all major MPI implementations (OpenMPI, Intel MPI, MPICH) for MVP.

**Technical Details**:
- The MPIJob controller should be implementation-agnostic, as MPI variants are container-level concerns
- For each supported MPI implementation, provide:
  - Reference container images with properly configured MPI environments
  - Documentation for required container configurations
  - Example job specifications
- Testing matrix:
  - Primary testing focus on OpenMPI 4.1.x (most common in ML workloads)
  - Secondary testing with Intel MPI 2021.x (important for enterprise customers)
  - Validation testing with MPICH 3.4.x (compatibility check)

**Implementation Approach**:
```yaml
# Container requirements documentation for each MPI variant
OpenMPI:
  - Base image with OpenMPI 4.1.x installed
  - SSH daemon configured for inter-worker communication
  - Proper environment variables (OMPI_MCA_*)
  - User permissions for MPI processes

Intel MPI:
  - Base image with Intel MPI 2021.x installed
  - Process management daemon (PMI)
  - Intel MPI environment variables
  - License considerations

MPICH:
  - Base image with MPICH 3.4.x
  - Hydra process manager configuration
  - Environment variables
```

### Q-ARCH-002: Networking and IPC

**Question**: What network configuration is required for MPI communication between workers?

**Recommendation**: Provide standard network configuration patterns and NetworkPolicy templates, with RDMA support deferred to post-MVP.

**Technical Details**:
- Communication Requirements:
  - Worker-to-worker: TCP ports 1024-65535 (MPI dynamic port range)
  - Launcher-to-worker: SSH port 22 (for OpenMPI), PMI ports for other MPI variants
  - All workers must have DNS resolution of other workers

- NetworkPolicy Templates:
  - Default template: Allow all traffic within namespace between MPIJob pods
  - Restricted template: Allow only required ports between specific pods
  - Example policies for common security postures

- CNI Compatibility:
  - Primary testing: OpenShift SDN (OVS-based)
  - Validation testing: OVN-Kubernetes
  - Documentation for other CNI plugins (Calico, Cilium, etc.)

**Implementation Approach**:
```yaml
# Example NetworkPolicy for MPIJob
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mpijob-allow-internal
spec:
  podSelector:
    matchLabels:
      training.openshift.io/mpijob: ${JOB_NAME}
  ingress:
  - from:
    - podSelector:
        matchLabels:
          training.openshift.io/mpijob: ${JOB_NAME}
  policyTypes:
  - Ingress
```

### Q-ARCH-003: Resource Orchestration

**Question**: How does MPIJob handle pod scheduling for multi-worker topologies?

**Recommendation**: Use Volcano Scheduler for gang scheduling in production environments, with a fallback to default scheduler for simple deployments.

**Technical Details**:
- Gang Scheduling:
  - All worker pods must be scheduled together or not at all
  - Prevents deadlocks where partial job placement blocks other jobs
  - Critical for resource-intensive jobs (multiple GPUs)

- Volcano Integration:
  - MPIJob CRD extended with scheduling annotations
  - PodGroups created for MPIJob workers
  - Queue-based scheduling for fair resource allocation

- Eviction Handling:
  - When a worker pod is evicted: entire job fails (fail-fast)
  - Clear error messages indicating which worker was evicted and why
  - Future enhancement: checkpointing and recovery (post-MVP)

**Implementation Approach**:
```yaml
# Example MPIJob with Volcano scheduling
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: tensorflow-training
spec:
  runPolicy:
    schedulingPolicy:
      priorityClass: training-critical
      minAvailable: 8  # All 8 workers must be scheduled together
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      template:
        metadata:
          annotations:
            scheduling.k8s.io/group-name: tf-training-job
    Worker:
      replicas: 8
      template:
        metadata:
          annotations:
            scheduling.k8s.io/group-name: tf-training-job
```

### Q-ARCH-004: Hardware Acceleration

**Question**: What GPU types and accelerator configurations are supported?

**Recommendation**: Primarily support NVIDIA GPUs for MVP, with documentation for AMD GPUs and a roadmap for expanded accelerator support.

**Technical Details**:
- Primary Support (MVP):
  - NVIDIA GPUs with CUDA 11.x/12.x
  - Tested configurations: V100, A100, T4, L4
  - Multi-GPU configurations: 1-8 GPUs per worker

- Secondary Support (Documented):
  - AMD GPUs with ROCm 5.x
  - Documentation for container requirements
  - Example job specifications

- Mixed GPU Support:
  - Not supported in MVP
  - Document as limitation with workaround (separate job runs)

- Future Roadmap:
  - Intel accelerators (post-MVP)
  - Heterogeneous worker configurations (post-MVP)
  - Custom resource definition for specialized accelerators

**Implementation Approach**:
```yaml
# Example worker spec with GPU resources
Worker:
  replicas: 4
  template:
    spec:
      containers:
      - name: training
        resources:
          limits:
            nvidia.com/gpu: 2  # 2 NVIDIA GPUs per worker
          requests:
            nvidia.com/gpu: 2
            memory: "16Gi"
            cpu: "8"
```

### Q-ARCH-005: Container Image Requirements

**Question**: Must users install specific MPI versions? Are there SSH or communication daemon requirements?

**Recommendation**: Provide clear container requirements documentation and reference images for common scenarios.

**Technical Details**:
- Base Requirements:
  - MPI implementation installed and properly configured
  - SSH daemon for OpenMPI (standard approach)
  - Appropriate user permissions for MPI processes

- Reference Images:
  - Base images for common ML frameworks + MPI combinations
  - TensorFlow + Horovod + OpenMPI
  - PyTorch + OpenMPI
  - Intel MPI + common frameworks

- Security Considerations:
  - SSH keys generated and distributed automatically by launcher
  - No plaintext credentials in container images
  - Permission model for non-root execution

**Implementation Approach**:
```Dockerfile
# Example Dockerfile for TensorFlow + Horovod + OpenMPI
FROM tensorflow/tensorflow:2.9.0-gpu

# Install OpenMPI dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    openssh-server openssh-client \
    libopenmpi-dev \
    && rm -rf /var/lib/apt/lists/*

# Install Horovod with TensorFlow and MPI support
RUN pip install --no-cache-dir horovod[tensorflow,keras]

# Configure SSH for OpenMPI
RUN mkdir -p /var/run/sshd && \
    sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PubkeyAuthentication yes/PubkeyAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#RSAAuthentication yes/RSAAuthentication yes/' /etc/ssh/sshd_config && \
    sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config

# Create SSH keys directory and setup entrypoint
RUN mkdir -p /root/.ssh
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
```

## Non-Functional Requirements

### Q-NFR-001: Performance Baseline

**Question**: Expected speedup characteristics? Acceptable overhead?

**Recommendation**: Define clear performance targets with benchmarks, focusing on both speedup and overhead metrics.

**Technical Details**:
- Speedup Targets:
  - Linear scaling up to 8 workers (8x speedup with 8 workers) for compute-bound workloads
  - Minimum 2x speedup with 4 workers for mixed workloads
  - Reference benchmarks with ResNet50, BERT, and custom distributed training examples

- Overhead Metrics:
  - Job submission latency: < 5 seconds (p95)
  - Worker initialization: < 30 seconds (p95)
  - Status propagation: < 10 seconds for state changes

- Network Requirements:
  - Minimum 10 Gbps network between worker nodes for acceptable performance
  - 25+ Gbps recommended for NLP and large-scale training
  - Document performance degradation with slower networks

**Implementation Approach**:
```python
# Example performance testing script (pseudocode)
def benchmark_scaling(worker_counts=[1, 2, 4, 8, 16]):
    """Test scaling efficiency with different worker counts."""
    results = {}
    for workers in worker_counts:
        job_spec = create_benchmark_job(
            model="resnet50",
            dataset="imagenet",
            batch_size=64*workers,  # Scale batch size with workers
            workers=workers,
            gpus_per_worker=1
        )
        start_time = time.time()
        job = create_mpijob(job_spec)
        wait_for_completion(job)
        elapsed = time.time() - start_time
        results[workers] = {
            'time': elapsed,
            'throughput': calculate_throughput(job),
            'scaling_efficiency': calculate_efficiency(job, baseline=results.get(1))
        }
    return results
```

### Q-NFR-002: Scalability Limits

**Question**: Maximum number of workers per job? Maximum concurrent MPIJobs per cluster?

**Recommendation**: Define and test specific scale targets for both single job scale and multi-job concurrency.

**Technical Details**:
- Per-Job Scale:
  - MVP target: 100 workers per job tested and supported
  - Extended testing with standard benchmark jobs at 16, 32, 64, and 100 workers
  - Document performance characteristics at each scale point
  - Document resource requirements for launcher at different worker counts

- Cluster Concurrency:
  - MVP target: 50 concurrent jobs per cluster tested and supported
  - Test scenarios: many small jobs vs. few large jobs
  - Resource contention patterns and mitigation strategies
  - Recommendations for namespace quotas at different cluster sizes

**Implementation Approach**:
```bash
#!/bin/bash
# Scale testing script
# Test single large job
echo "Testing maximum worker scale..."
for WORKERS in 16 32 64 100; do
  echo "Creating job with $WORKERS workers..."
  odh training create mpijob scale-test-$WORKERS --workers $WORKERS --image benchmark:latest
  # Monitor and collect metrics
done

# Test concurrent jobs
echo "Testing maximum concurrency..."
for i in $(seq 1 50); do
  echo "Creating job $i of 50..."
  odh training create mpijob concurrency-test-$i --workers 4 --image benchmark:latest
done
# Monitor cluster health and job success rates
```

### Q-NFR-003: Reliability and Fault Tolerance

**Question**: What happens if one worker fails? Support for checkpoint/restart?

**Recommendation**: Implement fail-fast behavior for MVP with clear error reporting, with checkpoint/restart capabilities planned for post-MVP.

**Technical Details**:
- Worker Failure Handling:
  - For MVP: Entire job fails if any worker fails (fail-fast)
  - Clear error reporting showing which worker failed and why
  - Automatic cleanup of all resources on failure

- Checkpoint/Restart (Post-MVP):
  - Framework-agnostic checkpoint storage specification
  - Integration with S3/object storage for model state
  - Automatic restart from latest checkpoint

- Node Failure Handling:
  - Integration with OpenShift node health monitoring
  - Quick detection and reporting of node failures affecting jobs
  - Appropriate status updates and event generation

**Implementation Approach**:
```yaml
# Example error reporting structure
status:
  state: Failed
  reason: "WorkerFailed"
  message: "Worker-3 failed with exit code 137 (OOM Killed). Consider increasing memory allocation."
  workerStatuses:
    - id: "Worker-0"
      state: "Succeeded"
      restarts: 0
    - id: "Worker-1"
      state: "Succeeded"
      restarts: 0
    - id: "Worker-2"
      state: "Succeeded"
      restarts: 0
    - id: "Worker-3"
      state: "Failed"
      reason: "OOMKilled"
      message: "Container exceeded memory limit (16Gi)"
      restarts: 2
```

## Monitoring and Logging

### Q-OBS-001: Metrics Collection

**Question**: Are training metrics (loss, accuracy) collected, or only infrastructure metrics?

**Recommendation**: Focus on infrastructure metrics for MVP, with a well-defined approach for training metrics in post-MVP.

**Technical Details**:
- MVP Metrics (Infrastructure):
  - Job lifecycle: counts by state, duration histograms
  - Resource utilization: CPU, memory, GPU utilization per worker
  - Initialization time, completion time, error counts
  - Worker health metrics: restarts, state transitions

- Post-MVP (Training Metrics):
  - Framework-agnostic metrics collection API
  - Standard metrics: loss, accuracy, learning rate, etc.
  - Custom metric registration
  - Visualization in Dashboard with historical comparison

**Implementation Approach**:
```yaml
# Example Prometheus metrics
# Infrastructure metrics (MVP)
mpijob_count{namespace, state}  # Count of jobs by state
mpijob_duration_seconds{namespace, job_name}  # Job duration
mpijob_worker_restarts{namespace, job_name, worker_id}  # Worker restart count
mpijob_gpu_utilization{namespace, job_name, worker_id, gpu_id}  # GPU utilization

# Training metrics (Post-MVP)
mpijob_training_loss{namespace, job_name, step}  # Training loss
mpijob_training_accuracy{namespace, job_name, step}  # Training accuracy
mpijob_learning_rate{namespace, job_name, step}  # Learning rate
```

### Q-OBS-002: Distributed Log Aggregation

**Question**: Real-time log streaming vs. batch retrieval? Log correlation?

**Recommendation**: Implement real-time log streaming with correlation features and configurable retention policies.

**Technical Details**:
- Log Collection:
  - Real-time streaming API for immediate visibility
  - Batch retrieval for historical analysis
  - Log aggregation across all workers and launcher
  - Worker-specific filtering options

- Log Correlation:
  - Timestamp normalization across workers
  - Correlation IDs for tracking related events
  - Metadata enrichment (worker ID, job name, namespace)

- Retention Policies:
  - Default: 7 days retention for completed jobs
  - Configurable via annotation: `training.openshift.io/log-retention: "30d"`
  - Integration with OpenShift logging infrastructure

**Implementation Approach**:
```go
// Example log aggregation function
func (c *Controller) GetJobLogs(ctx context.Context, req *pb.GetJobLogsRequest) (*pb.LogStream, error) {
    // Find all pods for the specified job
    pods, err := c.findJobPods(req.Namespace, req.JobName)
    if err != nil {
        return nil, err
    }

    // Create log streams for each pod
    streams := make(map[string]io.ReadCloser)
    for _, pod := range pods {
        podLogs, err := c.kubeClient.CoreV1().Pods(req.Namespace).GetLogs(pod.Name, &corev1.PodLogOptions{
            Container: "training",
            Follow:    req.Follow,
            TailLines: req.TailLines,
            Previous:  req.Previous,
        }).Stream(ctx)
        if err != nil {
            continue  // Skip failed streams but continue with others
        }
        streams[pod.Labels["training.openshift.io/replica-id"]] = podLogs
    }

    // Multiplex streams with worker ID prefixing
    return newAggregatedLogStream(streams), nil
}
```

### Q-OBS-003: Debugging Tools

**Question**: Can users exec into running worker pods? Support for attaching debuggers?

**Recommendation**: Provide comprehensive debugging capabilities with appropriate permissions and documentation.

**Technical Details**:
- Interactive Debugging:
  - Allow `kubectl exec` access to worker pods (requires appropriate permissions)
  - Support attaching to specific processes within workers
  - Document debugging workflows for common scenarios

- Debugger Support:
  - Documentation for configuring common ML debuggers (TensorFlow Debugger, PyTorch Profiler)
  - Example job configurations with debug mode enabled
  - Remote debugger setup guides (VSCode, PyCharm)

- Post-Mortem Analysis:
  - Failure analysis tools in Dashboard
  - Log search and pattern recognition
  - Resource utilization timeline correlated with job events

**Implementation Approach**:
```yaml
# Example debug-enabled job configuration
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: debug-training
  annotations:
    training.openshift.io/debug-mode: "true"  # Enables debug features
spec:
  runPolicy:
    cleanPodPolicy: None  # Keep pods after completion for inspection
  mpiReplicaSpecs:
    Worker:
      replicas: 4
      template:
        spec:
          containers:
          - name: training
            image: training:debug
            ports:
            - containerPort: 8080
              name: debugger
            env:
            - name: TF_DEBUGGER_PORT
              value: "8080"
```

## Security Questions

### Q-SEC-001: RBAC Model

**Question**: Is MPIJob creation a privileged operation vs. other job types?

**Recommendation**: Implement a consistent RBAC model aligned with existing training job types, with fine-grained permissions for different operations.

**Technical Details**:
- Permission Model:
  - `training.openshift.io/mpijobs: create,update,patch,delete` - Job management
  - `training.openshift.io/mpijobs/status: get` - Status viewing
  - `training.openshift.io/mpijobs/log: get` - Log access

- Role Hierarchy:
  - Data Scientist: Create, manage own jobs, view status/logs
  - MLOps Engineer: Manage all jobs in namespace
  - Administrator: Full control across namespaces

- Creation Constraints:
  - Same privilege level as other job types
  - Additional validation for resource quotas
  - Optional: Additional validation for GPU requests

**Implementation Approach**:
```yaml
# Example RBAC configuration for data scientist
kind: Role
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: data-scientist
  namespace: team-ml
rules:
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs", "tfjobs", "pytorchjobs"]
  verbs: ["create", "get", "list", "watch", "update", "patch", "delete"]
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs/status", "tfjobs/status", "pytorchjobs/status"]
  verbs: ["get", "list", "watch"]
- apiGroups: ["kubeflow.org"]
  resources: ["mpijobs/log", "tfjobs/log", "pytorchjobs/log"]
  verbs: ["get"]
```

### Q-SEC-002: Network Security

**Question**: What network policies are required? Specific port ranges?

**Recommendation**: Provide detailed network security documentation with example policies for different security postures.

**Technical Details**:
- Required Connectivity:
  - Launcher to Workers: SSH (port 22) for OpenMPI
  - Worker to Worker: Dynamic ports (typical range: 1024-65535)
  - All pods: DNS resolution for worker hostnames

- NetworkPolicy Templates:
  - Default: Allow all traffic between pods with same job label
  - Restricted: Allow only specific ports between launcher and workers
  - Isolated: Strict policies for high-security environments

- Service Mesh Integration:
  - Optional mTLS for worker communication (post-MVP)
  - Traffic encryption recommendations
  - Integration with OpenShift Service Mesh

**Implementation Approach**:
```yaml
# Example restricted NetworkPolicy for MPIJob
apiVersion: networking.k8s.io/v1
kind: NetworkPolicy
metadata:
  name: mpijob-restricted
spec:
  podSelector:
    matchLabels:
      training.openshift.io/mpijob: ${JOB_NAME}
  ingress:
  - from:
    - podSelector:
        matchLabels:
          training.openshift.io/mpijob-role: launcher
          training.openshift.io/mpijob: ${JOB_NAME}
    ports:
    - port: 22
      protocol: TCP
  - from:
    - podSelector:
        matchLabels:
          training.openshift.io/mpijob-role: worker
          training.openshift.io/mpijob: ${JOB_NAME}
    ports:
    - port: 1024
      protocol: TCP
    - port: 65535
      protocol: TCP
  policyTypes:
  - Ingress
```

### Q-SEC-003: Secret Management

**Question**: How are credentials and sensitive configuration managed?

**Recommendation**: Leverage OpenShift secret management with appropriate mounting patterns and access controls.

**Technical Details**:
- Secret Types:
  - Registry credentials for private images
  - Dataset access credentials
  - License files for proprietary MPI implementations
  - SSH keys for worker communication

- Management Patterns:
  - Launcher generates and distributes communication keys
  - User-provided secrets mounted via Kubernetes secrets
  - Integration with external secret managers via CSI drivers (HashiCorp Vault, AWS Secrets Manager)

- Access Controls:
  - Least privilege access to secrets
  - Secret scoping to specific namespaces
  - Audit logging for secret access

**Implementation Approach**:
```yaml
# Example job with secrets
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: secure-training
spec:
  mpiReplicaSpecs:
    Worker:
      replicas: 4
      template:
        spec:
          containers:
          - name: training
            image: private-registry.example.com/training:latest
            volumeMounts:
            - name: dataset-credentials
              mountPath: /etc/credentials
              readOnly: true
          volumes:
          - name: dataset-credentials
            secret:
              secretName: s3-dataset-access
          imagePullSecrets:
          - name: private-registry-credentials
```

### Q-SEC-004: Maximum Job Duration

**Question**: Should there be a maximum job duration policy?

**Recommendation**: Implement configurable timeouts with defaults based on user roles and resource requests.

**Technical Details**:
- Default Timeouts:
  - Standard jobs: 24 hours default timeout
  - Extended jobs: 7 days with annotation
  - Administrative override for longer runs

- Configuration Options:
  - Job-level setting: `training.openshift.io/timeout: "72h"`
  - Namespace default: Configurable by administrator
  - Override permissions: Administrative role required

- Timeout Behavior:
  - Warning events at 80% of timeout
  - Graceful termination signal at timeout
  - Forced termination after grace period
  - Detailed reason code in status

**Implementation Approach**:
```yaml
# Example MPIJob with timeout configuration
apiVersion: kubeflow.org/v2
kind: MPIJob
metadata:
  name: long-running-job
  annotations:
    training.openshift.io/timeout: "72h"  # 72 hour timeout
    training.openshift.io/timeout-action: "terminate"  # or "notify-only"
spec:
  runPolicy:
    activeDeadlineSeconds: 259200  # 72 hours in seconds
  mpiReplicaSpecs:
    Launcher:
      replicas: 1
      # ...
```

## User Experience Questions

### Q-UX-001: Abstraction Level

**Question**: How much MPI complexity do we expose to data scientists?

**Recommendation**: Implement progressive disclosure pattern with sensible defaults for most users, and advanced options for power users.

**Technical Details**:
- Abstraction Levels:
  - Level 1 (Default): Simple job submission with minimal MPI knowledge required
    - Users specify: image, command, worker count, resources
    - System handles: MPI configuration, launcher setup

  - Level 2 (Advanced): Configuration options for experienced users
    - Fine-tuning of MPI parameters
    - Custom launcher resources
    - Network tuning options

  - Level 3 (Expert): Full control for power users
    - Custom entrypoints
    - Advanced MPI flags
    - Specialized container configurations

**Implementation Approach**:
```typescript
// Example UI code for progressive disclosure
const MPIJobForm: React.FC = () => {
  const [showAdvanced, setShowAdvanced] = useState(false);
  const [showExpert, setShowExpert] = useState(false);

  return (
    <Form>
      {/* Basic fields always visible */}
      <FormGroup label="Job Name" required>
        <TextInput name="jobName" />
      </FormGroup>
      <FormGroup label="Training Image" required>
        <TextInput name="image" />
      </FormGroup>
      <FormGroup label="Worker Count" required>
        <NumberInput name="workers" min={1} max={100} />
      </FormGroup>
      <FormGroup label="GPUs per Worker" required>
        <NumberInput name="gpus" min={0} max={8} />
      </FormGroup>

      {/* Advanced toggle */}
      <Button variant="link" onClick={() => setShowAdvanced(!showAdvanced)}>
        {showAdvanced ? "Hide Advanced Options" : "Show Advanced Options"}
      </Button>

      {showAdvanced && (
        <>
          <FormGroup label="MPI Implementation">
            <Select name="mpiImplementation" options={['OpenMPI', 'Intel MPI', 'MPICH']} />
          </FormGroup>
          <FormGroup label="Launcher Resources">
            <ResourceSelector name="launcherResources" />
          </FormGroup>

          {/* Expert toggle */}
          <Button variant="link" onClick={() => setShowExpert(!showExpert)}>
            {showExpert ? "Hide Expert Options" : "Show Expert Options"}
          </Button>

          {showExpert && (
            <>
              <FormGroup label="MPI Flags">
                <TextInput name="mpiFlags" />
              </FormGroup>
              <FormGroup label="Custom Entrypoint">
                <TextInput name="entrypoint" />
              </FormGroup>
            </>
          )}
        </>
      )}

      <ActionGroup>
        <Button type="submit">Create Job</Button>
        <Button variant="secondary">Cancel</Button>
      </ActionGroup>
    </Form>
  );
};
```

### Q-UX-002: Cross-Channel Symmetry

**Question**: Should CLI/SDK support all UI features, or can some be UI-exclusive?

**Recommendation**: Maintain feature parity across all channels (CLI, SDK, UI) with consistent models and terminology.

**Technical Details**:
- Core Capabilities (All Channels):
  - Job CRUD operations
  - Status monitoring
  - Log access
  - Resource configuration

- Channel-Optimized Features:
  - UI: Visual topology, resource calculators
  - CLI: Batch operations, scriptable workflows
  - SDK: Programmatic integration, event hooks

- Unified Models:
  - Shared API models across all channels
  - Consistent terminology and concepts
  - Same validation rules everywhere

**Implementation Approach**:
```typescript
// Shared TypeScript model used by UI, CLI (via OpenAPI), and SDK
export interface MPIJobSpec {
  // Core fields (available in all channels)
  name: string;
  namespace?: string;
  image: string;
  command: string[];
  workerCount: number;
  gpusPerWorker: number;
  cpuPerWorker: string;
  memoryPerWorker: string;

  // Advanced fields (available in all channels but possibly hidden in UI)
  mpiImplementation?: 'OpenMPI' | 'IntelMPI' | 'MPICH';
  launcherResources?: ResourceRequirements;
  networkPolicy?: 'default' | 'restricted' | 'isolated';

  // Expert fields (available in all channels, deep in UI)
  mpiFlags?: string[];
  entrypoint?: string;
  slotsPerWorker?: number;
}
```

### Q-UX-003: Training Metrics Integration

**Question**: How deeply should MPIJob integrate with training metrics?

**Recommendation**: Implement a two-phase approach with infrastructure metrics in MVP and opt-in training metrics in post-MVP.

**Technical Details**:
- MVP Approach:
  - Infrastructure metrics only (resource utilization, job status)
  - No changes required to user training code
  - Documentation for manual metrics collection

- Post-MVP Approach:
  - Optional metrics collection via standardized API
  - Support for common ML frameworks (TensorFlow, PyTorch)
  - Minimal code changes for users to enable metrics

- Integration Architecture:
  - Sidecar pattern for metrics collection
  - Prometheus endpoint for scraping
  - Dashboard visualization components

**Implementation Approach**:
```python
# Example user training script with metrics integration (post-MVP)
import tensorflow as tf
from openshift_ai.metrics import MetricsLogger

# Initialize metrics logger
metrics = MetricsLogger.for_current_job()

# Model training loop
model = create_model()
for epoch in range(epochs):
    for batch in dataset:
        loss = train_step(model, batch)
        accuracy = calculate_accuracy(model, batch)

        # Report metrics to OpenShift AI
        metrics.log_metrics({
            'loss': float(loss),
            'accuracy': float(accuracy),
            'learning_rate': get_learning_rate(),
            'epoch': epoch,
            'step': step
        })
```

## Business Questions

### Q-BIZ-001: Pricing and Packaging

**Question**: Does MPIJob require premium SKU or is it included in base OpenShift AI?

**Technical Recommendation**: Include basic MPIJob support in the base OpenShift AI offering, with premium features for enterprise needs.

**Technical Details**:
- Base Features (All Tiers):
  - MPIJob creation and management
  - Basic observability (status, logs)
  - Standard resource allocation
  - Up to 8 workers per job

- Premium Features:
  - Advanced observability (custom metrics, dashboards)
  - Scale beyond 8 workers per job
  - Enhanced security features
  - Priority scheduling

**Implementation Approach**:
```yaml
# Feature toggles based on subscription tier
apiVersion: config.openshift.io/v1
kind: FeatureGate
metadata:
  name: MPIJobFeatures
spec:
  featureSet: CustomNoUpgrade
  customNoUpgrade:
    enabled:
    # Base features - always enabled
    - BasicMPIJob
    - StandardWorkerCount
    - BasicObservability
    # Premium features - enabled based on subscription
    - AdvancedMPIMetrics
    - ExtendedWorkerCount
    - EnhancedSecurity
    - PriorityScheduling
```

### Q-BIZ-002: Framework Certification

**Question**: Which ML frameworks and versions do we officially support and test?

**Recommendation**: Define a tiered framework certification approach with comprehensive testing for primary frameworks.

**Technical Details**:
- Tier 1 (Fully Certified):
  - TensorFlow 2.x with Horovod
  - PyTorch 1.x/2.x with Horovod
  - Comprehensive testing across worker counts and GPU configurations

- Tier 2 (Validated):
  - MXNet with Horovod
  - JAX
  - Validation testing with standard configurations

- Tier 3 (Compatible):
  - Any framework with MPI support
  - Documentation for configuration
  - Community support model

**Implementation Approach**:
```yaml
# Example certification test matrix
frameworks:
  - name: TensorFlow
    versions: ["2.10", "2.11", "2.12"]
    mpi_implementations: ["OpenMPI", "Intel MPI"]
    worker_counts: [1, 2, 4, 8, 16]
    gpus_per_worker: [1, 2, 4, 8]
    tests:
      - benchmark: "ResNet50"
      - benchmark: "BERT-Large"
      - benchmark: "Custom CNN"

  - name: PyTorch
    versions: ["1.13", "2.0"]
    mpi_implementations: ["OpenMPI", "Intel MPI"]
    worker_counts: [1, 2, 4, 8, 16]
    gpus_per_worker: [1, 2, 4, 8]
    tests:
      - benchmark: "ResNet50"
      - benchmark: "BERT-Large"
      - benchmark: "Custom RNN"
```

### Q-BIZ-003: Migration Support

**Question**: What level of migration support do we provide for customers moving from other platforms?

**Recommendation**: Provide tiered migration support with documented patterns, tools, and services offerings.

**Technical Details**:
- Documentation:
  - Migration guides for common platforms (AWS SageMaker, Azure ML, Kubeflow)
  - Framework-specific conversion examples
  - Architecture patterns for equivalent functionality

- Migration Tools:
  - Configuration converters (SageMaker → MPIJob, etc.)
  - Validation tools for compatibility checking
  - Script generators for common patterns

- Professional Services:
  - Migration assessment service
  - Hands-on migration assistance
  - Custom integration development

**Implementation Approach**:
```python
# Example migration tool - SageMaker to MPIJob converter
def convert_sagemaker_to_mpijob(sagemaker_config_file):
    """Convert SageMaker distributed training config to MPIJob."""
    with open(sagemaker_config_file) as f:
        sagemaker_config = json.load(f)

    # Extract key parameters
    image = sagemaker_config['AlgorithmSpecification']['TrainingImage']
    instance_count = sagemaker_config['ResourceConfig']['InstanceCount']
    instance_type = sagemaker_config['ResourceConfig']['InstanceType']

    # Map SageMaker instance type to OpenShift resources
    resources = map_instance_type_to_resources(instance_type)

    # Generate MPIJob config
    mpijob_config = {
        "apiVersion": "kubeflow.org/v2",
        "kind": "MPIJob",
        "metadata": {
            "name": f"migrated-{sagemaker_config['TrainingJobName']}",
            "annotations": {
                "training.openshift.io/migrated-from": "sagemaker"
            }
        },
        "spec": {
            "mpiReplicaSpecs": {
                "Launcher": {
                    "replicas": 1,
                    # Launcher spec
                },
                "Worker": {
                    "replicas": instance_count,
                    "template": {
                        "spec": {
                            "containers": [{
                                "image": image,
                                "resources": resources
                            }]
                        }
                    }
                }
            }
        }
    }

    return mpijob_config
```

## Conclusion

These technical recommendations provide detailed guidance for addressing the clarification questions in the MPIJob support specification. The recommendations balance implementation complexity with user needs, focusing on a robust MVP with clear paths for post-MVP enhancements.

Key decisions include:
1. Supporting all major MPI implementations rather than selecting a single one
2. Using Volcano for gang scheduling in production environments
3. NVIDIA GPU focus for MVP with documented AMD support
4. Clear container requirements with reference implementations
5. Infrastructure metrics in MVP, training metrics in post-MVP
6. Progressive disclosure pattern for MPI complexity in the UI
7. Feature parity across CLI, SDK, and UI channels

These recommendations align with the overall architecture and development phases outlined in the implementation plan, providing a comprehensive technical foundation for the MPIJob implementation.