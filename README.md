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

- **Docker** (version 20.10+)
- **Docker Compose** (version 2.0+)
- **Git**
- **Curl** (for health checks)

## 🛠️ Quick Start

### 1. Clone the Repository
```bash\ngit clone <repository-url>\ncd ses-dashboard-monitoring\n```\n\n### 2. Run Installation Script
```bash
chmod +x install.sh
./install.sh
```

The installation script will:
- ✅ Check Docker and Docker Compose installation
- ✅ Create necessary directories
- ✅ Build and start all services
- ✅ Run database migrations automatically
- ✅ Display service URLs and credentials\n\n### 3. Access the Application\n\n| Service | URL | Description |\n|---------|-----|-------------|\n| **Application** | http://localhost | Complete SES Dashboard |\n| **API Documentation** | http://localhost/swagger/index.html | Swagger UI |\n| **Database** | localhost:5432 | PostgreSQL (admin access) |\n\n### 4. Default Credentials\n```\nUsername: admin\nPassword: admin123\n```\n\n**That's it!** The `install.sh` script handles everything - from building Docker images to database initialization. No manual steps required.\n\n## 🏗️ Architecture\n\n```\n┌─────────────────┐    ┌─────────────────┐\n│   Frontend      │    │   Backend       │\n│   (Node.js +    │◄──►│   (Go/Gin)      │\n│    Express)     │    │   Port: 8080    │\n│   Port: 80      │    └─────────────────┘\n└─────────────────┘              │\n          │                      │\n          │        ┌─────────────────┐\n          └────────┤   PostgreSQL    │\n                   │   Port: 5432    │\n                   └─────────────────┘\n```\n\n### Components\n\n- **Frontend**: React 19 + TypeScript + Tailwind CSS served by Express.js\n- **Backend**: Go 1.25 + Gin framework with clean architecture\n- **Database**: PostgreSQL 15 with optimized schema and indexes\n- **Proxy**: Express.js handles API routing and static file serving\n\n## 📁 Project Structure\n\n```\nses-dashboard-monitoring/\n├── ses-dashboard-monitoring/          # Backend (Go)\n│   ├── cmd/api/                      # Application entry point\n│   │   └── main.go                   # Main application file\n│   ├── internal/                     # Internal packages\n│   │   ├── config/                   # Configuration management\n│   │   ├── delivery/http/            # HTTP handlers and middleware\n│   │   ├── domain/                   # Business logic and entities\n│   │   │   ├── sesevent/             # SES event domain\n│   │   │   ├── settings/             # Settings domain\n│   │   │   ├── suppression/          # Suppression domain\n│   │   │   └── user/                 # User domain\n│   │   ├── infrastructure/           # External dependencies\n│   │   │   ├── aws/                  # AWS SES client\n│   │   │   ├── database/             # Database connection\n│   │   │   └── repository/           # Data access layer\n│   │   ├── services/                 # Background services\n│   │   │   ├── cleanup_service.go    # Data cleanup automation\n│   │   │   └── sync_service.go       # AWS sync automation\n│   │   └── usecase/                  # Business use cases\n│   ├── config/                       # Configuration files\n│   │   └── config.yaml               # Application configuration\n│   ├── docs/                         # Swagger documentation\n│   └── Dockerfile                    # Backend Docker image\n├── ses-dashboard-frontend/           # Frontend (React + Node.js)\n│   ├── src/                         # React source code\n│   │   ├── components/              # Reusable React components\n│   │   │   ├── Layout.tsx           # Main layout component\n│   │   │   ├── Charts.tsx           # Chart components\n│   │   │   └── EventsTable.tsx      # Events table component\n│   │   ├── pages/                   # Page components\n│   │   │   ├── DashboardPage.tsx    # Main dashboard\n│   │   │   ├── EventsPage.tsx       # Events listing\n│   │   │   ├── AnalyticsPage.tsx    # Analytics charts\n│   │   │   ├── SuppressionPage.tsx  # Suppression management\n│   │   │   ├── UsersPage.tsx        # User management\n│   │   │   └── SettingsPage.tsx     # System settings\n│   │   ├── services/                # API services\n│   │   │   └── api.ts               # API client with Axios\n│   │   └── types/                   # TypeScript type definitions\n│   ├── server.js                    # Express server with proxy\n│   └── Dockerfile                   # Frontend Docker image\n├── docker-compose.yml               # Docker Compose configuration\n├── init.sql                         # Database initialization\n├── install.sh                       # Installation script\n└── README.md                        # This file\n```\n\n## ⚙️ Configuration\n\n### Environment Variables\n\n| Variable | Description | Default |\n|----------|-------------|---------|\n| `DB_HOST` | Database host | `postgres` |\n| `DB_PORT` | Database port | `5432` |\n| `DB_USER` | Database username | `ses_user` |\n| `DB_PASSWORD` | Database password | `ses_password` |\n| `DB_NAME` | Database name | `ses_monitoring` |\n| `JWT_SECRET` | JWT signing secret | `your-super-secret-jwt-key` |\n| `PORT` | Backend server port | `8080` |\n| `BACKEND_URL` | Backend URL for frontend proxy | `http://backend:8080` |\n\n### Database Schema\n\nThe application uses 4 main tables:\n\n1. **users** - User accounts with role-based access\n2. **ses_events** - SES event logs with full event data\n3. **app_settings** - Configurable application settings\n4. **suppressions** - Email suppression list management\n\n### AWS SES Configuration\n\nConfigure AWS SES settings through the admin panel:\n\n1. Navigate to **Admin → Settings**\n2. Configure AWS credentials and region\n3. Test the connection\n4. Enable AWS integration\n\n### SNS Webhook Setup\n\nTo receive SES events via SNS:\n\n1. Create an SNS topic in AWS\n2. Subscribe your endpoint: `http://your-domain/sns/ses`\n3. Configure SES to publish events to the SNS topic\n4. Events will be automatically processed and stored\n\n## 🔧 Development

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

