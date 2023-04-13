package handler

import (
	"net/http"
)

// Documentation...
func HandlerHistory() func(http.ResponseWriter, *http.Request) {

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
func HandlerHistoryPost(w http.ResponseWriter, r *http.Request) {

}

// HandlerRenewablesGet
func HandlerHistoryGet(w http.ResponseWriter, r *http.Request) {

}
