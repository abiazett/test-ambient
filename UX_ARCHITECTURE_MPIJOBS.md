# UX Architecture Analysis: MPIJobs Support in RedHat OpenShift AI

**Feature**: MPIJobs Integration with KubeFlow Trainer V2
**Target Platform**: RedHat OpenShift AI
**Date**: 2025-10-29
**Document Owner**: UX Architecture Team

---

## Executive Summary

This document provides a comprehensive UX architecture analysis for integrating MPIJobs support into RedHat OpenShift AI using KubeFlow Trainer V2. The analysis focuses on creating a unified, consistent experience across CLI, SDK, and Dashboard UI touchpoints while maintaining alignment with existing KubeFlow Trainer job patterns. The key design principle is to abstract infrastructure complexity while providing the observability and control enterprise users require.

---

## 1. User Journey Mapping

### 1.1 Data Scientist Journey

**Primary Goal**: Run distributed training jobs efficiently without deep Kubernetes knowledge

#### Journey Stages:

**Stage 1: Discovery & Planning** (Pre-Job)
- **Touchpoint**: Documentation, Dashboard UI
- **Actions**:
  - Understand when to use MPIJobs vs other distributed training methods
  - Review available resources (GPUs, nodes, storage)
  - Determine resource requirements for their model
- **Pain Points**:
  - Uncertainty about MPI vs other distributed training approaches
  - Unclear resource availability and quotas
- **UX Needs**:
  - Contextual guidance on job type selection
  - Clear visibility of available cluster resources
  - Example templates and common configurations

**Stage 2: Job Configuration** (Creation)
- **Touchpoint**: Jupyter Notebook (SDK), Dashboard UI
- **Actions**:
  - Configure worker count and resource allocation (GPUs, memory, CPU)
  - Specify training script and container image
  - Set up persistent storage for model checkpoints
  - Define environment variables and dependencies
- **Pain Points**:
  - Complex YAML configurations
  - Uncertainty about optimal worker/resource ratios
  - Storage configuration complexity
- **UX Needs**:
  - Intelligent defaults based on workload patterns
  - Visual configuration wizard in Dashboard
  - SDK helper functions with validation
  - Clear error messages during configuration validation

**Stage 3: Job Submission** (Execution)
- **Touchpoint**: SDK, CLI, Dashboard UI
- **Actions**:
  - Submit job to cluster
  - Verify job acceptance and scheduling
- **Pain Points**:
  - Opaque submission failures
  - Lack of immediate feedback on resource availability
- **UX Needs**:
  - Immediate validation feedback
  - Clear queuing status
  - Estimated start time when resources unavailable

**Stage 4: Job Monitoring** (Active)
- **Touchpoint**: Dashboard UI, TensorBoard, CLI
- **Actions**:
  - Monitor training progress and metrics
  - Check resource utilization
  - Review logs from launcher and workers
  - Track model checkpoint progress
- **Pain Points**:
  - Fragmented monitoring across tools
  - Difficulty distinguishing launcher vs worker issues
  - Missing context on why jobs are slow or stuck
- **UX Needs**:
  - Unified dashboard view of all job components
  - Real-time metrics visualization
  - Topology view showing launcher-worker relationships
  - Integrated log aggregation with filtering
  - Resource utilization metrics per worker

**Stage 5: Troubleshooting** (When Issues Occur)
- **Touchpoint**: Dashboard UI, CLI logs
- **Actions**:
  - Diagnose failures or performance issues
  - Review error messages and logs
  - Adjust configuration and retry
- **Pain Points**:
  - Cryptic Kubernetes errors
  - Difficulty identifying root cause (network, resource, code)
  - Lost context when job fails
- **UX Needs**:
  - User-friendly error translations
  - Guided troubleshooting workflows
  - Failure pattern recognition with suggestions
  - Historical context preservation

**Stage 6: Completion & Iteration** (Post-Job)
- **Touchpoint**: Dashboard UI, Notebook
- **Actions**:
  - Review final metrics and results
  - Access trained models and artifacts
  - Compare with previous runs
  - Clone/modify configuration for next iteration
- **Pain Points**:
  - Manual tracking of experiment variations
  - Difficulty comparing distributed runs
- **UX Needs**:
  - Job history with searchable metadata
  - Side-by-side comparison views
  - Easy job cloning with configuration tweaks
  - Integration with model registry

---

### 1.2 MLOps Engineer Journey

**Primary Goal**: Operationalize, optimize, and scale distributed training workflows

#### Journey Stages:

**Stage 1: Infrastructure Planning**
- **Touchpoint**: Dashboard Admin UI, CLI
- **Actions**:
  - Assess cluster capacity for MPIJobs
  - Plan resource quotas and policies
  - Configure networking for MPI communication
- **Pain Points**:
  - Unclear MPI-specific infrastructure requirements
  - Network configuration complexity
- **UX Needs**:
  - Infrastructure readiness checklist
  - Network topology validation tools
  - Capacity planning calculator

**Stage 2: Template & Pipeline Creation**
- **Touchpoint**: SDK, CLI, Dashboard UI
- **Actions**:
  - Create reusable job templates
  - Build CI/CD pipelines for training
  - Define parameter ranges for optimization
- **Pain Points**:
  - Template versioning and sharing
  - Pipeline integration complexity
- **UX Needs**:
  - Template library with versioning
  - SDK integration patterns for common CI/CD tools
  - Parameter validation in templates

**Stage 3: Job Orchestration**
- **Touchpoint**: CLI, Dashboard UI, SDK
- **Actions**:
  - Schedule jobs with priority queues
  - Manage job dependencies
  - Implement resource quotas
- **Pain Points**:
  - Limited visibility into queue dynamics
  - Difficulty balancing fairness and priority
- **UX Needs**:
  - Queue management dashboard
  - Priority visualization and editing
  - Fair-share policy configuration UI

**Stage 4: Performance Optimization**
- **Touchpoint**: Dashboard Metrics, Monitoring Tools
- **Actions**:
  - Analyze job performance patterns
  - Identify bottlenecks (compute, network, I/O)
  - Optimize resource allocation
- **Pain Points**:
  - Fragmented metrics across systems
  - Difficulty attributing costs to jobs
- **UX Needs**:
  - Comprehensive performance analytics dashboard
  - MPI-specific metrics (communication overhead, load balancing)
  - Cost attribution and optimization recommendations
  - Comparative analysis across job configurations

**Stage 5: Governance & Compliance**
- **Touchpoint**: Dashboard Admin UI, Audit Logs
- **Actions**:
  - Monitor resource usage and costs
  - Enforce security policies
  - Track audit trails
- **Pain Points**:
  - Incomplete audit visibility
  - Manual policy enforcement
- **UX Needs**:
  - Policy-as-code interface
  - Automated compliance checking
  - Detailed audit reporting with filtering

---

### 1.3 OpenShift AI Administrator Journey

**Primary Goal**: Maintain platform stability, security, and optimal resource utilization

#### Journey Stages:

**Stage 1: Platform Setup**
- **Touchpoint**: OpenShift Console, CLI
- **Actions**:
  - Install/upgrade KubeFlow Trainer with MPIJob support
  - Configure cluster-wide settings
  - Set up monitoring and alerting
- **Pain Points**:
  - Complex installation dependencies
  - Version compatibility concerns
- **UX Needs**:
  - Automated compatibility checking
  - Installation wizard with validation steps
  - Clear upgrade path documentation

**Stage 2: User Enablement**
- **Touchpoint**: Dashboard Admin UI, Documentation
- **Actions**:
  - Create projects and assign permissions
  - Configure resource quotas per team
  - Provide training and documentation
- **Pain Points**:
  - Granular permission management
  - User onboarding overhead
- **UX Needs**:
  - Role-based access templates
  - Self-service user guides
  - Interactive onboarding tours

**Stage 3: Operations & Monitoring**
- **Touchpoint**: Dashboard Admin UI, Monitoring Systems
- **Actions**:
  - Monitor cluster health and utilization
  - Track job success/failure rates
  - Identify resource contention
- **Pain Points**:
  - Alert fatigue from too many signals
  - Difficulty prioritizing issues
- **UX Needs**:
  - Centralized operations dashboard
  - Intelligent alerting with severity levels
  - Trend analysis and predictive insights
  - MPIJob-specific health indicators

**Stage 4: Troubleshooting & Support**
- **Touchpoint**: Dashboard Admin UI, Logs, CLI
- **Actions**:
  - Investigate user-reported issues
  - Debug infrastructure problems
  - Optimize cluster configurations
- **Pain Points**:
  - Time-consuming root cause analysis
  - Limited visibility into past states
- **UX Needs**:
  - Historical state replay capability
  - Cross-component log correlation
  - Diagnostic runbook integration

**Stage 5: Capacity Planning**
- **Touchpoint**: Dashboard Analytics, Reporting
- **Actions**:
  - Analyze usage trends
  - Plan for growth and scaling
  - Optimize resource allocation
- **Pain Points**:
  - Manual data aggregation
  - Lack of predictive insights
- **UX Needs**:
  - Usage analytics dashboard
  - Capacity forecasting tools
  - Cost optimization recommendations

---

## 2. Primary Use Cases

### Use Case 1: Large-Scale Language Model Fine-Tuning

**User Persona**: Data Scientist (Mid-level experience)

**User Goal**: Fine-tune a 7B parameter LLM on enterprise data using multiple GPUs across several nodes

**Current Workarounds/Pain Points**:
- Manual configuration of MPI environment variables and networking
- Using standalone MPI installations outside OpenShift AI ecosystem
- Complex YAML files with cryptic MPI-specific settings
- No unified monitoring - switching between kubectl, TensorBoard, and custom scripts
- Difficulty troubleshooting communication failures between nodes
- Manual management of shared storage for model checkpoints

**How MPIJobs Addresses This**:
- **Simplified Configuration**: SDK provides high-level abstractions for MPI worker topology
- **Integrated Networking**: Automatic configuration of MPI communication channels within OpenShift
- **Unified Observability**: Dashboard shows launcher pod, worker pods, and their interconnections in one view
- **Smart Defaults**: Platform suggests optimal worker count and resource allocation based on model size
- **Checkpoint Management**: Automatic integration with persistent volumes using ReadWriteMany access
- **Error Handling**: User-friendly error messages translate MPI communication failures into actionable guidance

**Success Metrics**:
- Time from initial configuration to job submission reduced by 70%
- Troubleshooting time reduced by 60% through unified observability
- Successful job completion rate increased by 40%

---

### Use Case 2: Hyperparameter Optimization with Parallel Trials

**User Persona**: MLOps Engineer (Senior experience)

**User Goal**: Run parallel hyperparameter search across 50+ configurations using distributed training for each trial

**Current Workarounds/Pain Points**:
- Sequential execution due to complexity of orchestrating multiple distributed jobs
- Manual resource allocation across trials leading to underutilization
- No automated comparison of results across trials
- Difficulty tracking which hyperparameters were used in each job
- Resource contention when manually scheduling overlapping jobs
- Lost trial results due to lack of structured artifact storage

