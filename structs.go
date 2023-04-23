package Assignment2

type StatusData struct {
	CountriesAPI   string `json:"countries_api"`
	NotificationSB string `json:"notifications_db"`
	Webhooks       int    `json:"webhooks"`
	Version        string `json:"version"`
	Uptime         int    `json:"uptime"`
}

type CountryData struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Year       int     `json:"year"`
	Percentage float64 `json:"percentage"`
}
type CountryMean struct {
	Name       string  `json:"name"`
	IsoCode    string  `json:"isoCode"`
	Percentage float64 `json:"percentage"`
}

type WebhookGet struct {
	WebhookID string `json:"webhook_id"`
	Url       string `json:"url"`
	Country   string `json:"country"`
	Calls     int64  `json:"calls"`
	Counter   int64  `json:"counter"`
}
type WebhookPost struct {
	Url     string `json:"url"`
	Country string `json:"country"`
	Calls   int64  `json:"calls"`
}

type WebhookInvoke struct {
	WebhookID string `json:"webhook_id"`
	Country   string `json:"country"`
	Calls     int64  `json:"calls"`
}

// Country Struct for storing data from REST Countries API
type Country struct {
	Alpha3Code string   `json:"cca3"`
	Border     []string `json:"borders"`
}
