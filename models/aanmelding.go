package models

import (
	"time"

	"gorm.io/gorm"
)

// Aanmelding representeert een vrijwilliger aanmelding in de database
type Aanmelding struct {
	ID             string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"` // Unieke identifier
	CreatedAt      time.Time  `json:"created_at" gorm:"not null"`                                // Tijdstip van aanmaken
	UpdatedAt      time.Time  `json:"updated_at" gorm:"not null"`                                // Tijdstip van laatste update
	Naam           string     `json:"naam" gorm:"not null" validate:"required,min=2,max=100"`    // Naam van de vrijwilliger
	Email          string     `json:"email" gorm:"not null" validate:"required,email"`           // Email adres voor communicatie
	Telefoon       string     `json:"telefoon" gorm:"not null" validate:"required"`              // Telefoonnummer
	Rol            string     `json:"rol" gorm:"not null" validate:"required"`                   // Gewenste rol (bijv. chauffeur, bijrijder)
	Afstand        string     `json:"afstand" gorm:"not null" validate:"required"`               // Maximale reisafstand
	Ondersteuning  string     `json:"ondersteuning"`                                             // Benodigde ondersteuning
	Bijzonderheden string     `json:"bijzonderheden"`                                            // Eventuele bijzonderheden
	Terms          bool       `json:"terms" gorm:"not null" validate:"required"`                 // Akkoord met voorwaarden
	EmailVerzonden bool       `json:"email_verzonden" gorm:"default:false"`                      // Of de bevestigingsemail is verzonden
	EmailVerzondOp *time.Time `json:"email_verzonden_op"`                                        // Wanneer de email is verzonden
}

// AanmeldingFormulier representeert het aanmeldingsformulier zoals ontvangen van de frontend
type AanmeldingFormulier struct {
	Naam           string `json:"naam" validate:"required,min=2,max=100"` // Naam van de vrijwilliger
	Email          string `json:"email" validate:"required,email"`        // Email adres
	Telefoon       string `json:"telefoon" validate:"required"`           // Telefoonnummer
	Rol            string `json:"rol" validate:"required"`                // Gewenste rol
	Afstand        string `json:"afstand" validate:"required"`            // Maximale reisafstand
	Ondersteuning  string `json:"ondersteuning"`                          // Benodigde ondersteuning
	Bijzonderheden string `json:"bijzonderheden"`                         // Eventuele bijzonderheden
	Terms          bool   `json:"terms" validate:"required"`              // Akkoord met voorwaarden
}

// TableName override voor GORM
func (Aanmelding) TableName() string {
	return "aanmeldingen"
}

// BeforeCreate wordt aangeroepen voor het aanmaken van een nieuw record
func (a *Aanmelding) BeforeCreate(tx *gorm.DB) error {
	if !a.EmailVerzonden {
		a.EmailVerzondOp = nil
	}
	return nil
}

// ToDatabase converteert een formulier naar een database model
func (f *AanmeldingFormulier) ToDatabase() *Aanmelding {
	return &Aanmelding{
		Naam:           f.Naam,
		Email:          f.Email,
		Telefoon:       f.Telefoon,
		Rol:            f.Rol,
		Afstand:        f.Afstand,
		Ondersteuning:  f.Ondersteuning,
		Bijzonderheden: f.Bijzonderheden,
		Terms:          f.Terms,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
