package appserver

import "net/http"

const consolePluginManifestJSON = `{
  "name": "cyops-console",
  "version": "0.0.37",
  "displayName": "CYOps",
  "description": "CYOps OpenShift operational assistant plugin.",
  "dependencies": {
    "@console/pluginAPI": "*"
  },
  "extensions": [
    {
      "type": "console.navigation/href",
      "properties": {
        "id": "cyops-diagnostics",
        "name": "CYOps Diagnostics",
        "href": "/console-plugin/diagnostics",
        "section": "cyops",
        "perspective": "admin"
      }
    }
  ]
}`

const consolePluginEntryJS = `window.__CYOPS_CONSOLE_PLUGIN__ = {
  name: "cyops-console",
  version: "0.0.37",
  diagnosticsPath: "/console-plugin/diagnostics"
};`

const consoleDiagnosticsHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>CYOps Diagnostics</title>
  <link rel="stylesheet" href="/console-plugin/diagnostics.css">
</head>
<body>
  <main class="cyops-diagnostics" data-cyops-view="diagnostics">
    <header class="cyops-header">
      <div>
        <p class="cyops-kicker">CYOps</p>
        <h1>Diagnostics</h1>
      </div>
      <button id="refresh" class="cyops-button" type="button">Refresh</button>
    </header>
    <section id="status" class="cyops-status" aria-live="polite">Loading diagnostics...</section>
    <section class="cyops-grid" aria-label="CYOps diagnostics">
      <article class="cyops-panel">
        <h2>Retrieval</h2>
        <dl id="retrieval"></dl>
      </article>
      <article class="cyops-panel">
        <h2>Documents</h2>
        <dl id="documents"></dl>
      </article>
      <article class="cyops-panel">
        <h2>Re-embedding</h2>
        <dl id="reembedding"></dl>
      </article>
      <article class="cyops-panel">
        <h2>Admin Context</h2>
        <dl id="admin"></dl>
      </article>
    </section>
  </main>
  <script type="module" src="/console-plugin/diagnostics.js"></script>
</body>
</html>`

const consoleDiagnosticsCSS = `:root {
  color-scheme: light dark;
  font-family: "Red Hat Text", system-ui, sans-serif;
  background: Canvas;
  color: CanvasText;
}

body {
  margin: 0;
}

.cyops-diagnostics {
  box-sizing: border-box;
  min-height: 100vh;
  padding: 24px;
}

.cyops-header {
  align-items: center;
  border-bottom: 1px solid color-mix(in srgb, CanvasText 20%, transparent);
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 20px;
  padding-bottom: 16px;
}

.cyops-kicker {
  font-size: 12px;
  font-weight: 700;
  margin: 0 0 4px;
  text-transform: uppercase;
}

h1, h2 {
  margin: 0;
}

h1 {
  font-size: 24px;
}

h2 {
  font-size: 16px;
  margin-bottom: 12px;
}

.cyops-button {
  border: 1px solid color-mix(in srgb, CanvasText 30%, transparent);
  border-radius: 4px;
  background: ButtonFace;
  color: ButtonText;
  cursor: pointer;
  font: inherit;
  min-height: 36px;
  padding: 0 14px;
}

.cyops-status {
  margin-bottom: 16px;
  min-height: 20px;
}

.cyops-grid {
  display: grid;
  gap: 16px;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
}

.cyops-panel {
  border: 1px solid color-mix(in srgb, CanvasText 20%, transparent);
  border-radius: 6px;
  padding: 16px;
}

dl {
  display: grid;
  gap: 8px 12px;
  grid-template-columns: minmax(120px, auto) 1fr;
  margin: 0;
}

dt {
  color: color-mix(in srgb, CanvasText 70%, transparent);
}

dd {
  font-weight: 600;
  margin: 0;
  overflow-wrap: anywhere;
}`

const consoleDiagnosticsJS = `const statusEl = document.querySelector("#status");
const refreshButton = document.querySelector("#refresh");

function setStatus(message) {
  statusEl.textContent = message;
}

function renderDefinitionList(id, entries) {
  const target = document.querySelector(id);
  target.replaceChildren();
  for (const [label, value] of entries) {
    const term = document.createElement("dt");
    term.textContent = label;
    const description = document.createElement("dd");
    description.textContent = String(value ?? "-");
    target.append(term, description);
  }
}

async function fetchJSON(path) {
  const response = await fetch(path, { credentials: "same-origin" });
  if (!response.ok) {
    throw new Error(path + " returned " + response.status);
  }
  return response.json();
}

async function loadDiagnostics() {
  setStatus("Loading diagnostics...");
  const [schema, diagnostics] = await Promise.all([
    fetchJSON("/api/ops/diagnostics/schema"),
    fetchJSON("/api/ops/diagnostics"),
  ]);

  renderDefinitionList("#retrieval", [
    ["Total", diagnostics.retrieval?.total],
    ["Slow", diagnostics.retrieval?.slow],
    ["Failures", diagnostics.retrieval?.failures],
    ["Last mode", diagnostics.retrieval?.last?.mode || "-"],
    ["Last results", diagnostics.retrieval?.last?.resultCount || 0],
  ]);
  renderDefinitionList("#documents", [
    ["Total", diagnostics.documents?.total],
    ["By status", JSON.stringify(diagnostics.documents?.byStatus || {})],
    ["By embedding", JSON.stringify(diagnostics.documents?.byEmbeddingStatus || {})],
  ]);
  renderDefinitionList("#reembedding", [
    ["Available", diagnostics.reembedding?.available ? "Yes" : "No"],
    ["Contract", schema.aggregateOnly ? "Aggregate only" : "Expanded"],
  ]);
  renderDefinitionList("#admin", [
    ["User", diagnostics.admin?.authorizedUser || "-"],
    ["Groups", (diagnostics.admin?.authorizedGroups || []).join(", ") || "-"],
    ["Primary entry", diagnostics.ui?.primaryEntry || schema.primaryEntry],
    ["Fallback route", diagnostics.ui?.fallbackRoute || "optional"],
  ]);
  setStatus("Diagnostics loaded from the OpenShift Web Console backend path.");
}

refreshButton.addEventListener("click", () => {
  loadDiagnostics().catch((error) => setStatus(error.message));
});

loadDiagnostics().catch((error) => setStatus(error.message));`

func (s *Server) consolePluginManifest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(consolePluginManifestJSON))
}

func (s *Server) consolePluginEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(consolePluginEntryJS))
}

func (s *Server) consoleDiagnostics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(consoleDiagnosticsHTML))
}

func (s *Server) consoleDiagnosticsJS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(consoleDiagnosticsJS))
}

func (s *Server) consoleDiagnosticsCSS(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(consoleDiagnosticsCSS))
}
