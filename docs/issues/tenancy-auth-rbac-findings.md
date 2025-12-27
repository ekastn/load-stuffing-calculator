# Tenancy / Auth / RBAC Findings

Date: 2025-12-27

This file captures investigation findings only (no fixes applied yet).

## Summary

### High impact

1) **Plan item endpoints are not workspace-scoped for non-trial users**
- Impact: cross-workspace plan item modification is possible if an attacker knows a `plan_id`.
- Root cause: plan item store queries (`load_items`) are keyed by `plan_id` only and the plan service does not verify that the referenced plan belongs to the caller’s active workspace before mutating items.

2) **`/invites/accept` handler maps all errors to 400**
- Impact: incorrect HTTP semantics and difficult client UX; unauthorized and forbidden cases are indistinguishable.

3) **`/auth/switch-workspace` handler maps all errors to 400**
- Impact: similar HTTP semantics issue; auth failures and invalid refresh tokens should generally be 401.

### Medium impact / design risks

4) **Admin/owner vs founder semantics are inconsistent across seed + middleware**
- Seed grants `*` permission to `owner` and `admin` (in addition to `founder`).
- `Role` middleware treats `admin` as universal bypass.
- Risk: depending on routing, workspace-admin roles may unintentionally gain platform-wide access.

## Detailed Findings

### 1) Plan item endpoints are not workspace-scoped (cross-tenant modification risk)

**Risk:** A non-trial authenticated user with `plan_item:*` in workspace A can potentially modify items for a plan in workspace B if they can guess/obtain a `plan_id`.

#### Service behavior

The plan service verifies plan ownership/workspace for **plan-level** endpoints (Get/List/Update/Delete/Calculate):
- `internal/service/plan_service.go:233` (`GetPlan`) uses `GetLoadPlan` with `WorkspaceID` for non-trial.
- `internal/service/plan_service.go:371` (`ListPlans`) uses `ListLoadPlans` with `WorkspaceID` for non-trial.
- `internal/service/plan_service.go:412` (`UpdatePlan`) uses `GetLoadPlan(..., WorkspaceID)` for non-trial.
- `internal/service/plan_service.go:684` (`CalculatePlan`) uses `GetLoadPlan(..., WorkspaceID)` for non-trial.

But **plan item** methods only enforce ownership for trial users (guest plans) and otherwise perform mutations purely by `plan_id`:
- `internal/service/plan_service.go:517` (`AddPlanItem`) only checks guest ownership when role is trial; no validation for non-trial.
- `internal/service/plan_service.go:560` (`GetPlanItem`) same.
- `internal/service/plan_service.go:588` (`UpdatePlanItem`) same.
- `internal/service/plan_service.go:658` (`DeletePlanItem`) same.

#### Store/query behavior

Plan item queries do not join against `load_plans` and do not contain workspace predicates:
- `internal/store/plan.sql.go:220` (`GetLoadItem`) is `WHERE plan_id = $1 AND item_id = $2`.
- `internal/store/plan.sql.go:180` (`DeleteLoadItem`) is `WHERE plan_id = $1 AND item_id = $2`.
- `internal/store/plan.sql.go:15` (`AddLoadItem`) inserts directly with `plan_id`.

This is fine only if the service verifies the plan belongs to the caller’s workspace before accessing items.

#### Notes

- Trial (`role=trial`) access is correctly constrained to `created_by_type='guest'` and `created_by_id` via `GetLoadPlanForGuest`.
- Non-trial users rely solely on presence of `workspace_id` in the JWT context, but the item methods do not use it.

### 2) Invite accept handler returns 400 for unauthorized

Service explicitly requires login (MVP):
- `internal/service/invite_service.go:173-178` returns `fmt.Errorf("unauthorized")` when context has no `user_id`.

Handler maps all errors to 400:
- `internal/handler/invite_handler.go:146-149` returns `http.StatusBadRequest` on any error.

### 3) Switch-workspace handler returns 400 for auth-related failures

The handler uses 400 for every error:
- `internal/handler/auth_handler.go:142-145`

Service can fail for reasons that are semantically 401 (missing/invalid access token, invalid refresh token, refresh token ownership mismatch):
- `internal/service/auth_service.go:255-283`

### 4) Admin/owner/founder semantics mismatch (RBAC design risk)

Seed grants global `*` permission to `founder`, `owner`, and `admin`:
- `cmd/db/scripts/default_user_permission.sql:79-101`

Role middleware grants `admin` universal bypass:
- `internal/middleware/role_middleware.go:27-31`

`Permission` middleware supports wildcard `*` and `resource:*` semantics and is likely the primary enforcement mechanism.
If `owner`/`admin` keep global `*`, they can potentially pass permission checks for platform-level resources unless those handlers also enforce tenant scoping via service-layer checks.

## Open Questions / Decisions Needed

1) **Should `POST /api/v1/invites/accept` require a JWT (login) or be public?**
- Current behavior requires login and requires user email to match invite email.

2) **Should global `*` be founder-only?**
- Current seed grants `*` to owner/admin as well. If owner/admin are intended to be workspace-scoped only, they likely should have explicit permissions rather than global `*`.

## Suggested Fix Approach (not implemented)

- In plan item service methods, always verify the plan belongs to the correct tenant before mutating items.
- Normalize handler error->HTTP status mapping for unauthorized vs forbidden vs bad request.
- Align `owner/admin/founder` semantics across seed data, permission middleware, and role middleware.
