DROP INDEX IF EXISTS idx_animes_mal_id;

ALTER TABLE animes
  DROP COLUMN IF EXISTS cover_url,
  DROP COLUMN IF EXISTS mal_id;
