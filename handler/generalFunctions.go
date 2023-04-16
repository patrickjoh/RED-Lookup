package handler

import (
	"Assignment2"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func convertCsvData() []Assignment2.CountryData {
	// Open CSV file
	file, err := os.Open(Assignment2.CSV_PATH)
	if err != nil {
		fmt.Println(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file: ", err)
		}
	}(file)

	// Read CSV file
	fileReader := csv.NewReader(file)
	// Read and skip the header row
	_, err = fileReader.Read()
	if err != nil {
		fmt.Println(err)
	}

	records, err := fileReader.ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	var csvData []Assignment2.CountryData
	for _, r := range records {
		if r[1] == "" { // Filtering out continents
			continue
		}
		year, err := strconv.Atoi(r[2])
		if err != nil {
			fmt.Println(err)
		}
		percentage, err := strconv.ParseFloat(r[3], 64)
		if err != nil {
			fmt.Println(err)
		}
		data := Assignment2.CountryData{
			Name:       r[0],
			IsoCode:    r[1],
			Year:       year,
			Percentage: percentage,
		}
		csvData = append(csvData, data)
	}

	return csvData
}

// Find all country data for a given ISO code
func findCountry(countries []Assignment2.CountryData, Iso string) []Assignment2.CountryData {
	var countryData []Assignment2.CountryData // empty list for the final data

	// Convert to upper case to avoid case sensitivity
	Iso = strings.ToUpper(Iso)
	for _, col := range countries {
		if strings.Contains(col.IsoCode, Iso) {
			newHisData := Assignment2.CountryData{
				Name:       col.Name,
				IsoCode:    col.IsoCode,
				Year:       col.Year,
				Percentage: col.Percentage,
			}
			countryData = append(countryData, newHisData)
		}
	}
	return countryData
}

/*
function that gets all countries, decodes them and puts them in a slice. Returns the slice
*/
func getCountries(isoCode []string) ([]Assignment2.Country, error) {

	countryUrl := Assignment2.COUNTRYAPI_CODES
	// Loop through each ISO code and append the code the URL
	// Append each code to the URL with a comma delimiter
	isoCodesStr := strings.Join(isoCode, ",")
	countryUrl += isoCodesStr

	countryResponse, err := http.Get(countryUrl)
	if err != nil {
		return nil, err
	}
	defer countryResponse.Body.Close()

	// Decode the JSON response into a slice of "Country" structs
	var countryData []Assignment2.Country
	err = json.NewDecoder(countryResponse.Body).Decode(&countryData)
	if err != nil {
		return nil, err
	}

	// Check if any countries were found
	if len(countryData) == 0 {
		return nil, errors.New("No countries found")
	}

	return countryData, nil
}

/*
function that gets all countries, checks for redundancy and returns a struct of all countries with mean percentage
*/
func getAllCountriesMean() []Assignment2.CountryMean {
	var retData []Assignment2.CountryMean
	countries := convertCsvData()
	lastCountry := ""
	// going through all countries
	for _, current := range countries {
		// makes sure there are no repeated countries
		if current.IsoCode != lastCountry {
			lastCountry = current.IsoCode
			// variable that stores mean percentage
			var mean float64
			// calls findCountry functions to get all instances of one country
			countryHistory := findCountry(countries, current.IsoCode)
			// loops through all instances and adds their percentage to the mean variable
			for _, currentYear := range countryHistory {
				mean += currentYear.Percentage
			}
			// converts length of countryHistory slice into a float/number of instances
			numberInstances := float64(len(countryHistory))
			// calculates mean percentage of a country
			mean = mean / numberInstances
			// initiates a country struct with the mean percentage
			countryMean := Assignment2.CountryMean{
				Name:       current.Name,
				IsoCode:    current.IsoCode,
				Percentage: mean,
			}
			// appends country to slice of countries
			retData = append(retData, countryMean)
		}
	}
	return retData
}
