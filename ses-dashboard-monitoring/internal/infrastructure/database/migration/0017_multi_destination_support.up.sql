ALTER TABLE ses_events ADD COLUMN IF NOT EXISTS destination_type VARCHAR(10) DEFAULT 'To';
ALTER TABLE ses_message_summaries ADD COLUMN IF NOT EXISTS recipients_dict JSONB DEFAULT '{}'::jsonb;

UPDATE ses_message_summaries
SET recipients_dict = jsonb_build_object(
    email, jsonb_build_object(
        'email', email,
        'status', latest_status,
        'diagnostic_code', '-',
        'type', 'To',
        'status_priority', status_priority
    )
)
WHERE recipients_dict = '{}'::jsonb OR recipients_dict IS NULL;
