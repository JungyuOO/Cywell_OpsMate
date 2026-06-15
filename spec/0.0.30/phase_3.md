# v0.0.30 Phase 3 - Local OLM Install Manifests

## Work

- [x] Added `deploy/olm/install/catalogsource.yaml`.
- [x] Updated OLM install README to apply CatalogSource before Subscription.
- [x] Kept `installPlanApproval: Manual`.

## Verification

- `oc apply --dry-run=client -f deploy/olm/install/catalogsource.yaml`
- `oc apply --dry-run=client -f deploy/olm/install/namespace.yaml -f deploy/olm/install/operatorgroup.yaml -f deploy/olm/install/subscription.yaml`

## Remaining Scope

- Apply these manifests to CRC after catalog image publication.

## Linked Issue

- #147
