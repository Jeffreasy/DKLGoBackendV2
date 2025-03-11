package mocks

import (
	"dklautomationgo/handlers"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

// MockAanmeldingHandler is een mock implementatie van de IAanmeldingHandler interface
type MockAanmeldingHandler struct {
	mock.Mock
}

// Controleer of MockAanmeldingHandler de IAanmeldingHandler interface implementeert
var _ handlers.IAanmeldingHandler = (*MockAanmeldingHandler)(nil)

// CreateAanmelding is een mock implementatie van de CreateAanmelding methode
func (m *MockAanmeldingHandler) CreateAanmelding(c *gin.Context) {
	m.Called(c)
}

// GetAanmeldingen is een mock implementatie van de GetAanmeldingen methode
func (m *MockAanmeldingHandler) GetAanmeldingen(c *gin.Context) {
	m.Called(c)
}

// GetAanmeldingByID is een mock implementatie van de GetAanmeldingByID methode
func (m *MockAanmeldingHandler) GetAanmeldingByID(c *gin.Context) {
	m.Called(c)
}

// UpdateAanmelding is een mock implementatie van de UpdateAanmelding methode
func (m *MockAanmeldingHandler) UpdateAanmelding(c *gin.Context) {
	m.Called(c)
}

// DeleteAanmelding is een mock implementatie van de DeleteAanmelding methode
func (m *MockAanmeldingHandler) DeleteAanmelding(c *gin.Context) {
	m.Called(c)
}
