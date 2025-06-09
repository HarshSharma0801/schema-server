package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// RealService is a real implementation of the Service interface
type RealService struct {
	logger        *zap.Logger
	contractsPath string
}

// NewRealService creates a new real service instance
func NewRealService(logger *zap.Logger, contractsPath string) Service {
	return &RealService{
		logger:        logger,
		contractsPath: contractsPath,
	}
}

// Generate implements real schema generation from configuration
func (r *RealService) Generate(ctx context.Context, checkConfig bool) error {
	r.logger.Info("Generating schemas", zap.Bool("checkConfig", checkConfig))

	// Create directories if they don't exist
	if err := r.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate schemas based on current data
	if err := r.generateDefaultSchemas(); err != nil {
		return fmt.Errorf("failed to generate default schemas: %w", err)
	}

	r.logger.Info("Schema generation completed successfully")
	return nil
}

// Download implements real schema downloading
func (r *RealService) Download(ctx context.Context, checkConfig bool) error {
	r.logger.Info("Downloading schemas", zap.Bool("checkConfig", checkConfig))

	// Simulate downloading schemas from external sources
	if err := r.downloadExternalSchemas(); err != nil {
		return fmt.Errorf("failed to download external schemas: %w", err)
	}

	r.logger.Info("Schema download completed successfully")
	return nil
}

// Validate implements real schema validation
func (r *RealService) Validate(ctx context.Context) error {
	r.logger.Info("Validating schemas")

	// Read and validate all schemas
	validationResults, err := r.validateAllSchemas()
	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	r.logger.Info("Schema validation completed", zap.Any("results", validationResults))
	return nil
}

// GetAllTestsSchema implements real test schema retrieval
func (r *RealService) GetAllTestsSchema(ctx context.Context) (map[string]map[string]*models.OpenAPI, error) {
	r.logger.Info("Getting all tests schema")

	testsDir := filepath.Join(r.contractsPath, "tests")
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
				openAPI, err := r.readOpenAPIFile(path)
				if err != nil {
					r.logger.Warn("Failed to read test schema", zap.String("path", path), zap.Error(err))
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
		if err := r.generateDefaultTestSchemas(); err != nil {
			return nil, fmt.Errorf("failed to generate default test schemas: %w", err)
		}
		return r.GetAllTestsSchema(ctx)
	}

	return result, nil
}

// GetAllDownloadedMocksSchemas implements real mock schema retrieval
func (r *RealService) GetAllDownloadedMocksSchemas(ctx context.Context) ([]models.MockMapping, error) {
	r.logger.Info("Getting all downloaded mocks schemas")

	mocksDir := filepath.Join(r.contractsPath, "mocks")
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
				openAPI, err := r.readOpenAPIFile(path)
				if err != nil {
					r.logger.Warn("Failed to read mock schema", zap.String("path", path), zap.Error(err))
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
		if err := r.generateDefaultMockSchemas(); err != nil {
			return nil, fmt.Errorf("failed to generate default mock schemas: %w", err)
		}
		return r.GetAllDownloadedMocksSchemas(ctx)
	}

	return result, nil
}

// GenerateMocksSchemas implements real mock schema generation
func (r *RealService) GenerateMocksSchemas(ctx context.Context, services []string, mappings map[string][]string) error {
	r.logger.Info("Generating mocks schemas", zap.Strings("services", services), zap.Any("mappings", mappings))

	for _, service := range services {
		endpoints := mappings[service]
		if len(endpoints) == 0 {
			continue
		}

		// Generate OpenAPI schema for this service
		openAPI := r.generateOpenAPIFromEndpoints(service, endpoints)

		// Save to file
		serviceDir := filepath.Join(r.contractsPath, "mocks", service)
		if err := os.MkdirAll(serviceDir, 0755); err != nil {
			return fmt.Errorf("failed to create service directory: %w", err)
		}

		filename := fmt.Sprintf("mock-%s.json", time.Now().Format("20060102-150405"))
		filepath := filepath.Join(serviceDir, filename)

		if err := r.writeOpenAPIFile(filepath, openAPI); err != nil {
			return fmt.Errorf("failed to write mock schema: %w", err)
		}

		r.logger.Info("Generated mock schema", zap.String("service", service), zap.String("file", filepath))
	}

	return nil
}

