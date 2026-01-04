package core

import (
	"context"
	"net/http"
)

// App represents the core application
type App struct {
	config     Config
	router     Router
	database   Database
	renderer   Renderer
	middleware []Middleware
}

// Config interface for configuration
type Config interface {
	GetPort() string
	GetHost() string
	GetDatabaseURL() string
	GetEnvironment() string
	IsHotReload() bool
}

// Router interface for HTTP routing
type Router interface {
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	Resource(path string, controller Controller)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Use(middleware Middleware)
}

// Database interface for data persistence
type Database interface {
	Connect(ctx context.Context) error
	Close() error
	Migrate(ctx context.Context) error
	Health() error
}

// Renderer interface for template and JSON rendering
type Renderer interface {
	RenderHTML(w http.ResponseWriter, template string, data interface{}) error
	RenderJSON(w http.ResponseWriter, data interface{}) error
	RenderError(w http.ResponseWriter, message string, status int) error
}

// Controller interface for HTTP controllers
type Controller interface {
	Index(w http.ResponseWriter, r *http.Request)
	Show(w http.ResponseWriter, r *http.Request)
	New(w http.ResponseWriter, r *http.Request)
	Create(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

// Middleware type for HTTP middleware
type Middleware func(http.Handler) http.Handler

// NewApp creates a new application instance
func NewApp(config Config, router Router, database Database, renderer Renderer) *App {
	return &App{
		config:   config,
		router:   router,
		database: database,
		renderer: renderer,
	}
}

// Start starts the application server
func (a *App) Start() error {
	// Connect to database if configured
	if a.config.GetDatabaseURL() != "" {
		if err := a.database.Connect(context.Background()); err != nil {
			return err
		}
	}
	
	// Apply middleware
	for _, mw := range a.middleware {
		a.router.Use(mw)
	}
	
	port := a.config.GetPort()
	if port == "" {
		port = "3000"
	}
	
	return http.ListenAndServe(":"+port, a.router)
}

// AddMiddleware adds middleware to the application
func (a *App) AddMiddleware(middleware Middleware) {
	a.middleware = append(a.middleware, middleware)
}

// Router returns the router instance
func (a *App) Router() Router {
	return a.router
}

// Database returns the database instance
func (a *App) Database() Database {
	return a.database
}

// Renderer returns the renderer instance
func (a *App) Renderer() Renderer {
	return a.renderer
}
