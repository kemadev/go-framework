package database

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kemadev/go-framework/pkg/config"
)

func NewClient(conf config.DatabaseConfig) (*pgxpool.Pool, error) {
	ctx := context.TODO()

	cfg, err := pgxpool.ParseConfig(conf.ConnectionURL.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing database config: %w", err)
	}

	cfg.ConnConfig.Tracer = otelpgx.NewTracer()

	pool, err := pgxpool.NewWithConfig(
		ctx,
		cfg,
	)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return pool, nil
}
