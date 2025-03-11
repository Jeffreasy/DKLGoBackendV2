-- database/migrations/000002_add_auth_tables.down.sql
-- Verwijder tabellen en types in omgekeerde volgorde

-- Verwijder refresh_tokens tabel
DROP TABLE IF EXISTS refresh_tokens;

-- Verwijder users tabel
DROP TABLE IF EXISTS users;

-- Verwijder ENUM types
DROP TYPE IF EXISTS user_status;
DROP TYPE IF EXISTS user_role; 