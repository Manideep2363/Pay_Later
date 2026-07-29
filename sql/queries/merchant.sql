-- name: CreateMerchant :execresult
INSERT INTO merchants (
    name,
    phone,
    commission_percentage
)
VALUES (
    ?,
    ?,
    ?
);

-- name: ListMerchants :many
SELECT
    merchant_id,
    name,
    phone,
    commission_percentage
FROM merchants
ORDER BY merchant_id;

-- name: GetMerchantByID :one
SELECT
    merchant_id,
    name,
    phone,
    commission_percentage
FROM merchants
WHERE merchant_id = ?;

-- name: UpdateMerchantCommission :execresult
UPDATE merchants
SET commission_percentage = ?
WHERE merchant_id = ?;


-- name: GetMerchantByPhone :one
SELECT
    merchant_id,
    name,
    phone,
    commission_percentage
FROM merchants
WHERE phone = ?;