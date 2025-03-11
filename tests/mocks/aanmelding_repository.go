package mocks

import (
	"dklautomationgo/database/repository"
	"dklautomationgo/models"

	"github.com/stretchr/testify/mock"
)

// MockAanmeldingRepository is a mock implementation of the AanmeldingRepository
type MockAanmeldingRepository struct {
	mock.Mock
}

// Create mocks the Create method
func (m *MockAanmeldingRepository) Create(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// GetAll mocks the GetAll method
func (m *MockAanmeldingRepository) GetAll(params *repository.QueryParams) ([]models.Aanmelding, error) {
	args := m.Called(params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).([]models.Aanmelding), args.Error(1)
}

// GetByID mocks the GetByID method
func (m *MockAanmeldingRepository) GetByID(id string) (*models.Aanmelding, error) {
	args := m.Called(id)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Aanmelding), args.Error(1)
}

// Update mocks the Update method
func (m *MockAanmeldingRepository) Update(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *MockAanmeldingRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// GetByEmail mocks the GetByEmail method
func (m *MockAanmeldingRepository) GetByEmail(email string) (*models.Aanmelding, error) {
	args := m.Called(email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Aanmelding), args.Error(1)
}

// GetByStartnummer mocks the GetByStartnummer method
func (m *MockAanmeldingRepository) GetByStartnummer(startnummer string) (*models.Aanmelding, error) {
	args := m.Called(startnummer)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Aanmelding), args.Error(1)
}

// Count mocks the Count method
func (m *MockAanmeldingRepository) Count(params *repository.QueryParams) (int64, error) {
	args := m.Called(params)
	return args.Get(0).(int64), args.Error(1)
}
