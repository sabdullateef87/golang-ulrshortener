package response

// CreateShortURLResponse represents the response for creating short URLs
type CreateShortURLResponse struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
	ShortURL  string `json:"short_url"`
	LongURL   string `json:"long_url"`
	Title     string `json:"title,omitempty"`
	CreatedAt string `json:"created_at"`
	ExpiresAt string `json:"expires_at,omitempty"`
}

// GetShortURLResponse represents the response for getting URL details
type GetShortURLResponse struct {
	ID          string `json:"id"`
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	LongURL     string `json:"long_url"`
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	ClickCount  int64  `json:"click_count"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	ExpiresAt   string `json:"expires_at,omitempty"`
}

// URLAnalyticsResponse represents the response for URL analytics
type URLAnalyticsResponse struct {
	ShortCode    string `json:"short_code"`
	ClickCount   int64  `json:"click_count"`
	CreatedAt    string `json:"created_at"`
	LastAccessed string `json:"last_accessed,omitempty"`
}
