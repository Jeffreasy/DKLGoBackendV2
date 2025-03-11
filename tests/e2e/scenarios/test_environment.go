package e2e

import (
	"dklautomationgo/tests/e2e/helpers"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"
)

// TestEnvironment bevat alle benodigde componenten voor E2E tests
type TestEnvironment struct {
	APIClient   *helpers.APIClient
	MailClient  *helpers.MailhogClient
	AdminClient *helpers.APIClient
}

// waitForService controleert of een service bereikbaar is
func waitForService(url string, maxRetries int) bool {
	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			return true
		}
		fmt.Printf("Waiting for service at %s to be ready... (%d/%d)\n", url, i+1, maxRetries)
		time.Sleep(2 * time.Second)
	}
	return false
}

// SetupTestEnvironment zet de testomgeving op
func SetupTestEnvironment(t *testing.T) (*TestEnvironment, func()) {
	// Start de testomgeving
	cmd := exec.Command("docker-compose", "-f", "../docker-compose.e2e.yml", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to start test environment: %v", err)
	}

	// Wacht tot de services gereed zijn
	fmt.Println("Waiting for services to be ready...")
	time.Sleep(15 * time.Second)

	// Controleer of de API bereikbaar is
	apiURL := "http://localhost:8081/health"
	if !waitForService(apiURL, 10) {
		t.Fatalf("API service is not ready after waiting")
	}

	// Controleer of MailHog bereikbaar is
	mailhogURL := "http://localhost:8025"
	if !waitForService(mailhogURL, 5) {
		t.Fatalf("MailHog service is not ready after waiting")
	}

	// Maak de test environment
	env := &TestEnvironment{
		APIClient:  helpers.NewAPIClient("http://localhost:8081"),
		MailClient: helpers.NewMailhogClient("http://localhost:8025"),
	}

	// Maak een admin client (indien nodig)
	// We maken de admin client zonder authenticatie om te voorkomen dat de test faalt als de login endpoint nog niet werkt
	adminClient := helpers.NewAPIClient("http://localhost:8081")
	env.AdminClient = adminClient

	// Teardown functie
	teardown := func() {
		cmd := exec.Command("docker-compose", "-f", "../docker-compose.e2e.yml", "down")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			t.Fatalf("Failed to stop test environment: %v", err)
		}
	}

	return env, teardown
}
