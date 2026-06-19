# v0.0.44 Planner

## 1. Goal

- Make the CYOps launcher mount as soon as OpenShift Console loads `/plugin-entry.js`.
- Keep the supported callback registration, but stop depending on `console.flag` handler invocation for visible UI.
- Add a browser-visible execution marker so the loaded entry can be diagnosed from the DOM.

## 2. Architecture Overview

- `cyops-gateway` remains the ConsolePlugin backend Service.
- `cyops-appserver` continues to serve `/plugin-manifest.json` and `/plugin-entry.js`.
- Manifest still includes `console.flag` with `$codeRef: cyopsLauncherFlag`.
- `plugin-entry.js` registers `cyops-console@0.0.44`, marks `document.documentElement`, and mounts CYOps immediately after registration.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | OpenShift callback dynamic plugin | Keeps callback registration |
| Launcher mount | Browser DOM side effect | Runs when entry script loads |
| Gateway | nginx unprivileged container | Existing v0.0.42 topology |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Entry-load mount fix | done | appserver tests |
| Phase 2 | OLM packaging and CRC upgrade | in progress | Go tests and OLM dry-run passed |
| Phase 3 | Issue/PR handoff | planned | pending |

Tracking issue: #198

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.43` to `cywell-opsmate.v0.0.44`.
- Preserve `cyops-gateway`, `cyops-appserver`, and PostgreSQL topology.
- Keep catalog graph as `v0.0.44 -> v0.0.43 -> v0.0.42`.

## 6. Message, Communication, And Data Protocol

- No backend API payload changes.
- Console still loads `/plugin-manifest.json` and `/plugin-entry.js`.
- Entry load now sets `data-cyops-plugin-entry="0.0.44"` and calls the CYOps launcher mount path directly.

## 7. Security Considerations

- No OAuth, credential, or RBAC changes.
- The new marker contains only the plugin version and no user or cluster data.
- The launcher continues to call backend APIs with same-origin Console proxy credentials.

## 8. Completion Criteria

- [x] Entry registers `cyops-console@0.0.44`.
- [x] Entry mounts CYOps immediately after load.
- [x] Entry leaves `data-cyops-plugin-entry="0.0.44"` for diagnosis.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [ ] v0.0.44 installs on CRC.
- [ ] Gateway endpoint smoke confirms the direct mount content.
