package database

import (
	"context"
	"database/sql"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"github.com/emarcey/data-vault/common/logger"
	"github.com/emarcey/data-vault/common/tracer"
)

type initFunc func(dbMock *MockDatabase)

type MockDatabase struct {
	db            *sql.DB
	mock          sqlmock.Sqlmock
	logger        *logrus.Logger
	tracerCreator tracer.TracerCreator
}

func (m *MockDatabase) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return m.db.ExecContext(ctx, query, args...)
}

func (m *MockDatabase) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return m.db.QueryContext(ctx, query, args...)
}

func (m *MockDatabase) CreateTrace(ctx context.Context, operation string) tracer.Tracer {
	return m.tracerCreator(ctx, operation)
}

func (m *MockDatabase) GetLogger() *logrus.Logger {
	return m.logger
}

func NewMockDatabase() (*MockDatabase, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, err
	}
	logger, err := logger.MakeLogger("text", "local")
	if err != nil {
		return nil, err
	}
	return &MockDatabase{
		db:            db,
		mock:          mock,
		logger:        logger,
		tracerCreator: tracer.NewNoOpTracerMaker(),
	}, nil
}
