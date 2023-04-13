package Assignment2

type StatusData struct {
	countries_api   string
	notification_db string
	webhooks        string
	version         string
	uptime          string
}

type CountData struct {
	name       string
	isoCode    string
	year       int
	percentage float64
}

type WebhookData struct {
	webhook_id string
	url        string
	country    string
	calls      int
}
