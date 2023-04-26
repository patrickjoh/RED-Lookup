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
	"os"
	"testing"
	"time"
)

// webhookID of the newly created webhook
var WebHookID string

// struct to avoid time variables
type webhookNoTimeVar struct {
	WebhookID string `json:"webhook_id" omitempty:"true"`
	Url       string `json:"url"`
	Country   string `json:"country"`
	Calls     int64  `json:"calls"`
	Counter   int64  `json:"counter" omitempty:"true"`
	Modified  bool   `json:"modified" omitempty:"true"`
}

// test post value
var sampleBody = map[string]interface{}{
	"webhook_id": "",
	"url":        "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704",
	"country":    "col",
	"calls":      69,
	"created":    time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC),
}

// TestInitFirebase Testing the initialization of firebase
func TestInitFirebase(t *testing.T) {
	//initialize firebase values
	err := InitFirebase()
	assert.Nil(t, err, "expected no error, error: %v", err)
}

// TestInitCache Testing the initialization of cache
func TestInitCache(t *testing.T) {
	//initialize cache values
	InitCache()

	// retrieve all data from Firebase and convert it all to structs.WebhookGet format
	docRef, err := Client.Collection(collection).Documents(ctx).GetAll()
	require.NoError(t, err)
	registeredWebhook := make(map[string]webhookNoTimeVar)
	for _, doc := range docRef {
		data := doc.Data()
		webhookID := doc.Ref.ID
		if calls, ok := data["Calls"].(int64); ok {
			newHook := webhookNoTimeVar{
				WebhookID: webhookID,
				Url:       data["Url"].(string),
				Country:   data["Country"].(string),
				Calls:     calls,
				Counter:   data["Counter"].(int64),
			}
			registeredWebhook[webhookID] = newHook
		}
	}

	//convert struct.WebhookGet to noTimeVar
	updatedCache := make(map[string]webhookNoTimeVar)
	for _, cache := range webhookCache.cache {
		updatedCache[cache.WebhookID] = webhookNoTimeVar{
			WebhookID: cache.WebhookID,
			Url:       cache.Url,
			Country:   cache.Country,
			Calls:     cache.Calls,
			Counter:   cache.Counter,
			Modified:  cache.Modified,
		}
	}

	//check if one of the value in the webhook is registered in the cache
	assert.Equal(t, registeredWebhook, updatedCache)
}

// TestRegisterWebhook test successful registerWebhook add to firebase successfully
func TestRegisterWebhook(t *testing.T) {

	// Change current working directory to the directory where the test file is located
	err := os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		// Reset the current working directory after the test has completed
		err := os.Chdir("./handler")
		if err != nil {
			t.Fatal(err)
		}
	}()
	// get all the csv-data
	Assignment2.CSVData = Assignment2.ConvertCsvData()

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
	sampleBody["webhook_id"] = webhookID

	// Use the webhook ID to retrieve the webhook from Firebase
	docRef := Client.Collection(collection).Doc(webhookID)
	docSnapshot, err := docRef.Get(ctx)
	require.NoError(t, err)
	var registeredWebhook webhookNoTimeVar
	err = docSnapshot.DataTo(&registeredWebhook)
	require.NoError(t, err)

	// Check that the retrieved webhook has the expected properties as the registered webhook
	assert.Equal(t, "https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704", registeredWebhook.Url)
	assert.Equal(t, "col", registeredWebhook.Country)
	assert.Equal(t, int64(69), registeredWebhook.Calls)
	assert.Equal(t, int64(0), registeredWebhook.Counter)
	assert.Equal(t, webhookID, registeredWebhook.WebhookID)
}

// TestSyncCacheToFirebase test if firebase is updated
func TestSyncCacheToFirebase(t *testing.T) {
	// start retrieval of all data from Firebase and convert it all to structs.WebhookGet format
	docRef, err := Client.Collection(collection).Documents(ctx).GetAll()
	require.NoError(t, err)

	startWebhook := make(map[string]webhookNoTimeVar)
	for _, doc := range docRef {
		data := doc.Data()
		webhookID := doc.Ref.ID
		if calls, ok := data["Calls"].(int64); ok {
			newHook := webhookNoTimeVar{
				WebhookID: webhookID,
				Url:       data["Url"].(string),
				Country:   data["Country"].(string),
				Calls:     calls,
				Counter:   data["Counter"].(int64),
			}
			startWebhook[webhookID] = newHook
		}
	}

	// create new data for the SyncCacheToFirebase to be updated
	//covert sampleData to []bytes
	data, err := json.MarshalIndent(sampleBody, "", " ")
	assert.Nil(t, err)
	//create request
	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(data))
	assert.Nil(t, err)
	resp := httptest.NewRecorder()
	handler := http.HandlerFunc(NotificationsHandler)
	handler.ServeHTTP(resp, request)

	// sync with SyncCacheToFirebase
	SyncCacheToFirebase()

	// updated retrieve of all data from Firebase and convert it all to structs.WebhookGet format
	docRef2, err := Client.Collection(collection).Documents(ctx).GetAll()
	require.NoError(t, err)
	endWebhook := make(map[string]webhookNoTimeVar)
	for _, doc := range docRef2 {
		data := doc.Data()
		webhookID := doc.Ref.ID
		if calls, ok := data["Calls"].(int64); ok {
			newHook := webhookNoTimeVar{
				WebhookID: webhookID,
				Url:       data["Url"].(string),
				Country:   data["Country"].(string),
				Calls:     calls,
				Counter:   data["Counter"].(int64),
			}
			endWebhook[webhookID] = newHook
		}
	}

	//convert struct.WebhookGet to noTimeVar
	updatedCache := make(map[string]webhookNoTimeVar)
	for _, cache := range webhookCache.cache {
		updatedCache[cache.WebhookID] = webhookNoTimeVar{
			WebhookID: cache.WebhookID,
			Url:       cache.Url,
			Country:   cache.Country,
			Calls:     cache.Calls,
			Counter:   cache.Counter,
			Modified:  cache.Modified,
		}
	}

	//check if the firebase is updated correctly
	require.NotEqual(t, startWebhook, endWebhook)

	//check if cache is successfully synced in the firebase
	assert.Equal(t, endWebhook, updatedCache)

}

