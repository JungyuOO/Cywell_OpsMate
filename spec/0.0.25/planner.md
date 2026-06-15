# v0.0.25 Planner

## 1. Goal

- Add a local appserver executable for browser verification.
- Capture local diagnostics view evidence without an extra OAuth redirect path.
- Prepare OpenShift Web Console evidence as the next live-cluster step.

## 2. Architecture Overview

- `cmd/appserver` runs the appserver on `CYOPS_LISTEN_ADDRESS` or `:8080`.
- Local smoke uses `CYOPS_ADMIN_USERS=admin` and `X-Forwarded-User: admin`.
- Local browser verification may set `CYOPS_DEV_ADMIN_USER=admin`; this injects `X-Forwarded-User` only for loopback requests so the direct browser path can emulate the OpenShift Console proxy without adding an OAuth redirect.
- Browser verification opens `/console-plugin/diagnostics` and checks the served static diagnostics view.

## 3. Technical Stack

- Go appserver CLI.
- PowerShell local smoke script.
- Browser-based local diagnostics verification.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Local appserver executable | done | `cmd/appserver` |
| Phase 2 | Local diagnostics smoke | done | `deploy/scripts/local-v025-diagnostics-smoke.ps1` |
| Phase 3 | Browser evidence | done | local browser verification notes |
| Phase 4 | v0.0.26 handoff | done | next scope handoff |

## 5. Migration Or Operation Strategy

- Local browser verification uses memory storage by default.
- PostgreSQL remains optional for the local view smoke.
- OpenShift Console evidence remains live-cluster work.

## 6. Message, Communication, And Data Protocol

| Path Or Setting | Purpose |
| --- | --- |
| `CYOPS_LISTEN_ADDRESS` | appserver listen address |
| `CYOPS_ADMIN_USERS=admin` | local diagnostics admin allowlist |
| `CYOPS_DEV_ADMIN_USER=admin` | loopback-only development user injection |
| `/console-plugin/diagnostics` | browser diagnostics view |
| `/api/ops/diagnostics` | aggregate diagnostics data |
| `/api/ops/diagnostics/schema` | diagnostics schema |

## 7. Security Considerations

- Local smoke uses a development-only forwarded user header.
- `CYOPS_DEV_ADMIN_USER` is ignored for non-loopback requests and must not be set in production/OpenShift deployments.
- The browser path must not contain OAuth Route handling.
- Diagnostics stays aggregate-only.

## 8. Completion Criteria

- [x] Local appserver binary builds.
- [x] Local smoke script validates diagnostics view/API/schema.
- [x] Browser verification opens the diagnostics page.
- [x] Tests pass.
