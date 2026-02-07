# Agent Notes (load-stuffing-calculator)

## Quick Reference

Multi-tenant 3D container load planning system with Go backend, Next.js web client, and Flutter mobile client. Implements bin packing algorithms, role-based access control, workspace management, and real-time 3D visualization.

**Tech Stack:**
- Backend: Go 1.24+ (Gin), PostgreSQL 14+ (SQLC + Goose), Python 3.10+ (packing service)
- Web: Next.js 16, React 19, TypeScript 5, Three.js, Tailwind CSS 4, shadcn/ui
- Mobile: Flutter 3.9+, Dart 3.9+, Material 3, Provider state management

**Key Metrics:**
- Backend test coverage: 95.6% (cache/config/env/middleware/response/types: 100%)
- Web test coverage: Service layer only (8 test files, no component/E2E tests)
- Mobile test coverage: 0% (no tests implemented)
- StuffingVisualizer: 1731 lines, 14 files (Three.js visualization engine)
- Total Dart code (mobile): ~10,262 lines across 72 source files

## Project Overview

This is a SaaS load planning system that solves the 3D bin packing problem for warehouse container loading. It supports multi-tenant workspaces, comprehensive role-based access control, and provides real-time 3D visualization of container stuffing calculations.

**Three main components:**

1. **Go API Backend** - REST API server with JWT authentication, RBAC middleware, plan calculation service, and packing algorithm integration
2. **Next.js Web Client** - Multi-tenant web application with native Three.js 3D visualization, trial calculator, and admin features
3. **Flutter Mobile Client** - Native mobile app for iOS/Android with WebView-based 3D visualization

**Target users:**
- Platform administrators (founder, admin) - Manage users, workspaces, roles across system
- Workspace owners - Manage organization workspaces and team members
- Planners - Create and optimize load plans
- Operators - View plans and execute loading instructions
- Trial users - Anonymous users exploring features via guest sessions

**Primary workflow:**
1. Planner creates shipment/plan with container selection
2. Planner adds items (from catalog or manual entry)
3. System calculates 3D placements using bin packing algorithm
4. Results visualized in 3D with step-by-step playback
5. Operator follows loading sequence on mobile/web
6. System tracks execution (planned, not implemented)

## Repository Structure

```
load-stuffing-calculator/
├── cmd/
│   ├── api/                    # Go API server entrypoint (main.go)
│   ├── db/
│   │   ├── migrations/         # Goose SQL migrations (versioned)
│   │   ├── queries/            # SQLC query definitions (.sql files)
│   │   └── scripts/            # Default users, permissions SQL
│   ├── packing/                # Python packing microservice (Flask)
│   └── seed/                   # Database seeding utilities
├── internal/                   # Go application code (private)
│   ├── api/                    # Router setup, CORS config
│   ├── auth/                   # JWT generation/validation, password hashing (bcrypt)
│   ├── cache/                  # Permission caching (in-memory)
│   ├── config/                 # Configuration loading (env vars)
│   ├── dto/                    # Data transfer objects (API request/response)
│   ├── gateway/                # External service clients (packing service)
│   ├── handler/                # Gin HTTP handlers (bind, validate, call services)
│   ├── middleware/             # JWT auth, permission checks, CORS
│   ├── packer/                 # Bin packing algorithm wrapper (boxpacker3)
│   ├── response/               # HTTP response helpers (Success, Error)
│   ├── service/                # Business logic (auth, plan, user, role, etc.)
│   ├── store/                  # SQLC generated code (queries + models)
│   └── types/                  # Domain types
├── web/                        # Next.js web client (primary UI)
│   ├── app/                    # Next.js App Router pages
│   │   ├── (app)/             # Authenticated routes (dashboard, shipments, etc.)
│   │   ├── (auth)/            # Login, register
│   │   ├── (public)/          # Landing page, embed viewer
│   │   └── (platform)/        # Platform admin routes (users, workspaces, roles)
│   ├── components/            # React components (ui/, forms, viewers)
│   ├── hooks/                 # Data fetching hooks (use-*.ts)
│   ├── lib/
│   │   ├── services/          # API service layer (11 services)
│   │   ├── types/             # TypeScript type definitions (418 lines)
│   │   ├── StuffingVisualizer/ # 3D visualization engine (1731 lines, 14 files)
│   │   ├── api.ts             # HTTP client with token refresh
│   │   ├── auth-context.tsx   # Authentication context provider
│   │   └── route-guard.tsx    # Route protection component
│   └── test/                  # Vitest tests (8 service tests)
├── mobile/                     # Flutter mobile client
│   ├── lib/
│   │   ├── config/            # Theme, routes, constants, assets
│   │   ├── components/        # Reusable UI widgets
│   │   ├── dtos/              # API JSON objects (json_serializable)
│   │   ├── mappers/           # DTO to Model transformers
│   │   ├── models/            # Domain models (Freezed)
│   │   ├── pages/             # Screens (auth, dashboard, plans, products, containers, profile)
│   │   ├── providers/         # State management (ChangeNotifier)
│   │   ├── services/          # API services, storage service
│   │   └── utils/             # UI helpers
│   ├── android/               # Android native config
│   ├── ios/                   # iOS native config
│   ├── linux/                 # Linux desktop config
│   ├── macos/                 # macOS desktop config
│   ├── windows/               # Windows desktop config
│   └── test/                  # Tests (empty, not implemented)
├── docs/                       # Documentation
│   ├── architecture.dot       # System architecture diagram (includes roadmap items)
│   ├── erd.dot               # Database ERD
│   ├── flowchart/            # Process flows
│   └── plans/                # Development plans
├── docker-compose.yml         # Docker services (db, api, packing, web)
├── Makefile                   # Go build/test/run commands
├── sqlc.yaml                  # SQLC configuration
├── .env.example               # Environment variables template
└── README.md                  # Main documentation
```

