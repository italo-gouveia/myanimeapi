-- ============================================================
-- 000001 — Rollback initial schema
-- Drop order respects foreign-key dependencies.
-- ============================================================

DROP TABLE IF EXISTS user_favorites;
DROP TABLE IF EXISTS password_reset_tokens;
DROP TABLE IF EXISTS media_attachments;
DROP TABLE IF EXISTS reviews;
DROP TABLE IF EXISTS animes;
DROP TABLE IF EXISTS users;
