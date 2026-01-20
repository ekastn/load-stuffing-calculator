-- name: GetDashboardStats :one
SELECT
    (SELECT COUNT(*) FROM plans p WHERE p.user_id = $1) AS total_plans,
    (SELECT COUNT(*) FROM containers) AS total_containers,
    (SELECT COUNT(*) FROM products) AS total_products,
    (SELECT COALESCE(SUM(quantity), 0) FROM plan_items pi
     JOIN plans p ON pi.plan_id = p.id
     WHERE p.user_id = $1) AS total_items_shipped;
