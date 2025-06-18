package server

import (
	"context"
	"net/http"

	"go.uber.org/zap"
)

// Test-related handlers

// getAllTestsSchema handles GET requests for retrieving all test schemas
func (s *Server) getAllTestsSchema(w http.ResponseWriter, r *http.Request) {
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

// generateTestsSchemas handles POST requests for generating test schemas
func (s *Server) generateTestsSchemas(w http.ResponseWriter, r *http.Request) {
	var request struct {
		SelectedTests []string `json:"selectedTests" bson:"selected_tests"`
	}

	if err := s.decodeJSONRequest(r, &request); err != nil {
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

// downloadTests handles POST requests for downloading tests
func (s *Server) downloadTests(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Path string `json:"path" bson:"path"`
	}

	if err := s.decodeJSONRequest(r, &request); err != nil {
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
