package fixtures

import (
	"dklautomationgo/models"
	"time"

	"github.com/google/uuid"
)

// GetTestUser returns a test user with admin role
func GetTestAdmin() *models.User {
	user := &models.User{
		ID:           uuid.MustParse("11111111-1111-1111-1111-111111111111"),
		Email:        "admin@dekoninklijkeloop.nl",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password123
		Role:         models.RoleBeheerder,
		Status:       models.StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return user
}

// GetTestUser returns a regular test user
func GetTestUser() *models.User {
	user := &models.User{
		ID:           uuid.MustParse("22222222-2222-2222-2222-222222222222"),
		Email:        "user@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password123
		Role:         models.RoleVrijwilliger,
		Status:       models.StatusActive,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return user
}

// GetTestPendingUser returns a test user with pending status
func GetTestPendingUser() *models.User {
	user := &models.User{
		ID:           uuid.MustParse("33333333-3333-3333-3333-333333333333"),
		Email:        "pending@example.com",
		PasswordHash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // password123
		Role:         models.RoleVrijwilliger,
		Status:       models.StatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	return user
}

// GetTestTokenResponse returns a test token response
func GetTestTokenResponse() *models.TokenResponse {
	return &models.TokenResponse{
		AccessToken:  "test-access-token",
		RefreshToken: "test-refresh-token",
		ExpiresIn:    900,
		TokenType:    "Bearer",
	}
}
