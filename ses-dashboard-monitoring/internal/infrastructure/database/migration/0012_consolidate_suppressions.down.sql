-- Rollback: Restore suppression_list table (data will be lost if suppressions was populated after migration)
-- This is a best-effort rollback - manual records added after migration cannot be fully restored

-- Step 1: Recreate suppression_list table
CREATE TABLE suppression_list (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    suppression_type VARCHAR(50) NOT NULL DEFAULT 'manual',
    reason TEXT,
    aws_status VARCHAR(50) DEFAULT 'unknown',
    is_active BOOLEAN DEFAULT TRUE,
    added_by INT,
    removed_by INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    removed_at TIMESTAMP
);

-- Step 2: Migrate data back from suppressions to suppression_list
INSERT INTO suppression_list (email, suppression_type, reason, aws_status, is_active, added_by, created_at, updated_at)
SELECT 
    email,
    suppression_type,
    reason,
    aws_status,
    is_active,
    added_by,
    created_at,
    updated_at
FROM suppressions
WHERE is_active = true;

-- Step 3: Remove added columns from suppressions (keep basic structure)
ALTER TABLE suppressions 
    DROP COLUMN IF EXISTS id,
    DROP COLUMN IF EXISTS suppression_type,
    DROP COLUMN IF EXISTS aws_status,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS added_by,
    DROP COLUMN IF EXISTS removed_by,
    DROP COLUMN IF EXISTS removed_at,
    DROP COLUMN IF EXISTS synced_at;

-- Step 4: Recreate indexes
CREATE INDEX idx_suppression_email ON suppression_list(email);
CREATE INDEX idx_suppression_active ON suppression_list(is_active);
CREATE INDEX idx_suppression_type ON suppression_list(suppression_type);
