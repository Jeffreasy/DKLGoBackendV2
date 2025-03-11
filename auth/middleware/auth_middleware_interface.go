package middleware

import (
	"dklautomationgo/models"

	"github.com/gin-gonic/gin"
)

// IAuthMiddleware definieert de interface voor de AuthMiddleware
type IAuthMiddleware interface {
	RequireAuth() gin.HandlerFunc
	RequireRole(roles ...models.UserRole) gin.HandlerFunc
}
