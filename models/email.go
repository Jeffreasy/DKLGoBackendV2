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
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
	Content     []byte `json:"content"`
	ContentID   string `json:"content_id"` // For inline images
}

// Email represents a processed email message
type Email struct {
	ID          string            `json:"id"`
	Sender      string            `json:"sender"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	HTML        string            `json:"html"`
	Account     string            `json:"account"`
	MessageID   string            `json:"message_id"`
	CreatedAt   string            `json:"created_at"`
	Read        bool              `json:"read"`
	Metadata    map[string]string `json:"metadata"`
	To          []string          `json:"to"`
	Cc          []string          `json:"cc"`
	Bcc         []string          `json:"bcc"`
	ReplyTo     []string          `json:"reply_to"`
	InReplyTo   string            `json:"in_reply_to"`
	References  []string          `json:"references"`
	Attachments []EmailAttachment `json:"attachments"`
	Headers     map[string]string `json:"headers"`
}

// EmailFetchOptions bevat opties voor het ophalen van emails
type EmailFetchOptions struct {
	Limit  int   `json:"limit"`
	Offset int   `json:"offset"`
	Read   *bool `json:"read"`
}

// EmailResponse is de response voor email requests
type EmailResponse struct {
	Data  []*Email `json:"data"`
	Error error    `json:"error,omitempty"`
}
