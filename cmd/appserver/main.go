package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/JungyuOO/Cywell_OpsMate/internal/appserver"
)

const (
	defaultListenAddress = ":8080"
	envListenAddress     = "CYOPS_LISTEN_ADDRESS"
	envTLSCertFile       = "TLS_CERT_FILE"
	envTLSKeyFile        = "TLS_KEY_FILE"
)

type serveConfig struct {
	ListenAddress string
	TLSCertFile   string
	TLSKeyFile    string
}

func main() {
	config := loadServeConfig(os.Getenv)

	server, err := appserver.NewServerFromConfig(context.Background(), appserver.LoadConfigFromEnv())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	protocol := "http"
	if config.tlsEnabled() {
		protocol = "https"
	}
	fmt.Fprintf(os.Stdout, "cyops appserver listening on %s://%s\n", protocol, config.ListenAddress)
	if err := listenAndServe(config, server); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func loadServeConfig(getenv func(string) string) serveConfig {
	listenAddress := strings.TrimSpace(getenv(envListenAddress))
	if listenAddress == "" {
		listenAddress = defaultListenAddress
	}
	return serveConfig{
		ListenAddress: listenAddress,
		TLSCertFile:   strings.TrimSpace(getenv(envTLSCertFile)),
		TLSKeyFile:    strings.TrimSpace(getenv(envTLSKeyFile)),
	}
}

func listenAndServe(config serveConfig, handler http.Handler) error {
	if err := config.validate(); err != nil {
		return err
	}
	if config.tlsEnabled() {
		return http.ListenAndServeTLS(config.ListenAddress, config.TLSCertFile, config.TLSKeyFile, handler)
	}
	return http.ListenAndServe(config.ListenAddress, handler)
}

func (config serveConfig) tlsEnabled() bool {
	return config.TLSCertFile != "" && config.TLSKeyFile != ""
}

func (config serveConfig) validate() error {
	if (config.TLSCertFile == "") != (config.TLSKeyFile == "") {
		return fmt.Errorf("%s and %s must be set together", envTLSCertFile, envTLSKeyFile)
	}
	return nil
}
