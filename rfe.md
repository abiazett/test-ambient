# MPIJobs Support for RedHat OpenShift AI using KubeFlow Trainer V2

**Feature Overview:**

RedHat OpenShift AI will integrate native support for MPIJobs through KubeFlow Trainer V2, enabling data scientists and MLOps engineers to run distributed MPI-based training workloads across multiple nodes with full observability and enterprise-grade management capabilities. This feature unlocks the ability to efficiently train large-scale models (including LLMs) using Message Passing Interface (MPI) protocols on hybrid cloud infrastructure, addressing a critical gap that currently forces customers to choose between single-node training bottlenecks or moving workloads to competitive cloud platforms.

**The Value Statement**: By providing integrated MPIJobs support through CLI, SDK, and ODH Dashboard UI, this feature reduces distributed training time by 60%, unlocks $2M+ of stranded on-premises GPU infrastructure value per customer, and enables OpenShift AI to compete directly with AWS SageMaker, Azure ML, and Google Vertex AI in the enterprise MLOps market. Data scientists gain the ability to scale training from single-node to multi-node workloads without leaving their familiar OpenShift AI environment, while platform administrators benefit from unified observability, resource management, and governance across all KubeFlow Trainer job types.

---

## Goals

**Primary Goals:**

1. **Democratize Distributed Training**: Enable data scientists of all skill levels to launch and manage multi-node MPI training jobs through intuitive interfaces (CLI, SDK, Dashboard UI) without requiring deep MPI or Kubernetes expertise.

2. **Accelerate Model Development**: Reduce training times by 60% or more by enabling efficient multi-node distributed training with 85%+ scaling efficiency, allowing data scientists to iterate faster on large-scale models.

3. **Unlock Hybrid Cloud Value**: Enable customers to fully utilize existing on-premises GPU infrastructure for distributed training, avoiding forced migrations to public cloud platforms and supporting data sovereignty requirements.

4. **Unified Observability**: Provide seamless monitoring and management of MPIJobs alongside other KubeFlow Trainer job types (TFJob, PyTorchJob, etc.) through a single pane of glass in the ODH Dashboard.

5. **Enterprise-Grade Operations**: Deliver MPIJobs support with production-ready capabilities including security, multi-tenancy, resource quotas, audit logging, and high availability.

**Who Benefits and How:**

- **Data Scientists**: Can scale training workloads from single-node to multi-node with minimal code changes, using familiar frameworks (Horovod, DeepSpeed, PyTorch DDP with MPI backend) through simple CLI commands or Python SDK calls. They gain 3-5x faster training iterations and can tackle larger models and datasets.

- **MLOps Engineers**: Gain programmatic control over distributed training workflows, enabling integration with CI/CD pipelines, automated hyperparameter tuning at scale, and continuous training systems. They benefit from standardized APIs and improved resource utilization (80%+ GPU utilization vs. 40-50% with single-node training).

- **OpenShift AI Administrators**: Achieve unified governance and observability across all training workloads, with centralized monitoring, resource management, and security policies. They can maximize ROI on GPU infrastructure and reduce operational complexity by consolidating distributed training on the OpenShift AI platform.

**Current State vs. Future State:**

| Dimension | Current State (Without MPIJobs) | Future State (With MPIJobs) |
|-----------|--------------------------------|----------------------------|
| **Training Scale** | Single-node only, limited by single GPU/node memory | Multi-node distributed training across 2-100+ nodes |
| **Training Speed** | 3-5x slower for large models | 60% reduction in training time with near-linear scaling |
| **GPU Utilization** | 40-50% utilization due to memory constraints | 80%+ utilization across all nodes |
| **Framework Support** | Limited to non-MPI frameworks | Full support for Horovod, MPI-based PyTorch DDP, DeepSpeed |
| **Hybrid Cloud** | Distributed training requires public cloud migration | Leverage on-prem GPUs with same capabilities as cloud |
| **Observability** | Fragmented monitoring across multiple tools | Unified dashboard for all KubeFlow Trainer job types |
| **User Experience** | Manual Kubernetes resource creation, complex debugging | Simple CLI/SDK/UI interfaces with guided workflows |
| **Competitive Position** | Feature gap vs. AWS/Azure/GCP | Feature parity with leading cloud ML platforms |

---

## Out of Scope

**Explicitly Out of Scope for MVP:**

1. **Non-MPI Distributed Training Patterns**: Alternative distributed training approaches (parameter servers, gossip protocols) that don't use MPI are out of scope. Users should use existing KubeFlow Trainer job types (TFJob, PyTorchJob with native distributed strategies) for those patterns.

2. **Custom MPI Implementations**: Support is limited to standard MPI implementations (Open MPI, Intel MPI, MPICH) as provided by KubeFlow MPI Operator. Custom or proprietary MPI libraries are not supported.

3. **Automatic Hyperparameter Optimization**: While MPIJobs can be used as part of HPO workflows, automatic HPO orchestration (e.g., Katib integration) is deferred to a future phase.

4. **Multi-Cluster Federation**: Running a single MPIJob across multiple OpenShift clusters is out of scope. Jobs are constrained to a single cluster.

5. **Windows Worker Nodes**: MPIJobs support is limited to Linux-based worker nodes. Windows containers are not supported for MPI workloads.

6. **Legacy KubeFlow Versions**: Integration is specifically for KubeFlow Trainer V2 (v2.0+). Backward compatibility with KubeFlow Training Operator V1 or standalone MPI Operator is not provided.

7. **Non-Kubernetes Deployment Models**: Integration with non-Kubernetes environments (e.g., Slurm, PBS, bare-metal MPI clusters) is out of scope.

**Out of Scope Personas:**

- **HPC System Administrators**: Users managing traditional HPC infrastructure with Slurm/PBS are not the target persona. This feature targets Kubernetes-native workflows.

- **Non-OpenShift Kubernetes Users**: While technically possible to install KubeFlow Trainer V2 on any Kubernetes distribution, this RFE specifically targets RedHat OpenShift AI customers. Upstream KubeFlow community support is handled separately.

---

## Requirements

**MVP Requirements** (Must-Have for Initial Release):

**R1. Core MPIJob Orchestration** [MVP]
- Integrate KubeFlow MPI Operator v2 as a managed component within OpenShift AI operator lifecycle
- Support creation, monitoring, and deletion of MPIJobs with launcher/worker pod architecture
- Enable configuration of worker replica count (2-100 workers), resource requests/limits (CPU, memory, GPU), and MPI-specific parameters (slots per worker)
- Provide automatic injection of MPI hostfile and SSH key management for inter-worker communication

**R2. CLI Interface** [MVP]
- Extend existing OpenShift AI CLI (`oc` plugin or dedicated `odh` CLI) with MPIJob commands: `create`, `get`, `describe`, `logs`, `delete`
- Support declarative job definition via YAML files compatible with KubeFlow Trainer V2 MPIJob schema
- Provide imperative commands for common use cases (e.g., `odh mpijob create --name my-training --workers 4 --image my-image:latest --command "python train.py"`)
- Enable real-time log streaming from launcher and worker pods with filtering capabilities

**R3. Python SDK** [MVP]
- Provide Python SDK (`odh-sdk` or `openshift-ai-sdk`) with `MPIJob` class supporting full lifecycle management
- Enable Pythonic job definition with type hints and validation (e.g., `job = MPIJob(name="my-job", workers=4, image="...", command=[...]); job.create()`)
- Support synchronous and asynchronous job submission, status checking, and log retrieval
- Integrate with Jupyter notebooks for interactive distributed training workflows

**R4. ODH Dashboard UI** [MVP]
- Add "MPIJobs" section to ODH Dashboard with list view showing all MPIJobs in user's namespace(s)
- Implement job creation wizard with guided form (name, image, command, workers, resources) and advanced YAML editor
- Provide job details page with 6 tabs: Overview (status, duration, progress), Configuration (YAML view), Workers (pod list with status), Logs (launcher + worker logs with filtering), Metrics (resource usage), Events (Kubernetes events)
- Display unified job list combining MPIJobs with other KubeFlow Trainer job types (TFJob, PyTorchJob) with filtering and sorting

**R5. Unified Observability** [MVP]
- Integrate MPIJob metrics (job status, pod status, resource usage) into existing OpenShift AI monitoring stack (Prometheus/Grafana)
- Provide real-time job status updates in Dashboard UI (Pending → Running → Succeeded/Failed) with automatic refresh
- Enable aggregated log viewing across launcher and all worker pods with timestamp synchronization
- Expose key metrics: job duration, worker scaling efficiency, GPU utilization per worker, network bandwidth between workers

**R6. Security and Multi-Tenancy** [MVP]
- Enforce namespace isolation with RBAC policies (users can only create/view MPIJobs in authorized namespaces)
- Support ServiceAccount-based authentication for launcher and worker pods
- Integrate with OpenShift NetworkPolicies to restrict inter-pod communication to authorized MPIJob workers
- Audit logging for all MPIJob lifecycle events (create, delete, completion) via OpenShift audit logs

**R7. Resource Management** [MVP]
- Support Kubernetes ResourceQuotas and LimitRanges for MPIJob pods (CPU, memory, GPU limits per namespace)
- Integrate with Kueue (if available) for fair-share scheduling and job queuing across multiple users
- Enable node affinity and taints/tolerations for GPU node targeting
- Provide clear error messages when resource requests exceed quota or no suitable nodes are available

