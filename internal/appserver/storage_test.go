package appserver

import (
	"context"
	"strings"
	"testing"
)

func TestLocalDocumentStorageBuildsSafeURI(t *testing.T) {
	storage := LocalDocumentStorage{BasePath: "/var/lib/cyops/documents"}

	object, err := storage.Store(context.Background(), "doc-001", "../runbook.txt", strings.NewReader("hello"))
	if err != nil {
		t.Fatal(err)
	}
	if object.URI != "/var/lib/cyops/documents/doc-001/runbook.txt" {
		t.Fatalf("uri = %q, want sanitized document path", object.URI)
	}
	if object.SizeBytes != 5 {
		t.Fatalf("size = %d, want 5", object.SizeBytes)
	}
}

func TestLocalDocumentStorageRequiresBasePath(t *testing.T) {
	storage := LocalDocumentStorage{}

	if _, err := storage.Store(context.Background(), "doc-001", "runbook.txt", strings.NewReader("hello")); err == nil {
		t.Fatal("expected error for empty base path")
	}
}
