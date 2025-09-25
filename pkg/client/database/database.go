package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kemadev/go-framework/pkg/config"
)

func NewClient(conf config.DatabaseConfig) (*pgxpool.Pool, error) {
	ctx := context.TODO()

	pool, err := pgxpool.New(
		ctx,
		conf.ConnectionURL.String(),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating database client: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error pinging database: %w", err)
	}

	return pool, nil
}
