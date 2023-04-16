package handler

import (
	"Assignment2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllCountries(t *testing.T) {
	// Sample data
	data := []Assignment2.CountryData{
		{Name: "Country1", Year: 2000, IsoCode: "20", Percentage: 1.5},
		{Name: "Country1", Year: 2005, IsoCode: "20", Percentage: 54.2},
		{Name: "Country2", Year: 2000, IsoCode: "40", Percentage: 13.2},
		{Name: "Country2", Year: 2010, IsoCode: "40", Percentage: 42.2},
		{Name: "Country3", Year: 2010, IsoCode: "50", Percentage: 99.1},
	}

	// Expected output
	expected := []Assignment2.CountryData{
		{Name: "Country1", Year: 2005, IsoCode: "20", Percentage: 54.2},
		{Name: "Country2", Year: 2010, IsoCode: "40", Percentage: 42.2},
		{Name: "Country3", Year: 2010, IsoCode: "50", Percentage: 99.1},
	}

	result := getAllCountries(data)

	assert.Equal(t, expected, result)
}

func TestGetOneCountry(t *testing.T) {
	// Sample data
	data := []Assignment2.CountryData{
		{Name: "Country1", Year: 2000, IsoCode: "20", Percentage: 1.5},
		{Name: "Country1", Year: 2005, IsoCode: "20", Percentage: 54.2},
		{Name: "Country2", Year: 2000, IsoCode: "40", Percentage: 13.2},
		{Name: "Country2", Year: 2010, IsoCode: "40", Percentage: 42.2},
		{Name: "Country3", Year: 2010, IsoCode: "50", Percentage: 99.1},
	}

	// Expected output
	expected := []Assignment2.CountryData{
		{Name: "Country1", Year: 2005, IsoCode: "20", Percentage: 54.2},
		{Name: "Country2", Year: 2010, IsoCode: "40", Percentage: 42.2},
	}

	isoCodes := []string{"20", "40"}

	result := getOneCountry(data, isoCodes)

	assert.Equal(t, expected, result)
}
