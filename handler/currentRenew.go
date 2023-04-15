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

	/*var neighbors = false

	// If neighbor bool parameter is provided
	if len(keywords) > 5 {
		neighbors = true
	}*/

	var countryData []Assignment2.CountryData

	// If country code is provided
	if len(keywords) >= 7 {
		// countryData = getOneCountry(keywords[5], neighbors)
		log.Println("getOneCountry() not implemented") // TODO
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
	for i, current := range data {
		if current.Name == currCount { // Still same country?
			if current.Year > currHighYear { // New highest year found
				highestRecord = current
				fmt.Println(i) // UWU remove
			}
		} else { // New country entered
			currCount = current.Name
			currHighYear = current.Year
			retData = append(retData, highestRecord)
		}
	}

	for i := 0; i < len(retData); i++ {
		log.Println(retData[i].Name)
	}

	return retData
}

func getOneCountry(keyword string, neighbor bool) {

}

func getNeighborCountry() {

}
