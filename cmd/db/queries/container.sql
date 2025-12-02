-- name: CreateContainer :one
INSERT INTO containers (
    name, 
    inner_length_mm, 
    inner_width_mm, 
    inner_height_mm, 
    max_weight_kg, 
    description
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetContainer :one
SELECT *
FROM containers
WHERE container_id = $1;

-- name: ListContainers :many
SELECT *
FROM containers
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateContainer :exec
UPDATE containers 
SET 
    name = $2,
    inner_length_mm = $3,
    inner_width_mm = $4,
    inner_height_mm = $5,
    max_weight_kg = $6,
    description = $7,
    updated_at = NOW()
WHERE container_id = $1;

-- name: DeleteContainer :exec
DELETE FROM containers 
WHERE container_id = $1;