**How MPIJobs Addresses This**:
- **Batch Job Templates**: Create reusable MPIJob templates with parameterized hyperparameters
- **Intelligent Scheduling**: Integration with Kueue for fair-share scheduling and resource optimization
- **Automated Tracking**: Each trial automatically tagged with hyperparameter values and linked to results
- **Resource Pooling**: Dynamic resource allocation across trials based on priority
- **Results Dashboard**: Comparative view showing metrics across all trials with filtering/sorting
- **Artifact Management**: Automatic capture and organization of model artifacts per trial

**Success Metrics**:
- Parallel trial execution enables 10x faster hyperparameter optimization
- Resource utilization increased from 45% to 85%
- Zero lost trial results through automated artifact management
- 50% reduction in time to identify optimal configurations

---

### Use Case 3: Multi-Node Computer Vision Model Training

**User Persona**: Data Scientist (Junior-Mid level experience)

**User Goal**: Train a large computer vision model on a dataset that doesn't fit on a single GPU, requiring data parallelism across 8 GPUs on 2 nodes

**Current Workarounds/Pain Points**:
- Limited understanding of distributed training concepts (data parallel vs model parallel)
- Struggles with MPI rank configuration and GPU affinity
- Difficulty debugging when training stalls or exhibits poor scaling
- Uncertainty about network bandwidth impact on training speed
- Manual calculation of effective batch size across workers
- No visibility into per-GPU utilization and balance

**How MPIJobs Addresses This**:
- **Guided Setup**: Dashboard wizard asks simple questions (dataset size, model size, desired training time) and recommends configuration
- **Automatic GPU Binding**: Platform handles MPI rank to GPU affinity automatically
- **Visual Topology**: Dashboard shows data flow from storage through each worker to understand parallelism
- **Performance Insights**: Real-time metrics showing per-GPU utilization, network bandwidth, and training throughput
- **Scaling Validation**: Pre-flight checks warn if configuration unlikely to scale efficiently
- **Educational Context**: In-dashboard tips explain MPI concepts relevant to current job state

**Success Metrics**:
- Junior data scientists can successfully configure distributed jobs without MLOps support
- 80% reduction in support tickets related to distributed training configuration
- 90% of jobs achieve >70% scaling efficiency on first attempt
- Training time reduced from 48 hours (single GPU) to 7 hours (8 GPUs distributed)

---

### Use Case 4: Continuous Training Pipeline for Production Models

**User Persona**: MLOps Engineer (Senior experience)

**User Goal**: Implement automated retraining pipeline that triggers MPIJobs weekly with updated data and deploys models to production after validation

**Current Workarounds/Pain Points**:
- Manual job submission and monitoring in retraining cycles
- Complex orchestration scripts to trigger jobs from CI/CD systems
- No automated validation of distributed training results before deployment
- Difficult to maintain consistency in job configuration across retraining cycles
- Poor visibility into historical training runs for debugging regressions
- Manual model versioning and rollback procedures

**How MPIJobs Addresses This**:
- **SDK Integration**: Programmatic job creation and monitoring APIs for pipeline integration
- **Job Templates**: Version-controlled templates ensure consistent configuration
- **Automated Validation**: Post-training validation gates with configurable success criteria
- **Event Webhooks**: Job state change notifications for pipeline orchestration
- **Historical Tracking**: Automatic lineage tracking linking training jobs to deployed models
- **GitOps Support**: Job definitions manageable through Git workflows with audit trail

**Success Metrics**:
- 100% automated retraining pipeline with zero manual intervention
- Model retraining cadence increased from monthly to weekly
- Mean time to detect training regressions reduced from 3 days to 4 hours
- Full audit trail from training data through deployed model for compliance

---

### Use Case 5: Research Team Collaboration on Shared Cluster Resources

**User Persona**: Team Lead / Administrator managing multiple Data Scientists

**User Goal**: Enable team of 10 data scientists to share cluster resources fairly while allowing urgent jobs to be prioritized

**Current Workarounds/Pain Points**:
- Informal "scheduling" through Slack messages to avoid conflicts
- Resource hoarding by users submitting jobs "just in case"
- No visibility into who's using what resources
- Inability to prioritize business-critical training jobs
- Difficulty attributing compute costs to specific projects or users
- Manual intervention required to kill runaway or forgotten jobs

**How MPIJobs Addresses This**:
- **Resource Quotas**: Per-user and per-project limits visible in Dashboard
- **Queue Management**: Fair-share scheduling with priority overrides for urgent jobs
- **Usage Dashboard**: Real-time and historical view of resource utilization by user/project
- **Cost Attribution**: Automatic cost calculation and chargeback reporting
- **Job Policies**: Configurable timeout and resource limit policies
- **Self-Service**: Users can see queue status and adjust their job priorities within limits

**Success Metrics**:
- Cluster utilization increased from 40% to 82%
- Time-to-compute for urgent jobs reduced from 6 hours to 30 minutes
- 95% reduction in resource conflict escalations
- Complete cost transparency enabling budget planning

---

## 3. Interface Design Principles by Integration Point

### 3.1 Command Line Interface (CLI)

#### Design Philosophy
The CLI should empower power users and enable automation while maintaining consistency with OpenShift and Kubernetes conventions. It should provide both imperative commands for quick actions and declarative patterns for GitOps workflows.

#### User Experience Principles

**1. Progressive Disclosure**
- Simple commands for common operations
- Advanced flags available but not required
- Contextual help that adapts to user's command history

**2. Consistency with OpenShift Patterns**
- Follow `oc` command conventions
- Support `-o yaml|json|wide` output formats
- Use familiar flags like `--namespace`, `--selector`

**3. Intelligent Defaults**
- Infer namespace from current context
- Apply sensible resource defaults based on cluster capacity
- Support configuration files for complex scenarios

#### Information Architecture

**Command Structure**:
```
oc-ai trainjob [action] [resource] [options]

Actions: create, get, describe, delete, logs, list
Resources: mpijob, pytorchjob, tfjob (unified under trainjob)
```

**Key Commands**:

```bash
# Create from file (declarative)
oc-ai trainjob create -f mpijob.yaml

# Create interactively (imperative with intelligent prompts)
oc-ai trainjob create mpijob --name my-training \
  --image myregistry/training:latest \
  --workers 4 \
  --gpu 1 \
  --script train.py

# Get status with MPI-specific information
oc-ai trainjob get mpijob my-training
oc-ai trainjob get mpijob my-training -o wide  # Shows launcher + all workers

# Stream logs (with automatic aggregation)
oc-ai trainjob logs mpijob my-training  # Launcher logs by default
oc-ai trainjob logs mpijob my-training --all-workers  # Interleaved worker logs
oc-ai trainjob logs mpijob my-training --worker 2  # Specific worker

# Describe with topology
oc-ai trainjob describe mpijob my-training
# Output includes:
#   - Overall status
#   - Launcher pod details
#   - Worker pod details and status
#   - Resource utilization summary
#   - Recent events and errors

# List with filtering
oc-ai trainjob list --type mpijob
oc-ai trainjob list --status running
oc-ai trainjob list --user alice
```

#### Observability in CLI

**Status Output Format**:
```
NAME            TYPE      STATUS    LAUNCHER   WORKERS   AGE
my-training     MPIJob    Running   Ready      3/4       5m

Launcher: my-training-launcher-abc123
  Status: Running
  Node: worker-node-1
  Resources: 2 CPU, 4Gi RAM

Workers:
  my-training-worker-0: Running (node: worker-node-2, GPU: 0)
  my-training-worker-1: Running (node: worker-node-3, GPU: 0)
  my-training-worker-2: Running (node: worker-node-4, GPU: 0)
  my-training-worker-3: Pending (waiting for GPU resources)

Recent Events:
  2m ago: Worker 3 pending due to insufficient GPU resources
  5m ago: Launcher started successfully
  5m ago: MPIJob created
```

**Error Handling**:
- Human-readable error messages that explain Kubernetes errors in ML context
- Suggested remediation steps for common issues
- Link to relevant documentation
- Error codes for programmatic parsing

Example:
```
Error: MPIJob creation failed
Reason: InsufficientResources
Details: Requested 4 GPUs but only 2 available in namespace 'datascience'

Suggestions:
  1. Reduce worker count: --workers 2
  2. Check quota: oc-ai quota get
  3. Wait for resources: oc-ai trainjob create --wait
  4. Use priority: oc-ai trainjob create --priority high

Documentation: https://docs.openshift.com/ai/mpijobs/troubleshooting#insufficient-resources
```

---

### 3.2 Python SDK

#### Design Philosophy
The SDK should abstract Kubernetes complexity while providing flexibility for advanced users. It should feel Pythonic and integrate naturally into Jupyter notebook workflows common in data science.

#### User Experience Principles

**1. Sensible Defaults with Explicit Overrides**
- Minimal required parameters for simple cases
- Rich configuration options for advanced scenarios
- Type hints and autocomplete support

**2. Context-Aware**
- Automatic authentication when running in OpenShift
- Environment detection (local vs cluster)
- Namespace inference from current context

**3. Notebook-Friendly**
- Rich output in Jupyter (HTML tables, progress bars)
- Async support for long-running operations
- Interactive progress tracking

#### Information Architecture

**SDK Structure**:
```python
from openshift_ai_training import TrainingClient, MPIJobConfig, ResourceSpec

# Client initialization
client = TrainingClient()  # Auto-detects credentials

# Simple job creation
job = client.create_mpijob(
    name="my-training",
    image="myregistry/training:latest",
    command=["python", "train.py"],
    workers=4,
    gpu_per_worker=1
)

# Advanced configuration
config = MPIJobConfig(
    name="advanced-training",
    image="myregistry/training:latest",
    command=["python", "train.py"],
    launcher=ResourceSpec(
        cpu="2",
        memory="4Gi"
    ),
    workers=4,
    worker_resources=ResourceSpec(
        cpu="8",
        memory="32Gi",
        gpu=1,
        gpu_type="nvidia.com/gpu"
    ),
    slots_per_worker=1,
    env_vars={
        "NCCL_DEBUG": "INFO",
        "NCCL_SOCKET_IFNAME": "eth0"
    },
    volumes=[
        VolumeMount(
            name="training-data",
            mount_path="/data",
            pvc_name="shared-data"
        )
    ],
    scheduling=SchedulingConfig(
        priority="high",
        queue="ml-training",
        min_available=5  # Gang scheduling
    )
)

job = client.create_mpijob(config)
```

**Job Monitoring**:
```python
# Blocking wait with progress bar (notebook-friendly)
job.wait_for_completion(show_progress=True)

# Async monitoring
async for status in job.watch_status():
    print(f"Launcher: {status.launcher.state}")
    print(f"Workers ready: {status.workers_ready}/{status.workers_total}")
    if status.is_failed:
        print(f"Error: {status.error_message}")
        break

# Log streaming
for log_line in job.stream_logs(follow=True):
    print(log_line)

# Aggregated logs from all workers
for source, log_line in job.stream_logs(all_workers=True, follow=True):
    print(f"[{source}] {log_line}")

# Metrics retrieval
metrics = job.get_metrics()
print(f"GPU utilization: {metrics.avg_gpu_utilization}%")
print(f"Network bandwidth: {metrics.network_throughput_gbps} Gbps")
```

