package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/lib/pq"
)

// PostgresDatabase implements Database interface for PostgreSQL
type PostgresDatabase struct {
	db    *sql.DB
	debug bool
}

// NewPostgresDatabase creates a new PostgreSQL database adapter
func NewPostgresDatabase() *PostgresDatabase {
	return &PostgresDatabase{}
}

// Connect connects to PostgreSQL database
func (d *PostgresDatabase) Connect(ctx context.Context) error {
	return nil // Will be implemented when DSN is provided
}

// ConnectWithDSN connects to PostgreSQL with DSN
// DSN format: postgres://user:password@host:port/dbname?sslmode=disable
func (d *PostgresDatabase) ConnectWithDSN(dsn string, debug bool) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open postgres database: %w", err)
	}
	
	d.db = db
	d.debug = debug
	
	// Test connection
	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping postgres database: %w", err)
	}
	
	if debug {
		log.Println("✅ PostgreSQL database connected (debug mode enabled)")
	}
	
	return nil
}

// Close closes the database connection
func (d *PostgresDatabase) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Migrate runs database migrations
func (d *PostgresDatabase) Migrate(ctx context.Context) error {
	// TODO: Implement migration logic
	return nil
}

// Health checks database connection health
func (d *PostgresDatabase) Health() error {
	if d.db != nil {
		return d.db.Ping()
	}
	return fmt.Errorf("database not connected")
}

// DB returns the underlying database/sql instance
func (d *PostgresDatabase) DB() interface{} {
	return d.db
}
