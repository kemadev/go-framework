package database

import (
	"context"
	"fmt"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kemadev/go-framework/pkg/config"
	"github.com/kemadev/go-framework/pkg/monitoring"
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

	return pool, nil
}

func Check(p *pgxpool.Pool) monitoring.StatusCheck {
	err := p.Ping(context.Background())
	if err != nil {
		return monitoring.StatusCheck{
			Status:  monitoring.StatusDown,
			Message: "ping failed",
		}
	}

	return monitoring.StatusCheck{
		Status:  monitoring.StatusOK,
		Message: monitoring.StatusOK.String(),
	}
}
