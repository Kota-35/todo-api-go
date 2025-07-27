# Build stage
FROM golang:1.23-alpine AS builder

# Install necessary tools for Prisma
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Set Go module mode
ENV GO111MODULE=on
ENV CGO_ENABLED=0

# Copy mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Debug: Check environment and generate Prisma client
RUN echo "=== Working directory ===" && \
    pwd && \
    echo "=== Files in current directory ===" && \
    ls -la && \
    echo "=== go.mod content ===" && \
    cat go.mod && \
    echo "=== Module list containing prisma ===" && \
    go list -m all | grep -i prisma || echo "No prisma modules found" && \
    echo "=== Go environment ===" && \
    go env GOMOD GOPATH GOROOT && \
    echo "=== Attempting to run prisma-client-go ===" && \
    go run github.com/steebchen/prisma-client-go generate --schema=./prisma/schema.prisma

# Build application
RUN go build -o /main ./cmd/server


FROM golang:1.23-alpine AS runner
WORKDIR /app
COPY --from=builder /main /app/main
CMD [ "/app/main" ]
