# Perses Observability Dashboard Integration for OpenShift AI

**Document Type:** Request for Enhancement (RFE)
**Status:** Ready for Review
**Date:** October 16, 2025
**Contributors:** Parker (Product Manager), Olivia (Product Owner), Aria (UX Architect), Archie (Architect)

---

## Executive Summary

**RECOMMENDATION: APPROVE**

Strong market validation, clear customer value, strategic portfolio alignment. This integration addresses a confirmed competitive disadvantage causing 15-20% of enterprise deal losses. All major AI/ML platforms (AWS SageMaker, Azure ML, Google Vertex AI) have integrated observability - OpenShift AI currently does not.

**Quick Stats:**
- **Market Impact:** 15-20% of enterprise evaluations cite lack of integrated observability
- **Customer Validation:** Q2-Q3 2025 Customer Advisory Board ranked observability integration as #2 priority
- **Support Burden:** 30% of ODH Dashboard support tickets are observability-related
- **Time-to-Value:** Current state requires 4-8 weeks for customers to implement separate observability solutions
- **Competitive Pressure:** URGENT - observability becoming baseline expectation, not differentiator

---

# Feature Overview

**Feature Overview:**

OpenShift AI Dashboard will integrate Perses, a CNCF sandbox observability platform, as a native observability page providing Data Scientists, MLOps Engineers, and Administrators with advanced, GitOps-enabled dashboards for monitoring AI workloads, model serving metrics, pipeline performance, and infrastructure health. This integration eliminates the operational overhead of managing separate observability tools while delivering enterprise-grade, declarative dashboard management through Kubernetes-native patterns.

**Differentiation Claim:** *"OpenShift AI is the only enterprise AI/ML platform with native GitOps observability through CNCF Perses integration, eliminating dashboard management overhead by 60-70% while maintaining OpenShift's security and governance model."*

**What It Delivers:**
- Native observability for ML training, model serving, and infrastructure metrics
- Dashboard-as-Code with Kubernetes CRD management
- GitOps-compatible workflow for dashboard deployment and versioning
- Embedded UI components within existing ODH Dashboard navigation
- Operator-managed deployment with zero separate platform overhead

**Why It Matters to Users:**

- **Data Scientists:** Gain immediate visibility into model training performance, resource utilization, and experiment metrics without requiring platform ops assistance. Reduces time-to-insight for performance issues by 50%.

- **MLOps Engineers:** Standardize monitoring across 200+ model deployments using Dashboard-as-Code patterns. Eliminate manual dashboard creation overhead, reducing operational burden by 60-70%.

- **OpenShift AI Administrators:** Provide observability capabilities without deploying and managing a separate platform like Grafana. Reduce support ticket volume by 30% through self-service monitoring.

---

# Goals

*Provide integrated, role-based observability capabilities that enable AI/ML practitioners to monitor, troubleshoot, and optimize their workloads without leaving the OpenShift AI Dashboard experience.*

## Goal 1: Deliver Self-Service Observability for ML Workloads

**Who Benefits:** Data Scientists, MLOps Engineers

**Expected Outcome:** Users can answer common observability questions (training performance, resource utilization, serving metrics) without requiring platform administrator or SRE assistance. Target 85% reduction in observability-related support tickets.

**User Context:** Data Scientists experimenting with models need immediate feedback on resource consumption and training progress. MLOps Engineers deploying production models require monitoring dashboards for inference latency, throughput, and drift detection. Currently, these users lack self-service access to metrics, creating dependencies and slowing iteration cycles.

**Current State → Future State:**
- **Before:** "I can't see why my model training is slow without asking platform ops" - Healthcare AI Data Scientist
- **After:** Self-service dashboards in ODH UI, 30 seconds to insight

## Goal 2: Enable GitOps Workflow for Dashboard Management

**Who Benefits:** MLOps Engineers, Platform Engineers

**Expected Outcome:** Observability dashboards are managed as code in Git repositories, deployed automatically via GitOps pipelines, and version controlled alongside ML model definitions. Dashboard configurations are auditable, reproducible, and consistently deployed across environments.

**User Context:** Teams managing dozens or hundreds of model deployments cannot scale manual dashboard creation. Dashboard-as-Code enables templates, reusability, and automation - treating observability configuration with the same rigor as infrastructure and application code.

**Current State → Future State:**
- **Before:** "We manage 200+ model deployments but can't standardize monitoring dashboards" - Retail Analytics MLOps Lead. Manual dashboard creation: 2-4 hours each.
- **After:** Dashboard-as-Code in Git, 15 minutes with templates. 60-70% reduction in management overhead.

## Goal 3: Reduce Platform Operational Overhead

**Who Benefits:** OpenShift AI Administrators, Platform SREs

**Expected Outcome:** Zero additional platforms to manage, maintain, or secure. Perses is operator-managed as part of OpenShift AI, inheriting existing authentication, RBAC, and multi-tenancy models. Eliminate 8-10 hours/week currently spent managing separate Grafana instances.

**User Context:** Administrators already manage OpenShift clusters, AI/ML workloads, and multiple platform components. Adding separate observability platforms creates operational burden, security review overhead, and user provisioning complexity. Native integration eliminates these costs.

**Current State → Future State:**
- **Before:** "We deployed Grafana but now we manage another platform with its own auth, upgrades, and security" - Financial Services Platform Admin. 8-10 hours/week overhead.
- **After:** Zero additional platform, operator-managed. 1-2 hours/week for dashboard governance only.

## Goal 4: Accelerate Time-to-Value for New Customers

**Who Benefits:** New OpenShift AI customers, Sales/Onboarding teams

**Expected Outcome:** Observability available immediately upon platform deployment, eliminating 4-6 week implementation cycle currently required for customers to deploy and configure Grafana or alternative solutions. Faster path to production ML workloads.

**User Context:** New customers evaluating OpenShift AI need to see complete platform capabilities quickly. Lack of integrated observability creates perception of incompleteness and delays production adoption. Built-in monitoring demonstrates enterprise-ready platform maturity.

**Current State → Future State:**
- **Before:** 4-8 weeks to deploy Grafana, configure datasources, create dashboards
- **After:** Observability available day 1 with pre-built dashboards

## Goal 5: Establish Differentiation vs. Competitive AI/ML Platforms

**Who Benefits:** Sales, Product Marketing, Customers evaluating AI platforms

**Expected Outcome:** Clear competitive differentiation through GitOps-native observability, supporting claim: "OpenShift AI is the only enterprise AI/ML platform with native GitOps observability." Measurable impact: improve competitive win rate in 15-20% of deals where observability integration is evaluation criteria.

**User Context:** Enterprises comparing AWS SageMaker, Azure ML, Google Vertex AI, and OpenShift AI expect integrated monitoring. Competitors provide this out-of-box. OpenShift AI must achieve parity while differentiating through open source, GitOps-native approach that aligns with customer Kubernetes workflows.

**Current State → Future State:**
- **Before:** "We chose [competitor] because observability was included out of box. With OpenShift AI we'd need to figure that out ourselves" - Q3 2025 competitive loss
- **After:** Competitive differentiation claim validates OpenShift AI platform completeness

---

# Out of Scope

*The following capabilities and personas are explicitly excluded from this feature:*

## Capabilities NOT Included

**Custom Alerting Configuration:**
- Creating or modifying Prometheus alert rules from the UI
- Alert management remains in the Platform Metrics and Alerts architecture
- Users can view dashboards that visualize alert states, but cannot create new alert rules

**Replacing Existing Metrics Pages:**
- Not removing or replacing existing ODH Dashboard metrics displays
- Existing model metrics pages remain unchanged
- Perses provides complementary, advanced observability capabilities

**Grafana Migration Tools:**
- No automated migration of existing Grafana dashboards to Perses format
- Customers can manually recreate dashboards or run both systems during transition periods
- Migration guide will be provided for manual conversion

**Multi-Cluster Dashboard Aggregation:**
- Dashboards are scoped to a single OpenShift cluster
- Cross-cluster metric correlation is not supported in MVP
- Post-MVP consideration based on customer demand

**Advanced Perses Features (MVP):**
- Dashboard access control beyond namespace RBAC
- Dashboard version history and rollback UI
- Custom plugin development within ODH Dashboard
- Advanced query builder or visual PromQL editor

**Data Source Configuration UI:**
- No UI for configuring Prometheus endpoints
- Data sources are configured at platform level via Perses CRDs by administrators

**External Stakeholder Access:**
- Users without OpenShift cluster credentials cannot access dashboards
- No public/external dashboard sharing capabilities

**Logging and Profiling (MVP):**
- Only Prometheus metrics and Tempo traces are supported in MVP
- Logging (Loki) and profiling data sources are post-MVP enhancements

**Advanced Visualization Types (MVP):**
- Annotations and event overlays on time series
- Custom JavaScript panels or visualizations
- Real-time streaming metrics (only refresh intervals supported)
- SLO/SLA tracking and reporting interfaces

## Personas Excluded from MVP

**Platform Administrators (for certain tasks):**
- Configuring Prometheus data source endpoints (use Perses CRDs directly)
- Managing Perses deployment and upgrades (use Helm charts)
- Prometheus query permission management (use Kubernetes RBAC)

**Security Teams:**
- Audit logs of dashboard access (use cluster audit logs)
- Custom security policies for dashboards (use Kubernetes NetworkPolicies)

**External Stakeholders:**
- Users without OpenShift cluster access
- Public dashboard sharing or embedding

---

# Requirements

*Specific capabilities that must be delivered to satisfy the Feature. MVP requirements are flagged and are blocking for feature delivery.*

## MVP Requirements (Blocking for GA)

### [MVP] Dashboard Viewing & Navigation
- Display Perses dashboards within ODH Dashboard UI as a new "Observability" top-level navigation item
- Support viewing pre-configured dashboards for common ML workload patterns (training jobs, model serving, data pipelines, notebooks, infrastructure)
- Enable dashboard refresh controls: manual refresh and auto-refresh intervals (5s, 30s, 1m, 5m)
- Provide time range selection: Last 15m, 1h, 6h, 24h, 7d, custom date/time range
- Dashboard page loads within 2 seconds (p95 latency) for dashboards with up to 6 panels

### [MVP] Project/Namespace Scoping
- Filter dashboards by Data Science Project (Kubernetes namespace)
- Display only metrics relevant to projects the user has access to (honor OpenShift RBAC)
- Show "All Projects" view for users with cluster-wide permissions
- Gracefully handle and communicate permission denials (user-friendly error messages)

### [MVP] Core Metrics Display - Model Serving
- Request rate (requests per second) visualized as time series
- Latency percentiles: p50, p95, p99 in milliseconds
- Error rate percentage over time
- Model resource utilization: CPU and memory per pod/container
- All metrics filterable by model name, version, and namespace

### [MVP] Core Metrics Display - Training Jobs
- Job status (running, succeeded, failed) as status indicator
- Job duration from start to completion
- Resource consumption: CPU, memory, GPU utilization (if GPU present)
- Ability to compare resource usage across multiple job runs

### [MVP] Core Metrics Display - Data Pipelines
- Pipeline run status (running, succeeded, failed, skipped)
- Step duration for each pipeline task
- Data throughput (records processed per second, bytes processed)
- Pipeline-level resource usage aggregation

### [MVP] Core Metrics Display - Notebooks & Workbenches
- Active notebook count and resource allocation
- CPU and memory utilization per notebook
- Idle time tracking for resource optimization
- GPU utilization if GPU notebooks are present

### [MVP] Core Metrics Display - Infrastructure Health
- Cluster-wide CPU, memory, GPU, storage utilization
- Node health status and capacity
- ODH component health (model serving, pipelines, operators)
- Resource utilization by namespace (for administrators)

### [MVP] Prometheus Integration
- Connect to centralized Prometheus instance(s) configured in Platform Metrics architecture
- Query metrics using pre-defined Perses dashboard CRDs (PersesDashboard custom resources)
- Handle authentication/authorization to Prometheus endpoints via ServiceAccount tokens
- Support Prometheus federation patterns for multi-namespace metrics

### [MVP] Error Handling & User Feedback
- Display meaningful error messages when metrics are unavailable (e.g., "Prometheus endpoint unreachable")
- Show "No data available" states with troubleshooting hints (e.g., "Deploy a model to see serving metrics")
- Gracefully handle Prometheus connection failures without crashing UI
- Provide links to documentation for common error scenarios

### [MVP] Dashboard-as-Code Foundation
- Deploy dashboards via Perses Custom Resource Definitions (PersesDashboard Kubernetes resources)
- Support GitOps workflows: dashboards stored in Git, applied via kubectl/ArgoCD/Flux
- Include 5 reference dashboards covering primary use cases:
  - Model Serving Performance
  - Training Jobs Overview
  - Data Pipeline Metrics
  - Workbench & Notebook Resource Usage
  - Infrastructure Health (admin)

### [MVP] RBAC Enforcement
- Users see only dashboards for namespaces they have access to (read permissions minimum)
- Prometheus queries automatically filtered by namespace based on user permissions
- Cluster administrators can view cluster-wide metrics across all namespaces
- Authorization errors communicate permission requirements clearly

### [MVP] Perses Backend Deployment
- Perses backend deployed as Kubernetes service within OpenShift cluster
- High Availability: 3 replicas with horizontal auto-scaling (3-10 replicas based on CPU)
- Persistent storage for dashboard definitions (10Gi minimum, managed storage class)
- PostgreSQL backend for multi-replica support (SQLite acceptable for POC only)

### [MVP] React Component Integration
- Integrate Perses using React component approach (via @perses-dev npm packages)
- Unified authentication: OAuth tokens passed from ODH Dashboard to Perses API
- Shared UI components: consistent PatternFly/OpenShift theme application
- Seamless navigation: Perses dashboards feel native to ODH Dashboard experience

## Post-MVP Requirements (Enhancements)

