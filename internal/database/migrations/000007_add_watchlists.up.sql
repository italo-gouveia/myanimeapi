CREATE TABLE IF NOT EXISTS watchlists (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    anime_id   BIGINT       NOT NULL REFERENCES animes(id) ON DELETE CASCADE,
    status     VARCHAR(20)  NOT NULL DEFAULT 'plan_to_watch',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT idx_watchlist_user_anime UNIQUE (user_id, anime_id)
);
CREATE INDEX IF NOT EXISTS idx_watchlist_anime ON watchlists (anime_id);
CREATE INDEX IF NOT EXISTS idx_watchlist_status ON watchlists (status);
