package appserver

import (
	"context"
	"database/sql"
	"os"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	envPostgresDSN         = "CYOPS_POSTGRES_DSN"
	envNamespace           = "CYOPS_NAMESPACE"
	envDocumentStoragePath = "CYOPS_DOCUMENT_STORAGE_PATH"
	envLightspeedEndpoint  = "CYOPS_LIGHTSPEED_ENDPOINT"
)

type AppConfig struct {
	PostgresDSN         string
	Namespace           string
	DocumentStoragePath string
	LightspeedEndpoint  string
}

func LoadConfigFromEnv() AppConfig {
	namespace := strings.TrimSpace(os.Getenv(envNamespace))
	if namespace == "" {
		namespace = "default"
	}
	return AppConfig{
		PostgresDSN:         strings.TrimSpace(os.Getenv(envPostgresDSN)),
		Namespace:           namespace,
		DocumentStoragePath: strings.TrimSpace(os.Getenv(envDocumentStoragePath)),
		LightspeedEndpoint:  strings.TrimSpace(os.Getenv(envLightspeedEndpoint)),
	}
}

func NewServerFromConfig(ctx context.Context, config AppConfig) (*Server, error) {
	documents := DocumentRepository(NewMemoryDocumentRepository())
	if config.PostgresDSN != "" {
		db, err := sql.Open("pgx", config.PostgresDSN)
		if err != nil {
			return nil, err
		}
		if err := db.PingContext(ctx); err != nil {
			_ = db.Close()
			return nil, err
		}
		if err := ApplyMigrations(ctx, db); err != nil {
			_ = db.Close()
			return nil, err
		}
		documents = NewPostgresDocumentRepository(db, config.Namespace)
	}

	var storage DocumentStorage
	if config.DocumentStoragePath != "" {
		storage = LocalDocumentStorage{BasePath: config.DocumentStoragePath}
	}

	return NewServerWithOptions(ServerOptions{
		Provider: LightspeedProvider{
			Config: LightspeedProviderConfig{EndpointURL: config.LightspeedEndpoint},
		},
		Documents: documents,
		Storage:   storage,
	}), nil
}
