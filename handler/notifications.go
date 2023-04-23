package handler

import (
	"Assignment2"
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

// Webhook DB
var webhooks = []Assignment2.WebhookInvoke{}

// var Mac hash.Hash
var Secret []byte

// Initialize signature (via init())
var SignatureKey = "X-SIGNATURE"

/*
initFirebase initializes the Firebase client and context.
taken from code example 13
*/
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

/*
Calls given URL with given content and awaits response (status and body).
*/ /*
func CallUrl(url string, method string, content string) {
	log.Println("Attempting invocation of url " + url + " with content '" + content + "'.")
	//res, err := http.Post(url, "text/plain", bytes.NewReader([]byte(content)))
	req, err := http.NewRequest(method, url, bytes.NewReader([]byte(content)))
	if err != nil {
		log.Printf("%v", "Error during request creation. Error:", err)
		return
	}

	/// BEGIN: HEADER GENERATION FOR CONTENT-BASED VALIDATION

	// Hash content (for content-based validation; not relevant for URL-based validation)
	mac := hmac.New(sha256.New, Secret)
	_, err = mac.Write([]byte(content))
	if err != nil {
		log.Printf("%v", "Error during content hashing. Error:", err)
		return
	}
	// Convert hash to string & add to header to transport to client for validation
	req.Header.Add(SignatureKey, hex.EncodeToString(mac.Sum(nil)))

	/// END: CONTENT-BASED VALIDATION

	// Perform invocation
	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("Error in HTTP request. Error:", err)
		return
	}

	// Read the response
	response, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("Something is wrong with invocation response. Error:", err)
		return
	}

	log.Println("Webhook " + url + " invoked. Received status code " + strconv.Itoa(res.StatusCode) +
		" and body: " + string(response))
}*/

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
	/*
		log.Println("Received POST request...")
		// Iterate through registered webhooks and invoke based on registered URL, method, and with received content
		for _, v := range webhooks {
			log.Println("Trigger event: Call to service endpoint with method POST" +
				" and content '" + string(text) + "'.")
			go CallUrl(v.WebhookID, v.Calls, v.Country)
		}*/

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
		/*_, err3 := fmt.Fprintln(w, m["text"])
		if err3 != nil {
			log.Println("Error while writing response body of message " + id)
			http.Error(w, "Error while writing response body of message "+id, http.StatusInternalServerError)
			return
		}*/
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

			/*_, err = fmt.Fprintln(w, m)
			if err != nil {
				log.Println("Error while writing response body (Error: " + err.Error() + ")")
				http.Error(w, "Error while writing response body (Error: "+err.Error()+")", http.StatusInternalServerError)
			}*/

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

func InvokeWebhook(id string) {
	// 1. Get webhook from id and parse into struct

	// 2. Get url from webhook

	// 3. Populate response struct with webhook data

	// 4. Send post request to url
}

/*		UWU pseudocode webhook invocation

Initializing call counting: (for method1)
	1. Create counting collection
	2. Copy iso codes from counting collection to local struct
	2. When new webhook is created:
		-> Go through counting collection
		-> If country does not exist in collection
			-> Add country to collection
			-> Add country to local struct

Counting calls: (method1)
	1. For each get request
		-> Check if country in request is in counting struct
		-> If it is
			-> Update count in counting collection for said country
			-> Call function that checks if any webhook should be invoked
				-> If webhook should be invoked
					-> Send relevant webhook id (or maybe webhook data) to invocation function

			NB1 problem: if count for countryX is 10, and new webhook is invoked for countryX
						for every 5 calls, should the webhook be invoked immediately?
						Should the webhook be invoked after the count has become 15?

			NB2 problem: requires separate collection for keeping track of call count

Counting calls: (method2)
	1. For each get request
		-> Save iso3code
		-> Loop through webhook struct (local copy of DB)
		-> For every webhook with matching iso3code
			-> Decrement calls by 1
			-> If the count for any webhook has become 0
				-> Call invocation function for said webhook (send webhook data or webhook id)
				-> Reset count to original number
			-> Update count in collection and struct

				NB problem: requires an extra field in webhook struct to remember original count


Invoking webhooks:
	1. Parse webhook into struct (or maybe webhook data is sent as parameter?)
	2. Get url from webhook
	3. Populate response struct with webhook data
	4. Send post request to url
*/
