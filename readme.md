# Go Boilerplate

[![CI](https://github.com/chud-lori/go-boilerplate/actions/workflows/ci.yaml/badge.svg)](https://github.com/chud-lori/go-boilerplate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chud-lori/go-boilerplate)](https://goreportcard.com/report/github.com/chud-lori/go-boilerplate)
![Go Version](https://img.shields.io/badge/go-1.23+-blue)
<!-- ![License](https://img.shields.io/github/license/chud-lori/go-boilerplate) -->

A modern, production-ready Go boilerplate for building scalable web APIs and microservices. This project follows Clean Architecture principles and includes best practices for modularity, testing, observability, and maintainability.

---

## ✨ Features

- **Clean Architecture**: Separation of concerns with `domain`, `adapters`, and `infrastructure` layers.
- **REST API**: User CRUD endpoints with DTOs, controllers, and routing.
- **gRPC Support**: Example gRPC service (`Mail`) with protobuf definitions and testable client/server.
- **Caching (Redis)**: In-memory caching with Redis through a `Cache` interface for performance optimization.
- **PostgreSQL Integration**: Repository pattern with transaction support, migrations, and test containers for DB testing.
- **Database Migrations**: Built-in support with [golang-migrate](https://github.com/golang-migrate/migrate).
- **Middleware**: Logging, API key authentication, and request context propagation.
- **Logging**: Structured logging with Logrus, configurable log levels.
- **Error Handling**: Centralized error types and helpers.
- **Testing**: Extensive unit and integration tests with mocks and test containers.
- **Dockerized**: Dockerfile and `docker-compose.yml` for local development and deployment.
- **Observability**: Loki/Promtail/Grafana stack for log aggregation and visualization.
- **Swagger Docs**: Built-in support for Swagger API documentation with [Swag CLI](https://github.com/swaggo/swag).

---

## 🗂️ Project Structure

```
├── adapters/                  # Interfaces connecting the app to the outside world
│   ├── controllers/           # HTTP handlers implementing input ports
│   ├── middleware/            # HTTP middleware (e.g., logging, API key auth)
│   ├── repositories/          # DB implementation of domain repositories
│   └── web/                   # Web utilities
│       ├── dto/               # Request/response DTO structs
│       ├── helper/            # Helper functions for web layer
│       └── routes.go          # HTTP route registration
│
├── cmd/
│   └── api/                   # Application entry point (`main.go`)
│
├── config/                    # Application configuration loading and parsing
│
├── docs/                      # Swagger documentation (auto-generated)
│
├── domain/                   # Core business logic layer (Clean Architecture)
│   ├── entities/              # Domain entities (e.g., User)
│   ├── ports/                 # Interfaces for services, repositories, cache, etc.
│   └── services/              # Application use case implementations
│
├── grpc_service/              # gRPC server/client implementation
│
├── infrastructure/
│   ├── cache/                 # Redis cache implementation
│   └── datastore/             # Database setup and connection logic
│
├── internal/
│   ├── testutils/             # Helpers and setup for tests
│   └── utils/                 # Internal utilities like graceful shutdown
│
├── migrations/                # SQL migration files for `golang-migrate`
│
├── mocks/                     # Mocks for interfaces used in tests
│
├── pkg/                       # Reusable utilities across layers
│   ├── auth/                  # Passcode generator and auth helpers
│   ├── errors/                # Custom error definitions and wrappers
│   └── logger/                # Logrus setup and log configuration
│
├── proto/                     # Generated protobuf files for gRPC
│
├── .env.example               # Example environment variables
├── .env.docker                # Env vars used in Docker deployments
├── docker-compose.yml         # Docker Compose file for multi-service setup
├── Dockerfile                 # Docker build instructions for the API service
├── grafana-datasources.yml    # Grafana configuration for Loki log data source
├── init.sql                   # Optional DB init script for Postgres service
├── mail.proto                 # Protobuf service definition for gRPC
├── Makefile                   # Developer automation (test, build, run, migrate)
├── promtail.yml               # Promtail config for log shipping to Loki
├── go.mod                     # Go module definition
├── go.sum                     # Go module checksums
└── readme.md                  # Project documentation (you’re here!)
```


---

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose

---

### 🧑‍💻 Local Development

1. **Clone the repo**
2. **Configure environment variables**
   ```sh
   cp .env.example .env
   cp keys.env.example keys.env
   ```
3. **Start services**
   ```sh
   make up
   ```

4. **API Endpoints**
   - REST: `POST /api/user`, `PUT /api/user/{userId}`, etc.
   - gRPC: See [`grpc_service/`](grpc_service/) and [`mail.proto`](mail.proto)

5. **Swagger Docs**
   ```sh
   make swagger
   ```
   Access at: `http://localhost:8080/docs/index.html`

---

### 🐳 Running with Docker

1. **Create Docker environment file**
   ```sh
   cp .env.example .env.docker
   ```

2. **Edit `.env.docker`** to set environment variables specifically for your Docker deployment, e.g.:
   ```
   DB_NAME=service_db
   PSQL_USER=postgres
   PSQL_PASSWORD=root
   DATABASE_URL=postgres://postgres:root@service-postgres:5432/service_db?sslmode=disable
   REDIS_URL=redis://service-redis:6379
   ```

3. **Start containers**
   ```sh
   docker-compose up --build
   ```

> 📝 Note: `.env.docker` will be used in `docker-compose.yml` via the `env_file` section. The application container will load this file at runtime by copying it as `.env`.

---

## 🧪 Running Tests

```sh
make test
```

- Unit tests for services and helpers.
- Integration tests using [testcontainers-go](https://github.com/testcontainers/testcontainers-go) for PostgreSQL and Redis.

---

## 🔁 Caching (Redis)

- **Redis** is integrated as a caching layer via the `Cache` interface in `domain/ports/cache.go`.
- **Usage**:
  - Store and retrieve values via `Set`, `Get`, `Delete`.
  - Injected directly into services for cache-first logic (e.g., `GetUser()` → check cache → fallback to DB).
- **Implementation**: `infrastructure/cache/redis_cache.go`
- **Tested via**: [testcontainers-go](https://github.com/testcontainers/testcontainers-go)

---

## 📊 Logging & Observability

- **Structured Logs**: Using Logrus, output to console or file.
- **Grafana + Loki + Promtail**:
  - View logs via `http://localhost:3000` (Grafana)
  - Configured through `docker-compose.yml` and `promtail.yml`

---

## 🧱 Extending

- Add new entities and repositories in `domain/` and `adapters/repositories/`.
- Add gRPC services in `proto/` and regenerate stubs.
- Add middleware in `adapters/middleware/`.

---

## 🪪 License

[Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)
