POST http://localhost:8888/api/notification-channels
Content-Type: application/json

{
  "name": "Test Slack Channel",
  "type": "slack",
  "config": {
    "webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
  }
}

HTTP/1.1 201
GET http://localhost:8888/api/notification-channels
HTTP/1.1 200
[Asserts]
jsonpath "$[0].id" exists
jsonpath "$[0].type" == "slack"
jsonpath "$[0].name" exists
jsonpath "$[0].created_at" exists

