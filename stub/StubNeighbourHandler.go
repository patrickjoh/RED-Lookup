package stub

import (
	"log"
	"net/http"
)

// Neighbour Define a struct that represents the structure of your JSON data
type Neighbour struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

// NeighbourHandler is a function that handles the /neighbour stub endpoint
func NeighbourHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// If the request is a GET, then call the stubHandler
	case http.MethodGet:
		handlerNeighbour(w)
	default:
		// Otherwise, return an error
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return

}

// handlerNeighbour Handles json input and returns it
func handlerNeighbour(w http.ResponseWriter) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON file
	output := parseFile("stub/res/NOR_Neighbours.json")

	// Finally, let's write the JSON to the response
	_, err := w.Write(output)
	if err != nil {
		log.Println("Error sending response: ", err)
		http.Error(w, "Error sending response", http.StatusInternalServerError)
		return
	}
}
