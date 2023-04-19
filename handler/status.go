package handler

import (
	"Assignment2"
	"encoding/json"
	"fmt"
	"google.golang.org/api/iterator"
	"log"
	"net/http"
	"os"
	"time"
)

var startTime time.Time

func init() {
	startTime = time.Now()
}

/*
StatusHandler: Entry point handler for Status handler
*/
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
	fd, err := os.Open("handler/data/renewable-share-energy.csv")
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
	var fireStoreAvail = firestoreStatus()

	var numOfHooks = 0
	// Only attempt to retrieve number of webhooks if a collection is found
	if fireStoreAvail == 200 {
		numOfHooks = GetNumWebhooks()
	}

	// Get status codes from response structs
	stData := Assignment2.StatusData{
		CountriesAPI:   restResp.Status,
		NotificationSB: fireStoreAvail,
		Webhooks:       numOfHooks,
		Version:        "v1",
		Uptime:         time.Since(startTime).String(),
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

// firestoreStatus checks availability of Firestore db and returns a status code
func firestoreStatus() int {

	// Check if client is nil
	if Client == nil {
		log.Println("Client is nil")
		return http.StatusInternalServerError
	}

	// Attempt to retrieve all collection references from Firestore
	collections, err := Client.Collections(Ctx).GetAll()

	// Return error if collection cannot be found
	if err != nil || collections == nil || len(collections) < 1 {
		return http.StatusServiceUnavailable
	}

	// Attempt to get document from any collection in Firestore
	for _, collectionRef := range collections {
		// Check if iter is nil
		iter := collectionRef.Limit(1).Documents(Ctx)
		if iter == nil {
			log.Println("Iter is nil")
			continue
		}

		_, err := iter.Next()
		if err == nil {
			// Return a status code indicating that Firestore service is available
			return http.StatusOK
		} else if err == iterator.Done {
			continue
		} else {
			log.Printf("Error while iterating through documents: %v", err)
			continue
		}
	}

	// Return error status code if no document found
	return http.StatusServiceUnavailable
}

// GetNumWebhooks retrieves and returns the number of registered webhooks from Firestore
func GetNumWebhooks() int {

	// Create reference to webhook collection in Firestore
	webhooksCollection := Client.Collection("webhooks")

	// Retrieve all webhooks from db
	iter := webhooksCollection.Documents(Ctx)
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
