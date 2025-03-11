package service

import (
	"dklautomationgo/models"

	"github.com/google/uuid"
)

// IAuthService definieert de interface voor de AuthService
type IAuthService interface {
	Login(email, password string) (*models.TokenResponse, error)
	RefreshToken(refreshToken string) (*models.TokenResponse, error)
	Logout(refreshToken string) error
	LogoutAll(userID uuid.UUID) error
	CreateUser(email, password string, role models.UserRole) (*models.User, error)
	ApproveUser(userID, approverID uuid.UUID) error
	UpdateUser(userID uuid.UUID, updates *models.UpdateUserRequest) error
	ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error
	ForgotPassword(email string) (string, error)
	ResetPassword(token, newPassword string) error
	GetUserByID(id uuid.UUID) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers() ([]models.User, error)
	DeleteUser(userID uuid.UUID, deleterID uuid.UUID) error
	AdminChangePassword(userID uuid.UUID, adminID uuid.UUID, newPassword string) error
	GetUserRepository() interface{}
}
