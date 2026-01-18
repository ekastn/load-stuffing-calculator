-- name: GetPlan :one
SELECT * FROM plans WHERE id = $1;

-- name: CreatePlan :one
INSERT INTO plans (container_id, status)
VALUES ($1, 'draft')
RETURNING *;

-- name: UpdatePlanStatus :one
UPDATE plans SET status = $2, calculated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: AddPlanItem :one
INSERT INTO plan_items (plan_id, product_id, quantity)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetPlanItems :many
SELECT pi.*, p.label, p.length_mm, p.width_mm, p.height_mm, p.weight_kg
FROM plan_items pi
JOIN products p ON pi.product_id = p.id
WHERE pi.plan_id = $1;

-- name: SavePlacement :exec
INSERT INTO placements (plan_id, product_id, pos_x, pos_y, pos_z, rotation, step_number)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: GetPlacements :many
SELECT pl.*, p.label
FROM placements pl
JOIN products p ON pl.product_id = p.id
WHERE pl.plan_id = $1
ORDER BY pl.step_number ASC;

-- name: DeletePlanPlacements :exec
DELETE FROM placements WHERE plan_id = $1;
