# v0.0.16 Planner

## 1. Goal

- Add an explicit admin approval field for pgvector migration.
- Expose operator status conditions for DSN configuration, pgvector requirement, pgvector readiness, migration approval, and retrieval mode readiness.
- Keep reconciliation from applying destructive pgvector migration automatically.
- Document the next re-embedding and console diagnostics work.

## 2. Architecture Overview

1. `OpsMateConfig.spec.database.pgVectorMigrationApproved` records admin approval intent.
2. Reconciliation still only manages Kubernetes resources and status; it does not execute `ApplyPGVectorEmbeddingMigration`.
3. Status conditions summarize config readiness without logging DSNs, tokens, prompts, or document text.
4. `PGVectorReady` is `Unknown` when runtime DB verification must happen in appserver/live smoke.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| Approval surface | CR spec boolean | simple explicit admin-controlled contract |
| Status | Kubernetes conditions | standard operator UX and OpenShift visibility |
| Runtime DB check | appserver startup/live smoke | avoids controller-side DSN reads in this version |
| Migration execution | future explicit job/tool | avoids silent destructive schema conversion |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | migration approval field | done | `pgVectorMigrationApproved` |
| Phase 2 | readiness status conditions | done | DSN, pgvector, retrieval mode |
| Phase 3 | CRD/sample/runbook update | done | approval field documented |
| Phase 4 | v0.0.17 handoff | done | re-embedding workflow next |

## Linked Issues

- Phase 1: #75
- Phase 2: #76
- Phase 3: #77
- Phase 4: #78

## 5. Migration Or Operation Strategy

- `pgVectorMigrationApproved=false` is the default and reconciliation will not migrate.
- `pgVectorMigrationApproved=true` only records approval; it still does not run migration in this version.
- Future migration automation must check this field before creating a migration job.
- Re-embedding must happen after migration because fallback BYTEA embeddings are reset.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| CR spec -> status | `PostgresDSNConfigured` | done |
| CR spec -> status | `PGVectorRequired` | done |
| CR spec -> status | `PGVectorMigrationApproved` | done |
| CR spec/runtime boundary -> status | `PGVectorReady` | done |
| CR spec -> status | `RetrievalModeReady` | done |

## 7. Security Considerations

- Conditions never include DSN values, tokens, prompts, filenames, or document text.
- The controller does not read database credentials in this version.
- Migration approval is visible in the CR, but migration execution remains separate and explicit.

## 8. Completion Criteria

- [x] Migration approval field exists in API, CRD, and sample.
- [x] Operator status exposes DSN and pgvector readiness conditions.
- [x] Invalid pgvector retrieval config degrades status.
- [x] Reconciliation still does not apply pgvector migration automatically.
- [x] v0.0.17 handoff captures re-embedding workflow scope.