**R8. Documentation and Examples** [MVP]
- Provide quickstart guide with end-to-end example (training a PyTorch model with Horovod on 4 workers)
- Document CLI commands, SDK API reference, and Dashboard UI workflows
- Include troubleshooting guide for common issues (SSH connectivity failures, MPI version mismatches, GPU allocation problems)
- Publish reference architectures for typical MPIJob configurations (2-node, 8-node, 16-node setups)

**Post-MVP Requirements** (Future Enhancements):

**R9. Advanced Scheduling** [Post-MVP]
- Gang scheduling support to ensure all workers are scheduled atomically (avoid partial allocations)
- Priority-based preemption for high-priority training jobs
- Elastic training support (dynamic worker scaling during job execution)

**R10. Hyperparameter Optimization Integration** [Post-MVP]
- Native integration with Katib for multi-trial HPO using MPIJobs
- Dashboard UI for tracking HPO experiments with multiple MPIJob trials

**R11. Model Registry Integration** [Post-MVP]
- Automatic model artifact registration in ODH Model Registry upon MPIJob completion
- Metadata tracking linking trained models to originating MPIJob configuration

**R12. Advanced Networking** [Post-MVP]
- Support for RDMA (InfiniBand, RoCE) for high-performance inter-worker communication
- Network topology awareness for optimal worker placement

**R13. Checkpoint and Fault Tolerance** [Post-MVP]
- Automatic checkpoint saving to persistent storage at configurable intervals
- Job restart from last checkpoint on transient failures (e.g., node eviction)

---

## Done - Acceptance Criteria

**The feature is considered complete and shippable when all of the following criteria are met:**

**AC1. Data Scientist Persona - End-to-End Job Execution**
- A data scientist with basic Python knowledge can submit an MPIJob using the SDK from a Jupyter notebook with ≤10 lines of code
- The submitted job with 4 workers and 1 GPU per worker successfully trains a PyTorch model using Horovod
- Job completes in ≤30% of the time required for equivalent single-node training (demonstrating distributed training speedup)
- Data scientist can view real-time job status and logs in ODH Dashboard without leaving the browser
- Upon job failure, error messages clearly indicate root cause (e.g., "Worker 2 failed: CUDA out of memory") and suggest remediation

**AC2. MLOps Engineer Persona - Programmatic Workflow Integration**
- An MLOps engineer can create an MPIJob using CLI commands in a CI/CD pipeline script with declarative YAML configuration
- The CLI provides exit codes and structured output (JSON) that can be parsed by automation tools
- The engineer can query job status, retrieve logs, and delete jobs programmatically using SDK or CLI
- SDK provides async/await patterns for non-blocking job submission and monitoring in Python-based orchestration frameworks

**AC3. Administrator Persona - Multi-Tenant Governance**
- A platform administrator can set namespace-level resource quotas (e.g., max 32 GPUs) that apply to MPIJobs
- When a user attempts to exceed quota, they receive a clear error message before any pods are created
- Administrator can view all MPIJobs across all namespaces in ODH Dashboard admin view with filtering by user/namespace
- Audit logs capture all MPIJob creation/deletion events with user identity and timestamp

**AC4. Unified Observability**
- ODH Dashboard "Training Jobs" page displays MPIJobs in the same list/table as TFJobs, PyTorchJobs, and other KubeFlow Trainer job types
- Users can filter jobs by type (MPIJob, TFJob, etc.), status (Running, Succeeded, Failed), and date range
- Clicking on an MPIJob opens a details page with 6 tabs (Overview, Configuration, Workers, Logs, Metrics, Events) populated with live data
- Logs tab displays synchronized logs from launcher and all worker pods with timestamp alignment and keyword search
- Metrics tab shows GPU utilization, memory usage, and network bandwidth for each worker with auto-refreshing graphs

**AC5. End-to-End Use Case - LLM Fine-Tuning**
- A data scientist can fine-tune a 7B parameter LLM (e.g., Llama 2) on 8 nodes with 4 GPUs each (32 GPUs total) using DeepSpeed with MPI backend
- Job configuration specifies worker count (8), GPUs per worker (4), container image with DeepSpeed + MPI, and training script command
- Job successfully launches all 32 worker pods, establishes MPI communication mesh, and begins training
- Training achieves ≥85% scaling efficiency (speedup on 32 GPUs is ≥27× vs. single GPU)
- Upon completion, model checkpoint is saved to persistent storage (PVC) and accessible for inference deployment

**AC6. Error Handling and Resilience**
- If a worker pod fails during training (e.g., OOMKilled), the entire MPIJob transitions to "Failed" state within 30 seconds
- Dashboard UI displays clear failure indication with root cause (e.g., "Worker 3 OOMKilled - requested memory: 16Gi, max usage: 18Gi before termination")
- Logs from the failed worker are retained and accessible via UI and CLI for debugging
- User can edit and resubmit the job with increased memory allocation without losing previous configuration

**AC7. Security and Isolation**
- MPIJob worker pods in namespace A cannot communicate with worker pods in namespace B (enforced by NetworkPolicies)
- SSH keys used for MPI inter-worker communication are unique per job and ephemeral (deleted with job)
- Non-root container images can be used for MPIJob workers (no requirement for privileged containers)
- ServiceAccount tokens used by launcher/worker pods have minimal RBAC permissions (no cluster-wide access)

**AC8. Documentation and Usability**
- Quickstart guide enables a new user to run their first successful MPIJob within 30 minutes (measured via user testing)
- API documentation includes runnable code examples for all common SDK operations (create, monitor, delete)
- Troubleshooting guide includes solutions for 10+ common error scenarios with step-by-step remediation
- Dashboard UI includes in-context help tooltips and links to relevant documentation for each configuration field

**AC9. Performance and Scalability**
- Platform can support 50+ concurrent MPIJobs across multiple namespaces without API server degradation
- Job submission via CLI/SDK completes within 5 seconds (time from command execution to "Job Created" confirmation)
- Dashboard UI loads job list page (100 jobs) within 2 seconds and updates job status every 5 seconds
- Worker pods start within 30 seconds of job submission when resources are available (excluding image pull time)

**AC10. Integration with Existing OpenShift AI Components**
- MPIJobs can be launched from notebooks running in ODH Jupyter environment with SDK pre-installed
- MPIJob pods can mount existing PVCs for training data and model checkpoint storage
- MPIJob metrics are visible in existing Prometheus/Grafana dashboards alongside other workload metrics
- MPIJob creation/deletion events appear in OpenShift Events tab and ODH Dashboard notifications

---

## Use Cases - i.e. User Experience & Workflow

**Use Case 1: Large-Scale LLM Fine-Tuning for Domain Adaptation**

**Persona**: Senior Data Scientist at Financial Services Company

**Goal**: Fine-tune a 13B parameter language model on proprietary financial documents (100GB training data) to create a domain-specific model for regulatory document analysis.

**Current Pain Point**: Single-node training with 4 A100 GPUs takes 18 days per training run, limiting experimentation to 1-2 runs per month. Model barely fits in memory with batch size of 1, resulting in poor training efficiency.

**Workflow with MPIJobs**:

1. **Setup Phase** (ODH Dashboard UI)
   - Navigate to ODH Dashboard → "Training Jobs" → "Create MPIJob"
   - Select template: "LLM Fine-Tuning (DeepSpeed + MPI)"
   - Configure: Name="llm-finetune-v3", Workers=8, GPUs per worker=4, Image="registry.redhat.io/llm-training:latest"
   - Advanced settings: Mount PVC "financial-docs-data", set memory=320Gi per worker, enable fp16 mixed precision
   - Click "Create Job" → Job submitted

2. **Monitoring Phase** (ODH Dashboard + Slack Notifications)
   - Dashboard shows job transitioning: Pending (10s) → Running (workers starting) → Training
   - Overview tab displays: 8/8 workers running, GPU utilization 92% avg, estimated time remaining: 2h 15m
   - Metrics tab shows: loss curve decreasing, throughput=45k tokens/sec (vs. 2k tokens/sec single-node)
   - Slack integration sends notification when job completes or fails

3. **Validation Phase** (SDK in Jupyter Notebook)
   ```python
   from odh_sdk import MPIJob
   job = MPIJob.get("llm-finetune-v3")
   job.wait_for_completion()  # Blocks until job done
   print(f"Training completed in {job.duration}")  # "2h 18m"
   # Load checkpoint and run evaluation
   ```

**Outcome**: Training time reduced from 18 days → 2.3 hours (190x speedup). Team can now iterate on model architecture and hyperparameters multiple times per day. Model achieves 15% higher F1 score on financial entity extraction due to increased batch size and training stability.

**Success Metrics**: Training time ≤3 hours, ≥85% scaling efficiency, successful completion rate >95%, zero data egress to public cloud.

---

**Use Case 2: Continuous Training Pipeline for Autonomous Vehicle Perception**

**Persona**: MLOps Engineer at Automotive AI Company

**Goal**: Build a CI/CD pipeline that automatically retrains object detection models nightly using the latest 500GB of dashcam footage data collected from test vehicles.

**Current Pain Point**: Manual job submission via kubectl YAML files is error-prone and lacks visibility. Training jobs frequently fail overnight due to resource contention with other teams, but failures aren't detected until morning. No programmatic way to track model lineage from training job to deployed model.

**Workflow with MPIJobs**:

