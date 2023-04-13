package Assignment2

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// HandlerStudent main handler for route related to `/students` requests
// Note: here we are using a higher-order function with closure, to propagate a reference
// to the DB down the processing pipeline.
func HandlerStudent(db StudentsStorage) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handleStudentPost(w, r, db)
		case http.MethodGet:
			handleStudentGet(w, r, db)
		}
	}
}

// handleStudentPost utility function, package level, for handling POST request
func handleStudentPost(w http.ResponseWriter, r *http.Request, db StudentsStorage) {
	var s Student

	// TODO: Rewrite to harden decoding
	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Println("Error decoding JSON", err)
		return
	}

	// Validation: Check whether student is properly populated
	if s.StudentID == "" || s.Name == "" || s.Age == 0 {
		http.Error(w, "Input did not contain complete student specification. Recheck posted student information and resubmit.", http.StatusBadRequest)
		log.Println("Empty ID on student", s)
		return
	}
	// Check if the student is new
	_, ok := db.Get(s.StudentID)
	if ok {
		http.Error(w, "Student already exists. Use "+http.MethodPut+" or "+http.MethodPatch+" to modify.", http.StatusBadRequest)
		log.Println("Student with ID " + s.StudentID + "already exists.")
		return
	}
	// Add new student
	fmt.Println("Adding student to db ...")
	err = db.Add(s)
	if err != nil {
		http.Error(w, "Error when adding student to DB: "+err.Error(), http.StatusInternalServerError)
		log.Println("Error when adding student to DB", http.StatusInternalServerError)
		return
	}
	_, err = fmt.Fprint(w, http.StatusOK) // 200 by default
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("Error when returning error code 200: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleStudentGet utility function, package level, to handle GET request to student route
func handleStudentGet(w http.ResponseWriter, r *http.Request, db StudentsStorage) {

	http.Header.Add(w.Header(), "content-type", "application/json")
	// alternative way:
	// w.Header().Add("content-type", "application/json")

	parts := strings.Split(r.URL.Path, "/")

	// error handling
	if len(parts) != 3 || parts[1] != STUDENTS_PATH {
		http.Error(w, "Malformed URL", http.StatusBadRequest)
		log.Println("Malformed URL in request.")
		return
	}

	// handle the request /students/ which will return ALL students as array of JSON objects
	if parts[2] == "" {
		replyWithAllStudents(w, db)
	} else {
		replyWithStudent(w, db, parts[2])
	}
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
