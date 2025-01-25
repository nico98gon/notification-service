FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/air-verse/air@latest

COPY . .

RUN go env
RUN go build -v -o notification-service ./cmd/main.go

# FROM debian:buster-slim
FROM golang:1.23-alpine

WORKDIR /app

COPY --from=builder /app/notification-service .
COPY --from=builder /go/bin/air /usr/local/bin/air

RUN chmod +x ./notification-service

EXPOSE 8082

CMD ["air", "-c", ".air.toml"]
