package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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
			processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags, destination_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
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
		e.DestinationType,
	)
	return err
}

func (r *sesEventRepo) UpsertMessageSummary(ctx context.Context, e *sesevent.Event) error {
	query := `
		INSERT INTO ses_message_summaries (
			message_id, email, subject, source, latest_event, latest_status, status_priority,
			first_event_at, last_event_at, recipients_dict, created_at, updated_at
		)
		VALUES (
			$1, $2, $3, $4, $5::text,
			CASE
				WHEN $5::text = 'Complaint' THEN 'Complaint'
				WHEN $5::text = 'Bounce' THEN 'Bounce'
				WHEN $5::text = 'Delivery' THEN 'Delivery'
				WHEN $5::text = 'Send' THEN 'Pending'
				ELSE COALESCE($5::text, 'Unknown')
			END,
			CASE
				WHEN $5::text = 'Complaint' THEN 40
				WHEN $5::text = 'Bounce' THEN 30
				WHEN $5::text = 'Delivery' THEN 20
				WHEN $5::text = 'Send' THEN 10
				ELSE 0
			END,
			$6::timestamp, $6::timestamp,
			jsonb_build_object($2::text, jsonb_build_object(
				'email', $2::text,
				'status', CASE
					WHEN $5::text = 'Complaint' THEN 'Complaint'
					WHEN $5::text = 'Bounce' THEN 'Bounce'
					WHEN $5::text = 'Delivery' THEN 'Delivery'
					WHEN $5::text = 'Send' THEN 'Pending'
					ELSE COALESCE($5::text, 'Unknown')
				END,
				'diagnostic_code', COALESCE(NULLIF($7::text, ''), '-'),
				'type', COALESCE(NULLIF($8::text, ''), 'To'),
				'status_priority', CASE
					WHEN $5::text = 'Complaint' THEN 40
					WHEN $5::text = 'Bounce' THEN 30
					WHEN $5::text = 'Delivery' THEN 20
					WHEN $5::text = 'Send' THEN 10
					ELSE 0
				END
			)),
			NOW(), NOW()
		)
		ON CONFLICT (message_id) DO UPDATE SET
			email = CASE WHEN EXCLUDED.last_event_at >= ses_message_summaries.last_event_at THEN EXCLUDED.email ELSE ses_message_summaries.email END,
			subject = CASE WHEN EXCLUDED.last_event_at >= ses_message_summaries.last_event_at THEN EXCLUDED.subject ELSE ses_message_summaries.subject END,
			source = CASE WHEN EXCLUDED.last_event_at >= ses_message_summaries.last_event_at THEN EXCLUDED.source ELSE ses_message_summaries.source END,
			latest_event = CASE WHEN EXCLUDED.last_event_at >= ses_message_summaries.last_event_at THEN EXCLUDED.latest_event ELSE ses_message_summaries.latest_event END,
			latest_status = CASE WHEN EXCLUDED.status_priority >= ses_message_summaries.status_priority THEN EXCLUDED.latest_status ELSE ses_message_summaries.latest_status END,
			status_priority = GREATEST(ses_message_summaries.status_priority, EXCLUDED.status_priority),
			first_event_at = LEAST(ses_message_summaries.first_event_at, EXCLUDED.first_event_at),
			last_event_at = GREATEST(ses_message_summaries.last_event_at, EXCLUDED.last_event_at),
			recipients_dict = ses_message_summaries.recipients_dict || jsonb_build_object(
				$2::text,
				COALESCE(ses_message_summaries.recipients_dict->($2::text), '{}'::jsonb) || jsonb_build_object(
					'email', $2::text,
					'status', CASE
						WHEN EXCLUDED.status_priority >= COALESCE((ses_message_summaries.recipients_dict->($2::text)->>'status_priority')::int, 0)
						THEN EXCLUDED.latest_status
						ELSE COALESCE(ses_message_summaries.recipients_dict->($2::text)->>'status', EXCLUDED.latest_status)
					END,
					'diagnostic_code', CASE
						WHEN $7::text != '' THEN $7::text
						ELSE COALESCE(ses_message_summaries.recipients_dict->($2::text)->>'diagnostic_code', '-')
					END,
					'type', CASE
						WHEN $8::text != '' THEN $8::text
						ELSE COALESCE(ses_message_summaries.recipients_dict->($2::text)->>'type', 'To')
					END,
					'status_priority', GREATEST(EXCLUDED.status_priority, COALESCE((ses_message_summaries.recipients_dict->($2::text)->>'status_priority')::int, 0))
				)
			),
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, e.MessageID, e.Email, e.Subject, e.Source, e.EventType, e.EventTimestamp, e.DiagnosticCode, e.DestinationType)
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

func (r *sesEventRepo) GetEventGroupsPaginated(ctx context.Context, limit, offset int) ([]*sesevent.MessageGroup, error) {
	return r.getEventGroups(ctx, limit, offset, "", "", "")
}

func (r *sesEventRepo) GetEventGroupsWithFilter(ctx context.Context, limit, offset int, search, startDate, endDate string) ([]*sesevent.MessageGroup, error) {
	return r.getEventGroups(ctx, limit, offset, search, startDate, endDate)
}

func (r *sesEventRepo) getEventGroups(ctx context.Context, limit, offset int, search, startDate, endDate string) ([]*sesevent.MessageGroup, error) {
	// Use deferred join (late row lookup) to prevent PostgreSQL from doing a sequential scan + sort
	// on the bloated table, which contains the wide JSONB column recipients_dict.
	innerQuery := `
		SELECT message_id
		FROM ses_message_summaries
		WHERE 1=1
	`
	args := []interface{}{}
	argIndex := 0

	if search != "" {
		argIndex++
		innerQuery += fmt.Sprintf(" AND (message_id || ' ' || email || ' ' || subject || ' ' || source) ILIKE $%d", argIndex)
		args = append(args, "%"+search+"%")
	}

	if startDate != "" {
		argIndex++
		innerQuery += fmt.Sprintf(" AND last_event_at >= $%d", argIndex)
		args = append(args, startDate)
	}

	if endDate != "" {
		argIndex++
		innerQuery += fmt.Sprintf(" AND last_event_at <= $%d", argIndex)
		args = append(args, endDate+" 23:59:59")
	}

	innerQuery += " ORDER BY last_event_at DESC"
	innerQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex+1, argIndex+2)
	args = append(args, limit, offset)

	query := fmt.Sprintf(`
		SELECT s.message_id, s.email, s.subject, s.source, s.latest_status, s.latest_event, s.first_event_at, s.last_event_at, COALESCE(s.recipients_dict, '{}'::jsonb)
		FROM (%s) AS page
		JOIN ses_message_summaries s ON s.message_id = page.message_id
		ORDER BY s.last_event_at DESC
	`, innerQuery)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	groups := []*sesevent.MessageGroup{}
	for rows.Next() {
		group := &sesevent.MessageGroup{}
		var recipientsJSON string
		err := rows.Scan(
			&group.MessageID,
			&group.Email,
			&group.Subject,
			&group.Source,
			&group.LatestStatus,
			&group.LatestEvent,
			&group.FirstEventAt,
			&group.LastEventAt,
			&recipientsJSON,
		)
		if err != nil {
			return nil, err
		}
		
		var recDict map[string]sesevent.RecipientDetail
		if err := json.Unmarshal([]byte(recipientsJSON), &recDict); err == nil {
			for _, detail := range recDict {
				group.RecipientsDetail = append(group.RecipientsDetail, detail)
			}
		}

		groups = append(groups, group)
	}
	return groups, nil
}

func (r *sesEventRepo) GetEventGroupCount(ctx context.Context, search, startDate, endDate string) (int, error) {
	if search == "" && startDate == "" && endDate == "" {
		// Use table statistics for fast unfiltered count estimation
		estimateQuery := `SELECT reltuples::bigint FROM pg_class WHERE relname = 'ses_message_summaries'`
		var count int
		err := r.db.QueryRowContext(ctx, estimateQuery).Scan(&count)
		if err == nil && count > 1000 { // Use estimate if it's reasonably large, avoiding slow COUNT(*)
			return count, nil
		}
	}

	query := `SELECT COUNT(*) FROM ses_message_summaries WHERE 1=1`
	args := []interface{}{}
	argIndex := 0

	if search != "" {
		argIndex++
		query += fmt.Sprintf(" AND (message_id || ' ' || email || ' ' || subject || ' ' || source) ILIKE $%d", argIndex)
		args = append(args, "%"+search+"%")
	}

	if startDate != "" {
		argIndex++
		query += fmt.Sprintf(" AND last_event_at >= $%d", argIndex)
		args = append(args, startDate)
	}

	if endDate != "" {
		argIndex++
		query += fmt.Sprintf(" AND last_event_at <= $%d", argIndex)
		args = append(args, endDate+" 23:59:59")
	}

	var count int
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

func (r *sesEventRepo) GetEventsByMessageID(ctx context.Context, messageID string) ([]*sesevent.Event, error) {
	query := `
		SELECT message_id, email, subject, event_type, status, reason, source, recipients,
			   event_timestamp, bounce_type, bounce_sub_type, diagnostic_code,
			   processing_time_millis, smtp_response, remote_mta_ip, reporting_mta, tags
		FROM ses_events
		WHERE message_id = $1
		ORDER BY event_timestamp ASC
	`
	rows, err := r.db.QueryContext(ctx, query, messageID)
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
	query := `SELECT COALESCE(SUM(total_events), 0) FROM mv_ses_daily_summary`
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
		WITH stats AS (
			SELECT COALESCE(SUM(total_events), 0) as total, COALESCE(SUM(bounce_count), 0) as bounces
			FROM mv_ses_daily_summary
		)
		SELECT CASE WHEN total = 0 THEN 0 ELSE (bounces * 100.0 / total) END
		FROM stats
	`
	var rate float64
	err := r.db.QueryRowContext(ctx, query).Scan(&rate)
	return rate, err
}

