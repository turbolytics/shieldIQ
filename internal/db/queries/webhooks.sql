-- name: CreateWebhook :one
INSERT INTO webhooks (tenant_id, name, secret, source)
VALUES ($1, $2, $3, $4)
    RETURNING *;