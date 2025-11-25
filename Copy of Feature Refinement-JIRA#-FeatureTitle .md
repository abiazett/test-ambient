# Feature Refinement Template

> **Instructions**: Create a copy of this template. Name the document following this format: `FeatureRefinement - Jira# - TITLE OF STRAT`. Add your feature refinement document to the [feature refinement documents folder](https://drive.google.com/drive/folders/1SZjUobSqm4t8etx9nVY0agPhH4Wu0IuS?usp=drive_link).
>
> *Guidance for the Feature owner & Delivery owner is located [here](?tab=t.ecd9znwfr865).*

---

## Feature Metadata

| Field | Value |
|-------|-------|
| **Feature Jira Link** | <Add link to the feature (XXXSTRAT-####)> |
| **Status** | Not started |
| **Slack Channel / Thread** | <Add slack channel or discussion thread if applicable> |
| **[Feature Owner](?tab=t.ecd9znwfr865#bookmark=id.f39w0xy1rfqd)** | Person |
| **[Delivery Owner](?tab=t.ecd9znwfr865#bookmark=id.f39w0xy1rfqd)** | Person |
| **RFE Council Reviewer** | Person |
| **Product** | <indicate product: RHOAI (specify managed or self-managed), RHAIIS, RHEL AI> |

---

## Feature Details

### Feature Overview

<Description of feature goes here. Who (which persona/user role) benefits from this feature, and how? What is the difference between today's current state and a world with this feature? If applicable, share a high level user narrative: Tell how the user would use this functionality and when they would use it.>

### The Why

<Why are we doing this feature now? What will this bring to the platform, how will it help us win customers, and what types of customers will benefit from it? What supporting data do we have? Provide details under supporting documentation. If this is assumptive and not data-driven, provide clear success metrics for validating assumptions.>

### High Level Requirements

<List of functionality delivered into the product through this feature, with a focus on only including those requirements that are mandatory to achieve initial release.>

*Suggested format: As a [user role/persona], I want [capability/feature], So that [benefit/business value].*

### Non-Functional Requirements

<List non-functional requirements such as performance parameters, security concerns or areas to watch, special considerations for disconnected environments, user expectations, upgrade from previous release considerations, etc.>

*Reference: [Non-Functional Requirements](https://www.altexsoft.com/blog/non-functional-requirements/)*

### Out-of-Scope

<Define the boundaries of the feature. List items that are Out-of-Scope.>

### Acceptance Criteria

<What must be true for this feature to be considered complete? Document clear, testable outcomes that define "done.">

*Suggested format: Given [context or precondition] When [action taken] Then [expected outcome]*

### Risks & Assumptions

<Document Risks and Assumptions.>

- **Risks**: Potential blockers or threats to successful delivery.
- **Assumptions**: Conditions believed to be true but needing validation.

### Supporting Documentation

<Provide links (not copies) to designs, workflows, wireframes, discovery notes, or technical documentation.>

### Additional Clarifying Information

<Add any other information that would further clarify the requirement and aid the delivery teams in understanding the full context of this request.>

---

## New Feature / Component Prerequisites & Dependencies

### ODH/RHOAI Build Process Onboarding

**Question**: Will this feature require onboarding of a new container Image or component? **YES** or **NO**

<If this feature requires onboarding of a new container image or new product component, [follow these instructions](?tab=t.658ekv1v7j02) to add the new container image or product component to the ODH/RHOAI build process. The process requires a 2 week lead time. Reach out to [Jay Koehler](mailto:jkoehler@redhat.com) for any questions or concerns.>

### License Validation

**Question**: Will this feature require bringing in new upstream projects or sub-projects into the product? **YES** or **NO**

<If this feature will require bringing in new upstream projects or sub-projects into the product, check the licensing for the project. The preference is Apache 2.0. If the Apache License can not be obtained, post the specific reason why not and which license is needed in forum-openshift-ai-architecture to get approval.>

### Accelerator/Package Support

**Question**: Does this feature require support from the AIPCC team? **YES** or **NO**

The questions below should help identify this if you're unsure:

- Are there any new or updated package requests? If yes, please follow the instructions in [AIPCC-1](https://issues.redhat.com/browse/AIPCC-1) to clone and populate the ticket accordingly. Please attach the new epic to the STRAT feature.
- Which accelerators should be supported or is this for a new accelerator?
- If there is any other work needed from AIPCC (i.e. a new variant, a new image for RHAIIS, etc.), please open a new Feature in the AIPCC project and fill in the template accordingly. Attach this Feature to the RHAISTRAT issue with the "depends on" issue link.

### Architecture Review Check

**Questions**:
- Does the feature have the label "requires_architecture_review"? **YES** or **NO**
- Does the related RFE indicate "Requires architecture review: YES"? **YES** or **NO**

If the answer is **YES** to either question, the feature design **must** be reviewed at the OpenShift AI Architecture Forum before the team commits to a specific solution. This ensures the proposed design aligns with the overall product vision and provides architectural visibility so that any concerns can be identified and addressed proactively.

Teams may still conduct any spikes or research activities to explore design options prior to the formal review.

### Additional Dependencies

<Include any other prerequisites or dependencies needed to deliver this feature.>

---

## High Level Plan

> **Note**: Include teams that are actively involved in delivering the feature in the table below. If you're unsure whether a team should be involved, check with the team. The platform team is included by default and will confirm if their support is needed.

### Team Delivery Plan

| Team(s) Involved in Delivery | Start Date | Work to Deliver (EPIC) | Team Dependencies | [T-Shirt Size Estimate](?tab=t.2od5zi1gdoq) | Approval/Comments |
|------------------------------|------------|------------------------|-------------------|---------------------------------------------|-------------------|
| [team-ai-core-platform](mailto:team-ai-core-platform@redhat.com) *(DO NOT REMOVE. Email the AI Core Platform TEAM using this alias)* | | | | | |
| | | | | | |
| | | | | | |

---

## How to Engage The Documentation and UXD Teams

### Documentation Team

Engage the Documentation team by:

1. Adding the "Documentation" component to the feature
2. Setting the "Product Documentation Required" field to **Yes** in the feature
3. Adding the docs team to the table in the section above
4. Reviewing the [Docs Intake Process](https://docs.google.com/document/d/1G_LKipII0DMX3UxpkxVEpgM9Pk5tHcfZdvnkjn9E1mI/edit?tab=t.0)
5. Making sure to flag any new features and enhancements for Release Notes

### UXD Team

Engage the UXD team by:

1. Adding the "UXD" component to the feature
2. Adding the UXD team to the table above
3. Reaching out to [Jenn Giardino](mailto:jgiardin@redhat.com) or [Beau Morley](mailto:bmorley@redhat.com)