## Architecture

### System Architecture

**Component interaction:**
```
┌─────────────┐
│  Web Client │ ─── REST API ───┐
└─────────────┘                  │
                                 ▼
┌─────────────┐           ┌──────────────┐         ┌─────────────────┐
│Mobile Client│ ─ REST ─→ │   Go API     │ ─ HTTP ─→│Packing Service │
└─────────────┘           │  (Gin + JWT) │         │   (Python)     │
                          └──────────────┘         └─────────────────┘
                                 │
                                 ▼
                          ┌──────────────┐
                          │ PostgreSQL   │
                          │   (SQLC)     │
                          └──────────────┘
```

**Multi-tenant workspace scoping:**
- JWT tokens embed `workspace_id` in claims
- Backend middleware extracts workspace context from token
- All database queries automatically filtered by workspace
- No custom headers required (no `X-Workspace-ID`)
- Users in multiple workspaces switch via `/auth/switch-workspace` endpoint
- Returns new JWT tokens scoped to target workspace

**Authentication flow:**
1. User logs in → backend generates JWT access token (2h) + refresh token (30d)
2. Tokens stored in client (localStorage for web, flutter_secure_storage for mobile)
3. All requests include `Authorization: Bearer <access_token>` header
4. On 401 response, client calls `/auth/refresh` with refresh token
5. New access token returned, original request retried
6. If refresh fails, user redirected to login

### Backend (Go API)

**Layer structure:**
```
HTTP Request
    ↓
Routes (internal/api/routes.go)
    ↓
Middleware (JWT auth + permission check)
    ↓
Handler (bind/validate input)
    ↓
Service (business logic, DB operations)
    ↓
Store (SQLC generated queries)
    ↓
PostgreSQL Database
```

**Packing engine:**
- Location: `internal/packer/packer.go`
- Wraps `github.com/bavix/boxpacker3` library
- Supports multiple strategies: bestfitdecreasing, minimizeboxes, greedy, parallel
- Goal options: minimizeboxes, maximizevolume
- Gravity simulation optional
- Units: dimensions in millimeters, weights in kg (converted to grams for boxpacker3)
- Returns placements with position (x, y, z), rotation code (0-5), and step number

**Database:**
- PostgreSQL 14+ with SQLC for type-safe queries
- Migrations managed by Goose (`cmd/db/migrations/`)
- Schema documented in `docs/erd.dot` (ERD diagram)
- Key tables:
  - Auth/RBAC: users, roles, permissions, role_permissions, refresh_tokens
  - Multi-tenant: workspaces, members, invites, platform_members
  - Master data: products, containers
  - Planning: load_plans, load_items, plan_results, plan_placements
- Bulk insert placements via SQLC's CopyFrom (efficient for large datasets)

**SQLC code generation:**
- Query definitions: `cmd/db/queries/*.sql`
- Generated code: `internal/store/*.sql.go`
- Type-safe Go structs matching database schema
- Automatic query method generation

### Web Client (Next.js)

**App Router structure:**
```
Route Groups:
(auth) → Login, Register (minimal layout)
(public) → Landing page, Embed viewer (public access)
(app) → Authenticated routes (dashboard layout with sidebar)
(platform) → Platform admin routes (nested under app)
```

**Service layer pattern:**
```
Component → Hook (use-*.ts) → Service (services/*.ts) → API Client (api.ts) → Backend
```

**StuffingVisualizer engine:**
- Location: `lib/StuffingVisualizer/` (1731 lines, 14 files)
- Native Three.js implementation with OrthographicCamera
- Manager pattern: Camera, Renderer, Scene, Lights, Controls, Animation, Interaction
- Features: step-by-step playback, hover tooltips, auto-fit camera, screenshots, PDF export
- Performance: proper disposal, requestAnimationFrame, efficient raycasting
- Used in: shipment detail page, embed viewer, trial calculator

**Multi-tenant features:**
- Workspace switching dropdown in sidebar
- Calls `POST /auth/switch-workspace` → new JWT tokens
- Navigation filtered by user permissions
- Platform admin routes visible only to founder/admin roles

**Trial calculator:**
- Landing page feature for unauthenticated users
- Guest session creation via `POST /auth/guest`
- Guest token stored separately in localStorage
- Rate limited (429 redirects to login)
- Full calculation and 3D visualization available
- Encourages signup for unlimited access

**State management:**
- React Context API (no external library)
- AuthContext: user, permissions, workspace, login, logout, switch
- PlanningContext: legacy local planning state (unused)
- StorageContext, ExecutionContext, AuditContext: stub implementations

### Mobile Client (Flutter)

