# v0.0.35 Phase 2 - CRC CatalogSource Refresh

## Work

- [x] Applied `deploy/olm/install/catalogsource.yaml` to CRC.
- [x] Confirmed CatalogSource moved from `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.33` to `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.34`.
- [x] Confirmed the new CatalogSource pod was `1/1 Running`.
- [x] Confirmed OLM created InstallPlan `install-nvtvs` for `cywell-opsmate.v0.0.34`.

## Verification

- `oc whoami --show-server` returned `https://api.crc.testing:6443`.
- `oc apply -f deploy\olm\install\catalogsource.yaml` returned `catalogsource.operators.coreos.com/cywell-opsmate-catalog configured`.
- `oc get catalogsource cywell-opsmate-catalog -n openshift-marketplace -o jsonpath='{.spec.image}'` returned `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.34`.
- `oc get pods -n openshift-marketplace -l olm.catalogSource=cywell-opsmate-catalog` showed `cywell-opsmate-catalog-6vdcl 1/1 Running`.
- `oc get installplan -n cywell-opsmate-olm` showed `install-nvtvs cywell-opsmate.v0.0.34 Manual false`.

## Remaining Scope

- Approve and verify the v0.0.34 InstallPlan.

## Linked Issue

- #170
