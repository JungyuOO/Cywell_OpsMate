package appserver

import "net/http"

const consolePluginManifestJSON = `{
  "name": "cyops-console",
  "version": "0.0.48",
  "baseURL": "/api/plugins/cyops-console/",
  "loadScripts": [
    "plugin-entry.js"
  ],
  "registrationMethod": "callback",
  "dependencies": {
    "@console/pluginAPI": "*"
  },
  "customProperties": {
    "console": {
      "displayName": "CYOps",
      "description": "CYOps OpenShift operational assistant plugin."
    }
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
    },
    {
      "type": "console.flag",
      "properties": {
        "handler": {
          "$codeRef": "cyopsLauncherFlag"
        }
      }
    }
  ]
}`

const consolePluginEntryJS = `window.__CYOPS_CONSOLE_PLUGIN__ = {
  name: "cyops-console",
  version: "0.0.48",
  diagnosticsPath: "/console-plugin/diagnostics"
};

(function () {
  const pluginName = "cyops-console";
  const launcherID = "cyops-console-launcher";
  const drawerID = "cyops-console-drawer";
  const styleID = "cyops-console-style";
  const pluginScript = document.currentScript;
  const pluginSource = pluginScript && pluginScript.src ? new URL(pluginScript.src, window.location.href) : null;
  const pluginProxyBase = "/api/plugins/" + pluginName;
  const apiBase = pluginSource && pluginSource.pathname.includes(pluginProxyBase) ? pluginProxyBase : "";

  function apiPath(path) {
    return apiBase + path;
  }

  function ensureStyle() {
    if (document.getElementById(styleID)) {
      return;
    }
    const style = document.createElement("style");
    style.id = styleID;
    style.textContent = ".cyops-launcher{position:fixed!important;right:22px!important;bottom:22px!important;z-index:2147483647!important;min-width:76px!important;height:48px!important;border-radius:8px!important;border:1px solid rgba(255,255,255,.58)!important;background:#151515!important;color:#fff!important;box-shadow:0 10px 30px rgba(0,0,0,.42)!important;font:700 14px 'Red Hat Text',system-ui,sans-serif!important;cursor:pointer!important;padding:0 16px!important;letter-spacing:0!important}.cyops-launcher:hover{background:#222!important}.cyops-drawer{position:fixed;right:22px;bottom:86px;z-index:2147483646;width:min(520px,calc(100vw - 32px));height:min(690px,calc(100vh - 120px));display:none;grid-template-rows:auto 1fr auto;background:#1f1f1f;color:#f5f5f5;border:1px solid rgba(255,255,255,.35);border-radius:8px;box-shadow:0 18px 48px rgba(0,0,0,.46);font:14px 'Red Hat Text',system-ui,sans-serif}.cyops-drawer[open]{display:grid}.cyops-drawer header{display:flex;align-items:center;justify-content:space-between;padding:16px 18px;border-bottom:1px solid rgba(255,255,255,.16)}.cyops-drawer h2{font-size:20px;margin:0;letter-spacing:0}.cyops-icon-button{width:34px;height:34px;border-radius:4px;border:1px solid rgba(255,255,255,.25);background:transparent;color:#fff;cursor:pointer}.cyops-body{display:grid;grid-template-columns:minmax(0,1fr) 172px;min-height:0}.cyops-chat{display:grid;grid-template-rows:1fr auto;min-width:0;border-right:1px solid rgba(255,255,255,.14)}.cyops-messages{overflow:auto;padding:16px;display:flex;flex-direction:column;gap:12px}.cyops-message{max-width:92%;padding:10px 12px;border-radius:6px;line-height:1.42;overflow-wrap:anywhere}.cyops-message.user{align-self:flex-end;background:#0066cc;color:#fff}.cyops-message.assistant{align-self:flex-start;background:#2d2d2d;color:#f5f5f5}.cyops-compose{display:grid;grid-template-columns:1fr 42px;gap:8px;padding:12px;border-top:1px solid rgba(255,255,255,.14)}.cyops-compose textarea{min-height:58px;max-height:120px;resize:vertical;border:1px solid #73bcf7;border-radius:6px;background:#262626;color:#fff;padding:10px;font:inherit}.cyops-send{width:42px;height:42px;border:0;border-radius:6px;background:#73bcf7;color:#111;font-weight:700;cursor:pointer}.cyops-docs{min-width:0;padding:14px;display:grid;grid-template-rows:auto auto 1fr;gap:10px}.cyops-docs h3{font-size:14px;margin:0}.cyops-file{display:block}.cyops-file input{width:100%;font-size:12px;color:#ddd}.cyops-doc-list{display:flex;flex-direction:column;gap:8px;overflow:auto}.cyops-doc-item{border:1px solid rgba(255,255,255,.14);border-radius:6px;padding:8px;background:#262626}.cyops-doc-name{font-weight:700;overflow-wrap:anywhere}.cyops-doc-meta{font-size:12px;color:#c7c7c7;margin-top:3px}.cyops-status{min-height:18px;padding:0 18px 14px;color:#c7c7c7;font-size:12px}@media (max-width:720px){.cyops-launcher{right:16px!important;bottom:16px!important}.cyops-drawer{right:10px;left:10px;width:auto}.cyops-body{grid-template-columns:1fr}.cyops-chat{border-right:0}.cyops-docs{border-top:1px solid rgba(255,255,255,.14);grid-template-rows:auto auto auto}.cyops-doc-list{max-height:112px}}";
    document.head.appendChild(style);
  }

  function createMessage(text, role) {
    const node = document.createElement("div");
    node.className = "cyops-message " + role;
    node.textContent = text;
    return node;
  }

  async function requestJSON(path, options) {
    const requestOptions = Object.assign({ credentials: "same-origin" }, options || {});
    const method = (requestOptions.method || "GET").toUpperCase();
    if (method !== "GET" && method !== "HEAD" && method !== "OPTIONS") {
      const headers = new Headers(requestOptions.headers || {});
      if (!headers.has("X-CSRFToken")) {
        headers.set("X-CSRFToken", "1");
      }
      if (!headers.has("X-CSRF-Token")) {
        headers.set("X-CSRF-Token", "1");
      }
      if (!headers.has("X-Requested-With")) {
        headers.set("X-Requested-With", "XMLHttpRequest");
      }
      requestOptions.headers = headers;
    }
    const response = await fetch(apiPath(path), requestOptions);
    if (!response.ok) {
      throw new Error(path + " returned " + response.status);
    }
    return response.json();
  }

  async function refreshDocuments(root) {
    const list = root.querySelector("[data-cyops-doc-list]");
    list.replaceChildren();
    const payload = await requestJSON("/api/documents");
    const items = payload.items || [];
    if (items.length === 0) {
      const empty = document.createElement("div");
      empty.className = "cyops-doc-item";
      empty.textContent = "No documents";
      list.appendChild(empty);
      return;
    }
    for (const item of items) {
      const row = document.createElement("div");
      row.className = "cyops-doc-item";
      const name = document.createElement("div");
      name.className = "cyops-doc-name";
      name.textContent = item.filename || item.id;
      const meta = document.createElement("div");
      meta.className = "cyops-doc-meta";
      meta.textContent = (item.status || "uploaded") + " / " + (item.embeddingStatus || "pending");
      row.append(name, meta);
      list.appendChild(row);
    }
  }

  function mountUI() {
    if (document.getElementById(launcherID)) {
      return;
    }
    ensureStyle();

    const launcher = document.createElement("button");
    launcher.id = launcherID;
    launcher.className = "cyops-launcher";
    launcher.type = "button";
    launcher.setAttribute("aria-label", "Open CYOps");
    launcher.setAttribute("data-cyops-launcher", "true");
    launcher.textContent = "CYOps";
    Object.assign(launcher.style, {
      position: "fixed",
      right: "22px",
      bottom: "22px",
      zIndex: "2147483647",
      minWidth: "76px",
      height: "48px",
      borderRadius: "8px",
      border: "1px solid rgba(255,255,255,.58)",
      background: "#151515",
      color: "#fff",
      boxShadow: "0 10px 30px rgba(0,0,0,.42)",
      font: "700 14px 'Red Hat Text', system-ui, sans-serif",
      cursor: "pointer",
      padding: "0 16px",
      letterSpacing: "0",
    });

    const drawer = document.createElement("section");
    drawer.id = drawerID;
    drawer.className = "cyops-drawer";
    drawer.setAttribute("aria-label", "CYOps chat");
    drawer.innerHTML = '<header><h2>CYOps</h2><button class="cyops-icon-button" type="button" aria-label="Close CYOps" data-cyops-close>x</button></header><div class="cyops-body"><div class="cyops-chat"><div class="cyops-messages" data-cyops-messages></div><form class="cyops-compose" data-cyops-compose><textarea name="message" placeholder="Ask a question..." aria-label="Ask CYOps"></textarea><button class="cyops-send" type="submit" aria-label="Send">Send</button></form></div><aside class="cyops-docs"><h3>Documents</h3><label class="cyops-file"><input type="file" data-cyops-upload></label><div class="cyops-doc-list" data-cyops-doc-list></div></aside></div><div class="cyops-status" data-cyops-status>Ready</div>';

    document.body.append(launcher, drawer);

    const messages = drawer.querySelector("[data-cyops-messages]");
    const status = drawer.querySelector("[data-cyops-status]");
    messages.appendChild(createMessage("CYOps is ready.", "assistant"));

    launcher.addEventListener("click", async () => {
      const open = drawer.hasAttribute("open");
      if (open) {
        drawer.removeAttribute("open");
        return;
      }
      drawer.setAttribute("open", "");
      try {
        await refreshDocuments(drawer);
      } catch (error) {
        status.textContent = error.message;
      }
    });

    drawer.querySelector("[data-cyops-close]").addEventListener("click", () => {
      drawer.removeAttribute("open");
    });

    drawer.querySelector("[data-cyops-compose]").addEventListener("submit", async (event) => {
      event.preventDefault();
      const textarea = drawer.querySelector("textarea[name='message']");
      const message = textarea.value.trim();
      if (!message) {
        return;
      }
      textarea.value = "";
      messages.appendChild(createMessage(message, "user"));
      status.textContent = "Thinking";
      try {
        const response = await requestJSON("/api/chat", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ message: message, provider: "lightspeed", rag: { enabled: true } }),
        });
        messages.appendChild(createMessage(response.answer || "No answer returned.", "assistant"));
        status.textContent = "Ready";
      } catch (error) {
        messages.appendChild(createMessage(error.message, "assistant"));
        status.textContent = "Error";
      }
      messages.scrollTop = messages.scrollHeight;
    });

    drawer.querySelector("[data-cyops-upload]").addEventListener("change", async (event) => {
      const file = event.target.files && event.target.files[0];
      if (!file) {
        return;
      }
      const form = new FormData();
      form.append("file", file);
      status.textContent = "Uploading";
      try {
        await requestJSON("/api/documents", { method: "POST", body: form });
        await refreshDocuments(drawer);
        status.textContent = "Ready";
      } catch (error) {
        status.textContent = error.message;
      } finally {
        event.target.value = "";
      }
    });
  }

  function start() {
    if (document.readyState === "loading") {
      document.addEventListener("DOMContentLoaded", mountUI, { once: true });
      return;
    }
    mountUI();
  }

  function markEntryLoaded() {
    document.documentElement.setAttribute("data-cyops-plugin-entry", "0.0.48");
  }

  function cyopsLauncherFlag() {
    markEntryLoaded();
    start();
    return Promise.resolve({ CYOPS_CONSOLE_LAUNCHER: true });
  }

  const pluginEntry = {
    cyopsLauncherFlag: function () {
      return Promise.resolve(function () {
        return cyopsLauncherFlag;
      });
    },
  };
  const registerPluginEntry = typeof loadPluginEntry === "function"
    ? loadPluginEntry
    : typeof window.loadPluginEntry === "function"
      ? window.loadPluginEntry
      : null;

  if (registerPluginEntry) {
    registerPluginEntry("cyops-console@0.0.48", pluginEntry);
  }
  markEntryLoaded();
  start();
  window.setTimeout(start, 250);
  window.setTimeout(start, 1000);
})();`

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
