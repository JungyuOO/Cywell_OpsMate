# v0.0.11 Planner

## 1. 목표

- pgvector `VECTOR(n)` schema migration path를 코드와 문서로 검증한다.
- Operator `OpsMateConfig`에 embedding provider auth/config boundary를 추가한다.
- appserver Deployment가 embedding endpoint/model/dimensions/Secret/pgvector required env를 받을 수 있게 한다.
- pgvector SQL retrieval은 schema 전환 이후 사용할 contract로 고정하고, 현재 runtime은 BYTEA fallback ranking을 유지한다.

## 2. 아키텍처 개요

1. `OpsMateConfig.spec.embedding`이 endpoint, model, dimensions, credentials Secret, pgvector required flag를 표현한다.
2. appserver Deployment builder가 해당 값을 환경변수와 SecretKeyRef로 전달한다.
3. appserver migration helper가 `VECTOR(n)` 전환 SQL을 생성하고 dimensions validation을 수행한다.
4. 기본 runtime migration은 기존 BYTEA schema를 유지해 dev/test DB를 깨지 않는다.

## 3. 기술 스택

| 영역 | 선택 | 이유 |
| --- | --- | --- |
| API | `EmbeddingSpec` | Operator config boundary 명확화 |
| Secret | Kubernetes `SecretKeyRef` | token value 노출 방지 |
| Migration path | generated SQL helper | pgvector extension 없는 test DB에서도 schema path 검증 |
| Retrieval | BYTEA fallback 유지 | pgvector runtime 전환 전 안정성 유지 |

## 4. 구현 단계

| Phase | 범위 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | EmbeddingSpec API/env boundary | 완료 | CRD type, Deployment env tests |
| Phase 2 | pgvector migration helper | 완료 | SQL generation, validation tests |
| Phase 3 | runtime docs/operating conditions | 완료 | migration README, spec completion |
| Phase 4 | v0.0.12 handoff | 완료 | pgvector runtime activation scope |

## 5. 마이그레이션 또는 운영 전략

- v0.0.11은 기본 `ApplyMigrations`에 `VECTOR(n)` 전환 SQL을 자동 포함하지 않는다.
- 운영자가 pgvector extension readiness와 dimensions를 확정한 뒤 generated SQL path를 적용한다.
- 기존 BYTEA embedding rows는 reset 또는 재생성 대상으로 본다.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 데이터 | 상태 |
| --- | --- | --- |
| OpsMateConfig -> Deployment | embedding endpoint/model/dimensions/secret/pgvector required | 구현 예정 |
| Secret -> appserver | `CYOPS_EMBEDDING_TOKEN` | 구현 예정 |
| migration helper -> SQL | `VECTOR(n)` transition statements | 구현 예정 |

## 7. 보안 고려사항

- embedding token은 SecretKeyRef로만 전달한다.
- token value는 CR status, logs, response에 기록하지 않는다.
- pgvector migration SQL은 dimensions validation을 통과한 값만 생성한다.

## 8. 완료 기준

- [x] `EmbeddingSpec`과 Deployment env/SecretKeyRef tests가 있다.
- [x] `VECTOR(n)` migration SQL helper와 validation tests가 있다.
- [x] migration README가 BYTEA fallback과 pgvector activation 조건을 설명한다.
- [x] v0.0.12 handoff가 문서화된다.
