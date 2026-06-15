# v0.0.31 Phase 2 - CRC Catalog Verification

## Work

- [x] Reused `deploy/olm/install/catalogsource.yaml`.
- [x] Kept catalog source in `openshift-marketplace`.
- [x] Applied CatalogSource to CRC and confirmed registry readiness.

## Verification

- `oc apply --dry-run=client -f deploy/olm/install/catalogsource.yaml`
- `oc apply -f deploy/olm/install/catalogsource.yaml`
- `oc get catalogsource cywell-opsmate-catalog -n openshift-marketplace -o yaml`
- `oc get pods -n openshift-marketplace | Select-String cywell`

Observed CRC evidence:

- CatalogSource `cywell-opsmate-catalog` reported connection state `READY`.
- Registry pod `cywell-opsmate-catalog-dwpln` reported `1/1 Running`.

## Remaining Scope

- Publish and apply the v0.0.31 CatalogSource image before upgrade testing.

## Linked Issue

- #152
