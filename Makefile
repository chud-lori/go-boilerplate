ifneq (,$(wildcard ./.env))
    include .env
    export
endif

APP_NAME := service-app
BUILD_DIR := bin
DOCKER_IMAGE := service-app
DOCKER_COMPOSE := docker-compose.yml

all: migration-reset migration-up test build run

migration-create:
	migrate create -ext sql -dir migrations -seq $(name)

migration-up:
	migrate -path migrations -database "${DATABASE_URL}" up

migration-down:
	migrate -path migrations -database "${DATABASE_URL}" down 1

# Rollback all migrations without prompt (for local reset only!)
migration-reset:
	yes | migrate -path migrations -database "${DATABASE_URL}" down

test:
	@echo "Running tests..."
	@go test ./... -v -cover
	@echo "Tests passed."

deps:
	@echo "Downloading dependencies..."
	@go mod download

build:
	@echo "Building binary..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api/main.go
	@echo "Build successful. Binary: $(BUILD_DIR)/$(APP_NAME)"

run: build
	@echo "Running application..."
	@./$(BUILD_DIR)/$(APP_NAME)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g cmd/api/main.go

# ====== GRPC ======
run-grpc:
	@go run ./cmd/grpcserver/main.go

# ====== UPLOAD CONSUMER ======
run-upload-consumer:
	@go run ./cmd/upload_consumer/main.go

# ====== DOCKER COMMANDS ======

docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE) .
	@echo "Docker image built: $(DOCKER_IMAGE)"

docker-test:
	@echo "Running tests inside Docker container..."
	@docker run --rm -v $$PWD:/app -w /app golang:1.22.2-alpine \
		sh -c "apk add --no-cache git && go mod download && go test ./... -v -cover"
	@echo "Docker test run completed."

up:
	@echo "Starting Docker Compose stack..."
	@docker compose -f $(DOCKER_COMPOSE) up --build

down:
	@echo "Stopping Docker Compose stack..."
	@docker compose -f $(DOCKER_COMPOSE) down -v

rebuild:
	@echo "Rebuilding Docker stack..."
	@make down
	@make docker-build
	@make up

# ====== HELP ======

help:
	@echo "Usage:"
	@echo "  make            Run tests, build, and start the app"
	@echo "  make test       Run all unit tests"
	@echo "  make deps       Download Go dependencies"
	@echo "  make build      Build the Go binary into $(BUILD_DIR)/$(APP_NAME)"
	@echo "  make run        Run the application locally"
	@echo "  make clean      Remove built binaries"
	@echo "  make swagger    Generate Swagger documentation"
	@echo "  make run-grpc   Run gRPC Server"
	@echo "  make run-upload-consumer   Run Upload Consumer (async worker)"
	@echo ""
	@echo "Migration targets:"
	@echo "  make migration-create name=your_migration_name   Create a new migration file"
	@echo "  make migration-up                                Apply all up migrations"
	@echo "  make migration-down                              Revert the last migration"
	@echo ""
	@echo "Docker targets:"
	@echo "  make docker-build   Build Docker image (also runs tests inside Dockerfile)"
	@echo "  make docker-test    Run tests in clean Go container"
	@echo "  make up             Start full Docker Compose stack"
	@echo "  make down           Stop Docker Compose stack"
	@echo "  make rebuild        Full Docker rebuild and restart"

.PHONY: all test deps build swagger run clean help run-grpc \
        docker-build docker-test up down rebuild migration-create \
		migration-up migration-down run-upload-consumer
