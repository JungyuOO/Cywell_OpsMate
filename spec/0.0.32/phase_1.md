# v0.0.32 Phase 1 - CatalogSource Refresh

## Work

- [x] Applied the v0.0.31 CatalogSource manifest to CRC.
- [x] Confirmed the registry pod and CatalogSource were healthy.

## Verification

- `oc apply -f deploy\olm\install\catalogsource.yaml`
- `oc get catalogsource cywell-opsmate-catalog -n openshift-marketplace -o jsonpath="{.spec.image}{' '}{.status.connectionState.lastObservedState}"`

Observed evidence:

- CatalogSource image: `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.31`
- CatalogSource state: `READY`

## Remaining Scope

- None for v0.0.31 CatalogSource refresh.

## Linked Issue

- #155
