package main

import (
	"Assignment2"
	"Assignment2/handler"
	"log"
	"net/http"
	"os"
)

func main() {

	// Retrieve potential environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: " + Assignment2.DEFAULT_PORT)
		port = Assignment2.DEFAULT_PORT
	}

	http.HandleFunc("/"+Assignment2.DEFAULT_PATH, handler.EmptyHandler)
	http.HandleFunc("/"+Assignment2.HISTORY_PATH, handler.HandlerHistory)
	http.HandleFunc("/"+Assignment2.NOTIFICATION_PATH, handler.HandlerNotifications)
	http.HandleFunc("/"+Assignment2.STATUS_PATH, handler.HandlerStaus)
	http.HandleFunc("/"+Assignment2.CURRENT_PATH, handler.HandlerRenewables)

	log.Println("Listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
