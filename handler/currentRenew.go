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
	fd, error := os.Open("renewable-share-energy.csv")

}

// HandlerRenewablesPost
func HandlerRenewablesPost(w http.ResponseWriter, r *http.Request) {

}