1. **Pipeline Definition** (Python SDK in Airflow DAG)
   ```python
   from airflow import DAG
   from odh_sdk import MPIJob, ResourceRequirements

   def submit_training_job(**context):
       job = MPIJob(
           name=f"perception-train-{context['ds_nodash']}",
           workers=16,
           image="quay.io/acme/yolo-training:v2.1",
           command=["python", "train.py", "--data-path", "/mnt/data/nightly"],
           resources=ResourceRequirements(gpu=2, memory="64Gi", cpu=8),
           volume_mounts=[{"name": "training-data", "mount_path": "/mnt/data"}]
       )
       job.create()
       return job.name

   def monitor_training_job(**context):
       job_name = context['ti'].xcom_pull(task_ids='submit_job')
       job = MPIJob.get(job_name)
       job.wait_for_completion(timeout=7200)  # 2 hour timeout
       if job.status == "Failed":
           raise Exception(f"Training failed: {job.failure_reason}")
       # Register model in registry
       register_model(job.get_checkpoint_path())
   ```

2. **Execution and Monitoring** (Airflow UI + ODH Dashboard)
   - Airflow triggers job submission at 2 AM
   - MPIJob created with automatic resource allocation from shared GPU pool
   - Kueue queues job if GPUs unavailable, provisions when resources free up
   - ODH Dashboard shows job progress; Airflow polls status via SDK
   - Upon completion, Airflow task extracts model checkpoint and registers in Model Registry

3. **Debugging and Iteration** (CLI)
   - If job fails, MLOps engineer receives PagerDuty alert
   - Uses CLI to investigate: `odh mpijob logs perception-train-20251028 --tail 100 --worker 8`
   - Identifies issue: Worker 8 crashed due to corrupt data file
   - Fixes data validation in pre-processing pipeline
   - Re-triggers Airflow DAG to retry training

**Outcome**: Zero-touch nightly training with automatic failure detection and alerting. Model training success rate improved from 60% (manual) to 97% (automated with SDK). Time from model training to production deployment reduced from 3 days (manual) to 4 hours (automated). Model Registry contains full lineage from training job to deployed model version.

**Success Metrics**: ≥95% training job success rate, ≤5 minutes from job failure to alert, zero manual job submissions, 100% model lineage traceability.

---

**Use Case 3: Multi-Node Computer Vision Training for Medical Imaging**

**Persona**: Research Scientist at Medical AI Startup

**Goal**: Train a 3D U-Net model for brain tumor segmentation using 10,000 high-resolution MRI scans (2TB dataset) to meet FDA validation timelines.

**Current Pain Point**: Limited to 2-node training with manual Kubernetes YAML files. No visibility into per-worker GPU utilization, leading to frequent OOM failures that require trial-and-error memory tuning. HIPAA compliance requires data to stay on-premises, but current tooling lacks audit logging for training jobs.

**Workflow with MPIJobs**:

1. **Job Configuration** (ODH Dashboard UI with Compliance Settings)
   - Create MPIJob: Name="tumor-segmentation-3dunet", Workers=12, GPUs per worker=1
   - Select compliant image from internal registry (HIPAA-approved base image)
   - Configure: Mount PVC "mri-dataset-encrypted" (encrypted at rest), set CPU=16, Memory=128Gi per worker
   - Enable audit logging: All job events logged to SIEM system for compliance
   - Set node affinity: Schedule on-prem GPU nodes only (not cloud burst nodes)

2. **Training Execution** (SDK in Jupyter + Dashboard Monitoring)
   ```python
   from odh_sdk import MPIJob

   job = MPIJob(
       name="tumor-segmentation-3dunet",
       workers=12,
       image="internal-registry/3dunet-horovod:v1.2-hipaa",
       command=["horovodrun", "-np", "12", "python", "train_3dunet.py"],
       resources={"gpu": 1, "memory": "128Gi", "cpu": 16},
       volume_mounts=[{"name": "mri-data", "mount_path": "/data"}],
       node_selector={"node-type": "on-prem-gpu"}
   )
   job.create()
   print(f"Job created: {job.name}")
   ```

   - Dashboard shows: 12/12 workers running on on-prem nodes, GPU utilization 88% avg
   - Logs tab displays synchronized output from all workers with MPI reduction logs
   - Metrics tab: Network bandwidth between workers = 8 Gbps (confirming MPI communication healthy)

3. **Compliance and Audit** (Admin Dashboard)
   - Administrator views audit logs: "MPIJob 'tumor-segmentation-3dunet' created by user 'scientist-jane' at 2025-10-29T14:32:01Z"
   - All worker pods scheduled on compliant on-prem nodes (verified via node labels)
   - Data never left on-premises boundary (confirmed via network policy logs)
   - Model checkpoint saved to encrypted PVC, access logged to SIEM

**Outcome**: Training time reduced from 6 days (2-node) → 14 hours (12-node), enabling team to meet FDA submission deadline. HIPAA compliance maintained with full audit trail. OOM failures eliminated via Dashboard metrics guiding memory tuning (128Gi per worker). Model achieved 92% Dice score on validation set.

**Success Metrics**: Training time ≤16 hours, 100% data residency on-prem, complete audit trail for compliance, zero OOM failures after initial tuning.

---

**Use Case 4: Hyperparameter Optimization at Scale for Recommendation System**

**Persona**: Senior MLOps Engineer at E-Commerce Platform

**Goal**: Optimize hyperparameters (learning rate, embedding dimensions, regularization) for a deep learning recommendation model serving 10M+ users. Need to evaluate 200+ hyperparameter combinations to maximize click-through rate.

**Current Pain Point**: Single-node training takes 4 hours per hyperparameter trial. Running 200 trials sequentially would take 800 hours (33 days). No way to parallelize trials across multiple nodes efficiently. Manual tracking of trial results in spreadsheets is error-prone.

**Workflow with MPIJobs** (Future: Katib Integration):

1. **HPO Experiment Definition** (Python Script with Parallel Job Submission)
   ```python
   from odh_sdk import MPIJob
   import itertools

   # Define hyperparameter search space
   learning_rates = [1e-4, 5e-4, 1e-3, 5e-3]
   embedding_dims = [64, 128, 256, 512]
   l2_regs = [1e-5, 1e-4, 1e-3]

   # Generate all combinations
   trials = list(itertools.product(learning_rates, embedding_dims, l2_regs))

   # Submit MPIJob for each trial (parallel execution)
   jobs = []
   for i, (lr, emb_dim, l2) in enumerate(trials):
       job = MPIJob(
           name=f"rec-hpo-trial-{i}",
           workers=4,  # Each trial uses 4 workers
           image="quay.io/ecommerce/rec-model:latest",
           command=["python", "train.py", f"--lr={lr}", f"--emb-dim={emb_dim}", f"--l2={l2}"],
           resources={"gpu": 1, "memory": "32Gi"}
       )
       job.create()
       jobs.append(job)

   print(f"Submitted {len(jobs)} HPO trials")

   # Monitor all trials
   for job in jobs:
       job.wait_for_completion()
       metric = job.get_metric("val_ctr")  # Get validation CTR from logs
       print(f"{job.name}: lr={lr}, emb_dim={emb_dim}, l2={l2} → val_ctr={metric}")
   ```

2. **Parallel Execution** (Dashboard Monitoring)
   - Kueue schedules trials based on resource availability (e.g., 10 trials running concurrently with 40 total workers)
   - Dashboard shows: 10 MPIJobs "Running", 190 "Queued", 0 "Completed"
   - As GPU nodes free up, queued jobs automatically start
   - Estimated completion time: 80 hours (10× speedup via parallelization)

3. **Results Analysis** (SDK in Notebook)
   ```python
   # Collect results from all completed jobs
   results = []
   for i in range(len(trials)):
       job = MPIJob.get(f"rec-hpo-trial-{i}")
       if job.status == "Succeeded":
           results.append({
               "trial": i,
               "lr": trials[i][0],
               "emb_dim": trials[i][1],
               "l2": trials[i][2],
               "val_ctr": job.get_metric("val_ctr"),
               "duration": job.duration
           })

   # Find best hyperparameters
   best = max(results, key=lambda x: x["val_ctr"])
   print(f"Best config: lr={best['lr']}, emb_dim={best['emb_dim']}, l2={best['l2']} → val_ctr={best['val_ctr']}")
   ```

**Outcome**: HPO experiment completed in 80 hours (vs. 800 hours sequential), enabling team to iterate on model architecture weekly instead of monthly. Best hyperparameters improved production CTR by 8%, translating to $2M+ annual revenue increase. Programmatic result tracking eliminated manual spreadsheet errors.

**Success Metrics**: HPO experiment completion ≤5 days, ≥10× speedup via parallelization, zero manual result tracking, production metric improvement ≥5%.

---

**Use Case 5: Shared Cluster Resource Management for Multi-Team AI Organization**

**Persona**: Platform Administrator at AI/ML SaaS Company

**Goal**: Manage a shared OpenShift AI cluster with 128 GPUs serving 5 teams (Computer Vision, NLP, Forecasting, Recommender, Research). Ensure fair resource allocation, prevent resource hogging, and maintain high GPU utilization.

**Current Pain Point**: Teams manually coordinate GPU access via Slack, leading to conflicts and idle GPUs. No visibility into cluster-wide resource usage. High-priority production retraining jobs are delayed by low-priority research experiments. No way to enforce team quotas.

**Workflow with MPIJobs** (Admin Perspective):

