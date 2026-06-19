# v0.0.38 Phase 1 - Dynamic Plugin Bundle Contract

## 작업 항목

- [x] Add required OpenShift dynamic plugin manifest fields.
- [x] Serve root `/plugin-entry.js` for Console backend proxy loading.
- [x] Make the entry script call the callback registration API.
- [x] Keep the diagnostics view as a fallback route.

## 검증

- `go test ./...` passed.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed.
- Local appserver smoke confirmed `/plugin-manifest.json` includes `registrationMethod: callback` and `loadScripts`.
- Local appserver smoke confirmed `/plugin-entry.js` includes `window.loadPluginEntry`.

## 남은 범위

- CRC OLM install smoke remains for Phase 3.
