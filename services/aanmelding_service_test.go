package services_test

import (
	"dklautomationgo/database/repository"
	"dklautomationgo/models"
	"dklautomationgo/services"
	"dklautomationgo/tests/fixtures"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockAanmeldingRepository is een mock implementatie van de IAanmeldingRepository interface
type MockAanmeldingRepository struct {
	mock.Mock
}

// Create is een mock implementatie van de Create methode
func (m *MockAanmeldingRepository) Create(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// FindAll is een mock implementatie van de FindAll methode
func (m *MockAanmeldingRepository) FindAll(limit, offset int) ([]*models.Aanmelding, error) {
	args := m.Called(limit, offset)
	if args.Get(0) == nil {
		return []*models.Aanmelding{}, args.Error(1)
	}
	return args.Get(0).([]*models.Aanmelding), args.Error(1)
}

// FindByID is een mock implementatie van de FindByID methode
func (m *MockAanmeldingRepository) FindByID(id string) (*models.Aanmelding, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Aanmelding), args.Error(1)
}

// Update is een mock implementatie van de Update methode
func (m *MockAanmeldingRepository) Update(aanmelding *models.Aanmelding) error {
	args := m.Called(aanmelding)
	return args.Error(0)
}

// Count is een mock implementatie van de Count methode
func (m *MockAanmeldingRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

// MockEmailService is een mock implementatie van de IEmailService interface
type MockEmailService struct {
	mock.Mock
}

// SendAanmeldingEmail is een mock implementatie van de SendAanmeldingEmail methode
func (m *MockEmailService) SendAanmeldingEmail(data *models.AanmeldingEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

// SendContactEmail is een mock implementatie van de SendContactEmail methode
func (m *MockEmailService) SendContactEmail(data *models.ContactEmailData) error {
	args := m.Called(data)
	return args.Error(0)
}

func setupAanmeldingServiceTest() (*services.AanmeldingService, *MockAanmeldingRepository, *MockEmailService) {
	mockRepo := new(MockAanmeldingRepository)
	mockEmailService := new(MockEmailService)

	// Gebruik de interfaces in plaats van concrete types
	service := services.NewAanmeldingService(mockRepo, mockEmailService)

	return service, mockRepo, mockEmailService
}

func TestCreateAanmelding_Success(t *testing.T) {
	// Setup
	service, mockRepo, mockEmailService := setupAanmeldingServiceTest()
	testAanmelding := fixtures.GetTestAanmelding()

	// Mock verwachtingen
	mockRepo.On("Create", testAanmelding).Return(nil)
	mockEmailService.On("SendAanmeldingEmail", mock.AnythingOfType("*models.AanmeldingEmailData")).Return(nil)
	mockRepo.On("Update", testAanmelding).Return(nil)

	// Voer de test uit
	err := service.CreateAanmelding(testAanmelding)

	// Controleer het resultaat
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestCreateAanmelding_RepositoryError(t *testing.T) {
	// Setup
	service, mockRepo, _ := setupAanmeldingServiceTest()
	testAanmelding := fixtures.GetTestAanmelding()

	// Mock verwachtingen
	mockRepo.On("Create", testAanmelding).Return(errors.New("repository error"))

	// Voer de test uit
	err := service.CreateAanmelding(testAanmelding)

	// Controleer het resultaat
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fout bij opslaan aanmelding")
	mockRepo.AssertExpectations(t)
}

func TestCreateAanmelding_EmailError(t *testing.T) {
	// Setup
	service, mockRepo, mockEmailService := setupAanmeldingServiceTest()
	testAanmelding := fixtures.GetTestAanmelding()

	// Mock verwachtingen
	mockRepo.On("Create", testAanmelding).Return(nil)
	mockEmailService.On("SendAanmeldingEmail", mock.AnythingOfType("*models.AanmeldingEmailData")).Return(errors.New("email error"))

	// Voer de test uit
	err := service.CreateAanmelding(testAanmelding)

	// Controleer het resultaat
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "fout bij versturen bevestigingsmail")
	mockRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}

func TestGetAanmeldingen_Success(t *testing.T) {
	// Setup
	service, mockRepo, _ := setupAanmeldingServiceTest()
	testAanmelding := fixtures.GetTestAanmelding()
	testAanmeldingen := []*models.Aanmelding{testAanmelding}
	params := repository.NewQueryParams()

	// Mock verwachtingen
	mockRepo.On("FindAll", params.GetLimit(), params.GetOffset()).Return(testAanmeldingen, nil)

	// Voer de test uit
	aanmeldingen, err := service.GetAanmeldingen(params)

	// Controleer het resultaat
	assert.NoError(t, err)
	assert.Len(t, aanmeldingen, 1)
	assert.Equal(t, testAanmelding.ID, aanmeldingen[0].ID)
	mockRepo.AssertExpectations(t)
}

func TestGetAanmeldingen_RepositoryError(t *testing.T) {
	// Setup
	service, mockRepo, _ := setupAanmeldingServiceTest()
	params := repository.NewQueryParams()

	// Mock verwachtingen
	mockRepo.On("FindAll", params.GetLimit(), params.GetOffset()).Return(nil, errors.New("repository error"))

	// Voer de test uit
	aanmeldingen, err := service.GetAanmeldingen(params)

	// Controleer het resultaat
	assert.Error(t, err)
	assert.Nil(t, aanmeldingen)
	assert.Contains(t, err.Error(), "fout bij ophalen aanmeldingen")
	mockRepo.AssertExpectations(t)
}

func TestGetAanmeldingByID_Success(t *testing.T) {
	// Setup
	service, mockRepo, _ := setupAanmeldingServiceTest()
	testAanmelding := fixtures.GetTestAanmelding()
	testID := testAanmelding.ID

	// Mock verwachtingen
	mockRepo.On("FindByID", testID).Return(testAanmelding, nil)

	// Voer de test uit
	aanmelding, err := service.GetAanmeldingByID(testID)

	// Controleer het resultaat
	assert.NoError(t, err)
	assert.Equal(t, testAanmelding.ID, aanmelding.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetAanmeldingByID_RepositoryError(t *testing.T) {
	// Setup
	service, mockRepo, _ := setupAanmeldingServiceTest()
	testID := "non-existent-id"

	// Mock verwachtingen
	mockRepo.On("FindByID", testID).Return(nil, errors.New("repository error"))

	// Voer de test uit
	aanmelding, err := service.GetAanmeldingByID(testID)

	// Controleer het resultaat
	assert.Error(t, err)
	assert.Nil(t, aanmelding)
	assert.Contains(t, err.Error(), "fout bij ophalen aanmelding")
	mockRepo.AssertExpectations(t)
}

func TestSendBevestigingsEmail_Success(t *testing.T) {
	// Setup
	service, mockRepo, mockEmailService := setupAanmeldingServiceTest()
	testAanmelding := fixtures.GetTestAanmelding()

	// Mock verwachtingen
	mockEmailService.On("SendAanmeldingEmail", mock.AnythingOfType("*models.AanmeldingEmailData")).Return(nil)
	mockRepo.On("Update", testAanmelding).Return(nil)

	// Voer de test uit
	err := service.SendBevestigingsEmail(testAanmelding)

	// Controleer het resultaat
	assert.NoError(t, err)
	assert.True(t, testAanmelding.EmailVerzonden)
	assert.NotNil(t, testAanmelding.EmailVerzondOp)
	mockRepo.AssertExpectations(t)
	mockEmailService.AssertExpectations(t)
}
