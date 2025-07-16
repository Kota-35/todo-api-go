# Build stage
FROM golang:1.24.4-alpine AS builder
WORKDIR /app
COPY . /app/
RUN go mod download && go build -o main /app/cmd/server/main.go

FROM golang:1.24.4-alpine AS runner
WORKDIR /app
COPY --from=builder /app/main .
USER 1001
CMD [ "/app/main" ]
