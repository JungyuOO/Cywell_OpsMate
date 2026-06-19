# v0.0.52 Planner

## 1. Goal

- Fix CYOps chat input so Enter sends the message and Shift+Enter keeps a newline.
- Avoid OpenShift Console plugin proxy POST failures for chat by adding a GET-compatible chat endpoint used by the browser plugin path.
- Add a visible `자료` entry near the OpenShift admin navigation so customers can open document management from the left menu.

## 2. Architecture Overview

- Backend keeps `POST /api/chat` for service-to-service and direct gateway clients.
- Backend also supports `GET /api/chat?message=...&provider=lightspeed&rag=true` for the Console plugin proxy path where POST can be rejected by the console backend.
- The Console plugin keeps the floating CYOps launcher and injects a `자료` button near the OpenShift left navigation management area.
- Mutating document requests now use a best-effort real Console CSRF token from `window.SERVER_FLAGS`, meta tags, or cookies.
- Document metadata and upload APIs remain backed by the existing CYOps appserver document repository.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console UI | OpenShift ConsolePlugin JS | Keyboard submit and nav injection |
| Appserver | Go HTTP API | GET-compatible chat fallback |
| Packaging | OLM catalog | `v0.0.52 -> v0.0.51 -> v0.0.50` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Chat input and proxy-safe chat call | done | appserver tests |
| Phase 2 | Left navigation document entry | done | plugin manifest/entry tests |
| Phase 3 | Packaging and CRC smoke | done | CRC v0.0.52 CSV and backend smoke |

Tracking issue: #214

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.51` to `cywell-opsmate.v0.0.52`.
- No CRD or data migration.
- Existing `OpsMateConfig` Lightspeed settings remain unchanged.

## 6. Message, Communication, And Data Protocol

- UI chat through Console plugin proxy uses GET query parameters.
- Direct backend clients can continue to use JSON POST.
- Both paths call the same provider handling and still route through Lightspeed.

## 7. Security Considerations

- GET chat is intended for Console plugin proxy compatibility and does not expose provider secrets.
- The endpoint still sends only the user message to Lightspeed; model/provider selection remains in OLSConfig.
- Future production hardening should use the official Console plugin fetch/auth helper once the plugin is moved from injected JS to a bundled frontend.

## 8. Completion Criteria

- [x] Enter sends chat, Shift+Enter remains available for multiline input.
- [x] Console proxy chat path avoids POST-only 403.
- [x] Left navigation exposes `자료` document management entry.
- [x] Go tests pass: `go test ./...`.
- [x] OLM dry-run passes: bundle manifests and CatalogSource.
- [x] v0.0.52 installs on CRC.
- [x] CRC backend smoke confirms GET and POST chat return Lightspeed responses.
- [ ] Browser smoke confirms chat no longer returns 403.
