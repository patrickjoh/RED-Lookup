package stub

import (
	"fmt"
	"net/http"
	"os"
)

// Define a struct that represents the structure of your JSON data
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func HandlerStub(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// If the request is a GET, then call the stubHandler
	case http.MethodGet:
		stubHandler(w)
	default:
		// Otherwise, return an error
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return

}

func stubHandler(w http.ResponseWriter) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON file
	output := parseFile("./res/RESTCountry.json")

	// Finally, let's write the JSON to the response
	w.Write(output)
}

func parseFile(filename string) []byte {
	// Open and
	file, e := os.ReadFile(filename)
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}
	return file
}
