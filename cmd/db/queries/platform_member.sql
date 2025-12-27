-- name: GetPlatformRoleByUserID :one
SELECT r.name AS role_name
FROM platform_members pm
JOIN roles r ON r.role_id = pm.role_id
WHERE pm.user_id = $1;

-- name: UpsertPlatformMember :exec
INSERT INTO platform_members (user_id, role_id)
VALUES ($1, $2)
ON CONFLICT (user_id) DO UPDATE
SET role_id = EXCLUDED.role_id;
