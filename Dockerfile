# ===============================================================
# Стадия 1: Сборка приложения
# ===============================================================
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/server ./cmd/api

# ===============================================================
# Стадия 2: Создание минимального исполняемого образа
# ===============================================================
FROM alpine:latest

RUN apk --no-cache add ca-certificates tini

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server .
COPY config ./config
COPY migrations ./migrations

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 8080

ENTRYPOINT ["/sbin/tini", "--", "./server"]

HEALTHCHECK --interval=15s --timeout=5s --start-period=10s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/healthz || exit 1
