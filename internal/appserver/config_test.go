package appserver

import (
	"context"
	"testing"
)

func TestLoadConfigFromEnv(t *testing.T) {
	t.Setenv(envPostgresDSN, "postgres://cyops:cyops@postgres/cyops")
	t.Setenv(envNamespace, "opsmate")
	t.Setenv(envDocumentStoragePath, "/var/lib/cyops/documents")
	t.Setenv(envLightspeedEndpoint, "https://lightspeed.example/api/chat")
	t.Setenv(envEmbeddingEndpoint, "https://embedding.example/embed")
	t.Setenv(envEmbeddingModel, "embedding-model")
	t.Setenv(envEmbeddingDimensions, "8")
	t.Setenv(envEmbeddingToken, "secret-token")
	t.Setenv(envPGVectorRequired, "true")

	config := LoadConfigFromEnv()

	if config.PostgresDSN != "postgres://cyops:cyops@postgres/cyops" {
		t.Fatalf("postgres dsn = %q", config.PostgresDSN)
	}
	if config.Namespace != "opsmate" {
		t.Fatalf("namespace = %q, want opsmate", config.Namespace)
	}
	if config.DocumentStoragePath != "/var/lib/cyops/documents" {
		t.Fatalf("storage path = %q", config.DocumentStoragePath)
	}
	if config.LightspeedEndpoint != "https://lightspeed.example/api/chat" {
		t.Fatalf("lightspeed endpoint = %q", config.LightspeedEndpoint)
	}
	if config.EmbeddingEndpoint != "https://embedding.example/embed" {
		t.Fatalf("embedding endpoint = %q", config.EmbeddingEndpoint)
	}
	if config.EmbeddingModel != "embedding-model" {
		t.Fatalf("embedding model = %q", config.EmbeddingModel)
	}
	if config.EmbeddingDimensions != 8 {
		t.Fatalf("embedding dimensions = %d, want 8", config.EmbeddingDimensions)
	}
	if config.EmbeddingToken != "secret-token" {
		t.Fatalf("embedding token = %q", config.EmbeddingToken)
	}
	if !config.PGVectorRequired {
		t.Fatal("pgvector required = false, want true")
	}
}

func TestNewServerFromConfigUsesMemoryWithoutPostgres(t *testing.T) {
	server, err := NewServerFromConfig(context.Background(), AppConfig{
		Namespace:           "opsmate",
		DocumentStoragePath: t.TempDir(),
	})
	if err != nil {
		t.Fatal(err)
	}
	if server == nil {
		t.Fatal("server is nil")
	}
}
