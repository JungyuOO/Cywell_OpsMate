# Phase 1 - Retrieval Metrics Aggregate

## Scope

- [x] Add a concrete sink for `RetrievalObserver`.
- [x] Track total retrievals, slow retrievals, failures, mode counts, failure reason counts, average duration, and the last observation summary.
- [x] Keep the sink free of prompt text, chunk text, document metadata, tokens, and DSNs.
- [x] Link GitHub Issue #60.

## Work Completed

- Added `RetrievalMetrics` and `RetrievalMetricsSnapshot`.
- Added a nil-safe `ObserveRetrieval` implementation.
- Added snapshot cloning so callers cannot mutate internal state.

## Verification

- `TestRetrievalMetricsAggregatesObservations` covers counts, slow flag, failure reason count, average duration, and last observation.

## Remaining Scope

- Prometheus/OpenTelemetry export remains a later production-hardening task.
