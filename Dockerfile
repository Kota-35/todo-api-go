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

# Add missing dependencies and generate Prisma client using go install
RUN go get github.com/joho/godotenv@v1.5.1 && \
    go get github.com/shopspring/decimal@v1.4.0 && \
    go get github.com/steebchen/prisma-client-go@v0.47.0 && \
    go mod tidy && \
    go install github.com/steebchen/prisma-client-go@v0.47.0 && \
    prisma-client-go generate --schema=./prisma/schema.prisma

# Build application
RUN go build -o /main ./cmd/server


FROM golang:1.23-alpine AS runner
WORKDIR /app
COPY --from=builder /main /app/main
CMD [ "/app/main" ]
