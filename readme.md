# Go Boilerplate

A modern, production-ready Go boilerplate for building scalable web APIs and microservices. This project includes best practices for clean architecture, modularity, testing, and observability.

---

## Features

- **Clean Architecture**: Separation of concerns with `domain`, `adapters`, and `infrastructure` layers.
- **REST API**: User CRUD endpoints with DTOs, controllers, and routing.
- **gRPC Support**: Example gRPC service (`Mail`) with protobuf definitions and testable client/server.
- **PostgreSQL Integration**: Repository pattern with transaction support, migrations, and test containers for DB testing.
- **Middleware**: Logging, API key authentication, and request context propagation.
- **Logging**: Structured logging with Logrus, configurable log levels.
- **Error Handling**: Centralized error types and helpers.
- **Testing**: Extensive unit and integration tests with mocks and test containers.
- **Dockerized**: Dockerfile and `docker-compose.yml` for local development and deployment.
- **Observability**: Loki/Promtail/Grafana stack for log aggregation and visualization.
- **Swagger Docs**: Built-in support for Swagger API documentation.

---

## Project Structure

```
.
├── adapters/           # Controllers, middleware, repositories, web routes & DTOs
├── cmd/                # Entrypoint (main.go)
├── config/             # Configuration files
├── domain/             # Entities, ports (interfaces), and services
├── grpc_service/       # gRPC client/server logic
├── infrastructure/     # Database implementation
├── internal/           # Internal utilities and test helpers
├── mocks/              # Auto-generated and hand-written mocks for testing
├── pkg/                # Shared packages (auth, logger, errors)
├── proto/              # Protobuf generated files
├── .github/            # CI workflows
├── Dockerfile
├── docker-compose.yml
├── db.sql              # DB schema
├── mail.proto          # Protobuf definition
├── Makefile
├── .env, keys.env      # Environment variables
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
   Then access via `/swagger/index.html` (if enabled).

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