package handler

import (
	"Assignment2"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var DocRefID string

// test post value
var sampleBody = map[string]interface{}{
	"url":     "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704",
	"country": "col",
	"calls":   69,
}

// code adds to firebase successfully
func TestRegisterWebhook(t *testing.T) {
	log.Println("test register webhook currently running")

	//initialise firebase connection
	InitFirebase()

	data, err := json.MarshalIndent(sampleBody, "", " ")
	if err != nil {
		log.Println("Error marshalling body.")
	}
	assert.Nil(t, err)

	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(data))
	if err != nil {
		log.Println("Error making request.")
	}

	responseR := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(responseR, request)

	if responseR.Code != http.StatusCreated {
		log.Println("Status code: ", responseR.Code)
	}
	log.Println("response here:", responseR.Code)

	// Check that the webhook was created successfully and get its ID
	var responseBody map[string]string
	require.Equal(t, http.StatusCreated, responseR.Code)
	_ = json.NewDecoder(responseR.Body).Decode(&responseBody)
	webhookID := responseBody["webhookId"]
	DocRefID = webhookID
	// Use the webhook ID to retrieve the webhook from Firestore
	docRef := Client.Collection(collection).Doc(webhookID)
	docSnapshot, err := docRef.Get(ctx)
	require.NoError(t, err)
	var retrievedWebhook Assignment2.WebhookGet
	err = docSnapshot.DataTo(&retrievedWebhook)
	require.NoError(t, err)

	// Check that the retrieved webhook has the expected properties
	assert.Equal(t, "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704", retrievedWebhook.Url)
	assert.Equal(t, "col", retrievedWebhook.Country)
	assert.Equal(t, int64(69), retrievedWebhook.Calls)
	assert.Equal(t, int64(0), retrievedWebhook.Counter)
	assert.Equal(t, "", retrievedWebhook.WebhookID)
}

// if correct id is given
func TestRetrieveWebhookWithID(t *testing.T) {
	log.Println("test retrive webhook currently running")

	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + DocRefID
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	if resp.StatusCode != http.StatusOK {
		log.Println("Status code: ", resp.StatusCode)
	}
	log.Println("response here:", resp.StatusCode)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	log.Println("hello: ")

	// Use the webhook ID to retrieve the webhook from Firestore
	docRef := Client.Collection(collection).Doc(DocRefID)
	docSnapshot, err := docRef.Get(ctx)
	require.NoError(t, err)
	var retrievedWebhook Assignment2.WebhookGet
	err = docSnapshot.DataTo(&retrievedWebhook)
	require.NoError(t, err)

	calls, err := strconv.Atoi(sampleBody["calls"].(string))
	if err != nil {
		fmt.Println(err)
	}

	// Check that the retrieved webhook has the expected properties
	assert.Equal(t, sampleBody["url"], retrievedWebhook.Url)
	assert.Equal(t, sampleBody["col"], retrievedWebhook.Country)
	assert.Equal(t, int64(calls), retrievedWebhook.Calls)
	assert.Equal(t, int64(0), retrievedWebhook.Counter)
	assert.Equal(t, "", retrievedWebhook.WebhookID)
}

// if no id is given
func TestRetrieveWebhookNoID(t *testing.T) {
	log.Println("test retrive webhook currently running")

	request, err := http.NewRequest(http.MethodGet, Assignment2.NOTIFICATION_PATH, nil)
	if err != nil {
		log.Println("Error making request.")
	}

	responseR := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(responseR, request)

	if responseR.Code != http.StatusOK {
		log.Println("Status code: ", responseR.Code)
	}
	log.Println("response here:", responseR.Code)
	assert.Equal(t, http.StatusOK, responseR.Code)
}

// if id is not found
func TestRetrieveWebhookNonExisting(t *testing.T) {
	log.Println("test retrive webhook currently running")

	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + "iiiiiiiii"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	if resp.StatusCode != http.StatusInternalServerError {
		log.Println("Status code: ", resp.StatusCode)
	}
	log.Println("response here:", resp.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// if correct id is given
func TestDeleteWebhookWithID(t *testing.T) {
	log.Println("test delete webhook currently running")

	//InitFirebase()

	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + DocRefID
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	if resp.StatusCode != http.StatusOK {
		log.Println("Status code: ", resp.StatusCode)
	}
	log.Println("response here:", resp.StatusCode)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

// if no id is given
func TestDeleteWebhookNoID(t *testing.T) {
	log.Println("test delete webhook currently running")

	request, err := http.NewRequest(http.MethodDelete, Assignment2.NOTIFICATION_PATH, nil)
	if err != nil {
		log.Println("Error making request.")
	}

	responseR := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(responseR, request)

	if responseR.Code != http.StatusBadRequest {
		log.Println("Status code: ", responseR.Code)
	}
	log.Println("response here:", responseR.Code)
	assert.Equal(t, http.StatusBadRequest, responseR.Code)
}

// if id is not found
func TestDeleteWebhookNonExistingID(t *testing.T) {
	log.Println("test delete webhook currently running")

	//InitFirebase()

	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + "iiiiiiiii"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	if resp.StatusCode != http.StatusInternalServerError {
		log.Println("Status code: ", resp.StatusCode)
	}
	log.Println("response here:", resp.StatusCode)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
