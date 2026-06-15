# v0.0.22 Planner

## 1. 목표

- Add a repeatable OpenShift smoke procedure for the admin auth Route and pgvector migration Job.
- Add a CYOps admin diagnostics endpoint that links retrieval metrics, document state, and re-embedding capability without exposing secrets.
- Keep Red Hat OpenShift Lightspeed viewer untouched and keep CYOps diagnostics separate.

## 2. 아키텍처 개요

- `deploy/scripts/openshift-v022-smoke.ps1` verifies the OpenShift login context, creates the OAuth cookie Secret, checks the admin Route OAuth redirect, waits for migration Job completion when present, and rejects sensitive log patterns.
- `GET /api/ops/diagnostics` is admin-protected and returns aggregate operational state only.
- The diagnostics endpoint references existing operational endpoints instead of introducing a second UI implementation in this version.

## 3. 기술 스택

- PowerShell smoke script invoking `oc` and `curl.exe`.
- Go appserver HTTP endpoint.
- Existing retrieval metrics, document repository, and admin auth plumbing.

## 4. 구현 단계

| Phase | 내용 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | OpenShift admin auth smoke | done | `deploy/scripts/openshift-v022-smoke.ps1` |
| Phase 2 | Migration Job operational smoke | done | Job wait/log/status checks in smoke script |
| Phase 3 | Admin diagnostics endpoint | done | `GET /api/ops/diagnostics` |
| Phase 4 | v0.0.23 handoff | done | next scope handoff |

## 5. 마이그레이션 또는 운영 전략

- The smoke script is intentionally explicit; it does not create database DSN Secrets or approve migrations by itself.
- Operators must prepare the `OpsMateConfig`, DSN Secret, admin groups, and cookie Secret scope before live cluster validation.
- Diagnostics are aggregate only: no document content, DSNs, tokens, prompts, or raw provider payloads.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 입력 | 결과 |
| --- | --- | --- |
| Smoke -> OpenShift | `oc whoami`, namespace, sample manifest | verifies cluster context |
| Smoke -> Route | unauthenticated `curl -I` | expects OAuth redirect/protection |
| Smoke -> Job | `oc wait job/<name> --for=condition=complete` | verifies migration completion |
| Smoke -> logs | `oc logs job/<name>` | fails if secret-like patterns appear |
| Admin -> diagnostics | `GET /api/ops/diagnostics` with admin auth | aggregate metrics and links |

## 7. 보안 고려사항

- The smoke script checks output for DSN/token/password patterns.
- Diagnostics require admin authorization and do not include customer document contents.
- Route smoke checks for OAuth protection before treating the admin surface as valid.

## 8. 완료 기준

- [x] Admin Route smoke procedure is scripted.
- [x] OAuth cookie Secret creation is scripted.
- [x] Migration Job completion/status/log smoke is scripted.
- [x] Diagnostics endpoint is admin-protected.
- [x] Diagnostics response exposes only aggregate operational state.
- [x] PowerShell script parses successfully.
