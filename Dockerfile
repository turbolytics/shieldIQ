# syntax=docker/dockerfile:1

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o /app/bin/sqlsec ./cmd/sqlsec

# Runtime stage
FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/bin/sqlsec /usr/local/bin/sqlsec
EXPOSE 8888
ENTRYPOINT ["/usr/local/bin/sqlsec"]

