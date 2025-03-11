package handlers_test

import (
	"bytes"
	"dklautomationgo/handlers"
	"dklautomationgo/models"
	"dklautomationgo/tests/fixtures"
	"dklautomationgo/tests/mocks"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupAanmeldingTest() (*gin.Engine, *mocks.MockAanmeldingService, *handlers.AanmeldingHandler) {
	// Zet Gin in test modus
	gin.SetMode(gin.TestMode)

	// Maak een mock service
	mockService := new(mocks.MockAanmeldingService)

	// Maak een handler met de mock service
	handler := handlers.NewAanmeldingHandler(mockService)

	// Maak een router
	router := gin.Default()

	return router, mockService, handler
}

func TestCreateAanmelding_Success(t *testing.T) {
	// Setup
	router, mockService, handler := setupAanmeldingTest()
	router.POST("/aanmeldingen", handler.CreateAanmelding)

	// Test data
	testAanmelding := fixtures.GetTestAanmelding()

	// Mock verwachtingen
	mockService.On("CreateAanmelding", mock.AnythingOfType("*models.Aanmelding")).Return(nil)

	// Maak een request body
	jsonData, _ := json.Marshal(testAanmelding)
	req, _ := http.NewRequest("POST", "/aanmeldingen", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Voer de request uit
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Controleer het resultaat
	assert.Equal(t, http.StatusCreated, w.Code)

	// Controleer de response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "message")
	assert.Contains(t, response, "aanmelding")

	// Verifieer dat de mock werd aangeroepen
	mockService.AssertExpectations(t)
}

func TestCreateAanmelding_InvalidInput(t *testing.T) {
	// Setup
	router, _, handler := setupAanmeldingTest()
	router.POST("/aanmeldingen", handler.CreateAanmelding)

	// Ongeldige JSON
	req, _ := http.NewRequest("POST", "/aanmeldingen", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	// Voer de request uit
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Controleer het resultaat
	assert.Equal(t, http.StatusBadRequest, w.Code)

	// Controleer de response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "error")
}

func TestCreateAanmelding_ServiceError(t *testing.T) {
	// Setup
	router, mockService, handler := setupAanmeldingTest()
	router.POST("/aanmeldingen", handler.CreateAanmelding)

	// Test data
	testAanmelding := fixtures.GetTestAanmelding()

	// Mock verwachtingen
	mockService.On("CreateAanmelding", mock.AnythingOfType("*models.Aanmelding")).Return(errors.New("service error"))

	// Maak een request body
	jsonData, _ := json.Marshal(testAanmelding)
	req, _ := http.NewRequest("POST", "/aanmeldingen", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Voer de request uit
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Controleer het resultaat
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Controleer de response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "error")

	// Verifieer dat de mock werd aangeroepen
	mockService.AssertExpectations(t)
}

func TestGetAanmeldingen_Success(t *testing.T) {
	// Setup
	router, mockService, handler := setupAanmeldingTest()
	router.GET("/aanmeldingen", handler.GetAanmeldingen)

	// Test data
	testAanmeldingen := []models.Aanmelding{*fixtures.GetTestAanmelding()}
	var count int64 = 1

	// Mock verwachtingen
	mockService.On("GetAanmeldingen", mock.AnythingOfType("*repository.QueryParams")).Return(testAanmeldingen, nil)
	mockService.On("CountAanmeldingen", mock.AnythingOfType("*repository.QueryParams")).Return(count, nil)

	// Maak een request
	req, _ := http.NewRequest("GET", "/aanmeldingen", nil)

	// Voer de request uit
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Controleer het resultaat
	assert.Equal(t, http.StatusOK, w.Code)

	// Controleer de response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "aanmeldingen")
	assert.Contains(t, response, "total")
	assert.Contains(t, response, "page")
	assert.Contains(t, response, "page_size")

	// Verifieer dat de mocks werden aangeroepen
	mockService.AssertExpectations(t)
}

func TestGetAanmeldingByID_Success(t *testing.T) {
	// Setup
	router, mockService, handler := setupAanmeldingTest()
	router.GET("/aanmeldingen/:id", handler.GetAanmeldingByID)

	// Test data
	testAanmelding := fixtures.GetTestAanmelding()
	testID := testAanmelding.ID

	// Mock verwachtingen
	mockService.On("GetAanmeldingByID", testID).Return(testAanmelding, nil)

	// Maak een request
	req, _ := http.NewRequest("GET", "/aanmeldingen/"+testID, nil)

	// Voer de request uit
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Controleer het resultaat
	assert.Equal(t, http.StatusOK, w.Code)

	// Controleer de response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "aanmelding")

	// Verifieer dat de mock werd aangeroepen
	mockService.AssertExpectations(t)
}

func TestGetAanmeldingByID_NotFound(t *testing.T) {
	// Setup
	router, mockService, handler := setupAanmeldingTest()
	router.GET("/aanmeldingen/:id", handler.GetAanmeldingByID)

	// Test data
	testID := "non-existent-id"

	// Mock verwachtingen
	mockService.On("GetAanmeldingByID", testID).Return(nil, errors.New("not found"))

	// Maak een request
	req, _ := http.NewRequest("GET", "/aanmeldingen/"+testID, nil)

	// Voer de request uit
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Controleer het resultaat
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Controleer de response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.Contains(t, response, "error")

	// Verifieer dat de mock werd aangeroepen
	mockService.AssertExpectations(t)
}
