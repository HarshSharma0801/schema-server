package main

import (
	"log"
	"os"

	"schema-server/internal/config"
	"schema-server/internal/contract"
	"schema-server/server"

	"go.uber.org/zap"
)

func main() {
	// Initialize logger
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
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

	// Create real contract service (not mock)
	contractService := contract.NewRealService(logger, contractsPath)

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "7080"
	}

	// Create and start the HTTP server
	httpServer := server.New(logger, contractService, cfg, port)

	logger.Info("Starting Schema Server with Real Service",
		zap.String("port", port),
		zap.String("contractsPath", contractsPath))

	if err := httpServer.Start(); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
