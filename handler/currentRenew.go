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
	case http.MethodPost:
		HandlerRenewablesPost(w, r)
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

	// Error handling
	if len(keywords) != 5 || keywords[3] != "uniinfo" {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	// if {country} -> do something

	// if {?neighbours=bool?} -> do something more

	var countryData []Assignment2.CountData

	for i, record := range records {
		countryData[i].Name = record[0]
		countryData[i].IsoCode = record[1]
		countryData[i].Year = record[2]
		countryData[i].Percentage = record[3]
	}

	jsonResponse, err := json.Marshal(countryData)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	w.WriteHeader(http.StatusOK)
}

// HandlerRenewablesPost
func HandlerRenewablesPost(w http.ResponseWriter, r *http.Request) {

}
