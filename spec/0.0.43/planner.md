# v0.0.43 Planner

## 1. Goal

- Fix CYOps ConsolePlugin callback registration so OpenShift Console executes CYOps codeRefs.
- Add a supported `console.flag` extension that forces the Console runtime to resolve the CYOps callback module map.
- Keep the v0.0.42 nginx gateway topology.

## 2. Architecture Overview

- `cyops-gateway` remains the ConsolePlugin backend Service.
- `cyops-appserver` still serves `/plugin-manifest.json` and `/plugin-entry.js`.
- Manifest includes a `console.flag` extension with `$codeRef: cyopsLauncherFlag`.
- `plugin-entry.js` registers `cyops-console@0.0.43` with a callback module map matching working OpenShift plugins.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | OpenShift callback dynamic plugin | Uses `loadPluginEntry("name@version", moduleMap)` |
| Launcher trigger | `console.flag` codeRef | Executes early and mounts CYOps launcher |
| Gateway | nginx unprivileged container | Existing v0.0.42 topology |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Callback module map fix | done | appserver tests |
| Phase 2 | OLM packaging and CRC upgrade | done | CRC CSV/gateway endpoint smoke |
| Phase 3 | Issue/PR handoff | done | PR #197 merged and issue #196 closed |

Tracking issue: #196

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.42` to `cywell-opsmate.v0.0.43`.
- No Service topology changes.
- Existing `cyops-gateway` and `cyops-appserver` continue to run.

## 6. Message, Communication, And Data Protocol

- No backend API payload changes.
- Console loads `/plugin-manifest.json`, sees `console.flag`, then resolves `cyopsLauncherFlag` from `/plugin-entry.js`.

## 7. Security Considerations

- No OAuth, credential, or RBAC changes.
- The flag handler returns only a local boolean flag and mounts local CYOps UI.

## 8. Completion Criteria

- [x] Manifest contains `console.flag` with `$codeRef: cyopsLauncherFlag`.
- [x] Entry uses `loadPluginEntry("cyops-console@0.0.43", moduleMap)`.
- [x] Go tests pass.
- [x] v0.0.43 installs on CRC.
- [x] Gateway endpoint smoke confirms callback module map content.
