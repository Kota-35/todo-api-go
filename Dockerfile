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

# Debug and generate Prisma client
RUN echo "=== Module list before tidy ===" && \
    go list -m all | grep prisma || echo "No prisma modules found" && \
    echo "=== Running go mod tidy ===" && \
    go mod tidy && \
    echo "=== Module list after tidy ===" && \
    go list -m all | grep prisma || echo "No prisma modules found" 
    
    
RUN echo "=== Trying to run prisma-client-go ===" && \
    go run github.com/steebchen/prisma-client-go generate --schema=./prisma/schema.prisma

# Build application
RUN go build -o /main ./cmd/server


FROM golang:1.23-alpine AS runner
WORKDIR /app
COPY --from=builder /main /app/main
CMD [ "/app/main" ]
