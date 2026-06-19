# CYOps Console Plugin Bundle

This directory holds the frontend source contract for the CYOps OpenShift
ConsolePlugin.

v0.0.44 keeps the no-toolchain approach but now serves a standard OpenShift
dynamic plugin manifest and callback entry bundle. The appserver serves:

- `/plugin-manifest.json`
- `/plugin-entry.js`
- `/console-plugin/diagnostics`
- `/console-plugin/diagnostics.js`
- `/console-plugin/diagnostics.css`

The manifest includes `baseURL`, `loadScripts`, and
`registrationMethod: "callback"`. It also exposes a `console.flag` extension
with `$codeRef: "cyopsLauncherFlag"` so the Console runtime resolves the CYOps
callback module. The entry script calls
`window.loadPluginEntry("cyops-console@0.0.44", ...)`, injects the CYOps
launcher, and uses only the Web Console plugin backend paths:

- `/api/chat`
- `/api/documents`
- `/api/ops/diagnostics`
- `/api/ops/diagnostics/schema`

The normal user journey is OpenShift Web Console -> CYOps ConsolePlugin ->
appserver backend. The fallback admin Route is not part of this path.

The bundle keeps the Red Hat OpenShift Lightspeed UI separate. CYOps owns its
own branding, diagnostics surface, chat drawer, and document-management
workflow.
