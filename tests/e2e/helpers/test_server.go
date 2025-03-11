package helpers

import (
	"fmt"
	"net/http"
	"os/exec"
	"time"
)

// TestServer is een helper voor het starten en stoppen van de test server
type TestServer struct {
	cmd *exec.Cmd
}

// NewTestServer maakt een nieuwe test server
func NewTestServer() *TestServer {
	return &TestServer{}
}

// Start start de test server
func (s *TestServer) Start() error {
	// In een echte implementatie zou je hier de applicatie starten
	// Voor nu gebruiken we de docker-compose setup
	fmt.Println("Test server is gestart via docker-compose")

	// Wacht tot de server gereed is
	for i := 0; i < 10; i++ {
		resp, err := http.Get("http://localhost:8080/health")
		if err == nil && resp.StatusCode == http.StatusOK {
			fmt.Println("Test server is gereed")
			return nil
		}
		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("timeout waiting for test server to be ready")
}

// Stop stopt de test server
func (s *TestServer) Stop() error {
	// In een echte implementatie zou je hier de applicatie stoppen
	// Voor nu gebruiken we de docker-compose setup
	fmt.Println("Test server is gestopt via docker-compose")
	return nil
}
