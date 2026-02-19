-- name: CreateProduct :one
INSERT INTO products (
    workspace_id,
    name,
    sku,
    length_mm,
    width_mm,
    height_mm,
    weight_kg,
    color_hex
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: ListProductsAll :many
SELECT *
FROM products
ORDER BY (workspace_id IS NULL) DESC, name
LIMIT $1 OFFSET $2;

-- name: GetProductAny :one
SELECT *
FROM products
WHERE product_id = $1;

-- name: GetProduct :one
SELECT *
FROM products
WHERE product_id = $1
  AND (workspace_id = $2 OR workspace_id IS NULL);

-- name: ListProducts :many
SELECT *
FROM products
WHERE workspace_id = $1 OR workspace_id IS NULL
ORDER BY (workspace_id IS NULL) DESC, name
LIMIT $2 OFFSET $3;

-- name: UpdateProduct :exec
UPDATE products
SET
    name = $3,
    sku = $4,
    length_mm = $5,
    width_mm = $6,
    height_mm = $7,
    weight_kg = $8,
    color_hex = $9,
    updated_at = NOW()
WHERE product_id = $1
  AND workspace_id = $2;

-- name: UpdateProductAny :exec
UPDATE products
SET
    name = $2,
    sku = $3,
    length_mm = $4,
    width_mm = $5,
    height_mm = $6,
    weight_kg = $7,
    color_hex = $8,
    updated_at = NOW()
WHERE product_id = $1;

-- name: DeleteProduct :exec
DELETE FROM products
WHERE product_id = $1
  AND workspace_id = $2;

-- name: DeleteProductAny :exec
DELETE FROM products
WHERE product_id = $1;
