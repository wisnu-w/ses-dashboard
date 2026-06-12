CREATE TABLE IF NOT EXISTS ses_message_summaries (
    message_id VARCHAR(100) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    subject TEXT DEFAULT '',
    source VARCHAR(255) DEFAULT '',
    latest_event VARCHAR(50) NOT NULL,
    latest_status VARCHAR(50) NOT NULL,
    status_priority INTEGER NOT NULL DEFAULT 0,
    first_event_at TIMESTAMP NOT NULL,
    last_event_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_ses_message_summaries_last_event_at ON ses_message_summaries(last_event_at DESC);
CREATE INDEX IF NOT EXISTS idx_ses_message_summaries_status_priority ON ses_message_summaries(status_priority DESC);
CREATE INDEX IF NOT EXISTS idx_ses_message_summaries_email ON ses_message_summaries(email);
CREATE INDEX IF NOT EXISTS idx_ses_message_summaries_source ON ses_message_summaries(source);

INSERT INTO ses_message_summaries (
    message_id, email, subject, source, latest_event, latest_status, status_priority,
    first_event_at, last_event_at, created_at, updated_at
)
SELECT DISTINCT ON (message_id)
    se.message_id,
    COALESCE(se.email, ''),
    COALESCE(se.subject, ''),
    COALESCE(se.source, ''),
    COALESCE(se.event_type, ''),
    summary.latest_status,
    summary.status_priority,
    summary.first_event_at,
    summary.last_event_at,
    NOW(),
    NOW()
FROM ses_events se
JOIN (
    SELECT
        message_id,
        MIN(event_timestamp) AS first_event_at,
        MAX(event_timestamp) AS last_event_at,
        CASE
            WHEN BOOL_OR(event_type = 'Complaint') THEN 'Complaint'
            WHEN BOOL_OR(event_type = 'Bounce') THEN 'Bounce'
            WHEN BOOL_OR(event_type = 'Delivery') THEN 'Delivery'
            WHEN BOOL_OR(event_type = 'Send') THEN 'Pending'
            ELSE COALESCE(MAX(event_type), 'Unknown')
        END AS latest_status,
        CASE
            WHEN BOOL_OR(event_type = 'Complaint') THEN 40
            WHEN BOOL_OR(event_type = 'Bounce') THEN 30
            WHEN BOOL_OR(event_type = 'Delivery') THEN 20
            WHEN BOOL_OR(event_type = 'Send') THEN 10
            ELSE 0
        END AS status_priority
    FROM ses_events
    WHERE message_id IS NOT NULL AND message_id <> ''
    GROUP BY message_id
) summary ON summary.message_id = se.message_id
WHERE se.message_id IS NOT NULL AND se.message_id <> ''
ORDER BY se.message_id, se.event_timestamp DESC
ON CONFLICT (message_id) DO NOTHING;
