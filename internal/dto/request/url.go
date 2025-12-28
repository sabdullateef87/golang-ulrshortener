package request

// CreateShortURLRequest represents the request payload for creating short URLs
type CreateShortURLRequest struct {
	URL         string `json:"url" validate:"required,url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"` // ISO 8601 format
}

// UpdateShortURLRequest represents the request payload for updating short URLs
type UpdateShortURLRequest struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ExpiresAt   string `json:"expires_at,omitempty"`
	IsActive    *bool  `json:"is_active,omitempty"` // Pointer to handle false values
}
