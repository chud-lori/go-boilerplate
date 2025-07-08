# Go Boilerplate

[![CI](https://github.com/chud-lori/go-boilerplate/actions/workflows/ci.yaml/badge.svg)](https://github.com/chud-lori/go-boilerplate/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/chud-lori/go-boilerplate)](https://goreportcard.com/report/github.com/chud-lori/go-boilerplate)
![Go Version](https://img.shields.io/badge/go-1.23+-blue)
![License](https://img.shields.io/github/license/chud-lori/go-boilerplate)

A modern, production-ready Go boilerplate for building scalable web APIs and microservices. This project follows Clean Architecture principles and includes best practices for modularity, testing, observability, and maintainability.

---

## ✨ Features

- **Clean Architecture**: Separation of concerns with domain, adapters, and infrastructure layers.
- **REST API**: User CRUD endpoints with DTOs, controllers, and routing.
- **gRPC Support**: Example gRPC service (Mail) with protobuf definitions and testable client/server.
- **Caching (Redis)**: In-memory caching with Redis through a Cache interface for performance optimization.
- **PostgreSQL Integration**: Repository pattern with transaction support, migrations, and test containers for DB testing.
- **Database Migrations**: Built-in support with [golang-migrate](https://github.com/golang-migrate/migrate).
- **Middleware**: Logging, API key authentication, and request context propagation.
- **Logging**: Structured logging with Logrus, configurable log levels.
- **Error Handling**: Centralized error types and helpers.
- **Testing**: Extensive unit and integration tests with mocks and test containers.
- **Dockerized**: Dockerfile and docker-compose.yml for local development and deployment.
- **Observability**: Loki/Promtail/Grafana stack for log aggregation and visualization.
- **Swagger Docs**: Built-in support for Swagger API documentation with [Swag CLI](https://github.com/swaggo/swag).

---

## 🗂️ Project Structure

```
├── adapters/                  # Interfaces connecting the app to the outside world
│   ├── controllers/           # HTTP handlers implementing input ports
│   ├── middleware/            # HTTP middleware (e.g., logging, API key auth)
│   ├── repositories/          # DB implementation of domain repositories
│   └── web/                   # Web utilities including DTOs and helpers
│       ├── dto/               # Request/response DTO structs
│       ├── helper/            # Helper functions for the web layer
│       └── routes.go          # HTTP route registration
│
├── cmd/                      # Application entry points
│   ├── api/                   # Main REST API entry point
│   └── grpcserver/            # Main gRPC mail server entry point
│
├── config/                   # Application configuration loading and parsing
│
├── docs/                     # Swagger documentation (auto-generated)
│
├── domain/                   # Core business logic layer (Clean Architecture)
│   ├── entities/              # Domain entities (e.g., User, Post)
│   ├── ports/                 # Interfaces for controllers, services, repos, etc.
│   └── services/              # Application use case implementations
│
├── infrastructure/           # External infrastructure implementations
│   ├── cache/                 # Redis cache implementation
│   ├── datastore/             # PostgreSQL DB setup and connection logic
│   └── grpc_clients/          # gRPC clients used by the application
│
├── internal/                 # Internal packages
│   ├── testutils/             # Helpers and setup for tests
│   └── utils/                 # Internal utilities (e.g., graceful shutdown)
│
├── migrations/               # SQL migration files for golang-migrate
│
├── mocks/                    # Mocks for interfaces used in unit tests
│
├── pkg/                      # Reusable utilities across layers
│   ├── auth/                  # Encryption, JWT, and passcode helpers
│   ├── errors/                # Custom error definitions and validation logic
│   └── logger/                # Logrus setup and logger abstraction
│
├── proto/                    # Generated protobuf files for gRPC
│
├── Dockerfile                # Docker build instructions for API service
├── docker-compose.yml        # Docker Compose for service orchestration
├── grafana-datasources.yml   # Grafana configuration for Loki
├── init.sql                  # Optional DB init script for Postgres service
├── mail.proto                # Protobuf definition for gRPC Mail service
├── Makefile                  # Developer automation (build, run, test)
├── promtail.yml              # Promtail config for log shipping to Loki
├── .env.example              # Example environment variables file
├── go.mod                    # Go module definition
├── go.sum                    # Go module checksums
├── LICENSE                   # License information
└── readme.md                 # Project documentation (you’re here!)
```

---

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- **k6**: For running load tests. [Install k6](https://k6.io/docs/getting-started/installation/)

---

### 🧑‍💻 Local Development

1. **Clone the repo**

2. **Configure environment variables**

```sh
cp .env.example .env
```

3. **Start services**

```sh
make up
```

4. **API Endpoints**
   - REST: `POST /api/user`, `PUT /api/user/{userId}`, etc.
   - gRPC: See [`cmd/grpcserver/`](cmd/grpcserver/) and [`mail.proto`](mail.proto)

5. **Swagger Docs**

```sh
make swagger
```

Access at: [http://localhost:8080/docs/index.html](http://localhost:8080/docs/index.html)

> 💡 Tip: You can use this Node.js command to generate a `JWT_SECRET`:

```sh
node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"
```

---

### ⚡️ Running the gRPC Mail Server

> **Note:** The Auth service now sends mail using gRPC and expects a mail gRPC server running at `localhost:50051`.
> You must start the gRPC mail server for the boilerplate to function correctly.

#### 1. Start the gRPC Mail Server

Run the gRPC Server using this command:

```sh
go run cmd/grpcserver/main.go
```

Or, using Makefile:

```sh
make run-grpc
```

#### 2. Confirm the server is running

The gRPC server should listen on `localhost:50051`.
The Auth service will connect to this address to send mail.

---


### 🐳 Running with Docker

1. **Create Docker environment file**

```sh
cp .env.example .env.docker
```

2. **Edit `.env.docker`** to set environment variables for Docker deployment:

```env
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

> 📝 `.env.docker` will be used by `docker-compose.yml` via the `env_file` section.
> The app will treat it as `.env` at runtime.

---

## 🧪 Running Tests

```sh
make test
```

- Unit tests for services and helpers
- Integration tests using [testcontainers-go](https://github.com/testcontainers/testcontainers-go) for PostgreSQL and Redis

---

## 📈 Performance Testing (k6)

This boilerplate includes a k6 script for load testing the REST API endpoints. This helps in understanding the system's performance and stability under various loads.

1.  **Ensure k6 is installed** (see [Prerequisites](#prerequisites)).
2.  **Make sure the Go API service is running** (e.g., via `make up` or `docker-compose up --build`).
3.  **Run the k6 script**:
    ```sh
    k6 run loadtest.js
    ```
    This script simulates a user flow including:
    * Creating a new post (`POST /api/post`)
    * Fetching posts twice (`GET /api/post`)

    Adjust the `vus` (Virtual Users) and `duration` in `loadtest.js` to simulate different load scenarios.

---

## 🔁 Caching (Redis)

- **Redis** is integrated via the `Cache` interface (`domain/ports/cache.go`)
- **Usage**:
    - Cache values with `Set`, retrieve with `Get`, delete with `Delete`
    - Used in services for cache-first logic (e.g., `GetAll()` → check cache → fallback to DB)
- **Implementation**: `infrastructure/cache/redis_cache.go`
- **Tested via**: [testcontainers-go](https://github.com/testcontainers/testcontainers-go)

---

## 📊 Logging & Observability

- **Structured Logging**: Logrus used for consistent, leveled logs
- **Grafana + Loki + Promtail** stack:
    - Access Grafana at: [http://localhost:3000](http://localhost:3000)
    - Logs shipped by Promtail and stored by Loki
    - Configured in `promtail.yml` and `grafana-datasources.yml`

---

## 🧱 Extending

- Add new entities in `domain/entities/` and update related ports/services
- Add repositories in `adapters/repositories/`
- Add gRPC services to `proto/`, regenerate stubs using `protoc`
- Add middleware in `adapters/middleware/`

---

## 🪪 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
