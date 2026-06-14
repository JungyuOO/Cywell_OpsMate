# OpenShift Web Console + AIOps Operator v0.0.4 계획

## 1) 목표

- Red Hat OpenShift Lightspeed UI를 그대로 노출하지 않고, Cywell OpsMate Operator가 제공하는 `CYOps` ConsolePlugin UI를 OpenShift Web Console 안에 띄운다.
- Lightspeed Operator/API는 답변 provider 중 하나로만 사용하고, 사용자에게 보이는 chat drawer/viewer는 CYOps가 소유한다.
- 고객사 문서 업로드, 문서 목록, 문서 관리, RAG 질의는 BYOKnowledge에 의존하지 않고 OpsMate 내부 기능으로 설계한다.
- v0.0.4는 구현 착수 전 아키텍처 결정을 고정하고, 이후 버전에서 frontend/backend/RAG를 나눠 구현할 수 있는 계약을 만든다.

## 2) 아키텍처 개요

1. OpenShift Web Console은 `ConsolePlugin`을 통해 `CYOps` 플러그인을 로드한다.
2. CYOps UI는 Lightspeed처럼 우측 하단 floating launcher와 chat drawer를 제공한다.
3. Red Hat Lightspeed Viewer는 CYOps의 UI가 아니며, OpsMate Operator가 이를 대체하거나 숨기는 것을 목표로 하지 않는다.
4. Lightspeed Operator가 설치되어 있더라도 CYOps는 자체 ConsolePlugin으로 별도 등록된다.
5. CYOps appserver는 질문을 받고 필요 시 Lightspeed REST API를 호출한다.
6. 고객 문서는 CYOps UI의 upload action으로 업로드되고, CYOps appserver가 metadata, chunk, embedding 상태를 관리한다.
7. 문서 목록/관리 패널은 chat drawer 안에서 `+` 버튼 또는 문서 drawer로 노출한다.
8. RAG는 BYOKnowledge 대신 OpsMate-owned ingestion/index/query pipeline으로 처리한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| Console UI | OpenShift ConsolePlugin dynamic plugin | Web Console 내부에 CYOps UI를 직접 노출 |
| Chat backend | Go appserver | Operator와 같은 언어/배포 경계 유지 |
| Lightspeed 연동 | REST API client/provider | Red Hat UI를 쓰지 않고 답변 provider로만 사용 |
| 문서 metadata | PostgreSQL | 고객 문서 목록, 상태, 권한, audit metadata 관리 |
| Vector index | PostgreSQL pgvector 우선 검토 | 별도 vector DB 없이 운영 단순화 |
| Embedding | 내부 embedding worker/server | BYOKnowledge 회피, 고객 문서 RAG 자체 구현 |
| 파일 원본 | PVC 또는 S3-compatible object store | 문서 원본 저장은 DB metadata와 분리 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | CYOps UX/API 아키텍처 결정 | 완료 | architecture decision, UI behavior contract |
| Phase 2 | Backend API contract | 완료 | chat/documents/ingestion API spec |
| Phase 3 | Data model 및 RAG pipeline contract | 완료 | PostgreSQL tables, vector/index strategy |
| Phase 4 | v0.0.5 구현 범위 정리 | 완료 | frontend/backend implementation handoff |

## 5) 마이그레이션 또는 운영 전략

- Lightspeed Operator 설치 여부는 CYOps UI 표시와 독립적으로 둔다.
- Lightspeed UI 자체를 숨기려면 Red Hat Lightspeed ConsolePlugin 설치/활성화 정책을 별도로 제어해야 하며, CYOps Operator가 타사 ConsolePlugin을 임의로 삭제하지 않는다.
- CYOps는 Lightspeed API endpoint/credential을 `OpsMateConfig.spec.lightspeed`로 받아 provider로 호출한다.
- 고객 문서 RAG는 PostgreSQL metadata와 vector index를 우선 사용하고, 대규모 고객사에서만 별도 embedding/vector service 분리를 검토한다.
- 문서 원본은 DB에 직접 저장하지 않고 PVC/object storage에 저장하는 방향을 기본으로 한다.

## 6) 메시지/통신/데이터 프로토콜

| 경로 | 용도 | v0.0.4 상태 |
| --- | --- | --- |
| CYOps ConsolePlugin -> appserver `/api/chat` | 사용자 질문과 RAG 옵션 전달 | 계약 정의 예정 |
| appserver -> Lightspeed REST API | OpenShift 운영 답변 provider 호출 | 계약 정의 예정 |
| CYOps ConsolePlugin -> appserver `/api/documents` | 문서 목록/업로드/삭제/상태 조회 | 계약 정의 예정 |
| appserver -> ingestion worker | 문서 chunk/embedding 생성 | 계약 정의 예정 |
| appserver -> PostgreSQL | metadata, chat session, document state, vector search | 계약 정의 예정 |

## 7) 보안 고려사항

- Lightspeed credential은 Secret 참조로만 사용하고 UI/API 응답에 노출하지 않는다.
- 고객 문서는 namespace/tenant 경계를 가져야 하며, 업로드/삭제/audit metadata를 기록한다.
- 문서 원본과 embedding/vector index는 고객별 격리를 기본으로 설계한다.
- CYOps chat은 민감정보 경고와 입력 정책을 자체 UI에 표시한다.
- Red Hat Lightspeed로 전달되는 질의에는 고객 문서 원문을 무조건 포함하지 않고, provider policy에 따라 최소 context만 보낸다.

## 8) 완료 기준

- [x] CYOps가 Lightspeed Viewer를 대체하는 별도 UI라는 결정이 문서화된다.
- [x] Lightspeed REST API provider 경계가 정의된다.
- [x] 고객 문서 upload/list/manage UI 경계가 정의된다.
- [x] BYOKnowledge 없이 OpsMate-owned RAG pipeline을 쓰는 방향이 정의된다.
- [x] backend API contract가 정의된다.
- [x] PostgreSQL/vector/embedding 데이터 경계가 정의된다.
- [x] v0.0.5 구현 우선순위가 정리된다.
