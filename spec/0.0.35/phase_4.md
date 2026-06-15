# v0.0.35 Phase 4 - Workload Readiness

## Work

- [x] Re-triggered the existing `cyops` `OpsMateConfig` reconcile.
- [x] Verified `cyops-appserver-tls` serving cert Secret was created.
- [x] Verified appserver deployment reached `1/1`.
- [x] Verified PostgreSQL deployment reached `1/1` under restricted SCC.
- [x] Captured controller and workload logs.

## Verification

- `oc annotate opsmateconfig cyops -n cywell-opsmate-olm cyops.cywell.io/reconcile-request=20260615174455 --overwrite` returned `opsmateconfig.opsmate.cywell.io/cyops annotated`.
- `oc get deploy,svc,secret,pod -n cywell-opsmate-olm` showed `cyops-appserver`, `cyops-postgres`, and `cywell-opsmate-controller-manager` all `1/1`.
- `oc get deploy,svc,secret,pod -n cywell-opsmate-olm` showed `secret/cyops-appserver-tls kubernetes.io/tls 2`.
- `oc get service cyops-appserver -n cywell-opsmate-olm -o yaml` preserved `service.beta.openshift.io/serving-cert-secret-name: cyops-appserver-tls`.
- `oc logs deployment/cyops-appserver -n cywell-opsmate-olm --tail=100` returned `cyops appserver listening on https://:8443`.
- `oc logs deployment/cyops-postgres -n cywell-opsmate-olm --tail=80` showed `database system is ready to accept connections`.
- `oc get opsmateconfig cyops -n cywell-opsmate-olm -o yaml` showed `status.overallStatus: Ready`.
- Controller logs showed one transient update conflict on `cyops-postgres`; resources reached Ready afterward.

## Remaining Scope

- PostgreSQL still logs a non-fatal `chmod: /var/run/postgresql: Operation not permitted` warning before startup. v0.0.36 should decide whether to clean this up or leave it as acceptable restricted-SCC noise.

## Linked Issue

- #174
