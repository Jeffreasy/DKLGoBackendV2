package repository

import (
	"dklautomationgo/models"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handelt database operaties voor gebruikers
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository maakt een nieuwe UserRepository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create maakt een nieuwe gebruiker aan
func (r *UserRepository) Create(user *models.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		log.Printf("[UserRepository] Error creating user: %v", result.Error)
		return result.Error
	}
	return nil
}

// FindByID zoekt een gebruiker op ID
func (r *UserRepository) FindByID(id uuid.UUID) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[UserRepository] Error finding user by ID: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// FindByEmail zoekt een gebruiker op email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[UserRepository] Error finding user by email: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// FindAll haalt alle gebruikers op
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	result := r.db.Find(&users)
	if result.Error != nil {
		log.Printf("[UserRepository] Error finding all users: %v", result.Error)
		return nil, result.Error
	}
	return users, nil
}

// Update werkt een gebruiker bij
func (r *UserRepository) Update(user *models.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		log.Printf("[UserRepository] Error updating user: %v", result.Error)
		return result.Error
	}
	return nil
}

// UpdateLastLogin werkt het laatste login tijdstip bij
func (r *UserRepository) UpdateLastLogin(id uuid.UUID) error {
	now := time.Now()
	result := r.db.Model(&models.User{}).Where("id = ?", id).Update("last_login", now)
	if result.Error != nil {
		log.Printf("[UserRepository] Error updating last login: %v", result.Error)
		return result.Error
	}
	return nil
}

// SetPasswordResetToken stelt een wachtwoord reset token in
func (r *UserRepository) SetPasswordResetToken(id uuid.UUID, token uuid.UUID, expires time.Time) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password_reset_token":   token,
		"password_reset_expires": expires,
	})
	if result.Error != nil {
		log.Printf("[UserRepository] Error setting password reset token: %v", result.Error)
		return result.Error
	}
	return nil
}

// FindByPasswordResetToken zoekt een gebruiker op wachtwoord reset token
func (r *UserRepository) FindByPasswordResetToken(token string) (*models.User, error) {
	var user models.User
	result := r.db.Where("password_reset_token = ? AND password_reset_expires > ?", token, time.Now()).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[UserRepository] Error finding user by reset token: %v", result.Error)
		return nil, result.Error
	}
	return &user, nil
}

// ClearPasswordResetToken wist een wachtwoord reset token
func (r *UserRepository) ClearPasswordResetToken(id uuid.UUID) error {
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"password_reset_token":   nil,
		"password_reset_expires": nil,
	})
	if result.Error != nil {
		log.Printf("[UserRepository] Error clearing password reset token: %v", result.Error)
		return result.Error
	}
	return nil
}

// ApproveUser keurt een gebruiker goed
func (r *UserRepository) ApproveUser(id uuid.UUID, approvedBy uuid.UUID) error {
	now := time.Now()
	result := r.db.Model(&models.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":      models.StatusActive,
		"approved_by": approvedBy,
		"approved_at": now,
	})
	if result.Error != nil {
		log.Printf("[UserRepository] Error approving user: %v", result.Error)
		return result.Error
	}
	return nil
}

// CreateRefreshToken maakt een nieuw refresh token aan
func (r *UserRepository) CreateRefreshToken(userID uuid.UUID, token string, expiresAt time.Time) (*models.RefreshToken, error) {
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     token,
		ExpiresAt: expiresAt,
	}
	result := r.db.Create(refreshToken)
	if result.Error != nil {
		log.Printf("[UserRepository] Error creating refresh token: %v", result.Error)
		return nil, result.Error
	}
	return refreshToken, nil
}

// FindRefreshToken zoekt een refresh token
func (r *UserRepository) FindRefreshToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	result := r.db.Where("token = ? AND expires_at > ? AND revoked = false", token, time.Now()).First(&refreshToken)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Printf("[UserRepository] Error finding refresh token: %v", result.Error)
		return nil, result.Error
	}
	return &refreshToken, nil
}

// RevokeRefreshToken herroept een refresh token
func (r *UserRepository) RevokeRefreshToken(token string) error {
	now := time.Now()
	result := r.db.Model(&models.RefreshToken{}).Where("token = ?", token).Updates(map[string]interface{}{
		"revoked":    true,
		"revoked_at": now,
	})
	if result.Error != nil {
		log.Printf("[UserRepository] Error revoking refresh token: %v", result.Error)
		return result.Error
	}
	return nil
}

// RevokeAllUserRefreshTokens herroept alle refresh tokens van een gebruiker
func (r *UserRepository) RevokeAllUserRefreshTokens(userID uuid.UUID) error {
	now := time.Now()
	result := r.db.Model(&models.RefreshToken{}).Where("user_id = ? AND revoked = false", userID).Updates(map[string]interface{}{
		"revoked":    true,
		"revoked_at": now,
	})
	if result.Error != nil {
		log.Printf("[UserRepository] Error revoking all user refresh tokens: %v", result.Error)
		return result.Error
	}
	return nil
}

// DeleteByID verwijdert een gebruiker op basis van ID
func (r *UserRepository) DeleteByID(id uuid.UUID) error {
	// Eerst alle refresh tokens verwijderen
	if err := r.db.Where("user_id = ?", id).Delete(&models.RefreshToken{}).Error; err != nil {
		log.Printf("[UserRepository] Error deleting user refresh tokens: %v", err)
		return err
	}

	// Daarna de gebruiker verwijderen
	result := r.db.Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		log.Printf("[UserRepository] Error deleting user: %v", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("[UserRepository] No user found with ID: %v", id)
		return errors.New("gebruiker niet gevonden")
	}

	return nil
}
