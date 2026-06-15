# v0.0.20 Planner

## 1. 목표

- Move `/api/ops/reembed` from token-only admin authorization toward OpenShift production authorization.
- Add OpenShift migration Job evidence handling so completed `PGVectorMigrationJob` results can update runtime status.
- Keep Red Hat OpenShift Lightspeed viewer untouched; CYOps remains a separate ConsolePlugin surface.

## 2. 아키텍처 개요

- Appserver accepts either the existing `X-CYOps-Admin-Token` or OpenShift OAuth proxy forwarded identity headers.
- `OpsMateConfig.spec.console.adminUsers/adminGroups` declares users and groups allowed to call admin operations.
- `PGVectorMigrationJob` remains explicit and approval-gated; Job completion/failure evidence can be applied to `status.pgVectorReady`.

## 3. 기술 스택

- Go appserver HTTP handlers.
- Kubernetes batch Job status conditions.
- OpenShift ConsolePlugin and OAuth proxy style forwarded identity headers.
- OpsMateConfig CRD schema and sample manifest.

## 4. 구현 단계

| Phase | 내용 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | OpenShift admin identity auth | done | `AdminAuthConfig`, forwarded user/group checks |
| Phase 2 | Migration Job runtime evidence | done | `ApplyPGVectorMigrationJobStatus` |
| Phase 3 | CRD/sample/spec verification | done | CRD schema, sample, tests |
| Phase 4 | v0.0.21 handoff | done | next scope handoff |

## 5. 마이그레이션 또는 운영 전략

- Token auth remains available for local/dev and emergency fallback.
- Production should configure `console.adminGroups` for an OpenShift group managed by the cluster administrator.
- pgvector migration execution remains explicit; operators should approve and run or reconcile the migration Job before marking pgvector ready.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 입력 | 결과 |
| --- | --- | --- |
| Admin API | `X-CYOps-Admin-Token` | legacy admin authorization |
| Admin API | `X-Forwarded-User` | allowed when in `spec.console.adminUsers` |
| Admin API | `X-Forwarded-Groups` CSV | allowed when any group is in `spec.console.adminGroups` |
| Job status -> CR status | `JobComplete=True` | `status.pgVectorReady=true` |
| Job status -> CR status | `JobFailed=True` | `status.pgVectorReady=false`, `status.pgVectorLastError` recorded |

## 7. 보안 고려사항

- Header-based user/group authorization assumes the appserver is reachable only behind a trusted OpenShift OAuth/proxy layer in production.
- Direct public exposure of the appserver is not acceptable with forwarded-header auth.
- DSNs and admin tokens remain Secret-backed and are not rendered as plain environment values in generated resources.

## 8. 완료 기준

- [x] `/api/ops/reembed` supports token, forwarded-user, and forwarded-group admin authorization.
- [x] `OpsMateConfig.spec.console` can configure admin users and groups.
- [x] Appserver Deployment passes admin allowlists as environment variables.
- [x] Migration Job completion/failure evidence updates runtime status fields.
- [x] CRD schema and sample manifest include the new admin allowlists.
- [x] Verification commands are recorded in phase docs.
