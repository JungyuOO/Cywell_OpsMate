# v0.0.39 Phase 1 - Manifest Schema Correction

## 작업 항목

- [x] Remove top-level legacy `displayName` and `description` fields from the standard manifest.
- [x] Add `customProperties.console.displayName` and `customProperties.console.description`.
- [x] Add regression coverage that rejects top-level legacy fields.

## 검증

- `go test ./...` passed.
- Local appserver smoke returned:
  - `version: 0.0.39`
  - `hasTopLevelDisplayName: false`
  - `hasTopLevelDescription: false`
  - `displayName: CYOps`
  - `registrationMethod: callback`
  - `loadScript: plugin-entry.js`

## 남은 범위

- v0.0.39 OLM packaging and CRC smoke remain for Phase 2.
