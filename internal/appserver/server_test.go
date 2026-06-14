package appserver

import (
	"bytes"
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
