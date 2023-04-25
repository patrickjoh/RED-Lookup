package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	opt := option.WithCredentialsJSON(Assignment2.FirebaseCredentials)
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
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodPost+
			"', '"+http.MethodGet+"' and '"+http.MethodDelete+"' are supported.", http.StatusNotImplemented)
	}
	return
}

// registerWebhook adds a webhook to Firestore db
func registerWebhook(w http.ResponseWriter, r *http.Request) {

	// Decode the request body into a webhook struct
	var newWebhook structs.WebhookGet
	err := json.NewDecoder(r.Body).Decode(&newWebhook)
	if err != nil {
		log.Println("Error in decoding request body", err.Error())
		http.Error(w, "Error in processing request body", http.StatusBadRequest)
		return
	}

	// Check if the POST request contains a valid json body
	if newWebhook.Url == "" {
		log.Println("Invalid URL: empty")
		http.Error(w, "Invalid URL: empty", http.StatusBadRequest)
		return
	}

	matchingCountry := findCountry(Assignment2.CSVData, newWebhook.Country)
	validCountry := false
	// Check if the country is valid
	if len(matchingCountry) > 0 {
		validCountry = true
	}
	// Returns an error if the country is invalid
	if validCountry == false {
		log.Println("Invalid Country: " + newWebhook.Country)
		http.Error(w, "Invalid Country: "+newWebhook.Country, http.StatusBadRequest)
		return
	}

	// Check if the calls is valid
	if newWebhook.Calls <= 0 {
		log.Println("Invalid Calls: must be greater than 0")
		http.Error(w, "Invalid Calls: must be greater than 0", http.StatusBadRequest)
		return
	}

	// Initialize counter for invocation
	newWebhook.Counter = 0

	// Generate a new document reference
	docRef := Client.Collection(collection).NewDoc()

	// Set the generated ID in the webhook data struct
	newWebhook.WebhookID = docRef.ID

	// Add the webhook to the database with the generated ID
	_, err = docRef.Set(ctx, newWebhook)
	if err != nil {
		// Error handling
		log.Println("Error when adding Webhook to database: ", err.Error())
		http.Error(w, "Error when adding Webhook to database: ", http.StatusBadRequest)
		return
	}

	// Create a response body with the newly created webhook ID
	// Return the newly created webhook ID in the response
	resp := struct {
		// WebhookID is the ID of the newly created webhook
		WebhookID string `json:"webhookId"`
	}{
		WebhookID: docRef.ID,
	}

	// Encode the response body and send it to the client
	jsonResponse, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error in encoding response body", err.Error())
		http.Error(w, "Error in encoding response body", http.StatusBadRequest)
		return
	}
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 201 (Created)
	w.WriteHeader(http.StatusCreated)
	// Write the response body
	w.Write(jsonResponse)
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

	// Extract the id from the URL
	id := parts[4]

	// Retrieve specific message based on id (Firestore-generated hash)
	docRef := Client.Collection(collection).Doc(id)

	// Attempt to retrieve reference to document
	doc, err := docRef.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Println("Document not found with ID: " + id)
			http.Error(w, "Document not found with ID: "+id, http.StatusNotFound)
			return
		} else {
			log.Println("Error extracting body of returned document of message " + id)
			http.Error(w, "Error extracting body of returned document of message "+id, http.StatusInternalServerError)
			return
		}
	}

	// Create a buffer to store the document data
	var data structs.WebhookGet

	// Get webhook to be deleted
	err = doc.DataTo(&data)
	if err != nil {
		log.Println("Error in decoding request body")
		http.Error(w, "Error in decoding body.", http.StatusBadRequest)
		return
	}

	// Attempt to delete webhook from Firestore
	_, err = docRef.Delete(ctx)
	if err != nil {
		log.Println("Error deleting document " + id)
		http.Error(w, "Error deleting document "+id, http.StatusInternalServerError)
		return
	}

	// Marshal the data and write it to the response
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling document data: ", err.Error())
		http.Error(w, "Error marshaling document data", http.StatusInternalServerError)
		return
	}
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Set the status code to 200 (OK)
	w.WriteHeader(http.StatusOK)
	// Write the response body
	w.Write(jsonData)
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

		var hooks []structs.WebhookGet

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

				newHook := structs.WebhookGet{
					WebhookID: m["WebhookID"].(string),
					Url:       m["Url"].(string),
					Country:   m["Country"].(string),
					Calls:     m["Calls"].(int64),
				}

				hooks = append(hooks, newHook)
			}

		}

		// Encode the response body and send it to the client
		marshallResponse, err := json.Marshal(hooks)
		if err != nil {
			log.Println("Error in encoding response body", err.Error())
			http.Error(w, "Error in encoding response body", http.StatusBadRequest)
			return
		}
		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")
		// Set the status code to 201 (Created)
		w.WriteHeader(http.StatusOK)
		// Write the response body
		w.Write(marshallResponse)
	}
}

func UpdateAndInvoke(isoCode string) {

	// Get all webhooks from Firestore
	iter := Client.Collection(collection).Documents(ctx) // Loop through all entries in collection "messages"

	var hooks []structs.WebhookGet

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

			newHook := structs.WebhookGet{
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

func invokeWebhook(invoke structs.WebhookGet) {

	id := invoke.WebhookID

	data := structs.WebhookInvoke{
		WebhookID: id,
		Country:   invoke.Country,
		Calls:     invoke.Calls,
	}

	payload, _ := json.Marshal(data)
	log.Println("Attempting invocation of url " + invoke.Url + " with content '" + "payload" + "'.")
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	_, err := http.Post(invoke.Url, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Println("Error during request creation. Error:", err)
		return
	}

}
