# Agent Commission Management System

**Agent Incentive, Commission and Producer Management System for India Post PLI**

A comprehensive Golang-based microservice for managing agent onboarding, commission calculation, disbursement, and lifecycle management for insurance agents.

## 📋 Overview

This system handles:
- ✅ Agent onboarding (Advisors, Coordinators, Departmental Employees, Field Officers)
- ✅ Agent profile and license management
- ✅ Commission rate configuration and processing
- ✅ Commission batch calculation with 6-hour SLA
- ✅ Trial and final statement generation
- ✅ Disbursement processing (Cheque & EFT) with 10-day SLA
- ✅ Commission clawback and suspense account management
- ✅ Comprehensive reporting and audit trails

## 🏗️ Architecture

- **Framework**: N-API Template (Golang)
- **Database**: PostgreSQL 16 with pgx driver
- **Workflow Engine**: Temporal.io for long-running processes
- **Architecture Pattern**: Hexagonal/Ports-and-Adapters
- **Dependency Injection**: Uber FX
- **HTTP Framework**: Gin
- **Query Builder**: Squirrel

## 📊 Project Statistics

- **Total APIs**: 105 endpoints
- **Temporal Workflows**: 8 workflows
- **Database Tables**: 20+ tables
- **Implementation Phases**: 12 phases
- **Estimated Timeline**: 60 days

## 🚀 Quick Start

### Prerequisites

- Go 1.25+
- Docker & Docker Compose
- PostgreSQL 16
- (Optional) golang-migrate for migrations

### 1. Clone and Setup

```bash
cd agent-commission
cp configs/config.yaml configs/config.local.yaml
# Edit config.local.yaml with your settings
```

### 2. Start Infrastructure

```bash
# Start PostgreSQL + Temporal + Redis
docker-compose up -d

# Verify services are running
docker-compose ps

# View logs
docker-compose logs -f
```

### 3. Run Database Migrations

```bash
# Install golang-migrate (if not already installed)
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
make migrate-up

# Or manually:
migrate -path db/migrations \
  -database "postgresql://postgres:postgres@localhost:5432/agent_commission_dev?sslmode=disable" \
  up
```

### 4. Build and Run Application

```bash
# Build
make build

# Run
make run

# Or directly with go
go run main.go
```

### 5. Access Services

- **Application**: http://localhost:8080
- **Health Check**: http://localhost:8080/health
- **Temporal UI**: http://localhost:8080 (Temporal Web)
- **PgAdmin**: http://localhost:5050 (with `--profile dev-tools`)
- **API Documentation**: http://localhost:8080/swagger (auto-generated)

## 📁 Project Structure

```
agent-commission/
├── main.go                    # Application entry point
├── go.mod                     # Go module dependencies
├── go.sum                     # Dependency checksums
├── docker-compose.yml         # Infrastructure setup
├── Makefile                   # Build and utility commands
├── bootstrap/
│   └── bootstrapper.go        # FX dependency injection modules
├── configs/
│   ├── config.yaml            # Base configuration
│   ├── config.dev.yaml        # Development config
│   ├── config.test.yaml       # Test config
│   └── config.prod.yaml       # Production config
├── core/
│   ├── domain/                # Domain models
│   └── port/                  # Port interfaces (request/response contracts)
├── handler/                   # HTTP handlers (controllers)
│   ├── request.go             # Request DTOs
│   └── response/              # Response DTOs
├── repo/
│   └── postgres/              # PostgreSQL repositories (data access)
├── workflows/                 # Temporal workflows
│   └── activities/            # Temporal activities
├── db/
│   └── migrations/            # Database migration scripts
├── docs/                      # API documentation (auto-generated)
└── test/                      # Integration tests
```

## 🔧 Development

### Available Make Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run unit tests
make test-coverage # Run tests with coverage report
make lint          # Run linter
make fmt           # Format code
make migrate-up    # Apply database migrations
make migrate-down  # Rollback last migration
make docker-up     # Start docker services
make docker-down   # Stop docker services
make clean         # Clean build artifacts
```

### Running Tests

```bash
# Unit tests
make test

# Integration tests (requires docker services running)
make test-integration

# End-to-end tests
make test-e2e

# Coverage report
make test-coverage
open coverage.html
```

### Code Generation

```bash
# Generate validators (auto-generated from struct tags)
make generate

