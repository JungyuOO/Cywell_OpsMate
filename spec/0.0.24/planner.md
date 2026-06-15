# v0.0.24 Planner

## 1. 목표

- Add a runnable CYOps diagnostics view served through the ConsolePlugin backend path.
- Keep the normal user journey inside OpenShift Web Console without an additional OAuth redirect.
- Preserve the fallback admin Route as optional direct operational access only.

## 2. 아키텍처 개요

- The appserver serves `/console-plugin/diagnostics`, `/console-plugin/diagnostics.js`, and `/console-plugin/diagnostics.css`.
- The diagnostics JavaScript calls `/api/ops/diagnostics` and `/api/ops/diagnostics/schema` with same-origin credentials.
- `web/console-plugin/src` records the frontend entry contract without introducing a Node build toolchain.

## 3. 기술 스택

- Go appserver static handlers.
- Plain browser ES module source contract.
- Existing appserver diagnostics JSON endpoints.

## 4. 구현 단계

| Phase | 내용 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | Console diagnostics served view | done | `/console-plugin/diagnostics` |
| Phase 2 | Diagnostics frontend source contract | done | `web/console-plugin/src/*.js` |
| Phase 3 | Verification and no-extra-OAuth checks | done | appserver tests |
| Phase 4 | v0.0.25 handoff | done | next scope handoff |

## 5. 마이그레이션 또는 운영 전략

- No new frontend dependency is introduced in this version.
- The served diagnostics view is intentionally small and can be replaced by a full ConsolePlugin bundle in a later version.
- Fallback admin Route checks remain opt-in through the smoke script.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 용도 |
| --- | --- |
| `GET /console-plugin/diagnostics` | CYOps diagnostics view shell |
| `GET /console-plugin/diagnostics.js` | diagnostics UI module |
| `GET /console-plugin/diagnostics.css` | diagnostics UI styles |
| `GET /api/ops/diagnostics` | aggregate diagnostics data |
| `GET /api/ops/diagnostics/schema` | aggregate-only UI contract |

## 7. 보안 고려사항

- The view does not handle OAuth redirects.
- API calls use same-origin credentials and rely on the ConsolePlugin backend path.
- Diagnostics content remains aggregate-only and secret-free.

## 8. 완료 기준

- [x] Diagnostics HTML view is served.
- [x] Diagnostics JS calls only same-origin appserver diagnostics APIs.
- [x] Tests verify no OAuth route handling exists in the console path.
- [x] Frontend source contract is recorded under `web/console-plugin/src`.
- [x] README documents Web Console as the normal path.
