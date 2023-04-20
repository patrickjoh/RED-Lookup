package handler

import (
	"Assignment2"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"strings"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

// Collection name in Firestore
const collection = "webhooks"

// Message counter to produce some variation in content
var ct = 0

func initFirebase() {
	ctx = context.Background()

	sa := option.WithCredentialsFile(Assignment2.FIRESTORE_CREDS)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)

	// Check whether there is an error when connecting to Firestore
	if err != nil {
		log.Fatalln(err)
	}
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	initFirebase()

	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", err)
		}
	}()

	switch r.Method {
	case http.MethodPost:
		addDocument(w, r)
	case http.MethodGet:
		retrieveDocument(w, r)
	case http.MethodDelete:
		deleteDocument(w, r)
	}
}

func addDocument(w http.ResponseWriter, r *http.Request) {
	// Read body
	text, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Reading payload from body failed.")
		http.Error(w, "Reading payload failed.", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	log.Println("Received request to add document for content ", string(text))
	if len(string(text)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(w, "Your payload appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
		return
	}

	var newWebhook Assignment2.WebhookData
	err = json.Unmarshal(text, &newWebhook)
	if err != nil {
		log.Println("Error in decoding request body")
		http.Error(w, "Error in decoding body.", http.StatusBadRequest)
		return
	}

	// Add element in embedded structure.
	docRef, _, err := client.Collection(collection).Add(ctx, newWebhook)
	if err != nil {
		// Error handling
		log.Println("Error when adding document " + string(text) + ", Error: " + err.Error())
		http.Error(w, "Error when adding document "+string(text)+", Error: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Store the generated ID in the webhook data struct
	newWebhook.WebhookID = docRef.ID

	// Return the newly created webhook ID in the response
	resp := struct {
		WebhookID string `json:"webhookId"`
	}{
		WebhookID: docRef.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func deleteDocument(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Error handling if no provided id
	if len(parts) < 5 || len(parts[4]) == 0 {
		http.Error(w, "Id must be specified. Try /energy/v1/notifications/{id}", http.StatusBadRequest)
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	id := parts[4]

	// Retrieve specific message based on id (Firestore-generated hash)
	res := client.Collection(collection).Doc(id)

	// Attempt to retrieve reference to document
	_, err2 := res.Get(ctx)
	if err2 != nil {
		log.Println("Error extracting body of returned document of message " + id)
		http.Error(w, "Error extracting body of returned document of message "+id, http.StatusInternalServerError)
		return
	}

	// Attempt to delete webhook from Firestore
	_, err3 := res.Delete(ctx)
	if err3 != nil {
		log.Println("Error deleting document " + id)
		http.Error(w, "Error deleting document "+id, http.StatusInternalServerError)
		return
	} else {
		log.Println("Document '" + id + "' deleted successfully")
		http.Error(w, "Document '"+id+"' deleted successfully", http.StatusOK)
	}
}

func postPayload(payload interface{}) {

}

// retrieveDocument a webhook specified by an id, or all webhooks if no id is provided
func retrieveDocument(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Retrieve individual webhook if id is provided
	if len(parts) > 4 && len(parts[4]) != 0 {

		id := parts[4]

		// Retrieve specific webhook based on id (Firestore-generated hash)
		res := client.Collection(collection).Doc(id)

		// Retrieve reference to document
		doc, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document of message " + id)
			http.Error(w, "Error extracting body of returned document of message "+id, http.StatusInternalServerError)
			return
		}

		// A message map with string keys. Each key is one field, like "text" or "timestamp"
		m := doc.Data()
		_, err3 := fmt.Fprintln(w, m["text"])
		if err3 != nil {
			log.Println("Error while writing response body of message " + id)
			http.Error(w, "Error while writing response body of message "+id, http.StatusInternalServerError)
			return
		}
	} else {
		// Retrieve all webhooks if no id is provided
		iter := client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"

		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}

			// A message map with string keys. Each key is one field, like "text" or "timestamp"
			m := doc.Data()
			_, err = fmt.Fprintln(w, m)
			if err != nil {
				log.Println("Error while writing response body (Error: " + err.Error() + ")")
				http.Error(w, "Error while writing response body (Error: "+err.Error()+")", http.StatusInternalServerError)
			}
		}
	}
}
