package Assignment2

func init() {
	// Find current working directory

}

// Url Paths as ant
const (
	// ENDPOINTS
	DEFAULT_PATH      = "/"
	CURRENT_PATH      = "/energy/v1/renewables/current/"
	HISTORY_PATH      = "/energy/v1/renewables/history/"
	NOTIFICATION_PATH = "/energy/v1/notifications/"
	STATUS_PATH       = "/energy/v1/status/"

	// PORTS
	DEFAULT_PORT = "8080"
	STUB_PORT    = "8081"

	// EXTERNAL API
	COUNTRYAPI_CODES = "http://129.241.150.113:8080/v3.1/alpha?codes="

	// STUB ENDPOINTS
	//STUB_COUNTRY      = "http://localhost:8081/Country/" // Stubbed handler for country data
	//STUB_NEIGHBOURS   = "http://localhost:8081/Neighbour/" // Stubbed handler for neighbour data

	// DATA-FILES
	CSV_PATH        = "handler/data/renewable-share-energy.csv" // CSV_PATH is the path to the CSV file
	FIRESTORE_CREDS = "firebase.json"                           // FIRESTORE_CREDS is the path to the firestore credentials

	// TIME CONSTANTS
	WEBHOOK_EXPIRATION = 30 * 24 // WEBHOOK_EXPIRATION is the number of hours a webhook is valid for - days * hours
	WEBHOOK_AGE_CHECK  = 24      // WEBHOOK_AGE_CHECK is the number of hours between webhook age checks
	WEBHOOK_SYNC       = 5       // WEBHOOK_SYNC is the number of minutes between webhook syncs

)
