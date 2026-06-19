package appserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzReturnsOK(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("content type = %q, want application/json", contentType)
	}
	if !strings.Contains(recorder.Body.String(), `"status":"ok"`) {
		t.Fatalf("body = %q, want status ok", recorder.Body.String())
	}
}

func TestHealthzRejectsNonGet(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}

func TestRetrievalMetricsEndpointReturnsSnapshot(t *testing.T) {
	metrics := NewRetrievalMetrics()
	metrics.ObserveRetrieval(RetrievalObservation{
		Mode:        "pgvector",
		ResultCount: 2,
	})
	server := NewServerWithOptions(ServerOptions{Metrics: metrics})
	request := httptest.NewRequest(http.MethodGet, "/api/ops/retrieval-metrics", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`"total":1`,
		`"byMode":{"pgvector":1}`,
		`"last":{"mode":"pgvector"`,
		`"resultCount":2`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}
}

func TestDiagnosticsEndpointRequiresAdmin(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{
		AdminAuth: AdminAuthConfig{Users: []string{"admin"}},
	})
	request := httptest.NewRequest(http.MethodGet, "/api/ops/diagnostics", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestDiagnosticsEndpointAllowsLoopbackDevUser(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{
		AdminAuth: AdminAuthConfig{Users: []string{"admin"}},
		DevUser:   "admin",
	})
	request := httptest.NewRequest(http.MethodGet, "/api/ops/diagnostics", nil)
	request.RemoteAddr = "127.0.0.1:54321"
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), `"authorizedUser":"admin"`) {
		t.Fatalf("body = %q, want local dev admin user", recorder.Body.String())
	}
}

func TestDiagnosticsEndpointRejectsDevUserFromNonLoopback(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{
		AdminAuth: AdminAuthConfig{Users: []string{"admin"}},
		DevUser:   "admin",
	})
	request := httptest.NewRequest(http.MethodGet, "/api/ops/diagnostics", nil)
	request.RemoteAddr = "192.0.2.10:54321"
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestConsoleDiagnosticsViewUsesConsoleBackendPath(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodGet, "/console-plugin/diagnostics", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "text/html; charset=utf-8" {
		t.Fatalf("content type = %q, want text/html", contentType)
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`CYOps Diagnostics`,
		`/console-plugin/diagnostics.js`,
		`data-cyops-view="diagnostics"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}
}

func TestConsolePluginManifestIsServed(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodGet, "/plugin-manifest.json", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("content type = %q, want application/json", contentType)
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`"name": "cyops-console"`,
		`"version": "0.0.49"`,
		`"customProperties":`,
		`"displayName": "CYOps"`,
		`"baseURL": "/api/plugins/cyops-console/"`,
		`"loadScripts":`,
		`"plugin-entry.js"`,
		`"registrationMethod": "callback"`,
		`"type": "console.flag"`,
		`"$codeRef": "cyopsLauncherFlag"`,
		`"console.navigation/href"`,
		`"/console-plugin/diagnostics"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}

	var manifest map[string]any
	if err := json.Unmarshal([]byte(body), &manifest); err != nil {
		t.Fatalf("manifest is not json: %v", err)
	}
	if _, ok := manifest["displayName"]; ok {
		t.Fatalf("manifest has top-level displayName, want displayName under customProperties.console")
	}
	if _, ok := manifest["description"]; ok {
		t.Fatalf("manifest has top-level description, want description under customProperties.console")
	}
}

func TestConsolePluginEntryIsServed(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodGet, "/console-plugin/plugin-entry.js", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if contentType := recorder.Header().Get("Content-Type"); contentType != "text/javascript; charset=utf-8" {
		t.Fatalf("content type = %q, want text/javascript", contentType)
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`cyops-console`,
		`0.0.49`,
		`loadPluginEntry`,
		`cyops-console@0.0.49`,
		`cyopsLauncherFlag`,
		`data-cyops-plugin-entry`,
		`pluginProxyBase`,
		`apiBase + path`,
		`X-CSRFToken`,
		`X-CSRF-Token`,
		`X-Requested-With`,
		`window.setTimeout(start, 1000)`,
		`data-cyops-launcher`,
		`CYOps`,
		`right: "22px"`,
		`bottom: "22px"`,
		`2147483647`,
		`/api/chat`,
		`/api/documents`,
		`provider: "lightspeed"`,
		`/console-plugin/diagnostics`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}
}

