-- name: CreateProduct :one
INSERT INTO products (
    name, 
    length_mm, 
    width_mm, 
    height_mm, 
    weight_kg, 
    color_hex
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetProduct :one
SELECT *
FROM products
WHERE product_id = $1;

-- name: ListProducts :many
SELECT *
FROM products
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateProduct :exec
UPDATE products 
SET 
    name = $2,
    length_mm = $3,
    width_mm = $4,
    height_mm = $5,
    weight_kg = $6,
    color_hex = $7,
    updated_at = NOW()
WHERE product_id = $1;

-- name: DeleteProduct :exec
DELETE FROM products 
WHERE product_id = $1;
