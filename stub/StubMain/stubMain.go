package StubMain

import (
	"Assignment2/stub"
	"log"
	"net/http"
	"os"
)

func StubMain() {

	// Define port
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 6969")
		port = "6969"
	}

	// Standard http server with reference to stubbed handler
	http.HandleFunc("/", stub.HandlerStub)

	log.Println("Running on port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}
