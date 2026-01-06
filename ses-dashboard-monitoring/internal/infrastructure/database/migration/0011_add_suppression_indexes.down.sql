-- Remove performance indexes for suppressions table
DROP INDEX IF EXISTS idx_suppressions_updated_at_desc;
DROP INDEX IF EXISTS idx_suppressions_source_updated_at;
DROP INDEX IF EXISTS idx_suppressions_email_gin;
DROP INDEX IF EXISTS idx_suppressions_reason_gin;
DROP INDEX IF EXISTS idx_suppressions_search_composite;

-- Remove indexes for suppression_list table
DROP INDEX IF EXISTS idx_suppression_list_active_updated_at;
DROP INDEX IF EXISTS idx_suppression_list_email_gin;
DROP INDEX IF EXISTS idx_suppression_list_reason_gin;
DROP INDEX IF EXISTS idx_suppression_list_search_composite;