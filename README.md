# Load Stuffing Calculator

A 3D container load planning system with multi-tenant support.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Quick Start](#quick-start)
- [Development](#development)
- [Testing](#testing)
- [Authentication](#authentication)
- [API Documentation](#api-documentation)
- [Database](#database)
- [Mobile App](#mobile-app)
- [Deployment](#deployment)
- [Troubleshooting](#troubleshooting)

## Overview

This application solves the 3D bin packing problem for container loading. It calculates optimal placement of items in containers, provides 3D visualization of the results, and manages user workspaces with role-based permissions.

**Primary users:**
- Warehouse planners (create and manage load plans)
- Operators (view plans and loading instructions)
- Administrators (manage users and master data)

## Features

**Core functionality:**
- 3D bin packing calculation using multiple algorithms (boxpacker3, py3dbp)
- Interactive 3D visualization with Three.js
- Step-by-step loading animation
- Product and container catalog management
- PDF report generation
- Multi-tenant workspace system
- Role-based access control (5 roles, 50+ permissions)
- Trial mode (no registration, limited to 3 plans)

**User roles:**

| Role | Description | Key Permissions |
|------|-------------|-----------------|
| `founder` | Platform superuser (SaaS founder) | Global `*` (all permissions) |
| `owner` | Workspace owner/CEO | Full workspace, member, invite, product, container, plan management |
| `personal` | Personal workspace owner | Product, container, plan management (no member/invite access) |
| `admin` | Workspace admin | Member and invite management, product, container, plan management |
| `planner` | Plan creator | Full plan access, read-only product/container |
| `operator` | Loading validator | Read plans, manage plan items, read-only product/container |
| `trial` | Anonymous trial user | Create plans (max 3), read-only product/container, no workspace |

## Architecture

### System components

```mermaid
graph LR
    A[Web UI<br/>Next.js] -->|REST| B[API Server<br/>Go + Gin]
    C[Mobile App<br/>Flutter] -->|REST| B
    B -->|SQLC| D[(PostgreSQL)]
    B -->|HTTP| E[Packing Service<br/>Python]
```

### Backend layers

```mermaid
graph TB
    A[HTTP Routes] --> B[Middleware<br/>JWT + RBAC]
    B --> C[Handlers]
    C --> D[Services<br/>Business Logic]
    D --> E[Store<br/>SQLC Queries]
    D --> F[Gateway<br/>Packing Service]
    E --> G[(Database)]
```

### Request flow for plan calculation

```mermaid
sequenceDiagram
    participant C as Client
    participant A as API
    participant D as Database
    participant P as Packing Service
    
    C->>A: POST /plans/:id/calculate
    A->>D: Fetch plan + items
    D-->>A: Plan data
    A->>P: POST /pack
    P-->>A: Placements
    A->>D: Insert placements
    A-->>C: Result with positions
```

## Tech Stack

**Backend:**
- Go 1.24+ (Gin framework)
- PostgreSQL 14+ (database)
- SQLC (type-safe SQL queries)
- Goose (migrations)
- JWT authentication
- Python 3.10+ (packing microservice)

**Frontend:**
- Next.js 16 (App Router)
- React 19
- TypeScript 5
- Tailwind CSS 4
- Three.js (3D visualization)
- Radix UI components

**Mobile:**
- Flutter 3.9+
- Dart 3.9+

**DevOps:**
- Docker + Docker Compose

## Project Structure

```
load-stuffing-calculator/
├── cmd/
│   ├── api/                    # API server entrypoint
│   ├── db/
│   │   ├── migrations/         # Goose SQL migrations
│   │   └── queries/            # SQLC query definitions
│   ├── packing/                # Python packing microservice
│   └── seed/                   # Database seeding
│
├── internal/                   # Go packages
│   ├── api/                    # Router setup
│   ├── auth/                   # JWT, password hashing
│   ├── cache/                  # Permission cache
│   ├── config/                 # Configuration
│   ├── dto/                    # Data transfer objects
│   ├── gateway/                # External service clients
│   ├── handler/                # HTTP handlers
│   ├── middleware/             # JWT, RBAC middleware
│   ├── packer/                 # Packing algorithm wrapper
│   ├── response/               # HTTP response helpers
│   ├── service/                # Business logic
│   ├── store/                  # SQLC generated code
│   └── types/                  # Domain types
│
├── web/                        # Next.js frontend
│   ├── app/                    # Pages and layouts
│   ├── components/             # React components
│   ├── hooks/                  # Data fetching hooks
│   └── lib/                    # API client, 3D visualizer
│
├── mobile/                     # Flutter app
│   ├── lib/
│   │   ├── pages/             # Screens
│   │   ├── services/          # API clients
│   │   ├── providers/         # State management
│   │   └── models/            # Data models
│   ├── android/
│   └── ios/
│
├── docs/                       # Documentation
│   ├── architecture1.dot       # Architecture diagram
│   ├── erd1.dot               # Database ERD
│   ├── flowchart/             # Process flows
│   └── plans/                 # Development plans
│
├── docker-compose.yml
├── Makefile
└── sqlc.yaml
```

## Quick Start

### Prerequisites

**Option A:** Docker + Docker Compose

**Option B:**
- Go 1.24+
- PostgreSQL 14+
- Python 3.10+
- Node.js 20+

### Using Docker Compose

```bash
# Clone repository
git clone https://github.com/ekastn/load-stuffing-calculator.git
cd load-stuffing-calculator

# Initialize submodules
git submodule update --init --recursive

# Configure
cp .env.example .env

# Start services
docker compose up --build

# Access
# API: http://localhost:8080
# Swagger: http://localhost:8080/docs/index.html
# Web: http://localhost:3000
# Login: admin@example.com / admin123
```

### Manual setup

<details>
<summary>Click to expand</summary>

**Database:**
```bash
createdb stuffing
export DATABASE_URL="postgres://postgres:password@localhost:5432/stuffing?sslmode=disable"
```

**Backend:**
```bash
go mod tidy
cp .env.example .env
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up
sqlc generate
swag init -g cmd/api/main.go -o internal/docs
make build
./bin/api
```

**Packing service:**
```bash
python3 -m venv cmd/packing/.venv
source cmd/packing/.venv/bin/activate
pip install -r cmd/packing/requirements.txt
python cmd/packing/app.py
```

**Frontend:**
```bash
cd web
pnpm install
pnpm dev
```

</details>

## Development

### Setup tools

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/swaggo/swag/cmd/swag@latest
go install go.uber.org/mock/mockgen@latest
```

### Backend workflow

```bash
make fmt                # Format code
make swag              # Generate Swagger docs
sqlc generate          # Generate DB code
make test              # Run tests
make coverage          # Test coverage
make build             # Build binary
make run               # Run server
```

### Frontend workflow

```bash
cd web
pnpm dev               # Dev server
pnpm test              # Run tests
pnpm lint              # Lint code
pnpm build             # Production build
```

### Database operations

```bash
# Check migration status
goose -dir cmd/db/migrations postgres "$DATABASE_URL" status

# Apply migrations
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up

# Rollback one migration
goose -dir cmd/db/migrations postgres "$DATABASE_URL" down

# Create new migration
goose -dir cmd/db/migrations create migration_name sql
```

### Common commands

```bash
# Run specific test
go test ./internal/service -run '^TestAuthService_Login$' -v

# Frontend specific test
pnpm -C web test -- -t "test name"

# Coverage report
make coverage-html

# Check per-package coverage
make coverage-summary
```

## Testing

**Current coverage:** 95.6% overall

**Coverage by package:**
- `cache`: 100%
- `config`: 100%
- `env`: 100%
- `middleware`: 100%
- `response`: 100%
- `types`: 100%
- `handler`: 97.8%
- `packer`: 96.5%
- `gateway`: 94.1%
- `service`: 94.3%
- `auth`: 93.6%

**Run tests:**
```bash
make test                    # All tests
make coverage               # With coverage report
make coverage-html          # HTML report
go test -v ./internal/service
cd web && pnpm test
```

**Test structure:**
- Backend: Table-driven tests with `t.Run()`
- Mocks: `go.uber.org/mock` for SQLC store
- Handler tests: `httptest` + `testify/assert`
- Frontend: Vitest + React Testing Library

## Authentication

### JWT tokens

**Access token:**
- Expiry: 2 hours
- Claims: `user_id`, `role`, `workspace_id`
- Algorithm: HS256 (HMAC-SHA256)
- Header: `Authorization: Bearer <token>`

**Refresh token:**
- Expiry: 30 days
- Format: `rf_YYYYMMDD_HHMMSS.microsec_random12`
- Stored in database with revocation tracking
- Rotated on each use (old token revoked, new token issued)

**Guest/trial token:**
- Expiry: 30 days
- Special access token for trial users
- No refresh token provided
- Limited to 3 plans per guest session

**Token lifecycle:**

1. **Login**: User provides credentials → API verifies and returns access token (2h) + refresh token (30d)
2. **Normal requests**: Client includes access token in Authorization header → API validates and processes request
3. **Token expires**: After 2 hours, API returns 401 Unauthorized
4. **Refresh**: Client sends refresh token → API validates, revokes old token, issues new token pair
5. **Switch workspace**: Client requests workspace change → API verifies membership and issues new access token with updated workspace_id

### Password security

- Algorithm: bcrypt with default cost (10)
- Passwords hashed on registration and password change
- Password verification uses constant-time comparison
- No plaintext passwords stored or logged

### Role-based access control

**Permission format:** `resource:action`

**Permission matching:**
- `*` - Global wildcard, matches all permissions (founder only)
- `resource:*` - Resource wildcard, matches all actions on resource (e.g., `plan:*` matches `plan:create`, `plan:read`, etc.)
- `resource:action` - Exact match for specific permission

**Examples:**
- `plan:create` - Create plans
- `product:read` - View products
- `user:*` - All user operations (create, read, update, delete)
- `*` - Global access to everything (founder only)

**Permission caching:**
- Permissions cached in-memory by role
- Cache lookup on each request
- Reduces database queries for permission checks

**Middleware usage:**
```go
// Single permission required
router.GET("/plans", middleware.Permission("plan:read"), handler)

// Any of multiple permissions (OR logic)
router.POST("/plans", middleware.Any("plan:create", "admin:*"), handler)

// All permissions required (AND logic)
router.DELETE("/users/:id", middleware.All("user:delete", "admin:access"), handler)
```

### Multi-tenancy

**Workspace types:**
- Personal: Auto-created per user
- Organization: Shared team workspace

**Data isolation:**
- All master data scoped by `workspace_id`
- Load plans scoped by `workspace_id`
- Queries automatically filter by workspace from JWT

**Switch workspace:**
```bash
POST /api/v1/auth/switch-workspace
{
  "workspace_id": "uuid",
  "refresh_token": "rf_..."
}
```

**Founder override:**
```bash
GET /api/v1/plans?workspace_id=other-workspace-uuid
```

### Trial mode

Trial users receive a special access token without registration.

**Characteristics:**
- Token expiry: 30 days (longer than regular access tokens)
- No refresh token provided
- No workspace association
- Plans limited to 3 per guest session

**Limitations:**
- Max 3 plans per guest
- Plans scoped by `created_by_type='guest'` and `created_by_id`
- Read-only access to product and container catalogs
- Cannot create workspaces or manage members

**Get trial token:**
```bash
POST /api/v1/auth/guest
Response: { "access_token": "..." }
```

## API Documentation

**Swagger UI:** http://localhost:8080/docs/index.html

### Authentication endpoints

| Method | Path | Description | Auth |
|--------|------|-------------|------|
| POST | `/api/v1/auth/register` | Create account | No |
| POST | `/api/v1/auth/login` | Login | No |
| POST | `/api/v1/auth/guest` | Get trial token | No |
| POST | `/api/v1/auth/refresh` | Refresh token | No |
| GET | `/api/v1/auth/me` | Current user | Yes |
| POST | `/api/v1/auth/switch-workspace` | Switch workspace | Yes |

### Load planning endpoints

| Method | Path | Permission |
|--------|------|------------|
| GET | `/api/v1/plans` | `plan:read` |
| POST | `/api/v1/plans` | `plan:create` |
| GET | `/api/v1/plans/:id` | `plan:read` |
| PUT | `/api/v1/plans/:id` | `plan:update` |
| DELETE | `/api/v1/plans/:id` | `plan:delete` |
| POST | `/api/v1/plans/:id/calculate` | `plan:calculate` |
| POST | `/api/v1/plans/:id/items` | `plan_item:create` |
| PUT | `/api/v1/plans/:id/items/:itemId` | `plan_item:update` |
| DELETE | `/api/v1/plans/:id/items/:itemId` | `plan_item:delete` |

### Master data endpoints

| Resource | Methods | Permission |
|----------|---------|------------|
| `/api/v1/containers` | GET, POST, PUT, DELETE | `container:{read,create,update,delete}` |
| `/api/v1/products` | GET, POST, PUT, DELETE | `product:{read,create,update,delete}` |

### Admin endpoints

| Resource | Methods | Permission |
|----------|---------|------------|
| `/api/v1/users` | GET, POST, PUT, DELETE | `user:{read,create,update,delete}` |
| `/api/v1/roles` | GET, POST, PUT, DELETE | `role:{read,create,update,delete}` |
| `/api/v1/permissions` | GET | `permission:read` |
| `/api/v1/workspaces` | GET, POST, PUT, DELETE | `workspace:{read,create,update,delete}` |
| `/api/v1/dashboard` | GET | `dashboard:read` |

### Example requests

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@example.com", "password": "admin123"}'
```

**Create plan:**
```bash
curl -X POST http://localhost:8080/api/v1/plans \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "container": {
      "length_mm": 12000,
      "width_mm": 2400,
      "height_mm": 2600,
      "max_weight_kg": 28000
    },
    "items": [{
      "label": "Box A",
      "length_mm": 600,
      "width_mm": 400,
      "height_mm": 300,
      "weight_kg": 25,
      "quantity": 10,
      "allow_rotation": true
    }],
    "auto_calculate": true
  }'
```

**Calculate plan:**
```bash
curl -X POST http://localhost:8080/api/v1/plans/{id}/calculate \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "strategy": "parallel",
    "goal": "minimizeboxes",
    "gravity": true
  }'
```

## Database

### Schema

```mermaid
erDiagram
    users ||--o{ refresh_tokens : has
    users ||--|| user_profiles : has
    users }o--|| roles : has
    roles ||--o{ role_permissions : has
    permissions ||--o{ role_permissions : has
    
    users ||--o{ workspaces : owns
    workspaces ||--o{ members : has
    workspaces ||--o{ invites : has
    workspaces ||--o{ containers : contains
    workspaces ||--o{ products : contains
    workspaces ||--o{ load_plans : contains
    
    load_plans ||--o{ load_items : has
    load_plans ||--o| plan_results : has
    plan_results ||--o{ plan_placements : has
    load_items ||--o{ plan_placements : references
```

### Tables

**Authentication:**
- `users` - User accounts
- `user_profiles` - User details
- `roles` - Role definitions
- `permissions` - Permission definitions
- `role_permissions` - Role-permission mapping
- `refresh_tokens` - Active tokens

**Multi-tenancy:**
- `workspaces` - Tenant containers
- `members` - Workspace membership
- `invites` - Workspace invitations
- `platform_members` - Founder access

**Master data:**
- `containers` - Container catalog
- `products` - Product catalog

**Planning:**
- `load_plans` - Plan header
- `load_items` - Items to pack
- `plan_results` - Calculation summary
- `plan_placements` - Item positions

### Management

**Migrations:**
```bash
goose -dir cmd/db/migrations postgres "$DATABASE_URL" status
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up
goose -dir cmd/db/migrations postgres "$DATABASE_URL" down
```

**SQLC:**
```bash
sqlc generate  # Generate Go code from queries
sqlc verify    # Verify query syntax
```

**Backup:**
```bash
# Docker
docker compose exec db pg_dump -U postgres stuffing > backup.sql

# Local
pg_dump -U postgres stuffing > backup.sql
```

**Restore:**
```bash
# Docker
docker compose exec -T db psql -U postgres stuffing < backup.sql

# Local
psql -U postgres stuffing < backup.sql
```

## Mobile App

Flutter application for iOS and Android.

**Features:**
- Login (username/password)
- Dashboard with statistics
- Plan management (list, create, view details)
- Product catalog (list, create, edit)
- Container catalog (list, create, edit)
- Profile management
- 3D visualization (embedded WebView from web app)

**Not implemented:**
- User registration
- Workspace switching
- Native 3D rendering

**Setup:**
```bash
cd mobile
flutter pub get
flutter run -d android   # or: ios, linux, macos, windows
```

**Build:**
```bash
flutter build apk --release        # Android APK
flutter build appbundle            # Android Bundle
flutter build ios --release        # iOS
```

For detailed documentation see [mobile/README.md](mobile/README.md).

## Deployment

### Docker

```bash
# Start services
docker compose up -d --build

# View logs
docker compose logs -f api

# Stop services
docker compose down
```

### Configuration

**Required environment variables:**

| Variable | Description | Example |
|----------|-------------|---------|
| `JWT_SECRET` | JWT signing key | `openssl rand -base64 32` |
| `DATABASE_URL` | Postgres connection | `postgres://user:pass@host:5432/db` |
| `PACKING_SERVICE_URL` | Packing service | `http://packing:5051` |
| `FOUNDER_EMAIL` | Initial admin email | `admin@company.com` |
| `FOUNDER_PASSWORD` | Initial admin password | Strong password |
| `SRV_ENV` | Environment | `dev` or `production` |

See `.env.example` for complete list.

### Health checks

```bash
# API
curl http://localhost:8080/api/v1/health

# Packing service
curl http://localhost:5051/health

# Database
docker compose exec db pg_isready -U postgres
```

### Monitoring

```bash
# View logs
docker compose logs -f
docker compose logs -f api

# Container status
docker compose ps

# Resource usage
docker stats
```

## Troubleshooting

### "Authorization header is required"

Missing JWT token. Include header:
```http
Authorization: Bearer <token>
```

### "Invalid or expired token"

Token expired (2h TTL). Refresh via:
```bash
POST /api/v1/auth/refresh
```

### Packing calculation fails

Check service is running:
```bash
curl http://localhost:5051/health
docker compose logs packing
```

Verify `PACKING_SERVICE_URL` environment variable.

### Database connection errors

Check PostgreSQL is running:
```bash
docker compose ps db
pg_isready -h localhost -p 5432
```

Verify `DATABASE_URL` format and credentials.

Check migrations applied:
```bash
goose -dir cmd/db/migrations postgres "$DATABASE_URL" status
```

### CORS errors

Check API CORS config in `internal/api/api.go`.
Verify `NEXT_PUBLIC_API_URL` in frontend.

### Permission denied

Check user role and permissions in database.
Verify role has required permission in `role_permissions` table.

### Port conflicts

Change ports in `.env`:
```bash
API_EXPOSE_PORT=8081
WEB_EXPOSE_PORT=3001
```

Or kill process:
```bash
lsof -ti:8080 | xargs kill
```

