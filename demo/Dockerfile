# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build basic-setup
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/basic-setup ./basic-setup/cmd

# Build super-worker
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/super-worker ./super-worker/cmd

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /app/bin /app/bin

# Expose ports
# 3010 for basic-setup web server
# 6060 for pprof
EXPOSE 3010 6060

# Default command (can be overridden in docker-compose)
CMD ["/app/bin/basic-setup"]

