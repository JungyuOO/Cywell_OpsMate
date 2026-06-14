# CYOps Appserver Migrations

This directory holds SQL skeletons for the CYOps appserver data model.

v0.0.5 records schema shape only. Runtime migration application and PostgreSQL
repository wiring are deferred to v0.0.6.

`0001_cyops_rag_schema.sql` intentionally uses `BYTEA` for the embedding
placeholder. Replace it with a pgvector `VECTOR(n)` column when the embedding
model and dimensions are fixed.
