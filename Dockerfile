# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application and migration binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/api ./cmd/api && \
    CGO_ENABLED=0 GOOS=linux go build -o /app/bin/migrate ./cmd/migrate

# Final stage
FROM alpine:3.19

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy the binaries from builder
COPY --from=builder /app/bin/api /app/api
COPY --from=builder /app/bin/migrate /app/migrate

# Copy configuration files
COPY config/config.yaml /app/config/config.yaml
COPY config/public.pem /app/config/public.pem

# Copy migrations
COPY migrations /app/migrations

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["/app/api"] 