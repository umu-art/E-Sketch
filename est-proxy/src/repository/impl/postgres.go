package impl

import (
	"context"
	"est-proxy/src/config"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.elastic.co/apm/v2"
	"log"
	"net/url"
)

type PostgresServiceImpl struct {
	db *pgxpool.Pool
}

func NewPostgresServiceImpl() *PostgresServiceImpl {
	repositoryAddress := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		config.POSTGRES_USERNAME,
		url.QueryEscape(config.POSTGRES_PASSWORD),
		config.POSTGRES_HOST,
		config.POSTGRES_PORT,
		config.POSTGRES_DATABASE)

	pgxConfig, err := pgxpool.ParseConfig(repositoryAddress)
	if err != nil {
		log.Fatalf("Unable to parse config: %v", err)
	}

	pgxConfig.MaxConns = 20

	db, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	return &PostgresServiceImpl{db: db}
}

func (p *PostgresServiceImpl) Release() {
	if p.db != nil {
		p.db.Close()
		p.db = nil
	}
}

func (p *PostgresServiceImpl) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	span, _ := apm.StartSpan(ctx, "sql query", "service")
	span.Context.SetLabel("query", sql)
	defer span.End()

	resp, err := p.db.Exec(ctx, sql, arguments...)
	if err != nil {
		log.Printf("failed execute sql query: %v", err)
		return resp, err
	}

	return resp, nil
}

func (p *PostgresServiceImpl) QueryRow(ctx context.Context, sql string, arguments ...any) pgx.Row {
	span, _ := apm.StartSpan(ctx, "sql query", "service")
	span.Context.SetLabel("query", sql)
	defer span.End()

	return p.db.QueryRow(ctx, sql, arguments...)
}

func (p *PostgresServiceImpl) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	span, _ := apm.StartSpan(ctx, "sql query", "service")
	span.Context.SetLabel("query", sql)
	defer span.End()

	return p.db.Query(ctx, sql, arguments...)
}
