package email

import (
	"log"
	"os"
	"strconv"
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
	DevMode      bool // Ontwikkelingsmodus voor testen
}

func GetDefaultConfig() *ServiceConfig {
	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "smtp.hostnet.nl"
	}

	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpPort := 587 // Default port for STARTTLS
	if smtpPortStr != "" {
		if port, err := strconv.Atoi(smtpPortStr); err == nil {
			smtpPort = port
		}
	}

	// Check development mode
	devMode := false
	devModeStr := os.Getenv("DEV_MODE")
	if devModeStr == "true" || devModeStr == "1" {
		devMode = true
		log.Printf("[GetDefaultConfig] Running in DEVELOPMENT mode - emails to external domains will be simulated")
	}

	// Log the SMTP configuration
	log.Printf("[GetDefaultConfig] Using SMTP configuration - Host: %s, Port: %d", smtpHost, smtpPort)
	if smtpPort == 465 {
		log.Printf("[GetDefaultConfig] Using implicit SSL/TLS for SMTP")
	} else if smtpPort == 587 {
		log.Printf("[GetDefaultConfig] Using STARTTLS for SMTP")
	} else {
		log.Printf("[GetDefaultConfig] Warning: Unusual SMTP port %d, please verify configuration", smtpPort)
	}

	return &ServiceConfig{
		Accounts: map[string]*EmailConfig{
			"info": {
				Email:    os.Getenv("SMTP_USER"),
				Password: os.Getenv("SMTP_PASSWORD"),
				IMAPHost: "imap.hostnet.nl",
				IMAPPort: 993,
				SMTPHost: smtpHost,
				SMTPPort: smtpPort,
			},
			"inschrijving": {
				Email:    "inschrijving@dekoninklijkeloop.nl",
				Password: os.Getenv("INSCHRIJVING_EMAIL_PASSWORD"),
				IMAPHost: "imap.hostnet.nl",
				IMAPPort: 993,
				SMTPHost: smtpHost,
				SMTPPort: smtpPort,
			},
			"noreply": {
				Email:    "noreply@dekoninklijkeloop.nl",
				Password: os.Getenv("NOREPLY_EMAIL_PASSWORD"),
				IMAPHost: "imap.hostnet.nl",
				IMAPPort: 993,
				SMTPHost: smtpHost,
				SMTPPort: smtpPort,
			},
		},
		Cache: CacheConfig{
			Enabled:    true,
			Duration:   5 * time.Minute,
			MaxEntries: 1000,
		},
		FetchTimeout: 2 * time.Minute,
		DevMode:      devMode,
	}
}
