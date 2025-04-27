# MyAnimeAPI

MyAnimeAPI is a RESTful API for managing anime, users, reviews, and authentication. It is built using Go, Gorilla Mux for routing, GORM for database interactions, and supports PostgreSQL as the database backend.

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
  - [API Endpoints](#api-endpoints)
    - [Authentication](#authentication)
    - [Users](#users)
    - [Anime](#anime)
    - [Reviews](#reviews)
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
  - [Testing](#testing)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [End-to-End Tests](#end-to-end-tests)
    - [Smoke Tests](#smoke-tests)
    - [Generating Test Coverage](#generating-test-coverage)
  - [Contributing](#contributing)
  - [Changelog](#changelog)
  - [Badges](#badges)
  - [License](#license)
  - [Acknowledgments](#acknowledgments)
  - [Contact](#contact)

---

## Features

- **Anime Management**: Create, read, update, and delete anime entries.
- **User Management**: Register, authenticate, and manage users.
- **Review Management**: Add, update, and delete reviews for anime.
- **Authentication**: JWT-based authentication for secure access.
- **Pagination**: Paginated responses for large datasets.
- **Rate Limiting**: Protect endpoints from abuse with rate limiting.
- **Swagger Documentation**: Auto-generated API documentation.
- **Favorite Anime**: Add and manage favorite anime entries.
- **Error Handling**: Structured error responses with detailed information.
- **Security**: Regular security scanning and vulnerability checks.

## Technologies Used

- **Go**: Backend programming language (v1.23.0).
- **Gorilla Mux**: HTTP router and dispatcher (v1.8.1).
- **GORM**: ORM for database interactions (v1.25.12).
- **PostgreSQL**: Relational database.
- **JWT**: JSON Web Tokens for authentication (v5.2.2).
- **Swagger**: API documentation (v1.16.4).
- **Docker**: Containerization for easy deployment and development.
- **Validator**: Input validation (v10.25.0).

## Getting Started

### Prerequisites

- Go 1.23.0 or higher
- PostgreSQL
- Docker (optional)

### Installation

Clone the repository:

```bash
git clone https://github.com/{{yourusername}}/myanimeapi.git
cd myanimeapi
```

Set up the database:

- Ensure PostgreSQL is running.
- Create a database named `myanimeapi`.
- Update the `.env` file with your PostgreSQL credentials.

Run the application:

```bash
go run cmd/main.go
```

The API will be available at `http://localhost:8080/v1`.

-------
How do I generate Swagger documentation?

```bash
swag init --dir ./cmd,./pkg/handlers,./pkg/models,./internal/errors --output ./cmd/docs
```

This command it will generate the swagger docs. And then, after you run the application locally, you will be abble to go to `http://localhost:8080/swagger/index.html`.

----------
Start the godoc server:

```bash
godoc -http=:6060
```
Open your browser and navigate to:

```bash
http://localhost:6060/api/myanimeapi/
```

----------
How do I generate mocks for testing?

Run the following command to generate a mock for the DBInterface:

```bash
mockgen -source=internal/db/db_interface.go -destination=internal/mocks/mock_db_interface.go -package=mocks
```
This will create a mock_db_interface.go file in the internal/db package.
If you are on the powershell, try this:
```shell
mockgen -source="$PWD\internal\db\db_interface.go" -destination="$PWD\internal\mocks\mock_db_interface.go" -package=mocks
```

----------
Run the Tests
Run the tests using the following command:

```bash
go test -v ./api/handlers
```

### Docker Setup

Build and run the Docker containers:

```bash
docker-compose up --build
```

This will start both the PostgreSQL database and the Go API server.

**Access the API:**  
The API will be available at `http://localhost:8080/v1`.
The Swagger documentation will be available at `http://localhost:8080/swagger/index.html`.

## API Endpoints

### Authentication

- **POST** `/v1/auth/register`: Register a new user.
- **POST** `/v1/auth/login`: Authenticate a user and receive a JWT token.

### Users

- **GET** `/v1/users`: Retrieve a paginated list of users (Admin only).
- **GET** `/v1/users/{id}`: Retrieve a specific user by ID.
- **POST** `/v1/users`: Create a new user.
- **PUT** `/v1/users/{id}`: Update an existing user.
- **DELETE** `/v1/users/{id}`: Delete a user.

### Anime

- **GET** `/v1/anime`: Get all anime entries with pagination.
- **GET** `/v1/anime/{id}`: Retrieve a specific anime by ID.
- **POST** `/v1/anime`: Create a new anime entry (Authenticated users only).
- **PUT** `/v1/anime/{id}`: Update an existing anime entry (Authenticated users only).
- **DELETE** `/v1/anime/{id}`: Delete an anime entry (Authenticated users only).
- **POST** `/v1/anime/{id}/favorite`: Add anime to favorites (Authenticated users only).
- **DELETE** `/v1/anime/{id}/favorite`: Remove anime from favorites (Authenticated users only).
- **GET** `/v1/anime/favorites`: Get user's favorite anime (Authenticated users only).

### Reviews

- **GET** `/v1/reviews/{id}`: Get a specific review by ID.
- **POST** `/v1/reviews`: Create a new review (Authenticated users only).
- **PUT** `/v1/reviews/{id}`: Update an existing review (Authenticated users only).
- **DELETE** `/v1/reviews/{id}`: Delete a review (Authenticated users only).

## Diagrams

This section provides visual representations of the application's architecture, database schema, component interactions, deployment flow, and key workflows.

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

## Testing

To run the tests, use the following command:

```bash
go test ./...
```

### Unit Tests
Run unit tests for the handlers:
```bash
go test -v ./api/handlers
```

### Integration Tests
Run integration tests to test database interactions:

```bash
go test -v ./tests/integration
```

### End-to-End Tests
Run end-to-end tests to test the API as a whole:

```bash
go test -v ./tests/e2e
```

### Smoke Tests
Run smoke tests to verify basic functionality:
```bash
go test -v ./tests/smoke
```

### Generating Test Coverage
Generate a test coverage report:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Commit your changes with clear and descriptive messages.
4. Submit a pull request.

## Changelog

See the [CHANGELOG.md](CHANGELOG.md) file for a detailed list of changes.

## Badges
[![Build Status](https://github.com/italo-gouveia/myAnimeAPI/actions/workflows/ci.yml/badge.svg)](https://github.com/italo-gouveia/myAnimeAPI/actions)
[![Test Coverage](https://codecov.io/gh/italo-gouveia/myAnimeAPI/branch/main/graph/badge.svg)](https://codecov.io/gh/italo-gouveia/myAnimeAPI)
[![Security Scan](https://github.com/italo-gouveia/myAnimeAPI/actions/workflows/security.yml/badge.svg)](https://github.com/italo-gouveia/myAnimeAPI/actions/workflows/security.yml)
[![Swagger](https://img.shields.io/badge/docs-swagger-blue)](https://github.com/italo-gouveia/myAnimeAPI/blob/main/cmd/docs/swagger.yaml)
[![GitHub release](https://img.shields.io/github/release/italo-gouveia/myAnimeAPI.svg)](https://github.com/italo-gouveia/myAnimeAPI/releases)
[![Go Version](https://img.shields.io/github/go-mod/go-version/italo-gouveia/myAnimeAPI)](https://github.com/italo-gouveia/myAnimeAPI)
[![GitHub stars](https://img.shields.io/github/stars/italo-gouveia/myAnimeAPI.svg?style=social)](https://github.com/italo-gouveia/myAnimeAPI/stargazers)
[![GitHub forks](https://img.shields.io/github/forks/italo-gouveia/myAnimeAPI.svg?style=social)](https://github.com/italo-gouveia/myAnimeAPI/network/members)
[![GitHub last commit](https://img.shields.io/github/last-commit/italo-gouveia/myAnimeAPI)](https://github.com/italo-gouveia/myAnimeAPI/commits/main)
[![GitHub issues](https://img.shields.io/github/issues/italo-gouveia/myAnimeAPI)](https://github.com/italo-gouveia/myAnimeAPI/issues)
[![GitHub pull requests](https://img.shields.io/github/issues-pr/italo-gouveia/myAnimeAPI)](https://github.com/italo-gouveia/myAnimeAPI/pulls)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
- [Gorilla Mux](https://github.com/gorilla/mux) for routing.
- [GORM](https://gorm.io/) for database interactions.
- [JWT](https://jwt.io/) for authentication.
- [Swagger](https://swagger.io/) for API documentation.
- [Docker](https://www.docker.com/) for Containerization.
- [Validator](https://github.com/go-playground/validator) for input validation.

## Contact
For questions or feedback, please reach out to italogouveiadev@outlook.com.