func TestConsoleDiagnosticsScriptCallsDiagnosticsAPIs(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodGet, "/console-plugin/diagnostics.js", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`/api/ops/diagnostics/schema`,
		`/api/ops/diagnostics`,
		`credentials: "same-origin"`,
		`OpenShift Web Console backend path`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}
	if strings.Contains(strings.ToLower(body), "oauth") {
		t.Fatalf("script = %q, want no oauth route handling in console path", body)
	}
}

func TestConsolePluginEntryRootAliasIsServed(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodGet, "/plugin-entry.js", nil)
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if !strings.Contains(recorder.Body.String(), `window.loadPluginEntry`) {
		t.Fatalf("body = %q, want callback entry", recorder.Body.String())
	}
}

func TestDiagnosticsEndpointReturnsSecretFreeOperationalSummary(t *testing.T) {
	metrics := NewRetrievalMetrics()
	metrics.ObserveRetrieval(RetrievalObservation{Mode: "pgvector", ResultCount: 3})
	repository := NewMemoryDocumentRepository()
	repository.Create("runbook.md", 12, "admin")
	server := NewServerWithOptions(ServerOptions{
		Documents: repository,
		Metrics:   metrics,
		AdminAuth: AdminAuthConfig{Groups: []string{"cyops-admins"}},
	})
	request := httptest.NewRequest(http.MethodGet, "/api/ops/diagnostics", nil)
	request.Header.Set("X-Forwarded-User", "admin")
	request.Header.Set("X-Forwarded-Groups", "system:authenticated, cyops-admins")
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`"retrieval":{"total":1`,
		`"documents":{"total":1`,
		`"byStatus":{"uploaded":1}`,
		`"byEmbeddingStatus":{"pending":1}`,
		`"authorizedUser":"admin"`,
		`"authorizedGroups":["system:authenticated","cyops-admins"]`,
		`"available":false`,
		`"retrievalMetrics":"/api/ops/retrieval-metrics"`,
		`"primarySurface":"openshift-web-console"`,
		`"ui":{"title":"CYOps Diagnostics","primaryEntry":"openshift-web-console","fallbackRoute":"optional"}`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}
	if strings.Contains(strings.ToLower(body), "token") || strings.Contains(strings.ToLower(body), "dsn") {
		t.Fatalf("body = %q, want no token or dsn fields", body)
	}
}

func TestDiagnosticsSchemaEndpointReturnsAggregateOnlyContract(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{
		AdminAuth: AdminAuthConfig{Users: []string{"admin"}},
	})
	request := httptest.NewRequest(http.MethodGet, "/api/ops/diagnostics/schema", nil)
	request.Header.Set("X-Forwarded-User", "admin")
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := recorder.Body.String()
	for _, want := range []string{
		`"version":"v0.0.23"`,
		`"primaryEntry":"openshift-web-console"`,
		`"aggregateOnly":true`,
		`"documentContent"`,
		`"postgresDsn"`,
		`"adminToken"`,
	} {
		if !strings.Contains(body, want) {
			t.Fatalf("body = %q, want %q", body, want)
		}
	}
}

func TestReembedEndpointRequiresPostgresRepository(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{AdminToken: "admin-token"})
	request := httptest.NewRequest(http.MethodPost, "/api/ops/reembed", strings.NewReader(`{"limit":1}`))
	request.Header.Set("X-CYOps-Admin-Token", "admin-token")
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusConflict, recorder.Body.String())
	}
	if !strings.Contains(recorder.Body.String(), "postgres document repository") {
		t.Fatalf("body = %q, want postgres repository message", recorder.Body.String())
	}
}

func TestReembedEndpointRequiresAdminToken(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{AdminToken: "admin-token"})
	request := httptest.NewRequest(http.MethodPost, "/api/ops/reembed", strings.NewReader(`{"limit":1}`))
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusForbidden)
	}
}

