package contract

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// Mock-related operations

// GetAllDownloadedMocksSchemas implements real mock schema retrieval
func (s *SchemaManagerService) GetAllDownloadedMocksSchemas(ctx context.Context) ([]models.MockMapping, error) {
	s.logger.Info("Getting all downloaded mocks schemas")

	mocksDir := filepath.Join(s.contractsPath, "mocks")
	var result []models.MockMapping

	// Read mock schemas from files
	if err := filepath.Walk(mocksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".json") && !info.IsDir() {
			// Parse file path to extract service and test set
			relPath, _ := filepath.Rel(mocksDir, path)
			parts := strings.Split(relPath, string(filepath.Separator))

			if len(parts) >= 2 {
				service := parts[0]
				testSetID := strings.TrimSuffix(parts[len(parts)-1], ".json")

				// Read and parse the mock file
				openAPI, err := s.readOpenAPIFile(path)
				if err != nil {
					s.logger.Warn("Failed to read mock schema", zap.String("path", path), zap.Error(err))
					return nil
				}

				mockMapping := models.MockMapping{
					Service:   service,
					TestSetID: testSetID,
					Mocks:     []*models.OpenAPI{openAPI},
				}
				result = append(result, mockMapping)
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to read mock schemas: %w", err)
	}

	// If no mocks found, generate some
	if len(result) == 0 {
		if err := s.generateDefaultMockSchemas(); err != nil {
			return nil, fmt.Errorf("failed to generate default mock schemas: %w", err)
		}
		return s.GetAllDownloadedMocksSchemas(ctx)
	}

	return result, nil
}

// GenerateMocksSchemas implements real mock schema generation
func (s *SchemaManagerService) GenerateMocksSchemas(ctx context.Context, services []string, mappings map[string][]string) error {
	s.logger.Info("Generating mocks schemas", zap.Strings("services", services), zap.Any("mappings", mappings))

	for _, service := range services {
		endpoints := mappings[service]
		if len(endpoints) == 0 {
			continue
		}

		// Generate OpenAPI schema for this service
		openAPI := s.generateOpenAPIFromEndpoints(service, endpoints)

		// Save to file
		serviceDir := filepath.Join(s.contractsPath, "mocks", service)
		if err := os.MkdirAll(serviceDir, 0755); err != nil {
			return fmt.Errorf("failed to create service directory: %w", err)
		}

		filename := fmt.Sprintf("mock-%s.json", time.Now().Format("20060102-150405"))
		filePath := filepath.Join(serviceDir, filename)

		if err := s.writeOpenAPIFile(filePath, openAPI); err != nil {
			return fmt.Errorf("failed to write mock schema: %w", err)
		}

		s.logger.Info("Generated mock schema", zap.String("service", service), zap.String("file", filePath))
	}

	return nil
}

// DownloadMocks implements mock downloading
func (s *SchemaManagerService) DownloadMocks(ctx context.Context, path string) error {
	s.logger.Info("Downloading mocks", zap.String("path", path))

	// Simulate downloading mocks by copying from generated to downloaded location
	downloadPath := filepath.Join(s.contractsPath, "mocks", "downloaded")
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}

	return s.copyDirectory(path, downloadPath)
}

// generateDefaultMockSchemas creates default mock schemas for demo purposes
func (s *SchemaManagerService) generateDefaultMockSchemas() error {
	services := []string{
		"user-service",
		"order-service",
		"payment-service",
		"inventory-service",
		"notification-service",
		"payment-gateway",
	}

	for _, service := range services {
		// Generate default endpoints for this service
		endpoints := s.getDefaultEndpointsForService(service)
		openAPI := s.generateOpenAPIFromEndpoints(service, endpoints)

		// Create service directory
		serviceDir := filepath.Join(s.contractsPath, "mocks", service)
		if err := os.MkdirAll(serviceDir, 0755); err != nil {
			return fmt.Errorf("failed to create service directory: %w", err)
		}

		// Save with timestamp
		filename := fmt.Sprintf("mock-%s.json", time.Now().Format("20060102-150405"))
		if service == "payment-service" || service == "notification-service" {
			filename = "mock-default.json"
		}
		filePath := filepath.Join(serviceDir, filename)

		if err := s.writeOpenAPIFile(filePath, openAPI); err != nil {
			return fmt.Errorf("failed to write mock schema %s: %w", service, err)
		}

		s.logger.Info("Generated default mock schema", zap.String("service", service))
	}

	return nil
}

// getDefaultEndpointsForService returns default endpoints for a given service
func (s *SchemaManagerService) getDefaultEndpointsForService(service string) []string {
	switch service {
	case "user-service":
		return []string{"/users", "/users/{id}", "/auth/login", "/auth/logout"}
	case "order-service":
		return []string{"/orders", "/orders/{id}", "/orders/{id}/status"}
	case "payment-service":
		return []string{"/payments", "/payments/{id}", "/payments/process"}
	case "inventory-service":
		return []string{"/inventory", "/inventory/{id}", "/inventory/check"}
	case "notification-service":
		return []string{"/notifications", "/notifications/send"}
	case "payment-gateway":
		return []string{"/gateway/charge", "/gateway/refund"}
	default:
		return []string{fmt.Sprintf("/%s", service), fmt.Sprintf("/%s/{id}", service)}
	}
}
