CREATE MATERIALIZED VIEW IF NOT EXISTS mv_ses_daily_summary AS
SELECT 
    DATE(event_timestamp) AS event_date,
    COUNT(DISTINCT message_id) AS total_events,
    COUNT(DISTINCT message_id) FILTER (WHERE event_type = 'Send') AS send_count,
    COUNT(DISTINCT message_id) FILTER (WHERE event_type = 'Delivery') AS delivery_count,
    COUNT(DISTINCT message_id) FILTER (WHERE event_type = 'Bounce') AS bounce_count,
    COUNT(DISTINCT message_id) FILTER (WHERE event_type = 'Complaint') AS complaint_count,
    COUNT(DISTINCT message_id) FILTER (WHERE event_type = 'Open') AS open_count,
    COUNT(DISTINCT message_id) FILTER (WHERE event_type = 'Click') AS click_count
FROM ses_events
GROUP BY DATE(event_timestamp);

CREATE UNIQUE INDEX IF NOT EXISTS idx_mv_ses_daily_date ON mv_ses_daily_summary (event_date);
