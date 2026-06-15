# v0.0.34 Phase 3 - ConsolePlugin Verification And Pod Blocker Fix

## Work

- [x] Confirmed ConsolePlugin update no longer fails with missing `resourceVersion`.
- [x] Confirmed manager logs stay clean after a reconcile trigger.
- [x] Confirmed ConsolePlugin points to `cyops-appserver` in `cywell-opsmate-olm`.
- [x] Captured next blockers: appserver waits for missing serving cert Secret and PostgreSQL fails under restricted SCC due unwritable data/run directories.
- [x] Fixed Service reconcile to preserve annotations, including OpenShift serving-cert annotation.
- [x] Added writable `emptyDir` mounts, `PGDATA`, and restricted-compatible security context to PostgreSQL.
- [x] Bumped manager, bundle, and catalog references to v0.0.34.

## Verification

- `oc apply -k deploy\olm\local-crc`
- `oc annotate opsmateconfig cyops -n cywell-opsmate-olm cyops.cywell.io/reconcile-request=<timestamp> --overwrite`
- `oc logs deployment/cywell-opsmate-controller-manager -n cywell-opsmate-olm --tail=160`
- `oc get consoleplugin cyops-console -o jsonpath="{.metadata.resourceVersion}{' '}{.spec.displayName}{' '}{.spec.backend.service.name}{' '}{.spec.backend.service.namespace}"`
- `oc get deploy,svc,pod,opsmateconfig -n cywell-opsmate-olm`
- `oc logs deploy/cyops-postgres -n cywell-opsmate-olm --previous --tail=120`
- `oc describe pod -l app.kubernetes.io/component=postgres -n cywell-opsmate-olm`
- `oc describe pod -l app.kubernetes.io/component=appserver -n cywell-opsmate-olm`
- `$env:GOCACHE = (Join-Path (Get-Location) '.cache\go-build'); go test ./...`
- `kubectl kustomize config/default`
- `oc apply --dry-run=client -f deploy\olm\bundle\manifests`
- `oc apply --dry-run=client -k deploy\olm\local-crc`
- `docker build -f deploy\containerfiles\manager.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate:v0.0.34 .`
- `docker build -f deploy\containerfiles\bundle.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.34 .`
- `docker build -f deploy\containerfiles\catalog.Containerfile -t ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.34 .`
- `docker run --rm -v ${PWD}\deploy\olm\catalog:/configs:ro quay.io/operator-framework/opm:latest validate /configs`

Observed evidence:

- Manager logs had startup entries only and no ConsolePlugin `resourceVersion` errors after reconcile trigger.
- ConsolePlugin state: `CYOps cyops-appserver cywell-opsmate-olm`.
- PostgreSQL previous logs showed `initdb: error: could not change permissions of directory "/var/lib/postgresql/data": Operation not permitted`.
- Appserver pod events showed `MountVolume.SetUp failed for volume "serving-cert" : secret "cyops-appserver-tls" not found`.

## Remaining Scope

- Publish v0.0.34 images, upgrade CRC to `cywell-opsmate.v0.0.34`, and verify PostgreSQL/appserver pods become ready.

## Linked Issue

- #167
