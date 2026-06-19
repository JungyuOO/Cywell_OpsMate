# v0.0.38 Planner

## 1. 목표

- Implement a real OpenShift dynamic plugin manifest and callback entry bundle for `cyops-console`.
- Show a CYOps launcher in the OpenShift Web Console and open a CYOps chat/document drawer from that launcher.
- Keep CYOps separate from the Red Hat OpenShift Lightspeed viewer while using the appserver API as the backend surface.

## 2. 아키텍처 개요

- `ConsolePlugin/cyops-console` continues to point at the appserver Service.
- The appserver serves a standard `plugin-manifest.json` with `baseURL`, `loadScripts`, and `registrationMethod: callback`.
- The entry script calls `window.loadPluginEntry("cyops-console", ...)` and injects the CYOps launcher/drawer into the Web Console page.
- The drawer calls the appserver through the Console plugin backend path for chat and document operations.

## 3. 기술 스택

| Area | Tooling | Notes |
| --- | --- | --- |
| Dynamic plugin | OpenShift Console callback plugin contract | No Node build toolchain in this version |
| Backend serving | Go appserver static handlers | Manifest and entry script are served from appserver |
| Packaging | OLM bundle/catalog | v0.0.38 replaces v0.0.37 |
| Smoke | Go tests, OCP curl/browser checks | Browser smoke may require trusted CRC console session |

## 4. 구현 단계

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Dynamic plugin manifest and callback entry | done | appserver tests |
| Phase 2 | CYOps launcher, chat, and document drawer | done | entry script browser smoke |
| Phase 3 | v0.0.38 packaging and CRC install smoke | done | OLM and endpoint evidence |
| Phase 4 | Issue/PR handoff | done | GitHub closure evidence |

Tracking issue: #186

## 5. 마이그레이션 또는 운영 전략

- Upgrade through the existing OLM channel graph from `v0.0.37` to `v0.0.38`.
- Keep the existing appserver Service and ConsolePlugin names stable.
- Do not remove or hide Red Hat Lightspeed ConsolePlugin resources.

## 6. 메시지/통신/데이터 프로토콜

| Caller | Endpoint | Purpose |
| --- | --- | --- |
| OpenShift Console | `/plugin-manifest.json` | Dynamic plugin manifest |
| OpenShift Console | `/plugin-entry.js` | Callback entry bundle |
| CYOps drawer | `/api/chat` | Lightspeed-compatible chat provider path |
| CYOps drawer | `/api/documents` | Customer document list/upload metadata path |

## 7. 보안 고려사항

- The plugin uses same-origin credentials through the OpenShift Console plugin backend path.
- No separate OAuth redirect is introduced in the frontend.
- Document content and provider tokens are not exposed in the drawer.
- The Red Hat Lightspeed viewer remains untouched.

## 8. 완료 기준

- [x] Manifest includes standard OpenShift dynamic plugin fields.
- [x] Entry script fires the `window.loadPluginEntry` callback.
- [x] CYOps launcher and chat/document drawer are present in the entry bundle.
- [x] Go tests pass.
- [x] v0.0.38 OLM packaging validates and installs on CRC.
- [x] GitHub Issue/PR workflow records the completed scope.
