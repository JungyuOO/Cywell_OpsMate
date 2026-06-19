# v0.0.49 Phase 2 - CRC Lightspeed Wiring

## Tasks

- [x] Add CRC Lightspeed Operator install manifest.
- [x] Add CRC NetworkPolicy allowing CYOps to call Lightspeed appserver.
- [x] Apply CYWELL internal LLM `OLSConfig` on CRC.

## Verification

- `lightspeed-operator.v1.1.1` CSV is `Succeeded`.
- `OLSConfig/cluster` reports `Ready`.
- CYOps gateway can call `https://lightspeed-app-server.openshift-lightspeed.svc:8443/readiness`.
- Lightspeed OpenAPI exposes `POST /v1/query` with `LLMRequest.query`.

## Remaining Scope

- Upgrade CYOps and patch `OpsMateConfig.spec.lightspeed.apiBaseURL`.
