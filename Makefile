.PHONY: test test-unit test-integration test-e2e clean build run

# Standaard target
all: build

# Bouw de applicatie
build:
	@echo "🔨 Bouwen van de applicatie..."
	go build -o bin/dklautomationgo

# Voer de applicatie uit
run: build
	@echo "🚀 Starten van de applicatie..."
	./bin/dklautomationgo

# Voer alle tests uit
test: test-unit test-integration test-e2e

# Voer unit tests uit
test-unit:
	@echo "🧪 Unit tests uitvoeren..."
	go test -v ./... -short

# Voer integratie tests uit
test-integration:
	@echo "🧪 Integratie tests uitvoeren..."
	@if [ "$(OS)" = "Windows_NT" ]; then \
		powershell -ExecutionPolicy Bypass -File ./run_integration_tests.ps1; \
	else \
		chmod +x ./run_integration_tests.sh && ./run_integration_tests.sh; \
	fi

# Voer end-to-end tests uit
test-e2e:
	@echo "🧪 End-to-End tests uitvoeren..."
	@if [ "$(OS)" = "Windows_NT" ]; then \
		powershell -ExecutionPolicy Bypass -File ./run_e2e_tests.ps1; \
	else \
		chmod +x ./run_e2e_tests.sh && ./run_e2e_tests.sh; \
	fi

# Schoon de build directory op
clean:
	@echo "🧹 Opruimen van build artifacts..."
	rm -rf bin/
	@echo "🧹 Opruimen van test containers..."
	docker-compose -f tests/integration/docker-compose.test.yml down -v --remove-orphans
	docker-compose -f tests/e2e/docker-compose.e2e.yml down -v --remove-orphans

# Help informatie
help:
	@echo "Beschikbare commando's:"
	@echo "  make build            - Bouw de applicatie"
	@echo "  make run              - Voer de applicatie uit"
	@echo "  make test             - Voer alle tests uit"
	@echo "  make test-unit        - Voer alleen unit tests uit"
	@echo "  make test-integration - Voer alleen integratie tests uit"
	@echo "  make test-e2e         - Voer alleen end-to-end tests uit"
	@echo "  make clean            - Ruim build artifacts en test containers op"
	@echo "  make help             - Toon deze help informatie" 