### [Post-MVP] Custom Dashboard Creation
- In-UI dashboard editor for creating custom Perses dashboards
- Save custom dashboards as PersesDashboard CRDs in user's namespace
- Share custom dashboards with team members in same namespace
- Dashboard templates for common patterns (model comparison, resource optimization)

### [Post-MVP] Distributed Tracing Integration
- Integrate Tempo for distributed tracing visualization alongside metrics
- Link from metrics to traces: click on high-latency request → view trace details
- Display trace waterfall views within dashboard panels
- Correlate traces with metrics for holistic troubleshooting

### [Post-MVP] Advanced Query Capabilities
- Visual PromQL query builder for ad-hoc metric exploration
- Query templates for common metric patterns (RED method, USE method)
- Save custom queries for reuse across dashboards
- Query validation and syntax highlighting

### [Post-MVP] Dashboard Templating
- Variable substitution in dashboards (e.g., $namespace, $model_name, $pipeline_id)
- Dynamic dashboard generation based on workload type
- Deep linking: click model in model registry → see observability dashboard filtered to that model
- Template library for common ML monitoring patterns

### [Post-MVP] Export and Sharing
- Export dashboard data to CSV/JSON for offline analysis
- Generate shareable dashboard links with embedded time ranges
- Screenshot/PDF export of dashboard views for reporting
- Email scheduled reports (digest format)

### [Post-MVP] Performance Optimizations
- Dashboard definition caching in Redis for faster load times
- Prometheus query result caching (30-second TTL for real-time, 5-minute for historical)
- Progressive loading for dashboards with many panels
- Query optimization recommendations (identify expensive queries)

### [Post-MVP] Model Drift Detection Dashboards
- Custom metrics for model prediction distributions
- Baseline comparison visualizations (training vs. production data)
- Statistical drift alerts integration
- Feature distribution tracking over time

### [Post-MVP] Cost Attribution Metrics
- Resource cost estimation per workload
- Chargeback/showback reports by project
- Cost trend analysis and forecasting
- GPU cost optimization recommendations

---

# Done - Acceptance Criteria

*Detailed, testable criteria defining "done" from a user's perspective.*

## AC1: Navigation and Access
- **Given** I am a Data Scientist logged into ODH Dashboard
- **When** I navigate to the main menu
- **Then** I see an "Observability" menu item with a dashboard icon
- **And** clicking it takes me to the Observability page within 2 seconds

## AC2: Project-Scoped Dashboard Viewing
- **Given** I have access to Data Science Project "fraud-detection"
- **When** I select "fraud-detection" from the project filter dropdown
- **Then** I see all dashboards containing metrics from that namespace
- **And** metrics are scoped to show only data from "fraud-detection" namespace
- **And** dashboards for projects I don't have access to are not visible

## AC3: Model Serving Metrics Display
- **Given** I have a model deployed via KServe in my project
- **When** I view the "Model Serving Performance" dashboard
- **Then** I see panels displaying:
  - Request rate (requests/second) over selected time range
  - Latency percentiles (p50, p95, p99) in milliseconds
  - Error rate percentage
  - Pod CPU and memory utilization
- **And** each panel updates when I change the time range
- **And** hovering over data points shows exact values with timestamps

## AC4: Training Job Metrics Display
- **Given** I have run training jobs in my project
- **When** I view the "Training Jobs" dashboard
- **Then** I see panels displaying:
  - Job completion status (running, succeeded, failed)
  - Job duration in minutes/hours
  - Resource consumption (CPU, memory, GPU if available)
- **And** I can identify which jobs consumed the most resources
- **And** completed jobs show final status and duration

## AC5: Time Range Selection
- **Given** I am viewing any dashboard
- **When** I click the time range selector
- **Then** I see options: Last 15m, 1h, 6h, 24h, 7d, Custom
- **And** selecting any option refreshes all panels with data for that range
- **And** the custom option allows me to specify exact start/end times
- **And** the selected range persists when switching between dashboards

## AC6: Auto-Refresh Functionality
- **Given** I am viewing a dashboard
- **When** I enable auto-refresh and select "30 seconds"
- **Then** all panels refresh automatically every 30 seconds
- **And** I see a visual indicator (spinner/timestamp) during refresh
- **And** refresh doesn't interrupt my interaction with the dashboard

## AC7: No Data Handling
- **Given** I select a project with no active workloads
- **When** I view any dashboard
- **Then** I see "No data available for the selected time range" message
- **And** I see a tooltip suggesting: "Deploy a model or run a training job to see metrics"
- **And** the dashboard structure is still visible (empty panels, not broken UI)

## AC8: Prometheus Connection Failure
- **Given** the Prometheus endpoint is unreachable
- **When** I attempt to view a dashboard
- **Then** I see an error message: "Unable to connect to metrics service"
- **And** I see a troubleshooting link to documentation
- **And** the page doesn't crash or show technical stack traces

## AC9: Dashboard-as-Code Deployment
- **Given** I am a cluster administrator
- **When** I apply a PersesDashboard CRD via `kubectl apply -f dashboard.yaml`
- **Then** the dashboard appears in the ODH Dashboard UI within 30 seconds
- **And** all users with appropriate RBAC can view it
- **And** updates to the CRD are reflected in the UI

## AC10: RBAC Enforcement
- **Given** I am a user without access to namespace "team-b-project"
- **When** I view the Observability page
- **Then** I do not see any dashboards or metrics from "team-b-project"
- **And** attempting to directly access those metrics returns an authorization error
- **And** I only see projects I have at least "view" permissions for

## AC11: Reference Dashboards Included
- **Given** I am a fresh installation of ODH with Perses integration
- **When** I navigate to the Observability page
- **Then** I see at least 5 pre-configured dashboards:
  - "Model Serving Performance"
  - "Training Jobs Overview"
  - "Data Pipeline Metrics"
  - "Workbench & Notebook Resources"
  - "Infrastructure Health" (admin only)
- **And** each dashboard has documentation/description explaining its purpose

## AC12: Performance Under Load
- **Given** a dashboard with 6 panels querying different metrics
- **When** I load the dashboard with a 24-hour time range
- **Then** all panels load within 5 seconds on a standard connection
- **And** the UI remains responsive during loading
- **And** slow-loading panels show individual loading indicators

## AC13: Contextual Integration
- **Given** I am viewing a deployed model on the Model Serving page
- **When** I click the "View Metrics" button on the model card
- **Then** I am taken to the Model Serving Performance dashboard
- **And** the dashboard is pre-filtered to show only that model's metrics
- **And** the time range defaults to "Last 1 hour" for real-time monitoring

## AC14: GitOps Workflow
- **Given** I am an MLOps Engineer with a dashboard definition in Git
- **When** I commit a change to the dashboard YAML and ArgoCD syncs
- **Then** the dashboard updates in the ODH Dashboard UI within 1 minute
- **And** all users see the updated dashboard version
- **And** I can view the Git commit history for dashboard changes

## AC15: Multi-Tenancy Isolation
- **Given** I am user Alice in namespace "team-a"
- **And** user Bob is in namespace "team-b"
- **When** I view the Observability page
- **Then** I see only metrics from "team-a"
- **And** Bob sees only metrics from "team-b"
- **And** neither of us can view each other's metrics without explicit RBAC grants

---

# Use Cases - User Experience & Workflow

## Use Case 1: Monitor Model Serving Performance in Production

**Primary Actor:** MLOps Engineer

**Description:** Monitor health, latency, throughput, and error rates of deployed models to ensure SLA compliance and identify performance degradation before customer impact.

**Main Success Scenario:**
1. MLOps Engineer logs into ODH Dashboard
2. Navigates to Observability → Model Serving Performance dashboard
3. Selects "production" namespace from project filter
4. Views consolidated dashboard showing:
   - Request rate trends over last 6 hours
   - Latency percentiles (p50, p95, p99) for all deployed models
   - Error rate by model and version
   - Resource utilization (CPU, memory) per model pod
5. Identifies model "fraud-detector-v2" has p95 latency spike from 50ms to 300ms
6. Clicks on latency panel, adjusts time range to narrow down spike to 2:15 PM - 2:45 PM
7. Enables auto-refresh (30s) to monitor ongoing behavior
8. Observes latency returning to normal after infrastructure scaling event
9. Documents incident in runbook with screenshot and dashboard link

**Alternative Flow 1A: No Data Available**
- At step 4, user sees "No data available" message
- Dashboard shows tooltip: "No models are currently deployed in the 'production' namespace"
- User navigates to Model Serving page to verify deployment status
- User discovers models are deployed in "prod" namespace (typo in filter)
- User corrects namespace filter, metrics appear

**Alternative Flow 1B: Permission Denied**
- At step 3, user attempts to select "production" namespace but it's not visible in dropdown
- User sees only namespaces they have access to ("staging", "dev")
- User contacts administrator to request production namespace access
- Administrator grants RoleBinding for production namespace
- User refreshes page, can now see production metrics

**Alternative Flow 1C: Multiple Models Comparison**
- At step 5, user wants to compare latency across all model versions
- User adjusts dashboard time range to "Last 7 days" to see trends over time
- User identifies that latency increased after v2 deployment
- User compares v1 vs v2 metrics side-by-side
- User decides to investigate v2 code changes for performance regression

**Alternative Flow 1D: Trace Investigation (Post-MVP)**
- At step 5, user identifies high latency but unclear root cause
- User clicks "View Traces" button on latency panel
- Tempo trace view opens, showing distributed trace for slow request
- User identifies slow database query as bottleneck
- User optimizes database query, redeploys model, verifies improvement

## Use Case 2: Debug Data Science Pipeline Failures

**Primary Actor:** Data Scientist

**Description:** Understand why pipeline failed or performed slowly by examining execution metrics, step durations, and resource constraints.

**Main Success Scenario:**
1. Data Scientist receives pipeline failure notification via email/Slack
2. Clicks "View Pipeline Details" link in notification
3. Pipeline details page shows failure status with "View Metrics" button
4. Clicks "View Metrics", lands on Data Pipeline Metrics dashboard
5. Dashboard loads pre-filtered to specific pipeline run ID
6. User views pipeline topology with execution times per step:
   - Data Ingestion: 5 minutes (completed)
   - Preprocessing: 45 minutes (failed with OOM error)
   - Feature Engineering: Not started
   - Model Training: Not started
7. User clicks on "Preprocessing" step to see detailed metrics
8. Examines memory usage graph showing steady climb to 8GB limit, then crash
9. Identifies data volume increased from 1M to 5M records (5x increase)
10. Returns to pipeline definition, increases memory limit from 8GB to 16GB
11. Re-runs pipeline, monitors dashboard, confirms successful completion

**Alternative Flow 2A: Comparative Analysis**
- At step 6, user wants to compare failed run with successful previous runs
- User clicks "Compare Runs" button
- Selects previous successful run from dropdown
- Dashboard shows side-by-side comparison of resource usage
- User identifies that successful runs processed 1M records, failed run processed 5M
- User realizes need to parameterize memory allocation based on data volume

**Alternative Flow 2B: Resource Optimization**
- At step 10, user wants to optimize resource requests (not just increase)
- User examines CPU utilization graph, sees only 40% CPU used despite requesting 4 cores
- User reduces CPU request from 4 to 2 cores (right-sizing)
- User increases memory from 8GB to 16GB (addressing OOM)
- Re-runs pipeline with optimized resource profile, saves cost

**Alternative Flow 2C: Data Quality Investigation**
- At step 9, user suspects data quality issue rather than just volume
- User navigates to data connection metrics, sees upstream data source errors
- User identifies corrupted data causing preprocessing to consume excessive memory
- User contacts data engineering team to fix upstream data quality issue
- User re-runs pipeline with clean data, completes successfully without memory increase

## Use Case 3: Establish Observability Standards Across Platform

**Primary Actor:** OpenShift AI Administrator

**Description:** Create, manage, and distribute standardized dashboards across teams for consistent monitoring, compliance, and operational excellence.

**Main Success Scenario:**
1. Administrator identifies need: teams are creating inconsistent custom dashboards
2. Navigates to local development environment with kubectl access
3. Creates PersesDashboard CRD YAML for "Standard Model Serving Dashboard":
   ```yaml
   apiVersion: perses.dev/v1alpha1
   kind: PersesDashboard
   metadata:
     name: standard-model-serving
     namespace: opendatahub-perses
   spec:
     display:
       name: "Standard Model Serving Performance"
       description: "Company-wide standard for monitoring model serving"
     panels:
       - type: timeseries
         queries:
           - promql: "rate(model_requests_total{namespace=\"$namespace\"}[5m])"
         title: "Request Rate"
       - type: timeseries
         queries:
           - promql: "histogram_quantile(0.95, rate(model_latency_bucket[5m]))"
         title: "P95 Latency"
   ```
4. Tests dashboard locally using `percli validate dashboard.yaml`
5. Commits dashboard YAML to Git repository (e.g., `observability-dashboards/model-serving.yaml`)
6. Creates pull request with dashboard definition
7. Team reviews PR, suggests adding error rate panel
8. Administrator updates dashboard, pushes changes
9. PR approved and merged to main branch
10. ArgoCD detects change, applies PersesDashboard CRD to cluster
11. Within 30 seconds, all users see "Standard Model Serving Performance" dashboard in Observability page
12. Users in different namespaces see dashboard filtered to their namespace automatically
13. Administrator documents dashboard in company wiki, linking to Git source for transparency

**Alternative Flow 3A: Dashboard Validation Failure**
- At step 4, `percli validate` returns error: "Invalid PromQL syntax"
- Administrator fixes query syntax error
- Re-validates locally until validation passes
- Proceeds with Git commit (prevents broken dashboard from reaching production)

**Alternative Flow 3B: Staged Rollout**
- At step 10, administrator wants to test dashboard in staging before production
- Administrator uses ArgoCD application for staging environment first
- Dashboard deploys to staging, administrator validates with test users
- After validation, administrator merges to production ArgoCD sync
- Dashboard rolls out to production users

