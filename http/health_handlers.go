package server

import (
	"net/http"
)

// Health check handlers

// healthCheck handles GET requests for health status
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	s.writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "healthy",
		"service": "schema-server",
	})
}
