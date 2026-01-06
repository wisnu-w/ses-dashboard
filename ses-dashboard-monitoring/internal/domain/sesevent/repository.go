package sesevent

import (
	"context"
	"time"
)

type Repository interface {
	Save(ctx context.Context, event *Event) error
	GetEvents(ctx context.Context) ([]*Event, error)
	GetEventsPaginated(ctx context.Context, limit, offset int) ([]*Event, error)
	GetEventsWithFilter(ctx context.Context, limit, offset int, search, startDate, endDate string) ([]*Event, error)
	GetFilteredEventCount(ctx context.Context, search, startDate, endDate string) (int, error)
	GetEventCount(ctx context.Context) (int, error)
	GetEventsByType(ctx context.Context, eventType string) ([]*Event, error)
	GetBounceRate(ctx context.Context) (float64, error)
	GetDeliveryRate(ctx context.Context) (float64, error)
	GetDailyMetrics(ctx context.Context, start, end *time.Time) ([]*DailyMetrics, error)
	GetMonthlyMetrics(ctx context.Context, start, end *time.Time) ([]*MonthlyMetrics, error)
	GetHourlyMetrics(ctx context.Context, start, end *time.Time) ([]*HourlyMetrics, error)
	GetEventTypeCounts(ctx context.Context) (map[string]int, error)
	DeleteOldEvents(ctx context.Context, cutoffDate time.Time) (int64, error)
}
