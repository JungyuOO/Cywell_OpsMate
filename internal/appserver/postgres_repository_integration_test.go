package appserver

import (
	"context"
	"database/sql"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresDocumentRepositoryIntegration(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.now = func() time.Time { return time.Date(2026, 6, 14, 8, 0, 0, 0, time.UTC) }
	repository.newID = func() (string, error) { return "00000000-0000-4000-8000-000000000001", nil }

	created, err := repository.CreateStored(ctx, CreateStoredDocumentInput{
		Filename:   "runbook.pdf",
		SizeBytes:  42,
		ObjectURI:  "/var/lib/cyops/documents/runbook.pdf",
		UploadedBy: "admin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Status != "uploaded" {
		t.Fatalf("created status = %q, want uploaded", created.Status)
	}
	if created.ObjectURI == "" {
		t.Fatal("created object uri is empty")
	}

	list, err := repository.ListContext(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("list len = %d, want 1", len(list))
	}

	found, err := repository.GetContext(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.Filename != "runbook.pdf" {
		t.Fatalf("filename = %q, want runbook.pdf", found.Filename)
	}
	if found.ObjectURI != "/var/lib/cyops/documents/runbook.pdf" {
		t.Fatalf("object uri = %q", found.ObjectURI)
	}

	deleting, err := repository.MarkDeletingContext(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deleting.Status != "deleting" {
		t.Fatalf("deleting status = %q, want deleting", deleting.Status)
	}
}

func TestPostgresStorageUploadIntegration(t *testing.T) {
	dsn := os.Getenv("CYOPS_POSTGRES_TEST_DSN")
	if dsn == "" {
		t.Skip("CYOPS_POSTGRES_TEST_DSN is not set")
	}

	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}
	if err := ApplyMigrations(ctx, db); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	server := NewServerWithOptions(ServerOptions{
		Provider:  MockProvider{},
		Documents: NewPostgresDocumentRepository(db, "opsmate"),
		Storage:   LocalDocumentStorage{BasePath: t.TempDir()},
	})

	upload := multipartRequest(t, "file", "../runbook.txt", "check pod status")
	upload.Header.Set("X-Forwarded-User", "admin")
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, upload)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("upload status = %d, want %d: %s", recorder.Code, http.StatusCreated, recorder.Body.String())
	}

	var filename string
	var objectURI string
	var sizeBytes int64
	if err := db.QueryRowContext(ctx, "SELECT filename, object_uri, size_bytes FROM cyops_documents WHERE namespace = $1", "opsmate").Scan(&filename, &objectURI, &sizeBytes); err != nil {
		t.Fatal(err)
	}
	if filename != "runbook.txt" {
		t.Fatalf("filename = %q, want multipart filename", filename)
	}
	if objectURI == "" || strings.Contains(objectURI, "..") {
		t.Fatalf("object uri = %q, want sanitized stored path", objectURI)
	}
	if sizeBytes != int64(len("check pod status")) {
		t.Fatalf("size = %d, want uploaded content size", sizeBytes)
	}
}
