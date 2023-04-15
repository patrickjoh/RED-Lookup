package handler

import (
	"Assignment2"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

/*
HandlerStatus: Entry point handler for Status handler
*/
func HandlerStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleStatus(w)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// handleStatus: shows availability for all third party services
func handleStatus(w http.ResponseWriter) {
	// Define the URLs
	restURL := "http://129.241.150.113:8080/"

	// open CSV file
	fd, err := os.Open("data/renewable-share-energy-csv.csv")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening CSV file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer fd.Close()

	// read CSV file
	fileReader := csv.NewReader(fd)
	records, err := fileReader.ReadAll()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading CSV file: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	fmt.Println(records)

	restResp, err := http.Get(restURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer restResp.Body.Close()

	// Get status codes from response structs
	stData := Assignment2.StatusData{
		CountriesAPI: restResp.Status,
		//Webhooks: , // TODO
		Version: "v1",
		Uptime:  time.Since(startTime).String(),
	}

	// Encode struct as JSON
	data, err := json.Marshal(stData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
