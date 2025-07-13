-- name: CreateWebhook :one
INSERT INTO webhooks (id, tenant_id, name, secret, source, events, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
    RETURNING *;

-- name: GetWebhook :one
SELECT * FROM webhooks WHERE id = $1 AND tenant_id = $2;
