# v0.0.31 Phase 3 - CRC Subscription, InstallPlan, And CSV Smoke

## Work

- [x] Documented local CRC apply commands.
- [x] Kept Subscription `installPlanApproval` as `Manual`.
- [x] Preserved real-server deployment as a later step.
- [x] Applied the local CRC overlay in isolated namespace `cywell-opsmate-olm`.
- [x] Reviewed and approved generated InstallPlan `install-xb9gw`.
- [x] Confirmed CSV `cywell-opsmate.v0.0.29` reached `Succeeded`.

## Verification

- `oc apply --dry-run=client -k deploy/olm/local-crc`
- `kubectl kustomize deploy/olm/local-crc | oc apply -f -`
- `oc get installplan -n cywell-opsmate-olm -o yaml`
- `oc patch installplan install-xb9gw -n cywell-opsmate-olm --type=merge --patch '{\"spec\":{\"approved\":true}}'`
- `oc get subscription,installplan,csv,deploy,pod -n cywell-opsmate-olm`

Observed CRC evidence:

- Subscription generated an InstallPlan for `cywell-opsmate.v0.0.29`.
- Approved InstallPlan installed CSV `cywell-opsmate.v0.0.29`.
- Manager deployment initially reached `1/1` with pod `Running`.

## Remaining Scope

- Upgrade the installed CSV to v0.0.31 after the manager, bundle, and catalog images are published.

## Linked Issue

- #153
