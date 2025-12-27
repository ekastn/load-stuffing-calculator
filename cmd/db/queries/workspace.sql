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
