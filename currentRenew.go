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
			handleStudentPost(w, r)
		case http.MethodGet:
			handleStudentGet(w, r)
		}
	}
}

// HandlerRenewablesPost
func HandlerRenewablesPost(w http.ResponseWriter, r *http.Request) {

}

// HandlerRenewablesGet
func HandlerRenewablesGet(w http.ResponseWriter, r *http.Request) {

}

// Reply handlers with DB interaction

// replyWithAllStudents prepares a response with all students from the student storage
func replyWithAllStudents(w http.ResponseWriter, db StudentsStorage) {
	if db.Count() == 0 {
		err := json.NewEncoder(w).Encode([]Student{})
		if err != nil {
			// this should never happen
			http.Error(w, "Error during JSON encoding of an empty array: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error during JSON encoding of an empty array", err)
			return
		}
	} else {
		a := make([]Student, 0, db.Count())
		a = append(a, db.GetAll()...)
		err := json.NewEncoder(w).Encode(a)
		if err != nil {
			http.Error(w, "Error during JSON encoding: "+err.Error(), http.StatusInternalServerError)
			log.Println("Error during JSON encoding", err)
			return
		}
	}
}

// replyWithStudent prepares a response with a single student from the student storage
func replyWithStudent(w http.ResponseWriter, db StudentsStorage, id string) {
	// make sure that id is valid
	s, ok := db.Get(id)
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	// handle /student/<id>
	err := json.NewEncoder(w).Encode(s)
	if err != nil {
		http.Error(w, "Error during JSON encoding: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error during JSON encoding", err)
		return
	}
}
