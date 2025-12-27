-- name: CreateContainer :one
INSERT INTO containers (
    workspace_id,
    name,
    inner_length_mm,
    inner_width_mm,
    inner_height_mm,
    max_weight_kg,
    description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetContainer :one
SELECT *
FROM containers
WHERE container_id = $1
  AND workspace_id = $2;

-- name: ListContainers :many
SELECT *
FROM containers
WHERE workspace_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateContainer :exec
UPDATE containers
SET
    name = $3,
    inner_length_mm = $4,
    inner_width_mm = $5,
    inner_height_mm = $6,
    max_weight_kg = $7,
    description = $8,
    updated_at = NOW()
WHERE container_id = $1
  AND workspace_id = $2;

-- name: DeleteContainer :exec
DELETE FROM containers
WHERE container_id = $1
  AND workspace_id = $2;
