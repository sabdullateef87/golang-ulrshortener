package routes

import (
	"net/http"
	"urlshortener/internal/http/handlers"
	"urlshortener/internal/http/middleware"
)

// SetupRoutes configures all application routes
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler()
	urlHandler := handlers.NewURLHandler()

	// Health check routes
	mux.Handle("/health", middleware.Logger(healthHandler))
	mux.Handle("/ping", middleware.Logger(healthHandler))

	// API routes
	apiMux := http.NewServeMux()
	apiMux.HandleFunc("/shorten", urlHandler.CreateShortURL)

	// Mount API routes under /api prefix
	mux.Handle("/api/", http.StripPrefix("/api", middleware.Logger(apiMux)))

	// Short URL redirect routes (catch-all for short codes)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			// Root path - show service info
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"service":"URL Shortener","status":"running"}`))
			return
		}
		// Assume it's a short code redirect
		urlHandler.RedirectToLongURL(w, r)
	})

	// Apply global middleware
	handler := middleware.CORS(mux)
	handler = middleware.Logger(handler)

	return handler
}
