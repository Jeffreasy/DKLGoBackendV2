#!/bin/bash
# Script om E2E tests uit te voeren

echo -e "\033[0;36m🧪 E2E tests starten...\033[0m"

# Zorg ervoor dat de test omgeving schoon is
echo -e "\033[0;33m🧹 Opruimen van vorige test omgeving...\033[0m"
docker-compose -f tests/e2e/docker-compose.e2e.yml down -v --remove-orphans

# Start de test omgeving
echo -e "\033[0;33m🚀 Starten van de test omgeving...\033[0m"
docker-compose -f tests/e2e/docker-compose.e2e.yml up -d

# Wacht tot de services klaar zijn
echo -e "\033[0;33m⏳ Wachten tot services gereed zijn...\033[0m"
sleep 10

# Voer de tests uit
echo -e "\033[0;32m🧪 E2E tests uitvoeren...\033[0m"
go test -v ./tests/e2e/scenarios/...

# Sla de exit code op
TEST_RESULT=$?

# Ruim de test omgeving op
echo -e "\033[0;33m🧹 Opruimen van test omgeving...\033[0m"
docker-compose -f tests/e2e/docker-compose.e2e.yml down -v --remove-orphans

# Geef de juiste exit code terug
if [ $TEST_RESULT -ne 0 ]; then
    echo -e "\033[0;31m❌ E2E tests gefaald!\033[0m"
    exit $TEST_RESULT
else
    echo -e "\033[0;32m✅ E2E tests succesvol afgerond!\033[0m"
    exit 0
fi 