package handlers

import (
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// IAanmeldingHandler definieert de interface voor aanmelding handlers
type IAanmeldingHandler interface {
	CreateAanmelding(c *gin.Context)
	GetAanmeldingen(c *gin.Context)
	GetAanmeldingByID(c *gin.Context)
	UpdateAanmelding(c *gin.Context)
	DeleteAanmelding(c *gin.Context)
}

// Controleer of AanmeldingHandler de IAanmeldingHandler interface implementeert
var _ IAanmeldingHandler = (*AanmeldingHandler)(nil)

// AanmeldingHandler bevat handlers voor aanmeldingen
type AanmeldingHandler struct {
	service services.IAanmeldingService
}

// NewAanmeldingHandler maakt een nieuwe AanmeldingHandler
func NewAanmeldingHandler(service services.IAanmeldingService) *AanmeldingHandler {
	return &AanmeldingHandler{
		service: service,
	}
}

// CreateAanmelding handelt het aanmaken van een aanmelding af
func (h *AanmeldingHandler) CreateAanmelding(c *gin.Context) {
	var aanmelding models.Aanmelding
	if err := c.ShouldBindJSON(&aanmelding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	if err := h.service.CreateAanmelding(&aanmelding); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Aanmelding succesvol aangemaakt",
		"aanmelding": aanmelding,
	})
}

// GetAanmeldingen handelt het ophalen van aanmeldingen af
func (h *AanmeldingHandler) GetAanmeldingen(c *gin.Context) {
	// Parse query parameters
	params := repository.NewQueryParams()

	// Paginering
	if page, err := strconv.Atoi(c.DefaultQuery("page", "1")); err == nil {
		params.WithPage(page)
	}

	if pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", "10")); err == nil {
		params.WithPageSize(pageSize)
	}

	// Sortering
	if sortField := c.Query("sort_field"); sortField != "" {
		params.WithSort(sortField, c.DefaultQuery("sort_order", "asc"))
	}

	// Zoeken
	if search := c.Query("search"); search != "" {
		params.WithSearch(search)
	}

	// Haal aanmeldingen op
	aanmeldingen, err := h.service.GetAanmeldingen(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Tel totaal aantal aanmeldingen
	total, err := h.service.CountAanmeldingen(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"aanmeldingen": aanmeldingen,
		"total":        total,
		"page":         params.Page,
		"page_size":    params.PageSize,
	})
}

// GetAanmeldingByID handelt het ophalen van een aanmelding op basis van ID af
func (h *AanmeldingHandler) GetAanmeldingByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is verplicht"})
		return
	}

	aanmelding, err := h.service.GetAanmeldingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aanmelding niet gevonden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"aanmelding": aanmelding})
}

// UpdateAanmelding handelt het bijwerken van een aanmelding af
func (h *AanmeldingHandler) UpdateAanmelding(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is verplicht"})
		return
	}

	// Controleer of aanmelding bestaat
	aanmelding, err := h.service.GetAanmeldingByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aanmelding niet gevonden"})
		return
	}

	// Bind JSON naar aanmelding
	if err := c.ShouldBindJSON(aanmelding); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	// Werk aanmelding bij
	if err := h.service.UpdateAanmelding(aanmelding); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Aanmelding succesvol bijgewerkt",
		"aanmelding": aanmelding,
	})
}

// DeleteAanmelding handelt het verwijderen van een aanmelding af
func (h *AanmeldingHandler) DeleteAanmelding(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is verplicht"})
		return
	}

	if err := h.service.DeleteAanmelding(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Aanmelding niet gevonden"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Aanmelding succesvol verwijderd"})
}
