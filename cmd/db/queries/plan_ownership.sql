-- name: CountPlansByCreator :one
SELECT COUNT(*)
FROM load_plans
WHERE created_by_type = sqlc.arg(created_by_type)
  AND created_by_id = sqlc.arg(created_by_id);

-- name: ClaimPlansFromGuest :exec
UPDATE load_plans
SET created_by_type = 'user',
    created_by_id = sqlc.arg(user_id),
    workspace_id = sqlc.arg(workspace_id)
WHERE created_by_type = 'guest'
  AND created_by_id = sqlc.arg(guest_id);
