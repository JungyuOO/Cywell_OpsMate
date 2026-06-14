package appserver

import (
	"context"
	"strings"
	"testing"
)

func TestPGVectorEmbeddingMigrationSQLValidatesDimensions(t *testing.T) {
	if _, err := PGVectorEmbeddingMigrationSQL(0); err == nil {
		t.Fatal("expected error for zero dimensions")
	}
	if _, err := PGVectorEmbeddingMigrationSQL(20000); err == nil {
		t.Fatal("expected error for excessive dimensions")
	}
}

func TestPGVectorEmbeddingMigrationSQLBuildsVectorMigration(t *testing.T) {
	sql, err := PGVectorEmbeddingMigrationSQL(768)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(sql, "CREATE EXTENSION IF NOT EXISTS vector") {
		t.Fatalf("sql = %q, want vector extension", sql)
	}
	if !strings.Contains(sql, "VECTOR(768)") {
		t.Fatalf("sql = %q, want VECTOR(768)", sql)
	}
	if !strings.Contains(sql, "USING NULL::VECTOR(768)") {
		t.Fatalf("sql = %q, want reset strategy", sql)
	}
}

func TestApplyPGVectorEmbeddingMigrationValidatesDimensionsBeforeExec(t *testing.T) {
	if err := ApplyPGVectorEmbeddingMigration(context.Background(), nil, 0); err == nil {
		t.Fatal("expected dimensions error")
	}
}
