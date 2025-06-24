package main

import (
	"fmt"
	"os"

	server "schema-server/http"
	"schema-server/internal/contract"

	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	contractsPath := os.Getenv("CONTRACTS_PATH")
	if contractsPath == "" {
		contractsPath = "."
	}

	service := contract.NewSchemaService(logger, contractsPath)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7080"
	}

	httpServer := server.New(logger, service, port)

	logger.Info("Starting Minimal Schema Server",
		zap.String("port", port),
		zap.String("contractsPath", contractsPath))

	if err := httpServer.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
