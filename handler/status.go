package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

/*
HandlerStaus: Entry point handler for Status handler
*/
func HandlerStaus(r http.ResponseWriter, w *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandelStatus(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// HandelStatusGet
func HandelStatus(w http.ResponseWriter, r *http.Request) {
	// Define the URLs
	url1 := "http://universities.hipolabs.com/"
	url2 := "https://restcountries.com/"

	// Instantiate client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// ...

	status := statusData{
		countries_api:   statusMap[url1],
		notification_db: statusMap[url2],
		version:         "v1",
		uptime:          time.Since(startTime).String(),
	}

	// Encode struct as JSON
	data, err := json.Marshal(diagnostics)
	if err != nil {
		fmt.Errorf("Error encoding JSON: %s", err.Error())
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(http.StatusOK)

}
