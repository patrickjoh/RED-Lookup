package handler

import (
	"Assignment2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Local URL
const URL = "http://localhost:" + Assignment2.DEFAULT_PORT

//test for status kode/error handling
/*
func TestHandleIsoCode(t *testing.T) {
	res, err := http.Get(URL + "/" + Assignment2.HISTORY_PATH + "/" + "NORE")
	if err != nil {
		t.Fatal("Isocode not accepted:", err.Error())
	}
	StatCode := httptest.NewRecorder()
}

/*
TestHandleHistoryGet tests the handleHistoryGet function
*/
/*
func TestHandleHistoryGet2(t *testing.T) {
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

	//Expected output
	expected := []structs.CountryData{
		{Name: "Norway", IsoCode: "NOR", Year: 1990, Percentage: 72.44774},
		{Name: "Norway", IsoCode: "NOR", Year: 1991, Percentage: 71.44005},
		{Name: "Norway", IsoCode: "NOR", Year: 1992, Percentage: 71.865555},
	}

	req, err := http.NewRequest(http.MethodGet, URL+Assignment2.HISTORY_PATH+"NOR"+"?begin=1990&end=1992", nil)
	assert.Nil(t, err, "Error creating request")
	rRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(HistoryHandler)
	handler.ServeHTTP(rRecorder, req)
	log.Println("UWU 7")
	var result []structs.CountryData
	log.Println("UWU 8: ", req.Body)
	jsondata := json.NewDecoder(req.Body).Decode(&result)
	log.Println("UWU 11")

	assert.Equal(t, http.StatusOK, rRecorder.Code, "Wrong status code")
	assert.Equal(t, expected, jsondata, "Wrong response body")
}

*/

/*
func TestHandleHistoryGet(t *testing.T) {
	// Create some sample country data
	// Create client instance
	client := http.Client{}
	// Retrieve content from server
	res, err := client.Get(URL + Assignment2.HISTORY_PATH + "NOR" + "?begin=1990&end=1992")
	if err != nil {
		t.Fatal("Get request to URL failed. Check whether server has been started manually! Error:", err.Error())
	}

	// Create a test request
	req, err := http.NewRequest("GET", URL+Assignment2.HISTORY_PATH+"NOR"+"?begin=1990&end=1992", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Decode array
	result := []structs.CountryData{}
	err := json.NewDecoder(req.Body).Decode(&result)
	if err2 != nil {
		t.Fatal("Error during decoding:", err2.Error())
	}

	// Call the handler function with the test request and response recorder
	handler := http.HandlerFunc(handleHistoryGet)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check that the status code is 200 OK
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	//Expected output
	expected := []structs.CountryData{
		{Name: "Norway", IsoCode: "NOR", Year: 1990, Percentage: 72.44774},
		{Name: "Norway", IsoCode: "NOR", Year: 1991, Percentage: 71.44005},
		{Name: "Norway", IsoCode: "NOR", Year: 1992, Percentage: 71.865555},
	}
	// Check the results
	assert.Equal(t, expected, result)
}
*/

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

/*
TestHandleMeanGet tests the mean functionality of handleHistoryGet
*/

/*
func TestHandleMeanGet(t *testing.T) {


		// Create some sample country data
		// Create client instance
		client := http.Client{}

		// URL under which server is instantiated
		fmt.Println("URL: ", server.URL)

		// Retrieve content from server
		res, err := client.Get(server.URL + "/" + STUDENTS_PATH + "/")
		if err != nil {
			t.Fatal("Get request to URL failed:", err.Error())
		}


	countries := []Assignment2.CountryMean{
		{
			Name:       "Australia",
			IsoCode:    "AUS",
			Percentage: 5.30004816,
		},
		{
			Name:       "Canada",
			IsoCode:    "CA",
			Percentage: 75.0,
		},
		{
			Name:       "Mexico",
			IsoCode:    "MX",
			Percentage: 85.0,
		},
	}

	// Call the function to get the country means
	triesMean(countries)

	// Check the results
	if len(result) != 3 {
		t.Errorf("Expected 3 countries, got %d", len(result))
	}

	for _, country := range result {
		if country.Name == "United States" {
			if country.Percentage != 75.0 {
				t.Errorf("Expected mean percentage of 75.0 for United States, got %f", country.Percentage)
			}
		} else if country.Name == "Canada" {
			if country.Percentage != 75.0 {
				t.Errorf("Expected mean percentage of 75.0 for Canada, got %f", country.Percentage)
			}
		} else if country.Name == "Mexico" {
			if country.Percentage != 87.5 {
				t.Errorf("Expected mean percentage of 87.5 for Mexico, got %f", country.Percentage)
			}
		}
	}
}

*/
