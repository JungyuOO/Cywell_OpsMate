package appserver

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	envPostgresDSN         = "CYOPS_POSTGRES_DSN"
	envNamespace           = "CYOPS_NAMESPACE"
	envDocumentStoragePath = "CYOPS_DOCUMENT_STORAGE_PATH"
	envLightspeedEndpoint  = "CYOPS_LIGHTSPEED_ENDPOINT"
	envEmbeddingEndpoint   = "CYOPS_EMBEDDING_ENDPOINT"
	envEmbeddingModel      = "CYOPS_EMBEDDING_MODEL"
	envEmbeddingDimensions = "CYOPS_EMBEDDING_DIMENSIONS"
	envEmbeddingToken      = "CYOPS_EMBEDDING_TOKEN"
	envPGVectorRequired    = "CYOPS_PGVECTOR_REQUIRED"
)

type AppConfig struct {
	PostgresDSN         string
	Namespace           string
	DocumentStoragePath string
	LightspeedEndpoint  string
	EmbeddingEndpoint   string
	EmbeddingModel      string
	EmbeddingDimensions int
	EmbeddingToken      string
	PGVectorRequired    bool
}

func LoadConfigFromEnv() AppConfig {
	namespace := strings.TrimSpace(os.Getenv(envNamespace))
	if namespace == "" {
		namespace = "default"
	}
	dimensions, _ := strconv.Atoi(strings.TrimSpace(os.Getenv(envEmbeddingDimensions)))
	return AppConfig{
		PostgresDSN:         strings.TrimSpace(os.Getenv(envPostgresDSN)),
		Namespace:           namespace,
		DocumentStoragePath: strings.TrimSpace(os.Getenv(envDocumentStoragePath)),
		LightspeedEndpoint:  strings.TrimSpace(os.Getenv(envLightspeedEndpoint)),
		EmbeddingEndpoint:   strings.TrimSpace(os.Getenv(envEmbeddingEndpoint)),
		EmbeddingModel:      strings.TrimSpace(os.Getenv(envEmbeddingModel)),
		EmbeddingDimensions: dimensions,
		EmbeddingToken:      strings.TrimSpace(os.Getenv(envEmbeddingToken)),
		PGVectorRequired:    parseBool(os.Getenv(envPGVectorRequired)),
	}
}

func NewServerFromConfig(ctx context.Context, config AppConfig) (*Server, error) {
	documents := DocumentRepository(NewMemoryDocumentRepository())
	var retriever Retriever
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
		if config.PGVectorRequired {
			if err := CheckPGVectorReady(ctx, db); err != nil {
				_ = db.Close()
				return nil, err
			}
		}
		repository := NewPostgresDocumentRepository(db, config.Namespace)
		documents = repository
		retriever = PostgresRetriever{
			Repository: repository,
			Embedder:   NewEmbeddingProviderFromConfig(config),
		}
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
		Retriever: retriever,
	}), nil
}

func NewEmbeddingProviderFromConfig(config AppConfig) EmbeddingProvider {
	if config.EmbeddingEndpoint == "" {
		return DeterministicEmbeddingProvider{
			Model:      config.EmbeddingModel,
			Dimensions: config.EmbeddingDimensions,
		}
	}
	return HTTPEmbeddingProvider{
		EndpointURL: config.EmbeddingEndpoint,
		Model:       config.EmbeddingModel,
		Dimensions:  config.EmbeddingDimensions,
		Token:       config.EmbeddingToken,
	}
}

func parseBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y", "on":
		return true
	default:
		return false
	}
}
