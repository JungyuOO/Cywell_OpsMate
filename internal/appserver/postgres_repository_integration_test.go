package appserver

import (
	"context"
	"database/sql"
	"os"
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
	migration, err := os.ReadFile("migrations/0001_cyops_rag_schema.sql")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, string(migration)); err != nil {
		t.Fatal(err)
	}
	if _, err := db.ExecContext(ctx, "TRUNCATE cyops_chat_messages, cyops_chat_sessions, cyops_document_embeddings, cyops_document_chunks, cyops_documents"); err != nil {
		t.Fatal(err)
	}

	repository := NewPostgresDocumentRepository(db, "opsmate")
	repository.now = func() time.Time { return time.Date(2026, 6, 14, 8, 0, 0, 0, time.UTC) }
	repository.newID = func() (string, error) { return "00000000-0000-4000-8000-000000000001", nil }

	created, err := repository.CreateContext(ctx, "runbook.pdf", 42, "admin")
	if err != nil {
		t.Fatal(err)
	}
	if created.Status != "uploaded" {
		t.Fatalf("created status = %q, want uploaded", created.Status)
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

	deleting, err := repository.MarkDeletingContext(ctx, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if deleting.Status != "deleting" {
		t.Fatalf("deleting status = %q, want deleting", deleting.Status)
	}
}
