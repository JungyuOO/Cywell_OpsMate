# Phase 1 - Re-embedding Batch Workflow

## Scope

- [x] Add ready document selection for re-embedding.
- [x] Add explicit batch re-embedding service method.
- [x] Preserve failure visibility through `embedding_status=failed` and `last_error`.
- [x] Link GitHub Issue #80.

## Work Completed

- Added `ListReadyDocumentsForReembedding`.
- Added `EmbeddingService.ReembedReadyDocuments`.
- Added `ReembeddingRequest` and `ReembeddingResult`.

## Verification

- `CYOPS_POSTGRES_TEST_DSN=postgres://cyops:cyops@localhost:55432/cyops?sslmode=disable go test ./internal/appserver -count=1`

## Remaining Scope

- Add an admin API or operator-triggered Job to invoke re-embedding in cluster.
