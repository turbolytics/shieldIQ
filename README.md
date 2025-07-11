# sqlsec
sqlsec is a lightweight SIEM for modern SaaS


# Quick Start

## Setup a Source

- Start the backing services:
```
make start-backing-services
```

- Create a Webhook to receive events from a source like GitHub:
```
curl -X POST http://localhost:8888/api/webhooks/create \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test github",
    "source": "github",
    "events": ["pull_request"]
  }'
  
  
 {
  "id": "eab1a2f2-4eac-4eff-af41-18ee303dfa59",
  "tenant_id": "00000000-0000-0000-0000-000000000000",
  "name": "Test github",
  "secret": "d297cd5b54e52ff4c604a0915323f8a588813313301c67b4d408d3450d008140",
  "source": "github",
  "created_at": "2025-07-08T11:09:10.269661Z",
  "events": [
    "event1",
    "event2"
  ]
}
```

- Add the Webhook to your GitHub repository
- Send a test Event / Verify the Webhook

## Create a Notification Channel 

```
curl -X POST http://localhost:8888/api/notification-channels \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Slack Channel",
    "type": "slack",
    "config": {
      "webhook_url": "<your-slack-webhook-url>"
    }
  }'
```

## Register a Rule


# Development

```
docker exec -i sqlsec_postgres psql -U sqlsec -d sqlsec < schema.sql
```

```
SQLSEC_DB_DSN=postgres://sqlsec:sqlsec@localhost:5432/sqlsec?sslmode=disable go run cmd/sqlsec/main.go serve -p 8888
```

```
docker exec -it sqlsec_postgres psql -U sqlsec
```