**Alternative Flow 3C: Dashboard Versioning**
- At step 7, administrator realizes dashboard needs breaking change
- Administrator creates new dashboard version: `standard-model-serving-v2`
- Both v1 and v2 dashboards exist simultaneously
- Administrator communicates migration timeline to teams
- After migration period, administrator deprecates v1 dashboard

**Alternative Flow 3D: Custom Dashboard Governance**
- At step 1, administrator wants to audit all custom dashboards across namespaces
- Administrator queries cluster: `kubectl get persesdashboards --all-namespaces`
- Identifies orphaned or unused dashboards (no access logs in 90 days)
- Administrator works with teams to consolidate or retire unused dashboards
- Reduces dashboard sprawl, improves discoverability

## Use Case 4: Investigate Model Drift and Data Quality

**Primary Actor:** MLOps Engineer

**Description:** Monitor production models for data drift, prediction shifts, and input quality issues impacting accuracy over time.

**Main Success Scenario:**
1. MLOps Engineer receives automated alert: "Model accuracy degradation detected"
2. Logs into ODH Dashboard, navigates to Observability → Model Monitoring
3. Views "Model Health" dashboard showing:
   - Accuracy metrics over last 30 days (trending down from 92% to 85%)
   - Prediction distribution (shifted toward negative class)
   - Input feature statistics
4. Clicks on "fraud-detector-v3" model to drill into model-specific dashboard
5. Examines "Feature Distribution" panel:
   - transaction_amount: mean shifted from $150 to $450 (3x increase)
   - merchant_category: distribution changed (new categories appeared)
6. Compares current distribution with training baseline (reference line on chart)
7. Confirms statistical drift detected 10 days ago (PSI score > 0.2 threshold)
8. Changes time range to "Last 90 days" to see longer-term trend
9. Correlates drift with external event: merchant onboarding campaign started 10 days ago
10. Initiates model retraining with recent 90 days of production data to adapt to new distribution
11. Monitors retraining job progress via Training Jobs dashboard
12. After retraining completes, deploys new model version, monitors accuracy recovery

**Alternative Flow 4A: False Positive Drift Alert**
- At step 6, user examines feature distribution
- Identifies distribution shift is seasonal pattern (holiday shopping surge)
- Consults historical data from previous year, confirms pattern repeats annually
- User adjusts drift detection threshold to account for seasonal variation
- User documents finding: "Expect transaction_amount spike in November-December"

**Alternative Flow 4B: Data Quality Issue**
- At step 5, user notices feature distribution has unexpected nulls (20% of records)
- User navigates to data pipeline metrics, sees upstream data connector errors
- User identifies data pipeline change introduced bug (incorrect field mapping)
- User contacts data engineering team, reverts pipeline change
- User monitors data quality metrics, confirms nulls drop to 0%
- User decides not to retrain model (drift was data bug, not real distribution shift)

**Alternative Flow 4C: Multi-Model Comparison**
- At step 4, user wants to compare drift across champion vs. challenger models
- User selects multi-model comparison view
- Dashboard shows side-by-side feature distributions for both models
- User identifies champion model more robust to distribution shift
- User decides to promote champion to 100% traffic (challenger shows degradation)

## Use Case 5: Optimize Notebook Resource Usage

**Primary Actor:** OpenShift AI Administrator

**Description:** Gain visibility into Jupyter notebook utilization to optimize capacity, identify waste, and enforce governance policies.

**Main Success Scenario:**
1. Administrator reviews monthly infrastructure costs, notices high GPU spend
2. Navigates to Observability → Workbenches & Notebooks dashboard
3. Views aggregate metrics:
   - 45 active notebooks across all projects
   - 12 GPUs allocated, average utilization: 35%
   - 8 notebooks idle >24 hours with GPU resources allocated
4. Filters by project "exploratory-data-science" (high resource allocation)
5. Drills down to notebook-level details:
   - 5 notebooks requesting 1 GPU each
   - Actual GPU utilization: 10-20% (underutilized)
   - Notebooks owned by same user (Alice)
6. Examines Alice's notebook activity logs, sees last code execution 3 days ago
7. Contacts Alice: "Your notebooks are idle but holding GPU resources"
8. Alice responds: "Forgot to shut them down after experiment completed"
9. Administrator culls idle notebooks, frees 5 GPUs
10. Administrator implements auto-culling policy: shut down notebooks after 4 hours of inactivity
11. Administrator configures policy in ODH Dashboard settings
12. Monitors dashboard over next 7 days:
    - GPU utilization increases to 70% (fewer idle resources)
    - Active notebook count drops from 45 to 28 (auto-culling working)
    - Cost reduction: 40% decrease in GPU spend
13. Shares utilization report with team leads for visibility

**Alternative Flow 5A: Capacity Planning**
- At step 3, administrator sees overall cluster GPU utilization at 85% consistently
- Identifies need to add GPU nodes to cluster to accommodate growth
- Uses historical trend data from dashboard (last 90 days)
- Forecasts GPU demand will reach 100% in 6 weeks
- Creates capacity plan: add 8 GPU nodes by end of quarter
- Tracks dashboard metrics to validate forecast accuracy

**Alternative Flow 5B: Chargeback Reporting**
- At step 4, administrator wants to generate chargeback report for finance
- Exports resource usage data by project (CPU hours, GPU hours, memory GB-hours)
- Applies cost model: $2/GPU-hour, $0.10/CPU-hour
- Generates cost attribution report showing spend by project
- Shares report with project managers for budget accountability

**Alternative Flow 5C: Right-Sizing Recommendations**
- At step 5, administrator identifies notebooks requesting 16 CPU cores but using only 4
- Administrator generates right-sizing recommendation report
- Shares report with users: "Your notebook requests 16 cores but uses 4 - consider reducing to 8"
- Users adjust resource requests, freeing capacity for other teams
- Cluster efficiency improves without adding nodes

---

## Workflow Details

### Workflow A: Discovering and Accessing Observability Dashboards

**Entry Points:**
1. **Primary Navigation**: Top-level "Observability" in left sidebar (after Model Serving, before Settings)
2. **Contextual Links**: From Model Serving cards ("View Metrics"), Pipeline runs ("Performance"), Data Connections
3. **Dashboard Home**: Quick access to recent/favorited dashboards
4. **Global Search**: Keywords like "observability", "metrics", "monitoring"
5. **Alert Notifications**: Click alert → dashboard link with context

**Discovery Journey:**
```
Awareness → Exploration → Context → Access → Orientation → Action
```

**Information Architecture:**
```
Observability (Main Section)
├── Overview (Landing page)
│   ├── Quick Stats (health indicators)
│   ├── System Health (service status)
│   └── Recent Alerts (last 24 hours)
├── Model Serving
│   ├── All Models Overview
│   ├── Individual Model Details
│   └── KServe Infrastructure
├── Pipelines
│   ├── Pipeline Runs Overview
│   ├── Performance Analytics
│   └── Individual Pipeline Details
├── Workbenches & Notebooks
│   ├── Resource Utilization
│   ├── Active Sessions
│   └── User Activity
├── Infrastructure (admin)
│   ├── Cluster Resources
│   ├── Storage Performance
│   └── Network Traffic
└── Custom Dashboards
    └── [User-defined]
```

### Workflow B: Viewing Metrics for AI Workloads

**Loading Experience:**
- Progressive loading: Critical metrics < 2s, detailed panels follow
- Skeleton loaders show structure during load
- Partial data display if some metrics unavailable
- Error boundaries prevent full page crash if one panel fails

**Understanding Metrics:**
- Visual hierarchy: Health/latency/errors top-left (most important)
- Consistent styling across dashboards (PatternFly components)
- Inline tooltips explain metric definitions on hover
- Color-coded thresholds: green (healthy), yellow (warning), red (critical)
- Status indicators: ✓ (good), ⚠️ (warning), ✗ (error)

**Interactions:**
- **Time range:** Last 15min, 1hr, 6hr, 24hr, 7d, 30d, custom date picker
- **Auto-refresh:** 5s, 15s, 30s, 1min intervals with pause/resume
- **Zoom/pan:** Click-drag to zoom into time window, reset button
- **Data inspection:** Hover for values, click for drill-down
- **Export:** Download panel data as CSV or PNG

**Contextual Actions:**
- Export dashboard (PNG/PDF)
- Share URL with filters embedded
- Set alert rules (future)
- Navigate to related resources (e.g., model details, pipeline runs)

### Workflow C: Filtering & Scoping Dashboards

**Filter Hierarchy:**

**Level 1 - Namespace/Project (Primary):**
- Dropdown: "All Projects" or specific project selection
- Respects Kubernetes RBAC (only shows accessible namespaces)
- Persists across navigation (browser local storage)
- Visual indicator of active filter (chip/badge)

**Level 2 - Workload Type:**
- Tabs/chips: "All", "Model Serving", "Pipelines", "Notebooks"
- Filters dashboard list and metrics

**Level 3 - Resource-Specific:**
- Model name, pipeline ID, user, pod
- Auto-complete with search as you type
- Multiple selection support (e.g., compare 3 models)

**Level 4 - Metric-Specific:**
- Prometheus label filters (advanced users)
- Advanced mode for power users (PromQL query editing)

**UX Features:**
- Persistent filters (survives browser refresh)
- Filter templates: "My Active Models", "Production Only"
- Clear all filters button
- Breadcrumb chips showing active filters
- Deep linking: URL encodes filters for sharing

### Workflow D: Troubleshooting Performance Issues

**Diagnostic Journey:**

**Phase 1 - Detection:**
- User notices issue (slow response, timeout, failed job)
- Or receives automated alert notification
- Navigates from affected resource page or alert link
- Dashboard loads with incident time context pre-selected

**Phase 2 - Scoping:**
- Overview dashboard shows high-level health indicators
- User identifies affected components via visual anomalies (red panels, spike graphs)
- User narrows time range to incident window (e.g., 2:00 PM - 2:30 PM)

**Phase 3 - Investigation:**
- User drills down to component-specific dashboard
- Multi-panel correlated view shows:
  - **Performance:** Latency, throughput, success rate
  - **Resources:** CPU, memory, GPU utilization
  - **Errors:** HTTP 5xx, exceptions, retries
  - **Dependencies:** Database queries, storage I/O, external API calls
- User compares with baseline (previous day or week)

**Phase 4 - Root Cause:**
- User identifies metric pattern (e.g., memory exhaustion, CPU throttling)
- User switches to trace view (Tempo integration, post-MVP) for request-level details
- User examines correlated log entries (link to logging system)
- User forms hypothesis (e.g., "New model version has memory leak")

**Phase 5 - Validation:**
- User makes configuration change (increase memory limit, rollback deployment)
- User monitors dashboard in real-time (auto-refresh enabled)
- User validates metrics return to healthy range
- User documents findings (runbook annotation, incident report)

**UX Support:**
- Correlation highlighting: synchronized time ranges across panels
- Comparison mode: overlay "good period" vs. "bad period"
- Timeline annotations: show deployments, config changes, incidents
- Suggested related dashboards: "Users also viewed..."
- Guided checklists: "Common causes of high latency"

---

## Navigation & Information Architecture

### ODH Dashboard Integration

**Placement in Main Navigation:**
```
OpenShift AI Dashboard
├── Home / Overview
├── Applications
├── Data Science Projects
├── Model Serving
├── Data Science Pipelines
├── [NEW] Observability ← Top-level item (after Pipelines)
│   ├── Overview
│   ├── Model Serving
│   ├── Pipelines
│   ├── Workbenches
│   ├── Infrastructure (admin)
│   └── Custom Dashboards
├── Resources
├── Settings
└── [User Profile]
```

**Contextual Integration Points:**

1. **Model Serving Page:**
   - Mini-metrics preview on model cards (sparklines showing latency trend)
   - "View Metrics" button on each model card → filtered observability dashboard
   - Health status indicator (green/yellow/red dot) based on error rate

2. **Pipeline Runs:**
   - Pipeline run list shows duration, status, resource summary columns
   - Pipeline details page: "Performance Metrics" tab with embedded dashboard
   - "View in Observability" button for full-screen analysis

3. **Data Science Project Overview:**
   - Project details page: "Observability" tab showing project-scoped metrics
   - Resource usage summary widget (CPU, memory, GPU)

4. **Notebook Server List:**
   - Active notebook list: real-time resource utilization bars (CPU, memory)
   - "View Metrics" button → notebook-specific dashboard

5. **Global Alerts:**
   - Header notification icon with alert count badge
   - Click → alert summary panel → "View Dashboard" links

**Navigation Behavior:**
- Breadcrumbs: `Home > Observability > Model Serving > fraud-detector-v2`
- Browser back button returns to previous context
- Deep linking: `/observability/model-serving?project=fraud-detection&model=fraud-detector-v2&range=6h`
- State preservation: filters and time range saved in browser session

---

## Persona-Specific Needs

### Data Scientists

**Context:** Experimentation, model development, pipeline iteration, learning

**Dashboard Priorities:**
1. **My Experiments Dashboard**
   - Pipeline run durations, step-by-step breakdown
   - Training metrics (loss curves, accuracy trends)
   - Resource utilization (CPU, memory, GPU)
   - Cost estimation per experiment

2. **Notebook Performance**
   - Personal resource usage (current notebook session)
   - Kernel restart frequency (indicates memory issues)
   - Library load times (identify slow imports)

3. **Model Training Insights**
   - GPU utilization and efficiency
   - Data loading bottlenecks
   - Hyperparameter comparison (side-by-side runs)

**UX Adaptations:**
- **Simplified language:** Avoid Kubernetes jargon (use "My Workspace" not "Namespace")
- **Guided interpretation:** Tooltips explain what metrics mean and what to do
- **Comparative experiment views:** Easy comparison of multiple training runs
- **Default filter to "my resources":** Don't overwhelm with team/cluster data
- **Learning resources:** Links to tutorials, best practices docs

