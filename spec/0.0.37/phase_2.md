# v0.0.37 Phase 2 - Catalog Graph Fix

## Work

- [x] Applied v0.0.36 CatalogSource to CRC and confirmed no v0.0.36 InstallPlan was produced.
- [x] Identified the channel graph break: v0.0.36 replaced `v0.0.33`, but CRC was installed at `v0.0.34`.
- [x] Updated v0.0.37 packaging references.
- [x] Updated catalog head to `cywell-opsmate.v0.0.37 replaces cywell-opsmate.v0.0.34`.
- [x] Ran tests and manifest dry-run validation.

## Verification

- `oc apply -f deploy\olm\install\catalogsource.yaml` configured the v0.0.36 CatalogSource.
- `oc get installplan -n cywell-opsmate-olm` showed no new InstallPlan after catalog refresh.
- `oc get subscription cywell-opsmate -n cywell-opsmate-olm -o yaml` stayed at `state: AtLatestKnown` with `currentCSV: cywell-opsmate.v0.0.34`.
- `oc logs -n openshift-marketplace -l olm.catalogSource=cywell-opsmate-catalog --tail=80` showed the v0.0.36 catalog served successfully, so the issue was graph topology rather than catalog pod readiness.
- `go test ./...` passed after the v0.0.37 graph/version update.
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests` passed for CSV `cywell-opsmate.v0.0.37` and the CRD.
- `oc apply --dry-run=client -f deploy\olm\install\catalogsource.yaml` passed.
- `kubectl kustomize deploy/olm/local-crc` still failed in this Codex sandbox with `evalsymlink failure ... Access is denied`.

## Remaining Scope

- Merge v0.0.37, wait for image/catalog publish, then apply v0.0.37 CatalogSource to CRC and approve the resulting InstallPlan.

## Linked Issue

- pending