1. **Resource Policy Configuration** (OpenShift CLI + Kueue)
   ```bash
   # Create ResourceQuota per team namespace
   oc create quota team-cv-quota --hard=requests.nvidia.com/gpu=32 -n team-cv
   oc create quota team-nlp-quota --hard=requests.nvidia.com/gpu=32 -n team-nlp
   oc create quota team-forecasting-quota --hard=requests.nvidia.com/gpu=24 -n team-forecasting
   oc create quota team-recommender-quota --hard=requests.nvidia.com/gpu=24 -n team-recommender
   oc create quota team-research-quota --hard=requests.nvidia.com/gpu=16 -n team-research

   # Configure Kueue ClusterQueue for fair-share scheduling
   kubectl apply -f - <<EOF
   apiVersion: kueue.x-k8s.io/v1beta1
   kind: ClusterQueue
   metadata:
     name: gpu-cluster-queue
   spec:
     resources:
       - name: nvidia.com/gpu
         total: 128
         flavors:
           - name: a100-80gb
             resources:
               - name: nvidia.com/gpu
                 total: 64
           - name: a100-40gb
             resources:
               - name: nvidia.com/gpu
                 total: 64
     fairSharing:
       preemption: true
       weights:
         - name: team-cv
           weight: 32
         - name: team-nlp
           weight: 32
         - name: team-research
           weight: 16
   EOF
   ```

2. **Cluster Monitoring** (ODH Dashboard Admin View)
   - Navigate to ODH Dashboard → "Admin" → "Training Jobs" (cluster-wide view)
   - Dashboard displays:
     - Active MPIJobs: 8 running, 4 queued
     - GPU Utilization: 118/128 GPUs allocated (92% utilization)
     - Per-Team Usage: CV=30 GPUs, NLP=28 GPUs, Forecasting=24 GPUs, Recommender=20 GPUs, Research=16 GPUs
     - Queued Jobs: "team-cv/object-detection-large" (waiting for 8 GPUs), "team-research/llm-experiment" (waiting for 16 GPUs)
   - Alerts: "team-cv approaching quota limit (30/32 GPUs used)"

3. **Dynamic Resource Allocation** (Automatic via Kueue)
   - team-cv submits MPIJob requesting 8 GPUs → Exceeds quota → Queued by Kueue
   - team-recommender's job completes, freeing 8 GPUs
   - Kueue automatically allocates freed GPUs to team-cv's queued job (fair-share policy)
   - Dashboard shows: team-cv job transitions from "Queued" → "Running"

4. **Priority Handling** (High-Priority Production Retraining)
   ```bash
   # team-nlp submits high-priority job with preemption
   odh mpijob create production-retrain \
     --workers 16 \
     --gpu 1 \
     --image registry/nlp-model:v3 \
     --priority high \
     --preempt true
   ```
   - Kueue preempts 16 GPUs from low-priority research jobs
   - Research jobs checkpointed and paused, production job starts immediately
   - Dashboard shows: 2 research jobs "Preempted", 1 production job "Running"
   - After production job completes, research jobs automatically resume from checkpoint

**Outcome**: GPU utilization increased from 62% (manual coordination) to 92% (automated scheduling). Average job queue time reduced from 6 hours to 15 minutes. Zero Slack coordination needed, eliminating team conflicts. High-priority production jobs meet SLAs with <5 minute start time via preemption.

**Success Metrics**: GPU utilization ≥90%, average queue time ≤30 minutes, zero manual resource coordination, 100% SLA compliance for high-priority jobs.

---

## Documentation Considerations

**Documentation Scope and Audience:**

The MPIJobs feature requires comprehensive documentation targeting three distinct audiences with varying technical expertise and information needs:

1. **Data Scientists**: Focus on practical "how-to" guides with minimal Kubernetes/MPI jargon. Emphasize SDK and Dashboard UI workflows. Provide runnable code examples and notebooks.

2. **MLOps Engineers**: Focus on integration patterns, CI/CD examples, programmatic APIs (SDK, CLI), and troubleshooting. Include architecture diagrams and API reference documentation.

3. **Platform Administrators**: Focus on installation, configuration, resource management, security, monitoring, and multi-tenancy. Include YAML configuration examples and operational runbooks.

**Documentation Structure:**

**1. Getting Started Guide** (Target: Data Scientists, 15-20 pages)
- **Prerequisites**: RedHat OpenShift AI version requirements, KubeFlow Trainer V2 dependency check, namespace setup
- **Quickstart Tutorial**: End-to-end example training a PyTorch model with Horovod on 4 workers in <30 minutes
  - Step 1: Prepare training script and container image
  - Step 2: Submit MPIJob via Dashboard UI (with screenshots)
  - Step 3: Monitor job progress and view logs
  - Step 4: Retrieve trained model from PVC
- **Common Patterns**: Template configurations for popular frameworks (PyTorch + Horovod, TensorFlow + MPI, DeepSpeed, Megatron-LM)
- **Jupyter Notebook Examples**: Downloadable notebooks demonstrating SDK usage

**2. User Guide** (Target: Data Scientists + MLOps Engineers, 40-50 pages)

*Chapter 1: Concepts*
- What is MPI and why use it for distributed training?
- MPIJob architecture: Launcher, Workers, MPI hostfile, SSH communication
- How MPIJobs integrate with KubeFlow Trainer V2
- Resource model: Workers, slots, GPUs, CPUs, memory

*Chapter 2: Creating MPIJobs*
- Dashboard UI: Job creation wizard, advanced YAML editor
- CLI: Imperative commands, declarative YAML files, command reference
- SDK: Python API, type hints, async patterns, error handling
- Configuration reference: All supported parameters (workers, resources, volumes, env vars, etc.)

*Chapter 3: Monitoring and Observability*
- Job lifecycle states (Pending, Running, Succeeded, Failed)
- Dashboard UI: Job list, details view (6 tabs), filtering/sorting
- Logs: Viewing launcher and worker logs, log aggregation, keyword search
- Metrics: GPU utilization, memory usage, network bandwidth, custom metrics

*Chapter 4: Advanced Topics*
- Multi-node configurations (2-node, 8-node, 16-node, 32+ node)
- GPU optimization: NCCL, RDMA (future), topology awareness
- Storage patterns: PVCs for datasets and checkpoints, S3/object storage
- Framework-specific guides: Horovod, DeepSpeed, PyTorch DDP with MPI backend, TensorFlow MPI
- Performance tuning: Batch size selection, gradient accumulation, mixed precision (fp16/bf16)

*Chapter 5: Integration Patterns*
- CI/CD pipelines: Airflow, Tekton, Jenkins examples
- Programmatic workflows: SDK in Python scripts, async job submission
- Model lifecycle: Training → Checkpoint → Model Registry → Deployment
- HPO integration: Using MPIJobs with Katib (future)

**3. Administrator Guide** (Target: Platform Administrators, 30-40 pages)

*Chapter 1: Installation and Configuration*
- Prerequisites: OpenShift AI operator, KubeFlow Trainer V2 installation
- Enabling MPIJobs feature: Operator configuration, CRD verification
- Network requirements: Pod-to-pod communication, NetworkPolicies
- Storage provisioning: Recommendations for shared PVCs (NFS, CephFS, etc.)

*Chapter 2: Resource Management*
- ResourceQuotas and LimitRanges: Per-namespace GPU/CPU/memory limits
- Kueue integration: Fair-share scheduling, priority classes, preemption
- Node affinity: Targeting GPU nodes, taints/tolerations
- Capacity planning: Estimating cluster size for concurrent MPIJobs

*Chapter 3: Security and Multi-Tenancy*
- Namespace isolation: RBAC policies, ServiceAccounts
- NetworkPolicies: Restricting inter-pod communication
- SSH key management: Ephemeral keys, rotation policies
- Image security: Trusted registries, vulnerability scanning
- Audit logging: Capturing MPIJob lifecycle events

*Chapter 4: Monitoring and Operations*
- Prometheus metrics: Exposed metrics for MPIJobs, alerting rules
- Grafana dashboards: Pre-built dashboards for cluster-wide view
- Logging: Integration with OpenShift logging stack (EFK/Loki)
- Health checks: Monitoring MPI Operator, KubeFlow Trainer components

*Chapter 5: Troubleshooting*
- Common issues: SSH failures, MPI version mismatches, GPU allocation failures, OOM errors
- Diagnostic commands: CLI and kubectl commands for debugging
- Log analysis: Identifying root causes from launcher/worker logs
- Performance issues: Slow training, poor scaling efficiency

*Chapter 6: Upgrade and Migration*
- Upgrading KubeFlow Trainer V2: Version compatibility matrix
- Migrating from standalone MPI Operator (if applicable)
- Breaking changes and deprecations
- Rollback procedures

**4. API Reference** (Target: MLOps Engineers, Auto-Generated)
- CLI command reference: All commands, flags, examples
- Python SDK API reference: Classes, methods, parameters, return types (auto-generated from docstrings)
- MPIJob CRD schema: YAML field reference with descriptions and examples
- REST API (if exposed): Endpoints, request/response formats

**5. Troubleshooting Guide** (Target: All Audiences, 20-25 pages)

*Organized by error category with symptoms → diagnosis → resolution:*

- **Job submission errors**: Quota exceeded, invalid configuration, image pull failures
- **SSH connectivity errors**: Hostfile issues, SSH key permissions, NetworkPolicy blocking
- **MPI errors**: Version mismatches, library not found, MPI_Init failures
- **Resource errors**: OOMKilled, GPU allocation failures, insufficient CPU/memory
- **Performance issues**: Slow training, poor scaling efficiency (<70%), network bottlenecks
- **Platform errors**: Node failures, storage issues, API server timeouts

*Each error includes:*
- Symptoms (error messages, observed behavior)
- Root cause explanation
- Step-by-step resolution
- Prevention tips

**6. Release Notes and Migration Guides**
- What's new in each release
- Breaking changes and deprecations
- Migration guides for users upgrading from previous versions
- Known issues and workarounds

