package mocks

import (
	"dklautomationgo/auth/middleware"
	"dklautomationgo/models"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// Voeg hier een commentaar toe om te bevestigen dat MockAuthMiddleware de IAuthMiddleware interface implementeert
// Zorg ervoor dat de compiler dit controleert
var _ middleware.IAuthMiddleware = (*MockAuthMiddleware)(nil)

// MockAuthMiddleware is a mock implementation of the AuthMiddleware
type MockAuthMiddleware struct {
	mock.Mock
}

// RequireAuth mocks the RequireAuth method
func (m *MockAuthMiddleware) RequireAuth() gin.HandlerFunc {
	args := m.Called()
	return args.Get(0).(gin.HandlerFunc)
}

// RequireRole mocks the RequireRole method
func (m *MockAuthMiddleware) RequireRole(roles ...models.UserRole) gin.HandlerFunc {
	args := m.Called(roles)
	return args.Get(0).(gin.HandlerFunc)
}
