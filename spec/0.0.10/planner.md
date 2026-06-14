# v0.0.10 Planner

## 1. 목표

- pgvector extension readiness를 명확한 error로 검증할 수 있게 한다.
- external embedding provider endpoint/model/dimensions config boundary를 추가한다.
- 기존 BYTEA embedding rows를 이용해 `/api/chat` RAG path에서 ranked retrieval을 수행한다.
- pgvector `VECTOR(n)` migration은 readiness와 dimension이 확정된 후속 버전으로 넘긴다.

## 2. 아키텍처 개요

1. appserver config가 embedding endpoint, model, dimensions, pgvector required flag를 읽는다.
2. pgvector required flag가 켜진 경우 startup에서 `CREATE EXTENSION IF NOT EXISTS vector` readiness를 확인한다.
3. embedding provider는 mock 또는 configured HTTP endpoint를 사용한다.
4. retriever는 query embedding과 stored chunk embeddings를 비교해 top-k context를 만든다.
5. chat response citations는 rank와 score를 포함해 source를 안정적으로 표현한다.

## 3. 기술 스택

| 영역 | 선택 | 이유 |
| --- | --- | --- |
| pgvector readiness | explicit startup validation | extension 권한/설치 문제를 조기 노출 |
| External embedding | HTTP JSON provider | internal service 또는 BYO endpoint 모두 대응 |
| Ranking | byte-vector cosine fallback | pgvector migration 전에도 retrieval behavior 검증 |
| Citations | rank/score metadata | UI source 표시와 debug 가능성 확보 |

## 4. 구현 단계

| Phase | 범위 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | pgvector readiness/config | 완료 | env loader, validation |
| Phase 2 | external embedding provider | 완료 | HTTP provider, tests |
| Phase 3 | ranked retrieval/citations | 완료 | ranking, citation metadata |
| Phase 4 | v0.0.11 handoff | 완료 | pgvector migration next scope |

## 5. 마이그레이션 또는 운영 전략

- v0.0.10은 schema column을 `VECTOR(n)`으로 변경하지 않는다.
- 운영에서 `CYOPS_PGVECTOR_REQUIRED=true`를 켜면 extension readiness failure가 startup error로 드러난다.
- BYTEA fallback ranking은 pgvector migration 전 API와 RAG UX를 검증하기 위한 과도기 경로다.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 데이터 | 상태 |
| --- | --- | --- |
| env -> appserver | embedding endpoint/model/dimensions, pgvector required | 구현 예정 |
| appserver -> embedding provider | chunk id, text | 구현 예정 |
| retriever -> repository | document ids, embeddings | 구현 예정 |
| retriever -> chat response | citations with rank/score | 구현 예정 |

## 7. 보안 고려사항

- external embedding provider endpoint와 model config는 응답에 노출하지 않는다.
- provider에는 selected chunk/query text만 전달한다.
- pgvector readiness error는 secret/DSN을 포함하지 않는다.
- BYOKnowledge는 계속 사용하지 않는다.

## 8. 완료 기준

- [x] pgvector readiness validation test 또는 integration evidence가 있다.
- [x] HTTP embedding provider config와 test가 있다.
- [x] `/api/chat` RAG path가 ranked citations를 반환한다.
- [x] v0.0.11 pgvector migration handoff가 문서화된다.
