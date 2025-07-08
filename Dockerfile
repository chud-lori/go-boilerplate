# Builder stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Install necessary tools for build & migration
RUN apk add --no-cache bash make curl ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Overwrite .env with docker env before running anything
# Consider managing environment variables via Docker Compose directly for better flexibility.
COPY .env.docker .env

# Build API service binary
# This assumes your main function for the web service is in cmd/api/main.go
RUN go build -o bin/api-service ./cmd/api

# Build gRPC server binary
# This assumes your main function for the gRPC server is in cmd/grpcserver/main.go
RUN go build -o bin/grpc-server ./cmd/grpcserver

# Production stage
FROM alpine:latest

WORKDIR /app

# Just copy the built binaries and env file
COPY --from=builder /app/bin/api-service ./api-service
COPY --from=builder /app/bin/grpc-server ./grpc-server
COPY --from=builder /app/.env .
