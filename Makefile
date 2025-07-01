
.PHONY: start-backing-services
start-backing-services:
	@echo "Starting backing services..."
	@docker-compose -f dev/docker-compose.yml up -d

.PHONY: stop-backing-services
stop-backing-services:
	@echo "Stopping backing services..."
	@docker-compose -f dev/docker-compose.yml down --remove-orphans

.PHONY: build
build:
	go build -o bin/sqlsec ./cmd/sqlsec

.PHONY: test-hurl
test-hurl:
	hurl --test --glob "tests/hurl/**/*.hurl"

.PHONY: sqlc-generate
sqlc-generate:
	@echo "Generating SQLC code..."
	sqlc generate