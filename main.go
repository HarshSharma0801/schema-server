package main

import (
	"fmt"
	"os"

	server "schema-server/http"
	"schema-server/internal/config"
	"schema-server/internal/contract"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	contractsPath := os.Getenv("CONTRACTS_PATH")
	if contractsPath == "" {
		contractsPath = "./contracts"
	}

	cfg := &config.Config{
		Contract: config.Contract{
			Path: contractsPath,
		},
	}

	// Create schema manager service
	contractService := contract.NewSchemaManagerService(logger, contractsPath)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "7080"
	}

	// Create and start the HTTP server
	httpServer := server.New(logger, contractService, cfg, port)

	logger.Info("Starting Schema Server with Schema Manager Service",
		zap.String("port", port),
		zap.String("contractsPath", contractsPath))

	if err := httpServer.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
