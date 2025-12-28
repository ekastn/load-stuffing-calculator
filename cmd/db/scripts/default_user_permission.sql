-- Default roles
INSERT INTO roles (name, description) VALUES
('founder', 'Platform superuser (SaaS founder)'),
('owner', 'Workspace owner / CEO'),
('personal', 'Personal workspace owner (single-user)'),
('admin', 'Workspace admin (can manage members)'),
('planner', 'Can create and manage shipment plans'),
('operator', 'Can validate loading steps & manage items'),
('trial', 'Anonymous trial users (limited)')
ON CONFLICT (name) DO NOTHING;

-- Default permissions (dev-mode)
INSERT INTO permissions (name, description) VALUES
('*', 'Global access'),

('user:*', 'Full access to users'),
('user:read', 'Read users'),
('user:create', 'Create users'),
('user:update', 'Update users'),
('user:delete', 'Delete users'),

('role:*', 'Full access to roles'),
('role:read', 'Read roles'),
('role:create', 'Create roles'),
('role:update', 'Update roles'),
('role:delete', 'Delete roles'),

('permission:*', 'Full access to permissions'),
('permission:read', 'Read permissions'),
('permission:create', 'Create permissions'),
('permission:update', 'Update permissions'),
('permission:delete', 'Delete permissions'),

('product:*', 'Full access to products'),
('product:read', 'Read products'),
('product:create', 'Create products'),
('product:update', 'Update products'),
('product:delete', 'Delete products'),

('container:*', 'Full access to containers'),
('container:read', 'Read containers'),
('container:create', 'Create containers'),
('container:update', 'Update containers'),
('container:delete', 'Delete containers'),

('workspace:*', 'Full access to workspaces'),
('workspace:read', 'Read workspaces'),
('workspace:create', 'Create workspaces'),
('workspace:update', 'Update workspaces'),
('workspace:delete', 'Delete workspaces'),

('member:*', 'Full access to workspace members'),
('member:read', 'Read workspace members'),
('member:create', 'Add workspace members'),
('member:update', 'Update workspace members'),
('member:delete', 'Remove workspace members'),

('invite:*', 'Full access to invites'),
('invite:read', 'Read invites'),
('invite:create', 'Create invites'),
('invite:delete', 'Revoke invites'),
('invite:accept', 'Accept invites'),

('plan:*', 'Full access to plans'),
('plan:read', 'Read plans'),
('plan:create', 'Create plans'),
('plan:update', 'Update plans'),
('plan:delete', 'Delete plans'),
('plan:calculate', 'Calculate plan placements'),

('plan_item:*', 'Full access to plan items'),
('plan_item:read', 'Read plan items'),
('plan_item:create', 'Create plan items'),
('plan_item:update', 'Update plan items'),
('plan_item:delete', 'Delete plan items'),

('dashboard:read', 'View dashboard')
ON CONFLICT (name) DO NOTHING;

-- Founder gets global "*"
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name = '*'
WHERE r.name = 'founder'
ON CONFLICT DO NOTHING;

-- Owner permissions (workspace-scoped; no global "*")
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN (
  'workspace:*',
  'member:*',
  'invite:*',
  'product:*',
  'container:*',
  'plan:*',
  'plan_item:*',
  'dashboard:read'
)
WHERE r.name = 'owner'
ON CONFLICT DO NOTHING;

-- Personal permissions (workspace-scoped; no members/invites/workspace:create)
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN (
  'workspace:read',
  'product:*',
  'container:*',
  'plan:*',
  'plan_item:*',
  'dashboard:read'
)
WHERE r.name = 'personal'
ON CONFLICT DO NOTHING;

-- Admin permissions (workspace-scoped; no global "*")
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN (
  'workspace:read',
  'member:*',
  'invite:*',
  'product:*',
  'container:*',
  'plan:*',
  'plan_item:*',
  'dashboard:read'
)
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Planner permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN ('workspace:read', 'plan:*', 'plan_item:*', 'product:read', 'container:read', 'dashboard:read')
WHERE r.name = 'planner'
ON CONFLICT DO NOTHING;

-- Operator permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN ('workspace:read', 'plan:read', 'plan_item:*', 'product:read', 'container:read', 'dashboard:read')
WHERE r.name = 'operator'
ON CONFLICT DO NOTHING;

-- Trial permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN ('workspace:read', 'plan:*', 'plan_item:*', 'product:read', 'container:read')
WHERE r.name = 'trial'
ON CONFLICT DO NOTHING;
