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
const collection = "webhooks"

// Message counter to produce some variation in content
var ct = 0

func initFirebase() {
	ctx = context.Background()

	sa := option.WithCredentialsFile("assignment2-group11-firebase-adminsdk-lo4w1-716a00f3f5.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err = app.Firestore(ctx)
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
