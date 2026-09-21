CREATE INDEX IF NOT EXISTS idx_ses_message_summaries_search_composite 
ON ses_message_summaries USING gin((message_id || ' ' || email || ' ' || subject || ' ' || source) gin_trgm_ops);
