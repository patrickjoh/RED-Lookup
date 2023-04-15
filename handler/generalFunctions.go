package handler

import (
	"Assignment2"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

func convertCsvData() []Assignment2.CountryData {
	// Open CSV file
	file, error := os.Open(Assignment2.CSV_PATH)
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file") // UWU remove when done
	defer file.Close()

	// Read CSV file
	fileReader := csv.NewReader(file)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
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