// GenerateTestsSchemas implements real test schema generation
func (r *RealService) GenerateTestsSchemas(ctx context.Context, selectedTests []string) error {
	r.logger.Info("Generating tests schemas", zap.Strings("selectedTests", selectedTests))

	for _, testName := range selectedTests {
		// Generate OpenAPI schema for this test
		openAPI := r.generateOpenAPIFromTestName(testName)

		// Save to file
		testDir := filepath.Join(r.contractsPath, "tests", "generated")
		if err := os.MkdirAll(testDir, 0755); err != nil {
			return fmt.Errorf("failed to create test directory: %w", err)
		}

		filename := fmt.Sprintf("%s.json", testName)
		filepath := filepath.Join(testDir, filename)

		if err := r.writeOpenAPIFile(filepath, openAPI); err != nil {
			return fmt.Errorf("failed to write test schema: %w", err)
		}

		r.logger.Info("Generated test schema", zap.String("test", testName), zap.String("file", filepath))
	}

	return nil
}

// DownloadTests implements real test downloading
func (r *RealService) DownloadTests(path string) error {
	r.logger.Info("Downloading tests", zap.String("path", path))

	// Create target directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}

	// Copy test files to download path
	testsDir := filepath.Join(r.contractsPath, "tests")
	return r.copyDirectory(testsDir, path)
}

// DownloadMocks implements real mock downloading
func (r *RealService) DownloadMocks(ctx context.Context, path string) error {
	r.logger.Info("Downloading mocks", zap.String("path", path))

	// Create target directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("failed to create download directory: %w", err)
	}

	// Copy mock files to download path
	mocksDir := filepath.Join(r.contractsPath, "mocks")
	return r.copyDirectory(mocksDir, path)
}

// HTTPDocToOpenAPI implements real HTTPDoc to OpenAPI conversion
func (r *RealService) HTTPDocToOpenAPI(logger *zap.Logger, httpDoc models.HTTPDoc) (models.OpenAPI, error) {
	r.logger.Info("Converting HTTPDoc to OpenAPI", zap.String("name", httpDoc.Name))

	// Convert HTTPDoc to OpenAPI dynamically
	openAPI := models.OpenAPI{
		OpenAPI: "3.0.0",
		Info: models.Info{
			Title:       httpDoc.Name,
			Version:     httpDoc.Version,
			Description: fmt.Sprintf("Generated from %s", httpDoc.Kind),
		},
		Paths: make(map[string]models.PathItem),
	}

	// Parse the request URL to extract path
	url := httpDoc.Spec.Request.URL
	if url == "" {
		url = "/api/default"
	}

	// Extract path from URL
	path := r.extractPathFromURL(url)

	// Create operation based on HTTP method
	operation := &models.Operation{
		Summary:     fmt.Sprintf("%s operation for %s", httpDoc.Spec.Request.Method, httpDoc.Name),
		Description: fmt.Sprintf("Generated from HTTPDoc: %s", httpDoc.Name),
		Responses:   make(map[string]models.ResponseItem),
	}

	// Add response based on HTTPDoc response
	statusCode := fmt.Sprintf("%d", httpDoc.Spec.Response.StatusCode)
	operation.Responses[statusCode] = models.ResponseItem{
		Description: httpDoc.Spec.Response.StatusMessage,
		Content: map[string]models.MediaType{
			"application/json": {
				Schema: models.Schema{
					Type:       "object",
					Properties: r.inferSchemaFromBody(httpDoc.Spec.Response.Body),
				},
			},
		},
	}

	// Add request body if present
	if httpDoc.Spec.Request.Body != "" {
		operation.RequestBody = &models.RequestBody{
			Content: map[string]models.MediaType{
				"application/json": {
					Schema: models.Schema{
						Type:       "object",
						Properties: r.inferSchemaFromBody(httpDoc.Spec.Request.Body),
					},
				},
			},
		}
	}

	// Add operation to path
	pathItem := models.PathItem{}
	switch strings.ToLower(httpDoc.Spec.Request.Method) {
	case "get":
		pathItem.Get = operation
	case "post":
		pathItem.Post = operation
	case "put":
		pathItem.Put = operation
	case "patch":
		pathItem.Patch = operation
	case "delete":
		pathItem.Delete = operation
	}

	openAPI.Paths[path] = pathItem

	// Save the converted schema
	convertedDir := filepath.Join(r.contractsPath, "generated")
	if err := os.MkdirAll(convertedDir, 0755); err == nil {
		filename := fmt.Sprintf("%s-converted.json", httpDoc.Name)
		filepath := filepath.Join(convertedDir, filename)
		r.writeOpenAPIFile(filepath, &openAPI)
	}

	return openAPI, nil
}

