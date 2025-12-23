-- name: CreateLoadPlan :one
INSERT INTO load_plans (
    plan_code,
    status,
    cont_label,
    length_mm,
    width_mm,
    height_mm,
    max_weight_kg,
    created_by_type,
    created_by_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: AddLoadItem :one
INSERT INTO load_items (
    plan_id,
    item_label,
    length_mm, width_mm, height_mm,
    weight_kg,
    quantity,
    allow_rotation,
    color_hex
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetLoadPlan :one
SELECT * FROM load_plans WHERE plan_id = $1;

-- name: GetLoadPlanForGuest :one
SELECT *
FROM load_plans
WHERE plan_id = $1
  AND created_by_type = 'guest'
  AND created_by_id = $2;

-- name: ListLoadPlans :many
SELECT * FROM load_plans
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: ListLoadPlansForGuest :many
SELECT *
FROM load_plans
WHERE created_by_type = 'guest'
  AND created_by_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdatePlanStatus :exec
UPDATE load_plans
SET status = $2
WHERE plan_id = $1;

-- name: ListLoadItems :many
SELECT * FROM load_items
WHERE plan_id = $1;

-- name: GetLoadItem :one
SELECT * FROM load_items
WHERE plan_id = $1 AND item_id = $2;

-- name: UpdateLoadItem :exec
UPDATE load_items
SET
    item_label = $3,
    length_mm = $4,
    width_mm = $5,
    height_mm = $6,
    weight_kg = $7,
    quantity = $8,
    allow_rotation = $9,
    color_hex = $10
WHERE plan_id = $1 AND item_id = $2;

-- name: DeleteLoadItem :exec
DELETE FROM load_items
WHERE plan_id = $1 AND item_id = $2;

-- name: UpdateLoadPlan :exec
UPDATE load_plans
SET
    plan_code = $2,
    cont_label = $3,
    length_mm = $4,
    width_mm = $5,
    height_mm = $6,
    max_weight_kg = $7,
    status = $8
WHERE plan_id = $1;

-- name: DeleteLoadPlan :exec
DELETE FROM load_plans
WHERE plan_id = $1;

-- name: CreatePlanResult :one
INSERT INTO plan_results (
    plan_id,
    total_loaded_weight_kg,
    volume_utilization_pct,
    is_feasible
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: DeletePlanResults :exec
DELETE FROM plan_results WHERE plan_id = $1;

-- name: CreatePlanPlacement :copyfrom
INSERT INTO plan_placements (
    result_id,
    item_id,
    pos_x,
    pos_y,
    pos_z,
    rotation_code,
    step_number
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
);

-- name: GetPlanResult :one
SELECT * FROM plan_results WHERE plan_id = $1;

-- name: ListPlanPlacements :many
SELECT * FROM plan_placements WHERE result_id = $1 ORDER BY step_number ASC;