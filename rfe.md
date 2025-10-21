# Perses Observability Dashboard Integration for OpenShift AI

**Feature Overview:**

OpenShift AI users currently lack unified visibility into their AI workload performance, model serving health, and infrastructure resource consumption - forcing them to context-switch between multiple tools or build custom monitoring solutions. By embedding Perses observability dashboards directly into the OpenShift AI Dashboard, we enable data scientists and MLOps teams to monitor their AI workloads alongside their model development workflows, reducing mean-time-to-detection for performance issues from hours to minutes and accelerating troubleshooting by 3-5x.

This integration adds Perses as a new observability page within the existing ODH Dashboard UI, providing users with advanced dashboard capabilities for monitoring AI workloads, model serving metrics, pipeline performance, and infrastructure health - all within a unified, CNCF-backed open source framework that consolidates metrics (Prometheus), traces (Tempo), logs (Loki), and profiling data into a single interface.

**Goals:**

**G1: Reduce Time-to-Insight for AI Workload Issues**
- **User Context:** Data scientists and MLOps engineers currently export metrics manually or switch to Grafana/Prometheus UIs, losing workflow context and spending 30-60 minutes on average to diagnose production issues
- **Expected Outcome:** Users can diagnose model serving latency, pipeline failures, or resource bottlenecks within 2 clicks from their current workflow, without leaving the OpenShift AI Dashboard
- **Success Metric:** 70% reduction in average time from "model behaving unexpectedly" to "root cause identified"
- **Who Benefits:** Data scientists gain faster feedback loops during experimentation; MLOps engineers resolve production incidents faster; administrators identify platform issues proactively

**G2: Enable Self-Service Observability for Data Science Teams**
- **User Context:** Today, data science teams depend on platform admins or SREs to access infrastructure metrics, creating bottlenecks (average 2-4 day turnaround for metric access requests) and limiting their ability to optimize workloads independently
- **Expected Outcome:** Data scientists view pre-configured dashboards for their workloads without requiring cluster admin access or Prometheus query expertise; dashboards automatically scope to their authorized projects
- **Success Metric:** 50%+ of observability questions answered without admin escalation (measured via support ticket reduction)
- **Who Benefits:** Data scientists become self-sufficient in performance troubleshooting; platform admins reduce operational burden; organizations accelerate AI development velocity

**G3: Establish Foundation for AI-Specific Observability**
- **User Context:** Generic Kubernetes/infrastructure dashboards don't surface AI-relevant metrics (model drift indicators, inference throughput patterns, GPU utilization correlation with training progress), requiring users to manually correlate disparate data sources
- **Expected Outcome:** Purpose-built dashboards that correlate AI workload metrics with infrastructure health, making observability actionable for data science personas who lack deep infrastructure expertise
- **Success Metric:** 80% of users rate observability features as "very useful" in quarterly NPS surveys
- **Who Benefits:** All personas gain AI workload-specific insights; competitive positioning improves against platforms like Databricks, SageMaker, and Vertex AI which provide integrated monitoring

**Current State vs. Future State:**

| Dimension | Current State | With This Feature |
|-----------|--------------|-------------------|
| **Access Pattern** | Navigate to separate Grafana instance, lose context, requires manual URL bookmarking | Embedded dashboards within ODH workflow, contextual navigation from notebooks/models/pipelines |
| **User Capability** | Admin-dependent for metrics access, requires Prometheus/PromQL knowledge | Self-service for data scientists, pre-built dashboards with intuitive interfaces |
| **Dashboard Coverage** | Generic K8s metrics, no AI workload focus, custom setup per customer | AI-specific views (model serving, training pipelines, notebooks, GPU utilization) out-of-the-box |
| **Multi-tenancy** | Single cluster view, no namespace isolation, security concerns | Scoped dashboards based on user project access, RBAC-enforced metric visibility |
| **Extensibility** | Custom Grafana setup per customer, no standardization, difficult to maintain | Standard dashboard-as-code with GitOps workflows, customer customization supported via CUE/YAML |

**Out of Scope:**

The following items are explicitly NOT included in this feature to maintain focus on core observability capabilities and manage complexity:

1. **Custom Alerting Rules Configuration:** Users cannot create/modify Prometheus alerting rules through this UI. Alerting configuration remains platform admin responsibility via the Platform Metrics and Alerts architecture. *Rationale: Alert sprawl risk, requires separate RBAC model, Platform Metrics RFE owns alerting capabilities.*

2. **Direct PromQL Query Interface:** No ad-hoc PromQL query builder for end users in MVP. Pre-built dashboards and guided query building only. *Rationale: Steep learning curve for data scientists, security risk with unrestricted query access to multi-tenant metrics.*

3. **Historical Metric Data Export:** Bulk metric data export or long-term metric storage beyond platform defaults (90-day retention). *Rationale: Storage cost implications, should be platform-level capability not feature-specific; customers requiring longer retention configure at Platform Metrics layer.*

4. **Third-Party Data Source Integration:** Initial scope limited to centralized Prometheus/Tempo/Loki infrastructure. No direct integration with external APM tools (Datadog, New Relic, Dynatrace, Splunk). *Rationale: Perses capability exists but explodes integration testing matrix and introduces vendor dependencies. Defer to post-MVP based on customer demand signals.*

5. **Dashboard Sharing Across Organizations:** Dashboards scoped to single OpenShift AI instance, no cross-cluster sharing or public dashboard marketplace. *Rationale: Security/isolation model complexity, multi-cluster identity federation not solved. Future enhancement if customer demand materializes.*

6. **Custom Visualization Plugin Development:** Users consume existing Perses visualization types (time series, gauges, tables, heatmaps, stat panels), cannot develop custom plugin types through ODH UI. *Rationale: Requires development environment, plugin version management complexity, security review overhead. Advanced users can contribute upstream to Perses project.*

7. **Real-Time Collaborative Dashboard Editing:** No Google Docs-style multi-user editing of dashboards. Dashboard-as-code changes follow GitOps workflow with standard version control. *Rationale: Operational conflict resolution complexity, not a customer pain point based on 15+ customer interviews.*

8. **Log Search and Analysis Interface:** While Perses supports Loki integration, full log search/analysis UI is out of scope. Focus remains on metric visualization. *Rationale: Log analysis requires different UX patterns (search syntax, large result sets, context viewers); adds 4-6 months to delivery timeline. Can be phased in post-GA.*

**Requirements:**

### MVP (Must-Have for Initial Release)

**R1: Embedded Dashboard Viewing (P0 - Critical)**
- **Capability:** Render Perses dashboards within ODH Dashboard UI using React component integration via `@perses-dev/components` and `@perses-dev/dashboards` npm packages
- **Acceptance Criteria:**
  - Users access observability page without leaving ODH context
  - Dashboards load within 3 seconds on typical network conditions (1Mbps+)
  - Visual consistency with ODH Dashboard PatternFly theme (colors, typography, spacing, interactive states)
  - Smooth navigation transitions between ODH pages and Perses content
- **Priority:** P0 - Core functionality; without this, feature doesn't exist
- **Customer Impact:** Foundation for all observability capabilities

**R2: Multi-Tenant Dashboard Scoping (P0 - Critical)**
- **Capability:** Dashboards automatically filter metrics to projects/namespaces user has RBAC access to, enforced at Prometheus query layer
- **Acceptance Criteria:**
  - User in project A cannot view metrics from project B even with URL manipulation
  - Scoping aligns with ODH Dashboard existing RBAC model (same project access = same metric visibility)
  - Adversarial testing validates no metric leakage across namespace boundaries
  - Performance impact of RBAC enforcement < 200ms per query
- **Priority:** P0 - Enterprise requirement; blocker for multi-tenant environments
- **Customer Impact:** Critical for 80% of enterprise deployments; financial services, healthcare, and government customers have regulatory requirements for data isolation

**R3: Pre-Built AI Workload Dashboards (P0 - Critical)**
- **Capability:** Ship 4-6 curated dashboards covering key AI workload scenarios:
  1. **Model Serving Overview:** Inference requests/sec, latency (p50/p95/p99), error rates by model/version, request payload sizes
  2. **Model Training Metrics:** Training loss curves, throughput (samples/sec), GPU utilization, memory consumption, epoch duration
  3. **Pipeline Execution Health:** Run duration by pipeline, success rate, step-level performance, resource utilization trends
  4. **Notebook Environment Resources:** Per-user CPU/memory/GPU usage, active sessions, kernel restart frequency, storage consumption
  5. **Cluster Infrastructure Health:** Node capacity (CPU/memory/GPU), pod health by namespace, persistent volume usage, network throughput
  6. **GPU Performance Analysis:** GPU utilization correlation with training progress, memory allocation patterns, multi-GPU scaling efficiency
- **Acceptance Criteria:**
  - Zero configuration required post-install, dashboards render data immediately if Platform Metrics deployed
  - Each dashboard includes contextual help (tooltips explaining metrics, typical healthy ranges)
  - Dashboards support common filtering variables (namespace, model name, pipeline ID, user)
  - Mobile-responsive layouts (readable on 1024px+ screens)
- **Priority:** P0 - Without pre-built content, adoption will be <10% based on Grafana usage patterns
- **Customer Impact:** 90% of customers in user research want "batteries included" observability; custom dashboard creation is advanced use case

**R4: Dashboard-as-Code Support (P0 - Critical)**
- **Capability:** Dashboards defined in Perses CUE format or YAML, stored in Kubernetes ConfigMaps, support GitOps workflows via Kustomize/ArgoCD
- **Acceptance Criteria:**
  - Platform admins can version control dashboard definitions in Git repositories
  - Apply dashboards via `kubectl apply -f dashboard.yaml` or ArgoCD sync
  - Dashboard updates propagate to ODH UI within 30 seconds
  - Rollback capability via standard GitOps practices (revert commit → sync)
- **Priority:** P0 - Enterprise operational model requirement
- **Customer Impact:** Customers require auditable, repeatable configuration for SOC2, ISO 27001, FedRAMP compliance

**R5: Integration with Platform Metrics Architecture (P0 - Critical)**
- **Capability:** Consume metrics from centralized Prometheus federation, no separate metric collection infrastructure
- **Acceptance Criteria:**
  - Zero additional metric storage overhead beyond Platform Metrics baseline
  - Respects existing retention policies (default 90 days)
  - Query performance meets SLA: simple queries <500ms, complex aggregations <2s at 1000 concurrent users
  - Graceful degradation if Prometheus temporarily unavailable (cached data + clear error messaging)
- **Priority:** P0 - Architectural dependency; feature is non-functional without metric backend
- **Customer Impact:** Customers won't accept separate monitoring infrastructure (doubles operational cost, storage requirements, complexity)

**R6: Basic Dashboard Management UI (P1 - High Priority)**
- **Capability:** Admin UI within ODH Dashboard to list, enable/disable, preview, and validate dashboards
- **Acceptance Criteria:**
  - CRUD operations on dashboard configs without requiring kubectl CLI access
  - Preview dashboard before making it available to users
  - Validate dashboard syntax/queries before deployment
  - Role-based access: cluster admins manage cluster-wide dashboards, project admins manage project-scoped dashboards
- **Priority:** P1 - Usability for non-Kubernetes-expert admins
- **Customer Impact:** 40% of ODH admins are application-focused, not K8s experts; CLI-only management creates adoption barriers

**R7: Single Sign-On Integration (P0 - Critical)**
- **Capability:** Perses components authenticate using ODH Dashboard session, no separate login required
- **Acceptance Criteria:**
  - Seamless authentication flow (users never see Perses login screen)
  - Session expiration aligned with ODH Dashboard policies
  - Token refresh handled transparently
- **Priority:** P0 - Enterprise SSO requirement
- **Customer Impact:** Universal enterprise customer requirement; separate authentication creates security audit concerns

**R8: Contextual Navigation Integration (P1 - High Priority)**
- **Capability:** Navigate to relevant Perses dashboards directly from ODH workflow contexts
- **Acceptance Criteria:**
  - Model serving page includes "View Metrics" button → Model Serving dashboard pre-filtered to that model
  - Training job page includes "View Training Metrics" → Training dashboard pre-filtered to that job
  - Data science project overview includes embedded summary metrics from Perses
  - Pipeline run details include "View Pipeline Metrics" → Pipeline dashboard pre-filtered to that run
