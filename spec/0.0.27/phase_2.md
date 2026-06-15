# v0.0.27 Phase 2 - GHCR Image Publish Workflow

## Work

- [x] Built the appserver image from `deploy/containerfiles/appserver.Containerfile`.
- [x] Tagged the image as `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27`.
- [x] Added `.github/workflows/appserver-image.yml` to publish `v0.0.27` and `latest` with GitHub Actions `packages: write`.

## Verification

- `docker build -f deploy/containerfiles/appserver.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27 .`
- Local `docker push ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.27` was attempted but the available local tokens did not have `write:packages`; publish is delegated to the repository workflow.

## Remaining Scope

- Confirm the publish workflow run after merge.

## Linked Issue

- #131
