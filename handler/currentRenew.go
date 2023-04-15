package handler

import (
	"net/http"
	"os"
)

// Documentation...
func HandlerRenewables(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandlerRenewablesPost(w, r)
	case http.MethodGet:
		HandlerRenewablesGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return

}

// HandlerRenewablesGet
func HandlerRenewablesGet(w http.ResponseWriter, r *http.Request) {
	// open CSV file
	fd, error := os.Open("data.csv")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	// read CSV file
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println(records)
}

// HandlerRenewablesPost
func HandlerRenewablesPost(w http.ResponseWriter, r *http.Request) {

}
