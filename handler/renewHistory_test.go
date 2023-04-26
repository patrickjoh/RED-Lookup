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

// test for not supported method
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

// Test if handler returns 200 ok with valid request
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

// tester for for errorhandling for om isokode som er skrevet inn enten er 3 eller 0
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
