package email

import (
	"os"
	"time"
)

// EmailConfig bevat de configuratie voor een email account
type EmailConfig struct {
	Email    string
	Password string
	IMAPHost string
	IMAPPort int
	SMTPHost string
	SMTPPort int
}

// CacheConfig bevat de configuratie voor email caching
type CacheConfig struct {
	Enabled    bool
	Duration   time.Duration
	MaxEntries int
}

// ServiceConfig bevat alle configuratie voor de email service
type ServiceConfig struct {
	Accounts     map[string]*EmailConfig
	Cache        CacheConfig
	FetchTimeout time.Duration
}

func GetDefaultConfig() *ServiceConfig {
	return &ServiceConfig{
		Accounts: map[string]*EmailConfig{
			"info": {
				Email:    "info@dekoninklijkeloop.nl",
				Password: os.Getenv("INFO_EMAIL_PASSWORD"),
				IMAPHost: "imap.hostnet.nl",
				IMAPPort: 993,
				SMTPHost: "mailout.hostnet.nl",
				SMTPPort: 587,
			},
			"inschrijving": {
				Email:    "inschrijving@dekoninklijkeloop.nl",
				Password: os.Getenv("INSCHRIJVING_EMAIL_PASSWORD"),
				IMAPHost: "imap.hostnet.nl",
				IMAPPort: 993,
				SMTPHost: "mailout.hostnet.nl",
				SMTPPort: 587,
			},
		},
		Cache: CacheConfig{
			Enabled:    true,
			Duration:   5 * time.Minute,
			MaxEntries: 1000,
		},
		FetchTimeout: 2 * time.Minute,
	}
}
