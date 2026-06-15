# v0.0.26 Planner

## 1. Goal

- Make the appserver image runnable in the OpenShift Web Console backend path.
- Align `cmd/appserver` with the existing service-ca TLS deployment contract.
- Keep live OpenShift screenshot/network evidence as an explicit follow-up if no cluster deployment is available in this turn.

## 2. Architecture Overview

- The appserver Deployment starts `/appserver` directly.
- The appserver listens on `:8443` in-cluster through `CYOPS_LISTEN_ADDRESS`.
- `cmd/appserver` serves HTTPS when `TLS_CERT_FILE` and `TLS_KEY_FILE` are both set, and serves HTTP for local development when neither is set.
- The appserver image contains `/appserver` and `cyops-pgvector-migrate` so the migration Job can keep using the same image family.

## 3. Technical Stack

- Go `net/http` server with optional TLS.
- Kubernetes Deployment env/command wiring.
- Multi-stage Containerfile with a UBI minimal runtime image.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Appserver TLS runtime | done | `cmd/appserver` TLS config |
| Phase 2 | Deployment command/listen wiring | done | appserver Deployment builder |
| Phase 3 | Appserver image packaging | done | `deploy/containerfiles/appserver.Containerfile` |
| Phase 4 | Live evidence handoff | done | v0.0.27 handoff |

## 5. Migration Or Operation Strategy

- Local verification remains plain HTTP unless TLS env vars are set.
- In OpenShift, service-ca mounts `tls.crt` and `tls.key`; appserver uses those files directly.
- The migration command remains available as `cyops-pgvector-migrate` in the appserver image.

## 6. Message, Communication, And Data Protocol

| Setting | Purpose |
| --- | --- |
| `CYOPS_LISTEN_ADDRESS=:8443` | in-cluster HTTPS listener |
| `TLS_CERT_FILE` | service-ca certificate path |
| `TLS_KEY_FILE` | service-ca key path |
| `/appserver` | appserver image entrypoint |
| `cyops-pgvector-migrate` | migration Job command |

## 7. Security Considerations

- TLS env vars must be set together; partial TLS config fails startup.
- No OAuth redirect is added to the Web Console path.
- The appserver Deployment still relies on the trusted OpenShift Console/backend service path for forwarded identity headers.
- The image runs as non-root and the Deployment keeps privilege escalation disabled.

## 8. Completion Criteria

- [x] `cmd/appserver` supports the service-ca TLS contract.
- [x] appserver Deployment command and listen address are explicit.
- [x] appserver image packaging includes the appserver and migration binaries.
- [x] Tests and builds pass.
