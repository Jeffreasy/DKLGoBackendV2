package email

import (
	"crypto/tls"
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

// IEmailService definieert de interface voor email services
type IEmailService interface {
	SendAanmeldingEmail(data *models.AanmeldingEmailData) error
	SendContactEmail(data *models.ContactEmailData) error
}

// Controleer of EmailService de IEmailService interface implementeert
var _ IEmailService = (*EmailService)(nil)

type EmailService struct {
	templates     map[string]*template.Template
	config        *ServiceConfig
	accountCaches map[string]*AccountCache
}

func NewEmailService() (*EmailService, error) {
	templates := make(map[string]*template.Template)

	// Get the current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %v", err)
	}
	log.Printf("[NewEmailService] Current working directory: %s", cwd)

	// Load contact email templates
	contactAdminTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/contact_admin_email.html", cwd))
	if err != nil {
		log.Printf("[NewEmailService] Failed to parse contact admin template: %v", err)
		return nil, fmt.Errorf("failed to parse contact admin template: %v", err)
	}
	templates["contact_admin_email.html"] = contactAdminTemplate
	log.Printf("[NewEmailService] Successfully loaded contact_admin_email.html template")

	contactUserTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/contact_email.html", cwd))
	if err != nil {
		log.Printf("[NewEmailService] Failed to parse contact user template: %v", err)
		return nil, fmt.Errorf("failed to parse contact user template: %v", err)
	}
	templates["contact_email.html"] = contactUserTemplate
	log.Printf("[NewEmailService] Successfully loaded contact_email.html template")

	// Load aanmelding email templates
	aanmeldingAdminTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/aanmelding_admin_email.html", cwd))
	if err != nil {
		log.Printf("[NewEmailService] Failed to parse aanmelding admin template: %v", err)
		return nil, fmt.Errorf("failed to parse aanmelding admin template: %v", err)
	}
	templates["aanmelding_admin_email.html"] = aanmeldingAdminTemplate
	log.Printf("[NewEmailService] Successfully loaded aanmelding_admin_email.html template")

	aanmeldingUserTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/aanmelding_email.html", cwd))
	if err != nil {
		log.Printf("[NewEmailService] Failed to parse aanmelding user template: %v", err)
		return nil, fmt.Errorf("failed to parse aanmelding user template: %v", err)
	}
	templates["aanmelding_email.html"] = aanmeldingUserTemplate
	log.Printf("[NewEmailService] Successfully loaded aanmelding_email.html template")

	// Get configuration
	config := GetDefaultConfig()
	log.Printf("[NewEmailService] Loaded email configuration with %d accounts", len(config.Accounts))

	// Initialize cache for each account
	accountCaches := make(map[string]*AccountCache)
	for accountName := range config.Accounts {
		accountCaches[accountName] = NewAccountCache()
		log.Printf("[NewEmailService] Initialized cache for account: %s", accountName)
	}

	return &EmailService{
		templates:     templates,
		config:        config,
		accountCaches: accountCaches,
	}, nil
}

// MarkEmailAsRead marks an email as read in the IMAP server and updates the cache
func (s *EmailService) MarkEmailAsRead(emailID string) error {
	// Parse the email ID to get account name and message number
	parts := strings.Split(emailID, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid email ID format")
	}
	accountName := parts[0]
	messageNum, err := strconv.ParseUint(parts[1], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid message number: %w", err)
	}

	// Get account config
	config, ok := s.config.Accounts[accountName]
	if !ok {
		return fmt.Errorf("unknown account: %s", accountName)
	}

	// Connect to IMAP server
	tlsConfig := &tls.Config{
		ServerName:         config.IMAPHost,
		InsecureSkipVerify: true,
		MinVersion:         tls.VersionTLS12,
	}

	c, err := client.DialTLS(fmt.Sprintf("%s:%d", config.IMAPHost, config.IMAPPort), tlsConfig)
	if err != nil {
		return fmt.Errorf("IMAP connection failed: %w", err)
	}
	defer c.Logout()

	// Login
	if err := c.Login(config.Email, config.Password); err != nil {
		return fmt.Errorf("IMAP login failed: %w", err)
	}

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		return fmt.Errorf("IMAP select inbox failed: %w", err)
	}

	// Create sequence set for the message
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uint32(messageNum))

	// Add the \Seen flag
	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.SeenFlag}
	err = c.Store(seqSet, item, flags, nil)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	// Update cache if enabled
	if s.config.Cache.Enabled {
		cache := s.accountCaches[accountName]
		cache.cacheMutex.Lock()
		for _, email := range cache.emails {
			if email.ID == emailID {
				email.Read = true
				break
			}
		}
		cache.cacheMutex.Unlock()
	}

	return nil
}
