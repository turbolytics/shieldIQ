# Monitor GitHub


- Start sqlsec
- Create a Webhook to receive events from GitHub








# Commands

### Create a Webhook
```
SQLSEC_API_BASE_URL=http://localhost:8888 go run cmd/sqlsec/main.go api webhooks create --name=github-dco-1 --source=github

+------------+------------------------------------------------------------------+
| Attribute  | Value                                                            |
+------------+------------------------------------------------------------------+
| source     | github                                                           |
| created_at | 2025-07-25T10:47:17.212476Z                                      |
| events     | <nil>                                                            |
| id         | a15f65dc-35d6-41fb-a7f4-583153a08af4                             |
| tenant_id  | 00000000-0000-0000-0000-000000000000                             |
| name       | github-dco-1                                                     |
| secret     | 1ddd56351c3bae833a57f107abe14516bdb1eafb43e33e55632d0bf817fedb25 |
+------------+------------------------------------------------------------------+
```
