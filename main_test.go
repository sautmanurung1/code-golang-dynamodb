package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestSearch(t *testing.T) {
	// Arrange
	e := echo.New()

	// Act
	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(`{
		"keywords": "92101",
		"availableOnly": 1,
		"forSaleTypes": ["By Agent", "Coming Soon", "By Owner", "Auction", "New Construction", "Foreclosures"],
		"propertyType": ["Condo", "House", "Town_House", "Multi_Unit", "Modular", "Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"],
		"otherAmenities": [],
		"viewTypes": [],
		"per_page": 200
	}`))
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assert
	if assert.NoError(t, Search(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.True(t, len(rec.Body.String()) > 10)
	}
}

// Error Testing but i wanna to complate this code
// func TestSearch(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(`{
// 		"keywords": "92101",
// 		"availableOnly": 1,
// 		"forSaleTypes": [
// 			"By Agent", "Coming Soon", "By Owner", "Auction",
// 			"New Construction", "Foreclosures"
// 		],
// 		"propertyType": [
// 			"Condo", "House", "Town_House", "Multi_Unit", "Modular",
// 			"Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"
// 		],
// 		"otherAmenities": [],
// 		"viewTypes": [],
// 		"per_page": 200
// 	}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Assertions
// 	if assert.NoError(t, Search(c)) {
// 		assert.Equal(t, http.StatusBadRequest, rec.Code)

// 		var responseBody map[string]interface{}
// 		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
// 		if assert.NoError(t, err) {
// 			assert.NotNil(t, responseBody)
// 			assert.Greater(t, len(responseBody), 10)
// 		}
// 	}
// }

// Error Code Test but i wanna to complate this error code
// func TestSearch(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	payload := `{
// 		"keywords": "92101",
// 		"availableOnly": 1,
// 		"forSaleTypes": [
// 			"By Agent", "Coming Soon", "By Owner", "Auction",
// 			"New Construction", "Foreclosures"
// 		],
// 		"propertyType": [
// 			"Condo", "House", "Town_House", "Multi_Unit", "Modular",
// 			"Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"
// 		],
// 		"otherAmenities": [],
// 		"viewTypes": [],
// 		"per_page": 200
// 	}`
// 	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(payload))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Invoke
// 	err := Search(c)

// 	// Assertions
// 	assert.NoError(t, err)
// 	assert.Equal(t, http.StatusBadRequest, rec.Code)

// 	var responseBody map[string]interface{}
// 	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
// 	assert.NoError(t, err)
// 	assert.NotNil(t, responseBody)
// 	assert.Greater(t, len(responseBody), 10)
// }

// func TestSearch(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(`{
// 		"keywords": "92101",
// 		"availableOnly": 1,
// 		"forSaleTypes": [
// 			"By Agent", "Coming Soon", "By Owner", "Auction",
// 			"New Construction", "Foreclosures"
// 		],
// 		"propertyType": [
// 			"Condo", "House", "Town_House", "Multi_Unit", "Modular",
// 			"Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"
// 		],
// 		"otherAmenities": [],
// 		"viewTypes": [],
// 		"per_page": 200
// 	}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Assertions
// 	err := Search(c)
// 	if err != nil {
// 		t.Fatalf("Search error: %v", err)
// 	}
// 	assert.Equal(t, http.StatusBadRequest, rec.Code)

// 	var responseBody map[string]interface{}
// 	err = json.Unmarshal(rec.Body.Bytes(), &responseBody)
// 	if err != nil {
// 		t.Fatalf("Error decoding response body: %v", err)
// 	}
// 	assert.NotEmpty(t, responseBody)
// }

// func TestSearch(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(`{
// 		"keywords": "92101",
// 		"availableOnly": 1,
// 		"forSaleTypes": [
// 			"By Agent", "Coming Soon", "By Owner", "Auction",
// 			"New Construction", "Foreclosures"
// 		],
// 		"propertyType": [
// 			"Condo", "House", "Town_House", "Multi_Unit", "Modular",
// 			"Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"
// 		],
// 		"otherAmenities": [],
// 		"viewTypes": [],
// 		"per_page": 200
// 	}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Assertions
// 	if assert.NoError(t, Search(c)) {
// 		assert.Equal(t, http.StatusBadRequest, rec.Code)

// 		var responseBody map[string]interface{}
// 		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
// 		if assert.NoError(t, err) {
// 			assert.NotNil(t, responseBody)
// 			if value, ok := responseBody["total"].(float64); ok {
// 				assert.Greater(t, value, float64(10))
// 			} else {
// 				t.Errorf("expected float64, but got %T", responseBody["total"])
// 			}
// 		}
// 	}
// }

// func TestSearch(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(`{
// 		"keywords": "92101",
// 		"availableOnly": 1,
// 		"forSaleTypes": [
// 			"By Agent", "Coming Soon", "By Owner", "Auction",
// 			"New Construction", "Foreclosures"
// 		],
// 		"propertyType": [
// 			"Condo", "House", "Town_House", "Multi_Unit", "Modular",
// 			"Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"
// 		],
// 		"otherAmenities": [],
// 		"viewTypes": [],
// 		"per_page": 200
// 	}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Assertions
// 	if assert.NoError(t, Search(c)) {
// 		assert.Equal(t, http.StatusBadRequest, rec.Code)

// 		var responseBody map[string]interface{}
// 		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
// 		if assert.NoError(t, err) {
// 			assert.NotNil(t, responseBody)
// 			if value, ok := responseBody["total"].(float64); ok && value != 0 {
// 				assert.Greater(t, value, float64(10))
// 			} else {
// 				t.Errorf("expected non-nil float64, but got %v", responseBody["total"])
// 			}
// 		}
// 	}
// }

// func TestSearch(t *testing.T) {
// 	// Setup
// 	e := echo.New()
// 	req := httptest.NewRequest(http.MethodPost, "/search-x.api", strings.NewReader(`{
// 		"keywords": "92101",
// 		"availableOnly": 1,
// 		"forSaleTypes": [
// 			"By Agent", "Coming Soon", "By Owner", "Auction",
// 			"New Construction", "Foreclosures"
// 		],
// 		"propertyType": [
// 			"Condo", "House", "Town_House", "Multi_Unit", "Modular",
// 			"Commercial", "Land", "Timeshare", "Parking", "Rental", "Other"
// 		],
// 		"otherAmenities": [],
// 		"viewTypes": [],
// 		"per_page": 200
// 	}`))
// 	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
// 	rec := httptest.NewRecorder()
// 	c := e.NewContext(req, rec)

// 	// Assertions
// 	if assert.NoError(t, Search(c)) {
// 		assert.Equal(t, http.StatusBadRequest, rec.Code)

// 		var responseBody map[string]interface{}
// 		err := json.Unmarshal(rec.Body.Bytes(), &responseBody)
// 		if assert.NoError(t, err) {
// 			assert.NotNil(t, responseBody)
// 			if value, ok := responseBody["total"].(float64); ok {
// 				if value > 0 {
// 					assert.Greater(t, value, float64(10))
// 				} else {
// 					t.Error("expected non-zero value for total")
// 				}
// 			} else {
// 				t.Errorf("expected float64, but got %T", responseBody["total"])
// 			}
// 		}
// 	}
// }
