package fixtures

import (
	"dklautomationgo/models"
	"time"
)

// GetTestAanmelding retourneert een test aanmelding
func GetTestAanmelding() *models.Aanmelding {
	now := time.Now()
	return &models.Aanmelding{
		ID:             "test-id-123",
		Naam:           "Test Gebruiker",
		Email:          "test@example.com",
		Telefoon:       "0612345678",
		Rol:            "chauffeur",
		Afstand:        "10 km",
		Ondersteuning:  "Geen",
		Bijzonderheden: "Geen",
		Terms:          true,
		EmailVerzonden: false,
		EmailVerzondOp: nil,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// GetTestAanmeldingMetEmail retourneert een test aanmelding met verzonden email
func GetTestAanmeldingMetEmail() *models.Aanmelding {
	aanmelding := GetTestAanmelding()
	now := time.Now()
	aanmelding.EmailVerzonden = true
	aanmelding.EmailVerzondOp = &now
	return aanmelding
}

// GetTestAanmeldingFormulier retourneert een test aanmeldingsformulier
func GetTestAanmeldingFormulier() *models.AanmeldingFormulier {
	return &models.AanmeldingFormulier{
		Naam:           "Test Gebruiker",
		Email:          "test@example.com",
		Telefoon:       "0612345678",
		Rol:            "chauffeur",
		Afstand:        "10 km",
		Ondersteuning:  "Geen",
		Bijzonderheden: "Geen",
		Terms:          true,
	}
}
