package handler

import (
	"Assignment2"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
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

// Secret
var Secret []byte

// SignatureKey initializes signature (via init())
var SignatureKey = "X-SIGNATURE"

// InitFirebase initializes the Firebase client and context.
// taken from code example 13
func InitFirebase() {
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

	InitFirebase()

	defer func() {
		err := client.Close()
		if err != nil {
			log.Fatal("Closing of the Firebase client failed. Error:", err)
		}
	}()

	switch r.Method {
	case http.MethodPost:
		registerWebhook(w, r)
	case http.MethodGet:
		retrieveWebhook(w, r)
	case http.MethodDelete:
		deleteWebhook(w, r)
	}
}

// addDocument adds a webhook to Firestore db
func registerWebhook(w http.ResponseWriter, r *http.Request) {

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

	var newWebhook Assignment2.WebhookGet
	err = json.Unmarshal(text, &newWebhook)
	if err != nil {
		log.Println("Error in decoding request body")
		http.Error(w, "Error in decoding body.", http.StatusBadRequest)
		return
	}

	// Initialize counter for invocation
	newWebhook.Counter = 0

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

// deleteDocument deletes a webhook from Firestore db
func deleteWebhook(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Error handling if no provided id
	if len(parts) < 5 || len(parts[4]) == 0 {
		http.Error(w, "Id must be specified. Try /energy/v1/notifications/{id}", http.StatusBadRequest)
		log.Println("Malformed URL", http.StatusBadRequest)
		return
	}

	id := parts[4]

	// Retrieve specific message based on id (Firestore-generated hash)
	res := client.Collection(collection).Doc(id)

	// Attempt to retrieve reference to document
	doc, err := res.Get(ctx)
	if err != nil {
		log.Println("Error extracting body of returned document of message " + id)
		http.Error(w, "Error extracting body of returned document of message "+id, http.StatusInternalServerError)
		return
	}

	var data Assignment2.WebhookGet

	// Get webhook to be deleted
	err = doc.DataTo(&data)
	if err != nil {
		log.Println("Error in decoding request body")
		http.Error(w, "Error in decoding body.", http.StatusBadRequest)
		return
	}

	data.WebhookID = id

	// Attempt to delete webhook from Firestore
	_, err = res.Delete(ctx)
	if err != nil {
		log.Println("Error deleting document " + id)
		http.Error(w, "Error deleting document "+id, http.StatusInternalServerError)
		return
	}

	log.Println("Document '" + id + "' deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func webhookInvocation(w http.ResponseWriter, r *http.Request) {

}

// retrieveDocument retrieves a webhook specified by an id, or all webhooks if no id
// is provided from firestore db
func retrieveWebhook(w http.ResponseWriter, r *http.Request) {
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

		m["WebhookID"] = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(m)
	} else {
		// Retrieve all webhooks if no id is provided
		iter := client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"

		var hooks []Assignment2.WebhookGet

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
			m["WebhookID"] = doc.Ref.ID

			if m["Calls"] != nil {

				new := Assignment2.WebhookGet{
					WebhookID: m["WebhookID"].(string),
					Url:       m["Url"].(string),
					Country:   m["Country"].(string),
					Calls:     m["Calls"].(int64),
				}

				hooks = append(hooks, new)
			}

		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hooks)
	}
}

func updateAndInvoke(isoCode string) {

	// Get all webhooks from Firestore
	iter := client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"

	var hooks []Assignment2.WebhookGet

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
		m["WebhookID"] = doc.Ref.ID

		if m["Calls"] != nil {

			newHook := Assignment2.WebhookGet{
				WebhookID: m["WebhookID"].(string),
				Url:       m["Url"].(string),
				Country:   m["Country"].(string),
				Calls:     m["Calls"].(int64),
			}
			hooks = append(hooks, newHook)
		}
	}

	// Loop through all webhooks
	for _, currentHook := range hooks {
		// If current webhook == isoCode
		if currentHook.Country == isoCode {
			currentHook.Counter++
			// If conditions for invocation are met
			if (currentHook.Calls%currentHook.Counter == 0) && currentHook.Counter > 0 {
				invokeWebhook(currentHook)
			}
		}
		updateWebhook(currentHook.WebhookID, currentHook.Counter)
	}
}

func updateWebhook(id string, counter int64) {

}

func invokeWebhook(invoke Assignment2.WebhookGet) {

	id := invoke.WebhookID

	data := Assignment2.WebhookInvoke{
		WebhookID: id,
		Country:   invoke.Country,
		Calls:     invoke.Calls,
	}

	payload, _ := json.Marshal(data)
	log.Println("Attempting invocation of url " + invoke.Url + " with content '" + "payload" + "'.")
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	_, err := http.Post(invoke.Url, "application/json", bytes.NewReader([]byte(payload)))
	if err != nil {
		log.Println("%v", "Error during request creation. Error:", err)
		return
	}

}
