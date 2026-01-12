# SES Dashboard Monitoring

[![Docker](https://img.shields.io/badge/Docker-Ready-blue?logo=docker)](https://docker.com)
[![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go)](https://golang.org)
[![React](https://img.shields.io/badge/React-19-61DAFB?logo=react)](https://reactjs.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?logo=postgresql)](https://postgresql.org)

A comprehensive monitoring dashboard for AWS SES (Simple Email Service) events with real-time analytics, suppression list management, and automated cleanup features.

## 🚀 Features

### 📊 **Dashboard & Analytics**
- Real-time SES event monitoring (Send, Delivery, Bounce, Complaint, Open, Click)
- Interactive charts and metrics visualization using Recharts
- Daily, monthly, and hourly analytics
- Bounce and delivery rate tracking
- Event filtering and search capabilities
- Responsive design with Tailwind CSS

### 🛡️ **Suppression List Management**
- View and manage AWS SES suppression list
- **Bulk operations**: Add/remove multiple emails at once
- Automatic sync with AWS SES suppression list
- Manual suppression with custom reasons
- Real-time AWS status checking
- Email selection with checkboxes

### ⚙️ **Administration**
- **User Management**: Role-based access control (Admin/User)
- **AWS SES Configuration**: Manage credentials and settings via UI
- **Data Retention**: Configurable log cleanup policies
- **Automated Cleanup**: Background service for old data removal
- **Settings Management**: Persistent configuration storage

### 🔧 **Technical Features**
- **RESTful API** with comprehensive Swagger documentation
- **JWT Authentication** with role-based authorization
- **PostgreSQL Database** with optimized indexes
- **Docker Containerization** for easy deployment
- **Express.js Proxy Server** for API routing
- **Background Services**: Sync and cleanup automation
- **SNS Webhook Integration** for real-time event processing
- **CORS Support** for cross-origin requests

## 📋 Prerequisites

- **Docker** (version 20.10+) or **Podman** (version 4.0+)
- **Docker Compose** (version 2.0+) or **Podman Compose**
- **Go** (version 1.25+ for local development)
- **Make** (for running build commands)
- **Git**
- **Curl** (for health checks)

## 🛠️ Quick Start

### 1. Clone the Repository
```bash
git clone <repository-url>
cd ses-dashboard
```

### 2. Prepare Environment (Optional)
The stack reads config from `ses-dashboard-monitoring/config/config.yaml` and Docker Compose reads `.env` if present.

If you want to generate a `.env` from `config.yaml`:
```bash
chmod +x generate-env.sh
./generate-env.sh
```

Optional: edit `.env` to override values for Docker Compose.

### 3. Run Installation Script
```bash
chmod +x install.sh
./install.sh
```

The installation script will:
- ✅ Check Docker and Docker Compose installation
- ✅ Create necessary directories
- ✅ Build and start all services
- ✅ Run database migrations automatically
- ✅ Display service URLs and credentials

### 4. Access the Application

| Service | URL | Description |
|---------|-----|-------------|
| **Application** | http://localhost | Complete SES Dashboard |
| **API Documentation** | http://localhost/swagger/index.html | Swagger UI (proxied) |
| **API Documentation (direct)** | http://localhost:8080/swagger/index.html | Swagger UI (backend) |
| **Database** | localhost:5432 | PostgreSQL (admin access) |
| **SNS SES Webhook** | http://localhost/sns/ses | SES events webhook (POST from SNS) |

### 5. Default Credentials
```
Username: admin
Password: password
```

**⚠️ IMPORTANT: Change the default password after first login for security!**

**That's it!** The `install.sh` script handles everything - from building Docker images to database migrations. No manual steps required.

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

### Components

- **Frontend**: React 19 + TypeScript + Tailwind CSS served by Express.js
- **Backend**: Go 1.25 + Gin framework with clean architecture
- **Database**: PostgreSQL 15 with optimized schema and indexes
- **Proxy**: Express.js handles API routing and static file serving

## 📁 Project Structure

```
ses-dashboard/
├── ses-dashboard-monitoring/          # Backend (Go)
│   ├── cmd/api/                      # Application entry point
│   │   └── main.go                   # Main application file
│   ├── cmd/migrate/                  # Migration tool
│   │   └── main.go                   # Database migration utility
│   ├── internal/                     # Internal packages
│   │   ├── config/                   # Configuration management
│   │   ├── delivery/http/            # HTTP handlers and middleware
│   │   ├── domain/                   # Business logic and entities
│   │   │   ├── sesevent/             # SES event domain
│   │   │   ├── settings/             # Settings domain
│   │   │   ├── suppression/          # Suppression domain
│   │   │   └── user/                 # User domain
│   │   ├── infrastructure/           # External dependencies
│   │   │   ├── aws/                  # AWS SES client
│   │   │   ├── database/             # Database connection & migrations
│   │   │   └── repository/           # Data access layer
│   │   ├── services/                 # Background services
│   │   │   ├── cleanup_service.go    # Data cleanup automation
│   │   │   └── sync_service.go       # AWS sync automation
│   │   └── usecase/                  # Business use cases
│   ├── config/                       # Configuration files
│   │   └── config.yaml               # Application configuration
│   ├── docs/                         # Swagger documentation
│   ├── Makefile                      # Build and development commands
│   └── Dockerfile                    # Backend Docker image
├── ses-dashboard-frontend/           # Frontend (React + Node.js)
│   ├── src/                         # React source code
│   │   ├── components/              # Reusable React components
│   │   │   ├── Layout.tsx           # Main layout component
│   │   │   ├── Charts.tsx           # Chart components
│   │   │   └── EventsTable.tsx      # Events table component
│   │   ├── pages/                   # Page components
│   │   │   ├── DashboardPage.tsx    # Main dashboard
│   │   │   ├── EventsPage.tsx       # Events listing
│   │   │   ├── AnalyticsPage.tsx    # Analytics charts
│   │   │   ├── SuppressionPage.tsx  # Suppression management
│   │   │   ├── UsersPage.tsx        # User management
│   │   │   └── SettingsPage.tsx     # System settings
│   │   ├── services/                # API services
│   │   │   └── api.ts               # API client with Axios
│   │   └── types/                   # TypeScript type definitions
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
| `APP_NAME` | App name | `ses-monitoring` |
| `APP_ENV` | Environment name | `local` |
| `APP_PORT` | Backend server port | `8080` |
| `ENABLE_SWAGGER` | Toggle Swagger | `true` |
| `JWT_SECRET` | JWT signing secret | `your-super-secret-jwt-key` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_USER` | Database username | `ses_user` |
| `DB_PASSWORD` | Database password | `password123!` |
| `DB_NAME` | Database name | `ses_dashboard` |
| `DB_SSLMODE` | Database SSL mode | `disable` |
| `AWS_REGION` | AWS region | `ap-southeast-1` |
| `AWS_ACCESS_KEY` | AWS access key | (empty) |
| `AWS_SECRET_KEY` | AWS secret key | (empty) |
| `SNS_TOPIC_ARN` | Allowlist SNS topic for `/sns/ses` | (empty = allow any) |
| `BACKEND_URL` | Backend URL for frontend proxy | `http://backend:8080` |
| `VITE_API_URL` | Frontend dev API base URL | `http://localhost:8080` |

### Config Files
- `ses-dashboard-monitoring/config/config.yaml` is the primary config file.
- `.env` is used by Docker Compose if present. You can generate it via `./generate-env.sh` or manage it manually.
  - Add `SNS_TOPIC_ARN` manually to `.env` if you want to restrict webhook intake.

Environment variables take precedence over YAML values.

When running via Docker Compose, the backend container overrides `DB_HOST` to `postgres` so it can reach the database service.

### Database Schema

The application uses 4 main tables:

1. **users** - User accounts with role-based access
2. **ses_events** - SES event logs with full event data
3. **app_settings** - Configurable application settings
4. **suppressions** - Email suppression list management

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
4. Confirm the subscription (HTTPS confirmation) so SNS can deliver messages
5. Events will be automatically processed and stored

#### AWS Console Steps (SNS)
1. Open **AWS SNS Console → Topics → Create topic** (type **Standard**).
2. After creating the topic, copy the **Topic ARN**.
3. Go to **Subscriptions → Create subscription**:
   - Protocol: **HTTPS**
   - Endpoint: `https://your-domain/sns/ses` (use your public domain; local `http://localhost` will not work from AWS)
4. Confirm the subscription (AWS sends a confirmation request to your endpoint).
5. In **Amazon SES → Configuration → Event destinations**, add this SNS topic as the destination for SES event types you need.

Optional hardening: set `SNS_TOPIC_ARN` so the API only accepts messages from that topic.

## 🔧 Development

> **Note:** For production deployment, simply use `./install.sh`. The sections below are for development purposes only.

### Database Migrations

The application uses golang-migrate for database schema management:

```bash
# Run migrations
cd ses-dashboard-monitoring
make migrate-up

# Rollback migrations
make migrate-down

# Check migration status
make migrate-version
```

### Local Development Setup

1. **Backend Development:**
```bash
cd ses-dashboard-monitoring
go mod download
export APP_PORT=8080
export JWT_SECRET=your-secret
go run cmd/api/main.go
```

2. **Frontend Development:**
```bash
cd ses-dashboard-frontend
npm install
export VITE_API_URL=http://localhost:8080
npm run dev
```

3. **Database Setup:**
```bash
# Start PostgreSQL only
docker-compose up postgres -d

# Run migrations
cd ses-dashboard-monitoring
make migrate-up
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

#### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/login` | User authentication |
| `PUT` | `/api/change-password` | Change user password |

#### Events & Metrics
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/events` | Get SES events with pagination |
| `GET` | `/api/metrics` | Get dashboard metrics |
| `GET` | `/api/metrics/daily` | Get daily analytics |
| `GET` | `/api/metrics/monthly` | Get monthly analytics |
| `GET` | `/api/metrics/hourly` | Get hourly analytics |

All metrics endpoints accept optional query parameters:
- `start_date=YYYY-MM-DD`
- `end_date=YYYY-MM-DD`

Default ranges:
- Daily: last 30 days
- Hourly: last 48 hours
- Monthly: last 12 months

#### Suppression Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/suppression` | Get suppression list |
| `POST` | `/api/suppression` | Add single email to suppression |
| `POST` | `/api/suppression/bulk` | **Bulk add** multiple emails |
| `DELETE` | `/api/suppression/bulk` | **Bulk remove** multiple emails |
| `DELETE` | `/api/suppression/:email` | Remove single email |
| `GET` | `/api/suppression/:email/status` | Check email AWS status |
| `POST` | `/api/suppression/sync` | Trigger AWS sync |
| `GET` | `/api/suppression/sync/status` | Get sync status |

#### Administration (Admin Only)
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users` | Get all users |
| `POST` | `/api/users` | Create new user |
| `PUT` | `/api/users/:id/reset-password` | Reset user password |
| `PUT` | `/api/users/:id/disable` | Disable user account |
| `DELETE` | `/api/users/:id` | Delete user account |
| `GET` | `/api/settings/aws` | Get AWS settings |
| `PUT` | `/api/settings/aws` | Update AWS settings |
| `POST` | `/api/settings/aws/test` | Test AWS connection |
| `GET` | `/api/settings/retention` | Get retention settings |
| `PUT` | `/api/settings/retention` | Update retention settings |

## 🛠️ Management Commands

### Docker Compose Commands

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend

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

# View database size
docker-compose exec postgres psql -U ses_user -d ses_monitoring -c "\l+"
```

## 🔒 Security Considerations

### Production Deployment

1. **Change Default Credentials:**
   - Update admin password after first login
   - Use strong JWT secret key (change `JWT_SECRET` environment variable)

2. **Environment Variables:**
   - Store sensitive data in environment variables
   - Use Docker secrets for production
   - Never commit credentials to version control

3. **Network Security:**
   - Use HTTPS in production
   - Configure firewall rules
   - Limit database access to application only
   - Use reverse proxy (nginx/traefik) for SSL termination

4. **AWS Security:**
   - Use IAM roles instead of access keys when possible
   - Limit SES permissions to minimum required
   - Enable AWS CloudTrail for audit logging
   - Rotate AWS credentials regularly

5. **Database Security:**
   - Use strong database passwords
   - Enable SSL connections
   - Regular security updates
   - Database connection pooling

## 📈 Monitoring & Logging

### Application Logs

```bash
# View all logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres

# Follow logs with timestamps
docker-compose logs -f -t
```

### Background Services

The application runs two background services:

1. **Cleanup Service**: Automatically removes old event logs based on retention settings
2. **Sync Service**: Periodically syncs suppression list with AWS SES

### Performance Monitoring

Monitor application performance through:
- Dashboard analytics page
- Database query performance
- Container resource usage
- API response times

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

# Test API endpoint
curl http://localhost/health
```

5. **AWS Integration Issues:**
```bash
# Check AWS settings in database
docker-compose exec postgres psql -U ses_user -d ses_monitoring -c "SELECT * FROM app_settings WHERE key LIKE 'aws_%';"
```

### Performance Optimization

1. **Database Optimization:**
   - Configure retention policies to manage data size
   - Monitor query performance with `EXPLAIN ANALYZE`
   - Add custom indexes for specific query patterns
   - Regular `VACUUM` and `ANALYZE` operations

2. **Container Resources:**
   - Adjust memory limits in docker-compose.yml
   - Monitor container resource usage with `docker stats`
   - Scale services horizontally if needed
   - Use multi-stage builds to reduce image size

3. **Application Performance:**
   - Enable Go profiling for performance analysis
   - Use connection pooling for database
   - Implement caching for frequently accessed data
   - Optimize React components with memoization

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- **Backend**: Follow Go best practices and clean architecture
- **Frontend**: Use TypeScript and functional components
- **Database**: Use migrations for schema changes
- **Testing**: Write tests for new features
- **Documentation**: Update API documentation for changes
- **Commits**: Follow conventional commit messages

### Code Style

- **Go**: Use `gofmt` and `golint`
- **TypeScript**: Use ESLint and Prettier
- **SQL**: Use consistent naming conventions
- **Docker**: Multi-stage builds and minimal base images

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **AWS SES** for email service integration
- **PostgreSQL** for reliable data storage
- **React** and **Go** communities for excellent frameworks
- **Docker** for containerization platform
- **Tailwind CSS** for utility-first styling
- **Recharts** for beautiful data visualization

## 📞 Support

For support and questions:

1. Check the [Issues](../../issues) page
2. Review the troubleshooting section
3. Create a new issue with detailed information
4. Include logs and system information

---

**Made with ❤️ by Wisnu**
