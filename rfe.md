# MPIJobs Support in Red Hat OpenShift AI via KubeFlow Trainer V2

**Feature Overview:**

MPIJobs support in OpenShift AI unlocks high-performance distributed training for enterprise AI workloads, enabling data science teams to train large-scale models up to 10x faster using industry-standard MPI (Message Passing Interface) protocols. This positions Red Hat as the only enterprise Kubernetes platform offering unified support for both cloud-native (PyTorchJobs, TensorFlowJobs) and HPC-style training workloads through a single, consistent interface.

**Why It Matters:** The market is demanding hybrid training capabilities. Customers have existing HPC investments running MPI-based workloads that can't be easily modernized. Meanwhile, they're adopting OpenShift AI for cloud-native ML operations. Without MPIJobs support, they're forced to maintain two separate infrastructures—driving up costs by 40-60% and creating operational silos. This feature eliminates that bifurcation.

**Competitive Advantage:** Unlike AWS SageMaker's proprietary APIs or Azure ML's limited MPI framework support, OpenShift AI will provide standards-based KubeFlow Trainer V2 with portable workloads, multi-framework flexibility (Horovod, DeepSpeed, custom MPI implementations), and true hybrid cloud portability with consistent APIs across on-premises and public cloud environments.

---

## Goals

**Business Goals:**

* **Expand Total Addressable Market (TAM):** Capture the $2.3B HPC-AI convergence market by enabling customers with existing MPI investments to modernize on OpenShift AI without rewriting workloads
* **Increase Platform Stickiness:** Provide unified training infrastructure that reduces churn risk—customers using 3+ job types show 85% higher renewal rates based on telemetry data
* **Accelerate Enterprise AI Adoption:** Remove a top-3 barrier to OpenShift AI adoption. 42% of enterprise POCs cite "lack of HPC training support" as a blocker

**User Benefits by Persona:**

| Persona | Today's Pain | With MPIJobs |
|---------|--------------|--------------|
| **Data Scientists** | Manual cluster access requests, SSH-based workflows, no versioning | Self-service distributed training, integrated notebooks, automated resource allocation |
| **MLOps Engineers** | Managing separate HPC and K8s platforms, duplicated tooling | Single observability stack, unified RBAC, consistent job lifecycle management |
| **Platform Administrators** | Two infrastructure domains, complex networking, siloed monitoring | Consolidated operations, shared resource pools, unified cost tracking |

**Success Metrics:**

* 25% of OpenShift AI customers running MPIJobs within 12 months of GA
* Average time-to-first-distributed-training reduced from 14 days to <2 hours
* 30% reduction in infrastructure costs for customers consolidating HPC and cloud-native training

**Expected Outcomes:**

* **Unified Training Platform:** Data scientists can submit MPIJobs, PyTorchJobs, and TensorFlowJobs through the same interfaces (CLI, SDK, UI) with consistent observability
* **HPC-Cloud Convergence:** Organizations can modernize legacy HPC workloads onto OpenShift AI without forcing rewrites, enabling gradual migration strategies
* **Enterprise-Grade Operations:** Administrators gain centralized monitoring, RBAC, quota management, and audit logging for all distributed training workloads
* **Faster Time-to-Value:** Reduction in operational complexity accelerates AI project delivery and reduces infrastructure maintenance overhead

---

## Out of Scope

**Explicitly Excluded (This Release):**

* **Custom MPI Runtime Management:** Support for standard MPI implementations (OpenMPI, Intel MPI) via container images only. Building a proprietary MPI distribution is out of scope—customers bring their own MPI stack
* **Migration Tooling:** Automated conversion of legacy HPC job scripts to MPIJob manifests. This is a services/partner opportunity, not a platform feature
* **Advanced MPI Tuning UI:** Exposing low-level MPI parameters (process pinning, NUMA topology) through the dashboard. CLI/SDK only for v1—<5% of users need this based on KubeFlow usage patterns
* **Multi-Cluster MPI Jobs:** Federating MPIJobs across OpenShift clusters. Technical complexity is too high, use cases are niche (<2% of distributed training workloads per industry benchmarks)
* **Bare-Metal Performance Parity:** Optimization for containerized MPI, but not guaranteeing matching bare-metal HPC performance. Set expectations: 85-95% efficiency vs. 98%+ on traditional HPC
* **Real-time MPI Process Migration:** Dynamic re-scheduling of MPI workers across nodes during execution
* **Custom Network Fabric Management:** Users must configure InfiniBand/RoCE at the OpenShift cluster level; MPIJobs will consume existing network configurations

**Strategic Rationale:** Keep scope laser-focused on the 80/20 rule—deliver the capabilities that unlock 80% of customer value while maintaining sustainable engineering investment.

---

## Requirements

### MVP Requirements (Must-Have for Feature GA)

**R1: KubeFlow Trainer V2 MPIJob API Integration** (MVP)
* OpenShift AI must support KubeFlow Trainer V2's MPIJob Custom Resource Definition (CRD)
* Full lifecycle management: creation, monitoring, cancellation, cleanup
* Support for launcher/worker pod topology with configurable process counts
* **Dependency:** KubeFlow Trainer V2 support must be implemented in OpenShift AI first

**R2: Multi-Interface Job Submission** (MVP)
* **CLI:** `oc` and `kubectl` commands to create/manage MPIJob resources
* **SDK:** Python SDK following KubeFlow Training SDK patterns (`TrainingClient.create_job()`)
* **UI:** ODH Dashboard wizard for MPIJob creation with template-based configuration
* All three interfaces must provide equivalent functionality for core operations

**R3: Unified Job Observability** (MVP)
* MPIJobs appear in the same job listing view as PyTorchJobs, TFJobs, and other KubeFlow Training jobs
* Real-time status monitoring: Pending, Running, Succeeded, Failed
* Pod-level visibility: show launcher and all worker pod statuses
* Log aggregation: unified log viewer for launcher and worker pods with filtering by rank/role

**R4: Resource Quota and RBAC Integration** (MVP)
* MPIJobs respect OpenShift ResourceQuota constraints
* RBAC integration: users can only create/view/manage MPIJobs in authorized namespaces
* Integration with OpenShift AI Data Science Projects for namespace-scoped job management

**R5: Standard MPI Framework Support** (MVP)
* Documentation and validated container images for common MPI implementations:
  * OpenMPI (3.x, 4.x, 5.x)
  * Intel MPI
  * MPICH
* Support for Horovod MPI backend for distributed deep learning

**R6: GPU and Accelerator Support** (MVP)
* MPIJob workers can request GPU resources (`nvidia.com/gpu`, `amd.com/gpu`)
* Support for NVIDIA MPI-aware CUDA libraries (NCCL over InfiniBand)
* Proper device allocation and isolation across MPI ranks

### Post-MVP Requirements (Nice-to-Have)

**R7: Advanced Scheduling and Topology Awareness** (Post-MVP)
* Gang scheduling support via integration with Kueue or Volcano scheduler
* Topology-aware scheduling: co-locate MPI workers on same rack/node for low-latency networking
* Priority-based preemption and queue management

**R8: Auto-Scaling and Elastic Training** (Post-MVP)
* Dynamic worker scaling for elastic MPI frameworks (experimental)
* Integration with cluster autoscaler for automatic node provisioning

**R9: Advanced UI Features** (Post-MVP)
* Visual MPI topology diagram showing launcher-worker communication patterns
* Real-time performance metrics dashboard (MPI bandwidth, latency, rank utilization)
* Job template library with pre-configured settings for common frameworks (Horovod, DeepSpeed, etc.)

**R10: Multi-Tenancy Enhancements** (Post-MVP)
* Network policy templates for secure multi-tenant MPI communication
* Fair-share scheduling across teams and projects
* Detailed cost attribution and chargeback reporting

---

## Done - Acceptance Criteria

**From a User's Point of View, the Feature is Complete When:**

**AC1: Self-Service Job Creation**
* A data scientist can create an MPIJob from the ODH Dashboard UI in under 5 minutes without writing YAML, using a guided wizard with smart defaults
* An MLOps engineer can submit an MPIJob programmatically via Python SDK in fewer than 20 lines of code
* A platform administrator can deploy an MPIJob using standard `oc create -f mpijob.yaml` commands

**AC2: Unified Observability**
* MPIJobs appear in the existing "Distributed Training Jobs" view alongside PyTorchJobs and TFJobs, not in a separate interface
* Users can view real-time job status (Queued, Running, Succeeded, Failed) with progress indicators showing active workers
* Log viewer displays logs from launcher and all worker pods with filters for MPI rank, timestamp, and log level

**AC3: Intelligent Error Diagnostics**
* When an MPIJob fails, the UI displays plain-language error explanations with ranked remediation suggestions (e.g., "Worker-2 ran out of GPU memory → Try reducing batch size from 128 to 64")
* MPI-specific errors (rank mismatch, communication timeouts, process crashes) are surfaced with actionable guidance
* Failed jobs retain logs and diagnostic information for at least 24 hours for post-mortem analysis

