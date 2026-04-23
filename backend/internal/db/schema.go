package db

import "context"

const schema = `
CREATE TABLE IF NOT EXISTS players (
    account_id     BIGINT PRIMARY KEY,
    last_synced_at TIMESTAMPTZ,
    total_matches  INT DEFAULT 0,
    steam_name     TEXT,
    avatar_url     TEXT,
    created_at     TIMESTAMPTZ DEFAULT NOW()
);

ALTER TABLE players ADD COLUMN IF NOT EXISTS steam_name TEXT;
ALTER TABLE players ADD COLUMN IF NOT EXISTS avatar_url TEXT;
CREATE INDEX IF NOT EXISTS idx_players_steam_name ON players(steam_name);

CREATE TABLE IF NOT EXISTS matches (
    match_id      BIGINT PRIMARY KEY,
    game_mode     TEXT,
    match_mode    TEXT,
    duration_s    INT,
    start_time    TIMESTAMPTZ,
    winning_team  TEXT,
    match_outcome TEXT,
    created_at    TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS player_matches (
    account_id         BIGINT,
    match_id           BIGINT REFERENCES matches(match_id),
    hero_id            INT,
    team               TEXT,
    kills              INT,
    deaths             INT,
    assists            INT,
    net_worth          INT,
    last_hits          INT,
    denies             INT,
    player_level       INT,
    assigned_lane      INT,
    abandon_match_time_s INT,
    won                BOOLEAN,
    PRIMARY KEY (account_id, match_id)
);

CREATE INDEX IF NOT EXISTS idx_player_matches_account ON player_matches(account_id);
CREATE INDEX IF NOT EXISTS idx_matches_start_time ON matches(start_time);
`

func (db *DB) Migrate(ctx context.Context) error {
    _, err := db.pool.Exec(ctx, schema)
    return err
}
