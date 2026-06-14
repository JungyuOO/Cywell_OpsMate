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
