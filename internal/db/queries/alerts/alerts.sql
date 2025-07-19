-- name: CreateAlert :one
INSERT INTO alerts (    id, tenant_id, rule_id, event_id, triggered_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, tenant_id, rule_id, event_id, triggered_at;

-- name: InsertAlertProcessingQueue :one
INSERT INTO alert_processing_queue (alert_id)
VALUES ($1)
RETURNING id, alert_id, status, locked_at, locked_by, processed_at, error;

-- name: FetchNextAlertForProcessing :one
UPDATE alert_processing_queue
SET status = 'processing',
    locked_at = now(),
    locked_by = $1
WHERE id = (
    SELECT id
    FROM alert_processing_queue
    WHERE status = 'pending'
    ORDER BY id
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
RETURNING alert_id;

-- name: MarkAlertProcessingFailed :exec
UPDATE alert_processing_queue
SET status = 'failed', error = $2, processed_at = now()
WHERE alert_id = $1;

-- name: MarkAlertProcessingDelivered :exec
UPDATE alert_processing_queue
SET status = 'delivered', processed_at = now()
WHERE alert_id = $1;
