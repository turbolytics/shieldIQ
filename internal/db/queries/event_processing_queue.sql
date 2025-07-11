-- name: InsertEventProcessingQueue :one
INSERT INTO event_processing_queue (event_id)
VALUES ($1)
RETURNING id, event_id, status, locked_at, locked_by, processed_at, error;

-- name: FetchNextEventForProcessing :one
UPDATE event_processing_queue
SET status = 'processing',
    locked_at = now(),
    locked_by = $1
WHERE id = (
    SELECT id
    FROM event_processing_queue
    WHERE status = 'pending'
    ORDER BY id
    FOR UPDATE SKIP LOCKED
    LIMIT 1
)
RETURNING event_id;

-- name: MarkEventProcessingFailed :exec
UPDATE event_processing_queue
SET status = 'failed', error = $2, processed_at = now()
WHERE event_id = $1;

-- name: MarkEventProcessingDone :exec
UPDATE event_processing_queue
SET status = 'done', processed_at = now()
WHERE event_id = $1;

