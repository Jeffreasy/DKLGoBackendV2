package repository

import (
	"dklautomationgo/models"
	"dklautomationgo/tests"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AanmeldingRepositoryTestSuite struct {
	suite.Suite
	db         *gorm.DB
	repository *AanmeldingRepository
}

func TestAanmeldingRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping repository test in short mode")
	}
	suite.Run(t, new(AanmeldingRepositoryTestSuite))
}

func (s *AanmeldingRepositoryTestSuite) SetupSuite() {
	// Setup test database
	var err error
	s.db, err = tests.SetupTestDB()
	if err != nil {
		s.T().Fatalf("Failed to setup test database: %v", err)
	}

	// Create repository
	s.repository = NewAanmeldingRepository(s.db)
}

func (s *AanmeldingRepositoryTestSuite) TearDownSuite() {
	// Cleanup test database
	tests.TeardownTestDB(s.db)
}

func (s *AanmeldingRepositoryTestSuite) SetupTest() {
	// Clean up data before each test
	tests.CleanupTestData(s.db)
}

func (s *AanmeldingRepositoryTestSuite) TestCreate() {
	// Create test aanmelding
	aanmelding := &models.Aanmelding{
		Naam:           "Test User",
		Email:          "test@example.com",
		Telefoon:       "0612345678",
		Rol:            "chauffeur",
		Afstand:        "10 km",
		Ondersteuning:  "Geen",
		Bijzonderheden: "Geen",
		Terms:          true,
	}

	// Test create
	err := s.repository.Create(aanmelding)
	s.Require().NoError(err)
	s.Require().NotEmpty(aanmelding.ID)

	// Verify created aanmelding
	var result models.Aanmelding
	err = s.db.First(&result, "id = ?", aanmelding.ID).Error
	s.Require().NoError(err)
	s.Assert().Equal(aanmelding.Naam, result.Naam)
	s.Assert().Equal(aanmelding.Email, result.Email)
	s.Assert().Equal(aanmelding.Telefoon, result.Telefoon)
	s.Assert().Equal(aanmelding.Rol, result.Rol)
	s.Assert().Equal(aanmelding.Afstand, result.Afstand)
}

func (s *AanmeldingRepositoryTestSuite) TestFindAll() {
	// Create test aanmeldingen
	aanmeldingen := []models.Aanmelding{
		{
			Naam:           "Test User 1",
			Email:          "test1@example.com",
			Telefoon:       "0612345678",
			Rol:            "chauffeur",
			Afstand:        "10 km",
			Ondersteuning:  "Geen",
			Bijzonderheden: "Geen",
			Terms:          true,
		},
		{
			Naam:           "Test User 2",
			Email:          "test2@example.com",
			Telefoon:       "0687654321",
			Rol:            "bijrijder",
			Afstand:        "5 km",
			Ondersteuning:  "Geen",
			Bijzonderheden: "Geen",
			Terms:          true,
		},
	}

	for i := range aanmeldingen {
		err := s.db.Create(&aanmeldingen[i]).Error
		s.Require().NoError(err)
		s.Require().NotEmpty(aanmeldingen[i].ID)
	}

	// Test FindAll
	result, err := s.repository.FindAll(10, 0)
	s.Require().NoError(err)
	s.Assert().Len(result, 2)

	// Test FindAll with pagination
	result, err = s.repository.FindAll(1, 0)
	s.Require().NoError(err)
	s.Assert().Len(result, 1)
}

func (s *AanmeldingRepositoryTestSuite) TestFindByID() {
	// Create test aanmelding
	aanmelding := models.Aanmelding{
		Naam:           "Test User",
		Email:          "test@example.com",
		Telefoon:       "0612345678",
		Rol:            "chauffeur",
		Afstand:        "10 km",
		Ondersteuning:  "Geen",
		Bijzonderheden: "Geen",
		Terms:          true,
	}

	err := s.db.Create(&aanmelding).Error
	s.Require().NoError(err)
	s.Require().NotEmpty(aanmelding.ID)

	// Test FindByID
	result, err := s.repository.FindByID(aanmelding.ID)
	s.Require().NoError(err)
	s.Assert().Equal(aanmelding.ID, result.ID)
	s.Assert().Equal(aanmelding.Naam, result.Naam)
	s.Assert().Equal(aanmelding.Email, result.Email)
	s.Assert().Equal(aanmelding.Telefoon, result.Telefoon)

	// Test FindByID with non-existent ID
	result, err = s.repository.FindByID("non-existent-id")
	s.Require().Error(err)
	s.Assert().Equal(gorm.ErrRecordNotFound, err)
}

func (s *AanmeldingRepositoryTestSuite) TestUpdate() {
	// Create test aanmelding
	aanmelding := models.Aanmelding{
		Naam:           "Test User",
		Email:          "test@example.com",
		Telefoon:       "0612345678",
		Rol:            "chauffeur",
		Afstand:        "10 km",
		Ondersteuning:  "Geen",
		Bijzonderheden: "Geen",
		Terms:          true,
	}

	err := s.db.Create(&aanmelding).Error
	s.Require().NoError(err)
	s.Require().NotEmpty(aanmelding.ID)

	// Update aanmelding
	aanmelding.Naam = "Updated User"
	aanmelding.Afstand = "15 km"

	// Test Update
	err = s.repository.Update(&aanmelding)
	s.Require().NoError(err)

	// Verify updated aanmelding
	var result models.Aanmelding
	err = s.db.First(&result, "id = ?", aanmelding.ID).Error
	s.Require().NoError(err)
	s.Assert().Equal("Updated User", result.Naam)
	s.Assert().Equal("15 km", result.Afstand)
}

func (s *AanmeldingRepositoryTestSuite) TestCount() {
	// Create test aanmeldingen
	aanmeldingen := []models.Aanmelding{
		{
			Naam:           "Test User 1",
			Email:          "test1@example.com",
			Telefoon:       "0612345678",
			Rol:            "chauffeur",
			Afstand:        "10 km",
			Ondersteuning:  "Geen",
			Bijzonderheden: "Geen",
			Terms:          true,
		},
		{
			Naam:           "Test User 2",
			Email:          "test2@example.com",
			Telefoon:       "0687654321",
			Rol:            "bijrijder",
			Afstand:        "5 km",
			Ondersteuning:  "Geen",
			Bijzonderheden: "Geen",
			Terms:          true,
		},
	}

	for i := range aanmeldingen {
		err := s.db.Create(&aanmeldingen[i]).Error
		s.Require().NoError(err)
	}

	// Test Count
	count, err := s.repository.Count()
	s.Require().NoError(err)
	s.Assert().Equal(int64(2), count)

	// Test CountByRol
	count, err = s.repository.CountByRol("chauffeur")
	s.Require().NoError(err)
	s.Assert().Equal(int64(1), count)
}
