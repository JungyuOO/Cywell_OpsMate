package main

import "testing"

func TestDefaultListenAddress(t *testing.T) {
	if defaultListenAddress != ":8080" {
		t.Fatalf("default listen address = %q, want :8080", defaultListenAddress)
	}
}

func TestMainDoesNotRequirePostgresByDefault(t *testing.T) {
	t.Setenv("CYOPS_POSTGRES_DSN", "")
	t.Setenv("CYOPS_LISTEN_ADDRESS", "127.0.0.1:0")
	config := loadServeConfig(func(name string) string {
		switch name {
		case envListenAddress:
			return "127.0.0.1:0"
		default:
			return ""
		}
	})

	if config.ListenAddress != "127.0.0.1:0" {
		t.Fatalf("listen address = %q, want 127.0.0.1:0", config.ListenAddress)
	}
	if config.tlsEnabled() {
		t.Fatal("tls enabled without cert/key")
	}
}

func TestLoadServeConfigReadsTLSFiles(t *testing.T) {
	config := loadServeConfig(func(name string) string {
		switch name {
		case envListenAddress:
			return ":8443"
		case envTLSCertFile:
			return "/tls/tls.crt"
		case envTLSKeyFile:
			return "/tls/tls.key"
		default:
			return ""
		}
	})

	if config.ListenAddress != ":8443" {
		t.Fatalf("listen address = %q, want :8443", config.ListenAddress)
	}
	if !config.tlsEnabled() {
		t.Fatal("tls enabled = false, want true")
	}
	if err := config.validate(); err != nil {
		t.Fatal(err)
	}
}

func TestServeConfigRejectsPartialTLS(t *testing.T) {
	config := serveConfig{
		ListenAddress: ":8443",
		TLSCertFile:   "/tls/tls.crt",
	}

	if err := config.validate(); err == nil {
		t.Fatal("expected partial tls config error")
	}
}
