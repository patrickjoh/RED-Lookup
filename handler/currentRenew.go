package handler

import (
	"Assignment2"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Documentation...
func HandlerRenewables(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandlerRenewablesGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return

}

// HandlerRenewablesGet
func HandlerRenewablesGet(w http.ResponseWriter, r *http.Request) {
	// open CSV file
	fd, error := os.Open("data.csv")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	// read CSV file
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}

	// Split url to get keyword
	keywords := strings.Split(r.URL.Path, "/")
	fmt.Println("\nlen(keywords): ", len(keywords))
	fmt.Println("\nkeywords[4]: ", keywords[4])
	fmt.Println("\nkeywords[5]: ", keywords[5])
	fmt.Println("\nkeywords[6]: ", keywords[6])

	// Error handling
	if len(keywords) < 5 || keywords[4] != "current" {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	var countryData []Assignment2.CountData

	var neighbors = false

	// If neighbor bool parameter is provided
	if len(keywords) > 5 {
		neighbors = true
	}

	// If country code is provided
	if len(keywords) == 5 {
		getOneCountry(keywords[5], neighbors)
	} else { // If no country code is provided
		getAllCountries(countryData, records)
	}

	jsonResponse, err := json.Marshal(countryData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	w.WriteHeader(http.StatusOK)
}

func getAllCountries(data []Assignment2.CountData, records [][]string) []Assignment2.CountData {
	for i, record := range records {
		data[i].Name = record[0]
		data[i].IsoCode = record[1]
		data[i].Year = record[2]
		data[i].Percentage = record[3]
	}
	return data
}

func getOneCountry(keyword string, neighbor bool) {

}

func getNeighborCountry() {

}
