package handler

import (
	"Assignment2"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Sample data for test functions
var sampleData = []Assignment2.CountryData{
	{Name: "Azerbajan", IsoCode: "AZE", Year: 2000, Percentage: 1.5},
	{Name: "Azerbajan", IsoCode: "AZE", Year: 2021, Percentage: 54.2},
	{Name: "Chad", IsoCode: "TCD", Year: 2021, Percentage: 99.1},
	{Name: "Germany", IsoCode: "DEU", Year: 2000, Percentage: 13.2},
	{Name: "Germany", IsoCode: "DEU", Year: 2021, Percentage: 42.2},
	{Name: "Thailand", IsoCode: "THA", Year: 2021, Percentage: 99.1},
}

// TestGetAllCountries tests the getAllCountries function*/
func TestGetAllCountries(t *testing.T) {
	// Expected output
	expected := []Assignment2.CountryData{
		{Name: "Azerbajan", IsoCode: "AZE", Year: 2021, Percentage: 54.2},
		{Name: "Chad", IsoCode: "TCD", Year: 2021, Percentage: 99.1},
		{Name: "Germany", IsoCode: "DEU", Year: 2021, Percentage: 42.2},
		{Name: "Thailand", IsoCode: "THA", Year: 2021, Percentage: 99.1},
	}

	result := getAllCountries(sampleData)

	assert.Equal(t, expected, result)
}

// TestGetOneCountry tests the getOneCountry function
func TestGetOneCountry(t *testing.T) {
	// Expected output
	expected := []Assignment2.CountryData{
		{Name: "Azerbajan", Year: 2021, IsoCode: "AZE", Percentage: 54.2},
		{Name: "Germany", Year: 2021, IsoCode: "DEU", Percentage: 42.2},
	}

	isoCodes := []string{"AZE", "DEU"}

	result := getOneCountry(sampleData, isoCodes)

	assert.Equal(t, expected, result)
}
