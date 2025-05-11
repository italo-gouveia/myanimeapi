-- Drop junction tables first due to foreign key constraints
DROP TABLE IF EXISTS anime_tags;
DROP TABLE IF EXISTS anime_genres;

-- Drop main tables
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS genres; 