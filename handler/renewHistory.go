package handler

import (
	"Assignment2"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Documentation...
func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleHistoryPost(w, r)
	case http.MethodGet:
		handleHistoryGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// handleHistoryPost
func handleHistoryPost(w http.ResponseWriter, r *http.Request) {

}

// handleHistoryGet
func handleHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Split url to get keyword
	urlKeywords := strings.Split(r.URL.Path, "/")

	iso := urlKeywords[5]   // Get country isoCode from url
	query := r.URL.RawQuery // Get the queries from url
	fmt.Println(len(iso))
	if len(iso) != 3 && len(iso) != 0 {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
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

	if begin > end {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} else if begin == "" && end == "" {
		begin = "0"
		end = "3000"
	} else if begin == "" {
		begin = end
	} else if end == "" {
		end = begin
	}

	// Send a response with the extracted values
	fmt.Fprintf(w, "Searching for: country = %s, year (%s - %s)\n", iso, begin, end)

	if iso == "" {
		var countMean []Assignment2.CountryMean // empty list for the final data
		countMean = getAllCountriesMean()

		fmt.Println(countMean)

		jsonResponse, err := json.Marshal(countMean)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	} else {
		var countData []Assignment2.CountryData       // empty list for the final data
		startYear, _ := strconv.Atoi(begin)           // beginning year
		endYear, _ := strconv.Atoi(end)               // end year
		country := findCountry(convertCsvData(), iso) // Acquire data from csv-file
		// loops through the csv-file and return the data that was asked for
		for _, col := range country {
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

		fmt.Println(countData)

		jsonResponse, err := json.Marshal(countData)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)
	}
}
