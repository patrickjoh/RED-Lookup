package handler

import (
	"Assignment2"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Method handler for handleRenewablesGet
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
	parts := strings.Split(r.URL.Path, "/") // parts[4]  'current'

	// Error handling
	if len(parts) < 5 || parts[4] != "current" {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	// Get the neighbor parameter from the URL query parameters
	neighbourStr := r.URL.Query().Get("neighbours")
	neighbourBool, err := strconv.ParseBool(neighbourStr)
	if err != nil && neighbourStr != "" {
		http.Error(w, "Invalid neighbours parameter", http.StatusBadRequest)
		log.Println("Invalid neighbours parameter")
		return
	}

	// If neighbourStr is empty, neighbourBool is false
	if neighbourStr == "" {
		neighbourBool = false
	}
	var countryData []Assignment2.CountryData
	// If country code is provided
	if len(parts) >= 6 {
		if len(parts[5]) != 3 {
		} else {
			var isoCodes []string
			isoCodes = append(isoCodes, parts[5])
			if neighbourBool {
				bordering, err := getNeighborCountry(w, parts[5])
				if err != nil {
					log.Fatal(err)
				} // Adding ISO codes for neighboring countries to list
				for _, borders := range bordering {
					isoCodes = append(isoCodes, borders)
				}
			}
			// Fetching data for one country, possibly with neighbors
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
func getOneCountry(data []Assignment2.CountryData, isoCodes []string) []Assignment2.CountryData {
	var returnData []Assignment2.CountryData
	currentHighestYear := 0                         // The currently highest year found
	var currentHighestEntry Assignment2.CountryData // The struct with the currently highest year

	for _, iso := range isoCodes {
		relCountries := findCountry(data, iso)
		for _, current := range relCountries {
			if current.Year > currentHighestYear { // New highest year found
				currentHighestEntry = current
				currentHighestYear = current.Year
			}
		}
		returnData = append(returnData, currentHighestEntry)
		currentHighestYear = 0
	}

	return returnData
}

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
