-- name: CreatePermission :one
INSERT INTO permissions (
    name, 
    description
) VALUES (
    $1, 
    $2
)
RETURNING *;

-- name: GetPermission :one
SELECT *
FROM permissions
WHERE permission_id = $1;

-- name: ListPermissions :many
SELECT *
FROM permissions
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdatePermission :exec
UPDATE permissions 
SET 
    name = $2,
    description = $3
WHERE permission_id = $1;

-- name: DeletePermission :exec
DELETE FROM permissions 
WHERE permission_id = $1;
