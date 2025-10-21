# MPIJob Support for Red Hat OpenShift AI

**Feature Overview:**

MPIJob support in Red Hat OpenShift AI enables enterprise data science teams to execute distributed training workloads using the industry-standard Message Passing Interface (MPI) protocol. This capability is critical for customers running large-scale deep learning models that cannot fit on single nodes or require massive computational resources to train in reasonable timeframes. By integrating KubeFlow Trainer V2's latest MPI capabilities, this feature provides data scientists, MLOps engineers, and administrators with unified observability and management for MPIJobs alongside other KubeFlow Trainer job types through CLI, SDK, and ODH Dashboard UI interfaces.

**Market Context:** Our customers are hitting walls with single-node training—models that take weeks to train, can't process datasets at scale, and force them to compromise on model quality. The competition (Azure ML, Google Vertex AI, AWS SageMaker) all provide native distributed training capabilities. Without MPIJob support, we're losing deals to hyperscalers and forcing customers into fragmented toolchains. This feature closes a critical competitive gap while leveraging our hybrid cloud differentiation.

**Business Impact:**
- **Time-to-Model Reduction**: Customers report 5-10x faster training times with distributed workloads
- **Model Quality Improvement**: Enables training on larger datasets and more complex architectures
- **TCO Optimization**: Better resource utilization across cluster infrastructure
- **Competitive Parity**: Closes feature gap with major cloud AI platforms

---

## Goals

**For Data Scientists:**
- **Today**: Must export models to external distributed training systems (Horovod standalone, Ray, custom scripts), losing OpenShift AI observability and workflow integration
- **With MPIJobs**: Launch distributed training directly from OpenShift AI with unified monitoring, logging, and artifact management
- **Outcome**: 40% reduction in training job setup time, seamless integration with existing pipelines

**For MLOps Engineers:**
- **Today**: Maintain separate infrastructure and tooling for distributed vs. single-node training, increasing operational complexity
- **With MPIJobs**: Unified interface for all KubeFlow Trainer job types through CLI, SDK, and Dashboard
- **Outcome**: Single pane of glass for all training workloads, reduced operational overhead

**For Platform Administrators:**
- **Today**: Custom integrations and security policies for multiple training frameworks
- **With MPIJobs**: Native OpenShift AI support with RBAC, quota management, and audit compliance
- **Outcome**: Consistent governance across all training job types

**Success Metrics:**
- **Adoption**: 30% of active users launch at least one MPIJob within 90 days of GA
- **Performance**: Distributed training jobs show measurable speedup (>2x) compared to single-node baseline
- **Integration**: 95% feature parity with other KubeFlow Trainer job types (observability, CLI, SDK, UI)
- **Customer Satisfaction**: NPS score increase of 10+ points among users running distributed workloads

---

## Out of Scope

### Excluded Job Types
- **Spark Jobs**: Different execution model, addressed by separate Spark operator
- **Dask Distributed**: Python-specific, lower priority based on customer demand data
- **Ray Jobs**: Emerging framework, evaluate for future roadmap post-MPIJobs GA

### Framework-Specific Features
- **Framework Optimizations**: No custom configurations for specific ML frameworks (TensorFlow, PyTorch, etc.)—users provide their own containers
- **Hyperparameter Tuning**: Handled by separate KubeFlow Katib integration, not MPIJob scope
- **Auto-scaling**: Dynamic worker scaling not in MVP—fixed topology only

### Use Cases
- **Non-ML Workloads**: MPIJobs focused on ML training, not general-purpose HPC or scientific computing
- **Streaming/Online Training**: Batch training only, no real-time model updates
- **Federated Learning**: Cross-cluster distributed training out of scope

### Personas
- **Hobbyist/Individual Users**: Focus on enterprise teams with multi-node infrastructure needs
- **Edge Deployments**: Distributed training assumes datacenter/cloud infrastructure

### Migration Tools
- **Automated Migration**: No conversion tools from other distributed training frameworks (Horovod standalone, DeepSpeed, etc.)
- **Job Translation**: Users must manually adapt existing distributed training code to MPIJob CRD format

**Rationale**: Customer data shows 82% of distributed training demand centers on MPI protocol. We need MVP velocity to close competitive gap—scope expansion comes post-GA based on adoption metrics.

---

## Requirements

### MVP Requirements (Must-Have for GA)

#### CLI Requirements
- **MPIJob-001**: CLI must support create, delete, list, and describe operations for MPIJob resources
- **MPIJob-002**: CLI must provide real-time job status monitoring including worker pod health and completion state
- **MPIJob-003**: CLI must enable log retrieval from both launcher and worker pods
- **MPIJob-004**: CLI must provide feature parity with existing TFJob and PyTorchJob CLI capabilities
- **MPIJob-005**: CLI must support YAML-based job specification with schema validation

#### SDK Requirements
- **MPIJob-006**: Python SDK must provide full CRUD operations for MPIJob resources
- **MPIJob-007**: SDK must use strongly-typed MPIJob specifications with validation
- **MPIJob-008**: SDK must support asynchronous job submission and status polling
- **MPIJob-009**: SDK must maintain API consistency with other KubeFlow Trainer job types (PyTorchJob, TFJob)
- **MPIJob-010**: SDK must include comprehensive API reference documentation with distributed training examples

#### UI Requirements (ODH Dashboard)
- **MPIJob-011**: Dashboard must provide form-based MPIJob creation with worker topology configuration
- **MPIJob-012**: MPIJobs must appear in unified job list view alongside other Trainer jobs with type filtering
- **MPIJob-013**: Job detail view must display worker pod status, resource utilization, and logs
- **MPIJob-014**: UI must maintain visual and interaction consistency with existing job type management
- **MPIJob-015**: UI must provide progressive disclosure pattern—essential fields visible, advanced MPI options collapsed

#### Observability Requirements
- **MPIJob-016**: System must collect and expose metrics for job duration, worker pod lifecycle, and resource consumption (CPU, memory, GPU)
- **MPIJob-017**: System must provide centralized log aggregation from all MPIJob pods (launcher and workers)
- **MPIJob-018**: System must generate Kubernetes events for all job state transitions
- **MPIJob-019**: Metrics must be exposed to OpenShift monitoring stack for unified observability
- **MPIJob-020**: UI must provide multi-level status visibility: job-level, worker aggregate, and individual worker detail

#### Integration Requirements
- **MPIJob-021**: MPIJobs must respect OpenShift RBAC with namespace-level permissions
- **MPIJob-022**: MPIJobs must honor OpenShift resource quotas and limit ranges
- **MPIJob-023**: System must support network policies for restricted network environments
- **MPIJob-024**: MPIJobs must support image pulls from OpenShift internal and external container registries
- **MPIJob-025**: MPIJobs must integrate with OpenShift authentication and service account mechanisms

#### KubeFlow Trainer V2 Integration
- **MPIJob-026**: Implementation must use KubeFlow Training Operator V2 API for MPIJob management
- **MPIJob-027**: MPIJob CRDs must align with KubeFlow Training Operator upstream specifications
- **MPIJob-028**: System must maintain compatibility with future KubeFlow Trainer V2 updates

**Dependency**: KubeFlow Trainer V2 support must be implemented in Red Hat OpenShift AI before MPIJob work begins.

### Nice-to-Have Requirements (Post-MVP)

#### Advanced CLI Features
- **MPIJob-029**: Support for job templates to save and reuse MPIJob configurations
- **MPIJob-030**: Bulk operations to manage multiple jobs simultaneously
- **MPIJob-031**: Job comparison capabilities to analyze configurations and results across runs

#### Enhanced SDK Features
- **MPIJob-032**: Job orchestration primitives to chain MPIJobs with other pipeline steps
- **MPIJob-033**: Automatic retry logic with exponential backoff
- **MPIJob-034**: Webhook callbacks for job state changes

#### Advanced UI Features
- **MPIJob-035**: Topology visualization showing worker node distribution
- **MPIJob-036**: Resource recommendations based on historical training data
- **MPIJob-037**: Cost estimation preview before job submission
- **MPIJob-038**: UI-based template creation and team sharing

#### Observability Enhancements
- **MPIJob-039**: Distributed tracing for cross-pod communication visibility
- **MPIJob-040**: Performance profiling to identify MPI communication bottlenecks
- **MPIJob-041**: Predictive alerts for early warning of job failures
- **MPIJob-042**: Integration with training metrics (loss, accuracy) in job monitoring views

**Prioritization Rationale**: MVP focuses on feature parity with existing job types—this is critical for customer adoption. Nice-to-have items require 6+ months of usage data to validate ROI.

---

## Done - Acceptance Criteria

### Functional Acceptance Criteria

**AC-001: Job Creation via Dashboard UI**
- GIVEN a data scientist accessing the ODH Dashboard
- WHEN they navigate to "Create Training Job" and select "MPIJob"
- THEN they can successfully configure and launch an MPIJob with worker count, resources, training image, and command
- AND the job appears in the unified job list with correct status

**AC-002: Job Creation via CLI**
- GIVEN an MLOps engineer with OpenShift AI CLI installed
- WHEN they execute `odh training create mpijob --config job-spec.yaml`
- THEN the MPIJob is created successfully
- AND they can retrieve status with `odh training describe <job-name>`
- AND they can access logs with `odh training logs <job-name>`

**AC-003: Job Creation via SDK**
- GIVEN a Python environment with OpenShift AI SDK installed
- WHEN a developer executes:
  ```python
  from odh.training import MPIJob
  job = MPIJob(name="test-training", workers=4, image="training:latest")
  job.create()
  ```
- THEN the MPIJob is created in the cluster
- AND the developer can retrieve status programmatically

