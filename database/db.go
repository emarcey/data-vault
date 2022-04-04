package database

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/emarcey/data-vault/common"
	"github.com/emarcey/data-vault/common/tracer"
)

type DatabaseOpts struct {
	Driver          string `yaml:"driver"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Host            string `yaml:"host"`
	DefaultDatabase string `yaml:"defaultDatabase"`
}

type Database interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	CreateTrace(ctx context.Context, operation string) tracer.Tracer
	GetLogger() *logrus.Logger
}

type DatabaseEngine struct {
	db            *sql.DB
	logger        *logrus.Logger
	tracerCreator tracer.TracerCreator
}

type DatabaseTransaction struct {
	tx            *sql.Tx
	logger        *logrus.Logger
	tracerCreator tracer.TracerCreator
}

func (db *DatabaseEngine) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	childTracer := db.tracerCreator(ctx, "execContext")
	defer childTracer.Close()

	childTracer.AddBreadcrumb(map[string]interface{}{"query": query, "args": args})

	result, err := db.db.ExecContext(ctx, query, args...)
	if err != nil {
		db.logger.Errorf("Error in ExecContext: %v", err)
		childTracer.CaptureException(err)
		return nil, err
	}

	return result, nil
}
func (db *DatabaseEngine) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	childTracer := db.tracerCreator(ctx, "queryContext")
	defer childTracer.Close()

	childTracer.AddBreadcrumb(map[string]interface{}{"query": query, "args": args})

	result, err := db.db.QueryContext(ctx, query, args...)
	if err != nil {
		db.logger.Errorf("Error in QueryContext: %v", err)
		childTracer.CaptureException(err)
		return nil, err
	}

	return result, nil
}

func (db *DatabaseEngine) CreateTrace(ctx context.Context, operation string) tracer.Tracer {
	return db.tracerCreator(ctx, operation)
}

func (db *DatabaseEngine) GetLogger() *logrus.Logger {
	return db.logger
}

func (db *DatabaseEngine) Close() error {
	return db.db.Close()
}

func (db *DatabaseEngine) StartTransaction(ctx context.Context) (*DatabaseTransaction, error) {
	tx, err := db.db.BeginTx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
	})
	if err != nil {
		return nil, common.NewDatabaseError(err, "StartTransaction", "")
	}
	return &DatabaseTransaction{
		tx:            tx,
		logger:        db.logger,
		tracerCreator: db.tracerCreator,
	}, nil
}

func (db *DatabaseTransaction) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	childTracer := db.tracerCreator(ctx, "execContext")
	defer childTracer.Close()

	childTracer.AddBreadcrumb(map[string]interface{}{"query": query, "args": args})

	result, err := db.tx.ExecContext(ctx, query, args...)
	if err != nil {
		db.logger.Errorf("Error in ExecContext: %v", err)
		childTracer.CaptureException(err)
		return nil, err
	}

	return result, nil
}
func (db *DatabaseTransaction) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	childTracer := db.tracerCreator(ctx, "queryContext")
	defer childTracer.Close()

	childTracer.AddBreadcrumb(map[string]interface{}{"query": query, "args": args})

	result, err := db.tx.QueryContext(ctx, query, args...)
	if err != nil {
		db.logger.Errorf("Error in QueryContext: %v", err)
		childTracer.CaptureException(err)
		return nil, err
	}

	return result, nil
}

func (db *DatabaseTransaction) CreateTrace(ctx context.Context, operation string) tracer.Tracer {
	return db.tracerCreator(ctx, operation)
}

func (db *DatabaseTransaction) GetLogger() *logrus.Logger {
	return db.logger
}

func (db *DatabaseTransaction) Commit() error {
	return db.tx.Commit()
}

func (db *DatabaseTransaction) Rollback() {
	err := db.tx.Rollback()
	if err != nil {
		db.logger.Errorf("Error in Rollback: %v", err)
	}

}

func NewDatabase(logger *logrus.Logger, tracerCreator tracer.TracerCreator, opts DatabaseOpts) (*DatabaseEngine, error) {
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
	return &DatabaseEngine{
		db:            db,
		logger:        logger,
		tracerCreator: tracerCreator,
	}, nil
}
