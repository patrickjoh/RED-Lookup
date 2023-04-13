package handler

import (
	"net/http"
)

// Documentation...
func HandlerStaus(r http.ResponseWriter, w *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandelStatusGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// HandelStatusGet
func HandelStatusGet(w http.ResponseWriter, r *http.Request) {

}
