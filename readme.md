# Go Boilerplate

A modern, production-ready Go boilerplate for building scalable web APIs and microservices. This project includes best practices for clean architecture, modularity, testing, and observability.

---

## Features

- **Clean Architecture**: Separation of concerns with `domain`, `adapters`, and `infrastructure` layers.
- **REST API**: User CRUD endpoints with DTOs, controllers, and routing.
- **gRPC Support**: Example gRPC service (`Mail`) with protobuf definitions and testable client/server.
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

## Project Structure

```
.
├── adapters                  # Application layer interfaces to external world
│   ├── controllers           # HTTP request handlers
│   ├── middleware            # HTTP middleware (e.g., auth, logging)
│   ├── repositories          # Database access interfaces
│   └── web
│       ├── dto               # Request/response data transfer objects
│       └── helper            # Web helpers (e.g., response formatting)
├── bin                       # Compiled binaries
├── cmd
│   └── api                   # Application entrypoint (main.go)
├── config                    # Application configuration management
├── docs                      # Generated documentation (Swagger, etc.)
├── domain                    # Business logic layer
│   ├── entities              # Core business entities
│   ├── ports                 # Interfaces for input/output layers
│   └── services              # Business use cases
├── grpc_service              # gRPC server/client logic and implementation
├── infrastructure
│   └── datastore             # Concrete implementation of DB or other infrastructure
├── internal
│   ├── testutils             # Shared testing utilities
│   └── utils                 # Internal helper functions
├── migrations                # SQL migration files (used by golang-migrate)
├── mocks                     # Auto-generated and custom mocks for testing
├── pkg
│   ├── auth                  # Authentication-related utilities
│   ├── errors                # Custom error types and wrappers
│   └── logger                # Logging setup and utilities
├── proto                     # Protobuf definitions and generated files
├── .github/            # CI workflows
├── Dockerfile
├── docker-compose.yml
├── db.sql              # DB schema
├── mail.proto          # Protobuf definition
├── Makefile
├── .env.example      # Environment variables
└── readme.md
```

---

## Getting Started

### Prerequisites

- Go 1.23+ (see `go.mod`)
- Docker & Docker Compose

### Running Locally

1. **Clone the repo**
2. **Configure environment variables**
   Copy `.env` and `keys.env` and adjust as needed.
3. **Start services**
   ```sh
   make up
   ```
   This will build the app, start PostgreSQL, Grafana, Loki, and Promtail.

4. **API Endpoints**
   - REST: `POST /api/user`, `PUT /api/user/{userId}`, etc.
   - gRPC: See [`grpc_service/`](grpc_service/) and [`mail.proto`](mail.proto)

5. **Swagger Docs**
   Generate docs with:
   ```sh
   make swagger
   ```
   Then access via `/docs/index.html` (if enabled).

6. **Run Migrations**
   Migrations are managed using `golang-migrate` via the Makefile:

   #### Install `migrate` CLI:
   ```sh
   brew install golang-migrate  # macOS
   ```
   Or download from: https://github.com/golang-migrate/migrate/releases

   #### Create a new migration
   ```sh
   make migration-create name=create_users_table
   ```

   #### Apply migrations
   ```sh
   make migration-up
   ```

   #### Rollback migrations
   ```sh
   make migration-down
   ```

**Regenerate Swagger Docs**
   If you update route annotations:

   #### Install `swag` CLI:
   ```sh
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

   #### Then run:
   ```sh
   make swagger
   ```

### Running Tests

```sh
make test
```

- Unit and integration tests use mocks and test containers for DB isolation.

---

## Logging & Observability

- **Logs**: Structured JSON logs via Logrus.
- **Loki/Promtail/Grafana**:
  - Logs are scraped from containers and visualized in Grafana (`localhost:3000`).

---

## Extending

- Add new entities and repositories in `domain/` and `adapters/repositories/`.
- Add new gRPC services in `mail.proto` and regenerate with `protoc`.
- Add new middleware in `adapters/middleware/`.

---

## License

[Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)