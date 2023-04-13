package Assignment2

type statusData struct {
	countries_api   string
	notification_db string
	webhooks        string
	version         string
	uptime          string
}

type countData struct {
	name       string
	isoCode    string
	year       int
	percentage float64
}

type webhookData struct {
	webhook_id string
	url        string
	country    string
	calls      int
}
