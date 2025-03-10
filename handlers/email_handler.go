package handlers

import (
	"dklautomationgo/models"
	"dklautomationgo/services/email"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *email.EmailService
}

func NewEmailHandler(emailService *email.EmailService) *EmailHandler {
	return &EmailHandler{
		emailService: emailService,
	}
}

// GetEmails handles GET /api/emails
func (h *EmailHandler) GetEmails(c *gin.Context) {
	// Parse query parameters
	options := &models.EmailFetchOptions{}

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			log.Printf("Invalid limit parameter: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		options.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			log.Printf("Invalid offset parameter: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
		options.Offset = offset
	}

	if readStr := c.Query("read"); readStr != "" {
		read, err := strconv.ParseBool(readStr)
		if err != nil {
			log.Printf("Invalid read parameter: %v", err)
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid read parameter"})
			return
		}
		options.Read = &read
	}

	// Fetch emails
	emails, err := h.emailService.FetchEmails(options)
	if err != nil {
		log.Printf("Error fetching emails: %v", err)
		if strings.Contains(err.Error(), "authentication failed") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email authentication failed"})
			return
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial tcp") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Could not connect to email server"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch emails: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": emails})
}

// GetEmailStats handles GET /api/emails/stats
func (h *EmailHandler) GetEmailStats(c *gin.Context) {
	// Get all emails to count
	emails, err := h.emailService.FetchEmails(nil)
	if err != nil {
		log.Printf("Error fetching emails for stats: %v", err)
		if strings.Contains(err.Error(), "authentication failed") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email authentication failed"})
			return
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial tcp") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Could not connect to email server"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch email stats: %v", err)})
		return
	}

	// Count total and unread
	total := len(emails)
	unread := 0
	for _, email := range emails {
		if !email.Read {
			unread++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"total":  total,
		"unread": unread,
	})
}

// MarkEmailAsRead handles PUT /api/emails/:id/read
func (h *EmailHandler) MarkEmailAsRead(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email ID is required"})
		return
	}

	// TODO: Implement marking email as read in IMAP
	// For now, we'll just return success
	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *EmailHandler) HandleContactEmail(c *gin.Context) {
	// Parse simplified contact form data
	var formData struct {
		Naam           string `json:"naam"`
		Email          string `json:"email"`
		Bericht        string `json:"bericht"`
		PrivacyAkkoord bool   `json:"privacy_akkoord"`
	}

	if err := c.BindJSON(&formData); err != nil {
		log.Printf("Error parsing contact form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Create full contact form data
	contact := models.ContactFormulier{
		Naam:           formData.Naam,
		Email:          formData.Email,
		Bericht:        formData.Bericht,
		PrivacyAkkoord: formData.PrivacyAkkoord,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Status:         "nieuw",
		EmailVerzonden: false,
	}

	log.Printf("Sending admin email to: %s", os.Getenv("ADMIN_EMAIL"))

	// Stuur email naar admin
	adminEmailData := &models.ContactEmailData{
		ToAdmin:    true,
		Contact:    &contact,
		AdminEmail: os.Getenv("ADMIN_EMAIL"),
	}
	if err := h.emailService.SendContactEmail(adminEmailData); err != nil {
		log.Printf("Error sending admin email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin notification"})
		return
	}

	log.Printf("Successfully sent admin email to: %s", os.Getenv("ADMIN_EMAIL"))

	// Stuur bevestigingsemail naar gebruiker
	userEmailData := &models.ContactEmailData{
		ToAdmin: false,
		Contact: &contact,
	}
	if err := h.emailService.SendContactEmail(userEmailData); err != nil {
		log.Printf("Error sending user email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}

	log.Printf("Successfully sent confirmation email to: %s", contact.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
}

func (h *EmailHandler) HandleAanmeldingEmail(c *gin.Context) {
	var aanmelding models.AanmeldingFormulier
	if err := c.BindJSON(&aanmelding); err != nil {
		log.Printf("Error parsing aanmelding form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	log.Printf("Sending admin email to: %s", os.Getenv("ADMIN_EMAIL"))

	// Stuur email naar admin
	adminEmailData := &models.AanmeldingEmailData{
		ToAdmin:    true,
		Aanmelding: &aanmelding,
		AdminEmail: os.Getenv("ADMIN_EMAIL"),
	}
	if err := h.emailService.SendAanmeldingEmail(adminEmailData); err != nil {
		log.Printf("Error sending admin email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin notification"})
		return
	}

	log.Printf("Successfully sent admin email to: %s", os.Getenv("ADMIN_EMAIL"))

	// Stuur bevestigingsemail naar gebruiker
	userEmailData := &models.AanmeldingEmailData{
		ToAdmin:    false,
		Aanmelding: &aanmelding,
	}
	if err := h.emailService.SendAanmeldingEmail(userEmailData); err != nil {
		log.Printf("Error sending user email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send confirmation email"})
		return
	}

	log.Printf("Successfully sent confirmation email to: %s", aanmelding.Email)

	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
}
