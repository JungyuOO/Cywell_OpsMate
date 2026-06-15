# CYOps Console Plugin Bundle

This directory holds the frontend source contract for the CYOps OpenShift
ConsolePlugin.

v0.0.24 introduces the first diagnostics view contract without adding a
JavaScript toolchain. The appserver serves the runnable view through:

- `/console-plugin/diagnostics`
- `/console-plugin/diagnostics.js`
- `/console-plugin/diagnostics.css`

The JavaScript source in `src/` documents the dynamic plugin entry shape and
calls only the Web Console backend paths:

- `/api/ops/diagnostics`
- `/api/ops/diagnostics/schema`

The normal user journey is OpenShift Web Console -> CYOps ConsolePlugin ->
appserver backend. The fallback admin Route is not part of this path.

The bundle keeps the Red Hat OpenShift Lightspeed UI separate. CYOps owns its
own branding, diagnostics surface, chat drawer, and document-management
workflow.
