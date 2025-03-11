// database/repository/aanmelding_repository.go
package repository

import (
	"dklautomationgo/models"
	"time"

	"gorm.io/gorm"
)

// AanmeldingRepository bevat methoden voor het werken met aanmeldingen in de database
type AanmeldingRepository struct {
	db *gorm.DB
}

// NewAanmeldingRepository maakt een nieuwe AanmeldingRepository
func NewAanmeldingRepository(db *gorm.DB) *AanmeldingRepository {
	return &AanmeldingRepository{db: db}
}

// Create slaat een nieuwe aanmelding op in de database
func (r *AanmeldingRepository) Create(aanmelding *models.Aanmelding) error {
	return r.db.Create(aanmelding).Error
}

// FindByID zoekt een aanmelding op basis van ID
func (r *AanmeldingRepository) FindByID(id string) (*models.Aanmelding, error) {
	var aanmelding models.Aanmelding
	err := r.db.Where("id = ?", id).First(&aanmelding).Error
	return &aanmelding, err
}

// FindAll haalt alle aanmeldingen op
func (r *AanmeldingRepository) FindAll(limit, offset int) ([]*models.Aanmelding, error) {
	var aanmeldingen []*models.Aanmelding
	err := r.db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&aanmeldingen).Error
	return aanmeldingen, err
}

// FindByRol haalt aanmeldingen op basis van rol op
func (r *AanmeldingRepository) FindByRol(rol string, limit, offset int) ([]*models.Aanmelding, error) {
	var aanmeldingen []*models.Aanmelding
	err := r.db.Where("rol = ?", rol).Order("created_at DESC").Limit(limit).Offset(offset).Find(&aanmeldingen).Error
	return aanmeldingen, err
}

// FindByAfstand haalt aanmeldingen op basis van afstand op
func (r *AanmeldingRepository) FindByAfstand(afstand string, limit, offset int) ([]*models.Aanmelding, error) {
	var aanmeldingen []*models.Aanmelding
	err := r.db.Where("afstand = ?", afstand).Order("created_at DESC").Limit(limit).Offset(offset).Find(&aanmeldingen).Error
	return aanmeldingen, err
}

// Update werkt een aanmelding bij
func (r *AanmeldingRepository) Update(aanmelding *models.Aanmelding) error {
	aanmelding.UpdatedAt = time.Now()
	return r.db.Save(aanmelding).Error
}

// MarkEmailSent markeert een aanmelding als verzonden
func (r *AanmeldingRepository) MarkEmailSent(id string) error {
	now := time.Now()
	return r.db.Model(&models.Aanmelding{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"email_verzonden":    true,
			"email_verzonden_op": now,
			"updated_at":         now,
		}).Error
}

// Count telt het aantal aanmeldingen
func (r *AanmeldingRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Aanmelding{}).Count(&count).Error
	return count, err
}

// CountByRol telt het aantal aanmeldingen per rol
func (r *AanmeldingRepository) CountByRol(rol string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Aanmelding{}).Where("rol = ?", rol).Count(&count).Error
	return count, err
}

// CountByAfstand telt het aantal aanmeldingen per afstand
func (r *AanmeldingRepository) CountByAfstand(afstand string) (int64, error) {
	var count int64
	err := r.db.Model(&models.Aanmelding{}).Where("afstand = ?", afstand).Count(&count).Error
	return count, err
}
