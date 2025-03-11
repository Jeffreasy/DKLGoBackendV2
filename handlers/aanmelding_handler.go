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

	"github.com/gin-gonic/gin"
)

type AanmeldingHandler struct {
	emailService   *email.EmailService
	aanmeldingRepo *repository.AanmeldingRepository
}

func NewAanmeldingHandler(emailService *email.EmailService, aanmeldingRepo *repository.AanmeldingRepository) *AanmeldingHandler {
	return &AanmeldingHandler{
		emailService:   emailService,
		aanmeldingRepo: aanmeldingRepo,
	}
}

func (h *AanmeldingHandler) HandleAanmeldingEmail(c *gin.Context) {
	var aanmeldingForm models.AanmeldingFormulier
	if err := c.BindJSON(&aanmeldingForm); err != nil {
		log.Printf("[HandleAanmeldingEmail] Error parsing aanmelding form: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	log.Printf("[HandleAanmeldingEmail] Successfully parsed aanmelding form for %s", aanmeldingForm.Naam)

	// Converteer naar database model
	aanmelding := aanmeldingForm.ToDatabase()

	// Sla de aanmelding op in de database
	if err := h.aanmeldingRepo.Create(aanmelding); err != nil {
		log.Printf("[HandleAanmeldingEmail] Error saving aanmelding: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save registration"})
		return
	}
	log.Printf("[HandleAanmeldingEmail] Successfully saved aanmelding with ID: %s", aanmelding.ID)

	adminEmail := os.Getenv("ADMIN_EMAIL")
	if adminEmail == "" {
		log.Printf("[HandleAanmeldingEmail] ADMIN_EMAIL environment variable not set")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Admin email not configured"})
		return
	}
	log.Printf("[HandleAanmeldingEmail] Admin email configured: %s", adminEmail)

	// Stuur email naar admin
	adminEmailData := &models.AanmeldingEmailData{
		ToAdmin:    true,
		Aanmelding: &aanmeldingForm,
		AdminEmail: adminEmail,
	}
	log.Printf("[HandleAanmeldingEmail] Sending admin email to: %s", adminEmail)
	if err := h.emailService.SendAanmeldingEmail(adminEmailData); err != nil {
		log.Printf("[HandleAanmeldingEmail] Error sending admin email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send admin notification"})
		return
	}
	log.Printf("[HandleAanmeldingEmail] Successfully sent admin email")

	// Stuur bevestigingsemail naar gebruiker
	userEmailData := &models.AanmeldingEmailData{
		ToAdmin:    false,
		Aanmelding: &aanmeldingForm,
	}
	log.Printf("[HandleAanmeldingEmail] Sending confirmation email to: %s", aanmeldingForm.Email)

	// Probeer de gebruikersmail te verzenden, maar ga door zelfs als het mislukt
	var userEmailError error
	if err := h.emailService.SendAanmeldingEmail(userEmailData); err != nil {
		log.Printf("[HandleAanmeldingEmail] Error sending user email: %v", err)
		userEmailError = err
	} else {
		log.Printf("[HandleAanmeldingEmail] Successfully sent confirmation email")

		// Update de database dat de email is verzonden
		if err := h.aanmeldingRepo.MarkEmailSent(aanmelding.ID); err != nil {
			log.Printf("[HandleAanmeldingEmail] Error marking email as sent in database: %v", err)
		}
	}

	// Bepaal de juiste respons op basis van of de gebruikersmail is verzonden
	if userEmailError != nil {
		// Als we in ontwikkelingsmodus zijn of als het een testdomein is, beschouwen we het als een succes
		if strings.Contains(aanmeldingForm.Email, "@example.com") ||
			strings.Contains(aanmeldingForm.Email, "@test.com") ||
			os.Getenv("DEV_MODE") == "true" ||
			os.Getenv("DEV_MODE") == "1" {
			log.Printf("[HandleAanmeldingEmail] Ignoring email error for test domain or in dev mode")
			c.JSON(http.StatusOK, gin.H{
				"message": "Registration successful. Admin notification sent. User email simulated.",
				"warning": "User email would normally be sent, but was simulated for test domain.",
				"id":      aanmelding.ID,
			})
			return
		}

		// In productie geven we een foutmelding terug, maar de aanmelding is wel verwerkt
		c.JSON(http.StatusOK, gin.H{
			"message": "Registration successful. Admin notification sent.",
			"warning": "Could not send confirmation email to user.",
			"id":      aanmelding.ID,
		})
		return
	}

	// Alles is succesvol
	c.JSON(http.StatusOK, gin.H{
		"message": "Registration successful. Confirmation emails sent.",
		"id":      aanmelding.ID,
	})
}

