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

func TestRegisterWebhook(t *testing.T) {

	InitFirebase()

	//https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704
	body := map[string]interface{}{
		"url":     "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704",
		"country": "NOR",
		"calls":   3,
	}

	log.Println("bod", body)

	data, err := json.MarshalIndent(body, "", " ")
	if err != nil {
		log.Println("Error marshalling body.")
	}
	assert.Nil(t, err)
	log.Println("data", data)

	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(data))
	if err != nil {
		log.Println("Error making request.")
	}
	log.Println("resp", request)

	responseR := httptest.NewRecorder()
	log.Println("resp", responseR)

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(responseR, request)

	if responseR.Code != http.StatusCreated {
		log.Println("Status code: ", responseR.Code)
	}
	log.Println("response here:", responseR.Code)
	assert.Equal(t, http.StatusCreated, responseR.Code)
}

/*
func TestDeleteWebhook() {

}

func TestRetrieveWebhook() {

}*/
