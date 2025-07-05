-- name: CreateRuleDestination :one
INSERT INTO rule_destinations (rule_id, channel_id)
VALUES ($1, $2)
RETURNING rule_id, channel_id;

-- name: ListRuleDestinationChannelIDs :many
SELECT channel_id FROM rule_destinations WHERE rule_id = $1;

-- name: ListNotificationChannelsForRule :many
SELECT nc.id, nc.tenant_id, nc.name, nc.type, nc.config, nc.created_at
FROM notification_channels nc
JOIN rule_destinations rd ON nc.id = rd.channel_id
WHERE rd.rule_id = $1;

-- name: DeleteRuleDestination :exec
DELETE FROM rule_destinations WHERE rule_id = $1 AND channel_id = $2;
