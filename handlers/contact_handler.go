package handlers

import (
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"dklautomationgo/services/email"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	emailService *email.EmailService
	contactRepo  *repository.ContactRepository
}

func NewContactHandler(emailService *email.EmailService, contactRepo *repository.ContactRepository) *ContactHandler {
	return &ContactHandler{
		emailService: emailService,
		contactRepo:  contactRepo,
	}
}

func (h *ContactHandler) HandleContactEmail(c *gin.Context) {
	// Parse simplified contact form data
	var formData struct {
		Naam           string `json:"naam"`
		Email          string `json:"email"`
		Bericht        string `json:"bericht"`
		PrivacyAkkoord bool   `json:"privacy_akkoord"`
	}

	if err := c.BindJSON(&formData); err != nil {
		log.Printf("[HandleContactEmail] Error parsing contact form: %v", err)
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

	// Sla het contactformulier op in de database
	if err := h.contactRepo.Create(&contact); err != nil {
		log.Printf("[HandleContactEmail] Error saving contact form: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save contact form"})
		return
	}
	log.Printf("[HandleContactEmail] Successfully saved contact form with ID: %s", contact.ID)

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Printf("[HandleContactEmail] ADMIN_EMAIL environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin email not configured"})
		return
	}
	log.Printf("[HandleContactEmail] Admin email configured: %s", adminEmail)

	// Stuur email naar admin
	adminEmailData := &models.ContactEmailData{
		ToAdmin:    true,
		Contact:    &contact,
		AdminEmail: adminEmail,
	}
	log.Printf("[HandleContactEmail] Sending admin email to: %s", adminEmail)
	if err := h.emailService.SendContactEmail(adminEmailData); err != nil {
		log.Printf("[HandleContactEmail] Error sending admin email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin notification"})
		return
	}
	log.Printf("[HandleContactEmail] Successfully sent admin email")

	// Stuur bevestigingsemail naar gebruiker
	userEmailData := &models.ContactEmailData{
		ToAdmin: false,
		Contact: &contact,
	}
	log.Printf("[HandleContactEmail] Sending confirmation email to: %s", contact.Email)

	// Probeer de gebruikersmail te verzenden, maar ga door zelfs als het mislukt
	var userEmailError error
	if err := h.emailService.SendContactEmail(userEmailData); err != nil {
		log.Printf("[HandleContactEmail] Error sending user email: %v", err)
		userEmailError = err
	} else {
		log.Printf("[HandleContactEmail] Successfully sent confirmation email")

		// Update de database dat de email is verzonden
		if err := h.contactRepo.MarkEmailSent(contact.ID); err != nil {
			log.Printf("[HandleContactEmail] Error marking email as sent in database: %v", err)
		}
	}

	// Bepaal de juiste respons op basis van of de gebruikersmail is verzonden
	if userEmailError != nil {
		// Als we in ontwikkelingsmodus zijn of als het een testdomein is, beschouwen we het als een succes
		if strings.Contains(contact.Email, "@example.com") ||
			strings.Contains(contact.Email, "@test.com") ||
			os.Getenv("DEV_MODE") == "true" ||
			os.Getenv("DEV_MODE") == "1" {
			log.Printf("[HandleContactEmail] Ignoring email error for test domain or in dev mode")
			c.JSON(http.StatusOK, gin.H{
				"message": "Contact form submitted successfully. Admin notification sent. User email simulated.",
				"warning": "User email would normally be sent, but was simulated for test domain.",
				"id":      contact.ID,
			})
			return
		}

		// In productie geven we een foutmelding terug, maar het contactformulier is wel verwerkt
		c.JSON(http.StatusOK, gin.H{
			"message": "Contact form submitted successfully. Admin notification sent.",
			"warning": "Could not send confirmation email to user.",
			"id":      contact.ID,
		})
		return
	}

	// Alles is succesvol
	c.JSON(http.StatusOK, gin.H{
		"message": "Contact form submitted successfully. Confirmation emails sent.",
		"id":      contact.ID,
	})
}

// GetContacts haalt contactformulieren op (admin endpoint)
func (h *ContactHandler) GetContacts(c *gin.Context) {
	// Parse query parameters
	limit := 10
	offset := 0
	status := c.Query("status")

	if limitParam := c.Query("limit"); limitParam != "" {
		if _, err := fmt.Sscanf(limitParam, "%d", &limit); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
	}

	if offsetParam := c.Query("offset"); offsetParam != "" {
		if _, err := fmt.Sscanf(offsetParam, "%d", &offset); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
	}

	var contacts []*models.ContactFormulier
	var err error
	var total int64

	// Haal contacten op basis van status of alle contacten
	if status != "" {
		contacts, err = h.contactRepo.FindByStatus(status, limit, offset)
		if err != nil {
			log.Printf("[GetContacts] Error fetching contacts by status: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
			return
		}
		total, err = h.contactRepo.CountByStatus(status)
	} else {
		contacts, err = h.contactRepo.FindAll(limit, offset)
		if err != nil {
			log.Printf("[GetContacts] Error fetching all contacts: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch contacts"})
			return
		}
		total, err = h.contactRepo.Count()
	}

	if err != nil {
		log.Printf("[GetContacts] Error counting contacts: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count contacts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     contacts,
		"total":    total,
		"has_more": (int64(offset) + int64(len(contacts))) < total,
	})
}

// UpdateContactStatus werkt de status van een contactformulier bij
func (h *ContactHandler) UpdateContactStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Contact ID is required"})
		return
	}

	var updateData struct {
		Status        string `json:"status" binding:"required"`
		BehandeldDoor string `json:"behandeld_door" binding:"required"`
	}

	if err := c.BindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.contactRepo.UpdateStatus(id, updateData.Status, updateData.BehandeldDoor); err != nil {
		log.Printf("[UpdateContactStatus] Error updating contact status: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update contact status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
