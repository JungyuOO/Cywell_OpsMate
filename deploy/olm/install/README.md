# CYOps OLM Install

This directory describes the intended customer-facing install path.

1. Publish the bundle image:

   `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.50`

2. Publish the catalog image:

   `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.50`

3. Create a `CatalogSource` named `cywell-opsmate-catalog` in `openshift-marketplace`.

4. Apply:

   ```powershell
   oc apply -f deploy/olm/install/catalogsource.yaml
   oc apply -f deploy/olm/install/namespace.yaml
   oc apply -f deploy/olm/install/operatorgroup.yaml
   oc apply -f deploy/olm/install/subscription.yaml
   ```

The direct `config/default` apply path remains a development bootstrap path only.

For CRC/local validation, prefer `deploy/olm/local-crc`. It installs the
Subscription into `cywell-opsmate-olm` so the OLM path can be tested without
removing or mutating the direct-bootstrap namespace.
