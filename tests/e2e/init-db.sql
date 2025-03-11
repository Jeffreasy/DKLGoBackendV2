-- Installeer de uuid-ossp extensie
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Maak de user_role en user_status types aan
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role') THEN
        CREATE TYPE user_role AS ENUM ('BEHEERDER', 'ADMIN', 'VRIJWILLIGER');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status') THEN
        CREATE TYPE user_status AS ENUM ('PENDING', 'ACTIVE', 'INACTIVE');
    END IF;
END
$$;

-- Geef alle rechten aan de postgres gebruiker
GRANT ALL PRIVILEGES ON DATABASE dklautomationgo_e2e TO postgres;

-- Maak de users tabel aan als deze nog niet bestaat
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role user_role NOT NULL DEFAULT 'VRIJWILLIGER',
    status user_status NOT NULL DEFAULT 'ACTIVE',
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Maak de refresh_tokens tabel aan als deze nog niet bestaat
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Voeg een admin gebruiker toe voor tests
-- Wachtwoord: admin123 (bcrypt hash)
INSERT INTO users (email, password_hash, role, status)
VALUES ('admin@example.com', '$2a$10$3QxDjD1ylgPnRgQLhBrTaeqdsNaLxkk2Y3p5mZ4QmY7NWz0.7ZwJ2', 'ADMIN', 'ACTIVE')
ON CONFLICT (email) DO NOTHING; 