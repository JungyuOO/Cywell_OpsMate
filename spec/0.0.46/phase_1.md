# v0.0.46 Phase 1 - CSRF And Provider Wiring

## Tasks

- [x] Add `X-CSRFToken` to non-GET drawer fetch requests.
- [x] Read `CYOPS_LIGHTSPEED_ENDPOINT`, token, provider, and model in appserver config.
- [x] Keep legacy `LIGHTSPEED_*` env fallback support.
- [x] Inject `CYOPS_LIGHTSPEED_*` env vars from `OpsMateConfig.spec.lightspeed`.
- [x] Send bearer token to configured provider endpoint when present.
- [x] Update regression tests.

## Verification

- `go test ./...` passed.
- Tests verify appserver config reads CYOps and legacy Lightspeed env names.
- Tests verify the provider sends bearer auth when `CYOPS_LIGHTSPEED_TOKEN` is configured.
- Tests verify the appserver Deployment injects `CYOPS_LIGHTSPEED_*` env vars and token SecretKeyRef.
- Tests verify the Console plugin entry includes `X-CSRFToken`.

## Remaining Scope

- OLM packaging, image publication, CRC upgrade, and browser chat smoke remain for later phases.