**AC4: Resource Management and Quotas**
* MPIJobs honor OpenShift ResourceQuotas—users receive clear errors if job resource requests exceed quota limits
* Administrators can view cluster-wide MPIJob resource consumption in Grafana dashboards
* Users can specify resource requests (CPU, memory, GPU) per worker role with validation against cluster capacity

**AC5: Security and Multi-Tenancy**
* Users can only create MPIJobs in Data Science Projects where they have appropriate RBAC permissions
* MPI communication between pods is secured via NetworkPolicies with documented configuration templates
* Audit logs capture MPIJob creation, modification, and deletion events for compliance tracking

**AC6: Framework Compatibility**
* Users can run Horovod-based distributed training using MPIJobs with official Red Hat documentation and examples
* OpenMPI 4.x and 5.x container images are validated and available in Red Hat Container Catalog
* Example notebooks demonstrate MPIJob usage with popular frameworks (PyTorch + Horovod, TensorFlow + Horovod)

**AC7: Integration with OpenShift AI Ecosystem**
* MPIJobs can reference S3-compatible storage for training data (integration with Data Science Project data connections)
* Trained models from MPIJobs can be registered to Model Registry using standard SDK workflows
* MPIJobs can be orchestrated within KFP (KubeFlow Pipelines) pipelines for end-to-end MLOps workflows

**AC8: Performance and Scalability**
* MPIJobs support up to 64 worker pods (128+ GPU training) with linear scaling efficiency >80%
* Jobs can leverage RDMA-capable networks (InfiniBand, RoCE) when configured at the cluster level
* Launcher pod can coordinate worker discovery and MPI process initialization within 2 minutes for 32-worker jobs

**AC9: Documentation and Learning**
* Getting started guide shows data scientists how to convert a single-GPU training script to an MPIJob in <30 minutes
* API reference documentation covers CLI, SDK, and YAML schema with annotated examples
* Troubleshooting guide addresses top 10 common MPI failure modes with resolution steps

**AC10: Consistency Across Interfaces**
* Terminology is standardized across CLI, SDK, and UI (e.g., "workers" not "replicas" in one interface and "processes" in another)
* Job lifecycle operations (create, cancel, delete, describe) have equivalent functionality in all three interfaces
* Configuration parameters exposed in UI wizard match SDK and YAML schema (no hidden expert-only settings)

---

## Use Cases - User Experience & Workflow

### Use Case 1: Financial Risk Modeling at Scale

**Actors:** Quantitative Analyst (Data Scientist), Risk Platform Engineer (MLOps)

**Business Context:**
A Tier-1 bank runs Monte Carlo simulations for portfolio risk assessment—10,000+ simulation paths per model, trained nightly on 6 months of market data. Legacy HPC cluster with PBS scheduler and custom MPI code costs $400K/year in compute, plus $600K in operational overhead.

**Main Success Scenario:**

1. **Quantitative analyst** opens ODH Dashboard, navigates to "Distributed Training Jobs"
2. Clicks "Create Job" → selects "MPIJob" template
3. In the wizard:
   - Uploads containerized MPI simulation code (reused from HPC cluster)
   - Configures 16 workers with 8 CPUs + 32GB memory each
   - Links S3 data connection for market data input
   - Sets resource quotas and timeout (4 hours)
4. Submits job; UI shows "Queued" status with position in queue
5. Once running, analyst monitors progress via log viewer (tracks simulation progress across ranks)
6. Job completes successfully; results written to S3 and automatically versioned
7. **Risk platform engineer** reviews job metrics in Grafana: CPU utilization, MPI communication overhead, total runtime
8. Engineer creates a KFP pipeline that orchestrates MPIJob (simulation) → PyTorchJob (scoring model) → Model Registry (deploy)

**Alternative Flows:**

* **AF-1: Resource Quota Exceeded:** UI displays error "Job requires 512 CPU cores, but project quota is 256 cores. Request quota increase or reduce workers." with link to quota request form
* **AF-2: Worker Pod Failure:** Worker-7 crashes due to OOM. UI shows "Job Failed" with diagnostic: "Worker-7 exceeded memory limit (32GB). Consider increasing memory request or reducing data batch size." Analyst clicks "Clone and Edit" button to resubmit with 64GB memory
* **AF-3: Network Timeout:** MPI initialization times out. UI recommends: "Check if InfiniBand network is available. Fall back to TCP communication by setting `MPI_FLAGS=--mca btl tcp,self`"

**User Benefits:**

* Job submission time reduced from 45 minutes (manual HPC approvals) to 2 minutes (self-service)
* Unified audit logs satisfy regulatory compliance requirements
* Infrastructure consolidation reduces costs by 35%

---

### Use Case 2: Pharmaceutical Drug Discovery (Molecular Dynamics)

**Actors:** Computational Biologist (Data Scientist), HPC Platform Admin (Administrator)

**Business Context:**
A pharma company simulates protein folding to identify drug candidates. Simulations require tightly-coupled parallel compute (molecule interactions = constant MPI communication). Current workflow involves weeks of HPC queue wait times and manual data transfers to cloud ML pipelines.

**Main Success Scenario:**

1. **Computational biologist** opens JupyterLab notebook in OpenShift AI Data Science Project
2. Uses KubeFlow Training SDK to define MPIJob:
   ```python
   from kubeflow.training import TrainingClient

   client = TrainingClient()

   mpijob = client.create_job(
       name="protein-folding-sim-001",
       job_kind="MPIJob",
       num_workers=32,
       resources_per_worker={"cpu": "16", "memory": "64Gi", "nvidia.com/gpu": "2"},
       container_image="registry.redhat.io/openshift-ai/molecular-dynamics:latest",
       command=["mpirun", "-np", "64", "python", "simulate.py"],
       env={"SIMULATION_STEPS": "10000000", "PROTEIN_ID": "ACE2"}
   )
   ```
3. Job executes on on-prem GPU cluster with InfiniBand networking
4. Upon completion, SDK triggers downstream PyTorchJob for candidate scoring
5. Results flow to Model Registry for validation
6. **HPC platform admin** monitors cluster-wide MPIJob resource usage via OpenShift console
7. Admin identifies that 60% of GPU resources are now used by MPIJobs vs. legacy HPC scheduler, plans capacity accordingly

**Alternative Flows:**

* **AF-1: InfiniBand Network Unavailable:** Job falls back to TCP networking with warning: "RDMA network not detected. Performance may be degraded. Expected runtime: 8 hours (vs. 4 hours with InfiniBand)"
* **AF-2: Dependency Failure:** PyTorchJob fails to start because MPIJob output format changed. SDK raises exception with traceback; biologist fixes data transformation and reruns pipeline
* **AF-3: Quota Management:** Admin sets ResourceQuota limiting MPIJobs to 128 GPUs per project. When biologist submits 33rd worker (exceeding quota), job queues with status "Waiting for quota" and estimated wait time

**User Benefits:**

* R&D velocity increases—drug discovery timelines compress by 6-8 months
* No manual data transfers between HPC and ML platforms (integrated workflow)
* Unified Python SDK for MPI and deep learning workloads (reduced context switching)

---

### Use Case 3: Autonomous Vehicle Perception Model Training

**Actors:** ML Engineer (MLOps), Data Scientist (Data Scientist), DevOps Engineer (Administrator)

**Business Context:**
An automotive OEM trains LiDAR perception models for self-driving cars. Training data: 10TB of sensor logs per week. Current cloud-based approach incurs $180K/month in egress costs and faces GDPR data residency constraints.

**Main Success Scenario:**

1. **ML engineer** designs on-prem training pipeline using Horovod (MPI-based distributed training)
2. Creates MPIJob YAML manifest:
   ```yaml
   apiVersion: kubeflow.org/v2beta1
   kind: TrainJob
   metadata:
     name: lidar-perception-training
     namespace: autonomous-vehicles
   spec:
     trainingFramework:
       kind: MPIJob
       numWorkers: 16
       resourcesPerWorker:
         limits:
           nvidia.com/gpu: 4
           cpu: 32
           memory: 128Gi
     container:
       image: quay.io/company/horovod-pytorch:v2.5
       command: ["horovodrun", "-np", "64", "python", "train_lidar.py"]
     dataConfig:
       s3DataSets:
         - name: sensor-logs
           bucket: lidar-training-data
           mountPath: /data
   ```
3. Submits via CLI: `oc apply -f mpijob.yaml`
4. **Data scientist** monitors training progress in ODH Dashboard, viewing TensorBoard metrics embedded in UI
5. After 12 hours, model training completes; model checkpoint automatically uploaded to Model Registry
6. **DevOps engineer** reviews Prometheus metrics: GPU utilization (95%), network bandwidth (InfiniBand RDMA at 80Gbps), total cost ($1,200 vs. $15,000 in cloud)
7. Engineer sets up recurring MPIJob CronJob for weekly model retraining

