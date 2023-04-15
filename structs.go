package Assignment2

type StatusData struct {
	CountriesAPI   string `json:"countries_api"`
	NotificationDB string `json:"notifications_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         string `json:"uptime"`
}

type CountryData struct {
	Name       string `json:"name"`
	IsoCode    string `json:"isoCode"`
	Year       int    `json:"year"`
	Percentage string `json:"percentage"`
}

type WebhookData struct {
	WebhookID string `json:"webhook_id"`
	Url       string `json:"url"`
	Country   string `json:"country"`
	Calls     int    `json:"calls"`
}
