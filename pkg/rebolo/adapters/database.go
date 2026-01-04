package adapters

import (
	"context"
	"database/sql"
	"fmt"
	
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
	_ "github.com/lib/pq"
)

// BunDatabase implements Database interface
type BunDatabase struct {
	db *bun.DB
}

func NewBunDatabase() *BunDatabase {
	return &BunDatabase{}
}

func (d *BunDatabase) Connect(ctx context.Context) error {
	return nil // Will be implemented when DSN is provided
}

func (d *BunDatabase) ConnectWithDSN(dsn string, debug bool) error {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	d.db = bun.NewDB(sqldb, pgdialect.New())
	
	if debug {
		d.db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithVerbose(true)))
	}
	
	return d.db.Ping()
}

func (d *BunDatabase) Close() error {
	if d.db != nil {
		return d.db.Close()
	}
	return nil
}

func (d *BunDatabase) Migrate(ctx context.Context) error {
	// TODO: Implement migration logic
	return nil
}

func (d *BunDatabase) Health() error {
	if d.db != nil {
		return d.db.Ping()
	}
	return fmt.Errorf("database not connected")
}

func (d *BunDatabase) DB() *bun.DB {
	return d.db
}
