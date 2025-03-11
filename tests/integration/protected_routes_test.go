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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type ProtectedRoutesTestSuite struct {
	suite.Suite
	db             *gorm.DB
	router         *gin.Engine
	authMiddleware *middleware.AuthMiddleware
	tokenService   *service.TokenService
	userRepo       *repository.UserRepository
	adminUser      *models.User
	regularUser    *models.User
}

func TestProtectedRoutesSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	suite.Run(t, new(ProtectedRoutesTestSuite))
}

func (s *ProtectedRoutesTestSuite) SetupSuite() {
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
	s.userRepo = repository.NewUserRepository(s.db)

	// Setup services
	s.tokenService = service.NewTokenService()
	authService := service.NewAuthService(s.userRepo, s.tokenService)

	// Setup middleware
	s.authMiddleware = middleware.NewAuthMiddleware(s.tokenService, s.userRepo)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(authService, s.authMiddleware)

	// Register auth routes
	auth := s.router.Group("/api/auth")
	{
		auth.POST("/login", authHandler.Login)
	}

	// Register protected routes
	protected := s.router.Group("/api/protected")
	protected.Use(s.authMiddleware.RequireAuth())
	{
		protected.GET("/user-only", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "User only content"})
		})
	}

	// Register admin routes
	admin := s.router.Group("/api/admin")
	admin.Use(s.authMiddleware.RequireAuth(), s.authMiddleware.RequireRole(models.RoleBeheerder))
	{
		admin.GET("/admin-only", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "Admin only content"})
		})
	}
}

func (s *ProtectedRoutesTestSuite) TearDownSuite() {
	// Cleanup test database
	tests.TeardownTestDB(s.db)
}

func (s *ProtectedRoutesTestSuite) SetupTest() {
	// Clean up data before each test
	tests.CleanupTestData(s.db)

	// Create test users
	s.createTestUsers()
}

func (s *ProtectedRoutesTestSuite) createTestUsers() {
	// Create admin user
	s.adminUser = &models.User{
		Email:  "admin@example.com",
		Role:   models.RoleBeheerder,
		Status: models.StatusActive,
	}
	err := s.adminUser.SetPassword("admin123")
	s.Require().NoError(err)
	result := s.db.Create(s.adminUser)
	s.Require().NoError(result.Error)
	s.Require().NotZero(s.adminUser.ID)

	// Create regular user
	s.regularUser = &models.User{
		Email:  "user@example.com",
		Role:   models.RoleGebruiker,
		Status: models.StatusActive,
	}
	err = s.regularUser.SetPassword("user123")
	s.Require().NoError(err)
	result = s.db.Create(s.regularUser)
	s.Require().NoError(result.Error)
	s.Require().NotZero(s.regularUser.ID)
}

func (s *ProtectedRoutesTestSuite) login(email, password string) string {
	// Login and get access token
	loginRequest := models.LoginRequest{
		Email:    email,
		Password: password,
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

	return response.AccessToken
}

func (s *ProtectedRoutesTestSuite) TestProtectedRouteWithoutToken() {
	// Test accessing protected route without token
	req, _ := http.NewRequest("GET", "/api/protected/user-only", nil)
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusUnauthorized, w.Code)
}

func (s *ProtectedRoutesTestSuite) TestProtectedRouteWithToken() {
	// Login as regular user
	token := s.login("user@example.com", "user123")

	// Test accessing protected route with token
	req, _ := http.NewRequest("GET", "/api/protected/user-only", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().Equal("User only content", response["message"])
}

func (s *ProtectedRoutesTestSuite) TestAdminRouteWithRegularUser() {
	// Login as regular user
	token := s.login("user@example.com", "user123")

	// Test accessing admin route with regular user token
	req, _ := http.NewRequest("GET", "/api/admin/admin-only", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusForbidden, w.Code)
}

func (s *ProtectedRoutesTestSuite) TestAdminRouteWithAdminUser() {
	// Login as admin user
	token := s.login("admin@example.com", "admin123")

	// Test accessing admin route with admin token
	req, _ := http.NewRequest("GET", "/api/admin/admin-only", nil)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, req)

	s.Assert().Equal(http.StatusOK, w.Code)

	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	s.Require().NoError(err)
	s.Assert().Equal("Admin only content", response["message"])
}
