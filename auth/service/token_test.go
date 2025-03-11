package service

import (
	"dklautomationgo/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	// Setup
	tokenService := NewTokenService()
	user := &models.User{
		ID:    uuid.New(),
		Email: "test@example.com",
		Role:  models.RoleBeheerder,
	}

	// Test
	token, err := tokenService.GenerateAccessToken(user)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Validate token
	claims, err := tokenService.ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, user.ID.String(), claims.UserID)
	assert.Equal(t, user.Email, claims.Email)
	assert.Equal(t, string(user.Role), claims.Role)
}

func TestGetUserIDFromToken(t *testing.T) {
	// Setup
	tokenService := NewTokenService()
	userID := uuid.New()
	user := &models.User{
		ID:    userID,
		Email: "test@example.com",
		Role:  models.RoleBeheerder,
	}

	// Generate token
	token, err := tokenService.GenerateAccessToken(user)
	assert.NoError(t, err)

	// Test
	extractedID, err := tokenService.GetUserIDFromToken(token)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, userID, extractedID)
}
