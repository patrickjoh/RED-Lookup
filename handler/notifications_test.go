package handler

import (
	"Assignment2"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestAddDocument(t *testing.T) {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", Assignment2.FIRESTORE_CREDS)
	// Initialize Firebase project configuration
	initFirebase()

	// Set up sample payload for testing
	payload := `{"testVal1":"testVal2"}`

	// Create HTTP request with sample payload
	req, err := http.NewRequest("POST", Assignment2.NOTIFICATION_PATH, strings.NewReader(payload))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Create mock HTTP response recorder
	rr := httptest.NewRecorder()

	// Execute HTTP request handler with mock request and response
	handler := http.HandlerFunc(addDocument)
	handler.ServeHTTP(rr, req)

	// Verify HTTP response matches expected output
	expectedResponse := `{"webhookId":"random-webhook-id"}`
	if rr.Body.String() != expectedResponse {
		t.Errorf("Handler returned unexpected body: got %v want %v",
			rr.Body.String(), expectedResponse)
	}
}

func testDeleteDocument() {

}

func testGetDocument() {

}