func (r *sesEventRepo) GetDeliveryRate(ctx context.Context) (float64, error) {
	query := `
		WITH stats AS (
			SELECT COALESCE(SUM(total_events), 0) as total, COALESCE(SUM(delivery_count), 0) as deliveries
			FROM mv_ses_daily_summary
		)
		SELECT CASE WHEN total = 0 THEN 0 ELSE (deliveries * 100.0 / total) END
		FROM stats
	`
	var rate float64
	err := r.db.QueryRowContext(ctx, query).Scan(&rate)
	return rate, err
}

func (r *sesEventRepo) GetDailyMetrics(ctx context.Context, start, end *time.Time) ([]*sesevent.DailyMetrics, error) {
	query := `
		SELECT 
			event_date as date,
			total_events,
			send_count,
			delivery_count,
			bounce_count,
			complaint_count,
			open_count,
			click_count,
			CASE WHEN total_events = 0 THEN 0 ELSE (bounce_count * 100.0 / total_events) END as bounce_rate,
			CASE WHEN total_events = 0 THEN 0 ELSE (delivery_count * 100.0 / total_events) END as delivery_rate
		FROM mv_ses_daily_summary
	`
	args := []interface{}{}
	conditions := []string{}
	if start != nil {
		args = append(args, *start)
		conditions = append(conditions, fmt.Sprintf("event_date >= DATE($%d)", len(args)))
	}
	if end != nil {
		args = append(args, *end)
		conditions = append(conditions, fmt.Sprintf("event_date < DATE($%d)", len(args)))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += `
		ORDER BY event_date DESC
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
			TO_CHAR(DATE_TRUNC('month', event_date), 'YYYY-MM') as month,
			COALESCE(SUM(total_events), 0) as total_events,
			COALESCE(SUM(send_count), 0) as send_count,
			COALESCE(SUM(delivery_count), 0) as delivery_count,
			COALESCE(SUM(bounce_count), 0) as bounce_count,
			COALESCE(SUM(complaint_count), 0) as complaint_count,
			COALESCE(SUM(open_count), 0) as open_count,
			COALESCE(SUM(click_count), 0) as click_count,
			CASE WHEN COALESCE(SUM(total_events), 0) = 0 THEN 0 ELSE (COALESCE(SUM(bounce_count), 0) * 100.0 / COALESCE(SUM(total_events), 1)) END as bounce_rate,
			CASE WHEN COALESCE(SUM(total_events), 0) = 0 THEN 0 ELSE (COALESCE(SUM(delivery_count), 0) * 100.0 / COALESCE(SUM(total_events), 1)) END as delivery_rate
		FROM mv_ses_daily_summary
	`
	args := []interface{}{}
	conditions := []string{}
	if start != nil {
		args = append(args, *start)
		conditions = append(conditions, fmt.Sprintf("event_date >= DATE($%d)", len(args)))
	}
	if end != nil {
		args = append(args, *end)
		conditions = append(conditions, fmt.Sprintf("event_date < DATE($%d)", len(args)))
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += `
		GROUP BY DATE_TRUNC('month', event_date)
		ORDER BY DATE_TRUNC('month', event_date) DESC
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
		SELECT 
			COALESCE(SUM(send_count), 0),
			COALESCE(SUM(delivery_count), 0),
			COALESCE(SUM(bounce_count), 0),
			COALESCE(SUM(complaint_count), 0),
			COALESCE(SUM(open_count), 0),
			COALESCE(SUM(click_count), 0)
		FROM mv_ses_daily_summary
	`
	var send, delivery, bounce, complaint, open, click int
	err := r.db.QueryRowContext(ctx, query).Scan(&send, &delivery, &bounce, &complaint, &open, &click)
	if err != nil {
		return nil, err
	}

	counts := map[string]int{
		"Send":      send,
		"Delivery":  delivery,
		"Bounce":    bounce,
		"Complaint": complaint,
		"Open":      open,
		"Click":     click,
	}
	return counts, nil
}

// DeleteOldEvents menghapus event logs yang lebih lama dari cutoff date
func (r *sesEventRepo) RefreshDailySummary(ctx context.Context) error {
	query := `REFRESH MATERIALIZED VIEW CONCURRENTLY mv_ses_daily_summary`
	_, err := r.db.ExecContext(ctx, query)
	return err
}

func (r *sesEventRepo) DeleteOldEvents(ctx context.Context, cutoffDate time.Time) (int64, error) {
	query := `DELETE FROM ses_events WHERE event_timestamp < $1`
	result, err := r.db.ExecContext(ctx, query, cutoffDate)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
