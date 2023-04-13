package handler

import (
	"net/http"
)

// Documentation...
func HandlerStaus() func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			HandlerRenewablesPost(w, r)
		case http.MethodGet:
			HandlerRenewablesGet(w, r)
		}
	}
}

// HandlerRenewablesPost
func HandlerStatusPost(w http.ResponseWriter, r *http.Request) {

}

// HandlerRenewablesGet
func HandlerStatusGet(w http.ResponseWriter, r *http.Request) {

}
