# v0.0.37 Planner

## 1. Goal

- Fix the v0.0.36 catalog graph break found during CRC upgrade smoke.
- Publish and install the v0.0.37 catalog on CRC.
- Verify the OLM upgrade to `cywell-opsmate.v0.0.37`.
- Verify the appserver plugin manifest endpoints through the live pod and, where possible, through the OpenShift Console plugin proxy.
- Record the remaining scope for a full CYOps dynamic plugin launcher/chat drawer.

## 2. Architecture Overview

- v0.0.36 added plugin manifest and entry endpoints to the appserver image, but its catalog head replaced `v0.0.33` while CRC was already installed at `v0.0.34`.
- v0.0.37 keeps the plugin endpoints and corrects the channel head to replace `v0.0.34`.
- CRC already has `ConsolePlugin/cyops-console` and should keep `cyops-console` enabled in Console operator `spec.plugins` for this smoke.
- OLM upgrades the manager, and the manager reconciles the appserver deployment to the v0.0.36 appserver image.
- Direct pod/service checks prove backend readiness even if browser authentication or CRC certificate behavior blocks full UI smoke.

## 3. Technical Stack

- OpenShift CLI.
- CRC OpenShift cluster with OLM.
- GHCR v0.0.36 manager, appserver, bundle, and catalog images.
- ConsolePlugin backend service.
- In-app browser smoke when the CRC console session can be reached safely.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Confirm v0.0.36 artifacts and CRC readiness | done | workflow/image/cluster evidence |
| Phase 2 | Fix and publish v0.0.37 graph | done | graph fix and validation evidence |
| Phase 3 | Verify plugin manifest backend endpoints | planned | appserver and console proxy evidence |
| Phase 4 | Handoff | planned | next-version scope |

## 5. Migration Or Operation Strategy

- Prefer in-place OLM upgrade from the existing CRC Subscription.
- Preserve existing Console operator plugins and keep `cyops-console` enabled only for CRC smoke.
- Do not hide or delete Red Hat Lightspeed UI.
- If authentication or certificate behavior blocks browser smoke, capture the blocker and verify backend endpoints directly.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| CatalogSource | `openshift-marketplace/cywell-opsmate-catalog` |
| Upgrade target | `cywell-opsmate.v0.0.37` |
| Manager image | `ghcr.io/jungyuoo/cywell-opsmate:v0.0.37` |
| Appserver image | `ghcr.io/jungyuoo/cywell-opsmate-appserver:v0.0.37` |
| Plugin manifest | `/plugin-manifest.json` |
| Plugin entry | `/console-plugin/plugin-entry.js` |

## 7. Security Considerations

- Do not print tokens or kubeconfig contents.
- Do not bypass browser security interstitials without explicit user action.
- Do not broaden SCC or mutate unrelated Console capabilities for smoke.

## 8. Completion Criteria

- [x] v0.0.36 artifacts are confirmed.
- [x] v0.0.37 graph fix is prepared and locally validated.
- [ ] CRC OLM upgrade reaches `cywell-opsmate.v0.0.37` `Succeeded`.
- [ ] Manager and appserver deployments run v0.0.37 images.
- [ ] Plugin manifest and entry endpoints return expected CYOps content.
- [ ] Browser/proxy smoke result or blocker is documented.
- [ ] Remaining launcher/chat drawer scope is assigned to the next version.
