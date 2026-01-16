-- name: CreateBill :one
INSERT INTO bills (
    user_id,
    type,
    status,
    barcode,
    amount_cents,
    fee_cents,
    final_amount_cents,
    recipient_name,
    due_date
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
)
RETURNING *;

-- name: GetBillByID :one
SELECT * FROM bills
WHERE id = $1
LIMIT 1;

-- name: GetBillByBarcode :one
SELECT * FROM bills
WHERE barcode = $1
LIMIT 1;

-- name: GetBillForUpdate :one
SELECT * FROM bills
WHERE id = $1
FOR UPDATE;

-- name: ListUserBills :many
SELECT * FROM bills
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListUserBillsByStatus :many
SELECT * FROM bills
WHERE user_id = $1 AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUserBills :one
SELECT COUNT(*) FROM bills
WHERE user_id = $1;

-- name: UpdateBillStatus :one
UPDATE bills
SET status = $2
WHERE id = $1
RETURNING *;

-- name: MarkBillAsPaid :one
UPDATE bills
SET
    status = 'paid',
    payment_date = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteBill :exec
UPDATE bills
SET status = 'cancelled'
WHERE id = $1;

-- name: ListOverdueBills :many
SELECT * FROM bills
WHERE status = 'pending'
  AND due_date < CURRENT_DATE
ORDER BY due_date ASC
LIMIT $1 OFFSET $2;

-- name: GetUserBillsStats :one
SELECT
    COUNT(*) as total_bills,
    COUNT(*) FILTER (WHERE status = 'pending') as pending_bills,
    COUNT(*) FILTER (WHERE status = 'paid') as paid_bills,
    COALESCE(SUM(final_amount_cents) FILTER (WHERE status = 'pending'), 0)::bigint as pending_amount_cents,
    COALESCE(SUM(final_amount_cents) FILTER (WHERE status = 'paid'), 0)::bigint as paid_amount_cents
FROM bills
WHERE user_id = $1;