func TestReembedEndpointAllowsForwardedAdminUser(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{
		AdminAuth: AdminAuthConfig{Users: []string{"cluster-admin"}},
	})
	request := httptest.NewRequest(http.MethodPost, "/api/ops/reembed", strings.NewReader(`{"limit":1}`))
	request.Header.Set("X-Forwarded-User", "cluster-admin")
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusConflict, recorder.Body.String())
	}
}

func TestReembedEndpointAllowsForwardedAdminGroup(t *testing.T) {
	server := NewServerWithOptions(ServerOptions{
		AdminAuth: AdminAuthConfig{Groups: []string{"cyops-admins"}},
	})
	request := httptest.NewRequest(http.MethodPost, "/api/ops/reembed", strings.NewReader(`{"limit":1}`))
	request.Header.Set("X-Forwarded-Groups", "system:authenticated, cyops-admins")
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusConflict, recorder.Body.String())
	}
}

func TestChatRoutesToProvider(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":"Why is this pod not ready?","provider":"lightspeed"}`))
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	body := recorder.Body.String()
	if !strings.Contains(body, `"provider":"lightspeed"`) {
		t.Fatalf("body = %q, want provider lightspeed", body)
	}
	if !strings.Contains(body, `Mocked CYOps response`) {
		t.Fatalf("body = %q, want mocked answer", body)
	}
}

func TestChatRejectsMissingMessage(t *testing.T) {
	server := NewServer()
	request := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":""}`))
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestChatRAGAddsRetrieverContextAndCitations(t *testing.T) {
	provider := &capturingProvider{answer: "Use the cited runbook."}
	server := NewServerWithOptions(ServerOptions{
		Provider:  provider,
		Documents: NewMemoryDocumentRepository(),
		Retriever: staticRetriever{
			result: RetrievalResult{
				Context: []ProviderContext{{
					Type:   "rag_chunk",
					Text:   "Check pod status before restart.",
					Source: "doc-001/chunk-001",
				}},
				Citations: []Citation{{
					DocumentID: "doc-001",
					Title:      "runbook.md",
					ChunkID:    "chunk-001",
				}},
			},
		},
	})

	request := httptest.NewRequest(http.MethodPost, "/api/chat", strings.NewReader(`{"message":"What should I check?","rag":{"enabled":true,"documentIds":["doc-001"]}}`))
	recorder := httptest.NewRecorder()

	server.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d: %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}
	if len(provider.request.Context) != 1 {
		t.Fatalf("provider context len = %d, want 1", len(provider.request.Context))
	}
	if provider.request.Context[0].Text != "Check pod status before restart." {
		t.Fatalf("provider context text = %q", provider.request.Context[0].Text)
	}
	if !strings.Contains(recorder.Body.String(), `"citations":[{"documentId":"doc-001","title":"runbook.md","chunkId":"chunk-001"}]`) {
		t.Fatalf("body = %q, want citation", recorder.Body.String())
	}
}

type capturingProvider struct {
	answer  string
	request ProviderRequest
}

func (p *capturingProvider) Chat(request ProviderRequest) (ProviderResponse, error) {
	p.request = request
	return ProviderResponse{Answer: p.answer, RawProvider: "test"}, nil
}

type staticRetriever struct {
	result RetrievalResult
	err    error
}

func (r staticRetriever) Retrieve(context.Context, RetrievalRequest) (RetrievalResult, error) {
	return r.result, r.err
}

