package main

import (
	"Assignment2"
	"Assignment2/handler"
	"Assignment2/stub/StubMain"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Stubbed handler in separate go routine
	go StubMain.StubMain("6969") // Pass the port number for the stub as an argument
	time.Sleep(2 * time.Second)

	log.Println("Main service starting...")
	// Retrieve potential environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: " + Assignment2.DEFAULT_PORT)
		port = Assignment2.DEFAULT_PORT
	}

	mux := http.NewServeMux() // Create a new ServeMux for the main application
	mux.HandleFunc(Assignment2.DEFAULT_PATH, handler.EmptyHandler)
	mux.HandleFunc(Assignment2.HISTORY_PATH, handler.HandlerHistory)
	mux.HandleFunc(Assignment2.NOTIFICATION_PATH, handler.HandlerNotifications)
	mux.HandleFunc(Assignment2.STATUS_PATH, handler.HandlerStatus)
	mux.HandleFunc(Assignment2.CURRENT_PATH, handler.HandlerRenewables)

	// Log that the main application is listening on the specified port
	log.Println("Main service listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, mux)) // Pass the ServeMux to ListenAndServe
}
