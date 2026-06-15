# v0.0.33 Phase 1 - CatalogSource v0.0.32 Refresh

## Work

- [x] Applied v0.0.32 CatalogSource to CRC.
- [x] Confirmed CatalogSource READY.

## Verification

- `oc apply -f deploy\olm\install\catalogsource.yaml`
- `oc get catalogsource cywell-opsmate-catalog -n openshift-marketplace -o jsonpath="{.spec.image}{' '}{.status.connectionState.lastObservedState}"`

Observed evidence:

- CatalogSource image: `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.32`
- CatalogSource state: `READY`

## Remaining Scope

- None for CatalogSource refresh.

## Linked Issue

- #160
