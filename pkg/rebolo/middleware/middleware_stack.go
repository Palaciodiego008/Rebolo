package middleware

import (
	"net/http"
	"path/filepath"
	"strings"
)

// MiddlewareFunc is a function that wraps an http.Handler
type MiddlewareFunc func(http.Handler) http.Handler

// MiddlewareConfig holds middleware configuration
type MiddlewareConfig struct {
	handler     MiddlewareFunc
	skipPaths   []string
	skipMethods []string
}

// MiddlewareStack manages a stack of middleware with skip patterns
type MiddlewareStack struct {
	middlewares []*MiddlewareConfig
}

// NewMiddlewareStack creates a new middleware stack
func NewMiddlewareStack() *MiddlewareStack {
	return &MiddlewareStack{
		middlewares: make([]*MiddlewareConfig, 0),
	}
}

// Use adds a middleware to the stack
func (ms *MiddlewareStack) Use(middleware MiddlewareFunc) *MiddlewareConfig {
	config := &MiddlewareConfig{
		handler:     middleware,
		skipPaths:   make([]string, 0),
		skipMethods: make([]string, 0),
	}
	ms.middlewares = append(ms.middlewares, config)
	return config
}

// Skip adds paths to skip for this middleware
func (mc *MiddlewareConfig) Skip(paths ...string) *MiddlewareConfig {
	mc.skipPaths = append(mc.skipPaths, paths...)
	return mc
}

// SkipMethod adds HTTP methods to skip for this middleware
func (mc *MiddlewareConfig) SkipMethod(methods ...string) *MiddlewareConfig {
	mc.skipMethods = append(mc.skipMethods, methods...)
	return mc
}

// shouldSkip checks if middleware should be skipped for this request
func (mc *MiddlewareConfig) shouldSkip(r *http.Request) bool {
	// Check path patterns
	for _, pattern := range mc.skipPaths {
		if matchPath(r.URL.Path, pattern) {
			return true
		}
	}
	
	// Check methods
	for _, method := range mc.skipMethods {
		if strings.EqualFold(r.Method, method) {
			return true
		}
	}
	
	return false
}

// matchPath checks if a path matches a pattern (supports wildcards)
func matchPath(path, pattern string) bool {
	// Exact match
	if path == pattern {
		return true
	}
	
	// Prefix match with wildcard
	if strings.HasSuffix(pattern, "/*") {
		prefix := strings.TrimSuffix(pattern, "/*")
		return strings.HasPrefix(path, prefix)
	}
	
	// Glob pattern match
	matched, _ := filepath.Match(pattern, path)
	return matched
}

// Apply applies all middleware in the stack to a handler
func (ms *MiddlewareStack) Apply(handler http.Handler) http.Handler {
	// Apply middleware in reverse order (last registered = outermost)
	for i := len(ms.middlewares) - 1; i >= 0; i-- {
		config := ms.middlewares[i]
		handler = ms.wrapWithSkip(config, handler)
	}
	return handler
}

// wrapWithSkip wraps a handler with middleware that can be skipped
func (ms *MiddlewareStack) wrapWithSkip(config *MiddlewareConfig, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.shouldSkip(r) {
			next.ServeHTTP(w, r)
			return
		}
		config.handler(next).ServeHTTP(w, r)
	})
}

// Group creates a middleware group for sub-routes
type MiddlewareGroup struct {
	stack       *MiddlewareStack
	middlewares []MiddlewareFunc
}

// NewMiddlewareGroup creates a new middleware group
func NewMiddlewareGroup(stack *MiddlewareStack) *MiddlewareGroup {
	return &MiddlewareGroup{
		stack:       stack,
		middlewares: make([]MiddlewareFunc, 0),
	}
}

// Use adds middleware to the group
func (mg *MiddlewareGroup) Use(middleware MiddlewareFunc) *MiddlewareGroup {
	mg.middlewares = append(mg.middlewares, middleware)
	return mg
}

// Apply applies group middleware to a handler
func (mg *MiddlewareGroup) Apply(handler http.Handler) http.Handler {
	// Apply group middleware
	for i := len(mg.middlewares) - 1; i >= 0; i-- {
		handler = mg.middlewares[i](handler)
	}
	
	// Apply global middleware
	return mg.stack.Apply(handler)
}

