-- name: CreateBudget :one
INSERT INTO budgets (
    user_id,
    category,
    period,
    limit_cents,
    current_spent_cents,
    alert_threshold,
    alerts_enabled,
    start_date,
    end_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetBudgetByID :one
SELECT * FROM budgets
WHERE id = $1
LIMIT 1;

-- name: GetBudgetByCategoryAndPeriod :one
SELECT * FROM budgets
WHERE user_id = $1
  AND category = $2
  AND period = $3
LIMIT 1;

-- name: GetBudgetForUpdate :one
SELECT * FROM budgets
WHERE id = $1
FOR UPDATE;

-- name: ListUserBudgets :many
SELECT * FROM budgets
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUserBudgetsByCategory :many
SELECT * FROM budgets
WHERE user_id = $1 AND category = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: ListUserBudgetsByPeriod :many
SELECT * FROM budgets
WHERE user_id = $1 AND period = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUserBudgets :one
SELECT COUNT(*) FROM budgets
WHERE user_id = $1;

-- name: UpdateBudget :one
UPDATE budgets
SET
    limit_cents = $2,
    alert_threshold = $3,
    alerts_enabled = $4
WHERE id = $1
RETURNING *;

-- name: UpdateBudgetSpent :one
UPDATE budgets
SET current_spent_cents = $2
WHERE id = $1
RETURNING *;

-- name: IncrementBudgetSpent :one
UPDATE budgets
SET current_spent_cents = current_spent_cents + $2
WHERE id = $1
RETURNING *;

-- name: ResetBudgetSpent :exec
UPDATE budgets
SET current_spent_cents = 0
WHERE user_id = $1;

-- name: DeleteBudget :exec
DELETE FROM budgets
WHERE id = $1;

-- name: GetBudgetsNearLimit :many
SELECT * FROM budgets
WHERE user_id = $1
  AND alerts_enabled = true
  AND (current_spent_cents::float / limit_cents::float * 100) >= alert_threshold
ORDER BY created_at DESC;

-- name: GetOverBudgets :many
SELECT * FROM budgets
WHERE user_id = $1
  AND current_spent_cents > limit_cents
ORDER BY created_at DESC;

-- name: GetUserBudgetsAnalytics :one
SELECT
    COUNT(*) as total_budgets,
    COALESCE(SUM(limit_cents), 0)::bigint as total_budget_cents,
    COALESCE(SUM(current_spent_cents), 0)::bigint as total_spent_cents,
    COUNT(*) FILTER (WHERE current_spent_cents > limit_cents) as over_budget_count
FROM budgets
WHERE user_id = $1;
