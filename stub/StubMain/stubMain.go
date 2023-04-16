package StubMain

import (
	"Assignment2/stub"
	"log"
	"net/http"
)

func StubMain(port string) { // Accept the port number as an argument
	log.Println("Stub service starting...")

	mux := http.NewServeMux()                        // Create a new ServeMux for the stub
	mux.HandleFunc("/Country/", stub.CountryHandler) // Register the stub handler
	mux.HandleFunc("/Neighbour/", stub.NeighbourHandler)

	log.Println("Stub service listening on port", port) // Log that the stub is listening on the specified port
	log.Fatal(http.ListenAndServe(":"+port, mux))       // Pass the ServeMux to ListenAndServe
}