**Job Management**:
```python
# List jobs with filtering
jobs = client.list_mpijobs(
    namespace="datascience",
    status="running",
    labels={"team": "nlp"}
)

# Get specific job
job = client.get_mpijob("my-training")

# Job details
print(job.status)  # "Running", "Succeeded", "Failed"
print(job.launcher_pod)
print(job.worker_pods)
print(job.duration)
print(job.cost_estimate)

# Delete job
job.delete()

# Clone and resubmit with modifications
new_job = job.clone(
    name="my-training-v2",
    workers=8  # Scale up
)
```

#### Observability in SDK

**Rich Status Objects**:
```python
status = job.get_status()

# Status object structure
status.phase  # "Pending", "Running", "Succeeded", "Failed"
status.launcher.state  # Pod state
status.launcher.node  # Scheduled node
status.launcher.resources  # Actual resource allocation

status.workers  # List of worker status objects
status.workers[0].state
status.workers[0].gpu_id
status.workers[0].metrics.gpu_utilization
status.workers[0].metrics.memory_usage

status.conditions  # Kubernetes conditions with timestamps
status.events  # Recent events
status.error_message  # User-friendly error if failed
```

**Jupyter Integration**:
```python
# Rich HTML display in notebooks
job.display()  # Renders interactive job status card

# Progress bar for training
with job.progress_tracker() as tracker:
    # Automatically updates based on logs or metrics
    pass

# Inline TensorBoard
job.launch_tensorboard()  # Opens TensorBoard in notebook
```

#### Error Handling

```python
from openshift_ai_training.exceptions import (
    InsufficientResourcesError,
    JobCreationError,
    AuthenticationError
)

try:
    job = client.create_mpijob(config)
except InsufficientResourcesError as e:
    print(f"Not enough resources: {e.details}")
    print(f"Available: {e.available_resources}")
    print(f"Requested: {e.requested_resources}")
    print(f"Suggestions: {e.suggestions}")
except JobCreationError as e:
    print(f"Failed to create job: {e}")
    print(f"Kubernetes error: {e.k8s_error}")
```

---

### 3.3 OpenShift AI Dashboard UI

#### Design Philosophy
The Dashboard should be the primary interface for visual exploration, monitoring, and management of training jobs. It should serve users with varying technical expertise, from junior data scientists to platform administrators. The UI must balance simplicity for common tasks with depth for complex scenarios.

#### User Experience Principles

**1. Progressive Disclosure**
- Show essential information by default
- Provide expandable sections for details
- Use drill-down navigation for deep exploration

**2. Task-Oriented Design**
- Design around user jobs-to-be-done
- Minimize cognitive load with clear information hierarchy
- Provide contextual actions based on job state

**3. Unified Observability**
- Single-pane-of-glass for job lifecycle
- Consistent patterns across job types
- Real-time updates without page refresh

**4. Responsive Guidance**
- Contextual help and tooltips
- Error messages with remediation guidance
- Intelligent defaults based on usage patterns

#### Information Architecture

**Navigation Structure**:
```
Dashboard Home
└── Distributed Training
    ├── Jobs (List View)
    │   ├── Active Jobs
    │   ├── Completed Jobs
    │   └── Failed Jobs
    ├── Job Details (Single Job View)
    │   ├── Overview
    │   ├── Topology
    │   ├── Logs
    │   ├── Metrics
    │   └── Configuration
    ├── Create Job (Wizard)
    │   ├── Job Type Selection
    │   ├── Configuration
    │   ├── Resources
    │   └── Review & Submit
    └── Templates
        ├── My Templates
        ├── Team Templates
        └── Platform Templates
```

---

#### 3.3.1 Jobs List View

**Layout Components**:

**Top Bar**:
- Filter by job type: All | MPIJob | PyTorchJob | TFJob
- Filter by status: All | Running | Pending | Succeeded | Failed
- Search by name/label
- Date range selector
- "Create Job" button (primary action)

**Job Cards / Table View** (Toggle between views):

**Card View** (Default for Data Scientists):
```
┌─────────────────────────────────────────────────────────────┐
│ [MPI] bert-fine-tuning-2025-10-29              [⋮ Actions]   │
│ Status: Running | Started: 15 minutes ago                    │
├─────────────────────────────────────────────────────────────┤
│ Launcher: ● Running                                          │
│ Workers:  ●●●○ 3/4 Ready                                     │
│                                                              │
│ Resources: 4 GPUs | 32 CPUs | 128Gi RAM                      │
│ Progress:  ████████████░░░░░░░░  65% (estimated)            │
├─────────────────────────────────────────────────────────────┤
│ [View Details] [Logs] [Metrics] [Stop Job]                   │
└─────────────────────────────────────────────────────────────┘
```

**Table View** (Denser, preferred by MLOps):
```
Name                   Type    Status    Workers  GPUs  Age    Actions
bert-fine-tuning       MPIJob  Running   3/4      4     15m    [⋮]
llama-pretrain         MPIJob  Succeeded 8/8      8     2h     [⋮]
cv-model-training      MPIJob  Failed    0/4      0     1h     [⋮]
```

**Status Indicators**:
- Green dot: Healthy/Running
- Yellow dot: Warning (degraded performance, worker failures)
- Red dot: Failed
- Gray dot: Pending
- Checkmark: Succeeded

**Contextual Actions Menu (⋮)**:
- View Details
- View Logs
- View Metrics
- Clone Job
- Export Configuration
- Stop Job (if running)
- Delete Job

---

#### 3.3.2 Job Details View

**Page Header**:
```
← Back to Jobs

[MPI] bert-fine-tuning-2025-10-29                    [Clone] [Stop] [Delete]
Status: Running for 15 minutes

Created by: alice@company.com | Project: NLP Research | Cost: $2.45 (estimated)
```

**Tabbed Navigation**:
1. Overview
2. Topology
3. Logs
4. Metrics
5. Configuration
6. Events

---

**Tab 1: Overview**

**Job Summary Section**:
```
┌─────────────────────────────────────────────────────┐
│ Job Status                                          │
├─────────────────────────────────────────────────────┤
│ Phase: Running                                      │
│ Started: 2025-10-29 14:15:32 UTC                    │
│ Duration: 15 minutes 43 seconds                     │
│ Estimated Completion: 18 minutes remaining          │
└─────────────────────────────────────────────────────┘
```

**Component Status Section**:
```
┌─────────────────────────────────────────────────────┐
│ Launcher Pod                                        │
├─────────────────────────────────────────────────────┤
│ Name: bert-fine-tuning-launcher-xyz123              │
│ Status: ● Running                                   │
│ Node: worker-node-1.compute.internal                │
│ Resources: 2 CPU | 4Gi RAM                          │
│ [View Logs] [View in OpenShift Console]             │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ Worker Pods (3/4 Ready)                             │
├─────────────────────────────────────────────────────┤
│ ● worker-0  Running   worker-node-2  1 GPU  [Logs]  │
│ ● worker-1  Running   worker-node-3  1 GPU  [Logs]  │
│ ● worker-2  Running   worker-node-4  1 GPU  [Logs]  │
│ ○ worker-3  Pending   -              1 GPU  [Logs]  │
│                                                     │
│ ⚠ Worker 3 pending: Insufficient GPU resources      │
│   Waiting for node with available GPU               │
│   [Troubleshooting Guide]                           │
└─────────────────────────────────────────────────────┘
```

**Resource Utilization Section** (Real-time):
```
┌─────────────────────────────────────────────────────┐
│ Resource Utilization                                │
├─────────────────────────────────────────────────────┤
│ GPUs:     ████████████████░░  85% avg across 3      │
│ CPU:      ██████████████░░░░  72%                   │
│ Memory:   ████████████░░░░░░  65% (83Gi / 128Gi)    │
│ Network:  ████████░░░░░░░░░░  42% (2.1 Gbps)        │
│                                                     │
│ [View Detailed Metrics]                             │
└─────────────────────────────────────────────────────┘
```

**Quick Actions**:
- Stream Logs (all workers)
- Open TensorBoard
- Download Checkpoints
- Clone Job with Same Config

---

**Tab 2: Topology**

**Visual Topology View**:

Interactive diagram showing:
- Launcher pod at center
- Worker pods arranged around it
- Lines showing MPI communication paths
- Color-coded by health status
- Tooltips on hover showing metrics

```
           [PVC: training-data]
                    │
         ┌──────────┴──────────┐
         │                     │
    [Launcher]            [Metrics Store]
         │
    ┌────┼────┬────┬────┐
    │    │    │    │    │
 [W-0] [W-1] [W-2] [W-3 Pending]
   │    │    │    │
   └────┴────┴────┴─── [Shared Storage]
```

**Node Placement View**:
- Visual representation of physical node distribution
- Shows which workers are co-located
- Highlights potential network bottlenecks
- Shows GPU affinity

**MPI Communication Metrics**:
- Message passing frequency
- Data transfer rates between workers
- Synchronization overhead
- Load balance across workers

---

**Tab 3: Logs**

**Log Viewer Interface**:

**Component Selector** (Tabs or Dropdown):
- All (interleaved with source labels)
- Launcher
- Worker 0
- Worker 1
- Worker 2
- Worker 3

**Toolbar**:
- [Live Tail] toggle
- [Download Logs] button
- [Filter by Level] dropdown: All | INFO | WARNING | ERROR
- [Search] input field
- [Wrap Lines] toggle
- [Timestamps] toggle

**Log Display**:
```
┌─────────────────────────────────────────────────────┐
│ 14:15:35 [launcher] INFO: Starting MPI job          │
│ 14:15:38 [worker-0] INFO: Initializing rank 0       │
│ 14:15:38 [worker-1] INFO: Initializing rank 1       │
│ 14:15:38 [worker-2] INFO: Initializing rank 2       │
│ 14:15:42 [launcher] INFO: All workers connected     │
│ 14:15:45 [worker-0] INFO: Epoch 1/10 started        │
│ 14:16:12 [worker-0] INFO: Epoch 1 loss: 0.453       │
│ 14:16:30 [worker-1] WARNING: GPU memory 92%         │
│ ...                                                 │
└─────────────────────────────────────────────────────┘
```

**Log Intelligence Features**:
- Automatic error highlighting
- Pattern recognition (e.g., "out of memory" gets special treatment)
- Suggested actions for common errors
- Ability to share log snippets with permalinks

---

**Tab 4: Metrics**

**Metrics Dashboard**:

**Time Range Selector**: Last 5m | 15m | 1h | 6h | 24h | Custom

**Metric Categories** (Expandable sections):

**1. Training Progress**:
- Loss curve (if available from logs/TensorBoard)
- Accuracy/metrics over time
- Throughput (samples/second)
- Estimated completion time

**2. GPU Utilization**:
- Per-worker GPU utilization (line chart)
- GPU memory usage (line chart)
- GPU temperature (line chart)
- CUDA errors/warnings

**3. System Resources**:
- CPU utilization per pod
- Memory usage per pod
- Network I/O (send/receive)
- Disk I/O (read/write)

**4. MPI-Specific Metrics**:
- Communication overhead percentage
- All-reduce operation latency
- Worker synchronization delays
- Load imbalance metrics

**5. Cost & Efficiency**:
- Estimated cost (running total)
- Cost per epoch
- Resource efficiency score
- Scaling efficiency (actual vs ideal speedup)

