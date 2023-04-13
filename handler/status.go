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
		HandelStatus(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// HandelStatusGet: shows availability for all third party services
func HandelStatus(w http.ResponseWriter, r *http.Request) {
	// Define the URLs
	//url1 := ... // TODO
	//url2 := ... // TODO

	// Instantiate client
	client := &http.Client{}
	defer client.CloseIdleConnections()

	// TODO: implement get request

	stData := Assignment2.StatusData{
		//Countries_api:   ..., // TODO
		//Notification_db: ..., // TODO
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
