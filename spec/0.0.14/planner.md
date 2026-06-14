# v0.0.14 Planner

## 1. Goal

- Make pgvector activation executable as a controlled migration step.
- Add a live pgvector RAG smoke test that validates migration, SQL ranking, `/api/chat` citations, and retrieval metrics together.
- Keep bytea fallback tests intact for dev/test PostgreSQL instances that do not provide pgvector.

## 2. Architecture Overview

1. Base migrations create the RAG schema with BYTEA embedding storage.
2. `ApplyPGVectorEmbeddingMigration` validates pgvector readiness and applies the generated `VECTOR(n)` migration.
3. `CYOPS_PGVECTOR_TEST_DSN` gates the live pgvector smoke test.
4. The live smoke seeds ready documents, chunks, and vector literals, then calls `/api/chat` with `PostgresRetriever` in `pgvector` mode.
5. Retrieval metrics are checked during the same smoke flow.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| Migration | explicit Go helper | prevents accidental pgvector conversion in default migrations |
| Live smoke | env-gated integration test | runs only when a pgvector-enabled DB exists |
| Ranking | pgvector `<=>` SQL path | verifies production retrieval mode |
| Metrics | v0.0.13 JSON snapshot | proves retrieval observations during live RAG |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | controlled pgvector migration helper | done | validates dimensions before DB mutation |
| Phase 2 | pgvector live RAG smoke | done | `CYOPS_PGVECTOR_TEST_DSN` gated |
| Phase 3 | local pgvector verification | done | Docker `pgvector/pgvector:pg16` smoke passed |
| Phase 4 | v0.0.15 handoff | done | OpenShift deployment/runbook next |

## Linked Issues

- Phase 1: #65
- Phase 2: #66
- Phase 3: #67
- Phase 4: #68

## 5. Migration Or Operation Strategy

- Default `ApplyMigrations` remains BYTEA-safe.
- Operators apply pgvector conversion only after embedding dimensions are fixed.
- `ApplyPGVectorEmbeddingMigration` calls `CheckPGVectorReady` before applying `VECTOR(n)`.
- Existing BYTEA fallback data is intentionally reset by the generated migration path; re-embedding is required after conversion.
- Development clusters without pgvector continue to use `CYOPS_RETRIEVAL_MODE=bytea`.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| appserver -> PostgreSQL | `CREATE EXTENSION IF NOT EXISTS vector` readiness | done |
| migration helper -> PostgreSQL | `ALTER COLUMN embedding TYPE VECTOR(n)` | done |
| `/api/chat` RAG | pgvector-ranked citations | done |
| `/api/ops/retrieval-metrics` | pgvector retrieval count/result count | done |

## 7. Security Considerations

- The live smoke does not log prompts, chunk text beyond test assertions, DSNs, tokens, or uploaded customer documents.
- Migration helper validates dimensions before touching the database.
- Production migration execution still needs operator-controlled rollout approval and backup/rollback procedure.
- The metrics endpoint remains behind the appserver boundary and contains only aggregate fields.

## 8. Completion Criteria

- [x] Controlled pgvector migration helper exists.
- [x] pgvector live RAG smoke test is present and env-gated.
- [x] Local Docker pgvector smoke passed.
- [x] Existing non-pgvector PostgreSQL integration tests still pass.
- [x] v0.0.15 handoff captures OpenShift deployment validation.
