package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	// Arrange
	e := echo.New()

	// Act
	requestBody := map[string]interface{}{
		"keywords":      "92101",
		"availableOnly": 1,
		"forSaleTypes": []string{
			"By Agent",
			"Coming Soon",
			"By Owner",
			"Auction",
			"New Construction",
			"Foreclosures",
		},
		"propertyType": []string{
			"Condo",
			"House",
			"Town_House",
			"Multi_Unit",
			"Modular",
			"Commercial",
			"Land",
			"Timeshare",
			"Parking",
			"Rental",
			"Other",
		},
		"otherAmenities": []interface{}{},
		"viewTypes":      []interface{}{},
		"per_page":       200,
	}
	reqBody, err := json.Marshal(requestBody)
	if err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/search-x.api", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	err = Search(c)

	// Assert
	if assert.NoError(t, err, nil) {
		if strings.Contains(rec.Body.String(), "error") {
			assert.Equal(t, http.StatusBadRequest, rec.Code)
		} else {
			assert.Equal(t, http.StatusOK, rec.Code)
		}
	} else {
		assert.Error(t, err, errors.New("error the code"))
	}
}
