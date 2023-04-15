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
