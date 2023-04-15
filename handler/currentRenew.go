package handler

import (
	"Assignment2"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Documentation...
func HandlerRenewables(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handelRenewablesGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return

}

// handleRenewablesGet
func handelRenewablesGet(w http.ResponseWriter, r *http.Request) {

	// Open CSV file
	fd, error := os.Open("handler/data/renewable-share-energy.csv")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file") // UWU remove when done
	defer fd.Close()

	// Read CSV file
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}

	// -----------------------------------------------------------------------

	// Split url to get keyword
	keywords := strings.Split(r.URL.Path, "/")
	// keywords[4]  'current'

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

	var countryData []Assignment2.CountData

	// If country code is provided
	if len(keywords) == 5 {
		// countryData = getOneCountry(keywords[5], neighbors)
		log.Println("getOneCountry() not implemented")
	} else { // If no country code is provided
		countryData = getAllCountries(records)
	}

	jsonResponse, err := json.Marshal(countryData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	w.WriteHeader(http.StatusOK)
}

func getAllCountries(records [][]string) []Assignment2.CountData {
	var data []Assignment2.CountData

	currCount := records[0][0] // Current country
	currHighYear := 0          // Current highest year
	var highestRecord []string
	var entryIndices [][]string

	// Finding entries with the most recent year
	for i, record := range records {
		// UWU check if entry has iso code
		if record[0] == currCount { // Still same country?
			year, _ := strconv.Atoi(record[1])
			if year > currHighYear { // New highest year found
				highestRecord = record
				fmt.Println(i)
			}
		} else { // New country entered
			currCount = record[0]
			currHighYear, _ = strconv.Atoi(record[2])
			entryIndices = append(entryIndices, highestRecord)
		}
	}

	/*for i, entry := range entryIndecies {
		for record := range records {
			data[i].Name = entry.record[0]
			data[i].IsoCode = entry.record[1]
			data[i].Year = entry.record[2]
			data[i].Percentage = entry.record[3]
		}
	}*/
	return data
}

func getOneCountry(keyword string, neighbor bool) {

}

func getNeighborCountry() {

}
