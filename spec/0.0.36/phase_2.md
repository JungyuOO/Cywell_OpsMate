# v0.0.36 Phase 2 - CRC Browser Smoke

## Work

- [x] Enabled `cyops-console` in CRC Console operator plugin list for smoke.
- [x] Confirmed console rollout completed and logs listed `cyops-console` as enabled.
- [x] Attempted in-app browser smoke against the CRC console route.
- [x] Attempted token-backed console plugin endpoint smoke.

## Verification

- `oc patch console.operator.openshift.io cluster --type merge --patch-file .tmp-console-plugins.json` returned `console.operator.openshift.io/cluster patched`.
- `oc get console.operator.openshift.io cluster -o jsonpath='{.spec.plugins}'` returned `["networking-console-plugin","monitoring-plugin","cyops-console"]`.
- `oc rollout status deployment/console -n openshift-console --timeout=180s` returned `deployment "console" successfully rolled out`.
- `oc logs deployment/console -n openshift-console --tail=160` listed `cyops-console` in the enabled plugin order.
- Console metrics showed `unknown: enabled`, which is consistent with CYOps being enabled but not yet fully identified by a proper dynamic plugin surface.
- Browser navigation to `https://console-openshift-console.apps-crc.testing` was blocked by `ERR_CERT_AUTHORITY_INVALID`; the browser safety interstitial was not bypassed.
- `curl.exe -k` calls to `/api/plugins/cyops-console/...` returned `401` with bearer token only, so browser-session-cookie verification remains pending.

## Remaining Scope

- CRC API later became unreachable at `127.0.0.1:6443`, so the `cyops-console` addition to Console operator `spec.plugins` could not be reverted or re-verified in this phase.
- v0.0.37 should re-run the browser smoke after the v0.0.36 appserver image is published and installed.

## Linked Issue

- #176
