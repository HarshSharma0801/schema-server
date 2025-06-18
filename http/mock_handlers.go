package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

// Mock-related handlers

// getAllDownloadedMocksSchemas handles GET requests for retrieving all downloaded mock schemas
func (s *Server) getAllDownloadedMocksSchemas(w http.ResponseWriter, r *http.Request) {
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

// generateMocksSchemas handles POST requests for generating mock schemas
func (s *Server) generateMocksSchemas(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Services []string            `json:"services" bson:"services"`
		Mappings map[string][]string `json:"mappings" bson:"mappings"`
	}

	if err := s.decodeJSONRequest(r, &request); err != nil {
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

// downloadMocks handles POST requests for downloading mocks
func (s *Server) downloadMocks(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path string `json:"path" bson:"path"`
	}

	if err := s.decodeJSONRequest(r, &request); err != nil {
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
