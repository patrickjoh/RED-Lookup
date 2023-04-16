package handler

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
)

// Firebase context and client used by Firestore functions throughout the program.
var ctx context.Context
var client *firestore.Client

// Collection name in Firestore
const collection = "messages"

// Message counter to produce some variation in content
var ct = 0

func HandlerNotifications(w http.ResponseWriter, r *http.Request) {
	// Firebase initialisation
	ctx = context.Background()

	// We use a service account, load credentials file that you downloaded from your project's settings menu.
	// It should reside in your project directory.
	// Make sure this file is git-ignored, since it is the access token to the database.
	sa := option.WithCredentialsFile("assignment2-group11-firebase-adminsdk-lo4w1-716a00f3f5.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	// Instantiate client
	client, err = app.Firestore(ctx)

	// Check whether there is an error when connecting to Firestore
	if err != nil {
		log.Fatalln(err)
	}

	// Close down client at the end of the function
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

// handleNotificationsPost utility function, package level, for handling POST request
func handleNotificationsPost(w http.ResponseWriter, r *http.Request) {

}

// handleNotificationsGet utility function, package level, to handle GET request to student route
func handleNotificationsGet(w http.ResponseWriter, r *http.Request) {
	// Remove the trailing slash and split the URL into parts
	parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/"), "/")

	// Create webhook
	if len(parts) == 4 {
		// ...
	} else if len(parts) == 5 { // View webook
		// Get {id} or something idk
	} else {
		// some fucking error
	}
}

// handleNotificationsDelete utility function, package level, to handle GET request to student route
func handleNotificationsDelete(w http.ResponseWriter, r *http.Request) {

}
