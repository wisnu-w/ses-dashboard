CREATE TABLE IF NOT EXISTS suppressions (
    email VARCHAR(255) PRIMARY KEY,
    reason TEXT NOT NULL,
    source VARCHAR(50) NOT NULL DEFAULT 'AWS',
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_suppressions_source ON suppressions(source);
CREATE INDEX IF NOT EXISTS idx_suppressions_updated_at ON suppressions(updated_at);