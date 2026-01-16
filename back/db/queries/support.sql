-- Support Tickets Queries

-- name: CreateTicket :one
INSERT INTO support_tickets (
    user_id,
    ticket_number,
    category,
    priority,
    status,
    subject,
    description
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;

-- name: GetTicketByID :one
SELECT * FROM support_tickets
WHERE id = $1
LIMIT 1;

-- name: GetTicketByNumber :one
SELECT * FROM support_tickets
WHERE ticket_number = $1
LIMIT 1;

-- name: ListUserTickets :many
SELECT * FROM support_tickets
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserTickets :one
SELECT COUNT(*) FROM support_tickets
WHERE user_id = $1;

-- name: ListUserTicketsByStatus :many
SELECT * FROM support_tickets
WHERE user_id = $1
  AND status = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: CountUserTicketsByStatus :one
SELECT COUNT(*) FROM support_tickets
WHERE user_id = $1
  AND status = $2;

-- name: UpdateTicketStatus :one
UPDATE support_tickets
SET status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateTicket :one
UPDATE support_tickets
SET status = $2,
    priority = $3,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetTicketForUpdate :one
SELECT * FROM support_tickets
WHERE id = $1
FOR UPDATE;

-- name: DeleteTicket :exec
DELETE FROM support_tickets
WHERE id = $1;

-- Ticket Messages Queries

-- name: CreateTicketMessage :one
INSERT INTO ticket_messages (
    ticket_id,
    user_id,
    message,
    is_staff
) VALUES (
    $1, $2, $3, $4
)
RETURNING *;

-- name: GetTicketMessageByID :one
SELECT * FROM ticket_messages
WHERE id = $1
LIMIT 1;

-- name: ListTicketMessages :many
SELECT * FROM ticket_messages
WHERE ticket_id = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;

-- name: CountTicketMessages :one
SELECT COUNT(*) FROM ticket_messages
WHERE ticket_id = $1;

-- name: GetLatestTicketMessage :one
SELECT * FROM ticket_messages
WHERE ticket_id = $1
ORDER BY created_at DESC
LIMIT 1;

-- name: DeleteTicketMessage :exec
DELETE FROM ticket_messages
WHERE id = $1;

-- Admin/Staff Queries

-- name: ListAllTickets :many
SELECT * FROM support_tickets
ORDER BY 
    CASE 
        WHEN priority = 'urgent' THEN 1
        WHEN priority = 'high' THEN 2
        WHEN priority = 'medium' THEN 3
        WHEN priority = 'low' THEN 4
    END,
    created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountAllTickets :one
SELECT COUNT(*) FROM support_tickets;

-- name: ListTicketsByStatus :many
SELECT * FROM support_tickets
WHERE status = $1
ORDER BY 
    CASE 
        WHEN priority = 'urgent' THEN 1
        WHEN priority = 'high' THEN 2
        WHEN priority = 'medium' THEN 3
        WHEN priority = 'low' THEN 4
    END,
    created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountTicketsByStatus :one
SELECT COUNT(*) FROM support_tickets
WHERE status = $1;

-- name: GetTicketStats :one
SELECT 
    COUNT(*) FILTER (WHERE status = 'open') as open_count,
    COUNT(*) FILTER (WHERE status = 'in_progress') as in_progress_count,
    COUNT(*) FILTER (WHERE status = 'waiting') as waiting_count,
    COUNT(*) FILTER (WHERE status = 'resolved') as resolved_count,
    COUNT(*) FILTER (WHERE status = 'closed') as closed_count,
    COUNT(*) FILTER (WHERE priority = 'urgent') as urgent_count,
    COUNT(*) FILTER (WHERE priority = 'high') as high_count,
    COUNT(*) as total_count
FROM support_tickets;
