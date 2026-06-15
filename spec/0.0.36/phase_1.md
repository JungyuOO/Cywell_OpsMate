# v0.0.36 Phase 1 - CYOps Launcher Absence Root Cause

## Work

- [x] Inspected `ConsolePlugin/cyops-console`.
- [x] Inspected `console.operator.openshift.io/cluster.spec.plugins`.
- [x] Identified that the plugin was registered but not enabled before the CRC smoke patch.
- [x] Identified that the appserver served diagnostics pages but did not expose dynamic plugin manifest/entry endpoints.

## Verification

- `oc get consoleplugin cyops-console -o yaml` showed `spec.displayName: CYOps` and backend service `cyops-appserver` in namespace `cywell-opsmate-olm`.
- `oc get console.operator.openshift.io cluster -o yaml` showed `spec.plugins` contained only `networking-console-plugin` and `monitoring-plugin` before the smoke patch.
- `oc get consoleplugins.console.openshift.io` showed `cyops-console`, `monitoring-plugin`, and `networking-console-plugin`.
- Code inspection showed appserver routes for `/console-plugin/diagnostics`, `/console-plugin/diagnostics.js`, and `/console-plugin/diagnostics.css`, but no `/plugin-manifest.json` or plugin entry route before this version.

## Remaining Scope

- OpenShift Console can only load CYOps after the cluster Console operator includes `cyops-console` in `spec.plugins`.
- CYOps does not replace the Red Hat Lightspeed floating button by registration alone; that button is controlled by the console `LightspeedButton` capability and Red Hat Lightspeed UI.

## Linked Issue

- #178
