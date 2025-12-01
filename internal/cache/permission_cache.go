package cache

import (
	"sync"
)

type PermissionCache struct {
	mu          sync.RWMutex
	permissions map[string][]string
}

func NewPermissionCache() *PermissionCache {
	return &PermissionCache{
		permissions: make(map[string][]string),
	}
}

func (c *PermissionCache) Get(role string) ([]string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	perms, ok := c.permissions[role]
	return perms, ok
}

func (c *PermissionCache) Set(role string, permissions []string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.permissions[role] = permissions
}

// Invalidate clears the cache (useful if permissions change at runtime)
func (c *PermissionCache) Invalidate() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.permissions = make(map[string][]string)
}
