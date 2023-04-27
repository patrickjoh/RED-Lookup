package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// HistoryHandler handles get requests for history endpoint
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleHistoryGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// handleHistoryGet returns renewables for a given country in a provided range, or
// if no country is specified returns the mean of the renewables for all countries by calling.
// It also sorts the data by percentage in descending order
func handleHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	query := r.URL.RawQuery // Get the queries from url

	// Parse the query string into a map
	params, err := url.ParseQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract values from the parameters
	begin := params.Get("begin")
	end := params.Get("end")
	sortByValue := params.Get("sortByValue")

	var rangedCountries []structs.CountryData // Slice to hold countries in range
	var startYear = 0                         // Default value for start year
	var endYear = 3000                        // Default value for end year

	// If no iso or name is given, print all countries mean percentage, else print one country's history
	if len(parts) == 5 {

		if begin != "" && end == "" { // If only begin year is given
			startYear, _ = strconv.Atoi(begin)
			rangedCountries = getFromBeginToEnd(startYear, endYear, Assignment2.CSVData)
		} else if end != "" && begin == "" { // If only end year is given
			endYear, _ = strconv.Atoi(end)
			rangedCountries = getFromBeginToEnd(startYear, endYear, Assignment2.CSVData)
		} else if end != "" && begin != "" { // If both begin and end year is given
			startYear, _ = strconv.Atoi(begin)
			endYear, _ = strconv.Atoi(end)
			if startYear < endYear { // Check if begin year is smaller than end year
				rangedCountries = getFromBeginToEnd(startYear, endYear, Assignment2.CSVData)
			} else { // Else return error
				http.Error(w, "Incorrect use of year. 'begin' must be smaller than 'end'. Try history/{country?}{?begin=year&end=year?}", http.StatusBadRequest)
				return
			}
		} else {
			rangedCountries = Assignment2.CSVData // If no country is given, get all countries
		}
		resp := getAllCountriesMean(rangedCountries) // Get mean percentage for  all countries

		// If user want to sort by percentage
		if sortByValue == "true" {
			// Sorting the countData slice from lowest to highest by country percentage
			sort.Slice(resp, func(i, j int) bool {
				return resp[i].Percentage < resp[j].Percentage
			})
		}
		// Marshall the response into a JSON string
		jsonResponse, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Println("Error sending response")
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
	} else if len(parts) == 6 { // If country name or isoCode is given
		countrySearch := parts[5] // Get country isoCode from url

		// Input must either be nothing or consist of more than three letters
		if len(countrySearch) < 3 {
			http.Error(w, "Input must either be nothing or consist of three or more letters", http.StatusBadRequest)
			return
		}

		// Find country data for the given country
		rangedCountries = findCountry(Assignment2.CSVData, countrySearch)
		if begin != "" && end == "" { // If only begin year is given
			startYear, _ = strconv.Atoi(begin)
			rangedCountries = getFromBeginToEnd(startYear, endYear, rangedCountries)
		} else if end != "" && begin == "" { // If only end year is given
			endYear, _ = strconv.Atoi(end)
			rangedCountries = getFromBeginToEnd(startYear, endYear, rangedCountries)
		} else if begin != "" && end != "" { // If both begin and end year is given
			startYear, _ = strconv.Atoi(begin)
			endYear, _ = strconv.Atoi(end)
			if startYear < endYear { // Check if begin year is smaller than end year
				rangedCountries = getFromBeginToEnd(startYear, endYear, rangedCountries)
			} else { // Else return error
				http.Error(w, "Incorrect use of year. 'begin' must be smaller than 'end'. Try history/{country?}{?begin=year&end=year?}", http.StatusBadRequest)
				return
			}
		}

		// If user want to sort by percentage
		if sortByValue == "true" {
			// Sorting the countData slice from lowest to highest by country percentage
			sort.Slice(rangedCountries, func(i, j int) bool {
				return rangedCountries[i].Percentage < rangedCountries[j].Percentage
			})
		}

		// If no country is found
		if len(rangedCountries) < 1 {
			http.Error(w, "No entry with matching credentials found", http.StatusNotFound)
			return
		} else {
			// Update counter for webhook invocation
			UpdateAndInvoke(rangedCountries[0].IsoCode)
		}
		// Marshall the response into a JSON string
		jsonResponse, err := json.Marshal(rangedCountries)
		if err != nil {
			log.Fatal(err)
		}
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Println("Error sending response")
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
	} else { // If url is malformed
		log.Println("Malformed URL")
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}
}

// getFromBeginToEnd finds all country data for all countries within a range
func getFromBeginToEnd(begin int, end int, countryIterators []structs.CountryData) []structs.CountryData {

	var returnData []structs.CountryData
	// Find country's history from year(begin to end)
	for _, col := range countryIterators {
		if col.Year <= end && col.Year >= begin && len(col.IsoCode) == 3 {
			newHisData := structs.CountryData{
				Name:       col.Name,
				IsoCode:    col.IsoCode,
				Year:       col.Year,
				Percentage: col.Percentage,
			}
			returnData = append(returnData, newHisData)
		}
	}
	return returnData
}

// getAllCountriesMean gets all countries, checks for redundancy and returns a struct of
// all countries with mean percentage of their renewable energy
func getAllCountriesMean(countries []structs.CountryData) []structs.CountryMean {

	var retData []structs.CountryMean
	lastCountry := ""
	// Going through all countries
	for _, current := range countries {
		// Makes sure there are no repeated countries
		if current.IsoCode != lastCountry && len(current.IsoCode) == 3 {
			lastCountry = current.IsoCode
			// Variable that stores mean percentage
			var mean float64
			// Calls findCountry functions to get all instances of one country
			countryHistory := findCountry(countries, current.IsoCode)
			// Loops through all instances and adds their percentage to the mean variable
			for _, currentYear := range countryHistory {
				mean += currentYear.Percentage
			}
			// Converts length of countryHistory slice into a float/number of instances
			numberInstances := float64(len(countryHistory))
			// Calculates mean percentage of a country
			mean = mean / numberInstances
			// Initiates a country struct with the mean percentage
			countryMean := structs.CountryMean{
				Name:       current.Name,
				IsoCode:    current.IsoCode,
				Percentage: mean,
			}
			// Appends country to slice of countries
			retData = append(retData, countryMean)
		}
	}
	return retData
}
