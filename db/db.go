package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"emarcey/data-vault/common"
	"emarcey/data-vault/common/tracer"
)

type DatabaseOpts struct {
	Driver          string
	Username        string
	Password        string
	Host            string
	DefaultDatabase string
}

type Database struct {
	db            *sql.DB
	logger        *logrus.Logger
	tracerCreator tracer.TracerCreator
}

func (db *Database) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	childTracer := db.tracerCreator(ctx, "execContext")
	defer childTracer.Close()

	childTracer.AddBreadcrumb(map[string]interface{}{"query": query, "args": args})

	result, err := db.db.ExecContext(ctx, query, args)
	if err != nil {
		db.logger.Errorf("Error in ExecContext: %v", err)
		childTracer.CaptureException(err)
		return nil, err
	}

	childTracer.AddBreadcrumb(map[string]interface{}{"result": result})
	return result, nil
}
func (db *Database) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	childTracer := db.tracerCreator(ctx, "queryContext")
	defer childTracer.Close()

	childTracer.AddBreadcrumb(map[string]interface{}{"query": query, "args": args})

	result, err := db.db.QueryContext(ctx, query, args)
	if err != nil {
		db.logger.Errorf("Error in QueryContext: %v", err)
		childTracer.CaptureException(err)
		return nil, err
	}

	childTracer.AddBreadcrumb(map[string]interface{}{"result": result})
	return result, nil
}

func (db *Database) Close() error {
	return db.db.Close()
}

func NewDatabase(logger *logrus.Logger, tracerCreator tracer.TracerCreator, opts DatabaseOpts) (*Database, error) {
	connStr := fmt.Sprintf("%s://%s:%s@%s/%s?sslmode=disable",
		opts.Driver,
		opts.Username,
		opts.Password,
		opts.Host,
		opts.DefaultDatabase,
	)
	db, err := sql.Open(opts.Driver, connStr)
	if err != nil {
		return nil, common.NewInitializationError("database", "Error during sql.Open: %v", err)
	}
	return &Database{
		db:            db,
		logger:        logger,
		tracerCreator: tracerCreator,
	}, nil
}
