# v0.0.34 Phase 1 - CatalogSource v0.0.33 Refresh

## Work

- [x] Applied v0.0.33 CatalogSource to CRC.
- [x] Confirmed CatalogSource READY.

## Verification

- `oc whoami --show-server`
- `oc apply -f deploy\olm\install\catalogsource.yaml`
- `oc get catalogsource cywell-opsmate-catalog -n openshift-marketplace -o jsonpath="{.spec.image}{' '}{.status.connectionState.lastObservedState}"`

Observed evidence:

- Active server: `https://api.crc.testing:6443`
- CatalogSource image: `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.33`
- CatalogSource state: `READY`

## Remaining Scope

- None for CatalogSource refresh.

## Linked Issue

- #166
