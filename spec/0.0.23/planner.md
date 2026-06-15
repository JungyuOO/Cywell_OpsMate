# v0.0.23 Planner

## 1. 목표

- Correct the admin UX boundary: CYOps ConsolePlugin inside OpenShift Web Console is the primary entry point.
- Mark the OAuth-protected admin Route as an optional fallback path, not the normal Web Console path.
- Harden diagnostics with a schema contract for future ConsolePlugin UI work.

## 2. 아키텍처 개요

- Users already authenticate when they enter OpenShift Web Console.
- CYOps ConsolePlugin calls the appserver backend through the ConsolePlugin backend service path.
- `GET /api/ops/diagnostics` and `GET /api/ops/diagnostics/schema` are the backend contract for the CYOps diagnostics view.
- `adminAuthProxyEnabled` remains opt-in fallback for direct access outside Web Console.

## 3. 기술 스택

- OpenShift ConsolePlugin metadata annotations.
- Go appserver diagnostics schema endpoint.
- Existing PowerShell OpenShift smoke script with fallback Route opt-in.

## 4. 구현 단계

| Phase | 내용 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | ConsolePlugin primary entry contract | done | ConsolePlugin annotations |
| Phase 2 | Diagnostics schema hardening | done | `/api/ops/diagnostics/schema` |
| Phase 3 | Fallback Route opt-in cleanup | done | sample and smoke script updates |
| Phase 4 | v0.0.24 handoff | done | next scope handoff |

## 5. 마이그레이션 또는 운영 전략

- Default sample keeps `adminAuthProxyEnabled=false`.
- Installers should rely on Web Console session flow for the normal CYOps UX.
- Enable fallback Route only for direct operational access outside Web Console.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 용도 | 비고 |
| --- | --- | --- |
| OpenShift Web Console -> CYOps ConsolePlugin | Primary UI entry | no extra OAuth redirect |
| ConsolePlugin -> appserver backend | Chat, documents, diagnostics | backend service path |
| `GET /api/ops/diagnostics` | aggregate diagnostics data | admin protected |
| `GET /api/ops/diagnostics/schema` | UI consumption contract | admin protected |
| Fallback Route -> auth proxy -> appserver | Direct operational access | optional only |

## 7. 보안 고려사항

- Diagnostics remains aggregate-only.
- Fallback Route must not be treated as the normal in-console UX.
- Forwarded identity headers are trusted only behind the OpenShift Console/backend boundary or the optional auth proxy.

## 8. 완료 기준

- [x] ConsolePlugin metadata marks OpenShift Web Console as primary entry.
- [x] ConsolePlugin metadata records diagnostics backend path.
- [x] Diagnostics schema endpoint documents aggregate-only forbidden fields.
- [x] Sample config disables fallback admin Route by default.
- [x] Smoke script checks fallback Route only when explicitly requested.
