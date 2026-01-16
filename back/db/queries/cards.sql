-- ========================================
-- CARDS QUERIES
-- ========================================

-- name: CreateCard :one
INSERT INTO cards (
    user_id,
    type,
    brand,
    status,
    card_number_encrypted,
    cvv_encrypted,
    pin_hash,
    last_four_digits,
    holder_name,
    expiry_month,
    expiry_year,
    daily_limit_cents,
    monthly_limit_cents,
    current_daily_spent_cents,
    current_monthly_spent_cents,
    is_contactless,
    is_international,
    block_international,
    block_online,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10,
    $11, $12, $13, $14, $15, $16, $17, $18, $19, $20
)
RETURNING *;

-- name: GetCardByID :one
SELECT * FROM cards
WHERE id = $1
LIMIT 1;

-- name: GetCardForUpdate :one
SELECT * FROM cards
WHERE id = $1
FOR UPDATE;

-- name: ListUserCards :many
SELECT * FROM cards
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: ListActiveUserCards :many
SELECT * FROM cards
WHERE user_id = $1 AND status = 'active'
ORDER BY created_at DESC;

-- name: UpdateCardStatus :exec
UPDATE cards
SET
    status = $2,
    blocked_at = CASE
        WHEN $2 = 'blocked' THEN NOW()
        WHEN $2 = 'active' THEN NULL
        ELSE blocked_at
    END,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateCardLimits :exec
UPDATE cards
SET
    daily_limit_cents = $2,
    monthly_limit_cents = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateCardSecuritySettings :exec
UPDATE cards
SET
    is_contactless = $2,
    is_international = $3,
    block_international = $4,
    block_online = $5,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateCardPIN :exec
UPDATE cards
SET
    pin_hash = $2,
    updated_at = NOW()
WHERE id = $1;

-- name: UpdateCardSpentAmounts :exec
UPDATE cards
SET
    current_daily_spent_cents = $2,
    current_monthly_spent_cents = $3,
    updated_at = NOW()
WHERE id = $1;

-- name: ResetDailySpent :exec
UPDATE cards
SET
    current_daily_spent_cents = 0,
    updated_at = NOW()
WHERE id = $1;

-- name: ResetMonthlySpent :exec
UPDATE cards
SET
    current_monthly_spent_cents = 0,
    updated_at = NOW()
WHERE id = $1;

-- name: ResetAllDailySpent :exec
UPDATE cards
SET
    current_daily_spent_cents = 0,
    updated_at = NOW()
WHERE status = 'active';

-- name: ResetAllMonthlySpent :exec
UPDATE cards
SET
    current_monthly_spent_cents = 0,
    updated_at = NOW()
WHERE status = 'active';

-- name: DeleteCard :exec
UPDATE cards
SET
    status = 'cancelled',
    updated_at = NOW()
WHERE id = $1;

-- name: CountUserCards :one
SELECT COUNT(*) FROM cards
WHERE user_id = $1;

-- name: CountUserActiveCards :one
SELECT COUNT(*) FROM cards
WHERE user_id = $1 AND status = 'active';

-- name: GetUserCardsByStatus :many
SELECT * FROM cards
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC;
