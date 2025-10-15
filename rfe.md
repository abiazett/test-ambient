# Perses Observability Dashboard Integration for OpenShift AI

**Feature Overview:**

This feature integrates Perses, a CNCF sandbox observability dashboard platform, into the OpenShift AI Dashboard as a new "Observability" navigation item. It provides Data Scientists, MLOps Engineers, and Platform Administrators with advanced, customizable dashboards for monitoring AI workloads, model serving metrics, pipeline performance, and infrastructure health. By embedding Perses within the ODH Dashboard UI, users gain access to powerful visualization capabilities backed by the centralized Prometheus/Tempo infrastructure without leaving their primary workspace.

**Why this matters:** Today, users must leave the OpenShift AI Dashboard to access observability data through separate tools like Grafana or the OpenShift Console, creating context switching and workflow friction. This integration brings comprehensive observability directly into the AI workflow, enabling faster troubleshooting, proactive monitoring, and data-driven optimization of AI workloads.

**Goals:**

* **Unified User Experience:** Provide seamless access to observability dashboards within the OpenShift AI Dashboard without requiring users to navigate to external tools or interfaces.

* **AI Workload Visibility:** Enable comprehensive monitoring of model serving endpoints, pipeline executions, notebook resources, and distributed training jobs through purpose-built dashboards.

* **Self-Service Observability:** Empower Data Scientists and MLOps Engineers to create, customize, and share dashboards tailored to their specific AI workloads and KPIs without requiring deep observability expertise.

* **Enterprise-Ready Integration:** Ensure the solution supports multi-tenancy, RBAC, scalability, and works seamlessly with the centralized Platform Metrics and Alerts architecture.

* **GitOps and Dashboard-as-Code:** Leverage Perses' native support for dashboard definitions as code, enabling version control, CI/CD integration, and declarative dashboard management aligned with modern MLOps practices.

**Who benefits:** Data Scientists gain visibility into model performance and resource utilization; MLOps Engineers can monitor and troubleshoot production AI services; Platform Administrators can oversee infrastructure health and capacity planning across all AI workloads.

**Today vs. Future State:** Currently, users must cobble together observability from disparate tools, manually correlating metrics across systems. With this feature, OpenShift AI becomes a complete platform with integrated observability, reducing time-to-insight and improving operational excellence.

**Out of Scope:**

* **Standalone Perses Deployment:** This RFE covers embedded integration only, not deploying Perses as a separate application outside the ODH Dashboard context.

* **Custom Alerting UI:** Alert configuration and management will remain in the Platform Metrics and Alerts architecture; Perses integration focuses on visualization and dashboard capabilities.

* **Non-Prometheus/Tempo Data Sources:** Initial scope supports only Prometheus metrics and Tempo traces aligned with the platform metrics strategy. Additional data sources (Loki logs, etc.) are future enhancements.

* **Generic Kubernetes Monitoring:** This integration focuses on AI/ML workload observability. General cluster monitoring remains the responsibility of OpenShift Console and cluster administrators.

* **External User Access:** Dashboards are for authenticated OpenShift AI Dashboard users only; public or unauthenticated dashboard sharing is out of scope.

* **Migration of Existing Grafana Dashboards:** Automated conversion or migration tools for existing Grafana dashboards to Perses format are not included in initial scope.

**Requirements:**

**Prerequisites (Must Be Met Before MVP Development):**

* **[PRE-MVP] Platform Metrics Readiness Validation:**
  - Centralized Prometheus deployed with HA configuration (≥2 replicas)
  - Prometheus retention ≥7 days, can handle 1M+ active time series
  - Following metrics available and stable:
    - KServe: request_duration, request_count, request_errors (by namespace, model)
    - Pipelines: run_duration, run_status, task_duration (by namespace, pipeline)
    - Notebooks: pod_cpu_usage, pod_memory_usage, pod_gpu_usage (by namespace, notebook)
    - Training: job_status, training_loss, worker_count (by namespace, job)
  - Prometheus accessible from ODH Dashboard backend pod
  - Service account with read permissions to Prometheus API created
  - Recording rules for common aggregations defined (reduce query load)
  - **GO/NO-GO Decision Gate:** If validation fails, RFE cannot proceed until Platform Metrics is production-ready

* **[PRE-MVP] Perses Technical Validation:**
  - 2-week technical spike to validate Perses capabilities
  - Load testing: Support 100+ concurrent users without degradation
  - Performance: Render 50-panel dashboard in <15 seconds
  - Security: No critical vulnerabilities found
  - Embedding approach (iframe/React/Module Federation) technically feasible
  - **GO/NO-GO Decision Gate:** Proceed only if Perses meets technical requirements

**Core Integration (MVP):**

* **[MVP] Navigation Integration:** Add "Observability" as a new primary navigation item in the ODH Dashboard left navigation menu.

* **[MVP] Perses UI Embedding:** Embed the Perses dashboard viewing and editing interface within the ODH Dashboard. Recommended approach: iframe with backend proxy (MVP), evolve to Module Federation in 18 months. Implementation must ensure:
  - Consistent ODH Dashboard navigation, header, and styling maintained
  - PatternFly v5 components used for all non-Perses UI elements
  - Design token mapping between Perses and PatternFly defined
  - Keyboard navigation and accessibility preserved across iframe boundary

* **[MVP] Prometheus Data Source Configuration:** Automatically configure Perses to connect to the centralized Prometheus instance with:
  - Backend proxy handling authentication (token-mediated flow)
  - RBAC enforcement via query rewriting (inject namespace filters)
  - No user-facing credential configuration required
  - Support for custom CA certificates (air-gapped environments)

* **[MVP] Pre-Built AI Workload Dashboards:** Include curated set of default dashboards covering:
  - Model serving metrics (inference latency p50/p95/p99, throughput, error rates per model)
  - Pipeline execution metrics (success rates, duration, resource consumption)
  - Notebook resource utilization (CPU, memory, GPU usage per notebook)
  - Infrastructure health (cluster-wide GPU utilization, capacity metrics)
  - Note: Distributed training dashboard deferred to Non-MVP pending usage validation

* **[MVP] Multi-Tenancy Support:**
  - Users can only view dashboards and metrics for authorized namespaces (OpenShift RBAC enforcement)
  - Namespace context injection: All queries automatically scoped to user's permitted namespaces
  - Query rewriting engine: Inject `namespace=~"ns1|ns2|ns3"` filters into all PromQL queries
  - Platform admin exception: Cluster-wide access for users with admin ClusterRole
  - Clear error messages identifying unauthorized namespace access attempts

* **[MVP] Dashboard CRUD Operations:** Enable users to create, read, update, and delete custom dashboards:
  - Dashboard storage: Kubernetes CRDs (namespace-scoped) for GitOps compatibility
  - Dashboard list view with search/filter capabilities
  - Dashboard editor with query validation and error highlighting
  - Confirmation modal for destructive actions (delete dashboard)
  - Unsaved changes warning when navigating away from edit mode

* **[MVP] Dashboard Discovery and Onboarding:**
  - First-time user welcome modal explaining Observability purpose
  - Pre-built dashboard gallery with descriptions and preview thumbnails
  - Dashboard list distinguishes pre-built vs custom dashboards (badges/sections)
  - "Create Dashboard" prominent CTA with wizard/query builder for beginners
  - Contextual help links throughout UI ("?" icons, tooltips)

* **[MVP] Variable and Filtering Support:** Allow dashboard creators to define variables (e.g., namespace, model name) that viewers can use to filter dashboard views dynamically. *(Moved from Non-MVP based on Product Manager feedback - critical for production use)*

**Accessibility (MVP - WCAG 2.1 Level AA Compliance):**

* **[MVP] Keyboard Navigation:**
  - All dashboard CRUD operations accessible without mouse
  - Tab order follows logical reading sequence (dashboard list → panels → controls)
  - Focus indicators visible with 3:1 contrast ratio
  - Enter/Space activates buttons and links
  - Arrow keys navigate within time range picker and variable selectors
  - Escape key dismisses overlays/modals

* **[MVP] Screen Reader Support:**
  - Dashboard list announces count: "4 dashboards available"
  - Each dashboard panel has descriptive aria-label
  - Charts include text summary or data table alternative (WCAG 1.1.1)
  - Status updates (loading, error, success) announced via aria-live regions
  - Panel edit controls have clear labels (not icon-only)

* **[MVP] Visual Accessibility:**
  - Text contrast ≥4.5:1, UI components ≥3:1 (WCAG 1.4.3, 1.4.11)
  - Charts use color-blind safe palette (not red/green only)
  - Charts distinguish data series by line style (solid/dashed) + color
  - UI scales to 200% browser zoom without horizontal scrolling (WCAG 1.4.10)
  - Focus indicators have 3:1 contrast ratio

* **[MVP] Time-Based Interactions:**
  - Dashboard auto-refresh is pausable (WCAG 2.2.2)
  - Pause control keyboard accessible and screen reader labeled
  - No time-limited interactions that cannot be extended

**Empty States and Error Handling (MVP):**

* **[MVP] Empty State Patterns:**
  - **Dashboard List Empty:** "No custom dashboards yet" with "Create Dashboard" CTA and "View pre-built dashboards" link
  - **Dashboard No Data:** Context-specific guidance per panel (e.g., "Deploy a model to see serving metrics [Deploy Model button]")
  - **Permission Denied:** "You don't have access to namespace 'prod-ml'" with "Learn about permissions" link
  - **Prometheus Connection Failure:** "Unable to connect to metrics service" with "Try again" button and troubleshooting link
  - **Query Timeout:** "Query timed out. Try shorter time range or simpler query" with adjustable time range picker
  - **Invalid PromQL Query:** Specific error with position indicator and example correction

