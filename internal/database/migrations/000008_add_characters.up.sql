-- Characters table
CREATE TABLE IF NOT EXISTS characters (
    id          BIGSERIAL    PRIMARY KEY,
    name        VARCHAR(150) NOT NULL,
    description TEXT,
    voice_actor VARCHAR(150),
    image_url   VARCHAR(500),
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    deleted_at  TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_characters_name       ON characters (name);
CREATE INDEX IF NOT EXISTS idx_characters_deleted_at ON characters (deleted_at);

-- Anime ↔ Character many-to-many junction
CREATE TABLE IF NOT EXISTS anime_characters (
    anime_id      BIGINT NOT NULL REFERENCES animes(id)      ON DELETE CASCADE,
    character_id  BIGINT NOT NULL REFERENCES characters(id)  ON DELETE CASCADE,
    PRIMARY KEY (anime_id, character_id)
);

CREATE INDEX IF NOT EXISTS idx_anime_characters_character_id ON anime_characters (character_id);
