# v0.0.28 Phase 2 - Manager Image Publish Workflow

## Work

- [x] Added `.github/workflows/manager-image.yml`.
- [x] Configured publish of `v0.0.28` and `latest`.
- [x] Scoped workflow permissions to `contents: read` and `packages: write`.

## Verification

- Workflow syntax is committed for post-merge execution.
- Local manager image build validates the same Containerfile path.

## Remaining Scope

- Confirm GitHub Actions run after merge.

## Linked Issue

- #136
