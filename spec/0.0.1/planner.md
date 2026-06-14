# OpenShift Web Console + AIOps Operator v0.0.1 계획

## 1) 목표

- v0.0.0에서 정의한 OpsMate Operator 구조 초안을 실제 Kubebuilder/controller-runtime 기반 프로젝트로 진전시킨다.
- `OpsMateConfig` API를 Kubernetes CRD로 생성할 수 있는 형태로 정리한다.
- manager, scheme 등록, RBAC, reconciler 기본 골격을 추가한다.
- OpenShift Web Console Plugin, PostgreSQL, Lightspeed API, AIOps, RAG 구현은 실제 배포 리소스 생성 전 단계의 경계와 설정만 다룬다.

## 2) 아키텍처 개요

1. `cmd/manager`는 controller-runtime manager를 시작한다.
2. `api/v1alpha1`는 `OpsMateConfig` 타입과 scheme 등록 코드를 제공한다.
3. `internal/controller`는 `OpsMateConfig` reconciler를 등록하고 관찰 가능한 기본 상태 전환만 담당한다.
4. `config/crd`, `config/rbac`, `config/manager`, `config/default`, `config/samples`는 Operator 배포에 필요한 최소 Kustomize 구조를 가진다.
5. Console, appserver, postgres, aiops, rag 패키지는 후속 reconcile 구현을 위한 내부 경계로 유지한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| Runtime | Go, controller-runtime | Operator reconciliation 표준 구현 |
| Scaffold | Kubebuilder v4 계열 layout | v0.0.0에서 정한 OpenShift Lightspeed Operator 유사 구조 |
| API | Kubernetes CRD `opsmate.cywell.io/v1alpha1` | OpenShift에서 선언형 설정 관리 |
| Config | Kustomize manifests | Operator SDK/Kubebuilder 기본 배포 방식 |
| Verification | `go fmt`, `go test`, `go build`, manifest generation check | 코드와 생성 산출물 일관성 확인 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | controller-runtime 의존성 및 API scheme 정리 | 완료 | `go.mod`, `groupversion_info.go`, 타입 marker |
| Phase 2 | manager 및 reconciler 등록 | 완료 | `cmd/manager`, reconciler setup |
| Phase 3 | CRD/RBAC/manager Kustomize 기본 매니페스트 | 완료 | `config/*`, sample CR |
| Phase 4 | 검증 및 v0.0.2 이관 범위 정리 | 완료 | phase 문서 체크, 검증 결과 |

## 5) 마이그레이션 또는 운영 전략

- v0.0.1은 Kubernetes 리소스를 실제로 생성하는 reconcile 로직을 최소화한다.
- CRD와 controller 실행 가능성을 먼저 확보하고, appserver/postgres/console 배포 생성은 다음 버전으로 넘긴다.
- Secret 값, API token, DB password는 샘플과 문서에 실제 값을 기록하지 않는다.

## 6) 메시지/통신/데이터 프로토콜

| 경로 | 용도 | v0.0.1 상태 |
| --- | --- | --- |
| Kubernetes API -> OpsMateConfig controller | CR 변경 감지와 reconcile 트리거 | controller 등록 |
| OpsMateConfig -> status.conditions | 기본 reconcile 결과 표현 | 최소 조건 필드 검토 |
| Operator -> appserver/postgres/console | 리소스 생성 | 후속 버전 이관 |
| Backend -> Lightspeed API | 외부 API 호출 | 후속 버전 이관 |

## 7) 보안 고려사항

- Lightspeed API credential은 `Secret` 참조 필드로만 표현한다.
- 샘플 CR에는 실제 token, password, endpoint secret을 넣지 않는다.
- RBAC는 `OpsMateConfig`와 향후 관리 대상 리소스에 필요한 권한만 단계적으로 추가한다.
- local `.env`와 캐시 산출물은 Git에 포함하지 않는다.

## 8) 완료 기준

- [x] controller-runtime 의존성이 추가되고 `go mod tidy`가 성공한다.
- [x] `OpsMateConfig` 타입이 scheme에 등록된다.
- [x] manager가 controller-runtime 기반으로 실행 가능한 형태가 된다.
- [x] reconciler가 `OpsMateConfig`에 대해 setup된다.
- [x] 최소 CRD/RBAC/sample manifest 구조가 작성된다.
- [x] `go fmt`, `go test`, `go build`가 성공한다.
- [x] 다음 버전으로 넘길 실제 리소스 reconcile 범위가 phase 문서에 기록된다.
