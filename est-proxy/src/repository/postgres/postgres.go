package postgres

import (
	"context"
	"est-proxy/src/repository/postgres/impl"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresService interface {
	Release()
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)
}

func NewPostgresService() PostgresService {
	return impl.NewPostgresServiceImpl()
}