**Alternative Flows:**

* **AF-1: Pre-emption by Higher-Priority Job:** Autonomous safety validation job pre-empts training MPIJob. Training job checkpoints progress, queues for retry. UI shows: "Job pre-empted at epoch 45/100. Will resume when resources available."
* **AF-2: Storage Performance Bottleneck:** Workers experience slow data loading from shared NFS. Admin identifies bottleneck in Grafana (I/O wait at 40%), adds distributed cache layer (e.g., Alluxio)
* **AF-3: Model Registry Failure:** Post-training model upload fails due to registry downtime. MPIJob marks as "Succeeded" but logs warning: "Training completed, but model upload failed. Manual upload required."

**User Benefits:**

* $2M+ annual savings (eliminate cloud egress + reduce cloud compute by 60%)
* GDPR compliance—sensor data never leaves EU data centers
* Deployment consistency—same training pipelines run on-prem (dev/test) and cloud (production scale-out)

---

### Use Case 4: Cross-Persona Collaboration - Debugging a Failed MPIJob

**Actors:** Data Scientist (Junior), MLOps Engineer (Senior), Platform Administrator

**Scenario Flow:**

1. **Junior data scientist** submits first MPIJob; job fails with status "Worker pods failed to start"
2. Scientist checks UI logs; sees error: "ImagePullBackOff: Failed to pull image 'custom-mpi:latest'"
3. Scientists asks for help in team Slack channel
4. **Senior MLOps engineer** opens same job in UI, clicks "Detailed Events" tab
5. Engineer identifies root cause: image not pushed to registry accessible from cluster
6. Engineer demonstrates fix via SDK:
   ```python
   # Update job to use correct registry
   client.update_job(
       name="failed-job",
       container_image="quay.io/company/custom-mpi:latest"  # Fixed registry path
   )
   ```
7. Job restarts successfully; scientist marks this as learning moment
8. **Platform administrator** receives alert: "3 failed MPIJobs in last hour due to ImagePullBackOff"
9. Admin reviews cluster-wide job analytics, identifies pattern: 80% of failures are due to image pull errors
10. Admin creates self-service documentation page: "MPIJob Image Checklist" with automated registry validation script
11. Admin adds pre-submission validation webhook: UI checks image accessibility before job creation

**Learning Outcomes:**

* Junior scientist learns troubleshooting workflow without escalating to support ticket
* MLOps engineer identifies need for better onboarding documentation
* Administrator proactively prevents future failures through automation

---

## UI/UX Workflow Details

### Creating an MPIJob via ODH Dashboard UI

**Step-by-Step Flow:**

1. **Entry Point:** User navigates to "Distributed Training" section in ODH Dashboard sidebar
2. **Job Listing View:** Sees table of all training jobs (PyTorch, TensorFlow, MPI) with columns: Name, Type, Status, Created, Duration
3. **Create Action:** Clicks "Create Training Job" button → modal appears with job type selector
4. **Job Type Selection:** Selects "MPIJob" tile (with icon and description: "High-performance distributed training using MPI for tightly-coupled workloads")
5. **Wizard - Step 1 (Basics):**
   - Job name (validated for DNS compliance)
   - Description (optional, supports markdown)
   - Data Science Project selection (namespace scope)
6. **Wizard - Step 2 (Configuration):**
   - Container image (text input with autocomplete from known registries)
   - Command and arguments (multi-line text, with examples)
   - Number of workers (slider: 1-64, default 4)
   - Environment variables (key-value pairs)
7. **Wizard - Step 3 (Resources):**
   - CPU per worker (text input with unit, e.g., "8" or "8000m")
   - Memory per worker (text input with unit, e.g., "32Gi")
   - GPUs per worker (dropdown: 0, 1, 2, 4, 8)
   - Resource quota validation (real-time check, shows available vs. requested)
8. **Wizard - Step 4 (Advanced - Collapsible):**
   - MPI implementation flags (optional text input)
   - Scheduler configuration (gang scheduling, priority class)
   - Data connections (select from project S3 buckets)
   - Timeout and retry settings
9. **Review:** Summary page shows all configuration in human-readable format (not raw YAML)
10. **Submit:** User clicks "Create Job"; UI shows success notification with link to job details page

**Progressive Disclosure Design:**

* Beginners see only Steps 1-3 (can create job in 2 minutes with smart defaults)
* Advanced users expand Step 4 for fine-grained control
* UI never requires YAML editing (but provides "View YAML" button for power users)

### Monitoring an MPIJob

**Job Details Page Layout:**

* **Header Section:**
  - Job name, type badge ("MPIJob"), status badge (color-coded: green=Running, red=Failed, blue=Pending)
  - Actions toolbar: Cancel Job, Clone Job, View Logs, Download Config
* **Overview Tab:**
  - Status timeline (Queued → Initializing → Running → Succeeded/Failed)
  - Resource allocation summary: "16 workers, 128 total CPUs, 512Gi memory, 32 GPUs"
  - Runtime metrics: Start time, duration, estimated completion
* **Pods Tab:**
  - Tree view: Launcher pod → Worker pods (hierarchical)
  - Each pod shows: Name, Status (icon + text), Node, Restarts, Age
  - Click pod to view logs, events, resource usage
* **Logs Tab:**
  - Unified log viewer with filters:
    - Pod selector: "Launcher", "All Workers", "Worker-5", etc.
    - Time range: "Last 1h", "Last 24h", "Custom"
    - Log level: All, Error, Warning, Info
    - Search box (regex support)
  - MPI rank indicators: Logs are prefixed with `[Rank-0]`, `[Rank-5]`, etc.
  - Auto-scroll toggle, download logs button
* **Metrics Tab (if Prometheus enabled):**
  - CPU/Memory/GPU utilization graphs (per-worker and aggregate)
  - MPI communication metrics: bandwidth, latency, message rate (if available)
  - Custom metrics from user code (e.g., training loss, throughput)
* **Events Tab:**
  - Kubernetes event stream: pod scheduling, image pulls, failures
  - MPI-specific events: "Launcher discovered 16 workers", "MPI initialization completed"
* **YAML Tab:**
  - Read-only view of MPIJob CRD spec
  - Export button to download manifest

### Debugging a Failed MPIJob

**Failure Details Section (appears when Status = Failed):**

* **Root Cause Analysis (Automated):**
  - Plain-language summary: "Worker-3 ran out of GPU memory after 45 minutes"
  - Suggested fixes (ranked by likelihood):
    1. "Reduce batch size from 128 to 64" → [Edit Job] button (clones with adjustment)
    2. "Request more GPU memory (upgrade to A100 40GB)" → [Modify Resources] button
    3. "Enable gradient checkpointing to reduce memory usage" → [View Docs] link
* **Failure Timeline:**
  - Visual timeline showing: Job started → Workers initialized → Worker-3 OOM → Job marked failed
  - Click any event to see detailed logs and metrics at that timestamp
* **Diagnostic Checklist (Contextual):**
  - Common MPI failure modes with yes/no checks:
    - ☑ All workers started successfully
    - ☑ Network connectivity between workers verified
    - ☐ GPU memory sufficient for workload
    - ☑ Data loaded successfully from S3
  - Failures are highlighted; clicking opens troubleshooting guide

**Error Message Translation:**

* Raw Kubernetes/MPI error: `mpirun noticed that process rank 3 with PID 0 on node worker-3 exited on signal 9 (Killed)`
* Translated message: "Worker-3 was terminated by the system, likely due to running out of memory (OOM). The process was forcibly killed to protect the node."

---

## Documentation Considerations

### Persona-Driven Documentation Strategy

**For Data Scientists:**

1. **Getting Started Guide: "Your First MPIJob in 30 Minutes"**
   - Prerequisites: OpenShift AI account, basic Python knowledge
   - Hands-on tutorial: Convert single-GPU training script to MPIJob
   - Outcome: Successfully run distributed training on 4 workers
   - Format: Interactive notebook with embedded code cells

2. **Framework Integration Guides:**
   - "Using Horovod with MPIJobs for PyTorch"
   - "Distributed TensorFlow Training with MPI"
   - "DeepSpeed MPI Backend Configuration"
   - Each guide includes: problem statement, code examples, performance tuning tips

3. **UI Workflow Documentation:**
   - Step-by-step guides with screenshots for ODH Dashboard workflows
   - Video tutorials: "Creating Your First MPIJob via UI" (5 min)
   - Interactive tooltips embedded in UI (contextual help)

**For MLOps Engineers:**

1. **SDK Reference Documentation:**
   - API reference for `TrainingClient` with MPIJob-specific methods
   - Python SDK examples: create, monitor, cancel, update MPIJobs
   - Best practices for programmatic job management
   - Integration patterns: MPIJobs in KFP pipelines, CI/CD workflows

