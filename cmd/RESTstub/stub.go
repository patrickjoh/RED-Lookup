package main

import (
	"Assignment2"
	"Assignment2/stub"
	"log"
	"net/http"
	"os"
)

func main() { // Accept the port number as an argument
	log.Println("Stub service starting...")

	// Retrieve potential environment variable
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: " + Assignment2.STUB_PORT)
		port = Assignment2.STUB_PORT
	}

	mux := http.NewServeMux()                        // Create a new ServeMux for the stub
	mux.HandleFunc("/Country/", stub.CountryHandler) // Register the stub handler
	mux.HandleFunc("/Neighbour/", stub.NeighbourHandler)

	log.Println("Stub service listening on port", port) // Log that the stub is listening on the specified port
	log.Fatal(http.ListenAndServe(":"+port, mux))       // Pass the ServeMux to ListenAndServe
}
