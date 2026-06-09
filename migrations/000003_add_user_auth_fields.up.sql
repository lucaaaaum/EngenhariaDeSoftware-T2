ALTER TABLE users ADD COLUMN IF NOT EXISTS email TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN IF NOT EXISTS password_hash TEXT NOT NULL DEFAULT '';

-- Adiciona constraint de email único depois de adicionar a coluna
CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique ON users (email);
