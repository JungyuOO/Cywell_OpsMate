# v0.0.46 Planner

## 1. Goal

- Fix CYOps drawer POST requests that can be rejected by the OpenShift Console plugin proxy with 403.
- Wire `OpsMateConfig.spec.lightspeed` into the runtime `LightspeedProvider` used by `/api/chat`.
- Allow a customer internal LLM or OpenShift Lightspeed-compatible endpoint to be used behind CYOps.

## 2. Architecture Overview

- The CYOps ConsolePlugin continues to call the appserver through the Console plugin backend proxy.
- Non-GET drawer requests add an `X-CSRFToken` header before crossing the Console proxy.
- The Operator injects `CYOPS_LIGHTSPEED_ENDPOINT`, provider, model, and optional token env vars into the appserver.
- The appserver posts chat requests to the configured Lightspeed-compatible HTTP endpoint.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | OpenShift callback dynamic plugin | Existing launcher and drawer |
| Chat provider | HTTP Lightspeed-compatible provider | Configured by `OpsMateConfig.spec.lightspeed` |
| Secret handling | Kubernetes SecretKeyRef | `credentialsSecretRef` reads key `token` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | CSRF and Lightspeed runtime wiring | done | appserver/controller tests |
| Phase 2 | OLM packaging and CRC upgrade | done | CRC CSV and gateway smoke passed |
| Phase 3 | Issue/PR handoff | done | PR #203 merged and issue #202 closed |

Tracking issue: #202

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.45` to `cywell-opsmate.v0.0.46`.
- Keep catalog graph as `v0.0.46 -> v0.0.45 -> v0.0.44`.
- Configure an internal LLM by setting `spec.lightspeed.apiBaseURL` to its chat endpoint.
- If the endpoint requires bearer auth, create a Secret with key `token` and set `spec.lightspeed.credentialsSecretRef`.

## 6. Message, Communication, And Data Protocol

- CYOps `/api/chat` payload remains unchanged.
- Appserver sends `{message, context, clusterContext, provider, model}` to the configured provider endpoint.
- Provider response expects JSON with an `answer` field.

## 7. Security Considerations

- Bearer token values are read from Kubernetes Secret and are not returned in API responses.
- CSRF header is static and contains no secret data.
- Provider errors are returned as generic `provider failed` responses.

## 8. Completion Criteria

- [x] Non-GET drawer requests include `X-CSRFToken`.
- [x] Appserver reads CYOps and legacy Lightspeed env names.
- [x] Operator injects Lightspeed endpoint/provider/model/token env vars.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [x] v0.0.46 installs on CRC.
- [x] Gateway smoke confirms v0.0.46 entry includes `X-CSRFToken`.
- [ ] Browser chat smoke no longer gets 403.
