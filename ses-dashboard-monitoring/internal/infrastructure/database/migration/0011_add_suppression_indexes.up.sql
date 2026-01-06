-- Add performance indexes for suppressions table
CREATE INDEX IF NOT EXISTS idx_suppressions_updated_at_desc ON suppressions(updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_suppressions_source_updated_at ON suppressions(source, updated_at DESC);
CREATE INDEX IF NOT EXISTS idx_suppressions_email_gin ON suppressions USING gin(email gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_suppressions_reason_gin ON suppressions USING gin(reason gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_suppressions_search_composite ON suppressions USING gin((email || ' ' || reason || ' ' || source) gin_trgm_ops);

-- Add indexes for suppression_list table if it exists
CREATE INDEX IF NOT EXISTS idx_suppression_list_active_updated_at ON suppression_list(is_active, updated_at DESC) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_suppression_list_email_gin ON suppression_list USING gin(email gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_suppression_list_reason_gin ON suppression_list USING gin(reason gin_trgm_ops);
CREATE INDEX IF NOT EXISTS idx_suppression_list_search_composite ON suppression_list USING gin((email || ' ' || reason) gin_trgm_ops) WHERE is_active = true;