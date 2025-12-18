-- name: CountTotalUsers :one
SELECT COUNT(*) FROM users;

-- name: CountActivePlans :one
SELECT COUNT(*) FROM load_plans WHERE status NOT IN ('COMPLETED', 'FAILED');

-- name: CountContainers :one
SELECT COUNT(*) FROM containers;

-- name: GetPlanStatusDistribution :many
SELECT status, COUNT(*) as count 
FROM load_plans 
GROUP BY status;

-- name: CountCompletedPlansToday :one
SELECT COUNT(*) FROM load_plans 
WHERE status = 'COMPLETED' 
AND created_at >= CURRENT_DATE;

-- name: CountTotalItems :one
SELECT COALESCE(SUM(quantity), 0)::BIGINT FROM load_items;

-- name: GetAvgVolumeUtilization :one
SELECT COALESCE(AVG(volume_utilization_pct), 0)::FLOAT
FROM plan_results;

-- name: CountCompletedPlans :one
SELECT COUNT(*) FROM load_plans WHERE status = 'COMPLETED';