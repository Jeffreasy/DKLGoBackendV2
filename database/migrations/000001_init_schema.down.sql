-- database/migrations/000001_init_schema.down.sql

-- Verwijder triggers
DROP TRIGGER IF EXISTS update_contact_formulieren_updated_at ON contact_formulieren;
DROP TRIGGER IF EXISTS update_aanmeldingen_updated_at ON aanmeldingen;

-- Verwijder functies
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Verwijder tabellen
DROP TABLE IF EXISTS contact_formulieren;
DROP TABLE IF EXISTS aanmeldingen;

-- Verwijder extensies (optioneel, kan andere applicaties beïnvloeden)
-- DROP EXTENSION IF EXISTS "uuid-ossp";