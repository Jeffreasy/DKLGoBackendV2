package email

import (
	"dklautomationgo/models"
	"fmt"
	"html/template"
	"os"
)

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

	// Load contact email templates
	contactAdminTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/contact_admin_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact admin template: %v", err)
	}
	templates["contact_admin"] = contactAdminTemplate

	contactUserTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/contact_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contact user template: %v", err)
	}
	templates["contact_user"] = contactUserTemplate

	// Load aanmelding email templates
	aanmeldingAdminTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/aanmelding_admin_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse aanmelding admin template: %v", err)
	}
	templates["aanmelding_admin"] = aanmeldingAdminTemplate

	aanmeldingUserTemplate, err := template.ParseFiles(fmt.Sprintf("%s/templates/aanmelding_email.html", cwd))
	if err != nil {
		return nil, fmt.Errorf("failed to parse aanmelding user template: %v", err)
	}
	templates["aanmelding_user"] = aanmeldingUserTemplate

	// Get configuration
	config := GetDefaultConfig()

	// Initialize cache for each account
	accountCaches := make(map[string]*AccountCache)
	for accountName := range config.Accounts {
		accountCaches[accountName] = &AccountCache{
			emails: make([]*models.Email, 0),
		}
	}

	return &EmailService{
		templates:     templates,
		config:        config,
		accountCaches: accountCaches,
	}, nil
}
