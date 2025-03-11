package email

import (
	"bytes"
	"crypto/tls"
	"dklautomationgo/models"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
)

func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error {
	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "contact_admin_email.html"
		subject = "Nieuw contactformulier ontvangen"
		recipient = data.AdminEmail
		log.Printf("Sending admin email to: %s using template: %s", recipient, templateName)
	} else {
		templateName = "contact_email.html"
		subject = "Bedankt voor je bericht"
		recipient = data.Contact.Email
		log.Printf("Sending user email to: %s using template: %s", recipient, templateName)
	}

	template := s.templates[templateName]
	if template == nil {
		log.Printf("Template not found: %s", templateName)
		return fmt.Errorf("template not found: %s", templateName)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		log.Printf("Failed to execute template: %v", err)
		return fmt.Errorf("failed to execute template: %v", err)
	}

	log.Printf("Successfully generated email body for template: %s", templateName)
	return s.sendEmail(recipient, subject, body.String())
}

func (s *EmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "aanmelding_admin_email.html"
		subject = "Nieuwe aanmelding ontvangen"
		recipient = data.AdminEmail
		log.Printf("[SendAanmeldingEmail] Preparing admin email - Template: %s, Recipient: %s", templateName, recipient)
	} else {
		templateName = "aanmelding_email.html"
		subject = "Bedankt voor je aanmelding"
		recipient = data.Aanmelding.Email
		log.Printf("[SendAanmeldingEmail] Preparing user email - Template: %s, Recipient: %s", templateName, recipient)
	}

	template := s.templates[templateName]
	if template == nil {
		log.Printf("[SendAanmeldingEmail] Template not found: %s", templateName)
		return fmt.Errorf("template not found: %s", templateName)
	}
	log.Printf("[SendAanmeldingEmail] Found template: %s", templateName)

	// Log template data for debugging
	log.Printf("[SendAanmeldingEmail] Template data: ToAdmin=%v, Naam=%s, Email=%s, Rol=%s, Afstand=%s",
		data.ToAdmin, data.Aanmelding.Naam, data.Aanmelding.Email, data.Aanmelding.Rol, data.Aanmelding.Afstand)

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		log.Printf("[SendAanmeldingEmail] Failed to execute template: %v", err)
		return fmt.Errorf("failed to execute template: %v", err)
	}
	log.Printf("[SendAanmeldingEmail] Successfully executed template, generated body length: %d", body.Len())

	if err := s.sendEmail(recipient, subject, body.String()); err != nil {
		log.Printf("[SendAanmeldingEmail] Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}
	log.Printf("[SendAanmeldingEmail] Successfully sent email to %s", recipient)

	return nil
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	log.Printf("[sendEmail] Starting email send process to: %s with subject: %s", to, subject)

	// Check if we should simulate email delivery
	if s.config.DevMode {
		// In development mode, only send to allowed domains
		allowedDomains := []string{
			"dekoninklijkeloop.nl",
			"localhost",
			"127.0.0.1",
		}

		shouldSimulate := true
		for _, domain := range allowedDomains {
			if strings.HasSuffix(to, "@"+domain) {
				shouldSimulate = false
				break
			}
		}

		if shouldSimulate {
			log.Printf("[sendEmail] DEV MODE: Simulating email delivery to %s", to)
			log.Printf("[sendEmail] DEV MODE: Email subject: %s", subject)
			log.Printf("[sendEmail] DEV MODE: Email body length: %d bytes", len(body))
			return nil
		}
	} else {
		// In production mode, still simulate for obvious test domains
		testDomains := []string{"@example.com", "@test.com", "@example.org"}
		for _, domain := range testDomains {
			if strings.HasSuffix(to, domain) {
				log.Printf("[sendEmail] Detected test email address: %s. Simulating successful delivery.", to)
				return nil
			}
		}
	}

	m := gomail.NewMessage()
	// Use the same email address for From header as the SMTP authentication
	emailConfig := s.config.Accounts["info"]
	if emailConfig == nil {
		log.Printf("[sendEmail] Email configuration not found for info account")
		return fmt.Errorf("email configuration not found for info account")
	}

	m.SetHeader("From", emailConfig.Email) // Use the authenticated email address
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	log.Printf("[sendEmail] Using SMTP Configuration - Host: %s, Port: %d, Username: %s, Password length: %d",
		emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.Email, len(emailConfig.Password))

	d := gomail.NewDialer(emailConfig.SMTPHost, emailConfig.SMTPPort, emailConfig.Email, emailConfig.Password)

	// Configure TLS based on port
	if emailConfig.SMTPPort == 465 {
		// Port 465 uses implicit SSL/TLS
		d.SSL = true
	} else {
		// Port 587 uses STARTTLS
		d.SSL = false
	}

	// TLS configuration
	d.TLSConfig = &tls.Config{
		ServerName:         emailConfig.SMTPHost,
		InsecureSkipVerify: true,             // Allow invalid certificates for testing
		MinVersion:         tls.VersionTLS10, // Allow older TLS versions
	}

	// Add retry logic for transient errors
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		var connectionType string
		if d.SSL {
			connectionType = "SSL"
		} else {
			connectionType = "STARTTLS"
		}

		log.Printf("[sendEmail] Attempt %d/%d: Connecting to SMTP server %s:%d with %s...",
			i+1, maxRetries, emailConfig.SMTPHost, emailConfig.SMTPPort, connectionType)

		if err := d.DialAndSend(m); err != nil {
			log.Printf("[sendEmail] Attempt %d/%d failed: %v", i+1, maxRetries, err)

			// Check if it's a network error
			if netErr, ok := err.(net.Error); ok {
				log.Printf("[sendEmail] Network error details - Type: %T, Timeout: %v, Temporary: %v",
					netErr, netErr.Timeout(), netErr.Temporary())
			}

			// Check if it's a TLS error
			if tlsErr, ok := err.(tls.RecordHeaderError); ok {
				log.Printf("[sendEmail] TLS error details: %v", tlsErr)
			}

			// Check if it's an authentication error
			if strings.Contains(err.Error(), "authentication") {
				log.Printf("[sendEmail] Authentication error detected. Please verify SMTP credentials.")
			}

			if i == maxRetries-1 {
				return fmt.Errorf("failed to send email after %d attempts: %v", maxRetries, err)
			}

			// Exponential backoff with a maximum of 5 seconds
			backoff := time.Duration(i+1) * 2 * time.Second
			if backoff > 5*time.Second {
				backoff = 5 * time.Second
			}
			log.Printf("[sendEmail] Waiting %v before next attempt", backoff)
			time.Sleep(backoff)
			continue
		}

		log.Printf("[sendEmail] Successfully sent email to: %s", to)
		return nil
	}

	return fmt.Errorf("failed to send email after %d retries", maxRetries)
}
