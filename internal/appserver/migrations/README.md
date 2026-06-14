# CYOps Appserver Migrations

This directory holds SQL skeletons for the CYOps appserver data model.

v0.0.7 applies these SQL files through the appserver migration runner when a
PostgreSQL DSN is configured. Keep SQL idempotent because startup may execute it
more than once.

`0001_cyops_rag_schema.sql` intentionally uses `BYTEA` for the embedding
placeholder. Replace it with a pgvector `VECTOR(n)` column when the embedding
model and dimensions are fixed.
