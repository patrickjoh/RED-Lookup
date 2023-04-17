package handler

import (
	"Assignment2"
	"cloud.google.com/go/firestore"
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

	var fireStoreAvail = firestoreStatus()

	var numOfHooks = 0
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

// firestoreStatus checks availability of Firestore db and returns status code
func firestoreStatus() int {
	ctx, client := GetContextAndClient()

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic: %v", r)
		}
	}()

	collections, err := func() ([]*firestore.CollectionRef, error) {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic in GetAll: %v", r)
			}
		}()
		return client.Collections(ctx).GetAll()
	}()

	if err != nil || collections == nil || len(collections) < 1 {
		return http.StatusServiceUnavailable
	}

	// Attempt to get document from any collection in Firestore
	for _, collectionRef := range collections {
		iter := collectionRef.Limit(1).Documents(ctx)
		_, err := iter.Next()
		if err == nil {
			// Return a status code indicating that Firestore service is available
			return http.StatusOK
		} else {
			log.Println("Error: %v", err)
		}
	}

	// Return error status code if no document found
	return http.StatusServiceUnavailable
}

// GetNumWebhooks retrieves and returns the number of registered webhooks from Firestore
func GetNumWebhooks() int {
	// Get context and client
	ctx, client = GetContextAndClient()

	// Create reference to webhook collection in Firestore
	webhooksCollection := client.Collection("webhooks")

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