// TestRegisterWebhookNoValue test unsuccessful testRegisterWebhook add fail
func TestRegisterWebhookNoValue(t *testing.T) {
	//create request
	sample := []byte(`{""}`)
	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(sample))
	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	// Check that the webhook was created successfully and get its ID
	require.Equal(t, http.StatusBadRequest, resp.Code)
}

// TestRegisterWebhookNoUrl test when no URL is given
func TestRegisterWebhookNoUrl(t *testing.T) {
	//create request
	sample := []byte(`{"url":"",
	"country": "COL",
	"calls":   69}`)
	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(sample))
	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	// Check that the webhook was created successfully and get its ID
	require.Equal(t, http.StatusBadRequest, resp.Code)
}

// TestRegisterWebhookInvalidCountry test registration of a invalid country
func TestRegisterWebhookInvalidCountry(t *testing.T) {
	//create request
	sample := []byte(`{"url":"https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704",
	"country": "GGG",
	"calls":   69}`)
	request, err := http.NewRequest(http.MethodPost, Assignment2.NOTIFICATION_PATH, bytes.NewReader(sample))
	assert.Nil(t, err)

	resp := httptest.NewRecorder()

	handler := http.HandlerFunc(NotificationsHandler)

	handler.ServeHTTP(resp, request)

	// Check that the webhook was created successfully and get its ID
	require.Equal(t, http.StatusBadRequest, resp.Code)
}

// TestRegisterWebhookNoValue test invalid number of calls
func TestRegisterWebhookInvalidCalls(t *testing.T) {
	//create request
	sample := []byte(`{"url":"https://webhook.site/63e2fb75-0742-44c1-9f14-fdd327649704",
	"country": "COL",
	"calls":   0}`)
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

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var recievedResponse webhookNoTimeVar
	_ = json.NewDecoder(resp.Body).Decode(&recievedResponse)

	// Convert sampleBody to JSON format
	expectedJsonData, err := json.Marshal(sampleBody)
	assert.Nil(t, err)

	// Unmarshal JSON data into WebhookGet struct
	var expectedWebhook webhookNoTimeVar
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

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
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

	var recievedResponse webhookNoTimeVar
	_ = json.NewDecoder(resp.Body).Decode(&recievedResponse)

	// Convert sampleBody to JSON format
	expectedJsonData, err := json.Marshal(sampleBody)
	assert.Nil(t, err)

	// Unmarshal JSON data into WebhookGet struct
	var expectedWebhook webhookNoTimeVar
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

// TestUpdateAndInvoke tests the counter mechanism in UpdateAndInvoke
func TestUpdateAndInvoke(t *testing.T) {
	webhookCache.cache = make(map[string]structs.WebhookGet)

	testHook := structs.WebhookGet{Country: "TCD", Counter: 0, Calls: 2}
	webhookCache.cache["TCD"] = testHook
	UpdateAndInvoke("TCD")
	UpdateAndInvoke("TCD")
	if webhookCache.cache["TCD"].Counter != 1 {
		t.Error("Counter incremented incorrectly")
	}

	if webhookCache.cache["TCD"].Counter != webhookCache.cache["TCD"].Counter {
		t.Error("Counter incremented incorrectly")
	}
}

// TestInvokeWebhook test if Invoked
func TestInvokeWebhook(t *testing.T) {
	// Create a mock HTTP server that will receive the POST request
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Check that the received content type is "application/json"
		if req.Header.Get("Content-Type") != "application/json" {
			t.Errorf("unexpected content type: %s", req.Header.Get("Content-Type"))
		}
		// Decode the received JSON payload into a WebhookInvoke struct
		var receivedData structs.WebhookInvoke
		err := json.NewDecoder(req.Body).Decode(&receivedData)
		assert.Nil(t, err, "error decoding received data: %v", err)

		// Check that the received data matches the expected data
		expectedData := structs.WebhookInvoke{
			WebhookID: "test",
			Country:   "test",
			Calls:     1,
		}
		if receivedData != expectedData {
			t.Errorf("unexpected received data: %+v", receivedData)
		}
	}))
	defer server.Close()

	// Call the function with test data
	testInvoke := structs.WebhookGet{
		WebhookID: "test",
		Country:   "test",
		Counter:   1,
		Url:       server.URL,
	}
	invokeWebhook(testInvoke)
}
