FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o job-website-backend ./cmd/api

FROM alpine:latest
WORKDIR /app
RUN adduser -D appuser
COPY --from=builder /app/job-website-backend .
COPY config ./config
COPY migrations ./migrations
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget --spider -q http://localhost:8080/ || exit 1
CMD ["./job-website-backend"] 