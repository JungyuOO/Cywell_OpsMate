package appserver

import (
	"context"
	"database/sql"
	"os"
	"strconv"
	"strings"
	"time"

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
	envRetrievalMode       = "CYOPS_RETRIEVAL_MODE"
	envRetrievalSlowMS     = "CYOPS_RETRIEVAL_SLOW_THRESHOLD_MS"
	envAdminToken          = "CYOPS_ADMIN_TOKEN"
	envAdminUsers          = "CYOPS_ADMIN_USERS"
	envAdminGroups         = "CYOPS_ADMIN_GROUPS"
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
	RetrievalMode       string
	RetrievalSlow       time.Duration
	AdminToken          string
	AdminUsers          []string
	AdminGroups         []string
}

func LoadConfigFromEnv() AppConfig {
	namespace := strings.TrimSpace(os.Getenv(envNamespace))
	if namespace == "" {
		namespace = "default"
	}
	dimensions, _ := strconv.Atoi(strings.TrimSpace(os.Getenv(envEmbeddingDimensions)))
	slowMS, _ := strconv.Atoi(strings.TrimSpace(os.Getenv(envRetrievalSlowMS)))
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
		RetrievalMode:       retrievalModeOrDefault(os.Getenv(envRetrievalMode)),
		RetrievalSlow:       time.Duration(slowMS) * time.Millisecond,
		AdminToken:          strings.TrimSpace(os.Getenv(envAdminToken)),
		AdminUsers:          splitCSV(os.Getenv(envAdminUsers)),
		AdminGroups:         splitCSV(os.Getenv(envAdminGroups)),
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
		metrics := NewRetrievalMetrics()
		documents = repository
		retriever = PostgresRetriever{
			Repository: repository,
			Embedder:   NewEmbeddingProviderFromConfig(config),
			Mode:       config.RetrievalMode,
			Observer:   metrics,
			SlowAfter:  config.RetrievalSlow,
		}
		embedder := NewEmbeddingProviderFromConfig(config)
		return NewServerWithOptions(ServerOptions{
			Provider: LightspeedProvider{
				Config: LightspeedProviderConfig{EndpointURL: config.LightspeedEndpoint},
			},
			Documents: documents,
			Storage:   storageFromConfig(config),
			Retriever: retriever,
			Metrics:   metrics,
			Embedder:  embedder,
			AdminAuth: adminAuthFromConfig(config),
		}), nil
	}

	return NewServerWithOptions(ServerOptions{
		Provider: LightspeedProvider{
			Config: LightspeedProviderConfig{EndpointURL: config.LightspeedEndpoint},
		},
		Documents: documents,
		Storage:   storageFromConfig(config),
		Retriever: retriever,
		AdminAuth: adminAuthFromConfig(config),
	}), nil
}

func adminAuthFromConfig(config AppConfig) AdminAuthConfig {
	return AdminAuthConfig{
		Token:  config.AdminToken,
		Users:  config.AdminUsers,
		Groups: config.AdminGroups,
	}
}

func storageFromConfig(config AppConfig) DocumentStorage {
	if config.DocumentStoragePath == "" {
		return nil
	}
	return LocalDocumentStorage{BasePath: config.DocumentStoragePath}
}

func retrievalModeOrDefault(value string) string {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "pgvector":
		return "pgvector"
	default:
		return "bytea"
	}
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

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}
