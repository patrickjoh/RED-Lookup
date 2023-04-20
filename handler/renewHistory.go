package handler

import (
	"Assignment2"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

// Documentation...
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

// handleHistoryGet
func handleHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Split url to get keyword
	urlKeywords := strings.Split(r.URL.Path, "/")

	iso := urlKeywords[5]   // Get country isoCode from url
	query := r.URL.RawQuery // Get the queries from url

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

	// Error and logic check for beginning and end of year
	if begin == "" && end == "" {
		begin = "0"
		end = "3000"
	} else if end == "" {
		end = begin
	} else if begin == "" {
		begin = end
	} else if begin > end {
		log.Printf("begining year (%s) > ending year(%s)", begin, end)
		http.Error(w, "Incorrect use of year", http.StatusBadRequest)
		return
	}

	var countData []Assignment2.CountryData // Empty list for the final data
	startYear, _ := strconv.Atoi(begin)     // Convert beginning year to int
	endYear, _ := strconv.Atoi(end)         // Convert end year to int
	countryIterators := convertCsvData()    // Read all countries data from csv

	// if Iso Code is given only one country's data is checked
	if iso != "" {
		countryIterators = findCountry(countryIterators, iso) // splice of one country's history
	}
	// Find country's history from year(begin to end)
	for _, col := range countryIterators {
		if col.Year <= endYear && col.Year >= startYear {
			newHisData := Assignment2.CountryData{
				Name:       col.Name,
				IsoCode:    col.IsoCode,
				Year:       col.Year,
				Percentage: col.Percentage,
			}
			countData = append(countData, newHisData)
		}
	}

	if len(countData) < 1 { // check if list is empty
		http.Error(w, "No country found", http.StatusNotFound)
		return
	}

	// if user want to sort by percentage
	if sortByValue == "true" {
		// Sorting the countData slice from lowest to highest by country percentage
		sort.Slice(countData, func(i, j int) bool {
			return countData[i].Percentage < countData[j].Percentage
		})
	}

	// if no Iso is given print all countries mean percentage else print one country's history
	if iso == "" {
		countMean := getAllCountriesMean(countData) // get all countries mean percentage

		jsonResponse, err := json.Marshal(countMean)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		jsonResponse, err := json.Marshal(countData)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}

/*
function that gets all countries, checks for redundancy and returns a struct of all countries with mean percentage
*/
func getAllCountriesMean(countries []Assignment2.CountryData) []Assignment2.CountryMean {
	var retData []Assignment2.CountryMean
	lastCountry := ""
	// going through all countries
	for _, current := range countries {
		// makes sure there are no repeated countries
		if current.IsoCode != lastCountry {
			lastCountry = current.IsoCode
			// variable that stores mean percentage
			var mean float64
			// calls findCountry functions to get all instances of one country
			countryHistory := findCountry(countries, current.IsoCode)
			// loops through all instances and adds their percentage to the mean variable
			for _, currentYear := range countryHistory {
				mean += currentYear.Percentage
			}
			// converts length of countryHistory slice into a float/number of instances
			numberInstances := float64(len(countryHistory))
			// calculates mean percentage of a country
			mean = mean / numberInstances
			// initiates a country struct with the mean percentage
			countryMean := Assignment2.CountryMean{
				Name:       current.Name,
				IsoCode:    current.IsoCode,
				Percentage: mean,
			}
			// appends country to slice of countries
			retData = append(retData, countryMean)
		}
	}
	return retData
}
