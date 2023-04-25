package handler

import (
	"Assignment2"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// RenewablesHandler handles get requests for current renewables
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

// handleRenewablesGet calls either getAllCountries or getOneCountry, depending on the parameters
// provided by the user, to retrieve the most recent renewables percentages for countries.
func handleRenewablesGet(w http.ResponseWriter, r *http.Request) {

	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Error handling
	if len(parts) < 5 || parts[4] != "current" {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		log.Println("Malformed URL")
		return
	}

	// Get the neighbour parameter from the URL query parameters
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

	// If the country code is provided
	if len(parts) == 6 {
		if len(parts[5]) != 3 {
			http.Error(w, "Iso code must be 3 letter long", http.StatusBadRequest)
			log.Println("Iso code not 3 letters long: ", parts[5])
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

			if len(isoCodes) < 1 {
				log.Println("No iso code found")
				http.Error(w, "No iso code found", http.StatusBadRequest)
				return
			}

			// Fetching data for one country, possibly with neighbors
			countryData = getOneCountry(convertCsvData(), isoCodes)

			// No country found with matching data
			if len(countryData) < 2 && countryData[0].Name == "" {
				log.Println("No country found")
				http.Error(w, "MNo country found", http.StatusBadRequest)
				return
			}
		}
	} else if len(parts) == 5 { // If no country code is provided
		countryData = getAllCountries(convertCsvData())

	} else { // If the URL is malformed
		log.Println("Malformed URL")
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(countryData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// getAllCountries retrieves the most recent entries for all countries.
func getAllCountries(data []Assignment2.CountryData) []Assignment2.CountryData {
	mostRecentData := make(map[string]Assignment2.CountryData)

	// Finding entries with the most recent year
	for _, current := range data {
		existingRecord, exists := mostRecentData[current.Name]
		if !exists || current.Year > existingRecord.Year {
			mostRecentData[current.Name] = current
		}
	}

	var retData []Assignment2.CountryData // Slice to hold the most recent entry for each country
	for _, record := range mostRecentData {
		//UpdateAndInvoke(record.IsoCode)   // UWU not certain if work??
		retData = append(retData, record) // Adding the most recent entry for each country to a slice
	}

	// Sorting the retData slice alphabetically by country name
	sort.Slice(retData, func(i, j int) bool {
		return retData[i].Name < retData[j].Name
	})

	return retData
}

// getOneCountry retrieves the latest entries for a given country, and if the neighbours
// parameter is set to true, it also retrieves the latest entries for the countries
// that share a border with the given country.
func getOneCountry(data []Assignment2.CountryData, isoCodes []string) []Assignment2.CountryData {

	if len(isoCodes) < 1 {
		return []Assignment2.CountryData{}
	}

	var returnData []Assignment2.CountryData
	currentHighestYear := 0                         // The currently highest year found
	var currentHighestEntry Assignment2.CountryData // The struct with the currently highest year

	for _, iso := range isoCodes {
		relCountries := findCountry(data, iso)
		for _, current := range relCountries {
			// New highest year found
			if current.Year > currentHighestYear {
				currentHighestEntry = current
				currentHighestYear = current.Year
			}
		}
		returnData = append(returnData, currentHighestEntry)
		currentHighestYear = 0

		// Update counter for webhook invocation
		//UpdateAndInvoke(iso)
	}

	return returnData
}

// getNeighborCountry searches for and returns ISO codes of all countries that share a
// border with the countries specified by the IsoCode parameter. Bordering countries
// are returned as a slice of strings containing their ISO codes.
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
