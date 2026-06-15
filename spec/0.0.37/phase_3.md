# v0.0.37 Phase 3 - Plugin Manifest Backend Smoke

## Work

- [x] Verified appserver serves `/plugin-manifest.json`.
- [x] Verified appserver serves `/console-plugin/plugin-entry.js`.
- [x] Verified appserver still serves `/console-plugin/diagnostics`.
- [x] Attempted Console plugin proxy smoke.

## Verification

- `oc get deploy cyops-appserver -n cywell-opsmate-olm -o jsonpath='{.spec.template.spec.containers[0].image}'` returned `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.37`.
- `oc get pod -n cywell-opsmate-olm -l app.kubernetes.io/component=appserver -o wide` showed `cyops-appserver-7b96fff8b9-6hznh 1/1 Running`.
- `oc get opsmateconfig cyops -n cywell-opsmate-olm -o jsonpath='{.status.overallStatus}'` returned `Ready`.
- `curl.exe -k -s https://127.0.0.1:18443/plugin-manifest.json` through `oc port-forward` returned `name: cyops-console`, `version: 0.0.37`, `displayName: CYOps`, and `console.navigation/href`.
- `curl.exe -k -s https://127.0.0.1:18443/console-plugin/plugin-entry.js` returned `window.__CYOPS_CONSOLE_PLUGIN__` with version `0.0.37`.
- `curl.exe -k -s https://127.0.0.1:18443/console-plugin/diagnostics` returned `CYOps Diagnostics`.
- `curl.exe -k` against the Console proxy endpoint returned `401` with bearer token only, so browser-session-cookie verification remains a UI smoke follow-up.
- `oc logs deploy/console -n openshift-console --tail=80` showed `cyops-console` in the enabled plugin order, but ConsolePlugin metrics still reported `unknown`.

## Remaining Scope

- Implement a real OpenShift dynamic plugin bundle/extension so the Console identifies CYOps as CYOps rather than `unknown` and can expose the intended launcher/chat UI.

## Linked Issue

- #182
