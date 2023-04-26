package handler

import (
	"Assignment2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func changeWorkingDirectory(t *testing.T) {
	// Change current working directory to the directory where the test file is located
	err := os.Chdir("..")
	if err != nil {
		t.Fatal(err)
	}
}

// TestStatusHandler tests the status endpoint with valid values
func TestStatusHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(StatusHandler))
	defer server.Close()

	changeWorkingDirectory(t)
	defer func() {
		// Reset the current working directory after the test has completed
		err := os.Chdir("./handler")
		if err != nil {
			t.Fatal(err)
		}
	}()

	url := server.URL + Assignment2.STATUS_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// TestCurrentRenewTooManyParts tests current endpoint if there are too many parts
func TestStatusHandlerWrongPath(t *testing.T) {
	InitFirebase()

	server := httptest.NewServer(http.HandlerFunc(StatusHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// TestStatusHandlerWrongMethod tests status endpoint for non-implemented methods
func TestStatusHandlerWrongMethod(t *testing.T) {
	InitFirebase()

	server := httptest.NewServer(http.HandlerFunc(StatusHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH
	response, err := http.Post(url, "application/json", nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotImplemented, response.StatusCode)
}
