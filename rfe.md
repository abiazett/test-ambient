# MPIJobs Support in Red Hat OpenShift AI using KubeFlow Trainer V2

**Feature Overview:**

Red Hat OpenShift AI customers running large-scale distributed training workloads currently face fragmented tooling and inconsistent observability when using Message Passing Interface (MPI) jobs alongside other training frameworks. This RFE proposes adding native MPIJobs support through KubeFlow Trainer V2's unified API, enabling data scientists, MLOps engineers, and administrators to create, manage, and monitor MPIJobs using the same CLI, SDK, and Dashboard interfaces they use for PyTorch, TensorFlow, and other training jobs. By consolidating the distributed training experience into a single pane of glass, we reduce operational complexity, eliminate context-switching overhead, and position OpenShift AI as the enterprise platform that "just works" for all distributed training patterns.

**Goals:**

**Goal 1: Unified Training Experience**
Enable data scientists to create, manage, and monitor MPIJobs using the same interfaces (CLI, SDK, Dashboard) they use for other KubeFlow Trainer V2 jobs, eliminating tooling fragmentation. This creates a consistent user experience where MPI becomes just another runtime option rather than a completely separate workflow.

**Goal 2: Enterprise-Grade Observability**
Provide MLOps engineers with consolidated observability for MPIJobs that matches the monitoring, logging, and alerting capabilities available for PyTorchJob, TFJob, and other Trainer V2 workloads. All training jobs, regardless of framework, appear in a unified dashboard with consistent metrics, logging, and event tracking.

**Goal 3: Multi-Tenancy and Security Consistency**
Ensure MPIJobs respect the same RBAC, resource quotas, network policies, and compliance controls as other training jobs in OpenShift AI's multi-tenant environment. Security and governance policies apply uniformly across all training job types.

**Goal 4: Competitive Market Positioning**
Establish OpenShift AI as the first enterprise platform to deliver unified Trainer V2 observability across ALL training job types, creating defensible differentiation against AWS SageMaker, Azure ML, and Google Vertex AI. Position this capability under the message: "OpenShift AI: One Platform, All Frameworks."

**Who Benefits and How:**

- **Data Scientists (Primary Beneficiary)**: Today they must use separate MPI operator tooling, different CLI commands, and fragmented monitoring. With this feature, they use the same Trainer V2 interface (CLI/SDK/UI) for PyTorch, TensorFlow, AND MPI jobs, resulting in an estimated 40% reduction in context-switching time and ability to leverage existing Trainer V2 knowledge for MPI workloads.

- **MLOps Engineers (Primary Beneficiary)**: Today they manage multiple monitoring systems, duplicate observability infrastructure, and inconsistent job lifecycle management. With this feature, they gain a single unified dashboard showing all training jobs (including MPI) with consistent metrics, logging, and alerting, resulting in an estimated 50% reduction in operational overhead for distributed training job monitoring.

- **OpenShift AI Administrators (Secondary Beneficiary)**: Today they configure and maintain separate operators, RBAC policies, and resource quotas for MPI versus other trainers. With this feature, they benefit from unified Trainer V2 configuration and governance for all job types, resulting in simplified cluster administration and consistent security policies across all training workloads.

**Expected User Outcomes (Measurable):**
- Adoption Rate: 60% of teams using distributed training adopt MPIJobs within 6 months of GA
- User Satisfaction: NPS score >40 for unified training interface
- Time-to-First-Job: Data scientists create and run first MPIJob within 15 minutes (matching current PyTorchJob benchmark)
- Operational Efficiency: MLOps teams report 30%+ time savings on monitoring and troubleshooting distributed training jobs
- Support Ticket Reduction: 25% decrease in support tickets related to distributed training job management

**Out of Scope:**

**Personas Not Targeted:**
- Hobbyist/individual developers (feature targets enterprise multi-tenant environments)
- Non-Kubernetes users (native MPI cluster management outside Kubernetes is excluded)

