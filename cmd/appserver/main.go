package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/JungyuOO/Cywell_OpsMate/internal/appserver"
)

const defaultListenAddress = ":8080"

func main() {
	listenAddress := strings.TrimSpace(os.Getenv("CYOPS_LISTEN_ADDRESS"))
	if listenAddress == "" {
		listenAddress = defaultListenAddress
	}

	server, err := appserver.NewServerFromConfig(context.Background(), appserver.LoadConfigFromEnv())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stdout, "cyops appserver listening on %s\n", listenAddress)
	if err := http.ListenAndServe(listenAddress, server); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
