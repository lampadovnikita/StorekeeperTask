package database

import (
	"context"
	"fmt"

	"github.com/lampadovnikita/StorekeeperTask/pkg/config"

	"github.com/jackc/pgx/v5/pgxpool"
)

const dsnFormat string = "postgresql://%s:%s@%s:%s/%s"

func NewPGXPool(ctx context.Context, pgcfg *config.PGConfig) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf(dsnFormat, pgcfg.Username, pgcfg.Password, pgcfg.Host, pgcfg.Port, pgcfg.Database)

	pgxCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	pool, err = pgxpool.NewWithConfig(ctx, pgxCfg)
	if err != nil {
		return nil, err
	}

	return pool, err
}
