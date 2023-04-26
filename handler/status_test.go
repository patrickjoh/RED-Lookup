package handler

import (
	"Assignment2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestStatusHandler tests the status endpoint with valid values
func TestStatusHandler(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(StatusHandler))
	defer server.Close()

	url := server.URL + Assignment2.STATUS_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// TestCurrentRenewTooManyParts tests current endpoint if there are too many parts
func TestStatusHandlerWrongPath(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(StatusHandler))
	defer server.Close()

	url := server.URL + Assignment2.CURRENT_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// TestStatusHandlerWrongMethod tests status endpoint for non-implemented methods
func TestStatusHandlerWrongMethod(t *testing.T) {

	server := httptest.NewServer(http.HandlerFunc(StatusHandler))
	defer server.Close()

	url := server.URL + Assignment2.NOTIFICATION_PATH
	response, err := http.Post(url, "application/json", nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotImplemented, response.StatusCode)
}
