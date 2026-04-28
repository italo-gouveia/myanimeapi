# MyAnimeAPI

MyAnimeAPI is a multi-protocol API for managing anime, users, reviews, and authentication. It exposes both a **REST** interface and a **GraphQL** interface, built on a hexagonal architecture so new adapters (gRPC, WebSocket, etc.) can be added without touching business logic. The backend is written in Go, uses Gorilla Mux for REST routing, gqlgen for GraphQL, GORM for database interactions, and PostgreSQL as the database.

---

## Table of Contents

- [MyAnimeAPI](#myanimeapi)
  - [Table of Contents](#table-of-contents)
  - [Features](#features)
  - [Technologies Used](#technologies-used)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Docker Setup](#docker-setup)
  - [Project Structure](#project-structure)
  - [API Endpoints](#api-endpoints)
    - [Health Check](#health-check)
    - [Version](#version)
    - [Authentication](#authentication)
    - [Users](#users)
    - [Anime](#anime)
    - [Genres](#genres)
    - [Tags](#tags)
    - [Reviews](#reviews)
  - [Documentation](#documentation)
    - [Swagger Documentation](#swagger-documentation)
    - [GoDoc Documentation](#godoc-documentation)
  - [Testing](#testing)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [End-to-End Tests](#end-to-end-tests)
    - [Smoke Tests](#smoke-tests)
    - [Test Coverage](#test-coverage)
- [TODO: Need to be updated the images/diagrams](#todo-need-to-be-updated-the-imagesdiagrams)
  - [Diagrams](#diagrams)
    - [**1. Architecture Diagram**](#1-architecture-diagram)
    - [**2. Database Schema (ER Diagram)**](#2-database-schema-er-diagram)
    - [**3. Component Diagram**](#3-component-diagram)
    - [**4. Deployment Diagram**](#4-deployment-diagram)
    - [**5. Client Request Flow**](#5-client-request-flow)
    - [**6. Authentication Flow**](#6-authentication-flow)
    - [**7. Rate Limiting Flow**](#7-rate-limiting-flow)
    - [**8. User Registration Sequence**](#8-user-registration-sequence)
    - [**9. User Authentication Sequence**](#9-user-authentication-sequence)
    - [**10. Anime Creation Sequence**](#10-anime-creation-sequence)
    - [**How to Generate Diagrams**](#how-to-generate-diagrams)
    - [**Notes**](#notes)
  - [Contributing](#contributing)
  - [Changelog](#changelog)
  - [Badges](#badges)
  - [License](#license)
  - [Acknowledgments](#acknowledgments)
  - [Contact](#contact)

---

## Features

- **Anime Management**: Create, read, update, and delete anime entries with support for genres and tags
- **User Management**: Register, authenticate, and manage users with enhanced security
- **Review Management**: Add, update, and delete reviews for anime with media attachments
- **Genre Management**: Manage anime genres with CRUD operations and bulk operations
- **Tag Management**: Manage anime tags with CRUD operations and bulk operations
- **Authentication**: JWT-based authentication with enhanced security features
- **Pagination**: Improved pagination with cursor-based navigation
- **Rate Limiting**: Enhanced rate limiting with different limits for auth and non-auth endpoints
- **Swagger Documentation**: Auto-generated API documentation with detailed examples
- **Favorite Anime**: Add and manage favorite anime entries with bulk operations
- **Error Handling**: Structured error responses with detailed information and context
- **Health Monitoring**: Comprehensive health check endpoint with dependency status
- **Security**: Regular security scanning with OWASP Dependency-Check, Semgrep, and Gitleaks
- **Dependency Management**: Proper dependency injection and service initialization
- **Media Storage**: Support for both local and S3 storage of media files
- **Bulk Operations**: Enhanced bulk operations for genres and tags with validation
- **Request Validation**: Comprehensive request validation and sanitization
- **Logging**: Structured logging with request ID tracking
- **Metrics**: Prometheus metrics for monitoring and alerting

## Technologies Used

- **Go**: Backend programming language (v1.23+)
- **Gorilla Mux**: HTTP router and dispatcher
- **gqlgen**: Code-first GraphQL server (`github.com/99designs/gqlgen`)
- **GORM**: ORM for database interactions
- **PostgreSQL**: Relational database (v15)
- **golang-migrate**: SQL-based schema migrations (embedded in binary)
- **JWT**: JSON Web Tokens for authentication
- **Swagger**: REST API documentation (swaggo)
- **Docker**: Containerization for easy deployment and development
- **Validator**: Input validation
- **AWS SDK**: For S3 storage integration (optional)
- **Bluemonday**: HTML sanitization
- **GolangCI-Lint**: Code quality and style checking
- **SonarCloud**: Code quality and security analysis
- **Semantic Release**: Automated version management

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Docker (optional)
- AWS Account (for S3 storage, optional)

### Installation

Clone the repository:

```bash
git clone https://github.com/italo-gouveia/myanimeapi.git
cd myanimeapi
```

Set up the database:

- Ensure PostgreSQL is running
- Create a database named `myanimeapi`
- Update the `.env` file with your PostgreSQL credentials
- For features like "Forgot Password", ensure you also configure email-related environment variables in the `.env` file (e.g., `EMAIL_FROM`, `EMAIL_PASSWORD`, `SMTP_HOST`, `SMTP_PORT`, `FRONTEND_URL`).

Run the application:

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080/v1`.

### Docker Setup

Build and run the Docker containers:

```bash
docker-compose up --build
```

This will start both the PostgreSQL database and the Go API server.

**Access the API:**  
The API will be available at `http://localhost:8080/v1`.
The Swagger documentation will be available at `http://localhost:8080/swagger/index.html`.

## Project Structure

The project follows a **hexagonal architecture** (Ports & Adapters), where business logic in `api/services/` and `api/repositories/` is completely decoupled from the protocol adapters (`http`, `graphql`). Adding a new adapter (gRPC, WebSocket, etc.) only requires creating a new directory under `api/adapters/`.

```
myanimeapi/
├── api/
│   ├── adapters/
│   │   ├── http/         # REST input adapter (package httphandler)
│   │   │   ├── anime_handlers.go
│   │   │   ├── auth_handler.go
│   │   │   ├── favorite_handler.go
│   │   │   ├── genre_handler.go
│   │   │   ├── review_handlers.go
│   │   │   ├── tag_handler.go
│   │   │   └── user_handler.go
│   │   └── graphql/      # GraphQL input adapter (package graphql)
│   │       ├── schema/schema.graphql
│   │       ├── generated/generated.go
│   │       ├── model/models_gen.go
│   │       ├── resolver.go
│   │       └── schema.resolvers.go
│   ├── models/           # Shared domain models and DTOs
│   ├── services/         # Application layer / use cases (input ports)
│   ├── repositories/     # Data access interfaces + implementations (output ports)
│   ├── middleware/       # HTTP middleware (auth, logging, rate limiting)
│   ├── auth/             # JWT helpers and password hashing
│   ├── utils/            # Shared utilities
│   ├── database/         # Admin user seeding
│   ├── mocks/            # Generated mocks for testing
│   └── routes/           # Composition root — wires all adapters + services
├── cmd/
│   ├── main.go           # Application entry point
│   └── docs/             # Generated Swagger documentation
├── internal/
│   ├── config/           # Configuration loading
│   ├── database/         # golang-migrate runner (embedded SQL migrations)
│   │   └── migrations/   # SQL migration files (*.up.sql / *.down.sql)
│   ├── db/               # GORM abstraction interface
│   ├── errors/           # Custom application error types
│   ├── logger/           # Structured logging (logrus)
│   └── utils/            # Internal utilities
├── tests/
│   ├── e2e/              # End-to-end tests (real DB required)
│   ├── integration/      # Integration tests (mock services or real DB)
│   ├── smoke/            # Smoke tests against running server
│   ├── fixtures/         # Shared test data factories
│   ├── suites/           # Reusable test suite base types
│   └── config/           # Test-specific configuration
│   (unit tests live co-located with source: api/adapters/http/*_test.go etc.)
├── scripts/
│   └── hooks/pre-push    # Git pre-push validation (lint + tests)
├── frontend/             # Next.js frontend (developed separately)
├── docker-compose.yml    # Production stack (PostgreSQL 15 + API)
├── docker-compose.test.yml # Test stack (PostgreSQL 15 on :5433 + Redis + MinIO)
├── gqlgen.yml            # gqlgen code generation config
└── assets/               # Architecture diagrams
```

## API Endpoints

### Health Check

- **GET** `/v1/health`: Check the health of the API and its dependencies

**Response:**
```json
{
  "status": "success",
  "message": "Service is healthy",
  "data": {
    "database": "OK",
    "cache": "OK",
    "storage": "OK"
  }
}
```

### Version

- **GET** `/v1/version`: Get the current API version

**Response:**
```json
{
  "version": "1.7.0",
  "build_time": "2025-05-12T10:00:00Z",
  "git_commit": "a9e16fe"
}
```

### Authentication

- **POST** `/v1/auth/register`: Register a new user
- **POST** `/v1/auth/login`: Authenticate a user and receive a JWT token
- **POST** `/v1/auth/refresh`: Refresh an expired JWT token
- **POST** `/v1/auth/logout`: Invalidate the current JWT token

### Users

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/v1/users/register` | Register a new user | No |
| POST | `/v1/users/login` | Authenticate a user | No |
| GET | `/v1/users/profile` | Get user profile | Yes |
| PUT | `/v1/users/profile` | Update user profile | Yes |
| POST | `/v1/users/forgot-password` | Request password reset | No |
| POST | `/v1/users/reset-password` | Reset password with token | No |
| POST | `/v1/users/change-password` | Change password | Yes |
| POST | `/v1/users/deactivate` | Deactivate account | Yes |
| GET | `/v1/users/{id}/reviews` | Get user's reviews | No |
| GET | `/v1/users/{id}/favorites` | Get user's favorites | No |

### Anime

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/v1/anime` | Get all anime with pagination | No |
| GET | `/v1/anime/{id}` | Get anime by ID | No |
| POST | `/v1/anime` | Create anime | Yes |
| PUT | `/v1/anime/{id}` | Update anime | Yes |
| DELETE | `/v1/anime/{id}` | Delete anime | Yes |
| GET | `/v1/anime/search` | Search anime | No |
| GET | `/v1/anime/{id}/reviews` | Get anime reviews | No |
| POST | `/v1/anime/{id}/favorite` | Add to favorites | Yes |
| DELETE | `/v1/anime/{id}/favorite` | Remove from favorites | Yes |
| GET | `/v1/anime/favorites` | Get user's favorites | Yes |

### Genres

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/v1/genres` | Get all genres | No |
| GET | `/v1/genres/{id}` | Get genre by ID | No |
| POST | `/v1/genres` | Create genre | Yes |
| PUT | `/v1/genres/{id}` | Update genre | Yes |
| DELETE | `/v1/genres/{id}` | Delete genre | Yes |
| GET | `/v1/genres/search` | Search genres | No |
| POST | `/v1/genres/bulk` | Bulk create genres | Yes |
| DELETE | `/v1/genres/bulk` | Bulk delete genres | Yes |
| GET | `/v1/genres/{id}/anime` | Get anime by genre | No |

### Tags

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/v1/tags` | Get all tags | No |
| GET | `/v1/tags/{id}` | Get tag by ID | No |
| POST | `/v1/tags` | Create tag | Yes |
| PUT | `/v1/tags/{id}` | Update tag | Yes |
| DELETE | `/v1/tags/{id}` | Delete tag | Yes |
| GET | `/v1/tags/search` | Search tags | No |
| POST | `/v1/tags/bulk` | Bulk create tags | Yes |
| DELETE | `/v1/tags/bulk` | Bulk delete tags | Yes |
| GET | `/v1/tags/{id}/anime` | Get anime by tag | No |

### Reviews

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/v1/reviews` | Get all reviews with pagination | No |
| GET | `/v1/reviews/{id}` | Get review by ID | No |
| POST | `/v1/reviews` | Create review | Yes |
| PUT | `/v1/reviews/{id}` | Update review | Yes |
| DELETE | `/v1/reviews/{id}` | Delete review | Yes |
| GET | `/v1/reviews/user/{id}` | Get user's reviews | No |
| GET | `/v1/reviews/anime/{id}` | Get anime's reviews | No |

### GraphQL

The GraphQL endpoint runs alongside REST and exposes the same business logic through a typed schema.

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/graphql` | POST | GraphQL query / mutation endpoint |
| `/v1/graphql/playground` | GET | Interactive GraphQL Playground (dev) |

**Available Queries:**

```graphql
query {
  anime(id: "1") { id title rating genres { name } }
  animes(page: 1, limit: 20) { total data { id title } }
  animesByTitle(title: "Naruto", page: 1, limit: 10) { total data { id title } }
  genre(id: "1") { id name }
  genres(page: 1, limit: 50) { id name }
}
```

**Available Mutations:**

```graphql
mutation {
  login(username: "user", password: "pass") { token }
  register(username: "new", email: "new@example.com", password: "pass") { id username }
  createAnime(title: "Attack on Titan", episodes: 75, status: "Completed") { id title }
  deleteAnime(id: "1")
}
```

Schema source: [`api/adapters/graphql/schema/schema.graphql`](api/adapters/graphql/schema/schema.graphql)  
Code generation config: [`gqlgen.yml`](gqlgen.yml)

## Documentation

### Swagger Documentation

Generate Swagger documentation:

```bash
swag init -g cmd/main.go --dir ./cmd,./api/adapters/http,./api/models,./internal/errors -o ./cmd/docs
```

Or simply:

```bash
make swagger
```

The REST documentation is available at `/swagger/index.html` when running the server. It provides detailed information about:

- Available endpoints
- Request/response schemas
- Authentication requirements
- Example requests and responses
- Error codes and their meanings

### GoDoc Documentation

Start the godoc server:

```bash
godoc -http=:6060
```

Access the documentation at `http://localhost:6060/pkg/myanimeapi/`
The GoDoc documentation is available at [pkg.go.dev](https://pkg.go.dev/github.com/italo-gouveia/myanimeapi).

## Testing

Unit tests live **co-located** with the source they test (`api/adapters/http/*_test.go`, etc.) following Go convention. Integration, E2E, and smoke tests live under `tests/`.

### Unit Tests (co-located)

```bash
# All packages (includes co-located unit tests)
go test ./api/... -v

# Specific adapter
go test ./api/adapters/http/... -v
```

### Integration Tests

```bash
go test ./tests/integration/... -v
```

### End-to-End Tests (requires PostgreSQL)

```bash
# Start test DB first, then run
make test-db
```

### Smoke Tests (requires running server)

```bash
go test ./tests/smoke/... -v
```

### All Tests with DB

```bash
# Starts postgres:5433 via Docker, runs all tests, stops container
make test-db
```

### Test Coverage

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Diagrams

Visual representations of the architecture, database schema, and key workflows. Images are stored in `./assets/`. Diagrams reflect the **pre-hexagonal** state and will be updated in a future iteration.

---

### **1. Architecture Diagram**
The high-level architecture of the MyAnimeAPI application, showing the interaction between components.

![Architecture Diagram](./assets/architecture_diagram.png)

---

### **2. Database Schema (ER Diagram)**
The Entity-Relationship (ER) diagram for the database, illustrating the relationships between `users`, `anime`, and `reviews` tables.

![Database Schema](./assets/entity_model_relationship.png)

---

### **3. Component Diagram**
A breakdown of the application's components, including handlers, middleware, and database interactions.

![Component Diagram](./assets/component_diagram.png)

---

### **4. Deployment Diagram**
The deployment flow, showing how the application is built, tested, and deployed using Docker and GitHub Actions.

![Deployment Diagram](./assets/deployment_diagram.png)

---

### **5. Client Request Flow**
The flow of a typical client request through the API, including middleware and handlers.

![Client Request Flow](./assets/client_request_flow.png)

---

### **6. Authentication Flow**
The JWT-based authentication flow, detailing how users authenticate and access protected routes.

![Authentication Flow](./assets/flow_diagram_for_auth.png)

---

### **7. Rate Limiting Flow**
The rate-limiting mechanism, showing how requests are throttled to prevent abuse.

![Rate Limiting Flow](./assets/flow_diagram_for_rate_limiting.png)

---

### **8. User Registration Sequence**
The sequence of steps involved in registering a new user.

![User Registration Sequence Diagram](./assets/sequence_diagram_for_user_registration.png)

---

### **9. User Authentication Sequence**
The sequence of steps involved in authenticate an user.

![User Authentication Sequence Diagram](./assets/sequence_diagram_for_user_authentication.png)

---

### **10. Anime Creation Sequence**
The sequence of steps involved in creating a new anime entry.

![Anime Creation Sequence Diagram](./assets/sequence_diagram_for_anime_creation.png)

---

### **How to Generate Diagrams**
1. Use tools like [Mermaid Live Editor](https://mermaid-js.github.io/mermaid-live-editor/) or [dbdiagram.io](https://dbdiagram.io/) to create the diagrams.
2. Export the diagrams as `.png` or `.svg` files.
3. Save the images in the `./assets/` directory.
4. Reference the images in the `README.md` as shown above.

---

### **Notes**
- Ensure all diagram images are stored in the `./assets/` directory.
- Use consistent naming conventions for the diagram files (e.g., `architecture_diagram.png`, `entity_model_relationship.png`).
- If you update the diagrams, make sure to update the corresponding images in the `resources` folder.


## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for a list of changes.

## Badges

[![Go Report Card](https://goreportcard.com/badge/github.com/italo-gouveia/myanimeapi)](https://goreportcard.com/report/github.com/italo-gouveia/myanimeapi)
[![GoDoc](https://godoc.org/github.com/italo-gouveia/myanimeapi?status.svg)](https://godoc.org/github.com/italo-gouveia/myanimeapi)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/italo-gouveia/myanimeapi/workflows/CI/badge.svg)](https://github.com/italo-gouveia/myanimeapi/actions)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=italo-gouveia_myanimeapi&metric=coverage)](https://sonarcloud.io/dashboard?id=italo-gouveia_myanimeapi)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=italo-gouveia_myanimeapi&metric=alert_status)](https://sonarcloud.io/dashboard?id=italo-gouveia_myanimeapi)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Gorilla Mux](https://github.com/gorilla/mux)
- [GORM](https://gorm.io/)
- [Swagger](https://swagger.io/)
- [JWT-Go](https://github.com/golang-jwt/jwt)
- [Validator](https://github.com/go-playground/validator)
- [Bluemonday](https://github.com/microcosm-cc/bluemonday)

## Contact
For questions or feedback, please reach out to italogouveiadev@outlook.com.


Project Link: [https://github.com/italo-gouveia/myanimeapi](https://github.com/italo-gouveia/myanimeapi)