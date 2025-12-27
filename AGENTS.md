# Agent Notes (load-stuffing-calculator)

## Project Overview
This repo implements a Load Planning / Container Stuffing system:
- Admin manages users/roles/permissions and master data (products, containers).
- Planner creates load plans (shipments), adds items, and runs 3D bin packing.
- Operator can read plans/results and replay loading steps in the UI.

Primary output of the packing engine is a placement list (position + rotation + step number), persisted as `plan_results` and `plan_placements`, and visualized in the `web/` Next.js UI.

## What’s Implemented vs Planned
Implemented (in code):
- REST API for auth + RBAC-protected CRUD + plan calculation.
- 3D packing calculation via `internal/packer` (boxpacker3).
- `web/` Next.js UI that calls the API and renders 3D placements with step playback.

Planned / diagram-only (not implemented in Go code yet):
- IoT ingestion (MQTT), realtime updates (WebSocket hub), execution validator + tolerance checks, execution logs.
- These appear in `docs/architecture.dot`, `docs/erd.dot`, `docs/usecase.plantuml` but do not currently exist in Go backend code.

## Repo Layout
- Go API (Gin): `cmd/api` (entrypoint), `internal/*` (app code)
- Web UI (primary): `web/` (Next.js + TypeScript + Vitest)
- Additional UI/prototype: `frontend/` (Vite + React)
- DB migrations (Goose): `cmd/db/migrations`
- SQLC config + generated code: `sqlc.yaml`, `internal/store/*.sql.go`

## Architecture (Implemented)
Backend layering (typical request path):
- Routes: `internal/api/routes.go`
- Middleware: `internal/middleware/*` (JWT auth + role checks)
- Handlers: `internal/handler/*` (Gin handlers; bind/validate input, call services)
- Services: `internal/service/*` (business logic; DB writes/reads; packing calculation)
- Store: `internal/store/*` (SQLC generated queries + models; bulk placement insert via CopyFrom)

Packing engine:
- `internal/packer/packer.go` wraps `github.com/bavix/boxpacker3`.
- Supports strategy/goal/gravity options (see `internal/packer/types.go`).
- Units: dimensions are in millimeters; weights passed in kg and converted to grams for boxpacker3.

Data model (implemented tables):
- Auth/RBAC: `users`, `roles`, `permissions`, `role_permissions`, `refresh_tokens`
- Master data: `products`, `containers`
- Planning/calculation:
  - `load_plans` (status: DRAFT/IN_PROGRESS/COMPLETED/FAILED/PARTIAL/CANCELLED)
  - `load_items`
  - `plan_results` (summary stats)
  - `plan_placements` (pos_x/pos_y/pos_z, rotation_code, step_number)

Frontend architecture (`web/`):
- API client + token refresh: `web/lib/api.ts`
- Data hooks: `web/hooks/use-*.ts`
- 3D visualization / step playback: `web/app/shipments/[id]/page.tsx` + `web/components/stuffing-viewer.tsx` + `web/lib/StuffingVisualizer/*`

## Core Flows
Auth:
- Login: `POST /api/v1/auth/login`
- Refresh: `POST /api/v1/auth/refresh`
- Protected endpoints require `Authorization: Bearer <access_token>`.

Admin (role: `admin`):
- CRUD: `/api/v1/users`, `/api/v1/roles`, `/api/v1/permissions`, `/api/v1/products`, `/api/v1/containers`

Planning & Calculation:
- Read access: planner + operator (admin implicit)
- Write access: planner (admin implicit)
- Typical lifecycle:
  1) Create plan: `POST /api/v1/plans` (can include items; may auto-calculate)
  2) Manage items: `POST/PUT/DELETE /api/v1/plans/:id/items`
  3) Calculate placements: `POST /api/v1/plans/:id/calculate`
     - Service fetches plan + items, calls packer, saves `plan_results` and bulk inserts `plan_placements`.
  4) Fetch for UI: `GET /api/v1/plans/:id` includes items + calculation + placements (if result exists)

Frontend rendering (shipments detail):
- The page converts backend `placements[]` + item dimensions into `PackedItem[]` and replays steps by filtering by `step_number`.

## Key Files to Start From
- Routes: `internal/api/routes.go`
- Plan calculation + persistence: `internal/service/plan_service.go`
- Packing algorithm adaptor: `internal/packer/packer.go`
- SQLC models/queries: `internal/store/plan.sql.go`, `cmd/db/queries/*.sql`
- UI placement mapping + step playback: `web/app/shipments/[id]/page.tsx`
- Architecture/ERD/usecases (includes roadmap items): `docs/architecture.dot`, `docs/erd.dot`, `docs/usecase.plantuml`

## Build/Lint/Test
- Go fmt: `make fmt` (runs `go fmt ./...`)
- Go tests: `make test` (runs `go test -v ./...`)
- Single Go test: `go test ./... -run '^TestName$' -v` (or `go test ./internal/service -run '^TestAuthService_Login$' -v`)
- Build API binary: `make build` (outputs `bin/api`), run: `make run`
- Swagger: `make swagfmt` / `make swag` (requires `swag`)
- Web (Next): `pnpm -C web dev|build|lint|test` (npm works too)
- Single Vitest test: `pnpm -C web test -- -t "test name"` or `pnpm -C web test -- path/to/file.test.ts`

## Code Style / Conventions
- Go: keep `gofmt` output; imports grouped by gofmt (stdlib, third-party, internal).
- Prefer service methods to return wrapped errors: `fmt.Errorf("context: %w", err)`; avoid panics.
- HTTP handlers: validate with `c.ShouldBindJSON`, use `response.Success` / `response.Error`, return early on errors.
- Tests: use table-driven tests and `t.Run("case_name", ...)`; Gin handler tests use `httptest` + `testify/assert`.
- TS/React (`web/`): `strict: true`; prefer `import type ...`; components `PascalCase`, functions/vars `camelCase`.
- Follow existing file-local formatting (no Prettier config found); rely on ESLint + platform defaults.

## Cursor / Copilot Rules
- No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` found.
