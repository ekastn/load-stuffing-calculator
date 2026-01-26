-- PLATFORM / FOUNDER QUERIES

-- name: CountGlobalUsers :one
SELECT COUNT(*) FROM users;

-- name: CountGlobalActivePlans :one
SELECT COUNT(*) FROM load_plans WHERE status NOT IN ('COMPLETED', 'FAILED');

-- name: CountGlobalContainers :one
SELECT COUNT(*) FROM containers;

-- name: GetGlobalPlanStatusDistribution :many
SELECT status, COUNT(*) as count 
FROM load_plans 
GROUP BY status;

-- name: CountGlobalCompletedPlansToday :one
SELECT COUNT(*) FROM load_plans 
WHERE status = 'COMPLETED' 
AND created_at >= CURRENT_DATE;

-- name: CountGlobalTotalItems :one
SELECT COALESCE(SUM(quantity), 0)::BIGINT FROM load_items;

-- name: GetGlobalAvgVolumeUtilization :one
SELECT COALESCE(AVG(volume_utilization_pct), 0)::FLOAT
FROM plan_results;

-- name: CountGlobalCompletedPlans :one
SELECT COUNT(*) FROM load_plans WHERE status = 'COMPLETED';


-- WORKSPACE SCOPED QUERIES

-- name: CountWorkspaceMembers :one
SELECT COUNT(*) 
FROM members 
WHERE workspace_id = $1;

-- name: CountWorkspaceActivePlans :one
SELECT COUNT(*) 
FROM load_plans 
WHERE workspace_id = $1 
AND status NOT IN ('COMPLETED', 'FAILED');

-- name: CountWorkspaceContainers :one
SELECT COUNT(*) 
FROM containers 
WHERE workspace_id = $1 OR workspace_id IS NULL;

-- name: GetWorkspacePlanStatusDistribution :many
SELECT status, COUNT(*) as count 
FROM load_plans 
WHERE workspace_id = $1
GROUP BY status;

-- name: CountWorkspaceCompletedPlansToday :one
SELECT COUNT(*) 
FROM load_plans 
WHERE workspace_id = $1 
AND status = 'COMPLETED' 
AND created_at >= CURRENT_DATE;

-- name: CountWorkspaceItems :one
SELECT COALESCE(SUM(li.quantity), 0)::BIGINT 
FROM load_items li
JOIN load_plans lp ON li.plan_id = lp.plan_id
WHERE lp.workspace_id = $1;

-- name: GetWorkspaceAvgVolumeUtilization :one
SELECT COALESCE(AVG(pr.volume_utilization_pct), 0)::FLOAT
FROM plan_results pr
JOIN load_plans lp ON pr.plan_id = lp.plan_id
WHERE lp.workspace_id = $1;

-- name: CountWorkspaceCompletedPlans :one
SELECT COUNT(*) 
FROM load_plans 
WHERE workspace_id = $1 
AND status = 'COMPLETED';
