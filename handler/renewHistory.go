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
		log.Println(w, "Incorrect use or year", http.StatusBadRequest)
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
	// Find country's history from beginning to end
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
