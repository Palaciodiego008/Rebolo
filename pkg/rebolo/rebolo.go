package rebolo

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Palaciodiego008/rebololang/pkg/rebolo/adapters"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/core"
	"github.com/Palaciodiego008/rebololang/pkg/rebolo/ports"
)

// Application represents the main application facade
type Application struct {
	*core.App
	config   *ConfigAdapter
	router   *adapters.MuxRouter
	database adapters.DatabaseAdapter
	renderer *adapters.HTMLRenderer
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

	app := &Application{
		App:      coreApp,
		config:   config,
		router:   router,
		database: database,
		renderer: renderer,
	}

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

func (a *Application) Resource(path string, controller core.Controller) {
	a.router.Resource(path, controller)
}

// Convenience methods for rendering
func (a *Application) RenderHTML(w http.ResponseWriter, template string, data interface{}) error {
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
