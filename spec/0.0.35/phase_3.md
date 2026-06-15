# v0.0.35 Phase 3 - CRC OLM Upgrade

## Work

- [x] Located InstallPlan `install-nvtvs` for `cywell-opsmate.v0.0.34`.
- [x] Approved the manual InstallPlan.
- [x] Confirmed CSV `cywell-opsmate.v0.0.34` reached `Succeeded`.
- [x] Confirmed manager deployment uses `ghcr.io/jungyuoo/cywell-opsmate:v0.0.34`.

## Verification

- `oc patch installplan install-nvtvs -n cywell-opsmate-olm --type merge --patch-file .tmp-installplan-approve.json` returned `installplan.operators.coreos.com/install-nvtvs patched`.
- `oc get installplan install-nvtvs -n cywell-opsmate-olm -o jsonpath='{.spec.approved}'` returned `true`.
- `oc get csv -n cywell-opsmate-olm` showed `cywell-opsmate.v0.0.34` with phase `Succeeded`.
- `oc get pods -n cywell-opsmate-olm -o wide` showed `cywell-opsmate-controller-manager-68bd8f8d57-pgp6t 1/1 Running`.
- Deployment inspection showed manager image `ghcr.io/jungyuoo/cywell-opsmate:v0.0.34`.

## Remaining Scope

- Re-trigger smoke reconcile and verify workload readiness.

## Linked Issue

- #171
