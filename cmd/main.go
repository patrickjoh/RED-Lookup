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
	log.Println("Main service starting...")
	// Stubbed handler in separate go routine
	go StubMain.StubMain()
	time.Sleep(2 * time.Second)

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

	log.Println("Main service listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
