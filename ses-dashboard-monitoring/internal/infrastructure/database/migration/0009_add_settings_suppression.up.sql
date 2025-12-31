-- Create app settings table for UI-configurable settings
CREATE TABLE app_settings (
    key VARCHAR(100) PRIMARY KEY,
    value TEXT,
    description TEXT,
    is_encrypted BOOLEAN DEFAULT FALSE,
    updated_by INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Create suppression list table
CREATE TABLE suppression_list (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    suppression_type VARCHAR(50) NOT NULL, -- 'bounce', 'complaint', 'manual'
    reason TEXT,
    aws_status VARCHAR(50) DEFAULT 'unknown', -- 'unknown', 'suppressed', 'not_suppressed'
    is_active BOOLEAN DEFAULT TRUE,
    added_by INT,
    removed_by INT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    removed_at TIMESTAMP
);

-- Add foreign key constraints after tables exist
ALTER TABLE app_settings ADD CONSTRAINT app_settings_updated_by_fkey FOREIGN KEY (updated_by) REFERENCES users(id);
ALTER TABLE suppression_list ADD CONSTRAINT suppression_list_added_by_fkey FOREIGN KEY (added_by) REFERENCES users(id);
ALTER TABLE suppression_list ADD CONSTRAINT suppression_list_removed_by_fkey FOREIGN KEY (removed_by) REFERENCES users(id);

-- Create indexes for performance
CREATE INDEX idx_suppression_email ON suppression_list(email);
CREATE INDEX idx_suppression_active ON suppression_list(is_active);
CREATE INDEX idx_suppression_type ON suppression_list(suppression_type);

-- Insert default settings (use NULL for updated_by initially)
INSERT INTO app_settings (key, value, description, updated_by) VALUES
('aws_enabled', 'false', 'Enable AWS SES integration features', NULL),
('aws_region', 'us-east-1', 'AWS region for SES', NULL),
('aws_access_key', '', 'AWS access key (encrypted)', NULL),
('aws_secret_key', '', 'AWS secret key (encrypted)', NULL),
('auto_suppress_bounces', 'true', 'Automatically add bounced emails to suppression list', NULL),
('auto_suppress_complaints', 'true', 'Automatically add complained emails to suppression list', NULL);