**Architecture pattern:**
```
UI (Pages) → Providers (ChangeNotifier) → Services (API calls) → API
                                              ↓
                                        Storage Service
                                    (flutter_secure_storage)

Data transformation:
API JSON → DTO (json_serializable) → Mapper → Model (Freezed) → UI
```

**Dependency injection:**
- All services instantiated in `main.dart`
- Provided via `MultiProvider` at app root
- Constructor injection pattern (testable)

**WebView approach for 3D:**
- Does NOT implement native Three.js in Dart/Flutter
- Embeds Next.js web client's 3D viewer in WebView
- URL: `{webBaseUrl}/embed/shipments/{id}?token={accessToken}`
- Full support: Android (Chrome WebView), iOS (WKWebView), Web (iframe)
- Partial support: macOS (WKWebView), Windows (WebView2)
- Fallback: Linux shows "Open in Browser" button (no WebView support)
- Rationale: reuses existing web implementation, single source of truth, automatic updates

**Platform support:**
- Android/iOS: Production-ready, full features
- Linux: Limited (WebView fallback for 3D)
- macOS/Windows: Partial (WebView may have limitations)
- Web: Full support (Flutter web build)

**Code generation:**
- Freezed: generates immutable models with copyWith, equality, toString
- json_serializable: generates fromJson/toJson for DTOs
- Command: `flutter pub run build_runner build --delete-conflicting-outputs`

## Authentication & Authorization

### JWT Tokens

**Access token:**
- Expiry: 2 hours
- Claims: `user_id`, `username`, `role`, `workspace_id`
- Algorithm: HS256 (HMAC-SHA256)
- Header format: `Authorization: Bearer <token>`

**Refresh token:**
- Expiry: 30 days
- Format: `rf_YYYYMMDD_HHMMSS.microsec_random12`
- Stored in database (`refresh_tokens` table) with revocation tracking
- Rotated on each use (old token revoked, new token issued)

**Guest token:**
- Expiry: 30 days (longer than regular access tokens)
- Special access token for trial calculator users
- No refresh token provided
- Rate limited to encourage signup

**Token lifecycle:**
1. Login returns access + refresh tokens
2. Client stores tokens (localStorage for web, secure storage for mobile)
3. Access token included in all API requests
4. On 401, client calls `/auth/refresh` with refresh token
5. New access token returned, old refresh token revoked
6. Original request retried with new access token
7. If refresh fails (invalid/expired), user logged out

### Role System (7 Roles)

Defined in `cmd/db/scripts/default_user_permission.sql`:

| Role | Description | Key Permissions | Use Case |
|------|-------------|-----------------|----------|
| `founder` | Platform superuser (SaaS founder) | `*` (global wildcard, all permissions) | Platform owner, ultimate admin |
| `owner` | Workspace owner/CEO | Full workspace access: workspace, member, invite, product, container, plan management | Organization leader |
| `personal` | Personal workspace owner | Product, container, plan management (no member/invite access) | Individual user workspace |
| `admin` | Workspace administrator | Member/invite management, product, container, plan management | Team administrator |
| `planner` | Load planner | Full plan access, read-only product/container | Warehouse planner |
| `operator` | Loading validator | Read plans, manage plan items, read-only product/container | Warehouse operator |
| `trial` | Anonymous trial user | Create plans (max 3), read-only product/container, no workspace | Guest/trial user |

### Permission System

**Format:** `resource:action`

**Examples:**
- `plan:create` - Create new plans
- `plan:read` - View plans
- `plan:*` - All plan actions (wildcard)
- `product:read` - View products
- `user:*` - All user operations
- `*` - Global admin (founder only)

**Matching logic:**
- Exact match: `plan:read` matches `plan:read`
- Resource wildcard: `plan:*` matches `plan:read`, `plan:create`, `plan:update`, `plan:delete`
- Global wildcard: `*` matches all permissions
- Implementation: `lib/permissions.ts` (web), `internal/middleware/permission_middleware.go` (backend)

**Permission caching:**
- Backend: in-memory cache by role (prevents repeated DB queries)
- Web: stored in AuthContext after login
- Mobile: stored in AuthProvider after login

**Enforcement:**
- Backend: middleware checks before handler execution
- Web: RouteGuard component, navigation filtering
- Mobile: no UI enforcement (relies on backend rejection)

### Workspace Switching

Multi-tenant users (e.g., owner of multiple organizations) can switch workspaces:

1. User selects workspace from dropdown
2. Frontend calls `POST /auth/switch-workspace` with target `workspace_id`
3. Backend validates user membership in target workspace
4. Backend generates new access + refresh tokens with `workspace_id` in claims
5. Frontend stores new tokens, updates AuthContext
6. All subsequent requests scoped to new workspace

## Core Features

### Implemented Features

**Authentication:**
- User registration with automatic personal workspace creation
- Login with JWT access + refresh tokens
- Token refresh flow (automatic on 401 responses)
- Guest sessions for trial calculator
- Session persistence across browser/app restarts
- Logout with token revocation

**Multi-Tenant Workspaces:**
- Create and manage workspaces
- Invite members via email
- Accept invitations with role assignment
- Workspace switching for multi-tenant users
- Data isolation (all queries scoped by workspace_id)

