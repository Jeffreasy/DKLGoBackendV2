package handlers

import (
	"bytes"
	"dklautomationgo/models"
	"dklautomationgo/tests/fixtures"
	"dklautomationgo/tests/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupTest() (*mocks.MockAuthService, *mocks.MockAuthMiddleware, *AuthHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)

	mockAuthService := new(mocks.MockAuthService)
	mockAuthMiddleware := new(mocks.MockAuthMiddleware)

	handler := NewAuthHandler(mockAuthService, mockAuthMiddleware)

	router := gin.New()
	router.Use(gin.Recovery())

	return mockAuthService, mockAuthMiddleware, handler, router
}

func TestLogin_Success(t *testing.T) {
	// Setup
	mockAuthService, _, handler, router := setupTest()

	// Configure router
	router.POST("/api/auth/login", handler.Login)

	// Setup mock expectations
	tokenResponse := fixtures.GetTestTokenResponse()
	mockAuthService.On("Login", "test@example.com", "password123").Return(tokenResponse, nil)

	// Create request
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	jsonBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, tokenResponse.AccessToken, response.AccessToken)
	assert.Equal(t, tokenResponse.RefreshToken, response.RefreshToken)
	assert.Equal(t, tokenResponse.ExpiresIn, response.ExpiresIn)
	assert.Equal(t, tokenResponse.TokenType, response.TokenType)

	// Verify mock
	mockAuthService.AssertExpectations(t)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	// Setup
	mockAuthService, _, handler, router := setupTest()

	// Configure router
	router.POST("/api/auth/login", handler.Login)

	// Setup mock expectations
	mockAuthService.On("Login", "test@example.com", "wrongpassword").Return(nil, errors.New("ongeldige inloggegevens"))

	// Create request
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}
	jsonBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "ongeldige inloggegevens", response["error"])

	// Verify mock
	mockAuthService.AssertExpectations(t)
}

func TestLogin_InvalidInput(t *testing.T) {
	// Setup
	_, _, handler, router := setupTest()

	// Configure router
	router.POST("/api/auth/login", handler.Login)

	// Create request with invalid JSON
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
}
