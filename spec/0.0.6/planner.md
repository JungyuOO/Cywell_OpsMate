# OpenShift Web Console + AIOps Operator v0.0.6 계획

## 1) 목표

- v0.0.5의 in-memory document repository를 PostgreSQL repository로 교체 가능한 구현으로 확장한다.
- document metadata CRUD와 delete marker를 실제 PostgreSQL schema에 맞춰 구현한다.
- Docker PostgreSQL을 사용해 repository 동작을 검증한다.
- real storage, Lightspeed REST client, frontend API integration은 skeleton과 경계 보강까지 진행한다.

## 2) 아키텍처 개요

1. appserver handler는 `DocumentRepository` interface에만 의존한다.
2. PostgreSQL repository는 `database/sql`을 사용해 `cyops_documents`를 읽고 쓴다.
3. migration SQL은 v0.0.5 skeleton을 그대로 적용 가능해야 한다.
4. Docker PostgreSQL 검증은 로컬 개발 검증용이며 OpenShift 운영 배포와 분리한다.

## 3) 기술 스택

| 구성요소 | 기술 | 근거 |
| --- | --- | --- |
| DB access | Go `database/sql` | driver 교체 가능 |
| PostgreSQL driver | pgx stdlib | Go 생태계 표준 PostgreSQL driver |
| Test DB | Docker PostgreSQL | 실제 SQL schema 검증 |
| Migration | raw SQL skeleton | operator/appserver migration strategy 전 단계 |

## 4) 구현 단계

| Phase | 제목 | 상태 | 주요 산출물 |
| --- | --- | --- | --- |
| Phase 1 | PostgreSQL document repository | 완료 | repository implementation, tests |
| Phase 2 | Docker PostgreSQL integration verification | 완료 | docker-based test evidence |
| Phase 3 | storage/provider wiring refinements | 완료 | appserver dependency wiring notes |
| Phase 4 | v0.0.7 ingestion handoff | 완료 | ingestion/embedding implementation scope |

## 5) 마이그레이션 또는 운영 전략

- v0.0.6는 appserver가 runtime migration을 자동 적용하지 않는다.
- migration application은 운영 정책이 정해질 때까지 별도 단계로 유지한다.
- Docker PostgreSQL은 로컬 검증에만 사용한다.

## 6) 메시지/통신/데이터 프로토콜

| 경로 | 용도 | v0.0.6 상태 |
| --- | --- | --- |
| appserver -> PostgreSQL | document metadata persistence | 구현 예정 |
| upload handler -> repository | metadata create | 구현 예정 |
| document delete -> repository | deleting marker | 구현 예정 |

## 7) 보안 고려사항

- repository는 document 원문을 저장하지 않고 metadata/object URI만 저장한다.
- SQL query는 parameterized query만 사용한다.
- DB connection string은 Secret/env로 주입하는 방향을 유지한다.

## 8) 완료 기준

- [x] PostgreSQL document repository가 구현된다.
- [x] repository unit/integration test가 있다.
- [x] Docker PostgreSQL에서 migration과 repository 동작이 검증된다.
- [x] `go fmt`, `go test`, `go build`가 성공한다.
- [x] v0.0.7 ingestion/embedding 범위가 정리된다.
