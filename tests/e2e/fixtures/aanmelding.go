package fixtures

import "dklautomationgo/models"

// GetValidAanmelding retourneert een geldige aanmelding voor tests
func GetValidAanmelding() models.AanmeldingFormulier {
	return models.AanmeldingFormulier{
		Naam:     "E2E Test Persoon",
		Email:    "e2e-test@example.com",
		Telefoon: "0612345678",
		Rol:      "Vrijwilliger",
		Afstand:  "5 KM",
		Terms:    true,
	}
}

// GetInvalidAanmelding retourneert een ongeldige aanmelding voor tests
func GetInvalidAanmelding() models.AanmeldingFormulier {
	return models.AanmeldingFormulier{
		Naam:  "",              // Leeg, wat ongeldig is
		Email: "invalid-email", // Ongeldig email formaat
		Terms: false,           // Moet true zijn
	}
}