**Documentation Delivery and Maintenance:**

- **Format**: HTML (hosted on docs.redhat.com), PDF (downloadable), man pages (for CLI)
- **Versioning**: Documentation versioned alongside OpenShift AI releases (e.g., docs for RHOAI 2.15, 2.16, etc.)
- **Search and Navigation**: Full-text search, breadcrumb navigation, cross-references between sections
- **Code Examples**: All code examples tested in CI/CD pipeline before documentation release (no stale examples)
- **Localization**: Initial release in English, evaluate demand for localization (Japanese, Chinese, German) in future
- **Feedback Mechanism**: "Was this page helpful?" feedback widget on every page, direct links to file documentation issues

**Existing Documentation Integration:**

Since this feature extends existing OpenShift AI functionality, documentation must integrate with and reference:

- **Existing KubeFlow Trainer V2 Documentation**: Link to upstream KubeFlow docs for MPI Operator concepts, explain RedHat-specific enhancements
- **OpenShift AI Dashboard Documentation**: Update existing Dashboard guide to include MPIJobs section, maintain consistent screenshot style
- **OpenShift AI SDK Documentation**: Extend existing SDK docs with MPIJob class reference, maintain consistent API patterns
- **GPU Management Documentation**: Cross-reference existing GPU setup guides (NVIDIA GPU Operator, node configuration)
- **Storage Documentation**: Link to existing PVC provisioning guides, explain MPIJob-specific storage patterns

**Documentation Examples Repository:**

Create a separate Git repository (`openshift-ai-mpijobs-examples`) with runnable examples:

- `examples/quickstart/`: Basic PyTorch + Horovod example (2 workers, CPU-only)
- `examples/pytorch-horovod-gpu/`: PyTorch + Horovod with 4 GPUs
- `examples/tensorflow-mpi/`: TensorFlow with MPI Allreduce
- `examples/deepspeed/`: DeepSpeed ZeRO optimization
- `examples/llm-finetuning/`: LLM fine-tuning with 8+ nodes
- `examples/hpo-parallel/`: Parallel hyperparameter optimization
- `examples/ci-cd-airflow/`: Airflow DAG for automated training
- `examples/ci-cd-tekton/`: Tekton pipeline for GitOps workflow

Each example includes:
- `README.md`: Step-by-step instructions
- `Dockerfile`: Container image definition
- `train.py` (or similar): Training script
- `mpijob.yaml`: MPIJob resource definition
- `notebook.ipynb`: Jupyter notebook demonstrating SDK usage

**Success Metrics for Documentation:**

- Time-to-first-successful-job: Users can run first MPIJob within 30 minutes of reading quickstart
- Support ticket reduction: 50% reduction in MPIJob-related support tickets after documentation release
- Documentation feedback: >80% "helpful" rating on feedback widgets
- Search effectiveness: Users find answers within top 3 search results 90% of the time
- Example adoption: 70% of users start from provided examples (measured via survey)

---

## Questions to Answer

> **Note**: Many of the architectural questions below have been addressed in the companion **ARCHITECTURE.md** document, which provides detailed technical design guidance from the Architecture team. Key architectural decisions are documented there as Architecture Decision Records (ADRs).

**Technical Architecture and Design Questions:**

**Q1. Integration and Deployment Architecture** ✅ *Addressed in ARCHITECTURE.md Section 1.2*
- How will the KubeFlow MPI Operator v2 be packaged and deployed as part of the RedHat OpenShift AI operator?
  - **Decision**: Sub-controller pattern - MPI Operator deployed as managed component within RHOAI operator (see ADR-001)
  - Upgrade path: Independent versioning with compatibility matrix maintained
- What is the minimum OpenShift version required to support MPIJobs?
  - **Recommendation**: OpenShift 4.14+ (for stable GPU Operator support)
  - RDMA support (post-MVP) requires OpenShift 4.16+ with Multus CNI
- Version compatibility between KubeFlow Trainer V2, MPI Operator, and RHOAI:
  - **Yes** - Compatibility matrix documented in ARCHITECTURE.md Section 8.2
  - Automated testing in CI/CD for certified version combinations

**Q2. Networking and Inter-Process Communication** ✅ *Addressed in ARCHITECTURE.md Section 3*
- What networking configuration is required for MPI inter-worker communication?
  - **MVP**: OpenShift SDN/OVN-Kubernetes is sufficient for typical workloads (10Gbps+ network recommended)
  - **Post-MVP**: RDMA (RoCE/InfiniBand) for high-performance scenarios (see Section 3.4)
  - **NetworkPolicy templates**: Provided in ARCHITECTURE.md Section 6.4 for multi-tenant isolation
- How are SSH keys managed for MPI launcher-worker communication?
  - **Decision**: Ephemeral SSH keys auto-generated per job, stored in job-scoped Secrets (see ADR-003)
  - Automatic cleanup: Keys deleted when MPIJob is deleted (owner references)
  - **Alternative considered**: PMIx - deferred to post-MVP due to container runtime complexity
- Expected network bandwidth requirements:
  - **Minimum**: 10 Gbps for 8-node, 32-GPU training (typical DL workload)
  - **Recommended**: 25-100 Gbps for large-scale (16+ nodes) or parameter-heavy models
  - **Guidance**: Document network topology best practices in Admin Guide (same rack/AZ preferred)

**Q3. Resource Orchestration and Scheduling** ✅ *Addressed in ARCHITECTURE.md Section 4*
- Gang scheduling for atomic worker allocation:
  - **Decision**: Kueue integration required for MVP to prevent resource deadlock (see ADR-004)
  - MPI Operator alone does NOT provide gang scheduling - workers can be partially allocated
  - **Risk mitigation**: Document Kueue requirement in installation prerequisites
- Preventing resource deadlock with competing MPIJobs:
  - **Solution**: Kueue's ClusterQueue with fair-share weights and admission control
  - Example configuration in ARCHITECTURE.md Section 4.2 prevents the 40+40 GPU deadlock scenario
- PriorityClass configuration:
  - **Recommendation**: Three default PriorityClasses provided: `training-high` (1000), `training-normal` (500), `training-low` (100)
  - **Preemption**: Post-MVP feature - requires checkpoint/restart capability (see Section 9.4)
  - MVP: Preemption triggers job failure, user must resubmit

**Q4. Hardware Acceleration and Runtime**
- What GPU types and drivers are supported?
  - NVIDIA GPUs only, or AMD/Intel GPUs as well?
  - Does this require NVIDIA GPU Operator, or is it optional?
  - What CUDA/cuDNN/NCCL versions are recommended for optimal performance?
- How do we handle GPU topology awareness?
  - For multi-GPU nodes, should we expose GPU affinity (e.g., NVLink topology) to MPI processes?
  - Does NCCL auto-detect optimal communication paths, or do we need to configure `NCCL_TOPO_FILE`?
- What container runtime features are required?
  - Do MPIJob containers require privileged mode or specific capabilities (e.g., `IPC_LOCK` for RDMA)?
  - Can MPIJobs run with non-root user and restrictive security contexts (for enterprise security requirements)?

**Q5. Required Software and Data Services**
- What container image requirements should be documented for MPIJob compatibility?
  - Must images include specific MPI implementations (Open MPI, Intel MPI, MPICH)?
  - What MPI version compatibility matrix should we provide?
  - Should RedHat provide base images with MPI + popular ML frameworks pre-installed?
- How should training data and model checkpoints be stored?
  - Recommended PVC types: NFS, CephFS, Ceph RBD, Hostpath, CSI drivers?
  - For large datasets (1TB+), should we recommend S3-compatible object storage with CSI mountpoints?
  - What is the guidance on ReadWriteMany vs. ReadWriteOnce PVCs?
- How do MPIJobs integrate with Model Registry?
  - Should there be automatic model registration upon job completion (post-MVP)?
  - How is lineage tracked from MPIJob → trained model → deployed inference service?

**Q6. Operational and Non-Functional Requirements**
- What are the performance SLAs for MPIJob operations?
  - Job submission latency: <5 seconds from CLI/SDK to job created
  - Worker pod startup time: <30 seconds (excluding image pull)
  - Dashboard UI page load time: <2 seconds for job list (100 jobs)
  - Log streaming latency: <5 seconds from pod log generation to UI display
- What are the scalability limits for the platform?
  - Maximum concurrent MPIJobs per cluster: 100? 500?
  - Maximum workers per MPIJob: 100? 1000?
  - Maximum total GPU count per cluster for MPIJobs: 512? 1024?
- What is the disaster recovery and fault tolerance strategy?
  - If a worker pod fails mid-training, does the entire MPIJob fail, or can it continue with remaining workers (elastic training)?
  - Should MPIJobs support automatic checkpointing to PVCs for fault tolerance (post-MVP)?
  - If the MPI Operator pod crashes, what happens to running MPIJobs?

**Q7. Monitoring, Logging, and Debugging**
- What metrics should be exposed to Prometheus?
  - Job-level metrics: Total jobs created, running, succeeded, failed, queued
  - Worker-level metrics: GPU utilization per worker, memory usage, CPU usage, network bandwidth
  - Performance metrics: Scaling efficiency, training throughput (samples/sec, tokens/sec)
- What Grafana dashboards should be provided out-of-the-box?
  - Cluster-wide dashboard: GPU utilization, active MPIJobs, queue depth
  - Per-job dashboard: Worker pod status, resource usage time series, loss curves (if custom metrics exposed)
- How should logs be aggregated and retained?
  - Should logs from all worker pods be automatically aggregated and stored centrally (e.g., EFK stack)?
  - What is the log retention policy (30 days? 90 days?)?
  - How do users access logs of completed/deleted jobs?

