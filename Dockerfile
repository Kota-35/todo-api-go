# Build stage
FROM golang:1.23-alpine AS builder

# Install necessary tools for Prisma
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Set Go module mode explicitly and disable GOPATH
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOPATH=/dev/null
ENV GOFLAGS=-mod=mod
ENV GOCACHE=/tmp/go-cache
ENV GOTMPDIR=/tmp

# Copy mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Debug and verify module mode
RUN pwd && ls -la && cat go.mod
RUN go env GOMOD GO111MODULE GOPATH GOROOT

# Ensure we're in module mode and clean dependencies
RUN go clean -modcache || true
RUN go mod tidy -v

# Generate Prisma client
RUN go run github.com/steebchen/prisma-client-go generate --schema=./prisma/schema.prisma

# Verify Prisma client was generated
RUN ls -la prisma/db/

# Debug: Check if files exist
RUN ls -la internal/domain/valueobject/
RUN ls -la internal/application/dto/auth/

# Debug: Check Go module resolution
RUN unset GOPATH && GO111MODULE=on go list -m all
RUN unset GOPATH && GO111MODULE=on go list ./internal/domain/valueobject

# Build with explicit module mode (unsetting GOPATH completely)
RUN cd /app && unset GOPATH && CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -mod=mod -v -o /main ./cmd/server


FROM golang:1.23-alpine AS runner
WORKDIR /app
COPY --from=builder /main /app/main
CMD [ "/app/main" ]
