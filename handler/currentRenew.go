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
func HandlerRenewables(w http.ResponseWriter, r *http.Request) {
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

	var neighbors = false

	// If neighbor bool parameter is provided
	if len(keywords) > 7 {
		if len(keywords[8]) == 3 {
			fmt.Println("Length of iso code mus be 3")
		} else {
			neighbors = true
		}
	}

	var countryData []Assignment2.CountryData

	// If country code is provided
	if len(keywords) >= 6 {
		// UWU transform to upper case
		countryData = getOneCountry(convertCsvData(), keywords[5], neighbors)
	} else { // If no country code is provided
		countryData = getAllCountries(convertCsvData())
	}

	jsonResponse, err := json.Marshal(countryData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	w.WriteHeader(http.StatusOK)
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

func getOneCountry(data []Assignment2.CountryData, keyword string, neighbor bool) []Assignment2.CountryData {

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

func getNeighborCountry() {

}