**Shipments (Load Plans):**
- Create shipments with container selection
- Add items from product catalog or manual entry
- Calculate 3D placements with algorithm selection
- View shipment details with 3D visualization
- Step-by-step animation playback
- Hover tooltips showing item details
- PDF report generation with screenshots
- Delete shipments
- Status tracking (DRAFT, IN_PROGRESS, COMPLETED, FAILED, PARTIAL, CANCELLED)

**3D Visualization:**
- Web: Native Three.js with custom StuffingVisualizer engine
- Mobile: WebView embedding web client's 3D viewer
- Both: Step-by-step playback, hover detection, auto-fit camera
- PDF export with vector graphics and item tables
- Embeddable viewer at `/embed/shipments/[id]` with token auth

**Master Data:**
- Product catalog (name, dimensions in mm, weight in kg, color)
- Container profiles (name, dimensions, max weight)
- Full CRUD operations via API and UIs

**Platform Administration:**
- User management (list, create, update, delete)
- Workspace management across all tenants
- Role management (create, edit, assign permissions)
- Permission management (list available permissions)
- Platform-level access control (founder and admin roles only)

**Dashboard:**
- Role-specific statistics (total plans, products, containers, etc.)
- Recent shipments list
- Quick actions based on permissions
- Utilization metrics

**Trial Calculator (Web only):**
- Fully functional load calculator on landing page
- Guest session creation without registration
- Preset containers or custom dimensions
- Product catalog selection or manual entry
- Real-time calculation and 3D visualization
- Rate limiting to encourage signup

### Incomplete/Planned Features

**Reports Module (scaffolded, not functional):**
- Routes exist in web: `/reports/audit`, `/reports/execution`, `/reports/manifest`
- Pages created but no data integration
- Backend endpoints not implemented
- Expected features: audit logs, execution logs, loading manifests

**Loading Module (routes exist, not implemented):**
- Routes in web: `/loading`, `/loading/[sessionId]`
- Intended for warehouse operator step-by-step guidance
- No implementation in backend or frontend

**Real-time Features (not started):**
- IoT sensor integration for weight validation
- WebSocket for live updates during loading
- Execution validator with tolerance checks
- MQTT ingestion for sensor data
- These appear in architecture diagrams (`docs/architecture.dot`) but not in code

**Mobile Features:**
- Profile page incomplete (no user info display, no logout button)
- Token refresh not implemented (requires manual re-login after 2 hours)
- QR code scanning dependency installed but no UI integration
- No offline mode (network required for all operations)

**Testing Gaps:**
- Mobile: zero test coverage (empty test directory)
- Web: only service layer tests (8 files), no component/E2E tests
- Backend: good coverage (95.6%) but could use more integration tests

## API Endpoints

All endpoints prefixed with `/api/v1`:

**Authentication:**
- `POST /auth/register` - Create account with workspace
- `POST /auth/login` - Login (returns access + refresh tokens)
- `POST /auth/refresh` - Refresh access token
- `POST /auth/guest` - Create guest session (trial calculator)
- `GET /auth/me` - Get current user + permissions
- `POST /auth/switch-workspace` - Switch to different workspace
- `POST /auth/logout` - Logout (revoke tokens)

**Plans (Shipments):**
- `GET /plans` - List plans (filtered by workspace)
- `POST /plans` - Create plan (optionally with items + auto-calculate)
- `GET /plans/:id` - Get plan detail (includes items, calculation, placements)
- `PUT /plans/:id` - Update plan
- `DELETE /plans/:id` - Delete plan
- `POST /plans/:id/calculate` - Calculate 3D placements
- `POST /plans/:id/items` - Add item to plan
- `PUT /plans/:id/items/:itemId` - Update item
- `DELETE /plans/:id/items/:itemId` - Delete item

**Products:**
- `GET /products` - List products (workspace scoped)
- `POST /products` - Create product
- `GET /products/:id` - Get product detail
- `PUT /products/:id` - Update product
- `DELETE /products/:id` - Delete product

**Containers:**
- `GET /containers` - List containers (workspace scoped)
- `POST /containers` - Create container
- `GET /containers/:id` - Get container detail
- `PUT /containers/:id` - Update container
- `DELETE /containers/:id` - Delete container

**Users (platform admin only):**
- `GET /users` - List all users
- `POST /users` - Create user
- `GET /users/:id` - Get user detail
- `PUT /users/:id` - Update user
- `DELETE /users/:id` - Delete user

**Roles (platform admin only):**
- `GET /roles` - List all roles
- `POST /roles` - Create role
- `GET /roles/:id` - Get role detail
- `PUT /roles/:id` - Update role
- `DELETE /roles/:id` - Delete role
- `POST /roles/:id/permissions` - Assign permissions to role

**Permissions (platform admin only):**
- `GET /permissions` - List all available permissions
- `POST /permissions` - Create permission
- `GET /permissions/:id` - Get permission detail
- `PUT /permissions/:id` - Update permission
- `DELETE /permissions/:id` - Delete permission

**Workspaces (platform admin only):**
- `GET /workspaces` - List all workspaces
- `POST /workspaces` - Create workspace
- `GET /workspaces/:id` - Get workspace detail
- `PUT /workspaces/:id` - Update workspace
- `DELETE /workspaces/:id` - Delete workspace

**Members (workspace scoped):**
- `GET /members` - List workspace members
- `POST /members` - Add member to workspace
- `GET /members/:id` - Get member detail
- `PUT /members/:id` - Update member (change role)
- `DELETE /members/:id` - Remove member from workspace

