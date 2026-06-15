package main

import (
	"os"
	"testing"
)

func TestDefaultListenAddress(t *testing.T) {
	if defaultListenAddress != ":8080" {
		t.Fatalf("default listen address = %q, want :8080", defaultListenAddress)
	}
}

func TestMainDoesNotRequirePostgresByDefault(t *testing.T) {
	t.Setenv("CYOPS_POSTGRES_DSN", "")
	t.Setenv("CYOPS_LISTEN_ADDRESS", "127.0.0.1:0")
	if os.Getenv("CYOPS_POSTGRES_DSN") != "" {
		t.Fatal("postgres dsn is set")
	}
}
