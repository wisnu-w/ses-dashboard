-- Seed admin user with proper hashed password
INSERT INTO users (username, password, email, role, active) 
VALUES ('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin@example.com', 'admin', true)
ON CONFLICT (username) DO UPDATE SET 
  password = EXCLUDED.password,
  email = EXCLUDED.email,
  role = EXCLUDED.role,
  active = EXCLUDED.active;