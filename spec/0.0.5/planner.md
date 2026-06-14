# OpenShift Web Console + AIOps Operator v0.0.5 계획

## 1) 목표

- v0.0.4에서 정의한 CYOps backend API contract를 컴파일 가능한 Go handler skeleton으로 구현한다.
- `/api/chat`, `/api/documents`, `/api/documents/{documentId}` 요청/응답 DTO와 handler tests를 추가한다.
- 문서 metadata 표현, mocked provider routing, upload boundary를 구현 준비 상태로 만든다.
- PostgreSQL schema, storage adapter, Lightspeed provider, frontend bundle은 skeleton 수준까지 순차 진행한다.

## 2) 아키텍처 개요

1. ConsolePlugin frontend는 CYOps appserver API를 호출한다.
2. appserver는 chat/document API DTO를 validation하고 JSON 응답을 반환한다.
3. chat endpoint는 provider interface를 통해 mocked provider에 라우팅한다.
4. document endpoint는 in-memory repository skeleton으로 metadata 표현을 고정한다.
5. 이후 PostgreSQL repository와 storage adapter로 교체 가능한 경계를 둔다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| HTTP API | Go `net/http` | 현재 appserver skeleton과 일관 |
| DTO | Go struct + JSON tags | contract와 테스트 고정 |
| Provider | interface | Lightspeed provider 교체 가능성 |
| Document repo | interface + memory skeleton | PostgreSQL 전환 전 handler 테스트 가능 |
| Verification | unit tests, `go test`, `go build` | handler 동작과 빌드 검증 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | Backend API DTO 및 handler skeleton | 진행 예정 | `/api/chat`, `/api/documents`, tests |
| Phase 2 | PostgreSQL schema/migration skeleton | 진행 예정 | SQL schema, migration docs |
| Phase 3 | Storage 및 Lightspeed provider skeleton | 진행 예정 | storage/provider interfaces |
| Phase 4 | Frontend bundle skeleton 및 v0.0.6 handoff | 진행 예정 | frontend shell, handoff |

## 5) 마이그레이션 또는 운영 전략

- v0.0.5는 실제 DB 연결 없이 API와 interface 경계를 먼저 고정한다.
- 문서 업로드 파일 저장은 local/PVC adapter skeleton까지만 구현한다.
- Lightspeed provider는 실제 API 호출 전 mocked provider로 handler contract를 안정화한다.
- v0.0.6에서 PostgreSQL repository와 ingestion runtime을 연결한다.

## 6) 메시지/통신/데이터 프로토콜

| 경로 | 용도 | v0.0.5 상태 |
| --- | --- | --- |
| `POST /api/chat` | chat request/response | 구현 예정 |
| `GET /api/documents` | document list | 구현 예정 |
| `POST /api/documents` | document upload metadata acceptance | 구현 예정 |
| `GET /api/documents/{id}` | document detail | 구현 예정 |
| `DELETE /api/documents/{id}` | delete marker | 구현 예정 |

## 7) 보안 고려사항

- handler는 Secret 값이나 문서 원문 전문을 log/response에 넣지 않는다.
- upload skeleton은 파일명과 metadata만 반환한다.
- mocked provider는 외부 호출을 하지 않는다.
- 실제 provider 구현 전 customer context 최소 전달 정책을 유지한다.

## 8) 완료 기준

- [ ] `/api/chat` DTO와 handler test가 있다.
- [ ] `/api/documents` list/upload/detail/delete DTO와 handler test가 있다.
- [ ] chat endpoint가 mocked provider로 라우팅된다.
- [ ] document upload가 metadata를 표현한다.
- [ ] PostgreSQL schema/migration skeleton이 있다.
- [ ] storage/provider interface skeleton이 있다.
- [ ] frontend bundle skeleton 위치가 정해진다.
- [ ] `go fmt`, `go test`, `go build`가 성공한다.
- [ ] v0.0.6 scope가 정리된다.
