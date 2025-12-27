-- name: CreateInvite :one
INSERT INTO invites (
    workspace_id,
    email,
    role_id,
    token_hash,
    invited_by_user_id,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: ListInvitesByWorkspace :many
SELECT
    i.invite_id,
    i.workspace_id,
    i.email,
    i.role_id,
    r.name AS role_name,
    i.invited_by_user_id,
    u.username AS invited_by_username,
    i.expires_at,
    i.accepted_at,
    i.revoked_at,
    i.created_at
FROM invites i
JOIN roles r ON r.role_id = i.role_id
JOIN users u ON u.user_id = i.invited_by_user_id
WHERE i.workspace_id = $1
ORDER BY i.created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetInviteByTokenHash :one
SELECT *
FROM invites
WHERE token_hash = $1
  AND revoked_at IS NULL
  AND accepted_at IS NULL
  AND expires_at > NOW();

-- name: AcceptInvite :exec
UPDATE invites
SET accepted_at = NOW()
WHERE invite_id = $1
  AND workspace_id = $2
  AND accepted_at IS NULL
  AND revoked_at IS NULL;

-- name: RevokeInvite :exec
UPDATE invites
SET revoked_at = NOW()
WHERE invite_id = $1
  AND workspace_id = $2
  AND accepted_at IS NULL
  AND revoked_at IS NULL;
