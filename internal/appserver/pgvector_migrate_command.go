package appserver

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGVectorMigrationCommandConfig struct {
	PostgresDSN string
	Dimensions  int
}

func RunPGVectorMigrationCommand(ctx context.Context, config PGVectorMigrationCommandConfig) error {
	if strings.TrimSpace(config.PostgresDSN) == "" {
		return fmt.Errorf("postgres dsn is required")
	}
	if config.Dimensions <= 0 {
		return fmt.Errorf("embedding dimensions are required")
	}

	db, err := sql.Open("pgx", config.PostgresDSN)
	if err != nil {
		return fmt.Errorf("open postgres connection: %w", err)
	}
	defer db.Close()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect to postgres: %w", err)
	}
	if err := ApplyPGVectorEmbeddingMigration(ctx, db, config.Dimensions); err != nil {
		return err
	}
	return nil
}

func PGVectorMigrationCommandConfigFromEnv(getenv func(string) string) (PGVectorMigrationCommandConfig, error) {
	dimensions, err := strconv.Atoi(strings.TrimSpace(getenv(envEmbeddingDimensions)))
	if err != nil || dimensions <= 0 {
		return PGVectorMigrationCommandConfig{}, fmt.Errorf("embedding dimensions are required")
	}
	return PGVectorMigrationCommandConfig{
		PostgresDSN: strings.TrimSpace(getenv(envPostgresDSN)),
		Dimensions:  dimensions,
	}, nil
}