**Visualization Options**:
- Line charts for time-series
- Heatmaps for per-worker comparisons
- Gauge charts for current values
- Comparison overlays for multiple jobs

**Export Options**:
- Download as CSV
- Export to Prometheus/Grafana
- Share dashboard link

---

**Tab 5: Configuration**

**Read-Only Configuration Display** (with copy button):

**Expandable Sections**:

**Job Metadata**:
```
Name:        bert-fine-tuning-2025-10-29
Namespace:   nlp-research
Created By:  alice@company.com
Labels:      team=nlp, experiment=fine-tuning, version=v2
```

**Job Spec**:
```yaml
apiVersion: kubeflow.org/v2
kind: TrainJob
metadata:
  name: bert-fine-tuning-2025-10-29
spec:
  runtime:
    mpi:
      launcher:
        resources:
          cpu: 2
          memory: 4Gi
      workers: 4
      slotsPerWorker: 1
      workerResources:
        cpu: 8
        memory: 32Gi
        gpu: 1
        gpuType: nvidia.com/gpu
  trainer:
    image: myregistry/bert-training:latest
    command:
      - python
      - train.py
      - --data-path=/data
      - --model-name=bert-base
    env:
      - name: NCCL_DEBUG
        value: INFO
  storage:
    volumes:
      - name: training-data
        persistentVolumeClaim:
          claimName: shared-data
        mountPath: /data
```

**Actions**:
- [Copy as YAML]
- [Copy as JSON]
- [Copy as SDK Code]
- [Download Configuration]
- [Clone Job]

**SDK Code Example** (Expandable):
Shows equivalent Python SDK code to recreate this job

---

**Tab 6: Events**

**Event Timeline**:

Chronological list of Kubernetes events related to the job:

```
┌─────────────────────────────────────────────────────┐
│ ▼ 2 minutes ago                                     │
│   Worker-3 still pending due to insufficient GPU    │
│   Type: Warning | Reason: FailedScheduling          │
│                                                     │
│ ▼ 5 minutes ago                                     │
│   Worker-2 successfully started on worker-node-4    │
│   Type: Normal | Reason: Started                    │
│                                                     │
│ ▼ 5 minutes ago                                     │
│   Worker-1 successfully started on worker-node-3    │
│   Type: Normal | Reason: Started                    │
│                                                     │
│ ▼ 15 minutes ago                                    │
│   Launcher pod assigned to worker-node-1            │
│   Type: Normal | Reason: Scheduled                  │
│                                                     │
│ ▼ 15 minutes ago                                    │
│   MPIJob created successfully                       │
│   Type: Normal | Reason: Created                    │
└─────────────────────────────────────────────────────┘
```

**Event Filters**:
- Show: All | Normal | Warning | Error
- Source: All | Scheduler | Controller | Kubelet

---

#### 3.3.3 Create Job Wizard

**Multi-Step Form with Progress Indicator**:

```
Step 1: Job Type → Step 2: Configuration → Step 3: Resources → Step 4: Review
  ●              →          ○           →         ○          →       ○
```

---

**Step 1: Job Type Selection**

**Visual Card Selection**:

```
┌────────────────────┐  ┌────────────────────┐  ┌────────────────────┐
│   [MPI Icon]       │  │  [PyTorch Icon]    │  │ [TensorFlow Icon]  │
│                    │  │                    │  │                    │
│   MPIJob          │  │  PyTorchJob        │  │  TFJob             │
│                    │  │                    │  │                    │
│ Best for:          │  │ Best for:          │  │ Best for:          │
│ • MPI-based models │  │ • PyTorch models   │  │ • TensorFlow       │
│ • HPC workloads    │  │ • Custom training  │  │ • Legacy workflows │
│ • Multi-node       │  │ • Flexible setups  │  │                    │
│                    │  │                    │  │                    │
│  [Select MPIJob]   │  │  [Select]          │  │  [Select]          │
└────────────────────┘  └────────────────────┘  └────────────────────┘
```

**Help Text**:
"Not sure which to choose? MPIJob is recommended for distributed training that requires tight coupling and high-performance communication between nodes."

**Template Option**:
"Or start from a template: [Dropdown: My Templates | Team Templates]"

---

**Step 2: Configuration**

**Form Sections**:

**Basic Information**:
```
Job Name: [bert-fine-tuning________________]
         ℹ Must be unique, lowercase, alphanumeric with dashes

Description: [Fine-tuning BERT on customer data_______]

Labels (optional):
  Key: [experiment________]  Value: [fine-tuning_____] [+ Add]
  Key: [team_____________]  Value: [nlp____________] [+ Add]
```

**Container Configuration**:
```
Container Image: [myregistry/training:latest_________] [Browse Registry]

Command: [python train.py --data-path=/data --epochs=10_____________]
         ℹ The command to run in the container

Working Directory: [/workspace_____________] (optional)
```

**Environment Variables** (Expandable):
```
[+ Add Environment Variable]

Name: [NCCL_DEBUG________]  Value: [INFO__________]  [×]
Name: [BATCH_SIZE________]  Value: [32____________]  [×]
```

**Secrets & Config Maps** (Expandable):
```
[+ Add Secret as Environment]
[+ Add ConfigMap as Environment]
[+ Mount Secret as Volume]
[+ Mount ConfigMap as Volume]
```

---

**Step 3: Resources**

**Worker Configuration**:
```
Number of Workers: [4_] ℹ How many parallel workers to run

⚙ Auto-suggest based on: [Dropdown: Dataset Size | Model Size | Desired Training Time]
```

**Resource Allocation Per Worker**:

```
┌─────────────────────────────────────────────────────┐
│ CPU:    [8__] cores                                 │
│         ▓▓▓▓▓▓▓▓░░░░░░░░ 8/16 available per node   │
│                                                     │
│ Memory: [32__] GiB                                  │
│         ▓▓▓▓▓▓▓▓░░░░░░░░ 32/64 available per node  │
│                                                     │
│ GPU:    [1__] units                                 │
│         Type: [nvidia.com/gpu ▼]                    │
│         ▓▓▓▓▓▓▓▓▓▓▓▓░░░░ 4/6 GPUs available         │
└─────────────────────────────────────────────────────┘
```

**Launcher Pod Resources** (Collapsible, with defaults):
```
CPU:    [2__] cores
Memory: [4__] GiB
(Launcher typically needs fewer resources than workers)
```

**Storage**:
```
[+ Add Volume]

Existing Persistent Volume Claims:
┌─────────────────────────────────────────────────────┐
│ ☑ shared-training-data                              │
│   Mount Path: [/data________________]               │
│   Access Mode: ReadWriteMany | Size: 500Gi          │
└─────────────────────────────────────────────────────┘

[+ Create New PVC]
```

**Advanced Scheduling** (Collapsible):
```
Priority Class: [Default ▼]
Queue: [ml-training ▼]
Min Available Pods: [5_] (Gang scheduling)
Node Selector: [+ Add]
Tolerations: [+ Add]
Affinity Rules: [+ Add]
```

**Cost Estimate** (Dynamic):
```
┌─────────────────────────────────────────────────────┐
│ Estimated Cost                                      │
│ Based on current configuration:                     │
│   4 workers × 1 GPU × $1.50/hour = $6.00/hour       │
│   Launcher resources: $0.20/hour                    │
│                                                     │
│   Total: $6.20/hour                                 │
│                                                     │
│ For 2-hour training run: $12.40                     │
└─────────────────────────────────────────────────────┘
```

**Validation Messages**:
```
✓ Configuration is valid
⚠ Warning: Requesting more than 50% of available GPUs
  Consider lowering worker count for better resource sharing
```

---

**Step 4: Review & Submit**

**Configuration Summary**:

```
┌─────────────────────────────────────────────────────┐
│ Job Configuration Summary                           │
├─────────────────────────────────────────────────────┤
│ Job Type:    MPIJob                                 │
│ Name:        bert-fine-tuning                       │
│ Image:       myregistry/training:latest             │
│                                                     │
│ Workers:     4                                      │
│ Resources per Worker:                               │
│   - CPU: 8 cores                                    │
│   - Memory: 32 GiB                                  │
│   - GPU: 1 (nvidia.com/gpu)                         │
│                                                     │
│ Launcher:    2 CPU, 4 GiB RAM                       │
│                                                     │
│ Storage:     shared-training-data → /data           │
│                                                     │
│ Environment Variables:                              │
│   - NCCL_DEBUG=INFO                                 │
│   - BATCH_SIZE=32                                   │
│                                                     │
│ Estimated Cost: $6.20/hour                          │
│                                                     │
│ [← Back]  [Save as Template]  [Create Job]          │
└─────────────────────────────────────────────────────┘
```

**Post-Submission**:
```
┌─────────────────────────────────────────────────────┐
│ ✓ Job Created Successfully!                         │
│                                                     │
│ Job "bert-fine-tuning" has been submitted.          │
│ Current status: Pending (waiting for resources)     │
│                                                     │
│ [View Job Details]  [Create Another Job]            │
└─────────────────────────────────────────────────────┘
```

---

#### 3.3.4 Admin Dashboard

**Platform-Level Overview** (For Administrators):

**Cluster Resource Dashboard**:
```
┌─────────────────────────────────────────────────────┐
│ Cluster Utilization                                 │
├─────────────────────────────────────────────────────┤
│ GPUs:     ████████████████░░  85% (17/20 in use)    │
│ CPU:      ████████████░░░░░░  65% (650/1000 cores)  │
│ Memory:   ███████████░░░░░░░  58% (2.3/4.0 TiB)     │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ Active Training Jobs                                │
├─────────────────────────────────────────────────────┤
│ MPIJobs:     12 running, 3 pending                  │
│ PyTorchJobs: 8 running, 1 pending                   │
│ TFJobs:      2 running, 0 pending                   │
└─────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────┐
│ Job Success Rate (Last 7 Days)                      │
├─────────────────────────────────────────────────────┤
│ MPIJobs:     87% (132/152)                          │
│ PyTorchJobs: 92% (185/201)                          │
│ TFJobs:      89% (42/47)                            │
└─────────────────────────────────────────────────────┘
```

**Per-Project Breakdown**:
```
Project          Jobs  GPUs  Cost/Day  Quota Usage
nlp-research      12    8    $456      80% GPU
computer-vision    8    6    $342      75% GPU
recommender-sys    3    3    $171      30% GPU
```

**Queue Management**:
```
Priority Queue Status:
High Priority:   2 running, 1 pending
Normal:         18 running, 4 pending
Low Priority:    5 running, 8 pending

[Manage Quotas]  [Configure Policies]
```

---

## 4. Consistency Requirements: Aligning with Existing KubeFlow Trainer Jobs

### 4.1 Visual Consistency

**Status Indicators**:
- Reuse existing color scheme for job states (green/yellow/red/gray)
- Consistent iconography for job types (add MPI icon to existing set)
- Uniform card layouts and spacing

**Typography & Styling**:
- Follow OpenShift AI PatternFly design system
- Maintain existing font sizes, weights, and hierarchies
- Consistent spacing and padding

**Navigation Patterns**:
- MPIJobs appear in same "Distributed Training" section as PyTorchJob/TFJob
- Identical navigation structure within job details
- Consistent breadcrumb patterns

