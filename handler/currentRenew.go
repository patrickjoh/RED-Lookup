package handler

import (
	"Assignment2"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Documentation...
func RenewablesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleRenewablesGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// handleRenewablesGet
func handleRenewablesGet(w http.ResponseWriter, r *http.Request) {

	// Split url to get keyword
	keywords := strings.Split(r.URL.Path, "/") // keywords[4]  'current'

	// Error handling
	if len(keywords) < 5 || keywords[4] != "current" {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	var countryData []Assignment2.CountryData

	// If country code is provided
	if len(keywords) >= 6 {
		if len(keywords[8]) == 3 {
			fmt.Println("Length of iso code mus be 3")
		} else {
			var isoCodes []string
			isoCodes = append(isoCodes, keywords[5])
			// if (neighbor bool is set)
			// 	isoCodes = append(isoCodes, getNeighborCountry(keywords[5]))
			countryData = getOneCountry(convertCsvData(), isoCodes)
		}
	} else { // If no country code is provided
		countryData = getAllCountries(convertCsvData())
	}

	jsonResponse, err := json.Marshal(countryData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func getAllCountries(data []Assignment2.CountryData) []Assignment2.CountryData {
	currCount := data[0].Name // Current country
	currHighYear := 0         // Current highest year
	var highestRecord Assignment2.CountryData
	var retData []Assignment2.CountryData

	// Finding entries with the most recent year
	for _, current := range data {
		if current.Name == currCount { // Still same country?
			if current.Year > currHighYear { // New highest year found
				highestRecord = current
			}
		} else { // New country entered
			currCount = current.Name
			currHighYear = current.Year
			retData = append(retData, highestRecord)
		}
	}

	return retData
}

// getOneCountry recursively gets countries from a slice of strings containing iso3codes
func getOneCountry(data []Assignment2.CountryData, keywords []string) []Assignment2.CountryData {
	var retData []Assignment2.CountryData
	currHighYear := 0 // Current highest year
	var highestRecord Assignment2.CountryData

	relCountries := findCountry(data, keywords[len(keywords)-1])
	for _, current := range relCountries {
		if current.Year > currHighYear { // New highest year found
			highestRecord = current
		}
	}

	// Creating new keyword list for next iteration
	newKeywords := keywords[:len(keywords)-1]
	retData = append(retData, highestRecord)

	// Potentially appending remaining neighboring countries recursively
	if len(keywords) > 0 {
		retData = getOneCountry(data, newKeywords)
	}
	return retData
}

/*	NON RECURSIVE IMPLEMENTATION
func getOneCountry(data []Assignment2.CountryData, keyword string) []Assignment2.CountryData {

	var retData []Assignment2.CountryData
	currHighYear := 0 // Current highest year
	var highestRecord Assignment2.CountryData

	relCountries := findCountry(data, keyword)
	for _, current := range relCountries {
		if current.Year > currHighYear { // New highest year found
			highestRecord = current
		}
	}
	retData = append(retData, highestRecord)

	return retData
}
*/

func getNeighborCountry(w http.ResponseWriter, IsoCode string) ([]string, error) {
	var borderCountries []string
	// Get bordering countries data from "REST_Countries" API
	specCountryURL := Assignment2.COUNTRYAPI_CODES + IsoCode
	countryResponse, err := http.Get(specCountryURL)
	if err != nil {
		http.Error(w, "Error during request to CountryAPI", http.StatusInternalServerError)
		log.Println("Failed to get bordering country data from CountryAPI")
		return nil, err
	}
	// Close the response body when the function returns
	defer countryResponse.Body.Close()

	// Struct to hold the response for the specified country
	var specCountryData []Assignment2.Country
	// Decode the response body into the struct
	err = json.NewDecoder(countryResponse.Body).Decode(&specCountryData)
	if err != nil {
		http.Error(w, "Error during request to CountryAPI", http.StatusInternalServerError)
		log.Println("Failed to decode country data from CountryAPI, about the specified country")
		return nil, err
	}

	// Get country data from "REST_Countries" API based on the list of border countries
	borders := specCountryData[0].Border

	// Check if the country has any bordering countries
	if len(specCountryData[0].Border) == 0 {
		http.Error(w, "No bordering countries", http.StatusNotFound)
		log.Println("No bordering countries")
		return nil, err
	}
	// Get country data from "REST_Countries" API based on the list of border countries
	countryData, err := getCountries(borders)
	if err != nil {
		http.Error(w, "Error during request to CountryAPI", http.StatusInternalServerError)
		log.Println("Failed to get country data from CountryAPI, about the bordering countries")
		return nil, err
	}
	// Append the border countries to the list of border countries
	for _, country := range countryData {
		borderCountries = append(borderCountries, country.Alpha3Code)
	}
	// Return the list of border countries
	return borderCountries, err
}
