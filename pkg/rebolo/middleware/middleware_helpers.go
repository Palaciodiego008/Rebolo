package middleware

import (
	"net/http"
)

// Common middleware examples

// CORSMiddleware adds CORS headers
func CORSMiddleware(allowOrigin string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			
			next.ServeHTTP(w, r)
		})
	}
}

// AuthMiddleware checks if user is authenticated (example)
func AuthMiddleware(redirectTo string) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// This is a simple example - you'd implement your own auth logic
			// For now, we'll check if there's a session with "authenticated" = true
			
			// You'd need to get the session store from context or app
			// For simplicity, we'll skip this check for now
			// In a real implementation, you'd inject the app or session store
			
			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements simple rate limiting (placeholder)
func RateLimitMiddleware(requestsPerMinute int) MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implement rate limiting logic here
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}

// GzipMiddleware adds gzip compression (placeholder)
func GzipMiddleware() MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Implement gzip compression here
			// For now, just pass through
			next.ServeHTTP(w, r)
		})
	}
}
