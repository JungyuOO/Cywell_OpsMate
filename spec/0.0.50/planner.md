# v0.0.50 Planner

## 1. Goal

- Make CYOps `/api/chat` fully compatible with OpenShift Lightspeed `LLMRequest`.
- Remove all non-OLS fields from the backend-to-Lightspeed request body.
- Complete CRC smoke where CYOps calls Lightspeed and Lightspeed calls the CYWELL internal LLM.

## 2. Architecture Overview

- CYOps appserver calls `POST /v1/query` on Lightspeed.
- Request body contains only `query`.
- Lightspeed owns authentication, OpenShift context, provider/model selection, and internal LLM calls.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| CYOps provider | Go HTTP provider | Sends strict OLS `LLMRequest` |
| Lightspeed | `lightspeed-operator.v1.1.1` | Ready on CRC |
| Packaging | OLM catalog | `v0.0.50 -> v0.0.49 -> v0.0.48` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Strict OLS request body | done | provider test |
| Phase 2 | Package and CRC upgrade | planned | pending |
| Phase 3 | CYOps-to-Lightspeed smoke | planned | pending |

Tracking issue: #210

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.49` to `cywell-opsmate.v0.0.50`.
- Keep the existing CRC `OpsMateConfig.spec.lightspeed.apiBaseURL` and token Secret.
- No CRD or data migration.

## 6. Message, Communication, And Data Protocol

- CYOps sends `{"query": "<user message>"}`.
- CYOps parses Lightspeed `response` into its chat answer.
- CYOps does not send `message`, `context`, `clusterContext`, `model`, or `provider` to Lightspeed.

## 7. Security Considerations

- CRC uses a short-lived ServiceAccount token Secret for backend-to-Lightspeed smoke.
- Real deployments should decide whether CYOps forwards user tokens or uses a service account with approved Lightspeed access.

## 8. Completion Criteria

- [x] Provider request body contains only `query`.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [ ] v0.0.50 installs on CRC.
- [ ] CYOps `/api/chat` returns a Lightspeed answer on CRC.
