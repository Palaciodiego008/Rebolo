package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/mattn/go-sqlite3"
)

// SQLiteDatabase implements Database interface for SQLite
type SQLiteDatabase struct {
	db    *sql.DB
	debug bool
}

// NewSQLiteDatabase creates a new SQLite database adapter
func NewSQLiteDatabase() *SQLiteDatabase {
	return &SQLiteDatabase{}
}

// Connect connects to SQLite database
func (d *SQLiteDatabase) Connect(ctx context.Context) error {
	return nil // Will be implemented when DSN is provided
}

// ConnectWithDSN connects to SQLite with DSN (file path)
func (d *SQLiteDatabase) ConnectWithDSN(dsn string, debug bool) error {
	// Open SQLite database
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return fmt.Errorf("failed to open sqlite database: %w", err)
	}
	
	d.db = db
	d.debug = debug
	
	// Test connection
	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping sqlite database: %w", err)
	}
	
	if debug {
		log.Println("✅ SQLite database connected (debug mode enabled)")
	}
	
	return nil
}

// Close closes the database connection
func (d *SQLiteDatabase) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Migrate runs database migrations
func (d *SQLiteDatabase) Migrate(ctx context.Context) error {
	// TODO: Implement migration logic
	return nil
}

// Health checks database connection health
func (d *SQLiteDatabase) Health() error {
	if d.db != nil {
		return d.db.Ping()
	}
	return fmt.Errorf("database not connected")
}

// DB returns the underlying database/sql instance
func (d *SQLiteDatabase) DB() interface{} {
	return d.db
}