**AC-004: Unified Observability**
- GIVEN multiple training jobs running (MPIJob, PyTorchJob, TFJob)
- WHEN a user views the job list in Dashboard
- THEN all job types appear in a single unified view with consistent status indicators
- AND users can filter by job type
- AND all jobs provide comparable monitoring capabilities (logs, metrics, events)

**AC-005: Worker Status Visibility**
- GIVEN a running MPIJob with 8 workers
- WHEN a user views the job detail page
- THEN they can see aggregated worker status (e.g., "7/8 Running, 1 Pending")
- AND they can drill down to individual worker pod status
- AND they can access logs from any specific worker

**AC-006: Resource Quota Enforcement**
- GIVEN a namespace with resource quotas configured
- WHEN a user attempts to create an MPIJob exceeding quota limits
- THEN the job creation fails with a clear error message
- AND the error indicates the quota constraint and current usage
- AND the user receives actionable guidance on resolving the issue

**AC-007: RBAC Permission Enforcement**
- GIVEN a user without MPIJob creation permissions in a namespace
- WHEN they attempt to create an MPIJob
- THEN the operation is denied with appropriate authorization error
- AND a user with appropriate permissions can successfully create MPIJobs

**AC-008: Job Lifecycle Management**
- GIVEN an MPIJob in any state
- WHEN a user with appropriate permissions deletes the job
- THEN the job and all associated resources (launcher, workers) are cleaned up
- AND no orphaned pods or resources remain in the cluster

### User Experience Acceptance Criteria

**AC-UX-001: Discoverability**
- WHEN a data scientist views the "Create Training Job" interface
- THEN MPIJob appears as an option with clear differentiation from other job types
- AND a comparison guide helps users select the appropriate job type
- AND the interface explains "Use MPIJob when: distributed training, multi-node, MPI-based frameworks"

**AC-UX-002: Quick Start Success**
- GIVEN a data scientist with no prior MPI experience
- WHEN they follow the "Getting Started with MPIJobs" guide
- THEN they can successfully launch their first distributed training job within 30 minutes
- AND the guide requires no Kubernetes knowledge beyond basic concepts

**AC-UX-003: Configuration Validation**
- GIVEN a user configuring an MPIJob
- WHEN they specify invalid or unavailable resources
- THEN validation errors appear inline before job submission
- AND errors explain the issue and suggest corrections
- EXAMPLE: "Requested 8 GPUs but only 4 available. Reduce worker count or adjust GPU allocation."

**AC-UX-004: Status Transparency**
- AT ALL TIMES during MPIJob lifecycle
- WHEN a user checks job status
- THEN they can see: overall job phase, individual worker states with counts, current step in lifecycle, time in current state

**AC-UX-005: Accessibility Compliance**
- ALL Dashboard UI elements related to MPIJob management
- MUST meet WCAG 2.1 Level AA standards
- INCLUDING keyboard navigation, screen reader support, and sufficient color contrast
- AND status must be communicated through multiple channels (color, icons, text labels)

**AC-UX-006: Terminology Consistency**
- ACROSS all documentation, UI, CLI, and SDK
- WHEN referring to MPI concepts
- THEN consistent terminology is used: "Worker" (not "replica"/"node"), "Launcher" (not "master"), "Slots per worker" (not "processes")

**AC-UX-007: Error Recovery**
- GIVEN a failed MPIJob
- WHEN a user views the failure details
- THEN the interface offers recovery options: "Retry with same configuration", "Clone and modify", "View detailed logs"
- AND previous job attempts are preserved for comparison

### Technical Acceptance Criteria

**AC-TECH-001: KubeFlow Training Operator V2 Integration**
- MPIJob implementation must use Training Operator V2 API
- MPIJob CRD must align with upstream KubeFlow specifications
- Changes to Training Operator API must not break MPIJob functionality

**AC-TECH-002: Multi-Tenancy Isolation**
- MPIJobs in different namespaces must be fully isolated
- Network traffic between MPIJob pods must respect namespace network policies
- Resource consumption must be tracked per namespace for chargeback

**AC-TECH-003: Performance Baseline**
- Distributed training using MPIJob must show measurable speedup (minimum 2x) compared to single-node baseline for reference workloads
- Job submission latency must be < 5 seconds under normal cluster load
- Status updates must reflect in UI within 10 seconds of state changes

**AC-TECH-004: Audit Logging**
- All MPIJob operations (create, delete, modify) must generate audit events
- Audit logs must include: user identity, timestamp, operation type, job name, namespace
- Audit events must integrate with OpenShift audit logging infrastructure

**AC-TECH-005: Container Image Support**
- MPIJobs must support standard container images with MPI installed
- Common MPI implementations must be validated: OpenMPI, Intel MPI
- Documentation must provide guidelines for building MPI-compatible training images

---

## Use Cases - User Experience & Workflow

### Use Case 1: Data Scientist - First Distributed Training Job

**Actor**: Maria, Data Scientist with 3 years ML experience, limited Kubernetes knowledge

**Goal**: Scale single-node PyTorch training to 8 GPUs across 4 nodes using Horovod

**Preconditions**:
- Maria has a working single-node training script
- OpenShift AI environment with GPU nodes available
- Maria has access to ODH Dashboard

**Main Success Scenario**:

1. Maria logs into ODH Dashboard and navigates to "Training Jobs"
2. She clicks "Create Training Job" and sees options: PyTorchJob, TFJob, MPIJob
3. She clicks "Which job type should I use?" and sees a decision guide:
   - "Use MPIJob for: Multi-node distributed training with MPI-based frameworks (Horovod, DeepSpeed with MPI backend)"
4. Maria selects "MPIJob" and sees a configuration form with progressive disclosure:
   - **Essential Fields (visible)**:
     - Job Name: "sentiment-model-distributed"
     - Training Image: "myregistry.com/training:horovod-latest"
     - Number of Workers: 4
     - Resources per Worker: 2 GPUs, 16 GB RAM, 8 CPUs
     - Command: `["python", "/workspace/train.py"]`
   - **Advanced Options (collapsed)**: Slots per worker, MPI implementation, environment variables
5. Maria fills in the essential fields, leaving advanced options at defaults
6. She clicks "Create Job" and sees validation: "Resources requested: 8 GPUs, 64 GB RAM—Available: 12 GPUs, 128 GB RAM ✓"
7. Job creation succeeds, Maria is redirected to job detail view
8. She sees job status: "Initializing - Starting launcher pod (1/5 pods ready)"
9. Status updates to "Running - 4/4 workers training" with a visual topology showing worker health
10. Maria clicks on Worker-2 to view logs, sees Horovod successfully initialized 8 ranks
11. After 45 minutes, status updates to "Succeeded"
12. Maria clicks "Download Model" to retrieve trained artifacts

**Alternative Flow 1 - Resource Constraint**:
- At step 6, validation shows: "Requested 8 GPUs but only 4 available. Estimated wait time: 25 minutes"
- Maria adjusts to 2 workers (4 GPUs total) and successfully creates job

**Alternative Flow 2 - Configuration Error**:
- At step 10, Maria sees status "Failed - Worker-2 crashed"
- She clicks on Worker-2, sees error: "MPI initialization failed: No network connectivity between workers"
- Error message includes link: "Troubleshooting MPI network issues →"
- Documentation guides her to check network policies

**Postconditions**:
- MPIJob completes successfully and Maria's model is trained
- Job execution is logged for audit purposes
- Maria can compare this job with previous single-node runs in Dashboard

**UI Mockup Notes**:
- Progressive disclosure: Simple interface by default, "Advanced Options" expander for power users
- Validation feedback: Real-time resource availability checks
- Visual topology: Grid/map showing worker health with color coding

---

### Use Case 2: MLOps Engineer - Automated Pipeline Integration

**Actor**: James, MLOps Engineer responsible for CI/CD pipelines

**Goal**: Integrate MPIJob into automated training pipeline triggered by Git commits

**Preconditions**:
- James has existing Tekton pipeline for single-node training
- Pipeline includes: data preprocessing, training, model evaluation, deployment
- James has OpenShift AI SDK installed in pipeline environment

**Main Success Scenario**:

1. James reviews MPIJob SDK documentation and identifies required API:
   ```python
   from odh.training import MPIJob, ResourceSpec
   ```
2. He modifies his pipeline script to create MPIJob programmatically:
   ```python
   # Pipeline Task: Distributed Training
   job = MPIJob(
       name=f"training-{git_commit_sha[:7]}",
       namespace="ml-team",
       workers=4,
       slots_per_worker=2,
       image=f"registry.com/models:build-{git_commit_sha}",
       command=["horovodrun", "-np", "8", "python", "train.py"],
       launcher_resources=ResourceSpec(cpu="2", memory="8Gi"),
       worker_resources=ResourceSpec(cpu="8", memory="32Gi", gpus=2),
       env={"DATASET_PATH": "/data/training-set", "MODEL_OUTPUT": "/models"}
   )
   ```
3. James adds error handling:
   ```python
   try:
       job.create(wait=False)  # Non-blocking submission
       logger.info(f"MPIJob {job.name} created successfully")
   except ResourceQuotaExceeded as e:
       logger.error(f"Quota exceeded: {e.details}")
       raise PipelineFailure("Insufficient resources for training")
   ```
4. He implements status polling loop:
   ```python
   while job.status != "Succeeded":
       if job.status == "Failed":
           logs = job.get_logs(component="launcher")
           raise PipelineFailure(f"Training failed: {logs}")
       time.sleep(30)
   ```
5. James configures webhook callback for job completion:
   ```python
   job.add_webhook(url="https://pipeline.com/training-complete", event="Succeeded")
   ```
6. He tests pipeline locally using dry-run mode:
   ```python
   validation_errors = job.validate()
   if validation_errors:
       logger.error(f"Invalid configuration: {validation_errors}")
   ```
