-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    user_id,
    action,
    resource_type,
    resource_id,
    old_values,
    new_values,
    ip_address,
    user_agent,
    request_id,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
RETURNING *;

-- name: GetAuditLogsByUserID :many
SELECT * FROM audit_logs
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAuditLogsByResource :many
SELECT * FROM audit_logs
WHERE resource_type = $1 AND resource_id = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: GetAuditLogsByRequestID :one
SELECT * FROM audit_logs
WHERE request_id = $1
LIMIT 1;

-- name: CountAuditLogsByUser :one
SELECT COUNT(*) FROM audit_logs
WHERE user_id = $1;