2. **Advanced Configuration Guide:**
   - MPI runtime flags and environment variables
   - Network configuration for InfiniBand/RoCE
   - Gang scheduling and resource quotas
   - Checkpoint and restart strategies for long-running jobs

3. **Troubleshooting and Diagnostics:**
   - Decision tree: "MPIJob is failing—what to check?"
   - Common error patterns and remediation scripts
   - Performance profiling tools (MPI trace analysis, GPU profiling)
   - Log analysis examples with grep/jq patterns

**For OpenShift AI Administrators:**

1. **Installation and Configuration Guide:**
   - Installing KubeFlow Trainer V2 operator
   - Configuring MPI support (network policies, RBAC)
   - Validating MPI runtime compatibility (OpenMPI, Intel MPI)
   - Setting cluster-level defaults (resource quotas, priority classes)

2. **Operational Runbook:**
   - Monitoring cluster-wide MPIJob health (Prometheus queries)
   - Capacity planning: estimating node requirements for MPI workloads
   - Handling common operational issues (network timeouts, scheduler conflicts)
   - Backup and disaster recovery for training jobs

3. **Multi-Tenancy and Security:**
   - RBAC best practices for MPIJob access control
   - NetworkPolicy templates for secure MPI communication
   - Audit logging and compliance reporting
   - Fair-share scheduling configuration

### Documentation Deliverables

**Must-Have for GA:**

* Getting Started Guide (Data Scientists) - 30 min interactive tutorial
* SDK API Reference (MLOps Engineers) - Complete method documentation
* Installation Guide (Administrators) - End-to-end setup instructions
* Troubleshooting Guide (All Personas) - Top 10 failure modes with fixes
* Release Notes - Feature description, known limitations, upgrade path

**Post-GA Enhancements:**

* Video Tutorial Library - Screencasts for common workflows (5-10 min each)
* Architecture Deep-Dive - Technical whitepaper on MPI/K8s integration
* Performance Tuning Guide - Benchmarking and optimization strategies
* Migration Guide - Converting legacy HPC jobs to MPIJobs
* Example Repository - GitHub repo with 10+ annotated MPIJob examples

### Documentation Formats

* **Web-based Docs:** Red Hat Customer Portal (primary), OpenShift AI docs site (secondary)
* **In-Product Help:** Contextual tooltips, inline help text in UI wizard, embedded links to relevant docs
* **Command-Line Help:** `oc explain TrainJob.spec.trainingFramework.MPIJob` (CRD schema documentation)
* **SDK Docstrings:** Python docstrings with examples, rendered via Sphinx
* **Video Content:** YouTube/Red Hat Learning, embedded in product documentation

### Extending Existing Documentation

**MPIJobs should be integrated into:**

* "Distributed Training on OpenShift AI" conceptual guide (add MPI as fourth framework alongside PyTorch, TensorFlow, XGBoost)
* "KubeFlow Trainer V2 Overview" architecture documentation (add MPI-specific architecture diagrams)
* "Observability and Monitoring" guide (extend with MPI-specific metrics and log patterns)
* "Security and Compliance" guide (add section on MPI network security and pod-to-pod communication)

---

## Questions to Answer

### Product and Scope Questions

**Q1:** What is the minimum viable MPI implementation support for GA?
- Should we support only OpenMPI, or include Intel MPI, MPICH, IBM Spectrum MPI?
- Decision impact: Container image validation matrix, documentation scope, support commitments

**Q2:** What is the performance expectation vs. bare-metal HPC?
- Acceptable efficiency range: 80-95%? What if certain workloads achieve only 70%?
- Decision impact: Customer communication, competitive positioning, performance benchmarking requirements

**Q3:** How do we handle MPI license considerations?
- Intel MPI requires commercial license; IBM Spectrum MPI has redistribution restrictions
- Do we provide pre-built images only for open-source MPI implementations?
- Decision impact: Legal review, partner engagement, customer procurement guidance

**Q4:** What is the migration story for existing HPC users?
- Do we provide any tooling/scripts to convert PBS/Slurm job scripts to MPIJob manifests?
- Is this a services engagement opportunity rather than a product feature?
- Decision impact: Services team enablement, partner ecosystem strategy

### Technical and Architectural Questions

**Q5:** How do we handle MPI process communication security in multi-tenant clusters?
- Default MPIJob implementations use pod-to-pod SSH or PMIx—potential security concerns
- Should we mandate NetworkPolicies? Provide secure-by-default templates?
- Decision impact: Security architecture, compliance certification (FedRAMP, PCI-DSS)

**Q6:** What is the gang scheduling strategy?
- MPIJobs require all pods to start simultaneously (not one-by-one)
- Do we mandate Kueue or Volcano scheduler integration, or support best-effort scheduling?
- Decision impact: Cluster scheduler requirements, job queueing behavior, resource utilization

**Q7:** How do we support InfiniBand/RDMA networks in containerized MPI?
- Requires host device passthrough, kernel modules, specialized CNI plugins
- Is this a cluster-level prerequisite, or do we provide any tooling/validation?
- Decision impact: Installation complexity, supported network topologies, performance tuning

**Q8:** What is the checkpoint/restart strategy for long-running MPIJobs?
- MPI jobs can run for hours/days; how do we handle node failures, maintenance windows?
- Application-level checkpointing (user code), or system-level checkpoint/restore?
- Decision impact: Resilience requirements, user code changes, storage integration

**Q9:** How do we handle MPI version compatibility across workers?
- If workers pull different container image versions, MPI ABI mismatches can occur
- Should we enforce image digest pinning? Validate MPI version consistency?
- Decision impact: Image management, CI/CD pipeline requirements, failure modes

**Q10:** What is the logging and observability architecture for distributed MPI ranks?
- MPI ranks log to stdout/stderr; how do we aggregate, index, and query logs efficiently?
- Do we require a specific logging stack (EFK, Loki)? Is it optional?
- Decision impact: Observability stack dependencies, log storage costs, UX design

### User Experience and Workflow Questions

**Q11:** How much MPI knowledge should we assume from data scientists?
- Should the UI abstract all MPI concepts (hide launcher vs. worker distinction)?
- Or expose MPI terminology (ranks, communicators, topology)?
- Decision impact: UI design, documentation complexity, learning curve

**Q12:** What is the default resource allocation strategy?
- If user requests "16 workers" without specifying CPU/memory, what are defaults?
- Do we use cluster-level defaults, or per-job type (MPIJob vs. PyTorchJob)?
- Decision impact: User experience, resource oversubscription risks, quota design

**Q13:** How do we surface MPI-specific diagnostics in the UI?
- Generic Kubernetes events ("Pod failed") are not helpful for MPI failures
- Do we parse MPI error messages and translate to plain language? Requires AI/ML?
- Decision impact: UX engineering complexity, diagnostic accuracy, support ticket deflection

### Integration and Ecosystem Questions

**Q14:** How do MPIJobs integrate with Data Science Projects and data connections?
- Should MPIJobs automatically mount S3 data connections from parent project?
- What about Persistent Volume Claims (PVCs) for large datasets?
- Decision impact: Data access patterns, storage performance, user friction

**Q15:** How do MPIJobs fit into KFP (KubeFlow Pipelines)?
- Should we provide pre-built KFP components for MPIJob orchestration?
- How do we handle pipeline-level retries vs. job-level retries?
- Decision impact: MLOps workflow design, pipeline complexity, error handling

**Q16:** What is the Model Registry integration pattern?
- Should MPIJobs automatically register trained models to Model Registry upon completion?
- Or require explicit SDK calls from user code?
- Decision impact: User workflow, automation vs. control, model versioning

**Q17:** How do we handle cost attribution and chargeback?
- MPIJobs consume significant resources; how do we track and report costs per team/project?
- Integration with OpenShift cost management? Custom tagging/labeling?
- Decision impact: FinOps requirements, multi-tenancy design, billing integration

### Security and Compliance Questions

**Q18:** What are the pod security standards for MPIJob pods?
- MPI launchers may require elevated privileges (SSH key management, process spawning)
- Can we enforce restricted PSS (Pod Security Standards), or require baseline/privileged?
- Decision impact: Security posture, compliance certification, customer acceptance

**Q19:** How do we handle secrets management for MPI communication?
- SSH keys for inter-pod communication, or certificate-based PMIx authentication?
- Do secrets auto-rotate? What is the key lifecycle management strategy?
- Decision impact: Security architecture, operational overhead, compliance requirements

**Q20:** What is the audit trail for MPIJob operations?
- Regulatory requirements (SOX, HIPAA, GDPR) mandate audit logging
- What events do we log? Job creation, modification, deletion, access to logs?
- Decision impact: Compliance certification, log retention policies, storage costs

---

## Background & Strategic Fit

### Market Context: AI/ML Infrastructure Convergence

The AI/ML infrastructure market is undergoing a fundamental shift toward **platform convergence**. Gartner predicts 65% of organizations will unify HPC and cloud-native ML workloads on shared platforms by 2027 (up from 12% in 2023). Three key drivers are accelerating this trend:

**1. Model Scale Economics**

Modern AI workloads span a spectrum from cloud-native (microservices, auto-scaling) to HPC-style (tightly-coupled, high-performance interconnects). Large language models (LLMs) and multimodal models require 100+ GPUs with low-latency communication—traditional HPC schedulers lack cloud-native elasticity, while Kubernetes wasn't designed for tightly-coupled parallel compute. The market demands platforms that support both paradigms.

**2. Talent Consolidation**

Organizations struggle to hire ML engineers with both data science and HPC expertise. Customers need platforms that enable generalist data scientists to leverage HPC resources without becoming MPI experts. Abstraction layers (like KubeFlow Trainer V2) that hide infrastructure complexity are becoming competitive differentiators.

**3. FinOps Pressure**

Running separate HPC clusters and Kubernetes platforms creates 40-60% resource waste due to independent capacity planning. CFOs are demanding infrastructure consolidation to reduce capital expenditure and operational overhead. Unified platforms that share resource pools across workload types deliver immediate cost savings.

### Red Hat OpenShift AI Strategic Positioning

MPIJobs support is **critical** to Red Hat's "Hybrid AI Platform" strategy for three reasons:

**1. Open Source Differentiation**

KubeFlow Trainer V2 is the only vendor-neutral, portable training API for Kubernetes. Competitors (AWS SageMaker, Azure ML, Google Vertex AI) lock customers into proprietary SDKs. By leading with open standards, Red Hat enables workload portability—customers can develop on OpenShift AI and deploy to any Kubernetes environment.

**Competitive Advantage:** Customers value avoiding vendor lock-in. In Red Hat's FY24 Win/Loss analysis, "open standards" was cited as the #1 decision factor in 58% of OpenShift AI wins.

**2. Hybrid Cloud Reality**

IDC's 2024 Enterprise AI Survey shows 78% of enterprise AI workloads span on-premises and public cloud. MPI is table-stakes for on-premises GPU clusters—without it, OpenShift AI can't compete in hybrid scenarios. Customers with edge deployments, data residency requirements, or existing HPC investments will choose competitors if Red Hat doesn't support MPI.

**Market Risk:** Losing hybrid AI deals costs Red Hat access to financial services (35% of pipeline), healthcare (18%), and government (12%)—verticals with strict on-premises requirements.

**3. Ecosystem Leverage**

Intel, NVIDIA, AWS (EFA), and Azure (InfiniBand) are investing in containerized MPI to modernize HPC workloads. Red Hat can leverage this ecosystem innovation rather than building from scratch. Partnerships with hardware vendors (NVIDIA DGX systems, Intel oneAPI) and cloud providers (AWS Outposts, Azure Arc) become force multipliers.

**Strategic Opportunity:** Position OpenShift AI as the integration layer between legacy HPC investments and modern cloud-native tooling—Red Hat becomes the bridge, not a competitor to existing infrastructure.

### Competitive Landscape

**Google Anthos + Vertex AI:**
- Added MPI support for distributed training in Q2 2024
- Limited to Google Cloud and Anthos-managed clusters
- Red Hat Advantage: True hybrid cloud (any Kubernetes), open standards (not Google-specific APIs)

**NVIDIA DGX Cloud:**
- Offers MPI-based training as a managed service
- Attracts customers wanting turnkey solutions for GPU-intensive workloads
- Red Hat Risk: Losing GPU-rich customers to proprietary closed platforms
- Red Hat Response: Partner with NVIDIA—position OpenShift AI as the open alternative that runs on DGX systems

**Microsoft Azure ML + AKS:**
- Previewing MPIJob integration via KubeFlow (similar to Red Hat approach)
- Tightly integrated with Azure infrastructure (InfiniBand, Azure Batch)
- Red Hat Advantage: Multi-cloud portability, stronger enterprise Linux/container heritage
- Red Hat Risk: If Microsoft reaches GA first, Red Hat appears behind in innovation

**AWS SageMaker HyperPod:**
- Supports MPI via custom distributed training APIs
- Forces customers into SageMaker SDK (not portable)
- Red Hat Advantage: KubeFlow portability, hybrid cloud deployment
- Red Hat Challenge: AWS operational simplicity vs. Red Hat's "bring your own cluster" model

### Industry Trends: HPC-AI Convergence

**Convergence of HPC and AI Workloads:**

Traditional HPC (scientific computing, simulations) and AI (training, inference) are merging. Examples:
- **Weather forecasting:** Combine physics-based simulations (MPI) with ML-based predictions (TensorFlow)
- **Drug discovery:** Molecular dynamics simulations (MPI) feed data to generative models (PyTorch)
- **Financial modeling:** Monte Carlo simulations (MPI) integrated with risk prediction models (XGBoost)

**Infrastructure Implications:**

Organizations can't afford separate infrastructure stacks. They need platforms that support:
- Tightly-coupled MPI communication (InfiniBand, RDMA)
- Cloud-native orchestration (Kubernetes, auto-scaling)
- Unified observability (Prometheus, Grafana)
- Consistent RBAC and security (OpenShift)

**Red Hat's Opportunity:**

OpenShift AI is uniquely positioned to be the convergence platform:
- **Strong HPC heritage:** Red Hat Enterprise Linux dominates HPC (70% of Top500 supercomputers)
- **Cloud-native leadership:** OpenShift is the leading enterprise Kubernetes platform
- **Open standards commitment:** KubeFlow, Prometheus, Kubeflow Pipelines are all open source

### Customer Evidence: Why MPIJobs Matter

**From Red Hat's FY24 Enterprise AI Survey (n=240 customers):**

- **42% of enterprise POCs** cited "lack of HPC training support" as a blocker to OpenShift AI adoption
- **67% of customers** with existing HPC infrastructure said they would "consolidate onto OpenShift AI" if MPI support were available
- **35% reduction in infrastructure costs** was the most cited benefit of platform consolidation
- **3+ training job types** (PyTorch, TensorFlow, MPI) correlated with **85% higher renewal rates**

**Customer Quote (Tier-1 Financial Institution):**

> "We run quantitative risk models using MPI on a legacy HPC cluster that costs $1M/year to operate. We want to modernize onto OpenShift AI, but can't rewrite 10 years of MPI code. If Red Hat adds MPI support, we'll migrate 80% of our training workloads within 12 months."

**What's at Stake:**

Without MPIJobs, Red Hat loses credibility in regulated industries (FSI, healthcare, government) where on-premises HPC is non-negotiable. That's **35% of the OpenShift AI pipeline** at risk—approximately **$800M+ in TAM expansion** opportunity.

### Red Hat's Go-to-Market Strategy

**Positioning:**

"OpenShift AI: The Only Enterprise Platform That Runs HPC and Cloud-Native AI Workloads Side-by-Side"

**Value Proposition:**

- **For CFOs:** Reduce infrastructure costs by 30-40% by consolidating HPC and Kubernetes onto shared resource pools
- **For CIOs:** Accelerate AI adoption by removing barriers—data scientists use familiar tools (Jupyter, Python) to access HPC resources
- **For Data Science Leaders:** Unify team workflows—no more context switching between HPC schedulers and cloud platforms

**Sales Play:**

Target accounts with:
1. Existing HPC investments ($500K+ annual spend)
2. Active OpenShift deployments (licensed nodes > 100)
3. Regulatory/compliance requirements (on-prem GPU clusters)

**Vertical Focus:**

- **Financial Services:** Risk modeling, fraud detection, algorithmic trading
- **Healthcare/Life Sciences:** Drug discovery, genomics, medical imaging
- **Government/Defense:** Scientific research, climate modeling, cybersecurity

---

## Customer Considerations

### Enterprise Requirements

**1. Procurement and Licensing Implications**

**Challenge:** Customers with existing commercial MPI implementations (Intel MPI, IBM Spectrum MPI) face license portability questions when moving to containers.

**Considerations:**
- **Intel MPI:** Requires commercial license; not redistributable in container images without Intel partnership. Customers must provide their own base images.
- **IBM Spectrum MPI:** Licensed with IBM hardware (Power systems); redistribution restrictions apply.
- **OpenMPI:** Open source (BSD license); Red Hat can provide validated images without licensing concerns.

**Red Hat Action:**
- Provide pre-built, supported container images for **OpenMPI only** (versions 3.x, 4.x, 5.x)
- Document procedures for customers to build custom images with commercial MPI implementations
- Partner with Intel/IBM to provide guidance on license compliance for containerized MPI
- Legal review required: Ensure Red Hat indemnification doesn't cover third-party MPI implementations

**Customer Impact:**
- Procurement teams need clear guidance: "Which MPI licenses transfer to container environments?"
- Budget implications: Some customers may need to procure new Intel MPI licenses for containerized usage

---

**2. Compliance and Security Requirements**

**Challenge:** Regulated industries (FSI, healthcare, government) have strict security and audit requirements that must be addressed.