**Invites (workspace scoped):**
- `GET /invites` - List workspace invitations
- `POST /invites` - Create invitation
- `GET /invites/:id` - Get invitation detail
- `DELETE /invites/:id` - Delete/cancel invitation
- `POST /invites/accept` - Accept invitation

**Dashboard:**
- `GET /dashboard` - Get statistics (role-based)

**Health:**
- `GET /health` - Health check

## Database

**Technology:** PostgreSQL 14+

**Migration system:** Goose
- Migrations: `cmd/db/migrations/*.sql`
- Commands:
  - `goose -dir cmd/db/migrations postgres "$DATABASE_URL" up` - Apply migrations
  - `goose -dir cmd/db/migrations postgres "$DATABASE_URL" down` - Rollback one
  - `goose -dir cmd/db/migrations postgres "$DATABASE_URL" status` - Check status
  - `goose -dir cmd/db/migrations create name sql` - Create new migration

**Code generation:** SQLC
- Configuration: `sqlc.yaml`
- Query definitions: `cmd/db/queries/*.sql`
- Generated code: `internal/store/*.sql.go`
- Command: `sqlc generate`
- Type-safe Go structs and query methods

**Schema:** See `docs/erd.dot` for entity relationship diagram

**Key tables:**
- `users` - User accounts
- `user_profiles` - User details (name, email, etc.)
- `roles` - Role definitions
- `permissions` - Permission definitions
- `role_permissions` - Role-permission mapping (many-to-many)
- `refresh_tokens` - Active refresh tokens with expiry tracking
- `workspaces` - Tenant containers (organization or personal)
- `members` - Workspace membership with role assignment
- `invites` - Workspace invitation records
- `platform_members` - Founder access across all workspaces
- `products` - Product catalog (workspace scoped)
- `containers` - Container profiles (workspace scoped)
- `load_plans` - Plan/shipment header with status
- `load_items` - Items to pack in plan
- `plan_results` - Calculation summary (total items, weight, utilization)
- `plan_placements` - Individual item placements (position, rotation, step)

**Indexes:**
- Foreign key columns
- Workspace ID for scoped queries
- User email for login lookups
- Plan status for filtering
- Created/updated timestamps for sorting

**Constraints:**
- Foreign keys with cascade delete where appropriate
- Unique constraints on email, workspace names
- Check constraints on dimensions (positive values)
- Not null constraints on required fields

## 3D Visualization

**Two implementations:**

### Web (Native Three.js)

**StuffingVisualizer engine:**
- Location: `web/lib/StuffingVisualizer/` (1731 lines, 14 files)
- Main class: `StuffingVisualizer` (facade pattern)
- Manager pattern separates concerns:
  - CameraManager: OrthographicCamera positioning, auto-fit
  - RendererManager: WebGL renderer lifecycle
  - SceneManager: Scene graph management
  - LightManager: Ambient + directional lighting
  - ControlsManager: OrbitControls integration
  - AnimationManager: Step-by-step playback (133 lines)
  - InteractionManager: Hover detection via raycasting (100 lines)

**Features:**
- Step-by-step animation (play, pause, reset, step forward/backward)
- Hover tooltips showing item name, dimensions, weight
- Auto-fit camera to container size
- Screenshot capture (1920x1080 for reports)
- PDF generation with jsPDF (580 lines)
- Event system (onStepChange, onPlayStateChange, onItemHover)
- Proper disposal (prevents memory leaks)

**Usage:**
- Shipment detail page: full viewer with controls
- Embed viewer: token-secured embeddable component
- Trial calculator: inline visualization

### Mobile (WebView Embed)

**Approach:**
- Does NOT implement Three.js natively in Dart
- Embeds web client's 3D viewer in WebView component
- URL: `{webBaseUrl}/embed/shipments/{id}?token={accessToken}`
- Component: `PlanVisualizerView` (lib/components/viewers/plan_visualizer_view.dart)
- Library: `webview_flutter` v4.13.1

**Platform support:**
- Android: Full support (Chrome WebView)
- iOS: Full support (WKWebView)
- Linux: Fallback to external browser (no WebView)
- macOS/Windows: Partial support (may have limitations)

**Rationale:**
- Reuses existing web implementation (single source of truth)
- Avoids porting 1731 lines of Three.js code to Dart
- Automatic updates when web client improves
- Reduces mobile app complexity and maintenance
- Tradeoff: requires internet connection, desktop limitations

## Testing

**Backend (Go):**
- Overall coverage: 95.6%
- Package coverage:
  - cache: 100%
  - config: 100%
  - env: 100%
  - middleware: 100%
  - response: 100%
  - types: 100%
  - handler: 97.8%
  - packer: 96.5%
  - gateway: 94.1%
  - service: 94.3%
  - auth: 93.6%
- Pattern: table-driven tests with `t.Run()`
- Mocks: `go.uber.org/mock` for SQLC store
- Handler tests: `httptest` + `testify/assert`
- Commands:
  - `make test` - Run all tests
  - `make coverage` - Generate coverage report
  - `make coverage-html` - HTML coverage report
  - `go test ./internal/service -run '^TestName$' -v` - Run specific test

