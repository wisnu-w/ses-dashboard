# SES Dashboard Monitoring

[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)](https://docker.com)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8?logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-18-61DAFB?logo=react)](https://reactjs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql)](https://postgresql.org)

A comprehensive monitoring dashboard for AWS SES (Simple Email Service) events with real-time analytics, suppression list management, and automated cleanup features.

## 🚀 Features

### 📊 **Dashboard & Analytics**
- Real-time SES event monitoring (Send, Delivery, Bounce, Complaint, Open, Click)
- Interactive charts and metrics visualization
- Daily, monthly, and hourly analytics
- Bounce and delivery rate tracking
- Event filtering and search capabilities

### 🛡️ **Suppression List Management**
- View and manage AWS SES suppression list
- Bulk add/remove email addresses
- Automatic sync with AWS SES
- Manual suppression with custom reasons
- Real-time AWS status checking

### ⚙️ **Administration**
- User management with role-based access control
- AWS SES configuration management
- Data retention policy settings
- Automated log cleanup
- System settings management

### 🔧 **Technical Features**
- RESTful API with Swagger documentation
- JWT-based authentication
- PostgreSQL database with migrations
- Docker containerization
- Express.js proxy server
- Background sync services
- SNS webhook integration

## 📋 Prerequisites

- **Docker** (version 20.10+)
- **Docker Compose** (version 2.0+)
- **Git**
- **Curl** (for health checks)

## 🛠️ Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd ses-dashboard-monitoring
```

### 2. Run Installation Script
```bash
chmod +x install.sh
./install.sh
```

The installation script will:
- ✅ Check Docker and Docker Compose installation
- ✅ Create necessary directories and configuration files
- ✅ Build and start all services
- ✅ Run database migrations
- ✅ Create default admin user
- ✅ Display service URLs and credentials

### 3. Access the Application

| Service | URL | Description |
|---------|-----|-------------|
| **Application** | http://localhost | Complete SES Dashboard |
| **API Docs** | http://localhost/swagger/index.html | Swagger documentation |
| **Database** | localhost:5432 | PostgreSQL (admin access) |

### 4. Default Credentials
```
Username: admin
Password: admin123
```

**That's it!** The `install.sh` script handles everything - from building Docker images to running database migrations. No manual steps required.

## 🏗️ Architecture

```
┌─────────────────┐    ┌─────────────────┐
│   Frontend      │    │   Backend       │
│   (Node.js +    │◄──►│   (Go/Gin)      │
│    Express)     │    │   Port: 8080    │
│   Port: 80      │    └─────────────────┘
└─────────────────┘              │
          │                      │
          │        ┌─────────────────┐
          └────────┤   PostgreSQL    │
                   │   Port: 5432    │
                   └─────────────────┘
```

## 📁 Project Structure

```
ses-dashboard-monitoring/
├── ses-dashboard-monitoring/          # Backend (Go)
│   ├── cmd/api/                      # Application entry point
│   ├── internal/                     # Internal packages
│   │   ├── config/                   # Configuration
│   │   ├── delivery/http/            # HTTP handlers
│   │   ├── domain/                   # Business logic
│   │   ├── infrastructure/           # Database, AWS, etc.
│   │   ├── services/                 # Background services
│   │   └── usecase/                  # Use cases
│   ├── config/                       # Configuration files
│   └── Dockerfile                    # Backend Docker image
├── ses-dashboard-frontend/           # Frontend (Node.js + React)
│   ├── src/                         # React source code
│   │   ├── components/              # React components
│   │   ├── pages/                   # Page components
│   │   ├── services/                # API services
│   │   └── types/                   # TypeScript types
│   ├── server.js                    # Express server with proxy
│   └── Dockerfile                   # Frontend Docker image
├── docker-compose.yml               # Docker Compose configuration
├── install.sh                       # Installation script
└── README.md                        # This file
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_HOST` | Database host | `postgres` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `ses_user` |
| `DB_PASSWORD` | Database password | `ses_password` |
| `DB_NAME` | Database name | `ses_monitoring` |
| `JWT_SECRET` | JWT signing secret | `your-super-secret-jwt-key` |
| `PORT` | Backend server port | `8080` |

### AWS SES Configuration

Configure AWS SES settings through the admin panel:

1. Navigate to **Admin → Settings**
2. Configure AWS credentials and region
3. Test the connection
4. Enable AWS integration

### SNS Webhook Setup

To receive SES events via SNS:

1. Create an SNS topic in AWS
2. Subscribe your endpoint: `http://your-domain/sns/ses`
3. Configure SES to publish events to the SNS topic

## 🔧 Development

> **Note:** For production deployment, simply use `./install.sh`. The sections below are for development purposes only.

### Local Development Setup

1. **Backend Development:**
```bash
cd ses-dashboard-monitoring
go mod download
go run cmd/api/main.go
```

2. **Frontend Development:**
```bash
cd ses-dashboard-frontend
npm install
npm run dev
```

3. **Database Setup:**
```bash
# Start PostgreSQL only
docker-compose up postgres -d

# Run migrations
# Add your migration commands here
```

### Building Docker Images

> **Note:** This is handled automatically by `install.sh`. Manual build only needed for development.

```bash
# Build all services
docker-compose build

# Build specific service
docker-compose build backend
docker-compose build frontend
```

## 📊 API Documentation

The API documentation is available via Swagger UI at:
```
http://localhost/swagger/index.html
```

### Key API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/login` | User authentication |
| `GET` | `/api/events` | Get SES events |
| `GET` | `/api/metrics` | Get dashboard metrics |
| `GET` | `/api/suppression` | Get suppression list |
| `POST` | `/api/suppression/bulk` | Bulk add suppressions |
| `DELETE` | `/api/suppression/bulk` | Bulk remove suppressions |
| `GET` | `/api/settings/aws` | Get AWS settings |
| `PUT` | `/api/settings/retention` | Update retention settings |

## 🛠️ Management Commands

### Docker Compose Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Restart specific service
docker-compose restart backend

# View service status
docker-compose ps

# Execute command in container
docker-compose exec backend sh
```

### Database Management

```bash
# Connect to database
docker-compose exec postgres psql -U ses_user -d ses_monitoring

# Backup database
docker-compose exec postgres pg_dump -U ses_user ses_monitoring > backup.sql

# Restore database
docker-compose exec -T postgres psql -U ses_user ses_monitoring < backup.sql
```

## 🔒 Security Considerations

### Production Deployment

1. **Change Default Credentials:**
   - Update admin password after first login
   - Use strong JWT secret key

2. **Environment Variables:**
   - Store sensitive data in environment variables
   - Use Docker secrets for production

3. **Network Security:**
   - Use HTTPS in production
   - Configure firewall rules
   - Limit database access

4. **AWS Security:**
   - Use IAM roles instead of access keys when possible
   - Limit SES permissions to minimum required
   - Enable AWS CloudTrail for audit logging

## 📈 Monitoring & Logging

### Application Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### Health Checks

The application includes health check endpoints:

- Application: `http://localhost`
- API Health: `http://localhost/api/health`

### Metrics

Monitor application metrics through:
- Dashboard analytics page
- Database query performance
- Container resource usage

## 🚨 Troubleshooting

### Common Issues

1. **Port Already in Use:**
```bash
# Check what's using the port
lsof -i :80
lsof -i :8080

# Stop conflicting services
sudo systemctl stop apache2
sudo systemctl stop nginx
```

2. **Database Connection Issues:**
```bash
# Check database logs
docker-compose logs postgres

# Verify database is running
docker-compose ps postgres

# Test database connection
docker-compose exec postgres psql -U ses_user -d ses_monitoring -c "SELECT 1;"
```

3. **Frontend Build Issues:**
```bash
# Clear npm cache
npm cache clean --force

# Rebuild frontend
docker-compose build --no-cache frontend
```

4. **Backend API Issues:**
```bash
# Check backend logs
docker-compose logs backend

# Verify configuration
docker-compose exec backend env | grep DB_
```

### Performance Optimization

1. **Database Optimization:**
   - Configure retention policies
   - Add database indexes for large datasets
   - Monitor query performance

2. **Container Resources:**
   - Adjust memory limits in docker-compose.yml
   - Monitor container resource usage
   - Scale services horizontally if needed

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go best practices for backend
- Use TypeScript for frontend development
- Write tests for new features
- Update documentation for API changes
- Follow conventional commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- AWS SES for email service integration
- PostgreSQL for reliable data storage
- React and Go communities for excellent frameworks
- Docker for containerization platform

## 📞 Support

For support and questions:

1. Check the [Issues](../../issues) page
2. Review the troubleshooting section
3. Create a new issue with detailed information

---

**Made with ❤️ by Wisnu**