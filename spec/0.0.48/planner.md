# v0.0.48 Planner

## 1. Goal

- Enforce the integration boundary where CYOps calls only the OpenShift Lightspeed API.
- Keep the internal CYWELL LLM URL and model selection in Lightspeed `OLSConfig`.
- Provide an OLSConfig example for `http://cllm.cywell.co.kr/v1` with model `gemma-4-26b-a4b-it-awq-8bit`.

## 2. Architecture Overview

- `OLSConfig.spec.llm` owns LLM provider URL, provider type, model list, default provider, and default model.
- `OpsMateConfig.spec.lightspeed.apiBaseURL` points to the Lightspeed appserver/API endpoint only.
- CYOps appserver forwards chat requests to Lightspeed without sending LLM provider or model fields.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Lightspeed config | `ols.openshift.io/v1alpha1` `OLSConfig` | Owns CYWELL internal LLM provider and model |
| CYOps config | `opsmate.cywell.io/v1alpha1` `OpsMateConfig` | Owns only the Lightspeed API endpoint |
| Runtime proxy | Go appserver HTTP provider | Sends message/context to Lightspeed |
| Packaging | OLM catalog | `v0.0.48 -> v0.0.47 -> v0.0.46` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Runtime boundary cleanup | done | appserver/provider tests |
| Phase 2 | Lightspeed OLSConfig example and samples | done | manifest review |
| Phase 3 | Packaging, issue, PR, CRC upgrade | done | CRC CSV `cywell-opsmate.v0.0.48` Succeeded |

Tracking issue: #206

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.47` to `cywell-opsmate.v0.0.48`.
- Apply or update Lightspeed `OLSConfig` separately with the CYWELL internal LLM provider.
- Keep `OpsMateConfig.spec.lightspeed.apiBaseURL` set to the Lightspeed API route/service, not `http://cllm.cywell.co.kr/v1`.

## 6. Message, Communication, And Data Protocol

- CYOps sends `message`, `context`, and `clusterContext` to the configured Lightspeed API.
- CYOps does not send `model` or LLM `provider` fields.
- Lightspeed selects `cywell-cllm` and `gemma-4-26b-a4b-it-awq-8bit` from `OLSConfig`.

## 7. Security Considerations

- Internal LLM credentials remain in the Lightspeed namespace Secret referenced by `OLSConfig`.
- CYOps Lightspeed token behavior remains limited to the Lightspeed API hop.
- CYOps responses expose only the answer text returned by Lightspeed.

## 8. Completion Criteria

- [x] CYOps provider payload no longer includes model/provider fields.
- [x] Appserver Deployment no longer injects Lightspeed model/provider env vars from `OpsMateConfig`.
- [x] CYWELL Lightspeed `OLSConfig` example exists.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [x] v0.0.48 installs on CRC.
