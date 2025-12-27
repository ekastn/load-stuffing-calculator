-- name: CreateUser :one
INSERT INTO users (
    role_id, 
    username, 
    email, 
    password_hash
) VALUES (
    $1, 
    $2, 
    $3, 
    $4
)
RETURNING  *;

-- name: ListUsers :many
SELECT 
    u.user_id,
    u.username,
    u.email,
    r.name AS role_name,
    up.full_name,
    up.phone_number,
    u.created_at
FROM users u
JOIN roles r ON u.role_id = r.role_id
LEFT JOIN user_profiles up ON up.user_id = u.user_id
ORDER BY u.created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetUserByID :one
SELECT 
    u.user_id,
    u.username,
    u.email,
    u.role_id,
    r.name AS role_name,
    u.created_at,
    up.full_name,
    up.gender,
    up.date_of_birth,
    up.phone_number,
    up.address,
    up.avatar_url
FROM users u
JOIN roles r ON u.role_id = r.role_id
LEFT JOIN user_profiles up ON up.user_id = u.user_id
WHERE u.user_id = $1;

-- name: UpdateUser :exec
UPDATE users 
SET 
    role_id = $2,
    email = $3,
    username = $4,
    updated_at = NOW()
WHERE user_id = $1;

-- name: GetUserByUsername :one
SELECT u.user_id, u.username, u.email, u.password_hash, u.role_id, r.name AS role_name
FROM users u
JOIN roles r ON u.role_id = r.role_id
WHERE u.username = $1;

-- name: GetUserByEmail :one
SELECT u.user_id, u.username, u.email, u.password_hash, u.role_id, r.name AS role_name
FROM users u
JOIN roles r ON u.role_id = r.role_id
WHERE u.email = $1;

-- name: GetRoleByName :one
SELECT role_id, name, description FROM roles
WHERE name = $1;

-- name: DeleteUser :exec
DELETE FROM users WHERE user_id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $2, updated_at = NOW() WHERE user_id = $1;