---

### 4.2 Terminology Consistency

**Unified Language Across Job Types**:

| Concept | Consistent Term | Notes |
|---------|----------------|-------|
| Job instance | Job | Not "training run" or "workload" |
| Job state | Status | Values: Pending, Running, Succeeded, Failed |
| Worker processes | Workers | Not "replicas" or "processes" |
| Primary controller | Launcher (MPI) / Master (PyTorch) | Type-specific but clearly labeled |
| Resource request | Resources | Consistent format: "8 CPU, 32Gi RAM, 1 GPU" |
| Job lifetime | Duration | Format: "2h 15m 43s" |

**Avoid Framework-Specific Jargon in Shared UI**:
- Use "distributed training job" not "MPI job" in general contexts
- Framework-specific terms only in job-type-specific views
- Tooltips explain technical terms (e.g., "MPI: Message Passing Interface")

---

### 4.3 Interaction Patterns

**Job Lifecycle Actions** (Consistent across types):
- Create → Configure → Submit → Monitor → Complete
- Same action buttons: [Clone], [Stop], [Delete], [View Logs]
- Same keyboard shortcuts where applicable
- Consistent confirmation dialogs

**List View Behaviors**:
- Same filters available: Type, Status, Date, User
- Same sorting options: Name, Status, Age, Resources
- Same bulk actions: Delete, Export
- Same search functionality

**Detail View Structure**:
- Identical tab structure: Overview | Topology | Logs | Metrics | Configuration | Events
- Same expandable sections
- Consistent action menus

---

### 4.4 Data Model Consistency

**Job Metadata** (Same fields across types):
- Name, Namespace, Labels, Annotations
- Created By, Created At, Updated At
- Owner References
- Status, Conditions, Events

**Resource Representation** (Uniform format):
```yaml
resources:
  cpu: "8"
  memory: "32Gi"
  gpu: 1
  gpuType: "nvidia.com/gpu"
```

**Status Schema** (Consistent structure):
```yaml
status:
  phase: "Running"  # Pending | Running | Succeeded | Failed
  startTime: "2025-10-29T14:15:32Z"
  completionTime: null
  conditions:
    - type: "Running"
      status: "True"
      lastTransitionTime: "2025-10-29T14:15:35Z"
  componentStatuses:
    launcher:
      state: "Running"
      ready: true
    workers:
      total: 4
      ready: 3
      failed: 0
```

---

### 4.5 MPI-Specific Adaptations

While maintaining consistency, MPIJobs require some specific accommodations:

**Topology Visualization**:
- **Unique to MPI**: Launcher-worker architecture diagram
- **Consistent with**: Similar topology views exist for PyTorch master-worker
- **Implementation**: Same visualization framework, MPI-specific layout

**MPI Communication Metrics**:
- **Unique to MPI**: MPI-specific metrics (all-reduce latency, message passing)
- **Consistent with**: Metrics tab structure same as other job types
- **Implementation**: Additional metric categories in same UI framework

**Configuration Options**:
- **Unique to MPI**: `slotsPerWorker`, MPI-specific environment variables
- **Consistent with**: Same configuration form structure, MPI fields appear conditionally
- **Implementation**: Dynamic form that adapts to job type

**Log Aggregation**:
- **Unique to MPI**: Clear launcher vs worker distinction
- **Consistent with**: Same log viewer interface as other multi-pod jobs
- **Implementation**: Source selector includes "Launcher" alongside workers

---

### 4.6 Documentation & Help Consistency

**Contextual Help**:
- Same tooltip icon and style
- Consistent help text structure: Brief explanation → Example → Link to docs
- Same "Learn More" link pattern

**Error Messages**:
- Same error message format: [Icon] Brief description → Details → Suggested action
- Consistent severity levels: Info, Warning, Error
- Same documentation link structure

**Onboarding**:
- MPIJob introduction follows same tutorial pattern as existing job types
- Consistent use of interactive walkthroughs
- Same documentation structure and TOC

---

### 4.7 API Consistency

**SDK Method Naming**:
```python
# Consistent pattern across job types
client.create_mpijob(...)     # Not create_mpi_job or createMPIJob
client.get_mpijob(...)
client.list_mpijobs(...)
client.delete_mpijob(...)

# Same as
client.create_pytorchjob(...)
client.create_tfjob(...)
```

**CLI Command Structure**:
```bash
# Consistent pattern
oc-ai trainjob create mpijob ...
oc-ai trainjob create pytorchjob ...
oc-ai trainjob create tfjob ...

# Not
oc-ai mpi create ...  # Inconsistent!
```

**Status Fields**:
- Same status structure across job types
- Consistent field names and types
- Compatible with existing monitoring tools

---

### 4.8 Migration Path for Existing Users

**Backwards Compatibility**:
- Existing MPIJob CRDs (if any) continue to work
- Gradual migration to unified TrainJob API
- Clear migration documentation

**Familiar Patterns**:
- Users familiar with PyTorchJob will recognize MPIJob patterns
- Same learning curve and mental models
- Transferable skills across job types

**Upgrade Experience**:
- New MPIJob features appear alongside existing jobs
- No disruption to existing workflows
- Clear "What's New" documentation

---

## 5. Error Handling & User Feedback

### 5.1 Error Taxonomy

Categorize errors by user impact and resolution:

**Category 1: User Configuration Errors** (Most Common)
- Invalid resource requests
- Malformed YAML/configuration
- Missing required fields
- Invalid container image references
- Storage access issues

**Category 2: Infrastructure Limitations**
- Insufficient cluster resources
- Quota exceeded
- Scheduling constraints not satisfiable
- Network policy restrictions

**Category 3: Runtime Failures**
- Container crashes
- Out of memory errors
- Application-level errors
- MPI communication failures
- Timeout errors

**Category 4: Platform Issues**
- API server errors
- Controller failures
- Operator issues
- Infrastructure outages

---

### 5.2 Error Presentation Principles

**1. User-Centric Language**
- Translate technical errors into user-friendly explanations
- Focus on impact and resolution, not internal details
- Avoid Kubernetes jargon where possible

**2. Actionable Guidance**
- Always provide next steps
- Offer multiple resolution paths when available
- Link to relevant documentation

**3. Contextual Error Messages**
- Show errors in context of user's action
- Preserve user input when errors occur
- Indicate which field/section has the issue

**4. Progressive Disclosure**
- Show summary error first
- Provide expandable details for technical users
- Include raw error for support escalation

---

### 5.3 Job Lifecycle Feedback

**Job Creation Phase**:

**Validation Feedback** (Real-time as user types):
```
Job Name: [bert-fine-tuning!]
❌ Invalid character '!'. Use lowercase letters, numbers, and dashes only.
```

**Pre-Submit Validation** (Before submission):
```
⚠ Configuration Warning
  Requesting 8 GPUs but only 6 available in cluster.

  Options:
  • Reduce workers to 6 [Apply]
  • Submit anyway and wait for resources [Continue]
  • Save as draft [Save Draft]
```

**Submission Success**:
```
✓ Job Submitted Successfully
  "bert-fine-tuning" is now pending.

  Current queue position: 2 of 5
  Estimated start time: ~5 minutes

  [View Job Details]
```

**Submission Failure**:
```
❌ Job Submission Failed
  Quota exceeded: GPU requests exceed namespace limit.

  Your request: 8 GPUs
  Current usage: 14 GPUs
  Namespace limit: 20 GPUs
  Available: 6 GPUs

  Suggestions:
  • Reduce worker count to 6 or fewer [Adjust]
  • Wait for other jobs to complete [View Active Jobs]
  • Request quota increase from administrator [Contact Admin]

  [View Details] [Edit Configuration]
```

---

**Pending Phase**:

**Queue Status** (Dashboard):
```
Status: Pending (3 minutes)

⏳ Waiting for Resources
   Required: 4 GPUs
   Available: 2 GPUs (2 GPUs in use by other jobs)

   Queue position: 1 of 3
   Estimated start: ~5 minutes

   Jobs ahead of you:
   • nlp-training-xyz (requires 2 GPUs, 2min remaining)

   [View Queue] [Increase Priority] [Cancel Job]
```

**Scheduling Issues**:
```
⚠ Scheduling Delayed
  Worker pods pending for 10 minutes.

  Reason: Insufficient GPU resources
  Details: Cluster has 6 GPUs, all currently allocated.

  Options:
  • Wait for resources to free up (estimated: 15min)
  • Reduce worker count [Edit Configuration]
  • Use CPU-only mode (slower) [Switch to CPU]
  • Cancel job [Cancel]

  [Troubleshooting Guide]
```

---

**Running Phase**:

**Healthy Running State** (Dashboard):
```
Status: Running (15 minutes)

✓ Job is healthy
  Launcher: Running
  Workers: 4/4 Ready

  Progress: 65% (estimated)
  GPU utilization: 85% avg

  [View Logs] [View Metrics]
```

**Degraded Performance Warning**:
```
⚠ Performance Warning
  GPU utilization lower than expected (35% vs 80% typical).

  Possible causes:
  • Data loading bottleneck
  • Network bandwidth limitation
  • Suboptimal batch size

  Recommendations:
  • Check data loader performance [View Logs]
  • Review network metrics [View Metrics]
  • Consider increasing batch size

  [Troubleshooting Guide] [Dismiss]
```

**Worker Failure** (Non-Fatal):
```
⚠ Worker Failed
  worker-2 crashed and is restarting.

  Status: 3/4 workers healthy
  Impact: Training continues with reduced capacity

  Last log from worker-2:
  "CUDA out of memory. Tried to allocate 12.00 GiB..."

  This is retry 1 of 3.

  [View Worker Logs] [Terminate Job]
```

**Job Hanging Detection**:
```
⚠ Possible Hang Detected
  No log output for 10 minutes. Job may be stuck.

  Common causes:
  • Distributed training deadlock
  • Network communication failure
  • Application waiting for external resource

  Suggestions:
  • Check worker logs [View All Logs]
  • Verify MPI communication [View Topology]
  • Consider restarting job [Restart]

  [Dismiss] [Get Help]
```

---

**Completion Phase**:

**Successful Completion**:
```
✓ Job Completed Successfully
  "bert-fine-tuning" finished in 2h 15m.

  Final Metrics:
  • Training loss: 0.124
  • Validation accuracy: 94.3%
  • Cost: $14.35

  Artifacts:
  • Model checkpoint: /data/checkpoints/final.pt
  • TensorBoard logs: /data/logs/

  [View Results] [Deploy Model] [Clone Job]
```

**Job Failed**:
```
❌ Job Failed
  "bert-fine-tuning" failed after 1h 23m.

  Failure Reason: Launcher pod terminated with error
  Exit Code: 137 (Out of Memory)

  Last logs:
  "RuntimeError: CUDA out of memory. Tried to allocate 12.00 GiB..."

  Suggestions:
  • Reduce batch size [Clone with Smaller Batch]
  • Increase worker memory [Clone with More Memory]
  • Enable gradient checkpointing [Documentation]

  Similar successful jobs for reference:
  • bert-fine-tuning-v1 (batch size: 16, success)

  [View Full Logs] [Clone & Fix] [Get Help]
```

