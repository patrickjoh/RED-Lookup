package handler

import (
	"Assignment2"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
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
func TestHandleHistoryGet(t *testing.T) {
	// Create some sample country data
	// Create client instance
	client := http.Client{}

	// Retrieve content from server
	res, err := client.Get(URL + Assignment2.HISTORY_PATH + "NOR" + "?begin=1990&end=1992")
	if err != nil {
		t.Fatal("Get request to URL failed. Check whether server has been started manually! Error:", err.Error())
	}

	// Decode array
	result := []Assignment2.CountryData{}
	err2 := json.NewDecoder(res.Body).Decode(&result)
	if err2 != nil {
		t.Fatal("Error during decoding:", err2.Error())
	}

	//Expected output
	expected := []Assignment2.CountryData{
		{Name: "Norway", IsoCode: "NOR", Year: 1990, Percentage: 72.44774},
		{Name: "Norway", IsoCode: "NOR", Year: 1991, Percentage: 71.44005},
		{Name: "Norway", IsoCode: "NOR", Year: 1992, Percentage: 71.865555},
	}
	// Check the results
	assert.Equal(t, expected, result)

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
