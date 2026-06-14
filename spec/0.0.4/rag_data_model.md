# CYOps RAG Data Model and Pipeline Contract

## 기본 결정

CYOps RAG는 BYOKnowledge를 사용하지 않는다. 고객 문서 ingestion, metadata, chunk, embedding, retrieval은 OpsMate appserver와 worker가 소유한다.

초기 구현은 PostgreSQL metadata와 pgvector를 기본 후보로 둔다. 별도 vector DB는 대규모 문서량, 고성능 검색, 멀티 테넌트 분리 요구가 확인될 때 옵션으로 분리한다.

## 저장소 경계

| 데이터 | 기본 저장 위치 | 설명 |
| --- | --- | --- |
| document metadata | PostgreSQL | 파일명, 크기, 상태, 업로드 사용자, scope |
| document original | PVC 또는 S3-compatible object store | 원본 파일. PostgreSQL에는 경로/키만 저장 |
| chunks | PostgreSQL | chunk text, token count, source offset |
| embeddings | PostgreSQL pgvector | chunk embedding vector |
| chat sessions | PostgreSQL | session/message metadata |
| audit events | PostgreSQL | upload/delete/query/admin actions |

## 테이블 초안

### cyops_documents

| column | type | notes |
| --- | --- | --- |
| id | uuid | primary key |
| namespace | text | tenant/scope boundary |
| filename | text | original filename |
| content_type | text | uploaded MIME type |
| size_bytes | bigint | file size |
| object_uri | text | PVC/object-store path |
| status | text | uploaded, processing, ready, failed, deleting |
| embedding_status | text | pending, processing, ready, failed |
| uploaded_by | text | OpenShift username |
| created_at | timestamptz | server timestamp |
| updated_at | timestamptz | server timestamp |
| deleted_at | timestamptz nullable | soft delete |
| last_error | text | sanitized error |

### cyops_document_chunks

| column | type | notes |
| --- | --- | --- |
| id | uuid | primary key |
| document_id | uuid | references cyops_documents(id) |
| chunk_index | integer | stable ordering |
| text | text | chunk content |
| token_count | integer | approximate token count |
| source_start | integer | optional source offset |
| source_end | integer | optional source offset |
| created_at | timestamptz | server timestamp |

### cyops_document_embeddings

| column | type | notes |
| --- | --- | --- |
| chunk_id | uuid | primary key, references cyops_document_chunks(id) |
| model | text | embedding model id |
| dimensions | integer | vector dimensions |
| embedding | vector | pgvector value |
| created_at | timestamptz | server timestamp |

### cyops_chat_sessions

| column | type | notes |
| --- | --- | --- |
| id | uuid | primary key |
| namespace | text | tenant/scope boundary |
| user_name | text | OpenShift username |
| created_at | timestamptz | server timestamp |
| updated_at | timestamptz | server timestamp |

### cyops_chat_messages

| column | type | notes |
| --- | --- | --- |
| id | uuid | primary key |
| session_id | uuid | references cyops_chat_sessions(id) |
| role | text | user, assistant, system |
| provider | text | lightspeed, local, external |
| content | text | sanitized content |
| citations_json | jsonb | cited chunks/documents |
| created_at | timestamptz | server timestamp |

## Ingestion pipeline

1. Upload request creates `cyops_documents(status=uploaded, embedding_status=pending)`.
2. Original file is stored in PVC/object store.
3. Ingestion worker marks document `processing`.
4. Parser extracts text.
5. Chunker writes `cyops_document_chunks`.
6. Embedding worker writes `cyops_document_embeddings`.
7. Document status becomes `ready`.
8. Query-time retrieval filters chunks by namespace, selected document IDs, and vector similarity.

## Embedding service boundary

초기에는 appserver 내부 worker로 시작할 수 있다. 다음 조건 중 하나가 생기면 별도 embedding Deployment로 분리한다.

- embedding 작업이 appserver 요청 latency에 영향을 준다.
- GPU 또는 별도 runtime이 필요하다.
- 고객별 embedding model 분리가 필요하다.
- 대량 문서 처리 queue가 필요하다.

## Provider policy

- Lightspeed provider에는 검색된 customer context 중 최소 필요 chunk만 전달한다.
- provider request에는 document ID와 chunk ID citation metadata를 포함하되 Secret, 원본 파일 경로, 전체 문서 원문은 포함하지 않는다.
- provider가 실패하면 CYOps는 RAG 검색 결과와 실패 reason을 분리해 UI에 표시한다.
