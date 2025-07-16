# Build stage
FROM golang:1.24.4-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Force Go modules mode
ENV GO111MODULE=on
ENV GOPATH=""
ENV CGO_ENABLED=0

# Copy go.mod and go.sum first (for better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy entire source code
COPY . .

# Build the application
RUN go build -o main ./cmd/server

# Runtime stage
FROM alpine:latest

# Install ca-certificates
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Set working directory
WORKDIR /home/appuser

# Copy binary from builder stage
COPY --from=builder /app/main .

# Change ownership
RUN chown appuser:appgroup main

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
