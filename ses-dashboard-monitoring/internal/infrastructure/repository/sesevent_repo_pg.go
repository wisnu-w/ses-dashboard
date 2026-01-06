package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"ses-monitoring/internal/domain/sesevent"
)

type sesEventRepo struct {
	db *sql.DB
}

func NewSESEventRepository(db *sql.DB) sesevent.Repository {
	return &sesEventRepo{db: db}
}

func (r *sesEventRepo) Save(ctx context.Context, e *sesevent.Event) error {
	query := `
		INSERT INTO ses_events (
			message_id, email, subject, event_type, status, reason, source, recipients,
			event_timestamp, bounce_type, bounce_sub_type, diagnostic_code,
			processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		e.MessageID,
		e.Email,
		e.Subject,
		e.EventType,
		e.Status,
		e.Reason,
		e.Source,
		e.Recipients,
		e.EventTimestamp,
		e.BounceType,
		e.BounceSubType,
		e.DiagnosticCode,
		e.ProcessingTimeMillis,
		e.SmtpResponse,
		e.RemoteMtaIp,
		e.ReportingMTA,
		e.Tags,
	)
	return err
}

func (r *sesEventRepo) GetEvents(ctx context.Context) ([]*sesevent.Event, error) {
	query := `
		SELECT message_id, email, subject, event_type, status, reason, source, recipients,
			   event_timestamp, bounce_type, bounce_sub_type, diagnostic_code,
			   processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags
		FROM ses_events
		ORDER BY event_timestamp DESC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*sesevent.Event
	for rows.Next() {
		e := &sesevent.Event{}
		err := rows.Scan(
			&e.MessageID,
			&e.Email,
			&e.Subject,
			&e.EventType,
			&e.Status,
			&e.Reason,
			&e.Source,
			&e.Recipients,
			&e.EventTimestamp,
			&e.BounceType,
			&e.BounceSubType,
			&e.DiagnosticCode,
			&e.ProcessingTimeMillis,
			&e.SmtpResponse,
			&e.RemoteMtaIp,
			&e.ReportingMTA,
			&e.Tags,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *sesEventRepo) GetEventsPaginated(ctx context.Context, limit, offset int) ([]*sesevent.Event, error) {
	query := `
		SELECT message_id, email, subject, event_type, status, reason, source, recipients,
			   event_timestamp, bounce_type, bounce_sub_type, diagnostic_code,
			   processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags
		FROM ses_events
		ORDER BY event_timestamp DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*sesevent.Event
	for rows.Next() {
		e := &sesevent.Event{}
		err := rows.Scan(
			&e.MessageID,
			&e.Email,
			&e.Subject,
			&e.EventType,
			&e.Status,
			&e.Reason,
			&e.Source,
			&e.Recipients,
			&e.EventTimestamp,
			&e.BounceType,
			&e.BounceSubType,
			&e.DiagnosticCode,
			&e.ProcessingTimeMillis,
			&e.SmtpResponse,
			&e.RemoteMtaIp,
			&e.ReportingMTA,
			&e.Tags,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *sesEventRepo) GetEventsWithFilter(ctx context.Context, limit, offset int, search, startDate, endDate string) ([]*sesevent.Event, error) {
	query := `
		SELECT message_id, email, subject, event_type, status, reason, source, recipients,
			   event_timestamp, bounce_type, bounce_sub_type, diagnostic_code,
			   processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags
		FROM ses_events
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 0

	if search != "" {
		argIndex++
		query += fmt.Sprintf(" AND (email ILIKE $%d OR subject ILIKE $%d OR source ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+search+"%")
	}

	if startDate != "" {
		argIndex++
		query += fmt.Sprintf(" AND event_timestamp >= $%d", argIndex)
		args = append(args, startDate)
	}

	if endDate != "" {
		argIndex++
		query += fmt.Sprintf(" AND event_timestamp <= $%d", argIndex)
		args = append(args, endDate+" 23:59:59")
	}

	query += " ORDER BY event_timestamp DESC"
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex+1, argIndex+2)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*sesevent.Event
	for rows.Next() {
		e := &sesevent.Event{}
		err := rows.Scan(
			&e.MessageID,
			&e.Email,
			&e.Subject,
			&e.EventType,
			&e.Status,
			&e.Reason,
			&e.Source,
			&e.Recipients,
			&e.EventTimestamp,
			&e.BounceType,
			&e.BounceSubType,
			&e.DiagnosticCode,
			&e.ProcessingTimeMillis,
			&e.SmtpResponse,
			&e.RemoteMtaIp,
			&e.ReportingMTA,
			&e.Tags,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *sesEventRepo) GetFilteredEventCount(ctx context.Context, search, startDate, endDate string) (int, error) {
	query := `SELECT COUNT(*) FROM ses_events WHERE 1=1`
	args := []interface{}{}
	argIndex := 0

	if search != "" {
		argIndex++
		query += fmt.Sprintf(" AND (email ILIKE $%d OR subject ILIKE $%d OR source ILIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, "%"+search+"%")
	}

	if startDate != "" {
		argIndex++
		query += fmt.Sprintf(" AND event_timestamp >= $%d", argIndex)
		args = append(args, startDate)
	}

	if endDate != "" {
		argIndex++
		query += fmt.Sprintf(" AND event_timestamp <= $%d", argIndex)
		args = append(args, endDate+" 23:59:59")
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *sesEventRepo) GetEventCount(ctx context.Context) (int, error) {
	query := `SELECT COUNT(DISTINCT message_id) FROM ses_events`
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *sesEventRepo) GetEventsByType(ctx context.Context, eventType string) ([]*sesevent.Event, error) {
	query := `
		SELECT message_id, email, subject, event_type, status, reason, source, recipients,
			   event_timestamp, bounce_type, bounce_sub_type, diagnostic_code,
			   processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags
		FROM ses_events
		WHERE event_type = $1
		ORDER BY event_timestamp DESC
	`
	rows, err := r.db.QueryContext(ctx, query, eventType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*sesevent.Event
	for rows.Next() {
		e := &sesevent.Event{}
		err := rows.Scan(
			&e.MessageID,
			&e.Email,
			&e.Subject,
			&e.EventType,
			&e.Status,
			&e.Reason,
			&e.Source,
			&e.Recipients,
			&e.EventTimestamp,
			&e.BounceType,
			&e.BounceSubType,
			&e.DiagnosticCode,
			&e.ProcessingTimeMillis,
			&e.SmtpResponse,
			&e.RemoteMtaIp,
			&e.ReportingMTA,
			&e.Tags,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, e)
	}
	return events, nil
}

func (r *sesEventRepo) GetBounceRate(ctx context.Context) (float64, error) {
	query := `
		SELECT 
			CASE 
				WHEN (SELECT COUNT(DISTINCT message_id) FROM ses_events) = 0 THEN 0
				ELSE (SELECT COUNT(DISTINCT message_id) FROM ses_events WHERE event_type = 'Bounce') * 100.0 / (SELECT COUNT(DISTINCT message_id) FROM ses_events)
			END
	`
	var rate float64
	err := r.db.QueryRowContext(ctx, query).Scan(&rate)
	return rate, err
}

func (r *sesEventRepo) GetDeliveryRate(ctx context.Context) (float64, error) {
	query := `
		SELECT 
			CASE 
				WHEN (SELECT COUNT(DISTINCT message_id) FROM ses_events) = 0 THEN 0
				ELSE (SELECT COUNT(DISTINCT message_id) FROM ses_events WHERE event_type = 'Delivery') * 100.0 / (SELECT COUNT(DISTINCT message_id) FROM ses_events)
			END
	`
	var rate float64
	err := r.db.QueryRowContext(ctx, query).Scan(&rate)
	return rate, err
}

func (r *sesEventRepo) GetDailyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.DailyMetrics, error) {
	query := `
		SELECT 
			DATE(event_timestamp) as date,
			COUNT(DISTINCT message_id) as total_events,
			COUNT(DISTINCT CASE WHEN event_type = 'Send' THEN message_id END) as send_count,
			COUNT(DISTINCT CASE WHEN event_type = 'Delivery' THEN message_id END) as delivery_count,
			COUNT(DISTINCT CASE WHEN event_type = 'Bounce' THEN message_id END) as bounce_count,
			COUNT(DISTINCT CASE WHEN event_type = 'Complaint' THEN message_id END) as complaint_count,
			COUNT(DISTINCT CASE WHEN event_type = 'Open' THEN message_id END) as open_count,
			COUNT(DISTINCT CASE WHEN event_type = 'Click' THEN message_id END) as click_count,
			CASE WHEN COUNT(DISTINCT message_id) = 0 THEN 0 ELSE (COUNT(DISTINCT CASE WHEN event_type = 'Bounce' THEN message_id END) * 100.0 / COUNT(DISTINCT message_id)) END as bounce_rate,
			CASE WHEN COUNT(DISTINCT message_id) = 0 THEN 0 ELSE (COUNT(DISTINCT CASE WHEN event_type = 'Delivery' THEN message_id END) * 100.0 / COUNT(DISTINCT message_id)) END as delivery_rate
		FROM ses_events
	`
	args := []interface{}{}
	conditions := []string{}
	if start != nil {
		args = append(args, *start)
		conditions = append(conditions, fmt.Sprintf("event_timestamp >= $%d", len(args)))
	}
	if end != nil {
		args = append(args, *end)
		conditions = append(conditions, fmt.Sprintf("event_timestamp < $%d", len(args)))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += `
		GROUP BY DATE(event_timestamp)
		ORDER BY DATE(event_timestamp) DESC
	`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*sesevent.DailyMetrics
	for rows.Next() {
		m := &sesevent.DailyMetrics{}
		err := rows.Scan(&m.Date, &m.TotalEvents, &m.SendCount, &m.DeliveryCount, &m.BounceCount, &m.ComplaintCount, &m.OpenCount, &m.ClickCount, &m.BounceRate, &m.DeliveryRate)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (r *sesEventRepo) GetMonthlyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.MonthlyMetrics, error) {
	query := `
		SELECT 
			TO_CHAR(DATE_TRUNC('month', event_timestamp), 'YYYY-MM') as month,
			COUNT(*) as total_events,
			SUM(CASE WHEN event_type = 'Send' THEN 1 ELSE 0 END) as send_count,
			SUM(CASE WHEN event_type = 'Delivery' THEN 1 ELSE 0 END) as delivery_count,
			SUM(CASE WHEN event_type = 'Bounce' THEN 1 ELSE 0 END) as bounce_count,
			SUM(CASE WHEN event_type = 'Complaint' THEN 1 ELSE 0 END) as complaint_count,
			SUM(CASE WHEN event_type = 'Open' THEN 1 ELSE 0 END) as open_count,
			SUM(CASE WHEN event_type = 'Click' THEN 1 ELSE 0 END) as click_count,
			CASE WHEN COUNT(*) = 0 THEN 0 ELSE (SUM(CASE WHEN event_type = 'Bounce' THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) END as bounce_rate,
			CASE WHEN COUNT(*) = 0 THEN 0 ELSE (SUM(CASE WHEN event_type = 'Delivery' THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) END as delivery_rate
		FROM ses_events
	`
	args := []interface{}{}
	conditions := []string{}
	if start != nil {
		args = append(args, *start)
		conditions = append(conditions, fmt.Sprintf("event_timestamp >= $%d", len(args)))
	}
	if end != nil {
		args = append(args, *end)
		conditions = append(conditions, fmt.Sprintf("event_timestamp < $%d", len(args)))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += `
		GROUP BY DATE_TRUNC('month', event_timestamp)
		ORDER BY DATE_TRUNC('month', event_timestamp) DESC
	`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*sesevent.MonthlyMetrics
	for rows.Next() {
		m := &sesevent.MonthlyMetrics{}
		err := rows.Scan(&m.Month, &m.TotalEvents, &m.SendCount, &m.DeliveryCount, &m.BounceCount, &m.ComplaintCount, &m.OpenCount, &m.ClickCount, &m.BounceRate, &m.DeliveryRate)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (r *sesEventRepo) GetHourlyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.HourlyMetrics, error) {
	query := `
		SELECT 
			TO_CHAR(DATE_TRUNC('hour', event_timestamp), 'YYYY-MM-DD HH24:00') as hour,
			COUNT(*) as total_events,
			SUM(CASE WHEN event_type = 'Send' THEN 1 ELSE 0 END) as send_count,
			SUM(CASE WHEN event_type = 'Delivery' THEN 1 ELSE 0 END) as delivery_count,
			SUM(CASE WHEN event_type = 'Bounce' THEN 1 ELSE 0 END) as bounce_count,
			SUM(CASE WHEN event_type = 'Complaint' THEN 1 ELSE 0 END) as complaint_count,
			SUM(CASE WHEN event_type = 'Open' THEN 1 ELSE 0 END) as open_count,
			SUM(CASE WHEN event_type = 'Click' THEN 1 ELSE 0 END) as click_count,
			CASE WHEN COUNT(*) = 0 THEN 0 ELSE (SUM(CASE WHEN event_type = 'Bounce' THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) END as bounce_rate,
			CASE WHEN COUNT(*) = 0 THEN 0 ELSE (SUM(CASE WHEN event_type = 'Delivery' THEN 1 ELSE 0 END) * 100.0 / COUNT(*)) END as delivery_rate
		FROM ses_events
	`
	args := []interface{}{}
	conditions := []string{}
	if start != nil {
		args = append(args, *start)
		conditions = append(conditions, fmt.Sprintf("event_timestamp >= $%d", len(args)))
	}
	if end != nil {
		args = append(args, *end)
		conditions = append(conditions, fmt.Sprintf("event_timestamp < $%d", len(args)))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += `
		GROUP BY DATE_TRUNC('hour', event_timestamp)
		ORDER BY DATE_TRUNC('hour', event_timestamp) DESC
	`
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var metrics []*sesevent.HourlyMetrics
	for rows.Next() {
		m := &sesevent.HourlyMetrics{}
		err := rows.Scan(&m.Hour, &m.TotalEvents, &m.SendCount, &m.DeliveryCount, &m.BounceCount, &m.ComplaintCount, &m.OpenCount, &m.ClickCount, &m.BounceRate, &m.DeliveryRate)
		if err != nil {
			return nil, err
		}
		metrics = append(metrics, m)
	}
	return metrics, nil
}

func (r *sesEventRepo) GetEventTypeCounts(ctx context.Context) (map[string]int, error) {
	query := `
		SELECT event_type, COUNT(DISTINCT message_id)
		FROM ses_events
		GROUP BY event_type
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var eventType string
		var count int
		err := rows.Scan(&eventType, &count)
		if err != nil {
			return nil, err
		}
		counts[eventType] = count
	}
	return counts, nil
}

// DeleteOldEvents menghapus event logs yang lebih lama dari cutoff date
func (r *sesEventRepo) DeleteOldEvents(ctx context.Context, cutoffDate time.Time) (int64, error) {
	query := `DELETE FROM ses_events WHERE event_timestamp < $1`
	result, err := r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
