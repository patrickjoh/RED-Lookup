package handler

import (
	"Assignment2"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Documentation...
func HandlerHistory(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		HandelHistoryPost(w, r)
	case http.MethodGet:
		HandelHistoryGet(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' is supported.", http.StatusNotImplemented)
	}
	return
}

// HandelHistoryPost
func HandelHistoryPost(w http.ResponseWriter, r *http.Request) {

}

// HandelHistoryGet
func HandelHistoryGet(w http.ResponseWriter, r *http.Request) {
	// Split url to get keyword
	urlKeywords := strings.Split(r.URL.Path, "/")

	// Error handling
	if len(urlKeywords) < 5 {
		log.Println(w, "Malformed URL", http.StatusBadRequest)
		return
	}

	fmt.Println("\nlen(keywords): ", len(urlKeywords))
	fmt.Println("\nkeywords[5]: ", urlKeywords[5])

	count := urlKeywords[5]
	query := r.URL.RawQuery

	// Parse the query string into a map
	params, err := url.ParseQuery(query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract values from the parameters
	begin := params.Get("begin")
	end := params.Get("end")

	// Send a response with the extracted values
	fmt.Fprintf(w, "Search query: country = %s, year (%s - %s)", count, begin, end)

	// open CSV file
	fd, error := os.Open("handler/data/renewable-share-energy.csv")
	if error != nil {
		fmt.Println("err occured", error)
	}
	fmt.Println("Successfully opened the CSV file")
	defer fd.Close()

	// read CSV file
	fileReader := csv.NewReader(fd)
	records, error := fileReader.ReadAll()
	if error != nil {
		fmt.Println(error)
	}

	var countData []Assignment2.HisData
	startYear, _ := strconv.Atoi(begin)
	endYear, _ := strconv.Atoi(end)

	for _, col := range records {
		year, _ := strconv.Atoi(col[2])
		if col[0] == count && year < endYear && year > startYear {
			fmt.Println("country found: ")
			for _, row := range col {
				fmt.Println(row + " ")
			}
			newHisData := Assignment2.HisData{
				Name:       col[0],
				IsoCode:    col[1],
				Year:       year,
				Percentage: col[3],
			}
			countData = append(countData, newHisData)
		}
	}

	fmt.Println(startYear, " ", endYear)
	fmt.Println(countData)
}
