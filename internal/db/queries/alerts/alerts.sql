-- name: CreateAlert :one
INSERT INTO alerts (    id, tenant_id, rule_id, event_id, triggered_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, tenant_id, rule_id, event_id, triggered_at;
