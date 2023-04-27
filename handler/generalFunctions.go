package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

// findCountry finds all entries in a struct slice that match the provided iso code or country name
func findCountry(countries []structs.CountryData, country string) []structs.CountryData {
	var countryData []structs.CountryData // empty list for the final data

	country = strings.ToUpper(country) // Convert input to upper case to avoid case sensitivity

	useIsoCode := len(country) == 3 // If the input is 3 characters long, it is an ISO code

	for _, col := range countries {
		// Determine the search field based on the length of the input string
		searchField := strings.ToUpper(col.Name) // Convert current country name to match input
		if useIsoCode {                          // If the input is an ISO code, search for the ISO code
			searchField = col.IsoCode
		}

		// Check if the search field contains the input country string
		if strings.Contains(searchField, country) {
			newHisData := structs.CountryData{
				Name:       col.Name,
				IsoCode:    col.IsoCode,
				Year:       col.Year,
				Percentage: col.Percentage,
			}
			countryData = append(countryData, newHisData) // Add the matching entry to the final list
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
