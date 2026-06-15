# v0.0.30 Phase 1 - Catalog Image Metadata

## Work

- [x] Added file-based catalog metadata for package `cywell-opsmate`.
- [x] Added alpha channel entry for `cywell-opsmate.v0.0.29`.
- [x] Linked the entry to `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.29`.

## Verification

- `docker build -f deploy/containerfiles/catalog.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.30 .`

## Remaining Scope

- Confirm catalog publish workflow after merge.

## Linked Issue

- #145
