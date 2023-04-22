package handler

import (
	"Assignment2"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddDocument(t *testing.T) {

	body := map[string]interface{}{
		"url":     "https://localhost:8080/client/",
		"country": "Norway",
		"calls":   5,
	}
	data, err := json.Marshal(body)
	if err != nil {
		log.Println("Error marshalling body.")
	}

	log.Println("UWU1")

	assert.Nil(t, err)

	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewBuffer(data))
	if err != nil {
		log.Println("Error making request.")
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationsHandler)

	log.Println("UWU3")
	handler.ServeHTTP(rr, request)
	log.Println("UWU4")

	if rr.Code != http.StatusCreated {
		log.Println("Status code: ", rr.Code)
	}

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func testDeleteDocument() {

}

func testGetDocument() {

}
