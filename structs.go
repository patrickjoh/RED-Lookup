package Assignment2

type StatusData struct {
	Countries_api   string `json:"countries_api"`
	Notification_db string `json:"notifications_db"`
	Version         string `json:"version"`
	Uptime          string `json:"uptime"`
}

type CountData struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Year       int     `json:"year"`
	Percentage float64 `json:"percentage"`
}

type WebhookData struct {
	Webhook_id string `json:"webhook_id"`
	Url        string `json:"url"`
	Country    string `json:"country"`
	Calls      int    `json:"calls"`
}
