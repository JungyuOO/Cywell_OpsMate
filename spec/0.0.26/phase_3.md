# v0.0.26 Phase 3 - Appserver Image Packaging

## Work

- [x] Added `deploy/containerfiles/appserver.Containerfile`.
- [x] Built `/appserver` and `cyops-pgvector-migrate` in the same appserver image.
- [x] Added `.dockerignore` to keep local caches and secrets out of image build context.

## Verification

- `docker build -f deploy/containerfiles/appserver.Containerfile -t cyops-appserver:v0.0.26 .`

## Remaining Scope

- Publish the image to GHCR when live cluster deployment is ready.

## Linked Issue

- #127