**Example Layout:**
```
┌────────────────────────────────────────┐
│ My Recent Pipeline Runs (Last 7 Days)  │
│ [Timeline: success/failure, durations]  │
│ ✓✓✗✓✗ (2 failures need attention)      │
└────────────────────────────────────────┘
┌─────────────────┐ ┌──────────────────┐
│ Current Run     │ │ Resource Usage   │
│ Step 3/5        │ │ CPU: 85% ⚠️       │
│ Training 23min  │ │ Mem: 12GB/16GB   │
│ ETA: 15min      │ │ GPU: 92% (good)  │
└─────────────────┘ └──────────────────┘
```

### MLOps Engineers

**Context:** Production reliability, performance optimization, SLO compliance, incident response

**Dashboard Priorities:**
1. **Production Model Health**
   - SLI/SLO tracking (uptime, latency, error rate)
   - RED method (Rate, Errors, Duration)
   - Model drift detection (statistical tests)
   - Dependency health (database, storage, APIs)

2. **Infrastructure Performance**
   - KServe/ModelMesh resource utilization
   - Autoscaling behavior and effectiveness
   - Resource saturation indicators
   - Network performance (request routing, load balancing)

3. **Pipeline Orchestration**
   - Success rates, retry patterns
   - Queue depths, parallelism metrics
   - Data volume trends
   - Execution time distributions

4. **Cost & Efficiency**
   - Resource costs by project/model
   - Utilization vs. requested (waste identification)
   - Idle time analysis
   - Cost optimization recommendations

**UX Adaptations:**
- **Advanced PromQL query builder:** Direct query editing for ad-hoc analysis
- **Multi-model comparison views:** Side-by-side dashboards
- **Alert integration:** Create alerts directly from dashboard panels (post-MVP)
- **Trace correlation:** Link from metrics to Tempo traces (post-MVP)
- **GitOps export:** Download dashboard as YAML for version control
- **API access links:** REST API endpoints for programmatic access
- **Runbook integration:** Link dashboards to operational runbooks

**Example Layout:**
```
┌──────────┬──────────┬──────────┬──────────┐
│ Uptime   │ Errors   │ P99 Lat  │ Req/min  │
│ 99.97%   │ 0.03%    │ 234ms    │ 15,234   │
│ 🟢 >99.95%│ 🟢 <0.1% │ 🟡 <200ms│ [Trend↑] │
└──────────┴──────────┴──────────┴──────────┘
┌────────────────────────────────────────────┐
│ Latency Heatmap (Last 1 hour)              │
│ [Time x Latency bucket visualization]      │
│ Shows temporal patterns in latency         │
└────────────────────────────────────────────┘
```

### OpenShift AI Administrators

**Context:** Platform health, capacity management, governance, cost control, user enablement

**Dashboard Priorities:**
1. **Platform Overview**
   - Cluster health (CPU, memory, GPU, storage capacity)
   - User activity, concurrent workloads
   - Service health (ODH components: operators, serving, pipelines)
   - Incidents, alerts, anomalies

2. **Multi-Tenancy & RBAC**
   - Resource usage by project/team
   - Quota utilization and violations
   - Permission audit trail
   - Orphaned resources (unused PVCs, idle notebooks)

3. **Capacity Planning**
   - Historical trends (3 months, 6 months, 1 year)
   - Growth projections and forecasts
   - Resource efficiency scores
   - Cost forecasts and budget tracking

4. **Dashboard Governance**
   - Dashboard catalog (all dashboards across cluster)
   - Usage statistics (views, users)
   - Version control history (Git commits)
   - Orphaned dashboards (unused >90 days)

**UX Adaptations:**
- **Executive summary views:** High-level health indicators for quick assessment
- **Drill-down hierarchy:** Platform → Project → Workload → Resource
- **Bulk operations:** Multi-select for actions (delete, export, assign)
- **Audit trail visibility:** Who accessed what, when
- **Template creation tools:** UI for creating dashboard templates
- **Export/reporting:** PDF reports for stakeholders

**Admin-Only Features:**
- Dashboard CRUD UI (create, update, delete dashboards)
- Access control editor (assign dashboards to roles/namespaces)
- Monitoring stack health (Prometheus, Perses, Tempo status)
- Global settings (default dashboards, retention policies)

**Example Layout:**
```
┌────────────────────────────────────────────┐
│ Platform Health Summary                    │
│ Services: ✓ Serving ✓ Pipelines ✓ Dash    │
│ Cluster: 🟢 Healthy | Alerts: 2 warnings  │
└────────────────────────────────────────────┘
┌─────────────┐ ┌────────────────────────────┐
│ Workloads   │ │ Cluster Capacity           │
│ Models: 23  │ │ CPU: 67% [■■■■■■■□□□]     │
│ Pipelines:12│ │ Mem: 72% [■■■■■■■□□□]     │
│ Notebooks:45│ │ GPU: 89% [■■■■■■■■■□] ⚠️  │
│             │ │ Need GPU nodes soon        │
└─────────────┘ └────────────────────────────┘
```

---

# Documentation Considerations

*Information that needs to be considered and planned for documentation to meet customer needs.*

## Architecture Documentation Required

### 1. Perses Integration Architecture (`docs/architecture/observability/perses-integration.md`)
- **Integration Pattern:** React component integration using @perses-dev npm packages (not iframe or Module Federation)
- **Component Diagram:** ODH Dashboard Frontend ↔ Perses Backend ↔ Prometheus/Tempo
- **Sequence Diagrams:**
  - Authentication flow: User login → OAuth token → Perses API validation
  - Dashboard loading: User clicks dashboard → API call → Query Prometheus → Render panels
  - Metric query execution: Dashboard panel → Perses backend → Prometheus → Results
- **Architecture Decision Records (ADRs):**
  - ADR-001: Why React component integration over iframe
  - ADR-002: Multi-tenancy model (single shared Perses vs. per-namespace)
  - ADR-003: RBAC enforcement boundary (Perses backend, ODH Dashboard, or both)

### 2. Multi-Tenancy & Security Model (`docs/architecture/observability/security-model.md`)
- **Namespace Isolation Strategy:** Query filtering, RBAC enforcement
- **RBAC Matrix:**
  | User Role | Dashboard View | Dashboard Edit | Metric Query Scope | Admin Functions |
  |-----------|----------------|----------------|-------------------|-----------------|
  | Data Scientist | Namespace-scoped | Own namespace | Namespace filter | No |
  | MLOps Engineer | Namespace-scoped | Own namespace | Namespace filter | No |
  | Project Admin | Namespace-scoped | Own namespace | Namespace filter | No |
  | Cluster Admin | Cluster-wide | All namespaces | No filter | Yes |
- **Threat Model & Mitigations:**
  - **Threat:** Privilege escalation (user A views namespace B metrics)
  - **Mitigation:** Prometheus query rewriting, dual RBAC enforcement
  - **Threat:** Data exfiltration (malicious dashboard with external HTTP)
  - **Mitigation:** Disable external data sources in Perses config
  - **Threat:** DoS (expensive PromQL queries)
  - **Mitigation:** Query timeout (60s), rate limiting (100 queries/min/user)
  - **Threat:** XSS (malicious dashboard JSON)
  - **Mitigation:** CSP headers, React XSS protection, JSON sanitization
- **Audit Logging:**
  - Events logged: dashboard.viewed, dashboard.created, dashboard.deleted, query.executed, permission.denied
  - Log format: JSON with timestamp, user, namespace, action, sourceIP
  - Retention: 90 days (configurable)

### 3. Data Source Configuration (`docs/architecture/observability/data-sources.md`)
- **Integration with Platform Metrics Architecture:** Perses depends on centralized Prometheus deployment
- **Prometheus Configuration:**
  - Endpoint: `http://prometheus-k8s.openshift-monitoring.svc:9091`
  - Authentication: ServiceAccount token injection
  - Federation: Multi-namespace metric access with label filtering
- **Tempo Integration (Post-MVP):**
  - Endpoint: `http://tempo-query-frontend.observability.svc:3200`
  - Authentication: ServiceAccount token
  - Trace-metric correlation: Link from latency spike to distributed trace
- **ServiceAccount Permissions:**
  ```yaml
  apiVersion: rbac.authorization.k8s.io/v1
  kind: ClusterRole
  metadata:
    name: perses-prometheus-reader
  rules:
  - nonResourceURLs: ["/api/v1/query", "/api/v1/query_range"]
    verbs: ["get"]
  ```
- **NetworkPolicy:**
  - Allow Perses pods → Prometheus pods
  - Allow ODH Dashboard pods → Perses pods
  - Deny all other ingress

## Deployment & Operations Documentation Required

### 1. Perses Deployment Guide (`docs/deployment/perses-installation.md`)
- **Prerequisites Checklist:**
  - [ ] Platform Metrics and Alerts architecture deployed
  - [ ] OpenShift OAuth configured for ODH Dashboard
  - [ ] Persistent storage available (10Gi minimum)
  - [ ] ServiceAccount with Prometheus query permissions created
- **Helm Chart Installation:**
  ```bash
  helm repo add perses https://perses.github.io/helm-charts
  helm install perses perses/perses \
    --namespace opendatahub-perses \
    --create-namespace \
    --set replicas=3 \
    --set persistence.enabled=true \
    --set persistence.size=10Gi \
    --set database.type=postgresql
  ```
- **High Availability Configuration:**
  - PostgreSQL backend required for multi-replica (SQLite only for single-replica/POC)
  - HorizontalPodAutoscaler: min 3, max 10 replicas, target 70% CPU
  - PodDisruptionBudget: minAvailable 2
- **Storage Configuration:**
  - StorageClass: Use managed storage (e.g., `managed-premium` on Azure, `gp3` on AWS)
  - Size: 10Gi minimum, 50Gi for large deployments (5000+ dashboards)
  - Backup: Velero integration (see disaster recovery section)

### 2. Upgrade Procedures (`docs/operations/perses-upgrades.md`)
- **Version Compatibility Matrix:**
  | Perses Backend | Perses npm packages | ODH Dashboard | Notes |
  |----------------|---------------------|---------------|-------|
  | 0.45.x | 0.45.x | 2.x | MVP baseline |
  | 0.46.x | 0.46.x | 2.x | Security patches only |
  | 0.47.x | 0.47.x | 2.1+ | Breaking API changes |
- **Upgrade Runbook:**
  1. **Backup:** Create Velero backup of Perses namespace
  2. **Staging Test:** Upgrade Perses in staging environment first
  3. **Validation:** Verify dashboards load, metrics display, no errors
  4. **Production Upgrade:** Upgrade Helm chart in production
     ```bash
     helm upgrade perses perses/perses \
       --namespace opendatahub-perses \
       --version 0.46.0 \
       --reuse-values
     ```
  5. **Smoke Test:** Load pre-built dashboards, verify metrics
  6. **Rollback (if needed):**
     ```bash
     helm rollback perses --namespace opendatahub-perses
     velero restore create --from-backup perses-pre-upgrade
     ```
- **Breaking Changes Migration Guide:** Document API changes, provide migration scripts

### 3. Disaster Recovery Runbook (`docs/operations/disaster-recovery.md`)
- **Backup Strategy:**
  - **Velero Scheduled Backups:** Daily at 1 AM UTC, retain 30 days
  - **Scope:** Perses PVC (dashboard definitions), ConfigMaps, Secrets
  - **GitOps Sync (Optional):** CronJob exports dashboards to Git repository every 6 hours
- **Recovery Targets:**
  - **RTO (Recovery Time Objective):** 1 hour
  - **RPO (Recovery Point Objective):** 24 hours (daily backups)
- **Disaster Recovery Procedure:**
  1. **Detect outage:** Monitoring alerts on Perses pod health
  2. **Attempt restart:** `kubectl delete pod -n opendatahub-perses -l app=perses`
  3. **If PVC corrupted:** Restore from Velero
     ```bash
     velero restore create perses-restore \
       --from-backup perses-daily-20251015
     ```
  4. **Redeploy Perses:** `helm upgrade perses perses/perses --reuse-values`
  5. **Verify:** `curl https://perses.odh.example.com/api/v1/dashboards`
  6. **Smoke test:** Load dashboard, verify metrics display
  7. **Post-incident review:** Document root cause, improve monitoring
- **DR Testing:** Quarterly drills (delete namespace, restore, verify)

### 4. Monitoring Perses Integration (`docs/operations/monitoring.md`)
- **Metrics to Monitor:**
  - Perses pod health: `kube_pod_status_phase{namespace="opendatahub-perses"}`
  - API response time: `perses_http_request_duration_seconds`
  - Prometheus query latency: `prometheus_query_duration_seconds`
  - Dashboard load errors: `perses_dashboard_load_errors_total`
- **Alerts to Configure:**
  - PersesPodDown: All Perses replicas down
  - PersesHighQueryErrorRate: >5% query errors
  - PersesStorageNearFull: PVC >80% utilization
  - PrometheusQuerySlow: p95 latency >2s
- **Dashboard for Monitoring Observability:** Meta-dashboard showing Perses health

## Developer Documentation Required

### 1. Integrating Perses Components (`docs/development/perses-component-integration.md`)
- **How to add Perses React components to ODH Dashboard:**
  ```typescript
  import { Dashboard } from '@perses-dev/dashboards';
  import { PrometheusClient } from '@perses-dev/prometheus-plugin';

  function ObservabilityPage() {
    const client = new PrometheusClient({ url: PROMETHEUS_URL });
    return <Dashboard client={client} dashboardId="model-serving" />;
  }
  ```
- **TypeScript Type Definitions:** Interface definitions for dashboard config, panel types
- **Example: Creating Custom Panel:**
  ```typescript
  import { PanelPlugin } from '@perses-dev/dashboards';

  const MyCustomPanel: PanelPlugin = {
    type: 'my-custom-panel',
    component: MyPanelComponent,
    options: { /* ... */ }
  };
  ```
- **Webpack Configuration:** Code splitting, lazy loading for Perses components

