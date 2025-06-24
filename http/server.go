package server

import (
	"encoding/json"
	"net/http"
	"strings"

	"schema-server/internal/contract"
	"schema-server/internal/matcher"
	"schema-server/internal/models"

	"go.uber.org/zap"
)

type Server struct {
	mux             *http.ServeMux
	logger          *zap.Logger
	contractService contract.SchemaService
	port            string
}

func New(logger *zap.Logger, contractService contract.SchemaService, port string) *Server {
	mux := http.NewServeMux()
	server := &Server{
		mux:             mux,
		logger:          logger,
		contractService: contractService,
		port:            port,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	s.mux.HandleFunc("GET /api/v1/health", s.healthCheck)
	s.mux.HandleFunc("POST /api/v1/schemas", s.uploadSchema)
	s.mux.HandleFunc("GET /api/v1/schemas", s.listSchemas)
	s.mux.HandleFunc("GET /api/v1/schemas/", s.schemaRouter)
	s.mux.HandleFunc("POST /api/v1/schemas/", s.compareSchema)
}

func (s *Server) Start() error {
	s.logger.Info("Starting schema HTTP server", zap.String("port", s.port))
	return http.ListenAndServe(":"+s.port, s.mux)
}

func (s *Server) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) writeError(w http.ResponseWriter, status int, message string, details string) {
	s.writeJSON(w, status, map[string]interface{}{
		"error":   message,
		"details": details,
	})
}

func (s *Server) decodeJSONRequest(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// Handlers
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "schema-server",
	})
}

func (s *Server) uploadSchema(w http.ResponseWriter, r *http.Request) {
	var schema models.OpenAPI
	if err := s.decodeJSONRequest(r, &schema); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid schema body", err.Error())
		return
	}
	if err := s.contractService.Upload(r.Context(), &schema); err != nil {
		s.writeError(w, http.StatusBadRequest, "Failed to upload schema", err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{"message": "Schema uploaded successfully", "title": schema.Info.Title})
}

func (s *Server) listSchemas(w http.ResponseWriter, r *http.Request) {
	titles, err := s.contractService.List(r.Context())
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Failed to list schemas", err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, map[string]interface{}{"schemas": titles})
}

func (s *Server) schemaRouter(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/schemas/")
	if path == "" {
		s.writeError(w, http.StatusBadRequest, "Missing schema title", "")
		return
	}
	parts := strings.SplitN(path, "/", 2)
	title := parts[0]
	if len(parts) == 2 && parts[1] == "download" {
		
		data, err := s.contractService.Download(r.Context(), title)
		if err != nil {
			s.writeError(w, http.StatusNotFound, "Schema not found", err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\""+title+".json\"")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
		return
	}
	// Fetch schema as JSON
	schema, err := s.contractService.Fetch(r.Context(), title)
	if err != nil {
		s.writeError(w, http.StatusNotFound, "Schema not found", err.Error())
		return
	}
	s.writeJSON(w, http.StatusOK, schema)
}

func (s *Server) compareSchema(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/schemas/")
	if !strings.HasSuffix(path, "/compare") {
		s.writeError(w, http.StatusNotFound, "Not found", "")
		return
	}
	title := strings.TrimSuffix(path, "/compare")
	if strings.HasSuffix(title, "/") {
		title = title[:len(title)-1]
	}
	if title == "" {
		s.writeError(w, http.StatusBadRequest, "Missing schema title", "")
		return
	}

	// Fetch stored schema as OpenAPI struct
	storedSchema, err := s.contractService.Fetch(r.Context(), title)
	if err != nil {
		s.writeError(w, http.StatusNotFound, "Schema not found", err.Error())
		return
	}

	// Read posted schema as OpenAPI struct
	var postedSchema models.OpenAPI
	if err := json.NewDecoder(r.Body).Decode(&postedSchema); err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid posted schema", err.Error())
		return
	}

	// Use matcher logic for semantic comparison
	score, match, err := matcher.Match(*storedSchema, postedSchema, "test-set", "mock-set", s.logger, matcher.CompareMode)
	if err != nil {
		s.writeError(w, http.StatusInternalServerError, "Schema comparison failed", err.Error())
		return
	}

	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"match": match,
		"score": score,
	})
}
