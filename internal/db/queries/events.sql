-- name: InsertEvent :one
INSERT INTO events (
    tenant_id,
    webhook_id,
    source,
    event_type,
    action,
    raw_payload,
    dedup_hash
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
)
RETURNING id, tenant_id, webhook_id, source, event_type, action, raw_payload, dedup_hash, received_at;

