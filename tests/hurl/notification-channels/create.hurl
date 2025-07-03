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
[Asserts]
jsonpath "$.id" exists
jsonpath "$.type" == "slack"
jsonpath "$.name" exists
jsonpath "$.created_at" exists