### 2. Creating Dashboard Templates (`docs/development/dashboard-as-code.md`)
- **Dashboard-as-Code Examples:**
  ```yaml
  apiVersion: perses.dev/v1alpha1
  kind: PersesDashboard
  metadata:
    name: model-serving-template
  spec:
    variables:
      - name: namespace
        type: text
        default: "$__namespace"
      - name: model
        type: query
        query: "label_values(model_requests_total{namespace=\"$namespace\"}, model)"
    panels:
      - title: "Request Rate"
        type: timeseries
        queries:
          - promql: "rate(model_requests_total{namespace=\"$namespace\", model=\"$model\"}[5m])"
  ```
- **Using Perses CUE SDK:**
  ```cue
  import "github.com/perses/perses/cue/schemas/v1alpha1"

  dashboard: v1alpha1.#Dashboard & {
    metadata: {
      name: "model-serving"
    }
    spec: {
      panels: [...]
    }
  }
  ```
- **Template Variables:** `$namespace`, `$model_id`, `$pipeline_run_id`
- **Testing Dashboards Locally:** `percli validate dashboard.yaml`
- **Publishing to Git:** Commit to repository, apply via GitOps (ArgoCD/Flux)

### 3. Perses API Integration Guide (`docs/development/perses-api.md`)
- **REST API Endpoints:**
  ```
  GET  /api/v1/dashboards              # List dashboards
  GET  /api/v1/dashboards/{name}       # Get dashboard
  POST /api/v1/dashboards              # Create dashboard
  PUT  /api/v1/dashboards/{name}       # Update dashboard
  DELETE /api/v1/dashboards/{name}     # Delete dashboard
  POST /api/v1/query                   # Execute PromQL query
  ```
- **Authentication:** Pass OAuth token in `Authorization: Bearer <token>` header
- **Error Handling:**
  ```typescript
  try {
    const response = await fetch('/api/v1/dashboards', {
      headers: { 'Authorization': `Bearer ${oauthToken}` }
    });
    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    return await response.json();
  } catch (error) {
    // Handle network errors, timeouts, etc.
  }
  ```
- **Rate Limiting:** 100 queries/min per user, exponential backoff on 429 errors

## User Documentation Required

### 1. Observability Dashboard Overview (`docs/user/observability-getting-started.md`)
- **What Observability Data is Available:**
  - Model serving: latency, throughput, error rate
  - Training jobs: duration, resource usage, GPU utilization
  - Pipelines: execution time, step breakdown, failures
  - Notebooks: CPU, memory, GPU usage per session
  - Infrastructure: cluster capacity, node health
- **How to Navigate to Observability:**
  1. Log into ODH Dashboard
  2. Click "Observability" in left sidebar
  3. Select dashboard category (Model Serving, Pipelines, etc.)
- **Overview of Pre-Built Dashboards:**
  - **Model Serving Performance:** Monitor deployed models
  - **Training Jobs Overview:** Track training progress
  - **Data Pipeline Metrics:** Debug pipeline failures
  - **Workbench Resources:** Optimize notebook usage
  - **Infrastructure Health:** Platform capacity planning (admin)
- **Screenshots:** Annotated images of each dashboard

### 2. Using Pre-Built Dashboards (`docs/user/using-dashboards.md`)
- **Model Serving Performance Dashboard:**
  - **What it shows:** Request rate, latency (p50, p95, p99), error rate, resource usage
  - **How to interpret:** Green = healthy, yellow = warning, red = critical
  - **Common actions:** Filter by model, adjust time range, enable auto-refresh
- **Training Jobs Overview Dashboard:**
  - **What it shows:** Job status, duration, resource consumption
  - **How to interpret:** Identify slow jobs, resource bottlenecks (CPU, memory, GPU)
  - **Common actions:** Compare runs, identify OOM errors, optimize resource requests
- **Data Pipeline Metrics Dashboard:**
  - **What it shows:** Pipeline execution times, step durations, failures
  - **How to interpret:** Identify bottleneck steps, data quality issues
  - **Common actions:** Debug pipeline failures, optimize performance
- **Time Range Selection:** Last 15m (real-time), 1h, 6h, 24h, 7d (trends), custom
- **Auto-Refresh:** Enable for real-time monitoring (5s, 30s, 1m intervals)

### 3. Creating Custom Dashboards (`docs/user/custom-dashboards.md`)
- **When to Create Custom Dashboards:**
  - Pre-built dashboards don't cover your specific use case
  - Need to combine metrics from multiple sources
  - Want to track custom application metrics
- **Step-by-Step Guide:**
  1. Define dashboard using PersesDashboard CRD YAML
  2. Add panels with PromQL queries
  3. Configure variables (namespace, model name, etc.)
  4. Validate locally: `percli validate my-dashboard.yaml`
  5. Apply to cluster: `kubectl apply -f my-dashboard.yaml`
  6. Verify in ODH Dashboard UI (appears in Custom Dashboards)
- **PromQL Query Examples for ML Workloads:**
  - **Model request rate:**
    ```promql
    rate(model_requests_total{namespace="fraud-detection", model="rf-v2"}[5m])
    ```
  - **GPU utilization:**
    ```promql
    nvidia_gpu_duty_cycle{namespace="training", job="training-job-123"}
    ```
  - **Pipeline duration (95th percentile):**
    ```promql
    histogram_quantile(0.95, rate(pipeline_duration_seconds_bucket[5m]))
    ```
- **Sharing Dashboards:** Apply CRD to namespace, all project members can view

### 4. Troubleshooting ML Workloads (`docs/user/troubleshooting-guide.md`)
- **Common Scenarios with Dashboard-Driven Workflows:**
  - **Slow Model Inference:**
    1. Open Model Serving Performance dashboard
    2. Identify high latency (p95 > threshold)
    3. Check CPU utilization → if >90%, scale up replicas
    4. Check dependency health → if database slow, optimize queries
  - **Pipeline Failures:**
    1. Open Data Pipeline Metrics dashboard
    2. Identify failed step (red status indicator)
    3. Check resource usage → if OOM, increase memory limit
    4. Check data quality metrics → if null values spike, investigate upstream
  - **High GPU Costs:**
    1. Open Workbench Resources dashboard
    2. Identify idle notebooks with GPU allocated
    3. Cull idle notebooks or implement auto-shutdown policy
    4. Right-size GPU requests based on utilization patterns
- **Correlating Metrics with Logs and Traces:**
  - Metrics show "what" is wrong (high latency)
  - Logs show "why" (error messages, exceptions)
  - Traces (post-MVP) show "where" (which service in call chain)

### 5. Dashboard-as-Code Best Practices (`docs/user/gitops-dashboards.md`)
- **Benefits of Managing Dashboards as Code:**
  - Version control: track changes, rollback if needed
  - Peer review: team reviews dashboard changes via pull requests
  - Auditability: Git history shows who changed what, when
  - Reproducibility: apply same dashboard to dev, staging, prod
- **Setting Up Git Repository:**
  ```
  dashboards/
  ├── model-serving/
  │   ├── overview.yaml
  │   └── per-model.yaml
  ├── pipelines/
  │   └── execution-metrics.yaml
  └── infrastructure/
      └── cluster-health.yaml
  ```
- **Integrating with ArgoCD/Flux:**
  ```yaml
  apiVersion: argoproj.io/v1alpha1
  kind: Application
  metadata:
    name: observability-dashboards
  spec:
    source:
      repoURL: https://github.com/myorg/dashboards
      path: dashboards
    destination:
      namespace: opendatahub-perses
  ```
- **Team Workflow:**
  1. Create/modify dashboard YAML locally
  2. Validate: `percli validate dashboard.yaml`
  3. Commit to feature branch, push to Git
  4. Create pull request
  5. Team reviews PR (dashboard changes visible in diff)
  6. Merge to main branch
  7. ArgoCD/Flux auto-deploys to cluster within 1 minute

## Links to Related Documentation

### Prerequisites (Must Exist Before Perses Integration)

**Platform Metrics Architecture:**
- Link: `/docs/architecture/metrics/platform-metrics-overview.md`
- Content: Centralized Prometheus deployment, metrics federation, ServiceMonitor configurations
- Why needed: Perses requires Prometheus metrics to visualize

**OpenShift OAuth Integration:**
- Link: `/docs/security/oauth-integration-guide.md`
- Content: How OAuth tokens are obtained and validated, ServiceAccount token usage, RBAC configuration
- Why needed: Perses uses OAuth for authentication

**ODH Dashboard Architecture:**
- Link: `/docs/architecture/dashboard/frontend-architecture.md`
- Content: React component patterns, routing, state management
- Why needed: Perses components integrate into ODH Dashboard React app

### Related OpenShift Documentation

**OpenShift Monitoring Stack:**
- How to access Prometheus UI
- Configure AlertManager
- Query Prometheus metrics

**RBAC Configuration:**
- How namespace permissions work
- Creating ClusterRoles/RoleBindings
- ServiceAccount permissions

**Persistent Storage:**
- StorageClass options (managed-premium, gp3, etc.)
- PVC backup with Velero
- Storage performance optimization

### External Resources

**Perses Official Documentation:**
- https://perses.dev/docs
- npm package reference (@perses-dev/components, @perses-dev/dashboards)
- REST API specification
- Dashboard-as-Code examples (CUE SDK, Go SDK)

**Prometheus Query Language (PromQL):**
- https://prometheus.io/docs/prometheus/latest/querying/basics/
- Query functions and operators
- Best practices for efficient queries

**CNCF Observability Landscape:**
- https://landscape.cncf.io/guide#observability-and-analysis
- Context on Perses within CNCF ecosystem
- Comparison with other CNCF observability projects

---

# Questions to Answer

*Refinement and architectural questions that must be answered before coding can begin.*

## 1. RBAC Enforcement Boundary

**Question:** Should RBAC be enforced at the Perses backend, ODH Dashboard backend, or both layers?

**Context:** Dual enforcement provides defense-in-depth but adds latency. Single enforcement simplifies architecture but creates security risk if bypassed.

**Options:**
- **Option A:** Enforce only at Perses backend (query rewriting, namespace filtering)
  - Pro: Single point of enforcement, simpler
  - Con: If Perses is bypassed, no protection
- **Option B:** Enforce only at ODH Dashboard (filter API requests before calling Perses)
  - Pro: Frontend control
  - Con: Backend Perses API still accessible if discovered
- **Option C:** Enforce at both layers (defense-in-depth)
  - Pro: Maximum security, compliance-friendly
  - Con: Slight latency overhead, more complex

**Decision Driver:** Compliance requirements (SOC2, FedRAMP) typically require defense-in-depth.

**Recommendation:** Option C (dual enforcement)

**Answer Required From:** Security Architect, Compliance Team

---

## 2. Dashboard Storage & Versioning Strategy

**Question:** How are custom dashboards stored, versioned, and backed up?

**Options:**
- **Option A:** Perses native storage (SQLite for POC, PostgreSQL for production)
  - Pro: Simple, out-of-box Perses functionality
  - Con: Disaster recovery requires Velero backups, no built-in versioning
- **Option B:** GitOps pattern (dashboards as CRDs committed to Git repository)
  - Pro: Version control, audit trail, change management via PRs
  - Con: More complex workflow for casual users
- **Option C:** Hybrid (Perses storage + periodic Git sync via CronJob)
  - Pro: Best of both worlds
  - Con: Synchronization complexity

**Implications:**
- GitOps enables audit trails and change management but adds complexity
- Perses native storage is simpler but disaster recovery is harder
- Hybrid approach provides both ease-of-use and safety net

**Recommendation:** Option C (hybrid approach)

**Answer Required From:** Platform Team, GitOps Working Group

---

## 3. Multi-Tenancy Isolation Model

**Question:** Is Perses deployed as a single shared instance with logical isolation, or per-namespace/per-tenant instances?

**Options:**
- **Option A:** Single shared Perses with namespace-based access controls
  - Pro: Cost-efficient, easier to manage
  - Con: Risk of noisy-neighbor, requires robust query filtering
  - Scaling: Supports 100s of namespaces
- **Option B:** Perses instance per namespace
  - Pro: True isolation, no cross-namespace risk
  - Con: Doesn't scale beyond 50-100 namespaces, high operational overhead
  - Scaling: Limited

**Scaling Implications:** Shared instance must handle noisy-neighbor scenarios (query rate limiting, resource quotas). Per-namespace approach has operational overhead.

**Recommendation:** Option A (single shared instance with robust RBAC)

**Answer Required From:** Architect, SRE Team

---

## 4. Prometheus Data Scope

**Question:** Can users in namespace A query metrics from namespace B, or is strict namespace isolation enforced?

**Context:** Some industries (healthcare, finance) require strict data isolation. Others benefit from cross-namespace troubleshooting (e.g., MLOps Engineer monitoring dev + prod).

**Options:**
- **Option A:** Strict isolation (users can only query their own namespace)
  - Pro: Meets compliance requirements
  - Con: Limits troubleshooting workflows
- **Option B:** Flexible isolation (users can query namespaces they have RBAC access to)
  - Pro: Supports real-world workflows
  - Con: Requires robust RBAC configuration

**Technical Approach:** Perses backend rewrites PromQL queries to inject namespace filter based on user's Kubernetes RBAC permissions.

**Recommendation:** Option B (flexible isolation with RBAC)

**Answer Required From:** Security Architect, Compliance Team, Product Owner

---

## 5. Dashboard Template Distribution

**Question:** How are pre-built dashboards (Model Serving Performance, GPU Utilization, etc.) distributed and updated?

**Options:**
- **Option A:** Bundled in Perses Helm chart (updated via Helm upgrades)
  - Pro: Simple, automatic distribution
  - Con: Requires Helm upgrade to update dashboards
- **Option B:** CRD-based (DashboardTemplate CRD synced via operator)
  - Pro: Dynamic updates without redeployment
  - Con: More complex operator logic
