package cache

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPermissionCache(t *testing.T) {
	cache := NewPermissionCache()

	assert.NotNil(t, cache)
	assert.NotNil(t, cache.permissions)
	assert.Equal(t, 0, len(cache.permissions))
}

func TestPermissionCache_SetAndGet(t *testing.T) {
	tests := []struct {
		name                string
		role                string
		permissions         []string
		expectFound         bool
		expectedPermissions []string
	}{
		{
			name:                "set_and_get_admin_permissions",
			role:                "admin",
			permissions:         []string{"*"},
			expectFound:         true,
			expectedPermissions: []string{"*"},
		},
		{
			name:                "set_and_get_multiple_permissions",
			role:                "planner",
			permissions:         []string{"plan:create", "plan:read", "plan:update", "plan:delete"},
			expectFound:         true,
			expectedPermissions: []string{"plan:create", "plan:read", "plan:update", "plan:delete"},
		},
		{
			name:                "set_and_get_empty_permissions",
			role:                "viewer",
			permissions:         []string{},
			expectFound:         true,
			expectedPermissions: []string{},
		},
		{
			name:                "set_and_get_nil_permissions",
			role:                "none",
			permissions:         nil,
			expectFound:         true,
			expectedPermissions: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewPermissionCache()

			// Set permissions
			cache.Set(tt.role, tt.permissions)

			// Get permissions
			perms, found := cache.Get(tt.role)

			assert.Equal(t, tt.expectFound, found)
			assert.Equal(t, tt.expectedPermissions, perms)
		})
	}
}

func TestPermissionCache_GetNonexistent(t *testing.T) {
	tests := []struct {
		name        string
		role        string
		expectFound bool
	}{
		{
			name:        "get_nonexistent_role",
			role:        "nonexistent",
			expectFound: false,
		},
		{
			name:        "get_empty_string_role",
			role:        "",
			expectFound: false,
		},
		{
			name:        "get_from_empty_cache",
			role:        "any-role",
			expectFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewPermissionCache()

			perms, found := cache.Get(tt.role)

			assert.Equal(t, tt.expectFound, found)
			assert.Nil(t, perms)
		})
	}
}

func TestPermissionCache_Overwrite(t *testing.T) {
	tests := []struct {
		name                string
		role                string
		firstPermissions    []string
		secondPermissions   []string
		expectedPermissions []string
	}{
		{
			name:                "overwrite_with_more_permissions",
			role:                "editor",
			firstPermissions:    []string{"plan:read"},
			secondPermissions:   []string{"plan:read", "plan:create", "plan:update"},
			expectedPermissions: []string{"plan:read", "plan:create", "plan:update"},
		},
		{
			name:                "overwrite_with_fewer_permissions",
			role:                "admin",
			firstPermissions:    []string{"*", "plan:create", "user:create"},
			secondPermissions:   []string{"*"},
			expectedPermissions: []string{"*"},
		},
		{
			name:                "overwrite_with_empty",
			role:                "restricted",
			firstPermissions:    []string{"plan:read"},
			secondPermissions:   []string{},
			expectedPermissions: []string{},
		},
		{
			name:                "overwrite_with_nil",
			role:                "none",
			firstPermissions:    []string{"plan:read"},
			secondPermissions:   nil,
			expectedPermissions: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewPermissionCache()

			// Set first permissions
			cache.Set(tt.role, tt.firstPermissions)

			// Overwrite with second permissions
			cache.Set(tt.role, tt.secondPermissions)

			// Get permissions
			perms, found := cache.Get(tt.role)

			assert.True(t, found)
			assert.Equal(t, tt.expectedPermissions, perms)
		})
	}
}

func TestPermissionCache_Invalidate(t *testing.T) {
	tests := []struct {
		name                  string
		setupRoles            map[string][]string
		verifyAfterInvalidate bool
	}{
		{
			name: "invalidate_single_role",
			setupRoles: map[string][]string{
				"admin": {"*"},
			},
			verifyAfterInvalidate: true,
		},
		{
			name: "invalidate_multiple_roles",
			setupRoles: map[string][]string{
				"admin":    {"*"},
				"planner":  {"plan:create", "plan:read"},
				"operator": {"plan:read"},
			},
			verifyAfterInvalidate: true,
		},
		{
			name:                  "invalidate_empty_cache",
			setupRoles:            map[string][]string{},
			verifyAfterInvalidate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := NewPermissionCache()

			// Setup roles
			for role, perms := range tt.setupRoles {
				cache.Set(role, perms)
			}

			// Verify roles exist before invalidation
			for role := range tt.setupRoles {
				_, found := cache.Get(role)
				assert.True(t, found, "Role %s should exist before invalidation", role)
			}

			// Invalidate cache
			cache.Invalidate()

			// Verify all roles are gone
			if tt.verifyAfterInvalidate {
				for role := range tt.setupRoles {
					_, found := cache.Get(role)
					assert.False(t, found, "Role %s should not exist after invalidation", role)
				}
			}

			// Verify cache is empty
			assert.Equal(t, 0, len(cache.permissions))
		})
	}
}

func TestPermissionCache_MultipleInvalidations(t *testing.T) {
	cache := NewPermissionCache()

	// Set some data
	cache.Set("role1", []string{"perm1"})
	cache.Set("role2", []string{"perm2"})

	// Invalidate multiple times
	cache.Invalidate()
	cache.Invalidate()
	cache.Invalidate()

	// Should still be empty and functional
	_, found := cache.Get("role1")
	assert.False(t, found)

	// Should be able to set after multiple invalidations
	cache.Set("role3", []string{"perm3"})
	perms, found := cache.Get("role3")
	assert.True(t, found)
	assert.Equal(t, []string{"perm3"}, perms)
}

func TestPermissionCache_ConcurrentAccess(t *testing.T) {
	cache := NewPermissionCache()
	const numGoroutines = 100
	const numOperations = 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // readers, writers, invalidators

	// Concurrent readers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cache.Get("admin")
			}
		}(i)
	}

	// Concurrent writers
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cache.Set("admin", []string{"*"})
			}
		}(i)
	}

	// Concurrent invalidators
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				cache.Invalidate()
			}
		}(i)
	}

	wg.Wait()

	// If we get here without race conditions, test passes
	assert.NotNil(t, cache)
}

func TestPermissionCache_SetAfterInvalidate(t *testing.T) {
	cache := NewPermissionCache()

	// Set initial data
	cache.Set("admin", []string{"*"})
	cache.Set("planner", []string{"plan:create"})

	// Invalidate
	cache.Invalidate()

	// Set new data after invalidation
	cache.Set("operator", []string{"plan:read"})
	cache.Set("admin", []string{"*", "user:create"})

	// Verify new data
	perms, found := cache.Get("operator")
	assert.True(t, found)
	assert.Equal(t, []string{"plan:read"}, perms)

	perms, found = cache.Get("admin")
	assert.True(t, found)
	assert.Equal(t, []string{"*", "user:create"}, perms)

	// Verify old data is gone
	_, found = cache.Get("planner")
	assert.False(t, found)
}
