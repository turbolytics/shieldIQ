
.PHONY: start-backing-services
start-backing-services:
	@echo "Starting backing services..."
	@docker-compose -f dev/docker-compose.yml up -d

build:
	go build -o bin/sqlsec ./cmd/sqlsec