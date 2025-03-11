package models

// ContactEmailData bevat de data nodig voor het versturen van contact formulier emails
type ContactEmailData struct {
	ToAdmin    bool              `json:"to_admin"`
	Contact    *ContactFormulier `json:"contact"`
	AdminEmail string            `json:"admin_email,omitempty"`
}

// AanmeldingEmailData bevat de data nodig voor het versturen van aanmelding emails
type AanmeldingEmailData struct {
	ToAdmin    bool                 `json:"to_admin"`
	Aanmelding *AanmeldingFormulier `json:"aanmelding"`
	AdminEmail string               `json:"admin_email,omitempty"`
}

// EmailAttachment represents an email attachment or inline image
type EmailAttachment struct {
	Filename    string `json:"filename"`     // Naam van het bestand
	ContentType string `json:"content_type"` // MIME type van de attachment
	Size        int64  `json:"size"`         // Grootte in bytes
	Content     []byte `json:"content"`      // De binary content
	ContentID   string `json:"content_id"`   // Voor inline afbeeldingen
}

// Email represents a processed email message with all its metadata and content
type Email struct {
	ID          string            `json:"id"`          // Unique identifier (meestal IMAP UID)
	Sender      string            `json:"sender"`      // Email adres van de verzender
	Subject     string            `json:"subject"`     // Onderwerp van de email
	Body        string            `json:"body"`        // Platte tekst versie
	HTML        string            `json:"html"`        // HTML versie (indien beschikbaar)
	Account     string            `json:"account"`     // Email account waar dit bericht bij hoort
	MessageID   string            `json:"message_id"`  // Originele Message-ID header
	CreatedAt   string            `json:"created_at"`  // Timestamp in RFC3339 formaat
	Read        bool              `json:"read"`        // Of de email als gelezen is gemarkeerd
	Metadata    map[string]string `json:"metadata"`    // Extra metadata velden
	To          []string          `json:"to"`          // Lijst van ontvangers
	Cc          []string          `json:"cc"`          // Carbon copy ontvangers
	Bcc         []string          `json:"bcc"`         // Blind carbon copy ontvangers
	ReplyTo     []string          `json:"reply_to"`    // Reply-To adressen
	InReplyTo   string            `json:"in_reply_to"` // Message-ID van de email waar dit een antwoord op is
	References  []string          `json:"references"`  // Gerelateerde Message-IDs
	Attachments []EmailAttachment `json:"attachments"` // Lijst van bijlagen
	Headers     map[string]string `json:"headers"`     // Alle email headers
}

// EmailFetchOptions bevat de parameters voor het ophalen van emails
type EmailFetchOptions struct {
	Limit  int   `json:"limit"`  // Maximum aantal emails om op te halen
	Offset int   `json:"offset"` // Aantal emails om over te slaan (voor paginatie)
	Read   *bool `json:"read"`   // Filter op gelezen/ongelezen status
}

// EmailResponse is de gestandaardiseerde response voor email requests
type EmailResponse struct {
	Data    []*Email `json:"data"`            // Lijst van emails
	Error   string   `json:"error,omitempty"` // Eventuele foutmelding
	Total   int      `json:"total"`           // Totaal aantal beschikbare emails
	HasMore bool     `json:"has_more"`        // Of er meer emails beschikbaar zijn
}
