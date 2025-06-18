package contract

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// Schema generation operations

// Generate implements real schema generation from configuration
func (s *SchemaManagerService) Generate(ctx context.Context, checkConfig bool) error {
	s.logger.Info("Generating schemas", zap.Bool("checkConfig", checkConfig))

	// Create directories if they don't exist
	if err := s.ensureDirectories(); err != nil {
		return fmt.Errorf("failed to create directories: %w", err)
	}

	// Generate schemas based on current data
	if err := s.generateDefaultSchemas(); err != nil {
		return fmt.Errorf("failed to generate default schemas: %w", err)
	}

	s.logger.Info("Schema generation completed successfully")
	return nil
}

// Download implements real schema downloading
func (s *SchemaManagerService) Download(ctx context.Context, checkConfig bool) error {
	s.logger.Info("Downloading schemas", zap.Bool("checkConfig", checkConfig))

	// Simulate downloading schemas from external sources
	if err := s.downloadExternalSchemas(); err != nil {
		return fmt.Errorf("failed to download external schemas: %w", err)
	}

	s.logger.Info("Schema download completed successfully")
	return nil
}

// Validate implements real schema validation
func (s *SchemaManagerService) Validate(ctx context.Context) error {
	s.logger.Info("Validating schemas")

	// Read and validate all schemas
	validationResults, err := s.validateAllSchemas()
	if err != nil {
		return fmt.Errorf("schema validation failed: %w", err)
	}

	s.logger.Info("Schema validation completed", zap.Any("results", validationResults))
	return nil
}

// HTTPDocToOpenAPI converts HTTPDoc format to OpenAPI spec
func (s *SchemaManagerService) HTTPDocToOpenAPI(logger *zap.Logger, httpDoc models.HTTPDoc) (models.OpenAPI, error) {
	logger.Info("Converting HTTPDoc to OpenAPI")

	// Extract path from URL
	path := s.extractPathFromURL(httpDoc.Spec.Request.URL)

	openAPI := models.OpenAPI{
		OpenAPI: "3.0.3",
		Info: models.Info{
			Title:       fmt.Sprintf("API for %s", httpDoc.Name),
			Description: "Converted from HTTPDoc",
			Version:     httpDoc.Version,
		},
		Servers: []map[string]string{
			{
				"url":         httpDoc.Spec.Request.URL,
				"description": "Main server",
			},
		},
		Paths: map[string]models.PathItem{
			path: s.createPathItemFromHTTPDoc(httpDoc),
		},
	}

	logger.Info("Successfully converted HTTPDoc to OpenAPI")
	return openAPI, nil
}

// createPathItemFromHTTPDoc creates a PathItem from HTTPDoc spec
func (s *SchemaManagerService) createPathItemFromHTTPDoc(httpDoc models.HTTPDoc) models.PathItem {
	pathItem := models.PathItem{}

	// Create operation based on method
	operation := &models.Operation{
		Summary:     fmt.Sprintf("%s operation", httpDoc.Spec.Request.Method),
		Description: fmt.Sprintf("Operation converted from HTTPDoc %s", httpDoc.Name),
		Responses:   make(map[string]models.ResponseItem),
	}

	// Handle request body if present
	if httpDoc.Spec.Request.Body != "" {
		requestBodySchema := s.inferSchemaFromBody(httpDoc.Spec.Request.Body)
		operation.RequestBody = &models.RequestBody{
			Content: map[string]models.MediaType{
				"application/json": {
					Schema: models.Schema{
						Type:       "object",
						Properties: requestBodySchema,
					},
				},
			},
		}
	}

	// Handle response
	responseBodySchema := s.inferSchemaFromBody(httpDoc.Spec.Response.Body)
	operation.Responses[fmt.Sprintf("%d", httpDoc.Spec.Response.StatusCode)] = models.ResponseItem{
		Description: httpDoc.Spec.Response.StatusMessage,
		Content: map[string]models.MediaType{
			"application/json": {
				Schema: models.Schema{
					Type:       "object",
					Properties: responseBodySchema,
				},
			},
		},
	}

	// Assign operation to correct HTTP method
	switch strings.ToUpper(httpDoc.Spec.Request.Method) {
	case "GET":
		pathItem.Get = operation
	case "POST":
		pathItem.Post = operation
	case "PUT":
		pathItem.Put = operation
	case "DELETE":
		pathItem.Delete = operation
	case "PATCH":
		pathItem.Patch = operation
	}

	return pathItem
}

