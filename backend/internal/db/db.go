package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool    *pgxpool.Pool
	syncMu  sync.Mutex
	syncing map[uint32]bool
}

func New(ctx context.Context, url string) (*DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("creating connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	return &DB{pool: pool, syncing: make(map[uint32]bool)}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) tryLockSync(accountID uint32) bool {
	db.syncMu.Lock()
	defer db.syncMu.Unlock()
	if db.syncing[accountID] {
		return false
	}
	db.syncing[accountID] = true
	return true
}

func (db *DB) unlockSync(accountID uint32) {
	db.syncMu.Lock()
	defer db.syncMu.Unlock()
	delete(db.syncing, accountID)
}

func (db *DB) IsSyncing(accountID uint32) bool {
	db.syncMu.Lock()
	defer db.syncMu.Unlock()
	return db.syncing[accountID]
}
