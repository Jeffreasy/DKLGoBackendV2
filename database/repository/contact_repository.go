// database/repository/contact_repository.go
package repository

import (
	"dklautomationgo/models"
	"time"

	"gorm.io/gorm"
)

// ContactRepository bevat methoden voor het werken met contactformulieren in de database
type ContactRepository struct {
	db *gorm.DB
}

// NewContactRepository maakt een nieuwe ContactRepository
func NewContactRepository(db *gorm.DB) *ContactRepository {
	return &ContactRepository{db: db}
}

// Create slaat een nieuw contactformulier op in de database
func (r *ContactRepository) Create(contact *models.ContactFormulier) error {
	return r.db.Create(contact).Error
}

// FindByID zoekt een contactformulier op basis van ID
func (r *ContactRepository) FindByID(id string) (*models.ContactFormulier, error) {
	var contact models.ContactFormulier
	err := r.db.Where("id = ?", id).First(&contact).Error
	return &contact, err
}

// FindAll haalt alle contactformulieren op
func (r *ContactRepository) FindAll(limit, offset int) ([]*models.ContactFormulier, error) {
	var contacts []*models.ContactFormulier
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&contacts).Error
	return contacts, err
}

// FindByStatus haalt contactformulieren op basis van status op
func (r *ContactRepository) FindByStatus(status string, limit, offset int) ([]*models.ContactFormulier, error) {
	var contacts []*models.ContactFormulier
	err := r.db.Where("status = ?", status).Order("created_at DESC").Limit(limit).Offset(offset).Find(&contacts).Error
	return contacts, err
}

// Update werkt een contactformulier bij
func (r *ContactRepository) Update(contact *models.ContactFormulier) error {
	contact.UpdatedAt = time.Now()
	return r.db.Save(contact).Error
}

// MarkEmailSent markeert een contactformulier als verzonden
func (r *ContactRepository) MarkEmailSent(id string) error {
	now := time.Now()
	return r.db.Model(&models.ContactFormulier{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email_verzonden":    true,
			"email_verzonden_op": now,
			"updated_at":         now,
		}).Error
}

// UpdateStatus werkt de status van een contactformulier bij
func (r *ContactRepository) UpdateStatus(id, status, behandeldDoor string) error {
	now := time.Now()
	return r.db.Model(&models.ContactFormulier{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":         status,
			"behandeld_door": behandeldDoor,
			"behandeld_op":   now,
			"updated_at":     now,
		}).Error
}

// Count telt het aantal contactformulieren
func (r *ContactRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.ContactFormulier{}).Count(&count).Error
	return count, err
}

// CountByStatus telt het aantal contactformulieren per status
func (r *ContactRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&models.ContactFormulier{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