- **Priority:** P1 - Critical for workflow integration
- **Customer Impact:** Without contextual access, users must manually configure filters, reducing adoption by 50%+

### Post-MVP (Nice-to-Have - Phase 2+)

**R9: Custom Dashboard Creation via UI (P2 - Medium Priority)**
- **Capability:** Drag-and-drop dashboard builder for admins to create project-specific views without writing CUE/YAML
- **Priority:** P2 - Advanced feature, 20-30% of customers need this based on competitive analysis
- **Defer Rationale:** Dashboard-as-code covers MVP; adds 4-6 weeks development time; can validate demand in Tech Preview

**R10: Dashboard Template Library (P2 - Medium Priority)**
- **Capability:** Catalog of community/partner-contributed dashboard templates (e.g., TensorFlow Serving specific, PyTorch profiling, Ray cluster monitoring)
- **Priority:** P2 - Ecosystem play, requires community maturity
- **Defer Rationale:** Not blocking adoption, can build once base adoption hits 40%+ and community contribution model established

**R11: Correlation with Distributed Tracing (P2 - Medium Priority)**
- **Capability:** Click-through from metric spike to related trace spans in Tempo (e.g., click latency spike → view distributed traces for slow requests)
- **Priority:** P2 - Advanced troubleshooting, high value for 15-20% of power users (MLOps engineers)
- **Defer Rationale:** Requires Tempo integration maturity in Platform Metrics; Perses tracing features still evolving upstream

**R12: Anomaly Detection Overlays (P3 - Low Priority)**
- **Capability:** ML-powered anomaly highlighting on metric charts (e.g., training loss anomaly detection, serving latency outlier identification)
- **Priority:** P3 - "AI observing AI" narrative is compelling but technology immature
- **Defer Rationale:** No clear vendor solution, would require custom ML model development and validation; market not demanding this yet

**R13: Cost Visibility Dashboards (P2 - Medium Priority)**
- **Capability:** Per-project cloud cost attribution based on resource consumption (CPU-hours, GPU-hours, storage)
- **Priority:** P2 - FinOps trend is strong (35% of enterprises prioritize this), but requires billing integration complexity
- **Defer Rationale:** Needs product-wide cost tracking capability first; OpenShift cost management integration required

**R14: Mobile-Optimized Dashboard Views (P3 - Low Priority)**
- **Capability:** Responsive layouts optimized for mobile on-call monitoring via smartphones
- **Priority:** P3 - Limited user demand (8% requested in surveys); 95% of usage is desktop-based
- **Defer Rationale:** Desktop workflow is primary; tablet support sufficient for MVP

### Enterprise & Scalability Requirements

