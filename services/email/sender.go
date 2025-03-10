package email

import (
	"bytes"
	"crypto/tls"
	"dklautomationgo/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopkg.in/gomail.v2"
)

func (s *EmailService) SendContactEmail(data *models.ContactEmailData) error {
	var templateName string
	var subject string
	var recipient string

	if data.ToAdmin {
		templateName = "contact_admin"
		subject = "Nieuw contactformulier ontvangen"
		recipient = data.AdminEmail
		log.Printf("Sending admin email to: %s using template: %s", recipient, templateName)
	} else {
		templateName = "contact_user"
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
		templateName = "aanmelding_admin"
		subject = "Nieuwe aanmelding ontvangen"
		recipient = data.AdminEmail
	} else {
		templateName = "aanmelding_user"
		subject = "Bedankt voor je aanmelding"
		recipient = data.Aanmelding.Email
	}

	template := s.templates[templateName]
	if template == nil {
		return fmt.Errorf("template not found: %s", templateName)
	}

	var body bytes.Buffer
	if err := template.Execute(&body, data); err != nil {
		return fmt.Errorf("failed to execute template: %v", err)
	}

	return s.sendEmail(recipient, subject, body.String())
}

func (s *EmailService) sendEmail(to, subject, body string) error {
	log.Printf("Attempting to send email to: %s with subject: %s", to, subject)

	m := gomail.NewMessage()
	fromEmail := os.Getenv("SMTP_FROM")
	if fromEmail == "" {
		fromEmail = "noreply@dekoninklijkeloop.nl"
	}
	m.SetHeader("From", fromEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		smtpPortStr = "587" // Default port for TLS
	}
	smtpUsername := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	log.Printf("SMTP Configuration - Host: %s, Port: %s, Username: %s, From: %s", smtpHost, smtpPortStr, smtpUsername, fromEmail)

	if smtpHost == "" || smtpUsername == "" || smtpPassword == "" {
		return fmt.Errorf("missing SMTP configuration - Host: %s, Username: %s", smtpHost, smtpUsername)
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Printf("Invalid SMTP port number: %v, using default 587", err)
		smtpPort = 587
	}

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)
	d.TLSConfig = &tls.Config{
		ServerName: smtpHost,
	}

	log.Printf("Attempting to connect to SMTP server: %s:%d with TLS", smtpHost, smtpPort)
	if err := d.DialAndSend(m); err != nil {
		log.Printf("Failed to send email: %v", err)
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Successfully sent email to: %s", to)
	return nil
}
