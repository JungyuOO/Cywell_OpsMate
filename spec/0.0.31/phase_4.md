# v0.0.31 Phase 4 - OpsMateConfig Smoke Finding And Upgrade Handoff

## Work

- [x] Applied a sanitized local `OpsMateConfig` and PostgreSQL password Secret through the local CRC overlay.
- [x] Confirmed the existing manager image can crash during or after `OpsMateConfig` reconcile.
- [x] Confirmed controller-runtime logs were not initialized, leaving crash diagnostics incomplete.
- [x] Added zap logger initialization to the manager entrypoint.
- [x] Bumped manager, bundle, and catalog references to v0.0.31 so the next local OLM test can upgrade in place.
- [x] Kept local Web Console evidence as the next milestone after reconcile succeeds.

## Verification

- Handoff reviewed against `spec/0.0.31/planner.md`.
- `oc get opsmateconfig,deploy,svc,pod -n cywell-opsmate-olm`
- `oc logs deployment/cywell-opsmate-controller-manager -n cywell-opsmate-olm --previous`
- `oc auth can-i list opsmateconfigs.cyops.cywell.io --as system:serviceaccount:cywell-opsmate-olm:cywell-opsmate-controller-manager -n cywell-opsmate-olm`
- `oc auth can-i create deployments --as system:serviceaccount:cywell-opsmate-olm:cywell-opsmate-controller-manager -n cywell-opsmate-olm`
- `oc auth can-i create consoleplugins.console.openshift.io --as system:serviceaccount:cywell-opsmate-olm:cywell-opsmate-controller-manager`
- `$env:GOCACHE = (Join-Path (Get-Location) '.cache\go-build'); go test ./...`
- `kubectl kustomize config/default`
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests`
- `oc apply --dry-run=client -k deploy\olm\local-crc`
- `docker build -f deploy\containerfiles\manager.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate:v0.0.31 .`
- `docker build -f deploy\containerfiles\bundle.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.31 .`
- `docker build -f deploy\containerfiles\catalog.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.31 .`
- `docker run --rm -v ${PWD}\deploy\olm\catalog:/configs:ro quay.io/operator-framework/opm:latest validate /configs`

Observed CRC evidence:

- `cywell-opsmate.v0.0.29` installed successfully, but applying `OpsMateConfig/cyops` did not produce reconciled appserver, PostgreSQL, service, or ConsolePlugin resources.
- Manager pod entered `CrashLoopBackOff` after reconcile was triggered.
- Previous logs only showed controller-runtime logger was never initialized.
- ServiceAccount authorization checks returned `yes` for the required namespace and ConsolePlugin operations.

## Remaining Scope

- See `v0.0.32_handoff.md`.

## Linked Issue

- #151
