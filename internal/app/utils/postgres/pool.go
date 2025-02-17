package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

const errLabel = "postgres"

// ConnenctionPool defines the set of methods required for interacting with a PostgreSQL connection pool.
// This interface abstracts common database operations like executing queries, transactions, and copying data.
type ConnenctionPool interface {
	Exec(ctx context.Context, query string, options ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, query string, options ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, query string, options ...interface{}) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error)
	Ping(ctx context.Context) error
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Close()
}
