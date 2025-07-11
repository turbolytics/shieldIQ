-- name: CreateRule :one
INSERT INTO rules (
    id, tenant_id, name, description, source, event_type, sql, eval_type, alert_level, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING id, tenant_id, name, description, source, event_type, sql, eval_type, alert_level, created_at;

-- name: ListRules :many
SELECT id, tenant_id, name, description, source, event_type, sql, eval_type, alert_level, created_at
FROM rules
WHERE tenant_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: DeleteRule :exec
DELETE FROM rules WHERE id = $1 AND tenant_id = $2;

-- name: GetRuleByID :one
SELECT id, tenant_id, name, description, source, event_type, sql, eval_type, alert_level, created_at
FROM rules
WHERE id = $1 AND tenant_id = $2;

-- name: GetRulesForEvent :many
SELECT * FROM rules WHERE tenant_id = $1 AND source = $2 AND event_type = $3;
