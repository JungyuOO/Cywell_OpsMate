# v0.0.27 Phase 1 - Version-Pinned Image References

## Work

- [x] Pinned appserver Deployment default image to `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27`.
- [x] Pinned pgvector migration Job image to the same appserver image tag.

## Verification

- `go test ./internal/controller/appserver ./internal/controller/postgres`

## Remaining Scope

- Future releases should update both image defaults together.

## Linked Issue

- #130
