-- name: CreateWebhook :one
INSERT INTO webhooks (id, tenant_id, name, secret, source, created_at)
VALUES ($1, $2, $3, $4, $5, $6)
    RETURNING *;