package integration

import (
	"bytes"
	"dklautomationgo/auth/handlers"
	"dklautomationgo/auth/middleware"
	"dklautomationgo/auth/service"
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"dklautomationgo/tests"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuthIntegrationTestSuite struct {
	suite.Suite
	db     *gorm.DB
	router *gin.Engine
}

func TestAuthIntegrationSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	suite.Run(t, new(AuthIntegrationTestSuite))
}

func (s *AuthIntegrationTestSuite) SetupSuite() {
	// Setup test database
	var err error
	s.db, err = tests.SetupTestDB()
	if err != nil {
		s.T().Fatalf("Failed to setup test database: %v", err)
	}

	// Setup router
	gin.SetMode(gin.TestMode)
	s.router = gin.Default()

	// Setup repositories
	userRepo := repository.NewUserRepository(s.db)

	// Setup services
	tokenService := service.NewTokenService()
	authService := service.NewAuthService(userRepo, tokenService)

	// Setup middleware
	authMiddleware := middleware.NewAuthMiddleware(tokenService, userRepo)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authService, authMiddleware)

	// Register routes
	auth := s.router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/refresh-token", authHandler.RefreshToken)
		auth.POST("/logout", authHandler.Logout)
	}
}

func (s *AuthIntegrationTestSuite) TearDownSuite() {
	// Cleanup test database
	tests.TeardownTestDB(s.db)
}

func (s *AuthIntegrationTestSuite) SetupTest() {
	// Clean up data before each test
	tests.CleanupTestData(s.db)
}

func (s *AuthIntegrationTestSuite) TestAuthFlow() {
	// Create a test user
	s.createTestUser()

	// Test login
	tokenResponse := s.testLogin()

	// Test refresh token
	s.testRefreshToken(tokenResponse.RefreshToken)

	// Test logout
	s.testLogout(tokenResponse.RefreshToken)
}

func (s *AuthIntegrationTestSuite) createTestUser() {
	// Create a test user directly in the database
	user := &models.User{
		Email:  "test@example.com",
		Role:   models.RoleBeheerder,
		Status: models.StatusActive,
	}
	err := user.SetPassword("password123")
	s.Require().NoError(err)

	result := s.db.Create(user)
	s.Require().NoError(result.Error)
	s.Require().NotZero(user.ID)
}

func (s *AuthIntegrationTestSuite) testLogin() *models.TokenResponse {
	// Test login
	loginRequest := models.LoginRequest{
		Email:    "test@example.com",
		Password: "password123",
	}

	jsonBody, _ := json.Marshal(loginRequest)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var response models.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().NotEmpty(response.AccessToken)
	s.Assert().NotEmpty(response.RefreshToken)

	return &response
}

func (s *AuthIntegrationTestSuite) testRefreshToken(refreshToken string) {
	// Test refresh token
	refreshRequest := models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonBody, _ := json.Marshal(refreshRequest)
	req, _ := http.NewRequest("POST", "/api/auth/refresh-token", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var response models.TokenResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().NotEmpty(response.AccessToken)
	s.Assert().NotEmpty(response.RefreshToken)
}

func (s *AuthIntegrationTestSuite) testLogout(refreshToken string) {
	// Test logout
	logoutRequest := models.RefreshTokenRequest{
		RefreshToken: refreshToken,
	}

	jsonBody, _ := json.Marshal(logoutRequest)
	req, _ := http.NewRequest("POST", "/api/auth/logout", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().Contains(response, "message")
	s.Assert().Equal("Succesvol uitgelogd", response["message"])
}
