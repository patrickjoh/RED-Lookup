package stub

import (
	"net/http"
)

// Define a struct that represents the structure of your JSON data
type Neighbour struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

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

func handlerNeighbour(w http.ResponseWriter) {
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Parse the JSON file
	output := parseFile("stub/res/NOR_Neighbours.json")

	// Finally, let's write the JSON to the response
	w.Write(output)
}
