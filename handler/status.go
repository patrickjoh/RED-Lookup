package handler

import (
	"Assignment2"
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
func HandlerStaus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handelStatus(w)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// HandelStatusGet: shows availability for all third party services
func handelStatus(w http.ResponseWriter) {
	// Define the URLs
	restUrl := "http://129.241.150.113:8080/"
	dataUrl := "https://ourworldindata.org/energy"

	// Make requests using http.Get
	dataResp, err := http.Get(dataUrl)
	if err != nil {
		fmt.Errorf("Error in response: %s", err.Error())
	}

	// Close response bodies when done
	defer dataResp.Body.Close()

	restResp, err := http.Get(restUrl)
	if err != nil {
		fmt.Errorf("Error in response: %s", err.Error())
	}

	defer restResp.Body.Close()

	// Get status codes from response structs
	stData := Assignment2.StatusData{
		Countries_api:   restResp.Status,
		Notification_db: dataResp.Status,
		//Webhooks: , // TODO
		Version: "v1",
		Uptime:  time.Since(startTime).String(),
	}

	// Encode struct as JSON
	data, err := json.Marshal(stData)
	if err != nil {
		fmt.Errorf("Error encoding JSON: %s", err.Error())
		return
	}

	// Set content type
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
