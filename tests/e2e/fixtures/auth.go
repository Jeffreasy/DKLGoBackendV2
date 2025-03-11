package fixtures

// GetValidLoginCredentials retourneert geldige inloggegevens voor tests
func GetValidLoginCredentials() map[string]string {
	return map[string]string{
		"email":    "admin@example.com",
		"password": "admin123",
	}
}

// GetInvalidLoginCredentials retourneert ongeldige inloggegevens voor tests
func GetInvalidLoginCredentials() map[string]string {
	return map[string]string{
		"email":    "admin@example.com",
		"password": "verkeerd_wachtwoord",
	}
}
