# CYOps OLM Install

This directory describes the intended customer-facing install path.

1. Publish the bundle image:

   `ghcr.io/jungyuoo/cywell-opsmate-bundle:v0.0.29`

2. Publish the catalog image:

   `ghcr.io/jungyuoo/cywell-opsmate-catalog:v0.0.30`

3. Create a `CatalogSource` named `cywell-opsmate-catalog` in `openshift-marketplace`.

4. Apply:

   ```powershell
   oc apply -f deploy/olm/install/catalogsource.yaml
   oc apply -f deploy/olm/install/namespace.yaml
   oc apply -f deploy/olm/install/operatorgroup.yaml
   oc apply -f deploy/olm/install/subscription.yaml
   ```

The direct `config/default` apply path remains a development bootstrap path only.
