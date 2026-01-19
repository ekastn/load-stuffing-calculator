-- name: GetProduct :one
SELECT * FROM products WHERE id = $1;

-- name: GetProductBySKU :one
SELECT * FROM products WHERE sku = $1;

-- name: ListProducts :many
SELECT * FROM products ORDER BY label ASC;

-- name: CreateProduct :one
INSERT INTO products (label, sku, length_mm, width_mm, height_mm, weight_kg)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET label = $2, sku = $3, length_mm = $4, width_mm = $5, height_mm = $6, weight_kg = $7
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :exec
DELETE FROM products WHERE id = $1;
