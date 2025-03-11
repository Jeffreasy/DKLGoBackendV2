package mocks

import (
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"dklautomationgo/services"

	"github.com/stretchr/testify/mock"
)

// MockAanmeldingService is een mock implementatie van de IAanmeldingService interface
type MockAanmeldingService struct {
	mock.Mock
}

// Controleer of MockAanmeldingService de IAanmeldingService interface implementeert
var _ services.IAanmeldingService = (*MockAanmeldingService)(nil)

// CreateAanmelding is een mock implementatie van de CreateAanmelding methode
func (m *MockAanmeldingService) CreateAanmelding(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// GetAanmeldingen is een mock implementatie van de GetAanmeldingen methode
func (m *MockAanmeldingService) GetAanmeldingen(params *repository.QueryParams) ([]models.Aanmelding, error) {
	args := m.Called(params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Aanmelding), args.Error(1)
}

// GetAanmeldingByID is een mock implementatie van de GetAanmeldingByID methode
func (m *MockAanmeldingService) GetAanmeldingByID(id string) (*models.Aanmelding, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Aanmelding), args.Error(1)
}

// UpdateAanmelding is een mock implementatie van de UpdateAanmelding methode
func (m *MockAanmeldingService) UpdateAanmelding(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// DeleteAanmelding is een mock implementatie van de DeleteAanmelding methode
func (m *MockAanmeldingService) DeleteAanmelding(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// CountAanmeldingen is een mock implementatie van de CountAanmeldingen methode
func (m *MockAanmeldingService) CountAanmeldingen(params *repository.QueryParams) (int64, error) {
	args := m.Called(params)
	return args.Get(0).(int64), args.Error(1)
}

// GetAanmeldingByEmail is een mock implementatie van de GetAanmeldingByEmail methode
func (m *MockAanmeldingService) GetAanmeldingByEmail(email string) (*models.Aanmelding, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Aanmelding), args.Error(1)
}

// SendBevestigingsEmail is een mock implementatie van de SendBevestigingsEmail methode
func (m *MockAanmeldingService) SendBevestigingsEmail(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}
