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
	output := "<h1>Welcome!<h1><h3>This service does not provide any functionality on root path level." +
		" Please try one of the paths below<h3>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Search for current percentage of renewables:<br>" +
		"<a href=\"" + Assignment2.CURRENT_PATH + "\">" + Assignment2.CURRENT_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Search for historical percentage of renewables:<br>" +
		"<a href=\"" + Assignment2.HISTORY_PATH + "\">" + Assignment2.HISTORY_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">Search for notification endpoint:<br>" +
		"<a href=\"" + Assignment2.NOTIFICATION_PATH + "\">" + Assignment2.NOTIFICATION_PATH + "</a></h5>" +
		"<h5 style=\"background-color: lightblue; width: 250px;\">For status:<br>" +
		"<a href=\"" + Assignment2.STATUS_PATH + "\">" + Assignment2.STATUS_PATH + "</a></h5>"

	// Write output to client
	_, err := fmt.Fprintf(w, "%v", output)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		log.Println("Error when returning output.")
	}
}
