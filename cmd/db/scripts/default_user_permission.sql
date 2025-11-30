-- Default roles
INSERT INTO roles (name, description) VALUES
('admin', 'Full system access'),
('planner', 'Can create and manage shipment plans'),
('operator', 'Can execute loading using IoT devices')
ON CONFLICT (name) DO NOTHING;

-- Contoh permissions
INSERT INTO permissions (name, description) VALUES
('user:create', 'Create new user'),
('user:read', 'Read user data'),
('user:update', 'Update user'),
('user:delete', 'Delete/deactivate user'),
('plan:*', 'Full access to plans'),
('execution:*', 'Full access to execution')
ON CONFLICT (name) DO NOTHING;

-- Admin punya semua
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.role_id, p.permission_id
FROM roles r, permissions p
WHERE r.name = 'admin'
ON CONFLICT DO NOTHING;