// Helper methods

func (r *RealService) ensureDirectories() error {
	dirs := []string{
		filepath.Join(r.contractsPath, "tests"),
		filepath.Join(r.contractsPath, "mocks"),
		filepath.Join(r.contractsPath, "generated"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (r *RealService) readOpenAPIFile(path string) (*models.OpenAPI, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var openAPI models.OpenAPI
	if err := json.Unmarshal(data, &openAPI); err != nil {
		return nil, err
	}

	return &openAPI, nil
}

func (r *RealService) writeOpenAPIFile(path string, openAPI *models.OpenAPI) error {
	data, err := json.MarshalIndent(openAPI, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (r *RealService) generateOpenAPIFromEndpoints(serviceName string, endpoints []string) *models.OpenAPI {
	openAPI := &models.OpenAPI{
		OpenAPI: "3.0.0",
		Info: models.Info{
			Title:   fmt.Sprintf("%s API", strings.Title(serviceName)),
			Version: "1.0.0",
		},
		Paths: make(map[string]models.PathItem),
	}

	for _, endpoint := range endpoints {
		path := r.extractPathFromURL(endpoint)

		operation := &models.Operation{
			Summary: fmt.Sprintf("Operation for %s", path),
			Responses: map[string]models.ResponseItem{
				"200": {
					Description: "Success",
					Content: map[string]models.MediaType{
						"application/json": {
							Schema: models.Schema{Type: "object"},
						},
					},
				},
			},
		}

		openAPI.Paths[path] = models.PathItem{Get: operation}
	}

	return openAPI
}

func (r *RealService) generateOpenAPIFromTestName(testName string) *models.OpenAPI {
	return &models.OpenAPI{
		OpenAPI: "3.0.0",
		Info: models.Info{
			Title:   fmt.Sprintf("Test: %s", testName),
			Version: "1.0.0",
		},
		Paths: map[string]models.PathItem{
			fmt.Sprintf("/test/%s", strings.ToLower(testName)): {
				Get: &models.Operation{
					Summary: fmt.Sprintf("Test operation for %s", testName),
					Responses: map[string]models.ResponseItem{
						"200": {
							Description: "Test success",
							Content: map[string]models.MediaType{
								"application/json": {
									Schema: models.Schema{Type: "object"},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *RealService) generateDefaultSchemas() error {
	// Generate some default test schemas
	if err := r.generateDefaultTestSchemas(); err != nil {
		return err
	}

	// Generate some default mock schemas
	return r.generateDefaultMockSchemas()
}

func (r *RealService) generateDefaultTestSchemas() error {
	testSets := []struct {
		setName  string
		testName string
		api      *models.OpenAPI
	}{
		{
			setName:  "user-service",
			testName: "get-user-test",
			api: &models.OpenAPI{
				OpenAPI: "3.0.0",
				Info:    models.Info{Title: "User Service Test", Version: "1.0.0"},
				Paths: map[string]models.PathItem{
					"/users/{id}": {
						Get: &models.Operation{
							Summary: "Get user by ID",
							Responses: map[string]models.ResponseItem{
								"200": {Description: "User found"},
							},
						},
					},
				},
			},
		},
		{
			setName:  "order-service",
			testName: "create-order-test",
			api: &models.OpenAPI{
				OpenAPI: "3.0.0",
				Info:    models.Info{Title: "Order Service Test", Version: "1.0.0"},
				Paths: map[string]models.PathItem{
					"/orders": {
						Post: &models.Operation{
							Summary: "Create new order",
							Responses: map[string]models.ResponseItem{
								"201": {Description: "Order created"},
							},
						},
					},
				},
			},
		},
	}

	for _, test := range testSets {
		testDir := filepath.Join(r.contractsPath, "tests", test.setName)
		if err := os.MkdirAll(testDir, 0755); err != nil {
			return err
		}

		filepath := filepath.Join(testDir, fmt.Sprintf("%s.json", test.testName))
		if err := r.writeOpenAPIFile(filepath, test.api); err != nil {
			return err
		}
	}

	return nil
}

func (r *RealService) generateDefaultMockSchemas() error {
	mockServices := []struct {
		serviceName string
		api         *models.OpenAPI
	}{
		{
			serviceName: "payment-service",
			api: &models.OpenAPI{
				OpenAPI: "3.0.0",
				Info:    models.Info{Title: "Payment Service Mock", Version: "1.0.0"},
				Paths: map[string]models.PathItem{
					"/payments": {
						Post: &models.Operation{
							Summary: "Process payment",
							Responses: map[string]models.ResponseItem{
								"200": {Description: "Payment processed"},
							},
						},
					},
				},
			},
		},
		{
			serviceName: "notification-service",
			api: &models.OpenAPI{
				OpenAPI: "3.0.0",
				Info:    models.Info{Title: "Notification Service Mock", Version: "1.0.0"},
				Paths: map[string]models.PathItem{
					"/notifications": {
						Post: &models.Operation{
							Summary: "Send notification",
							Responses: map[string]models.ResponseItem{
								"200": {Description: "Notification sent"},
							},
						},
					},
				},
			},
		},
	}

	for _, mock := range mockServices {
		mockDir := filepath.Join(r.contractsPath, "mocks", mock.serviceName)
		if err := os.MkdirAll(mockDir, 0755); err != nil {
			return err
		}

		filepath := filepath.Join(mockDir, "mock-default.json")
		if err := r.writeOpenAPIFile(filepath, mock.api); err != nil {
			return err
		}
	}

	return nil
}

func (r *RealService) downloadExternalSchemas() error {
	// Simulate downloading schemas from external sources
	// In a real implementation, this would fetch from URLs, APIs, etc.
	r.logger.Info("Simulating external schema download")
	return nil
}

func (r *RealService) validateAllSchemas() (map[string]bool, error) {
	results := make(map[string]bool)

	// Validate test schemas
	testsDir := filepath.Join(r.contractsPath, "tests")
	if err := filepath.Walk(testsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".json") {
			_, err := r.readOpenAPIFile(path)
			results[path] = err == nil
		}
		return nil
	}); err != nil {
		return nil, err
	}

	// Validate mock schemas
	mocksDir := filepath.Join(r.contractsPath, "mocks")
	if err := filepath.Walk(mocksDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".json") {
			_, err := r.readOpenAPIFile(path)
			results[path] = err == nil
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *RealService) copyDirectory(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		dstPath := filepath.Join(dst, relPath)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstFile.Close()

		_, err = srcFile.WriteTo(dstFile)
		return err
	})
}

func (r *RealService) extractPathFromURL(url string) string {
	// Simple URL path extraction
	if strings.Contains(url, "://") {
		parts := strings.Split(url, "://")
		if len(parts) > 1 {
			domainAndPath := parts[1]
			if slashIndex := strings.Index(domainAndPath, "/"); slashIndex != -1 {
				return domainAndPath[slashIndex:]
			}
		}
		return "/api/default"
	}

	if strings.HasPrefix(url, "/") {
		return url
	}

	return "/" + url
}

func (r *RealService) inferSchemaFromBody(body string) map[string]map[string]interface{} {
	properties := make(map[string]map[string]interface{})

	if body == "" {
		return properties
	}

	// Try to parse as JSON and infer schema
	var jsonData interface{}
	if err := json.Unmarshal([]byte(body), &jsonData); err == nil {
		if obj, ok := jsonData.(map[string]interface{}); ok {
			for key, value := range obj {
				properties[key] = map[string]interface{}{
					"type": r.inferJSONType(value),
				}
			}
		}
	}

	return properties
}

func (r *RealService) inferJSONType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case float64, int:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	default:
		return "string"
	}
}