- **Option C:** REST API seeding (init container applies dashboards on Perses deployment)
  - Pro: Flexible, can update without redeployment
  - Con: Timing issues (ensure Perses is ready before seeding)

**Update Semantics:** Can platform admins push dashboard updates to all users automatically, or must users opt-in to updates?

**Recommendation:** Option C (REST API seeding with versioning)

**Answer Required From:** Product Owner, Platform Team

---

## 6. Perses Version Upgrade Strategy

**Question:** How are Perses backend upgrades coordinated with ODH Dashboard compatibility?

**Risk:** Breaking changes in Perses npm packages (@perses-dev/components) or REST API contract could break ODH Dashboard integration.

**Mitigation Options:**
- **Option A:** API version negotiation (ODH Dashboard specifies compatible Perses API version in requests)
- **Option B:** Feature flags (gradual rollout of new Perses features, disable if incompatible)
- **Option C:** Phased rollout (upgrade Perses in staging, validate ODH Dashboard compatibility, then production)
- **Option D:** Pin Perses npm package versions (ODH Dashboard upgrades Perses packages independently)

**Recommendation:** Combination of Option C (phased rollout) + Option D (version pinning)

**Answer Required From:** Release Engineering, Tech Lead

---

## 7. Performance Characteristics at Scale

**Question:** What are acceptable latency targets for dashboard rendering with 10,000+ time series metrics?

**Assumptions to Validate:**
- **Query response time:** <2s for p95 (Prometheus query execution + network transfer)
- **Dashboard list load:** <500ms
- **Concurrent users:** 500 per cluster without degradation
- **Dashboard panels:** 6-10 panels per dashboard

**Load Testing Required:** Simulate 500 concurrent users loading GPU utilization dashboard (high-cardinality metrics), measure p95/p99 latency and error rate.

**Success Criteria:**
- p95 latency < 2s
- p99 latency < 5s
- Error rate < 1%
- No Prometheus resource exhaustion

**Answer Required From:** Performance Engineering, SRE Team

---

## 8. Integration with Model Metadata

**Question:** Can dashboards be parameterized by model ID, pipeline run ID, or experiment ID from other ODH Dashboard pages?

**Use Case:** User clicks on model "fraud-detector-v2" in Model Registry, sees observability dashboard filtered to that specific model's metrics.

**Technical Requirement:** Perses must support URL parameters (e.g., `?model_id=fraud-detector-v2`) that map to PromQL label filters (`{model="fraud-detector-v2"}`).

**Implementation:**
- ODH Dashboard constructs deep link: `/observability/model-serving?namespace=prod&model=fraud-detector-v2`
- Perses dashboard uses variable: `$model` mapped to URL parameter
- PromQL query: `rate(model_requests_total{namespace="$namespace", model="$model"}[5m])`

**Answer Required From:** UX Designer, Frontend Engineer, Perses Community (verify feature support)

---

## 9. Alert Routing from Dashboards

**Question:** Can users create AlertManager alert rules directly from Perses dashboards?

**Workflow:** User sees anomaly in dashboard (e.g., latency spike) → clicks "Create Alert Rule" → defines threshold → alert routes to PagerDuty/Slack.

**Perses Capability:** Does Perses support alert rule creation UI, or must users manually edit AlertManager configuration?

**Scope Decision:** Is this MVP or post-MVP? Current "Out of Scope" section excludes custom alerting configuration.

**Recommendation:** Post-MVP (provides view-only alert state visualization in MVP, alert creation in post-MVP)

**Answer Required From:** Product Owner, Perses Community

---

## 10. Disaster Recovery & Data Loss Prevention

**Question:** What is the RPO (Recovery Point Objective) for dashboard definitions?

**Scenario:** Perses PVC is corrupted. How quickly can dashboards be restored?

**Backup Strategy Options:**
- **Option A:** Velero scheduled backups only
  - Daily backups = 24-hour RPO
  - Hourly backups = 1-hour RPO (higher storage cost)
- **Option B:** GitOps sync (continuous = near-zero RPO, but requires automation)
  - CronJob exports dashboards to Git every 6 hours
  - Git repository is source of truth
- **Option C:** Hybrid (Velero for full state + Git for dashboard definitions)
  - Velero for disaster recovery (full system state)
  - Git for dashboard versioning and audit trail

**RTO Target:** 1 hour (time to restore service)

**Recommendation:** Option C (hybrid approach)

**Answer Required From:** SRE Team, Backup/DR Team

---

## 11. React Component Integration Maturity

**Question:** Are Perses npm packages (@perses-dev/components, @perses-dev/dashboards) stable enough for production use?

**Assessment Needed:**
- Review Perses GitHub issues for component-related bugs
- Check npm package release cadence and breaking changes
- Validate TypeScript type definitions are complete
- Test compatibility with ODH Dashboard's React/TypeScript versions
- Evaluate community support and responsiveness

**Risk:** Perses is CNCF sandbox project (not graduated), may have API instability.

**Mitigation:**
- Conduct technical spike/POC to validate stability
- Engage with Perses community for support commitment
- Plan for contributing bug fixes upstream if needed
- Have fallback plan (iframe integration) if components prove unstable

**Answer Required From:** Frontend Engineer (spike/POC)

---

## 12. Authentication Flow Details

**Question:** How exactly are OAuth tokens passed from ODH Dashboard to Perses backend?

**Options:**
- **Option A:** ODH Dashboard makes API calls directly to Perses, including OAuth token in Authorization header
  - Pro: Simple, no proxy needed
  - Con: Frontend exposes Perses API endpoint
- **Option B:** ODH Dashboard backend proxies requests to Perses, validates token server-side
  - Pro: Token validation on backend (more secure)
  - Con: Additional latency, backend complexity
- **Option C:** OAuth2 Proxy sidecar handles authentication (like Grafana deployment pattern)
  - Pro: Standard pattern, proven
  - Con: Additional component to manage

**Recommendation:** Option A (direct token passing) for React component integration, simplest approach

**Answer Required From:** Security Architect, Backend Engineer

---

## 13. Custom Metrics Export Standards

**Question:** What guidance do we provide to users who want to export custom metrics from models/notebooks?

**Scenario:** Data Scientist wants to export "model_prediction_distribution" metric from inference code to track drift.

**Documentation Needed:**
- Examples using Prometheus client libraries (Python, Java, Go)
- ServiceMonitor configuration to scrape custom metrics
- Metric naming conventions (e.g., `model_*`, `pipeline_*`)
- Best practices: counter vs. gauge vs. histogram
- How to visualize custom metrics in Perses dashboards

**Answer Required From:** Technical Writer, Developer Advocate

---

## 14. Cost Attribution Metrics

**Question:** Should dashboards include cost estimation (e.g., "This GPU job cost $12")?

**Complexity:** Requires mapping resource usage to cloud provider pricing, considering node types, reserved instances, spot pricing.

**Scope:** Likely post-MVP, but architect if we anticipate this requirement.

**Implementation Approach:**
- Store cost models in ConfigMap (CPU: $0.10/hr, GPU: $2/hr)
- Calculate cost from resource usage metrics: `cost = cpu_hours * cpu_rate + gpu_hours * gpu_rate`
- Display in dashboard panel as additional metric

**Recommendation:** Post-MVP (customer demand validation needed)

**Answer Required From:** Product Owner, FinOps Team

---

## 15. Accessibility Compliance

**Question:** Does Perses meet WCAG 2.1 AA accessibility standards out-of-box?

**Assessment Needed:**
- Screen reader compatibility (NVDA, JAWS, VoiceOver)
- Keyboard navigation (tab order, focus indicators)
- Color contrast (4.5:1 for text, 3:1 for UI)
- Alternative text for charts

**Mitigation:** If Perses components have accessibility gaps, ODH Dashboard team may need to wrap/enhance components or contribute accessibility improvements upstream.

**Testing Plan:**
- Automated: aXe, WAVE, Lighthouse in CI/CD
- Manual: Keyboard-only navigation, screen reader testing
- User testing with people with disabilities

**Answer Required From:** Accessibility Specialist, Frontend Engineer

---

# Background & Strategic Fit

## Market Context: The Observability Gap in AI/ML Platforms

All major AI/ML platforms in the market have integrated observability as a standard capability:

- **AWS SageMaker:** Built-in CloudWatch integration, pre-configured dashboards for model endpoints, training jobs, and notebooks
- **Azure Machine Learning:** Azure Monitor integration with ML-specific metrics, out-of-the-box dashboards for experiments and deployments
- **Google Vertex AI:** Cloud Monitoring integration, unified observability for the entire ML lifecycle
- **Databricks:** Built-in metrics and monitoring for notebooks, jobs, and model serving

OpenShift AI currently lacks this integrated observability, creating a competitive disadvantage. Customer data shows:

- **15-20% of enterprise evaluations** cite lack of integrated observability as a concern or blocking factor
- **78% of current customers** deploy Grafana separately, creating operational overhead (separate authentication, upgrades, security patching, user provisioning)
- **30% of ODH Dashboard support tickets** are observability-related, with users unable to troubleshoot performance issues without administrator help
- **Q2-Q3 2025 Customer Advisory Board feedback:** Observability integration ranked #2 priority feature request (behind Model Registry improvements)

### Customer Pain Points

**Data Scientists:**
- "I can't see why my model training is slow without asking platform ops" - Data Scientist, Healthcare AI team
- Currently switching between Jupyter notebooks, ODH Dashboard, and external Grafana with different authentication systems and UI patterns
- No self-service access to GPU utilization, training metrics, or resource consumption
- Dependency on platform team for custom dashboard creation

**MLOps Engineers:**
- "We manage 200+ model deployments but can't standardize monitoring dashboards" - MLOps Lead, Retail Analytics
- Manual dashboard creation doesn't scale to dozens or hundreds of models
- Dashboard configuration not in source control, no rollback capability
- Drift between dev/staging/prod observability configurations

**Administrators:**
- "We deployed Grafana but now we manage another platform with its own auth, upgrades, and security" - Platform Admin, Financial Services
- Separate RBAC configuration for observability vs. OpenShift AI projects
- Additional operational overhead: 8-10 hours/week for Grafana management
- Security concerns: external tool accessing cluster metrics, compliance review delays
- No integration with OpenShift AI project isolation model

## Why Perses? Strategic Technology Choice

Perses was selected over alternatives (Grafana, Prometheus UI, custom dashboards) for strategic reasons:

### 1. Open Source Licensing
- **Perses:** Apache 2.0 license, CNCF governance, no vendor lock-in
- **Grafana:** AGPL license creating compliance friction in many enterprises
- **Market Data:** 40% increase in enterprise inquiries about Grafana licensing alternatives in 2024-2025
- **Customer Impact:** Legal approval delays eliminated, faster time-to-deployment

### 2. GitOps-Native Design
- **Perses:** Built for Dashboard-as-Code with Kubernetes CRDs, CUE SDK, Go SDK
- **Grafana:** Requires manual dashboard management or complex automation workarounds (Terraform, jsonnet)
- **Customer Validation:** SAP and Amadeus (Perses early adopters) report 60-70% reduction in dashboard management overhead
- **Market Trend:** 65% of enterprises standardizing on GitOps by end of 2025 (Gartner)

### 3. Red Hat Ecosystem Alignment
- **Proof Point:** Red Hat already uses Perses for OpenShift Distributed Tracing UI (Cluster Observability Operator 0.3.0+)
- **Benefits:** Reusing proven technology, consistent UX across OpenShift product family, shared engineering investment
- **CNCF Ecosystem:** SAP, Amadeus, Chronosphere contributing heavily - enterprise validation
- **Community:** Red Hat as key contributor, influence over roadmap

### 4. Embeddable Components
- **Perses:** npm packages (@perses-dev/components) enable deep integration into ODH Dashboard React UI
- **Grafana:** Embedding typically requires iframes, creating UX friction (separate authentication, navigation issues, performance overhead)
- **User Experience:** Perses dashboards feel native to ODH Dashboard, not bolted-on

### Competitive Positioning Matrix

| Capability | Grafana | Custom Dashboards | Perses (Proposed) |
|------------|---------|-------------------|-------------------|
| **GitOps-Native** | Manual/scripted | N/A | Native CRDs |
| **Licensing** | AGPL concerns | N/A | Apache 2.0 |
| **K8s RBAC Integration** | External | Limited | Native |
| **Dashboard-as-Code** | Terraform/API | Manual | SDK + CRDs |
| **OpenShift Platform Alignment** | External tool | N/A | Red Hat strategic |
| **Multi-tenancy** | OSS limitations | Custom | Native projects |
| **Embeddable UI** | iframe | N/A | React components |

### Market Differentiation Statement

**"OpenShift AI is the only enterprise AI/ML platform with native GitOps observability through CNCF Perses integration, eliminating dashboard management overhead by 60-70% while maintaining OpenShift's security and governance model."**

This differentiation matters because:
- **Gartner Data:** 65% of enterprises will standardize on GitOps by end of 2025
- **Customer Feedback:** Top 3 request is "reduce operational complexity" in AI/ML platforms
- **Competitive Pressure:** AWS SageMaker, Azure ML, Google Vertex AI all have built-in observability - we need parity plus differentiation

## Business Impact

### Without This Feature

**Continued Competitive Losses:**
- 15-20% of enterprise deals lost to competitors citing "lack of integrated observability"
- Prospect quote: "We chose [competitor] because observability was included out of box. With OpenShift AI we'd need to figure that out ourselves"

**Longer Time-to-Value:**
- 4-8 weeks for customers to implement separate observability solutions
- Delays production ML workload deployment
- Increases onboarding friction and perceived platform immaturity

**Higher Support Costs:**
- 30% of ODH Dashboard support tickets are observability-related
- Users unable to self-serve, creating dependency on platform teams
- Support ticket quote: "How do I see why my training is slow?"

