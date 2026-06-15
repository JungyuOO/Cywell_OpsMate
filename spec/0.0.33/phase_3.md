# v0.0.33 Phase 3 - OpsMateConfig Reconcile Retry And ConsolePlugin Fix

## Work

- [x] Confirmed v0.0.32 manager starts and no longer logs cluster-scope `OpsMateConfig` RBAC errors.
- [x] Confirmed appserver, PostgreSQL, service, and ConsolePlugin resources are created.
- [x] Captured the next blocker: repeated ConsolePlugin update fails because unstructured reconcile did not preserve existing `resourceVersion`.
- [x] Fixed unstructured reconcile to create an empty current object with the desired GVK and update only metadata/spec.
- [x] Added a regression test for updating an existing ConsolePlugin.
- [x] Bumped manager, bundle, and catalog references to v0.0.33.

## Verification

- `oc get deploy,svc,pod,opsmateconfig -n cywell-opsmate-olm`
- `oc get consoleplugin cyops-console -o yaml`
- `oc logs deployment/cywell-opsmate-controller-manager -n cywell-opsmate-olm --tail=120`
- `$env:GOCACHE = (Join-Path (Get-Location) '.cache\go-build'); go test ./...`
- `kubectl kustomize config/default`
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests`
- `oc apply --dry-run=client -k deploy\olm\local-crc`
- `docker build -f deploy\containerfiles\manager.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate:v0.0.33 .`
- `docker build -f deploy\containerfiles\bundle.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.33 .`
- `docker build -f deploy\containerfiles\catalog.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.33 .`
- `docker run --rm -v ${PWD}\deploy\olm\catalog:/configs:ro quay.io/operator-framework/opm:latest validate /configs`

Observed evidence:

- Manager pod `cywell-opsmate-controller-manager-5c45dbcd69-wxpnd` was `1/1 Running`.
- `cyops-appserver`, `cyops-postgres`, `cyops-appserver` service, `cyops-postgres` service, and `ConsolePlugin/cyops-console` were created.
- Logs showed `consoleplugins.console.openshift.io "cyops-console" is invalid: metadata.resourceVersion: Invalid value: 0: must be specified for an update`.

## Remaining Scope

- Publish v0.0.33 images, upgrade CRC to `cywell-opsmate.v0.0.33`, and verify ConsolePlugin update no longer errors.

## Linked Issue

- #163
