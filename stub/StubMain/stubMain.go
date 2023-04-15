package StubMain

import (
	"Assignment2/stub"
	"log"
	"net/http"
	"os"
)

func main() {

	// Define port
	port := os.Getenv("PORT")
	if port == "" {
		log.Println("$PORT has not been set. Default: 6969")
		port = "6969"
	}

	// Standard http server with reference to stubbed handler
	http.HandleFunc("countries/fetch", stub.HandlerStub)

	log.Println("Running on port", port)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

}