**Web (Next.js):**
- Framework: Vitest 4.0.16
- Coverage: Service layer only (8 test files)
- No component tests
- No integration/E2E tests
- Location: `web/test/lib/services/`
- Commands:
  - `pnpm -C web test` - Run all tests
  - `pnpm -C web test -- -t "test name"` - Run specific test
  - `pnpm -C web test -- path/to/file.test.ts` - Run file

**Mobile (Flutter):**
- Coverage: 0% (no tests implemented)
- Test directory empty: `mobile/test/`
- Recommended: Unit tests for services/mappers, widget tests, integration tests
- Commands (when implemented):
  - `flutter test` - Run all tests
  - `flutter test --coverage` - Generate coverage

## Development Workflows

### Backend (Go)

**Setup:**
```bash
# Install dependencies
go mod tidy

# Install tools
go install github.com/pressly/goose/v3/cmd/goose@latest
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
go install github.com/swaggo/swag/cmd/swag@latest

# Setup database
createdb stuffing
export DATABASE_URL="postgres://postgres:password@localhost:5432/stuffing?sslmode=disable"

# Run migrations
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up

# Generate SQLC code
sqlc generate

# Generate Swagger docs
swag init -g cmd/api/main.go -o internal/docs
```

**Development:**
```bash
make fmt          # Format code
make swag         # Generate Swagger docs
sqlc generate     # Generate DB code (after query changes)
make test         # Run tests
make coverage     # Test coverage
make build        # Build binary (outputs bin/api)
make run          # Run server (localhost:8080)
```

**Database workflow:**
```bash
# Create new migration
goose -dir cmd/db/migrations create migration_name sql

# Check migration status
goose -dir cmd/db/migrations postgres "$DATABASE_URL" status

# Apply migrations
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up

# Rollback one migration
goose -dir cmd/db/migrations postgres "$DATABASE_URL" down
```

**SQLC workflow:**
```bash
# 1. Add/modify query in cmd/db/queries/*.sql
# 2. Regenerate code
sqlc generate

# 3. New Go methods available in internal/store/*.sql.go
```

### Web (Next.js)

**Setup:**
```bash
cd web
pnpm install

# Configure environment (optional)
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > .env.local
```

**Development:**
```bash
pnpm dev          # Start dev server (localhost:3000)
pnpm build        # Production build
pnpm start        # Start production server
pnpm lint         # Run ESLint
pnpm test         # Run Vitest tests
```

**Adding features:**
```bash
# Add shadcn/ui component
pnpm dlx shadcn@latest add [component-name]

# Add new page: create file in app/(app)/feature/page.tsx
# Add new service: create file in lib/services/feature.ts
# Add new hook: create file in hooks/use-feature.ts
# Add types: create file in lib/types/feature.ts
```

### Mobile (Flutter)

**Setup:**
```bash
cd mobile
flutter pub get

# Generate code (Freezed + json_serializable)
flutter pub run build_runner build --delete-conflicting-outputs

# Check devices
flutter devices
```

**Development:**
```bash
# Run on device
flutter run -d <device_id>

# Run on Android
flutter run -d android

# Run on iOS
flutter run -d ios

# Run on desktop
flutter run -d linux  # or macos, windows

# Code generation (after DTO/Model changes)
flutter pub run build_runner build --delete-conflicting-outputs

# Watch mode (auto-regenerate)
flutter pub run build_runner watch --delete-conflicting-outputs
```

**Build production:**
```bash
# Android
flutter build apk --release
flutter build appbundle --release

# iOS
flutter build ios --release

# Desktop
flutter build linux --release
flutter build macos --release
flutter build windows --release
```

## Environment Configuration

### Backend (.env)

```bash
# Server
SRV_ENV=dev                           # dev or production
SRV_PORT=8080

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/stuffing?sslmode=disable

# JWT
JWT_SECRET=your-secret-key            # Use: openssl rand -base64 32
JWT_ACCESS_EXPIRE=2h
JWT_REFRESH_EXPIRE=720h              # 30 days

# Packing Service
PACKING_SERVICE_URL=http://localhost:5051

# Default Admin
FOUNDER_EMAIL=admin@example.com
FOUNDER_PASSWORD=changeme
```

### Web (.env.local)

```bash
# API URL (default: http://localhost:8080/api/v1)
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Mobile (lib/config/constants.dart)

```dart
static const apiBaseUrl = String.fromEnvironment(
  'API_URL',
  defaultValue: 'https://stuffing-api.irc-enter.tech/api/v1',
);

