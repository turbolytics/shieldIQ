FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    g++ cmake wget unzip pkg-config git && \
    rm -rf /var/lib/apt/lists/*

# Install toml++ headers
RUN wget https://github.com/marzer/tomlplusplus/archive/refs/tags/v3.4.0.zip && \
    unzip v3.4.0.zip && \
    cp -r tomlplusplus-3.4.0/include/* /usr/local/include/ && \
    rm -rf v3.4.0.zip tomlplusplus-3.4.0

ENV CGO_ENABLED=1
ENV CPATH=/usr/local/include

WORKDIR /app
COPY . .

RUN go build -v -o /app/bin/shieldIQ ./cmd/shieldIQ


# Runtime stage
FROM debian:bullseye-slim

WORKDIR /app
COPY --from=builder /app/bin/shieldIQ /usr/local/bin/shieldIQ

EXPOSE 8888
ENTRYPOINT ["/usr/local/bin/shieldIQ"]
