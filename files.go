package Assignment2

import (
	"Assignment2/structs"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// FirebaseCredentials is the credentials for the Firebase project
var FirebaseCredentials []byte

func init() {
	var err error // Error variable

	FirebaseCredentials, err = os.ReadFile(FIRESTORE_CREDS)
	if err != nil {
		log.Println("Error reading firestore credentials: ", err)
	}
}

// CSVData is the data from the CSV file
var CSVData []structs.CountryData

func init() {
	CSVData = ConvertCsvData()
}

// ConvertCsvData takes data from a csv file and converts it to a slice of structs
func ConvertCsvData() []structs.CountryData {
	// Open CSV file
	file, err := os.Open(CSV_PATH)
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

	var csvData []structs.CountryData
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
		data := structs.CountryData{
			Name:       r[0],
			IsoCode:    r[1],
			Year:       year,
			Percentage: percentage,
		}
		csvData = append(csvData, data)
	}

	return csvData
}
