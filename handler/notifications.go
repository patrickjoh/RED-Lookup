package handler

import (
	"net/http"
)

// HandlerStudent main handler for route related to `/students` requests
// Note: here we are using a higher-order function with closure, to propagate a reference
// to the DB down the processing pipeline.
func HandlerNotifications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		handleNotificationsPost(w, r)
	case http.MethodGet:
		handleNotificationsGet(w, r)
	}
}

// handleNotificationsPost utility function, package level, for handling POST request
func handleNotificationsPost(w http.ResponseWriter, r *http.Request) {

}

// handleNotificationsGet utility function, package level, to handle GET request to student route
func handleNotificationsGet(w http.ResponseWriter, r *http.Request) {

}