**Q8. Security and Access Control**
- What RBAC roles and permissions are required for different personas?
  - Data Scientist role: Create/view/delete MPIJobs in assigned namespaces, view logs
  - MLOps Engineer role: Same as data scientist + read access to cluster-wide metrics
  - Administrator role: Manage MPI Operator, view all MPIJobs across namespaces, configure resource quotas
- How is multi-tenancy enforced?
  - Can users in namespace A view or interfere with MPIJobs in namespace B?
  - Are NetworkPolicies automatically applied to isolate MPI traffic between namespaces?
- What security scans and certifications are required?
  - Must container images pass specific vulnerability scans (e.g., Red Hat Container Certification)?
  - Are there compliance requirements (FedRAMP, PCI-DSS, HIPAA) that affect MPIJob design?
- How are secrets managed for training jobs?
  - If training scripts need to access S3 credentials, database passwords, or API keys, how are these injected securely?
  - Should we integrate with OpenShift's secret management (Sealed Secrets, External Secrets Operator)?

**Q9. User Experience and Interface Design**
- What is the exact UI/UX design for the ODH Dashboard MPIJob integration?
  - Do we have wireframes or mockups for the job creation wizard, list view, and details page?
  - How do we ensure visual consistency with existing KubeFlow Trainer job types (TFJob, PyTorchJob)?
  - Are there specific accessibility requirements (WCAG 2.1 compliance)?
- What is the SDK API surface and naming convention?
  - Should the SDK use `odh_sdk.MPIJob` or `openshift_ai_sdk.training.MPIJob`?
  - How do we maintain API stability across RHOAI versions?
  - Are there breaking changes between KubeFlow Trainer V2 upstream and RHOAI's SDK?
- What CLI naming convention should be used?
  - Should MPIJob commands be under `oc` plugin (e.g., `oc odh mpijob create`), a dedicated CLI (`odh mpijob create`), or `kubectl` plugin (`kubectl mpijob create`)?
  - How do we ensure consistency with existing RHOAI CLI patterns?

**Q10. Testing and Validation Strategy**
- What is the testing strategy for MPIJob functionality?
  - Unit tests: SDK, CLI, Dashboard backend API
  - Integration tests: End-to-end job submission, monitoring, deletion
  - Performance tests: Scaling efficiency, resource utilization, latency benchmarks
  - Security tests: RBAC enforcement, network isolation, secret handling
- What test environments are required?
  - Do we need dedicated GPU clusters for CI/CD testing?
  - What is the minimum GPU node count for automated tests (2 nodes? 4 nodes?)?
- What are the acceptance criteria for MVP release?
  - Must all 10 Acceptance Criteria in the RFE be validated with automated tests?
  - What is the minimum passing rate for integration tests (95%? 98%?)?

**Q11. Upstream Contribution and Maintenance**
- What is the relationship between RedHat's MPIJob implementation and upstream KubeFlow?
  - Are we contributing enhancements back to KubeFlow Trainer V2 upstream?
  - Are there RedHat-specific patches that diverge from upstream?
- How do we handle upstream breaking changes?
  - If KubeFlow Trainer V3 introduces breaking API changes, what is our migration strategy?
  - Do we maintain long-term support (LTS) branches?

**Product and Go-to-Market Questions:**

**Q12. Feature Prioritization and Roadmap**
- What is the MVP feature set vs. post-MVP enhancements?
  - Is gang scheduling required for MVP, or can it be post-MVP?
  - Is Katib integration (HPO) required for MVP, or post-MVP?
  - Is RDMA support required for MVP, or post-MVP?
- What is the release timeline?
  - Target RHOAI version for MVP: 2.16? 2.17?
  - Private beta timeframe: 3 months before GA?
  - GA release date: Q1 2026? Q2 2026?

**Q13. Customer Validation and Beta Program**
- Which customers should be invited to private beta?
  - Prioritize strategic enterprise accounts ($500K+ ARR) with on-prem GPU infrastructure
  - Include at least one customer from each key vertical (financial services, healthcare, automotive)
- What are the beta exit criteria?
  - Minimum number of successful customer deployments: 5? 10?
  - Minimum number of MPIJobs executed in production: 100? 500?
  - Customer satisfaction score: NPS >50?

**Q14. Competitive Differentiation**
- How do we position MPIJobs vs. competitive offerings (AWS SageMaker, Azure ML, Google Vertex AI)?
  - What are the unique value propositions that only RHOAI + MPIJobs can deliver?
  - Should we publish benchmark comparisons (e.g., RHOAI MPIJob vs. SageMaker Distributed Training)?
- How do we message hybrid cloud advantages?
  - What customer case studies can we develop to demonstrate on-prem + cloud flexibility?

**Q15. Pricing and Packaging**
- Is MPIJob support included in base RHOAI license, or a premium add-on?
- Are there usage-based pricing implications (e.g., per-GPU-hour metering)?
- How do we handle bring-your-own-license (BYOL) scenarios for underlying MPI implementations?

---

## Background & Strategic Fit

**Market Context and Opportunity:**

The enterprise AI/ML landscape is undergoing a fundamental shift driven by three converging trends:

1. **LLM and Foundation Model Adoption**: Organizations are increasingly fine-tuning large language models (7B-70B+ parameters) and training custom foundation models for domain-specific applications (legal, medical, financial, industrial). These models are too large to train on single nodes, creating demand for efficient multi-node distributed training.

2. **Hybrid Cloud Imperative**: 73% of enterprises operate hybrid cloud environments due to data sovereignty requirements, regulatory compliance (GDPR, HIPAA, FedRAMP), and existing on-premises infrastructure investments. However, most cloud-based ML platforms (AWS SageMaker, Azure ML) do not extend their distributed training capabilities to on-premises or edge environments, forcing customers to choose between performance and compliance.

3. **GPU Infrastructure Stranded Value**: Enterprises have invested heavily in on-premises GPU infrastructure (NVIDIA A100, H100, AMD MI250) for AI/ML workloads. Without native distributed training support in their ML platforms, these GPUs are underutilized (40-50% utilization) because data scientists are limited to single-node training, leaving $2M+ per customer in stranded infrastructure value.

**The Business Problem:**

RedHat OpenShift AI currently lacks native support for MPI-based distributed training jobs, creating three critical business impacts:

1. **Revenue Gap**: We lose 10-15 enterprise deals annually ($5-8M in ARR) to competitors (AWS SageMaker, Azure ML, Databricks) that offer integrated distributed training. Sales feedback indicates distributed training is a top-3 requirement for 60% of enterprise AI/ML buyer conversations.

2. **Customer Churn Risk**: Existing customers with growing model sizes (e.g., migrating from BERT-base to GPT-class models) are forced to run distributed training workloads outside OpenShift AI, fragmenting their MLOps toolchain and increasing the risk of platform abandonment. Customer surveys indicate 40% of RHOAI users have evaluated competitive platforms specifically for distributed training capabilities.

3. **Competitive Disadvantage**: Every major cloud ML platform (AWS SageMaker, Azure ML, Google Vertex AI, Databricks) and emerging Kubernetes-native ML platforms (Run:ai, Domino, Determined AI) support MPI-based distributed training. OpenShift AI is the only enterprise ML platform in the market without this capability, ceding the high-growth LLM fine-tuning market to competitors.

**Strategic Fit with RedHat Portfolio:**

Adding MPIJobs support to OpenShift AI via KubeFlow Trainer V2 aligns with three core RedHat strategic pillars:

**1. Hybrid Cloud Leadership**
- RedHat is uniquely positioned to deliver distributed training across hybrid cloud (on-prem, edge, any public cloud) with OpenShift's consistent Kubernetes platform
- MPIJobs enable customers to leverage existing on-prem GPU investments while maintaining optionality for cloud bursting, avoiding cloud vendor lock-in
- This differentiates RHOAI from AWS/Azure/GCP, which are cloud-only platforms with proprietary APIs

**2. Open Source Foundation and Community Alignment**
- KubeFlow Trainer V2 is the industry-standard open-source framework for ML training on Kubernetes, with contributions from Google, AWS, Microsoft, NVIDIA, and 200+ organizations
- Integrating MPI Operator (part of KubeFlow) demonstrates RedHat's commitment to upstream-first development and avoids proprietary lock-in
- RedHat can contribute enhancements back to KubeFlow community, building ecosystem leadership

**3. Enterprise AI Democratization**
- Simplifying distributed training through CLI, SDK, and Dashboard UI removes barriers for data scientists without deep MPI/Kubernetes expertise
- This aligns with OpenShift AI's mission to make enterprise AI accessible to all skill levels, not just ML infrastructure specialists
- Unified observability across all training job types (MPIJob, TFJob, PyTorchJob) reduces operational complexity and accelerates time-to-value

**Customer Validation and Demand Signals:**

Demand for MPIJobs support is validated through multiple channels:

**Direct Customer Requests** (Anonymized):
- **Fortune 500 Financial Services**: *"We have 200 A100 GPUs on-premises but our data scientists can only use them for single-node training due to compliance requirements. We need distributed training in RHOAI to unlock $2M of stranded infrastructure value, or we'll be forced to migrate to SageMaker and export anonymized data to AWS."*

- **Medical AI Startup**: *"Training our 3D medical imaging models on-prem takes 6 days per experiment. Competitors using Azure ML with distributed training iterate 10x faster. We need MPI support in OpenShift AI to hit FDA submission deadlines and maintain HIPAA compliance."*

