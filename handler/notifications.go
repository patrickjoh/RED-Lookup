package handler

import (
	"Assignment2"
	"cloud.google.com/go/firestore"
	"context"
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

// GetContextAndClient returns the Firebase context and client to functions which require their use
func GetContextAndClient() (context.Context, *firestore.Client) {
	return ctx, client
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

	log.Println("Received request to add document for content ", string(text))
	if len(string(text)) == 0 {
		log.Println("Content appears to be empty.")
		http.Error(w, "Your payload (to be stored as document) appears to be empty. Ensure to terminate URI with /.", http.StatusBadRequest)
	} else {
		// Add element in embedded structure.
		id, _, err := client.Collection("webhooks").Add(ctx,
			map[string]interface{}{
				"text": string(text),
				"ct":   ct,
				"time": firestore.ServerTimestamp,
			})
		ct++
		if err != nil {
			// Error handling
			log.Println("Error when adding document " + string(text) + ", Error: " + err.Error())
			http.Error(w, "Error when adding document "+string(text)+", Error: "+err.Error(), http.StatusBadRequest)
			return
		} else {
			// Returns document ID in body
			log.Println("Document added to collection. Identifier of returned document: " + id.ID)
			http.Error(w, id.ID, http.StatusCreated)

			return
		}
	}
}

func deleteDocument(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")
	log.Println("UWU1 ")

	// Error handling
	if len(parts) < 4 || parts[3] != "notifications" {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	id := parts[4]

	if len(id) != 0 {
		// Retrieve specific message based on id (Firestore-generated hash)
		res := client.Collection("webhooks").Doc(id)

		// Retrieve reference to document
		_, err2 := res.Get(ctx)
		if err2 != nil {
			log.Println("Error extracting body of returned document of message " + id)
			http.Error(w, "Error extracting body of returned document of message "+id, http.StatusInternalServerError)
			return
		}

		// Delete webhook from Firestore
		_, err3 := res.Delete(ctx)
		if err3 != nil {
			log.Println("Error deleting document " + id)
			http.Error(w, "Error deleting document "+id, http.StatusInternalServerError)
			return
		}

	} else {
		http.Error(w, "Id must be specified. Try /energy/v1/notifications/{id}", http.StatusBadRequest)
	}
}

func postPayload(payload interface{}) {

}

// retrieveDocument a webhook specified by an id, or all webhooks if no id is provided
func retrieveDocument(w http.ResponseWriter, r *http.Request) {
	// Get id from RUL
	id := r.URL.Query().Get("id")

	// Extract individual message
	if len(id) != 0 {

		// Retrieve specific message based on id (Firestore-generated hash)
		res := client.Collection("webhooks").Doc(id)

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
		// Collective retrieval of messages
		iter := client.Collection("webhooks").Documents(ctx) // Loop through all entries in collection "messages"

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
