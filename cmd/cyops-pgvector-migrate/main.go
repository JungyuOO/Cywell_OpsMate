package main

import (
	"context"
	"fmt"
	"os"

	"github.com/JungyuOO/Cywell_OpsMate/internal/appserver"
)

func main() {
	config, err := appserver.PGVectorMigrationCommandConfigFromEnv(os.Getenv)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := appserver.RunPGVectorMigrationCommand(context.Background(), config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, "pgvector migration completed")
}
