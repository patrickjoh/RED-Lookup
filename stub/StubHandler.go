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

func StubHandler(w http.ResponseWriter) {
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
