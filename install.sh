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

# Build and start services
start_services() {
    print_status "Building and starting services..."
    
    # Stop existing containers
    docker-compose down 2>/dev/null || true
    
    # Build and start services
    docker-compose up --build -d
    
    print_success "Services started successfully"
}

# Run database migrations
run_migrations() {
    print_status "Running database migrations..."
    
    # Wait a bit more for database to be ready
    sleep 5
    
    # Run migrations inside backend container
    docker-compose exec -T backend sh -c "
        cd /app && make migrate-up
    " || print_warning "Migrations may have already been applied"
    
    print_success "Database migrations completed"
}

# Wait for services to be ready
wait_for_services() {
    print_status "Waiting for services to be ready..."
    
    # Wait for database
    print_status "Waiting for database..."
    sleep 15
    
    # Wait for backend
    print_status "Waiting for backend API..."
    for i in {1..30}; do
        if curl -s http://localhost:8080 > /dev/null 2>&1; then
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