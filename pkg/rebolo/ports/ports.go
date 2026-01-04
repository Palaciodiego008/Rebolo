package ports

import (
	"context"
)

// ConfigPort defines configuration operations
type ConfigPort interface {
	Load() (ConfigData, error)
	GetEnv(key, defaultValue string) string
}

// DatabasePort defines database operations
type DatabasePort interface {
	Connect(dsn string) error
	Execute(ctx context.Context, query string, args ...interface{}) error
	Query(ctx context.Context, query string, args ...interface{}) (interface{}, error)
	Close() error
}

// TemplatePort defines template operations
type TemplatePort interface {
	Parse(pattern string) error
	Execute(name string, data interface{}) ([]byte, error)
}

// FileSystemPort defines file system operations
type FileSystemPort interface {
	ReadFile(path string) ([]byte, error)
	WriteFile(path string, data []byte) error
	Exists(path string) bool
	MkdirAll(path string) error
}

// ConfigData represents configuration data
type ConfigData struct {
	App struct {
		Name string `yaml:"name"`
		Env  string `yaml:"env"`
	} `yaml:"app"`
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
	Database struct {
		Driver string `yaml:"driver"` // postgres, sqlite, mysql
		URL    string `yaml:"url"`    // Connection string/DSN or file path for sqlite
		Debug  bool   `yaml:"debug"`  // Enable query logging
	} `yaml:"database"`
	Assets struct {
		HotReload bool `yaml:"hot_reload"`
	} `yaml:"assets"`
}
