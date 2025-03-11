package handlers

import (
	"dklautomationgo/auth/middleware"
	"dklautomationgo/auth/service"
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthHandler bevat handlers voor authenticatie
type AuthHandler struct {
	authService    service.IAuthService
	authMiddleware middleware.IAuthMiddleware
}

// NewAuthHandler maakt een nieuwe AuthHandler
func NewAuthHandler(authService service.IAuthService, authMiddleware middleware.IAuthMiddleware) *AuthHandler {
	return &AuthHandler{
		authService:    authService,
		authMiddleware: authMiddleware,
	}
}

// RegisterRoutes registreert de auth routes
func (h *AuthHandler) RegisterRoutes(r *gin.Engine) {
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh-token", h.RefreshToken)
		auth.POST("/forgot-password", h.ForgotPassword)
		auth.POST("/reset-password", h.ResetPassword)
		auth.GET("/test-password", h.TestPassword)
		auth.GET("/create-admin", h.CreateAdminUser)
		auth.GET("/update-admin-password", h.UpdateAdminPassword)
		auth.GET("/create-beheerder", h.CreateBeheerderUser)
		auth.GET("/find-user", h.FindUserByEmail)
		auth.GET("/rename-beheerder", h.RenameBeheerder)

		// Beschermde routes
		secured := auth.Use(h.authMiddleware.RequireAuth())
		{
			secured.POST("/logout", h.Logout)
			secured.PUT("/password", h.ChangePassword)

			// Admin routes
			admin := auth.Group("/admin")
			admin.Use(h.authMiddleware.RequireAuth())
			admin.Use(h.authMiddleware.RequireRole(models.RoleBeheerder))
			{
				admin.POST("/users", h.CreateUser)
				admin.GET("/users", h.GetUsers)
				admin.GET("/users/:id", h.GetUser)
				admin.PUT("/users/:id", h.UpdateUser)
				admin.PUT("/users/:id/approve", h.ApproveUser)
				admin.DELETE("/users/:id", h.DeleteUser)
				admin.PUT("/users/:id/password", h.AdminChangePassword)
			}
		}
	}
}

// Login handelt login verzoeken af
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	tokens, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		log.Printf("[AuthHandler] Login error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshToken handelt token refresh verzoeken af
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	tokens, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		log.Printf("[AuthHandler] Refresh token error: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Logout handelt logout verzoeken af
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	if err := h.authService.Logout(req.RefreshToken); err != nil {
		log.Printf("[AuthHandler] Logout error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij uitloggen"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Succesvol uitgelogd"})
}

// ForgotPassword handelt wachtwoord vergeten verzoeken af
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req models.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	token, err := h.authService.ForgotPassword(req.Email)
	if err != nil {
		log.Printf("[AuthHandler] Forgot password error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij verwerken van verzoek"})
		return
	}

	// In een echte applicatie zou je hier een email sturen met de reset link
	// Voor nu geven we de token terug in de response (alleen voor ontwikkeling)
	c.JSON(http.StatusOK, gin.H{
		"message": "Wachtwoord reset link is verzonden",
		"token":   token, // Verwijder dit in productie!
	})
}

// ResetPassword handelt wachtwoord reset verzoeken af
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req models.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	if err := h.authService.ResetPassword(req.Token, req.NewPassword); err != nil {
		log.Printf("[AuthHandler] Reset password error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wachtwoord succesvol gewijzigd"})
}

// ChangePassword handelt wachtwoord wijziging verzoeken af
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	user := middleware.GetUserFromContext(c)
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niet ingelogd"})
		return
	}

	if err := h.authService.ChangePassword(user.ID, req.CurrentPassword, req.NewPassword); err != nil {
		log.Printf("[AuthHandler] Change password error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wachtwoord succesvol gewijzigd"})
}

// CreateUser handelt gebruiker aanmaak verzoeken af
func (h *AuthHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	user, err := h.authService.CreateUser(req.Email, req.Password, req.Role)
	if err != nil {
		log.Printf("[AuthHandler] Create user error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user.ToResponse())
}

// GetUsers haalt alle gebruikers op
func (h *AuthHandler) GetUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers()
	if err != nil {
		log.Printf("[AuthHandler] Get users error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij ophalen gebruikers"})
		return
	}

	// Converteer naar veilige response objecten
	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	c.JSON(http.StatusOK, responses)
}