### Local Development Setup\n\n1. **Backend Development:**\n```bash\ncd ses-dashboard-monitoring\ngo mod download\ngo run cmd/api/main.go\n```\n\n2. **Frontend Development:**\n```bash\ncd ses-dashboard-frontend\nnpm install\nnpm run dev\n```\n\n3. **Database Setup:**
```bash
# Start PostgreSQL only
docker-compose up postgres -d

# Run migrations
cd ses-dashboard-monitoring
make migrate-up
```\n\n### Building Docker Images\n\n> **Note:** This is handled automatically by `install.sh`. Manual build only needed for development.\n\n```bash\n# Build all services\ndocker-compose build\n\n# Build specific service\ndocker-compose build backend\ndocker-compose build frontend\n```\n\n## 📊 API Documentation\n\nThe API documentation is available via Swagger UI at:\n```\nhttp://localhost/swagger/index.html\n```\n\n### Key API Endpoints\n\n#### Authentication\n| Method | Endpoint | Description |\n|--------|----------|-------------|\n| `POST` | `/api/login` | User authentication |\n| `PUT` | `/api/change-password` | Change user password |\n\n#### Events & Metrics\n| Method | Endpoint | Description |\n|--------|----------|-------------|\n| `GET` | `/api/events` | Get SES events with pagination |\n| `GET` | `/api/metrics` | Get dashboard metrics |\n| `GET` | `/api/metrics/daily` | Get daily analytics |\n| `GET` | `/api/metrics/monthly` | Get monthly analytics |\n| `GET` | `/api/metrics/hourly` | Get hourly analytics |\n\n#### Suppression Management\n| Method | Endpoint | Description |\n|--------|----------|-------------|\n| `GET` | `/api/suppression` | Get suppression list |\n| `POST` | `/api/suppression` | Add single email to suppression |\n| `POST` | `/api/suppression/bulk` | **Bulk add** multiple emails |\n| `DELETE` | `/api/suppression/bulk` | **Bulk remove** multiple emails |\n| `DELETE` | `/api/suppression/:email` | Remove single email |\n| `GET` | `/api/suppression/:email/status` | Check email AWS status |\n| `POST` | `/api/suppression/sync` | Trigger AWS sync |\n| `GET` | `/api/suppression/sync/status` | Get sync status |\n\n#### Administration (Admin Only)\n| Method | Endpoint | Description |\n|--------|----------|-------------|\n| `GET` | `/api/users` | Get all users |\n| `POST` | `/api/users` | Create new user |\n| `PUT` | `/api/users/:id/reset-password` | Reset user password |\n| `PUT` | `/api/users/:id/disable` | Disable user account |\n| `DELETE` | `/api/users/:id` | Delete user account |\n| `GET` | `/api/settings/aws` | Get AWS settings |\n| `PUT` | `/api/settings/aws` | Update AWS settings |\n| `POST` | `/api/settings/aws/test` | Test AWS connection |\n| `GET` | `/api/settings/retention` | Get retention settings |\n| `PUT` | `/api/settings/retention` | Update retention settings |\n\n## 🛠️ Management Commands\n\n### Docker Compose Commands\n\n```bash\n# Start all services\ndocker-compose up -d\n\n# View logs\ndocker-compose logs -f\n\n# View specific service logs\ndocker-compose logs -f backend\ndocker-compose logs -f frontend\n\n# Stop all services\ndocker-compose down\n\n# Restart specific service\ndocker-compose restart backend\n\n# View service status\ndocker-compose ps\n\n# Execute command in container\ndocker-compose exec backend sh\n```\n\n### Database Management\n\n```bash\n# Connect to database\ndocker-compose exec postgres psql -U ses_user -d ses_monitoring\n\n# Backup database\ndocker-compose exec postgres pg_dump -U ses_user ses_monitoring > backup.sql\n\n# Restore database\ndocker-compose exec -T postgres psql -U ses_user ses_monitoring < backup.sql\n\n# View database size\ndocker-compose exec postgres psql -U ses_user -d ses_monitoring -c \"\\l+\"\n```\n\n## 🔒 Security Considerations\n\n### Production Deployment\n\n1. **Change Default Credentials:**\n   - Update admin password after first login\n   - Use strong JWT secret key (change `JWT_SECRET` environment variable)\n\n2. **Environment Variables:**\n   - Store sensitive data in environment variables\n   - Use Docker secrets for production\n   - Never commit credentials to version control\n\n3. **Network Security:**\n   - Use HTTPS in production\n   - Configure firewall rules\n   - Limit database access to application only\n   - Use reverse proxy (nginx/traefik) for SSL termination\n\n4. **AWS Security:**\n   - Use IAM roles instead of access keys when possible\n   - Limit SES permissions to minimum required\n   - Enable AWS CloudTrail for audit logging\n   - Rotate AWS credentials regularly\n\n5. **Database Security:**\n   - Use strong database passwords\n   - Enable SSL connections\n   - Regular security updates\n   - Database connection pooling\n\n## 📈 Monitoring & Logging\n\n### Application Logs\n\n```bash\n# View all logs\ndocker-compose logs -f\n\n# View specific service logs\ndocker-compose logs -f backend\ndocker-compose logs -f frontend\ndocker-compose logs -f postgres\n\n# Follow logs with timestamps\ndocker-compose logs -f -t\n```\n\n### Background Services\n\nThe application runs two background services:\n\n1. **Cleanup Service**: Automatically removes old event logs based on retention settings\n2. **Sync Service**: Periodically syncs suppression list with AWS SES\n\n### Performance Monitoring\n\nMonitor application performance through:\n- Dashboard analytics page\n- Database query performance\n- Container resource usage\n- API response times\n\n## 🚨 Troubleshooting\n\n### Common Issues\n\n1. **Port Already in Use:**\n```bash\n# Check what's using the port\nlsof -i :80\nlsof -i :8080\n\n# Stop conflicting services\nsudo systemctl stop apache2\nsudo systemctl stop nginx\n```\n\n2. **Database Connection Issues:**\n```bash\n# Check database logs\ndocker-compose logs postgres\n\n# Verify database is running\ndocker-compose ps postgres\n\n# Test database connection\ndocker-compose exec postgres psql -U ses_user -d ses_monitoring -c \"SELECT 1;\"\n```\n\n3. **Frontend Build Issues:**\n```bash\n# Clear npm cache\nnpm cache clean --force\n\n# Rebuild frontend\ndocker-compose build --no-cache frontend\n```\n\n4. **Backend API Issues:**\n```bash\n# Check backend logs\ndocker-compose logs backend\n\n# Verify configuration\ndocker-compose exec backend env | grep DB_\n\n# Test API endpoint\ncurl http://localhost:8080/api/health\n```\n\n5. **AWS Integration Issues:**\n```bash\n# Test AWS credentials\ndocker-compose exec backend aws ses describe-configuration-sets\n\n# Check AWS settings in database\ndocker-compose exec postgres psql -U ses_user -d ses_monitoring -c \"SELECT * FROM app_settings WHERE key LIKE 'aws_%';\"\n```\n\n### Performance Optimization\n\n1. **Database Optimization:**\n   - Configure retention policies to manage data size\n   - Monitor query performance with `EXPLAIN ANALYZE`\n   - Add custom indexes for specific query patterns\n   - Regular `VACUUM` and `ANALYZE` operations\n\n2. **Container Resources:**\n   - Adjust memory limits in docker-compose.yml\n   - Monitor container resource usage with `docker stats`\n   - Scale services horizontally if needed\n   - Use multi-stage builds to reduce image size\n\n3. **Application Performance:**\n   - Enable Go profiling for performance analysis\n   - Use connection pooling for database\n   - Implement caching for frequently accessed data\n   - Optimize React components with memoization\n\n## 🤝 Contributing\n\n1. Fork the repository\n2. Create a feature branch (`git checkout -b feature/amazing-feature`)\n3. Commit your changes (`git commit -m 'Add some amazing feature'`)\n4. Push to the branch (`git push origin feature/amazing-feature`)\n5. Open a Pull Request\n\n### Development Guidelines\n\n- **Backend**: Follow Go best practices and clean architecture\n- **Frontend**: Use TypeScript and functional components\n- **Database**: Use migrations for schema changes\n- **Testing**: Write tests for new features\n- **Documentation**: Update API documentation for changes\n- **Commits**: Follow conventional commit messages\n\n### Code Style\n\n- **Go**: Use `gofmt` and `golint`\n- **TypeScript**: Use ESLint and Prettier\n- **SQL**: Use consistent naming conventions\n- **Docker**: Multi-stage builds and minimal base images\n\n## 📄 License\n\nThis project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.\n\n## 🙏 Acknowledgments\n\n- **AWS SES** for email service integration\n- **PostgreSQL** for reliable data storage\n- **React** and **Go** communities for excellent frameworks\n- **Docker** for containerization platform\n- **Tailwind CSS** for utility-first styling\n- **Recharts** for beautiful data visualization\n\n## 📞 Support\n\nFor support and questions:\n\n1. Check the [Issues](../../issues) page\n2. Review the troubleshooting section\n3. Create a new issue with detailed information\n4. Include logs and system information\n\n---\n\n**Made with ❤️ by Wisnu**