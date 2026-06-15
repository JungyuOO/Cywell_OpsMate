# v0.0.15 Planner

## 1. Goal

- Make OpenShift pgvector deployment configuration explicit in `OpsMateConfig`.
- Wire appserver PostgreSQL DSN through a Secret-backed environment variable.
- Keep pgvector migration operator-approved instead of implicit reconciliation.
- Document the OpenShift validation runbook for pgvector startup, migration, RAG smoke, and diagnostics.

## 2. Architecture Overview

1. `OpsMateConfig.spec.database.image` selects the PostgreSQL image.
2. If pgvector is required and no image is specified, the managed PostgreSQL Deployment defaults to `pgvector/pgvector:pg16`.
3. `OpsMateConfig.spec.database.dsnSecretRef` and `dsnSecretKey` inject `CYOPS_POSTGRES_DSN` into the appserver.
4. `CYOPS_PGVECTOR_REQUIRED=true` and `CYOPS_RETRIEVAL_MODE=pgvector` remain the startup and retrieval activation switches.
5. Migration remains a controlled admin action using the v0.0.14 helper/runbook path.

## 3. Tech Stack

| Area | Choice | Reason |
| --- | --- | --- |
| DB image | CR-controlled image with pgvector default | lets OpenShift deploy a vector-enabled DB |
| DSN | SecretKeyRef | avoids embedding credentials in the CR or logs |
| Migration | operator-approved runbook | avoids destructive automatic schema conversion |
| Diagnostics | `/api/ops/retrieval-metrics` | confirms pgvector retrieval during smoke |

## 4. Implementation Steps

| Phase | Scope | Status | Notes |
| --- | --- | --- | --- |
| Phase 1 | database spec and Postgres image selection | done | custom image and pgvector default |
| Phase 2 | appserver DSN Secret wiring | done | `CYOPS_POSTGRES_DSN` from SecretKeyRef |
| Phase 3 | CRD/sample/runbook | done | OpenShift pgvector validation docs |
| Phase 4 | v0.0.16 handoff | done | operator automation next |

## Linked Issues

- Phase 1: #70
- Phase 2: #71
- Phase 3: #72
- Phase 4: #73

## 5. Migration Or Operation Strategy

- Reconciliation does not apply `VECTOR(n)` migration automatically.
- Admin prepares backup and confirms embedding dimensions.
- Admin applies controlled migration and then re-embeds documents.
- Startup validation uses `CYOPS_PGVECTOR_REQUIRED=true` to fail early if the extension is missing.
- Rollback for dev/test is to return to `retrievalMode=bytea` and restore a pre-migration database backup.

## 6. Message / Communication / Data Protocol

| Path | Payload | Status |
| --- | --- | --- |
| `OpsMateConfig.spec.database.image` -> Postgres Deployment | container image | done |
| `OpsMateConfig.spec.database.dsnSecretRef/key` -> appserver | `CYOPS_POSTGRES_DSN` SecretKeyRef | done |
| `OpsMateConfig.spec.embedding.requirePGVector` -> appserver | `CYOPS_PGVECTOR_REQUIRED` | existing |
| `OpsMateConfig.spec.embedding.retrievalMode` -> appserver | `CYOPS_RETRIEVAL_MODE` | existing |

## 7. Security Considerations

- DSN is referenced from a Secret and is not stored directly in `OpsMateConfig`.
- Migration runbook requires backup and explicit approval before schema conversion.
- Metrics remain aggregate-only and must not include prompts, chunks, tokens, filenames, or DSNs.
- Custom database images should be pinned and reviewed before production rollout.

## 8. Completion Criteria

- [x] PostgreSQL image can be configured from `OpsMateConfig`.
- [x] pgvector-required configs default managed PostgreSQL to a pgvector-enabled image.
- [x] Appserver receives PostgreSQL DSN through SecretKeyRef.
- [x] CRD and sample include pgvector OpenShift fields.
- [x] Migration runbook and v0.0.16 handoff are documented.
