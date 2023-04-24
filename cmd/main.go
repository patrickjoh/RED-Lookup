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

	// Close down client
	defer func() {
		err := handler.Client.Close()
		if err != nil {
			log.Fatal("Closing of the firebase client failed. Error:", err)
		}
	}()

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
