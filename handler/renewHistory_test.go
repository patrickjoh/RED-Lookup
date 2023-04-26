package handler

import (
	"Assignment2"
	"Assignment2/structs"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Local URL
const URL = "http://localhost:" + Assignment2.DEFAULT_PORT

// TestHandleHistoryGet Tests main functionality of the history handler via the getFromBeginToEnd function
func TestHandleHistoryGet(t *testing.T) {
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

	sampleData := []structs.CountryData{
		{Name: "Norway", IsoCode: "NOR", Year: 1990, Percentage: 72.44774},
		{Name: "Norway", IsoCode: "NOR", Year: 1991, Percentage: 71.44005},
		{Name: "Norway", IsoCode: "NOR", Year: 1992, Percentage: 71.865555},
		{Name: "Norway", IsoCode: "NOR", Year: 1993, Percentage: 71.17737},
	}

	//Expected output
	expected := []structs.CountryData{
		{Name: "Norway", IsoCode: "NOR", Year: 1990, Percentage: 72.44774},
		{Name: "Norway", IsoCode: "NOR", Year: 1991, Percentage: 71.44005},
		{Name: "Norway", IsoCode: "NOR", Year: 1992, Percentage: 71.865555},
	}

	result := getFromBeginToEnd(1990, 1992, sampleData)

	assert.Equal(t, expected, result)

}

// TestMethodNotImplemented tests if the correct statuscode is returned when the method is not implemented
func TestMethodNotImplemented(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(HistoryHandler))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH + "NOR"
	response, err := http.Post(url, "", nil)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotImplemented, response.StatusCode)
}

// TestFunctionality tests if the function returns 200 OK with correct url
func TestFunctionalityStatusOK(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(HistoryHandler))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH + "NOR"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)

}

// TestMalformedURL tests if the function returns 400 bad request with malformed url
func TestMalformedURL(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(handleHistoryGet))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH + "NORE"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)

}

// TestBeginEnd Tests error handling if begin is greater than end and that the correct statuscode is returned
func TestBeginEnd(t *testing.T) {
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
	server := httptest.NewServer(http.HandlerFunc(handleHistoryGet))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH + "NOR" + "?begin=1992&end=1990"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusBadRequest, response.StatusCode)
}

// Tests if the function returns 404 not found with invalid request and no countries are found
func TestNoCountryFound(t *testing.T) {
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
	server := httptest.NewServer(http.HandlerFunc(handleHistoryGet))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH + "HIJ"
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusNotFound, response.StatusCode)
}

// Test  if the function returns 200 ok with valid request of no isocodes in the url
func TestNoIsocode(t *testing.T) {
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
	server := httptest.NewServer(http.HandlerFunc(handleHistoryGet))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// TestPercentageSort tests if the function percentageSort returns 200 ok
func TestPercentageSort(t *testing.T) {
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
	server := httptest.NewServer(http.HandlerFunc(handleHistoryGet))
	defer server.Close()

	url := server.URL + Assignment2.HISTORY_PATH
	response, err := http.Get(url)

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, response.StatusCode)
}

// TestGetCountriesMean tests if the function getAllCountriesMean returns the correct mean
func TestGetCountriesMean(t *testing.T) {
	sampleData := []structs.CountryData{
		{Name: "Norway", IsoCode: "NOR", Year: 1990, Percentage: 72.44774},
		{Name: "Norway", IsoCode: "NOR", Year: 1991, Percentage: 71.44005},
	}

	expected := []structs.CountryMean{
		{Name: "Norway", IsoCode: "NOR", Percentage: 71.943895},
	}
	result := getAllCountriesMean(sampleData)

	assert.Equal(t, expected, result)
}