* **[MVP] Error Message Pattern:** All errors follow: [Problem] + [Impact] + [Action]
  - Example: "Dashboard failed to load. Changes not saved. [Retry] [Export JSON]"

**Performance and Scalability (MVP):**

* **[MVP] Performance Requirements:**
  - Dashboard list loads in <1 second (p95)
  - 10-panel dashboard renders in <5 seconds (p95)
  - 50-panel dashboard renders in <15 seconds (p95)
  - Support 500+ concurrent dashboard viewers
  - Prometheus queries complete in <3 seconds (p95)

* **[MVP] Caching Strategy:**
  - Browser cache: Static assets, 1 hour TTL, versioned by content hash
  - Backend query cache (Redis): 30-second TTL, keyed by query hash + namespace
  - Prometheus native caching: 1-minute TTL for duplicate queries

* **[MVP] Performance Testing:**
  - Load test: 500 concurrent users viewing different dashboards
  - Stress test: Identify breaking point and graceful degradation behavior
  - Soak test: 72-hour sustained load to detect memory leaks

**Lifecycle Management (MVP):**

* **[MVP] Dashboard Persistence:**
  - Dashboards stored as Kubernetes CRDs (namespace-scoped)
  - Automatic dashboard backup before every ODH upgrade
  - Dashboard definitions survive ODH pod restarts/upgrades
  - Dashboard size limit: 1MB (etcd limit), validation webhook enforces

* **[MVP] Version Compatibility:**
  - Document ODH version ↔ Perses version compatibility matrix
  - Test upgrade path: ODH N → N+1 with existing dashboards
  - Dashboard schema versioning (persesv1alpha1)

* **[MVP] Rollback and Recovery:**
  - Feature flag to disable Observability if production issues found
  - Circuit breaker: Auto-disable if Perses crashes 3 times in 5 minutes
  - Dashboard export/import for manual backup/restore
  - Graceful degradation: Show last cached dashboard state if Prometheus unavailable

**Support and Diagnostics (MVP):**

* **[MVP] Self-Service Health Check:**
  - Location: Observability → Settings → "Run Health Check"
  - Tests: Perses running, Prometheus connectivity, auth token valid, namespace access, sample query
  - Output: Pass/Fail per test with actionable remediation steps
  - "Copy diagnostics" button for support tickets

* **[MVP] Diagnostic Bundle:**
  - Perses pod logs and status
  - Dashboard definitions (JSON/YAML export)
  - Prometheus connectivity test results
  - User RBAC configuration
  - Recent error messages

* **[MVP] Support Documentation:**
  - Runbooks for top 5 issue categories (no data, slow dashboard, permission denied, query errors, upgrade issues)
  - Troubleshooting decision tree for "No data showing" scenarios
  - Escalation path documented: L1 → L2 → Engineering → Perses upstream

**Air-Gapped and Disconnected Environments (MVP):**

* **[MVP] Offline Support:**
  - All Perses UI assets bundled with ODH Dashboard (no CDN dependencies)
  - No external API calls at runtime (no version checking, telemetry to external services, DNS lookups)
  - All container images mirrorable to disconnected registries
  - Support for arbitrary registry URLs (no hardcoded docker.io)
  - Custom CA certificate support for internal Prometheus/Tempo HTTPS
  - Automated test runs in cluster with egress blocked (part of CI/CD)
  - Documentation includes air-gapped installation instructions

**Telemetry and Monitoring (MVP):**

* **[MVP] Feature Telemetry (with user opt-in/privacy compliance):**
  - observability_feature_daily_active_users
  - observability_dashboard_views_total (by dashboard_name, user_role)
  - observability_custom_dashboard_created_total
  - observability_dashboard_load_time_seconds (p50, p95, p99)
  - observability_query_errors_total (by error_type)
  - observability_empty_state_views_total (by scenario)
  - **Privacy:** DO NOT capture metric values, namespace names, model names, or user identities

* **[MVP] Feature Health Monitoring:**
  - Alert: observability_feature_availability < 99% over 5 minutes
  - Alert: observability_dashboard_load_time_seconds p95 > 10s over 5 minutes
  - Alert: observability_query_errors_rate > 10% over 5 minutes
  - Runbook for Red Hat SRE on incident response

**Post-MVP Enhancements:**

* **[Non-MVP] Tempo Trace Integration:** Connect Perses to the centralized Tempo instance for distributed tracing visualization of model serving requests and pipeline workflows.

* **[Non-MVP] Dashboard Sharing:** Allow users to share dashboards with other users or teams within their authorized scope, with read-only or edit permissions.

* **[Non-MVP] Dashboard Templates:** Provide a template library where users can instantiate common dashboard patterns with their specific resource selectors.

* **[Non-MVP] Dashboard-as-Code GitOps:** Enable synchronization between Git repositories and deployed dashboards with bi-directional sync.

* **[Non-MVP] Export and Import:** Support exporting dashboards as portable definitions and importing dashboards from files or URLs.

* **[Non-MVP] Dashboard Versioning:** Track dashboard version history and enable rollback to previous versions.

* **[Non-MVP] Distributed Training Dashboard:** Add pre-built dashboard for training jobs after usage validation.

* **[Non-MVP] Advanced Query Builder:** Visual query builder with AI-specific metrics library, template snippets, and natural language to PromQL conversion.

* **[Non-MVP] Contextual Deep Links:** Links from Model Serving detail page to Observability dashboard filtered to specific model, with URL parameters preserving context.

**Done - Acceptance Criteria:**

**Core Functionality:**

* A new "Observability" navigation item appears in the ODH Dashboard left menu for all authenticated users with appropriate permissions.

* Clicking "Observability" loads the Perses dashboard interface embedded within the ODH Dashboard layout, maintaining consistent navigation, header, and styling.

* Users can view at least three pre-built dashboards (model serving, pipelines, notebooks) populated with live metrics from their authorized namespaces.

* First-time users see onboarding guidance (welcome modal or inline tutorial) explaining the purpose of observability and highlighting pre-built dashboards.

* Dashboard list clearly distinguishes pre-built dashboards from custom dashboards (visual indicators like badges or section separators).

* Users with appropriate permissions can create a new dashboard, add panels with Prometheus queries, save the dashboard, and subsequently view it in their dashboard list.

* Users can edit an existing dashboard (add/remove panels, modify queries, adjust layouts) and changes persist correctly.

* Dashboard deletion requires confirmation modal following ODH destructive action pattern.

* Dashboard views respect OpenShift RBAC: users only see metrics and dashboards for namespaces they have access to; attempting to query unauthorized namespaces returns specific error messages identifying the unauthorized namespace.

* The Perses integration automatically authenticates to the centralized Prometheus instance without requiring users to configure credentials or data source connections manually.

**Accessibility (WCAG 2.1 AA Compliance):**

* Dashboard feature passes automated accessibility scan (aXe DevTools) with zero critical violations.

* Manual keyboard navigation test: All CRUD operations (create, read, update, delete dashboards; add/remove panels) completable without mouse.

* Tab order follows logical reading sequence (dashboard list → selected dashboard → panels → time controls → action buttons).

* Focus indicators visible on all interactive elements with 3:1 contrast ratio minimum.

* Screen reader test (NVDA or JAWS): Dashboard purpose, panel contents, and controls are announced and understandable.

* Dashboard list announces count: "4 dashboards available" or "No dashboards found".

* Each dashboard panel has descriptive aria-label (e.g., "Inference latency panel showing p95 latency over time").

* Charts provide text summary or data table alternative accessible to screen readers.

* Color contrast audit passes: All text ≥4.5:1 contrast, UI components ≥3:1 contrast.

* Charts use color-blind safe palette and distinguish data series by both color and line style (solid/dashed/dotted).

* UI remains functional at 200% browser zoom without horizontal scrolling or content overlap.

* Dashboard auto-refresh includes pause control that is keyboard accessible and screen reader announced.

**Empty States and Error Handling:**

* Empty dashboard list shows guided empty state with "Create Dashboard" CTA and link to pre-built dashboards.

* Dashboards with no data show per-panel empty states with contextual guidance and action links (e.g., "Deploy a model to see serving metrics [Deploy Model button]").

* Permission errors display specific messages identifying which namespace is unauthorized: "You don't have access to namespace 'prod-ml'. Contact administrator to request access."

* Prometheus connection failures show user-friendly error with retry mechanism: "Unable to connect to metrics service. [Try Again] [Troubleshooting Guide]".

* Query syntax errors show specific error position and example corrections: "Invalid query syntax at position 15. Expected 'by' keyword. Example: rate(metric[5m]) by (namespace)".

* Dashboard save failures preserve user work with export option: "Failed to save dashboard. Changes not saved. [Retry] [Export JSON]".

* All error messages follow pattern: [Problem] + [Impact] + [Action].

**Performance:**

* Dashboard list loads in <1 second (p95 latency).

* Dashboard with 10 panels renders in <5 seconds (p95 latency) on standard network conditions.

* Dashboard with 50 panels renders in <15 seconds (p95 latency).

* System supports 500 concurrent dashboard viewers without performance degradation (validated via load testing).

* Prometheus queries complete in <3 seconds (p95 latency).

* Query results appropriately cached (30-second TTL) to reduce load on Prometheus.

* Dashboard performance gracefully degrades under high load (shows loading skeletons, not blank screens or errors).

**Lifecycle and Reliability:**

* Dashboard definitions survive ODH platform upgrades (automatic backup tested).

* Rollback from new ODH version to previous version tested successfully with dashboards intact.

* Feature flag can disable Observability feature without affecting rest of ODH Dashboard.

