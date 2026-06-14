# v0.0.8 Planner

## 1. 목표

- 업로드된 문서를 ingestion 대상으로 처리해 `cyops_document_chunks`에 chunk를 저장한다.
- plain text와 markdown parser를 우선 구현한다.
- ingestion 성공/실패 상태를 `cyops_documents.status`, `embedding_status`, `last_error`에 반영한다.
- embedding provider와 vector retrieval은 interface와 다음 버전 handoff까지만 정리한다.

## 2. 아키텍처 개요

1. appserver 또는 worker가 document metadata의 `object_uri`를 조회한다.
2. parser가 원본 파일을 읽어 text section을 만든다.
3. chunker가 deterministic chunk로 section text를 나눈다.
4. repository가 기존 chunks를 교체하고 document 상태를 `ready` 또는 `failed`로 갱신한다.
5. embedding provider는 다음 버전에서 chunk 단위로 호출한다.

## 3. 기술 스택

| 영역 | 선택 | 이유 |
| --- | --- | --- |
| Parser | Go interface + plain text/markdown 구현 | dependency 없이 최소 RAG path 검증 |
| Chunker | deterministic rune chunker | 테스트 가능하고 tokenizer 교체가 쉬움 |
| Persistence | PostgreSQL `cyops_document_chunks` | v0.0.5 schema와 일관 |
| Worker | appserver 내부 service boundary | 별도 Deployment 전 단계 |

## 4. 구현 단계

| Phase | 범위 | 상태 | 산출물 |
| --- | --- | --- | --- |
| Phase 1 | parser/chunker boundary | 완료 | parser, chunker, tests |
| Phase 2 | PostgreSQL chunk repository | 완료 | chunk insert/delete/status methods |
| Phase 3 | ingestion service flow | 완료 | ingest document orchestration |
| Phase 4 | v0.0.9 embedding/retrieval handoff | 완료 | next scope docs |

## 5. 마이그레이션 또는 운영 전략

- 기존 `cyops_document_chunks` table을 사용한다.
- ingestion 재시도 시 같은 document의 chunks를 삭제 후 다시 insert한다.
- 실패 상태는 document row에 기록하고 원본 파일은 삭제하지 않는다.

## 6. 메시지/통신/데이터 프로토콜

| 경로 | 데이터 | 상태 |
| --- | --- | --- |
| repository -> ingestion | document id, filename, object URI | 구현 예정 |
| parser -> chunker | parsed sections | 구현 예정 |
| chunker -> repository | chunk index, text, source offsets | 구현 예정 |
| ingestion -> repository | processing/ready/failed status | 구현 예정 |

## 7. 보안 고려사항

- ingestion은 repository가 반환한 object URI만 읽는다.
- path traversal은 upload storage 단계에서 차단되며 ingestion은 추가로 빈 URI를 거부한다.
- parser error는 API 응답에 직접 노출하지 않고 `last_error`에 제한적으로 기록한다.

## 8. 완료 기준

- [x] plain text/markdown parser와 chunker tests가 있다.
- [x] PostgreSQL repository가 chunks 저장과 ingestion 상태 갱신을 지원한다.
- [x] Docker PostgreSQL integration test가 ingestion flow를 검증한다.
- [x] v0.0.9 embedding/retrieval handoff가 문서화된다.
