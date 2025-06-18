package contract

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// Test-related operations

// GetAllTestsSchema implements real test schema retrieval
func (s *SchemaManagerService) GetAllTestsSchema(ctx context.Context) (map[string]map[string]*models.OpenAPI, error) {
	s.logger.Info("Getting all tests schema")

	testsDir := filepath.Join(s.contractsPath, "tests")
	result := make(map[string]map[string]*models.OpenAPI)

	// Read test schemas from files
	if err := filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(path, ".json") && !info.IsDir() {
			// Parse file path to extract test set and test case
			relPath, _ := filepath.Rel(testsDir, path)
			parts := strings.Split(relPath, string(filepath.Separator))

			if len(parts) >= 2 {
				testSet := parts[0]
				testCase := strings.TrimSuffix(parts[len(parts)-1], ".json")

				// Read and parse the schema file
				openAPI, err := s.readOpenAPIFile(path)
				if err != nil {
					s.logger.Warn("Failed to read test schema", zap.String("path", path), zap.Error(err))
					return nil
				}

				if result[testSet] == nil {
					result[testSet] = make(map[string]*models.OpenAPI)
				}
				result[testSet][testCase] = openAPI
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("failed to read test schemas: %w", err)
	}

	// If no schemas found, generate some
	if len(result) == 0 {
		if err := s.generateDefaultTestSchemas(); err != nil {
			return nil, fmt.Errorf("failed to generate default test schemas: %w", err)
		}
		return s.GetAllTestsSchema(ctx)
	}

	return result, nil
}

// GenerateTestsSchemas implements real test schema generation
func (s *SchemaManagerService) GenerateTestsSchemas(ctx context.Context, selectedTests []string) error {
	s.logger.Info("Generating tests schemas", zap.Strings("selectedTests", selectedTests))

	for _, testName := range selectedTests {
		// Generate OpenAPI schema for this test
		openAPI := s.generateOpenAPIFromTestName(testName)

		// Create test directory
		testDir := filepath.Join(s.contractsPath, "tests", "generated")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			return fmt.Errorf("failed to create test directory: %w", err)
		}

		// Save to file
		filename := fmt.Sprintf("%s-test.json", testName)
		filePath := filepath.Join(testDir, filename)

		if err := s.writeOpenAPIFile(filePath, openAPI); err != nil {
			return fmt.Errorf("failed to write test schema: %w", err)
		}

		s.logger.Info("Generated test schema", zap.String("test", testName), zap.String("file", filePath))
	}

	return nil
}

// DownloadTests implements test downloading
func (s *SchemaManagerService) DownloadTests(path string) error {
	s.logger.Info("Downloading tests", zap.String("path", path))

	// Simulate downloading tests by copying from generated to downloaded location
	downloadPath := filepath.Join(s.contractsPath, "tests", "downloaded")
	if err := os.MkdirAll(downloadPath, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}

	return s.copyDirectory(path, downloadPath)
}

// generateDefaultTestSchemas creates default test schemas for demo purposes
func (s *SchemaManagerService) generateDefaultTestSchemas() error {
	testCases := []string{
		"user-authentication-test",
		"order-fulfillment-test",
		"payment-processing-test",
		"inventory-update-test",
	}

	testsDir := filepath.Join(s.contractsPath, "tests", "generated")
	if err := os.MkdirAll(testsDir, 0755); err != nil {
		return fmt.Errorf("failed to create tests directory: %w", err)
	}

	for _, testName := range testCases {
		openAPI := s.generateOpenAPIFromTestName(testName)
		filename := fmt.Sprintf("%s.json", testName)
		filePath := filepath.Join(testsDir, filename)

		if err := s.writeOpenAPIFile(filePath, openAPI); err != nil {
			return fmt.Errorf("failed to write test schema %s: %w", testName, err)
		}

		s.logger.Info("Generated default test schema", zap.String("test", testName))
	}

	return nil
}

// generateOpenAPIFromTestName generates an OpenAPI schema based on a test name
func (s *SchemaManagerService) generateOpenAPIFromTestName(testName string) *models.OpenAPI {
	// Parse test name to infer structure
	parts := strings.Split(testName, "-")
	serviceName := strings.Join(parts[:len(parts)-1], "-") // Everything except "test"

	return &models.OpenAPI{
		OpenAPI: "3.0.3",
		Info: models.Info{
			Title:       fmt.Sprintf("%s Test API", strings.Title(serviceName)),
			Description: fmt.Sprintf("Test schema for %s operations", serviceName),
			Version:     "1.0.0",
		},
		Servers: []map[string]string{
			{
				"url":         fmt.Sprintf("https://api.example.com/%s", serviceName),
				"description": fmt.Sprintf("%s test server", strings.Title(serviceName)),
			},
		},
		Paths: map[string]models.PathItem{
			fmt.Sprintf("/%s/test", serviceName): {
				Post: &models.Operation{
					Summary:     fmt.Sprintf("Execute %s test", serviceName),
					Description: fmt.Sprintf("Run test case for %s", serviceName),
					RequestBody: &models.RequestBody{
						Content: map[string]models.MediaType{
							"application/json": {
								Schema: models.Schema{
									Type: "object",
									Properties: map[string]map[string]interface{}{
										"testId": {
											"type":        "string",
											"description": "Unique test identifier",
										},
										"data": {
											"type":        "object",
											"description": "Test input data",
										},
									},
								},
							},
						},
					},
					Responses: map[string]models.ResponseItem{
						"200": {
							Description: "Test executed successfully",
							Content: map[string]models.MediaType{
								"application/json": {
									Schema: models.Schema{
										Type: "object",
										Properties: map[string]map[string]interface{}{
											"success": {
												"type":        "boolean",
												"description": "Test success status",
											},
											"result": {
												"type":        "object",
												"description": "Test result data",
											},
											"timestamp": {
												"type":        "string",
												"format":      "date-time",
												"description": "Test execution timestamp",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
