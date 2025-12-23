-- Default roles
INSERT INTO roles (name, description) VALUES
('admin', 'Full system access'),
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

-- Admin gets only global "*"
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name = '*'
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;

-- Planner permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN ('plan:*', 'plan_item:*', 'product:read', 'container:read', 'dashboard:read')
WHERE r.name = 'planner'
ON CONFLICT DO NOTHING;

-- Operator permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN ('plan:read', 'plan_item:*', 'product:read', 'container:read', 'dashboard:read')
WHERE r.name = 'operator'
ON CONFLICT DO NOTHING;

-- Trial permissions
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r
JOIN permissions p ON p.name IN ('plan:*', 'plan_item:*', 'product:read', 'container:read')
WHERE r.name = 'trial'
ON CONFLICT DO NOTHING;
