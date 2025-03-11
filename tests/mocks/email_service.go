package mocks

import (
	"dklautomationgo/models"

	"github.com/stretchr/testify/mock"
)

// MockEmailService is a mock implementation of the EmailService
type MockEmailService struct {
	mock.Mock
}

// SendAanmeldingBevestiging mocks the SendAanmeldingBevestiging method
func (m *MockEmailService) SendAanmeldingBevestiging(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// SendContactBevestiging mocks the SendContactBevestiging method
func (m *MockEmailService) SendContactBevestiging(contact *models.ContactFormulier) error {
	args := m.Called(contact)
	return args.Error(0)
}

// SendPasswordReset mocks the SendPasswordReset method
func (m *MockEmailService) SendPasswordReset(user *models.User, resetToken string) error {
	args := m.Called(user, resetToken)
	return args.Error(0)
}

// SendWelcomeEmail mocks the SendWelcomeEmail method
func (m *MockEmailService) SendWelcomeEmail(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// SendStartnummerEmail mocks the SendStartnummerEmail method
func (m *MockEmailService) SendStartnummerEmail(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// SendBulkEmail mocks the SendBulkEmail method
func (m *MockEmailService) SendBulkEmail(recipients []string, subject string, templateName string, data map[string]interface{}) error {
	args := m.Called(recipients, subject, templateName, data)
	return args.Error(0)
}
