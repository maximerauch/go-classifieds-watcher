# Build Stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Download dependencies first (caching layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build a static binary
# CGO_ENABLED=0 ensures the binary is fully linked and doesn't depend on system C libraries
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/watcher/main.go

# Final Runtime Stage
FROM alpine:latest

WORKDIR /app

# Install CA certificates to allow HTTPS requests to the API
RUN apk --no-cache add ca-certificates

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Run the application
CMD ["./main"]