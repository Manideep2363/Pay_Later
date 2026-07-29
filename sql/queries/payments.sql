-- name: CreatePayment :exec
INSERT INTO payments (
    user_id,
    amount
)
VALUES (?, ?);

-- name: ListUserPayments :many
SELECT
    payment_id,
    user_id,
    amount,
    paid_at
FROM payments
WHERE user_id = ?
ORDER BY paid_at DESC;

-- name: GetPaymentByID :one
SELECT
    payment_id,
    user_id,
    amount,
    paid_at
FROM payments
WHERE payment_id = ?;