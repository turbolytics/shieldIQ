# Create a rule
POST http://localhost:8888/api/rules/
Content-Type: application/json

{
  "name": "RuleDest Test Rule",
  "description": "A rule for destination testing",
  "source": "github",
  "event_type": "pull_request",
  "condition": "state = 'open'",
  "evaluation_type": "LIVE_TRIGGER",
  "alert_level": "LOW"
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id" exists
[Captures]
rule_id: jsonpath "$.id"

# Create a notification channel
POST http://localhost:8888/api/notification-channels
Content-Type: application/json

{
  "name": "DestTestChannel",
  "type": "slack",
  "config": {
    "webhook_url": "https://hooks.slack.com/services/XXX/YYY/ZZZ"
  }
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id" exists
[Captures]
channel_id: jsonpath "$.id"

# Attach the notification channel to the rule as a destination
POST http://localhost:8888/api/rules/{{rule_id}}/destinations/
Content-Type: application/json

{
  "rule_id": "{{rule_id}}",
  "channel_id": "{{channel_id}}"
}

HTTP/1.1 201
[Asserts]
jsonpath "$.id" exists
[Captures]
dest_id: jsonpath "$.id"

# List destinations for the rule (should be 1)
GET http://localhost:8888/api/rules/{{rule_id}}/destinations/

HTTP/1.1 200
[Asserts]
jsonpath "$" count == 1
jsonpath "$[0].id" exists

# Delete the destination
DELETE http://localhost:8888/api/rules/{{rule_id}}/destinations/{{dest_id}}

HTTP/1.1 204

# List destinations for the rule (should be 0)
GET http://localhost:8888/api/rules/{{rule_id}}/destinations/

HTTP/1.1 200
[Asserts]
jsonpath "$" count == 0
