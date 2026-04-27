-- ============================================================
-- 000003 — favorites table (Favorite model, table = "favorites")
-- ============================================================

CREATE TABLE IF NOT EXISTS favorites (
    id         BIGSERIAL    PRIMARY KEY,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    user_id    BIGINT        NOT NULL REFERENCES users  (id) ON DELETE CASCADE,
    anime_id   BIGINT        NOT NULL REFERENCES animes (id) ON DELETE CASCADE,
    CONSTRAINT uq_favorites_user_anime UNIQUE (user_id, anime_id)
);

CREATE INDEX IF NOT EXISTS idx_favorites_anime ON favorites (anime_id);
