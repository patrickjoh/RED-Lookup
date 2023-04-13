package Assignment2

import (
	"encoding/json"
	"log"
	"net/http"
)

// Documentation...
func HandlerRenewables() func(http.ResponseWriter, *http.Request) {

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
func HandlerRenewablesPost(w http.ResponseWriter, r *http.Request) {

}

// HandlerRenewablesGet
func HandlerRenewablesGet(w http.ResponseWriter, r *http.Request) {

}
