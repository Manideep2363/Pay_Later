-- name: CheckEmailExists :one
SELECT EXISTS(
    SELECT 1
    FROM users
    WHERE email = ?
);