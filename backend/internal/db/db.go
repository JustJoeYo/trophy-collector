package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
    pool *pgxpool.Pool
}

func New(ctx context.Context, url string) (*DB, error) {
    pool, err := pgxpool.New(ctx, url)
    if err != nil {
        return nil, fmt.Errorf("creating connection pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("pinging database: %w", err)
    }

    return &DB{pool: pool}, nil
}

func (db *DB) Close() {
    db.pool.Close()
}
