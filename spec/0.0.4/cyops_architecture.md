# CYOps Architecture Decision

## 결정

CYOps는 Red Hat OpenShift Lightspeed Viewer를 재사용하거나 iframe으로 감싸는 UI가 아니다. CYOps Operator가 설치되면 OpenShift Web Console에는 CYOps ConsolePlugin이 별도로 등록되고, 사용자는 CYOps floating launcher와 chat drawer를 통해 질문한다.

Lightspeed Operator/API는 CYOps backend가 호출할 수 있는 provider 중 하나다. 즉, Lightspeed의 답변 능력은 사용할 수 있지만, 사용자가 보는 브랜드, UI, 문서 관리, RAG 흐름은 CYOps가 소유한다.

## UI 계약

- 우측 하단 launcher는 CYOps 아이콘/브랜드로 표시한다.
- drawer 상단 title은 `CYOps`로 표시한다.
- chat composer의 `+` action은 문서 업로드와 문서 관리 drawer를 연다.
- 문서 관리 drawer는 Lightspeed Viewer 위치의 좌측 또는 chat drawer 내부 split panel로 표시한다.
- 고객 문서 목록은 업로드 파일명, 상태, chunk/embedding 상태, 업로드 사용자, 생성 시각, 삭제 action을 포함한다.

## Lightspeed 관계

- Red Hat Lightspeed ConsolePlugin이 설치되어 있으면 Red Hat UI도 별도로 보일 수 있다.
- CYOps Operator는 Red Hat ConsolePlugin을 임의로 숨기거나 삭제하지 않는다.
- 고객 환경에서 Red Hat UI를 보이지 않게 하려면 Lightspeed Operator의 console plugin 활성화 정책 또는 설치 구성을 별도로 정해야 한다.
- CYOps는 `OpsMateConfig.spec.lightspeed`로 받은 endpoint/Secret을 사용해 Lightspeed REST API만 호출한다.

## 고객 문서 RAG

CYOps는 BYOKnowledge를 사용하지 않는 자체 RAG pipeline을 목표로 한다.

기본 흐름:

1. 사용자가 CYOps UI에서 문서를 업로드한다.
2. appserver가 원본 파일을 object storage 또는 PVC에 저장한다.
3. appserver가 PostgreSQL에 document metadata를 기록한다.
4. ingestion worker가 문서를 chunk로 분할한다.
5. embedding worker/server가 chunk embedding을 생성한다.
6. PostgreSQL metadata와 pgvector index를 이용해 query-time context를 검색한다.
7. appserver가 검색 context와 사용자 질문을 provider policy에 따라 Lightspeed API 또는 다른 LLM provider로 전달한다.

## 우선순위

v0.0.5는 UI 전체 구현보다 API contract와 backend skeleton을 먼저 구현한다. 그 다음 frontend bundle과 문서 업로드 UI를 연결한다.

## 후속 결정 필요

- Lightspeed REST API endpoint와 auth scheme의 정확한 계약
- pgvector를 기본 포함할지, 별도 vector store를 옵션으로 둘지
- 문서 원본 저장소 기본값(PVC vs S3-compatible object store)
- 고객/namespace/tenant 권한 모델
- 업로드 파일 포맷과 최대 크기
