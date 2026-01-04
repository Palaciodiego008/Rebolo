package rebolo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Palaciodiego008/rebololang/pkg/rebolo/adapters"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/core"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/errors"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/middleware"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/ports"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/session"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/validation"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/watcher"
)

// Application represents the main application facade
type Application struct {
	*core.App
	config          *ConfigAdapter
	router          *adapters.MuxRouter
	database        adapters.DatabaseAdapter
	renderer        *adapters.HTMLRenderer
	watcher         *watcher.FileWatcher
	sessionStore    *session.SessionStore       // Session management
	errorHandlers   errors.ErrorHandlers        // Custom error handlers
	middlewareStack *middleware.MiddlewareStack // Middleware stack with skip patterns
	mu              sync.RWMutex                // For thread-safe template reloading
	ctx             context.Context
	cancelFunc      context.CancelFunc
	lastChangeTime  time.Time // Track last file change for polling
}

// ConfigAdapter adapts ports.ConfigData to core.Config
type ConfigAdapter struct {
	data ports.ConfigData
}

func (c *ConfigAdapter) GetPort() string           { return c.data.Server.Port }
func (c *ConfigAdapter) GetHost() string           { return c.data.Server.Host }
func (c *ConfigAdapter) GetDatabaseDriver() string { return c.data.Database.Driver }
func (c *ConfigAdapter) GetDatabaseURL() string    { return c.data.Database.URL }
func (c *ConfigAdapter) GetDatabaseDebug() bool    { return c.data.Database.Debug }
func (c *ConfigAdapter) GetEnvironment() string    { return c.data.App.Env }
func (c *ConfigAdapter) IsHotReload() bool         { return c.data.Assets.HotReload }

// New creates a new ReboloLang application
func New() *Application {
	// Load configuration
	configPort := adapters.NewYAMLConfig()
	configData, err := configPort.Load()
	if err != nil {
		log.Printf("Failed to load config: %v", err)
	}

	config := &ConfigAdapter{data: configData}
	router := adapters.NewMuxRouter()
	renderer := adapters.NewHTMLRenderer()

	// Create database adapter based on driver from config
	var database adapters.DatabaseAdapter
	if config.GetDatabaseURL() != "" {
		driver := config.GetDatabaseDriver()
		if driver == "" {
			driver = "postgres" // Default to postgres for backward compatibility
			log.Printf("⚠️  No database driver specified, defaulting to 'postgres'")
		}

		factory := adapters.NewDatabaseFactory()
		database, err = factory.CreateDatabase(driver)
		if err != nil {
			log.Printf("❌ Failed to create database adapter: %v", err)
			database = adapters.NewBunDatabase() // Fallback to postgres
		} else {
			// Connect to database
			debug := config.GetDatabaseDebug() || config.GetEnvironment() == "development"
			if err := database.ConnectWithDSN(config.GetDatabaseURL(), debug); err != nil {
				log.Printf("❌ Database connection failed: %v", err)
			} else {
				log.Printf("✅ Database connected successfully (driver: %s)", driver)
			}
		}
	} else {
		// No database configured, use a default instance
		database = adapters.NewBunDatabase()
	}

	// Create core app
	coreApp := core.NewApp(config, router, database, renderer)

	// Add default middleware
	coreApp.AddMiddleware(LoggingMiddleware)
	coreApp.AddMiddleware(RecoveryMiddleware)

	ctx, cancel := context.WithCancel(context.Background())

	// Generate a random secret key for sessions in development
	// In production, this should come from environment variable
	secretKey := []byte("rebolo-secret-key-change-in-production")
	sessionStore := session.NewCookieSessionStore("rebolo_session", secretKey)

	app := &Application{
		App:             coreApp,
		config:          config,
		router:          router,
		database:        database,
		renderer:        renderer,
		sessionStore:    sessionStore,
		errorHandlers:   errors.NewErrorHandlers(),
		middlewareStack: middleware.NewMiddlewareStack(),
		ctx:             ctx,
		cancelFunc:      cancel,
	}

	// Set custom error handlers on router
	router.Router.NotFoundHandler = app.NotFoundHandler()
	router.Router.MethodNotAllowedHandler = app.MethodNotAllowedHandler()

	return app
}

// Start starts the application
func (a *Application) Start() error {
	port := a.config.GetPort()
	if port == "" {
		port = "3000"
	}

	fmt.Printf("🚀 ReboloLang server starting on port %s\n", port)
	return a.App.Start()
}

// Convenience methods for routing
func (a *Application) GET(path string, handler http.HandlerFunc) {
	a.router.GET(path, handler)
}

func (a *Application) POST(path string, handler http.HandlerFunc) {
	a.router.POST(path, handler)
}

func (a *Application) PUT(path string, handler http.HandlerFunc) {
	a.router.PUT(path, handler)
}

func (a *Application) DELETE(path string, handler http.HandlerFunc) {
	a.router.DELETE(path, handler)
}

// ServeStatic serves static files from a directory
func (a *Application) ServeStatic(prefix, dir string) {
	fs := http.FileServer(http.Dir(dir))
	a.router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, fs))
}

func (a *Application) Resource(path string, controller core.Controller) {
	a.router.Resource(path, controller)
}

// createRenderer creates a new HTML renderer (used for hot reload)
func (a *Application) createRenderer() *adapters.HTMLRenderer {
	return adapters.NewHTMLRenderer()
}

