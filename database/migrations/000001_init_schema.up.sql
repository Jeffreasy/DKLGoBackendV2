-- database/migrations/000001_init_schema.up.sql
-- Maak de uuid-ossp extensie aan voor het genereren van UUIDs
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Contact formulieren tabel
CREATE TABLE IF NOT EXISTS contact_formulieren (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    naam VARCHAR(100) NOT NULL CHECK (LENGTH(TRIM(naam)) > 0),
    email VARCHAR(255) NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    bericht TEXT NOT NULL CHECK (LENGTH(TRIM(bericht)) > 0),
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMP WITH TIME ZONE,
    privacy_akkoord BOOLEAN NOT NULL DEFAULT FALSE,
    status VARCHAR(50) NOT NULL DEFAULT 'nieuw' CHECK (status IN ('nieuw', 'in behandeling', 'afgerond', 'gearchiveerd')),
    behandeld_door VARCHAR(255),
    behandeld_op TIMESTAMP WITH TIME ZONE,
    notities TEXT
);

-- Aanmeldingen tabel
CREATE TABLE IF NOT EXISTS aanmeldingen (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    naam VARCHAR(100) NOT NULL CHECK (LENGTH(TRIM(naam)) > 0),
    email VARCHAR(255) NOT NULL CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    telefoon VARCHAR(20) NOT NULL CHECK (LENGTH(TRIM(telefoon)) > 0),
    rol VARCHAR(50) NOT NULL CHECK (rol IN ('Deelnemer', 'Vrijwilliger', 'Chauffeur', 'Bijrijder', 'Verzorging')),
    afstand VARCHAR(50) NOT NULL CHECK (afstand IN ('2.5 KM', '5 KM', '10 KM', '15 KM', 'Halve marathon')),
    ondersteuning TEXT,
    bijzonderheden TEXT,
    terms BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden BOOLEAN NOT NULL DEFAULT FALSE,
    email_verzonden_op TIMESTAMP WITH TIME ZONE
);

-- Commentaar voor documentatie
COMMENT ON TABLE contact_formulieren IS 'Opslag van contactformulieren ingediend via de website';
COMMENT ON TABLE aanmeldingen IS 'Opslag van vrijwilligersaanmeldingen ingediend via de website';

-- Indexen voor contact_formulieren
CREATE INDEX idx_contact_email ON contact_formulieren(email);
CREATE INDEX idx_contact_status ON contact_formulieren(status);
CREATE INDEX idx_contact_created_at ON contact_formulieren(created_at DESC);

-- Indexen voor aanmeldingen
CREATE INDEX idx_aanmelding_email ON aanmeldingen(email);
CREATE INDEX idx_aanmelding_rol ON aanmeldingen(rol);
CREATE INDEX idx_aanmelding_afstand ON aanmeldingen(afstand);
CREATE INDEX idx_aanmelding_created_at ON aanmeldingen(created_at DESC);

-- Trigger voor automatische update van updated_at bij contact_formulieren
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_contact_formulieren_updated_at
BEFORE UPDATE ON contact_formulieren
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_aanmeldingen_updated_at
BEFORE UPDATE ON aanmeldingen
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();