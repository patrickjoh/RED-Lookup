package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Sample data for test functions
var sampleData = []structs.CountryData{
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
	expected := []structs.CountryData{
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
	// Expected output
	data := []structs.CountryData{}

	result := getAllCountries(data)

	assert.Equal(t, 0, len(result))
}

// TestGetOneCountry tests the getOneCountry function
func TestGetOneCountry(t *testing.T) {
	// Expected output
	expected := []structs.CountryData{
		{Name: "Azerbaijan", Year: 2021, IsoCode: "AZE", Percentage: 54.2},
		{Name: "Germany", Year: 2021, IsoCode: "DEU", Percentage: 42.2},
	}

	isoCodes := []string{"AZE", "DEU"}

	result := getOneCountry(sampleData, isoCodes)

	assert.Equal(t, expected, result)
}

// TestCurrentRenew tests current endpoint with no iso code
func TestCurrentRenew(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// TestCurrentRenewTooManyParts tests current endpoint if there are too many parts
func TestCurrentRenewTooManyParts(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH + "AAAAA/hdjfhdjfh"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// TestGetOneCountryWrongMethod tests current endpoint if the method is not implemented
func TestGetOneCountryWrongMethod(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH
	response, err := http.Post(url, "", nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotImplemented, response.StatusCode)
}

// TestGetOneCountryTooLongIso tests current endpoint if the iso code is not three letters long
func TestGetOneCountryTooLongIso(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH + "abababba"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// TestGetOneCountryNonExistingIso tests the current endpoint if a non-existing iso code is provided
func TestGetOneCountryNonExistingIso(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH + "aaa"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// TestGetOneCountryNonExistingIso tests the current endpoint if a non-existing iso code is provided
func TestGetOneCountryWrongPath(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// TestGetNeighborCountry tests the current endpoint if ...
func TestGetNeighborCountry(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH + "AZE?neighbours=true"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// TestRenewablesHandlerInvalidBool tests the current endpoint with invalid bool
func TestRenewablesHandlerInvalidBool(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(RenewablesHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH + "AZE?neighbours=adhfkdjshhsd"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}
