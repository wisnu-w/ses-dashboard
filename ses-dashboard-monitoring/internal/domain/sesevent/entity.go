package sesevent

import "time"

type Event struct {
	ID                   int64
	MessageID            string
	Email                string
	Subject              string
	EventType            string
	Status               string
	Reason               string
	CreatedAt            time.Time
	Source               string
	Recipients           string // JSON array
	EventTimestamp       time.Time
	BounceType           string
	BounceSubType        string
	DiagnosticCode       string
	ProcessingTimeMillis int
	SmtpResponse         string
	RemoteMtaIp          string
	ReportingMTA         string
	Tags                 string // JSON map
}

type DailyMetrics struct {
	Date           string  `json:"date"`
	TotalEvents    int     `json:"total_events"`
	SendCount      int     `json:"send_count"`
	DeliveryCount  int     `json:"delivery_count"`
	BounceCount    int     `json:"bounce_count"`
	ComplaintCount int     `json:"complaint_count"`
	OpenCount      int     `json:"open_count"`
	ClickCount     int     `json:"click_count"`
	BounceRate     float64 `json:"bounce_rate"`
	DeliveryRate   float64 `json:"delivery_rate"`
}

type MonthlyMetrics struct {
	Month          string  `json:"month"`
	TotalEvents    int     `json:"total_events"`
	SendCount      int     `json:"send_count"`
	DeliveryCount  int     `json:"delivery_count"`
	BounceCount    int     `json:"bounce_count"`
	ComplaintCount int     `json:"complaint_count"`
	OpenCount      int     `json:"open_count"`
	ClickCount     int     `json:"click_count"`
	BounceRate     float64 `json:"bounce_rate"`
	DeliveryRate   float64 `json:"delivery_rate"`
}

type HourlyMetrics struct {
	Hour           string  `json:"hour"`
	TotalEvents    int     `json:"total_events"`
	SendCount      int     `json:"send_count"`
	DeliveryCount  int     `json:"delivery_count"`
	BounceCount    int     `json:"bounce_count"`
	ComplaintCount int     `json:"complaint_count"`
	OpenCount      int     `json:"open_count"`
	ClickCount     int     `json:"click_count"`
	BounceRate     float64 `json:"bounce_rate"`
	DeliveryRate   float64 `json:"delivery_rate"`
}