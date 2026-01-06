package usecase

import (
	"context"
	"time"

	"ses-monitoring/internal/domain/sesevent"
)

type SESUsecase struct {
	repo sesevent.Repository
}

func NewSESUsecase(repo sesevent.Repository) *SESUsecase {
	return &SESUsecase{repo: repo}
}

func (uc *SESUsecase) HandleEvent(
	ctx context.Context,
	event *sesevent.Event,
) error {
	return uc.repo.Save(ctx, event)
}

func (uc *SESUsecase) GetEvents(ctx context.Context) ([]*sesevent.Event, error) {
	return uc.repo.GetEvents(ctx)
}

func (uc *SESUsecase) GetEventsPaginated(ctx context.Context, limit, offset int) ([]*sesevent.Event, error) {
	return uc.repo.GetEventsPaginated(ctx, limit, offset)
}

func (uc *SESUsecase) GetEventsWithFilter(ctx context.Context, limit, offset int, search, startDate, endDate string) ([]*sesevent.Event, error) {
	return uc.repo.GetEventsWithFilter(ctx, limit, offset, search, startDate, endDate)
}

func (uc *SESUsecase) GetFilteredEventCount(ctx context.Context, search, startDate, endDate string) (int, error) {
	return uc.repo.GetFilteredEventCount(ctx, search, startDate, endDate)
}

func (uc *SESUsecase) GetEventCount(ctx context.Context) (int, error) {
	return uc.repo.GetEventCount(ctx)
}

func (uc *SESUsecase) GetEventsByType(ctx context.Context, eventType string) ([]*sesevent.Event, error) {
	return uc.repo.GetEventsByType(ctx, eventType)
}

func (uc *SESUsecase) GetBounceRate(ctx context.Context) (float64, error) {
	return uc.repo.GetBounceRate(ctx)
}

func (uc *SESUsecase) GetDeliveryRate(ctx context.Context) (float64, error) {
	return uc.repo.GetDeliveryRate(ctx)
}

func (uc *SESUsecase) GetDailyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.DailyMetrics, error) {
	return uc.repo.GetDailyMetrics(ctx, start, end)
}

func (uc *SESUsecase) GetMonthlyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.MonthlyMetrics, error) {
	return uc.repo.GetMonthlyMetrics(ctx, start, end)
}

func (uc *SESUsecase) GetHourlyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.HourlyMetrics, error) {
	return uc.repo.GetHourlyMetrics(ctx, start, end)
}

func (uc *SESUsecase) GetEventTypeCounts(ctx context.Context) (map[string]int, error) {
	return uc.repo.GetEventTypeCounts(ctx)
}