- **Automotive OEM**: *"Our autonomous vehicle perception models require 32+ GPUs for training. We evaluated RHOAI but had to choose Databricks because it supports multi-node training. If RHOAI adds MPI support, we'll migrate our $800K/year MLOps spend."*

**Sales Feedback**:
- 23 enterprise opportunities ($12M+ pipeline) have cited distributed training as a "must-have" or "deal-breaker" requirement in FY2025
- 8 customers have explicitly asked "when will RHOAI support Horovod/MPI training?" in product roadmap discussions
- Competitive losses: 6 deals ($4.2M ARR) lost to AWS SageMaker in FY2025, with distributed training cited as the primary differentiator in loss analysis

**Market Research**:
- IDC AI Infrastructure Survey (2024): 84% of organizations are increasing AI infrastructure investment in 2025, with distributed training cited as the #2 priority (after GPU procurement)
- Gartner Hype Cycle for AI (2024): Distributed training platforms are in the "Peak of Inflated Expectations" phase, indicating high market demand and FOMO risk for vendors without this capability
- KubeFlow Community Survey (2024): 67% of KubeFlow users run distributed training workloads, with MPI Operator cited as one of the top 5 most-used components

**Ecosystem Alignment:**

**KubeFlow Trainer V2 as Foundation:**
- KubeFlow Trainer V2 (released July 2025) is the next-generation training framework for Kubernetes, unifying multiple training operators (TFJob, PyTorchJob, MPIJob, XGBoostJob, PaddleJob) under a single API
- RedHat is a KubeFlow Steering Committee member, ensuring influence over roadmap and early access to new features
- Adopting Trainer V2 positions RHOAI as a modern, upstream-aligned platform (vs. competitors building proprietary solutions)

**NVIDIA Ecosystem Integration:**
- MPIJobs with NCCL (NVIDIA Collective Communications Library) enable near-linear scaling efficiency (85-95%) for multi-GPU training, maximizing ROI on NVIDIA GPU investments
- NVIDIA provides reference architectures and validated configurations for KubeFlow MPI Operator with A100/H100 GPUs
- RedHat-NVIDIA partnership (announced 2023) includes joint go-to-market for AI infrastructure, with MPIJobs as a key capability

**OpenShift Platform Synergies:**
- MPIJobs leverage core OpenShift capabilities: RBAC, NetworkPolicies, ResourceQuotas, monitoring (Prometheus/Grafana), logging (EFK/Loki)
- Integration with OpenShift AI Data Science Projects (namespaces) enables seamless multi-tenant isolation for MPIJobs
- GPU Operator integration ensures automatic GPU driver management for MPI worker pods

**Competitive Landscape:**

| Platform | MPI Support | Hybrid Cloud | Open Source | Unified Observability | Maturity |
|----------|-------------|--------------|-------------|----------------------|----------|
| **AWS SageMaker** | ✅ (Native) | ❌ (Cloud-only) | ❌ (Proprietary) | ✅ | High |
| **Azure ML** | ✅ (Native) | ❌ (Cloud-only) | ❌ (Proprietary) | ✅ | High |
| **Google Vertex AI** | ✅ (Native) | ❌ (Cloud-only) | ❌ (Proprietary) | ✅ | High |
| **Databricks** | ✅ (Horovod) | ⚠️ (Limited) | ⚠️ (Partially) | ✅ | High |
| **Run:ai** | ✅ (KubeFlow) | ✅ | ⚠️ (Commercial) | ⚠️ | Medium |
| **Determined AI** | ✅ (Native) | ✅ | ✅ (Open Core) | ✅ | Medium |
| **RHOAI (Current)** | ❌ | ✅ | ✅ | ✅ | N/A |
| **RHOAI (With MPIJobs)** | ✅ | ✅ | ✅ | ✅ | Target: High |

**Key Differentiators with MPIJobs:**
- **Only enterprise platform** with MPI distributed training + hybrid cloud + 100% open source foundation
- **Unified observability** across all KubeFlow Trainer job types (competitors have fragmented UIs for different training frameworks)
- **Enterprise-grade** security, multi-tenancy, and compliance (vs. Run:ai/Determined which have limited enterprise features)
- **RedHat support** and SLA guarantees (vs. self-managed KubeFlow upstream)

**Return on Investment (ROI) for Customers:**

Based on customer interviews and market analysis, MPIJobs deliver measurable ROI across three dimensions:

**1. Training Time Reduction** → Faster Time-to-Market
- **Baseline**: Single-node training for 7B LLM fine-tuning = 120 hours (5 days)
- **With MPIJobs**: 8-node training (32 GPUs) = 4-5 hours
- **Speedup**: 25-30× (near-linear scaling)
- **Business Impact**: Data scientists iterate 25× faster, reducing model development cycle from months to weeks
- **Revenue Impact**: Faster model deployment = earlier revenue realization ($500K-2M per model for production use cases)

**2. GPU Utilization** → Cost Optimization
- **Baseline**: Single-node training utilizes 40-50% of GPU capacity (memory-bound)
- **With MPIJobs**: Distributed training achieves 80-90% GPU utilization
- **Cost Savings**: 2× more effective GPU capacity without hardware procurement
- **Example**: Customer with 128 GPUs ($5M infrastructure) unlocks additional 64 GPU-equivalent capacity ($2.5M value)

**3. Hybrid Cloud Flexibility** → Avoid Cloud Lock-In
- **Baseline**: Customers forced to migrate to AWS/Azure for distributed training pay 3-5× markup on cloud GPU hours vs. on-prem TCO
- **With MPIJobs**: Run distributed training on-prem, reserve cloud for bursting
- **Cost Avoidance**: $50K-200K/year in cloud GPU costs for typical enterprise workload
- **Compliance**: Avoid data egress fees and regulatory risk of moving sensitive data to public cloud

**Total Estimated ROI**: 210% over 3 years (based on Forrester Total Economic Impact framework applied to similar ML platforms)

**Strategic Risks of Not Delivering:**

**Market Share Erosion**: Without MPIJobs, RHOAI is excluded from the fastest-growing segment of the enterprise ML market (LLM/foundation model training), ceding market share to AWS, Azure, and emerging competitors.

**Customer Churn**: Existing customers with scaling model training needs will adopt hybrid architectures (RHOAI for inference, SageMaker for training), increasing platform fragmentation and churn risk.

**Talent Acquisition**: Data scientists and MLOps engineers prefer platforms with modern capabilities (distributed training, GPU optimization). Lack of MPIJobs makes RHOAI less attractive to top talent, hurting customer adoption.

**Partner Ecosystem**: NVIDIA, Intel, AMD, and cloud providers are investing heavily in distributed training ecosystems. Without MPIJobs, RedHat risks being excluded from co-marketing and joint solution development.

**Timeline and Dependencies:**

**Critical Path Dependencies:**
1. **KubeFlow Trainer V2 Integration**: RHOAI must first integrate KubeFlow Trainer V2 as the foundation layer (dependency for this RFE)
2. **MPI Operator Packaging**: MPI Operator v2 must be packaged as a RedHat-supported component with certified container images
3. **Dashboard UI Development**: ODH Dashboard must be extended with MPIJob-specific UI components
4. **Documentation and Examples**: Comprehensive docs and runnable examples must be created for GA release

**Estimated Timeline:**
- **Q4 2025**: Complete KubeFlow Trainer V2 integration (prerequisite)
- **Q1 2026**: MPIJobs private beta with 5-10 strategic customers
- **Q2 2026**: MPIJobs public beta, gather feedback and iterate
- **Q3 2026**: MPIJobs GA release as part of RHOAI 2.17 (target)
- **Q4 2026**: Post-MVP enhancements (gang scheduling, Katib integration, RDMA support)

---

## Customer Considerations

**Enterprise Deployment Scenarios:**

MPIJobs must support diverse enterprise deployment models and constraints:

**1. On-Premises GPU Clusters**
- **Scenario**: Large enterprises (financial services, healthcare, government) with dedicated on-prem GPU infrastructure for compliance and data sovereignty
- **Requirements**:
  - Support for air-gapped environments (disconnected from internet)
  - Integration with internal container registries (Artifactory, Harbor, Quay)
  - Compatibility with enterprise networking (firewalls, proxies, segmented networks)
  - GPU driver management for diverse hardware (NVIDIA A100, H100, A40, older V100)
- **Success Criteria**: MPIJobs run successfully in air-gapped environment with all dependencies available in internal registry

**2. Hybrid Cloud (On-Prem + Public Cloud)**
- **Scenario**: Enterprises running primary workloads on-prem with cloud bursting for peak demand or specific use cases
- **Requirements**:
  - MPIJobs must run identically on-prem and in cloud (AWS, Azure, GCP) with same configuration
  - Support for workload federation (e.g., data preprocessing on-prem, training in cloud)
  - Cross-cluster model artifact sharing (S3-compatible object storage)
  - Network latency tolerance (acknowledgment that inter-cluster MPI is out of scope, but data transfer workflows must be smooth)
- **Success Criteria**: Customer can develop MPIJob configuration on-prem, deploy to AWS/Azure OpenShift without modifications

