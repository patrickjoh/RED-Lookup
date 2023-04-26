package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// webhookID of the newly created webhook
var WebHookID string

// test post value
var sampleBody = map[string]interface{}{
	"url":     "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704",
	"country": "col",
	"calls":   69,
}

// mockFirestoreClient struct for testing
type mockFirestoreClient struct {
	data map[string]interface{}
}

// Create a mock Firestore client with test data
var mockClient = &mockFirestoreClient{
	data: map[string]interface{}{
		"webhook1": map[string]interface{}{
			"WebhookID": "webhook1",
			"Url":       "https://example.com",
			"Country":   "US",
			"Calls":     int64(2),
			"Counter":   int64(0),
		},
		"webhook2": map[string]interface{}{
			"WebhookID": "webhook2",
			"Url":       "https://example.org",
			"Country":   "CA",
			"Calls":     int64(3),
			"Counter":   int64(0),
		},
	},
}

// Testing the initialization of firebase
func TestInitFirebase(t *testing.T) {
	//initialize firebase values
	err := InitFirebase()
	assert.Nil(t, err, "expected no error, error: %v", err)
}

// test successful registerWebhook add to firebase successfully
func TestRegisterWebhook(t *testing.T) {
	//covert sampleData to []bytes
	data, err := json.MarshalIndent(sampleBody, "", " ")
	assert.Nil(t, err)

	//create request
	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(data))
	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	// Check that the webhook was created successfully and get its ID
	var responseBody map[string]string
	require.Equal(t, http.StatusCreated, resp.Code)
	_ = json.NewDecoder(resp.Body).Decode(&responseBody)
	webhookID := responseBody["webhookId"]
	WebHookID = webhookID

	// Use the webhook ID to retrieve the webhook from Firestore
	docRef := Client.Collection(collection).Doc(webhookID)
	docSnapshot, err := docRef.Get(ctx)
	require.NoError(t, err)
	var registeredWebhook structs.WebhookGet
	err = docSnapshot.DataTo(&registeredWebhook)
	require.NoError(t, err)

	// Check that the retrieved webhook has the expected properties as the registered webhook
	assert.Equal(t, "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704", registeredWebhook.Url)
	assert.Equal(t, "col", registeredWebhook.Country)
	assert.Equal(t, int64(69), registeredWebhook.Calls)
	assert.Equal(t, int64(0), registeredWebhook.Counter)
	assert.Equal(t, "", registeredWebhook.WebhookID)
}

// TestRegisterWebhookNoValue test unsuccessful testRegisterWebhook add fail
func TestRegisterWebhookNoValue(t *testing.T) {
	//create request
	sample := []byte(`{"url": "https://example.com"}`)
	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(sample))
	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	// Check that the webhook was created successfully and get its ID
	require.Equal(t, http.StatusBadRequest, resp.Code)
}

// TestRetrieveWebhookWithID test when existing id is given
func TestRetrieveWebhookWithID(t *testing.T) {
	//create request
	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + WebHookID
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

	var recievedResponse structs.WebhookGet
	_ = json.NewDecoder(resp.Body).Decode(&recievedResponse)

	// Convert sampleBody to JSON format
	expectedJsonData, err := json.Marshal(sampleBody)
	assert.Nil(t, err)

	// Unmarshal JSON data into WebhookGet struct
	var expectedWebhook structs.WebhookGet
	err = json.Unmarshal(expectedJsonData, &expectedWebhook)
	assert.Nil(t, err)

	assert.Equal(t, expectedWebhook, recievedResponse)
}

// TestRetrieveWebhookNoID test when no id is given, expect status code 200
func TestRetrieveWebhookNoID(t *testing.T) {
	request, err := http.NewRequest(http.MethodGet, Assignment2.NOTIFICATION_PATH, nil)
	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	assert.Equal(t, http.StatusOK, resp.Code)
}

// TestRetrieveWebhookNonExisting test when a not existing id is used
func TestRetrieveWebhookNonExisting(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + "iiiiiiiii"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// TestUpdateAndInvokeOK test if it works
func TestUpdateAndInvokeOK(t *testing.T) {
}

// TestDeleteWebhookWithID test when an existing id is given
func TestDeleteWebhookWithID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + WebHookID
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var recievedResponse structs.WebhookGet
	_ = json.NewDecoder(resp.Body).Decode(&recievedResponse)

	// Convert sampleBody to JSON format
	expectedJsonData, err := json.Marshal(sampleBody)
	assert.Nil(t, err)

	// Unmarshal JSON data into WebhookGet struct
	var expectedWebhook structs.WebhookGet
	err = json.Unmarshal(expectedJsonData, &expectedWebhook)
	assert.Nil(t, err)

	assert.Equal(t, expectedWebhook, recievedResponse)
}

// TestDeleteWebhookNoID test if no id is given
func TestDeleteWebhookNoID(t *testing.T) {
	request, err := http.NewRequest(http.MethodDelete, Assignment2.NOTIFICATION_PATH, nil)
	if err != nil {
		log.Println("Error making request.")
	}

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

// TestDeleteWebhookNonExistingID test when a not existing id is used
func TestDeleteWebhookNonExistingID(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(NotificationsHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH + "iiiiiiiii"
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	assert.Nil(t, err)

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
