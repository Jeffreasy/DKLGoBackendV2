package models

import (
	"time"

	"gorm.io/gorm"
)

// ContactFormulier representeert een contact formulier inzending met tracking informatie
type ContactFormulier struct {
	ID               string     `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"` // Unieke identifier
	CreatedAt        time.Time  `json:"created_at" gorm:"not null"`                                // Tijdstip van aanmaken
	UpdatedAt        time.Time  `json:"updated_at" gorm:"not null"`                                // Tijdstip van laatste update
	Naam             string     `json:"naam" gorm:"not null" validate:"required,min=2,max=100"`    // Naam van de contactpersoon
	Email            string     `json:"email" gorm:"not null" validate:"required,email"`           // Email adres voor communicatie
	Bericht          string     `json:"bericht" gorm:"type:text;not null" validate:"required"`     // Het bericht van de gebruiker
	EmailVerzonden   bool       `json:"email_verzonden" gorm:"default:false"`                      // Of de bevestigingsemail is verzonden
	EmailVerzondenOp *time.Time `json:"email_verzonden_op"`                                        // Wanneer de email is verzonden
	PrivacyAkkoord   bool       `json:"privacy_akkoord" gorm:"not null" validate:"required"`       // Of gebruiker akkoord is met privacy voorwaarden
	Status           string     `json:"status" gorm:"not null;default:'nieuw'"`                    // Status van de aanvraag (nieuw/in behandeling/afgerond)
	BehandeldDoor    *string    `json:"behandeld_door"`                                            // Wie de aanvraag heeft behandeld
	BehandeldOp      *time.Time `json:"behandeld_op"`                                              // Wanneer de aanvraag is behandeld
	Notities         *string    `json:"notities" gorm:"type:text"`                                 // Interne notities over de aanvraag
}

// TableName override voor GORM
func (ContactFormulier) TableName() string {
	return "contact_formulieren"
}

// BeforeCreate wordt aangeroepen voor het aanmaken van een nieuw record
func (c *ContactFormulier) BeforeCreate(tx *gorm.DB) error {
	if c.Status == "" {
		c.Status = "nieuw"
	}
	return nil
}
