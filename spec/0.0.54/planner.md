# v0.0.54 Planner

## 1. Goal

- Make the Documents surface feel like an OpenShift Console page instead of a floating overlay.
- Keep the CYOps chat launcher separate from document management.
- Preserve the v0.0.53 path-based chat fallback.

## 2. Architecture Overview

- The injected left navigation entry closes the chat drawer and navigates to `/console-plugin/documents`.
- `/console-plugin/documents` owns document upload/list management as a normal page surface.
- The document page uses a Console-like header, toolbar, and table layout.
- The plugin entry no longer mounts or owns document upload DOM.

## 3. Technology Stack

| Area | Tooling | Notes |
| --- | --- | --- |
| Console UI | OpenShift ConsolePlugin JS | Launcher plus navigation handoff |
| Documents page | Appserver-served HTML/CSS/JS | Page-style upload/list management |
| Packaging | OLM catalog | `v0.0.54 -> v0.0.53 -> v0.0.52` |

## 4. Implementation Steps

| Phase | Scope | Status | Evidence |
| --- | --- | --- | --- |
| Phase 1 | Remove injected document overlay and route navigation to the documents page | done | plugin entry tests |
| Phase 2 | Restyle documents page as a Console-style page/table | done | documents asset tests |
| Phase 3 | Package and CRC smoke | done | PR #219, CSV `cywell-opsmate.v0.0.54`, plugin asset smoke |

Tracking issue: #218

## 5. Migration Or Operations Strategy

- Upgrade from `cywell-opsmate.v0.0.53` to `cywell-opsmate.v0.0.54`.
- No database, CRD, or document storage migration.
- Browser users may need a hard refresh because OpenShift Console can cache plugin assets.

## 6. Message, Communication, And Data Protocol

- Chat remains on `/api/chat` for direct appserver calls and `/api/chat/message/<encoded-message>` through the plugin proxy fallback.
- Documents page continues to use `/api/documents` for list and upload.
- The left `자료` entry uses page navigation instead of in-place DOM overlay mutation.

## 7. Security Considerations

- Document upload continues to use the existing authenticated Console route and appserver document API.
- The page removes overlay z-index conflicts that could hide underlying Console controls.
- Future versions should move the navigation/page implementation to official Console extension surfaces when the richer page API is introduced.

## 8. Completion Criteria

- [x] Chat drawer does not contain document management.
- [x] Injected `자료` action navigates to `/console-plugin/documents`.
- [x] No injected `cyops-doc-workspace` overlay remains in plugin entry.
- [x] Documents page renders as a Console-like page with header, toolbar, and table rows.
- [x] Go tests pass.
- [x] OLM dry-run passes.
- [x] v0.0.54 is pushed through PR and CRC smoke is recorded.
