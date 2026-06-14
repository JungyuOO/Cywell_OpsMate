# Phase 2 - Appserver Metrics Endpoint

## Scope

- [x] Wire `RetrievalMetrics` into appserver options.
- [x] Connect `NewServerFromConfig` so PostgreSQL-backed retrievers publish observations into the server metrics sink.
- [x] Expose `GET /api/ops/retrieval-metrics`.
- [x] Link GitHub Issue #61.

## Work Completed

- `ServerOptions` now accepts a metrics sink.
- The appserver defaults to an empty sink when no custom sink is provided.
- PostgreSQL config path shares one sink between `PostgresRetriever` and the HTTP endpoint.
- Added a JSON endpoint for operator diagnostics.

## Verification

- `TestRetrievalMetricsEndpointReturnsSnapshot` verifies endpoint output for a recorded pgvector observation.

## Remaining Scope

- Endpoint authorization depends on the final OpenShift ConsolePlugin/backend deployment boundary.