7. James commits pipeline changes, triggering a test run
8. Pipeline executes: builds container → creates MPIJob → polls status → job succeeds → downloads artifacts → runs evaluation
9. James verifies in Dashboard that MPIJob was created by service account, shows correct ownership

**Alternative Flow 1 - Job Failure**:
- At step 8, job fails due to out-of-memory error in Worker-3
- Pipeline captures failure, attaches logs to pipeline run
- James receives alert with link to failed MPIJob in Dashboard
- He debugs, increases memory allocation in pipeline script, re-runs

**Alternative Flow 2 - CLI-Based Automation**:
- Instead of Python SDK, James uses CLI in shell script:
   ```bash
   odh training create mpijob \
     --name training-$(git rev-parse --short HEAD) \
     --config mpijob-template.yaml \
     --set workers=4 \
     --set image="registry.com/models:build-${GIT_COMMIT}" \
     --wait --timeout 3600 \
     --output-artifacts ./models/
   ```

**Postconditions**:
- MPIJob successfully integrates into automated CI/CD pipeline
- Training jobs are created programmatically with consistent configuration
- Audit trail shows service account as job creator

**Integration Points**:
- Python SDK with type hints and validation
- CLI with scriptable output (JSON for parsing)
- Webhook support for event-driven workflows

---

### Use Case 3: Administrator - Multi-Tenant Resource Management

**Actor**: Sarah, OpenShift AI Platform Administrator

**Goal**: Configure MPIJob quotas and permissions for multiple teams sharing infrastructure

**Preconditions**:
- Sarah has cluster-admin privileges
- Multiple teams (nlp-team, cv-team, research-team) use shared OpenShift AI environment
- Teams have different resource priorities and budgets

**Main Success Scenario**:

1. Sarah accesses OpenShift Admin Console and navigates to "Resource Quotas"
2. She creates namespace-specific quotas for MPIJob resources:
   ```yaml
   # nlp-team namespace
   apiVersion: v1
   kind: ResourceQuota
   metadata:
     name: mpijob-quota
     namespace: nlp-team
   spec:
     hard:
       requests.nvidia.com/gpu: "16"  # Max 16 GPUs total
       count/mpijobs.kubeflow.org: "5"  # Max 5 concurrent jobs
   ```
3. Sarah configures RBAC for MPIJob creation using role templates:
   - nlp-team: Full MPIJob creation rights
   - cv-team: Full MPIJob creation rights
   - research-team: Read-only (cannot create, can view)
4. She sets up network policies to isolate MPIJob traffic between namespaces:
   ```yaml
   # Allow MPI communication only within namespace
   apiVersion: networking.k8s.io/v1
   kind: NetworkPolicy
   metadata:
     name: mpijob-isolation
     namespace: nlp-team
   spec:
     podSelector:
       matchLabels:
         training.kubeflow.org/job-role: worker
     ingress:
     - from:
       - podSelector:
           matchLabels:
             training.kubeflow.org/job-name: "{{.JobName}}"
   ```
5. Sarah creates monitoring dashboard for MPIJob resource usage:
   - Panel 1: GPU utilization by namespace
   - Panel 2: Running MPIJobs per team
   - Panel 3: Job success/failure rate
   - Panel 4: Resource quota consumption
6. She documents team-specific guidelines:
   - "nlp-team: Max 4 workers per job, contact admin for exceptions"
   - "cv-team: Max 8 workers per job, GPU quota resets monthly"
7. Sarah receives alert: "nlp-team quota 90% consumed—15/16 GPUs in use"
8. She investigates in Dashboard: views all MPIJobs in nlp-team namespace, identifies long-running job
9. Sarah contacts nlp-team lead, recommends optimization or job termination

**Alternative Flow 1 - Quota Violation Troubleshooting**:
- Team member reports: "I can't create MPIJob, but quota should be available"
- Sarah checks resource quotas in namespace: 5/5 MPIJobs running
- She identifies 2 jobs in "Failed" state still counting against quota
- Sarah deletes failed jobs, quota freed, team member can proceed

**Alternative Flow 2 - Security Incident Response**:
- Security scanner detects suspicious MPIJob in research-team namespace
- Sarah reviews audit logs: identifies unauthorized service account created MPIJob
- She terminates job immediately, revokes service account permissions
- Audit trail preserved for security investigation

**Postconditions**:
- Multi-tenant environment with isolated, fair resource allocation
- Teams operate within quotas with clear visibility
- Administrator has operational monitoring and control

**Admin Interfaces**:
- OpenShift Admin Console for quota/RBAC configuration
- Monitoring dashboards (Prometheus/Grafana) for resource visibility
- Audit logs for compliance and security

---

### Use Case 4: Data Scientist - Troubleshooting Failed Job

**Actor**: Alex, Data Scientist with moderate distributed training experience

**Goal**: Diagnose and resolve MPIJob failure

**Preconditions**:
- Alex has created MPIJob that immediately fails
- Job status shows "Failed" after 2 minutes

**Main Success Scenario**:

1. Alex opens ODH Dashboard and navigates to his MPIJob "nlp-model-training"
2. Job detail page shows status: "Failed - 1/8 workers failed"
3. Visual topology highlights Worker-5 in red with error icon
4. Alex clicks on Worker-5, sees condensed error message:
   - "Container failed with exit code 137 (OOMKilled)"
   - "Memory limit (16 GB) exceeded"
5. Error card provides actionable suggestions:
   - "Increase memory allocation for workers"
   - "Reduce batch size in training script"
   - "View full logs for details →"
6. Alex clicks "View full logs", sees:
   ```
   [Worker-5] INFO: Loading training data...
   [Worker-5] INFO: Batch size: 512, Workers: 8
   [Worker-5] ERROR: Failed to allocate 18.5 GB for model parameters
   [Worker-5] FATAL: Out of memory
   ```
7. Alex identifies the issue: Model + data exceeds 16 GB memory limit
8. He clicks "Clone and Modify" button on job detail page
9. Pre-filled form appears with previous configuration, Alex changes:
   - Worker memory: 16 GB → 32 GB
10. Alex clicks "Create Job" to resubmit with updated configuration
11. New job runs successfully, all workers healthy
12. Alex compares both jobs side-by-side in Dashboard to verify fix

**Alternative Flow 1 - MPI Communication Failure**:
- Error message: "MPI initialization timeout - workers cannot communicate"
- Suggestions:
  - "Check network policies allow traffic between workers"
  - "Verify MPI implementation matches across all workers"
  - "Review MPI configuration guide →"
- Alex follows network troubleshooting guide, identifies restrictive NetworkPolicy

**Alternative Flow 2 - Image Pull Failure**:
- Error message: "Failed to pull image 'myregistry.com/training:v2.0'"
- Suggestions:
  - "Verify image exists in registry"
  - "Check image pull secrets are configured"
  - "Review image naming guide →"
- Alex realizes typo in image tag, corrects and resubmits

**Postconditions**:
- Alex successfully diagnoses and resolves job failure
- Knowledge gained helps prevent similar issues in future
- Failed job retained for comparison and learning

**UX Patterns**:
- **Error Surfacing**: Most relevant error shown immediately, no navigation required
- **Actionable Guidance**: Suggestions tailored to specific error type
- **Easy Retry**: Clone-and-modify workflow preserves configuration while allowing fixes
- **Learning Support**: Links to documentation for deeper understanding

---

### Use Case 5: Cross-Channel Workflow - Dashboard to CLI

**Actor**: Priya, Data Scientist who uses both UI and CLI

**Goal**: Create job in Dashboard, monitor/debug via CLI in terminal

**Preconditions**:
- Priya has access to ODH Dashboard and CLI installed locally
- She prefers UI for job creation, CLI for monitoring during active development

**Main Success Scenario**:

1. Priya creates MPIJob "experiment-42" in Dashboard with 4 workers
2. She opens terminal and authenticates CLI: `odh login --server https://openshift-ai.example.com`
3. Priya lists her jobs:
   ```bash
   $ odh training list --namespace my-project
   NAME           TYPE     STATUS    WORKERS  CREATED
   experiment-42  mpijob   Running   4/4      2m ago
   experiment-41  mpijob   Succeeded 8/8      1h ago
   ```
4. She watches job status in real-time:
   ```bash
   $ odh training describe experiment-42 --watch
   Name:       experiment-42
   Type:       MPIJob
   Status:     Running
   Workers:    4/4 running
   Duration:   3m 45s

   Worker Status:
   - Launcher: Running (45s)
   - Worker-1: Running (30s)
   - Worker-2: Running (30s)
   - Worker-3: Running (28s)
   - Worker-4: Running (30s)
   ```
5. Priya tails logs from specific worker:
   ```bash
   $ odh training logs experiment-42 --worker 2 --follow
   [Worker-2] Epoch 1/10 - Loss: 0.523
   [Worker-2] Epoch 2/10 - Loss: 0.412
   ```
6. She notices Worker-3 has higher loss than others, retrieves those logs:
   ```bash
   $ odh training logs experiment-42 --worker 3 --follow
   ```
7. After job completes, Priya exports job configuration for reuse:
   ```bash
   $ odh training get experiment-42 --output yaml > experiment-42-config.yaml
   ```
8. She modifies YAML for next experiment, creates new job via CLI:
   ```bash
   $ odh training create mpijob --config experiment-42-config.yaml --name experiment-43
   ```
9. Priya switches back to Dashboard to view visual metrics and download model artifacts

**Alternative Flow - Debugging with CLI then Viewing in UI**:
- CLI shows cryptic error, Priya wants visual context
- She opens Dashboard, navigates directly to job by name
- UI provides visual topology and formatted error messages
- Priya returns to CLI for log analysis with grep/filtering