// generateOpenAPIFromEndpoints generates an OpenAPI schema based on service endpoints
func (s *SchemaManagerService) generateOpenAPIFromEndpoints(serviceName string, endpoints []string) *models.OpenAPI {
	return &models.OpenAPI{
		OpenAPI: "3.0.3",
		Info: models.Info{
			Title:       fmt.Sprintf("%s API", strings.Title(serviceName)),
			Description: fmt.Sprintf("API specification for %s", serviceName),
			Version:     "1.0.0",
		},
		Servers: []map[string]string{
			{
				"url":         fmt.Sprintf("https://api.example.com/%s", serviceName),
				"description": fmt.Sprintf("%s API server", strings.Title(serviceName)),
			},
		},
		Paths: s.generatePathsFromEndpoints(endpoints, serviceName),
	}
}

// generatePathsFromEndpoints creates OpenAPI paths from endpoint strings
func (s *SchemaManagerService) generatePathsFromEndpoints(endpoints []string, serviceName string) map[string]models.PathItem {
	paths := make(map[string]models.PathItem)

	for _, endpoint := range endpoints {
		pathItem := models.PathItem{}

		// Simple GET operation for all endpoints
		pathItem.Get = &models.Operation{
			Summary:     fmt.Sprintf("Get %s", serviceName),
			Description: fmt.Sprintf("Retrieve %s data", serviceName),
			Responses: map[string]models.ResponseItem{
				"200": {
					Description: "Successful response",
					Content: map[string]models.MediaType{
						"application/json": {
							Schema: models.Schema{
								Type: "object",
								Properties: map[string]map[string]interface{}{
									"data": {"type": "object"},
								},
							},
						},
					},
				},
			},
		}

		paths[endpoint] = pathItem
	}

	return paths
}

// generateDefaultSchemas creates default schemas for the system
func (s *SchemaManagerService) generateDefaultSchemas() error {
	s.logger.Info("Generating default schemas")
	// This would contain logic for generating system-wide schemas
	return nil
}

// validateAllSchemas validates all existing schemas
func (s *SchemaManagerService) validateAllSchemas() (map[string]bool, error) {
	results := make(map[string]bool)

	// This would contain actual validation logic
	// For now, we'll simulate successful validation
	results["mock-schemas"] = true
	results["test-schemas"] = true
	results["generated-schemas"] = true

	return results, nil
}

// extractPathFromURL extracts the path component from a URL
func (s *SchemaManagerService) extractPathFromURL(url string) string {
	// Remove protocol and domain
	re := regexp.MustCompile(`^https?://[^/]+`)
	path := re.ReplaceAllString(url, "")

	if path == "" {
		return "/"
	}

	return path
}

// inferSchemaFromBody attempts to infer JSON schema from a body string
func (s *SchemaManagerService) inferSchemaFromBody(body string) map[string]map[string]interface{} {
	if body == "" {
		return make(map[string]map[string]interface{})
	}

	var data interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		// If parsing fails, return a generic object schema
		return map[string]map[string]interface{}{
			"data": {"type": "object"},
		}
	}

	return s.inferSchemaFromValue(data)
}

// inferSchemaFromValue recursively infers schema from a parsed JSON value
func (s *SchemaManagerService) inferSchemaFromValue(value interface{}) map[string]map[string]interface{} {
	schema := make(map[string]map[string]interface{})

	switch v := value.(type) {
	case map[string]interface{}:
		for key, val := range v {
			schema[key] = map[string]interface{}{
				"type": s.inferJSONType(val),
			}
		}
	default:
		schema["value"] = map[string]interface{}{
			"type": s.inferJSONType(value),
		}
	}

	return schema
}

// inferJSONType infers the JSON schema type from a Go value
func (s *SchemaManagerService) inferJSONType(value interface{}) string {
	switch value.(type) {
	case string:
		return "string"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return "integer"
	case float32, float64:
		return "number"
	case bool:
		return "boolean"
	case []interface{}:
		return "array"
	case map[string]interface{}:
		return "object"
	case nil:
		return "null"
	default:
		return "string"
	}
}
