package handler

import (
	"Assignment2"
	"fmt"
	"log"
	"net/http"
)

//DefaultHandler
/*Handler for the default path*/
func EmptyHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure interpretation as HTML by client (browser)
	w.Header().Set("content-type", "text/html")

	// Offer information for redirection to paths
	output := "<h1>!placeholder!</h1>" +
		"<p>This service does not provide any functionality on the root path level. Please try one of the paths below:</p>" +
		"<h3>Search for current percentage of renewables:</h3>" +
		"<ul>" +
		"<li>" +
		"<a href=\"" + Assignment2.CURRENT_PATH + "\">" + Assignment2.CURRENT_PATH + "</a>" +
		"</li>" +
		"</ul>" +
		"<h3>Search for historical percentage of renewables:</h3>" +
		"<ul>" +
		"<li>" +
		"<a href=\"" + Assignment2.HISTORY_PATH + "\">" + Assignment2.HISTORY_PATH + "</a>" +
		"</li>" +
		"</ul>" +
		"<h3>Search for notification endpoint:</h3>" +
		"<ul>" +
		"<li>" +
		"<a href=\"" + Assignment2.NOTIFICATION_PATH + "\">" + Assignment2.NOTIFICATION_PATH + "</a>" +
		"</li>" +
		"</ul>" +
		"<h3>Check the status:</h3>" +
		"<ul>" +
		"<li>" +
		"<a href=\"" + Assignment2.STATUS_PATH + "\">" + Assignment2.STATUS_PATH + "</a>" +
		"</li>" +
		"</ul>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when returning output.")
	}
}
