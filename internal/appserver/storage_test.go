package appserver

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLocalDocumentStorageBuildsSafeURI(t *testing.T) {
	basePath := t.TempDir()
	storage := LocalDocumentStorage{BasePath: basePath}

	object, err := storage.Store(context.Background(), "doc-001", "../runbook.txt", strings.NewReader("hello"))
	if err != nil {
		t.Fatal(err)
	}
	wantURI := filepath.ToSlash(filepath.Join(basePath, "doc-001", "runbook.txt"))
	if object.URI != wantURI {
		t.Fatalf("uri = %q, want sanitized document path", object.URI)
	}
	if object.SizeBytes != 5 {
		t.Fatalf("size = %d, want 5", object.SizeBytes)
	}
	content, err := os.ReadFile(filepath.Join(basePath, "doc-001", "runbook.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "hello" {
		t.Fatalf("stored content = %q, want hello", string(content))
	}
}

func TestLocalDocumentStorageRequiresBasePath(t *testing.T) {
	storage := LocalDocumentStorage{}

	if _, err := storage.Store(context.Background(), "doc-001", "runbook.txt", strings.NewReader("hello")); err == nil {
		t.Fatal("expected error for empty base path")
	}
}