* Circuit breaker isolates Perses failures from rest of ODH Dashboard (Perses crash doesn't bring down Dashboard).

* Graceful degradation when Prometheus unavailable: Shows cached dashboard data or clear error message with troubleshooting link.

* Dashboard size validation prevents dashboards exceeding 1MB (etcd limit).

**Support and Diagnostics:**

* Self-service health check tool available in Observability settings.

* Health check tests and reports: Perses running, Prometheus connectivity, authentication valid, RBAC permissions, sample query execution.

* Health check provides actionable remediation steps for each failure.

* "Copy diagnostics" button exports diagnostic bundle (logs, config, dashboard definitions, test results) for support tickets.

* Support runbooks created and reviewed for top 5 issue categories.

* Support team training completed (hands-on lab with all L1/L2 support engineers).

**Air-Gapped/Disconnected Environments:**

* Feature deploys successfully in cluster with egress blocked (no internet access).

* No 404 errors for external resources in browser console when deployed in disconnected environment.

* All fonts, icons, and assets render correctly without external CDN access.

* Dashboard rendering works with custom CA certificates for internal Prometheus HTTPS.

* Documentation includes air-gapped installation instructions with container image mirroring guide.

**Security:**

* RBAC penetration testing completed: Users cannot access unauthorized namespace metrics (attempts properly blocked).

* Query injection testing: Malicious PromQL queries blocked or sanitized by backend proxy.

* Audit logging: All dashboard views and queries logged with user identity for compliance.

* Security review sign-off from Red Hat Product Security team.

* No XSS vulnerabilities: All user input (dashboard names, descriptions, panel titles) properly sanitized.

**Documentation:**

* Quick Start Guide (3 pages, <15 min read) available covering viewing pre-built dashboards.

* Custom Dashboard Creation Tutorial (5 pages, <30 min to complete) with step-by-step instructions.

* Comprehensive Metrics Reference (15-20 pages) documenting all available metrics with examples.

* Troubleshooting Guide (10+ common scenarios) with decision trees and resolution steps.

* Architecture diagrams available: Component architecture, authentication flow, data flow, multi-tenancy model, deployment topology.

* At least 3 video tutorials published: (1) UI tour (5 min), (2) Creating custom dashboard (10 min), (3) Troubleshooting common issues (8 min).

* All video content has closed captions and transcripts.

* All documentation WCAG 2.1 AA compliant (screen reader compatible, proper heading structure, alt text for images).

* Documentation deliverables reviewed and approved by technical writing team.

* Documentation usability tested: 3 representative users per persona complete key tasks using only documentation (no external help).

**UX Consistency:**

* All UI components outside Perses-rendered charts use PatternFly v5 (Cards, Alerts, Modals, Forms, Empty States).

* Layout maintains ODH Dashboard navigation, header, and breadcrumbs (embedded content doesn't hide chrome).

* Design tokens follow PatternFly CSS variables (no hardcoded colors, spacing, or font sizes).

* UX review confirms consistency with ODH Dashboard interaction patterns (modals, toasts, confirmation flows).

* Success notifications use toast pattern (auto-dismiss after 5 seconds).

* Error notifications persist until user dismisses.

* Unsaved changes warning modal follows ODH pattern with "Save", "Discard", "Cancel" options.

**Telemetry:**

* Feature telemetry collection implemented (with user opt-in) tracking adoption, usage, performance, and errors.

* Telemetry respects privacy: No metric values, namespace names, model names, or user identities captured.

* Feature health monitoring alerts configured for SRE team (availability, performance, error rates).

**Pre-Launch Validation:**

* Beta program completed: 5-10 design partner customers used feature for 30+ days, feedback incorporated.

* Load testing passed: 500 concurrent users, 72-hour soak test, chaos engineering scenarios (Perses crash, Prometheus down, network partition).

* Platform Metrics prerequisites validated: All required metrics available and stable, Prometheus production-ready.

* Perses technical validation passed: Performance, security, and embedding approach meet requirements.

**Use Cases - i.e. User Experience & Workflow:**

**Use Case 1: Data Scientist Monitors Model Serving Performance**

*Main Success Scenario:*
1. Data Scientist logs into OpenShift AI Dashboard
2. Navigates to "Observability" in the left menu
3. Selects "Model Serving Overview" from the dashboard list
4. Views real-time metrics for their deployed model: inference latency p50/p95/p99, requests per second, error rate
5. Identifies a latency spike correlated with increased traffic
6. Uses dashboard time range selector to examine historical patterns
7. Shares dashboard URL with team members to collaborate on optimization

*Alternative Flow 1a - No Models Deployed:*
- Dashboard shows empty state with guidance on deploying models and what metrics will appear

*Alternative Flow 1b - Multi-Model Comparison:*
- Data Scientist uses dashboard variables to compare metrics across multiple model versions side-by-side

**Use Case 2: MLOps Engineer Troubleshoots Pipeline Failure**

*Main Success Scenario:*
1. MLOps Engineer receives alert about pipeline failure (from Platform Alerts)
2. Opens OpenShift AI Dashboard and navigates to "Observability"
3. Selects "Pipeline Execution Metrics" dashboard
4. Filters dashboard to specific pipeline and time range of failure
5. Observes elevated resource consumption and OOM errors in specific pipeline steps
6. Drills into trace data (if Tempo integration available) to see distributed trace of failed pipeline run
7. Identifies root cause and adjusts pipeline resource requests

*Alternative Flow 2a - Creating Custom Pipeline Dashboard:*
- Engineer clicks "Create Dashboard" to build custom dashboard for new pipeline
- Adds panels with specific Prometheus queries for custom metrics emitted by pipeline
- Saves dashboard with descriptive name and shares with pipeline development team

**Use Case 3: Platform Administrator Monitors Infrastructure Health**

*Main Success Scenario:*
1. Platform Administrator navigates to "Observability" in ODH Dashboard
2. Selects "Infrastructure Health" dashboard showing cluster-wide AI workload metrics
3. Views aggregated metrics across all namespaces: total GPU utilization, notebook pod counts, model serving capacity
4. Identifies overutilized nodes and plans capacity expansion
5. Exports dashboard definition to Git repository for version control
6. Updates dashboard to add new panel for recently deployed distributed training framework

**Use Case 4: Data Science Team Creates Shared Custom Dashboard**

*Main Success Scenario:*
1. Team lead navigates to "Observability" and clicks "Create Dashboard"
2. Uses Perses UI to add panels for team-specific KPIs (model accuracy metrics from custom exporters, business metrics)
3. Configures dashboard variables for model name and environment (dev/staging/prod)
4. Saves dashboard as "Team Model Performance Dashboard"
5. Shares dashboard with team members via built-in sharing mechanism
6. Team members can view dashboard but only team lead can edit
7. Dashboard definition is optionally exported to team's Git repository

**User Workflow Diagram (Conceptual):**

```
┌─────────────────────────────────────────────────┐
│         OpenShift AI Dashboard                  │
│  ┌─────────────────────────────────────────┐   │
│  │ Navigation                               │   │
│  │ • Projects                               │   │
│  │ • Applications                           │   │
│  │ • Model Serving                          │   │
│  │ ► Observability [NEW]                    │   │
│  │ • Resources                              │   │
│  └─────────────────────────────────────────┘   │
│                                                  │
│  ┌──────────────────────────────────────────┐  │
│  │ Perses Dashboard Interface (Embedded)    │  │
│  │                                           │  │
│  │ Dashboard List | Create Dashboard        │  │
│  │                                           │  │
│  │ ┌────────────────────────────────────┐  │  │
│  │ │ Model Serving Overview             │  │  │
│  │ │ ┌──────────┐ ┌──────────┐         │  │  │
│  │ │ │Latency   │ │Throughput│         │  │  │
│  │ │ │  Graph   │ │  Graph   │         │  │  │
│  │ │ └──────────┘ └──────────┘         │  │  │
│  │ │ ┌──────────┐ ┌──────────┐         │  │  │
│  │ │ │Error Rate│ │CPU Usage │         │  │  │
│  │ │ └──────────┘ └──────────┘         │  │  │
│  │ └────────────────────────────────────┘  │  │
│  └──────────────────────────────────────────┘  │
│                   ▲                             │
│                   │ Queries                     │
│                   ▼                             │
│  ┌──────────────────────────────────────────┐  │
│  │ Centralized Prometheus                   │  │
│  │ (Platform Metrics Architecture)          │  │
│  └──────────────────────────────────────────┘  │
└─────────────────────────────────────────────────┘
```

**Documentation Considerations:**

**Documentation Deliverables (MVP):**

| Document | Format | Audience | Length | Owner | Due Date |
|----------|--------|----------|--------|-------|----------|
| Observability Overview | HTML (docs site) | All users | 2 pages | Tech Writer | 2 weeks before GA |
| Quick Start Guide | HTML (docs site) | Data Scientists | 3 pages | Tech Writer | 2 weeks before GA |
| Creating Custom Dashboards Tutorial | HTML (docs site) | MLOps Engineers | 5 pages | Tech Writer | 2 weeks before GA |
| Metrics Reference | HTML (docs site) | All technical users | 15-20 pages | Engineer (generate) + Tech Writer (edit) | 1 week before GA |
| Troubleshooting Guide | HTML (docs site) | All users | 5-7 pages | Support Engineer + Tech Writer | 1 week before GA |
| Admin Integration Guide | HTML (docs site) | Platform Admins | 8-10 pages | Tech Writer | 2 weeks before GA |
| Architecture Diagrams | SVG/PNG + draw.io source | Admins, Architects | 5 diagrams | Tech Writer + Engineer | 2 weeks before GA |
| Video: UI Tour | Video (5 min) + Transcript | All users | 5 min | Video team + Tech Writer | 1 week before GA |
| Video: Creating Dashboard | Video (10 min) + Transcript | MLOps Engineers | 10 min | Video team + Tech Writer | 1 week before GA |
| Video: Troubleshooting | Video (8 min) + Transcript | All users | 8 min | Video team + Tech Writer | 1 week before GA |
| Release Notes | HTML (release notes site) | All users | 1-2 pages | PM + Tech Writer | At GA |

**Documentation Content Specifications:**

* **Observability Overview (2 pages):**
  - What is observability and why it matters for AI workloads
  - Quick introduction to the Observability feature and its purpose
  - Overview of pre-built dashboards and what they monitor
  - Link to Quick Start Guide for hands-on tutorial

* **Quick Start Guide by Persona (3 pages):**
  - **Data Scientist Quick Start:** "Monitor Your First Model in 5 Minutes"
  - **MLOps Engineer Quick Start:** "Troubleshoot Pipelines with Observability"
  - **Platform Admin Quick Start:** "Deploy and Configure Observability"
  - Each guide: Goal, Prerequisites, Step-by-step, Expected outcome, Next steps

* **Creating Custom Dashboards Tutorial (5 pages):**
  - Step-by-step: Create dashboard, add panels, write PromQL queries, save
  - Understanding Prometheus query language basics for AI metrics
  - Query builder/wizard for beginners (visual interface)
  - Adding and configuring panels (graphs, stat panels, tables)
  - Using variables for dynamic filtering (namespace, model name)
  - Dashboard layout and organization tips
  - Saving and sharing dashboards
  - Time to complete: <30 minutes

* **Comprehensive Metrics Reference (15-20 pages):**

  For each metric family, document:
  - Metric name and type (counter, gauge, histogram, summary)
  - Description: What does this metric measure?
  - Labels: All available label dimensions with examples
  - Unit of measurement
  - Typical values: What's normal vs. problematic?
  - Example queries: Common ways to use this metric
  - Related metrics: Other metrics to correlate with
  - Source: Which component exports this metric?

  Metric categories:
  - Model serving metrics (KServe, ModelMesh, TGIS)
  - Pipeline metrics (Kubeflow Pipelines, Tekton)
  - Notebook metrics (Jupyter, VS Code)
  - Distributed training metrics (PyTorch, TensorFlow operators)
  - Infrastructure metrics (GPU, CPU, memory, storage)

* **Troubleshooting Guide (5-7 pages, 10+ scenarios):**

  **Dashboard Issues:**
  - Dashboard shows "No data" or empty panels
    - Check: Prometheus data source configuration
    - Check: Time range selection (data may not exist for that period)
    - Check: Namespace/RBAC permissions
    - Check: Workload actually running and instrumented
    - Decision tree flowchart for diagnosis
  - Dashboard renders slowly or times out
    - Query optimization techniques
    - Reducing time range or panel count
    - Understanding cardinality impact
  - "Permission denied" or "Unauthorized" errors
    - RBAC troubleshooting flowchart
    - Namespace access verification steps
    - Request access procedures

  **Query and Metrics Issues:**
  - "Unknown metric" errors in custom queries
    - How to discover available metrics
    - Checking if workload is instrumented correctly
    - Metric naming conventions reference
  - Unexpected metric values or missing datapoints
    - Understanding scrape intervals and gaps
    - Prometheus retention policies
    - Metric staleness and resets

  **Performance Issues:**
  - Dashboard queries are slow
    - Query optimization best practices
    - Using recording rules for complex queries
    - Cardinality reduction techniques
  - Perses UI is unresponsive
    - Browser compatibility check
    - Network connectivity diagnostics
    - Perses component health verification

  **Integration Issues:**
  - Authentication/SSO failures
  - Prometheus connection errors
  - Tempo trace visualization not working (Non-MVP)

  **Self-Service Diagnostics:**
  - How to run health check tool
  - How to export diagnostic bundle for support tickets

* **Platform Admin Integration Guide (8-10 pages):**
  - Prerequisites: Platform Metrics architecture deployed
  - Installation and configuration steps
  - RBAC configuration for multi-tenancy
  - Prometheus data source configuration
  - Pre-built dashboard deployment
  - Backup and restore procedures
  - Upgrade procedures and version compatibility
  - Performance tuning and scaling
  - Monitoring the observability feature itself
  - Air-gapped deployment instructions

* **Architecture Diagrams (5 diagrams):**
  - **Component Architecture:** ODH Dashboard, embedded Perses, Prometheus, Tempo, data flow
  - **Authentication Flow:** User auth, token propagation, RBAC enforcement points
  - **Data Flow:** Metrics emission from workloads → Prometheus → Perses → User browser
  - **Multi-Tenancy Architecture:** Namespace isolation enforcement across components
  - **Deployment Topology:** Kubernetes resources, pods, services, ingress configuration

**Integration with Existing OpenShift AI Documentation:**

* **Navigation and Taxonomy:**
  - Create top-level section: "Monitoring and Observability"
  - Cross-reference from every feature section to relevant observability guidance
  - Bidirectional linking: Feature docs ↔ Observability docs

* **Inline Observability Guidance:**

  Update existing feature documentation with inline observability sections:

  Example for Model Serving Documentation:
  ```markdown
  ## Deploying a Model
  [Existing deployment steps]

  ## Monitoring Your Deployed Model (NEW)
  After deployment, monitor your model's performance:
  1. Navigate to Observability > Model Serving Overview
  2. Select your model from the dropdown
  3. Key metrics to watch: [link to metrics guide]
  4. [Screenshot of dashboard]

  For detailed monitoring guidance, see [Observability documentation].
  ```

  - Model Serving docs → "Monitoring Model Performance" section
  - Pipeline docs → "Tracking Pipeline Executions" section
  - Distributed Training docs → "Monitoring Training Jobs" section
  - Notebook docs → "Monitoring Notebook Resources" section

* **Searchability and Discovery:**
  - Ensure observability terms in doc search index
  - Add "related topics" suggestions to existing pages
  - "See also" sections linking observability and feature docs
  - Tag content appropriately for faceted search

* **Consistent Terminology:**
  - Shared glossary across all docs
  - Use "dashboard" not "panel set"
  - Use "Observability" not "Perses" in user-facing docs

**Accessibility and Localization (MVP):**

* **Accessibility (WCAG 2.1 AA Compliance):**
  - All documentation pages screen reader compatible
  - Proper heading hierarchy (h1 > h2 > h3)
  - All images have descriptive alt text
  - All video content has closed captions and transcripts
  - Links have descriptive text (not "click here")
  - Code examples have proper semantic markup
  - Color not the only means of conveying information
  - Documentation available in formats suitable for assistive technology (HTML)

* **Internationalization and Localization:**
  - **UI Localization:** Perses UI strings externalized for translation
  - **Initial Supported Languages:** English (primary), Spanish, French, German, Japanese, Chinese (Simplified)
  - **RTL Language Support:** Arabic, Hebrew (future consideration)
  - **Documentation Localization:**
    - Priority languages: English (primary), Spanish, Japanese, Chinese (Simplified)
    - Community contribution process for additional languages
    - Translation workflow and quality assurance
    - Keep translations in sync with English source
  - **Localization Quality:**
    - Cultural appropriateness review (colors, symbols, examples)
    - Number/date formatting respects locale
  - **Testing:**
    - Localization testing for each supported language
    - Usability testing with native speakers

**Post-MVP Documentation (Non-MVP):**

* **Dashboard-as-Code Guide:** Documentation for managing dashboards as code
  - YAML/JSON schema reference for Perses dashboards
  - GitOps workflow for dashboard deployment
  - CI/CD integration examples
  - Dashboard validation and testing
  - Version control best practices

* **Advanced PromQL Tutorial (10 pages):** Deep dive into Prometheus query language for AI workloads
  - Functions, operators, aggregations
  - Subqueries and recording rules
  - High-cardinality optimizations
  - Troubleshooting complex queries

* **Interactive Tutorial Environment:**
  - Sandbox with pre-deployed sample workloads
  - Guided exercises with instant feedback
  - "Try it yourself" embedded in docs

* **Dashboard Template Gallery:**
  - Searchable catalog of dashboard templates
  - User ratings and comments
  - "Used by similar teams" recommendations

* **Webinar Series:**
  - "Observability Best Practices for AI/ML Workloads" (Monthly)
  - "Customer Showcase: How [Company] Uses Observability" (Quarterly)

**Documentation Quality Assurance:**

* **Technical Accuracy Testing:**
  - Every procedure tested by someone who didn't write it
  - Test on clean environment matching target deployment
  - All code examples copied and executed (no typos)
  - All screenshots from current build (version-consistent)
  - All links verified (no 404s)

* **Usability Testing:**
  - Minimum 3 users per persona test documentation
  - Users complete key tasks using only documentation (no hints)
  - Measure: task completion rate, time, errors, satisfaction
  - Iterate docs based on testing feedback

* **Peer Review:**
  - Subject matter expert reviews technical content
  - Technical writer reviews clarity and structure
  - Product manager reviews alignment with feature goals
  - Accessibility specialist reviews compliance

* **Sign-off Criteria:**
  - All reviewers approve
  - All test users complete tasks successfully
  - Zero P0/P1 doc bugs open
  - Style guide compliance verified

**Documentation Maintenance:**

* **Ownership:** Technical Writing team owns documentation with SME input from Engineering
* **Update Process:** Documentation updated with each feature release
* **Review Cycle:** Quarterly review for accuracy, annual comprehensive review
* **Deprecation:** Clear deprecation warnings before removing features, maintain docs for N-2 versions
* **Feedback Loop:** "Was this helpful?" on every page, analytics reviewed monthly, high-impact updates prioritized

**Questions to Answer:**

**Architecture and Design (Must Answer Before Implementation):**

1. **Integration Approach:** What is the recommended technical approach for embedding Perses?
   - **Recommendation (from Architect agent):** iframe with backend proxy for MVP (fastest time to value), evolve to Module Federation in 18 months
   - **Option A (iframe):** Simplest integration, full isolation, but limited styling control and accessibility challenges (keyboard focus management, screen reader navigation)
   - **Option B (React Components):** Best UX but Perses doesn't currently expose stable React component library (would require 6-12 month upstream contribution)
   - **Option C (Module Federation):** Ideal long-term but requires Perses federated module support (6-9 month timeline)
   - **Decision needed:** 2-week technical spike to validate iframe approach feasibility (accessibility, performance, embedding)
   - **Mitigation:** Design token mapping (Perses → PatternFly), careful keyboard navigation handling across iframe boundary

2. **Authentication and Authorization Flow:** How does the embedded Perses instance authenticate to Prometheus and enforce RBAC?
   - **Recommendation (from Architect agent):** Token-mediated backend proxy with query rewriting engine
   - **Architecture:**
     ```
     User Browser → ODH Dashboard Frontend (Bearer token)
       ↓
     ODH Dashboard Backend (validate JWT, extract namespaces, inject filters)
       ↓
     Perses Backend (service account token + namespace context)
       ↓
     Prometheus (scoped queries: namespace=~"ns1|ns2|ns3")
     ```
   - **RBAC Enforcement:** Query rewriting at backend proxy level + Prometheus datasource config (defense in depth)
   - **Platform Admin Exception:** Users with ClusterRole bypass namespace filters for cluster-wide visibility
   - **Decision needed:** Detailed architecture design document for auth/authz flow
   - **Security Review:** RBAC penetration testing mandatory before GA

3. **Dashboard Storage:** Where are Perses dashboard definitions stored?
   - **Recommendation (from Architect agent):** Kubernetes CRDs (namespace-scoped) with database escape hatch for scale
   - **Phase 1 (MVP):** CRDs for GitOps compatibility, multi-tenancy, backup/restore with standard K8s tools
   - **Advantages:** Native RBAC, GitOps-ready, no new infrastructure, audit via K8s audit logs
   - **Limitations:** etcd 1MB size limit (validation webhook enforces), performance at 1000+ dashboards per namespace
   - **Phase 2 (Post-MVP):** Hybrid CRD + PostgreSQL if scale demands (hot: CRDs for active dashboards, cold: DB for archived)
   - **Cross-Namespace Sharing:** Implement DashboardTemplate/DashboardRef pattern for shared templates
   - **Decision needed:** Dashboard schema version strategy (v1alpha1 → v1beta1 → v1)

4. **Pre-Built Dashboard Maintenance:** Who owns and maintains the default dashboards?
   - **Ownership:** Engineering team creates, Product team approves, packaged as ConfigMaps with ODH releases
   - **User Editability:** Pre-built dashboards are read-only templates; users "duplicate" to customize (prevents accidental modification)
   - **Update Strategy:** Versioned (e.g., model-serving-v2), users notified of new versions, opt-in upgrade
   - **Metrics Schema Changes:** Recording rules provide abstraction layer, dashboards reference stable recording rule names not raw metrics
   - **Decision needed:** Dashboard versioning and update notification mechanism

5. **Scalability and Performance:** What are the performance characteristics and limits?
   - **Requirements (from Architect agent):**
     - 500+ concurrent viewers
     - <5s render for 10-panel dashboard (p95)
     - <15s render for 50-panel dashboard (p95)
   - **Caching Strategy:** 3-layer cache (browser 1h, backend Redis 30s, Prometheus 1min)
   - **Testing Plan:** Load test (500 users), stress test (breaking point), soak test (72h), chaos engineering
   - **Query Optimization:** Enforce query complexity limits, provide query optimization guide, warn users about expensive queries
   - **Decision needed:** Performance test environment and criteria for GO/NO-GO

6. **Perses Version and Maturity:** What version of Perses should be used?
   - **Risk Assessment (from PXE agent):** CNCF Sandbox = 62% attrition rate, 20% probability of abandonment over 3 years
   - **Mitigation Strategy:**
     - Pin specific Perses version, thorough testing before adopting new versions
     - Establish relationship with Perses maintainers, Red Hat engineer becomes reviewer/maintainer
     - Contribute AI/ML features upstream (model performance panels, pipeline visualizations)
     - **Contingency Plan:** Fork strategy documented, Grafana migration path designed
   - **Upstream Engagement:** Join Perses community meetings, propose embedding API improvements, influence roadmap
   - **Decision needed:** Which Perses version to pin for MVP, upgrade strategy, fork criteria

7. **Backward Compatibility and Migration:** How does this feature interact with existing observability tools?
   - **Grafana Coexistence:** Perses complements (not replaces) Grafana initially; clarify positioning in messaging
   - **Migration Path:** No automated Grafana → Perses conversion for MVP; provide manual migration guide for common patterns
   - **OpenShift Monitoring:** Observability feature focuses on AI workloads, OpenShift Console remains for cluster monitoring
   - **Messaging:** "Perses brings observability into your AI workflow; use OpenShift Console for cluster admin tasks"
   - **Decision needed:** Clear documentation on when to use Observability vs OpenShift Console vs external Grafana

8. **Multi-Cluster Considerations:** How does this work in multi-cluster OpenShift AI deployments?
   - **Out of Scope for MVP:** Multi-cluster support deferred to future RFE
   - **Single Cluster:** Each cluster has own Perses instance, no cross-cluster dashboard aggregation
   - **Future Architecture:** Thanos for multi-cluster Prometheus federation + centralized Perses instance
   - **Decision needed:** Document multi-cluster as future work, validate single-cluster approach sufficient for MVP customers

9. **Plugin Ecosystem:** Does ODH Dashboard leverage Perses plugin capabilities?
   - **MVP:** Use built-in Perses panel types (graph, stat, table, heatmap)
   - **Post-MVP:** Explore AI-specific visualizations:
     - Confusion matrix panel
     - ROC curve panel
     - Model drift detection visualization
     - Training loss curves with epoch markers
   - **Upstream Contribution:** Contribute AI/ML panels to Perses community if developed
   - **Decision needed:** Assess demand for custom panels during beta, prioritize based on customer feedback

10. **Documentation Tooling and Workflow:** How are docs built, reviewed, and published?
    - **Tooling:** Existing OpenShift AI docs infrastructure (Asciidoc, Hugo, or current toolchain)
    - **Workflow:** Tech writer drafts → SME review → Usability testing → Editorial review → Publish
    - **Review Process:** PRs for all doc changes, staged preview environment, accessibility scan in CI/CD
    - **Code Example Testing:** All code examples executed in test environment, screenshots automated where possible
    - **Decision needed:** Assign technical writer to project, establish doc review SLA

11. **Support Model and Escalation:** How is support structured for this feature?
    - **Level 1 Support:** Use runbooks for common issues (no data, slow dashboard, permission denied)
    - **Level 2 Support:** Deeper diagnostics, use health check tool output, engage engineering if needed
    - **Level 3 Engineering:** Complex issues, potential Perses bugs, integration problems
    - **Upstream Perses:** If Perses bug identified, escalate to Perses community with reproducer
    - **SLA:** Align with standard Red Hat support SLAs for OpenShift AI, clarify observability feature criticality
    - **Decision needed:** Support team training plan, runbook completion date

**Background & Strategic Fit:**

**Strategic Context:**

OpenShift AI aims to be a comprehensive, enterprise-grade AI/ML platform. Observability is a critical capability gap that impacts user experience, operational excellence, and competitive positioning. This feature addresses that gap by integrating best-in-class observability tooling directly into the platform.

**Why Perses:**

* **CNCF Project:** As a CNCF sandbox project, Perses aligns with Red Hat's open source and cloud-native strategy. Contributing to and building on CNCF projects strengthens the broader ecosystem.

* **Modern Architecture:** Perses was designed with modern observability needs in mind: Dashboard-as-Code, Kubernetes-native, plugin extensibility, and GitOps workflows align perfectly with how AI/ML teams operate today.

* **Prometheus-First:** Perses has deep Prometheus integration, which is the standard metrics system in the Kubernetes ecosystem and the foundation of the Platform Metrics architecture.

* **Embeddability:** Perses is designed to be embedded in other applications, unlike Grafana which is primarily deployed as a standalone service. This enables the integrated UX this RFE targets.

* **Innovation Opportunity:** OpenShift AI can influence Perses roadmap to ensure AI/ML use cases are first-class, potentially contributing features like model performance visualizations, pipeline flow diagrams, or experiment tracking integrations.

**Dependency on Platform Metrics Architecture:**

This feature is a consumer of the Platform Metrics and Alerts architecture, which establishes:
- Centralized Prometheus for metrics collection
- Instrumentation of AI workload components (KServe, Kubeflow, notebooks, training operators)
- Tempo for distributed tracing
- Standardized metrics schemas and labels

Perses integration must wait until Platform Metrics provides stable, comprehensive metrics from all relevant AI workload components. This RFE assumes that dependency is met or in progress.

**Competitive Landscape:**

* **AWS SageMaker:** Offers integrated CloudWatch dashboards for model monitoring and resource utilization
* **Azure ML:** Provides built-in metrics and monitoring dashboards within the Azure ML Studio
* **Google Vertex AI:** Features integrated observability with Cloud Monitoring dashboards
* **Databricks:** Has native monitoring dashboards for models and infrastructure

All major cloud AI platforms provide integrated observability. OpenShift AI must match this capability to be competitive, especially for enterprises evaluating platforms.

**Community and Ecosystem:**

* **CNCF Engagement:** This feature strengthens Red Hat's engagement with CNCF observability projects
* **AI/ML Observability Standards:** Opportunity to shape emerging standards for AI workload observability through Perses community
* **Partner Ecosystem:** Enables partners building on OpenShift AI to leverage standardized observability without building custom solutions

**Roadmap Alignment:**

* **Near-term:** This feature is a critical enabler for production AI workload adoption
* **Medium-term:** Sets foundation for advanced capabilities like automated anomaly detection, SLO tracking, and cost optimization dashboards
* **Long-term:** Positions OpenShift AI as an observability-first platform where telemetry drives intelligent automation, auto-scaling, and self-healing capabilities

**Customer Considerations:**

**Enterprise Requirements:**

* **Multi-Tenancy and Isolation:** Large enterprises run multiple data science teams on shared infrastructure. Dashboard access must be strictly isolated to prevent data leakage between teams or projects. The integration must honor OpenShift's existing RBAC and namespace isolation mechanisms.

* **Compliance and Audit:** Regulated industries (finance, healthcare) require audit trails of who accessed what metrics and when. Consider whether dashboard access should be logged and auditable.

* **Data Residency and Privacy:** Metrics may contain sensitive information (model names, customer IDs in labels). Ensure metrics data never leaves the cluster and storage complies with data residency requirements.

* **SSO and Identity Integration:** Enterprises use corporate identity providers (LDAP, Active Directory, SAML, OIDC). The integration must seamlessly work with ODH Dashboard's existing authentication without requiring separate Perses login.

**Scale Considerations:**

* **Large-Scale Deployments:** Some customers run hundreds of concurrent notebooks, dozens of model serving endpoints, and continuous pipeline executions. Dashboard queries must remain performant even with high-cardinality metrics.

* **Many Users:** Platforms may have 100+ data scientists. Ensure the integration can support concurrent dashboard viewing and editing without performance degradation.

* **Long Retention:** Customers often need to analyze historical trends over weeks or months. Ensure Prometheus retention policies and dashboard query performance accommodate long time ranges.

**Operational Considerations:**

* **Upgrade and Maintenance:** Minimize operational burden of maintaining Perses alongside OpenShift AI. Seamless upgrades, minimal configuration, and self-healing are critical.

* **Resource Overhead:** Be transparent about the infrastructure resources (CPU, memory, storage) required for the Perses integration. Customers need to plan capacity accordingly.

* **High Availability:** Production AI workloads require 24/7 observability. The integration should support HA deployment of any Perses components.

**Customization and Flexibility:**

* **Industry-Specific Dashboards:** Financial services, healthcare, retail, and other verticals have unique AI use cases. Provide extensibility for customers or partners to create industry-specific dashboard templates.

* **Custom Metrics:** Customers often instrument their models with custom metrics (business KPIs, domain-specific measures). Ensure the integration supports custom Prometheus metrics seamlessly.

* **Branding and Whitelabeling:** Some enterprise deployments require co-branding or whitelabeling. Consider whether the embedded Perses UI can be styled to match customer branding.

**Training and Adoption:**

* **User Skill Levels:** Data Scientists may not be familiar with Prometheus query language or observability best practices. Provide intuitive interfaces, query builders, and extensive documentation to lower the barrier to entry.

* **Change Management:** Introducing new observability tooling requires user training and change management. Provide resources (webinars, workshops, best practices guides) to drive adoption.

**Support and SLA:**

* **Support Model:** Clarify support boundaries between Red Hat support for ODH Dashboard integration and upstream Perses community support for Perses functionality.

* **SLA Requirements:** Understand if observability dashboards are considered critical infrastructure (impacting SLAs) or auxiliary tooling. This impacts testing rigor, HA requirements, and support priority.

**Migration from Existing Tools:**

* **Grafana Users:** Many customers already use Grafana for observability. Provide clear guidance on whether Perses replaces Grafana, coexists with it, or migrates from it. Offer migration paths or conversion tools if feasible.

* **OpenShift Monitoring Integration:** Some customers already use OpenShift's built-in monitoring stack. Clarify how this feature relates to and integrates with existing OpenShift observability capabilities.

**Air-Gapped and Restricted Environments:**

* **Disconnected Deployments:** Government, defense, and high-security customers run in air-gapped environments. Ensure all Perses components, UI assets, and dependencies are available for offline installation.

* **Restricted Networks:** Some environments have strict egress filtering. Verify no runtime dependencies on external services (CDNs, external APIs, telemetry beacons).

**Cost and Licensing:**

* **Open Source Licensing:** Perses is Apache 2.0 licensed, which is compatible with enterprise use. Ensure all dependencies have acceptable licenses.

* **Operational Costs:** Be transparent about the infrastructure costs of running the observability stack (Prometheus, Tempo, Perses) at scale. Provide sizing guidelines to help customers estimate costs.

**Success Metrics:**

**Product Metrics (Track via Telemetry):**

* **Adoption Rate:** Percentage of OpenShift AI users who access the Observability feature within 30 days of deployment
  - **Target:** >60% of active users
  - **Leading Indicator:** % who access within 7 days (target: >30%)

* **Usage Frequency:** Average number of dashboard views per user per week
  - **Target:** ≥3 views/user/week for active users
  - **Engagement Levels:** Curious (<1/week), Active (1-5/week), Power User (>5/week)

* **Dashboard Stickiness:** % of users who viewed dashboards 3+ times in 30 days
  - **Target:** >50% (indicates real value beyond curiosity)

* **Custom Dashboard Creation:** Number of custom dashboards created by users
  - **Target:** >30% of active users create at least one custom dashboard within 90 days
  - **Indicates:** Engagement beyond defaults, self-service capability

* **Pre-Built vs Custom Ratio:** Ratio of pre-built dashboard views to custom dashboard views
  - **Baseline:** Expect 80/20 initially (validates default dashboard quality)
  - **Maturity:** Over time, 60/40 as users create custom dashboards

* **Time to First Value (TTFV):** Time from enabling Observability to first custom dashboard creation
  - **Target:** <15 minutes for guided users, <45 minutes for unguided

* **Dashboard Load Performance:** p95 dashboard render time
  - **Target:** <5 seconds for 10-panel dashboards
  - **Monitor:** Alert if >10 seconds (poor performance kills adoption)

* **Error/Empty State Rate:** % of dashboard loads that show "no data" or errors
  - **Target:** <15% (indicates configuration/permission issues)
  - **Action Trigger:** >25% = investigate RBAC or Platform Metrics issues

**Business Impact Metrics:**

* **Revenue Protection:** Track deals where observability was discussed as requirement
  - **Win Rate:** % win rate when customer RFP includes integrated monitoring requirement
  - **Target:** Match or exceed AWS/Azure/GCP win rates in deals with observability criteria
  - **Risk Mitigation:** Identify "$X million ARR at risk" without this feature

* **Time to Resolution (MTTR):** Reduction in mean time to resolve model serving or pipeline issues
  - **Baseline:** Current MTTR via kubectl/logs (estimate: 2-4 hours)
  - **Target:** 60% reduction (45-90 minutes with dashboards)
  - **Measurement:** Support ticket timestamps, customer interviews

* **Support Ticket Reduction:** Decrease in support tickets related to "Why is my model slow?" or "Why did my pipeline fail?"
  - **Baseline:** Count of observability-related tickets pre-launch
  - **Target:** 40% reduction within 6 months post-GA
  - **Indicator:** Users self-service answers via dashboards

* **Customer Satisfaction:** NPS or CSAT scores specifically for observability capabilities
  - **Target:** NPS >40, CSAT >4.0/5.0
  - **Survey:** Quarterly survey of users who accessed Observability

* **Expansion Revenue:** Increase in node/GPU consumption post-observability launch
  - **Hypothesis:** Better visibility → more trust → more usage → more infrastructure purchased
  - **Measurement:** Track resource consumption changes for customers who adopt observability

**Competitive Metrics:**

* **Feature Parity Perception:** % of customers who rate our observability as "equivalent or better" than AWS/Azure/GCP
  - **Target:** >70% rate as equivalent or better
  - **Measurement:** Annual competitive survey

* **Competitive Win/Loss Analysis:** Track deals where observability was discussed
  - **Tag Deals:** "Observability requirement" tag in CRM
  - **Track Outcomes:** Win/loss rate, reasons for wins/losses
  - **Competitive Intelligence:** What observability features competitors highlight

**Leading Indicators (Predict Success Before Lagging Metrics):**

* **Week 1 Adoption:** % of users who access Observability within 7 days
  - **Target:** >30% (early indicator of awareness/discoverability)

* **Dashboard Interaction Rate:** % of dashboard viewers who interact (change time range, adjust filters)
  - **Target:** >60% (indicates engagement, not passive viewing)

* **Documentation Engagement:** % of users who access observability docs before or after using feature
  - **Target:** >40% access docs (indicates self-service success)

* **Beta Customer Satisfaction:** NPS from 5-10 beta customers
  - **Target:** NPS >50 before GA (validates product-market fit)

**Anti-Metrics (Watch for Warning Signs):**

* **Dashboard Abandonment:** Users who viewed dashboards once and never returned
  - **Warning Threshold:** >50% abandonment = investigate usability issues

* **Support Ticket Increase:** Spike in tickets about dashboard issues
  - **Warning Threshold:** >10% of total tickets = usability or reliability problem

* **Performance Degradation:** Dashboard load time p95 trends upward
  - **Warning Threshold:** Consistent >10 seconds = investigate performance issues

**Success Criteria for Launch:**

* **Beta Program:** 5-10 customers use feature for 30+ days, average NPS >40
* **Adoption:** >50% of beta customers access Observability ≥3 times/week
* **Custom Dashboards:** ≥30% of beta users create at least one custom dashboard
* **Support Load:** <5% of beta customer support tickets related to Observability issues
* **Performance:** All beta customers meet <5s dashboard render time p95

---

## Competitive Analysis

**Competitive Feature Matrix:**

| Capability | OpenShift AI (Today) | OpenShift AI (with Perses) | AWS SageMaker | Azure ML | Google Vertex AI | Databricks |
|------------|----------------------|-----------------------------|---------------|----------|------------------|------------|
| **Integrated Dashboards** | ❌ No | ✅ Yes | ✅ Yes (CloudWatch) | ✅ Yes (Azure Monitor) | ✅ Yes (Cloud Monitoring) | ✅ Yes (native) |
| **Custom Metrics Support** | ⚠️ Limited | ✅ Yes (Prometheus) | ✅ Yes | ✅ Yes | ✅ Yes | ✅ Yes |
| **Dashboard-as-Code** | ❌ No | ✅ Yes (YAML/GitOps) | ⚠️ Limited (CloudFormation) | ❌ No | ❌ No | ⚠️ Limited |
| **Model Drift Detection** | ❌ No | ❌ No (Future) | ✅ Yes (Model Monitor) | ✅ Yes (Data Collector) | ✅ Yes | ✅ Yes |
| **Distributed Tracing** | ❌ No | ⚠️ Non-MVP (Tempo) | ❌ No | ⚠️ Limited | ✅ Yes (Cloud Trace) | ❌ No |
| **Multi-Cloud Portable** | ✅ Yes (K8s) | ✅ Yes (K8s + Perses) | ❌ No (AWS-only) | ❌ No (Azure-only) | ❌ No (GCP-only) | ⚠️ Partial |
| **Open Source** | ✅ Yes | ✅ Yes | ❌ No | ❌ No | ❌ No | ⚠️ Partial |
| **Air-Gapped Support** | ✅ Yes | ✅ Yes | ❌ No | ❌ No | ❌ No | ⚠️ Limited |
| **RBAC/Multi-Tenancy** | ✅ Yes (K8s RBAC) | ✅ Yes (K8s RBAC) | ✅ Yes (IAM) | ✅ Yes (Azure AD) | ✅ Yes (IAM) | ✅ Yes |
| **Cost Visibility** | ❌ No | ⚠️ Partial (resource metrics) | ✅ Yes (Cost Explorer) | ✅ Yes (Cost Analysis) | ✅ Yes (Billing) | ✅ Yes |

**Differentiation Strategy:**

**Our Unique Value Props:**
1. **"The only AI platform with GitOps-native observability"** - Dashboard-as-Code with K8s CRDs enables version control, CI/CD, declarative management
2. **"Multi-cloud observability portability"** - Perses dashboards work across any Kubernetes cluster, no vendor lock-in to CloudWatch/Azure Monitor
3. **"Open source observability without vendor lock-in"** - CNCF Perses + Prometheus = community-driven, extensible, no proprietary metrics formats
4. **"Enterprise-grade for regulated industries"** - Air-gapped support, strict RBAC, audit logging = government/finance/healthcare ready

**Gaps vs Competitors (Acknowledge and Address):**

* **AWS SageMaker Model Monitor:** Has automatic drift detection, data quality monitoring, bias detection
  - **Our Response:** Roadmap item post-MVP; focus on "we enable you to build this with custom metrics" flexibility
  - **Positioning:** "SageMaker monitors what AWS decides; OpenShift AI lets you monitor what matters to YOUR business"

* **Azure ML Model Data Collector:** Automatic payload logging and analysis
  - **Our Response:** "We provide the platform, you control the data" - better for privacy-sensitive use cases
  - **Positioning:** Flexibility vs. automatic-but-opaque monitoring

* **Cost Visibility:** Competitors integrate billing/cost dashboards deeply
  - **Our Response:** Show resource consumption metrics, customers can correlate with their billing systems
  - **Post-MVP:** FinOps dashboard templates with cost estimation formulas

**Competitive Win Strategy:**

**When to win with Observability:**
- **Enterprise with mature MLOps:** Highlight Dashboard-as-Code, GitOps workflows, K8s-native approach
- **Multi-cloud/hybrid deployments:** Emphasize portability vs. cloud provider lock-in
- **Regulated industries:** Air-gapped support, audit logging, open source transparency
- **Teams with custom metrics:** Prometheus flexibility vs. pre-defined metrics in competitors

**Positioning Statements (for Sales):**
- **vs AWS:** "Don't let AWS dictate what you monitor. OpenShift AI gives you full control with open source Prometheus and GitOps-native dashboards."
- **vs Azure:** "Escape Azure vendor lock-in. Our observability works anywhere Kubernetes runs, with full air-gapped support for regulated environments."
- **vs GCP:** "Google's observability is great for Google Cloud. Ours works across all your clouds and on-prem, with enterprise RBAC and compliance built-in."
- **vs Databricks:** "Go beyond infrastructure metrics. OpenShift AI observability covers the full ML lifecycle from notebooks to production models."

**Sales Enablement Needed:**
- Competitive battlecards with feature comparisons
- Demo environment showcasing Dashboard-as-Code and GitOps workflows
- ROI calculator: "Save $X/year vs AWS CloudWatch costs" + "Reduce MTTR by 60%"
- Customer case studies (post-GA): "How [FinServ Co] uses OpenShift AI Observability in air-gapped environment"

---

## Risk Assessment and Mitigation

**Critical Risks (High Impact, Medium-High Probability):**

### Risk #1: Perses Project Maturity (CNCF Sandbox)
- **Impact:** HIGH - Major rework or migration required if Perses stalls/abandoned
- **Probability:** MEDIUM (20% over 3 years based on CNCF Sandbox attrition rates)
- **Blast Radius:** ALL customers using Observability feature
- **Indicators:** Commit velocity drops, maintainers leave, no releases for 6+ months

**Mitigation Strategy:**
- Pin specific Perses version with thorough testing before adopting newer versions
- Establish relationship with Perses maintainers (Red Hat engineer becomes reviewer/maintainer)
- Contribute to Perses project to influence roadmap and build goodwill
- **Contingency Plan A:** Fork Perses, maintain Red Hat-supported version (cost: 1-2 FTE ongoing)
- **Contingency Plan B:** Migrate to Grafana (cost: 3-6 months migration effort, lose Dashboard-as-Code benefits)
- **Contingency Plan C:** Build custom dashboard UI on Prometheus API (cost: 6-12 months development)
- **Decision Trigger:** If Perses shows abandonment signs, initiate contingency planning immediately

### Risk #2: Dashboard Performance at Scale
- **Impact:** HIGH - Slow dashboards = feature abandonment = customer dissatisfaction
- **Probability:** HIGH - Complex PromQL queries + high-cardinality metrics = known performance issue
- **Field Evidence:** Every observability feature struggles with this initially

**Mitigation Strategy:**
- Performance testing BEFORE GA: Simulate 500 concurrent users, 50+ models, 1000+ metrics
- Define and enforce query complexity limits (backend validates queries before execution)
- Implement aggressive caching: Browser (1h) → Backend Redis (30s) → Prometheus (1min)
- Provide query optimization guidance in documentation (recording rules, label filters)
- Monitor dashboard performance in production, proactive customer outreach if p95 >10s
- **Escalation:** If performance issues found in beta, delay GA until resolved

### Risk #3: RBAC Enforcement Gaps
- **Impact:** CRITICAL - Data leakage between tenants = security incident, compliance violation, customer trust loss
- **Probability:** MEDIUM - Multi-tenancy is complex, easy to misconfigure
- **Compliance Impact:** Regulated industries (finance, healthcare) = showstopper if RBAC broken

**Mitigation Strategy:**
- Security review with Red Hat Product Security team before beta
- Penetration testing: Dedicated security engineer attempts to access unauthorized namespace metrics
- Automated tests: RBAC enforcement across various scenarios (100+ test cases)
- Query rewriting engine with fail-safe defaults (deny unless explicitly authorized)
- Audit logging: Track who accessed what metrics for forensic analysis
- **Pre-GA Requirement:** Zero critical security findings, sign-off from security team

### Risk #4: Platform Metrics Not Ready
- **Impact:** BLOCKS MVP - Cannot build observability without metrics to observe
- **Probability:** MEDIUM - Dependency on parallel team/effort
- **Timeline Impact:** Could delay RFE by months if Platform Metrics blocked

**Mitigation Strategy:**
- **MANDATORY 2-week validation spike BEFORE design/implementation begins**
- Establish GO/NO-GO decision gate after validation
- Document required metrics clearly for Platform Metrics team (requirements in this RFE)
- If metrics <50% complete, DEFER this RFE until dependency ready
- Weekly sync with Platform Metrics team during development
- **Fallback:** Phase 1 with subset of metrics (model serving only), expand later

### Risk #5: Upgrade Path / Dashboard Loss
- **Impact:** HIGH - Dashboard loss during upgrade = angry customers, trust erosion, churn risk
- **Probability:** MEDIUM - First few upgrades always have issues
- **Customer Impact:** Hours of work creating dashboards lost = severe

**Mitigation Strategy:**
- Automated dashboard backup before every ODH upgrade (mandatory, not optional)
- Upgrade testing in staging with real customer-like dashboards before production rollout
- Dashboard schema versioning (v1alpha1 → v1beta1 → v1) with migration tools
- Rollback plan documented and tested (can revert to previous ODH version with dashboards intact)
- Clear communication: "Back up your dashboards before upgrading" even with automated backup
- **Insurance:** Dashboard export/import feature for manual backup

**Medium Risks (Medium Impact, Medium Probability):**

### Risk #6: Documentation Quality and Completeness
- **Impact:** MEDIUM - Poor docs = high support costs, low adoption, frustrated users
- **Probability:** MEDIUM - Docs often written last, rushed before release

**Mitigation Strategy:**
- Allocate dedicated technical writer from project start
- User testing: Watch 3 real users per persona try to use feature with only docs (no help)
- Iterative improvement: Monitor support tickets, update docs for common issues
- Video tutorials supplement written docs (learning style diversity)
- **Quality Gate:** Docs must pass usability testing before GA

### Risk #7: No Data / Empty State User Confusion
- **Impact:** MEDIUM-HIGH - Poor first impression = low adoption, high support costs
- **Probability:** HIGH - Based on field experience with similar features

**Mitigation Strategy:**
- Smart empty states with actionable guidance (already in requirements)
- Self-service health check tool (already in requirements)
- Onboarding tutorial for first-time users (already in requirements)
- **Validation:** Usability testing with users who have no prior observability experience

### Risk #8: Air-Gapped Deployment Failure
- **Impact:** MEDIUM-HIGH - Government, defense, high-security customers blocked
- **Probability:** LOW-MEDIUM - If not designed for from start, hard to retrofit

**Mitigation Strategy:**
- Test air-gapped deployment during development (already in requirements)
- Ensure all Perses UI assets bundled (no CDN dependencies) - already in requirements
- Automated test in CI/CD: Cluster with egress blocked
- **Validation:** At least 1 beta customer tests in air-gapped environment before GA

**Risk Summary Matrix:**

| Risk | Impact | Probability | Mitigation Status | Owner |
|------|--------|-------------|-------------------|-------|
| Perses Maturity | HIGH | MEDIUM | Contingency plans defined | Architect |
| Performance at Scale | HIGH | HIGH | Testing plan required | Engineering |
| RBAC Gaps | CRITICAL | MEDIUM | Security review mandatory | Security + Engineering |
| Platform Metrics Blocked | BLOCKS MVP | MEDIUM | Validation spike required | Product + Architect |
| Dashboard Loss on Upgrade | HIGH | MEDIUM | Backup automation required | Engineering |
| Documentation Quality | MEDIUM | MEDIUM | Dedicated writer assigned | Tech Writing |
| Empty State Confusion | MEDIUM-HIGH | HIGH | Requirements defined | UX + Engineering |
| Air-Gapped Failure | MEDIUM-HIGH | LOW-MEDIUM | Testing plan defined | Engineering |

---

## GO/NO-GO Decision Criteria

**4-Week Validation Phase (Before Design/Implementation Begins):**

After completing 4-week validation phase, evaluate against these criteria to make GO/NO-GO decision.

### **GO Criteria (ALL Must Pass):**

**Platform Metrics Readiness:**
- ✅ Centralized Prometheus deployed with HA configuration (≥2 replicas)
- ✅ Prometheus has ≥7 day retention and can handle 1M+ active time series
- ✅ Required metrics available and stable for ≥80% of planned pre-built dashboards:
  - Model serving metrics (KServe/ModelMesh): request_duration, request_count, errors
  - Pipeline metrics: run_duration, run_status, task metrics
  - Notebook metrics: CPU, memory, GPU utilization
- ✅ Prometheus accessible from ODH Dashboard backend pod (network policies allow)
- ✅ Service account created with read permissions to Prometheus API
- **Decision:** Platform Metrics team confirms production-readiness

**Perses Technical Validation:**
- ✅ Perses deployed successfully in test environment
- ✅ Load testing: Handles 100+ concurrent users without degradation
- ✅ Performance: Renders 50-panel dashboard in <15 seconds (p95)
- ✅ Security scan: Zero critical vulnerabilities found
- ✅ Embedding approach (iframe/React/Federation) technically validated and feasible
- ✅ Accessibility baseline: Perses UI meets minimum WCAG 2.1 AA requirements or gaps identified with mitigation plan
- **Decision:** Engineering team confirms technical feasibility

**Business Case:**
- ✅ Customer discovery completed: 10 customers interviewed, pain points validated with data
- ✅ Beta customer commitment: ≥5 customers committed to beta program (signed agreements)
- ✅ Competitive analysis: Feature gaps identified and positioning strategy defined
- ✅ ROI justifies investment: Revenue protection, MTTR reduction, support cost savings quantified
- **Decision:** Product Management and leadership approve business case

**Project Readiness:**
- ✅ Architecture design complete: Auth/authz flow, storage, multi-tenancy, caching strategy documented
- ✅ UX wireframes complete: Landing page, empty states, first-run experience, error scenarios designed
- ✅ Documentation plan: Deliverables table, owners, timelines, quality assurance process defined
- ✅ Technical writer assigned to project
- ✅ Support team engaged and training plan scheduled
- **Decision:** Cross-functional team (Product, Engineering, UX, Docs, Support) aligned and ready

**Perses Community Health:**
- ✅ Perses project shows active development (>10 commits/month over last 3 months)
- ✅ Perses maintainers engaged: Red Hat engineer introduced, communication channel established
- ✅ Perses roadmap reviewed: No major breaking changes planned that would impact integration
- **Decision:** Perses dependency risk acceptable

### **NO-GO Criteria (ANY Triggers Re-Evaluation):**

**Platform Metrics Insufficient:**
- ❌ Platform Metrics <50% complete (insufficient metrics for meaningful dashboards)
- ❌ Prometheus not production-ready (single replica, <3 day retention, <100K time series capacity)
- ❌ Metrics schemas unstable (changed 3+ times in last month)
- **Action:** DEFER RFE until Platform Metrics ready; revisit in 3 months

**Perses Technical Failure:**
- ❌ Perses performance unacceptable (<20 concurrent users causes issues, >30s render times)
- ❌ Critical security vulnerabilities found in Perses with no patch available
- ❌ Embedding approach technically infeasible (accessibility blockers, cannot integrate with ODH Dashboard)
- **Action:** Pivot to alternative (Grafana) or build custom solution; re-scope RFE

**Perses Project Health Concerns:**
- ❌ Perses shows signs of abandonment (no commits in 3 months, maintainers left)
- ❌ Perses community unresponsive (no response to engagement attempts)
- ❌ Perses roadmap conflicts with our needs (breaking changes planned, no embedding support)
- **Action:** Evaluate fork strategy or alternative dashboard solutions

**Business Case Fails:**
- ❌ Customer discovery reveals low demand (customers don't prioritize integrated observability)
- ❌ No beta customers willing to participate
- ❌ Competitive analysis shows feature parity not critical for wins
- **Action:** Deprioritize or cancel RFE; focus resources on higher-value features

**Resource Constraints:**
- ❌ Required engineering resources not available (team overcommitted)
- ❌ Technical writer not available for parallel documentation effort
- ❌ UX design resources not available for wireframes and usability testing
- **Action:** Defer RFE until resources available or secure additional resources

### **Decision Process:**

**Week 4 Decision Meeting:**
- **Attendees:** Product Manager, Engineering Lead, Architect, UX Lead, Technical Writer, Support Manager
- **Materials:** Validation spike results, customer discovery synthesis, business case, project plan
- **Outcome:** GO (proceed to implementation) or NO-GO (defer/pivot with documented reasons)
- **If GO:** Finalize architecture, begin sprint planning, kickoff implementation
- **If NO-GO:** Document reasons, define what would need to change for future GO, communicate to stakeholders

**Checkpoints During Implementation (if GO):**
- **End of Sprint 2:** Prototype review - validate technical approach with working code
- **End of Sprint 4:** Alpha review - internal dogfooding, collect feedback
- **End of Sprint 6:** Beta launch decision - GO if alpha feedback positive, NO-GO if major issues
- **Beta completion:** GA launch decision - GO if beta success criteria met, NO-GO if not

---

## Implementation Timeline (if GO)

**Phase 1: Validation (Weeks 1-4) - Before GO/NO-GO Decision:**
- Week 1-2: Platform Metrics validation spike + Perses technical spike (parallel)
- Week 2-3: Customer discovery (10 interviews) + Competitive analysis
- Week 3-4: Architecture design + UX wireframes + Documentation planning + Business case
- Week 4: GO/NO-GO decision meeting

**Phase 2: Implementation (Weeks 5-16) - IF GO:**
- Weeks 5-6: Sprint 0 (setup, infrastructure, development environment)
- Weeks 7-8: Sprint 1 (backend proxy, authentication, RBAC enforcement)
- Weeks 9-10: Sprint 2 (Perses embedding, first dashboard prototype)
- Weeks 11-12: Sprint 3 (pre-built dashboards, dashboard CRUD)
- Weeks 13-14: Sprint 4 (empty states, error handling, onboarding)
- Weeks 15-16: Sprint 5 (accessibility, performance optimization, support tooling)

**Phase 3: Beta Program (Weeks 17-28):**
- Week 17: Beta launch with 5-10 design partners
- Weeks 18-28: Beta feedback collection, iterative improvements, documentation refinement
- Weekly sync with beta customers, bi-weekly sprint demos
- Beta success criteria: NPS >40, >50% use ≥3 times/week, <5% support tickets

**Phase 4: GA Preparation (Weeks 29-32):**
- Week 29: Final security review, penetration testing
- Week 30: Support team training, runbook finalization
- Week 31: Documentation final review, video tutorial production
- Week 32: GA launch readiness review

**Phase 5: GA Launch (Week 33+):**
- Week 33: General Availability release
- Weeks 34-38: Hyper-care period (daily standups, rapid bug fixes)
- Month 2-3: Post-launch monitoring, success metrics tracking
- Month 4-6: Post-MVP planning based on feedback

**Total Timeline:** 8 months from validation start to GA (if no delays)
