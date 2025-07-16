# Build stage
FROM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

ENV GO111MODULE=on

COPY . .
RUN go build -trimpath -o /main ./cmd/server


FROM golang:1.24.4-alpine AS runner
WORKDIR /app
COPY --from=builder /main /app/main
CMD [ "/app/main" ]
