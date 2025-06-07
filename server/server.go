package server

import (
	"context"
	"encoding/json"
	"net/http"

	"schema-server/internal/config"
	"schema-server/internal/contract"
	"schema-server/internal/models"

	"go.uber.org/zap"
)

// Server represents the HTTP server for schema operations
type Server struct {
	mux             *http.ServeMux
	logger          *zap.Logger
	contractService contract.Service
	config          *config.Config
	port            string
}

// New creates a new instance of the schema HTTP server
func New(logger *zap.Logger, contractService contract.Service, config *config.Config, port string) *Server {
	mux := http.NewServeMux()

	server := &Server{
		mux:             mux,
		logger:          logger,
		contractService: contractService,
		config:          config,
		port:            port,
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all the HTTP routes for schema operations
func (s *Server) setupRoutes() {
	// Fetching schema routes
	s.mux.HandleFunc("/api/v1/schemas/tests", s.getAllTestsSchema)
	s.mux.HandleFunc("/api/v1/schemas/mocks/downloaded", s.getAllDownloadedMocksSchemas)

	// Generating/Uploading schema routes
	s.mux.HandleFunc("/api/v1/schemas/mocks/generate", s.generateMocksSchemas)
	s.mux.HandleFunc("/api/v1/schemas/tests/generate", s.generateTestsSchemas)
	s.mux.HandleFunc("/api/v1/schemas/generate", s.generateSchemas)

	// Downloading schema routes
	s.mux.HandleFunc("/api/v1/schemas/download", s.downloadSchemas)
	s.mux.HandleFunc("/api/v1/schemas/tests/download", s.downloadTests)
	s.mux.HandleFunc("/api/v1/schemas/mocks/download", s.downloadMocks)

	// Validation route
	s.mux.HandleFunc("/api/v1/schemas/validate", s.validateSchemas)

	// Convert HTTPDoc to OpenAPI
	s.mux.HandleFunc("/api/v1/schemas/convert", s.convertHTTPDocToOpenAPI)

	// Health check endpoint
	s.mux.HandleFunc("/api/v1/health", s.healthCheck)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting schema HTTP server", zap.String("port", s.port))
	return http.ListenAndServe(":"+s.port, s.mux)
}

// writeJSON is a helper function to write JSON responses
func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// writeError is a helper function to write error responses
func (s *Server) writeError(w http.ResponseWriter, status int, message string, details string) {
	s.writeJSON(w, status, map[string]interface{}{
		"error":   message,
		"details": details,
	})
}

// Health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "schema-server",
	})
}

// Get all tests schema
func (s *Server) getAllTestsSchema(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := context.Background()
	testsMapping, err := s.contractService.GetAllTestsSchema(ctx)
	if err != nil {
		s.logger.Error("Failed to get all tests schema", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to retrieve tests schema", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    testsMapping,
		"message": "Tests schema retrieved successfully",
	})
}

// Get all downloaded mocks schemas
func (s *Server) getAllDownloadedMocksSchemas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	ctx := context.Background()
	mocksMapping, err := s.contractService.GetAllDownloadedMocksSchemas(ctx)
	if err != nil {
		s.logger.Error("Failed to get all downloaded mocks schemas", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to retrieve downloaded mocks schemas", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    mocksMapping,
		"message": "Downloaded mocks schemas retrieved successfully",
	})
}

// Generate mocks schemas
func (s *Server) generateMocksSchemas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var request struct {
		Services []string            `json:"services"`
		Mappings map[string][]string `json:"mappings"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	ctx := context.Background()
	err := s.contractService.GenerateMocksSchemas(ctx, request.Services, request.Mappings)
	if err != nil {
		s.logger.Error("Failed to generate mocks schemas", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to generate mocks schemas", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Mocks schemas generated successfully",
	})
}

// Generate tests schemas
func (s *Server) generateTestsSchemas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var request struct {
		SelectedTests []string `json:"selectedTests"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	ctx := context.Background()
	err := s.contractService.GenerateTestsSchemas(ctx, request.SelectedTests)
	if err != nil {
		s.logger.Error("Failed to generate tests schemas", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to generate tests schemas", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Tests schemas generated successfully",
	})
}

// Generate schemas (general)
func (s *Server) generateSchemas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var request struct {
		CheckConfig bool `json:"checkConfig"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := s.contractService.Generate(context.Background(), request.CheckConfig)
	if err != nil {
		s.logger.Error("Failed to generate schemas", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to generate schemas", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Schemas generated successfully",
	})
}

// Download schemas
func (s *Server) downloadSchemas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var request struct {
		CheckConfig bool `json:"checkConfig"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := s.contractService.Download(context.Background(), request.CheckConfig)
	if err != nil {
		s.logger.Error("Failed to download schemas", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to download schemas", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Schemas downloaded successfully",
	})
}

// Download tests
func (s *Server) downloadTests(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var request struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	err := s.contractService.DownloadTests(request.Path)
	if err != nil {
		s.logger.Error("Failed to download tests", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to download tests", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Tests downloaded successfully",
	})
}

// Download mocks
func (s *Server) downloadMocks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var request struct {
		Path string `json:"path"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	ctx := context.Background()
	err := s.contractService.DownloadMocks(ctx, request.Path)
	if err != nil {
		s.logger.Error("Failed to download mocks", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to download mocks", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Mocks downloaded successfully",
	})
}

// Validate schemas
func (s *Server) validateSchemas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	err := s.contractService.Validate(context.Background())
	if err != nil {
		s.logger.Error("Failed to validate schemas", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Schema validation failed", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Schemas validated successfully",
	})
}

// Convert HTTPDoc to OpenAPI
func (s *Server) convertHTTPDocToOpenAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		s.writeError(w, http.StatusMethodNotAllowed, "Method not allowed", "")
		return
	}

	var httpDoc models.HTTPDoc

	if err := json.NewDecoder(r.Body).Decode(&httpDoc); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid HTTPDoc format", err.Error())
		return
	}

	openAPI, err := s.contractService.HTTPDocToOpenAPI(s.logger, httpDoc)
	if err != nil {
		s.logger.Error("Failed to convert HTTPDoc to OpenAPI", zap.Error(err))
		s.writeError(w, http.StatusInternalServerError, "Failed to convert HTTPDoc to OpenAPI", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"data":    openAPI,
		"message": "HTTPDoc converted to OpenAPI successfully",
	})
}
