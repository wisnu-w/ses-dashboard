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

type MessageGroup struct {
	MessageID    string    `json:"message_id"`
	Email        string    `json:"email"`
	Subject      string    `json:"subject"`
	Source       string    `json:"source"`
	LatestStatus string    `json:"latest_status"`
	LatestEvent  string    `json:"latest_event"`
	EventTypes   []string  `json:"event_types"`
	EventCount   int       `json:"event_count"`
	FirstEventAt time.Time `json:"first_event_at"`
	LastEventAt  time.Time `json:"last_event_at"`
	HasBounce    bool      `json:"has_bounce"`
	HasComplaint bool      `json:"has_complaint"`
	HasDelivery  bool      `json:"has_delivery"`
	HasOpen      bool      `json:"has_open"`
	HasClick     bool      `json:"has_click"`
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
