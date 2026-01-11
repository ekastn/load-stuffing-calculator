# Web UI Permission Flow Gaps (RouteGuard / Permissions)

Date: 2025-12-30

This file captures investigation findings only (no fixes applied yet).

## Summary

The web UI has been partially migrated to permission-based access control via `RouteGuard requiredPermissions`, but several pages either:

- have **no RouteGuard at all** (page access not permission gated), or
- are guarded for **read** permissions while still rendering **write** UI actions that should be gated by create/update/delete permissions.

This creates inconsistent UX (users see buttons that will fail) and increases the chance of missing authorization checks during future refactors.

## High Impact

### 1) Shipment detail page missing RouteGuard 

- File: `web/app/(app)/shipments/[id]/page.tsx`
- Behavior: No `RouteGuard` wrapper.
- Risks:
  - A user who can reach the route can trigger plan fetches and mutations until the API denies them.
  - UI logic uses `user.role` (planner/admin) to show actions instead of using permissions.

Recommended:
- Wrap the page with `<RouteGuard requiredPermissions={["plan:read"]}>`.
- Gate action buttons using permissions:
  - Recalculate / advanced calculation: `plan:calculate`
  - Add items / delete items: `plan_item:*` (or granular `plan_item:create/delete` if API is changed)

Backend alignment note:
- API uses `plan_item:*` for add/update/delete item routes (`internal/api/routes.go`). If the UI gates with granular permissions, the API should also enforce granular `plan_item:create/update/delete`.

## Medium Impact

### 2) Read-guarded pages render write actions without permission checks

These pages have a read-level `RouteGuard`, but render mutations without checking create/update/delete permissions.

#### Members

- File: `web/app/(app)/settings/members/page.tsx`
- Guard: `member:read`
- Actions shown:
  - Add member
  - Update member role
  - Remove member
- Recommended UI gating:
  - Add member: `member:create`
  - Update role: `member:update`
  - Remove member: `member:delete`

#### Invites

- File: `web/app/(app)/settings/invites/page.tsx`
- Guard: `invite:read`
- Actions shown:
  - New invite
  - Revoke invite
- Recommended UI gating:
  - New invite: `invite:create`
  - Revoke invite: `invite:delete`

#### Products

- File: `web/app/(app)/products/page.tsx`
- Guard: `product:read`
- Actions shown:
  - New product
  - Edit
  - Delete
- Recommended UI gating:
  - New: `product:create`
  - Edit: `product:update`
  - Delete: `product:delete`

#### Containers

- File: `web/app/(app)/containers/page.tsx`
- Guard: `container:read`
- Actions shown:
  - New container
  - Edit
  - Delete
- Recommended UI gating:
  - New: `container:create`
  - Edit: `container:update`
  - Delete: `container:delete`

### 3) Remaining role-based gating (should prefer permissions)

The UI still uses role strings (e.g. planner/admin) for some behavior that should be aligned with permissions.

- File: `web/app/(app)/shipments/page.tsx`
  - “Create Shipment” button is shown only for `role === "planner"`.
  - Recommended: show it if the user has `plan:create`.

- File: `web/app/(app)/shipments/[id]/page.tsx`
  - Planner/admin role gates actions.
  - Recommended: gate based on permissions as described above.

## Low Impact / Cleanup

### 4) Duplicated permission helper implementations

The same helper logic exists in multiple places:

- `web/lib/route-guard.tsx` has `permissionMatches` + `hasAnyPermission`.
- `web/components/dashboard-layout.tsx` duplicates the same logic.

Recommended:
- Move shared helpers to something like `web/lib/permissions.ts` and import both places.

### 5) RouteGuard debug logging

- File: `web/lib/route-guard.tsx`
- Contains `console.log("[v0] RouteGuard ...")` statements.

Recommended:
- Remove or guard behind a dev-only flag.

## Suggested Fix Approach (not implemented)

1. Add `RouteGuard requiredPermissions={["plan:read"]}` to shipment detail page.
2. Introduce permission checks for action buttons on the pages above.
3. Optional: refactor duplicated permission helper logic into a shared module.
4. Validate via `pnpm -C web build` and `pnpm -C web test`.
