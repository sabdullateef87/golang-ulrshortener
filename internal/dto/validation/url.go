package validation

import (
	"fmt"
	"net/url"
	"strings"
	"time"
	"urlshortener/internal/dto/request"
)

// ValidateCreateShortURLRequest validates the create short URL request
func ValidateCreateShortURLRequest(req *request.CreateShortURLRequest) error {
	// Validate URL
	if strings.TrimSpace(req.URL) == "" {
		return fmt.Errorf("URL is required")
	}

	// Parse and validate URL format
	parsedURL, err := url.Parse(req.URL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	// Check if URL has scheme
	if parsedURL.Scheme == "" {
		return fmt.Errorf("URL must include scheme (http:// or https://)")
	}

	// Only allow http and https
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("only HTTP and HTTPS URLs are allowed")
	}

	// Check if URL has host
	if parsedURL.Host == "" {
		return fmt.Errorf("URL must include a valid host")
	}

	// Validate expiration date if provided
	if req.ExpiresAt != "" {
		_, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return fmt.Errorf("invalid expires_at format, use ISO 8601 format (RFC3339)")
		}
	}

	// Validate title length
	if len(req.Title) > 255 {
		return fmt.Errorf("title cannot exceed 255 characters")
	}

	// Validate description length
	if len(req.Description) > 1000 {
		return fmt.Errorf("description cannot exceed 1000 characters")
	}

	return nil
}