**Postconditions**:
- Seamless workflow across UI and CLI based on task context
- Configuration portable between interfaces
- User can choose optimal tool for each activity

**Multi-Channel Requirements**:
- Consistent job naming and referencing
- Real-time status synchronization
- Configuration export/import between channels
- Equivalent functionality (user can accomplish any task in any channel)

---

## Documentation Considerations

### Documentation Strategy

MPIJob documentation must serve three distinct personas with different needs, knowledge levels, and learning styles. The documentation architecture balances comprehensive reference material with task-oriented quick-start guides, ensuring users can both achieve immediate success and develop deep expertise over time.

### Documentation Structure

```
OpenShift AI Documentation
├── Training & Models
│   ├── Training Job Types (Overview)
│   │   ├── Choosing the Right Job Type ⬅ Decision guide for job selection
│   │   ├── PyTorchJob
│   │   ├── TFJob
│   │   ├── MPIJob ⬅ New section
│   │   └── Job Comparison Matrix
│   │
│   ├── MPIJob Documentation ⬅ Primary landing page
│   │   ├── Overview & Use Cases
│   │   ├── Getting Started
│   │   │   ├── Quick Start Tutorial (30-minute first job)
│   │   │   ├── Prerequisites & Setup
│   │   │   └── Your First Distributed Training Job
│   │   ├── Configuration Guide
│   │   │   ├── Worker Topology & Sizing
│   │   │   ├── Resource Allocation
│   │   │   ├── MPI-Specific Configuration
│   │   │   └── Framework Integration (Horovod, DeepSpeed, Intel MPI)
│   │   ├── Monitoring & Troubleshooting
│   │   │   ├── Understanding Job Status
│   │   │   ├── Worker Health Monitoring
│   │   │   ├── Common Issues & Solutions
│   │   │   └── Log Analysis Guide
│   │   ├── Advanced Topics
│   │   │   ├── Performance Tuning
│   │   │   ├── Network Optimization
│   │   │   ├── Security Considerations
│   │   │   └── Integration with Data Science Pipelines
│   │   └── Reference
│   │       ├── MPIJob YAML Specification
│   │       ├── CLI Command Reference
│   │       ├── Python SDK API Reference
│   │       └── Configuration Options Complete List
│   │
│   └── Administration
│       └── MPIJob Administration
│           ├── Installation & Setup
│           ├── Resource Management & Quotas
│           ├── RBAC Configuration
│           ├── Network Policy Configuration
│           └── Monitoring & Governance
```

### Persona-Specific Documentation

#### Data Scientists

**Primary Documentation Needs:**
- Task-oriented guides: "How do I...?"
- Conceptual bridging: ML concepts → Kubernetes concepts
- Example-heavy: Copy-paste templates for common scenarios
- Visual aids: Diagrams, screenshots, topology visualizations

**Critical Documents:**

| Document | Purpose | Format | Success Metric |
|----------|---------|--------|----------------|
| **Quick Start Tutorial** | First successful job in <30 min | Step-by-step with screenshots | 80% completion rate |
| **Decision Guide** | "Should I use MPIJob?" | Flowchart with use cases | Users make correct choice |
| **Configuration Cookbook** | Common patterns (Horovod, DeepSpeed) | Template library with explanations | 70% use templates |
| **Troubleshooting Guide** | Self-service error resolution | Symptom → diagnosis → solution | 60% resolve without support |

**Bridging ML and Kubernetes Concepts:**

| User Concept | Technical Reality | Documentation Approach |
|--------------|------------------|----------------------|
| "Workers training together" | Launcher pod + Worker pods with MPI communication | Visual topology diagram with role explanations |
| "Training image" | Container image with code + dependencies + MPI | "Your training environment packaged as a container" |
| "Number of workers" | MPI replica count with launcher + workers | "How many machines work on training simultaneously" |
| "Slots per worker" | MPI process slots (usually matches GPU count) | "Parallel processes per worker—typically one per GPU" |

**Documentation Example - Quick Start Tutorial:**

```markdown
# Your First MPIJob: Distributed MNIST Training

**Time Required:** 30 minutes
**Prerequisites:** Access to OpenShift AI Dashboard, basic Python knowledge

## What You'll Build
Train a simple neural network on MNIST dataset using 2 workers with Horovod.

## Step 1: Access the Dashboard
1. Log in to OpenShift AI Dashboard: `https://your-cluster/`
2. Navigate to **Training Jobs** in the left sidebar
3. Click **Create Training Job**

[Screenshot: Dashboard navigation]

