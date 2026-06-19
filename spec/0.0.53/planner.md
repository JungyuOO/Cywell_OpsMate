# v0.0.53 Planner

## 1. Goal

- Remove document management from the CYOps chat drawer.
- Show document upload and management in a central OpenShift Console-style workspace when the left Documents entry is selected.
- Replace query-string chat fallback with path-based GET chat so Console proxy forwarding does not drop the message.

## 2. Architecture Overview

- Chat drawer remains focused on chat only.
- The Documents navigation entry opens a central CYOps document workspace overlay aligned with the OpenShift Console content area.
- Console plugin chat uses `GET /api/chat/message/<encoded-message>` when loaded through the OpenShift plugin proxy.
- Direct clients keep using `POST /api/chat`.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console UI | OpenShift ConsolePlugin JS | Chat drawer plus central document workspace |
| Appserver | Go HTTP API | Path-based GET chat fallback |
| Packaging | OLM catalog | `v0.0.53 -> v0.0.52 -> v0.0.51` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Chat drawer cleanup and path chat fallback | done | appserver tests |
| Phase 2 | Central document workspace | done | plugin entry tests |
| Phase 3 | Packaging and CRC smoke | done | PR #217, CSV `cywell-opsmate.v0.0.53`, path chat smoke |

Tracking issue: #216 (closed by PR #217)

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.52` to `cywell-opsmate.v0.0.53`.
- No CRD or data migration.
- Existing Lightspeed and document repository configuration remains unchanged.

## 6. Message, Communication, And Data Protocol

- Console plugin chat uses an encoded path segment for the message.
- Backend decodes the path segment and routes to the same provider handling.
- Document management continues to use `/api/documents`.

## 7. Security Considerations

- Path-based GET chat avoids proxy query forwarding issues but should still not carry secrets or large prompt payloads.
- Document upload still uses CSRF-aware POST handling through the console plugin.
- Future bundled frontend should replace DOM injection with official Console page/navigation APIs.

## 8. Completion Criteria

- [x] Chat drawer no longer contains document upload/list UI.
- [x] Documents navigation opens a central document workspace.
- [x] Console proxy chat uses path-based GET fallback.
- [x] Go tests pass: `go test ./...`.
- [x] OLM dry-run passes: bundle manifests and CatalogSource.
- [x] v0.0.53 installs on CRC: CSV `cywell-opsmate.v0.0.53` is `Succeeded`.
- [x] CRC backend smoke confirms path-based chat returns Lightspeed response through `/api/chat/message/<encoded-message>`.
