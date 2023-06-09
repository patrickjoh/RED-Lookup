package main

import (
	"Assignment2"
	"Assignment2/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	log.Println("Main service starting...")
	// Retrieve potential environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: " + Assignment2.DEFAULT_PORT)
		port = Assignment2.DEFAULT_PORT
	}
	// Initialize Firebase
	err := handler.InitFirebase()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase: %v", err)
	}
	defer handler.Client.Close()

	handler.InitCache()                 // Initialize the cache
	go handler.PeriodicSyncCache()      // Start a goroutine to periodically sync the cache to Firebase
	go handler.RemoveExpiredWebhooks()  // Start a goroutine to remove expired webhooks
	defer handler.SyncCacheToFirebase() // Sync the cache to Firebase before closing the application

	mux := http.NewServeMux() // Create a new ServeMux for the main application
	mux.HandleFunc(Assignment2.DEFAULT_PATH, handler.EmptyHandler)
	mux.HandleFunc(Assignment2.HISTORY_PATH, handler.HistoryHandler)
	mux.HandleFunc(Assignment2.NOTIFICATION_PATH, handler.NotificationsHandler)
	mux.HandleFunc(Assignment2.STATUS_PATH, handler.StatusHandler)
	mux.HandleFunc(Assignment2.CURRENT_PATH, handler.RenewablesHandler)

	// Log that the main application is listening on the specified port
	log.Println("Main service listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux)) // Pass the ServeMux to ListenAndServe
}
