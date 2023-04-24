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
	"log"
	"net/http"
	"strings"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var Client *firestore.Client

// Collection name in Firestore
const collection = "webhooks"

// InitFirebase initializes the Firebase client and context.
func InitFirebase() error {
	ctx = context.Background()

	// Replace "path/to/your-service-account-key.json" with the actual path to your service account key file.
	credentials := Assignment2.FIRESTORE_CREDS
	opt := option.WithCredentialsFile(credentials)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("Failed to create a new Firebase app: %v", err)
		return err
	}

	Client, err = app.Firestore(ctx)
	if err != nil {
		log.Fatalf("Failed to create a new Firestore client: %v", err)
		return err
	}

	return nil
}

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		registerWebhook(w, r)
	case http.MethodGet:
		retrieveWebhook(w, r)
	case http.MethodDelete:
		deleteWebhook(w, r)
	}
}

// registerWebhook adds a webhook to Firestore db
func registerWebhook(w http.ResponseWriter, r *http.Request) {

	// Decode the request body into a webhook struct
	var newWebhook Assignment2.WebhookGet
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		log.Println("Error in decoding request body", err.Error())
		http.Error(w, "Error in processing request body", http.StatusBadRequest)
		return
	}

	// Initialize counter for invocation
	newWebhook.Counter = 0

	// Add element in embedded structure.
	docRef, _, err := Client.Collection(collection).Add(ctx, newWebhook)
	if err != nil {
		// Error handling
		log.Println("Error when adding Webhook to database: ", err.Error())
		http.Error(w, "Error when adding Webhook to database: ", http.StatusBadRequest)
		return
	}

	// Store the generated ID in the webhook data struct
	newWebhook.WebhookID = docRef.ID

	// Update the document with the generated ID
	_, err = Client.Collection(collection).Doc(docRef.ID).Update(ctx, []firestore.Update{{Path: "webhookId", Value: docRef.ID}})
	if err != nil {
		log.Println("Error when updating document with generated ID, Error: " + err.Error())
		http.Error(w, "Error when updating document with generated ID, Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the newly created webhook ID in the response
	resp := struct {
		WebhookID string `json:"webhookId"`
	}{
		WebhookID: docRef.ID,
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.Marshal(resp)
	if err != nil {
		// handle error
	}
	w.Write(jsonData)

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
	res := Client.Collection(collection).Doc(id)

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

// retrieveDocument retrieves a webhook specified by an id, or all webhooks if no id
// is provided from firestore db
func retrieveWebhook(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Retrieve individual webhook if id is provided
	if len(parts) > 4 {

		id := parts[4]
		// Retrieve specific webhook based on id (Firestore-generated hash)
		res := Client.Collection(collection).Doc(id)

		// Retrieve reference to document
		doc, err := res.Get(ctx)
		if err != nil {
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
		iter := Client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"

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
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(hooks)
	}
}

func UpdateAndInvoke(isoCode string) {

	// Get all webhooks from Firestore
	iter := Client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"

	var hooks []Assignment2.WebhookGet

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Failed to iterate: %v", err)
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
				Counter:   m["Counter"].(int64),
			}
			hooks = append(hooks, newHook)
		}
	}

	// Loop through all webhooks
	for _, currentHook := range hooks {
		// If current webhook == isoCode
		if currentHook.Country == strings.ToUpper(isoCode) {
			currentHook.Counter++
			// If conditions for invocation are met
			if (currentHook.Counter%currentHook.Calls == 0) && currentHook.Counter > 0 {
				invokeWebhook(currentHook)
			}
		}
		// Update webhook in Firestore
		docRef := Client.Collection(collection).Doc(currentHook.WebhookID)

		// Set the "counter" field of the webhook to the new value
		_, err := docRef.Set(ctx, currentHook)

		// If error, log and return
		if err != nil {
			log.Printf("Failed updating document: %v", err)
		}
	}
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
		log.Println("Error during request creation. Error:", err)
		return
	}

}
