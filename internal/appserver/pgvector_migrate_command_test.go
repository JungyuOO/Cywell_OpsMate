package appserver

import (
	"strings"
	"testing"
)

func TestPGVectorMigrationCommandConfigFromEnv(t *testing.T) {
	config, err := PGVectorMigrationCommandConfigFromEnv(func(key string) string {
		switch key {
		case envPostgresDSN:
			return " postgres://cyops:cyops@postgres/cyops "
		case envEmbeddingDimensions:
			return "768"
		default:
			return ""
		}
	})
	if err != nil {
		t.Fatal(err)
	}
	if config.PostgresDSN != "postgres://cyops:cyops@postgres/cyops" {
		t.Fatalf("dsn = %q", config.PostgresDSN)
	}
	if config.Dimensions != 768 {
		t.Fatalf("dimensions = %d", config.Dimensions)
	}
}

func TestPGVectorMigrationCommandConfigFromEnvRequiresDimensions(t *testing.T) {
	_, err := PGVectorMigrationCommandConfigFromEnv(func(string) string { return "" })
	if err == nil {
		t.Fatal("expected dimensions error")
	}
	if strings.Contains(err.Error(), "postgres://") {
		t.Fatalf("error leaks dsn: %v", err)
	}
}
