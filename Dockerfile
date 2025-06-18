# Builder stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Install necessary tools for build & migration
RUN apk add --no-cache bash make curl ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Overwrite .env with docker env before running anything
COPY .env.docker .env

# Install golang-migrate
# Removed, the migration should not be in build image, move to docker-compose.yml instead
# RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.1/migrate.linux-amd64.tar.gz \
#     | tar -xz && mv migrate /usr/local/bin/migrate && chmod +x /usr/local/bin/migrate

# Run database migrations
# RUN make migration-reset
# RUN make migration-up

# Run tests, REMOVED to avoid error testcontainners
# pun test on CI pipeline instead
# RUN make test

# Build final binary
RUN make build

# Production stage
FROM alpine:latest

WORKDIR /app

# Just copy the built binary and env file
ARG APP_NAME=bin/service-app
COPY --from=builder /app/${APP_NAME} ./main
COPY --from=builder /app/.env .

# Run the app directly
CMD ["./main"]
