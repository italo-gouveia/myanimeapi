-- ============================================================
-- 000005 — Add cover_url and mal_id to animes
-- cover_url: direct CDN link (e.g. MyAnimeList / Jikan)
-- mal_id:    MyAnimeList numeric ID — used for ETL deduplication
-- ============================================================

ALTER TABLE animes
  ADD COLUMN IF NOT EXISTS cover_url TEXT,
  ADD COLUMN IF NOT EXISTS mal_id    INTEGER;

-- Partial unique index: only unique among rows where mal_id IS NOT NULL
-- (NULL values don't violate UNIQUE constraints in Postgres but this is explicit)
CREATE UNIQUE INDEX IF NOT EXISTS idx_animes_mal_id
  ON animes (mal_id)
  WHERE mal_id IS NOT NULL;
