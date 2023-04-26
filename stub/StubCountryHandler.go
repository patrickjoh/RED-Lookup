package stub

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// Country is a struct that represents the structure of your JSON data
type Country struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

func CountryHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// If the request is a GET, then call the stubHandler
	case http.MethodGet:
		handlerCountry(w)
	default:
		// Otherwise, return an error
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return

}

func handlerCountry(w http.ResponseWriter) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON file
	output := parseFile("stub/res/NOR_Country.json")

	// Finally, let's write the JSON to the response
	_, err := w.Write(output)
	if err != nil {
		log.Println("Error sending response: ", err)
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}
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
