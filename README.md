# De Koninklijke Loop - Technische Documentatie

## Inhoudsopgave
1. [Projectoverzicht](#projectoverzicht)
2. [Architectuur](#architectuur)
3. [Technische Stack](#technische-stack)
4. [Codestructuur](#codestructuur)
5. [API Endpoints](#api-endpoints)
6. [Database Schema](#database-schema)
7. [Authenticatie en Autorisatie](#authenticatie-en-autorisatie)
8. [Email Service](#email-service)
9. [Docker Setup](#docker-setup)
10. [Deployment](#deployment)
11. [Ontwikkelomgeving](#ontwikkelomgeving)
12. [Troubleshooting](#troubleshooting)
13. [Testen](#testen)

## Projectoverzicht

De Koninklijke Loop is een webapplicatie voor het beheren van een hardloopevenement. De applicatie biedt functionaliteit voor:
- Registratie van vrijwilligers
- Contact formulieren
- Email communicatie
- Beheer van gebruikers en rollen
- Admin dashboard

De backend is ontwikkeld in Go (Golang) met een RESTful API architectuur. De applicatie maakt gebruik van PostgreSQL voor dataopslag en is gecontaineriseerd met Docker voor eenvoudige deployment.

## Architectuur

De applicatie volgt een gelaagde architectuur:

1. **Presentatielaag**: API endpoints (handlers)
2. **Servicelaag**: Business logica
3. **Datalaag**: Database repositories
4. **Infrastructuurlaag**: Database connectie, email service, authenticatie

### Dataflow

1. HTTP request komt binnen via een API endpoint
2. Request wordt gevalideerd en verwerkt door een handler
3. Handler roept services aan voor business logica
4. Services gebruiken repositories voor database operaties
5. Response wordt teruggestuurd naar de client

## Technische Stack

- **Backend**: Go (Golang) met Gin web framework
- **Database**: PostgreSQL
- **ORM**: GORM
- **Authenticatie**: JWT (JSON Web Tokens)
- **Email**: SMTP
- **Containerisatie**: Docker
- **CI/CD**: Render.yaml configuratie

## Codestructuur

De codebase is georganiseerd in de volgende hoofdmappen:

### `/auth`
Bevat alle authenticatie-gerelateerde code:
- `/handlers`: API endpoints voor login, registratie, etc.
- `/middleware`: JWT authenticatie middleware
- `/service`: Authenticatie business logica en token management

### `/database`
Database-gerelateerde code:
- `/migrations`: SQL migratie scripts
- `/repository`: Data access layer voor verschillende entiteiten
- `connection.go`: Database connectie setup
- `models.go`: Algemene database modellen

### `/handlers`
API endpoint handlers:
- `aanmelding_handler.go`: Vrijwilligers registratie endpoints
- `contact_handler.go`: Contact formulier endpoints
- `email_handler.go`: Email gerelateerde endpoints

### `/models`
Datamodellen:
- `aanmelding.go`: Vrijwilligers registratie modellen
- `contact.go`: Contact formulier modellen
- `email.go`: Email gerelateerde modellen
- `user.go`: Gebruiker en authenticatie modellen

### `/nginx`
Nginx configuratie voor productie:
- `/conf.d`: Nginx configuratie bestanden
- `/www`: Statische bestanden voor Nginx

### `/scripts`
Helper scripts:
- `docker-helper.ps1/.sh`: Docker management scripts
- `generate-ssl-cert.ps1/.sh`: SSL certificaat generatie
- `migrate.ps1/.sh`: Database migratie scripts
- `healthcheck.sh`: Health check script voor Docker

### `/services`
Business logica services:
- `/email`: Email service implementatie

### `/templates`
HTML email templates:
- `aanmelding_admin_email.html`: Admin notificatie voor nieuwe aanmeldingen
- `aanmelding_email.html`: Bevestigingsmail voor vrijwilligers
- `contact_admin_email.html`: Admin notificatie voor nieuwe contactformulieren
- `contact_email.html`: Bevestigingsmail voor contactformulieren

## API Endpoints

### Publieke Endpoints

#### Contact Formulier
- **POST** `/api/contact`
  - Verwerkt een contactformulier inzending
  - Body: `{ "naam": string, "email": string, "bericht": string, "privacy_akkoord": boolean }`
  - Response: `{ "id": string, "message": string }`

#### Vrijwilligers Aanmelding
- **POST** `/api/aanmelding`
  - Verwerkt een vrijwilliger aanmelding
  - Body: `{ "naam": string, "email": string, "telefoon": string, "rol": string, "afstand": string, "ondersteuning": string, "bijzonderheden": string, "terms": boolean }`
  - Geldige waarden voor `rol`: "Deelnemer", "Vrijwilliger", "Chauffeur", "Bijrijder", "Verzorging"
  - Geldige waarden voor `afstand`: "2.5 KM", "5 KM", "10 KM", "15 KM", "Halve marathon"
  - Response: `{ "id": string, "message": string }`

### Authenticatie Endpoints

- **POST** `/api/auth/login`
  - Authenticatie endpoint
  - Body: `{ "email": string, "password": string }`
  - Response: `{ "access_token": string, "refresh_token": string, "expires_in": number, "token_type": string }`

- **POST** `/api/auth/refresh-token`
  - Vernieuw een verlopen toegangstoken
  - Body: `{ "refresh_token": string }`
  - Response: `{ "access_token": string, "refresh_token": string, "expires_in": number, "token_type": string }`

- **POST** `/api/auth/forgot-password`
  - Start het wachtwoord reset proces
  - Body: `{ "email": string }`
  - Response: `{ "message": string }`

- **POST** `/api/auth/reset-password`
  - Reset een wachtwoord met een reset token
  - Body: `{ "token": string, "new_password": string }`
  - Response: `{ "message": string }`

### Beveiligde Endpoints (Admin)

#### Email Management
- **GET** `/api/emails`
  - Haal alle emails op
  - Response: `{ "data": [Email], "total": number, "has_more": boolean }`

- **GET** `/api/emails/stats`
  - Haal email statistieken op
  - Response: `{ "total": number, "unread": number, "accounts": [{ "name": string, "total": number, "unread": number }] }`

- **PUT** `/api/emails/:id/read`
  - Markeer een email als gelezen
  - Response: `{ "success": boolean }`

#### Contact Management
- **GET** `/api/contacts`
  - Haal alle contactformulieren op
  - Response: `{ "data": [ContactFormulier], "total": number }`

- **PUT** `/api/contacts/:id/status`
  - Werk de status van een contactformulier bij
  - Body: `{ "status": string, "notities": string }`
  - Response: `{ "success": boolean }`

#### Aanmelding Management
- **GET** `/api/aanmeldingen`
  - Haal alle aanmeldingen op
  - Response: `{ "data": [Aanmelding], "total": number }`

- **GET** `/api/aanmeldingen/stats`
  - Haal aanmelding statistieken op
  - Response: `{ "total": number, "by_role": { "role": number }, "by_afstand": { "afstand": number } }`

- **GET** `/api/aanmeldingen/:id`
  - Haal een specifieke aanmelding op
  - Response: `Aanmelding`

### Health Check
- **GET** `/health`
  - Controleer de status van de applicatie
  - Response: `{ "status": "ok", "message": "Service is healthy" }`

## Database Schema

### `contact_formulieren`
Opslag van contactformulieren ingediend via de website.

| Kolom | Type | Beschrijving |
|-------|------|-------------|
| id | UUID | Primaire sleutel |
| created_at | TIMESTAMP | Tijdstip van aanmaken |
| updated_at | TIMESTAMP | Tijdstip van laatste update |
| naam | VARCHAR(100) | Naam van de contactpersoon |
| email | VARCHAR(255) | Email adres |
| bericht | TEXT | Het bericht van de gebruiker |
| email_verzonden | BOOLEAN | Of de bevestigingsemail is verzonden |
| email_verzonden_op | TIMESTAMP | Wanneer de email is verzonden |
| privacy_akkoord | BOOLEAN | Of gebruiker akkoord is met privacy voorwaarden |
| status | VARCHAR(50) | Status van de aanvraag (nieuw/in behandeling/afgerond/gearchiveerd) |
| behandeld_door | VARCHAR(255) | Wie de aanvraag heeft behandeld |
| behandeld_op | TIMESTAMP | Wanneer de aanvraag is behandeld |
| notities | TEXT | Interne notities over de aanvraag |

### `aanmeldingen`
Opslag van vrijwilligersaanmeldingen ingediend via de website.

| Kolom | Type | Beschrijving |
|-------|------|-------------|
| id | UUID | Primaire sleutel |
| created_at | TIMESTAMP | Tijdstip van aanmaken |
| updated_at | TIMESTAMP | Tijdstip van laatste update |
| naam | VARCHAR(100) | Naam van de vrijwilliger |
| email | VARCHAR(255) | Email adres |
| telefoon | VARCHAR(20) | Telefoonnummer |
| rol | VARCHAR(50) | Gewenste rol (Deelnemer, Vrijwilliger, Chauffeur, Bijrijder, Verzorging) |
| afstand | VARCHAR(50) | Maximale reisafstand (2.5 KM, 5 KM, 10 KM, 15 KM, Halve marathon) |
| ondersteuning | TEXT | Benodigde ondersteuning |
| bijzonderheden | TEXT | Eventuele bijzonderheden |
| terms | BOOLEAN | Akkoord met voorwaarden |
| email_verzonden | BOOLEAN | Of de bevestigingsemail is verzonden |
| email_verzonden_op | TIMESTAMP | Wanneer de email is verzonden |

### `users`
Gebruikers van het systeem.

| Kolom | Type | Beschrijving |
|-------|------|-------------|
| id | UUID | Primaire sleutel |
| email | VARCHAR(255) | Email adres (uniek) |
| password_hash | VARCHAR(255) | Gehashte wachtwoord |
| role | user_role | Rol (BEHEERDER, ADMIN, VRIJWILLIGER) |
| status | user_status | Status (PENDING, ACTIVE, INACTIVE) |
| approved_by | UUID | Wie de gebruiker heeft goedgekeurd |
| approved_at | TIMESTAMP | Wanneer de gebruiker is goedgekeurd |
| last_login | TIMESTAMP | Laatste login tijdstip |
| password_reset_token | UUID | Token voor wachtwoord reset |
| password_reset_expires | TIMESTAMP | Vervaldatum van reset token |
| created_at | TIMESTAMP | Tijdstip van aanmaken |
| updated_at | TIMESTAMP | Tijdstip van laatste update |

### `refresh_tokens`
Refresh tokens voor JWT authenticatie.

| Kolom | Type | Beschrijving |
|-------|------|-------------|
| id | UUID | Primaire sleutel |
| user_id | UUID | Gebruiker ID (foreign key) |
| token | VARCHAR(255) | Token waarde (uniek) |
| expires_at | TIMESTAMP | Vervaldatum |
| created_at | TIMESTAMP | Tijdstip van aanmaken |
| revoked | BOOLEAN | Of de token is ingetrokken |
| revoked_at | TIMESTAMP | Wanneer de token is ingetrokken |

## Authenticatie en Autorisatie

De applicatie gebruikt JWT (JSON Web Tokens) voor authenticatie:

1. Gebruiker logt in met email/wachtwoord
2. Server valideert credentials en genereert access token en refresh token
3. Access token wordt gebruikt voor API requests (Authorization header)
4. Refresh token wordt gebruikt om een nieuw access token te krijgen wanneer deze verloopt

### Rollen
- **BEHEERDER**: Volledige toegang tot alle functionaliteit
- **ADMIN**: Toegang tot beheer van aanmeldingen en contactformulieren
- **VRIJWILLIGER**: Beperkte toegang tot eigen gegevens

### Middleware
De `auth.middleware` package bevat middleware voor het valideren van JWT tokens en het controleren van gebruikersrollen.

## Email Service

De email service is verantwoordelijk voor:
1. Het versturen van bevestigingsmails naar gebruikers
2. Het versturen van notificaties naar admins
3. Het ophalen en verwerken van inkomende emails

### Email Accounts
De applicatie gebruikt drie email accounts:
- **info@dekoninklijkeloop.nl**: Algemene communicatie
- **inschrijving@dekoninklijkeloop.nl**: Aanmeldingen
- **noreply@dekoninklijkeloop.nl**: Automatische emails

### Email Templates
HTML templates voor emails zijn opgeslagen in de `/templates` map:
- `aanmelding_admin_email.html`: Admin notificatie voor nieuwe aanmeldingen
- `aanmelding_email.html`: Bevestigingsmail voor vrijwilligers
- `contact_admin_email.html`: Admin notificatie voor nieuwe contactformulieren
- `contact_email.html`: Bevestigingsmail voor contactformulieren

## Docker Setup

De applicatie is gecontaineriseerd met Docker voor eenvoudige deployment en ontwikkeling.

### Containers
- **app**: Go applicatie container
- **db**: PostgreSQL database container

### Docker Compose
De `docker-compose.yml` definieert de ontwikkelomgeving:
- Netwerk configuratie
- Volume mounts
- Omgevingsvariabelen
- Health checks
- Port mappings

### Productie Setup
De `docker-compose.prod.yml` bevat productie-specifieke configuratie:
- Nginx reverse proxy
- SSL/TLS configuratie
- Productie omgevingsvariabelen

### Helper Scripts
- `docker-helper.ps1/.sh`: Scripts voor Docker management (start, stop, logs, etc.)
- `generate-ssl-cert.ps1/.sh`: Scripts voor het genereren van SSL certificaten

## Deployment

### Deployment op Render

De applicatie is geconfigureerd voor eenvoudige deployment op [Render](https://render.com) via de `render.yaml` configuratie.

#### Vereisten
- Een Render account
- Een GitHub repository met de code

#### Stappen voor deployment

1. **Repository voorbereiden**
   - Zorg ervoor dat je repository de volgende bestanden bevat:
     - `render.yaml`
     - `Dockerfile`
     - Alle code en configuratiebestanden

2. **Deployment op Render**
   - Log in op je Render account
   - Ga naar "Blueprints" in het dashboard
   - Klik op "New Blueprint Instance"
   - Selecteer je GitHub repository
   - Render zal automatisch de `render.yaml` detecteren en de services configureren
   - Controleer de configuratie en klik op "Apply"

3. **Omgevingsvariabelen instellen**
   - Na het aanmaken van de services, ga naar de web service
   - Ga naar "Environment" tab
   - Vul de volgende geheime omgevingsvariabelen in:
     - `SMTP_USER`: Email gebruikersnaam
     - `SMTP_PASSWORD`: Email wachtwoord
     - `SMTP_FROM`: Email afzender
     - `INFO_EMAIL_PASSWORD`: Wachtwoord voor info@dekoninklijkeloop.nl
     - `INSCHRIJVING_EMAIL_PASSWORD`: Wachtwoord voor inschrijving@dekoninklijkeloop.nl
     - `NOREPLY_EMAIL_PASSWORD`: Wachtwoord voor noreply@dekoninklijkeloop.nl
     - `ADMIN_EMAIL`: Email adres van de beheerder

4. **Database initialiseren**
   - De database wordt automatisch aangemaakt door Render
   - De migraties worden automatisch uitgevoerd bij de eerste start van de applicatie

5. **SSL/TLS configuratie**
   - Render biedt automatisch SSL/TLS voor alle web services
   - Je kunt een aangepast domein configureren in de "Settings" tab van je web service

### Handmatige deployment

Voor handmatige deployment, volg deze stappen:

1. **Database voorbereiden**
   - Maak een PostgreSQL database aan
   - Voer de migraties uit met het commando:
     ```
     ./migrate -path=./database/migrations -database "postgres://user:password@host:port/dbname?sslmode=disable" up
     ```

2. **Applicatie bouwen**
   - Bouw de applicatie met:
     ```
     go build -o dklautomationgo .
     ```

3. **Configuratie**
   - Maak een `.env` bestand aan op basis van `.env.example`
   - Vul alle vereiste omgevingsvariabelen in

4. **Applicatie starten**
   - Start de applicatie met:
     ```
     ./dklautomationgo
     ```

### Docker deployment

Voor deployment met Docker, volg deze stappen:

1. **Docker Compose gebruiken**
   - Zorg ervoor dat Docker en Docker Compose zijn geïnstalleerd
   - Maak een `.env` bestand aan op basis van `.env.example`
   - Start de applicatie met:
     ```
     docker-compose up -d
     ```

2. **Logs bekijken**
   - Bekijk de logs met:
     ```
     docker-compose logs -f
     ```

3. **Applicatie stoppen**
   - Stop de applicatie met:
     ```
     docker-compose down
     ```

## Ontwikkelomgeving

### Vereisten
- Go 1.23+
- Docker Desktop
- PostgreSQL (lokaal of via Docker)
- Git

### Setup
1. Clone de repository
2. Kopieer `.env.example` naar `.env` en pas de waarden aan
3. Start de Docker containers: `.\scripts\docker-helper.ps1 start`
4. De applicatie is beschikbaar op http://localhost:8080

### Database Migraties
Database migraties worden automatisch uitgevoerd bij het starten van de container. Handmatige migraties kunnen worden uitgevoerd met:
```powershell
.\scripts\migrate.ps1 up
```

## Testen

De applicatie bevat verschillende soorten tests om de functionaliteit te valideren:

### Unit Tests

Unit tests testen individuele componenten in isolatie, zonder afhankelijkheid van externe systemen zoals de database. Deze tests zijn snel en betrouwbaar.

#### Beschikbare Unit Tests

- **Services Tests**: Testen voor de business logica in de services laag
- **Handlers Tests**: Testen voor de API endpoints
- **Auth Service Tests**: Testen voor de authenticatie service
- **Auth Handlers Tests**: Testen voor de authenticatie endpoints

#### Unit Tests Uitvoeren

Je kunt alle unit tests uitvoeren met het volgende commando:

```powershell
.\scripts\run_unit_tests.ps1
```

Of individuele test packages uitvoeren:

```bash
go test ./services -v
go test ./handlers -v
go test ./auth/service -v
go test ./auth/handlers -v
```

### Integratie Tests

Integratie tests testen de interactie tussen verschillende componenten, inclusief de database. Deze tests vereisen een testdatabase.

#### Testdatabase Setup

Voordat je integratie tests kunt uitvoeren, moet je een testdatabase aanmaken:

```powershell
.\scripts\setup_test_db.ps1
```

Dit script:
1. Controleert of Docker draait
2. Controleert of de database container actief is
3. Maakt een testdatabase aan (`dklautomationgo_test`)
4. Stelt de juiste omgevingsvariabelen in voor de tests

#### Integratie Tests Uitvoeren

Na het opzetten van de testdatabase kun je de integratie tests uitvoeren:

```bash
go test ./database/repository -v
go test ./tests/integration -v
```

**Opmerking**: Voor het uitvoeren van integratie tests moet de database container draaien en toegankelijk zijn.

#### Alleen Integratie Tests Uitvoeren

Om alleen de integratie tests uit te voeren, kun je het volgende PowerShell script gebruiken:

```powershell
.\scripts\run-tests.ps1 -Integration
```

Dit script:
1. Start de test database container
2. Wacht tot de database gereed is
3. Voert alleen de integratie tests uit
4. Stopt en verwijdert de test database container na afloop

Als je problemen ondervindt met de database verbinding, kun je ook het volgende script gebruiken:

```powershell
.\run_integration_tests.ps1
```

Dit script zorgt voor:
1. Het starten van de test database container
2. Het instellen van de juiste omgevingsvariabelen
3. Het aanmaken van de test database als deze nog niet bestaat
4. Het uitvoeren van de integratie tests
5. Het opruimen van resources na afloop

**Omgevingsvariabelen voor Integratie Tests**

Voor het uitvoeren van integratie tests zijn de volgende omgevingsvariabelen nodig:

```powershell
$env:TEST_DB_HOST="localhost"
$env:TEST_DB_PORT="5433"  # Let op: dit is de gemapte poort in docker-compose.test.yml
$env:TEST_DB_USER="postgres"
$env:TEST_DB_PASSWORD="Bootje@12"
$env:TEST_DB_NAME="dklautomationgo_test"
```

### Mocks

Voor unit tests worden mock implementaties gebruikt van repositories en services:

- `tests/mocks/aanmelding_repository.go`: Mock voor de aanmelding repository
- `tests/mocks/aanmelding_service.go`: Mock voor de aanmelding service
- `tests/mocks/auth_middleware.go`: Mock voor de authenticatie middleware
- `tests/mocks/auth_service.go`: Mock voor de authenticatie service
- `tests/mocks/email_service.go`: Mock voor de email service

### Test Fixtures

Test fixtures bevatten testdata die wordt gebruikt in de tests:

- `tests/fixtures/aanmeldingen.go`: Testdata voor aanmeldingen
- `tests/fixtures/users.go`: Testdata voor gebruikers

### Troubleshooting Tests

#### Database Connectie Problemen

Als je problemen ondervindt met de database connectie tijdens het testen:

1. Controleer of de Docker container draait: `docker ps`
2. Controleer of de testdatabase bestaat: `docker exec dklautomationgo-db psql -U postgres -c "\l"`
3. Controleer de omgevingsvariabelen: 
   ```powershell
   echo $env:TEST_DB_HOST
   echo $env:TEST_DB_PORT
   echo $env:TEST_DB_USER
   echo $env:TEST_DB_PASSWORD
   echo $env:TEST_DB_NAME
   ```

#### Testdatabase Opnieuw Aanmaken

Als je de testdatabase volledig opnieuw wilt aanmaken:

```powershell
docker exec -it dklautomationgo-db psql -U postgres -c "DROP DATABASE IF EXISTS dklautomationgo_test;"
docker exec -it dklautomationgo-db psql -U postgres -c "CREATE DATABASE dklautomationgo_test;"
```

## Troubleshooting

### Bekende Problemen

#### 1. Database Connectie Problemen
- **Symptoom**: Foutmelding "failed to connect to database"
- **Oplossing**: Controleer de database instellingen in `.env` en zorg ervoor dat de database draait

#### 2. Email Verzending Problemen
- **Symptoom**: Foutmelding "failed to send email"
- **Oplossing**: Controleer de SMTP instellingen en wachtwoorden in `.env`

#### 3. Aanmelding Validatie Fouten
- **Symptoom**: Foutmelding "violates check constraint"
- **Oplossing**: Zorg ervoor dat de waarden voor `rol` en `afstand` overeenkomen met de toegestane waarden:
  - Geldige waarden voor `rol`: "Deelnemer", "Vrijwilliger", "Chauffeur", "Bijrijder", "Verzorging"
  - Geldige waarden voor `afstand`: "2.5 KM", "5 KM", "10 KM", "15 KM", "Halve marathon"

#### 4. Docker Script Problemen
- **Symptoom**: Foutmelding "cannot execute: required file not found" voor `/app/scripts/migrate.sh`
- **Oplossing**: 
  - Controleer of de scripts uitvoerbaar zijn (`chmod +x`) en correct zijn gekopieerd naar de container
  - Als het probleem aanhoudt, pas de `docker-compose.yml` aan om de migrate.sh aanroep te verwijderen:
    ```yaml
    command: ["/bin/bash", "-c", "/app/dklautomationgo"]
    ```
  - Alternatief: voer de migraties handmatig uit na het starten van de container:
    ```powershell
    docker exec dklautomationgo-app /app/scripts/migrate.sh
    ```

#### 5. Testdatabase Problemen
- **Symptoom**: Foutmelding "database dklautomationgo_test does not exist" tijdens het uitvoeren van tests
- **Oplossing**: Maak de testdatabase aan met het script `.\scripts\setup_test_db.ps1` of handmatig:
  ```powershell
  docker exec -it dklautomationgo-db psql -U postgres -c "CREATE DATABASE dklautomationgo_test;"
  ```

### Logging
De applicatie logt informatie naar stdout, wat kan worden bekeken met:
```powershell
.\scripts\docker-helper.ps1 logs
```

### Health Check
De applicatie biedt een health check endpoint op `/health` om de status te controleren.

## Tests

De applicatie bevat verschillende soorten tests:

### Unit Tests

Unit tests testen individuele componenten in isolatie.

```bash
go test -v ./... -short
# of
make test-unit
```

### Integratie Tests

Integratie tests testen de interactie tussen verschillende componenten.

```bash
# Windows
.\run_integration_tests.ps1
# of
make test-integration

# Linux/macOS
chmod +x run_integration_tests.sh
./run_integration_tests.sh
# of
make test-integration
```

### End-to-End (E2E) Tests

End-to-End tests testen de volledige applicatie, inclusief de frontend, backend, database en externe services. De E2E tests zijn geïmplementeerd met behulp van Docker Compose om een geïsoleerde testomgeving op te zetten.

#### Huidige Status

De E2E tests zijn succesvol geïmplementeerd en kunnen worden uitgevoerd met het volgende commando:

```bash
# Windows
.\run_e2e_tests.ps1
# of
make test-e2e

# Linux/macOS
chmod +x run_e2e_tests.sh
./run_e2e_tests.sh
# of
make test-e2e
```

De volgende test scenario's zijn geïmplementeerd:
- **Aanmeldingsflow**: Testen van het indienen van een aanmelding
- **Validatie**: Testen van validatie bij ongeldige aanmeldingen
- **Login**: Testen van mislukte login pogingen

Opmerking: Sommige tests zijn momenteel uitgeschakeld vanwege beperkingen in de testomgeving. Zie [ENDTOENDDocumentatie.md](ENDTOENDDocumentatie.md) voor meer informatie over de huidige status en geplande verbeteringen.

#### Testomgeving

De E2E tests maken gebruik van:
- Een geïsoleerde PostgreSQL database
- MailHog voor het testen van email functionaliteit
- Een aparte instantie van de applicatie

#### Troubleshooting

Als je problemen ondervindt bij het uitvoeren van de E2E tests, raadpleeg dan de [ENDTOENDDocumentatie.md](ENDTOENDDocumentatie.md) voor troubleshooting tips.

### Alle Tests Uitvoeren

Om alle tests uit te voeren:

```bash
make test
```

Zie [TESTING.md](TESTING.md) voor meer informatie over het testen van de applicatie.

## Makefile

De applicatie bevat een Makefile met verschillende commando's:

```bash
# Bouw de applicatie
make build

# Voer de applicatie uit
make run

# Voer alle tests uit
make test

# Voer alleen unit tests uit
make test-unit

# Voer alleen integratie tests uit
make test-integration

# Voer alleen end-to-end tests uit
make test-e2e

# Ruim build artifacts en test containers op
make clean

# Toon help informatie
make help
``` 