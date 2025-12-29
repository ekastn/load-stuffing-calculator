-- name: CreateWorkspace :one
INSERT INTO workspaces (
    type,
    name,
    owner_user_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetWorkspace :one
SELECT *
FROM workspaces
WHERE workspace_id = $1;

-- name: GetPersonalWorkspaceByOwner :one
SELECT *
FROM workspaces
WHERE owner_user_id = $1
  AND type = 'personal'
LIMIT 1;

-- name: ListWorkspacesForUser :many
SELECT w.*
FROM workspaces w
JOIN members m ON m.workspace_id = w.workspace_id
WHERE m.user_id = $1
ORDER BY w.created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListWorkspacesAll :many
SELECT
    w.*,
    u.username AS owner_username,
    u.email AS owner_email
FROM workspaces w
JOIN users u ON u.user_id = w.owner_user_id
ORDER BY w.created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListWorkspacesByOwner :many
SELECT *
FROM workspaces
WHERE owner_user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateWorkspace :exec
UPDATE workspaces
SET
    name = $2,
    updated_at = NOW()
WHERE workspace_id = $1;

-- name: TransferWorkspaceOwnership :exec
UPDATE workspaces
SET
    owner_user_id = $2,
    updated_at = NOW()
WHERE workspace_id = $1;

-- name: DeleteWorkspace :exec
DELETE FROM workspaces
WHERE workspace_id = $1;
