-- name: CreateCardTransaction :one
INSERT INTO card_transactions (
    card_id,
    user_id,
    amount_cents,
    merchant_name,
    merchant_category,
    status,
    is_international,
    transaction_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: GetCardTransactionByID :one
SELECT * FROM card_transactions
WHERE id = $1
LIMIT 1;

-- name: ListCardTransactions :many
SELECT * FROM card_transactions
WHERE card_id = $1
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3;

-- name: ListUserCardTransactions :many
SELECT * FROM card_transactions
WHERE user_id = $1
ORDER BY transaction_date DESC
LIMIT $2 OFFSET $3;

-- name: CountCardTransactions :one
SELECT COUNT(*) FROM card_transactions
WHERE card_id = $1;

-- name: GetCardTransactionsByCategory :many
SELECT
    merchant_category,
    COUNT(*) as transaction_count,
    SUM(amount_cents) as total_amount_cents
FROM card_transactions
WHERE user_id = $1
  AND transaction_date >= sqlc.arg(start_date)
  AND transaction_date <= sqlc.arg(end_date)
GROUP BY merchant_category
ORDER BY total_amount_cents DESC;

-- name: GetCardTransactionsByDateRange :many
SELECT * FROM card_transactions
WHERE user_id = $1
  AND transaction_date >= sqlc.arg(start_date)
  AND transaction_date <= sqlc.arg(end_date)
ORDER BY transaction_date DESC;
