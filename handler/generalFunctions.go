package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// findCountry finds all entries in a struct slice that match the provided iso code
func findCountry(countries []structs.CountryData, Iso string) []structs.CountryData {
	var countryData []structs.CountryData // empty list for the final data

	// Convert to upper case to avoid case sensitivity
	Iso = strings.ToUpper(Iso)
	for _, col := range countries {
		if strings.Contains(col.IsoCode, Iso) {
			newHisData := structs.CountryData{
				Name:       col.Name,
				IsoCode:    col.IsoCode,
				Year:       col.Year,
				Percentage: col.Percentage,
			}
			countryData = append(countryData, newHisData)
		}
	}
	return countryData
}

// getCountries gets all countries from REST_countries and returns them as a slice of structs
func getCountries(isoCode []string) ([]structs.Country, error) {

	countryUrl := Assignment2.COUNTRYAPI_CODES
	// Loop through each ISO code and append the code the URL
	// Append each code to the URL with a comma delimiter
	isoCodesStr := strings.Join(isoCode, ",")
	countryUrl += isoCodesStr

	countryResponse, err := http.Get(countryUrl)
	if err != nil {
		return nil, err
	}
	defer countryResponse.Body.Close()

	// Decode the JSON response into a slice of "Country" structs
	var countryData []structs.Country
	err = json.NewDecoder(countryResponse.Body).Decode(&countryData)
	if err != nil {
		return nil, err
	}

	// Check if any countries were found
	if len(countryData) == 0 {
		return nil, errors.New("no countries found")
	}

	return countryData, nil
}
