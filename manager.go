package sqlmanager

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"golang.org/x/sync/semaphore"
)

// TODO - Currently using Background context with no cancellation. Should we implement this?

type SQLManager struct {
	db        *sql.DB
	semaphore *semaphore.Weighted
}

func (m *SQLManager) Exec(query string, args ...any) (sql.Result, error) {
	if err := m.semaphore.Acquire(context.Background(), 1); err != nil {
		return nil, err
	}
	defer m.semaphore.Release(1)

	return m.db.Exec(query, args...)
}

func (m *SQLManager) Query(query string, args ...any) (*sql.Rows, error) {
	if err := m.semaphore.Acquire(context.Background(), 1); err != nil {
		return nil, err
	}
	defer m.semaphore.Release(1)

	return m.db.Query(query, args...)
}

func (m *SQLManager) QueryRow(query string, args ...any) (*sql.Row, error) {
	if err := m.semaphore.Acquire(context.Background(), 1); err != nil {
		return nil, err
	}
	defer m.semaphore.Release(1)

	return m.db.QueryRow(query, args...), nil
}

func (m *SQLManager) Close() {
	// TODO - Wait for on-going operations to complete first?
	m.db.Close()
}

func NewSQLManager(driverName, dataSourceName string, maxConnections int64) (*SQLManager, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// * Postgres will throw "pq: sorry, too many clients already" when max_connections is REACHED,
	// * NOT when max_connections are IN USE.
	// * Meaning if max_connections=4, then only 3 connections may be active at a time.
	// * Removing 1 here ensures this will be handled.
	maxConnections = maxConnections - 1

	return &SQLManager{
		db:        db,
		semaphore: semaphore.NewWeighted(int64(maxConnections)),
	}, nil
}
