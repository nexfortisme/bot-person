# Use the official Go image as the base image
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY . .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

RUN mkdir -p /app/data /app/img && \
    chown -R appuser:appgroup /app

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/example.env .

RUN chown appuser:appgroup main

USER appuser

# Mount host paths at runtime (see docker-compose.yml):
#   ./.env              -> /app/.env
#   ./data/db.sqlite    -> /app/data/db.sqlite  (or mount ./data -> /app/data)
#   ./img               -> /app/img
VOLUME ["/app/.env", "/app/data", "/app/img"]

ENV DB_PATH=/app/data/db.sqlite
ENV DEV_MODE=false

CMD ["./main"]
