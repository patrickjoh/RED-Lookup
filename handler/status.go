package handler

import (
	"Assignment2"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

// StatusHandler is the entry point handler for Status handler
func StatusHandler(w http.ResponseWriter, r *http.Request) {
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
	fd, err := os.Open(Assignment2.CSV_PATH)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error opening CSV file: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer func(fd *os.File) {
		err := fd.Close()
		if err != nil {
			fmt.Println("Error closing file: ", err)
		}
	}(fd)

	restResp, err := http.Get(restURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error in response: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	defer restResp.Body.Close()

	// Attempt to find a collection
	var fireStoreAvail = FirestoreStatus()

	var numOfHooks = 0
	status := strconv.Itoa(http.StatusOK) + " " + http.StatusText(http.StatusOK)
	// Only attempt to retrieve number of webhooks if a collection is found
	if fireStoreAvail == status {
		numOfHooks = GetNumWebhooks()
	}

	elapsed := time.Since(startTime)
	elapsedSeconds := int(elapsed.Seconds())

	// Get status codes from response structs
	stData := Assignment2.StatusData{
		CountriesAPI:   restResp.Status,
		NotificationSB: fireStoreAvail,
		Webhooks:       numOfHooks,
		Version:        "v1",
		Uptime:         elapsedSeconds,
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

// FirestoreStatus checks availability of Firestore db and returns a status code
func FirestoreStatus() string {

	initFirebase()

	// Check if client is nil
	if client == nil {
		log.Println("Client is nil")
		status := strconv.Itoa(http.StatusInternalServerError) + " " + http.StatusText(http.StatusInternalServerError)
		return status
	}

	// Attempt to retrieve all collection references from Firestore
	collections, err := client.Collections(ctx).GetAll()

	// Return error if collection cannot be found
	if err != nil || collections == nil || len(collections) < 1 {

		status := strconv.Itoa(http.StatusServiceUnavailable) + " " + http.StatusText(http.StatusServiceUnavailable)
		return status
	}

	// Attempt to get document from any collection in Firestore
	for _, collectionRef := range collections {
		// Check if iter is nil
		iter := collectionRef.Limit(1).Documents(ctx)
		if iter == nil {
			log.Println("Iter is nil")
			continue
		}

		_, err := iter.Next()
		if err == nil {
			// Return a status code indicating that Firestore service is available
			status := strconv.Itoa(http.StatusOK) + " " + http.StatusText(http.StatusOK)
			return status
		} else if err == iterator.Done {
			continue
		} else {
			log.Printf("Error while iterating through documents: %v", err)
			continue
		}
	}

	// Return error status code if no document found
	status := strconv.Itoa(http.StatusServiceUnavailable) + " " + http.StatusText(http.StatusServiceUnavailable)
	return status
}

// GetNumWebhooks retrieves and returns the number of registered webhooks from Firestore
func GetNumWebhooks() int {

	// Create reference to webhook collection in Firestore
	webhooksCollection := client.Collection(collection)

	// Retrieve all webhooks from db
	iter := webhooksCollection.Documents(ctx)
	var numWebhooks int
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			panic(err)
		}
		numWebhooks++
	}

	return numWebhooks
}