func TestDocumentsUploadListDetailAndDelete(t *testing.T) {
	server := NewServer()

	upload := multipartRequest(t, "file", "runbook.txt", "check pod status")
	upload.Header.Set("X-Forwarded-User", "admin")
	uploadRecorder := httptest.NewRecorder()
	server.ServeHTTP(uploadRecorder, upload)
	if uploadRecorder.Code != http.StatusCreated {
		t.Fatalf("upload status = %d, want %d: %s", uploadRecorder.Code, http.StatusCreated, uploadRecorder.Body.String())
	}
	if !strings.Contains(uploadRecorder.Body.String(), `"id":"doc-001"`) {
		t.Fatalf("upload body = %q, want doc-001", uploadRecorder.Body.String())
	}

	list := httptest.NewRequest(http.MethodGet, "/api/documents", nil)
	listRecorder := httptest.NewRecorder()
	server.ServeHTTP(listRecorder, list)
	if listRecorder.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listRecorder.Code, http.StatusOK)
	}
	if !strings.Contains(listRecorder.Body.String(), `"filename":"runbook.txt"`) {
		t.Fatalf("list body = %q, want uploaded filename", listRecorder.Body.String())
	}

	detail := httptest.NewRequest(http.MethodGet, "/api/documents/doc-001", nil)
	detailRecorder := httptest.NewRecorder()
	server.ServeHTTP(detailRecorder, detail)
	if detailRecorder.Code != http.StatusOK {
		t.Fatalf("detail status = %d, want %d", detailRecorder.Code, http.StatusOK)
	}
	if !strings.Contains(detailRecorder.Body.String(), `"embeddingStatus":"pending"`) {
		t.Fatalf("detail body = %q, want pending embedding status", detailRecorder.Body.String())
	}

	deleteRequest := httptest.NewRequest(http.MethodDelete, "/api/documents/doc-001", nil)
	deleteRecorder := httptest.NewRecorder()
	server.ServeHTTP(deleteRecorder, deleteRequest)
	if deleteRecorder.Code != http.StatusOK {
		t.Fatalf("delete status = %d, want %d", deleteRecorder.Code, http.StatusOK)
	}
	if !strings.Contains(deleteRecorder.Body.String(), `"status":"deleting"`) {
		t.Fatalf("delete body = %q, want deleting", deleteRecorder.Body.String())
	}
}

func TestDocumentsUploadStoresFileAndPersistsObjectURI(t *testing.T) {
	repository := NewMemoryDocumentRepository()
	server := NewServerWithOptions(ServerOptions{
		Provider:  MockProvider{},
		Documents: repository,
		Storage:   LocalDocumentStorage{BasePath: t.TempDir()},
	})

	upload := multipartRequest(t, "file", "../runbook.txt", "check pod status")
	uploadRecorder := httptest.NewRecorder()
	server.ServeHTTP(uploadRecorder, upload)
	if uploadRecorder.Code != http.StatusCreated {
		t.Fatalf("upload status = %d, want %d: %s", uploadRecorder.Code, http.StatusCreated, uploadRecorder.Body.String())
	}

	items := repository.List()
	if len(items) != 1 {
		t.Fatalf("documents len = %d, want 1", len(items))
	}
	if items[0].ObjectURI == "" {
		t.Fatal("object uri is empty")
	}
	if strings.Contains(items[0].ObjectURI, "..") {
		t.Fatalf("object uri = %q, want sanitized path", items[0].ObjectURI)
	}
	if items[0].SizeBytes != int64(len("check pod status")) {
		t.Fatalf("size = %d, want uploaded content size", items[0].SizeBytes)
	}
}

func TestDocumentsUploadDoesNotCreateMetadataWhenStorageFails(t *testing.T) {
	repository := NewMemoryDocumentRepository()
	server := NewServerWithOptions(ServerOptions{
		Provider:  MockProvider{},
		Documents: repository,
		Storage:   failingStorage{},
	})

	upload := multipartRequest(t, "file", "runbook.txt", "check pod status")
	uploadRecorder := httptest.NewRecorder()
	server.ServeHTTP(uploadRecorder, upload)
	if uploadRecorder.Code != http.StatusInternalServerError {
		t.Fatalf("upload status = %d, want %d: %s", uploadRecorder.Code, http.StatusInternalServerError, uploadRecorder.Body.String())
	}
	if len(repository.List()) != 0 {
		t.Fatalf("documents len = %d, want 0", len(repository.List()))
	}
}

type failingStorage struct{}

func (failingStorage) Store(context.Context, string, string, io.Reader) (StoredObject, error) {
	return StoredObject{}, errors.New("storage failed")
}

func multipartRequest(t *testing.T, field string, filename string, content string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile(field, filename)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := part.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodPost, "/api/documents", &body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	return request
}