# Generate mocks for testing
make generate-mocks
```

## 📝 Configuration

Configuration is loaded from `configs/` directory based on environment:

- `config.yaml` - Base configuration
- `config.dev.yaml` - Development overrides
- `config.test.yaml` - Testing overrides
- `config.prod.yaml` - Production settings

Set environment with:
```bash
export APP_ENV=dev  # or test, prod
```

## 🗄️ Database

### Migrations

Migrations are located in `db/migrations/` and managed using golang-migrate.

```bash
# Apply all migrations
make migrate-up

# Rollback last migration
make migrate-down

# Create new migration
migrate create -ext sql -dir db/migrations -seq migration_name

# Check migration version
migrate -path db/migrations -database "postgresql://..." version
```

### Schema

See `db/migrations/README.md` for complete schema documentation.

## ⚡ Temporal Workflows

### Workflows Implemented

1. **Agent Onboarding** - Multi-step onboarding with HRMS integration
2. **Commission Batch Processing** - 6-hour SLA batch calculation
3. **Trial Statement Approval** - 7-day SLA approval workflow
4. **Disbursement Processing** - 10-day SLA payment processing
5. **License Renewal** - 3-day SLA renewal workflow
6. **License Reminders** - Scheduled reminders (T-30, T-15, T-7, T-0)
7. **Agent Termination** - Multi-step termination with commission settlement
8. **Clawback Recovery** - Graduated recovery workflow

### Monitoring Workflows

Access Temporal UI at http://localhost:8080 to monitor workflow execution, view history, and troubleshoot issues.

## 📊 API Endpoints

### Agent Management
- `POST /agents/new/init` - Initialize agent onboarding
- `POST /agents/search` - Search agents
- `GET /agents/{agentId}` - Get agent details
- `PUT /agents/{agentId}` - Update agent profile

### License Management
- `POST /agents/{agentId}/licenses` - Add license
- `GET /agents/{agentId}/licenses` - List licenses
- `POST /agents/{agentId}/licenses/{licenseId}/renew` - Renew license

### Commission Management
- `POST /commission/batch/trigger` - Trigger commission batch
- `GET /commission/trial-statements` - List trial statements
- `POST /commission/trial-statements/{id}/approve` - Approve statement
- `GET /commission/disbursements` - List disbursements

See full API documentation at `/swagger` when running the application.

## 🔐 Security

- JWT-based authentication (handled by n-api-server)
- Role-based access control
- Input validation using go-playground/validator
- SQL injection prevention via prepared statements
- Rate limiting and request throttling

## 📈 Monitoring

- **Metrics**: Prometheus metrics exposed at `/metrics`
- **Health Check**: `/health` endpoint
- **Profiling**: `/debug/pprof` (dev only)
- **Stats**: `/debug/statsviz` (dev only)
- **Tracing**: OpenTelemetry integration

## 🧪 Testing

### Unit Tests
Located alongside source files (`*_test.go`)

### Integration Tests
Located in `test/` directory, test complete workflows with real PostgreSQL and Temporal.

### Coverage Target
80%+ code coverage required for all phases.

## 📖 Documentation

- **Implementation Plan**: See `/plan.md`
- **Project Context**: See `/context.md`
- **Requirements**: See `/Incentive/analysis/IC_Incentive_Commission_Producer_Management_Analysis.md`
- **Database Schema**: See `/db/migrations/README.md`
- **API Documentation**: Auto-generated Swagger at `/swagger`

## 🛠️ Troubleshooting

### Application won't start
```bash
# Check if ports are available
lsof -i :8080  # Application port
lsof -i :5432  # PostgreSQL port
lsof -i :7233  # Temporal port

# Check docker services
docker-compose ps
docker-compose logs
```

### Database connection issues
```bash
# Test database connectivity
psql -h localhost -U postgres -d agent_commission_dev

# Check migrations
make migrate-version
```

### Temporal workflow issues
```bash
# View Temporal UI
open http://localhost:8080

# Check Temporal logs
docker-compose logs temporal
```

## 🤝 Contributing

### Code Style
- Follow Go best practices
- Use `gofmt` for formatting
- Run linter before committing: `make lint`
- Add traceability comments with FR/BR/VR/WF references

### Commit Messages
```
[Phase X.Y] Brief description

Detailed description of changes

Traceability:
- FR-IC-XXX-YYY
- BR-IC-XXX-YYY

https://claude.ai/code/session_ID
```

## 📄 License

Proprietary - India Post PLI

## 📧 Contact

For questions or support, contact the development team.

---

**Status**: Phase 0 Complete ✅  
**Last Updated**: 2026-01-28  
**Version**: 1.0.0