**3. Multi-Tenant Shared Clusters**
- **Scenario**: Large organizations with centralized OpenShift AI platform serving multiple business units, teams, or projects
- **Requirements**:
  - Strong namespace isolation (team A cannot view/modify team B's MPIJobs)
  - Fair-share scheduling (no single team monopolizes GPUs)
  - Chargeback/showback (track GPU-hours per team for cost allocation)
  - Admin visibility across all namespaces for capacity planning and troubleshooting
- **Success Criteria**: 5 teams run concurrent MPIJobs on shared 128-GPU cluster with fair allocation, zero interference, complete audit trail

**Security and Compliance Considerations:**

**Regulatory Compliance:**
- **HIPAA (Healthcare)**: MPIJob pods must support encrypted data at rest (PVC encryption), encrypted data in transit (TLS for S3 access), audit logging of all data access, and no PHI in container images or logs
- **GDPR (EU Data Privacy)**: Training data must remain in EU regions, audit logs must capture data lineage, right-to-deletion must extend to training jobs and model artifacts
- **FedRAMP (US Government)**: MPIJobs must run on FedRAMP-certified OpenShift clusters, all container images must pass NIST vulnerability scans, FIPS 140-2 compliance for cryptographic operations
- **PCI-DSS (Payment Card Industry)**: Training on transaction data requires network segmentation (NetworkPolicies), access control (RBAC), and quarterly security audits

**Security Best Practices:**
- **Least Privilege**: MPIJob pods run with minimal ServiceAccount permissions (no cluster-admin access)
- **Image Security**: Only signed container images from trusted registries, vulnerability scanning with Clair/Trivy, no images with critical CVEs
- **Secret Management**: Training script secrets (API keys, DB passwords) injected via external secret managers (Vault, AWS Secrets Manager) integrated with OpenShift
- **Network Isolation**: NetworkPolicies automatically applied to restrict MPI traffic to same-job pods, deny cross-namespace communication
- **Audit Logging**: All MPIJob lifecycle events (create, delete, complete, fail) logged to SIEM with user identity, timestamp, and resource details

**Data Governance:**
- **Data Lineage**: Track which datasets were used to train which models (integration with Model Registry and ODH Data Catalog)
- **Data Access Control**: Training jobs inherit user's data access permissions (e.g., RBAC on PVCs, S3 bucket policies)
- **Data Residency**: Ability to enforce geographic constraints (e.g., EU data must be trained on EU clusters, never replicated to US)

**Performance and Scalability Considerations:**

**Target Workload Profiles:**
- **Small-scale** (2-4 nodes, 4-16 GPUs): Rapid prototyping, hyperparameter tuning, research experiments
  - Expected completion time: 30 minutes - 4 hours
  - Resource allocation: CPU=8, Memory=64Gi, GPU=1-4 per worker
- **Medium-scale** (8-16 nodes, 32-64 GPUs): Production model training, LLM fine-tuning (7B-13B params), computer vision models
  - Expected completion time: 2-12 hours
  - Resource allocation: CPU=16, Memory=128Gi, GPU=2-4 per worker
- **Large-scale** (32+ nodes, 128+ GPUs): Foundation model training, 70B+ LLM fine-tuning, massive-scale hyperparameter optimization
  - Expected completion time: 12-72 hours
  - Resource allocation: CPU=32, Memory=256Gi, GPU=4-8 per worker

**Scalability Targets:**
- **Concurrent Jobs**: Platform must support 50+ concurrent MPIJobs across multiple namespaces without API server degradation or scheduling delays
- **Cluster Size**: Validated configurations for 64-GPU, 128-GPU, 256-GPU, 512-GPU clusters
- **Worker Count**: Single MPIJob must support up to 100 workers (post-MVP target: 500+ workers)
- **Data Scale**: Training datasets up to 10TB (stored on shared PVCs or object storage)

**Performance Benchmarking:**
- Publish reference benchmarks for common workloads:
  - ResNet-50 training (ImageNet): Baseline throughput (images/sec) on 2/4/8/16 nodes
  - BERT-Large fine-tuning: Baseline training time on 2/4/8 nodes
  - GPT-7B fine-tuning: Baseline training time and scaling efficiency on 8/16/32 nodes
- Enable customers to reproduce benchmarks with provided scripts and configurations
- Provide guidance on expected scaling efficiency (target: 85% scaling efficiency up to 16 nodes, 75% up to 32 nodes)

**Cost Optimization Guidance:**

**Cloud Cost Management:**
- For hybrid deployments, provide guidance on when to train on-prem vs. cloud:
  - On-prem for: Sensitive data, long-running jobs, repeated experiments
  - Cloud bursting for: Spike workloads, short experiments, hardware unavailability on-prem
- Document cost comparison methodology (on-prem TCO vs. cloud GPU-hour pricing)

**GPU Utilization Optimization:**
- Provide best practices for maximizing GPU utilization:
  - Right-sizing worker resources (avoid over-provisioning memory/CPU)
  - Mixed precision training (fp16/bf16) to double throughput
  - Gradient accumulation to simulate larger batch sizes without OOM
  - Data loading optimization (prefetch, multi-threaded loading) to avoid GPU starvation

**Usability and Adoption Considerations:**

**Skill Level Accommodation:**
- **Junior Data Scientists**: Focus on Dashboard UI wizard with templates, minimal Kubernetes knowledge required
- **Experienced Engineers**: Provide CLI and SDK for programmatic workflows, CI/CD integration
- **Administrators**: Offer advanced configuration options (resource policies, monitoring, security)

**Migration Path from Competitive Platforms:**
- Provide migration guides for customers moving from:
  - **AWS SageMaker**: Map SageMaker Distributed Training API to MPIJob SDK, config file converter
  - **Azure ML**: Map Azure ML MPI jobs to RHOAI MPIJob, similar patterns
  - **Standalone KubeFlow**: Guide for migrating from upstream KubeFlow MPI Operator to RHOAI-managed version
- Include automated migration tools where feasible (e.g., YAML converter from SageMaker training job to MPIJob)

**Support and Troubleshooting:**
- **RedHat Support Readiness**: Train support engineers on MPIJob troubleshooting, create internal runbooks for common issues
- **Customer Self-Service**: Provide diagnostic tools (e.g., `odh mpijob diagnose <job-name>` command that runs health checks)
- **Community Resources**: Establish RHOAI community forum for MPIJob questions, monitor and respond to issues

**Localization and Internationalization:**
- **UI Localization**: Dashboard UI should support localization (initial: English, future: Japanese, Chinese, German) based on customer demand
- **Documentation**: Evaluate demand for localized documentation for key markets (Asia-Pacific, EMEA)
- **Time Zones and Formats**: Dashboard displays timestamps in user's local timezone, supports date/time formats per locale

**Accessibility:**
- **WCAG 2.1 Level AA Compliance**: Dashboard UI must be accessible to users with disabilities (screen reader support, keyboard navigation, color contrast)
- **CLI Accessibility**: Provide structured output formats (JSON, YAML) parseable by assistive tools

**Training and Enablement:**

**Customer Onboarding:**
- **Quickstart Workshop**: 2-hour hands-on workshop for new customers (provided as part of RHOAI onboarding)
  - Module 1: Concepts (30 min)
  - Module 2: First MPIJob via Dashboard (30 min)
  - Module 3: SDK usage in Jupyter (30 min)
  - Module 4: Monitoring and troubleshooting (30 min)
- **Office Hours**: Weekly drop-in sessions with RedHat solutions architects for MPIJob questions
- **Certification**: Include MPIJob content in RHOAI certification exam (Red Hat Certified Specialist in OpenShift AI)

**Field Enablement:**
- **Sales Enablement Kit**: Slides, demo script, ROI calculator, competitive objection handling for account teams
- **Solutions Architect Training**: 1-day technical deep-dive on MPIJob architecture, deployment, troubleshooting
- **Partner Training**: Extend training to RedHat partners (system integrators, ISVs) for broader ecosystem reach

**Change Management:**

**Organizational Readiness:**
- Many enterprises have existing HPC/MPI expertise (computational scientists, HPC admins) but these teams may not be familiar with Kubernetes
- MPIJobs bridge HPC and cloud-native worlds - provide documentation that speaks to both audiences:
  - **For HPC users**: "How MPIJobs compare to Slurm/PBS", "Mapping MPI Allreduce to Kubernetes pods"
  - **For Kubernetes users**: "MPI fundamentals for cloud-native engineers", "How launcher-worker architecture maps to Kubernetes jobs"
- Encourage cross-functional collaboration between data science, ML engineering, and platform engineering teams

**Rollout Strategy:**
- **Phase 1**: Private beta with early adopters (high-touch support, frequent feedback cycles)
- **Phase 2**: Public beta with self-service documentation (broader audience, community support)
- **Phase 3**: GA release with full support, SLAs, and production readiness guarantees
- **Phase 4**: Post-GA enhancements based on customer feedback and competitive analysis

**Success Metrics and KPIs:**

**Adoption Metrics:**
- **30-day activation**: 50% of RHOAI customers run at least 1 MPIJob within 30 days of feature GA
- **90-day adoption**: 30% of customers have MPIJobs in production (>50 jobs executed)
- **Power users**: 10% of customers run 500+ MPIJobs within first 6 months (indicating CI/CD integration)

**Business Impact Metrics:**
- **Revenue influence**: MPIJobs feature influences $5M+ in new ARR within first year
- **Win rate improvement**: Competitive win rate vs. AWS/Azure improves by 15 percentage points in deals where distributed training is a requirement
- **Expansion revenue**: 20% of customers increase RHOAI spend (more users, more namespaces) after adopting MPIJobs

**Customer Satisfaction Metrics:**
- **NPS (Net Promoter Score)**: MPIJob feature NPS >50 (measured via in-product survey)
- **Support ticket reduction**: 50% reduction in distributed training-related support tickets (vs. customers using workarounds)
- **Retention**: Customers using MPIJobs have 10 percentage point higher retention rate vs. non-users

---

**End of RFE Document**