package handlers

import (
	"encoding/json"
	"net/http"
	"urlshortener/internal/dto/response"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// ServeHTTP implements the http.Handler interface
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errorResp := response.ErrorResponse{
			Error:   "Method not allowed",
			Message: "Only GET method is allowed",
			Code:    http.StatusMethodNotAllowed,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	healthResp := response.HealthResponse{
		Status:  "ok",
		Message: "URL Shortener Service is running",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(healthResp)
}

// Health is a convenience function for direct handler registration
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.ServeHTTP(w, r)
}
