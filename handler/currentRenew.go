package handler

import (
	"Assignment2"
	"Assignment2/structs"
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

// handleRenewablesGet calls either getAllCountries or getCountries, depending on the parameters
// provided by the user, to retrieve the most recent renewables percentages for countries.
func handleRenewablesGet(w http.ResponseWriter, r *http.Request) {

	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Get the neighbour parameter from the URL query parameters
	neighbourStr := r.URL.Query().Get("neighbours")
	neighbourBool, err := strconv.ParseBool(neighbourStr)
	if err != nil && neighbourStr != "" {
		http.Error(w, "Invalid neighbours parameter", http.StatusBadRequest)
		return
	}

	// If neighbourStr is empty, neighbourBool is false
	if neighbourStr == "" {
		neighbourBool = false
	}
	var countryData []structs.CountryData

	if len(parts) == 6 { // If country name or code is provided
		if len(parts[5]) < 3 {
			http.Error(w, "Input must me either a 3 letter ISO code or a country name", http.StatusBadRequest)
		} else {
			var countries []string
			countries = append(countries, parts[5])
			if neighbourBool {
				bordering, err := getNeighbourCountry(w, parts[5])
				if err != nil {
					http.Error(w, "No country found", http.StatusBadRequest)
					return
				} // Adding ISO codes for neighboring countries to list
				for _, borders := range bordering {
					countries = append(countries, borders)
				}
			}

			// Fetching data for one country, possibly with neighbors
			countryData = getCountries(Assignment2.CSVData, countries)

			// No country found with matching data
			if len(countryData) < 2 && countryData[0].Name == "" {
				log.Println("No country found")
				http.Error(w, "No country found", http.StatusBadRequest)
				return
			}
		}
	} else if len(parts) == 5 { // If no country code is provided
		countryData = getAllCountries(Assignment2.CSVData)

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
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Println("Error sending response")
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}
}

// getAllCountries retrieves the most recent entries for all countries.
func getAllCountries(data []structs.CountryData) []structs.CountryData {
	mostRecentData := make(map[string]structs.CountryData)

	// Finding entries with the most recent year
	for _, current := range data {
		existingRecord, exists := mostRecentData[current.Name]
		if !exists || current.Year > existingRecord.Year {
			mostRecentData[current.Name] = current
		}
	}

	var retData []structs.CountryData // Slice to hold the most recent entry for each country
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

// getCountries retrieves the latest entries for a given country, and if the neighbours
// parameter is set to true, it also retrieves the latest entries for the countries
// that share a border with the given country.
func getCountries(data []structs.CountryData, countrySearch []string) []structs.CountryData {

	if len(countrySearch) < 1 {
		return []structs.CountryData{}
	}

	var returnData []structs.CountryData        // Slice to hold the country data to be returned
	currentHighestYear := 0                     // The currently highest year found
	var currentHighestEntry structs.CountryData // The struct with the currently highest year

	for _, iso := range countrySearch {
		relCountries := findCountry(data, iso)
		for _, current := range relCountries {
			// New highest year found
			if current.Year > currentHighestYear {
				currentHighestEntry = current
				currentHighestYear = current.Year
			}
		}
		// Prevent duplicates
		if len(returnData) == 0 || currentHighestEntry.Name != returnData[len(returnData)-1].Name {
			returnData = append(returnData, currentHighestEntry)
		}
		currentHighestYear = 0
		if len(relCountries) > 1 { // If a country was found
			// Update counter for webhook invocation
			UpdateAndInvoke(relCountries[0].IsoCode)
		}
	}

	return returnData
}

// getNeighbourCountry searches for and returns ISO codes of all countries that share a
// border with the countries specified by the IsoCode parameter. Bordering countries
// are returned as a slice of strings containing their ISO codes.
func getNeighbourCountry(w http.ResponseWriter, searchCountry string) ([]string, error) {
	var borderCountries []string
	var searchCountryURL string
	// Get bordering countries data from "REST_Countries" API
	if len(searchCountry) == 3 {
		searchCountryURL = Assignment2.COUNTRYAPI_CODES + searchCountry
	} else {
		searchCountryURL = Assignment2.COUNTRYAPI_NAME + searchCountry
	}
	countryResponse, err := http.Get(searchCountryURL)
	if err != nil {
		http.Error(w, "Error during request to CountryAPI", http.StatusInternalServerError)
		log.Println("Failed to get bordering country data from CountryAPI: ")
		return nil, err
	}

	// If no country was found in the CountryAPI
	if countryResponse.StatusCode == 404 {
		return nil, err // Return empty slice
	}

	// Close the response body when the function returns
	defer func() {
		if countryResponse != nil {
			if err := countryResponse.Body.Close(); err != nil {
				log.Println("Error closing response body:", err)
			}
		}
	}()

	// Struct to hold the response for the specified country
	var specCountryData []structs.Country
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
		return nil, err
	}
	// Get country data from "REST_Countries" API based on the list of border countries
	countryData, err := getRESTCountries(borders)
	if err != nil {
		http.Error(w, "Error during request to CountryAPI", http.StatusInternalServerError)
		log.Println("Failed to get country data from CountryAPI, about the bordering countries")
		return nil, err
	}
	// Append the border countries to the list of border countries
	for _, country := range countryData {
		borderCountries = append(borderCountries, country.Alpha3Code)
		// Update counter for webhook invocation
		UpdateAndInvoke(country.Alpha3Code)
	}
	// Return the list of border countries
	return borderCountries, err
}
