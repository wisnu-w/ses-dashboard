ALTER TABLE ses_events DROP COLUMN IF NOT EXISTS destination_type;
ALTER TABLE ses_message_summaries DROP COLUMN IF NOT EXISTS recipients_dict;
