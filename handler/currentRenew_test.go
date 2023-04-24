package handler

import (
	"Assignment2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sample data for test functions
var sampleData = []Assignment2.CountryData{
	{Name: "Azerbaijan", IsoCode: "AZE", Year: 2000, Percentage: 1.5},
	{Name: "Azerbaijan", IsoCode: "AZE", Year: 2021, Percentage: 54.2},
	{Name: "Chad", IsoCode: "TCD", Year: 2021, Percentage: 99.1},
	{Name: "Germany", IsoCode: "DEU", Year: 2000, Percentage: 13.2},
	{Name: "Germany", IsoCode: "DEU", Year: 2021, Percentage: 42.2},
	{Name: "Thailand", IsoCode: "THA", Year: 2021, Percentage: 99.1},
}

// TestGetAllCountries tests the getAllCountries function
func TestGetAllCountries(t *testing.T) {
	InitFirebase()
	// Expected output
	expected := []Assignment2.CountryData{
		{Name: "Azerbaijan", IsoCode: "AZE", Year: 2021, Percentage: 54.2},
		{Name: "Chad", IsoCode: "TCD", Year: 2021, Percentage: 99.1},
		{Name: "Germany", IsoCode: "DEU", Year: 2021, Percentage: 42.2},
		{Name: "Thailand", IsoCode: "THA", Year: 2021, Percentage: 99.1},
	}

	result := getAllCountries(sampleData)

	assert.Equal(t, expected, result)
}

// TestGetAllCountries tests the getAllCountries function
func TestGetAllCountriesNoData(t *testing.T) {
	InitFirebase()
	// Expected output
	data := []Assignment2.CountryData{}

	result := getAllCountries(data)

	assert.Equal(t, 0, len(result))
}

// TestGetOneCountry tests the getOneCountry function
func TestGetOneCountry(t *testing.T) {
	InitFirebase()
	// Expected output
	expected := []Assignment2.CountryData{
		{Name: "Azerbaijan", Year: 2021, IsoCode: "AZE", Percentage: 54.2},
		{Name: "Germany", Year: 2021, IsoCode: "DEU", Percentage: 42.2},
	}

	isoCodes := []string{"AZE", "DEU"}

	result := getOneCountry(sampleData, isoCodes)

	assert.Equal(t, expected, result)
}

// TestCurrentRenewTooManyParts tests current endpoint if there are too many parts
func TestCurrentRenewTooManyParts(t *testing.T) {
	InitFirebase()
	request, err := http.NewRequest(http.MethodGet, Assignment2.CURRENT_PATH+"AAAAA/hdjfhdjfh", nil)

	assert.Nil(t, err)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RenewablesHandler)

	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

// TestGetOneCountryWrongMethod tests current endpoint if the method is not implemented
func TestGetOneCountryWrongMethod(t *testing.T) {
	InitFirebase()
	request, err := http.NewRequest(http.MethodPost, Assignment2.CURRENT_PATH, nil)

	assert.Nil(t, err)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RenewablesHandler)

	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusNotImplemented, responseRecorder.Code)
}

// TestGetOneCountryTooLongIso tests current endpoint if the iso code is not three letters long
func TestGetOneCountryTooLongIso(t *testing.T) {
	InitFirebase()
	request, err := http.NewRequest(http.MethodGet, Assignment2.CURRENT_PATH+"abababa", nil)

	assert.Nil(t, err)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RenewablesHandler)

	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}

// TestGetOneCountryNonExistingIso tests the current endpoint if a non-existing iso code is provided
func TestGetOneCountryNonExistingIso(t *testing.T) {
	InitFirebase()
	request, err := http.NewRequest(http.MethodGet, Assignment2.CURRENT_PATH+"aaa", nil)

	assert.Nil(t, err)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(RenewablesHandler)

	handler.ServeHTTP(responseRecorder, request)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
}
