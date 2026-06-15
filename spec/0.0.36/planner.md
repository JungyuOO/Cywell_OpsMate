# v0.0.36 Planner

## 1. Goal

- Explain and verify why CYOps does not appear in the OpenShift Web Console floating launcher after the Operator is running.
- Run a CRC browser smoke for the ConsolePlugin activation path.
- Tighten packaging verification so manager, bundle, catalog, and ConsolePlugin-serving assumptions are checked together.
- Keep the PostgreSQL restricted-SCC chmod warning documented without changing database image behavior unless it blocks readiness.

## 2. Architecture Overview

- `ConsolePlugin/cyops-console` registers CYOps with OpenShift, but the OpenShift Console only loads plugins listed in `console.operator.openshift.io/cluster.spec.plugins`.
- The appserver currently serves diagnostics pages and APIs; a full OpenShift dynamic plugin manifest/asset bundle is a separate frontend packaging concern.
- CRC smoke can enable `cyops-console` in the Console operator to test whether the Console attempts to load CYOps.
- OLM packaging validation remains local and GitHub Actions backed.

## 3. Technical Stack

- OpenShift CLI.
- CRC OpenShift Web Console.
- OpenShift Console operator `spec.plugins`.
- `ConsolePlugin` backend service.
- Go unit tests and kustomize/OLM validation.
- Browser smoke through the in-app browser.

## 4. Implementation Steps

| Phase | Scope | Status | Output |
| --- | --- | --- | --- |
| Phase 1 | Root-cause CYOps launcher absence | done | ConsolePlugin/Console operator evidence |
| Phase 2 | CRC ConsolePlugin browser smoke | done | browser and console evidence |
| Phase 3 | Packaging validation cleanup | done | tests and packaging changes |
| Phase 4 | Handoff | done | v0.0.37 scope |

## 5. Migration Or Operation Strategy

- Do not remove or hide Red Hat OpenShift Lightspeed resources.
- For CRC smoke only, patch `console.operator.openshift.io/cluster.spec.plugins` to include `cyops-console`.
- Preserve the pre-existing Console operator plugins and document the temporary CRC mutation.
- Treat PostgreSQL `/var/run/postgresql` chmod output as a warning while the pod remains Ready.

## 6. Message, Communication, And Data Protocol

| Surface | Contract |
| --- | --- |
| ConsolePlugin | `cyops-console` |
| Console operator | `console.operator.openshift.io/cluster.spec.plugins` |
| Appserver service | `cyops-appserver.cywell-opsmate-olm.svc:8443` |
| Diagnostics view | `/console-plugin/diagnostics` |
| Dynamic plugin assets | To be verified in this version |

## 7. Security Considerations

- CRC console patches are local validation steps only.
- Do not commit kubeconfig, tokens, screenshots with secrets, or generated TLS materials.
- Do not broaden SCC permissions to silence non-fatal PostgreSQL warnings.
- Do not disable Red Hat Lightspeed UI from the CYOps operator.

## 8. Completion Criteria

- [x] The absence of CYOps in the launcher is explained with cluster evidence.
- [x] CRC console plugin activation smoke is run and documented.
- [x] Packaging validation commands are run and documented.
- [x] Any required product change is implemented with tests, or the next version scope is explicitly recorded.