**Partial Checkpoint Recovery**:
```
⚠ Job Failed, Checkpoints Available
  Job failed but model checkpoints were saved.

  Latest checkpoint: epoch 7 of 10 (70% complete)
  Location: /data/checkpoints/epoch-7/

  Options:
  • Resume from checkpoint [Resume Job]
  • Use partial results [Download Checkpoint]
  • Start fresh [Clone Job]

  [View Checkpoints]
```

---

### 5.4 Specific Error Scenarios & Messages

**Scenario: Insufficient GPU Resources**

**Dashboard Message**:
```
❌ Scheduling Failed: Insufficient GPU Resources

Your request: 4 GPUs (1 per worker)
Available: 2 GPUs in namespace 'datascience'

Current GPU allocation:
• nlp-training-xyz: 1 GPU (running for 45min)
• cv-model-train: 1 GPU (running for 2h)

Options:
1. Reduce worker count to 2
   [Clone with 2 Workers]

2. Wait for resources (estimated: 15 minutes)
   [Keep Job Pending]

3. Request higher priority
   [Request Priority Increase]

4. Check other namespaces
   [View Cluster-Wide Resources]

[Contact Administrator] [View Documentation]
```

**CLI Message**:
```
Error: Insufficient GPU resources

Requested: 4 GPUs
Available: 2 GPUs

Active jobs using GPUs:
  nlp-training-xyz    1 GPU   45m
  cv-model-train      1 GPU   2h

Suggestions:
  1. Reduce workers: oc-ai trainjob create ... --workers 2
  2. Wait for resources: oc-ai trainjob create ... --wait
  3. Check quotas: oc-ai quota get -n datascience

Documentation: https://docs.openshift.com/ai/troubleshooting#insufficient-gpu
```

---

**Scenario: MPI Communication Failure**

**Dashboard Message**:
```
❌ Job Failed: MPI Communication Error

Worker-2 failed to connect to launcher via MPI.

Technical Details:
  Error: "MPI_Init failed: Connection refused"
  Affected Worker: worker-2 (10.0.3.45)
  Launcher: launcher-pod (10.0.3.12)

Possible Causes:
• Network policy blocking MPI communication ports
• Istio sidecar injection interfering (common issue)
• Firewall rules on worker nodes

Troubleshooting Steps:
1. Check network policy [View Network Policies]
2. Verify Istio configuration [Check Istio Settings]
   → Recommended: Add annotation "sidecar.istio.io/inject: false"
3. Review worker logs [View Worker-2 Logs]

[Recreate with Istio Disabled] [View Full Error] [Get Help]
```

**SDK Error**:
```python
MPIJobError: Job failed due to MPI communication error

Details:
  Job: bert-fine-tuning-2025-10-29
  Error: MPI_Init failed on worker-2

Suggested Fix:
  Add the following to your job configuration:

  metadata:
    annotations:
      sidecar.istio.io/inject: "false"

Code example:
  config = MPIJobConfig(
      name="my-job",
      annotations={"sidecar.istio.io/inject": "false"},
      ...
  )

Documentation: https://docs.openshift.com/ai/mpijobs/istio-compatibility
```

---

**Scenario: Container Image Pull Failure**

**Dashboard Message**:
```
❌ Job Failed to Start: Cannot Pull Container Image

Image: myregistry/training:latest
Error: ImagePullBackOff

Common Causes:
• Image does not exist or tag is incorrect
• Authentication required but not provided
• Registry is unreachable from cluster

Troubleshooting:
1. Verify image exists
   Command: podman pull myregistry/training:latest

2. Check image pull secrets
   [View Image Pull Secrets in Namespace]

3. Test registry connectivity
   [Test Registry Connection]

[Edit Image Name] [Add Pull Secret] [View Full Error]
```

---

**Scenario: Out of Memory Error**

**Dashboard Message**:
```
❌ Job Failed: Out of Memory (OOM)

Worker-0 terminated due to insufficient memory.

Memory Stats:
  Requested: 32 GiB
  Peak usage: 31.8 GiB (99%)
  Limit: 32 GiB

Analysis:
  Your job used almost all allocated memory before failing.
  This suggests the model or batch size is too large for current allocation.

Recommendations (choose one):
1. Increase worker memory
   [Clone with 64 GiB Memory]

2. Reduce batch size
   Current: 32 (inferred from logs)
   Suggested: 16
   [Clone with Smaller Batch]

3. Enable gradient checkpointing (if using PyTorch)
   [View Documentation]

4. Use mixed precision training (FP16)
   [View Documentation]

Similar successful configurations:
• bert-fine-tuning-v1: 64 GiB memory, batch size 16 ✓

[View Memory Usage Graph] [Clone & Fix] [Get Help]
```

---

**Scenario: Storage Access Denied**

**Dashboard Message**:
```
❌ Job Failed to Start: Storage Access Error

Cannot mount volume "training-data".

Error: "Access denied for PVC 'shared-data'"

Cause:
  The PVC 'shared-data' exists but your service account lacks permissions.

Resolution:
1. Check PVC exists in correct namespace
   [View PVCs in 'datascience']

2. Verify access permissions
   Required: ReadWriteMany access
   [Check PVC Permissions]

3. Request access from administrator
   [Contact Administrator]

Alternative:
  Use a different PVC: [Browse Available PVCs]

[View Full Error] [Edit Storage Configuration]
```

---

### 5.5 Proactive Feedback Mechanisms

**Pre-Flight Validation** (Before job submission):

```
Running pre-flight checks...

✓ Container image is accessible
✓ Requested resources within quota
✓ Storage volumes are accessible
⚠ Network policy may block MPI communication
  → Recommendation: Add Istio exclusion annotation
✓ GPU type is available in cluster

Configuration score: 4/5 (Good)

[View Recommendations] [Submit Anyway] [Apply Fixes]
```

**Predictive Warnings**:

```
⚠ Training May Be Slow
  Based on similar jobs, your configuration may result in slower training.

  Your config:
  • 4 workers
  • Network bandwidth: 1 Gbps

  Similar job "cv-training-xyz" with same config had:
  • Scaling efficiency: 45% (expected: >70%)
  • Training time: 8 hours (50% longer than expected)

  Suggestion: Use nodes with 10 Gbps networking for better performance.

  [View Network Topology] [Request Better Nodes] [Proceed Anyway]
```

**Resource Optimization Suggestions**:

```
💡 Optimization Opportunity
  Your job is running but could be more efficient.

  Observation: GPU utilization is only 40% on average.
  Typical for this workload: 75-85%

  Likely cause: Data loading bottleneck

  Suggestions:
  • Increase data loader workers
  • Pre-cache data to faster storage
  • Use data prefetching

  Estimated savings: 30% reduction in training time

  [View Optimization Guide] [Dismiss]
```

---

### 5.6 Feedback Timing & Delivery

**Immediate Feedback** (< 1 second):
- Form validation errors
- Client-side configuration checks
- Syntax errors

**Fast Feedback** (< 10 seconds):
- Server-side validation
- Resource availability checks
- Pre-flight validations
- Quota checks

**Deferred Feedback** (Minutes):
- Scheduling outcomes
- Initial pod starts
- Runtime configuration issues

**Long-Running Feedback** (During execution):
- Performance warnings
- Resource utilization alerts
- Progress updates
- Completion notifications

---

### 5.7 Notification Channels

**In-Dashboard Notifications**:
- Toast messages for immediate actions
- Alert banners for important warnings
- Badge indicators on navigation items
- Status changes in job cards

**Email Notifications** (Configurable):
- Job completion (success/failure)
- Long-pending jobs
- Critical errors requiring attention
- Scheduled maintenance affecting jobs

**Webhook Integration**:
- Job state changes
- Error events
- Completion events
- Custom event filters

**CLI Output**:
- Immediate error messages
- Progress indicators for long operations
- Final status on completion

---

## 6. Accessibility & Usability Considerations

### 6.1 Skill Level Accommodations

**For Junior Data Scientists**:

**Guided Workflows**:
- Step-by-step wizard for job creation
- Contextual help bubbles explaining concepts
- Example configurations with annotations
- "Recommended" badges on common options

**Educational Scaffolding**:
- Inline explanations of MPI concepts
- Visual diagrams showing distributed architecture
- Links to tutorials and getting-started guides
- Glossary of terms with tooltips

**Safe Defaults**:
- Pre-validated, working configurations
- Automatic resource suggestions based on workload
- Error prevention through validation
- Undo/rollback capabilities

**Example: Simplified Job Creation**:
```
┌─────────────────────────────────────────────────────┐
│ Create Distributed Training Job - Simple Mode      │
├─────────────────────────────────────────────────────┤
│                                                     │
│ What are you training?                              │
│ ◉ Language Model (e.g., BERT, GPT)                 │
│ ○ Computer Vision Model (e.g., ResNet, YOLO)       │
│ ○ Other                                             │
│                                                     │
│ How large is your model?                            │
│ ○ Small (< 1B parameters)  → Suggests 2 workers    │
│ ◉ Medium (1-10B parameters) → Suggests 4 workers   │
│ ○ Large (> 10B parameters) → Suggests 8 workers    │
│                                                     │
│ Container image with your training code:            │
│ [myregistry/training:latest_________________]       │
│                                                     │
│ [Advanced Options ▼]  [Create Job]                  │
│                                                     │
│ 💡 Tip: Not sure about configuration?              │
│    Start with our template: [Use BERT Template]     │
└─────────────────────────────────────────────────────┘
```

---

**For Experienced ML Engineers**:

**Power User Features**:
- Direct YAML editing with validation
- CLI with advanced flags and scripting support
- SDK with full programmatic control
- Batch operations and job templates

**Advanced Monitoring**:
- Detailed metrics and performance analytics
- Custom dashboards and alerts
- Integration with external monitoring tools
- Exportable data for offline analysis

**Optimization Tools**:
- Performance profiling
- Cost analysis and optimization suggestions
- Comparative analysis across experiments
- Historical trend analysis

**Example: Advanced Configuration**:
```
┌─────────────────────────────────────────────────────┐
│ Create MPIJob - Expert Mode                        │
├─────────────────────────────────────────────────────┤
│ [Visual Editor] [YAML Editor] ◉                     │
│                                                     │
│ apiVersion: kubeflow.org/v2                         │
│ kind: TrainJob                                      │
│ metadata:                                           │
│   name: optimized-training                          │
│   annotations:                                      │
│     sidecar.istio.io/inject: "false"                │
│     kueue.x-k8s.io/queue-name: high-priority        │
│ spec:                                               │
│   runtime:                                          │
│     mpi:                                            │
│       launcher:                                     │
│         resources: {cpu: 4, memory: 8Gi}            │
│       workers: 8                                    │
│       slotsPerWorker: 2                             │
│       workerResources:                              │
│         nvidia.com/gpu: 2                           │
│         rdma/hca: 1  # RDMA for high-perf           │
│   ...                                               │
│                                                     │
│ ✓ Valid | CPU: 36 | Memory: 256Gi | GPUs: 16       │
│ [Import from File] [Validate] [Submit]             │
└─────────────────────────────────────────────────────┘
```

---

**For Platform Administrators**:

**Governance & Control**:
- Policy enforcement (quotas, resource limits)
- User and team management
- Cost allocation and chargeback
- Audit logging and compliance

