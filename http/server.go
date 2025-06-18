package server

import (
	"encoding/json"
	"net/http"

	"schema-server/internal/config"
	"schema-server/internal/contract"

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

// setupRoutes configures all the HTTP routes with specific methods
func (s *Server) setupRoutes() {
	// Health check endpoint
	s.mux.HandleFunc("GET /api/v1/health", s.healthCheck)

	// Schema fetching routes (GET)
	s.mux.HandleFunc("GET /api/v1/schemas/tests", s.getAllTestsSchema)
	s.mux.HandleFunc("GET /api/v1/schemas/mocks/downloaded", s.getAllDownloadedMocksSchemas)

	// Schema generation routes (POST)
	s.mux.HandleFunc("POST /api/v1/schemas/generate", s.generateSchemas)
	s.mux.HandleFunc("POST /api/v1/schemas/mocks/generate", s.generateMocksSchemas)
	s.mux.HandleFunc("POST /api/v1/schemas/tests/generate", s.generateTestsSchemas)

	// Schema download routes (POST)
	s.mux.HandleFunc("POST /api/v1/schemas/download", s.downloadSchemas)
	s.mux.HandleFunc("POST /api/v1/schemas/tests/download", s.downloadTests)
	s.mux.HandleFunc("POST /api/v1/schemas/mocks/download", s.downloadMocks)

	// Schema validation route (POST)
	s.mux.HandleFunc("POST /api/v1/schemas/validate", s.validateSchemas)

	// HTTPDoc to OpenAPI conversion route (POST)
	s.mux.HandleFunc("POST /api/v1/schemas/convert", s.convertHTTPDocToOpenAPI)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.logger.Info("Starting schema HTTP server", zap.String("port", s.port))
	return http.ListenAndServe(":"+s.port, s.mux)
}

// Utility functions

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

// decodeJSONRequest is a helper function to decode JSON request bodies
func (s *Server) decodeJSONRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}
