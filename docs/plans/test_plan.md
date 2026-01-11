# Comprehensive Test Plan for Load Stuffing Calculator API

**Document Version:** 1.0  
**Date:** January 10, 2026  
**Author:** OpenCode  
**Status:** Draft - Awaiting Review

---

## Table of Contents

1. [Executive Summary](#executive-summary)
2. [Coverage Analysis](#coverage-analysis)
3. [Test Implementation Phases](#test-implementation-phases)
4. [Testing Patterns & Standards](#testing-patterns--standards)
5. [Expected Outcomes](#expected-outcomes)
6. [Open Questions](#open-questions)

---

## Executive Summary

### Current State

After running comprehensive coverage analysis (`go test -coverprofile=coverage.out ./...`), the project shows:

- **Overall Coverage:** ~45% (weighted average)
- **Total Test Lines:** 4,078 lines across service/handler tests
- **Test Framework:** Go testing + testify/assert + httptest + Gin test mode
- **Test Pattern:** Table-driven tests (well-established in service layer)

### Critical Gaps Identified

| Module | Current Coverage | Status | Priority |
|--------|------------------|--------|----------|
| `cmd/api` | 0% | No tests | CRITICAL |
| `internal/api` | 0% | No tests | CRITICAL |
| `internal/auth` | 29.8% | Partial (ValidateToken missing) | HIGH |
| `internal/middleware` | 59.3% | Partial (JWT middleware missing) | HIGH |
| `internal/cache` | 0% | No tests | MEDIUM |
| `internal/gateway` | 0% | No tests | MEDIUM |
| `internal/config` | 0% | No tests | MEDIUM |
| `internal/handler` | 44.7% | Many handlers missing | HIGH |
| `internal/service` | 58.1% | Good, but gaps remain | MEDIUM |
| `internal/response` | 100% | ✓ Complete | - |
| `internal/packer` | 72.9% | Good coverage | - |

### Key Findings

**Functions with 0% Coverage (High Priority):**

- **API Core:** `NewApp`, `setupRouter`, `Run`, `setupRoutes`, `HealthCheck`
- **Auth Context:** All 8 context helper functions
- **Auth JWT:** `ValidateToken` (critical for security)
- **Middleware:** `JWT` (entire auth middleware - 0%)
- **Cache:** All 4 permission cache methods
- **Gateway:** `NewHTTPPackingGateway`, `Pack` (external service integration)
- **Handlers (0%):** Register, GuestToken, RefreshToken, SwitchWorkspace, Dashboard, Workspace (all 4), Member (all 4), Invite (all 4)

**Error Paths Without Coverage:**
- 15+ error handling branches across auth/middleware/gateway
- All `if err != nil` blocks in untested modules

### Proposed Goal

**Target Overall Coverage:** 72-75% (+27-30% improvement)

Focus on:
1. Critical security components (auth, middleware)
2. API infrastructure (routing, initialization)
3. Core business logic gaps (plan calculation, gateway)
4. Handler completeness (workspace, member, invite handlers)

---

## Coverage Analysis

### Detailed Module Breakdown

#### 1. Main API (`cmd/api/main.go`) - 0%

**Lines:** 62 total  
**Untested Functions:**
- `main()` - Bootstraps DB, seeds founder, starts server

**Analysis:**  
- Hard to test due to `os.Exit`, signal handling, goroutines
- Focus should be on extracting testable logic to `internal/api`

---

#### 2. API Core (`internal/api/`) - 0%

**Files:**
- `api.go` (138 lines) - `NewApp`, `setupRouter`, `Run`
- `routes.go` (154 lines) - `setupRoutes`, `HealthCheck`

**Untested Functions:**
- `NewApp(cfg, db)` - Initializes all handlers, services, middleware
- `setupRouter()` - Configures Gin, CORS, routes
- `Run()` - HTTP server lifecycle with graceful shutdown
- `setupRoutes(r)` - Registers all REST endpoints
- `HealthCheck(c)` - Simple health endpoint

**Error Paths:** Minimal (panics are acceptable for initialization)

**Complexity:** Medium - dependency injection, many moving parts

---

#### 3. Auth Infrastructure (`internal/auth/`) - 29.8%

**Files:**
- `context.go` (65 lines) - **0% coverage**
- `jwt.go` (68 lines) - **71% coverage** (ValidateToken missing)
- `password.go` - **100% coverage** ✓

**Untested Functions (context.go):**
- `WithUserID`, `WithRole`, `WithWorkspaceID`, `WithWorkspaceOverrideID`
- `UserIDFromContext`, `RoleFromContext`, `WorkspaceIDFromContext`, `WorkspaceOverrideIDFromContext`

**Untested Functions (jwt.go):**
- `ValidateToken(tokenString, secret)` - **CRITICAL SECURITY FUNCTION**

**Error Paths:**
- ValidateToken: Parse error, invalid signature, expired token, invalid claims
- Context helpers: Type assertion failures

**Complexity:** Low-Medium (jwt.go has critical security logic)

---

#### 4. Middleware (`internal/middleware/`) - 59.3%

**Files:**
- `auth_middleware.go` (52 lines) - **0% coverage**
- `permission_middleware.go` - Has tests (59.3% overall)

**Untested Functions:**
- `JWT(secret)` - Entire JWT authentication middleware

**Error Paths (15 identified):**
1. Missing Authorization header
2. Empty Authorization header
3. Invalid format (not "Bearer <token>")
4. Wrong auth scheme (Basic, etc.)
5. Token parse error
6. Invalid token signature
7. Expired token
8. Missing user_id claim
9. Missing role claim
10. Context propagation failures

**Complexity:** Medium - Critical security component, multiple error branches

---

#### 5. Permission Cache (`internal/cache/`) - 0%

**Files:**
- `permission_cache.go` (37 lines) - **0% coverage**

**Untested Functions:**
- `NewPermissionCache()` - Constructor
- `Get(role)` - Retrieve cached permissions
- `Set(role, permissions)` - Store permissions
- `Invalidate()` - Clear cache

**Error Paths:** None (simple map operations)

**Complexity:** Low - Thread-safety critical (sync.RWMutex)

**Testing Focus:**
- Race condition detection (concurrent reads/writes)
- Cache hit/miss scenarios
- Invalidation correctness

---

#### 6. Packing Gateway (`internal/gateway/`) - 0%

**Files:**
- `packing_gateway.go` (174 lines) - **0% coverage**

**Untested Functions:**
- `NewHTTPPackingGateway(baseURL, timeout)` - Constructor
- `Pack(ctx, req)` - HTTP POST to external packing service

**Error Paths (9+ identified):**
1. Marshal request error
2. HTTP request creation error
3. Network error (connection refused)
4. Context timeout/cancellation
5. HTTP 4xx errors
6. HTTP 5xx errors
7. Response body read error
8. JSON unmarshal error
9. Success=false in response
10. Missing data field
11. Non-200 status with error message

**Complexity:** High - External service integration, many failure modes

**Testing Strategy:** Use `httptest.NewServer` to mock packing service

---

#### 7. Configuration (`internal/config/`) - 0%

**Files:**
- `config.go` (42 lines) - **0% coverage**

**Untested Functions:**
- `Load()` - Reads env vars, loads .env in dev mode

**Error Paths:** Minimal (godotenv.Load error is logged and ignored)

**Complexity:** Low

**Testing Focus:**
- Default values
- Env var overrides
- Dev vs prod mode
- Founder vs Admin backwards compatibility

---

#### 8. Handlers (`internal/handler/`) - 44.7%

**Existing Tests:** auth (partial), user (partial), role (partial), permission (partial), container (partial), product (partial), plan (partial)

**Missing Handler Tests (0% coverage):**

| Handler | Missing Functions | Lines |
|---------|-------------------|-------|
| `auth_handler.go` | Register, GuestToken, RefreshToken, SwitchWorkspace | ~80 |
| `dashboard_handler.go` | NewDashboardHandler, GetStats | ~25 |
| `workspace_handler.go` | All 4 handlers (List, Create, Update, Delete) | ~130 |
| `member_handler.go` | All 4 handlers (List, Add, Update, Delete) | ~180 |
| `invite_handler.go` | All 4 handlers (List, Create, Revoke, Accept) | ~170 |

**Total Untested Handler Lines:** ~585 lines

**Error Paths:** Each handler has 3-5 error paths:
- Invalid JSON binding
- Auth context missing (user_id)
- Service errors
- Validation errors

---

#### 9. Services (`internal/service/`) - 58.1%

**Existing Tests:** Good table-driven coverage for auth (partial), user, role, permission, container, product, plan (partial)

**Missing Service Tests (0% coverage):**

| Service | Missing Functions | Complexity |
|---------|-------------------|------------|
| `auth_service.go` | RefreshToken, SwitchWorkspace, Me, GuestToken | High |
| `plan_service.go` | GetPlanItem, **CalculatePlan** | **CRITICAL** |
| `role_service.go` | GetRolePermissions, UpdateRolePermissions | Medium |
| `user_service.go` | UpdateUser, DeleteUser, ChangePassword | Medium |
| `dashboard_service.go` | GetStats | Low |
| `workspace_service.go` | Some edge cases | Medium |
| `member_service.go` | AddMember (17.5%), UpdateMemberRole (0%) | High |
| `invite_service.go` | CreateInvite (21.2%), RevokeInvite (0%) | Medium |

**Most Critical:** `CalculatePlan` - Core business logic (3D bin packing calculation + DB persistence)

---

## Test Implementation Phases

### Phase 1: Core API Infrastructure (CRITICAL)

**Goal:** Test API bootstrap, routing, and HTTP layer  
**Target Coverage:** 0% → 80-85%  
**Effort:** 2-3 hours  
**Priority:** ⚠️ CRITICAL

#### Deliverables

##### 1.1 `internal/api/api_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestNewApp:
  ✓ happy_path_with_valid_config_and_pool
  ✓ verify_all_handlers_initialized
  ✓ verify_jwt_secret_set
  ✓ edge_empty_config_uses_defaults
  
TestSetupRouter:
  ✓ happy_path_creates_router_with_cors
  ✓ verify_gin_mode_set_correctly
  ✓ verify_cors_headers_configured
  ✓ verify_middleware_chain_order
  
TestHealthCheck:
  ✓ happy_path_returns_200_with_json
  ✓ verify_response_structure (status, time, version)
  ✓ verify_time_format_rfc3339
  ✓ edge_concurrent_requests (race detector)
```

**Testing Approach:**
- Use `httptest.NewRecorder` for HTTP testing
- Mock `pgxpool.Pool` with testcontainers or mock interface
- Use `gin.SetMode(gin.TestMode)`

**Error Paths:** Minimal (initialization panics acceptable)

---

##### 1.2 `internal/api/routes_test.go` (NEW)

```go
// Test Cases

TestSetupRoutes_AllEndpointsRegistered:
  ✓ verify_auth_routes: /login, /register, /refresh, /guest, /switch-workspace, /me
  ✓ verify_protected_routes: /users, /roles, /permissions, /containers, /products, /plans
  ✓ verify_workspace_routes: /workspaces, /members, /invites
  ✓ verify_dashboard_route
  ✓ verify_docs_routes: /docs, /docs/*
  ✓ verify_health_route
  
TestRouteMiddlewareChaining:
  ✓ verify_jwt_middleware_applied_to_protected_routes
  ✓ verify_permission_middleware_correctly_chained
  ✓ verify_auth_routes_not_protected_before_jwt
  ✓ verify_cors_middleware_global
```

**Testing Approach:**
- Inspect `gin.Engine.Routes()` to verify registration
- Test actual requests to ensure middleware order
- Use mock handlers to verify middleware execution

---

### Phase 2: Authentication Infrastructure (HIGH)

**Goal:** Achieve comprehensive auth coverage for security  
**Target Coverage:** 29.8% → 90%+  
**Effort:** 3-4 hours  
**Priority:** 🔴 HIGH

#### Deliverables

##### 2.1 `internal/auth/context_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestContextHelpers:
  ✓ happy_set_and_retrieve_user_id
  ✓ happy_set_and_retrieve_role
  ✓ happy_set_and_retrieve_workspace_id
  ✓ happy_set_and_retrieve_workspace_override_id
  ✓ edge_empty_context_returns_false
  ✓ edge_wrong_type_in_context_returns_false
  ✓ edge_empty_string_values
  ✓ edge_nil_values_for_optional_fields
```

**Testing Approach:**
- Use `context.Background()` as base
- Test each With* and *FromContext pair
- Verify type safety and nil handling

**Error Paths:** Type assertion failures (all paths covered)

---

##### 2.2 `internal/auth/jwt_test.go` (ENHANCE)

```go
// NEW Test Cases (Table-Driven)

TestValidateToken: // ← CRITICAL (0% → 100%)
  ✓ happy_valid_token_with_all_claims
  ✓ happy_valid_token_without_workspace_id
  ✓ error_expired_token
  ✓ error_invalid_signature (wrong secret)
  ✓ error_malformed_token_string
  ✓ error_empty_token_string
  ✓ error_missing_required_claims
  ✓ error_wrong_signing_algorithm
  ✓ edge_token_at_exact_expiration_boundary
  
TestGenerateAccessTokenWithTTL: // ← ENHANCE
  ✓ edge_ttl_zero (immediate expiration)
  ✓ edge_ttl_negative (should fail)
  ✓ edge_very_large_ttl
  ✓ edge_empty_secret
  ✓ edge_empty_user_id
  ✓ edge_nil_workspace_id_pointer
  
TestGenerateRefreshToken: // ← ENHANCE
  ✓ edge_concurrent_generation (uniqueness check)
  ✓ verify_format_rf_YYYYMMDD_HHMMSS.microseconds_random
```

**Testing Approach:**
- Generate tokens with `GenerateAccessToken`
- Validate with `ValidateToken` using correct/incorrect secrets
- Test expiration with custom TTLs
- Mock time if needed (via test helper)

**Error Paths:** All JWT validation errors (9+ paths)

---

##### 2.3 `internal/middleware/auth_middleware_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestJWT: // ← ENTIRE MIDDLEWARE (0% → 100%)
  ✓ happy_valid_bearer_token_populates_context
  ✓ happy_token_with_workspace_id
  ✓ happy_token_without_workspace_id
  ✓ error_missing_authorization_header → 401
  ✓ error_empty_authorization_header → 401
  ✓ error_invalid_format_no_bearer → 401
  ✓ error_invalid_format_wrong_scheme → 401
  ✓ error_too_many_parts → 401
  ✓ error_expired_token → 401
  ✓ error_invalid_signature → 401
  ✓ error_malformed_token → 401
  ✓ verify_gin_context_user_id_set
  ✓ verify_gin_context_role_set
  ✓ verify_gin_context_workspace_id_set
  ✓ verify_request_context_updated
  ✓ verify_next_called_on_success
  ✓ verify_abort_called_on_failure
```

**Testing Approach:**
- Use `httptest.NewRecorder` + `gin.CreateTestContext`
- Generate valid tokens with `auth.GenerateAccessToken`
- Test all error branches (15+ paths)
- Verify `c.Next()` vs `c.Abort()` behavior

**Error Paths:** ALL 15 identified error paths must be covered

---

### Phase 3: Supporting Infrastructure (MEDIUM)

**Goal:** Complete coverage for cache, gateway, config  
**Target Coverage:** 0% → 80-90%  
**Effort:** 3-4 hours  
**Priority:** 🟡 MEDIUM

#### Deliverables

##### 3.1 `internal/cache/permission_cache_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestNewPermissionCache:
  ✓ happy_creates_non_nil_cache
  
TestPermissionCache_GetSet:
  ✓ happy_set_then_get_returns_same_permissions
  ✓ happy_get_nonexistent_returns_nil_false
  ✓ edge_set_empty_permissions_slice
  ✓ edge_set_nil_permissions
  ✓ edge_overwrite_existing_role
  ✓ concurrent_multiple_readers (race detector)
  ✓ concurrent_multiple_writers (race detector)
  ✓ concurrent_mixed_read_write (race detector)
  
TestPermissionCache_Invalidate:
  ✓ happy_invalidate_clears_all
  ✓ edge_invalidate_empty_cache
  ✓ verify_get_after_invalidate_returns_false
  ✓ concurrent_invalidate_during_operations
```

**Testing Approach:**
- Test basic CRUD operations
- Use `go test -race` for concurrency testing
- Launch goroutines for concurrent access
- Use `sync.WaitGroup` to coordinate

**Error Paths:** None (simple map operations)

**Focus:** Thread-safety verification

---

##### 3.2 `internal/gateway/packing_gateway_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestNewHTTPPackingGateway:
  ✓ happy_creates_gateway_with_url_and_timeout
  ✓ edge_empty_baseurl_defaults_to_localhost
  ✓ edge_trailing_slash_trimmed
  ✓ edge_zero_timeout
  
TestHTTPPackingGateway_Pack: // ← CRITICAL BUSINESS LOGIC
  ✓ happy_successful_pack_with_placements
  ✓ happy_successful_pack_with_unfitted_items
  ✓ happy_zero_placements_zero_unfitted
  ✓ error_network_error_server_unreachable
  ✓ error_context_cancelled_timeout
  ✓ error_http_400_with_error_response
  ✓ error_http_500_with_error_message
  ✓ error_http_200_but_success_false
  ✓ error_http_200_but_missing_data_field
  ✓ error_invalid_json_response_body
  ✓ error_empty_response_body
  ✓ edge_very_large_response_many_placements
```

**Testing Approach:**
- Use `httptest.NewServer` to mock packing service
- Return controlled JSON responses
- Simulate network errors by closing server early
- Use `context.WithTimeout` for timeout tests

**Error Paths:** ALL 11+ error handling branches must be covered

---

##### 3.3 `internal/config/config_test.go` (NEW)

```go
// Test Cases (Table-Driven with env mocking)

TestLoad:
  ✓ happy_load_with_all_env_vars_set
  ✓ happy_load_with_defaults_no_env_vars
  ✓ edge_srv_env_dev_loads_dotenv_file
  ✓ edge_srv_env_prod_skips_dotenv_file
  ✓ edge_missing_dotenv_file_in_dev (should continue)
  ✓ edge_founder_vs_admin_env_var_backwards_compatibility
  ✓ verify_all_fields_populated_correctly
  ✓ verify_default_values
```

**Testing Approach:**
- Use `t.Setenv()` (Go 1.17+) or `os.Setenv` + cleanup
- Test with/without .env file
- Verify fallback logic for FOUNDER_* vs ADMIN_*

**Error Paths:** Minimal (godotenv error ignored)

---

### Phase 4: Missing Handler Coverage (HIGH)

**Goal:** Complete handler layer coverage  
**Target Coverage:** 44.7% → 75%+  
**Effort:** 4-5 hours  
**Priority:** 🟠 HIGH

#### Deliverables

##### 4.1 `internal/handler/auth_handler_test.go` (ENHANCE)

```go
// NEW Test Functions (Table-Driven)

TestAuthHandler_Register:
  ✓ happy_successful_registration
  ✓ happy_registration_with_guest_token_claims_plans
  ✓ error_invalid_request_body
  ✓ error_service_error_user_exists
  ✓ error_database_error
  ✓ verify_returns_201_on_success
  
TestAuthHandler_GuestToken:
  ✓ happy_generate_guest_token_successfully
  ✓ error_service_error
  ✓ verify_returns_200_with_token
  
TestAuthHandler_RefreshToken:
  ✓ happy_valid_refresh_token
  ✓ error_invalid_request_body
  ✓ error_invalid_refresh_token
  ✓ error_expired_refresh_token
  ✓ error_service_error
  
TestAuthHandler_SwitchWorkspace:
  ✓ happy_switch_to_valid_workspace
  ✓ error_invalid_request_body
  ✓ error_workspace_not_found
  ✓ error_no_permission_to_workspace
  ✓ error_service_error
```

**Pattern:** Follow existing `auth_handler_test.go` pattern with mocked `MockAuthService`

---

##### 4.2 `internal/handler/dashboard_handler_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestNewDashboardHandler:
  ✓ happy_constructor_returns_non_nil_handler
  
TestDashboardHandler_GetStats:
  ✓ happy_returns_dashboard_stats
  ✓ error_service_error
  ✓ error_unauthorized_no_user_id_in_context
  ✓ verify_correct_status_codes
```

**Mock Required:** `MockDashboardService`

---

##### 4.3 `internal/handler/workspace_handler_test.go` (NEW)

```go
// Test Cases (All 4 handlers, Table-Driven)

TestWorkspaceHandler_ListWorkspaces:
  ✓ happy_returns_workspace_list
  ✓ error_service_error
  ✓ error_unauthorized
  
TestWorkspaceHandler_CreateWorkspace:
  ✓ happy_creates_workspace
  ✓ error_invalid_json
  ✓ error_service_error
  ✓ verify_returns_201
  
TestWorkspaceHandler_UpdateWorkspace:
  ✓ happy_updates_workspace
  ✓ error_invalid_json
  ✓ error_workspace_not_found
  ✓ error_service_error
  
TestWorkspaceHandler_DeleteWorkspace:
  ✓ happy_deletes_workspace
  ✓ error_workspace_not_found
  ✓ error_service_error
  ✓ verify_returns_204_or_200
```

**Mock Required:** `MockWorkspaceService`

---

##### 4.4 `internal/handler/member_handler_test.go` (NEW)

```go
// Test Cases (All 4 handlers, Table-Driven)

TestMemberHandler_ListMembers:
  ✓ happy_returns_member_list
  ✓ error_service_error
  ✓ error_unauthorized
  
TestMemberHandler_AddMember:
  ✓ happy_adds_member
  ✓ error_invalid_json
  ✓ error_user_not_found
  ✓ error_already_member
  ✓ error_service_error
  
TestMemberHandler_UpdateMemberRole:
  ✓ happy_updates_member_role
  ✓ error_invalid_json
  ✓ error_member_not_found
  ✓ error_invalid_role
  ✓ error_service_error
  
TestMemberHandler_DeleteMember:
  ✓ happy_removes_member
  ✓ error_member_not_found
  ✓ error_cannot_remove_owner
  ✓ error_service_error
```

**Mock Required:** `MockMemberService`

---

##### 4.5 `internal/handler/invite_handler_test.go` (NEW)

```go
// Test Cases (All 4 handlers, Table-Driven)

TestInviteHandler_ListInvites:
  ✓ happy_returns_invite_list
  ✓ error_service_error
  ✓ error_unauthorized
  
TestInviteHandler_CreateInvite:
  ✓ happy_creates_invite
  ✓ error_invalid_json
  ✓ error_user_already_member
  ✓ error_invalid_email
  ✓ error_service_error
  
TestInviteHandler_RevokeInvite:
  ✓ happy_revokes_invite
  ✓ error_invite_not_found
  ✓ error_unauthorized
  ✓ error_service_error
  
TestInviteHandler_AcceptInvite:
  ✓ happy_accepts_invite
  ✓ error_invalid_json
  ✓ error_invalid_token
  ✓ error_expired_invite
  ✓ error_already_accepted
  ✓ error_service_error
```

**Mock Required:** `MockInviteService`

---

### Phase 5: Service Layer Gap Filling (MEDIUM)

**Goal:** Complete service coverage for critical business logic  
**Target Coverage:** 58.1% → 75%+  
**Effort:** 4-6 hours  
**Priority:** 🟡 MEDIUM-HIGH

#### Deliverables

##### 5.1 `internal/service/auth_service_test.go` (ENHANCE)

```go
// NEW Test Functions (Table-Driven)

TestAuthService_RefreshToken:
  ✓ happy_valid_refresh_token_returns_new_tokens
  ✓ error_refresh_token_not_found
  ✓ error_refresh_token_expired
  ✓ error_user_not_found
  ✓ error_database_error
  
TestAuthService_SwitchWorkspace:
  ✓ happy_switch_to_valid_workspace
  ✓ error_workspace_not_found
  ✓ error_user_not_member_of_workspace
  ✓ error_invalid_workspace_id
  ✓ error_database_error
  
TestAuthService_Me:
  ✓ happy_returns_user_info_with_permissions
  ✓ happy_platform_member_with_platform_role
  ✓ happy_workspace_member_with_workspace_role
  ✓ error_no_user_id_in_context
  ✓ error_user_not_found
  ✓ error_database_error
  
TestAuthService_GuestToken:
  ✓ happy_generates_guest_token_successfully
  ✓ verify_token_has_trial_role
  ✓ verify_token_has_no_workspace_id
  ✓ error_token_generation_fails
```

**Pattern:** Follow existing `TestAuthService_Login` pattern

---

##### 5.2 `internal/service/plan_service_test.go` (ENHANCE)

```go
// NEW Test Functions (Table-Driven)

TestPlanService_GetPlanItem:
  ✓ happy_returns_plan_item_by_id
  ✓ error_plan_not_found
  ✓ error_item_not_found
  ✓ error_no_permission
  ✓ error_database_error
  
TestPlanService_CalculatePlan: // ← CRITICAL BUSINESS LOGIC
  ✓ happy_successful_calculation_all_items_fit
  ✓ happy_partial_fit_some_items_unfitted
  ✓ happy_no_items_fit
  ✓ error_plan_not_found
  ✓ error_no_items_in_plan
  ✓ error_packing_service_unreachable
  ✓ error_packing_service_returns_error
  ✓ error_database_error_saving_results
  ✓ error_database_error_saving_placements
  ✓ verify_plan_status_updated_correctly
  ✓ verify_placement_coordinates_saved
  ✓ verify_rotation_codes_saved
  ✓ verify_step_numbers_saved
  ✓ edge_zero_items
  ✓ edge_container_too_small_for_all_items
```

**Mock Required:** `MockPackingService` (gateway)

**Complexity:** HIGH - Core 3D packing logic + DB persistence

---

##### 5.3 `internal/service/role_service_test.go` (ENHANCE)

```go
// NEW Test Functions (Table-Driven)

TestRoleService_GetRolePermissions:
  ✓ happy_returns_role_permissions
  ✓ error_role_not_found
  ✓ error_database_error
  
TestRoleService_UpdateRolePermissions:
  ✓ happy_updates_permissions
  ✓ error_role_not_found
  ✓ error_invalid_permission_ids
  ✓ error_database_transaction_error
  ✓ verify_cache_invalidated
```

---

##### 5.4 `internal/service/user_service_test.go` (ENHANCE)

```go
// NEW Test Functions (Table-Driven)

TestUserService_UpdateUser:
  ✓ happy_updates_user_details
  ✓ error_user_not_found
  ✓ error_username_already_exists
  ✓ error_email_already_exists
  ✓ error_database_error
  
TestUserService_DeleteUser:
  ✓ happy_deletes_user
  ✓ error_user_not_found
  ✓ error_cannot_delete_self
  ✓ error_database_error
  
TestUserService_ChangePassword:
  ✓ happy_changes_password
  ✓ error_user_not_found
  ✓ error_current_password_incorrect
  ✓ error_new_password_too_weak
  ✓ error_database_error
```

---

##### 5.5 `internal/service/dashboard_service_test.go` (NEW)

```go
// Test Cases (Table-Driven)

TestNewDashboardService:
  ✓ happy_constructor_returns_non_nil_service
  
TestDashboardService_GetStats:
  ✓ happy_returns_stats_for_user
  ✓ happy_returns_stats_for_workspace
  ✓ error_user_not_found
  ✓ error_workspace_not_found
  ✓ error_database_error
  ✓ verify_stat_counts_correct
```

---

##### 5.6 Additional Service Enhancements

```go
// internal/service/member_service_test.go (ENHANCE)
TestMemberService_AddMember:
  ✓ Add error path coverage (currently 17.5%)
  
TestMemberService_UpdateMemberRole:
  ✓ NEW: Full coverage (currently 0%)

// internal/service/invite_service_test.go (ENHANCE)
TestInviteService_CreateInvite:
  ✓ Add error path coverage (currently 21.2%)
  
TestInviteService_RevokeInvite:
  ✓ NEW: Full coverage (currently 0%)
```

---

### Phase 6: Integration Tests (OPTIONAL)

**Goal:** End-to-end confidence with real database  
**Target Coverage:** Supplement existing coverage  
**Effort:** 6-8 hours  
**Priority:** 🔵 LOW-MEDIUM (Nice-to-have)

#### Deliverables

##### 6.1 `cmd/api/main_test.go` (NEW - OPTIONAL)

```go
// Integration Test Cases

TestMain_Bootstrap: // ← Requires testcontainers or docker-compose
  ✓ happy_db_connection_successful
  ✓ happy_founder_seeded_successfully
  ✓ happy_app_initialization_successful
  ✓ error_db_connection_fails
  ✓ error_founder_seed_fails (should log warning, not exit)
```

**Approach:**
- Use `testcontainers-go` for PostgreSQL
- Test actual DB connection + migration
- Test founder seeding idempotency

**Complexity:** HIGH - Requires Docker, slow tests

**Recommendation:** Lower priority unless integration tests are critical

---

## Testing Patterns & Standards

### 1. Table-Driven Tests (MANDATORY)

All tests MUST follow the table-driven pattern:

```go
func TestFunctionName(t *testing.T) {
    tests := []struct {
        name    string        // Test case name (snake_case)
        input   InputType     // Function input
        mock    MockBehavior  // Mock setup function
        want    OutputType    // Expected output
        wantErr bool          // Expect error?
        errMsg  string        // Expected error message (optional)
    }{
        {
            name:    "happy_path_description",
            input:   InputType{...},
            mock:    func(m *MockService) { ... },
            want:    ExpectedValue,
            wantErr: false,
        },
        {
            name:    "error_case_description",
            input:   InputType{...},
            mock:    func(m *MockService) { ... },
            wantErr: true,
            errMsg:  "expected error message",
        },
        {
            name:    "edge_case_description",
            input:   InputType{...},
            want:    EdgeCaseValue,
            wantErr: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Arrange
            mockSvc := new(MockService)
            if tt.mock != nil {
                tt.mock(mockSvc)
            }
            
            // Act
            got, err := FunctionUnderTest(tt.input)
            
            // Assert
            if tt.wantErr {
                assert.Error(t, err)
                if tt.errMsg != "" {
                    assert.Contains(t, err.Error(), tt.errMsg)
                }
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
            
            // Verify mocks
            mockSvc.AssertExpectations(t)
        })
    }
}
```

**Test Case Naming Convention:**
- `happy_*` - Successful execution paths
- `error_*` - Error handling paths
- `edge_*` - Edge cases (nil, empty, boundary values)
- Use snake_case for readability

---

### 2. Error Path Coverage (MANDATORY)

**Rule:** Every `if err != nil` block in source code MUST have a corresponding test case.

**Example:**

```go
// Source code
func Process(data string) error {
    if data == "" {
        return errors.New("data is required")  // ← Test case needed
    }
    
    result, err := ExternalCall(data)
    if err != nil {
        return fmt.Errorf("external call: %w", err)  // ← Test case needed
    }
    
    if err := SaveResult(result); err != nil {
        return fmt.Errorf("save: %w", err)  // ← Test case needed
    }
    
    return nil
}

// Test cases MUST include:
// ✓ error_empty_data
// ✓ error_external_call_fails
// ✓ error_save_fails
```

**Verification:**
- Use coverage report to find uncovered error branches
- Prioritize error paths in security-critical code (auth, validation)

---

### 3. Branch Coverage (MANDATORY)

**Rule:** All `if/else`, `switch/case` branches must be executed by tests.

**Example:**

```go
// Source code
func GetStatus(code int) string {
    switch code {
    case 200:
        return "OK"      // ← Test case needed
    case 404:
        return "Not Found"  // ← Test case needed
    case 500:
        return "Error"   // ← Test case needed
    default:
        return "Unknown" // ← Test case needed
    }
}

// Test cases MUST include:
// ✓ case_200_returns_ok
// ✓ case_404_returns_not_found
// ✓ case_500_returns_error
// ✓ case_unknown_code_returns_unknown
```

---

### 4. Edge Case Testing (MANDATORY)

**Required edge cases for all tests:**

| Data Type | Edge Cases to Test |
|-----------|-------------------|
| **String** | Empty `""`, Whitespace `"   "`, Very long (>1000 chars), Special chars, Unicode |
| **Integer** | Zero `0`, Negative `-1`, Max `math.MaxInt64`, Min `math.MinInt64` |
| **Pointer** | Nil `nil`, Valid pointer |
| **Slice** | Nil `nil`, Empty `[]`, Single element `[x]`, Many elements |
| **Map** | Nil `nil`, Empty `map[]`, Single key, Many keys |
| **Time** | Zero `time.Time{}`, Past, Future, Boundary (midnight, year boundary) |

**Example:**

```go
tests := []struct {
    name  string
    input string
    want  string
}{
    {name: "happy_normal_input", input: "test", want: "processed"},
    {name: "edge_empty_string", input: "", want: ""},
    {name: "edge_whitespace_only", input: "   ", want: ""},
    {name: "edge_very_long_string", input: strings.Repeat("a", 10000), want: "..."},
    {name: "edge_unicode", input: "日本語", want: "processed"},
}
```

---

### 5. Mock Usage Patterns

#### Handler Tests: Mock Services

```go
// Example: Testing handler with mocked service
func TestUserHandler_CreateUser(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    tests := []struct {
        name       string
        payload    dto.CreateUserRequest
        mockSetup  func(*mocks.MockUserService)
        wantStatus int
    }{
        {
            name:    "happy_create_user",
            payload: dto.CreateUserRequest{Username: "test", ...},
            mockSetup: func(m *mocks.MockUserService) {
                m.On("CreateUser", mock.Anything, mock.Anything).
                    Return(&dto.UserResponse{ID: "123"}, nil)
            },
            wantStatus: http.StatusCreated,
        },
        {
            name:    "error_service_fails",
            payload: dto.CreateUserRequest{Username: "test", ...},
            mockSetup: func(m *mocks.MockUserService) {
                m.On("CreateUser", mock.Anything, mock.Anything).
                    Return(nil, errors.New("db error"))
            },
            wantStatus: http.StatusInternalServerError,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockSvc := new(mocks.MockUserService)
            tt.mockSetup(mockSvc)
            
            h := handler.NewUserHandler(mockSvc)
            
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            
            jsonBytes, _ := json.Marshal(tt.payload)
            c.Request = httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(jsonBytes))
            c.Request.Header.Set("Content-Type", "application/json")
            
            h.CreateUser(c)
            
            assert.Equal(t, tt.wantStatus, w.Code)
            mockSvc.AssertExpectations(t)
        })
    }
}
```

#### Service Tests: Mock Querier

```go
// Example: Testing service with mocked database
func TestUserService_CreateUser(t *testing.T) {
    tests := []struct {
        name      string
        req       dto.CreateUserRequest
        mockSetup func(*mocks.MockQuerier)
        wantErr   bool
    }{
        {
            name: "happy_create",
            req:  dto.CreateUserRequest{Username: "test", ...},
            mockSetup: func(m *mocks.MockQuerier) {
                m.On("CreateUser", mock.Anything, mock.Anything).
                    Return(store.User{UserID: uuid.New()}, nil)
            },
            wantErr: false,
        },
        {
            name: "error_duplicate_username",
            req:  dto.CreateUserRequest{Username: "existing", ...},
            mockSetup: func(m *mocks.MockQuerier) {
                m.On("CreateUser", mock.Anything, mock.Anything).
                    Return(store.User{}, errors.New("duplicate key"))
            },
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockQ := new(mocks.MockQuerier)
            tt.mockSetup(mockQ)
            
            svc := service.NewUserService(mockQ)
            
            _, err := svc.CreateUser(context.Background(), tt.req)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
            
            mockQ.AssertExpectations(t)
        })
    }
}
```

#### Middleware Tests: Mock Gin Context

```go
// Example: Testing middleware
func TestJWTMiddleware(t *testing.T) {
    gin.SetMode(gin.TestMode)
    
    tests := []struct {
        name       string
        authHeader string
        wantStatus int
        wantAbort  bool
    }{
        {
            name:       "happy_valid_token",
            authHeader: "Bearer " + validToken,
            wantStatus: http.StatusOK,
            wantAbort:  false,
        },
        {
            name:       "error_missing_header",
            authHeader: "",
            wantStatus: http.StatusUnauthorized,
            wantAbort:  true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            w := httptest.NewRecorder()
            c, _ := gin.CreateTestContext(w)
            c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
            if tt.authHeader != "" {
                c.Request.Header.Set("Authorization", tt.authHeader)
            }
            
            middleware.JWT("secret")(c)
            
            assert.Equal(t, tt.wantAbort, c.IsAborted())
            if tt.wantAbort {
                assert.Equal(t, tt.wantStatus, w.Code)
            }
        })
    }
}
```

---

### 6. Assertion Patterns

**Use `testify/assert` for all assertions:**

```go
import (
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

// Basic assertions
assert.NoError(t, err)                       // Error should be nil
assert.Error(t, err)                         // Error should NOT be nil
assert.Equal(t, expected, actual)            // Values should be equal
assert.NotEqual(t, unexpected, actual)       // Values should differ
assert.Nil(t, value)                         // Value should be nil
assert.NotNil(t, value)                      // Value should NOT be nil
assert.True(t, condition)                    // Condition should be true
assert.False(t, condition)                   // Condition should be false

// Collection assertions
assert.Len(t, slice, 5)                      // Slice should have length 5
assert.Contains(t, slice, item)              // Slice should contain item
assert.ElementsMatch(t, expected, actual)    // Slices have same elements (any order)
assert.Empty(t, slice)                       // Slice should be empty

// String assertions
assert.Contains(t, str, "substring")         // String should contain substring
assert.NotContains(t, str, "substring")      // String should NOT contain substring

// Use require.* for critical assertions (stops test immediately on failure)
require.NoError(t, err)                      // Stop if error occurs
require.NotNil(t, value)                     // Stop if value is nil
```

**Assertion Best Practices:**
- Use `assert.*` for regular checks (test continues on failure)
- Use `require.*` for setup/critical checks (test stops on failure)
- Provide custom messages for complex assertions: `assert.Equal(t, expected, actual, "should return correct user ID")`

---

### 7. Test File Organization

```go
package handler_test  // Use _test suffix for external testing

import (
    "testing"
    // Standard library
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    
    // External dependencies
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    
    // Internal packages
    "github.com/ekastn/load-stuffing-calculator/internal/dto"
    "github.com/ekastn/load-stuffing-calculator/internal/handler"
    "github.com/ekastn/load-stuffing-calculator/internal/mocks"
)

// Helper functions at top (if needed)
func stringPtr(s string) *string { return &s }

// Test functions below, grouped by function under test
func TestHandlerName_MethodName(t *testing.T) { ... }
```

---

### 8. Race Condition Testing

For concurrent code (e.g., `PermissionCache`), always test with race detector:

```bash
go test -race ./internal/cache
```

**Example concurrent test:**

```go
func TestPermissionCache_ConcurrentAccess(t *testing.T) {
    cache := cache.NewPermissionCache()
    
    var wg sync.WaitGroup
    
    // Launch 100 concurrent writers
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            role := fmt.Sprintf("role_%d", id)
            perms := []string{fmt.Sprintf("perm_%d", id)}
            cache.Set(role, perms)
        }(i)
    }
    
    // Launch 100 concurrent readers
    for i := 0; i < 100; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            role := fmt.Sprintf("role_%d", id)
            _, _ = cache.Get(role)
        }(i)
    }
    
    wg.Wait()
    
    // Verify no panics occurred (race detector will catch issues)
    assert.NotNil(t, cache)
}
```

---

### 9. HTTP Mocking for Gateway Tests

**Use `httptest.NewServer` to mock external HTTP services:**

```go
func TestPackingGateway_Pack(t *testing.T) {
    tests := []struct {
        name         string
        mockResponse string
        mockStatus   int
        wantErr      bool
    }{
        {
            name: "happy_successful_pack",
            mockResponse: `{
                "success": true,
                "data": {
                    "units": "mm",
                    "placements": [
                        {"item_id": "1", "pos_x": 0, "pos_y": 0, "pos_z": 0}
                    ],
                    "unfitted": [],
                    "stats": {"fitted_count": 1}
                }
            }`,
            mockStatus: http.StatusOK,
            wantErr:    false,
        },
        {
            name:         "error_http_500",
            mockResponse: `{"success": false, "error": {"message": "server error"}}`,
            mockStatus:   http.StatusInternalServerError,
            wantErr:      true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock HTTP server
            server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                w.WriteHeader(tt.mockStatus)
                w.Write([]byte(tt.mockResponse))
            }))
            defer server.Close()
            
            // Create gateway pointing to mock server
            gw := gateway.NewHTTPPackingGateway(server.URL, 5*time.Second)
            
            req := gateway.PackRequest{
                Units:     "mm",
                Container: gateway.PackContainerIn{Length: 1000, Width: 500, Height: 500},
                Items:     []gateway.PackItemIn{{ItemID: "1", Length: 100, Width: 100, Height: 100}},
            }
            
            resp, err := gw.Pack(context.Background(), req)
            
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
            }
        })
    }
}
```

---

## Expected Outcomes

### Coverage Improvements (Estimated)

| Module | Current | Target | Delta | Priority |
|--------|---------|--------|-------|----------|
| `cmd/api` | 0% | 40% | +40% | CRITICAL |
| `internal/api` | 0% | 85% | +85% | CRITICAL |
| `internal/auth` | 29.8% | 90% | +60.2% | HIGH |
| `internal/middleware` | 59.3% | 95% | +35.7% | HIGH |
| `internal/cache` | 0% | 90% | +90% | MEDIUM |
| `internal/gateway` | 0% | 80% | +80% | MEDIUM |
| `internal/config` | 0% | 70% | +70% | MEDIUM |
| `internal/handler` | 44.7% | 75% | +30.3% | HIGH |
| `internal/service` | 58.1% | 75% | +16.9% | MEDIUM |

**Overall Project Coverage:**  
45% → **72-75%** (+27-30%)

---

### Test File Count

**New Test Files:** 13  
**Enhanced Test Files:** 5  
**Total Test Files After:** 36+

**New Test Files:**
1. `internal/api/api_test.go`
2. `internal/api/routes_test.go`
3. `internal/auth/context_test.go`
4. `internal/middleware/auth_middleware_test.go`
5. `internal/cache/permission_cache_test.go`
6. `internal/gateway/packing_gateway_test.go`
7. `internal/config/config_test.go`
8. `internal/handler/dashboard_handler_test.go`
9. `internal/handler/workspace_handler_test.go`
10. `internal/handler/member_handler_test.go`
11. `internal/handler/invite_handler_test.go`
12. `internal/service/dashboard_service_test.go`
13. `cmd/api/main_test.go` (optional)

**Enhanced Test Files:**
1. `internal/auth/jwt_test.go`
2. `internal/handler/auth_handler_test.go`
3. `internal/service/auth_service_test.go`
4. `internal/service/plan_service_test.go`
5. `internal/service/role_service_test.go`
6. `internal/service/user_service_test.go`
7. `internal/service/member_service_test.go`
8. `internal/service/invite_service_test.go`

---

### Quality Improvements

**Security:**
- ✓ JWT validation fully tested (auth.ValidateToken)
- ✓ Auth middleware fully tested (all error paths)
- ✓ Token expiration edge cases covered
- ✓ Authorization header format validation tested

**Reliability:**
- ✓ All error paths in gateway tested (network, timeout, malformed responses)
- ✓ Plan calculation edge cases covered (empty, partial fit, no fit)
- ✓ Permission cache thread-safety verified with race detector
- ✓ API initialization and routing verified

**Maintainability:**
- ✓ Consistent table-driven test pattern across all tests
- ✓ Clear test case naming (happy/error/edge)
- ✓ Comprehensive mock usage documentation
- ✓ All tests follow project conventions

---

### Performance Impact

**Test Execution Time (Estimated):**

| Phase | Test Files | Estimated Runtime |
|-------|------------|-------------------|
| Phase 1 (API) | 2 | 5-10 seconds |
| Phase 2 (Auth) | 2 | 10-15 seconds |
| Phase 3 (Infrastructure) | 3 | 10-15 seconds |
| Phase 4 (Handlers) | 5 | 15-20 seconds |
| Phase 5 (Services) | 8 | 20-30 seconds |
| Phase 6 (Integration) | 1 | 30-60 seconds* |

**Total Runtime:** 60-90 seconds (without integration tests)  
**With Integration Tests:** 90-150 seconds*

*Integration tests with testcontainers are significantly slower

**CI/CD Considerations:**
- Use `-race` flag selectively (slower but catches concurrency bugs)
- Run integration tests separately in CI pipeline
- Parallelize test packages with `-p` flag

---

## Open Questions

Before proceeding with implementation, please clarify:

### 1. Phase Priority

**Question:** Which phases should be implemented first?

**Recommended Order:**
1. Phase 2 (Auth) - Security critical
2. Phase 1 (API Core) - Foundation for integration
3. Phase 3 (Infrastructure) - Supporting components
4. Phase 4 (Handlers) - API completeness
5. Phase 5 (Services) - Business logic depth
6. Phase 6 (Integration) - Optional confidence boost

**Your Preference:** _________

---

### 2. Integration vs Unit Testing

**Question:** For `internal/api/api_test.go`, which approach do you prefer?

**Option A: Pure Unit Tests (Recommended)**
- Mock `pgxpool.Pool` interface
- Fast execution (~2-3 seconds)
- No external dependencies
- Focus on initialization logic

**Option B: Integration-Style Tests**
- Use testcontainers for real PostgreSQL
- Slower execution (~30-60 seconds)
- More realistic but requires Docker
- Tests actual DB connection

**Your Preference:** _________

---

### 3. Gateway Testing Strategy

**Question:** How should we test the external packing service gateway?

**Option A: Mock HTTP Server (Recommended)**
- Use `httptest.NewServer`
- Fast, deterministic
- Tests error handling thoroughly
- No external dependencies

**Option B: Skip Testing**
- Not recommended
- Leaves 174 lines untested
- Risky for integration failures

**Option C: Real Service**
- Requires packing service running
- Slow, flaky
- Not suitable for unit tests

**Your Preference:** _________

---

### 4. Coverage Target

**Question:** What is your target overall coverage percentage?

**Current:** 45%  
**Realistic Target:** 72-75% (with plan)  
**Aggressive Target:** 80%+ (requires Phase 6 + additional work)

**Your Target:** _________%

---

### 5. Test Database for Services

**Question:** Do you have a test database setup, or should service tests use mocks exclusively?

**Option A: Mocks Only (Current Pattern)**
- Already established in codebase
- Fast execution
- No DB setup required
- Good for unit testing

**Option B: Real Test DB**
- More realistic
- Requires testcontainers or docker-compose
- Slower but catches DB-specific issues

**Your Preference:** _________

---

### 6. Race Detector in CI

**Question:** Should the CI pipeline always run with `-race` flag?

**Option A: Always Use -race (Recommended for Critical Code)**
- Catches concurrency bugs early
- ~2-10x slower test execution
- Essential for cache, auth, middleware tests

**Option B: Separate CI Job**
- Run normal tests fast
- Run race tests separately
- Best of both worlds

**Your Preference:** _________

---

### 7. Test File Naming

**Question:** Current pattern is `*_test.go` in same package. Continue or change?

**Current Pattern:**
- `internal/handler/auth_handler_test.go` → `package handler_test`
- External testing (can't access private functions)

**Alternative:**
- Same package (`package handler`) for testing private functions

**Your Preference:** _________

---

## Next Steps

1. **Review this plan** and answer open questions above
2. **Approve phases** to implement (or modify priority)
3. **Begin implementation** starting with highest priority phase
4. **Generate coverage reports** after each phase
5. **Iterate** based on coverage gaps found

---

## Appendix: Commands Reference

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...

# View coverage report (text)
go tool cover -func=coverage.out

# View coverage report (HTML)
go tool cover -html=coverage.out -o coverage.html

# Run tests with race detector
go test -race ./...

# Run specific package tests
go test ./internal/api -v

# Run specific test function
go test ./internal/handler -run '^TestAuthHandler_Login$' -v

# Run tests in parallel
go test -p 4 ./...

# Generate coverage for specific packages
go test -coverprofile=coverage.out -covermode=atomic ./internal/api ./internal/auth
```

### Coverage Analysis

```bash
# Find untested functions
go tool cover -func=coverage.out | grep "0.0%"

# Sort by coverage percentage
go tool cover -func=coverage.out | awk '{print $3,$0}' | sort -n

# Package-level coverage summary
go test -cover ./...

# Detailed coverage with branch info
go test -covermode=count -coverprofile=coverage.out ./...
```

### CI/CD Integration

```yaml
# Example GitHub Actions workflow
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      
      - name: Run tests
        run: go test -v -coverprofile=coverage.out ./...
      
      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.out
  
  race:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
      
      - name: Run race detector
        run: go test -race ./...
```

---

**Document Status:** Draft - Awaiting Review  
**Last Updated:** January 10, 2026
