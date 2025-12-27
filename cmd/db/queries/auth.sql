-- name: CreateRefreshToken :exec
INSERT INTO refresh_tokens (
    token,
    user_id,
    workspace_id,
    expires_at
) VALUES (
    $1,
    $2,
    $3,
    $4
);

-- name: GetRefreshToken :one
SELECT
    user_id,
    workspace_id,
    expires_at,
    revoked_at
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
DELETE FROM refresh_tokens 
WHERE token = $1;

-- name: UpdateRefreshTokenWorkspace :exec
UPDATE refresh_tokens
SET
    workspace_id = $2,
    updated_at = NOW()
WHERE token = $1;

-- name: GetPermissionsByRole :many
SELECT p.name
FROM permissions p
JOIN role_permissions rp ON p.permission_id = rp.permission_id
JOIN roles r ON rp.role_id = r.role_id
WHERE r.name = $1;