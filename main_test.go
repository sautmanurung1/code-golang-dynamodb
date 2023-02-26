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
