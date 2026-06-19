# v0.0.42 Planner

## 1. Goal

- Add an nginx gateway in front of the CYOps ConsolePlugin backend.
- Make the OpenShift ConsolePlugin point at the gateway Service instead of the appserver Service.
- Keep appserver as the API, Lightspeed, document, and RAG backend.

## 2. Architecture Overview

- `cyops-gateway` is the ConsolePlugin backend Service.
- `cyops-gateway` runs nginx with OpenShift service-serving TLS.
- nginx proxies plugin asset, diagnostics, and API requests to `cyops-appserver`.
- `cyops-appserver` remains responsible for `/plugin-manifest.json`, `/plugin-entry.js`, `/api/chat`, `/api/documents`, and diagnostics routes.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Gateway | nginx unprivileged container | Fronts OpenShift ConsolePlugin backend |
| TLS | OpenShift service-serving cert | Separate `cyops-gateway-tls` secret |
| Operator | controller-runtime | Reconciles ConfigMap, Deployment, Service, ConsolePlugin |
| Backend | Go appserver | Existing API and plugin asset source |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Gateway controller resources | done | Go tests |
| Phase 2 | OLM packaging and CRC upgrade | done | CRC CSV/gateway endpoint smoke |
| Phase 3 | Issue/PR handoff | done | PR #195 merged and issue #194 closed |

Tracking issue: #194

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.41` to `cywell-opsmate.v0.0.42`.
- The reconciler creates `cyops-gateway`, `cyops-gateway-tls`, and a gateway nginx ConfigMap.
- Existing appserver and PostgreSQL resources stay in place.
- The ConsolePlugin backend Service changes from `cyops-appserver` to `cyops-gateway`.

## 6. Message, Communication, And Data Protocol

- Browser to OpenShift ConsolePlugin backend uses HTTPS through the Console proxy.
- Console proxy to `cyops-gateway` uses Service backend TLS.
- `cyops-gateway` proxies HTTPS to `cyops-appserver`.
- No API payload contract changes.

## 7. Security Considerations

- Gateway runs as non-root with dropped Linux capabilities.
- Gateway TLS uses OpenShift service-serving certificate injection.
- nginx proxy disables upstream certificate verification only for in-cluster service-serving cert compatibility.
- No OAuth redirect or credential storage changes.

## 8. Completion Criteria

- [x] Reconciler creates gateway ConfigMap, Deployment, and Service.
- [x] ConsolePlugin backend points at `cyops-gateway`.
- [x] OLM permissions include ConfigMap reconciliation.
- [x] Go tests pass.
- [x] v0.0.42 installs on CRC and gateway endpoint serves plugin assets.