## Step 2: Select Job Type
1. You'll see three options: PyTorchJob, TFJob, MPIJob
2. Click **MPIJob** (we're using Horovod which requires MPI)
3. **Why MPIJob?** Use MPIJob when your training code uses MPI-based frameworks like Horovod

[Screenshot: Job type selection]

## Step 3: Configure Your Job
Fill in these fields:

- **Job Name:** `mnist-distributed` (must be unique)
- **Training Image:** `quay.io/opendatahub/horovod-mnist:latest` (pre-built example)
- **Number of Workers:** `2` (we'll use 2 workers for this demo)
- **CPUs per Worker:** `2`
- **Memory per Worker:** `4 GB`
- **Command:** `["horovodrun", "-np", "2", "python", "/examples/mnist.py"]`

Leave **Advanced Options** collapsed—defaults work for this example.

[Screenshot: Configuration form]

## Step 4: Launch Job
1. Click **Create Job**
2. You'll see a validation message: "Resources available ✓"
3. Job creation succeeds, you're redirected to job detail page

## Step 5: Monitor Progress
Watch your job train in real-time:

- **Status:** Should change from "Pending" → "Running" → "Succeeded"
- **Workers:** Visual grid shows 2 workers (should both be green/running)
- **Logs:** Click "View Logs" to see training output

[Screenshot: Job monitoring view]

Expected logs:
```
[Worker-1] Epoch 1/5 - Loss: 0.523
[Worker-2] Epoch 1/5 - Loss: 0.519
```

Training takes ~10 minutes.

## Step 6: View Results
Once status shows "Succeeded":

1. Click **Download Artifacts** to retrieve trained model
2. Model saved to `~/Downloads/mnist-distributed-model.h5`

**Congratulations!** You've run your first distributed training job.

## Next Steps
- [Run your own training code](./bring-your-own-code.md)
- [Scale to more workers](./scaling-workers.md)
- [Optimize performance](./performance-tuning.md)

## Troubleshooting
**Job stuck in "Pending"?** → Check resource availability in Dashboard
**Workers failing?** → View logs for specific worker that's red in topology view
```

#### MLOps Engineers

**Primary Documentation Needs:**
- API reference: Complete, searchable, with type signatures
- Integration patterns: How to embed in pipelines (Tekton, Argo, Kubeflow Pipelines)
- Automation examples: Reusable scripts for common workflows
- CLI command reference: Man-page style with all flags

**Critical Documents:**

| Document | Purpose | Format | Success Metric |
|----------|---------|--------|----------------|
| **Python SDK Reference** | Comprehensive API docs | Searchable with interactive examples | <2 min to find any API |
| **CLI Reference** | Complete command catalog | Man-page style | Support tickets ↓ 30% |
| **Integration Guide** | Pipeline embedding patterns | Code examples for common CI/CD tools | 50% use provided patterns |
| **Automation Cookbook** | Reusable scripts | GitHub repo with working examples | 100+ clones in 3 months |

**SDK Documentation Example:**

```python
"""
OpenShift AI Training SDK - MPIJob API Reference
"""

from odh.training import MPIJob, ResourceSpec

# === Quick Example ===
job = MPIJob.from_template("horovod-tensorflow")
job.create(wait=True)

# === Full API Reference ===

class MPIJob:
    """
    Represents an MPI-based distributed training job.

    Use MPIJob for training workloads that use Message Passing Interface (MPI)
    for inter-process communication, such as Horovod or DeepSpeed with MPI backend.

    Example:
        >>> job = MPIJob(
        ...     name="training-job",
        ...     workers=4,
        ...     slots_per_worker=2,
        ...     image="myregistry.com/training:latest",
        ...     command=["horovodrun", "-np", "8", "python", "train.py"]
        ... )
        >>> job.create()
        >>> job.wait_for_completion()
        >>> artifacts = job.get_artifacts()
    """

    def __init__(
        self,
        name: str,
        namespace: str = "default",
        workers: int = 1,
        slots_per_worker: int = 1,
        image: str,
        command: List[str],
        launcher_resources: Optional[ResourceSpec] = None,
        worker_resources: Optional[ResourceSpec] = None,
        env: Optional[Dict[str, str]] = None,
        volumes: Optional[List[Volume]] = None,
        mpi_implementation: Literal["OpenMPI", "IntelMPI", "MPICH"] = "OpenMPI",
        clean_pod_policy: Literal["Running", "All", "None"] = "Running"
    ):
        """
        Initialize MPIJob configuration.

        Args:
            name: Unique job identifier within namespace (must be DNS-1123 compliant)
            namespace: Kubernetes namespace for job (default: "default")
            workers: Number of worker replicas (minimum: 1, maximum: 100)
            slots_per_worker: MPI process slots per worker (typically matches GPU count)
            image: Container image with training code and MPI installed
            command: Entrypoint command for training script
            launcher_resources: CPU/memory/GPU allocation for launcher pod
            worker_resources: CPU/memory/GPU allocation for each worker pod
            env: Environment variables injected into all pods
            volumes: Persistent volumes mounted to pods (for datasets, models)
            mpi_implementation: MPI implementation (must match container image)
            clean_pod_policy: When to delete pods after job completion

        Raises:
            ValueError: If configuration is invalid (e.g., workers < 1)
            AuthenticationError: If user lacks namespace access

        Example:
            >>> from odh.training import MPIJob, ResourceSpec
            >>> job = MPIJob(
            ...     name="my-training",
            ...     workers=4,
            ...     worker_resources=ResourceSpec(cpu="8", memory="32Gi", gpus=2),
            ...     image="myregistry.com/training:v1.0",
            ...     command=["python", "train.py", "--epochs", "100"]
            ... )
        """
        pass

    def create(self, wait: bool = False, timeout: int = 3600) -> None:
        """
        Submit MPIJob to cluster for execution.

        Args:
            wait: Block until job completes (default: False, non-blocking)
            timeout: Maximum wait time in seconds if wait=True (default: 3600)

        Raises:
            ResourceQuotaExceeded: If job would exceed namespace quota
            ValidationError: If job specification is invalid
            TimeoutError: If wait=True and job doesn't complete within timeout

        Example:
            >>> job.create()  # Non-blocking submission
            >>> # OR
            >>> job.create(wait=True, timeout=7200)  # Wait up to 2 hours
        """
        pass

    def status(self) -> JobStatus:
        """
        Get current job execution status.

        Returns:
            JobStatus object with fields:
                - phase: str (Pending, Running, Succeeded, Failed)
                - workers_running: int
                - workers_pending: int
                - workers_failed: int
                - start_time: datetime
                - completion_time: Optional[datetime]

        Example:
            >>> status = job.status()
            >>> if status.phase == "Running":
            ...     print(f"{status.workers_running}/{job.workers} workers active")
        """
        pass

    def wait_for_completion(self, timeout: int = 3600, poll_interval: int = 10) -> JobStatus:
        """
        Block until job completes (Succeeded or Failed).

        Args:
            timeout: Maximum wait time in seconds
            poll_interval: Seconds between status checks

        Returns:
            Final JobStatus

        Raises:
            TimeoutError: If job doesn't complete within timeout
            JobFailedError: If job fails (includes logs and error details)

        Example:
            >>> try:
            ...     final_status = job.wait_for_completion(timeout=3600)
            ...     print(f"Job completed in {final_status.duration}")
            ... except JobFailedError as e:
            ...     print(f"Training failed: {e.message}")
            ...     print(f"Failed worker logs: {e.logs}")
        """
        pass

    def get_logs(
        self,
        component: Literal["launcher", "worker"] = "launcher",
        worker_id: Optional[int] = None,
        tail: Optional[int] = None
    ) -> str:
        """
        Retrieve logs from job pods.

        Args:
            component: "launcher" or "worker"
            worker_id: Specific worker index if component="worker" (0-indexed)
            tail: Return last N lines only (None = all logs)

        Returns:
            Log output as string

        Example:
            >>> # Get launcher logs
            >>> launcher_logs = job.get_logs(component="launcher")
            >>>
            >>> # Get specific worker logs
            >>> worker_logs = job.get_logs(component="worker", worker_id=2, tail=100)
        """
        pass

    def delete(self, wait: bool = False) -> None:
        """
        Delete job and associated resources from cluster.

        Args:
            wait: Block until all pods are terminated (default: False)

        Example:
            >>> job.delete(wait=True)  # Blocking deletion
        """
        pass

    @classmethod
    def from_template(cls, template_name: str, **overrides) -> 'MPIJob':
        """
        Create MPIJob from predefined template.

        Available templates:
            - "horovod-tensorflow": TensorFlow with Horovod
            - "horovod-pytorch": PyTorch with Horovod
            - "deepspeed-mpi": DeepSpeed with MPI backend
            - "intel-mpi-pytorch": PyTorch with Intel MPI

        Args:
            template_name: Name of template
            **overrides: Template parameters to override

        Example:
            >>> job = MPIJob.from_template(
            ...     "horovod-pytorch",
            ...     name="my-job",
            ...     workers=8,
            ...     worker_resources=ResourceSpec(gpus=4)
            ... )
        """
        pass

# === Error Reference ===

class ResourceQuotaExceeded(Exception):
    """Raised when job would exceed namespace resource quota."""
    def __init__(self, details: str):
        self.details = details  # Human-readable quota information

class ValidationError(Exception):
    """Raised when job configuration is invalid."""
    def __init__(self, field: str, message: str):
        self.field = field      # Which field is invalid
        self.message = message  # Why it's invalid

class JobFailedError(Exception):
    """Raised when job execution fails."""
    def __init__(self, message: str, logs: str):
        self.message = message  # High-level failure reason
        self.logs = logs       # Relevant log excerpts
```

#### Administrators

**Primary Documentation Needs:**
- Installation and setup guides
- Resource management strategies (quotas, node pools)
- Security and compliance configuration
- Monitoring and troubleshooting at system level

**Critical Documents:**

| Document | Purpose | Format | Success Metric |
|----------|---------|--------|----------------|
| **Installation Guide** | Enable MPIJob support | Step-by-step with verification | 95% successful installation |
| **Resource Management** | Capacity planning & quotas | Decision matrices with examples | Admins configure quotas correctly |
| **Security Guide** | Network policies, RBAC | Templates with customization instructions | Zero security incidents |
| **Troubleshooting** | System-level diagnostics | Symptom-based flowcharts | Mean time to resolution ↓ 40% |

### Documentation Integration Points

#### In-Product Documentation (Contextual Help)

| Location | Content Type | Example |
|----------|--------------|---------|
| Job Creation Form | Tooltip | "Slots per worker: Number of MPI processes per worker. Typically matches GPU count. [Learn more →]" |
| Status Page | Inline help | "Job Pending: Waiting for resources. [View resource availability] [Troubleshooting guide]" |
| Error Messages | Direct documentation links | "Error: ResourceQuotaExceeded [View quota documentation] [Contact administrator]" |
| CLI Help | Embedded examples | `odh training create mpijob --help` includes copy-paste examples |

#### External Documentation Links

**Links OUT (from OpenShift AI docs):**
- KubeFlow Training Operator (upstream reference)
- Horovod documentation (framework integration)
- OpenMPI/Intel MPI documentation (MPI implementation details)
- OpenShift general concepts (for users new to OpenShift)

**Links IN (to OpenShift AI docs):**
- From OpenShift AI main documentation hub
- From training job comparison matrix
- From Data Science Pipelines integration docs
- From SDK/CLI reference pages

### Documentation Maintenance Strategy

**Version Control:**
- Documentation versioned alongside OpenShift AI releases
- Clear version indicators on every page
- Legacy version documentation archived and accessible

**Content Updates:**
- Quarterly review of accuracy based on product changes
- User feedback loop: "Was this helpful?" on every page
- Analytics tracking: most viewed pages, search queries, bounce rates

**Community Contributions:**
- GitHub repository for docs with PR process
- Community-contributed examples and tutorials
- Documentation office hours for Q&A

---

## Questions to Answer

### Technical Architecture Questions

**Q-ARCH-001: MPI Implementation Support**
- **Question**: Should MPIJob support all MPI implementations (OpenMPI, Intel MPI, MPICH) or focus on a primary implementation for MVP?
- **Decision Needed By**: Architecture review phase
- **Impact**: Affects testing scope, documentation complexity, container image requirements
- **Stakeholders**: Engineering, Product, Technical Support

**Q-ARCH-002: Integration and Deployment Architecture**
- **Question**: How will MPIJob integrate with existing KubeFlow Training Operator V2 deployment in OpenShift AI?
- **Sub-questions**:
  - Is Training Operator V2 deployed as a separate component or embedded?
  - What are the upgrade path implications when Training Operator upstream changes?
  - How do we handle version skew between OpenShift AI and upstream KubeFlow?
- **Decision Needed By**: Design phase
- **Impact**: Deployment model, upgrade strategy, support boundaries

**Q-ARCH-003: Networking and Inter-process Communication (IPC)**
- **Question**: What network configuration is required for MPI communication between workers?
- **Sub-questions**:
  - Does MPI traffic require specific ports to be opened?
  - Are there CNI plugin requirements or limitations?
  - How do we handle network policies in restricted environments?
  - Is RDMA support required for high-performance scenarios?
- **Decision Needed By**: Architecture review
- **Impact**: Network policy configuration, performance characteristics, enterprise compatibility

**Q-ARCH-004: Resource Orchestration and Scheduling**
- **Question**: How does MPIJob handle pod scheduling for multi-worker topologies?
- **Sub-questions**:
  - Is gang scheduling required to prevent partial job placement?
  - How do we ensure workers are placed for optimal network performance?
  - Can workers span multiple node pools or GPU types?
  - What happens if a worker pod is evicted mid-training?
- **Decision Needed By**: Design phase
- **Impact**: Reliability, performance, resource utilization

**Q-ARCH-005: Hardware Acceleration and Runtime**
- **Question**: What GPU types and accelerator configurations are supported?
- **Sub-questions**:
  - Support for NVIDIA GPUs only, or also AMD/Intel accelerators?
  - Does GPU topology awareness matter for MPI communication?
  - Are GPU affinities and NUMA locality respected?
  - Support for mixed GPU types within a single job?
- **Decision Needed By**: Requirements refinement
- **Impact**: Hardware compatibility matrix, performance optimization

**Q-ARCH-006: Required Software and Data Services**
- **Question**: What dependencies must be present in training container images?
- **Sub-questions**:
  - Must users install specific MPI versions?
  - Are there SSH or communication daemon requirements?
  - How is distributed filesystem access handled (NFS, S3, etc.)?
  - Support for proprietary MPI implementations (e.g., IBM Spectrum MPI)?
- **Decision Needed By**: Documentation planning
- **Impact**: User onboarding complexity, container image guidelines

### Operational and Non-Functional Requirements

**Q-NFR-001: Performance Baseline and SLAs**
- **Question**: What performance guarantees can we provide for MPIJob?
- **Sub-questions**:
  - Expected speedup characteristics (2x with 2 workers, 4x with 4 workers, etc.)?
  - Acceptable overhead for job submission and worker initialization?
  - Network throughput requirements for acceptable training performance?
- **Decision Needed By**: Testing phase
- **Impact**: Customer expectations, competitive positioning, performance testing scope

**Q-NFR-002: Scalability Limits**
- **Question**: What are the maximum supported job sizes?
- **Sub-questions**:
  - Maximum number of workers per job?
  - Maximum number of concurrent MPIJobs per cluster?
  - Supported cluster sizes (node count, GPU count)?
- **Decision Needed By**: Architecture review
- **Impact**: Testing scope, documentation, quota recommendations

**Q-NFR-003: Reliability and Fault Tolerance**
- **Question**: How does MPIJob handle failures?
- **Sub-questions**:
  - What happens if one worker fails—does entire job fail?
  - Support for checkpoint/restart?
  - Automatic retry of failed jobs?
  - Node failure handling?
- **Decision Needed By**: Design phase
- **Impact**: User experience, training job success rates

### Monitoring, Logging, and Debugging

**Q-OBS-001: Metrics Collection**
- **Question**: What metrics are automatically collected vs. user-configured?
- **Sub-questions**:
  - Are training metrics (loss, accuracy) collected, or only infrastructure metrics?
  - Integration with Prometheus/OpenShift monitoring?
  - Custom metric support for user-defined KPIs?
- **Decision Needed By**: Design phase
- **Impact**: Observability capabilities, user experience

**Q-OBS-002: Distributed Log Aggregation**
- **Question**: How are logs from multiple workers aggregated and presented?
- **Sub-questions**:
  - Real-time log streaming vs. batch retrieval?
  - Log correlation across workers (timestamp synchronization)?
  - Log retention policies?
  - Integration with OpenShift logging stack (EFK/ELK)?
- **Decision Needed By**: Design phase
- **Impact**: Troubleshooting experience, storage requirements

**Q-OBS-003: Debugging Tools**
- **Question**: What debugging capabilities are provided for distributed training?
- **Sub-questions**:
  - Can users exec into running worker pods?
  - Support for attaching debuggers to training processes?
  - Distributed profiling tools (MPI profilers)?
  - Post-mortem analysis of failed jobs?
- **Decision Needed By**: Requirements refinement
- **Impact**: Developer productivity, support burden

### Security and Access Control

**Q-SEC-001: RBAC Model**
- **Question**: What permissions are required for MPIJob operations?
- **Sub-questions**:
  - Is MPIJob creation a privileged operation vs. other job types?
  - Separate permissions for create/read/update/delete?
  - Role templates for common use cases (data-scientist, mlops-engineer)?
- **Decision Needed By**: Design phase
- **Impact**: Security model, onboarding complexity

**Q-SEC-002: Network Security**
- **Question**: What network policies are required for MPIJob communication?
- **Sub-questions**:
  - Default-deny network policies compatible?
  - Specific port ranges that must be allowed?
  - East-west traffic encryption support?
  - Integration with service mesh (Istio)?
- **Decision Needed By**: Architecture review
- **Impact**: Enterprise adoption, compliance requirements

**Q-SEC-003: Secret Management**
- **Question**: How are credentials and sensitive configuration managed?
- **Sub-questions**:
  - Support for OpenShift secrets, ConfigMaps?
  - Integration with external secret management (Vault, AWS Secrets Manager)?
  - Environment variable injection vs. volume mounts for secrets?
- **Decision Needed By**: Design phase
- **Impact**: Security posture, compliance

**Q-SEC-004: Multi-Tenancy Isolation**
- **Question**: What isolation guarantees exist between MPIJobs in different namespaces?
- **Sub-questions**:
  - Can MPIJob traffic leak across namespace boundaries?
  - Resource consumption visibility across tenants?
  - Audit trail per tenant?
- **Decision Needed By**: Architecture review
- **Impact**: Enterprise multi-tenant deployments

### User Experience and Documentation

**Q-UX-001: Abstraction Level**
- **Question**: How much MPI complexity do we expose to data scientists?
- **Research Needed**: Mental model study with target users (8-10 data scientists)
- **Decision Criteria**: Task completion rate in usability testing
- **Impact**: UI design, documentation approach, learning curve

**Q-UX-002: Status Granularity**
- **Question**: What level of worker status detail do we show by default?
- **Research Needed**: Eye-tracking study, user interviews on information needs
- **Decision Criteria**: User comprehension and troubleshooting efficiency
- **Impact**: Dashboard UI design, information architecture

**Q-UX-003: Cross-Channel Symmetry**
- **Question**: Should CLI/SDK support all UI features, or can some be UI-exclusive?
- **Research Needed**: Usage analytics, MLOps workflow analysis
- **Decision Criteria**: Channel usage patterns, feature request frequency
- **Impact**: Development scope, API design

**Q-UX-004: Error Recovery Patterns**
- **Question**: Can users edit failed jobs, or must they create new ones?
- **Research Needed**: Mental model study, competitive analysis
- **Decision Criteria**: User expectations, audit trail requirements
- **Impact**: UI design, job immutability semantics

**Q-UX-005: Training Metrics Integration**
- **Question**: How deeply should MPIJob integrate with training metrics (loss, accuracy)?
- **Research Needed**: User journey mapping, tool usage analysis
- **Decision Criteria**: Implementation complexity vs. user value
- **Impact**: Observability feature scope, integration architecture

### Product and Business Questions

**Q-BIZ-001: Pricing and Packaging**
- **Question**: Does MPIJob require premium SKU or is it included in base OpenShift AI?
- **Stakeholders**: Product Management, Sales, Legal
- **Decision Needed By**: Pre-GA pricing finalization
- **Impact**: Revenue model, competitive positioning, go-to-market strategy

**Q-BIZ-002: Cloud Service Integration**
- **Question**: Should MPIJob be available in managed Red Hat OpenShift AI services (ROSA, ARO)?
- **Sub-questions**:
  - Different feature set for managed vs. self-hosted?
  - Cloud-specific optimizations (spot instances, accelerator types)?
- **Decision Needed By**: Roadmap planning
- **Impact**: Total addressable market, cloud partnerships

**Q-BIZ-003: Framework Certification**
- **Question**: Which ML frameworks and versions do we officially support and test?
- **Sub-questions**:
  - Certification program for "tested" vs. "compatible" frameworks?
  - Partner ecosystem (NVIDIA, Intel, framework maintainers)?
- **Decision Needed By**: Pre-GA testing
- **Impact**: Support commitments, partnership opportunities

**Q-BIZ-004: Migration Support**
- **Question**: What level of migration support do we provide for customers moving from other platforms?
- **Sub-questions**:
  - Professional services offering for migration?
  - Migration tools or just documentation?
  - Which source platforms prioritized (SageMaker, Azure ML, Horovod standalone)?
- **Decision Needed By**: Go-to-market planning
- **Impact**: Customer acquisition, services revenue

---

## Background & Strategic Fit

### Market Landscape and Competitive Position

The enterprise AI/ML platform market is experiencing rapid consolidation around end-to-end capabilities spanning data preparation, model training, deployment, and monitoring. Distributed training—specifically the ability to scale model training across multiple nodes and accelerators—has transitioned from a niche requirement for research teams to a mainstream expectation for production AI platforms.

**Competitive Analysis:**

| Capability | AWS SageMaker | Azure ML | Google Vertex AI | Red Hat OpenShift AI (Current) | OpenShift AI + MPIJobs |
|------------|---------------|----------|------------------|--------------------------------|------------------------|
| Distributed Training | Native, automatic | Native with Azure ML Compute | Native with Vertex AI Training | DIY with custom operators | Native via KubeFlow Trainer V2 |
| MPI Protocol Support | Yes (Horovod) | Yes (Horovod, DeepSpeed) | Yes (Horovod) | No (manual setup only) | Yes (standard MPI implementations) |
| On-Premises Deployment | No (cloud-only) | No (cloud-only) | No (cloud-only) | Yes (core value prop) | Yes |
| Multi-Cloud / Hybrid | Limited | Limited | Limited | Yes (run anywhere) | Yes |
| Open Source Foundation | Proprietary | Proprietary | Proprietary | KubeFlow, Kubernetes | KubeFlow, Kubernetes |
| GPU Utilization Tools | CloudWatch insights | Azure Monitor | Cloud Monitoring | Prometheus/Grafana | Prometheus/Grafana |

**Key Insight**: Hyperscalers (AWS, Azure, Google) have commoditized distributed training capabilities. Customers expect these features as table stakes—our differentiation lies in hybrid cloud flexibility and open source transparency, but only if we match their functional capabilities.

**Market Data**:
- **67% of enterprise ML teams** require distributed training for production workloads (Gartner, 2024)
- **$12.4B total addressable market** for enterprise ML platforms by 2026 (IDC)
- **82% of enterprises** cite "vendor lock-in concerns" as top barrier to cloud AI adoption (Forrester, 2024)

Red Hat's strategic advantage—on-premises and hybrid cloud deployment—resonates strongly with regulated industries (financial services, healthcare, government) that cannot or will not move training workloads to public clouds. However, this advantage evaporates if we lack feature parity on core training capabilities.

### Competitive Gap Analysis

**Current State Pain**:
We've lost 3 high-value deals ($800K+ ARR each) in the last two quarters specifically citing "lack of distributed training support" as a disqualifying factor. Customer feedback reveals:

1. **Financial Services**: "We need to train fraud detection models on 500M transactions. Single-node training takes 2 weeks—we can't iterate fast enough."
2. **Healthcare**: "Medical imaging models with 3D CNNs won't fit in single-GPU memory. We're stuck with Azure ML because OpenShift AI can't scale."
3. **Automotive**: "Autonomous vehicle training requires massive parallelism. We want OpenShift's hybrid cloud story, but we can't compromise on training performance."

These are not edge cases—they represent mainstream enterprise AI use cases. The absence of native distributed training forces customers into one of three suboptimal paths:

- **Path A**: Stay with hyperscaler platforms (we lose the deal)
- **Path B**: Build custom MPI integrations (increases operational complexity, defeats "managed platform" value proposition)
- **Path C**: Compromise model quality by reducing dataset size or model complexity (unacceptable for competitive industries)

### Strategic Rationale: Why MPIJobs Now?

**1. Competitive Necessity (Urgency: HIGH)**

The competitive gap is widening, not narrowing. Recent announcements:
- **AWS SageMaker** (October 2024): 40% training speedup with enhanced distributed training
- **Azure ML** (September 2024): Automated distributed training configuration (users specify goals, system determines topology)
- **Google Vertex AI** (August 2024): Integrated distributed training with TPU pods (50+ accelerators)

If we don't close this gap within 2 quarters, the market perception will solidify: "OpenShift AI is fine for small models, but not serious about large-scale deep learning."

**2. KubeFlow Trainer V2 Momentum (Timing: OPTIMAL)**

The upstream KubeFlow community has consolidated multiple training operators (TFJob, PyTorchJob, MPIJob, etc.) under the unified Training Operator V2 API. This architectural shift provides:

- **Consistent API Surface**: Same patterns across all job types—reduces implementation complexity
- **Upstream Maturity**: MPIJob is a stable, production-tested component (4+ years in CNCF)
- **Reduced Technical Debt**: Implementing MPIJob now aligns with Trainer V2 architecture; delaying creates refactoring burden

Red Hat is already committed to adopting Training Operator V2 for OpenShift AI. Adding MPIJob is a natural extension of that work, not a separate initiative. The incremental effort is relatively low while the customer value is high—this is the optimal window for delivery.

**3. Customer Demand Signal (Validation: STRONG)**

Quantified customer demand:
- **23 active customer requests** for distributed training (tracked in Salesforce)
- **14 specifically mention MPI** or Horovod
- **8 are from existing customers** willing to expand usage if capability is available
- **Estimated revenue impact**: $5M+ net-new ARR within 12 months of GA

Priority customer profiles:
- **Financial Services**: Fraud detection, risk modeling, algorithmic trading (models with 1B+ parameters)
- **Healthcare**: Medical imaging (3D CNNs), genomics (transformer models), drug discovery (molecular dynamics)
- **Automotive**: Autonomous vehicle perception (sensor fusion models), simulation (digital twin training)
- **Research Institutions**: Climate modeling, materials science, computational biology

These customers have budget, urgency, and clear use cases. They are not asking for experimental features—they need capabilities that competitors already provide.

**4. Open Source Differentiation (Market Positioning)**

Unlike hyperscaler proprietary distributed training solutions, KubeFlow MPIJob provides:

- **No Vendor Lock-In**: Customers can replicate OpenShift AI environment anywhere (cloud, on-prem, edge)
- **Full Transparency**: Users understand and control the underlying mechanics
- **Community Innovation**: Contributions flow back to upstream, accelerating ecosystem evolution
- **Cost Predictability**: No surprise costs for training job execution (vs. SageMaker per-minute billing)

This aligns perfectly with Red Hat's strategic narrative: "Open source innovation meets enterprise reliability." In customer conversations, this resonates strongly—especially post-procurement, when cloud bills become painful surprises.

### Red Hat AI/ML Strategy Alignment

MPIJob support directly advances Red Hat's stated strategic pillars for AI/ML:

**Pillar 1: Enterprise-Grade AI Platform**
- Completes training job portfolio (PyTorchJob, TFJob, MPIJob = comprehensive coverage)
- Matches hyperscaler capabilities while maintaining differentiation (hybrid cloud, open source)
- Positions OpenShift AI as credible for production-scale AI workloads

**Pillar 2: Hybrid Cloud Leadership**
- Enables distributed training on customer infrastructure (not forced cloud migration)
- Same experience on-prem, ROSA, ARO, or multi-cloud
- Addresses regulatory and data sovereignty requirements for distributed training

**Pillar 3: Open Source Innovation**
- Deepens Red Hat's leadership in KubeFlow community
- Contributions to Training Operator upstream benefit entire ecosystem
- Demonstrates commitment to upstream-first development

**Pillar 4: Kubernetes-Native AI**
- All AI workloads managed through standard Kubernetes APIs
- Unified RBAC, networking, storage, observability for ML and application workloads
- Reduces operational burden vs. siloed ML platforms

### Product Roadmap Context

**Prerequisites:**
- **KubeFlow Training Operator V2 Integration**: Must be completed before MPIJob work begins (estimated Q1 2025)
- **PyTorchJob and TFJob Migration**: Move existing job types to Trainer V2 API (parallel track)

**Parallel Initiatives:**
- **Data Science Pipelines V2**: Workflow orchestration for multi-step ML pipelines (MPIJob integrates as pipeline step)
- **Model Registry Enhancements**: Artifact management for trained models (MPIJob outputs integrate)
- **Observability Improvements**: Unified monitoring across all job types (MPIJob benefits from shared infrastructure)

**Follow-On Opportunities:**
- **Advanced Scheduling**: Gang scheduling, topology-aware placement for optimal distributed training performance
- **Elastic Training**: Dynamic worker scaling based on available resources
- **Framework-Specific Optimizations**: Pre-configured templates for Horovod, DeepSpeed, Intel MPI

### Market Timing and Urgency

**Why This Matters Now:**

1. **Model Scale Trajectory**: GPT-4 class models (175B+ parameters) cannot be trained on single nodes. As enterprises adopt transformer architectures and multimodal models, distributed training becomes mandatory—not optional.

2. **Cost Pressure**: Cloud training costs are driving enterprises toward on-premises alternatives. AWS bills for distributed training can exceed $10K/day for large models. Customers see ROI in owning infrastructure, but need equivalent tooling.

3. **Regulatory Environment**: EU AI Act, GDPR, CCPA, and sector-specific regulations (HIPAA, PCI-DSS) create compliance pressures. On-premises distributed training addresses data residency and auditability requirements that cloud solutions struggle to meet.

4. **Competitive Dynamics**: SageMaker, Azure ML, and Vertex AI are not standing still. Recent enhancements (automated distributed training, performance optimizations) raise the bar continuously. Delay means falling further behind.

**Risk of Delay:**

- **Immediate (Q1-Q2 2025)**: Continue losing deals to hyperscalers, revenue impact $2-3M/quarter
- **Medium-term (Q3-Q4 2025)**: Customer workarounds (DIY MPI, external Horovod) create technical debt that's hard to migrate, reducing lifetime value
- **Long-term (2026+)**: Market perception hardens: "OpenShift AI is for simple ML, not serious deep learning"—becomes difficult to overcome even with feature parity

Conversely, delivering MPIJobs in Q2 2025:
- **Closes competitive gap** during key budget cycles (many enterprises finalize AI platform decisions in Q2-Q3)
- **Enables customer success stories** showcasing hybrid cloud distributed training (powerful marketing differentiation)
- **Positions OpenShift AI** for inclusion in analyst reports as "comprehensive AI platform" (Gartner Magic Quadrant, Forrester Wave)

### Success Definition

**What Good Looks Like (12 Months Post-GA):**

- **Customer Adoption**: 40% of OpenShift AI customers with GPU infrastructure actively use MPIJobs
- **Deal Impact**: Distributed training no longer appears in "reasons we lost" deal reviews
- **Revenue**: $5M+ in net-new ARR attributable to distributed training capability
- **Market Position**: Analyst recognition (Gartner, Forrester) cites distributed training as differentiation factor
- **Community**: Red Hat recognized as top contributor to KubeFlow Training Operator upstream

**What Failure Looks Like:**

- <10% adoption due to complexity, performance issues, or poor integration
- Continued customer escalations about missing distributed training features
- Sales team still positioning distributed training as "roadmap item" rather than "available capability"
- Competitors extend lead with next-generation distributed training features

### Execution Recommendations

**Phase 1: MVP (Target Q2 2025)**
- Focus on feature parity with existing job types (TFJob, PyTorchJob)
- CLI, SDK, UI basics with unified observability
- Documentation for common patterns (Horovod, DeepSpeed)
- **Goal**: Enable customers to run distributed training with OpenShift AI; unblock deals

**Phase 2: Enterprise Hardening (Target Q3 2025)**
- Address enterprise customer considerations: multi-tenancy, compliance, security
- Migration documentation (Horovod standalone → MPIJob)
- Reference architectures for common deployment patterns
- Support training and troubleshooting runbooks
- **Goal**: Production-ready for regulated industries; drive adoption

**Phase 3: Optimization (Target Q4 2025)**
- Advanced features based on adoption data: topology visualization, cost optimization, performance tuning
- Framework-specific templates and optimizations
- Integration with Data Science Pipelines for end-to-end workflows
- **Goal**: Best-in-class distributed training experience; customer advocacy and case studies

### Key Success Factors

1. **Early Customer Access**: Beta program with 5-10 strategic accounts to validate requirements and gather feedback before GA
2. **Support Readiness**: Level 2/3 support training on distributed training concepts and common failure modes before GA launch
3. **Documentation Quality**: Assume users are ML experts but not Kubernetes experts—make it simple and example-driven
4. **Performance Validation**: Benchmark against cloud providers, publish results demonstrating comparable or superior performance
5. **Executive Sponsorship**: Visible commitment from product and engineering leadership signals priority to customers and internal teams

---

## Customer Considerations

### Multi-Tenancy

**Customer Need**: Multiple teams sharing OpenShift AI infrastructure without interference or resource contention.

**Requirements**:
- **Namespace Isolation**: MPIJobs in different namespaces must be fully isolated (network, resources, visibility)
- **Resource Quota Enforcement**: Per-namespace quotas for GPU count, CPU, memory, and concurrent job count
- **Fair Scheduling**: Prevent single team from monopolizing cluster resources with large distributed jobs
- **Chargeback/Showback**: Track resource consumption per team/project for cost allocation

**Implementation Considerations**:
- MPIJobs must respect OpenShift ResourceQuota and LimitRange objects
- Network policies must prevent cross-namespace communication between MPI workers
- Audit logs must capture namespace context for all operations
- Dashboard UI must provide namespace-scoped views and admin-level cross-namespace visibility

**Migration Consideration**: Customers running DIY distributed training must map existing access controls to OpenShift RBAC model. Many have informal "team allocation" schemes that need formalization.

**Customer Quote**: *"We have 5 data science teams sharing a 64-GPU cluster. If one team launches an 8-node job that blocks everyone else for hours, we have a political problem, not just a technical one."* — ML Platform Lead, Fortune 500 Retailer

---

### Compliance and Governance

**Customer Need**: Audit trails, data sovereignty, and regulatory compliance (GDPR, HIPAA, SOC2, PCI-DSS) for distributed training workloads.

**Requirements**:
- **Audit Logging**: All MPIJob operations (create, delete, modify) must generate audit events with user identity, timestamp, job details
- **Data Residency**: Support for node selectors and topology constraints to enforce geographic placement of training workloads
- **Secrets Management**: Secure handling of credentials for container registries, data sources, and APIs
- **Policy Enforcement**: Integration with OPA/Gatekeeper for validation of job specifications against organizational policies
- **Immutable Audit Trail**: Audit events cannot be modified or deleted, even by administrators

**Compliance Scenarios**:
- **GDPR (EU Data Protection)**: Training data must remain in EU datacenters; audit trail proves compliance
- **HIPAA (Healthcare)**: PHI accessed during training requires audit trail, encryption in transit, role-based access controls
- **SOC2 (SaaS Providers)**: Evidence of change management, access controls, and monitoring for auditor reviews
- **PCI-DSS (Financial Services)**: Training jobs accessing payment data require network segmentation and logging

**Implementation Considerations**:
- Integrate with OpenShift audit logging (forwarded to SIEM systems)
- Document compliance mappings: "How MPIJob satisfies HIPAA technical safeguards"
- Provide policy templates for common compliance frameworks

**Customer Quote**: *"We can't adopt distributed training unless we can prove to auditors exactly who ran what jobs, when, with what data, and on which infrastructure. A missing audit trail is a showstopper."* — CISO, Fortune 500 Financial Services

---

### Migration from Existing Solutions

**Common Migration Paths**:

1. **Horovod Standalone → MPIJob** (45% of customers)
   - **Current State**: Custom scripts launching `horovodrun` across SSH-connected nodes
   - **Migration Challenge**: Containerization, configuration as Kubernetes YAML
   - **Support Needed**: Side-by-side code examples, container image guidelines

2. **Custom MPI Scripts → MPIJob** (30% of customers)
   - **Current State**: Handwritten `mpirun` invocations, manual hostfile management
   - **Migration Challenge**: Understanding KubeFlow abstractions, launcher vs. worker roles
   - **Support Needed**: Conceptual mapping documentation, troubleshooting guide

3. **Ray/Dask → MPIJob** (15% of customers)
   - **Current State**: Python-native distributed frameworks
   - **Migration Challenge**: Code refactoring from Ray actors to MPI ranks
   - **Support Needed**: Clear messaging about when to migrate vs. stay with Ray/Dask

4. **Cloud Provider Native → MPIJob** (10% of customers)
   - **Current State**: SageMaker distributed training, Azure ML parallel jobs
   - **Migration Challenge**: Translating cloud-specific config to KubeFlow, cost modeling
   - **Support Needed**: Cloud-to-OpenShift AI comparison guide, TCO calculator

**Migration Support Deliverables**:
- **Documentation**: Step-by-step migration guides for each path
- **Code Examples**: Before/after comparisons showing cloud → OpenShift AI
- **Container Image Guidelines**: "Building MPI-Compatible Training Images" guide
- **Network Configuration**: CNI plugin requirements, performance tuning
- **Estimation Tool**: "How long will migration take?" assessment questionnaire

**Migration Risk**: Customers expect "drop-in replacement" simplicity. Reality requires container packaging, YAML authoring, and often code modifications. Under-communicating effort leads to adoption resistance.

**Recommendation**: Publish realistic effort estimates:
- **Simple Case** (Horovod with existing containers): 1-2 days
- **Moderate Case** (Custom MPI, need to containerize): 1-2 weeks
- **Complex Case** (Ray migration, code refactoring): 1-2 months

---

### Support and SLA Expectations

**Customer Need**: Production-grade reliability with enterprise support for distributed training workloads.

**Requirements**:

- **Availability SLA**: 99.9% uptime for control plane (job submission, status queries)
- **Recovery SLA**: Failed jobs must not orphan resources or leave cluster in degraded state
- **Debugging Support**: Clear error messages, troubleshooting runbooks, escalation paths to L2/L3
- **Performance Baseline**: Documented expectations for job submission latency, worker initialization time, training throughput

**Support Challenges**:

1. **Complexity of Distributed Training Failures**:
   - Network timeouts between workers (CNI plugin issue? Firewall? Application bug?)
   - MPI rank mismatches (configuration error? Container image mismatch?)
   - Out-of-memory on specific workers (insufficient resources? Memory leak? Data imbalance?)

2. **Boundary Between Platform and User Code**:
   - Customer expectation: "Red Hat support should debug my training script"
   - Reality: Support can diagnose infrastructure issues, not model training bugs
   - **Need**: Clear documentation of support boundaries

3. **Performance Troubleshooting**:
   - "My distributed job is slower than single-node training" (often user code issue: poor data parallelization, communication bottlenecks)
   - Support requires understanding of MPI communication patterns, framework internals
   - **Need**: L2/L3 training on distributed training concepts

**Support Readiness Plan**:
- **Pre-GA Training**: L2/L3 support engineers complete distributed training workshop (Horovod, MPI basics, common failure patterns)
- **Runbooks**: Documented troubleshooting procedures for top 10 failure scenarios
- **Known Good Configurations**: Validated reference architectures with tested container images
- **Escalation Path**: When to involve engineering vs. resolve at L2

**Customer Quote**: *"We're paying for enterprise support. If a distributed training job fails and support says 'this is a user code issue,' that's unacceptable—we need help diagnosing whether it's infrastructure or application."* — VP of Engineering, Healthcare AI Startup

---

### Cost and Resource Optimization

**Customer Need**: Maximize ROI on GPU and compute investments through efficient resource utilization.

**Requirements**:

- **Resource Packing**: Scheduler efficiently places multi-worker jobs to minimize fragmentation
- **GPU Utilization Visibility**: Metrics showing actual GPU usage across distributed workers
- **Job Preemption Support**: Ability to use spot/preemptible instances in cloud environments
- **Cost Attribution**: Chargeback/showback by team/project for MPIJob resource consumption
- **Idle Resource Alerts**: Notifications when GPUs are allocated but underutilized

**Market Context**: Customers report 30-50% idle GPU capacity due to poor job scheduling and resource fragmentation. An 8-GPU job may block cluster resources even though it only uses 60% GPU utilization. For customers with $2M+ GPU investments, optimization is a C-level concern.

**Optimization Strategies**:

1. **Gang Scheduling**: Ensure all workers for a job are scheduled simultaneously (prevents partial allocation tying up resources)
2. **Bin Packing**: Place jobs to minimize wasted GPU/CPU resources
3. **Topology-Aware Placement**: Co-locate workers on same network switch/rack to minimize communication latency
4. **Resource Right-Sizing**: Recommendations based on historical usage (e.g., "Your last 5 jobs used 40% of requested memory—consider reducing allocation")

**Implementation Considerations**:
- Integration with OpenShift scheduler for optimal placement
- Metrics dashboards showing resource efficiency trends
- Alerts for long-running jobs with low utilization

**Customer Quote**: *"We spent $2M on GPUs. If distributed training doesn't improve utilization beyond our current 55%, we can't justify the investment. Every percentage point matters."* — CTO, Healthcare AI Startup

---

### Enterprise Requirements Summary

| Requirement Category | MVP (Phase 1) | Post-MVP (Phase 2-3) |
|---------------------|---------------|---------------------|
| **Multi-Tenancy** | Namespace isolation, resource quotas | Fair scheduling, advanced chargeback |
| **Compliance** | Audit logging, RBAC | Policy enforcement (OPA), compliance reports |
| **Migration Support** | Documentation, examples | Migration assessment tool, professional services |
| **Support Readiness** | L2 training, runbooks | Predictive diagnostics, performance profiling tools |
| **Cost Optimization** | Basic resource metrics | Gang scheduling, right-sizing recommendations |

**Success Metric**: 80% of enterprise customers (defined as >1000 employees, regulated industry, or $100K+ ARR) successfully deploy MPIJobs in production within 6 months of GA.

---

*This RFE provides a comprehensive foundation for MPIJob support in Red Hat OpenShift AI. Implementation should proceed in phases, with continuous feedback loops from early adopters to refine requirements and address unforeseen challenges.*
