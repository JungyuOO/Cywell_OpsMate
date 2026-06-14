# v0.0.12 Planner

## 1. 목표

- pgvector runtime activation을 config로 분리한다.
- `/api/chat` RAG retrieval이 `bytea` fallback mode와 `pgvector` SQL mode를 구분할 수 있게 한다.
- pgvector SQL ranking query path를 코드로 검증한다.
- embedding token rotation과 retrieval observability 운영 조건을 문서화한다.

## 2. 아키텍처 개요

1. appserver가 `CYOPS_RETRIEVAL_MODE`를 읽어 `bytea` 또는 `pgvector` mode를 선택한다.
2. `bytea` mode는 기존 in-process cosine fallback을 유지한다.
3. `pgvector` mode는 SQL similarity query path를 사용한다.
4. slow retrieval threshold와 failure reason은 appserver 내부 관측 경계로 기록한다.

## 3. 기술 스택

| 영역 | 선택 | 이유 |
| --- | --- | --- |
| Retrieval mode | env config | OpenShift 배포에서 안전하게 전환 |
| pgvector query | SQL builder + repository method | pgvector DB 적용 전 query contract 검증 |
| Observability | lightweight sink interface | metrics dependency 없이 boundary 고정 |
| Auth operations | documentation | Secret rotation 운영 기준 명확화 |

## 4. 구현 단계

| Phase | 범위 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | retrieval mode config | 완료 | env loader, server wiring |
| Phase 2 | pgvector SQL ranking path | 완료 | query builder, repository method tests |
| Phase 3 | observability/auth operations | 완료 | retrieval observer, docs |
| Phase 4 | v0.0.13 handoff | 완료 | OpenShift live pgvector scope |

## 5. 마이그레이션 또는 운영 전략

- `CYOPS_RETRIEVAL_MODE=bytea`가 기본값이다.
- `CYOPS_RETRIEVAL_MODE=pgvector`는 pgvector schema activation 이후에만 사용한다.
- token rotation은 Secret update 후 appserver rollout restart를 기본 운영 절차로 둔다.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 데이터 | 상태 |
| --- | --- | --- |
| env -> appserver | `CYOPS_RETRIEVAL_MODE`, slow threshold | 구현 예정 |
| appserver -> PostgreSQL | pgvector top-k query | 구현 예정 |
| retrieval -> observer | duration, mode, failure reason | 구현 예정 |

## 7. 보안 고려사항

- embedding token은 logs/status/metrics에 기록하지 않는다.
- retrieval observer에는 query text와 chunk text를 기록하지 않는다.
- pgvector mode는 schema readiness와 dimensions가 맞는 경우에만 운영한다.

## 8. 완료 기준

- [x] retrieval mode config와 tests가 있다.
- [x] pgvector SQL ranking query builder와 tests가 있다.
- [x] retrieval observability boundary와 tests가 있다.
- [x] token rotation 및 v0.0.13 handoff가 문서화된다.
