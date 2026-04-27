-- ============================================================
-- 000004 — user_genres junction table (User.Genres many2many)
-- Depends on: 000001 (users), 000002 (genres)
-- ============================================================

CREATE TABLE IF NOT EXISTS user_genres (
    user_id  BIGINT NOT NULL REFERENCES users  (id) ON DELETE CASCADE,
    genre_id BIGINT NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, genre_id)
);
