# v0.0.13 Planner

## 1. Goal

- Connect retrieval observations to a concrete backend metrics sink.
- Expose a safe operator-facing metrics snapshot without prompt, chunk text, token, or DSN data.
- Document the pgvector live rollout boundary and keep OpenShift live validation as the next operational step when a pgvector-enabled database is available.

## 2. Architecture Overview

1. `PostgresRetriever` emits `RetrievalObservation` for bytea and pgvector modes.
2. `RetrievalMetrics` aggregates counts, slow retrievals, failure reasons, result counts, and average duration.
3. `NewServerFromConfig` wires the same metrics sink into the retriever and the HTTP server.
4. `/api/ops/retrieval-metrics` returns a JSON snapshot for operator-side diagnostics.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| Metrics sink | in-process Go aggregate | avoids new dependencies before final metrics backend selection |
| API | JSON endpoint | easy smoke testing and console integration |
| Retrieval | existing observer interface | keeps retrieval mode code unchanged |
| pgvector rollout | documented handoff | OpenShift live DB validation depends on target cluster resources |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | retrieval metrics aggregate | done | thread-safe snapshot and tests |
| Phase 2 | appserver endpoint and config wiring | done | `/api/ops/retrieval-metrics` |
| Phase 3 | pgvector rollout boundary | done | runbook and next validation scope |
| Phase 4 | v0.0.14 handoff | done | live pgvector/OpenShift smoke next |

## Linked Issues

- Phase 1: #60
- Phase 2: #61
- Phase 3: #62
- Phase 4: #63

## 5. Migration Or Operation Strategy

- The new metrics sink is in-memory and resets on pod restart.
- It is safe for dev and early operator diagnostics.
- Production Prometheus/OpenTelemetry export remains a follow-up once the console/operator deployment shape is fixed.
- pgvector live validation still requires a PostgreSQL image with the `vector` extension installed in the target OpenShift namespace.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| `RetrievalObserver` -> `RetrievalMetrics` | mode, duration, slow flag, failure reason, result count | done |
| `GET /api/ops/retrieval-metrics` | aggregate JSON snapshot | done |
| appserver -> PostgreSQL pgvector | SQL top-k query from v0.0.12 | unchanged |

## 7. Security Considerations

- Metrics do not include user prompts, uploaded document text, chunk text, tokens, DSNs, filenames, or object URIs.
- Failure reasons are controlled internal strings.
- The endpoint is intended to sit behind the same OpenShift console/plugin auth boundary as the rest of the appserver.
- A cluster-facing production metrics endpoint should add RBAC and scrape policy review before release.

## 8. Completion Criteria

- [x] Retrieval observations are aggregated by a concrete sink.
- [x] Slow retrieval and failure reason counts are observable.
- [x] Appserver exposes a JSON metrics snapshot.
- [x] Unit tests cover aggregation and endpoint behavior.
- [x] OpenShift live pgvector validation is explicitly handed off to the next version.
