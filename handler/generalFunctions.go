package handler

import (
	"Assignment2"
	"encoding/csv"
	"fmt"
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
		data := Assignment2.CountryData{
			Name:       r[0],
			IsoCode:    r[1],
			Year:       year,
			Percentage: r[3],
		}
		csvData = append(csvData, data)
	}

	return csvData
}

// Find all country data based on partial or full name of the country
func findCountry(countries []Assignment2.CountryData, partialName string) []Assignment2.CountryData {
	var countData []Assignment2.CountryData // empty list for the final data
	for _, col := range countries {
		if strings.Contains(col.Name, partialName) {
			newHisData := Assignment2.CountryData{
				Name:       col.Name,
				IsoCode:    col.IsoCode,
				Year:       col.Year,
				Percentage: col.Percentage,
			}
			countData = append(countData, newHisData)
		}
	}
	return countData
}
