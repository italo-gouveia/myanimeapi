-- ============================================================
-- 000001 — Initial schema
-- Creates: users, animes, reviews, media_attachments,
--          password_reset_tokens, user_favorites
-- ============================================================

-- ----------------------------------------------------------
-- users
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS users (
    id          BIGSERIAL PRIMARY KEY,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ,
    username    VARCHAR(255) NOT NULL,
    email       VARCHAR(255) NOT NULL,
    password    VARCHAR(255) NOT NULL,
    is_active   BOOLEAN      NOT NULL DEFAULT TRUE,
    profile_pic TEXT,
    bio         TEXT,
    social_links JSONB,
    is_admin    BOOLEAN      NOT NULL DEFAULT FALSE,
    CONSTRAINT uq_users_username UNIQUE (username),
    CONSTRAINT uq_users_email    UNIQUE (email)
);

CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users (deleted_at);

-- ----------------------------------------------------------
-- animes
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS animes (
    id          BIGSERIAL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    rating      DOUBLE PRECISION NOT NULL DEFAULT 0,
    episodes    INTEGER          NOT NULL DEFAULT 0,
    status      VARCHAR(50),
    start_date  TIMESTAMPTZ,
    end_date    TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_animes_title  ON animes (title);
CREATE INDEX IF NOT EXISTS idx_animes_status ON animes (status);

-- ----------------------------------------------------------
-- reviews
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS reviews (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    user_id    BIGINT NOT NULL REFERENCES users  (id),
    anime_id   BIGINT NOT NULL REFERENCES animes (id) ON DELETE CASCADE,
    content    TEXT   NOT NULL,
    rating     INTEGER NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_reviews_anime ON reviews (anime_id, user_id);
CREATE INDEX IF NOT EXISTS idx_reviews_user  ON reviews (user_id);

-- ----------------------------------------------------------
-- media_attachments
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS media_attachments (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    review_id  BIGINT        NOT NULL REFERENCES reviews (id) ON DELETE CASCADE,
    type       VARCHAR(50)   NOT NULL,
    url        TEXT          NOT NULL
);

-- ----------------------------------------------------------
-- password_reset_tokens
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id         BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    user_id    BIGINT        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    token      VARCHAR(255)  NOT NULL,
    expires_at TIMESTAMPTZ  NOT NULL,
    CONSTRAINT uq_password_reset_tokens_token UNIQUE (token)
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens (user_id);

-- ----------------------------------------------------------
-- user_favorites  (User.Favorites many2many junction)
-- ----------------------------------------------------------
CREATE TABLE IF NOT EXISTS user_favorites (
    user_id  BIGINT NOT NULL REFERENCES users  (id) ON DELETE CASCADE,
    anime_id BIGINT NOT NULL REFERENCES animes (id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, anime_id)
);
