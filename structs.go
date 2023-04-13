package Assignment2

type StatusData struct {
	Countries_api   string
	Notification_db string
	Version         string
	Uptime          string
}

type CountData struct {
	Name       string
	IsoCode    string
	Year       int
	Percentage float64
}

type WebhookData struct {
	Webhook_id string
	Url        string
	Country    string
	Calls      int
}
