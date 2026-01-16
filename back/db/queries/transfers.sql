-- name: CreateTransfer :one
INSERT INTO transfers (
    user_id,
    type,
    status,
    amount_cents,
    fee_cents,
    currency,
    pix_key,
    pix_key_type,
    recipient_name,
    recipient_document,
    recipient_bank,
    recipient_branch,
    recipient_account,
    recipient_account_type,
    recipient_user_id,
    scheduled_for,
    completed_at,
    failure_reason,
    authentication_code
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19
)
RETURNING *;

-- name: GetTransferByID :one
SELECT * FROM transfers
WHERE id = $1
LIMIT 1;

-- name: ListUserTransfers :many
SELECT * FROM transfers
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserTransfers :one
SELECT COUNT(*) FROM transfers
WHERE user_id = $1;

-- name: UpdateTransferStatus :one
UPDATE transfers
SET
    status = $2,
    completed_at = CASE
        WHEN $2 = 'completed' THEN NOW()
        ELSE completed_at
    END,
    failure_reason = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CancelTransfer :one
UPDATE transfers
SET
    status = 'cancelled',
    updated_at = NOW()
WHERE id = $1 AND status = 'pending'
RETURNING *;

-- name: GetTransferForUpdate :one
SELECT * FROM transfers
WHERE id = $1
FOR UPDATE;

-- name: GetDailyTransferSum :one
SELECT COALESCE(SUM(amount_cents + fee_cents), 0)::bigint as total
FROM transfers
WHERE user_id = $1
  AND type IN ('pix', 'ted', 'p2p')
  AND status IN ('completed', 'processing', 'pending')
  AND created_at >= CURRENT_DATE
  AND created_at < CURRENT_DATE + INTERVAL '1 day';

-- name: GetMonthlyTransferSum :one
SELECT COALESCE(SUM(amount_cents + fee_cents), 0)::bigint as total
FROM transfers
WHERE user_id = $1
  AND type IN ('pix', 'ted', 'p2p')
  AND status IN ('completed', 'processing', 'pending')
  AND created_at >= DATE_TRUNC('month', CURRENT_DATE)
  AND created_at < DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month';

-- name: ListUserTransfersByStatus :many
SELECT * FROM transfers
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;

