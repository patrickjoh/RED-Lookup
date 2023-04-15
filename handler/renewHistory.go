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
func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandelHistoryPost(w, r)
	case http.MethodGet:
		HandelHistoryGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// HandelHistoryPost
func HandelHistoryPost(w http.ResponseWriter, r *http.Request) {

}

// HandelHistoryGet
func HandelHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Split url to get keyword
	urlKeywords := strings.Split(r.URL.Path, "/")

	// Error handling
	if len(urlKeywords) < 5 {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	count := urlKeywords[5] // Get country name from url
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

	// Send a response with the extracted values
	fmt.Fprintf(w, "Searching for: country = %s, year (%s - %s)", count, begin, end)

	contries := convertCsvData() // Acquire data from csv-file

	var countData []Assignment2.CountryData // empty list for the final data
	startYear, _ := strconv.Atoi(begin)     // beginning year
	endYear, _ := strconv.Atoi(end)         // end year

	fmt.Println(contries[0].Name)
	// loops through the csv-file and return the data that was asked for
	for _, col := range contries {
		if col.Name == count && col.Year < endYear && col.Year > startYear {
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
