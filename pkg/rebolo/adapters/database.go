package adapters

import (
	"context"
	"fmt"
	"strings"
)

// DatabaseAdapter is a common interface for all database adapters
type DatabaseAdapter interface {
	Connect(ctx context.Context) error
	ConnectWithDSN(dsn string, debug bool) error
	Close() error
	Migrate(ctx context.Context) error
	Health() error
	DB() interface{} // Returns underlying database instance
}

// DatabaseFactory creates database adapters based on driver type
type DatabaseFactory struct{}

// NewDatabaseFactory creates a new database factory
func NewDatabaseFactory() *DatabaseFactory {
	return &DatabaseFactory{}
}

// CreateDatabase creates a database adapter based on the driver type
func (f *DatabaseFactory) CreateDatabase(driver string) (DatabaseAdapter, error) {
	driver = strings.ToLower(driver)

	switch driver {
	case "postgres", "postgresql":
		return NewPostgresDatabase(), nil
	case "sqlite", "sqlite3":
		return NewSQLiteDatabase(), nil
	case "mysql":
		return NewMySQLDatabase(), nil
	default:
		return nil, fmt.Errorf("unsupported database driver: %s (supported: postgres, sqlite, mysql)", driver)
	}
}

// BunDatabase is an alias for backward compatibility
// Deprecated: Use NewPostgresDatabase() instead
type BunDatabase = PostgresDatabase

// NewBunDatabase creates a new PostgreSQL database for backward compatibility
// Deprecated: Use NewPostgresDatabase() instead
func NewBunDatabase() *BunDatabase {
	return NewPostgresDatabase()
}
