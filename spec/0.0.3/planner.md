# OpenShift Web Console + AIOps Operator v0.0.3 계획

## 1) 목표

- v0.0.2에서 생성한 appserver `Service`를 OpenShift Web Console에 노출할 수 있도록 ConsolePlugin reconcile을 추가한다.
- appserver TLS/service-ca 연결을 위한 Kubernetes annotation과 Secret 참조 경계를 정리한다.
- Lightspeed backend API의 최소 HTTP service skeleton을 준비하되, 실제 Lightspeed API 호출은 후속 버전으로 이관한다.
- PostgreSQL Secret/PVC/HA 운영화와 Operator bundle/catalog 정리는 명시적 후속 범위로 유지한다.

## 2) 아키텍처 개요

1. `OpsMateConfigReconciler`가 appserver/PostgreSQL 리소스에 이어 ConsolePlugin desired object를 생성 또는 갱신한다.
2. ConsolePlugin은 OpenShift `console.openshift.io/v1` API를 unstructured object로 관리한다.
3. appserver `Service`는 service-ca serving cert annotation을 받는다.
4. appserver `Deployment`는 serving cert Secret을 읽을 수 있는 볼륨/환경 경계를 가진다.
5. backend HTTP skeleton은 health endpoint 중심으로 시작한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| ConsolePlugin | unstructured Kubernetes object | OpenShift API 의존성 추가 없이 reconcile 가능 |
| TLS | OpenShift service-ca annotation | ConsolePlugin과 backend service TLS 연결 기반 |
| Backend skeleton | Go net/http | 후속 Lightspeed API client 연결 전 최소 서버 표면 |
| Verification | unit tests, `go test`, `go build`, Kustomize | 리소스 shape와 빌드 회귀 검증 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | ConsolePlugin desired resource helper | 완료 | console plugin builder, tests |
| Phase 2 | ConsolePlugin reconcile 및 RBAC | 진행 예정 | reconciler 연결, RBAC update |
| Phase 3 | appserver TLS/service-ca 경계 | 진행 예정 | service annotation, deployment mount/env, tests |
| Phase 4 | backend HTTP skeleton 및 검증 | 진행 예정 | appserver package server skeleton, phase docs |

## 5) 마이그레이션 또는 운영 전략

- v0.0.3은 ConsolePlugin 리소스와 appserver TLS 경계를 잡는 버전이다.
- 실제 frontend bundle asset, plugin UI code, Lightspeed API 호출은 후속 버전에서 구현한다.
- ConsolePlugin disabled 상태에서는 ConsolePlugin 생성을 건너뛰는 방향을 유지한다.
- PostgreSQL Secret/PVC/HA는 v0.0.3에서 다루지 않는다.

## 6) 메시지/통신/데이터 프로토콜

| 경로 | 용도 | v0.0.3 상태 |
| --- | --- | --- |
| OpenShift Console -> ConsolePlugin | plugin registration | 구현 예정 |
| ConsolePlugin -> appserver Service | plugin backend endpoint | TLS 경계 |
| appserver -> Lightspeed API | 외부 API 호출 | 후속 버전 이관 |
| appserver -> PostgreSQL | conversation cache | 후속 버전 보강 |

## 7) 보안 고려사항

- ConsolePlugin backend service는 service-ca 기반 TLS를 전제로 한다.
- appserver TLS Secret은 Secret 이름만 참조하고 실제 인증서 값은 저장하지 않는다.
- RBAC는 ConsolePlugin 리소스 접근 권한만 추가한다.
- backend skeleton은 health endpoint 외 민감정보를 노출하지 않는다.

## 8) 완료 기준

- [x] ConsolePlugin desired object helper와 단위 테스트가 있다.
- [ ] reconciler가 ConsolePlugin을 create/update한다.
- [ ] RBAC에 ConsolePlugin 권한이 반영된다.
- [ ] appserver Service/Deployment에 TLS/service-ca 경계가 반영된다.
- [ ] backend HTTP skeleton이 빌드된다.
- [ ] `go fmt`, `go test`, `go build`가 성공한다.
- [ ] `kubectl kustomize config/default`가 성공한다.
- [ ] v0.0.4 이관 범위가 phase 문서에 기록된다.
