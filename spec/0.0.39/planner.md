# v0.0.39 Planner

## 1. 목표

- Fix the CYOps dynamic plugin manifest so OpenShift Console accepts it as a standard callback plugin.
- Move plugin display metadata from top-level legacy fields into `customProperties.console`.
- Preserve the v0.0.38 launcher/chat implementation and make it loadable by the live Console plugin validator.

## 2. 아키텍처 개요

- `ConsolePlugin/cyops-console` remains the OpenShift registration resource.
- The appserver continues to serve `/plugin-manifest.json` and `/plugin-entry.js`.
- The manifest now follows the standard schema shape only:
  - `name`
  - `version`
  - `baseURL`
  - `loadScripts`
  - `registrationMethod`
  - `dependencies`
  - `customProperties`
  - `extensions`

## 3. 기술 스택

| Area | Tooling | Notes |
| --- | --- | --- |
| Manifest | OpenShift Console standard dynamic plugin schema | No top-level legacy display fields |
| Appserver | Go handlers | Static manifest and entry bundle |
| Packaging | OLM bundle/catalog | v0.0.39 replaces v0.0.38 |

## 4. 구현 단계

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Manifest schema correction | done | appserver tests |
| Phase 2 | Packaging and CRC smoke | done | OLM and endpoint evidence |
| Phase 3 | Issue/PR handoff | done | GitHub evidence |

Tracking issue: #188

## 5. 마이그레이션 또는 운영 전략

- Upgrade from `cywell-opsmate.v0.0.38` to `cywell-opsmate.v0.0.39` through OLM.
- Keep the existing service, plugin name, and plugin entry paths stable.

## 6. 메시지/통신/데이터 프로토콜

| Caller | Endpoint | Purpose |
| --- | --- | --- |
| OpenShift Console | `/plugin-manifest.json` | Standard dynamic plugin manifest |
| OpenShift Console | `/plugin-entry.js` | Callback plugin entry |

## 7. 보안 고려사항

- No credential flow changes.
- The plugin still uses same-origin Console backend paths.
- Red Hat Lightspeed and other installed plugins remain untouched.

## 8. 완료 기준

- [x] Manifest no longer has top-level `displayName` or `description`.
- [x] Manifest exposes `customProperties.console.displayName`.
- [x] Go tests pass.
- [x] v0.0.39 packaging validates and installs on CRC.
- [x] CRC endpoint smoke confirms the accepted manifest shape.
