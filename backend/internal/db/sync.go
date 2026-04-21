package db

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/JustJoeYo/trophy-collector/internal/clients"
	"github.com/JustJoeYo/trophy-collector/internal/models"
)

func (db *DB) NeedsSync(ctx context.Context, accountID uint32) (bool, error) {
    var lastSynced *time.Time
    err := db.pool.QueryRow(ctx,
        "SELECT last_synced_at FROM players WHERE account_id = $1",
        accountID,
    ).Scan(&lastSynced)

    if err != nil {
        return true, nil
    }
    if lastSynced == nil {
        return true, nil
    }
    return time.Since(*lastSynced) > time.Hour, nil
}

func (db *DB) SyncPlayer(ctx context.Context, client clients.DeadlockClient, accountID uint32) error {
	if !db.tryLockSync(accountID) {
		slog.Info("sync already in progress, skipping", "account_id", accountID)
		return nil
	}
	defer db.unlockSync(accountID)
    slog.Info("syncing player", "account_id", accountID)

    var lastSynced *time.Time
    db.pool.QueryRow(ctx,
        "SELECT last_synced_at FROM players WHERE account_id = $1",
        accountID,
    ).Scan(&lastSynced)

    limit := 50
    total := 0
    var minMatchID *uint64

    for {
        matches, err := client.GetPlayerMatchesPage(ctx, accountID, minMatchID, limit, lastSynced)
        if err != nil {
            return fmt.Errorf("fetching page (min_match_id=%v): %w", minMatchID, err)
        }

        if len(matches) == 0 {
            break
        }

        inserted, err := db.insertMatches(ctx, accountID, matches)
        if err != nil {
            return fmt.Errorf("inserting matches: %w", err)
        }

        total += inserted
        lastID := matches[len(matches)-1].MatchID + 1
        minMatchID = &lastID
        slog.Info("synced matches", "account_id", accountID, "count", len(matches), "next_min_match_id", lastID)

        if len(matches) < limit {
            break
        }
    }

    _, err := db.pool.Exec(ctx, `
        INSERT INTO players (account_id, last_synced_at, total_matches)
        VALUES ($1, NOW(), $2)
        ON CONFLICT (account_id) DO UPDATE
        SET last_synced_at = NOW(),
            total_matches = players.total_matches + $2
    `, accountID, total)

    slog.Info("sync complete", "account_id", accountID, "new_matches", total)
    return err
}


func (db *DB) insertMatches(ctx context.Context, accountID uint32, matches []models.Match) (int, error) {
    tx, err := db.pool.Begin(ctx)
    if err != nil {
        return 0, err
    }
    defer tx.Rollback(ctx)

    inserted := 0
    for _, match := range matches {
        _, err := tx.Exec(ctx, `
            INSERT INTO matches (match_id, game_mode, match_mode, duration_s, start_time, winning_team, match_outcome)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            ON CONFLICT (match_id) DO NOTHING
        `, match.MatchID, match.GameMode, match.MatchMode, match.DurationS, match.StartTime, match.WinningTeam, match.MatchOutcome)
        if err != nil {
            return 0, err
        }

        for _, player := range match.Players {
            if player.AccountID != accountID {
                continue
            }
            tag, err := tx.Exec(ctx, `
                INSERT INTO player_matches
                (account_id, match_id, hero_id, team, kills, deaths, assists, net_worth, last_hits, denies, player_level, assigned_lane, abandon_match_time_s, won)
                VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
                ON CONFLICT (account_id, match_id) DO NOTHING
            `, accountID, match.MatchID, player.HeroID, player.Team, player.Kills, player.Deaths, player.Assists, player.NetWorth, player.LastHits, player.Denies, player.PlayerLevel, player.AssignedLane, player.AbandonMatchTimeS, player.Team == match.WinningTeam)
            if err != nil {
                return 0, err
            }
            inserted += int(tag.RowsAffected())
            break
        }
    }

    return inserted, tx.Commit(ctx)
}