**Multi-Tenancy:**
- MPI implementations often use pod-to-pod SSH for process communication—potential security audit concern
- **Mitigation:** Provide NetworkPolicy templates for pod-to-pod communication isolation; document alternatives (PMIx with TLS)
- **Compliance Impact:** Security teams will audit SSH key management, rotation policies, and access controls

**Data Residency:**
- Large training datasets (TB-scale) distributed across worker nodes must comply with data residency laws (GDPR, CCPA)
- **Mitigation:** Support node selectors and pod affinity to enforce data locality (e.g., "all workers in EU availability zone")
- **Compliance Impact:** Legal teams need documentation on how MPIJobs respect data residency constraints

**Audit Logging:**
- SOX, HIPAA, FedRAMP require comprehensive audit trails for AI workloads
- **Required Events:** MPIJob creation, modification, deletion, log access, resource consumption
- **Integration:** OpenShift audit logs must capture all MPIJob lifecycle events
- **Compliance Impact:** Non-negotiable for regulated industries—MPIJob audit logging must meet certification standards

**FIPS Compliance:**
- Government customers (FedRAMP High, DoD IL5) require FIPS 140-2 validated cryptographic modules
- **Challenge:** Not all MPI implementations have FIPS-validated builds
- **Red Hat Action:** Validate OpenMPI with FIPS-enabled OpenSSL; document FIPS compliance status for each supported MPI version
- **Customer Impact:** Government procurement requires FIPS compliance—missing this blocks entire vertical

---

**3. Support Model and Scope Boundaries**

**Challenge:** MPI is a complex distributed systems technology; customers need clarity on what Red Hat supports vs. third-party vendors.

**Support Scope:**

| Component | Red Hat Support | Third-Party Support |
|-----------|-----------------|---------------------|
| MPIJob orchestration (KubeFlow Trainer V2) | ✅ Red Hat | N/A |
| OpenMPI runtime (Red Hat-provided images) | ✅ Red Hat | N/A |
| Intel MPI runtime | ❌ Intel Support | ✅ Intel |
| Custom user MPI code | ❌ Customer | ✅ Customer/SI Partners |
| InfiniBand/RDMA network configuration | ❌ Mellanox/NVIDIA | ✅ Hardware Vendor |
| Container image build errors | ✅ Red Hat (if using RHEL base) | ❌ Customer (custom images) |

**Customer Communication:**
- **Critical:** Document support boundaries in legal language (SLA exclusions, third-party dependencies)
- Technical Account Managers (TAMs) must be trained to triage: "Is this an OpenShift AI issue or an MPI runtime issue?"
- Knowledge Base articles for common scenarios: "Intel MPI hangs on initialization → Contact Intel Support"

**Field Readiness:**
- **Gap:** 60% of Red Hat Account Executives lack HPC background
- **Training Required:** 2-day enablement program: "MPI Fundamentals for Red Hat Sales"
- **Content:** What is MPI? When to use MPIJobs vs. PyTorchJobs? Demo: Running first MPIJob
- **Risk:** Sales teams overselling capabilities or misunderstanding customer requirements

---

**4. Migration and Adoption Strategy**

**Challenge:** Customers won't rip-and-replace existing HPC infrastructure—they need a gradual migration path.

**Hybrid Operation Period:**
- **Reality:** 6-18 months of running legacy HPC and OpenShift AI in parallel
- **Requirement:** "Bridge patterns" documentation—how to gradually shift workloads without disrupting production

**Migration Patterns:**

**Pattern 1: Pilot Project Approach**
1. Identify low-risk MPI workload (non-production, short-running)
2. Containerize and run on OpenShift AI as proof-of-concept
3. Compare performance, cost, operational overhead vs. legacy HPC
4. Iteratively migrate additional workloads based on lessons learned

**Pattern 2: Hybrid Orchestration**
1. Keep legacy HPC for long-running, mission-critical jobs (short-term)
2. Route new/experimental MPI workloads to OpenShift AI
3. Gradually deprecate HPC as confidence builds
4. Decommission legacy cluster after 12-18 months

**Pattern 3: Data Co-Location**
1. Use legacy HPC for data-intensive workloads (to avoid network egress)
2. Use OpenShift AI for compute-intensive workloads (elastic GPU access)
3. Eventually consolidate storage onto OpenShift AI data layer (S3, PVCs)

**Documentation Required:**
- **Reference Architectures:** "On-Prem HPC + Cloud OpenShift AI Hybrid Design"
- **Decision Tree:** "Should This Workload Stay on HPC or Move to MPIJobs?"
- **Case Studies:** Customer migration stories with timelines, costs, lessons learned

**Services Opportunity:**
- Red Hat Consulting: "HPC-to-OpenShift AI Migration Assessment"
- Partner Ecosystem: Systems integrators (Accenture, Deloitte) can deliver migration projects

---

**5. Cost Optimization and FinOps**

**Challenge:** MPI workloads are resource-intensive; customers will hit quota limits quickly and face unexpected costs.

**Cost Visibility:**
- **Requirement:** "This MPIJob consumed $X in compute" integrated with OpenShift cost management
- **Implementation:** Tag MPIJob pods with cost allocation labels (team, project, cost center)
- **Reporting:** Grafana dashboards showing per-job cost, cumulative team spend, forecast vs. budget

**Quota Management:**
- **Reality:** 16-worker MPIJob with 4 GPUs each = 64 GPUs (easily exceeds project quotas)
- **User Experience:** Clear error messages: "Requested 64 GPUs, but project quota is 32. Contact admin or reduce workers."
- **Admin Tools:** Cluster-wide visibility into quota utilization by team/project

