package handlers

import (
	"encoding/json"
	"net/http"
	"urlshortener/internal/dto/request"
	"urlshortener/internal/dto/response"
	"urlshortener/internal/dto/validation"
)

// URLHandler handles URL shortening operations
type URLHandler struct {
	// urlService service.URLService // Will add this when we create the service layer
}

// NewURLHandler creates a new URL handler
func NewURLHandler() *URLHandler {
	return &URLHandler{}
}

// CreateShortURL handles POST /api/shorten requests
func (h *URLHandler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errorResp := response.ErrorResponse{
			Error:   "Method not allowed",
			Message: "Only POST method is allowed",
			Code:    http.StatusMethodNotAllowed,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	var req request.CreateShortURLRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errorResp := response.ErrorResponse{
			Error:   "Invalid JSON payload",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// Validate request
	if err := validation.ValidateCreateShortURLRequest(&req); err != nil {
		errorResp := response.ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// TODO: Call service layer to create short URL
	// TODO: Return proper response

	// Temporary response
	successResp := response.SuccessResponse{
		Message: "URL shortening endpoint - validation passed",
		Data: map[string]string{
			"url":   req.URL,
			"title": req.Title,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(successResp)
}

// RedirectToLongURL handles GET /{shortCode} requests
func (h *URLHandler) RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
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

	// Extract short code from URL path
	shortCode := r.URL.Path[1:] // Remove leading slash

	if shortCode == "" {
		errorResp := response.ErrorResponse{
			Error:   "Short code required",
			Message: "Short code is required in the URL path",
			Code:    http.StatusBadRequest,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorResp)
		return
	}

	// TODO: Call service layer to get original URL
	// TODO: Increment click count
	// TODO: Check if URL is expired
	// TODO: Redirect to original URL

	// Temporary response
	successResp := response.SuccessResponse{
		Message: "Redirect endpoint - coming soon",
		Data: map[string]string{
			"short_code": shortCode,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResp)
}
