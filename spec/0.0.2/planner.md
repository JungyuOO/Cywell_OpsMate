# OpenShift Web Console + AIOps Operator v0.0.2 계획

## 1) 목표

- v0.0.1의 controller-runtime scaffold 위에 실제 Kubernetes 리소스 reconcile을 시작한다.
- `OpsMateConfig` 기준으로 appserver `Deployment`/`Service`와 PostgreSQL 기본 리소스를 생성 또는 갱신한다.
- ConsolePlugin과 Lightspeed backend API 구현은 경계를 유지하되, 이 버전에서는 배포 가능한 내부 서비스 표면을 우선 확보한다.
- RAG와 AIOps의 실제 동작 구현은 별도 버전 범위로 유지한다.

## 2) 아키텍처 개요

1. `OpsMateConfigReconciler`가 CR을 조회한다.
2. appserver 패키지가 backend `Deployment`와 `Service` desired object를 만든다.
3. postgres 패키지가 PostgreSQL `Deployment`와 `Service` desired object를 만든다.
4. reconciler가 owner reference를 설정하고 server-side object create/update를 수행한다.
5. reconcile 결과는 `status.conditions`와 `overallStatus`에 최소 상태로 기록한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| Reconcile | controller-runtime client | Kubernetes object lifecycle 관리 |
| Appserver | Go helpers returning typed apps/v1/core/v1 objects | 테스트 가능한 desired state 분리 |
| PostgreSQL | Deployment/Service | v0.0.0에서 기본 DB 방향으로 PostgreSQL 명시 |
| Status | `metav1.Condition` | Kubernetes 표준 상태 표현 |
| Verification | Go unit tests, `go test`, `go build`, Kustomize check | 리소스 shape와 빌드 검증 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | appserver desired resource helper | 진행 예정 | appserver Deployment/Service builder, tests |
| Phase 2 | PostgreSQL desired resource helper | 진행 예정 | postgres Deployment/Service builder, tests |
| Phase 3 | reconciler create/update 및 status | 진행 예정 | controller reconcile flow, RBAC update |
| Phase 4 | 검증 및 v0.0.3 이관 범위 정리 | 진행 예정 | phase 문서 체크, 검증 결과 |

## 5) 마이그레이션 또는 운영 전략

- v0.0.2는 단일 namespace 안에서 appserver와 PostgreSQL을 배포하는 최소 경로만 다룬다.
- credential Secret 생성은 하지 않고, CR의 Secret 참조를 appserver 환경변수로 전달하는 경계까지만 구현한다.
- PostgreSQL password Secret 자동 생성, PVC, backup, HA, migration은 후속 버전으로 이관한다.
- ConsolePlugin 리소스는 appserver endpoint가 안정화된 뒤 후속 버전에서 연결한다.

## 6) 메시지/통신/데이터 프로토콜

| 경로 | 용도 | v0.0.2 상태 |
| --- | --- | --- |
| OpsMateConfig -> appserver Deployment | Lightspeed endpoint/Secret 참조 전달 | 구현 예정 |
| OpsMateConfig -> PostgreSQL Deployment | 기본 DB 설정 전달 | 구현 예정 |
| appserver Service -> PostgreSQL Service | 내부 DB 연결 | 환경변수 경계 |
| Console Plugin -> appserver Service | Web Console 질의 전달 | 후속 버전 이관 |
| appserver -> Lightspeed API | 외부 API 호출 | 후속 버전 이관 |

## 7) 보안 고려사항

- Secret 값은 읽거나 로그에 남기지 않는다.
- appserver에는 Secret 이름만 환경변수로 전달한다.
- PostgreSQL 기본 배포는 개발용 최소 리소스이며, 운영 Secret/PVC/HA는 후속 버전에서 분리한다.
- RBAC는 v0.0.2에서 생성하는 `Deployment`, `Service`, status update 권한만 추가한다.

## 8) 완료 기준

- [ ] appserver desired `Deployment`/`Service` helper와 단위 테스트가 있다.
- [ ] PostgreSQL desired `Deployment`/`Service` helper와 단위 테스트가 있다.
- [ ] reconciler가 appserver/PostgreSQL 리소스를 create/update한다.
- [ ] reconciler가 `OpsMateConfig` status를 최소 갱신한다.
- [ ] RBAC에 필요한 apps/core/status 권한이 반영된다.
- [ ] `go fmt`, `go test`, `go build`가 성공한다.
- [ ] `kubectl kustomize config/default`가 성공한다.
- [ ] v0.0.3 이관 범위가 phase 문서에 기록된다.
