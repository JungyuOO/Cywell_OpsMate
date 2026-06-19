# v0.0.45 Planner

## 1. Goal

- Fix CYOps chat and document API calls from the OpenShift Console plugin drawer.
- Preserve the visible v0.0.44 launcher behavior.
- Route browser requests through the Console plugin backend proxy instead of the Console root path.

## 2. Architecture Overview

- `cyops-gateway` remains the ConsolePlugin backend Service.
- Console loads CYOps entry from `/api/plugins/cyops-console/plugin-entry.js`.
- The entry captures that plugin proxy base during script execution and later uses `/api/plugins/cyops-console/api/...` for chat and document requests.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | OpenShift callback dynamic plugin | Existing launcher stays visible |
| API routing | Console plugin backend proxy | Avoids Console root `/api/chat` 404 |
| Gateway | nginx unprivileged container | Existing topology |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Proxy API path fix | done | appserver tests |
| Phase 2 | OLM packaging and CRC upgrade | done | CRC CSV and gateway smoke passed |
| Phase 3 | Issue/PR handoff | done | PR #201 merged and issue #200 closed |

Tracking issue: #200

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.44` to `cywell-opsmate.v0.0.45`.
- Preserve gateway, appserver, PostgreSQL, and ConsolePlugin names.
- Keep catalog graph as `v0.0.45 -> v0.0.44 -> v0.0.43`.

## 6. Message, Communication, And Data Protocol

- Chat and document payloads do not change.
- Browser path changes from `/api/chat` to `/api/plugins/cyops-console/api/chat` when the entry is loaded through OpenShift Console.
- Local diagnostics and direct appserver routes still use `/api/...`.

## 7. Security Considerations

- Same-origin credentials are still used.
- No OAuth, RBAC, token, or secret handling changes.

## 8. Completion Criteria

- [x] Entry captures plugin proxy base at script load time.
- [x] Event handlers reuse that base for chat and document requests.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [x] v0.0.45 installs on CRC.
- [x] Gateway smoke confirms proxy-base entry content.
- [ ] Browser smoke confirms chat no longer calls Console root `/api/chat`.
