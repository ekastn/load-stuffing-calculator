# SaaS Multi-Workspace Plan (Founder + Org + Personal)

## Goals

- Add SaaS-ready tenancy via **Workspaces**.
- Support two account usage modes with the same login:
  - **Personal workspace** (individual account)
  - **Organization workplaces** (multi-member)
- Support **workspace switching** in the Web UI by re-issuing JWT access tokens.
- Keep RBAC **non-advanced** (no org-custom roles):
  - Roles remain **global** in `roles`/`permissions`/`role_permissions`.
- Add a **Founder** “platform” capability that can do anything.
- Keep the public trial/guest flow working (role `trial`).

## Non-Goals

- Per-organization custom roles/permissions editor.
- Billing/quotas (beyond existing trial plan limit).
- Email delivery for invites (token response is enough for MVP).

---

## Terminology

- **User**: global identity (`users` table).
- **Workspace**: tenancy boundary. All master data + plans live inside a workspace.
- **Membership**: mapping `user <-> workspace` with a role.
- **Platform member**: mapping for platform-level actors (currently only `founder`).

---

## Core Invariants

### Workspace

- Every user has exactly **one personal** workspace.
- Organization workspaces can have multiple members.
- Every workspace has exactly one owner (`workspaces.owner_user_id`).

### Roles

- Roles remain global (seeded): `founder`, `owner`, `admin`, `planner`, `operator`, `trial`.
- `trial` remains for guest tokens; it is not a member role.

### Permissions

- RBAC is permission-based (existing middleware): role → permissions.
- **Founder** has global access.
- Workspace members are authorized by membership role + permission middleware.

### Founder behavior

- Founder can act on any workspace.
- Founder may optionally provide `workspace_id` to certain endpoints to operate on a different workspace without switching.

### Workspace scoping

- For non-founder requests, handlers/services MUST use the active `workspace_id` from JWT.
- For founder requests, handlers may accept `workspace_id` as an override.

---

## API (Flat Resources)

All endpoints live under `/api/v1`.

### Auth

- `POST /auth/login` (public)
- `POST /auth/register` (public)
- `POST /auth/guest` (public) → issues `trial` JWT
- `POST /auth/refresh` (public)
- `POST /auth/switch-workspace` (JWT required)
  - Request: `{ "workspace_id": "uuid" }`
  - Response: new access token (optionally refresh rotation).

### Workspaces

- `GET /workspaces` (JWT required)
  - Default: list only workspaces the user is a member of.
  - Founder: may list all (optional).
- `POST /workspaces` (JWT required)
  - Creates an **organization** workspace; creator becomes owner member.
- `PATCH /workspaces/:id` (JWT required)
  - Rename workspace
  - Transfer ownership by setting `{ "owner_user_id": "uuid" }` (owner/founder only)
- `DELETE /workspaces/:id` (JWT required)
  - Deletes workspace and ALL workspace-owned records (owner/founder only)

### Members (workspace memberships)

- `GET /members` (JWT required)
  - Lists members for active workspace.
  - Founder: may override with `?workspace_id=...`.
- `POST /members` (JWT required)
  - Adds an existing user to the target workspace.
  - Body: `{ "user_identifier": "email|username|uuid", "role": "admin|planner|operator" }`
- `PATCH /members/:member_id` (JWT required)
  - Updates member’s role (admin/promotions allowed).
- `DELETE /members/:member_id` (JWT required)
  - Removes a member.
  - Owner membership cannot be removed (409: transfer or delete workspace).

### Invites

- `GET /invites` (JWT required)
  - Lists invites for active workspace.
  - Founder: may override with `?workspace_id=...`.
- `POST /invites` (JWT required)
  - Creates an invite.
  - Body: `{ "email": "...", "role": "admin|planner|operator" }`
- `DELETE /invites/:invite_id` (JWT required)
  - Revokes invite.
- `POST /invites/accept` (JWT optional)
  - Body: `{ "token": "..." }`
  - If logged in and email matches, accept.
  - If not logged in: MVP can return 401/409 and require login/register.
  - Response should include a **new access token** for the invite workspace to switch the session.

---

## Data Model (Postgres)

### New tables

- `workspaces(workspace_id, type, name, owner_user_id, created_at, updated_at)`
- `members(member_id, workspace_id, user_id, role_id, created_at, updated_at)`
- `invites(invite_id, workspace_id, email, role_id, token_hash, invited_by_user_id, expires_at, accepted_at, revoked_at, created_at)`
- `platform_members(user_id, role_id, created_at)`

### Schema changes

- Add `workspace_id` to workspace-owned resources:
  - `containers.workspace_id`
  - `products.workspace_id`
  - `load_plans.workspace_id` (nullable to keep guest trial plans `workspace_id IS NULL`)
- Add `refresh_tokens.workspace_id` to remember active workspace for refresh.

### Cascades

- Workspace deletion must cascade to:
  - `members`, `invites`
  - all workspace-owned master data
  - all workspace-owned plans and associated rows

---

## Implementation Phases

### Phase 1: DB + Seeding

1. Add migration:
   - Create new tables.
   - Add workspace_id columns + FKs + cascades.
   - Add refresh_tokens.workspace_id.
2. Update `cmd/db/scripts/default_user_permission.sql`:
   - Add roles: `founder`, `owner`.
   - Add permissions for workspace/members/invites (if needed).
   - Map permissions:
     - `founder` should get `*`.
     - `owner` should get broad workspace permissions.
3. Backfill script:
   - Create personal workspace for existing users.
   - Create membership row for that personal workspace.

### Phase 2: Store (SQLC)

1. Add query files for workspaces/members/invites/platform members.
2. Update existing queries for containers/products/plans to scope by workspace_id.
3. Run `sqlc generate` to update `internal/store/*`.

### Phase 3: Auth + JWT

1. Add `workspace_id` claim to JWT.
2. Extend auth context to include workspace id.
3. Login/Register/Refresh:
   - Ensure personal workspace exists.
   - Determine active workspace id.
   - Determine role:
     - if platform_members has founder role → token role `founder`
     - else membership role in active workspace
4. Implement `POST /auth/switch-workspace`:
   - Update refresh token record workspace_id.
   - Return new access token.

### Phase 4: Tenant features

1. Implement handlers/services for:
   - Workspaces
   - Members
   - Invites
2. Add route wiring.
3. Enforce invariants:
   - admins cannot remove owner
   - invite cannot grant owner/founder
   - personal workspace blocks members/invites mutations

### Phase 5: Web

1. Add workspace switcher.
2. Add members/invites screens for org workspaces.
3. Implement invite accept to switch session (store returned JWT).

### Phase 6: Tests

- Add service tests:
  - membership enforcement
  - owner transfer
  - workspace delete cascade
  - founder workspace override behavior on members/invites

---

## Acceptance Criteria

- Users can belong to multiple org workspaces and switch via JWT.
- Individual users (personal workspace owner) can create products/containers/plans.
- Admins can manage members/invites in org workspace but cannot remove owner.
- Owner can transfer ownership or delete workspace (deleting everything).
- Founder can view/manage across workspaces and optionally override workspace_id on members/invites.
- Existing trial flow still works and trial limit enforcement remains.
