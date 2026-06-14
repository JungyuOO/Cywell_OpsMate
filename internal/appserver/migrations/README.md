# CYOps Appserver Migrations

This directory holds SQL skeletons for the CYOps appserver data model.

v0.0.7 applies these SQL files through the appserver migration runner when a
PostgreSQL DSN is configured. Keep SQL idempotent because startup may execute it
more than once.

`0001_cyops_rag_schema.sql` intentionally keeps the embedding column as `BYTEA`
for the default dev/test path. v0.0.11 adds `PGVectorEmbeddingMigrationSQL` as
the explicit activation path for `VECTOR(n)` once the target dimensions and
pgvector extension readiness are confirmed.

Existing `BYTEA` fallback embeddings should be treated as reset/rebuild data
when moving to `VECTOR(n)`. The generated migration uses `USING NULL::VECTOR(n)`
so operators can apply it only after planning a re-embedding pass.
