# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Set Go module mode explicitly
ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org,direct

# Copy mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy source code
COPY . .

# Clean and build
RUN go mod tidy
RUN go build -o /main ./cmd/server


FROM golang:1.23-alpine AS runner
WORKDIR /app
COPY --from=builder /main /app/main
CMD [ "/app/main" ]
