# Build stage - use latest Go
FROM golang:alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o app .

# Final stage - minimal image
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS/SSL connections
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/app .

# Expose port
EXPOSE 8080

# Run
CMD ["./app"]
