-- name: GetContainer :one
SELECT * FROM containers WHERE id = $1;

-- name: ListContainers :many
SELECT * FROM containers ORDER BY created_at DESC;

-- name: CreateContainer :one
INSERT INTO containers (name, length_mm, width_mm, height_mm, max_weight_kg)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateContainer :one
UPDATE containers
SET name = $2, length_mm = $3, width_mm = $4, height_mm = $5, max_weight_kg = $6
WHERE id = $1
RETURNING *;

-- name: DeleteContainer :exec
DELETE FROM containers WHERE id = $1;
