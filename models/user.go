package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserRole definieert de mogelijke rollen voor gebruikers
type UserRole string

const (
	RoleBeheerder    UserRole = "BEHEERDER"
	RoleAdmin        UserRole = "ADMIN"
	RoleVrijwilliger UserRole = "VRIJWILLIGER"
	RoleGebruiker    UserRole = "GEBRUIKER"
)

// UserStatus definieert de mogelijke statussen voor gebruikers
type UserStatus string

const (
	StatusPending  UserStatus = "PENDING"
	StatusActive   UserStatus = "ACTIVE"
	StatusInactive UserStatus = "INACTIVE"
)

// User representeert een gebruiker in het systeem
type User struct {
	ID                   uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Email                string     `json:"email" gorm:"type:varchar(255);unique;not null"`
	PasswordHash         string     `json:"-" gorm:"type:varchar(255);not null"` // Niet zichtbaar in JSON
	Role                 UserRole   `json:"role" gorm:"type:user_role;not null"`
	Status               UserStatus `json:"status" gorm:"type:user_status;not null;default:'PENDING'"`
	ApprovedBy           *uuid.UUID `json:"approved_by,omitempty" gorm:"type:uuid;references:id"`
	ApprovedAt           *time.Time `json:"approved_at,omitempty" gorm:"type:timestamp with time zone"`
	LastLogin            *time.Time `json:"last_login,omitempty" gorm:"type:timestamp with time zone"`
	PasswordResetToken   *uuid.UUID `json:"-" gorm:"type:uuid"`
	PasswordResetExpires *time.Time `json:"-" gorm:"type:timestamp with time zone"`
	CreatedAt            time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt            time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:now()"`
}

// UserResponse is een veilige versie van User voor API responses
type UserResponse struct {
	ID         uuid.UUID  `json:"id"`
	Email      string     `json:"email"`
	Role       UserRole   `json:"role"`
	Status     UserStatus `json:"status"`
	ApprovedAt *time.Time `json:"approved_at,omitempty"`
	LastLogin  *time.Time `json:"last_login,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

// ToResponse converteert een User naar een veilige UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:         u.ID,
		Email:      u.Email,
		Role:       u.Role,
		Status:     u.Status,
		ApprovedAt: u.ApprovedAt,
		LastLogin:  u.LastLogin,
		CreatedAt:  u.CreatedAt,
		UpdatedAt:  u.UpdatedAt,
	}
}

// LoginRequest representeert een login verzoek
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// CreateUserRequest representeert een verzoek om een nieuwe gebruiker aan te maken
type CreateUserRequest struct {
	Email    string   `json:"email" binding:"required,email"`
	Password string   `json:"password" binding:"required,min=8"`
	Role     UserRole `json:"role" binding:"required"`
}

// UpdateUserRequest representeert een verzoek om een gebruiker bij te werken
type UpdateUserRequest struct {
	Email  *string     `json:"email" binding:"omitempty,email"`
	Role   *UserRole   `json:"role"`
	Status *UserStatus `json:"status"`
}

// ChangePasswordRequest representeert een verzoek om een wachtwoord te wijzigen
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=8"`
}

// ForgotPasswordRequest representeert een verzoek om een wachtwoord reset
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest representeert een verzoek om een wachtwoord te resetten
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// AdminChangePasswordRequest representeert een verzoek van een beheerder om een wachtwoord te wijzigen
type AdminChangePasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// SetPassword stelt een nieuw wachtwoord in voor de gebruiker
func (u *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

// CheckPassword controleert of het opgegeven wachtwoord overeenkomt met de hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// RefreshToken representeert een refresh token voor JWT authenticatie
type RefreshToken struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `json:"user_id" gorm:"type:uuid;not null;references:id"`
	Token     string     `json:"token" gorm:"type:varchar(255);not null;unique"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"type:timestamp with time zone;not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	Revoked   bool       `json:"revoked" gorm:"type:boolean;not null;default:false"`
	RevokedAt *time.Time `json:"revoked_at,omitempty" gorm:"type:timestamp with time zone"`
}

// TokenResponse representeert een JWT token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // Seconden tot expiratie
	TokenType    string `json:"token_type"` // Meestal "Bearer"
}

// RefreshTokenRequest representeert een verzoek om een token te vernieuwen
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
