package appserver

import (
	"context"
	"database/sql"
	"fmt"
)

func CheckPGVectorReady(ctx context.Context, db *sql.DB) error {
	if _, err := db.ExecContext(ctx, "CREATE EXTENSION IF NOT EXISTS vector"); err != nil {
		return fmt.Errorf("pgvector extension is not ready: %w", err)
	}
	return nil
}

func ApplyPGVectorEmbeddingMigration(ctx context.Context, db *sql.DB, dimensions int) error {
	query, err := PGVectorEmbeddingMigrationSQL(dimensions)
	if err != nil {
		return err
	}
	if err := CheckPGVectorReady(ctx, db); err != nil {
		return err
	}
	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("pgvector embedding migration failed: %w", err)
	}
	return nil
}

func PGVectorEmbeddingMigrationSQL(dimensions int) (string, error) {
	if dimensions <= 0 {
		return "", fmt.Errorf("embedding dimensions must be positive")
	}
	if dimensions > 16384 {
		return "", fmt.Errorf("embedding dimensions exceed supported maximum")
	}
	return fmt.Sprintf(`CREATE EXTENSION IF NOT EXISTS vector;
ALTER TABLE cyops_document_embeddings
    ALTER COLUMN embedding DROP DEFAULT,
    ALTER COLUMN embedding TYPE VECTOR(%d)
    USING NULL::VECTOR(%d);`, dimensions, dimensions), nil
}
