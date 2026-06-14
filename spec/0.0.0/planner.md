# OpenShift Web Console + AIOps Operator v0.0.0 구축 계획

## 1) 목표

- OpenShift Web Console에서 사용할 수 있는 Cywell OpsMate Operator의 초기 구조를 만든다.
- OpenShift Lightspeed Operator 구조를 참고해 Go 기반 Operator 프로젝트로 방향을 확정한다.
- Lightspeed API를 사용할 수 있는 설정 경계를 `OpsMateConfig` 초안에 반영한다.
- Lightspeed Operator가 사용하는 PostgreSQL conversation cache 방식을 기준으로 우리 프로젝트의 기본 DB를 PostgreSQL로 확정한다.
- AIOps와 RAG는 후속 버전에서 구현할 수 있도록 패키지 경계만 둔다.
- `v0.0.0`에서는 실제 reconcile 로직을 구현하지 않고 프로젝트 구조 초안과 명세를 완료한다.

## 2) 아키텍처 개요

통신 흐름 초안

1. OpenShift Web Console 사용자가 OpsMate Console Plugin을 통해 질의한다.
2. Console Plugin은 OpsMate backend service로 요청을 전달한다.
3. OpsMate backend는 Lightspeed API 설정을 사용해 외부 또는 내부 Lightspeed 호환 API를 호출한다.
4. 대화/세션 cache와 후속 운영 이력 저장은 PostgreSQL을 기본 DB로 사용한다.
5. AIOps module은 OpenShift 이벤트, 경고, 로그, 운영 지표를 분석하는 별도 내부 경계로 둔다.
6. RAG module은 BYOK index image 또는 사내 운영 지식 index를 연결하는 선택 기능으로 둔다.
7. Operator controller는 `OpsMateConfig`를 reconcile해 backend, PostgreSQL, Console Plugin, AIOps/RAG 관련 리소스를 관리한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| Operator | Go | OpenShift Lightspeed Operator가 Go 기반이며 controller-runtime/Kubebuilder 구조를 사용 |
| Scaffold | Kubebuilder v4 / Operator SDK 구조 | Lightspeed Operator의 `PROJECT`, `api`, `cmd`, `internal/controller`, `config`, `bundle` 구조와 정합 |
| API | `opsmate.cywell.io/v1alpha1` `OpsMateConfig` | 우리 프로젝트 전용 CRD로 Lightspeed/AIOps/RAG 설정을 분리 |
| Backend | Go service boundary | Operator와 같은 언어로 관리하고 후속 버전에서 HTTP API 구현 가능 |
| DB | PostgreSQL | Lightspeed Operator의 conversation cache DB와 동일 계열 사용 |
| Console | OpenShift ConsolePlugin | Web Console 내부 메뉴/화면 노출 기준 |
| AIOps | Go package boundary | Operator lifecycle과 같은 저장소에서 시작하되 서비스 경계 유지 |
| RAG | BYOK index image / vector DB 경계 | Lightspeed Operator의 RAG image/index path 구조 참고 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | 프로젝트 규칙 및 버전 명세 체계 정리 | ✅ 완료 | `AGENTS.md`, 초기 `planner.md`, `phase_1.md` |
| Phase 2 | Lightspeed Operator 전수 구조 분석 | ✅ 완료 | `lightspeed_operator_review.md`, 기술 결정 |
| Phase 3 | Go Operator 프로젝트 구조 초안 | ✅ 완료 | `go.mod`, `PROJECT`, `cmd/manager`, `internal/controller` |
| Phase 4 | OpsMateConfig API 초안 | ✅ 완료 | `api/v1alpha1/opsmateconfig_types.go` |
| Phase 5 | v0.0.0 검증 및 후속 버전 이관 | ✅ 완료 | phase 문서 체크, 검증 명령 정의, 다음 버전 범위 |

## 5) 운영 전략

- `v0.0.0`은 구조 초안만 포함한다.
- 실제 Kubernetes reconcile 로직은 다음 버전에서 구현한다.
- 다음 버전 시작 시 `spec/<next-version>/` 폴더와 새 브랜치를 생성한다.
- phase 완료 시 `planner.md`와 `phase_N.md` 체크를 갱신하고 커밋/푸시한다.
- GitHub Issue는 phase 작업 단위로 생성하고 PR 또는 commit 이력을 통해 Closed한다.

## 6) 메시지/통신/데이터 프로토콜

`v0.0.0`에서는 프로토콜을 구현하지 않고 다음 경계를 고정한다.

| 경로 | 용도 | v0.0.0 상태 |
| --- | --- | --- |
| Console Plugin → OpsMate backend | Web Console 질의 전달 | 구조만 정의 |
| OpsMate backend → Lightspeed API | Lightspeed API 호출 | `LightspeedSpec` 초안 |
| OpsMate backend → PostgreSQL | 대화 cache 및 운영 이력 저장 | `DatabaseSpec` 초안, type 기본값은 postgres |
| OpsMate backend → AIOps module | 장애/이벤트/경고 분석 | package boundary |
| OpsMate backend → RAG module | 운영 지식 검색 | `RAGSpec` 초안 |
| Operator → OpenShift API | Deployment/Service/Secret/ConsolePlugin reconcile | 후속 버전 이관 |

## 7) 보안 고려사항

- Lightspeed API credential은 Kubernetes Secret 참조만 사용하고 코드에 저장하지 않는다.
- PostgreSQL password는 Secret으로 관리하며 초기화 이후 임의 변경으로 DB를 손상시키지 않는다.
- Console Plugin과 backend 사이의 통신은 OpenShift service-ca 또는 사용자 제공 TLS Secret을 사용한다.
- 고객사별 운영 데이터와 RAG index는 namespace, Secret, image pull secret, audit log 기준을 후속 버전에서 분리한다.
- GitHub token, API token, DB password는 문서와 commit에 기록하지 않는다.

## 8) 완료 기준

- [x] 루트 개발 규칙이 `AGENTS.md`에 기록된다.
- [x] OpenShift Lightspeed Operator 저장소 구조를 확인하고 분석 문서를 남긴다.
- [x] Go module과 Operator SDK 구조 초안을 만든다.
- [x] `OpsMateConfig` API 타입 초안을 만든다.
- [x] PostgreSQL을 기본 DB 방향으로 명시한다.
- [x] AIOps/RAG는 후속 버전 구현 범위로 분리한다.
- [x] Go 검증 명령을 `Makefile`에 정의한다.
- [x] `go fmt`, `go test`, `go build`는 로컬 Go 설치 후 수행한다.
- [x] 다음 버전으로 이관할 작업을 phase 문서에 남긴다.
