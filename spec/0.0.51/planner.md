# v0.0.51 Planner

## 1. Goal

- Unblock CYOps-to-Lightspeed chat on CRC when Lightspeed is exposed through the internal HTTPS service.
- Keep OpsMate as a Lightspeed client only; OLSConfig continues to own the CYWELL internal LLM URL and model.
- Complete CRC smoke where CYOps `/api/chat` returns the Lightspeed answer.

## 2. Architecture Overview

- CYOps appserver calls `POST /v1/query` on `https://lightspeed-app-server.openshift-lightspeed.svc:8443`.
- For OpenShift internal `.svc` HTTPS endpoints only, the appserver uses a transport that skips certificate verification because CRC service-serving certificates are not in the base image trust store.
- External HTTPS endpoints continue to use the default Go trust behavior.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| CYOps provider | Go HTTP client | Scoped internal service TLS handling |
| Lightspeed | `lightspeed-operator.v1.1.1` | OLSConfig owns internal LLM provider |
| Packaging | OLM catalog | `v0.0.51 -> v0.0.50 -> v0.0.49` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Scoped `.svc` TLS handling | done | provider tests |
| Phase 2 | Package v0.0.51 | planned | pending |
| Phase 3 | CRC chat smoke | planned | pending |

Tracking issue: #212

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.50` to `cywell-opsmate.v0.0.51`.
- Keep the existing `OpsMateConfig.spec.lightspeed.apiBaseURL` pointing to the Lightspeed `/v1/query` service URL.
- Keep the existing `cyops-lightspeed-token` Secret for CRC smoke.
- No CRD or data migration.

## 6. Message, Communication, And Data Protocol

- CYOps sends `{"query": "<user message>"}` to Lightspeed.
- Lightspeed chooses the configured `cywell-cllm` provider and `gemma-4-26b-a4b-it-awq-8bit` model from OLSConfig.
- CYOps parses the Lightspeed `response` field into the chat answer.

## 7. Security Considerations

- The TLS verification bypass is restricted to HTTPS endpoints whose hostname is an OpenShift internal service DNS name containing `.svc`.
- External endpoints keep normal certificate verification.
- Future hardening should mount the OpenShift service CA bundle into the appserver and trust it explicitly instead of skipping verification.
- CRC uses a short-lived ServiceAccount token Secret. Production must choose either projected service account credentials or user OAuth token forwarding.

## 8. Completion Criteria

- [x] TLS bypass is scoped to OpenShift internal service DNS names.
- [x] Go tests pass: `go test ./...`.
- [x] OLM dry-run passes: bundle manifests and CatalogSource.
- [ ] v0.0.51 installs on CRC.
- [ ] CYOps `/api/chat` returns a Lightspeed answer on CRC.
