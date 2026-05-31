-- Add role column with safe default
ALTER TABLE users ADD COLUMN IF NOT EXISTS role VARCHAR(50) NOT NULL DEFAULT 'user';
-- Backfill existing admins
UPDATE users SET role = 'admin' WHERE is_admin = true;
-- Index
CREATE INDEX IF NOT EXISTS idx_users_role ON users (role);
