package Assignment2

type StatusData struct {
	Countries_api   string `json:"countries_api"`
	Notification_db string `json:"notifications_db"`
	Webhooks        int    `json:"webhooks"`
	Version         string `json:"version"`
	Uptime          string `json:"uptime"`
}

type CountData struct {
	Name       string `json:"name"`
	IsoCode    string `json:"isoCode"`
	Year       string `json:"year"`
	Percentage string `json:"percentage"`
}

type WebhookData struct {
	Webhook_id string `json:"webhook_id"`
	Url        string `json:"url"`
	Country    string `json:"country"`
	Calls      int    `json:"calls"`
}
