# v0.0.37 Phase 1 - Artifact And CRC Readiness

## Work

- [x] Confirmed v0.0.36 GitHub Actions publish status.
- [x] Confirmed CRC API responds.
- [x] Confirmed v0.0.36 image manifests exist.
- [x] Confirmed CRC current CSV is `cywell-opsmate.v0.0.34`.
- [x] Confirmed `cyops-console` remains enabled in Console operator `spec.plugins`.

## Verification

- `gh run list --limit 4 --json workflowName,status,conclusion,headBranch,createdAt` showed v0.0.36 appserver, manager, bundle, and catalog workflows completed successfully.
- `oc get clusterversion` responded with OpenShift `4.21.8`.
- `docker manifest inspect` succeeded for:
  - `ghcr.io/jungyuoo/cywell-opsmate:v0.0.36`
  - `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.36`
  - `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.36`
  - `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.36`
- `oc get csv,installplan,subscription -n cywell-opsmate-olm` showed current CSV `cywell-opsmate.v0.0.34`.
- `oc get console.operator.openshift.io cluster -o jsonpath='{.spec.plugins}'` returned `["networking-console-plugin","monitoring-plugin","cyops-console"]`.

## Remaining Scope

- Apply the v0.0.37 graph fix because the v0.0.36 catalog did not produce an InstallPlan from the current `v0.0.34` install.

## Linked Issue

- #181
