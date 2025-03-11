package mocks

import (
	"dklautomationgo/auth/service"
	"dklautomationgo/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

// Voeg hier een commentaar toe om te bevestigen dat MockAuthService de IAuthService interface implementeert
// Zorg ervoor dat de compiler dit controleert
var _ service.IAuthService = (*MockAuthService)(nil)

// MockAuthService is a mock implementation of the AuthService
type MockAuthService struct {
	mock.Mock
}

// Login mocks the Login method
func (m *MockAuthService) Login(email, password string) (*models.TokenResponse, error) {
	args := m.Called(email, password)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenResponse), args.Error(1)
}

// RefreshToken mocks the RefreshToken method
func (m *MockAuthService) RefreshToken(refreshToken string) (*models.TokenResponse, error) {
	args := m.Called(refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TokenResponse), args.Error(1)
}

// Logout mocks the Logout method
func (m *MockAuthService) Logout(refreshToken string) error {
	args := m.Called(refreshToken)
	return args.Error(0)
}

// LogoutAll mocks the LogoutAll method
func (m *MockAuthService) LogoutAll(userID uuid.UUID) error {
	args := m.Called(userID)
	return args.Error(0)
}

// CreateUser mocks the CreateUser method
func (m *MockAuthService) CreateUser(email, password string, role models.UserRole) (*models.User, error) {
	args := m.Called(email, password, role)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// ApproveUser mocks the ApproveUser method
func (m *MockAuthService) ApproveUser(userID, approverID uuid.UUID) error {
	args := m.Called(userID, approverID)
	return args.Error(0)
}

// UpdateUser mocks the UpdateUser method
func (m *MockAuthService) UpdateUser(userID uuid.UUID, updates *models.UpdateUserRequest) error {
	args := m.Called(userID, updates)
	return args.Error(0)
}

// ChangePassword mocks the ChangePassword method
func (m *MockAuthService) ChangePassword(userID uuid.UUID, currentPassword, newPassword string) error {
	args := m.Called(userID, currentPassword, newPassword)
	return args.Error(0)
}

// ForgotPassword mocks the ForgotPassword method
func (m *MockAuthService) ForgotPassword(email string) (string, error) {
	args := m.Called(email)
	return args.String(0), args.Error(1)
}

// ResetPassword mocks the ResetPassword method
func (m *MockAuthService) ResetPassword(token, newPassword string) error {
	args := m.Called(token, newPassword)
	return args.Error(0)
}

// GetUserByID mocks the GetUserByID method
func (m *MockAuthService) GetUserByID(id uuid.UUID) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// GetUserByEmail mocks the GetUserByEmail method
func (m *MockAuthService) GetUserByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// GetAllUsers mocks the GetAllUsers method
func (m *MockAuthService) GetAllUsers() ([]models.User, error) {
	args := m.Called()
	return args.Get(0).([]models.User), args.Error(1)
}

// DeleteUser mocks the DeleteUser method
func (m *MockAuthService) DeleteUser(userID uuid.UUID, deleterID uuid.UUID) error {
	args := m.Called(userID, deleterID)
	return args.Error(0)
}

// AdminChangePassword mocks the AdminChangePassword method
func (m *MockAuthService) AdminChangePassword(userID uuid.UUID, adminID uuid.UUID, newPassword string) error {
	args := m.Called(userID, adminID, newPassword)
	return args.Error(0)
}

// GetUserRepository mocks the GetUserRepository method
func (m *MockAuthService) GetUserRepository() interface{} {
	args := m.Called()
	return args.Get(0)
}