**Operations Dashboard**:
- Cluster-wide resource view
- Job queue management
- Performance trends and analytics
- Capacity planning tools

**Configuration Management**:
- Platform-wide defaults and templates
- Integration with CI/CD pipelines
- Infrastructure-as-code support
- Version control and rollback

**Example: Admin Dashboard**:
```
┌─────────────────────────────────────────────────────┐
│ MPIJob Platform Administration                     │
├─────────────────────────────────────────────────────┤
│                                                     │
│ [Cluster Overview] [Quotas] [Policies] [Users]     │
│                                                     │
│ Resource Utilization by Team:                      │
│                                                     │
│ NLP Research        ████████████░░░░  75% of quota  │
│   12 jobs, 8 GPUs, $456/day                         │
│   [View Details] [Adjust Quota]                     │
│                                                     │
│ Computer Vision     ██████████░░░░░░  60% of quota  │
│   8 jobs, 6 GPUs, $342/day                          │
│   [View Details] [Adjust Quota]                     │
│                                                     │
│ Recent Policy Violations:                           │
│ ⚠ User alice@company attempted job requiring 12    │
│   GPUs (quota: 10). Job blocked.                    │
│   [Review] [Adjust Policy]                          │
│                                                     │
│ [Configure Platform Defaults] [View Audit Logs]     │
└─────────────────────────────────────────────────────┘
```

---

### 6.2 Web Accessibility (WCAG 2.1 Level AA Compliance)

**Keyboard Navigation**:
- Full keyboard support for all interactive elements
- Logical tab order through forms and controls
- Keyboard shortcuts for common actions (with visual indicators)
- Skip links to jump to main content

**Screen Reader Support**:
- Semantic HTML with proper ARIA labels
- Descriptive alt text for all images and icons
- Status announcements for dynamic content updates
- Proper heading hierarchy (h1 → h2 → h3)

**Visual Accessibility**:
- Minimum 4.5:1 contrast ratio for text
- Color not used as sole indicator (icons + color)
- Resizable text up to 200% without loss of functionality
- Focus indicators with clear visual contrast

**Content Clarity**:
- Plain language, avoiding unnecessary jargon
- Consistent terminology throughout interface
- Clear error messages and instructions
- Logical information hierarchy

**Responsive Design**:
- Usable on various screen sizes
- Touch-friendly targets (minimum 44×44 pixels)
- Zoom-friendly layouts
- Mobile-optimized views for monitoring

---

### 6.3 Internationalization (i18n)

**Language Support**:
- English as default
- Framework for additional languages
- Right-to-left (RTL) layout support when needed
- Locale-aware date/time formatting

**Cultural Considerations**:
- Culturally neutral iconography
- Locale-appropriate number formatting
- Time zone awareness for global teams
- Currency formatting for cost estimates

---

### 6.4 Performance & Responsiveness

**Fast Page Loads**:
- Initial page load < 2 seconds
- Lazy loading for heavy components
- Optimized asset delivery (CDN)
- Progressive rendering

**Real-Time Updates**:
- WebSocket connections for live job status
- Optimistic UI updates
- Background data refreshing
- Smooth animations and transitions

**Scalability**:
- Efficient rendering of large job lists (virtualization)
- Pagination for historical data
- Filtering and search at API level
- Caching strategies for frequently accessed data

---

### 6.5 Error Recovery & Data Preservation

**Autosave & Draft Support**:
- Automatic saving of in-progress job configurations
- "Resume where you left off" capability
- Draft jobs saved locally or server-side
- Conflict resolution for concurrent edits

**Undo/Redo**:
- Configuration changes are reversible
- Clear undo/redo controls
- History of configuration changes
- Version comparison

**Graceful Degradation**:
- Core functionality works without JavaScript
- Fallbacks for unsupported features
- Offline mode with limited functionality
- Clear messaging when features unavailable

---

### 6.6 Discoverability & Learnability

**Onboarding Experience**:
- Optional interactive tutorial for new users
- "What's New" announcements for feature updates
- Contextual tooltips on first use
- Sample jobs and templates for learning

**In-App Documentation**:
- Context-sensitive help links
- Embedded documentation snippets
- Video tutorials and demos
- Searchable knowledge base

**Progressive Disclosure**:
- Hide advanced features by default
- "Show advanced options" toggles
- Collapsible sections for detailed information
- Gradual exposure to complexity

**Search & Navigation**:
- Global search across jobs, documentation, and help
- Breadcrumb navigation
- Recently viewed jobs
- Bookmarking and favorites

---

### 6.7 Collaboration & Sharing

**Team Workflows**:
- Shared job templates across teams
- Commenting on job configurations
- @mentions for notifications
- Shared dashboards and views

**Knowledge Sharing**:
- Export job configurations as shareable links
- Embed job details in documentation
- Screenshot-friendly views for presentations
- Exportable reports (PDF, CSV)

**Access Control**:
- Role-based access (viewer, editor, admin)
- Project-based isolation
- Audit trails for sensitive actions
- Fine-grained permissions

---

## 7. RFE Sections: Synthesized Outputs

### 7.1 Use Cases Summary

**UC1: Large-Scale LLM Fine-Tuning**
- **User**: Data Scientist
- **Goal**: Fine-tune 7B parameter model across multiple GPUs/nodes
- **Value**: 70% reduction in configuration time, 40% increase in job success rate

**UC2: Hyperparameter Optimization at Scale**
- **User**: MLOps Engineer
- **Goal**: Run 50+ parallel distributed training trials
- **Value**: 10x faster optimization, 85% resource utilization (vs 45%)

**UC3: Multi-Node CV Model Training**
- **User**: Junior-Mid Data Scientist
- **Goal**: Train large CV model with data parallelism across 8 GPUs
- **Value**: Enables junior users to self-serve without MLOps support, 7 hours vs 48 hours training time

**UC4: Continuous Training Pipeline**
- **User**: MLOps Engineer
- **Goal**: Automated weekly retraining with MPIJobs in CI/CD
- **Value**: 100% automation, weekly vs monthly retraining cadence

**UC5: Shared Cluster Resource Management**
- **User**: Team Lead/Administrator
- **Goal**: Fair-share scheduling across 10 data scientists
- **Value**: 82% utilization (vs 40%), 95% reduction in resource conflicts

---

### 7.2 User Experience & Workflow

**Key Workflows**:

1. **Job Creation Flow**:
   - Discovery → Configuration → Validation → Submission → Monitoring
   - Multiple entry points: UI Wizard, SDK, CLI, Templates
   - Guided experience for beginners, power features for experts

2. **Monitoring Flow**:
   - Real-time status updates in Dashboard
   - Unified observability: topology, logs, metrics, events
   - Proactive alerting for issues
   - Historical tracking and comparison

3. **Troubleshooting Flow**:
   - User-friendly error messages with suggested fixes
   - Contextual help and documentation links
   - Diagnostic tools and runbooks
   - Support escalation paths

4. **Iteration Flow**:
   - Job cloning with configuration tweaks
   - Template creation from successful jobs
   - A/B comparison of different configurations
   - Integration with model registry

**Cross-Persona Consistency**:
- Same job object across CLI, SDK, Dashboard
- Consistent terminology and mental models
- Seamless handoffs between personas (e.g., Data Scientist creates, MLOps optimizes)

---

### 7.3 Requirements (UX Perspective)

**Dashboard UI Requirements**:
1. Unified job list view showing MPIJobs alongside PyTorch/TF jobs
2. Drill-down job details with topology, logs, metrics, configuration tabs
3. Visual job creation wizard with intelligent defaults
4. Real-time status updates via WebSocket
5. MPI-specific visualizations (launcher-worker topology, communication metrics)
6. Responsive design for mobile monitoring
7. WCAG 2.1 AA compliance
8. Template library with sharing capabilities

**SDK Requirements**:
1. Pythonic API with type hints and autocomplete
2. Jupyter notebook integration (rich displays, progress bars)
3. Blocking and async monitoring patterns
4. Comprehensive error handling with actionable exceptions
5. Auto-authentication in OpenShift environment
6. Job cloning and templating APIs
7. Metrics and log streaming APIs

**CLI Requirements**:
1. Consistent with OpenShift `oc` command patterns
2. Imperative and declarative job creation
3. Human-readable status output with MPI-specific details
4. Log aggregation across launcher and workers
5. JSON/YAML output for scripting
6. Contextual error messages with suggestions

**Observability Requirements**:
1. Real-time metrics: GPU utilization, MPI communication overhead, network bandwidth
2. Aggregated logs with source filtering (launcher vs workers)
3. Event timeline with Kubernetes and application events
4. Performance analytics and efficiency scoring
5. Cost attribution and tracking
6. Alerting for failures, performance degradation, resource issues

**Consistency Requirements**:
1. Visual consistency with existing KubeFlow Trainer jobs (PyTorch, TF)
2. Terminology alignment across job types
3. Unified status schema and lifecycle states
4. Same interaction patterns and navigation structure
5. Backwards-compatible API design

**Error Handling Requirements**:
1. User-friendly error messages translating Kubernetes errors
2. Suggested remediation steps for common issues
3. Pre-flight validation before job submission
4. Contextual error presentation with progressive disclosure
5. Multiple notification channels (dashboard, email, webhooks)

**Accessibility Requirements**:
1. Full keyboard navigation support
2. Screen reader compatibility
3. 4.5:1 minimum contrast ratio
4. Resizable text without functionality loss
5. Localization framework for i18n

**Skill Level Requirements**:
1. Guided wizard mode for junior users
2. Direct YAML editing for experts
3. Contextual help and educational content
4. Progressive disclosure of advanced features
5. Template system for common patterns

---

### 7.4 Acceptance Criteria (User Perspective)

**From Data Scientist Perspective**:

1. **Job Creation**:
   - AC: I can create an MPIJob through Dashboard wizard in < 5 minutes without reading documentation
   - AC: I can specify workers, GPUs, and container image in simple form
   - AC: System suggests resource configuration based on my model size
   - AC: I receive immediate validation errors with clear fix suggestions
   - AC: I can clone a previous job and modify parameters easily

2. **Job Monitoring**:
   - AC: I can see real-time status of launcher and all workers in single dashboard view
   - AC: I can view aggregated logs from all pods with filtering by source
   - AC: I see clear progress indicator (estimated % completion)
   - AC: I receive alerts when job fails with user-friendly explanation
   - AC: I can view GPU utilization and training metrics in real-time

3. **Troubleshooting**:
   - AC: Error messages are in plain language, not Kubernetes jargon
   - AC: When job fails, I see suggested fixes relevant to the error type
   - AC: I can compare failed job configuration with similar successful jobs
   - AC: I have access to comprehensive logs and events in one place
   - AC: I can get help through in-dashboard links to relevant documentation

4. **Results & Iteration**:
   - AC: I can access trained models and checkpoints easily after completion
   - AC: I can compare metrics across multiple training runs
   - AC: I can create new job based on previous successful configuration
   - AC: Training duration and cost are clearly displayed
   - AC: I can save job configuration as template for future use

---

**From MLOps Engineer Perspective**:

1. **Pipeline Integration**:
   - AC: I can create and monitor MPIJobs programmatically via Python SDK
   - AC: SDK provides blocking and async APIs for pipeline orchestration
   - AC: I can trigger jobs from CI/CD systems with clear success/failure signals
   - AC: Job state changes can trigger webhooks for downstream processes
   - AC: I can batch-create jobs with parameterized templates

2. **Resource Management**:
   - AC: I can configure priority queues and fair-share scheduling
   - AC: I can set resource quotas per team/project
   - AC: I receive alerts when cluster utilization is suboptimal
   - AC: I can see cost attribution per job, team, and project
   - AC: I can enforce resource limits through policies

3. **Performance Optimization**:
   - AC: I can view MPI-specific metrics (communication overhead, scaling efficiency)
   - AC: I can compare performance across different configurations
   - AC: System suggests optimizations based on historical data
   - AC: I can export metrics to external monitoring tools (Prometheus, Grafana)
   - AC: I can identify bottlenecks (compute, network, I/O) through detailed metrics

4. **Template & Governance**:
   - AC: I can create team-wide templates with validated configurations
   - AC: I can enforce mandatory parameters (e.g., cost center tags)
   - AC: I can version and audit template changes
   - AC: I can share best-practice configurations across organization
   - AC: I can restrict certain advanced options for non-admin users

---

**From Platform Administrator Perspective**:

1. **Installation & Configuration**:
   - AC: I can install MPIJob support via standard OpenShift operator pattern
   - AC: Installation validates cluster prerequisites (network, storage, GPU drivers)
   - AC: I can configure platform-wide defaults for MPIJobs
   - AC: I receive clear documentation on upgrade path from standalone MPI Operator
   - AC: I can test installation with sample jobs

2. **Operations & Monitoring**:
   - AC: I have a dedicated admin dashboard showing all active MPIJobs across projects
   - AC: I can monitor cluster-wide GPU, CPU, memory utilization
   - AC: I receive alerts for platform-level issues (operator failures, resource exhaustion)
   - AC: I can view job success/failure rates and trends over time
   - AC: I can identify resource-hogging jobs and take action

3. **User Management**:
   - AC: I can assign role-based access (viewer, editor, admin) per project
   - AC: I can set quotas per user or team
   - AC: I can audit who created which jobs and when
   - AC: I can provide self-service onboarding materials for users
   - AC: I can disable certain features for specific user groups

4. **Troubleshooting & Support**:
   - AC: I can view detailed error logs and diagnostics for failed jobs
   - AC: I can correlate job issues with cluster events
   - AC: I have runbooks for common MPIJob failure scenarios
   - AC: I can export audit logs for compliance
   - AC: I can simulate resource constraints to test policies

---

**Cross-Cutting Acceptance Criteria**:

1. **Consistency**:
   - AC: MPIJobs appear in same UI section as PyTorchJobs and TFJobs
   - AC: Job detail views have identical tab structure across job types
   - AC: Status indicators use same color scheme and iconography
   - AC: Terminology is consistent (e.g., "Workers" not "Replicas")
   - AC: SDK method names follow same pattern as other job types

2. **Accessibility**:
   - AC: All interactive elements are keyboard-navigable
   - AC: Screen readers can navigate and understand all UI elements
   - AC: Text contrast meets WCAG 2.1 AA standards (4.5:1 minimum)
   - AC: UI works with text zoom up to 200%
   - AC: Error messages are announced to screen readers

3. **Performance**:
   - AC: Dashboard loads in < 2 seconds
   - AC: Job status updates in < 5 seconds
   - AC: Log streaming is real-time with < 1 second lag
   - AC: Large job lists (100+) render smoothly with virtualization
   - AC: Metrics queries return in < 3 seconds

4. **Error Handling**:
   - AC: 100% of Kubernetes errors are translated to user-friendly messages
   - AC: Every error message includes at least one suggested remediation
   - AC: Critical errors (e.g., job failures) trigger notifications
   - AC: Users can access raw error details if needed
   - AC: Common error patterns are documented with solutions

5. **Documentation & Help**:
   - AC: Every form field has contextual help tooltip
   - AC: Error messages link to relevant documentation
   - AC: Getting-started tutorial is accessible from dashboard
   - AC: API documentation includes examples for all common operations
   - AC: In-app search finds relevant help articles

---

## 8. Ecosystem Integration Considerations

### 8.1 Integration with OpenShift AI Components

**Integration with ODH Dashboard**:
- MPIJobs listed alongside existing training jobs
- Unified navigation and branding
- Shared authentication and authorization
- Consistent design system (PatternFly)

**Integration with Model Registry**:
- Automatic registration of trained models
- Link training jobs to registered models
- Model versioning and lineage tracking
- Metadata tagging (hyperparameters, metrics, cost)

**Integration with Model Serving**:
- Streamlined deployment of trained models
- One-click "Deploy Model" from completed job
- Connection between training job and inference endpoint
- Performance comparison (training vs serving metrics)

**Integration with Data Science Projects**:
- MPIJobs scoped to DS Projects
- Shared storage and notebooks within projects
- Project-level quotas and permissions
- Unified project dashboard

**Integration with Notebooks (Jupyter)**:
- SDK available in notebook environments
- Launch jobs directly from notebooks
- Monitor jobs inline with rich displays
- TensorBoard integration within notebooks

---

### 8.2 Integration with Kubernetes Ecosystem

**Kueue Integration** (Job Queueing):
- Fair-share scheduling across teams
- Priority-based queue management
- Resource quotas and preemption
- Multi-cluster job distribution

**JobSet Integration**:
- Coordinated multi-component jobs
- Dependency management between jobs
- Batch job patterns

**Volcano Integration** (Gang Scheduling):
- All-or-nothing pod scheduling
- Efficient resource allocation for distributed jobs
- Prevention of partial job starts

**Prometheus/Grafana Integration**:
- Metrics export for custom dashboards
- Alerting rules for job failures
- Historical metrics analysis
- SLO/SLI tracking

**OpenTelemetry Integration**:
- Distributed tracing across job components
- Performance profiling
- Span correlation for debugging

---

### 8.3 Integration with ML Tooling

**TensorBoard**:
- Automatic TensorBoard launch from Dashboard
- Centralized TensorBoard for all experiments
- Metrics comparison across jobs
- Embedded TensorBoard in notebook environments

**MLflow / Kubeflow Pipelines**:
- Job tracking and experiment management
- Parameter and metric logging
- Model versioning and registry
- Pipeline step integration

**Git Integration**:
- Job configurations stored in Git
- GitOps workflows for job submission
- Version-controlled templates
- Audit trail through Git history

**CI/CD Integration**:
- Jenkins/Tekton pipeline steps for job creation
- Automated job submission on code changes
- Integration tests for training jobs
- Model deployment pipelines

---

## 9. Mental Models & Information Architecture

### 9.1 User Mental Models

**Data Scientist Mental Model**:
```
Training Job
  └─ Configuration (image, command, resources)
  └─ Execution (running on cluster)
  └─ Results (model, metrics, logs)

Workflow: Configure → Submit → Monitor → Iterate
```

**MLOps Mental Model**:
```
Training Pipeline
  └─ Templates (reusable configurations)
  └─ Jobs (instantiated from templates)
  └─ Metrics (performance, cost, efficiency)
  └─ Governance (quotas, policies, access)

Workflow: Design → Automate → Optimize → Govern
```

**Administrator Mental Model**:
```
Platform
  └─ Users & Teams (access control)
  └─ Resources (compute, storage, network)
  └─ Jobs (running workloads)
  └─ Policies (quotas, limits, priorities)

Workflow: Provision → Monitor → Optimize → Support
```

---

### 9.2 Conceptual Hierarchy

**Top Level**: Distributed Training Jobs
**Second Level**: Job Types (MPIJob, PyTorchJob, TFJob)
**Third Level**: Job Instance (specific training run)
**Fourth Level**: Components (Launcher, Workers)
**Fifth Level**: Resources (Pods, Volumes, Networks)

**User Navigation Path**:
```
Dashboard Home
  → Distributed Training
    → Jobs (filtered by type, status)
      → Job Details
        → Component Status (Launcher, Workers)
          → Logs / Metrics / Events
```

---

### 9.3 Information Scent

Users should be able to predict what they'll find before clicking:

**Good Information Scent**:
- "View Logs" → Shows aggregated logs from all pods
- "Metrics" → Shows performance metrics over time
- "Configuration" → Shows job YAML/configuration
- "Clone Job" → Creates copy with same configuration

**Poor Information Scent (Avoid)**:
- "Details" → Too vague, what kind of details?
- "More Info" → Unclear what additional information
- "Advanced" → Unclear what advanced features

---

## 10. Future Considerations & Extensibility

### 10.1 Planned Enhancements

**Enhanced Observability**:
- Distributed tracing across MPI ranks
- Flame graphs for performance profiling
- Anomaly detection for training metrics
- Predictive failure detection

**Advanced Scheduling**:
- Multi-cluster job distribution
- Spot instance integration with preemption handling
- Dynamic resource scaling during training
- Cost-optimized scheduling

**Collaboration Features**:
- Shared experiments and results
- Commenting and annotations on jobs
- Team dashboards and views
- Real-time collaboration on configurations

**AI-Assisted Configuration**:
- ML-powered configuration recommendations
- Automatic hyperparameter suggestions
- Resource optimization based on historical data
- Failure prediction and prevention

---

### 10.2 Extensibility Points

**Plugin Architecture**:
- Custom metrics exporters
- Third-party monitoring integrations
- Custom job validators
- Notification channel plugins

**API Extensibility**:
- Custom resource definitions for specialized workloads
- Webhook integrations for lifecycle events
- REST API for external tools
- GraphQL API for complex queries

**UI Extensibility**:
- Custom dashboard widgets
- Embedded third-party visualizations
- Themeable UI components
- Mobile app support

---

## Conclusion

This UX architecture analysis provides a comprehensive framework for integrating MPIJobs into RedHat OpenShift AI with KubeFlow Trainer V2. The design prioritizes:

1. **User-Centric Design**: Workflows tailored to each persona's goals and skill level
2. **Unified Experience**: Consistent patterns across CLI, SDK, and Dashboard
3. **Observability**: Comprehensive monitoring with actionable insights
4. **Accessibility**: Inclusive design for users of all abilities
5. **Ecosystem Integration**: Seamless fit within OpenShift AI and Kubernetes ecosystems

The approach balances simplicity for common tasks with depth for advanced use cases, ensuring that junior data scientists can self-serve while experts retain full control. By abstracting infrastructure complexity and providing intelligent guidance, this design reduces time-to-value and increases success rates for distributed training workloads.

The architecture is extensible and future-proof, with clear patterns for adding new features, integrations, and customizations as the platform evolves.

---

## Appendix: Key Files Referenced

- **Current workspace**: `/workspace/sessions/agentic-session-1761761022/workspace/test-ambient`
- **Branch**: `ambient-mpi-jobs-4`
- **Related repositories**:
  - `/workspace/sessions/agentic-session-1761761022/workspace/trainer` (KubeFlow Trainer)
  - `/workspace/sessions/agentic-session-1761761022/workspace/mpi-operator` (MPI Operator)

---

**Document Version**: 1.0
**Last Updated**: 2025-10-29
**Next Review**: Upon RFE approval
