# ----------------------------------------------------------------------
# Stage 1: Builder - ใช้สำหรับคอมไพล์โค้ด Go เท่านั้น
# ----------------------------------------------------------------------
FROM golang:1.25-alpine AS builder

ENV CGO_ENABLED=1
ENV GOOS=linux

RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o /bot-runner main.go database.go core_logic.go api_client.go config.go

# ----------------------------------------------------------------------
# Stage 2: Final Image - ใช้ Alpine เพื่อให้ Image เล็กที่สุด
# ----------------------------------------------------------------------
FROM alpine:latest

RUN apk --no-cache add sqlite-libs

WORKDIR /app

COPY --from=builder /bot-runner /app/bitkub-rebalance-bot

COPY web /app/web

CMD ["/app/bitkub-rebalance-bot"]