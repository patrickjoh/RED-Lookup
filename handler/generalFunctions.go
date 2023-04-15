package handler

import (
	"Assignment2"
	"encoding/csv"
	"fmt"
	"os"
)

func convertCsvData() []Assignment2.CountData {
	// Open CSV file
	fd, error := os.Open("/handler/data/renewable-share-energy.csv")
	if error != nil {
		fmt.Println(error)
	}
	fmt.Println("Successfully opened the CSV file") // UWU remove when done
	defer fd.Close()

	// Read CSV file
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}

	var csvData []Assignment2.CountData
	for _, r := range records {
		if r[1] != "" { // Filtering out continents
			data := Assignment2.CountData{
				Name:       r[0],
				IsoCode:    r[1],
				Year:       r[2],
				Percentage: r[3],
			}
			csvData = append(csvData, data)
		}
	}
	return csvData
}