**Auto-Scaling Challenges:**
- **Limitation:** MPI jobs require fixed worker counts (can't auto-scale during execution)
- **Customer Expectation Management:** Document that MPIJobs are "static scale" vs. PyTorchJobs (elastic)
- **Workaround:** Elastic Horovod (post-MVP) supports dynamic worker addition/removal

**Cost Optimization Best Practices:**
- **Spot Instances/Preemptible Nodes:** Can MPIJobs run on cheaper interruptible compute?
- **Trade-off:** Potential job preemption vs. 60-80% cost savings
- **Guidance:** "Use spot nodes for non-critical training; reserve on-demand for production"

**Chargeback Integration:**
- **Enterprise Requirement:** IT Finance teams need to bill back AI compute costs to business units
- **Integration:** OpenShift metering + MPIJob resource usage = per-team invoicing
- **Complexity:** Requires cross-functional collaboration (IT Finance, Platform Engineering, Data Science)

---

**6. Performance Expectations and Customer Communication**

**Challenge:** Customers expect containerized MPI to match bare-metal HPC performance—reality is 85-95% efficiency.

**Setting Expectations:**
- **Transparent Communication:** "MPIJobs on OpenShift AI achieve 85-95% of bare-metal performance due to container overhead and network virtualization."
- **Trade-offs:** Slight performance penalty in exchange for operational simplicity, elasticity, and cloud-native integration
- **Customer Segmentation:**
  - **Performance-Critical Customers** (99%+ efficiency required): Keep on bare-metal HPC; OpenShift AI not suitable
  - **Pragmatic Customers** (85-95% acceptable): Migrate to MPIJobs for operational benefits

**Benchmarking and Validation:**
- **Required:** Publish reference benchmarks: "Horovod MPI training on OpenShift AI vs. Bare-Metal Slurm"
- **Metrics:** Training time, GPU utilization, network bandwidth, cost per epoch
- **Transparency:** Document test conditions (hardware, network, MPI version, dataset) for reproducibility

**Performance Tuning Documentation:**
- **Content:** "Optimizing MPIJobs for InfiniBand Networks"
- **Topics:** MPI flags, NCCL tuning, RDMA configuration, GPU-Direct RDMA
- **Audience:** Advanced users and performance engineers

---

**7. Ecosystem and Partner Considerations**

**Hardware Vendor Partnerships:**
- **NVIDIA:** Validate MPIJobs on DGX systems; co-marketing "OpenShift AI on NVIDIA Infrastructure"
- **Intel:** Collaborate on Intel MPI containerization guidance; optimize for Xeon/Gaudi processors
- **AMD:** Validate MPIJobs with ROCm (AMD GPU stack); document configuration

**Cloud Provider Integration:**
- **AWS:** Support for EFA (Elastic Fabric Adapter) in MPIJobs on ROSA (Red Hat OpenShift on AWS)
- **Azure:** InfiniBand integration on ARO (Azure Red Hat OpenShift)
- **Google Cloud:** Multi-NIC support for GKE-based OpenShift AI deployments

**ISV and Framework Partners:**
- **Databricks:** Integrate MPIJobs with Databricks MLflow for experiment tracking
- **Weights & Biases:** Ensure W&B logging works seamlessly from distributed MPI ranks
- **Domino Data Lab:** Partner on enterprise MLOps platform integration

**Systems Integrators (SIs):**
- **Accenture, Deloitte, Cognizant:** Enable SIs to deliver "HPC Modernization on OpenShift AI" engagements
- **Training:** Provide SI partner training on MPIJobs architecture, best practices, troubleshooting
- **Incentives:** Co-sell agreements, reference architectures, joint customer engagements

---

## Technical Considerations

### 1. Integration and Deployment Architecture

**KubeFlow Trainer V2 Integration:**
- MPIJobs must be implemented as a `TrainJob` with `trainingFramework.kind: MPIJob` following KubeFlow Trainer V2 API standards
- Integration with KubeFlow Training Operator for lifecycle management (creation, monitoring, cleanup)
- CRD schema alignment with upstream KubeFlow for portability and community support

**Deployment Topology:**
- **Launcher Pod:** Single pod responsible for MPI process coordination and worker orchestration
- **Worker Pods:** N worker pods (configurable) running MPI rank processes
- **Communication:** Launcher discovers workers via Kubernetes Service/Headless Service for DNS-based resolution
- **Initialization:** Launcher spawns MPI processes on workers via SSH or PMIx protocol

**OpenShift AI Component Integration:**
- Integration with ODH Operator for MPIJob CRD installation and version management
- Dashboard backend API extensions for MPIJob CRUD operations
- SDK extensions (`kubeflow.training.TrainingClient`) for programmatic job management

**Namespace and RBAC:**
- MPIJobs scoped to Data Science Project namespaces
- RBAC roles: `mpijob-creator`, `mpijob-viewer`, `mpijob-editor` aligned with existing OpenShift AI roles
- Service accounts with minimal privileges (no cluster-admin for user-submitted jobs)

---

### 2. Networking and Inter-Process Communication (IPC)

**MPI Communication Mechanisms:**
- **SSH-based:** Traditional approach using pod-to-pod SSH for process spawning (requires SSH key management)
- **PMIx-based:** Modern process management interface with better security (TLS-based authentication)
- **Recommendation:** Support both; prefer PMIx for new deployments (better security posture)

**Network Requirements:**
- **TCP/IP:** Baseline support for standard Kubernetes networking (Calico, OVN-Kubernetes)
- **RDMA (InfiniBand/RoCE):** High-performance networking for low-latency MPI communication
  - Requires: RDMA-capable CNI plugin (e.g., Multus + SR-IOV, Macvlan)
  - Requires: Host device passthrough for InfiniBand HCAs (Host Channel Adapters)
  - Requires: Kernel modules (ib_core, rdma_cm) on OpenShift nodes

**NetworkPolicy Considerations:**
- MPIJob pods need east-west communication (launcher ↔ workers, worker ↔ worker)
- Default NetworkPolicy templates to allow MPI traffic while isolating from other namespaces
- Security: Prevent unauthorized pods from joining MPI communication groups

**DNS and Service Discovery:**
- Headless Service for worker pod discovery (predictable DNS names: `worker-0.mpijob.namespace.svc`)
- Launcher uses hostfile generation to map MPI ranks to pod IPs/DNS names
- Challenge: Handle pod rescheduling (IP changes); require stable network identities (StatefulSet pattern?)

**Bandwidth and Latency Optimization:**
- **GPU-Direct RDMA:** Enable direct GPU-to-GPU communication bypassing CPU (requires NVIDIA GPUDirect drivers)
- **NCCL Optimization:** Tune NCCL (NVIDIA Collective Communications Library) for multi-node scaling
- **MPI Tuning Flags:** Document recommended MPI flags for OpenShift environments (e.g., `--mca btl tcp,self` for TCP-only)

---

### 3. Resource Orchestration and Scheduling

**Gang Scheduling Requirement:**
- MPIJobs require **all pods to start simultaneously** (can't run with partial workers)
- If any worker pod fails to schedule, entire job must wait (not start some workers and hang)

**Scheduler Options:**
- **Option 1: Kueue Integration (Recommended)**
  - Kueue provides gang scheduling and resource quotas for batch workloads
  - Benefits: Unified queue management across MPIJobs, PyTorchJobs, TFJobs
  - Integration: Kueue `ClusterQueue` and `LocalQueue` with resource quotas
- **Option 2: Volcano Scheduler**
  - Alternative gang scheduler for Kubernetes
  - Benefits: Mature, widely used in HPC-on-Kubernetes projects
  - Trade-off: Additional component to install and maintain
- **Option 3: Best-Effort Scheduling (Not Recommended)**
  - Rely on default Kubernetes scheduler; accept risk of partial pod scheduling
  - Risk: Jobs hang indefinitely if some workers never schedule

**Recommendation:** Mandate Kueue for MPIJob support; document Kueue installation as prerequisite.

**Priority and Preemption:**
- Support Kubernetes `PriorityClass` for MPIJobs (high-priority safety models preempt low-priority experiments)
- Preemption handling: MPIJob should checkpoint progress before termination (requires application-level checkpointing)

**Resource Quotas:**
- MPIJobs respect OpenShift `ResourceQuota` at namespace level (CPU, memory, GPU limits)
- Validation: Reject job creation if total resource requests exceed quota
- User feedback: Clear error messages with quota usage and available capacity

**Topology-Aware Scheduling:**
- **Pod Affinity/Anti-Affinity:** Co-locate workers on same rack/availability zone for low-latency networking
- **Node Selectors:** Allow users to target specific node pools (e.g., InfiniBand-enabled nodes, GPU-rich nodes)
- **Topology Spread Constraints:** Distribute workers across failure domains for resilience (trade-off with network latency)

---

### 4. Hardware Acceleration and Runtime

**GPU Support:**
- MPIJob workers must support NVIDIA GPUs (`nvidia.com/gpu` resource requests)
- AMD GPUs (`amd.com/gpu`) and Intel GPUs (`intel.com/gpu`) for multi-vendor support
- GPU allocation: Each MPI rank can request 1 or more GPUs (e.g., 4 workers × 2 GPUs = 8-GPU training)

**NVIDIA-Specific Features:**
- **NCCL (NVIDIA Collective Communications Library):** MPI backend for multi-GPU/multi-node training
- **GPU-Direct RDMA:** Direct GPU-to-GPU communication over InfiniBand (bypasses CPU, reduces latency)
- **MIG (Multi-Instance GPU):** Partition A100/H100 GPUs into smaller instances (e.g., 1 A100 → 7 MIG instances)
  - Challenge: MPI rank-to-MIG instance mapping requires careful configuration

**MPI Runtime Selection:**
- Users specify MPI implementation via container image (e.g., `nvcr.io/nvidia/pytorch:24.01-py3` includes OpenMPI + NCCL)
- Red Hat provides validated base images with OpenMPI 4.x, 5.x on RHEL UBI (Universal Base Image)
- Custom MPI: Users can BYO (Bring Your Own) MPI implementation (Intel MPI, MPICH, custom builds)

**Container Runtime Considerations:**
- **cri-o:** Default OpenShift container runtime; must support device passthrough (GPUs, InfiniBand HCAs)
- **Privilege Requirements:** Some MPI configurations require `privileged` or `SYS_PTRACE` capabilities (security risk—document alternatives)

**Performance Optimization:**
- **CPU Pinning:** Bind MPI ranks to specific CPU cores (NUMA affinity) to reduce memory latency
- **Huge Pages:** Enable huge page support for large memory allocations (reduces TLB misses)
- **Kernel Parameters:** Tune sysctl settings for MPI workloads (e.g., `net.core.rmem_max` for socket buffer size)

---

### 5. Required Software and Data Services

**Data Input/Output:**
- **S3-Compatible Storage:** Integration with OpenShift Data Foundation (ODF), MinIO, AWS S3, Azure Blob
- **Persistent Volumes:** Support for RWX (ReadWriteMany) PVCs for shared training datasets across workers
- **Data Loading:** High-throughput data pipelines (avoid bottlenecks—workers shouldn't wait for data)

**Model Checkpointing:**
- **Shared Storage:** All workers must write checkpoints to shared PVC or S3 bucket
- **Checkpoint Frequency:** Configurable (e.g., every 1000 steps) to balance I/O overhead vs. recovery time
- **Failure Recovery:** If job fails, users can restart from last checkpoint (requires application-level support)

**Logging and Metrics:**
- **Log Aggregation:** Collect logs from launcher and all worker pods
- **Metrics Collection:** Prometheus integration for resource usage (CPU, memory, GPU, network)
- **Custom Metrics:** Support for user-defined metrics (training loss, accuracy) via Prometheus pushgateway

**Dependency Management:**
- **Container Images:** Users provide images with all dependencies (MPI runtime, training frameworks, libraries)
- **Image Registries:** Support for private registries (Red Hat Container Catalog, Quay, Harbor, Docker Hub)
- **Image Pull Policies:** Configurable (Always, IfNotPresent) with credential management via OpenShift Secrets

**Secrets and Configuration:**
- **Environment Variables:** Inject secrets (API keys, credentials) as env vars in MPIJob pods
- **ConfigMaps:** Mount configuration files (MPI hostfiles, training configs) as volumes
- **SSH Keys:** Launcher and workers share SSH keys for process spawning (requires Secret with key pair)

---

### 6. Operational and Non-Functional Requirements

**Scalability:**
- Support up to **64 workers per MPIJob** (128-256 GPUs for large-scale training)
- Handle **100+ concurrent MPIJobs** per OpenShift cluster (depends on cluster size)
- Control plane scalability: MPIJob operator should not become bottleneck

**Resilience and Fault Tolerance:**
- **Worker Pod Failures:** If a worker crashes, entire MPIJob fails (MPI characteristic—no automatic recovery)
- **Restart Policy:** Support job-level restart (recreate all pods) with configurable retry limit
- **Checkpointing:** Encourage users to implement application-level checkpointing for long-running jobs

**Performance:**
- **Job Startup Latency:** MPIJob should reach "Running" state within 2 minutes for 16-worker jobs (assuming images are pre-pulled)
- **Scheduling Delay:** Gang scheduling should not add >30 seconds overhead vs. individual pod scheduling
- **Network Initialization:** MPI process discovery and communication setup <1 minute

**Availability:**
- **Control Plane HA:** MPIJob operator should be highly available (multi-replica deployment)
- **Upgrade Safety:** Rolling upgrades of MPIJob operator should not disrupt running jobs

**Observability:**
- **Job Status:** Real-time status updates (Pending, Running, Succeeded, Failed) via Kubernetes API
- **Progress Tracking:** Integration with training frameworks for epoch/step progress (requires user instrumentation)
- **Historical Data:** Retain completed job metadata for 30 days (configurable) for analysis

---

### 7. Monitoring, Logging, and Debugging

**Monitoring:**
- **Prometheus Metrics:**
  - Job-level: creation time, duration, status, worker count
  - Pod-level: CPU/memory/GPU utilization per worker
  - MPI-level: communication bandwidth, latency, message rate (if available from MPI runtime)
- **Grafana Dashboards:**
  - Cluster-wide: "MPIJob Resource Usage by Project"
  - Job-specific: "MPIJob XYZ - Worker Performance"
  - Admin view: "MPIJob Scheduler Queue Depth"

**Logging Architecture:**
- **Log Collection:** FluentBit/Fluentd agents on nodes collect logs from MPIJob pods
- **Log Aggregation:** Logs sent to Elasticsearch, Loki, or CloudWatch (depending on deployment)
- **Log Indexing:** Index by job name, namespace, pod name, MPI rank for efficient querying
- **Log Retention:** Default 7 days (configurable); cost trade-off for long-term storage

**Structured Logging:**
- **Standardized Format:** JSON logs with fields: timestamp, rank, log level, message
- **MPI Rank Identification:** Each log line tagged with MPI rank (e.g., `"rank": 5`)
- **Error Correlation:** Aggregate errors from all workers to identify systemic failures (e.g., all ranks OOM)

**Debugging Tools:**
- **Interactive Shell:** Allow users to exec into failed launcher/worker pods for debugging (requires RBAC permissions)
- **Log Search:** UI provides full-text search across all worker logs (powered by Elasticsearch/Loki)
- **Event Viewer:** Kubernetes events timeline for MPIJob lifecycle (pod scheduling, failures, completions)

**Common Failure Modes and Diagnostics:**

| Failure Mode | Symptoms | Diagnostic Approach | Resolution |
|--------------|----------|---------------------|------------|
| **ImagePullBackOff** | Worker pods stuck in "Waiting" | Check image registry connectivity; validate image name/tag | Fix image path; add registry credentials |
| **OOMKilled** | Worker pod terminated with exit code 137 | Check pod logs for memory usage; review resource requests | Increase memory limit; optimize code |
| **MPI Initialization Timeout** | Launcher logs: "Timeout waiting for workers" | Check NetworkPolicy; verify DNS resolution | Fix NetworkPolicy; check Service configuration |
| **GPU Allocation Failure** | Worker pod pending with "Insufficient nvidia.com/gpu" | Check cluster GPU availability; review resource quotas | Add GPU nodes; increase quota |
| **SSH Connection Refused** | Launcher logs: "Connection refused on worker-3:22" | Check SSH service in worker containers; verify SSH keys | Fix Dockerfile; regenerate SSH keys |

---

### 8. Security and Access Control

**RBAC (Role-Based Access Control):**
- **Roles:**
  - `mpijob-creator`: Can create, delete MPIJobs in assigned namespaces
  - `mpijob-viewer`: Read-only access to MPIJob status and logs
  - `mpijob-editor`: Can create, modify, delete MPIJobs
  - `mpijob-admin`: Full control including RBAC management
- **Integration:** Leverage OpenShift AI Data Science Project roles (existing `edit`, `view`, `admin` roles extend to MPIJobs)

**NetworkPolicy:**
- **Ingress Rules:**
  - Allow launcher → worker communication (all ports, or specific SSH/PMIx ports)
  - Allow worker → worker communication (for MPI collective operations)
- **Egress Rules:**
  - Allow workers → S3 endpoints (for data loading)
  - Allow workers → internet (for pip installs, if needed)
- **Isolation:** Deny traffic from other namespaces (prevent lateral movement)

**Pod Security Standards:**
- **Preferred Level: Restricted**
  - No privileged containers, no host namespace sharing, no dangerous capabilities
  - Challenge: Some MPI configurations require `SYS_PTRACE` (for debugging), `IPC_LOCK` (for RDMA)
  - **Solution:** Document alternatives (e.g., PMIx instead of ptrace-based process control)
- **Acceptable Level: Baseline**
  - Allow specific capabilities on a case-by-case basis (with security review)
- **Avoid: Privileged**
  - Last resort for legacy MPI applications; requires security exception approval

**Secrets Management:**
- **SSH Keys:** Launcher and workers share SSH key pair (stored in Kubernetes Secret)
  - Key rotation: Manual (users regenerate keys periodically) or automated (Vault integration, future enhancement)
- **S3 Credentials:** Data connection secrets mounted as environment variables or files
- **Registry Credentials:** Pull secrets attached to MPIJob ServiceAccount

**Container Image Security:**
- **Image Scanning:** Recommend scanning user-provided images for vulnerabilities (Clair, Trivy)
- **Trusted Registries:** Encourage use of Red Hat Container Catalog, approved Quay repositories
- **Base Image Compliance:** Red Hat-provided MPI images based on RHEL UBI (FIPS-compliant, CVE patching)

**Audit Logging:**
- **Events to Log:**
  - MPIJob creation (who, when, resource requests, image)
  - MPIJob modification (parameter changes, user identity)
  - MPIJob deletion (who initiated, completion status)
  - Log access (who viewed logs, which workers)
- **Log Format:** Structured JSON with OpenShift audit log schema
- **Retention:** 90 days minimum (compliance requirement for regulated industries)

**Data Encryption:**
- **In-Transit:** MPI communication over TLS (if using PMIx) or SSH tunnel (if using SSH-based MPI)
- **At-Rest:** Training data encrypted in S3 (S3 server-side encryption), PVC encryption (OpenShift Data Foundation encryption)

---

## Summary

This RFE proposes adding MPIJobs support to Red Hat OpenShift AI via KubeFlow Trainer V2, enabling customers to run high-performance distributed training workloads using industry-standard MPI protocols. The feature addresses a critical market gap—customers with existing HPC investments currently face a forced choice between modernizing onto OpenShift AI (losing MPI capabilities) or maintaining separate infrastructure (incurring 40-60% cost overhead).

By delivering unified MPI and cloud-native training support through a single platform, Red Hat positions OpenShift AI as the only enterprise Kubernetes solution that truly supports hybrid AI workloads. This unlocks $800M+ in TAM expansion, strengthens Red Hat's competitive position against AWS, Azure, and Google, and removes a top-3 barrier to enterprise AI adoption.

The feature must deliver on three core user promises:

1. **Self-Service Simplicity:** Data scientists create MPIJobs in under 5 minutes via UI/SDK without writing YAML or understanding Kubernetes internals
2. **Unified Observability:** MPIJobs appear alongside PyTorchJobs and TFJobs with consistent monitoring, logging, and debugging experiences
3. **Enterprise-Grade Operations:** Administrators gain centralized RBAC, quota management, audit logging, and cost tracking across all distributed training workloads

Success will be measured by customer adoption (25% of OpenShift AI customers running MPIJobs within 12 months), operational efficiency (time-to-first-distributed-training reduced from 14 days to <2 hours), and cost reduction (30% infrastructure savings for customers consolidating HPC and Kubernetes platforms).

This RFE provides a comprehensive blueprint for delivering MPIJobs support—from product strategy and user experience design to technical architecture and customer considerations. The next steps involve cross-functional alignment (Product, Engineering, UX, Documentation), technical feasibility validation, and roadmap prioritization to bring this critical capability to market.
