# v0.0.9 Planner

## 1. 목표

- deterministic mock embedding provider로 chunk embedding persistence를 검증한다.
- PostgreSQL `cyops_document_embeddings`에 embedding result를 저장한다.
- `/api/chat` RAG 요청이 retriever를 호출하고 provider context와 citations를 응답에 포함한다.
- pgvector `VECTOR(n)` 전환 조건은 문서와 code boundary로 고정하되 production vector tuning은 하지 않는다.

## 2. 아키텍처 개요

1. ingestion으로 생성된 chunks를 embedding service가 조회한다.
2. embedding provider가 chunk별 vector payload를 만든다.
3. repository가 `cyops_document_embeddings`에 embedding rows를 교체 저장한다.
4. retriever가 document IDs 기준으로 ready chunks를 선택한다.
5. `/api/chat`은 retrieved citations를 provider context에 붙이고 응답에도 citations를 포함한다.

## 3. 기술 스택

| 영역 | 선택 | 이유 |
| --- | --- | --- |
| Embedding provider | deterministic mock | 외부 model 없이 persistence와 API path 검증 |
| Embedding storage | existing `BYTEA` column | pgvector dimension 확정 전 안전한 placeholder |
| Retrieval | repository-backed deterministic retriever | vector ranking 전 API contract 고정 |
| Chat integration | existing `ChatProvider` context | Lightspeed/API provider와 분리 유지 |

## 4. 구현 단계

| Phase | 범위 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | embedding provider boundary | 완료 | provider interface, mock tests |
| Phase 2 | embedding persistence | 완료 | repository methods, integration tests |
| Phase 3 | chat retrieval API | 완료 | retriever interface, `/api/chat` citations |
| Phase 4 | v0.0.10 handoff | 완료 | pgvector/retrieval next scope |

## 5. 마이그레이션 또는 운영 전략

- v0.0.9는 기존 `BYTEA` embedding column을 유지한다.
- pgvector 전환은 embedding dimensions와 extension 권한 조건이 확정된 뒤 `VECTOR(n)` migration으로 진행한다.
- embedding 재실행 시 기존 embedding rows를 삭제 후 다시 insert한다.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 데이터 | 상태 |
| --- | --- | --- |
| chunks -> embedding provider | chunk id, text | 구현 예정 |
| embedding provider -> repository | model, dimensions, bytes | 구현 예정 |
| `/api/chat` -> retriever | message, document ids | 구현 예정 |
| retriever -> provider context | citations and short context text | 구현 예정 |

## 7. 보안 고려사항

- provider context에는 retrieved chunk text만 제한적으로 포함한다.
- citations는 document id, filename, chunk id만 노출한다.
- BYOKnowledge는 사용하지 않고 CYOps-owned storage/retrieval path를 유지한다.
- mock embedding은 production model로 오인되지 않도록 model id를 명시한다.

## 8. 완료 기준

- [x] deterministic embedding mock test가 있다.
- [x] PostgreSQL embedding persistence integration test가 있다.
- [x] `/api/chat` RAG path가 retriever를 호출하고 citations를 반환한다.
- [x] v0.0.10 pgvector/retrieval handoff가 문서화된다.
