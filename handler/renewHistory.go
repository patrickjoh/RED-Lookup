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
	// Split url to get keyword
	urlKeywords := strings.Split(r.URL.Path, "/")

	iso := urlKeywords[5]   // Get country isoCode from url
	query := r.URL.RawQuery // Get the queries from url

	// Iso must either be nothing or consist of three letters
	if len(iso) != 3 && len(iso) != 0 {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		return
	}

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
	if begin > end && end != "" {
		log.Printf("begining year (%s) > ending year(%s)", begin, end)
		http.Error(w, "Incorrect use of year. Try history/{country?}{?begin=year&end=year?}", http.StatusBadRequest)
		return
	}

	var rangedCountries []structs.CountryData
	var startYear = 0
	var endYear = 3000

	// If no Iso is given print all countries mean percentage else print one country's history
	if iso == "" {
		rangedCountries = Assignment2.ConvertCsvData()
		if begin != "" && end == "" {
			startYear, _ = strconv.Atoi(begin)
			rangedCountries = getFromBeginToEnd(startYear, endYear, Assignment2.ConvertCsvData())
		} else if end != "" && begin == "" {
			endYear, _ = strconv.Atoi(end)
			rangedCountries = getFromBeginToEnd(startYear, endYear, Assignment2.ConvertCsvData())
		} else if end != "" && begin != "" {
			startYear, _ = strconv.Atoi(begin)
			endYear, _ = strconv.Atoi(end)
			rangedCountries = getFromBeginToEnd(startYear, endYear, Assignment2.ConvertCsvData())
		}

		resp := getAllCountriesMean(rangedCountries) // get all countries mean percentage

		// If user want to sort by percentage
		if sortByValue == "true" {
			// Sorting the countData slice from lowest to highest by country percentage
			sort.Slice(resp, func(i, j int) bool {
				return resp[i].Percentage < resp[j].Percentage
			})
		}

		jsonResponse, err := json.Marshal(resp)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Println("Error sending response")
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}

	} else { // Only data for one country is returned

		rangedCountries = findCountry(Assignment2.ConvertCsvData(), iso)
		if begin != "" && end == "" {
			startYear, _ = strconv.Atoi(begin)
			rangedCountries = getFromBeginToEnd(startYear, endYear, rangedCountries)
		} else if end != "" && begin == "" {
			endYear, _ = strconv.Atoi(end)
			rangedCountries = getFromBeginToEnd(startYear, endYear, rangedCountries)
		} else if begin != "" && end != "" {
			startYear, _ = strconv.Atoi(begin)
			endYear, _ = strconv.Atoi(end)
			rangedCountries = getFromBeginToEnd(startYear, endYear, rangedCountries)
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
		}

		jsonResponse, err := json.Marshal(rangedCountries)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(jsonResponse)
		if err != nil {
			log.Println("Error sending response")
			http.Error(w, "Error sending response", http.StatusInternalServerError)
			return
		}
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
			UpdateAndInvoke(countryMean.IsoCode) // UWU maybe work, maybe not???
			// Appends country to slice of countries
			retData = append(retData, countryMean)
		}
	}
	return retData
}
