-- name: CreateAdminUser :one
INSERT INTO admin_users (email,
                         password_hash,
                         display_name,
                         is_active)
VALUES ($1, $2, $3, $4)
RETURNING
    id,
    email,
    password_hash,
    display_name,
    is_active,
    last_login_at,
    created_at,
    updated_at;


-- name: GetAdminUserByEmail :one
SELECT id,
       email,
       password_hash,
       display_name,
       is_active,
       last_login_at,
       created_at,
       updated_at
FROM admin_users
WHERE email = $1;


-- name: GetAdminUserById :one
SELECT id,
       email,
       password_hash,
       display_name,
       is_active,
       last_login_at,
       created_at,
       updated_at
FROM admin_users
WHERE id = $1;


-- name: CreateAdminSession :one
INSERT INTO admin_sessions (admin_user_id,
                            token_hash,
                            expires_at)
VALUES ($1, $2, $3)
RETURNING
    id,
    admin_user_id,
    token_hash,
    expires_at,
    revoked_at,
    created_at,
    last_used_at;


-- name: GetAdminSessionByTokenHash :one
SELECT id,
       admin_user_id,
       token_hash,
       expires_at,
       revoked_at,
       created_at,
       last_used_at
FROM admin_sessions
WHERE token_hash = $1;


-- name: RevokeAdminSession :exec
UPDATE admin_sessions
SET revoked_at = $1
WHERE token_hash = $2;