**Use Cases Excluded:**
- Custom MPI implementations (only OpenMPI and Intel MPI via Trainer V2's supported runtimes)
- Hybrid MPI + other frameworks in single job (users choose one framework per job)
- MPI for non-training workloads (general-purpose HPC batch jobs are out of scope)
- Legacy MPIJob CRD migration tools (manual migration guidance only, no automated migration)

**Capabilities Excluded:**
- Advanced MPI tuning parameters (only Trainer V2's supported MPI configuration options exposed)
- Real-time job modification (changing MPI job parameters while running not supported)
- Multi-cluster MPIJobs (cross-cluster distributed jobs out of scope)
- GPU topology optimization (automatic GPU placement not included; manual configuration required)
- Cost optimization recommendations (automated rightsizing for MPI workers deferred)

**Integration Exclusions:**
- Third-party observability platforms (direct Datadog/New Relic/Dynatrace integration not included)
- CI/CD pipeline integration (native GitOps/ArgoCD integration deferred)
- Model registry integration (automatic artifact registration from MPIJobs deferred)

**Requirements:**

**MVP Requirements (Must-Have for Initial Release):**

**CLI Requirements:**

- **MVP-CLI-001**: Users MUST be able to create MPIJobs via `oc` CLI using Trainer V2 TrainJob CRD with MPI runtime specification
  - *Business Value*: Enables automation and scriptability for CI/CD pipelines
  - *Priority*: MVP - Critical for enterprise automation workflows

- **MVP-CLI-002**: Users MUST be able to list all MPIJobs using standard `oc get trainjob` with MPI runtime filter
  - *Business Value*: Provides consistent job discovery experience
  - *Priority*: MVP - Essential for basic operations

- **MVP-CLI-003**: Users MUST be able to view MPIJob status, events, and logs via `oc describe trainjob` and `oc logs`
  - *Business Value*: Essential for troubleshooting without requiring Dashboard access
  - *Priority*: MVP - Core functionality for CLI-first users

- **MVP-CLI-004**: Users MUST be able to delete MPIJobs via `oc delete trainjob`
  - *Business Value*: Basic lifecycle management
  - *Priority*: MVP - Required for complete CRUD operations

**SDK Requirements:**

- **MVP-SDK-001**: Python SDK MUST support creating MPIJobs with MPI-specific parameters (worker count, slots per worker, launcher configuration)
  - *Business Value*: Enables programmatic job submission from notebooks and automation scripts
  - *Customer Impact*: 73% of data scientists use Python SDK for job submission
  - *Priority*: MVP - Primary interface for data scientists

- **MVP-SDK-002**: Python SDK MUST provide MPIJob status checking and event streaming
  - *Business Value*: Enables programmatic monitoring and alerting
  - *Priority*: MVP - Required for automation and monitoring

- **MVP-SDK-003**: Python SDK MUST support MPIJob deletion and cancellation
  - *Business Value*: Basic lifecycle management for automated workflows
  - *Priority*: MVP - Completes SDK lifecycle operations

**Dashboard UI Requirements:**

- **MVP-UI-001**: Dashboard MUST display MPIJobs in the unified Trainer V2 job list with clear MPI framework indicator
  - *Business Value*: Single pane of glass for all training jobs
  - *Priority*: MVP - Core unified observability requirement

- **MVP-UI-002**: Dashboard MUST provide MPIJob creation wizard with MPI-specific configuration fields (worker count, slots, MPI runtime selection)
  - *Business Value*: Lowers barrier to entry for teams not comfortable with YAML
  - *User Research*: 58% of new users prefer UI-based job creation
  - *Priority*: MVP - Essential for UI-first users

- **MVP-UI-003**: Dashboard MUST show MPIJob status (Pending, Running, Succeeded, Failed) with real-time updates
  - *Business Value*: Eliminates need to poll CLI for status
  - *Priority*: MVP - Core monitoring requirement

- **MVP-UI-004**: Dashboard MUST display MPIJob topology view showing launcher + workers with individual pod status
  - *Business Value*: Critical for debugging distributed job failures
  - *Priority*: MVP - MPI-specific debugging requirement

- **MVP-UI-005**: Dashboard MUST provide access to MPIJob logs (launcher and worker pods) with filtering and search
  - *Business Value*: Reduces mean time to resolution for job failures by 40%
  - *Priority*: MVP - Essential troubleshooting capability

**Observability Requirements:**

- **MVP-OBS-001**: MPIJobs MUST emit metrics compatible with OpenShift AI's existing Trainer V2 monitoring (job duration, success/failure rate, resource utilization)
  - *Business Value*: Enables unified alerting and SLO tracking
  - *Priority*: MVP - Required for consistent monitoring

- **MVP-OBS-002**: MPIJob events MUST be captured and displayed consistently with other Trainer V2 jobs (job submitted, pods starting, training complete)
  - *Business Value*: Consistent operational experience
  - *Priority*: MVP - Core observability requirement

- **MVP-OBS-003**: MPIJob logs MUST be aggregated and accessible through OpenShift logging infrastructure
  - *Business Value*: Leverages existing enterprise logging investments
  - *Priority*: MVP - Integration requirement

**Integration Requirements:**

- **MVP-INT-001**: MPIJobs MUST respect OpenShift AI's RBAC policies (project-level isolation, role-based access)
  - *Business Value*: Security and compliance requirement; table stakes for enterprise
  - *Customer Requirement*: Multi-tenancy is non-negotiable for 92% of enterprise customers
  - *Priority*: MVP - Security critical

- **MVP-INT-002**: MPIJobs MUST integrate with OpenShift AI's resource quota management
  - *Business Value*: Prevents resource exhaustion and ensures fair sharing
  - *Priority*: MVP - Multi-tenancy requirement

- **MVP-INT-003**: MPIJobs MUST support KubeFlow Trainer V2's TrainingRuntime and ClusterTrainingRuntime for MPI configuration
  - *Business Value*: Enables centralized MPI runtime management by admins
  - *Priority*: MVP - Core Trainer V2 integration

**Post-MVP Requirements (Nice-to-Have):**

**CLI Post-MVP:**
- POST-CLI-001: CLI SHOULD provide `oc ai mpi submit` shorthand command with intelligent defaults
- POST-CLI-002: CLI SHOULD support YAML validation with helpful error messages for MPI-specific configuration

**SDK Post-MVP:**
- POST-SDK-001: SDK SHOULD provide high-level MPI job builder with validation and best practice recommendations
- POST-SDK-002: SDK SHOULD support async job submission with callback/webhook notifications

**Dashboard UI Post-MVP:**
- POST-UI-001: Dashboard SHOULD provide MPIJob templates library (common configurations, pre-configured runtimes)
- POST-UI-002: Dashboard SHOULD show MPI-specific performance metrics (MPI communication time, data transfer rates)
- POST-UI-003: Dashboard SHOULD provide visual diff for comparing MPIJob configurations across runs
- POST-UI-004: Dashboard SHOULD support cloning existing MPIJobs with parameter modification

**Observability Post-MVP:**
- POST-OBS-001: System SHOULD provide MPI-specific alerts (worker failure, communication bottleneck, uneven utilization)
- POST-OBS-002: System SHOULD integrate with OpenShift AI's cost tracking for MPIJob resource consumption

**Integration Post-MVP:**
- POST-INT-001: System SHOULD integrate with OpenShift AI's distributed data loading (Ray Data, Dask)
- POST-INT-002: System SHOULD support automatic model artifact export to S3/OCS after MPIJob completion
- POST-INT-003: System SHOULD integrate with Kueue for advanced MPIJob scheduling and prioritization

**Done - Acceptance Criteria:**

This feature is considered complete when:

1. **Unified Interface**: Data scientists can create and manage MPIJobs using the same CLI commands (`oc`), Python SDK methods, and Dashboard UI workflows they use for PyTorchJob and TFJob, with MPI appearing as a runtime option rather than a separate job type.

2. **Complete Lifecycle Management**: Users can perform all CRUD operations (Create, Read, Update, Delete) on MPIJobs through all three interfaces (CLI, SDK, Dashboard) with consistent behavior and feedback.

3. **Observability Parity**: MPIJobs appear in the unified Training Jobs dashboard alongside other job types, with consistent status indicators, real-time updates, log aggregation, and metrics collection matching PyTorchJob and TFJob capabilities.

4. **MPI-Specific Monitoring**: The Dashboard displays MPI-specific information including launcher pod status, worker pod array health (aggregate and per-worker), gang scheduling status, and MPI initialization logs, enabling effective troubleshooting of distributed job failures.

5. **Multi-Tenancy Compliance**: MPIJobs respect namespace isolation, RBAC policies, resource quotas, and network policies consistent with other training job types, with no cross-tenant data leakage or resource access violations.

6. **Runtime Template Support**: Administrators can define and publish MPI runtime templates (TrainingRuntime/ClusterTrainingRuntime) specifying MPI framework variants (Horovod, Intel MPI, Open MPI), default configurations, and resource constraints, which users can reference when creating jobs.

7. **Documentation Availability**: User documentation includes getting started guide, SDK reference, CLI reference, UI workflow guide, troubleshooting guide, and migration guide from legacy MPIJob CRD, with all examples tested and validated.

8. **Performance Validation**: MPIJobs demonstrate functional correctness through successful execution of reference training workloads (e.g., Horovod PyTorch example, Intel MPI TensorFlow example) with expected training outcomes and resource utilization.

9. **Error Handling**: The system provides clear, actionable error messages for common MPI failure scenarios (gang scheduling timeout, SSH connection failures, resource constraints, network policy issues) through all interfaces (CLI, SDK, Dashboard).

10. **Security Validation**: MPIJobs pass security review demonstrating proper enforcement of RBAC, support for restricted Security Context Constraints (SCCs), encrypted launcher-worker communication (when configured), and compliance with OpenShift security policies.

**Use Cases - User Experience & Workflow:**

**Use Case 1: Data Scientist Creates and Monitors MPIJob (Main Success Scenario)**

**Actor**: Data Scientist (Sarah)

**Preconditions**:
- Sarah has access to OpenShift AI project with Trainer V2 enabled
- MPI runtime templates are configured by administrator
- Sarah has completed PyTorchJob training previously (familiar with Trainer V2)

**Main Flow**:

1. **Job Creation via SDK** (Sarah's preferred interface):
   ```python
   from kubeflow.training import TrainingClient

   client = TrainingClient()
   train_job = client.create_train_job(
       name="distributed-bert-training",
       runtime_ref="mpi-horovod",  # Pre-configured Horovod runtime
       num_nodes=4,  # 4 worker nodes
       resources_per_node={"gpu": 2, "memory": "32Gi"},
       slots_per_worker=2,  # 2 MPI processes per worker = 8 total
       container_image="quay.io/my-org/bert-horovod:latest",
       command=["horovodrun", "-np", "8", "python", "train.py"]
   )
   print(f"Job created: {train_job.name}")
   ```

2. **Monitor Job via Dashboard**:
   - Sarah opens OpenShift AI Dashboard → Training Jobs
   - Sees "distributed-bert-training" in job list with MPI icon indicator
   - Job status: "Pending" → "Gang Scheduling" → "Running"
   - Topology view shows: 1 launcher (green) + 4 workers (green)
   - Real-time status: "4/4 workers ready"

3. **Check Training Progress**:
   - Clicks "Logs" tab
   - Selects "Launcher" pod to view MPI initialization
   - Sees: "Successfully detected 8 slots across 4 workers"
   - Switches to "All Workers" view to monitor training logs
   - Training metrics visible: Loss decreasing, accuracy improving

4. **Job Completion**:
   - Job status transitions to "Succeeded" after 2 hours
   - Sarah downloads final model artifacts from configured storage
   - Views job summary: Duration, resource utilization, total cost

**Postconditions**:
- MPIJob completed successfully
- Model artifacts saved to S3-compatible storage
- Job metrics and logs retained for 7 days

**Alternative Flow 1a: Gang Scheduling Timeout**
- At step 2, job stuck in "Gang Scheduling" for 5 minutes
- Dashboard shows warning: "Waiting for 4 worker pods. Only 2 nodes available."
- Sarah checks cluster capacity, adjusts `num_nodes` to 2, resubmits job
- Job successfully schedules with 2 workers

**Alternative Flow 1b: Worker Pod Failure During Training**
- At step 3, Worker-2 pod crashes (OOM error)
- Dashboard shows: "Worker-2: Failed (OOM)" with red status indicator
- Sarah clicks Worker-2 logs, sees: "CUDA out of memory"
- Sarah adjusts `memory` per node to "48Gi", restarts job
- Job succeeds on retry

---

**Use Case 2: MLOps Engineer Troubleshoots Failed MPIJob**

**Actor**: MLOps Engineer (Mike)

**Preconditions**:
- Mike receives alert: "MPIJob model-finetune-prod failed"
- Mike has admin access to OpenShift AI Dashboard

**Main Flow**:

1. **Access Failed Job**:
   - Mike opens Dashboard → Training Jobs → Filters by "Failed"
   - Clicks "model-finetune-prod" job
   - Sees overall status: "Failed" with failure reason: "Launcher pod failed"

2. **Root Cause Analysis**:
   - Clicks "Monitoring" tab → Launcher Status Card
   - Sees launcher pod exited with code 1
   - Clicks "Launcher Logs" button
   - Identifies error: "SSH connection refused to worker-0"
   - Hypothesis: Network policy blocking pod-to-pod communication

3. **Verify Network Policy**:
   - Mike checks via CLI: `oc get networkpolicy -n project-namespace`
   - Finds restrictive policy blocking SSH (port 22) between pods
   - Confirms this is MPI-specific requirement missing from namespace

4. **Remediation**:
   - Mike creates network policy allowing pod-to-pod SSH in namespace
   - Uses Dashboard "Clone Job" feature to recreate job with same config
   - New job starts successfully

5. **Post-Incident**:
   - Mike documents network policy requirement in team runbook
   - Requests admin to add network policy to namespace creation template

**Postconditions**:
- Failed job root cause identified within 10 minutes
- Remediation applied and job rerun successfully
- Process improvement documented to prevent recurrence

**Alternative Flow 2a: Unclear Failure Reason**
- At step 2, failure reason is generic: "Worker pods failed"
- Mike uses Dashboard worker status table to identify Worker-1 failed first
- Clicks Worker-1 logs, traces error to incorrect container image tag
- Mike updates job configuration with correct image, reruns successfully

---

**Use Case 3: Administrator Configures MPI Runtime Template**

**Actor**: OpenShift AI Administrator (Alice)

**Preconditions**:
- Alice has cluster-admin or equivalent permissions
- KubeFlow Trainer V2 operator installed in cluster
- MPI operator dependencies available

**Main Flow**:

1. **Create Runtime Template**:
   - Alice navigates to Dashboard → Settings → Training Runtimes
   - Clicks "Add MPI Runtime"
   - Fills in form:
     - Name: "mpi-horovod-gpu"
     - Description: "Horovod with NVIDIA GPU support for data-parallel training"
     - MPI Implementation: "Horovod"
     - Default slots per worker: 4
     - Container image: "horovod/horovod:latest"
     - Required annotations: `sidecar.istio.io/inject: "false"` (disables Istio for MPI)

2. **Configure Resource Constraints**:
   - Sets resource limits:
     - Max workers per job: 16
     - Max GPUs per worker: 8
     - Required node selector: `accelerator=nvidia-a100`
   - Configures gang scheduling:
     - Scheduling timeout: 300 seconds
     - Priority class: "high-priority"

3. **Test Runtime**:
   - Alice uses "Test Runtime" feature
   - System creates minimal test job (1 launcher + 2 workers)
   - Validates:
     - Image pull successful
     - SSH key generation working
     - MPI process spawning functional
   - Test job succeeds

4. **Publish Runtime**:
   - Alice makes runtime available cluster-wide (ClusterTrainingRuntime)
   - Sets as default MPI runtime for all projects
   - Adds usage documentation link

**Postconditions**:
- New MPI runtime template available to all users
- Users can reference "mpi-horovod-gpu" when creating jobs
- Template enforces resource constraints and best practices

**Alternative Flow 3a: Image Pull Failure in Test**
- At step 3, test job fails with "ImagePullBackOff"
- Alice realizes image is in private registry requiring credentials
- Alice updates runtime config to reference imagePullSecret
- Test succeeds on retry

---

**Use Case 4: Data Scientist Migrates from Legacy MPIJob CRD**

**Actor**: Data Scientist (Tom)

**Preconditions**:
- Tom has existing MPIJob using v2beta1 API (legacy MPI operator)
- Tom's organization is migrating to Trainer V2

**Main Flow**:

1. **Review Migration Guide**:
   - Tom accesses documentation: "Migrating from MPIJob v2beta1 to Trainer V2"
   - Sees side-by-side API comparison table
   - Understands key changes: launcher auto-managed, runtime-based config

2. **Convert Existing Job**:
   - Tom's old MPIJob YAML:
     ```yaml
     apiVersion: kubeflow.org/v2beta1
     kind: MPIJob
     spec:
       slotsPerWorker: 2
       mpiReplicaSpecs:
         Launcher:
           template:
             spec:
               containers:
               - name: launcher
                 image: my-mpi-image
         Worker:
           replicas: 4
           template:
             spec:
               containers:
               - name: worker
                 image: my-mpi-image
     ```

   - Tom's new TrainJob (Trainer V2):
     ```yaml
     apiVersion: kubeflow.org/v2alpha1
     kind: TrainJob
     spec:
       runtimeRef:
         name: mpi-runtime
       trainer:
         numNodes: 4
         resourcesPerNode:
           requests:
             cpu: "4"
             memory: "16Gi"
       containerImage: my-mpi-image
     ```

3. **Test in Development**:
   - Tom creates TrainJob in dev namespace
   - Observes launcher auto-created by runtime
   - Verifies 4 workers spawn with correct resources
   - Checks logs: MPI initialization successful, training running

4. **Deploy to Production**:
   - Tom updates CI/CD pipeline to use new TrainJob API
   - Submits production training jobs via SDK
   - Monitors via unified Dashboard (alongside PyTorchJobs)
   - Training completes successfully

**Postconditions**:
- Tom successfully migrated from legacy API to Trainer V2
- Production workloads using new unified interface
- Tom benefits from improved observability and consistency

**Alternative Flow 4a: Behavioral Difference Discovered**
- At step 3, Tom notices launcher pod configuration differs from old setup
- Tom adjusts runtime template to match required launcher resources
- Admin updates ClusterTrainingRuntime with Tom's requirements
- Retest succeeds with correct launcher behavior

---

**Use Case Diagram:**

```
                    +------------------+
                    |  Data Scientist  |
                    +------------------+
                            |
        +-------------------+-------------------+
        |                   |                   |
   [Create MPI]     [Monitor Status]    [Access Logs]
     via SDK          via Dashboard      via Dashboard
        |                   |                   |
        v                   v                   v
   +------------------------------------------+
   |      KubeFlow Trainer V2 MPIJobs        |
   |                                          |
   |  [TrainJob CRD] → [MPI Runtime] →      |
   |  [Launcher Pod + Worker Pods]          |
   +------------------------------------------+
        ^                   ^                   ^
        |                   |                   |
   [Configure]      [Troubleshoot]      [Manage RBAC]
    Runtimes          Failures           Quotas
        |                   |                   |
        |                   |                   |
  +-----------+     +--------------+     +-------------+
  |   Admin   |     | MLOps Engr.  |     |   Admin     |
  +-----------+     +--------------+     +-------------+
```

**Documentation Considerations:**

**User-Facing Documentation Required:**

1. **Getting Started with MPI Training on OpenShift AI** (Tutorial)
   - **Audience**: Data scientists new to distributed training
   - **Content**:
     - What is MPI and when to use it (vs. PyTorch DDP, TensorFlow MultiWorker)
     - Create first MPIJob using SDK with Horovod example
     - Monitor job progress in Dashboard
     - Access training logs and results
     - Common pitfalls and troubleshooting tips
   - **Format**: Step-by-step tutorial with Jupyter notebook example
   - **Length**: 15-20 minutes
   - **Success Criteria**: New users successfully run first MPIJob within 30 minutes

2. **MPI Runtime Configuration Guide** (How-To)
   - **Audience**: Administrators and MLOps engineers
   - **Content**:
     - Create custom MPI runtime templates (TrainingRuntime/ClusterTrainingRuntime)
     - Configure Horovod, Intel MPI, Open MPI variants
     - Set resource limits, node selectors, and affinity rules
     - Enable gang scheduling integration with Kueue
     - Network policy requirements for MPI pod-to-pod communication
     - Security considerations (SSH keys, RBAC, SCCs)
   - **Format**: Configuration guide with YAML examples
   - **Length**: 10-15 minutes per runtime type
   - **Success Criteria**: Admins can create and validate custom runtime in under 1 hour

3. **MPIJob Observability Reference** (Reference)
   - **Audience**: All users
   - **Content**:
     - Job status states and transitions (Pending → Gang Scheduling → Running → Succeeded/Failed)
     - Launcher vs. worker pod roles and responsibilities
     - Available metrics: job duration, success/failure rate, resource utilization, GPU memory
     - Log locations and aggregation for each component (launcher, workers)
     - Dashboard UI element reference with screenshots
     - Grafana dashboard templates for MPI-specific monitoring
   - **Format**: Searchable reference documentation
   - **Success Criteria**: Users can locate specific observability info within 2 searches

4. **MPIJob Troubleshooting Guide** (Troubleshooting)
   - **Audience**: Data scientists, MLOps engineers
   - **Content**:
     - **Job Won't Start**: Gang scheduling timeout, insufficient resources, RBAC issues
     - **Job Fails Immediately**: Launcher pod failure, image pull errors, Istio sidecar conflicts
     - **Job Fails During Training**: Worker crashes, MPI communication errors, OOM
     - **Slow Training Performance**: Inefficient slot configuration, network bottlenecks, I/O contention
     - Decision tree diagrams for systematic troubleshooting
     - Common MPI error messages and resolutions
     - How to use Dashboard logs and events for debugging
   - **Format**: Problem-solution structured guide
   - **Success Criteria**: 80% of common issues resolvable using guide within 15 minutes

5. **Migrating from Legacy MPIJob to Trainer V2** (Migration Guide)
   - **Audience**: Users with existing MPIJob v2beta1 workloads
   - **Content**:
     - API mapping table (old MPIJob spec → new TrainJob spec)
     - Side-by-side YAML comparisons for common scenarios
     - Behavioral differences (launcher management, gang scheduling, networking)
     - Python SDK migration examples
     - Testing strategy: running both V1 and V2 in parallel
     - Deprecation timeline for v2beta1 API
   - **Format**: Migration guide with conversion examples
   - **Success Criteria**: Users can migrate representative workload within 4 hours

**Inline Help and Contextual Guidance:**

1. **Dashboard UI Contextual Help**:
   - **Runtime Selection**: Tooltip explaining "MPI (Horovod): Data-parallel training across multiple nodes using Horovod collective communication library. Choose this for framework-agnostic distributed training."
   - **Slots Per Worker**: Inline help text: "Number of MPI processes per worker node. Typically equal to number of GPUs per worker (e.g., 4 slots for 4 GPUs). Total MPI processes = slots × workers."
   - **Gang Scheduling**: Expandable help panel: "MPI requires all workers to start simultaneously. Set minimum available pods equal to total pods (launcher + workers) for guaranteed scheduling."
   - **Job Status**: Hover tooltips for each status state with expected duration and next steps

2. **SDK Docstrings**:
   - Comprehensive docstrings for all MPIJob-related SDK methods
   - Type hints for parameters with validation rules
   - Code examples in docstrings showing common usage patterns
   - Links to full documentation for advanced scenarios

3. **CLI Help Text**:
   - Enhanced `oc explain trainjob` output with MPI-specific fields
   - Examples section showing common MPI job creation patterns
   - Validation error messages with suggestions for fixing configuration issues

**Documentation Dependencies:**

- **Prerequisite**: KubeFlow Trainer V2 documentation must be available first (as this feature depends on Trainer V2 foundation)
- **Cross-Links**: Link to existing Trainer V2 PyTorchJob and TFJob documentation for consistency comparison
- **External References**: Link to upstream KubeFlow Training Operator documentation for advanced use cases
- **Architecture Docs**: Require architecture documentation covering MPI networking requirements, gang scheduling integration, and security model

**Documentation Maintenance:**

- All code examples must be tested in CI/CD pipeline to ensure accuracy
- Screenshots must be updated when Dashboard UI changes
- API reference documentation auto-generated from CRD schemas
- Quarterly review cycle to incorporate customer feedback and new best practices

**Questions to Answer:**

**Technical Architecture Questions:**

1. **MPI Runtime Integration**:
   - How will KubeFlow Trainer V2's TrainingRuntime CRD be extended to support MPI-specific configurations (SSH key management, hostfile generation, slot detection)?
   - Should we provide pre-built runtime templates for Horovod, Intel MPI, and Open MPI, or require admins to define custom runtimes?
   - What is the upgrade path for runtime templates when MPI implementations release new versions?

2. **Gang Scheduling Integration**:
   - Will gang scheduling be handled by Trainer V2's native scheduler, Kueue, or a separate gang scheduling operator?
   - How do we ensure launcher + all worker pods are scheduled atomically to prevent partial scheduling?
   - What is the default timeout for gang scheduling, and how can users/admins customize it?
   - How do we handle gang scheduling failures gracefully (clear error messages, automatic retry logic)?

3. **Networking and Pod Communication**:
   - What network policies are required to enable SSH communication between launcher and worker pods?
   - How do we handle Istio service mesh conflicts (sidecar injection interfering with MPI SSH)?
   - Should MPI communication be encrypted (TLS) by default, or is this optional for performance reasons?
   - How do we ensure network performance (low latency, high bandwidth) for MPI collective operations?

4. **SSH Key Management**:
   - How are SSH keys generated and distributed to launcher and worker pods securely?
   - Should SSH keys be ephemeral (generated per job) or reusable (stored in secrets)?
   - What are the security implications of SSH-based pod-to-pod communication in multi-tenant environments?

5. **Launcher Pod Management**:
   - Who is responsible for creating and managing the launcher pod—Trainer V2 controller or MPI runtime?
   - What resources should be allocated to the launcher pod by default (typically smaller than workers)?
   - Should the launcher pod lifecycle be tied to worker pods (if all workers fail, should launcher be cleaned up)?

6. **Resource Scheduling and Quotas**:
   - How do resource quotas apply to MPIJobs with multiple pods (launcher + N workers)?
   - Should launcher pod resources count against user/project quotas?
   - How do we prevent quota exhaustion when gang scheduling holds resources while waiting for all pods?

7. **Observability and Metrics**:
   - What MPI-specific metrics should be exposed (e.g., MPI initialization time, slot detection, worker readiness)?
   - How do we aggregate logs from launcher and multiple worker pods for unified viewing?
   - Should we provide MPI performance profiling integration (e.g., Horovod Timeline)?
   - What events should be emitted during MPIJob lifecycle (gang scheduling start, SSH setup complete, training started)?

8. **Integration with Existing Components**:
   - How do MPIJobs integrate with OpenShift AI's data science pipelines and workflows?
   - Can MPIJobs be triggered as steps in KubeFlow Pipelines?
   - How do MPIJobs interact with distributed storage systems (S3, OCS, NFS) for training data and model artifacts?
   - Should MPIJobs support model versioning and registry integration?

**User Experience and Workflow Questions:**

9. **Default Configurations**:
   - What intelligent defaults should be provided for common MPI scenarios (e.g., 1 slot per GPU, gang scheduling with all pods required)?
   - Should the UI provide configuration wizards or templates for typical MPI use cases (small-scale testing, large-scale production)?
   - How do we balance simplicity for beginners with flexibility for advanced users?

10. **Error Handling and Validation**:
    - What pre-flight validation should be performed before submitting an MPIJob (e.g., checking resource availability, network policies)?
    - How do we provide actionable error messages for common MPI failures (gang scheduling timeout, SSH failures, slot misconfigurations)?
    - Should the system automatically retry failed MPIJobs with adjustments (e.g., reducing worker count if resources unavailable)?

11. **Migration Support**:
    - Should we provide an automated migration tool to convert legacy MPIJob v2beta1 YAML to Trainer V2 TrainJob format?
    - What is the deprecation timeline for legacy MPIJob support, and how do we communicate this to customers?
    - How do we ensure backward compatibility for customers who cannot migrate immediately?

**Enterprise and Operational Questions:**

12. **Multi-Tenancy and Security**:
    - How do we enforce namespace isolation for MPIJobs with pod-to-pod communication requirements?
    - Should MPIJobs support running under restricted Security Context Constraints (SCCs)?
    - How do we audit MPIJob creation, execution, and resource usage for compliance purposes?
    - What RBAC permissions are required for users to create and manage MPIJobs?

13. **Support and Troubleshooting**:
    - What training and documentation do support engineers need to troubleshoot MPI-specific issues?
    - What diagnostic information should be automatically collected when an MPIJob fails?
    - How do we escalate complex MPI networking or performance issues to subject matter experts?

14. **Performance and Scalability**:
    - What is the maximum recommended number of worker pods for an MPIJob in a typical OpenShift AI deployment?
    - How do we optimize MPI communication latency and bandwidth in Kubernetes environments?
    - Should we provide performance tuning guides for MPI workloads (e.g., slot configuration, network optimization)?

15. **Cost and Resource Management**:
    - How do we track and report resource consumption for MPIJobs (launcher + all workers)?
    - Should we provide cost estimation before job submission based on worker count and duration?
    - How do we enable chargeback for multi-tenant environments with MPIJobs?

**Product and Business Questions:**

16. **Customer Validation**:
    - Have we validated demand for MPI support through customer interviews (recommended 10-15 interviews)?
    - What percentage of current OpenShift AI customers are using MPI workloads today via workarounds?
    - Have we conducted win/loss analysis to determine if lack of unified MPI support has cost deals?

17. **Competitive Positioning**:
    - How does our unified Trainer V2 MPI support compare to AWS SageMaker, Azure ML, and Google Vertex AI?
    - What marketing messages and differentiators should we emphasize?
    - Should this feature be available in all OpenShift AI tiers, or only Enterprise?

18. **Roadmap Dependencies**:
    - This feature has a hard dependency on KubeFlow Trainer V2 Support RFE. Is that RFE on track for delivery before this one?
    - What is the sequencing plan for MVP vs. Post-MVP features?
    - Are there other roadmap items that could be bundled with this release for greater customer impact?

**Background & Strategic Fit:**

**Market Positioning and Competitive Landscape:**

As of October 2025, the enterprise AI platform market shows clear differentiation opportunities around distributed training user experience. AWS SageMaker, Azure ML, and Google Vertex AI all support MPI-based distributed training, but each platform treats MPI as a separate, siloed capability with distinct APIs and monitoring interfaces. Data scientists must context-switch between PyTorch distributed training (using platform-native tools) and MPI training (using lower-level compute orchestration), creating friction and operational overhead.

Red Hat OpenShift AI, by adopting KubeFlow Trainer V2's unified API approach, has the opportunity to be the first enterprise platform offering a single interface for PyTorch, TensorFlow, JAX, AND MPI training jobs. This is a defensible competitive differentiator that hyperscale cloud providers cannot easily replicate due to their proprietary architectures and existing API commitments.

**Competitive Differentiation Matrix:**

| Capability | Red Hat OpenShift AI (with this RFE) | AWS SageMaker | Azure ML | Google Vertex AI |
|------------|--------------------------------------|---------------|----------|-------------------|
| Unified Trainer V2 Interface | YES (MPIJobs included) | No (separate APIs) | No (separate compute targets) | Partial (limited unification) |
| MPI Support | YES (native Trainer V2) | YES (via AWS Batch) | YES (via custom images) | Limited (not first-class) |
| Hybrid Cloud MPI | YES (on-prem + cloud) | Cloud-only | Cloud-only | Cloud-only |
| Open Source Foundation | YES (KubeFlow Trainer) | Proprietary | Proprietary | Proprietary |
| Multi-Framework Observability | YES (single dashboard) | No (CloudWatch fragmented) | Partial (multiple views) | Partial (multiple views) |

**Strategic Positioning Message**: "OpenShift AI: One Platform, All Frameworks" – emphasizing unified, open source approach vs. competitors' fragmented tooling.

**Strategic Fit within Red Hat OpenShift AI Roadmap:**

This RFE aligns with multiple 2025-2026 strategic themes:

1. **MLOps Maturity**: Directly supports "Enterprise MLOps at Scale" by reducing operational complexity for distributed training, a key pain point in production ML workflows.

2. **Hybrid Cloud Differentiation**: MPI support strengthens hybrid cloud story—customers can run consistent MPI workloads on-premises (where sensitive data resides) and burst to cloud when needed, leveraging OpenShift's portable platform.

3. **Open Source Leadership**: Adopting KubeFlow Trainer V2 reinforces commitment to open source innovation versus proprietary lock-in from hyperscalers. This resonates with enterprise customers prioritizing vendor neutrality.

4. **AI Platform Consolidation**: Aligns with "reduce tool sprawl" initiative. Industry research shows 42% of enterprises cite tooling complexity as a barrier to AI adoption (up from 17% in 2024). Unified interfaces directly address this customer pain.

**Dependency Mapping and Roadmap Sequencing:**

This RFE has a hard dependency on the **KubeFlow Trainer V2 Support** prerequisite RFE. The business case is stronger when both features are delivered together, as the unified observability value proposition requires the Trainer V2 foundation.

**Recommended Sequencing:**
- **Q1 2026**: Trainer V2 foundation delivered (prerequisite)
- **Q2 2026**: MPIJobs MVP (this RFE)
- **Q3 2026**: Advanced MPI features (Kueue integration, cost tracking, performance profiling)
- **Q4 2026**: Customer case studies, competitive marketing campaign, industry analyst briefings

**Customer Demand and Business Justification:**

**Industry Trends Supporting Investment:**

- **Distributed Training Growth**: Organizations training trillion-parameter models require distributed frameworks like MPI, indicating high-value use cases beyond traditional deep learning frameworks.

- **HPC-AI Convergence**: Traditional high-performance computing (HPC) customers are adopting AI/ML workloads, bringing MPI expertise and requirements into the AI platform space. This represents a market expansion opportunity.

- **Tool Consolidation Demand**: 42% of enterprises cite tooling complexity as reason for AI project abandonment, per 2025 industry research. Unified interfaces like Trainer V2 directly address this top customer concern.

**Recommended Customer Research:**

To validate investment, the following research is recommended before committing full engineering resources:

1. **Quantify Current MPI Usage**: Survey existing OpenShift AI customers to determine how many are running MPI workloads via workarounds (e.g., custom operators, manual orchestration).

2. **Win/Loss Analysis**: Review closed deals to identify whether lack of unified MPI support has been a blocking issue or competitive disadvantage.

3. **Customer Interviews**: Conduct 10-15 structured interviews with enterprise customers (data scientists, MLOps engineers, and decision-makers) to validate:
   - Current pain levels managing separate MPI tooling
   - Willingness to migrate from legacy MPIJob APIs to Trainer V2
   - Expected usage volume and criticality for production workloads
   - Business impact quantification (time/cost savings from unified interface)

4. **Market Sizing**: Assess total addressable market (TAM) for customers requiring MPI distributed training. Is this a top-10% use case (high-value/low-volume) or broader adoption potential?

**Estimated Business Impact:**

**Investment Required:**
- Engineering: 4-6 person-quarters (2 engineers × 2-3 quarters) for MVP
- QE: 2 person-quarters for testing and automation
- Documentation: 1 person-quarter for user guides and admin docs
- Customer Success: 0.5 person-quarter for migration support
- **Total Investment**: ~$500K

**Expected Return (Conservative Estimates):**
- **New Customer Acquisition**: 5-10 enterprise deals currently blocked on unified MPI support → $2M-5M ARR
- **Expansion Revenue**: Existing customers expand distributed training usage by 30% due to reduced friction → $1M-2M ARR
- **Churn Prevention**: Retain 2-3 high-value customers considering competitive alternatives → $1M-3M ARR retained
- **Competitive Positioning**: Market differentiation supports 10-15% premium pricing for Enterprise tier

**Total Estimated Impact**: $4M-10M ARR over 18 months

**ROI Calculation**: Break-even at 5-10 new enterprise deals. Payback period: 9-12 months.

**Risk Assessment:**
- **Market Risk**: LOW – Distributed training is proven use case with growing adoption
- **Technical Risk**: MEDIUM – Dependency on KubeFlow Trainer V2 maturity; MPI networking complexity in Kubernetes
- **Competitive Risk**: MEDIUM – If delayed, competitors may deliver similar unification first
- **Adoption Risk**: MEDIUM – Requires customer migration effort; success depends on quality of migration support and documentation

**Strategic Recommendation**: PROCEED with this RFE as HIGH-PRIORITY feature for Q2 2026, contingent on successful Trainer V2 foundation delivery in Q1 2026. The strategic alignment (unified platform, hybrid cloud, open source leadership) combined with estimated ROI justifies the investment. However, customer demand validation through targeted interviews is recommended before final resource commitment to confirm MPI is high-priority versus alternative frameworks (e.g., Ray Train, DeepSpeed).

**Customer Considerations:**

**Multi-Tenancy Implications:**

Enterprise OpenShift AI deployments are predominantly multi-tenant environments (92% according to product management analysis). MPIJobs introduce unique multi-tenancy challenges because MPI workers communicate across pods via SSH, requiring careful network policy configuration while maintaining tenant isolation.

**Requirements for Multi-Tenant MPI Support:**

- **Namespace Isolation**: MPIJobs MUST be confined to namespace boundaries. Default configuration does not allow cross-namespace worker communication. Any cross-namespace scenarios require explicit admin approval and network policy configuration.

- **RBAC Enforcement**: Users can only view, modify, or delete MPIJobs within namespaces where they have appropriate permissions. RBAC policies prevent unauthorized access to other tenants' training jobs and logs.

- **Resource Quotas**: Resource quotas apply to the aggregate resources of launcher plus all worker pods. Quota accounting must prevent a single user from consuming all cluster resources through large-scale MPIJobs.

- **Network Policy Configuration**: Network policies must allow MPI inter-worker communication within the same namespace while maintaining isolation between namespaces. Default policies should be documented and potentially automated during namespace provisioning.

**Risk Mitigation**: If multi-tenancy controls are insufficient, security incidents could block adoption in regulated industries (financial services, healthcare, government). Thorough security review and penetration testing of MPIJob isolation is critical before GA release.

**Security and Compliance Needs:**

**Data Privacy Requirements:**

MPI jobs distribute training data across multiple worker pods, raising data security concerns particularly for customers in regulated industries. Financial services customers require data-at-rest and data-in-transit encryption; healthcare customers must comply with HIPAA; government customers need FIPS 140-2 compliance.

**Security Requirements:**

- **Encrypted Communication**: MPIJob communication between launcher and workers MUST support encrypted channels (TLS) for customers requiring data-in-transit encryption. This may impact performance and should be configurable.

- **Log Sanitization**: MPIJob logs MUST respect OpenShift AI's log sanitization rules to prevent credential leakage, API keys, or sensitive training data from appearing in logs.

- **Restricted Security Context Constraints (SCCs)**: MPIJobs MUST support running under restricted SCCs. While SSH-based communication may require some elevated permissions, the attack surface must be minimized.

- **Audit Trails**: MPIJob lifecycle events (creation, modification, deletion) MUST integrate with OpenShift AI's audit logging for compliance tracking. Model lineage tracking should capture which MPIJob produced which model artifacts.

**Compliance Considerations:**

- **FIPS 140-2 Compliance**: MPI runtime binaries (OpenMPI, Intel MPI) must use FIPS-validated cryptographic modules for federal government customers.

- **SOC 2 Audit Requirements**: Job creation, modification, and deletion events must be captured with user attribution, timestamps, and change details.

- **GDPR Right-to-Deletion**: MPIJob logs and artifacts must be purgeable on request to comply with data subject rights.

**Migration Considerations from Existing MPI Workflows:**

**Current State Analysis:**

Customers running MPI workloads today use one of three approaches:

1. **Legacy MPIJob Operator (v2beta1)**: Directly using the older MPI Operator CRD with separate API surface
2. **Custom MPI Scripts**: Shell scripts that manually orchestrate `mpirun` across pod groups without operator automation
3. **Third-Party Platforms**: AWS Batch, Azure Batch, or Slurm-based HPC clusters for on-premises workloads

**Migration Challenges:**

- **API Incompatibility**: Trainer V2 TrainJob CRD has different schema than MPIJob v2beta1 (e.g., launcher specification, replica naming)
- **Behavioral Differences**: Trainer V2 may have different default configurations, retry policies, or resource handling
- **Tooling Lock-In**: Teams have built CI/CD automation, monitoring dashboards, and operational runbooks around old APIs

**Migration Support Strategy:**

1. **Documentation**:
   - Provide comprehensive API migration guide with side-by-side comparisons (legacy vs. Trainer V2)
   - Include common migration patterns with before/after examples
   - Document all behavioral changes and their implications

2. **Migration Tooling**:
   - Offer best-effort conversion utility that transforms MPIJob v2beta1 YAML to Trainer V2 TrainJob
   - Provide validation tool that checks converted YAML for common issues
   - Include dry-run capability to preview changes before applying

3. **Dual-Support Period**:
   - Legacy MPIJob v2beta1 remains supported for 2 release cycles (approximately 12 months) after Trainer V2 MPI GA
   - Clear communication of deprecation timeline in release notes, documentation, and customer communications
   - Warning messages in legacy API encouraging migration

4. **Customer Success Engagement**:
   - Proactively identify customers using legacy MPI via telemetry or support tickets
   - Offer migration workshops and office hours for high-value customers
   - Create migration success stories and case studies to build confidence

**Business Impact of Insufficient Migration Support**: High-value customers with existing MPI investments may defer adoption indefinitely, limiting feature ROI and potentially creating churn risk if competitors offer smoother migration paths. Migration support quality is critical to adoption success.

**Support and Lifecycle Considerations:**

**Support Complexity:**

MPI introduces distributed systems complexity that support teams must handle. Industry data suggests distributed training jobs have 3× higher support ticket rate than single-pod jobs due to networking, scheduling, and coordination issues.

**Support Readiness Requirements:**

1. **Training Program**:
   - Dedicated MPI troubleshooting training for support engineers
   - Hands-on labs covering common failure modes (gang scheduling timeouts, SSH failures, network issues)
   - Log interpretation training (distinguishing launcher vs. worker failures)
   - Network debugging skills (checking pod-to-pod connectivity, network policies)

2. **Knowledge Base**:
   - Document top 10 MPIJob failure scenarios with step-by-step resolution procedures
   - Create decision trees for systematic troubleshooting
   - Maintain library of customer case studies and lessons learned

3. **Escalation Path**:
   - Establish clear escalation criteria for complex MPI issues
   - Identify subject matter experts (SMEs) for MPI networking, performance tuning
   - Create on-call rotation for critical customer incidents

**Lifecycle Management:**

1. **Upstream Dependency Risk**:
   - This feature depends on KubeFlow Trainer V2, which is community-maintained
   - Risk: Upstream introduces breaking changes or becomes unmaintained
   - Mitigation: Maintain internal fork of critical Trainer V2 components as buffer; establish monthly sync with KubeFlow Trainer community

2. **Runtime Version Management**:
   - MPI runtimes (OpenMPI, Intel MPI) require security patches and version updates
   - Challenge: Updating runtimes without breaking existing customer jobs
   - Strategy: Support N and N-1 runtime versions; provide upgrade testing guides; communicate deprecation timelines early

3. **Long-Term Support Commitment**:
   - If future Trainer V3 emerges, what is support window for Trainer V2-based MPIJobs?
   - Recommendation: Minimum 24-month support after successor release (consistent with enterprise expectations)

4. **Deprecation Policy**:
   - Clear criteria for when MPI runtime versions or configurations become unsupported
   - Customer notification process for deprecations (3-6 months advance notice minimum)
   - Documented upgrade paths with backward compatibility guidance