// GetAanmeldingen haalt aanmeldingen op (admin endpoint)
func (h *AanmeldingHandler) GetAanmeldingen(c *gin.Context) {
	// Parse query parameters
	limit := 10
	offset := 0
	rol := c.Query("rol")
	afstand := c.Query("afstand")

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

	var aanmeldingen []*models.Aanmelding
	var err error
	var total int64

	// Haal aanmeldingen op basis van rol of afstand
	if rol != "" {
		aanmeldingen, err = h.aanmeldingRepo.FindByRol(rol, limit, offset)
		if err != nil {
			log.Printf("[GetAanmeldingen] Error fetching aanmeldingen by rol: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registrations"})
			return
		}
		total, err = h.aanmeldingRepo.CountByRol(rol)
	} else if afstand != "" {
		aanmeldingen, err = h.aanmeldingRepo.FindByAfstand(afstand, limit, offset)
		if err != nil {
			log.Printf("[GetAanmeldingen] Error fetching aanmeldingen by afstand: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registrations"})
			return
		}
		total, err = h.aanmeldingRepo.CountByAfstand(afstand)
	} else {
		aanmeldingen, err = h.aanmeldingRepo.FindAll(limit, offset)
		if err != nil {
			log.Printf("[GetAanmeldingen] Error fetching all aanmeldingen: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch registrations"})
			return
		}
		total, err = h.aanmeldingRepo.Count()
	}

	if err != nil {
		log.Printf("[GetAanmeldingen] Error counting aanmeldingen: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count registrations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     aanmeldingen,
		"total":    total,
		"has_more": (int64(offset) + int64(len(aanmeldingen))) < total,
	})
}

// GetAanmeldingByID haalt een specifieke aanmelding op
func (h *AanmeldingHandler) GetAanmeldingByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Registration ID is required"})
		return
	}

	aanmelding, err := h.aanmeldingRepo.FindByID(id)
	if err != nil {
		log.Printf("[GetAanmeldingByID] Error fetching aanmelding: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Registration not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": aanmelding})
}

// GetAanmeldingStats haalt statistieken op over aanmeldingen
func (h *AanmeldingHandler) GetAanmeldingStats(c *gin.Context) {
	// Totaal aantal aanmeldingen
	total, err := h.aanmeldingRepo.Count()
	if err != nil {
		log.Printf("[GetAanmeldingStats] Error counting aanmeldingen: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count registrations"})
		return
	}

	// Haal unieke rollen op
	rollen := []string{"Deelnemer", "Vrijwilliger", "Chauffeur", "Bijrijder", "Verzorging"}
	rolStats := make(map[string]int64)

	for _, rol := range rollen {
		count, err := h.aanmeldingRepo.CountByRol(rol)
		if err != nil {
			log.Printf("[GetAanmeldingStats] Error counting aanmeldingen by rol %s: %v", rol, err)
			continue
		}
		rolStats[rol] = count
	}

	// Haal unieke afstanden op
	afstanden := []string{"2.5 KM", "5 KM", "10 KM", "15 KM", "Halve marathon"}
	afstandStats := make(map[string]int64)

	for _, afstand := range afstanden {
		count, err := h.aanmeldingRepo.CountByAfstand(afstand)
		if err != nil {
			log.Printf("[GetAanmeldingStats] Error counting aanmeldingen by afstand %s: %v", afstand, err)
			continue
		}
		afstandStats[afstand] = count
	}

	c.JSON(http.StatusOK, gin.H{
		"total":      total,
		"by_rol":     rolStats,
		"by_afstand": afstandStats,
	})
}
