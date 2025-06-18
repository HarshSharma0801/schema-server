package server

import (
	"context"
	"net/http"

	"schema-server/internal/models"

	"go.uber.org/zap"
)

// General schema handlers

// generateSchemas handles POST requests for general schema generation
func (s *Server) generateSchemas(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CheckConfig bool `json:"checkConfig" bson:"check_config"`
	}

	if err := s.decodeJSONRequest(r, &request); err != nil {
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

// downloadSchemas handles POST requests for downloading schemas
func (s *Server) downloadSchemas(w http.ResponseWriter, r *http.Request) {
	var request struct {
		CheckConfig bool `json:"checkConfig" bson:"check_config"`
	}

	if err := s.decodeJSONRequest(r, &request); err != nil {
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

// validateSchemas handles POST requests for validating schemas
func (s *Server) validateSchemas(w http.ResponseWriter, r *http.Request) {
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

// convertHTTPDocToOpenAPI handles POST requests for converting HTTPDoc to OpenAPI
func (s *Server) convertHTTPDocToOpenAPI(w http.ResponseWriter, r *http.Request) {
	var httpDoc models.HTTPDoc

	if err := s.decodeJSONRequest(r, &httpDoc); err != nil {
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
