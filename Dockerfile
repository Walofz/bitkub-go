# ----------------------------------------------------------------------
# Stage 1: Builder (ใช้ Debian Bookworm)
# ----------------------------------------------------------------------
FROM golang:1.25-bookworm AS builder

# ตั้งค่า Environment
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /app/main main.go database.go core_logic.go api_client.go config.go

# ----------------------------------------------------------------------
# Stage 2: Final Image (ใช้ Debian Bookworm Slim)
# ----------------------------------------------------------------------
FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    sqlite3 \
    libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

ENV TZ=Asia/Bangkok

WORKDIR /app

COPY --from=builder /app/main /app/bitkub-rebalance-bot

COPY web /app/web

RUN chmod +x /app/bitkub-rebalance-bot

CMD ["/app/bitkub-rebalance-bot"]