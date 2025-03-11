package middleware

import (
	"dklautomationgo/auth/service"
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware bevat middleware functies voor authenticatie
type AuthMiddleware struct {
	tokenService *service.TokenService
	userRepo     *repository.UserRepository
}

// NewAuthMiddleware maakt een nieuwe AuthMiddleware
func NewAuthMiddleware(tokenService *service.TokenService, userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		tokenService: tokenService,
		userRepo:     userRepo,
	}
}

// RequireAuth middleware controleert of de gebruiker is ingelogd
func (m *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Haal token uit Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authenticatie vereist"})
			return
		}

		// Controleer Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ongeldige authenticatie header"})
			return
		}

		tokenString := parts[1]

		// Valideer token
		claims, err := m.tokenService.ValidateToken(tokenString)
		if err != nil {
			log.Printf("[AuthMiddleware] Token validation error: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ongeldige of verlopen token"})
			return
		}

		// Haal gebruiker op
		userID, err := m.tokenService.GetUserIDFromToken(tokenString)
		if err != nil {
			log.Printf("[AuthMiddleware] Error getting user ID from token: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Ongeldige token"})
			return
		}

		user, err := m.userRepo.FindByID(userID)
		if err != nil {
			log.Printf("[AuthMiddleware] Error finding user: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Serverfout"})
			return
		}

		if user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Gebruiker niet gevonden"})
			return
		}

		// Controleer of gebruiker actief is
		if user.Status != models.StatusActive {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Gebruiker is niet actief"})
			return
		}

		// Sla gebruiker en claims op in context
		c.Set("user", user)
		c.Set("claims", claims)

		c.Next()
	}
}

// RequireRole middleware controleert of de gebruiker de vereiste rol heeft
func (m *AuthMiddleware) RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Haal gebruiker uit context
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authenticatie vereist"})
			return
		}

		userObj, ok := user.(*models.User)
		if !ok {
			log.Printf("[AuthMiddleware] Error casting user object")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Serverfout"})
			return
		}

		// Controleer of gebruiker een van de vereiste rollen heeft
		hasRole := false
		for _, role := range roles {
			if userObj.Role == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Onvoldoende rechten"})
			return
		}

		c.Next()
	}
}

// GetUserFromContext haalt de gebruiker uit de context
func GetUserFromContext(c *gin.Context) *models.User {
	user, exists := c.Get("user")
	if !exists {
		return nil
	}

	userObj, ok := user.(*models.User)
	if !ok {
		return nil
	}

	return userObj
}