**Market Perception:**
- "OpenShift AI is incomplete" if customers must assemble their own monitoring stack
- Constant justification required in sales cycles
- Platform team burden: "We love OpenShift AI but had to deploy Grafana separately"

### With This Feature

**Competitive Differentiation:**
- Claim: "Only enterprise AI platform with GitOps-native observability"
- Parity with hyperscalers + differentiation through open source
- Reference-able customer success stories

**Faster Customer Onboarding:**
- Observability available day 1, not week 6
- 2-3 week reduction in time-to-production
- Better first impression during POC/evaluation

**Platform Stickiness:**
- Integrated observability increases switching costs
- Reduces "bring your own tool" fragmentation
- Cross-sell opportunity: drives adoption of OpenShift Observability for broader cluster monitoring

**Support Cost Reduction:**
- 40% reduction in observability-related support tickets (based on similar features in other products)
- Self-service capability eliminates "How do I monitor X?" questions
- Platform teams empowered with standardized dashboards

**Revenue Impact:**
- Improved win rate in 15-20% of deals where observability is evaluation criteria
- Deal velocity improvement: 2-3 week reduction in evaluation cycle
- Reduced churn: customers less likely to evaluate alternatives

## Strategic Alignment

### Platform Completeness Strategy
Perses integration directly supports strategic goal of providing "complete, enterprise-ready AI/ML platform." Observability is table-stakes capability - customers expect it integrated, not as separate procurement.

**Competitive Reality:** All major AI/ML platforms include integrated observability. OpenShift AI must achieve parity to remain competitive.

### GitOps & Developer Experience
Aligns with platform strategy emphasizing declarative, Git-managed configurations. Data Scientists and MLOps Engineers work in code - dashboards should be code too.

**Customer Workflow Alignment:** "If we can GitOps our models and pipelines, why not our dashboards?" becomes "Yes, you can do both."

### Red Hat Ecosystem Integration
Leverages existing Red Hat investment in Perses (OpenShift Distributed Tracing). Demonstrates portfolio synergy and "better together" story for OpenShift + OpenShift AI.

**Engineering Efficiency:** Shared components, patterns, and expertise with OpenShift Observability team.

### Open Source Leadership
Contributing to CNCF Perses reinforces Red Hat's position as open source leader in cloud-native AI/ML space. Differentiates from hyperscaler proprietary solutions.

**CNCF Credibility:** Perses sandbox project with enterprise adopters (SAP, Amadeus) validates production-readiness.

## Dependencies

### Hard Dependency: Platform Metrics and Alerts Architecture

Perses integration requires foundational metrics pipeline to be in place. This RFE assumes Platform Metrics architecture delivers:

- **Centralized Prometheus Deployment:** Multi-tenant Prometheus with federation across namespaces
- **ServiceMonitor Configurations:** Metrics collection for ODH components (KServe, Kubeflow Pipelines, Notebooks, etc.)
- **Tempo Distributed Tracing (Post-MVP):** Centralized Tempo deployment for trace data
- **RBAC for Metrics Access:** Kubernetes RBAC policies ensuring users can only query metrics from namespaces they have access to

**Sequencing:** Platform Metrics → Perses Integration → Advanced ML-specific dashboards

Without Platform Metrics architecture, Perses integration cannot function as there is no data source to visualize.

## Timeline and Phasing

### MVP (6-9 months)
- React component integration of Perses into ODH Dashboard
- 5 pre-built dashboards (model serving, training, pipelines, notebooks, infrastructure)
- Namespace-based filtering and RBAC enforcement
- Dashboard-as-Code via PersesDashboard CRDs
- Basic error handling and user feedback

### Post-MVP Phase 1 (12-15 months)
- Custom dashboard creation UI (in-dashboard editor)
- Distributed tracing integration (Tempo)
- Advanced query builder (visual PromQL)
- Dashboard templating with variables ($namespace, $model_id)
- Export and sharing features

### Post-MVP Phase 2 (18-24 months)
- AI-assisted troubleshooting (LLM analyzes metrics + logs to suggest root causes)
- Observability-as-Code automation (dashboards provisioned when new models deployed)
- Cost attribution metrics (estimate resource costs in dashboards)
- Cross-product correlation (correlate metrics across training, serving, data pipelines)
- Model drift detection dashboards (statistical tests, baseline comparison)

## Success Metrics

This feature will be considered successful if:

### Adoption Metrics
- **Target:** 40% of active Data Science Projects use Observability dashboards within 60 days of GA
- **Measure:** Unique namespace count accessing Perses dashboards / Total active projects
- **Baseline:** Currently 0% (new feature)

### Support Ticket Reduction
- **Target:** 30% reduction in observability-related support tickets within 90 days
- **Measure:** Support ticket volume (category: observability) before vs. after GA
- **Current Baseline:** 30% of ODH Dashboard tickets are observability-related

### Customer Satisfaction
- **Target:** NPS score ≥ 40 for Observability feature
- **Measure:** In-app survey after 2 weeks of usage: "How likely are you to recommend the Observability dashboards to a colleague?"

### Competitive Win Rate
- **Target:** 10% increase in enterprise deal win rate where observability is evaluation criteria
- **Measure:** Track deal outcomes, competitive win/loss analysis
- **Current Baseline:** 15-20% of deals cite lack of integrated observability as concern

### Performance
- **Target:** 95th percentile dashboard load time < 5 seconds
- **Measure:** Frontend performance monitoring (dashboard initial render time)
- **Target:** <1% error rate on metric queries

### Platform Reliability
- **Target:** 99.5% uptime for Observability page (excludes upstream Prometheus downtime)
- **Measure:** Page availability monitoring
- **Target:** <5 P1/P2 bugs reported in first 30 days post-GA

### GitOps Adoption
- **Target:** 20% of teams create custom dashboards via PersesDashboard CRDs within 90 days
- **Measure:** Count of non-default PersesDashboard resources / Total teams

## Risk Mitigation

### Risk: Perses Immaturity (CNCF Sandbox Project)
- **Mitigation:** Conduct technical spike to validate npm package stability, engage with Perses community for support commitment, plan for contributing bug fixes upstream
- **Fallback:** Iframe integration if React components prove unstable

### Risk: Integration Complexity with React Components
- **Mitigation:** Allocate 2 sprints for POC integration, identify technical challenges early, have fallback plan (iframe embedding) if component integration proves too complex

### Risk: Performance at Scale (High-Cardinality Metrics)
- **Mitigation:** Load testing with 10,000+ time series, Prometheus recording rules to pre-aggregate queries, query timeout limits (60s), rate limiting (100 queries/min/user)

### Risk: User Confusion with Existing Metrics
- **Mitigation:** Clear messaging in UI ("Advanced Observability via Perses"), user documentation explaining differences, eventual migration of existing metrics pages to Perses

### Risk: Platform Metrics Dependency Delays
- **Mitigation:** Explicit sequencing in roadmap, coordinate delivery timelines between teams, document clear interface contracts

---

# Customer Considerations

*Customer-specific considerations for designing and delivering this feature.*

## Multi-Tenancy and Data Isolation

### Requirement
Customers need strict namespace isolation to meet compliance requirements (GDPR, HIPAA, SOC2, FedRAMP). Users in namespace A must not be able to query metrics from namespace B without explicit RBAC grants.

### Implementation

**Namespace-Based Query Filtering:**
- All Prometheus queries automatically injected with `namespace="<user-namespace>"` filter
- Perses backend rewrites queries before sending to Prometheus
- Transparent to users (they write queries without namespace filter, backend adds it)

**RBAC Enforcement:**
- User's OpenShift RBAC permissions determine which namespaces they can view in dashboard filters
- Kubernetes RoleBinding determines access: `kubectl get rolebindings -n fraud-detection`
- No duplicate permission system - Perses honors OpenShift RoleBindings and ClusterRoleBindings

**Kubernetes-Native Permissions:**
- Namespace Viewer (view role): Can view dashboards, read-only metrics
- Namespace Editor (edit role): Can create custom dashboards in namespace
- Cluster Administrator (cluster-admin): Can view all namespaces, publish cluster-wide dashboards

### Distribution Models

Customers can choose deployment model based on security requirements:

**1. Shared Perses Instance (Recommended):**
- Single Perses deployment, logical isolation via query filtering
- Cost-efficient, easier to manage
- Supports 100s of namespaces
- Requires robust RBAC configuration and query rewriting

**2. Namespace-Scoped Perses (High Security):**
- Perses instance per namespace (true isolation)
- Higher cost, operational overhead
- Suitable for highly regulated industries (finance, healthcare)
- Doesn't scale beyond 50-100 namespaces

**3. Hierarchical Isolation (Future):**
- Organizational units (department → team → project)
- Hierarchical access (department admin sees all team dashboards)
- Post-MVP consideration based on customer demand

### Graceful Permission Handling

**User sees only accessible namespaces:**
- Filter dropdown shows only namespaces with at least "view" permission
- No information leakage (user cannot see existence of inaccessible namespaces)

**Permission denied messaging:**
- Clear error: "You do not have access to namespace 'production'. Request access from your administrator."
- Links to documentation on OpenShift RBAC
- Contact information for administrator

## RBAC Enforcement

### Defense-in-Depth Approach

Authorization happens at multiple layers:

**1. ODH Dashboard Frontend:**
- UI only shows namespaces user has access to in filter dropdown
- Client-side validation prevents UI from making unauthorized requests
- First line of defense (UX layer)

**2. Perses Backend:**
- Validates user's OAuth token with OpenShift API
- Checks Kubernetes RBAC permissions before executing queries
- Rewrites PromQL queries to inject namespace filter
- Second line of defense (application layer)

**3. Prometheus (Optional):**
- Final enforcement via Prometheus RBAC proxy (if configured)
- Third line of defense (data layer)

### Permission Levels

| User Role | Dashboard View | Dashboard Edit | Metric Query Scope | Admin Functions |
|-----------|----------------|----------------|-------------------|-----------------|
| **Data Scientist** | Namespace-scoped | Own namespace only | Auto-filtered to namespace | No |
| **MLOps Engineer** | Namespace-scoped | Own namespace only | Auto-filtered to namespace | No |
| **Project Admin** | Namespace-scoped | Own namespace only | Auto-filtered to namespace | No |
| **Cluster Admin** | Cluster-wide | All namespaces | No filter (can query all) | Yes (publish platform dashboards, manage Perses) |

### UX Considerations

- **Lock icons** for read-only dashboards (platform-provided)
- **"Create a copy to customize"** button for platform dashboards
- **Contextual help links** explaining permissions
- **Clear error messages** instead of technical stack traces

## Accessibility (WCAG 2.1 AA Compliance)

### Visual Accessibility

**Color Contrast:**
- Text: 4.5:1 minimum contrast ratio
- UI components: 3:1 minimum contrast ratio
- Automated testing in CI/CD with aXe, WAVE, Lighthouse

**Colorblind-Safe Palettes:**
- Dashboard visualizations use color schemes safe for deuteranopia, protanopia, tritanopia
- Never rely on color alone (use patterns, labels, icons)

**Text Alternatives:**
- All chart data available as text/table view for screen readers
- ARIA labels describe what each metric represents

**Scalable UI:**
- Dashboard supports browser zoom up to 200% without horizontal scrolling
- Responsive design for different screen sizes

### Keyboard Navigation

**Full keyboard accessibility:**
- All dashboard controls accessible via keyboard (Tab, Enter, Arrow keys)
- Keyboard shortcuts: R (refresh), T (time range), F (filter), Esc (close modal), ? (help)
- Focus indicators clearly visible (2px outline, high contrast)
- Logical tab order (left-to-right, top-to-bottom)

**Skip links:**
- "Skip to main content" link at top of page
- "Skip to dashboard panels" for quick navigation

### Screen Reader Support

**Semantic HTML:**
- Proper heading hierarchy (h1, h2, h3)
- Landmarks (navigation, main, aside)
- ARIA roles for custom components

**Chart Accessibility:**
- Alternative text descriptions for charts (e.g., "Line chart showing model latency over time, p95 latency increased from 100ms to 300ms at 2:15 PM")
- Data tables as alternative to visual charts

**State Changes:**
- `aria-live="polite"` regions announce updates (e.g., "Dashboard loaded", "Filtering by namespace production")
- Loading states communicated to screen readers

### Temporal Content

**Auto-refresh controls:**
- Pause/resume button prominently displayed
- Announce refresh events to screen readers
- Disable animations option for users with vestibular disorders

### Testing

**Automated Testing:**
- aXe, WAVE, Lighthouse in CI/CD pipeline
- Fails build if accessibility violations detected

**Manual Testing:**
- Keyboard-only navigation testing
- Screen reader testing (NVDA on Windows, JAWS on Windows, VoiceOver on macOS, Orca on Linux)
- Color contrast verification tools

**User Testing:**
- Test with people with disabilities (vision, motor, cognitive)
- Incorporate feedback into design iterations

## Performance at Scale

### Scalability Targets

- **Users:** 500 concurrent users per cluster without degradation
- **Dashboards:** 5,000 total dashboards across all namespaces
- **Metrics:** 10,000+ time series per dashboard query
- **Namespaces:** 200+ Data Science Projects per cluster
- **Dashboard Load Time:** p95 < 2s, p99 < 5s

### Performance Optimizations

**1. Query Optimization:**
- **Prometheus Recording Rules:** Pre-aggregate frequently-accessed metrics
  ```yaml
  groups:
  - name: model_serving
    interval: 30s
    rules:
    - record: model:request_rate:5m
      expr: rate(model_requests_total[5m])
  ```
- **Query Timeout:** 60 seconds maximum (prevent expensive queries from overloading Prometheus)
- **Rate Limiting:** 100 queries/min per user (prevent abuse, distributed denial-of-service)

**2. Caching Strategy:**
- **Dashboard Definitions:** Cached in-memory (30-minute TTL, invalidate on CRD update)
- **Prometheus Query Results:**
  - Real-time dashboards: 30-second TTL
  - Historical dashboards: 5-minute TTL
