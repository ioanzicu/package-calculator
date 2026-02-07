# Build stage
FROM golang:1.25.4-alpine AS builder

WORKDIR /app

# Install goose
RUN GOBIN=/app/bin go install github.com/pressly/goose/v3/cmd/goose@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/gopher ./cmd/main.go

# Runtime stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/gopher .
COPY --from=builder /app/bin/goose .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/migrations ./migrations
# Copy .env file if it's meant to be used inside, but generally better to pass env vars
# COPY .env . 

# Expose port (document only)
EXPOSE 8080

# Run the binary
CMD ["./gopher"]
