-- name: GetMerchantCommissionSummary :many
SELECT
    merchant_id,
    CAST(
        COALESCE(SUM(commission_amount), 0.00)
        AS DECIMAL(10,2)
    ) AS total_commission
FROM transactions
GROUP BY merchant_id
ORDER BY total_commission DESC;

-- name: GetUserOutstandingDues :many
SELECT
    user_id,
    name,
    current_due
FROM users
ORDER BY current_due DESC;

-- name: GetUsersAtCreditLimit :many
SELECT
    user_id,
    name,
    credit_limit,
    current_due
FROM users
WHERE current_due >= credit_limit;

-- name: GetOutstandingBalance :one
SELECT
    CAST(COALESCE(SUM(current_due), 0.00) AS DECIMAL(10,2)) AS total_due
FROM users;