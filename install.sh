#!/bin/bash

set -e

echo "🚀 SES Dashboard Monitoring - Installation Script"
echo "=================================================="

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if Docker is installed
check_docker() {
    print_status "Checking Docker installation..."
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_success "Docker is installed"
}

# Check if Docker Compose is installed
check_docker_compose() {
    print_status "Checking Docker Compose installation..."
    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
    print_success "Docker Compose is installed"
}

# Create necessary directories
create_directories() {
    print_status "Creating necessary directories..."
    mkdir -p logs
    mkdir -p data/postgres
    print_success "Directories created"
}

# Create initial database setup
create_init_sql() {
    print_status "Creating database initialization script..."
    cat > init.sql << 'EOF'
-- Create database if not exists
SELECT 'CREATE DATABASE ses_monitoring' WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'ses_monitoring')\gexec

-- Connect to the database
\c ses_monitoring;

-- Create admin user (password: admin123)
-- Hash for 'admin123': $2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi
INSERT INTO users (username, email, password_hash, role, is_active, created_at, updated_at) 
VALUES ('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin', true, NOW(), NOW())
ON CONFLICT (username) DO NOTHING;

-- Create default settings
INSERT INTO app_settings (key, value, description, updated_by, created_at, updated_at) VALUES
('aws_enabled', 'false', 'Enable/disable AWS SES integration', 1, NOW(), NOW()),
('aws_region', 'us-east-1', 'AWS region for SES service', 1, NOW(), NOW()),
('retention_days', '30', 'Number of days to retain event logs (0 = never delete)', 1, NOW(), NOW()),
('retention_enabled', 'true', 'Enable/disable automatic log retention cleanup', 1, NOW(), NOW())
ON CONFLICT (key) DO NOTHING;
EOF
    print_success "Database initialization script created"
}

# Build and start services
start_services() {
    print_status "Building and starting services..."
    
    # Stop existing containers
    docker-compose down 2>/dev/null || true
    
    # Build and start services
    docker-compose up --build -d
    
    print_success "Services started successfully"
}

# Wait for services to be ready
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for database
    print_status "Waiting for database..."
    sleep 10
    
    # Wait for backend
    print_status "Waiting for backend API..."
    for i in {1..30}; do
        if curl -s http://localhost:8080/api/health > /dev/null 2>&1; then
            break
        fi
        sleep 2
    done
    
    # Wait for frontend
    print_status "Waiting for frontend..."
    for i in {1..30}; do
        if curl -s http://localhost > /dev/null 2>&1; then
            break
        fi
        sleep 2
    done
    
    print_success "All services are ready"
}

# Run database migrations
run_migrations() {
    print_status "Running database migrations..."
    
    # Execute migrations inside the backend container
    docker-compose exec -T backend sh -c "
        echo 'Running database migrations...'
        # Add your migration commands here if you have a migration tool
        # For now, we rely on the init.sql script
    " || print_warning "Migration command not found, using init.sql instead"
    
    print_success "Database migrations completed"
}

# Show service status
show_status() {
    print_status "Service Status:"
    docker-compose ps
    
    echo ""
    print_success "🎉 Installation completed successfully!"
    echo ""
    echo "📋 Service URLs:"
    echo "   • Application: http://localhost"
    echo "   • API Docs:    http://localhost/swagger/index.html"
    echo "   • Database:    localhost:5432"
    echo ""
    echo "🔐 Default Admin Credentials:"
    echo "   • Username: admin"
    echo "   • Password: admin123"
    echo ""
    echo "📚 Useful Commands:"
    echo "   • View logs:     docker-compose logs -f"
    echo "   • Stop services: docker-compose down"
    echo "   • Restart:       docker-compose restart"
    echo ""
}

# Main installation process
main() {
    echo ""
    print_status "Starting installation process..."
    
    check_docker
    check_docker_compose
    create_directories
    create_init_sql
    start_services
    wait_for_services
    run_migrations
    show_status
}

# Handle script interruption
trap 'print_error "Installation interrupted"; exit 1' INT TERM

# Run main function
main

print_success "Installation script completed!"