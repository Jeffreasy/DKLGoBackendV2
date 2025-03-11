# Testdocumentatie voor De Koninklijke Loop

Dit document beschrijft de testprocedures voor de De Koninklijke Loop applicatie.

## Inhoudsopgave
1. [Teststructuur](#teststructuur)
2. [Unit Tests](#unit-tests)
3. [Integratie Tests](#integratie-tests)
4. [Mocks en Fixtures](#mocks-en-fixtures)
5. [Testdatabase Setup](#testdatabase-setup)
6. [Testscripts](#testscripts)
7. [Troubleshooting](#troubleshooting)

## Teststructuur

De tests zijn georganiseerd volgens de volgende structuur:

```
dklautomationgo/
├── auth/
│   ├── handlers/
│   │   └── auth_handler_test.go
│   └── service/
│       └── token_test.go
├── database/
│   └── repository/
│       └── aanmelding_repository_test.go
├── handlers/
│   └── aanmelding_handler_test.go
├── services/
│   └── aanmelding_service_test.go
└── tests/
    ├── fixtures/
    │   ├── aanmeldingen.go
    │   └── users.go
    ├── integration/
    │   ├── auth_integration_test.go
    │   └── protected_routes_test.go
    ├── mocks/
    │   ├── aanmelding_handler.go
    │   ├── aanmelding_repository.go
    │   ├── aanmelding_service.go
    │   ├── auth_middleware.go
    │   ├── auth_service.go
    │   └── email_service.go
    └── setup_test_db.go
```

## Unit Tests

Unit tests testen individuele componenten in isolatie, zonder afhankelijkheid van externe systemen zoals de database. Deze tests zijn snel en betrouwbaar.

### Services Tests

De services tests testen de business logica in de services laag. Deze tests gebruiken mocks voor repositories en andere services om de afhankelijkheden te isoleren.

**Locatie**: `services/aanmelding_service_test.go`

**Geteste functionaliteit**:
- `CreateAanmelding`: Het aanmaken van een nieuwe aanmelding
- `GetAanmeldingen`: Het ophalen van aanmeldingen
- `GetAanmeldingByID`: Het ophalen van een specifieke aanmelding
- `SendBevestigingsEmail`: Het versturen van een bevestigingsmail

**Voorbeeld test**:
```go
func TestCreateAanmelding_Success(t *testing.T) {
    // Setup
    service, mockRepo, mockEmailService := setupAanmeldingServiceTest()
    aanmelding := &models.Aanmelding{
        Naam: "Test Persoon",
        Email: "test@example.com",
    }
    
    // Expectations
    mockRepo.On("Create", aanmelding).Return(nil)
    mockRepo.On("Update", mock.Anything).Return(nil)
    mockEmailService.On("SendAanmeldingEmail", mock.Anything).Return(nil)
    
    // Execute
    err := service.CreateAanmelding(aanmelding)
    
    // Assert
    assert.NoError(t, err)
    mockRepo.AssertExpectations(t)
    mockEmailService.AssertExpectations(t)
}
```

### Handlers Tests

De handlers tests testen de API endpoints. Deze tests gebruiken mocks voor services om de afhankelijkheden te isoleren.

**Locatie**: `handlers/aanmelding_handler_test.go`

**Geteste functionaliteit**:
- `CreateAanmelding`: Het verwerken van een aanmeldingsverzoek
- `GetAanmeldingen`: Het ophalen van aanmeldingen
- `GetAanmeldingByID`: Het ophalen van een specifieke aanmelding

**Voorbeeld test**:
```go
func TestCreateAanmelding_Success(t *testing.T) {
    // Setup
    mockService := new(mocks.MockAanmeldingService)
    handler := NewAanmeldingHandler(mockService)
    
    // Test data
    aanmelding := models.AanmeldingFormulier{
        Naam: "Test Persoon",
        Email: "test@example.com",
        Telefoon: "0612345678",
        Rol: "Vrijwilliger",
        Afstand: "5 KM",
        Terms: true,
    }
    
    // Expectations
    mockService.On("CreateAanmelding", mock.Anything).Return(nil)
    
    // Setup HTTP request
    jsonData, _ := json.Marshal(aanmelding)
    req, _ := http.NewRequest("POST", "/aanmeldingen", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    // Setup response recorder
    w := httptest.NewRecorder()
    
    // Setup router
    router := gin.Default()
    router.POST("/aanmeldingen", handler.CreateAanmelding)
    
    // Execute
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    mockService.AssertExpectations(t)
}
```

### Auth Service Tests

De auth service tests testen de authenticatie service. Deze tests testen de token generatie en validatie.

**Locatie**: `auth/service/token_test.go`

**Geteste functionaliteit**:
- `GenerateAccessToken`: Het genereren van een access token
- `GetUserIDFromToken`: Het extraheren van de gebruikers-ID uit een token

**Voorbeeld test**:
```go
func TestGenerateAccessToken(t *testing.T) {
    // Setup
    tokenService := NewTokenService()
    user := &models.User{
        ID: uuid.New(),
        Email: "test@example.com",
        Role: models.RoleAdmin,
    }
    
    // Execute
    token, err := tokenService.GenerateAccessToken(user)
    
    // Assert
    assert.NoError(t, err)
    assert.NotEmpty(t, token)
    
    // Validate token
    claims, err := tokenService.ValidateToken(token)
    assert.NoError(t, err)
    assert.Equal(t, user.ID.String(), claims["sub"])
    assert.Equal(t, user.Email, claims["email"])
    assert.Equal(t, string(user.Role), claims["role"])
}
```

### Auth Handlers Tests

De auth handlers tests testen de authenticatie endpoints. Deze tests gebruiken mocks voor de auth service.

**Locatie**: `auth/handlers/auth_handler_test.go`

**Geteste functionaliteit**:
- `Login`: Het inloggen van een gebruiker
- `RefreshToken`: Het vernieuwen van een token
- `Logout`: Het uitloggen van een gebruiker

**Voorbeeld test**:
```go
func TestLogin_Success(t *testing.T) {
    // Setup
    mockAuthService := new(mocks.MockAuthService)
    mockAuthMiddleware := new(mocks.MockAuthMiddleware)
    handler := NewAuthHandler(mockAuthService, mockAuthMiddleware)
    
    // Test data
    loginRequest := models.LoginRequest{
        Email: "test@example.com",
        Password: "password123",
    }
    tokenResponse := &models.TokenResponse{
        AccessToken: "access_token",
        RefreshToken: "refresh_token",
        ExpiresIn: 900,
        TokenType: "Bearer",
    }
    
    // Expectations
    mockAuthService.On("Login", loginRequest.Email, loginRequest.Password).Return(tokenResponse, nil)
    
    // Setup HTTP request
    jsonData, _ := json.Marshal(loginRequest)
    req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    
    // Setup response recorder
    w := httptest.NewRecorder()
    
    // Setup router
    router := gin.Default()
    router.POST("/auth/login", handler.Login)
    
    // Execute
    router.ServeHTTP(w, req)
    
    // Assert
    assert.Equal(t, http.StatusOK, w.Code)
    mockAuthService.AssertExpectations(t)
}
```

## Integratie Tests

Integratie tests testen de interactie tussen verschillende componenten, inclusief de database. Deze tests vereisen een testdatabase.

### Repository Tests

De repository tests testen de interactie met de database. Deze tests gebruiken een echte database verbinding.

**Locatie**: `database/repository/aanmelding_repository_test.go`

**Geteste functionaliteit**:
- `Create`: Het aanmaken van een nieuwe aanmelding in de database
- `FindAll`: Het ophalen van aanmeldingen uit de database
- `FindByID`: Het ophalen van een specifieke aanmelding uit de database
- `Update`: Het bijwerken van een aanmelding in de database
- `Count`: Het tellen van aanmeldingen in de database

### API Integratie Tests

De API integratie tests testen de volledige API, inclusief de interactie met de database. Deze tests gebruiken een echte database verbinding.

**Locatie**: `tests/integration/auth_integration_test.go`, `tests/integration/protected_routes_test.go`

**Geteste functionaliteit**:
- Authenticatie flow (login, token vernieuwing, logout)
- Toegang tot beveiligde routes met verschillende gebruikersrollen

## Mocks en Fixtures

### Mocks

Mocks worden gebruikt om afhankelijkheden te simuleren tijdens het testen. De volgende mocks zijn beschikbaar:

- `tests/mocks/aanmelding_repository.go`: Mock voor de aanmelding repository
- `tests/mocks/aanmelding_service.go`: Mock voor de aanmelding service
- `tests/mocks/auth_middleware.go`: Mock voor de authenticatie middleware
- `tests/mocks/auth_service.go`: Mock voor de authenticatie service
- `tests/mocks/email_service.go`: Mock voor de email service

### Fixtures

Fixtures bevatten testdata die wordt gebruikt in de tests:

- `tests/fixtures/aanmeldingen.go`: Testdata voor aanmeldingen
- `tests/fixtures/users.go`: Testdata voor gebruikers

## Testdatabase Setup

Voor het uitvoeren van integratie tests is een testdatabase nodig. De `tests/setup_test_db.go` file bevat code om een verbinding met de testdatabase op te zetten.

### Handmatige Setup

1. Maak een testdatabase aan:
   ```powershell
   docker exec -it dklautomationgo-db psql -U postgres -c "CREATE DATABASE dklautomationgo_test;"
   ```

2. Stel de omgevingsvariabelen in:
   ```powershell
   $env:TEST_DB_HOST="localhost"
   $env:TEST_DB_PORT="5432"
   $env:TEST_DB_USER="postgres"
   $env:TEST_DB_PASSWORD="Bootje@12"
   $env:TEST_DB_NAME="dklautomationgo_test"
   ```

### Geautomatiseerde Setup

Gebruik het `scripts/setup_test_db.ps1` script om de testdatabase automatisch aan te maken en de omgevingsvariabelen in te stellen.

## Testscripts

### `scripts/run_unit_tests.ps1`

Dit script voert alle unit tests uit die geen database verbinding nodig hebben:

```powershell
#!/usr/bin/env pwsh
# Dit script voert alleen de unit tests uit die geen database verbinding nodig hebben

Write-Host "Uitvoeren van unit tests voor services..."
go test ./services -v

Write-Host "Uitvoeren van unit tests voor handlers..."
go test ./handlers -v

Write-Host "Uitvoeren van unit tests voor auth/service..."
go test ./auth/service -v

Write-Host "Uitvoeren van unit tests voor auth/handlers..."
go test ./auth/handlers -v

Write-Host "Unit tests voltooid."
```

### `scripts/setup_test_db.ps1`

Dit script maakt een testdatabase aan en stelt de omgevingsvariabelen in:

```powershell
#!/usr/bin/env pwsh
# Dit script maakt een testdatabase aan voor de tests

# Configuratie
$DB_HOST = if ($env:TEST_DB_HOST) { $env:TEST_DB_HOST } else { "localhost" }
$DB_PORT = if ($env:TEST_DB_PORT) { $env:TEST_DB_PORT } else { "5432" }
$DB_USER = if ($env:TEST_DB_USER) { $env:TEST_DB_USER } else { "postgres" }
$DB_PASSWORD = if ($env:TEST_DB_PASSWORD) { $env:TEST_DB_PASSWORD } else { "Bootje@12" }
$DB_NAME = if ($env:TEST_DB_NAME) { $env:TEST_DB_NAME } else { "dklautomationgo_test" }

# Controleer of Docker draait
$dockerRunning = docker ps 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "Docker is niet actief. Start Docker en probeer het opnieuw."
    exit 1
}

# Controleer of de database container draait
$dbContainer = docker ps --filter "name=dklautomationgo-db" --format "{{.Names}}"
if (-not $dbContainer) {
    Write-Host "De database container 'dklautomationgo-db' is niet actief. Start de container en probeer het opnieuw."
    exit 1
}

# Maak de testdatabase aan
Write-Host "Aanmaken van testdatabase '$DB_NAME'..."
$createDbCmd = "docker exec dklautomationgo-db psql -U $DB_USER -c 'CREATE DATABASE $DB_NAME;'"
Invoke-Expression $createDbCmd

if ($LASTEXITCODE -eq 0) {
    Write-Host "Testdatabase '$DB_NAME' is succesvol aangemaakt."
} else {
    Write-Host "Er is een fout opgetreden bij het aanmaken van de testdatabase. Mogelijk bestaat deze al."
}

# Stel de omgevingsvariabelen in voor de tests
$env:TEST_DB_HOST = $DB_HOST
$env:TEST_DB_PORT = $DB_PORT
$env:TEST_DB_USER = $DB_USER
$env:TEST_DB_PASSWORD = $DB_PASSWORD
$env:TEST_DB_NAME = $DB_NAME

Write-Host "Omgevingsvariabelen zijn ingesteld voor de tests."
Write-Host "Je kunt nu de tests uitvoeren met: go test ./..."
```

## Troubleshooting

### Database Connectie Problemen

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

### Testdatabase Opnieuw Aanmaken

Als je de testdatabase volledig opnieuw wilt aanmaken:

```powershell
docker exec -it dklautomationgo-db psql -U postgres -c "DROP DATABASE IF EXISTS dklautomationgo_test;"
docker exec -it dklautomationgo-db psql -U postgres -c "CREATE DATABASE dklautomationgo_test;"
```

### Mock Verwachtingen Problemen

Als je problemen ondervindt met mock verwachtingen:

1. Controleer of alle verwachte methode aanroepen zijn gedefinieerd met `On()`
2. Controleer of de parameters overeenkomen
3. Gebruik `mock.Anything` voor parameters die je niet exact wilt matchen
4. Voeg `AssertExpectations(t)` toe aan het einde van je test om te controleren of alle verwachtingen zijn voldaan 