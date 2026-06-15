# Local CRC OLM Install

This path is for local CRC validation before any real-server deployment.

It intentionally installs the Subscription into `cywell-opsmate-olm` so it does
not conflict with the development bootstrap namespace `cywell-opsmate-system`.

Run in order:

```powershell
oc apply -f deploy/olm/install/catalogsource.yaml
kubectl kustomize deploy/olm/local-crc | oc apply -f -
oc get catalogsource cywell-opsmate-catalog -n openshift-marketplace
oc get subscription,installplan,csv -n cywell-opsmate-olm
```

Approve the generated InstallPlan only after confirming it references the CYOps
CSV and expected images.

After the CSV succeeds, the same overlay applies a local `OpsMateConfig` and
test PostgreSQL password Secret so the reconciler can create appserver,
PostgreSQL, and ConsolePlugin resources without depending on customer secrets.
