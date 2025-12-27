-- name: CreateProduct :one
INSERT INTO products (
    workspace_id,
    name,
    length_mm,
    width_mm,
    height_mm,
    weight_kg,
    color_hex
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetProduct :one
SELECT *
FROM products
WHERE product_id = $1
  AND workspace_id = $2;

-- name: ListProducts :many
SELECT *
FROM products
WHERE workspace_id = $1
ORDER BY name
LIMIT $2 OFFSET $3;

-- name: UpdateProduct :exec
UPDATE products
SET
    name = $3,
    length_mm = $4,
    width_mm = $5,
    height_mm = $6,
    weight_kg = $7,
    color_hex = $8,
    updated_at = NOW()
WHERE product_id = $1
  AND workspace_id = $2;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_id = $1
  AND workspace_id = $2;
