# v0.0.47 Planner

## 1. Goal

- Reduce remaining OpenShift Console plugin proxy 403 risk for CYOps drawer POST requests.
- Accept common internal LLM and OpenAI-compatible response shapes from the configured Lightspeed provider endpoint.
- Keep the v0.0.46 Lightspeed runtime wiring intact.

## 2. Architecture Overview

- CYOps drawer non-GET requests include CSRF and XMLHttpRequest headers.
- The appserver `LightspeedProvider` still sends one JSON chat payload to the configured endpoint.
- Provider response parsing accepts `answer`, `response`, `output`, `content`, `text`, `generated_text`, OpenAI-style `choices[].message.content`, and similar nested structures.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console plugin | Browser fetch through OpenShift plugin proxy | Adds defensive non-secret headers |
| Provider adapter | Go JSON traversal | Avoids hard-coding one LLM response schema |
| Packaging | OLM catalog | `v0.0.47 -> v0.0.46 -> v0.0.45` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Header and response adapter | done | appserver tests |
| Phase 2 | OLM packaging and CRC upgrade | done | CRC CSV `cywell-opsmate.v0.0.47` Succeeded |
| Phase 3 | Issue/PR handoff | done | PR #205 merged, issue #204 closed |

Tracking issue: #204

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.46` to `cywell-opsmate.v0.0.47`.
- No CRD or storage migration.
- Existing `spec.lightspeed.apiBaseURL` values continue to work.

## 6. Message, Communication, And Data Protocol

- Request payload is unchanged.
- Response parsing is expanded to common chat completion shapes.
- If no answer text is found, appserver returns a generic provider failure to the drawer.

## 7. Security Considerations

- `X-CSRFToken`, `X-CSRF-Token`, and `X-Requested-With` are static non-secret headers.
- Provider bearer token behavior from v0.0.46 is unchanged.
- Raw provider payloads are still not returned to users.

## 8. Completion Criteria

- [x] Non-GET drawer requests include CSRF and XMLHttpRequest headers.
- [x] Provider parses OpenAI-style and internal LLM response fields.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [x] v0.0.47 installs on CRC.
- [ ] Browser chat smoke no longer gets 403.

Browser chat smoke remains user-observed because the CRC console session must issue the proxied request from the authenticated browser context.
