package handler

import (
	"Assignment2"
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
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
		handleNotificationsPost(w, r)
	case http.MethodGet:
		handleNotificationsGet(w, r)
	case http.MethodDelete:
		handleNotificationsDelete(w, r)
	}
}

func handleNotificationsPost(w http.ResponseWriter, r *http.Request) {

}

func handleNotificationsGet(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	if len(parts) == 4 {
		// Create webhook
	} else if len(parts) == 5 {
		// View webhook
	} else {
		// Handle error
	}
}

func handleNotificationsDelete(w http.ResponseWriter, r *http.Request) {

}

func postPayload(payload interface{}) {

}

func retrieveDocument(id string) {

}

// GetNumWebhooks retrieves and returns the number of registered webhooks from Firestore
func GetNumWebhooks() (int, error) {

	// Set up Firebase app
	opt := option.WithCredentialsFile(Assignment2.FIRESTORE_CREDS)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return 0, err
	}

	// Initialize Firestore client
	client, err := app.Firestore(context.Background())
	if err != nil {
		return 0, err
	}
	defer client.Close()

	// Create reference to webhook collection in Firestore
	webhooksCollection := client.Collection("webhooks")

	// Retrieve all webhooks
	iter := webhooksCollection.Documents(context.Background())
	var numWebhooks int
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		numWebhooks++
	}

	return numWebhooks, nil
}