// GetUser haalt een specifieke gebruiker op
func (h *AuthHandler) GetUser(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige gebruiker ID"})
		return
	}

	user, err := h.authService.GetUserByID(id)
	if err != nil {
		log.Printf("[AuthHandler] Get user error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij ophalen gebruiker"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gebruiker niet gevonden"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// UpdateUser werkt een gebruiker bij
func (h *AuthHandler) UpdateUser(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige gebruiker ID"})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	if err := h.authService.UpdateUser(id, &req); err != nil {
		log.Printf("[AuthHandler] Update user error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := h.authService.GetUserByID(id)
	c.JSON(http.StatusOK, user.ToResponse())
}

// ApproveUser keurt een gebruiker goed
func (h *AuthHandler) ApproveUser(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige gebruiker ID"})
		return
	}

	approver := middleware.GetUserFromContext(c)
	if approver == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niet ingelogd"})
		return
	}

	if err := h.authService.ApproveUser(id, approver.ID); err != nil {
		log.Printf("[AuthHandler] Approve user error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, _ := h.authService.GetUserByID(id)
	c.JSON(http.StatusOK, user.ToResponse())
}

// TestPassword is een test endpoint om wachtwoord hashing te testen
func (h *AuthHandler) TestPassword(c *gin.Context) {
	// Maak een nieuwe gebruiker
	user := &models.User{}

	// Set wachtwoord
	password := "Admin123!"
	if err := user.SetPassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Controleer wachtwoord
	if !user.CheckPassword(password) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Wachtwoord verificatie mislukt"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Wachtwoord hash en verificatie succesvol",
		"hash":    user.PasswordHash,
	})
}

// CreateAdminUser is een test endpoint om een admin gebruiker aan te maken
func (h *AuthHandler) CreateAdminUser(c *gin.Context) {
	// Maak een nieuwe admin gebruiker
	user, err := h.authService.CreateUser("beheerder@dekoninklijkeloop.nl", "Admin123!", models.RoleBeheerder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Activeer de gebruiker
	user.Status = models.StatusActive
	if err := h.authService.UpdateUser(user.ID, &models.UpdateUserRequest{
		Status: &user.Status,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Admin gebruiker aangemaakt",
		"user":    user.ToResponse(),
	})
}

// UpdateAdminPassword is een test endpoint om het wachtwoord van de beheerder te updaten
func (h *AuthHandler) UpdateAdminPassword(c *gin.Context) {
	// Maak een nieuwe gebruiker met het wachtwoord
	user := &models.User{}
	password := "Admin123!"
	if err := user.SetPassword(password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Update het wachtwoord direct in de database via de repository
	userRepoInterface := h.authService.GetUserRepository()
	userRepo, ok := userRepoInterface.(*repository.UserRepository)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kan user repository niet casten"})
		return
	}

	// Haal de gebruiker op
	existingUser, err := userRepo.FindByEmail("beheerder@dekoninklijkeloop.nl")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if existingUser == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Beheerder niet gevonden"})
		return
	}

	// Update het wachtwoord
	existingUser.PasswordHash = user.PasswordHash
	if err := userRepo.Update(existingUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beheerder wachtwoord geüpdatet",
		"hash":    user.PasswordHash,
	})
}

// CreateBeheerderUser is een test endpoint om een beheerder gebruiker aan te maken
func (h *AuthHandler) CreateBeheerderUser(c *gin.Context) {
	// Maak een nieuwe beheerder gebruiker
	user, err := h.authService.CreateUser("beheerder2@dekoninklijkeloop.nl", "Admin123!", models.RoleBeheerder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Activeer de gebruiker
	user.Status = models.StatusActive
	if err := h.authService.UpdateUser(user.ID, &models.UpdateUserRequest{
		Status: &user.Status,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Beheerder gebruiker aangemaakt",
		"user":    user.ToResponse(),
	})
}

// DeleteUser verwijdert een gebruiker
func (h *AuthHandler) DeleteUser(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige gebruiker ID"})
		return
	}

	deleter := middleware.GetUserFromContext(c)
	if deleter == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niet ingelogd"})
		return
	}

	if err := h.authService.DeleteUser(id, deleter.ID); err != nil {
		log.Printf("[AuthHandler] Delete user error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Gebruiker succesvol verwijderd"})
}

// FindUserByEmail zoekt een gebruiker op basis van email
func (h *AuthHandler) FindUserByEmail(c *gin.Context) {
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email parameter is vereist"})
		return
	}

	user, err := h.authService.GetUserByEmail(email)
	if err != nil {
		log.Printf("[AuthHandler] Find user by email error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij zoeken gebruiker"})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Gebruiker niet gevonden"})
		return
	}

	c.JSON(http.StatusOK, user.ToResponse())
}

// RenameBeheerder hernoemt beheerder2@dekoninklijkeloop.nl naar beheerder@dekoninklijkeloop.nl
func (h *AuthHandler) RenameBeheerder(c *gin.Context) {
	// Zoek de originele beheerder
	originalBeheerder, err := h.authService.GetUserByEmail("beheerder@dekoninklijkeloop.nl")
	if err != nil {
		log.Printf("[AuthHandler] Error finding original beheerder: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij zoeken originele beheerder"})
		return
	}

	// Zoek beheerder2
	beheerder2, err := h.authService.GetUserByEmail("beheerder2@dekoninklijkeloop.nl")
	if err != nil {
		log.Printf("[AuthHandler] Error finding beheerder2: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij zoeken beheerder2"})
		return
	}

	if beheerder2 == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Beheerder2 niet gevonden"})
		return
	}

	// Als de originele beheerder bestaat, verwijder deze
	if originalBeheerder != nil {
		if err := h.authService.DeleteUser(originalBeheerder.ID, beheerder2.ID); err != nil {
			log.Printf("[AuthHandler] Error deleting original beheerder: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij verwijderen originele beheerder"})
			return
		}
	}

	// Update de email van beheerder2 naar beheerder
	newEmail := "beheerder@dekoninklijkeloop.nl"
	if err := h.authService.UpdateUser(beheerder2.ID, &models.UpdateUserRequest{
		Email: &newEmail,
	}); err != nil {
		log.Printf("[AuthHandler] Error updating beheerder2 email: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Fout bij updaten email"})
		return
	}

	// Haal de bijgewerkte gebruiker op
	updatedUser, _ := h.authService.GetUserByID(beheerder2.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Beheerder succesvol hernoemd",
		"user":    updatedUser.ToResponse(),
	})
}

// AdminChangePassword stelt een beheerder in staat om het wachtwoord van een gebruiker te wijzigen
func (h *AuthHandler) AdminChangePassword(c *gin.Context) {
	id, err := parseUUID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige gebruiker ID"})
		return
	}

	var req models.AdminChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ongeldige invoer"})
		return
	}

	admin := middleware.GetUserFromContext(c)
	if admin == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Niet ingelogd"})
		return
	}

	if err := h.authService.AdminChangePassword(id, admin.ID, req.NewPassword); err != nil {
		log.Printf("[AuthHandler] Admin change password error: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Wachtwoord succesvol gewijzigd"})
}
