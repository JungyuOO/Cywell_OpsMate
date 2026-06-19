# v0.0.51 Phase 3 - CRC Smoke

## Tasks

- [x] Confirm Lightspeed remains Ready.
- [x] Confirm CYOps `OpsMateConfig` points to Lightspeed `/v1/query`.
- [x] Confirm CYOps `/api/chat` returns a Lightspeed answer.
- [x] Record issue/PR/CRC evidence.

## Verification

- Created #212 for the v0.0.51 CRC Lightspeed service TLS smoke.
- PR #213 merged and closed #212.
- CRC CSV `cywell-opsmate.v0.0.51` is `Succeeded`.
- CRC appserver Deployment image is `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.51`.
- CRC `/api/chat` smoke from `deploy/cyops-gateway` returned HTTP 200:
  `Hello from OpenShift Lightspeed! It's great to connect with CYOps. How can I assist you with your OpenShift environment today?`

## Remaining Scope

- Future hardening should mount the OpenShift service CA bundle instead of using the scoped `.svc` TLS fallback.
