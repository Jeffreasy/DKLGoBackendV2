package handlers

import (
	"dklautomationgo/models"
	"dklautomationgo/services/email"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
			return
		}
		options.Limit = limit
	}

	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
			return
		}
		options.Offset = offset
	}

	if readStr := c.Query("read"); readStr != "" {
		read, err := strconv.ParseBool(readStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid read parameter"})
			return
		}
		options.Read = &read
	}

	// Fetch emails
	emails, err := h.emailService.FetchEmails(options)
	if err != nil {
		log.Printf("[ERROR] Failed to fetch emails: %v", err)
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
		log.Printf("[ERROR] Failed to fetch email stats: %v", err)
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

	err := h.emailService.MarkEmailAsRead(id)
	if err != nil {
		log.Printf("[ERROR] Failed to mark email as read: %v", err)

		if strings.Contains(err.Error(), "invalid email ID") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email ID format"})
			return
		}
		if strings.Contains(err.Error(), "unknown account") {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unknown email account"})
			return
		}
		if strings.Contains(err.Error(), "authentication failed") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email authentication failed"})
			return
		}
		if strings.Contains(err.Error(), "connection refused") || strings.Contains(err.Error(), "dial tcp") {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Could not connect to email server"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to mark email as read: %v", err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