**R15: Audit Logging (P0/P1 - Context Dependent)**
- **Capability:** Dashboard access and modification events logged to ODH audit trail (who accessed which dashboard, when, which projects' data viewed)
- **Priority:** P0 for regulated industries (healthcare/HIPAA, finance/SOX, government/FedRAMP = 35% of revenue); P1 for general availability
- **Timeline:** MVP if ODH audit framework exists, otherwise P1 for MVP+1
- **Customer Impact:** Required for SOC2, ISO 27001, HIPAA, FedRAMP compliance audits

**R16: High Availability (P0 - Critical)**
- **Capability:** Dashboard availability SLA matches ODH Dashboard (99.5% uptime target)
- **Acceptance Criteria:**
  - Perses backend components deployed with 2+ replicas
  - Support for pod disruption budgets
  - Graceful handling of backend failures (show cached data + warning)
  - Zero downtime upgrades
- **Priority:** P0 - Production workload dependency
- **Customer Impact:** Production AI workloads depend on observability for incident response; downtime unacceptable

**R17: Scalability to 1000+ Concurrent Users (P0 - Critical)**
- **Capability:** Support large enterprise deployments without performance degradation
- **Acceptance Criteria:**
  - Load testing validates 1000 concurrent users with <3s dashboard load time
  - Prometheus query performance optimized (query result caching, efficient PromQL generation)
  - Horizontal scaling of Perses backend components
  - Resource consumption documented (CPU, memory, storage per 100 users)
- **Priority:** P0 - Top 10 customers have 800-2000 data scientists
- **Customer Impact:** Performance degradation at scale creates negative perception of entire platform

**R18: Air-Gapped Installation Support (P1 - High Priority)**
- **Capability:** All Perses dependencies (npm packages, container images, dashboard definitions) available via disconnected install
- **Acceptance Criteria:**
  - All container images mirrorable to internal registries
  - No runtime dependencies on external CDNs or package registries
  - Documented disconnected installation procedure
  - Tested in true air-gapped environment (no internet access)
- **Priority:** P1 - 25% of enterprise customers require this (government, defense, financial services)
- **Timeline:** Must be resolved before GA, can defer for initial Tech Preview

**R19: Backup & Disaster Recovery (P1 - High Priority)**
- **Capability:** Dashboard configurations included in ODH backup/restore procedures
- **Acceptance Criteria:**
  - Dashboard definitions (ConfigMaps) backed up via Velero or equivalent
  - Restore procedure tested with <4 hour RTO
  - Documentation covers backup best practices
- **Priority:** P1 - Enterprise operational requirement
- **Customer Impact:** Custom dashboards represent significant configuration investment; loss creates major operational disruption

**Done - Acceptance Criteria:**

The Perses observability integration will be considered complete and ready for General Availability when the following criteria are met from a user's perspective:

### User Discovery & Navigation
- Users can discover the observability feature within 30 seconds of logging into ODH Dashboard (visible in left navigation as "Observability" or "Metrics & Dashboards")
- Clicking "Observability" navigation loads dashboard list page within 2-3 seconds
- Users can navigate to relevant dashboards directly from context (e.g., "View Metrics" button on model serving page opens appropriate dashboard pre-filtered)
- Browser back/forward navigation works predictably; breadcrumbs accurately reflect location

### Visual & Interactive Consistency
- Perses dashboards visually indistinguishable from native ODH pages (same colors, fonts, spacing, button styles)
- Loading states match ODH patterns (skeleton screens or spinners consistent with rest of application)
- Error messages follow ODH tone and formatting conventions
- Dashboard interactions (hover tooltips, zoom, pan, time range selection) feel responsive (<200ms feedback)

### Multi-Tenancy & Security
- Data scientists see only metrics for projects they have access to (automatic filtering by namespace RBAC)
- Attempting to access unauthorized metrics returns clear error message (not exposed via URL manipulation)
- Administrators can view cluster-wide dashboards; project members see project-scoped views
- No separate authentication required (SSO from ODH Dashboard session)

### Pre-Built Dashboard Functionality
- At least 4 pre-built dashboards available immediately after installation with zero configuration:
  1. Model Serving Performance (latency, throughput, errors)
  2. Model Training Metrics (loss, GPU utilization, throughput)
  3. Pipeline Execution Health (run duration, success rate, resource usage)
  4. Cluster Infrastructure (node health, capacity, pod status)
- Each dashboard displays data within 3 seconds if Platform Metrics is operational
- Dashboards include contextual help (tooltips explaining metrics and healthy ranges)
- Users can adjust time range (last 15m, 1h, 6h, 24h, 7d, custom) and see updated data within 2 seconds

### Dashboard Interactivity
- Users can click panels to expand to full-screen view for detailed analysis
- Users can apply filters (namespace, model name, job ID) that affect all panels consistently
- Zoom and pan operations on time series charts work smoothly
- Dashboard state (time range, filters, zoom level) persists when navigating away and returning

### Dashboard Management (Administrators)
- Platform admins can create new dashboards by applying YAML/CUE files via kubectl or GitOps
- Dashboard changes propagate to UI within 30 seconds
- Admin UI provides list of all dashboards with enable/disable capability
- Dashboard syntax errors caught before deployment with clear validation messages

### Performance & Reliability
- System supports 100 concurrent users viewing dashboards with <3 second load time (target: 1000+ users for enterprise scale)
- Dashboard remains functional if Prometheus temporarily unavailable (shows cached data + informative error)
- No memory leaks during extended use (8+ hour sessions)
- Perses backend components deployed in HA configuration (2+ replicas) with 99.5% uptime

### Accessibility (WCAG 2.1 AA)
- All dashboard interactions accessible via keyboard (tab navigation, focus indicators visible)
- Screen readers can interpret charts (text alternatives provided)
- Color contrast meets 4.5:1 ratio minimum
- Information not conveyed by color alone (patterns, labels, icons supplement)
- No keyboard traps in embedded Perses components

### Documentation Completeness
- Getting started guide shows data scientists how to access training/serving metrics in <5 minutes
- Troubleshooting guide covers "no data displayed" and "dashboard not loading" scenarios
- Administrator guide covers dashboard-as-code workflow with GitOps example
- At least one video walkthrough demonstrating key workflows (3-5 minutes)

### Integration Quality
- Platform Metrics architecture deployed and operational (dependency)
- Metric queries respect Platform Metrics retention policies (default 90 days)
- No separate metric storage infrastructure required
- Dashboard queries optimized for performance (query result caching enabled)

### Enterprise Readiness
- Audit logging operational (who accessed which dashboard, timestamp) for regulated industries
- Air-gapped installation documented and tested
- Backup/restore procedure documented with <4 hour RTO validated
- Upgrade path tested (dashboards survive ODH version upgrades)

### Success Metrics (Post-Launch)
- 70% of OpenShift AI users access observability dashboards within first 30 days
- Average time-to-diagnosis for performance issues reduced by 50% (measured via before/after comparison)
- Support tickets related to "how do I monitor X" reduced by 40%
- 80% of users rate observability as "useful" or "very useful" in quarterly survey

**Use Cases - i.e. User Experience & Workflow:**

### Use Case 1: Data Scientist - Monitoring Model Training Performance

**Primary Actor:** Sofia, Data Scientist
**Goal:** Understand why her LLM training job is progressing slower than expected
**Preconditions:** Sofia has a training job running via Jupyter notebook in OpenShift AI

**Main Success Scenario:**

1. **Entry Point:** Sofia is training a large language model using a notebook in ODH. After 2 hours, she notices training throughput is 30% lower than her previous experiment.
2. **Discovery:** She navigates to the ODH Dashboard home page and sees "Observability" in the left navigation panel (new menu item below "Data Science Projects").
3. **Dashboard Selection:** She clicks "Observability" and lands on a dashboard list page showing categorized dashboards:
   - **Training & Experimentation:** Model Training Metrics, GPU Performance Analysis, Notebook Resources
   - **Model Serving:** Serving Performance, Inference Latency
   - **Platform Health:** Cluster Infrastructure, Pipeline Execution
4. **Navigation:** She clicks "Model Training Metrics" and sees a comprehensive dashboard with panels:
   - Training loss curve (updating in real-time)
   - Throughput (samples/second) over time
   - GPU utilization (%) per GPU
   - Memory consumption (allocated vs. available)
   - Epoch duration trends
5. **Filtering:** At the top of the dashboard, she sees a namespace filter pre-set to her project ("sofia-llm-research") and a "Training Job" dropdown. She selects her current job from the dropdown.
6. **Analysis:** She uses the time range picker to focus on "Last 2 hours" (matching when she started the job). She immediately notices:
   - GPU utilization oscillates between 35-45%, not the expected 85-95%
   - Memory shows periodic spikes corresponding to GPU utilization drops
7. **Insight Discovery:** She clicks on the GPU utilization panel, which expands to full-screen view. A tooltip explains: "Low GPU utilization often indicates data loading bottlenecks. Check dataloader settings or increase prefetch workers."
8. **Cross-Reference:** She clicks back and scrolls to a "Data Loading Throughput" panel showing dataloader queue depth is frequently zero, confirming the bottleneck.
9. **Action:** Armed with this insight, she bookmarks the dashboard URL (which includes her filter selections) and returns to her notebook to increase dataloader worker count from 2 to 8.
10. **Validation:** After restarting training, she refreshes the dashboard and watches GPU utilization climb to 90%, validating her fix worked.

**Alternative Flow 1 - Contextual Access:**

1. Sofia is on the "Data Science Projects" page viewing her project details.
2. Under "Active Resources," she sees her running training job with a "View Metrics" icon/link next to it.
3. Clicking "View Metrics" navigates directly to the Model Training Metrics dashboard, already filtered to her specific job.
4. She sees job-specific metrics immediately without needing to configure filters.

**Alternative Flow 2 - Comparing Experiments:**

1. Sofia wants to compare her current training run with a previous successful run.
2. On the Model Training Metrics dashboard, she opens the "Training Job" dropdown which supports multi-select.
3. She selects both "current-run-20250115" and "baseline-run-20250110".
4. The dashboard panels update to show overlay charts (e.g., two loss curves on same axis with different colors, side-by-side GPU utilization).
5. She identifies that baseline run had consistent 90% GPU utilization, while current run oscillates, confirming a regression.

**Exception Flow - No Data Available:**

1. Sofia navigates to Model Training Metrics dashboard.
2. Instead of charts, she sees an informative message: "No training metrics available for your project in the selected time range. Ensure your training jobs are emitting Prometheus metrics. [View Setup Guide]"
3. Clicking "View Setup Guide" opens documentation explaining how to instrument training code to emit metrics.

---

### Use Case 2: MLOps Engineer - Troubleshooting Production Model Serving Latency

**Primary Actor:** Marcus, MLOps Engineer
**Goal:** Diagnose and resolve sudden latency spike in production model serving endpoint
**Preconditions:** Marcus has deployed a BERT model for production inference via KServe; monitoring alerts him to latency issues

**Main Success Scenario:**

1. **Alert Trigger:** Marcus receives a Slack alert from his monitoring system: "BERT model p95 latency exceeded 500ms (current: 1200ms)."
2. **Quick Access:** He opens OpenShift AI Dashboard on his laptop and navigates to "Model Serving" page.
3. **Contextual Navigation:** He finds the affected model deployment ("bert-sentiment-v2-prod") and clicks the "Observability" tab (new tab alongside "Metrics," "Logs," "YAML").
4. **Dashboard Load:** The "Model Serving Performance" dashboard loads, automatically filtered to "bert-sentiment-v2-prod" deployment.
5. **Immediate Visualization:** He sees a dashboard with key panels:
   - **Request Rate:** Requests/second (shows spike from 100 to 300 rps starting 15 minutes ago)
   - **Latency Percentiles:** p50, p95, p99 over time (p95 jumped from 200ms to 1200ms)
   - **Error Rate:** 4xx and 5xx errors (shows 2% 5xx errors starting with latency spike)
   - **Resource Utilization:** CPU and memory per pod
   - **Autoscaling Status:** Current replicas vs. desired replicas
6. **Time Correlation:** He notices request rate spike correlates exactly with latency increase. He adjusts time range to "Last 6 hours" to see broader context.
7. **Pattern Recognition:** He identifies this is not the first spike - similar pattern occurred 24 hours ago at the same time (daily batch job submission spike).
8. **Resource Investigation:** He clicks the "Resource Utilization" panel to expand it. He sees CPU is at 95% across all pods, and autoscaling is in progress (current: 3 pods, desired: 8 pods).
9. **Autoscaling Lag:** He switches to "Autoscaling Status" panel and sees new pods are starting but not yet ready (initializing for 3-4 minutes due to large model loading).
10. **Root Cause Identified:** The issue is insufficient autoscaling headroom for sudden traffic spikes combined with slow pod startup time.
11. **Sharing with Team:** He copies the dashboard URL (which includes current time range and filters) and pastes it into Slack: "Found the issue - autoscaling can't keep up with daily batch spike. Here's the data: [URL]. Proposing we increase min replicas to 5 and add pod affinity to keep warm standbys."
12. **Follow-Up:** He creates a custom dashboard (using dashboard-as-code) specifically for this model that includes additional panels for cold start tracking and replica scaling events.

**Alternative Flow - Multi-Model Comparison:**

1. Marcus needs to evaluate whether to roll out a new model version (v3) to replace the current production version (v2).
2. He navigates to "Observability" → "Model Serving Performance" dashboard.
3. He uses the "Model Deployment" multi-select filter to choose both "bert-sentiment-v2-prod" and "bert-sentiment-v3-canary".
4. The dashboard updates to show side-by-side comparison panels:
   - Latency comparison (v3 shows 15% lower p95 latency)
   - Throughput comparison (v3 handles same request rate)
   - Error rate comparison (v3 has 0.1% errors vs. v2's 0.5%)
   - Resource usage comparison (v3 uses 20% more memory but same CPU)
5. Based on this data, he decides to proceed with v3 rollout and documents the evidence in his change request.

**Exception Flow - Perses Backend Unavailable:**

1. Marcus navigates to the Model Serving Performance dashboard.
2. The page displays a warning banner: "Observability service temporarily unavailable. Showing cached data from 5 minutes ago. [Retry] [Troubleshoot]"
3. Panels display with a visual indicator (greyed out or watermark) showing they're cached.
4. He clicks "Retry" and the system attempts to reconnect. If successful, fresh data loads; if not, the error message updates with next retry time.

---

### Use Case 3: OpenShift AI Administrator - Platform Health Monitoring & Capacity Planning

**Primary Actor:** Priya, OpenShift AI Platform Administrator
**Goal:** Perform daily platform health check and identify capacity planning needs for upcoming quarter
**Preconditions:** Priya is responsible for maintaining OpenShift AI cluster health for 200+ data scientists

**Main Success Scenario (Daily Health Check):**

1. **Morning Routine:** Priya starts her day by logging into OpenShift AI Dashboard.
2. **Platform Overview:** She navigates to "Observability" → "Platform Overview" dashboard.
3. **Health At-a-Glance:** The dashboard shows high-level health indicators:
   - **Active Workloads:** 87 running (12 training jobs, 15 serving deployments, 45 active notebooks, 15 pipelines)
   - **Cluster Resource Utilization:** CPU 67%, Memory 72%, GPU 84% (color-coded: green <70%, yellow 70-85%, red >85%)
   - **Project Count:** 34 active data science projects
   - **Anomalies:** 2 items flagged (highlighted in yellow)
4. **Anomaly Investigation:** She clicks on the "Anomalies" indicator and sees:
   - "Project 'data-science-team-5' showing unusual network egress (5TB in last 24h vs. typical 50GB)"
   - "GPU utilization in 'ml-research' namespace dropped from 85% to 15% at 03:00 AM"
5. **Deep Dive - Network Issue:** She clicks the first anomaly, which navigates to "Namespace Resource Usage" dashboard pre-filtered to "data-science-team-5" project.
6. **Detailed Analysis:** She sees network egress chart showing sustained 500 Mbps outbound traffic starting at 2 AM. She cross-references with "Active Workloads" panel and identifies a specific notebook session.
7. **Action:** She clicks through to the Data Science Projects page, contacts the project owner, and discovers they're inadvertently syncing large datasets to external S3 repeatedly in a loop. She helps them fix the script.
8. **Validation:** She returns to the Platform Overview dashboard and adds the project to her watchlist for monitoring over the next 24 hours.

**Main Success Scenario (Capacity Planning):**

1. **Quarterly Planning:** Priya needs to forecast resource needs for Q2 2025.
2. **Historical Trends:** She navigates to "Observability" → "Historical Trends" dashboard.
3. **Time Range Adjustment:** She changes the time range to "Last 90 days" to see quarterly trends.
4. **Key Metrics Analysis:** She reviews trend panels:
   - **GPU Utilization Trend:** Shows steady climb from 60% (Dec) to 84% (Feb), linear projection suggests 95%+ by April
   - **Storage Growth:** Persistent volumes growing 15% per month
   - **User Activity:** Active notebook sessions increased 25% quarter-over-quarter
   - **Model Serving Load:** Inference request rate doubled since December
5. **Scenario Modeling:** She notes current GPU capacity is 24 A100 GPUs. At current growth rate, she'll need 6-8 additional GPUs by end of Q2.
6. **Cost Estimation:** She switches to "Resource Cost Attribution" panel (if available post-MVP) or manually calculates based on resource usage per project.
7. **Reporting:** She exports key charts (screenshot or download CSV if available) for her capacity planning presentation to leadership.
8. **Custom Dashboard Creation:** She decides to create a custom "Executive Summary" dashboard for monthly reports. She navigates to "Observability" → "Manage Dashboards" → "Create New Dashboard."
9. **Dashboard-as-Code:** Since she's comfortable with YAML, she clicks "Edit YAML" and uses a provided template to define a new dashboard with 6 high-level panels (total workloads, resource utilization trends, cost per project, top resource consumers).
10. **GitOps Workflow:** She saves the dashboard definition to her GitOps repository (`dashboards/executive-summary.yaml`) and applies it via ArgoCD. Within 30 seconds, the new dashboard appears in the ODH Dashboard observability section.

**Alternative Flow - Multi-Cluster View (Post-MVP):**

1. Priya manages OpenShift AI deployments across 3 clusters (dev, staging, prod).
2. She navigates to a "Multi-Cluster Overview" dashboard (post-MVP feature).
3. She sees aggregated metrics across all clusters with cluster selector filter.
4. She identifies that dev cluster has low utilization (30%) while prod is at 85%, informing workload rebalancing decisions.

**Exception Flow - Missing Platform Metrics:**

1. Priya navigates to Platform Overview dashboard immediately after ODH upgrade.
2. She sees a message: "Platform Metrics service not detected. Observability features require Platform Metrics & Alerts architecture. [Installation Guide] [Check Status]"
3. Clicking "Check Status" runs diagnostic checks and reports: "Prometheus operator: Not installed. Suggested action: Deploy Platform Metrics via [this procedure]."

---

### Use Case Diagram Summary

**Actors:**
- Data Scientist (Sofia)
- MLOps Engineer (Marcus)
- Platform Administrator (Priya)

**Primary Use Cases:**
1. Monitor Training Performance (Data Scientist)
2. Troubleshoot Serving Latency (MLOps Engineer)
3. Platform Health Check (Administrator)
4. Capacity Planning (Administrator)
5. Compare Model Versions (MLOps Engineer)
6. Custom Dashboard Creation (Administrator, MLOps Engineer)

**System Interactions:**
- OpenShift AI Dashboard (UI layer)
- Perses Embedded Components (visualization layer)
- Platform Metrics - Prometheus (data layer)
- Kubernetes RBAC (authorization layer)
- GitOps Repository (configuration management)

**Integration Points:**
- Model Serving page → Model Serving Performance dashboard (contextual link)
- Training Job page → Model Training Metrics dashboard (contextual link)
- Data Science Project page → Project Resource Usage dashboard (contextual link)
- Pipeline Run page → Pipeline Execution Health dashboard (contextual link)

---

### Workflow Success Criteria

**Navigation Efficiency:**
- Users reach relevant dashboard in 2 clicks or less from primary workflow context
- Contextual links (e.g., "View Metrics" on model serving page) reduce configuration burden

**Information Discovery:**
- Users identify performance issues within 2 minutes of viewing dashboard (vs. 30+ minutes with separate tools)
- Tooltips and contextual help enable self-service troubleshooting without admin escalation

**Workflow Continuity:**
- Navigating to observability and back to previous context feels seamless (no jarring visual/interaction differences)
- Dashboard state (filters, time range) persists appropriately during session

**Collaboration:**
- Users can share dashboard URLs with teammates that preserve context (time range, filters, zoom level)
- Custom dashboards enable team-specific workflows without fragmenting user experience

**Documentation Considerations:**

Comprehensive documentation is critical to drive adoption and enable self-service observability. Documentation must support the entire user journey from discovery to mastery, while maintaining consistency with the broader OpenShift AI documentation ecosystem.

### Documentation Structure & Deliverables

**Tier 1: Embedded Contextual Help (In-Product)**
- **Dashboard Panel Tooltips:** Every metric chart includes hover tooltip explaining:
  - What the metric measures (in plain language, not just technical definition)
  - Typical healthy ranges (e.g., "GPU utilization: 80-95% indicates efficient training; <50% suggests bottlenecks")
  - Common causes of anomalies (e.g., "Latency spikes often correlate with autoscaling lag or resource contention")
- **Empty State Guidance:** When dashboards show no data:
  - Clear explanation of why (e.g., "No training jobs running in selected time range" vs. "Platform Metrics not configured")
  - Actionable next steps (e.g., "Start a training job" or "Configure metric collection [link to setup guide]")
- **Error Message Context:** When queries fail or components error:
  - User-friendly error messages (not raw Prometheus error strings)
  - Troubleshooting links relevant to error type

**Tier 2: Task-Oriented Guides (Primary Documentation)**

**For Data Scientists:**

1. **"Monitoring Your Model Training"** (5-minute quickstart)
   - How to access training metrics from notebook or Data Science Projects page
   - Understanding key training metrics: loss curves, throughput, GPU utilization, memory consumption
   - Adjusting time ranges to focus on specific training runs or epochs
   - Identifying common bottlenecks:
     - Data loading (low GPU utilization with high memory)
     - Gradient computation (high GPU utilization but low throughput)
     - Distributed training issues (GPU utilization imbalance across nodes)
   - Screenshots showing each step with annotations

2. **"Model Serving Performance at a Glance"** (5-minute quickstart)
   - Accessing serving metrics for your deployed models (from Model Serving page)
   - Interpreting latency percentiles (p50, p95, p99): what's normal, what's concerning
   - Understanding request rates and error patterns (4xx vs. 5xx errors)
   - Correlating performance changes with deployment events (new version rollout, traffic spike)
   - When to scale your deployment based on metrics (CPU saturation, memory pressure, queue depth)

**For MLOps Engineers:**

3. **"Comprehensive Observability Workflows"** (15-minute deep dive)
   - Complete tour of available dashboards and their purposes:
     - When to use Training vs. Serving vs. Platform dashboards
     - How dashboards correlate with ML lifecycle stages
   - Correlating metrics across dashboards:
     - Training metrics → model quality signals
     - Serving metrics → infrastructure health
     - Cross-referencing namespace resources with workload performance
   - Using filters and variables effectively:
     - Multi-environment monitoring (dev, staging, prod)
     - Comparing model versions side-by-side
     - Project-level vs. cluster-wide views
   - Sharing dashboard states:
     - URL structure and parameters (how to craft shareable links)
     - Embedding dashboard links in incident reports, runbooks, status pages

4. **"Troubleshooting with Observability Dashboards"** (20-minute deep dive)
   - Troubleshooting decision tree:
     - Symptom: "Model serving latency increased" → Dashboard: "Model Serving Performance" → Check: Request rate spike? Resource saturation? Error rate increase? → Action: Scale deployment, optimize model, investigate failing requests
     - Symptom: "Training not progressing" → Dashboard: "Model Training Metrics" → Check: GPU utilization? Loss curve plateaued? OOM errors? → Action: Adjust batch size, increase dataloader workers, check for NaN gradients
     - Symptom: "Pipeline failures" → Dashboard: "Pipeline Execution Health" → Check: Which step failing? Resource constraints? Network issues? → Action: Increase step resources, check data availability, review logs
   - Common metric patterns and their meanings:
     - Sawtooth pattern in memory: likely data loading inefficiency
     - Request latency bimodal distribution: cold start issue or heterogeneous request complexity
     - GPU utilization oscillating: batch size too small or dataloader bottleneck
   - Integrating observability with incident response workflows (screenshots of real troubleshooting scenarios)

**For Platform Administrators:**

5. **"Platform Health Monitoring"** (10-minute guide)
   - Daily/weekly platform health check workflows:
     - Morning health check routine (5 minutes): Platform Overview dashboard review
     - Weekly deep dive (30 minutes): Historical trends, capacity forecasting, anomaly investigation
   - Capacity planning using historical trends:
     - Interpreting 90-day resource utilization trends
     - Forecasting future needs (linear projection techniques)
     - Identifying seasonal patterns (academic semester cycles, fiscal year-end data crunches)
   - Multi-tenant resource usage patterns:
     - Fair share enforcement monitoring
     - Identifying "noisy neighbor" scenarios
     - Project-level resource attribution for chargeback

6. **"Creating Custom Dashboards"** (30-minute guide)
   - When to create a custom dashboard vs. using provided ones:
     - Use case: organization-specific metrics (e.g., compliance-related tracking)
     - Use case: executive reporting (high-level KPIs)
     - Use case: specialized ML frameworks (Ray, Horovod) requiring custom panels
   - Dashboard-as-code approach:
     - Understanding Perses dashboard YAML/CUE structure
     - Anatomy of a dashboard definition (variables, panels, queries, layout)
     - Using provided templates as starting points
   - GitOps workflow for dashboards:
     - Storing dashboards in Git repository
     - Applying dashboards via kubectl or ArgoCD
     - Version control and rollback procedures
     - Peer review process for shared dashboards
   - Dashboard validation and testing:
     - Syntax validation before deployment
     - Preview mode for testing before making available to users
     - Performance testing for complex queries

7. **"Advanced Dashboard Customization"** (45-minute advanced guide)
   - Writing PromQL queries for custom metrics (with curated ML-specific examples):
     - Aggregating metrics across multiple pods/models
     - Rate calculations (requests per second from counters)
     - Histogram quantile queries (latency percentiles)
     - Joining metrics from different sources
   - Using variables for dynamic filtering:
     - Namespace, model name, job ID variables
     - Multi-select variables for comparison scenarios
     - Cascading variables (select project → available models update)
   - Creating composite panels:
     - Overlay charts (multiple metrics on same axis)
     - Side-by-side comparisons
     - Stat panels with thresholds and color coding
   - Dashboard organization best practices:
     - Information hierarchy (most critical metrics top-left)
     - Logical flow (overview → details)
     - Effective use of rows and panel sizing

**Tier 3: Reference Documentation**

8. **"Available Metrics Reference"**
   - Comprehensive catalog of metrics available in ODH environment:
     - **Training Metrics:** `odh_training_loss`, `odh_training_throughput`, `odh_gpu_utilization`, `odh_training_epoch_duration`
     - **Serving Metrics:** `odh_inference_requests_total`, `odh_inference_duration_seconds`, `odh_inference_errors_total`
     - **Pipeline Metrics:** `odh_pipeline_run_duration_seconds`, `odh_pipeline_step_status`
     - **Resource Metrics:** `container_cpu_usage_seconds_total`, `container_memory_working_set_bytes`, `nvidia_gpu_duty_cycle`
   - For each metric:
     - Description (what it measures)
     - Type (counter, gauge, histogram, summary)
     - Labels (available dimensions for filtering)
     - Typical queries (examples of useful PromQL)
     - Retention period (90 days default)
   - Metric naming conventions and best practices
   - How to request new metrics (contribution process)

9. **"Dashboard Management API Reference"** (for automation)
   - Programmatic dashboard creation, update, deletion via Kubernetes API
   - Dashboard schema reference (complete YAML/CUE spec)
   - Importing/exporting dashboards programmatically
   - Dashboard permission management via RBAC
   - Example automation scripts (Python, Go) for common tasks

10. **"Troubleshooting Reference"**
   - Comprehensive error code catalog with remediation steps:
     - `ERR_AUTH_001`: "Authentication failed" → Check ODH Dashboard session, try re-login
     - `ERR_QUERY_002`: "Query timeout" → Reduce time range, simplify query, check Prometheus health
     - `ERR_DATA_003`: "No data for selected time range" → Verify metric collection configured, adjust time range
   - Common issues and solutions:
     - Issue: "Dashboard shows no data"
       - Check 1: Is Platform Metrics deployed? (`kubectl get prometheus -n openshift-ai`)
       - Check 2: Are workloads emitting metrics? (verify ServiceMonitor configurations)
       - Check 3: Is user's RBAC correct? (test query access via `kubectl auth can-i`)
       - Check 4: Is time range appropriate? (recent data only available if within retention period)
     - Issue: "Dashboard loads slowly"
       - Check 1: How many panels on dashboard? (>20 panels can cause performance issues)
       - Check 2: Query complexity? (simplify PromQL, add recording rules for expensive queries)
       - Check 3: Time range too broad? (90-day queries are expensive; recommend <7 days for interactive use)
     - Issue: "Metrics not scoped to my namespace"
       - Check 1: Verify dashboard uses namespace variable correctly
       - Check 2: Check RBAC permissions (`kubectl get rolebindings -n <namespace>`)
       - Check 3: Verify Prometheus query includes namespace label filter

### Documentation Integration & Discoverability

**Consistency with ODH Documentation:**
- Observability docs live within existing OpenShift AI documentation structure (not separate Perses docs site)
- URL pattern: `docs.openshift.com/ai/observability/...` (not `docs.perses.dev/...`)
- Cross-reference existing ODH concepts throughout:
  - "See also: [Creating Data Science Projects]" when discussing project-scoped metrics
  - "Related: [Deploying Models with KServe]" in serving metrics documentation
- Consistent terminology (use "Data Science Project" not "namespace" in user-facing docs; use "Model Serving" not "KServe inference service")
- Same navigation structure and visual design as rest of ODH docs

**Progressive Disclosure Strategy:**
- **In-Product (Tier 1):** Tooltips, empty states, error messages → no need to leave UI for common questions
- **Task Guides (Tier 2):** Linked from in-product help or "Learn More" buttons → covers 80% of user needs
- **Reference (Tier 3):** Comprehensive details for power users and administrators → discoverable via search or task guide links

**Multimedia Learning Resources:**
- **Video Walkthroughs:** (5-minute screencasts)
  - "Your First 5 Minutes with Observability Dashboards" (data scientist perspective)
  - "Troubleshooting Model Serving with Dashboards" (MLOps engineer perspective)
  - "Creating a Custom Dashboard" (administrator perspective)
- **Interactive Tutorials:** (step-by-step guided experiences within ODH)
  - "Tutorial: Identify a Training Bottleneck" (uses sample data to teach dashboard navigation)
  - "Tutorial: Compare Two Model Versions" (interactive walkthrough of comparison workflow)
- **Screenshots and Diagrams:**
  - Every documentation page includes annotated screenshots showing real dashboards
  - Architectural diagrams showing how Perses, Prometheus, and ODH Dashboard integrate
  - Workflow diagrams for troubleshooting decision trees

**Search & Discoverability:**
- All observability documentation indexed in main ODH docs search
- Metadata tags by persona: `persona:data-scientist`, `persona:mlops`, `persona:administrator`
- Filtered discovery: "Show me docs for data scientists" filters to relevant guides
- Related links at bottom of each page (auto-generated based on topic similarity)

**Localization Considerations:**
- Documentation must support localization for international customers (initially English; plan for translations)
- Screenshots should use generic example data (not culturally specific content)
- Measurement units clearly labeled (e.g., "milliseconds" not "ms" for accessibility)

### Documentation Maintenance & Versioning

**Version Alignment:**
- Documentation versioned alongside OpenShift AI releases (not Perses upstream versions)
- Clear indicators when features are version-specific: "Available in OpenShift AI 2.5+" badges
- Deprecation notices for features being removed with migration guides

**Community Contribution:**
- Documentation source in Git repository (Markdown or AsciiDoc format)
- Contribution guidelines for community-submitted documentation improvements
- Review process for technical accuracy (engineering team review required)

**Feedback Mechanisms:**
- "Was this helpful?" buttons on each documentation page (track usefulness)
- "Report an issue" link to submit doc bugs or suggest improvements
- Analytics tracking which pages are most visited (identify gaps where users are searching but not finding answers)

**Questions to answer:**

The following refinement and architectural questions must be resolved before coding can begin:

### Architecture & Technical Design

1. **Perses Backend Deployment Model:**
   - Will Perses backend (API server, database) run as separate pods within ODH Dashboard namespace or as a cluster-wide service?
   - What is the recommended HA configuration (replica count, anti-affinity rules, resource requests/limits)?
   - How will Perses backend scale (horizontal pod autoscaling criteria, vertical scaling considerations)?
   - What database will Perses use for dashboard storage (embedded SQLite, PostgreSQL, ConfigMap-based)?

2. **React Component Integration Specifics:**
   - Which Perses npm packages are required (`@perses-dev/components`, `@perses-dev/dashboards`, `@perses-dev/plugin-system`, others)?
   - What version compatibility constraints exist between Perses packages and ODH Dashboard's React/MUI/React Query versions?
   - How will Perses providers (QueryClient, PluginRegistry, ChartsProvider, SnackbarProvider) integrate with ODH Dashboard's existing provider hierarchy?
   - Will Perses components be lazy-loaded (code splitting strategy), and if so, what's the bundle size impact?

3. **Authentication & Authorization:**
   - How will ODH Dashboard session tokens be passed to Perses components (bearer token, cookie, context provider)?
   - Will Perses backend require a service account with specific Kubernetes RBAC permissions (if so, which permissions)?
   - How will namespace-scoped RBAC translate to Prometheus query-level filtering (label enforcer sidecar, Prometheus RBAC proxy, custom middleware)?
   - What happens when user's session expires while viewing dashboards (automatic redirect to login, inline refresh prompt)?

4. **Data Source Configuration:**
   - How will Perses components connect to Prometheus (direct connection, ODH Dashboard backend proxy, Perses backend proxy)?
   - If using proxy, what's the endpoint structure (`/api/v1/observability/query` or similar)?
   - How will Tempo/Loki be configured for tracing/log integration (post-MVP, but architecture should accommodate future extension)?
   - Will datasource configurations be stored in ConfigMaps, Perses backend database, or ODH Dashboard configuration?

5. **Dashboard Storage & Management:**
   - Where are dashboard definitions stored (Kubernetes ConfigMaps with specific labels, Perses backend database, custom CRDs)?
   - How are default/pre-built dashboards shipped (baked into container images, applied via Kustomize overlays, managed by operator)?
   - How do custom user-created dashboards get stored and scoped (per-user, per-project, cluster-wide)?
   - What's the dashboard lifecycle management approach (GitOps-first with UI preview, UI-first with Git sync, hybrid)?

6. **Multi-Tenancy Enforcement:**
   - At what layer is namespace filtering enforced (Prometheus query injection, Perses backend middleware, ODH Dashboard frontend)?
   - How are dashboard permissions managed (Kubernetes RBAC for ConfigMaps, Perses-native RBAC, custom permission model)?
   - Can project admins create project-scoped dashboards, or only cluster admins create all dashboards?
   - How do we prevent users from crafting direct Prometheus queries that bypass namespace filtering?

7. **Performance & Scalability:**
   - What's the caching strategy for dashboard queries (frontend cache, backend cache, Prometheus recording rules)?
   - How many concurrent dashboard queries can a single Perses backend instance handle (load testing targets)?
   - What's the query timeout strategy (client-side, Perses backend, Prometheus query timeout configuration)?
   - How do we prevent expensive queries from degrading Prometheus for other consumers (query cost limits, rate limiting)?

### Integration & Dependencies

8. **Platform Metrics Dependency:**
   - What's the minimum viable Platform Metrics deployment (Prometheus Operator, federation, Thanos, all of the above)?
   - Can Perses integration function in "degraded mode" if Platform Metrics is partially deployed (e.g., Prometheus exists but not Tempo)?
   - How do we detect Platform Metrics availability at runtime (health check endpoint, Kubernetes API query for Prometheus CR)?
   - What error messaging should users see if Platform Metrics is missing vs. misconfigured?

9. **ODH Dashboard Integration Points:**
   - Which ODH Dashboard pages need "View Metrics" contextual links (Model Serving, Data Science Projects, Pipelines - which others)?
   - How are contextual links generated (static routes with query parameters, dynamic dashboard ID lookup)?
   - When embedding panels on detail pages (post-MVP), what's the component API (React component with props, iframe with parameters)?
   - How does navigation breadcrumb integrate (does "Observability" become part of breadcrumb trail, or separate top-level section)?

10. **Upgrade & Migration:**
    - How are dashboard definitions migrated across ODH versions (automated migration scripts, manual admin action, backward compatibility guaranteed)?
    - When Perses npm packages are updated, what's the testing and rollout process (canary deployment, blue-green, rolling update)?
    - How do we handle breaking changes in Perses upstream (pin to stable version, maintain compatibility layer, communicate breaking changes to users)?
    - What's the rollback strategy if a Perses integration update causes issues (separate deployment for Perses backend, feature flag to disable observability)?

### User Experience & Design

11. **Visual Theming:**
    - What specific PatternFly components/tokens will be used to theme Perses (color palette mapping, typography scale, spacing system)?
    - Are there Perses UI elements that can't be themed and will look inconsistent (if so, which ones and is that acceptable)?
    - Who owns the theming implementation (ODH Dashboard frontend team, joint effort, Perses upstream contribution)?
    - How is theming maintained across Perses version updates (CSS overrides, theme provider configuration, Perses theming API)?

12. **Navigation & Information Architecture:**
    - Should "Observability" be a top-level navigation item, or nested under "Administration" or "Platform"?
    - What's the URL structure (`/observability/dashboards/:dashboardId`, `/metrics/...`, other)?
    - How do users navigate between multiple dashboards (dashboard picker dropdown, sidebar list, breadcrumb-based)?
    - Should there be a dashboard "home page" or land directly on a default dashboard?

13. **Contextual Embedding:**
    - When clicking "View Metrics" from Model Serving page, should it open in new tab, navigate within same window, or open in modal overlay?
    - How do we preserve user's context when navigating to observability and back (browser history state, route parameters)?
    - For embedded panels on detail pages (post-MVP), what's the interaction model (inline chart, click to expand full dashboard, hover preview)?

14. **Error Handling & Empty States:**
    - What are all possible error states users might encounter (Prometheus down, query timeout, no data, auth failure, dashboard not found, syntax error in dashboard definition)?
    - For each error state, what's the user-visible message and recovery action (designs needed for each)?
    - How do we distinguish between "no data" (expected, e.g., no training jobs running) vs. "broken" (unexpected, e.g., metrics not collected)?

15. **Accessibility Implementation:**
    - Who is responsible for WCAG 2.1 AA compliance testing (ODH Dashboard QA team, dedicated accessibility consultant, automated tools)?
    - Are there known Perses accessibility issues from upstream that need remediation (if so, document and plan fixes)?
    - What's the keyboard navigation model for complex interactions (chart zoom, panel expand, variable selection)?
    - How will screen readers interpret charts (text alternatives, ARIA labels, data tables as fallback)?

### Operations & Support

16. **Deployment & Configuration:**
    - What Kubernetes resources are required (Deployments, Services, ConfigMaps, Secrets, CRDs, ServiceMonitors)?
    - Are there Helm chart values for configuration (enable/disable observability, Prometheus endpoint, resource limits)?
    - What's the installation flow (operator installs Perses automatically, admin manually deploys, hybrid)?
    - How do admins validate successful installation (health check endpoint, test dashboard, CLI command)?

17. **Monitoring & Alerting:**
    - How do we monitor the health of Perses integration itself (Prometheus metrics for Perses backend, ODH Dashboard integration health checks)?
    - What alerts should fire if observability is degraded (Perses backend down, query latency high, authentication failures)?
    - How do admins troubleshoot issues (logs to check, metrics to query, diagnostic tools)?

18. **Air-Gapped Requirements:**
    - What external dependencies exist at runtime (CDN resources, external APIs, container image pulls)?
    - How are npm packages bundled for disconnected environments (all dependencies in ODH Dashboard bundle, separate artifact)?
    - What documentation is needed for air-gapped installation (specific mirror procedures, image lists, prerequisites)?

19. **Resource Overhead:**
    - What are the CPU/memory requirements for Perses backend (idle, typical load, peak load)?
    - What's the storage requirement for dashboard definitions (assume 50 custom dashboards, 1000 users)?
    - How much additional load does this place on Prometheus (query rate increase, storage impact)?

20. **Support & Troubleshooting:**
    - What are the top 10 predicted support issues (based on Grafana/observability support history)?
    - What diagnostic information should users provide when filing support tickets (logs, screenshots, dashboard definitions, query examples)?
    - How do support engineers access user environments to troubleshoot (must-gather script, diagnostic dashboard)?

### Security & Compliance

21. **Data Access & Privacy:**
    - Can cluster admins view all users' metrics, or are there privacy controls (audit logging, restricted admin roles)?
    - Are there metrics that should never be exposed via dashboards (secrets, PII, sensitive infrastructure details)?
    - How are metric labels sanitized to prevent information leakage (automated label filtering, admin configuration)?

22. **Audit Logging:**
    - What events must be audited (dashboard access, dashboard modification, query execution, authentication failure)?
    - What's the audit log format and storage (Kubernetes audit log, custom audit service, SIEM integration)?
    - How long are audit logs retained (compliance requirement varies by industry - document options)?

23. **Security Scanning & CVEs:**
    - What's the process for responding to CVEs in Perses dependencies (upstream patches, forked fixes, version pinning)?
    - How are Perses container images scanned (ODH pipeline integration, separate process)?
    - What's the SLA for patching critical security issues in observability components?

### Success Metrics & Validation

24. **Adoption Metrics:**
    - How do we measure feature adoption (unique users accessing observability, dashboards viewed per user, time spent in observability section)?
    - What analytics/telemetry do we collect (with user consent/privacy controls)?
    - What's the success criteria for exiting Tech Preview to GA (X% of active users, Y% satisfaction score, <Z% bug rate)?

25. **Performance Validation:**
    - What are the load testing scenarios (100 users, 500 users, 1000 users viewing different dashboards)?
    - What's the acceptance criteria for performance (dashboard load time <3s, query response <2s, P95 latency thresholds)?
    - How do we validate performance doesn't degrade ODH Dashboard overall (baseline metrics, comparative testing)?

26. **User Validation:**
    - What user acceptance testing is planned (beta program, design partner engagement, usability studies)?
    - How many users from each persona must validate workflows (5 data scientists, 5 MLOps engineers, 3 admins minimum)?
    - What feedback mechanisms are built in for continuous improvement (in-app feedback, quarterly surveys, customer advisory board)?

**Background & Strategic Fit:**

### Market Timing & Competitive Landscape

The AI/ML platform market has matured significantly over the past three years, shifting from a focus on model experimentation to operationalizing AI at scale. According to Gartner's 2024 AI Platform Magic Quadrant, operational capabilities - including monitoring, governance, and lifecycle management - now account for 40% of evaluation criteria, up from just 15% in 2021. This shift reflects enterprise buyers' recognition that successful AI initiatives depend not just on training models, but on reliably deploying, monitoring, and maintaining them in production.

**Competitive Pressure:**
We are losing competitive evaluations to platforms with integrated observability capabilities:

- **Databricks:** Added unified monitoring to their ML platform in 2023, now cited as a top-3 differentiator in customer references. Their Lakehouse Monitoring provides automatic drift detection and performance tracking without configuration.

- **Amazon SageMaker:** Model Monitor provides real-time inference monitoring as a default capability, with automatic data quality and bias detection. Customers receive alerts proactively without setting up separate monitoring infrastructure.

- **Google Vertex AI:** Embeds monitoring directly in the model deployment workflow, with pre-built dashboards for common ML frameworks (TensorFlow, PyTorch, Scikit-learn). Users see model performance metrics within seconds of deployment.

Our sales engineering team reports that **observability is a "must-respond" line item in 70% of enterprise RFPs**. When we respond with "use Grafana separately" or "configure Prometheus manually," we score poorly on UX evaluation criteria. Buyers expect integrated, AI-workload-specific observability as a baseline capability, not as a separate tool requiring specialized expertise.

**Deal Impact Analysis:**
- **Q3 2024:** Lost 2 deals (combined $850K ACV) where observability gap was explicitly cited as a decision factor. Post-loss debriefs revealed buyers chose competitors because "their observability was built-in, not bolted-on."
- **Q4 2024:** 5 active opportunities ($2.1M pipeline) have observability as an explicit RFP requirement, with detailed acceptance criteria around dashboard accessibility for data scientists.
- **Win Rate Impact:** Our win rate for deals WITH observability requirements is 28%, compared to our overall win rate of 42%. This 14-percentage-point gap represents significant lost revenue.

**Revenue Opportunity:**
- Improving win rate on observability-dependent deals by 15 percentage points (28% → 43%) would unlock an estimated **$3-5M additional revenue in FY26**.
- Reducing professional services engagements for custom monitoring setup (currently ~$40K per customer, 15 engagements per year) would save **$600K annually** while improving customer satisfaction.

### Strategic Alignment with Product Roadmap

**FY25-26 Strategic Pillar: Enterprise-Grade MLOps**

This feature directly advances our strategic objective of becoming the leading enterprise MLOps platform. Our roadmap progression is:

1. ✅ **Completed (FY23-24):** Model training & experimentation
   - Jupyter notebook environment
   - Distributed training (Ray, Horovod)
   - Experiment tracking (MLflow integration)

2. ✅ **Completed (FY24-25):** Model serving infrastructure
   - KServe integration for scalable inference
   - Multi-model serving with autoscaling
   - A/B testing and canary deployments

3. 🎯 **In Progress (FY25-26):** Observability & lifecycle management ← **This RFE**
   - Unified observability (Perses integration)
   - Model registry with lineage tracking
   - Automated model retraining pipelines

4. 📅 **Planned (FY26-27):** Governance & compliance
   - Model risk management
   - Bias detection & explainability
   - Audit trails & regulatory reporting

**Observability is the foundational layer** that enables lifecycle management and governance. You cannot govern what you cannot measure; you cannot detect drift without monitoring; you cannot optimize what you cannot observe. By delivering observability now, we create the data foundation for the next phase of our strategic roadmap.

### Platform Strategy: Open Hybrid Cloud

**Why Perses Over Proprietary Solutions:**

Selecting Perses as our observability foundation aligns with Red Hat's core strategic differentiator: open hybrid cloud. This decision supports three critical objectives:

1. **Avoiding Vendor Lock-In:**
   - Enterprise customers increasingly resist proprietary monitoring solutions (Datadog, Dynatrace) due to cost unpredictability and data lock-in concerns.
   - A financial services customer (top 5 account) recently told us: "We spent $2M on Datadog last year. If OpenShift AI requires another proprietary tool, we're going with Databricks even though we prefer Red Hat."
   - Perses is a CNCF sandbox project with credible governance, reducing customer concerns about abandonment or forced migration.

2. **Hybrid Deployment Consistency:**
   - Customers run OpenShift AI across on-premises, AWS, Azure, and Google Cloud environments.
   - Perses dashboard definitions (CUE/YAML) are portable across environments, enabling consistent observability regardless of infrastructure.
   - Proprietary solutions often have cloud-specific integrations that don't work consistently across hybrid deployments.

3. **Extensibility for ISV Ecosystem:**
   - We're building an ecosystem of AI/ML tool vendors who integrate with OpenShift AI.
   - Perses's plugin architecture allows these partners to extend observability without fragmenting the user experience.
   - Example: Model monitoring vendors (Arize, Fiddler, WhyLabs) can develop Perses plugins to display their insights within the unified ODH Dashboard interface.

**Open Source Contribution Opportunity:**

By adopting Perses early (it's still in CNCF sandbox), we have an opportunity to influence its roadmap and contribute features valuable to the AI/ML community:
- AI workload-specific panel types (confusion matrix visualizations, ROC curve displays)
- Multi-GPU performance tracking panels
- Integration patterns for ML frameworks (TensorFlow Profiler, PyTorch Profiler)

This positions Red Hat as a leader in cloud-native AI observability, enhancing our brand within the CNCF ecosystem.

### Ecosystem Strategy: Partner Integration

We're seeing significant demand from our partner ecosystem for standardized observability integration points:

**NVIDIA Partnership:**
- NVIDIA wants to showcase GPU utilization dashboards for their DGX customers using OpenShift AI.
- They've requested a standard integration point for GPU-specific metrics (NVLink utilization, tensor core efficiency, multi-instance GPU allocation).
- Perses provides this integration point without requiring NVIDIA to customize Grafana or build proprietary dashboards.

**Model Monitoring Vendor Partnerships:**
- Companies like Arize, Fiddler, and WhyLabs provide specialized model monitoring (drift detection, bias analysis, explainability).
- Currently, each vendor requires customers to use their separate dashboard, creating fragmentation.
- With Perses integration, these vendors can build plugins that surface their insights within OpenShift AI's unified observability interface.
- This "plug-and-play" ecosystem model is attractive to both vendors (easier customer adoption) and customers (consistent UX).

**Red Hat Observability Portfolio Alignment:**
- Red Hat's broader observability strategy includes OpenShift Logging (Loki), OpenShift Distributed Tracing (Tempo), and OpenShift Monitoring (Prometheus).
- Perses natively integrates with all three, creating alignment between OpenShift AI and the broader OpenShift observability story.
- Customers using OpenShift Monitoring for infrastructure already have familiarity with Prometheus/Tempo/Loki; Perses extends that knowledge to AI workloads.

### Technical Foundation: Platform Metrics Architecture

This feature has a critical dependency on the **Platform Metrics and Alerts** architecture currently in development. That initiative provides:

**Centralized Metric Collection:**
- Prometheus federation for multi-tenant metric aggregation
- ServiceMonitor configurations for automatic metric discovery from AI workloads
- Standard metric naming conventions for consistency (e.g., `odh_training_*`, `odh_inference_*`, `odh_pipeline_*`)

**RBAC-Aware Metric Access:**
- Namespace-scoped metric queries enforced at Prometheus layer
- Integration with OpenShift's existing RBAC model
- Audit logging for metric access (compliance requirement)

**Scalable Storage Backend:**
- Thanos for long-term metric retention (90+ days)
- Query performance optimization (recording rules for expensive aggregations)
- Multi-cluster metric federation (supporting hybrid deployments)

**Risk: Platform Metrics Delivery Dependency**

If Platform Metrics architecture slips or is deprioritized, this Perses integration becomes "shelf-ware" - a beautiful UI with no data to display. **Mitigation strategies:**

1. **Phased Release:** Structure as Tech Preview that explicitly requires Platform Metrics, graduating to GA only after Platform Metrics reaches GA stability.

2. **Fallback Mode:** Design Perses integration to gracefully degrade if Platform Metrics is partially deployed (e.g., works with basic Prometheus Operator, degraded without Thanos long-term storage).

3. **Clear Messaging:** Documentation and release notes must clearly communicate the Platform Metrics dependency to avoid customer confusion.

### Customer Pain Points Driving Urgency

**Current-State Customer Frustrations (from 15+ customer interviews, Q3-Q4 2024):**

1. **"We waste hours context-switching between tools"** (8 customers mentioned this)
   - Data scientists start in Jupyter notebooks, jump to Grafana to check GPU utilization, then to Prometheus to query specific metrics, then back to notebooks.
   - "By the time I find the issue, I've lost my flow state and forgotten what I was optimizing for."

2. **"Our data scientists can't troubleshoot performance without help"** (12 customers)
   - Platform teams become bottlenecks because data scientists lack Prometheus/Grafana expertise.
   - "We get Slack messages at 2 AM asking 'is the GPU broken?' when really their batch size is too small. They need self-service observability."

3. **"Generic Kubernetes dashboards don't help us"** (10 customers)
   - Existing Grafana dashboards show pod CPU/memory but don't correlate with AI-relevant metrics like training loss or inference latency.
   - "I need to know 'is my model training efficiently?' not 'what's my pod's memory usage?' - those are different questions."

4. **"We've spent $40K+ on custom Grafana setup"** (5 enterprise customers)
   - Customers hire consultants (often Red Hat Professional Services) to build custom observability.
   - "This feels like something that should be included. We're already paying for OpenShift AI."

5. **"We can't share insights with our team"** (6 customers)
   - Grafana dashboards are admin-only; data scientists can't create or share dashboards.
   - "When I find a performance issue, I take screenshots and paste into Slack. There's no way to share the live dashboard with my team."

**Why Now:**

These pain points have always existed, but three factors make this the right time to address them:

1. **AI Adoption Scale:** Customers are moving from 10-20 data scientists (where manual support is feasible) to 200-500 data scientists (where self-service is mandatory).

2. **Production AI Workloads:** Customers are moving models to production at scale, where observability isn't "nice to have" but "mission critical" for SLA compliance.

3. **Competitive Market:** Alternatives (Databricks, SageMaker) have solved this problem, making it a baseline expectation rather than a differentiator.

### Business Case Summary

**Investment Required:**
- Engineering: 2-3 person-quarters for MVP (1 frontend engineer, 1 backend engineer, 0.5 QA engineer)
- Design: 3-4 weeks UX research & design (dashboard layouts, navigation integration, theming)
- Documentation: 2-3 weeks technical writing (guides, references, troubleshooting)
- Total estimated cost: ~$150-200K (loaded cost)

**Expected Return:**
- **Revenue Impact:** $3-5M incremental revenue in FY26 from improved win rates
- **Cost Savings:** $600K/year reduction in professional services for custom monitoring
- **Customer Satisfaction:** Estimated 8-10 point NPS increase among MLOps engineer persona (currently 45, target 55)
- **Churn Reduction:** Observability gap cited in 2 customer churn postmortems in FY24; addressing this reduces churn risk for ~$500K ARR

**ROI:** Estimated 15-25x return on investment over 2 years.

**Strategic Value (Non-Financial):**
- Closes top-3 competitive gap vs. Databricks/SageMaker
- Enables next phase of roadmap (lifecycle management, governance) which depends on observability foundation
- Strengthens Red Hat's position in CNCF AI/ML ecosystem
- Creates platform for ISV partnerships (model monitoring vendors)

### Key Risks to Monitor

1. **Platform Metrics Delivery Risk (HIGH):**
   - **Risk:** This feature is useless without the data layer. If Platform Metrics slips, this feature creates negative perception ("vaporware").
   - **Mitigation:** Tight coordination with Platform Metrics team; clear Tech Preview → GA gate tied to Platform Metrics maturity.

2. **Perses Project Maturity Risk (MEDIUM):**
   - **Risk:** Perses is CNCF sandbox (not graduated); if project stalls or governance issues arise, we're exposed.
   - **Mitigation:** Active contribution to Perses upstream; fork contingency plan; monitor project health metrics (commit frequency, maintainer engagement, issue resolution time).

3. **Performance at Scale Risk (MEDIUM):**
   - **Risk:** Prometheus query load from 1000+ users could degrade monitoring for entire cluster.
   - **Mitigation:** Comprehensive load testing; query result caching; Prometheus resource sizing guidance; query cost limits.

4. **Adoption Risk (LOW-MEDIUM):**
   - **Risk:** If users don't discover or understand observability features, adoption will be <20%, failing to justify investment.
   - **Mitigation:** Prominent navigation placement; in-product onboarding; contextual links from workflows; video walkthroughs; success team enablement.

5. **Technical Debt Risk (LOW):**
   - **Risk:** Quick integration shortcuts (e.g., iframe instead of React components) create long-term maintenance burden.
   - **Mitigation:** Invest in proper React component integration upfront; allocate time for PatternFly theming; architect for upgradability.

**Customer Considerations:**

Enterprise customers have specific requirements that must be addressed in the design and delivery of this feature to ensure adoption, compliance, and long-term success.

### Multi-Tenancy & Data Isolation (Critical Priority)

**Customer Context:**
Financial services, healthcare, and government customers (representing 35-40% of revenue) operate in highly regulated environments where **data isolation is a compliance requirement**, not just a best practice. These organizations face significant penalties for data breaches or unauthorized access (HIPAA violations: $50K+ per incident; SOX violations: criminal liability for executives; FedRAMP: contract termination).

**Specific Requirements:**

1. **Namespace-Scoped Metric Queries:**
   - A data scientist in project A must never see metrics from project B, even with URL manipulation or crafted API requests.
   - Enforcement must occur at the Prometheus query layer (not just UI hiding), with label matchers automatically injected: `{namespace="project-a"}`.
   - **Validation Required:** Adversarial testing where users attempt to escape namespace scope via crafted PromQL queries, dashboard definition manipulation, or API parameter injection.

2. **Project-Level Dashboard Scoping:**
   - Default dashboards should show only the user's authorized projects.
   - Custom dashboards created by project members should be scoped to that project (not visible cluster-wide unless explicitly shared by admins).
   - **Design Consideration:** Should project admins have ability to create project-scoped dashboards, or only cluster admins? Customer feedback suggests project-level dashboard creation increases adoption.

3. **Audit Trail for Cross-Namespace Access:**
   - If cluster admins view metrics across multiple namespaces, this access must be logged with: admin identity, timestamp, which namespaces' data accessed, duration of access.
   - Audit logs must be tamper-proof and retained per compliance requirements (HIPAA: 6 years; SOX: 7 years; FedRAMP: varies by agency).

**Customer Quote (Healthcare, 800-bed hospital system):**
> "We cannot have data scientists from our oncology research team seeing patient volume metrics from cardiology. That's a HIPAA violation. If your observability doesn't enforce namespace isolation at the query layer, we can't use it."

**Implication for Design:**
- Integration testing must include multi-tenant scenarios with adversarial access attempts.
- Documentation must clearly explain the isolation model and how it maps to Kubernetes RBAC.
- Performance impact of RBAC enforcement (query rewriting, label injection) must be <200ms per query to remain acceptable.

---

### Compliance & Audit Requirements (High Priority)

**Customer Context:**
Regulated industries require comprehensive audit trails for compliance frameworks including HIPAA (healthcare), SOX (finance), FedRAMP (government), ISO 27001 (international), and GDPR (EU data privacy).

**Specific Requirements:**

1. **Dashboard Access Logging:**
   - Log who accessed which dashboards, when, and which projects' data was viewed.
   - Example audit event: `{"timestamp": "2025-02-15T14:32:11Z", "user": "alice@example.com", "action": "view_dashboard", "dashboard": "model-serving-performance", "namespaces": ["project-alpha"], "duration_seconds": 127}`
   - Audit logs must integrate with existing ODH audit framework (if available) or provide standalone audit log stream.

2. **Dashboard Modification Logging:**
   - Log dashboard create/update/delete operations with full diff of changes.
   - Example: `{"timestamp": "2025-02-15T15:45:33Z", "user": "bob-admin@example.com", "action": "update_dashboard", "dashboard": "custom-gpu-metrics", "changes": {"added_panels": ["nvlink-utilization"], "removed_panels": []}}`
   - This supports "who changed what, when, and why" questions during audits.

3. **Metric Query Logging (Optional, High-Security Environments):**
   - Some government customers require logging every Prometheus query executed, including query text and result set size.
   - **Design Consideration:** This can generate massive log volume (1000 users × 10 queries/min = 10K log entries/min). Provide as opt-in capability with clear storage implications documented.

4. **Retention & Immutability:**
   - Audit logs must be stored in immutable storage (write-once, read-many) to prevent tampering.
   - Retention periods vary by regulation: HIPAA (6 years), SOX (7 years), FedRAMP (agency-specific, often 3 years minimum).
   - **Integration Point:** Can OpenShift Logging / Loki be configured for immutable audit log storage? Or does this require separate audit log infrastructure?

**Customer Quote (Financial Services, Fortune 100 bank):**
> "We get audited twice a year for SOX compliance. Auditors want to see who accessed what data and when. If we can't provide audit logs showing data scientist Alice only accessed her authorized projects, we fail the audit."

**Implication for Design:**
- Audit logging should be table-stakes for GA (not post-MVP), given prevalence of regulated industries in customer base.
- Performance impact of audit logging must be negligible (<10ms per operation).
- Documentation must provide sample audit queries (e.g., "show all dashboard access by user X in the last 30 days").

---

### Air-Gapped & Disconnected Environments (Medium-High Priority)

**Customer Context:**
Government (defense, intelligence), financial services (high-security trading floors), and some healthcare customers (25% of customer base) operate in **air-gapped environments** with no internet access. These environments require all dependencies to be mirrorable to internal registries.

**Specific Requirements:**

1. **Container Image Mirroring:**
   - All Perses-related container images (backend, UI, init containers) must be pullable from Red Hat's official registry and mirrorable to customer internal registries.
   - No runtime dependencies on Docker Hub, Quay.io public repos, or other external registries.
   - **Testing Requirement:** CI/CD pipeline must include air-gapped installation testing (cluster with no external network access).

2. **npm Package Bundling:**
   - Perses npm packages (`@perses-dev/*`) must be bundled into ODH Dashboard at build time, not fetched at runtime.
   - No runtime calls to npmjs.com or unpkg.com CDNs.
   - **Verification:** Inspect ODH Dashboard container for external URL references in JavaScript bundle.

3. **Documentation & Dependencies:**
   - Installation documentation must explicitly cover disconnected environments with step-by-step mirroring procedures.
   - List all dependencies (container images, file URLs, package versions) in machine-readable format for automated mirroring tools.
   - **Example:** Provide `image-list.txt` with all container images and their SHA256 digests for mirroring validation.

4. **Dashboard Definition Portability:**
   - Dashboard definitions (YAML/CUE) must not reference external resources (e.g., no `image: https://example.com/logo.png` in dashboard configs).
   - All assets (icons, fonts if any) bundled in ODH Dashboard or Perses containers.

**Customer Quote (US Department of Defense contractor):**
> "Our classification level doesn't allow internet access. If your installation requires pulling from Docker Hub or npm at runtime, we literally cannot install it. Everything must be mirrorable to our internal Nexus repository."

**Implication for Design:**
- Air-gapped support should be P1 for MVP (blocking for 25% of customers).
- Release process must include publishing a "disconnected install bundle" with all dependencies.
- QA must include true air-gapped testing (not just "internet disabled" but physically isolated network).

---

### Single Sign-On & Identity Integration (Critical Priority)

**Customer Context:**
Enterprise customers universally require SSO integration (SAML, OIDC) for all tools. **Separate authentication for Perses is a non-starter** - it creates security audit concerns, user friction, and operational overhead.

**Specific Requirements:**

1. **Seamless Authentication Flow:**
   - Users authenticated to ODH Dashboard should automatically have access to Perses components without seeing a login prompt.
   - Perses backend must accept bearer tokens from ODH Dashboard session (or equivalent authentication mechanism).
   - **User Experience:** Clicking "Observability" in navigation should load dashboards within 2-3 seconds with no intermediate auth screens.

2. **Session Lifecycle Alignment:**
   - Perses session expiration should align with ODH Dashboard session policies (typically 8-12 hours idle timeout).
   - If ODH Dashboard session expires, Perses should gracefully handle this (redirect to login, not cryptic 401 errors).
   - Token refresh should be transparent to users (no interruption to dashboard viewing).

3. **Role-Based Access Propagation:**
   - User's roles in OpenShift (cluster-admin, project-admin, project-viewer) should map to Perses permissions (cluster-wide dashboards, project-scoped dashboards, read-only).
   - **Design Question:** How does OpenShift RBAC map to Perses permission model? (See "Questions to Answer" section)

4. **Multi-Factor Authentication Support:**
   - If customer's IdP requires MFA, Perses integration must not bypass or break this.
   - MFA should be handled at ODH Dashboard layer; Perses integration trusts authenticated session.

**Customer Quote (Global pharmaceutical company, 15K employees):**
> "We use Okta SSO with MFA for everything. If data scientists have to log in separately to Perses, they'll just not use it. And our security team will never approve a separate auth system - that's an audit finding waiting to happen."

**Implication for Design:**
- SSO integration is P0 (critical) for MVP.
- Architecture must clearly define authentication flow (sequence diagram needed).
- Testing must include IdP integration testing (Keycloak, Okta, Azure AD, Ping Identity).

---

### Operational Considerations

#### Upgrade Path & Backward Compatibility (High Priority)

**Customer Context:**
Enterprise customers upgrade OpenShift AI **2-4 times per year** (major releases, security patches). Dashboards represent significant configuration investment (custom dashboards, tuned queries, organizational standards). **Loss of dashboard configurations during upgrade is unacceptable.**

**Specific Requirements:**

1. **Dashboard Definition Persistence:**
   - Dashboard definitions must be stored separately from Perses runtime (e.g., Kubernetes ConfigMaps, not baked into containers).
   - Upgrades must preserve existing dashboards with automatic schema migration if format changes.
   - **Rollback Consideration:** If upgrade fails, rollback should restore previous dashboard state without manual intervention.

2. **Backward Compatibility:**
   - Dashboard definitions created in version N should work in version N+1 without modification (or with automatic migration).
   - If breaking changes are unavoidable, provide migration tooling: `percli migrate --from 1.0 --to 2.0 dashboard.yaml`.
   - **Deprecation Policy:** Give customers at least 2 major versions' notice before removing deprecated dashboard features.

3. **Upgrade Testing:**
   - Automated testing must validate dashboard compatibility across upgrade paths (N → N+1, N → N+2).
   - Include complex dashboards with custom queries, variables, and multiple panels in test suite.

4. **Documentation:**
   - Release notes must clearly document dashboard schema changes and migration steps.
   - Provide upgrade runbook: pre-upgrade checklist, backup procedures, rollback steps.

**Customer Quote (Insurance company, 5K data scientists):**
> "We have 80 custom dashboards that took us 6 months to perfect. If an OpenShift AI upgrade wipes those out, we're not upgrading. And if we can't upgrade, we're stuck on an unsupported version, which our security team won't allow."

**Implication for Design:**
- Dashboard storage architecture must prioritize upgrade-safety (externalized configuration, schema versioning).
- QA must include upgrade testing as first-class concern (not just "fresh install works").

---

#### Resource Overhead & Capacity Planning (Medium Priority)

**Customer Context:**
Enterprise customers are sensitive to **control plane resource consumption**. They've sized OpenShift clusters based on expected workload, and unexpected overhead from new features can cause capacity crunches or budget overruns (cloud costs).

**Specific Requirements:**

1. **Resource Documentation:**
   - Clearly document Perses backend resource requirements (CPU, memory, storage) at different scales:
     - Small deployment (100 users): X CPUs, Y GB memory
     - Medium deployment (500 users): X CPUs, Y GB memory
     - Large deployment (1000+ users): X CPUs, Y GB memory
   - Provide resource request/limit recommendations for Kubernetes manifests.

2. **Incremental Cost Transparency:**
   - Document incremental impact on Prometheus (query load, storage growth).
   - Example: "1000 users viewing dashboards generates ~500 queries/min to Prometheus, requiring additional 2 CPU cores and 4 GB memory allocation."

3. **Capacity Planning Guidance:**
   - Provide sizing calculator or guidelines: "Estimate 50 MB memory per 100 concurrent dashboard viewers."
   - Include autoscaling recommendations (HPA configuration for Perses backend).

4. **Monitoring the Monitor:**
   - Provide dashboards (meta!) showing Perses backend health: query latency, error rate, resource usage.
   - Alert templates for Perses degradation (backend pod restarts, query timeouts >10% of requests).

**Customer Quote (SaaS company, cost-conscious engineering culture):**
> "We run OpenShift on AWS, so every CPU and GB of memory has a real dollar cost. Before we enable a new feature, we need to know: what's this going to cost us at our scale? If it's $10K/year, fine. If it's $100K/year, we need to budget for it."

**Implication for Design:**
- Performance benchmarking and resource profiling should be part of QA (not just functional testing).
- Documentation must include capacity planning section with cost implications.

---

#### Backup & Disaster Recovery (Medium Priority)

**Customer Context:**
Enterprise customers have RTO (Recovery Time Objective) and RPO (Recovery Point Objective) requirements for all systems. For non-production environments, typical targets are **RTO <4 hours, RPO <24 hours**. For production observability, even stricter (RTO <1 hour, RPO <1 hour).

**Specific Requirements:**

1. **Dashboard Configuration Backup:**
   - Dashboard definitions (ConfigMaps, CRDs, or database) must be includable in standard OpenShift backup tools (Velero, OADP).
   - Backup frequency: typically daily for dev/test, hourly for production.
   - **Testing:** Restore testing must validate dashboards are functional after restore (not just "files exist").

2. **Disaster Recovery Procedures:**
   - Document complete DR procedure: backup, restore, validation.
   - Include recovery from common scenarios: namespace deletion, cluster migration, database corruption.
   - **RTO Target:** Dashboard functionality restored within 4 hours of initiating recovery.

3. **State Management:**
   - Identify all stateful components: dashboard definitions, user preferences (if stored), cached data.
   - Document what's lost if backup is N hours old (e.g., dashboards created in last N hours not recovered).

**Customer Quote (Healthcare provider, disaster recovery focus):**
> "We test DR quarterly. If observability dashboards aren't included in our backup, and we lose the cluster, we'll be flying blind during recovery. That's unacceptable for a production system."

**Implication for Design:**
- Documentation must include backup/restore procedures as first-class content (not buried in appendix).
- QA should include restore testing (automate backup → destroy → restore → validate cycle).

---

### User Experience Considerations

#### Persona-Specific Defaults (High Priority)

**Customer Context:**
Data scientists, MLOps engineers, and administrators have **different mental models and information needs**. A one-size-fits-all dashboard approach creates cognitive overload for some personas while under-serving others.

**Specific Requirements:**

1. **Role-Based Default Dashboards:**
   - Data scientists logging in should see training and notebook dashboards prominently.
   - MLOps engineers should see serving and pipeline dashboards.
   - Administrators should see platform health and capacity dashboards.
   - **Implementation:** Dashboard metadata includes `roles: [data-scientist, mlops, admin]` tag; ODH Dashboard filters visible dashboards based on user's primary role.

2. **Customizable Home View:**
   - Users should be able to set a "home dashboard" that loads by default when clicking "Observability."
   - Teams should be able to set organization default (e.g., "all data scientists in finance dept start with Finance Team GPU Usage dashboard").

3. **Progressive Disclosure:**
   - Novice users see simplified dashboards with fewer panels and more contextual help.
   - Advanced users can enable "advanced mode" with more detailed metrics and technical terminology.
   - **Design Consideration:** Should this be automatic (based on usage patterns) or manual (user preference toggle)?

**User Research Validation (needed before final design):**
- Conduct usability testing with 5-10 users per persona to validate dashboard relevance and layout.
- Measure time-to-insight: "How long does it take a data scientist to identify their training bottleneck?"

**Implication for Design:**
- Dashboard metadata schema must support role tagging and priority ordering.
- ODH Dashboard navigation must filter/sort dashboards based on user context.

---

#### Learning Curve & Self-Service (High Priority)

**Customer Context:**
Data scientists are experts in statistics and machine learning, but **not Prometheus, PromQL, or infrastructure monitoring**. Steep learning curves reduce adoption; if users can't understand dashboards within 5 minutes, they'll abandon the feature.

**Specific Requirements:**

1. **Self-Explanatory Dashboards:**
   - Every chart must have a clear title that explains what it measures (not just metric name).
   - Example: ❌ "`container_cpu_usage_seconds_total`" → ✅ "CPU Usage per Training Job (cores)"
   - Tooltips must explain healthy ranges: "GPU utilization 80-95% indicates efficient training. <50% suggests data loading bottlenecks."

2. **Contextual Help:**
   - Dashboard pages include "How to interpret this dashboard" expandable section.
   - Link to documentation directly from dashboards (e.g., "Learn more about training metrics" button).

3. **No PromQL Required:**
   - Pre-built dashboards should cover 80% of use cases without requiring users to write queries.
   - For custom dashboards (post-MVP), provide query builder UI (select metric, aggregation, filters) rather than raw PromQL editor.

4. **Onboarding:**
   - First-time users see brief walkthrough (2-3 step overlay): "This is observability → Click here to see training metrics → Adjust time range here."
   - Sample dashboards with synthetic data for learning (users can explore without running real workloads).

**Usability Testing (required for MVP validation):**
- Test with 10+ data scientists unfamiliar with Prometheus.
- Success criteria: 80% can identify "is my training job using GPUs efficiently?" within 3 minutes.

**Implication for Design:**
- Significant investment in UX polish (tooltips, help text, labeling) is required - not just functional dashboard rendering.
- Technical writing must collaborate closely with UX design for in-product help.

---

### Integration & Ecosystem Considerations

#### Existing Workflow Integration (Critical Priority)

**Customer Context:**
Users spend **80% of their time in specific workflow contexts**: Jupyter notebooks (data scientists), model serving pages (MLOps), pipeline execution logs (data engineers). Observability must be **contextually accessible** from these workflows, not siloed in a separate section.

**Specific Requirements:**

1. **Contextual "View Metrics" Links:**
   - Model serving page: each deployed model has "Observability" tab or "View Metrics" button → opens Model Serving Performance dashboard filtered to that model.
   - Training job page: notebook or training job list includes "View Training Metrics" link → opens Model Training Metrics dashboard filtered to that job.
   - Pipeline run details: "View Pipeline Metrics" → opens Pipeline Execution Health dashboard filtered to that run.
   - **UX Goal:** Users reach relevant metrics in 1-2 clicks from their current context.

2. **Embedded Metrics Summary:**
   - Data Science Project overview page includes high-level metrics summary (e.g., "GPU utilization: 78%, Training jobs: 3 running, Serving latency p95: 120ms").
   - Clicking on summary panel navigates to detailed dashboard.
   - **Design Consideration:** Should summaries be live-updated (WebSocket/polling) or static (refresh on page load)?

3. **Bidirectional Navigation:**
   - From dashboard, users should be able to navigate back to source context (e.g., "View model deployment" link from serving dashboard → model serving page).
   - Breadcrumb trail should reflect journey: "Data Science Projects → my-project → bert-model → Observability".

**Customer Quote (Data scientist, 5+ years experience):**
> "I shouldn't have to remember to go check a separate monitoring tool. If I'm looking at my model deployment and it's slow, there should be a button right there that shows me why it's slow."

**Implication for Design:**
- ODH Dashboard team must implement contextual navigation (not just Perses integration team).
- URL structure must support deep linking with filter parameters (e.g., `/observability/serving?model=bert-v2&namespace=my-project`).

---

#### Custom Metric Integration (High Priority)

**Customer Context:**
Advanced users instrument **custom business metrics** specific to their domain (e.g., "loan approval model bias score," "fraud detection false positive rate," "recommendation click-through rate"). These metrics must be visible in dashboards alongside infrastructure metrics.

**Specific Requirements:**

1. **Custom Metric Support:**
   - Users can emit custom Prometheus metrics from their training/serving code (using Prometheus client libraries).
   - Custom metrics automatically discovered by Platform Metrics (via ServiceMonitor annotations).
   - Custom metrics available in dashboard queries just like built-in metrics.

2. **Documentation & Examples:**
   - End-to-end example showing:
     1. Instrumenting Python training code with `prometheus_client` library
     2. Configuring ServiceMonitor for custom metric scraping
     3. Creating dashboard panel displaying custom metric
   - Example for common ML frameworks: TensorFlow, PyTorch, Scikit-learn
   - Example for model serving: KServe predictor with custom metric

3. **Metric Naming Best Practices:**
   - Document naming conventions to avoid collision (e.g., customer-specific prefix: `mycompany_ml_*`).
   - Provide metric naming validator (CLI tool or web form) to check for conflicts.

**Customer Quote (MLOps engineer, financial services):**
> "We have compliance metrics like 'bias score by demographic group' that we must track for regulatory reporting. If we can't include those in our dashboards, we'll need a separate tool anyway."

**Implication for Design:**
- Custom metric support is not a Perses-specific requirement, but documentation must show how Perses + Platform Metrics + user code fit together.
- Sample code and templates should be provided (not just "read Prometheus docs").

---

### Support & Troubleshooting Considerations

#### Common Support Issues (High Priority)

**Customer Context:**
Based on Grafana and Prometheus support history, **"dashboard shows no data"** will be the #1 support ticket type (40-50% of tickets). Effective troubleshooting documentation reduces support burden and improves customer satisfaction.

**Specific Requirements:**

1. **Troubleshooting Decision Tree:**
   - Provide clear decision tree in documentation:
     ```
     Dashboard shows no data?
     ├─ Is Platform Metrics deployed? → No → Deploy Platform Metrics
     ├─ Are workloads emitting metrics?
     │  ├─ Check ServiceMonitor: kubectl get servicemonitor -n <namespace>
     │  ├─ Test Prometheus scrape: curl http://<pod-ip>:8080/metrics
     ├─ Is time range appropriate? → Adjust to "Last 6 hours"
     ├─ Is user's RBAC correct?
     │  ├─ Test: kubectl auth can-i get pods --namespace <namespace>
     │  ├─ Verify: User is member of project
     └─ Is Prometheus healthy?
        └─ Check: kubectl get prometheus -n openshift-ai-monitoring
     ```

2. **Diagnostic Tools:**
   - Provide must-gather script for support tickets: collects Perses logs, dashboard definitions, Prometheus status, RBAC bindings, sample query results.
   - CLI command for self-service diagnostics: `odh-diagnostics observability` → runs health checks, reports issues.

3. **Error Message Clarity:**
   - Replace technical errors with user-friendly messages:
     - ❌ "Query error: timeout exceeded (context deadline exceeded)"
     - ✅ "Dashboard query took too long. Try reducing the time range to less than 7 days or contact your administrator if this persists. [Troubleshooting Guide]"

**Support Runbook (for Red Hat support engineers):**
- Top 10 failure scenarios with resolution steps:
  1. "No data displayed" → Check Platform Metrics, ServiceMonitor, RBAC
  2. "Dashboard not loading" → Check Perses backend logs, browser console errors
  3. "Slow dashboard performance" → Check query complexity, time range, Prometheus load
  4. "Authentication failed" → Check ODH Dashboard session, SSO configuration
  5. ... (expand to 10 scenarios)

**Implication for Design:**
- Error handling must provide actionable guidance (not just error codes).
- Diagnostic tooling should be built into product (not ad-hoc scripts created during support incidents).

---

### Summary: Enterprise Readiness Checklist

Before declaring this feature GA-ready, validate all customer considerations are addressed:

**Security & Compliance:**
- ✅ Multi-tenant namespace isolation enforced at query layer (adversarial testing passed)
- ✅ Audit logging operational for regulated industries (integrated with ODH audit framework)
- ✅ Air-gapped installation tested and documented
- ✅ SSO integration working with major IdPs (Keycloak, Okta, Azure AD tested)

**Operational Excellence:**
- ✅ Upgrade path tested (N → N+1 version preserves dashboards)
- ✅ Resource requirements documented (CPU, memory, storage at 100/500/1000 user scales)
- ✅ Backup/restore procedures documented and tested (RTO <4 hours validated)
- ✅ HA configuration validated (99.5% uptime target achieved in testing)

**User Experience:**
- ✅ Persona-specific default dashboards validated with user research (5+ users per persona)
- ✅ Contextual navigation implemented (1-2 clicks from workflow to relevant dashboard)
- ✅ Self-explanatory dashboards with tooltips and help text (usability testing passed)
- ✅ Documentation complete (getting started, troubleshooting, reference)

**Integration & Ecosystem:**
- ✅ Platform Metrics integration validated (Prometheus/Tempo/Loki connectivity tested)
- ✅ Custom metric support documented with end-to-end example
- ✅ ODH Dashboard integration seamless (visual consistency, navigation coherence)

**Support Readiness:**
- ✅ Troubleshooting decision tree published
- ✅ Diagnostic tooling available (must-gather script, health check CLI)
- ✅ Support runbook complete for top 10 failure scenarios
- ✅ Error messages clear and actionable

Only when all checkboxes are validated should this feature graduate from Tech Preview to General Availability.