// Or override at runtime:
// flutter run --dart-define=API_URL=http://192.168.1.100:8080/api/v1
```

## Code Style & Conventions

### Go

- Format: `gofmt` (run via `make fmt`)
- Imports: grouped by stdlib, third-party, internal
- Errors: wrap with context using `fmt.Errorf("context: %w", err)`
- Avoid panics except in initialization
- HTTP handlers: validate with `c.ShouldBindJSON`, use `response.Success`/`response.Error`
- Tests: table-driven with `t.Run("case_name", ...)`
- Comments: godoc for exported functions

### TypeScript/React (Web)

- Strict mode: enabled in `tsconfig.json`
- Imports: prefer `import type ...` for type-only imports
- Naming: components `PascalCase`, functions/variables `camelCase`
- Components: functional components with hooks (no class components)
- Client components: add `"use client"` directive for interactive components
- Formatting: ESLint + platform defaults (no Prettier config)

### Dart/Flutter (Mobile)

- Format: `dart format lib/`
- Naming: files `snake_case`, classes `PascalCase`, variables `camelCase`
- Models: use Freezed for immutability
- DTOs: use json_serializable for JSON mapping
- State: ChangeNotifier pattern with Provider
- Tests: table-driven when appropriate (not implemented yet)

## Key Files Reference

**Backend:**
- `cmd/api/main.go` - Entrypoint
- `internal/api/routes.go` - Route definitions
- `internal/service/plan_service.go` - Plan calculation logic
- `internal/packer/packer.go` - Bin packing wrapper
- `internal/middleware/permission_middleware.go` - RBAC enforcement
- `internal/auth/jwt.go` - JWT generation/validation
- `cmd/db/queries/plan.sql` - Plan queries for SQLC
- `cmd/db/migrations/*.sql` - Database migrations

**Web:**
- `app/(app)/shipments/[id]/page.tsx` - Shipment detail with 3D viewer
- `components/stuffing-viewer.tsx` - 3D viewer component wrapper
- `lib/StuffingVisualizer/stuffing-visualizer.ts` - Main 3D engine class
- `lib/api.ts` - HTTP client with token refresh
- `lib/auth-context.tsx` - Authentication context provider
- `lib/route-guard.tsx` - Route protection component
- `lib/services/plans.ts` - Plan API service
- `components/trial-load-calculator.tsx` - Landing page calculator (739 lines)

**Mobile:**
- `lib/main.dart` - App entry point + provider setup
- `lib/pages/plans/plan_detail_page.dart` - Plan detail with 3D viewer
- `lib/components/viewers/plan_visualizer_view.dart` - WebView wrapper
- `lib/services/api_service.dart` - HTTP client with interceptor
- `lib/providers/auth_provider.dart` - Auth state management
- `lib/services/plan_service.dart` - Plan API service
- `lib/config/routes.dart` - GoRouter configuration

## Design Decisions

### Why JWT-embedded workspace_id instead of header?

**Decision:** Embed `workspace_id` in JWT claims instead of custom header (`X-Workspace-ID`)

**Rationale:**
- Simpler client implementation (no header management)
- Prevents workspace_id/token mismatch (single source of truth)
- Backend can trust workspace context (validated during JWT generation)
- Workspace switching requires new token (more secure)
- Eliminates need for middleware to check header vs token consistency

**Tradeoff:** Requires new JWT when switching workspaces (acceptable cost for security)

### Why SQLC over ORM?

**Decision:** Use SQLC for type-safe SQL queries instead of ORM (GORM, ent, etc.)

**Rationale:**
- Full SQL control (no ORM abstraction leaks)
- Type safety without code generation magic
- Predictable queries (no hidden N+1 problems)
- Better performance (direct SQL, no query builder overhead)
- Easier debugging (see exact SQL in query files)
- Migration flexibility (Goose + raw SQL)

**Tradeoff:** More verbose than ORM DSL, requires SQL knowledge

### Why WebView for mobile 3D instead of native?

**Decision:** Embed web client's Three.js viewer in WebView instead of implementing native 3D in Dart/Flutter

**Rationale:**
- Reuses existing implementation (1731 lines of tested code)
- Single source of truth for 3D visualization
- Automatic updates when web improves
- Avoids porting Three.js concepts to Dart
- Reduces mobile maintenance burden
- Faster time to market

**Tradeoff:**
- Requires internet connection (no offline 3D)
- Desktop platform limitations (Linux fallback, macOS/Windows partial)
- Slightly slower initial load (WebView startup)
- Less control over 3D performance

**Alternative considered:** Flutter 3D libraries (flutter_cube, vector_math)
- Rejected: immature ecosystem, would require reimplementation, ongoing maintenance

### Why React Context instead of Redux/Zustand?

**Decision:** Use React Context API for global state instead of external libraries

**Rationale:**
- Simple state needs (auth, few global contexts)
- No complex state transitions
- Built-in to React (no additional dependency)
- Good enough for current scale
- Can migrate later if needed

**Tradeoff:** Less tooling, more verbose than Zustand, potential re-renders if not optimized

### Why service layer pattern?

**Decision:** Introduce service layer between components and API client in web/mobile

**Rationale:**
- Consistent API abstraction across codebase
- Type-safe interfaces
- Easy to mock for testing
- Centralized error handling
- Business logic separation
- Can add caching/validation in one place

**Tradeoff:** More files, more boilerplate (acceptable for maintainability)

## Deployment

### Docker Compose (Development)

```bash
# Start all services
docker compose up --build

# Services:
# - db: PostgreSQL on 5432
# - api: Go API on 8080
# - packing: Python service on 5051
# - web: Next.js on 3000

# Access:
# Web: http://localhost:3000
# API: http://localhost:8080
# Swagger: http://localhost:8080/docs/index.html
# Login: admin@example.com / admin123

# Stop services
docker compose down

# Clean volumes
docker compose down -v
```

### Production Deployment

**Backend (Go API):**
```bash
# Build binary
make build

# Binary at: bin/api

# Run with environment variables
export DATABASE_URL="postgres://..."
export JWT_SECRET="..."
./bin/api

# Or use systemd service
```

**Web (Next.js):**
```bash
cd web

# Build production
pnpm build

# Start production server
pnpm start

# Or export static (if applicable)
pnpm build && pnpm export
```

**Mobile:**
```bash
# Android
flutter build appbundle --release
# Upload to Google Play Console

# iOS
flutter build ios --release
# Open Xcode, archive, upload to App Store Connect
```

**Database:**
```bash
# Run migrations in production
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up

# Backup
pg_dump -U postgres stuffing > backup.sql

# Restore
psql -U postgres stuffing < backup.sql
```

**Environment variables in production:**
- Use secret management (AWS Secrets Manager, HashiCorp Vault, etc.)
- Never commit secrets to git
- Set `SRV_ENV=production`
- Use strong `JWT_SECRET` (32+ random bytes)
- Configure CORS for production domains

**Reverse proxy (nginx example):**
```nginx
server {
    listen 80;
    server_name example.com;

    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
    }
}
```

**Monitoring:**
- API health check: `GET /api/v1/health`
- Database connection: check API startup logs
- Application metrics: consider adding Prometheus
- Error tracking: consider Sentry or similar

## Troubleshooting

### Backend won't start

**Check:**
- Database is running: `pg_isready -h localhost -p 5432`
- DATABASE_URL is correct
- Migrations applied: `goose -dir cmd/db/migrations postgres "$DATABASE_URL" status`
- JWT_SECRET is set
- Port 8080 not in use: `lsof -ti:8080`

**Solution:**
```bash
# Start database
docker compose up db

# Apply migrations
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up

# Run API
make run
```

### Web client can't connect to API

**Check:**
- API is running: `curl http://localhost:8080/api/v1/health`
- NEXT_PUBLIC_API_URL is correct (check .env.local)
- CORS is configured in backend (check `internal/api/api.go`)
- Browser console for errors

**Solution:**
```bash
# Web: check/set API URL
echo "NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1" > web/.env.local

# Backend: ensure CORS allows frontend origin
```

### Mobile can't connect to API

**Check:**
- API is accessible from mobile device
- Android emulator: use `10.0.2.2` instead of `localhost`
- iOS simulator: use computer's local IP
- Physical device: ensure same network or use public API URL

**Solution:**
```dart
// mobile/lib/config/constants.dart
static const apiBaseUrl = 'http://192.168.1.100:8080/api/v1'; // Your local IP
```

Or:
```bash
flutter run --dart-define=API_URL=http://192.168.1.100:8080/api/v1
```

### Token expired / authentication loops

**Check:**
- Access token expiry (default 2 hours)
- Refresh token exists in storage
- `/auth/refresh` endpoint works
- Clock sync between client/server

**Solution:**
```bash
# Web: clear localStorage and re-login
# Mobile: clear app data and re-login

# Backend: check JWT_SECRET hasn't changed
```

### 3D viewer blank/not loading

**Web:**
- Check browser console for Three.js errors
- Verify WebGL support: `document.createElement("canvas").getContext("webgl")`
- Check plan has placements (call `/plans/:id` API)
- Container dimensions valid (non-zero)

**Mobile:**
- Check platform supports WebView (Android/iOS yes, Linux no)
- Verify embed URL is correct
- Check token is valid
- Try opening embed URL in mobile browser

### Database migration conflicts

**Check:**
- Migration files are in correct order (numbered)
- No conflicting migrations from multiple branches
- Database is in expected state

**Solution:**
```bash
# Check migration status
goose -dir cmd/db/migrations postgres "$DATABASE_URL" status

# Rollback to known good state
goose -dir cmd/db/migrations postgres "$DATABASE_URL" down

# Reapply migrations
goose -dir cmd/db/migrations postgres "$DATABASE_URL" up
```

### SQLC generation errors

**Check:**
- Query syntax is valid PostgreSQL
- sqlc.yaml configuration correct
- Table/column names match database

**Solution:**
```bash
# Verify query syntax
psql -d stuffing -c "EXPLAIN SELECT ..."

# Regenerate
sqlc generate

# Check errors in output
```

### Build failures

**Go:**
```bash
go mod tidy        # Update dependencies
go clean -cache    # Clear build cache
make build         # Rebuild
```

**Web:**
```bash
rm -rf node_modules .next
pnpm install
pnpm build
```

**Mobile:**
```bash
flutter clean
flutter pub get
flutter pub run build_runner build --delete-conflicting-outputs
flutter build apk
```

### Rate limit on trial calculator

**Expected behavior:** Guests are rate limited to encourage signup

**Solutions:**
- Sign up for account (unlimited access)
- Wait for rate limit reset (typically 24 hours)
- Backend controls rate limit per IP

### Permission denied errors

**Check:**
- User role has required permission
- Check `role_permissions` table
- Permission format correct (`resource:action`)
- Backend logs for permission check failures

**Solution:**
```sql
-- Check user's role
SELECT role FROM users WHERE id = 'user-id';

-- Check role's permissions
SELECT p.resource, p.action 
FROM permissions p 
JOIN role_permissions rp ON p.id = rp.permission_id 
JOIN roles r ON r.id = rp.role_id 
WHERE r.name = 'planner';

-- Add missing permission
INSERT INTO role_permissions (role_id, permission_id) 
VALUES ('role-id', 'permission-id');
```