// EnableHotReload enables file watching and hot reload for development
func (a *Application) EnableHotReload() error {
	// Create file watcher
	fw := watcher.NewFileWatcher(a, []string{"views", "src", "public", "controllers"})

	// Start watching
	if err := fw.Start(); err != nil {
		return fmt.Errorf("failed to start file watcher: %v", err)
	}

	a.watcher = fw

	// Add hot reload middleware FIRST to inject script into HTML
	a.AddMiddleware(middleware.HotReloadMiddleware(true, "/__rebolo__/changes"))

	// Register polling endpoint for checking changes
	a.GET("/__rebolo__/changes", a.hotReloadChangesHandler)

	log.Printf("🔥 Hot reload enabled - watching files for changes")
	return nil
}

// hotReloadChangesHandler handles polling requests to check for file changes
func (a *Application) hotReloadChangesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Check if there are any pending changes
	// For simplicity, we'll use a timestamp-based approach
	a.mu.RLock()
	lastChange := a.lastChangeTime
	a.mu.RUnlock()

	response := map[string]interface{}{
		"changed": false,
		"time":    time.Now().Unix(),
	}

	// Check if there was a change in the last 2 seconds
	if time.Since(lastChange) < 2*time.Second {
		response["changed"] = true
		response["lastChange"] = lastChange.Unix()
	}

	JSON(w, response)
}

// GetSession retrieves the session for the current request
func (a *Application) GetSession(r *http.Request, w http.ResponseWriter) (*session.Session, error) {
	return a.sessionStore.Get(r, w)
}

// SetSessionStore allows custom session store configuration
func (a *Application) SetSessionStore(store *session.SessionStore) {
	a.sessionStore = store
}

// Shutdown gracefully shuts down the application
func (a *Application) Shutdown() {
	if a.watcher != nil {
		a.watcher.Close()
	}
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
}

// Convenience methods for rendering
func (a *Application) RenderHTML(w http.ResponseWriter, template string, data interface{}) error {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.renderer.RenderHTML(w, template, data)
}

func (a *Application) RenderJSON(w http.ResponseWriter, data interface{}) error {
	return a.renderer.RenderJSON(w, data)
}

func (a *Application) RenderError(w http.ResponseWriter, message string, status int) error {
	return a.renderer.RenderError(w, message, status)
}

// DB returns the underlying database/sql instance for convenience
func (a *Application) DB() *sql.DB {
	if a.database != nil {
		if db, ok := a.database.DB().(*sql.DB); ok {
			return db
		}
	}
	return nil
}

// Middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func RecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// Global convenience functions for backward compatibility
func Render(w http.ResponseWriter, template string, data interface{}) error {
	renderer := adapters.NewHTMLRenderer()
	return renderer.RenderHTML(w, template, data)
}

func JSON(w http.ResponseWriter, data interface{}) error {
	renderer := adapters.NewHTMLRenderer()
	return renderer.RenderJSON(w, data)
}

func JSONError(w http.ResponseWriter, message string, status int) error {
	renderer := adapters.NewHTMLRenderer()
	return renderer.RenderError(w, message, status)
}

// Use adds a middleware to the global stack
// Returns the MiddlewareConfig to allow chaining with Skip()
func (a *Application) Use(mw middleware.MiddlewareFunc) *middleware.MiddlewareConfig {
	return a.middlewareStack.Use(mw)
}

// Group creates a middleware group for specific routes
func (a *Application) Group(middlewares ...middleware.MiddlewareFunc) *middleware.MiddlewareGroup {
	group := middleware.NewMiddlewareGroup(a.middlewareStack)
	for _, mw := range middlewares {
		group.Use(mw)
	}
	return group
}

// UpdateLastChangeTime updates the last change time for hot reload
func (a *Application) UpdateLastChangeTime(t time.Time) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.lastChangeTime = t
}

// ReloadTemplates reloads HTML templates
func (a *Application) ReloadTemplates() {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.renderer = a.createRenderer()
}

// Bind binds request data to a struct
func (a *Application) Bind(r *http.Request, v interface{}) error {
	return validation.Bind(r, v)
}

// BindAndValidate binds and validates in one step
func (a *Application) BindAndValidate(r *http.Request, v interface{}) error {
	return validation.BindAndValidate(r, v)
}

// SetErrorHandler sets a custom error handler for a status code
func (a *Application) SetErrorHandler(code int, handler errors.ErrorHandler) {
	if a.errorHandlers == nil {
		a.errorHandlers = errors.NewErrorHandlers()
	}
	a.errorHandlers[code] = handler
}

// HandleError handles an error with the appropriate error handler
func (a *Application) HandleError(w http.ResponseWriter, r *http.Request, err error, code int) {
	if a.errorHandlers == nil {
		a.errorHandlers = errors.NewErrorHandlers()
	}

	// Try to render custom error page from views/errors/{code}.html
	templatePath := fmt.Sprintf("errors/%d.html", code)
	a.mu.RLock()
	renderer := a.renderer
	a.mu.RUnlock()

	if renderer != nil {
		renderErr := renderer.RenderHTML(w, templatePath, map[string]interface{}{
			"Code":  code,
			"Error": err,
			"Path":  r.URL.Path,
		})
		if renderErr == nil {
			return
		}
	}

	// Use custom handler if available
	if handler, ok := a.errorHandlers[code]; ok {
		handler(w, r, err, code)
		return
	}

	// Fallback to standard error
	http.Error(w, fmt.Sprintf("Error %d", code), code)
}

// NotFoundHandler is a custom 404 handler
func (a *Application) NotFoundHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.HandleError(w, r, fmt.Errorf("page not found: %s", r.URL.Path), 404)
	}
}

// MethodNotAllowedHandler is a custom 405 handler
func (a *Application) MethodNotAllowedHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.HandleError(w, r, fmt.Errorf("method not allowed: %s", r.Method), 405)
	}
}

// InternalErrorHandler handles 500 errors
func (a *Application) InternalErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("❌ Internal Server Error: %v", err)
	a.HandleError(w, r, err, 500)
}
