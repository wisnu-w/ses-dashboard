-- Migration to consolidate suppression tables
-- This migration adds columns to suppressions table to support all features from suppression_list
-- Then migrates data from suppression_list to suppressions
-- Finally drops suppression_list table

-- Step 1: Add new columns to suppressions table to match suppression_list features
ALTER TABLE suppressions 
    ADD COLUMN IF NOT EXISTS id SERIAL,
    ADD COLUMN IF NOT EXISTS suppression_type VARCHAR(50) NOT NULL DEFAULT 'manual',
    ADD COLUMN IF NOT EXISTS aws_status VARCHAR(50) DEFAULT 'unknown',
    ADD COLUMN IF NOT EXISTS is_active BOOLEAN DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS added_by INT,
    ADD COLUMN IF NOT EXISTS removed_by INT,
    ADD COLUMN IF NOT EXISTS removed_at TIMESTAMP,
    ADD COLUMN IF NOT EXISTS synced_at TIMESTAMP;

-- Step 2: Migrate data from suppression_list to suppressions
-- Only migrate active records that don't exist in suppressions
INSERT INTO suppressions (email, reason, source, suppression_type, aws_status, is_active, added_by, created_at, updated_at)
SELECT 
    sl.email,
    COALESCE(sl.reason, 'Migrated from suppression_list'),
    'MANUAL',
    sl.suppression_type,
    sl.aws_status,
    sl.is_active,
    sl.added_by,
    sl.created_at,
    sl.updated_at
FROM suppression_list sl
WHERE sl.is_active = true
    AND NOT EXISTS (SELECT 1 FROM suppressions s WHERE s.email = sl.email);

-- Step 3: Drop suppression_list table and its constraints
DROP TABLE IF EXISTS suppression_list;

-- Step 4: Add indexes for new columns
CREATE INDEX IF NOT EXISTS idx_suppressions_active ON suppressions(is_active);
CREATE INDEX IF NOT EXISTS idx_suppressions_type ON suppressions(suppression_type);
CREATE INDEX IF NOT EXISTS idx_suppressions_email ON suppressions(email);

-- Step 5: Add comment explaining the consolidation
COMMENT ON TABLE suppressions IS 'Consolidated suppression list (merged from suppression_list and suppressions tables)';
