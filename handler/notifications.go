package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"bytes"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
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

type WebhookCache struct {
	sync.RWMutex
	cache    map[string]structs.WebhookGet
	lastSync time.Time
}

var webhookCache WebhookCache

func InitCache() {

	webhookCache = WebhookCache{
		cache:    make(map[string]structs.WebhookGet),
		lastSync: time.Now(),
	}

	// Retrieve all documents from Firestore and add them to the cache
	docs, err := Client.Collection(collection).Documents(ctx).GetAll()
	if err != nil {
		log.Fatalf("Failed to get documents from Firestore: %v", err)
	}

	for _, doc := range docs {
		// A message map with string keys. Each key is one field, like "text" or "timestamp"
		m := doc.Data()
		m["WebhookID"] = doc.Ref.ID

		if m["WebhookID"] != nil {
			newHook := structs.WebhookGet{
				WebhookID: m["WebhookID"].(string),
				Url:       m["Url"].(string),
				Country:   m["Country"].(string),
				Calls:     m["Calls"].(int64),
				Counter:   m["Counter"].(int64),
				Modified:  false,
			}

			webhookCache.cache[newHook.WebhookID] = newHook
		}
	}

	log.Println("Cache initialized")
}

func SyncCacheToFirebase() {
	webhookCache.Lock()
	defer webhookCache.Unlock()

	updatedWebhooksCount := 0
	log.Println("Starting SyncCacheToFirebase")

	batch := Client.BulkWriter(ctx) // create a new batch write

	tmpCache := make(map[string]structs.WebhookGet) // creates a temp cache for webhooks that are modified
	for _, webhook := range webhookCache.cache {
		if !webhook.Modified {
			continue
		}
		tmpCache[webhook.WebhookID] = webhook
		updatedWebhooksCount++
		docRef := Client.Collection(collection).Doc(webhook.WebhookID)
		batch.Set(docRef, map[string]interface{}{
			"Counter": webhook.Counter,
		}, firestore.MergeAll) // add the update operation to the batch

	}
	if updatedWebhooksCount > 0 {
		// Sends the batch request to Firebase
		batch.Flush()
	}

	// Update webhook cache
	for webhookID, updatedHook := range tmpCache {
		updatedHook.Modified = false
		webhookCache.cache[webhookID] = updatedHook
	}
	log.Printf("Finished SyncCacheToFirebase. Updated %d webhooks", updatedWebhooksCount)
}

func PeriodicSyncCache() {
	// Sync the cache with Firebase periodically
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		SyncCacheToFirebase()
	}
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

// registerWebhook adds a webhook to Firestore db and in-memory cache
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

	newWebhook.Counter = 0                           // Initialize counter for invocation
	docRef := Client.Collection(collection).NewDoc() // Generate a new document reference

	newWebhook.WebhookID = docRef.ID // Set the generated ID in the webhook data struct
	newWebhook.Modified = false      // Initialize modified flag

	var firebaseWebhook structs.WebhookFirebase // Create a new webhook struct to store the data from the database

	// Set the data from the request body to the new webhook struct
	firebaseWebhook.WebhookID = newWebhook.WebhookID
	firebaseWebhook.Url = newWebhook.Url
	firebaseWebhook.Country = newWebhook.Country
	firebaseWebhook.Calls = newWebhook.Calls
	firebaseWebhook.Counter = newWebhook.Counter

	// Add the webhook to the database with the generated ID
	_, err = docRef.Set(ctx, firebaseWebhook)
	if err != nil {
		// Error handling
		log.Println("Error when adding Webhook to database: ", err.Error())
		http.Error(w, "Error when adding Webhook to database: ", http.StatusBadRequest)
		return
	}

	// Add the webhook to the in-memory cache
	webhookCache.Lock()
	webhookCache.cache[newWebhook.WebhookID] = newWebhook
	webhookCache.Unlock()

	// Create a response body with the newly created webhook ID
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

// deleteWebhook deletes a webhook from Firestore db and the in-memory cache
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

	// Attempt to retrieve webhook from the in-memory cache
	webhookCache.Lock()
	data, ok := webhookCache.cache[id]
	webhookCache.Unlock()
	if !ok {
		log.Println("Webhook not found with ID: " + id)
		http.Error(w, "Webhook not found with ID: "+id, http.StatusNotFound)
		return
	}

	// Attempt to delete webhook from Firestore
	docRef := Client.Collection(collection).Doc(id)
	_, err := docRef.Delete(ctx)
	if err != nil {
		log.Println("Error deleting document " + id)
		http.Error(w, "Error deleting document "+id, http.StatusInternalServerError)
		return
	}

	// Remove webhook from the in-memory cache
	webhookCache.Lock()
	delete(webhookCache.cache, id)
	webhookCache.Unlock()

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

func retrieveWebhook(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Retrieve individual webhook if id is provided
	if len(parts) > 4 {
		id := parts[4]

		// Retrieve webhook from the in-memory cache
		webhookCache.RLock()
		webhook, exists := webhookCache.cache[id]
		webhookCache.RUnlock()

		if !exists {
			http.Error(w, "Webhook not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(webhook)
	} else {
		// Retrieve all webhooks if no id is provided
		webhookCache.RLock()
		defer webhookCache.RUnlock()

		hooks := make([]structs.WebhookGet, 0, len(webhookCache.cache))
		for _, hook := range webhookCache.cache {
			hooks = append(hooks, hook)
		}

		// Set the content type to JSON
		w.Header().Set("Content-Type", "application/json")
		// Set the status code to 200 (OK)
		w.WriteHeader(http.StatusOK)
		// Write the response body
		json.NewEncoder(w).Encode(hooks)
	}
}

func UpdateAndInvoke(isoCode string) {
	webhookCache.Lock()         // Change RLock to Lock
	defer webhookCache.Unlock() // Add defer to release lock after the function returns

	for webhookID, currentHook := range webhookCache.cache {
		if currentHook.Country == strings.ToUpper(isoCode) {
			currentHook.Counter++       // Increment counter
			currentHook.Modified = true // Reset modified flag
			// Checks if the counter is a multiple of amount of calls in the webhook
			if currentHook.Counter%currentHook.Calls == 0 {
				go invokeWebhook(currentHook) // Invoke webhook that matches the country
			}
			webhookCache.cache[webhookID] = currentHook // Update webhook cache with the updated webhook
		}
	}
}

// invokeWebhook invokes a webhook
func invokeWebhook(invoke structs.WebhookGet) {
	// Creates a new webhook invoke struct
	data := structs.WebhookInvoke{
		WebhookID: invoke.WebhookID,
		Country:   invoke.Country,
		Calls:     invoke.Counter,
	}
	// Creates a payload with the data that will be sent to the webhook
	payload, _ := json.Marshal(data)
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	_, err := http.Post(invoke.Url, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Println("Error during request creation. Error:", err)
		return
	}

}
