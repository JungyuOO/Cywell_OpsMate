# v0.0.21 Planner

## 1. 목표

- Add a trusted OpenShift OAuth proxy and Route shape for CYOps admin operations.
- Reconcile approved pgvector migration Jobs and reflect live Job status in `OpsMateConfig.status`.
- Keep the Red Hat OpenShift Lightspeed viewer untouched; CYOps remains its own Operator/ConsolePlugin surface.

## 2. 아키텍처 개요

- The appserver remains a cluster-internal Service.
- A separate admin auth proxy Deployment/Service/Route is created only when `spec.console.adminAuthProxyEnabled=true`.
- The auth proxy uses an OpenShift ServiceAccount OAuth redirect annotation and forwards authenticated user/group headers to the appserver.
- The appserver still authorizes admin APIs with `spec.console.adminUsers/adminGroups`.
- The controller reconciles `PGVectorMigrationJob` only when migration is explicitly approved.

## 3. 기술 스택

- OpenShift OAuth proxy style Deployment.
- OpenShift Route as `route.openshift.io/v1` unstructured resource.
- Kubernetes ServiceAccount OAuth redirect annotations.
- Kubernetes batch Job status conditions.

## 4. 구현 단계

| Phase | 내용 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | OpenShift admin auth proxy resources | done | ServiceAccount, Deployment, Service, Route |
| Phase 2 | Live migration Job reconciliation | done | Job reconcile and status evidence application |
| Phase 3 | CRD/RBAC/sample updates | done | CRD schema, Role, sample |
| Phase 4 | v0.0.22 handoff | done | next scope handoff |

## 5. 마이그레이션 또는 운영 전략

- `adminAuthProxyEnabled` is opt-in so existing development installs remain unaffected.
- The OAuth cookie Secret must contain `session_secret`.
- Operators should bind the desired OpenShift users/groups and configure `adminUsers/adminGroups` before exposing the admin Route.
- `pgVectorMigrationApproved=true` remains the explicit gate for migration Job reconciliation.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 입력 | 결과 |
| --- | --- | --- |
| Route -> auth proxy | OpenShift login session | authenticated request to appserver |
| Auth proxy -> appserver | `X-Forwarded-User`, `X-Forwarded-Groups` | appserver admin allowlist check |
| Controller -> Job | approved `PGVectorMigrationJob` spec | Job created/updated |
| Job -> CR status | `JobComplete=True` | `status.pgVectorReady=true` |
| Job -> CR status | `JobFailed=True` | `status.pgVectorReady=false`, `status.pgVectorLastError` |

## 7. 보안 고려사항

- The appserver admin API must not be exposed directly outside the cluster.
- Forwarded identity headers are trusted only when traffic enters through the auth proxy Route.
- DSN, token, and cookie material remain Secret-backed.
- Admin user/group authorization remains appserver-side to prevent a successfully authenticated but unauthorized OpenShift user from invoking admin APIs.

## 8. 완료 기준

- [x] Auth proxy ServiceAccount includes an OpenShift OAuth redirect reference to the Route.
- [x] Auth proxy Deployment forwards user headers to the cluster-internal appserver Service.
- [x] Route targets the auth proxy Service, not the appserver Service.
- [x] Controller reconciles approved migration Jobs.
- [x] Controller applies live Job status to `OpsMateConfig.status.pgVectorReady`.
- [x] CRD/RBAC/sample resources include the new admin auth proxy fields and required permissions.