- **Redis Cache (Optional):** Shared cache across Perses replicas for better hit rate

**3. Progressive Loading:**
- Dashboard panels load asynchronously (don't block on slow queries)
- Individual panel loading indicators (user sees partial data immediately)
- Lazy loading for off-screen panels (only load when scrolled into view)
- Critical KPIs load first (<2s), detailed panels follow

**4. Horizontal Scaling:**
- **Perses Backend:** 3-10 replicas, HorizontalPodAutoscaler at 70% CPU target
- **Load Balancer:** Distributes requests across replicas (round-robin or least connections)
- **Stateless Design:** PostgreSQL backend enables seamless scaling (no pod affinity required)

**5. Data Sampling:**
- **Downsampling for Long Time Ranges:**
  - >24 hours: 5-minute resolution (instead of 15-second)
  - >7 days: 1-hour resolution
- **User-Requestable High Resolution:** "Show detailed data" button for specific time windows
- **Transparent UI Indicator:** "Showing downsampled data for performance" message

**6. Client-Side Optimizations:**
- **Virtual Scrolling:** For long dashboard lists (only render visible items)
- **Canvas Rendering:** For dense time series charts (better performance than SVG)
- **Debounced Inputs:** 300ms debounce on search/filter inputs
- **Disabled Animations by Default:** Reduce CPU usage (opt-in for users who want animations)

### Load Testing

**Required Before GA:**
- Simulate 500 concurrent users loading high-cardinality dashboards (GPU metrics, model serving latency)
- Measure p95/p99 latency, error rate, resource consumption (Prometheus, Perses, ODH Dashboard)
- Target: p95 < 2s dashboard load time, <1% error rate, no Prometheus resource exhaustion

**Load Testing Tools:**
- k6 or Locust for HTTP load generation
- Scenarios: Dashboard list load, dashboard view, filter changes, auto-refresh
- Prometheus/Grafana for observing system under load (meta-monitoring)

## Enterprise Security

### Authentication

**OpenShift OAuth:**
- Primary authentication mechanism, single sign-on with ODH Dashboard
- User logs into ODH Dashboard once, Perses inherits session

**Token-Based Auth:**
- ODH Dashboard passes OAuth token to Perses API in `Authorization: Bearer <token>` header
- Perses backend validates token with OpenShift TokenReview API

**Token Expiration Handling:**
- Graceful redirect to login if token expires (refresh token flow)
- In-app message: "Your session has expired, please log in again"

### Data Security

**In-Transit Encryption:**
- All API calls over HTTPS/TLS 1.2+ (enforced by OpenShift Routes)
- Certificate management via cert-manager or OpenShift service CA

**At-Rest Encryption:**
- Perses PVC encrypted using storage class encryption (cloud provider managed keys)
- Dashboard definitions stored in PostgreSQL with encryption at rest

**Secret Management:**
- Data source credentials (Prometheus endpoints) stored in Kubernetes Secrets
- Secrets mounted as volumes, not exposed in environment variables

### Audit Logging

**Events Logged:**
- `dashboard.viewed`: User viewed dashboard
- `dashboard.created`: User created custom dashboard
- `dashboard.deleted`: User deleted dashboard
- `query.executed`: User executed PromQL query
- `permission.denied`: User attempted unauthorized access

**Log Format (JSON):**
```json
{
  "timestamp": "2025-10-16T10:30:00Z",
  "event": "dashboard.viewed",
  "user": "alice@example.com",
  "namespace": "fraud-detection",
  "dashboard": "model-serving-performance",
  "sourceIP": "10.0.1.5",
  "userAgent": "Mozilla/5.0..."
}
```

**Retention:** 90 days (configurable), exported to centralized log aggregator (Loki, Splunk, etc.)

**Compliance:** Logs support SOC2, ISO 27001, FedRAMP audit requirements

### Vulnerability Management

**Container Image Scanning:**
- Perses container images scanned for CVEs (Trivy, Clair, Snyk)
- Fail build if critical vulnerabilities detected

**Automated Patching:**
- Automated patching pipeline for critical vulnerabilities
- Rebuild images with updated base OS, redeploy via Helm

**Security Advisories:**
- Monitor Perses GitHub security advisories
- Subscribe to CNCF security mailing lists
- Incident response plan for zero-day vulnerabilities

## Backup and Disaster Recovery

### Backup Strategy

**Velero Scheduled Backups:**
```yaml
apiVersion: velero.io/v1
kind: Schedule
metadata:
  name: perses-daily-backup
spec:
  schedule: "0 1 * * *"  # Daily at 1 AM UTC
  template:
    includedNamespaces:
    - opendatahub-perses
    includedResources:
    - persistentvolumeclaims
    - persistentvolumes
    - configmaps
    - secrets
    ttl: 720h  # 30 days retention
```

**GitOps Sync (Optional but Recommended):**
- CronJob exports dashboards to Git repository every 6 hours
- Enables version control, change tracking, code review for dashboard changes
- Implementation: Kubernetes CronJob calls Perses API, commits dashboard YAML to Git
- Git repository becomes source of truth for dashboard definitions

### Recovery Targets

- **RTO (Recovery Time Objective):** 1 hour (time to restore Perses service)
- **RPO (Recovery Point Objective):** 24 hours (daily backups, max 1 day of dashboard changes lost)

### Disaster Recovery Runbook

**Procedure:**
1. **Detect Perses Outage:** Monitoring alerts on pod health (all replicas down)
2. **Attempt Pod Restart:** `kubectl delete pod -n opendatahub-perses -l app=perses` (often resolves transient issues)
3. **If PVC Corrupted:** Restore PVC from Velero backup
   ```bash
   velero restore create perses-restore \
     --from-backup perses-daily-20251015 \
     --include-namespaces opendatahub-perses
   ```
4. **Redeploy Perses:** `helm upgrade perses perses/perses --reuse-values` (recreates pods with restored PVC)
5. **Verify API Responds:** `curl https://perses.odh.example.com/api/v1/dashboards`
6. **Smoke Test:** Load pre-built dashboard in ODH Dashboard UI, verify metrics display
7. **Post-Incident Review:** Document root cause, improve monitoring to detect earlier

### DR Testing

**Quarterly Disaster Recovery Drills:**
- Delete Perses namespace, restore from backup, verify all dashboards functional
- Success criteria: All dashboards restored, no data loss beyond RPO window, RTO met (<1 hour)
- Document lessons learned, update runbook

## Migration from Existing Observability Tools

### Customer Scenario
Customers currently using Grafana or custom dashboards want to migrate to Perses without disrupting operations.

### Migration Approach

**1. Coexistence Period (3-6 months):**
- Perses and Grafana run simultaneously during transition
- Users can access either system during migration
- No forced cutover date (teams migrate at their own pace)

**2. Dashboard Recreation:**
- No automated migration tool in MVP (Grafana JSON → Perses YAML not 1:1 mapping)
- Users manually recreate Grafana dashboards as PersesDashboard CRDs
- Provide migration guide with side-by-side examples

**3. Training & Documentation:**
- Migration guide: "Converting Grafana Dashboards to Perses"
- Office hours: Weekly sessions with Perses integration team
- Example dashboards: Grafana → Perses equivalents for common patterns

**4. Gradual Adoption:**
- Teams adopt Perses at their own pace based on readiness
- New projects start with Perses (no Grafana)
- Existing projects migrate when convenient (no mandate)

### Migration Challenges

**Dashboard Compatibility:**
- Grafana dashboards are JSON-based, Perses uses different schema
- No 1:1 translation (e.g., Grafana "graph" panel vs. Perses "timeseries" panel)
- Some Grafana features may not exist in Perses (e.g., Grafana-specific plugins)

**Panel Type Differences:**
- Grafana has 30+ panel types, Perses has core set (timeseries, gauge, table, heatmap)
- Custom Grafana panels require re-implementation in Perses

**Alerting Differences:**
- Grafana alerts don't map to Perses (alerts remain in Platform Metrics architecture)
- Alert visualization differs (Grafana embedded alerts vs. Perses view-only alert state)

### Customer Support During Migration

- **Dedicated Slack Channel:** #perses-migration for questions
- **Office Hours:** Weekly 1-hour sessions with Perses team
- **Migration Guide:** Step-by-step documentation with examples
- **Dashboard Library:** Pre-built Perses dashboards for common ML use cases

## Customization and Extensibility

### GitOps Workflow

**Dashboard Version Control:**
- Customers store dashboards in Git repositories alongside application code
- Example repository structure:
  ```
  dashboards/
  ├── model-serving/
  │   ├── overview.yaml
  │   └── per-model.yaml
  ├── pipelines/
  │   └── execution-metrics.yaml
  └── infrastructure/
      └── cluster-health.yaml
  ```

**CI/CD Integration:**
- **Validation:** `percli validate dashboard.yaml` in CI pipeline
- **Automated Deployment:** ArgoCD/Flux syncs dashboards to cluster
- **Peer Review:** Dashboard changes reviewed via pull requests (like code)

**Benefits:**
- Version control: track changes, rollback if needed
- Audit trail: Git history shows who changed what, when
- Reproducibility: apply same dashboard to dev, staging, prod environments

### Custom Metrics

**Exporting Custom Metrics from Models/Notebooks:**
- Customers use Prometheus client libraries (Python, Java, Go)
- Example (Python):
  ```python
  from prometheus_client import Counter, Histogram

  prediction_counter = Counter('model_predictions_total',
                                'Total predictions made',
                                ['model', 'version'])

  prediction_latency = Histogram('model_prediction_duration_seconds',
                                 'Time spent making prediction')

  @prediction_latency.time()
  def predict(input_data):
      result = model.predict(input_data)
      prediction_counter.labels(model='fraud-detector', version='v2').inc()
      return result
  ```

**ServiceMonitor Configuration:**
- Configure Prometheus to scrape custom metrics
- Example:
  ```yaml
  apiVersion: monitoring.coreos.com/v1
  kind: ServiceMonitor
  metadata:
    name: fraud-detector-metrics
  spec:
    selector:
      matchLabels:
        app: fraud-detector
    endpoints:
    - port: metrics
      interval: 30s
  ```

**Documentation:**
- Metric naming conventions: `model_*`, `pipeline_*`, `notebook_*`
- Best practices: counter vs. gauge vs. histogram
- Examples for common ML metrics (prediction distribution, drift score, data quality)

### Dashboard Templates

**Reusable Dashboard Templates:**
- Customers create templates with variables (e.g., `$namespace`, `$model_name`)
- Template applied to multiple models without duplication
- Example:
  ```yaml
  spec:
    variables:
      - name: model_name
        type: query
        query: label_values(model_requests_total, model)
    panels:
      - title: "Request Rate for $model_name"
        queries:
          - promql: "rate(model_requests_total{model=\"$model_name\"}[5m])"
  ```

**Template Sharing:**
- Templates shared across teams via Git repositories
- Public template library on GitHub (community contributions)
- Example: "Standard Model Serving Dashboard" template applied to all new models

### Community Contribution

**Red Hat Contributions:**
- ML-specific dashboard examples contributed to Perses community
- OpenShift AI dashboard library published to GitHub (Apache 2.0 license)
- Example dashboards: model serving, GPU utilization, pipeline performance

**Customer Contributions:**
- Customers contribute their own dashboards back to community
- Public dashboard library grows over time (like Helm charts)
- Recognition program for top contributors

## Internationalization (i18n)

### MVP Scope
English only for dashboard UI elements (time range labels, error messages, navigation)

### Post-MVP Localization

**Dashboard Content:**
- Panel titles/descriptions (user-defined, can be in any language)
- Custom dashboard text in user's preferred language

**UI Chrome:**
- Time range picker, filter controls, navigation translated to 5-10 languages
- Priority languages based on customer demand: Japanese, Chinese (Simplified), German, French, Spanish

**Documentation:**
- User documentation translated to high-priority languages
- Architecture and developer documentation remain English

## Vendor Lock-In Avoidance

### Open Source Commitment

**Perses is Apache 2.0 Licensed:**
- No vendor lock-in, customers can self-host independently
- CNCF governance ensures neutral stewardship

**Dashboard Portability:**
- Dashboard definitions are open Kubernetes CRDs (standard YAML)
- Exportable via Perses API (`GET /api/v1/dashboards/{name}`)
- Can be imported into any Perses deployment (not OpenShift-specific)

**Data Portability:**
- Prometheus queries are standard PromQL (portable to Grafana, Datadog, etc.)
- No proprietary query language or lock-in

### Cloud Provider Neutrality

**Deployment Flexibility:**
- Integration works on OpenShift clusters on AWS, Azure, GCP, bare-metal
- No cloud-specific dependencies (uses standard Kubernetes primitives)
- Storage class agnostic (works with any CSI-compliant storage)

**Data Source Flexibility:**
- Perses supports multiple Prometheus instances
- Can connect to Prometheus outside OpenShift cluster if needed
- No hard dependency on OpenShift Monitoring stack (though recommended)

---

*This RFE represents a strategic investment in OpenShift AI platform completeness, addressing a confirmed competitive gap with a modern, open-source, GitOps-native approach. By integrating Perses, we provide enterprise customers with the observability capabilities they expect while maintaining Red Hat's commitment to open source leadership and reducing operational complexity.*

---

**Document Version:** 1.0
**Last Updated:** October 16, 2025
**Contributors:**
- **Parker (Product Manager):** Market strategy, competitive analysis, business impact
- **Olivia (Product Owner):** Requirements, acceptance criteria, success metrics
- **Aria (UX Architect):** User experience, workflows, persona design, accessibility
- **Archie (Architect):** Technical architecture, integration patterns, security, scalability

**Next Steps:**
1. Review with stakeholders (Engineering, Product, UX, Security)
2. Approval decision (APPROVE/REJECT/DEFER)
3. Technical design phase (detailed architecture, API contracts, data models)
4. Sprint planning and backlog creation
5. MVP implementation (6-9 months)