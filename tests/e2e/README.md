# End-to-End (E2E) Tests voor De Koninklijke Loop

Dit directory bevat de End-to-End (E2E) tests voor de De Koninklijke Loop applicatie. Deze tests testen de volledige applicatie, inclusief de frontend, backend, database en externe services.

## Structuur

```
tests/e2e/
├── docker-compose.e2e.yml  # Docker Compose configuratie voor de test omgeving
├── fixtures/               # Test data fixtures
│   ├── aanmelding.go       # Testdata voor aanmeldingen
│   └── auth.go             # Testdata voor authenticatie
├── helpers/                # Helper functies voor de tests
│   ├── api_client.go       # API client voor het maken van HTTP requests
│   ├── assertions.go       # Assertions voor het valideren van responses
│   ├── mailhog_client.go   # MailHog client voor het testen van emails
│   └── test_server.go      # Helper voor het starten/stoppen van de test server
├── scenarios/              # Test scenario's
│   ├── aanmelding_flow_test.go  # Test voor de aanmeldingsflow
│   ├── login_flow_test.go       # Test voor de loginflow
│   └── test_environment.go      # Setup code voor de test omgeving
├── init-db.sql             # SQL script voor het initialiseren van de testdatabase
└── README.md               # Dit bestand
```

## Huidige Status

De volgende test scenario's zijn geïmplementeerd:

1. **Aanmeldingsflow**
   - Test voor het indienen van een geldige aanmelding (✅ Werkend)
   - Test voor validatie bij ongeldige aanmeldingen (✅ Werkend)
   - Email verificatie is momenteel uitgeschakeld in de tests omdat de SMTP instellingen niet correct zijn geconfigureerd in de testomgeving

2. **Loginflow**
   - Test voor mislukte login pogingen (✅ Werkend)
   - Test voor succesvolle login (⚠️ Uitgeschakeld) - De login endpoint werkt niet correct in de testomgeving

Voor meer informatie over de huidige status en geplande verbeteringen, zie [ENDTOENDDocumentatie.md](../../ENDTOENDDocumentatie.md) in de root van het project.

## Tests Uitvoeren

Om de E2E tests uit te voeren, gebruik je het volgende commando vanuit de root van het project:

**Windows (PowerShell):**
```powershell
.\run_e2e_tests.ps1
```

**Linux/macOS (Bash):**
```bash
chmod +x run_e2e_tests.sh
./run_e2e_tests.sh
```

## Test Omgeving

De test omgeving bestaat uit de volgende componenten:

1. **App Container**: De applicatie zelf, gebouwd vanuit de Dockerfile
2. **Database Container**: Een PostgreSQL database voor de tests
3. **MailHog Container**: Een SMTP server voor het testen van emails

De configuratie van deze containers is gedefinieerd in `docker-compose.e2e.yml`.

## Nieuwe Tests Toevoegen

Om een nieuwe E2E test toe te voegen:

1. Maak een nieuw bestand aan in de `scenarios/` directory
2. Implementeer een test suite die `suite.Suite` uitbreidt
3. Implementeer de `SetupSuite` en `TearDownSuite` methoden
4. Voeg test methoden toe die beginnen met `Test`
5. Voeg een test runner functie toe die de suite uitvoert

Voorbeeld:

```go
package scenarios

import (
	"testing"
	"github.com/stretchr/testify/suite"
	"dklautomationgo/tests/e2e/helpers"
)

type MijnTestSuite struct {
	suite.Suite
	env      *TestEnvironment
}

func (s *MijnTestSuite) SetupSuite() {
	var err error
	s.env, err = NewTestEnvironment()
	s.Require().NoError(err)
}

func (s *MijnTestSuite) TearDownSuite() {
	s.env.Cleanup()
}

func (s *MijnTestSuite) TestMijnFeature() {
	// Test implementatie
}

func TestMijnTestSuite(t *testing.T) {
	suite.Run(t, new(MijnTestSuite))
}
```

## Helpers

### API Client

De API client (`helpers/api_client.go`) biedt functionaliteit voor het maken van HTTP requests naar de API. Het ondersteunt zowel geauthenticeerde als niet-geauthenticeerde requests.

```go
// Nieuwe API client maken
client := helpers.NewAPIClient("http://localhost:8080")

// Geauthenticeerde client maken
authClient := helpers.NewAuthenticatedAPIClient("http://localhost:8080", "admin@example.com", "password")

// POST request maken
response, err := client.Post("/api/endpoint", data)

// GET request maken
response, err := client.Get("/api/endpoint")
```

### MailHog Client

De MailHog client (`helpers/mailhog_client.go`) biedt functionaliteit voor het ophalen van emails uit de MailHog server.

```go
// Nieuwe MailHog client maken
mailClient := helpers.NewMailhogClient("http://localhost:8025")

// Emails ophalen voor een specifiek adres
emails, err := mailClient.GetEmailsTo("user@example.com")
```

### Assertions

De assertions helper (`helpers/assertions.go`) biedt functionaliteit voor het valideren van responses en emails.

```go
// JSON response valideren
helpers.AssertJSONResponse(t, response, 200, &responseData)

// Email valideren
helpers.AssertEmailReceived(t, mailClient, "user@example.com", "Email onderwerp")
```

## Fixtures

Fixtures (`fixtures/`) bevatten voorgedefinieerde test data die gebruikt kan worden in de tests.

```go
// Geldige aanmelding ophalen
aanmelding := fixtures.GetValidAanmelding()

// Ongeldige aanmelding ophalen
invalidAanmelding := fixtures.GetInvalidAanmelding()
```

## Troubleshooting

Als je problemen ondervindt bij het uitvoeren van de E2E tests:

1. **Database Connectie Problemen**
   - Controleer of de Docker containers draaien: `docker-compose -f tests/e2e/docker-compose.e2e.yml ps`
   - Controleer de logs van de database container: `docker-compose -f tests/e2e/docker-compose.e2e.yml logs e2e-db`
   - Controleer of de database correct is geïnitialiseerd: `docker-compose -f tests/e2e/docker-compose.e2e.yml exec e2e-db psql -U postgres -d dklautomationgo_e2e -c "\dt"`

2. **App Container Problemen**
   - Controleer de logs van de app container: `docker-compose -f tests/e2e/docker-compose.e2e.yml logs app`
   - Controleer of de health endpoint bereikbaar is: `Invoke-WebRequest -Uri http://localhost:8081/health`

3. **Email Problemen**
   - Controleer of MailHog draait: `docker-compose -f tests/e2e/docker-compose.e2e.yml ps mailhog`
   - Controleer de MailHog web interface op http://localhost:8025
   - Controleer de logs van de app container voor email gerelateerde fouten: `docker-compose -f tests/e2e/docker-compose.e2e.yml logs app | Select-String -Pattern "email"` 