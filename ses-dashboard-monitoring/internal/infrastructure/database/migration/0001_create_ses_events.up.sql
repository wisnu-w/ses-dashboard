CREATE TABLE IF NOT EXISTS ses_events (
  id BIGSERIAL PRIMARY KEY,
  message_id VARCHAR(100),
  email VARCHAR(255),
  event_type VARCHAR(50),
  status VARCHAR(50),
  reason TEXT,
  created_at TIMESTAMP DEFAULT now(),
  source VARCHAR(255),
  recipients TEXT, -- JSON array
  event_timestamp TIMESTAMP,
  bounce_type VARCHAR(50),
  bounce_sub_type VARCHAR(50),
  diagnostic_code TEXT,
  processing_time_millis INTEGER,
  smtp_response TEXT,
  remote_mta_ip VARCHAR(50),
  reporting_mta VARCHAR(100),
  tags TEXT -- JSON map
);
