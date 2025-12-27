-- name: CreateMember :one
INSERT INTO members (
    workspace_id,
    user_id,
    role_id
) VALUES (
    $1, $2, $3
)
RETURNING *;

-- name: GetMember :one
SELECT *
FROM members
WHERE member_id = $1;

-- name: GetMemberByWorkspaceAndUser :one
SELECT *
FROM members
WHERE workspace_id = $1
  AND user_id = $2;

-- name: GetMemberRoleNameByWorkspaceAndUser :one
SELECT r.name AS role_name
FROM members m
JOIN roles r ON r.role_id = m.role_id
WHERE m.workspace_id = $1
  AND m.user_id = $2;

-- name: ListMembersByWorkspace :many
SELECT
    m.member_id,
    m.workspace_id,
    m.user_id,
    m.role_id,
    r.name AS role_name,
    u.username,
    u.email,
    m.created_at,
    m.updated_at
FROM members m
JOIN users u ON u.user_id = m.user_id
JOIN roles r ON r.role_id = m.role_id
WHERE m.workspace_id = $1
ORDER BY m.created_at ASC
LIMIT $2 OFFSET $3;

-- name: UpdateMemberRole :exec
UPDATE members
SET
    role_id = $3,
    updated_at = NOW()
WHERE member_id = $1
  AND workspace_id = $2;

-- name: DeleteMember :exec
DELETE FROM members
WHERE member_id = $1
  AND workspace_id = $2;
