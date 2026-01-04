package adapters

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	
	_ "github.com/go-sql-driver/mysql"
)

// MySQLDatabase implements Database interface for MySQL
type MySQLDatabase struct {
	db    *sql.DB
	debug bool
}

// NewMySQLDatabase creates a new MySQL database adapter
func NewMySQLDatabase() *MySQLDatabase {
	return &MySQLDatabase{}
}

// Connect connects to MySQL database
func (d *MySQLDatabase) Connect(ctx context.Context) error {
	return nil // Will be implemented when DSN is provided
}

// ConnectWithDSN connects to MySQL with DSN
// DSN format: user:password@tcp(host:port)/dbname?parseTime=true
func (d *MySQLDatabase) ConnectWithDSN(dsn string, debug bool) error {
	// Open MySQL database
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open mysql database: %w", err)
	}
	
	d.db = db
	d.debug = debug
	
	// Test connection
	if err := d.db.Ping(); err != nil {
		return fmt.Errorf("failed to ping mysql database: %w", err)
	}
	
	if debug {
		log.Println("✅ MySQL database connected (debug mode enabled)")
	}
	
	return nil
}

// Close closes the database connection
func (d *MySQLDatabase) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

// Migrate runs database migrations
func (d *MySQLDatabase) Migrate(ctx context.Context) error {
	// TODO: Implement migration logic
	return nil
}

// Health checks database connection health
func (d *MySQLDatabase) Health() error {
	if d.db != nil {
		return d.db.Ping()
	}
	return fmt.Errorf("database not connected")
}

// DB returns the underlying database/sql instance
func (d *MySQLDatabase) DB() interface{} {
	return d.db
}
