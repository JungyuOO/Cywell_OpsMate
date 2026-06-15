# v0.0.35 Phase 1 - Release Artifact Confirmation

## Work

- [x] Confirmed Issues #165, #166, #167, and #168 are closed after PR #169.
- [x] Confirmed main branch workflows for v0.0.34 manager, bundle, and catalog completed successfully.
- [x] Confirmed GHCR manifests exist for:
  - `ghcr.io/jungyuoo/cywell-opsmate:v0.0.34`
  - `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.34`
  - `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.34`

## Verification

- `gh issue list --state all --search "repo:JungyuOO/Cywell_OpsMate 165 OR 166 OR 167 OR 168" --json number,title,state`
- `gh run list --limit 12 --json databaseId,displayTitle,workflowName,status,conclusion,headBranch,createdAt`
- `docker manifest inspect ghcr.io/jungyuoo/cywell-opsmate:v0.0.34`
- `docker manifest inspect ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.34`
- `docker manifest inspect ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.34`

## Remaining Scope

- Refresh CRC CatalogSource to v0.0.34.
- Approve the resulting InstallPlan.
- Verify v0.0.34 workload readiness.

## Linked Issue

- #172
