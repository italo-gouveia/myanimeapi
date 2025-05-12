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
  - [Project Structure](#project-structure)
  - [API Endpoints](#api-endpoints)
    - [Health Check](#health-check)
    - [Version](#version)
    - [Authentication](#authentication)
    - [Users](#users)
      - [User Profile Update](#user-profile-update)
      - [Change Password](#change-password)
      - [Deactivate Account](#deactivate-account)
    - [Anime](#anime)
    - [Genres](#genres)
      - [Search Genres](#search-genres)
      - [Bulk Create Genres](#bulk-create-genres)
      - [Bulk Delete Genres](#bulk-delete-genres)
    - [Tags](#tags)
    - [Reviews](#reviews)
    - [Favorites](#favorites)
  - [Documentation](#documentation)
    - [Swagger Documentation](#swagger-documentation)
    - [GoDoc Documentation](#godoc-documentation)
  - [Testing](#testing)
    - [Unit Tests](#unit-tests)
    - [Integration Tests](#integration-tests)
    - [End-to-End Tests](#end-to-end-tests)
    - [Smoke Tests](#smoke-tests)
    - [Test Coverage](#test-coverage)
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
- **User Management**: Register, authenticate, and manage users
- **Review Management**: Add, update, and delete reviews for anime
- **Genre Management**: Manage anime genres with CRUD operations
- **Tag Management**: Manage anime tags with CRUD operations
- **Authentication**: JWT-based authentication for secure access
- **Pagination**: Paginated responses for large datasets
- **Rate Limiting**: Protect endpoints from abuse with rate limiting
- **Swagger Documentation**: Auto-generated API documentation
- **Favorite Anime**: Add and manage favorite anime entries
- **Error Handling**: Structured error responses with detailed information and consistent format
- **Health Monitoring**: Comprehensive health check endpoint with detailed service status
- **Security**: Regular security scanning and vulnerability checks
- **Dependency Management**: Proper dependency injection and service initialization

## Technologies Used

- **Go**: Backend programming language (v1.23.0)
- **Gorilla Mux**: HTTP router and dispatcher (v1.8.1)
- **GORM**: ORM for database interactions (v1.25.12)
- **PostgreSQL**: Relational database
- **JWT**: JSON Web Tokens for authentication (v5.2.2)
- **Swagger**: API documentation (v1.16.4)
- **Docker**: Containerization for easy deployment and development
- **Validator**: Input validation (v10.25.0)

## Getting Started

### Prerequisites

- Go 1.23.0 or higher
- PostgreSQL
- Docker (optional)

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

```
myanimeapi/
├── api/
│   ├── handlers/     # HTTP request handlers
│   ├── models/       # Data models and DTOs
│   ├── services/     # Business logic layer
│   ├── repositories/ # Data access layer
│   ├── middleware/   # HTTP middleware
│   ├── auth/         # Authentication related code
│   └── utils/        # Utility functions
├── cmd/
│   └── main.go       # Application entry point
├── internal/
│   ├── config/       # Configuration management
│   ├── database/     # Database connection and setup
│   ├── errors/       # Custom error types
│   └── utils/        # Internal utilities
├── tests/
│   ├── e2e/         # End-to-end tests
│   ├── integration/ # Integration tests
│   └── smoke/       # Smoke tests
└── frontend/        # Next.js frontend application
```

## API Endpoints

### Health Check

- **GET** `/v1/health`: Check the health of the API and its dependencies

**Response:**
```json
{
  "status": "success",
  "message": "Service is healthy",
  "data": "OK"
}
```

**Error Response:**
```json
{
  "status": "error",
  "error": {
    "code": "SERVICE_UNAVAILABLE",
    "message": "Database connection failed",
    "details": "The database connection is not healthy."
  }
}
```

### Version

- **GET** `/v1/version`: Get the current API version

**Response:**
```json
{
  "version": "1.6.0"
}
```

### Authentication

- **POST** `/v1/auth/register`: Register a new user
- **POST** `/v1/auth/login`: Authenticate a user and receive a JWT token

### Users

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| POST | `/v1/users/register` | Register a new user | No |
| POST | `/v1/users/login` | Authenticate a user and receive a JWT token | No |
| GET | `/v1/users/profile` | Get the authenticated user's profile | Yes |
| PUT | `/v1/users/profile` | Update the authenticated user's profile | Yes |
| POST | `/v1/users/change-password` | Change the authenticated user's password | Yes |
| POST | `/v1/users/deactivate` | Deactivate the authenticated user's account | Yes |

#### User Profile Update
Update the authenticated user's profile information.

**Request Body:**
```json
{
  "username": "new_username",
  "email": "new_email@example.com",
  "profile_pic": "https://example.com/profile.jpg",
  "bio": "User's biography",
  "social_links": {
    "twitter": "https://twitter.com/username",
    "instagram": "https://instagram.com/username"
  },
  "genre_ids": [1, 2, 3]
}
```

#### Change Password
Change the authenticated user's password.

**Request Body:**
```json
{
  "current_password": "old_password",
  "new_password": "new_password"
}
```

#### Deactivate Account
Deactivate the authenticated user's account.

**Request Body:**
```json
{
  "password": "current_password"
}
```

### Anime

- **GET** `/v1/anime`: Get all anime entries with pagination
- **GET** `/v1/anime/{id}`: Retrieve a specific anime by ID
- **POST** `/v1/anime`: Create a new anime entry (Authenticated users only)
- **PUT** `/v1/anime/{id}`: Update an existing anime entry (Authenticated users only)
- **DELETE** `/v1/anime/{id}`: Delete an anime entry (Authenticated users only)
- **POST** `/v1/anime/{id}/favorite`: Add anime to favorites (Authenticated users only)
- **DELETE** `/v1/anime/{id}/favorite`: Remove anime from favorites (Authenticated users only)
- **GET** `/v1/anime/favorites`: Get user's favorite anime (Authenticated users only)

### Genres

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/v1/genres` | Get all genres | No |
| GET | `/v1/genres/{id}` | Get genre by ID | No |
| POST | `/v1/genres` | Create a new genre | Yes |
| PUT | `/v1/genres/{id}` | Update a genre | Yes |
| DELETE | `/v1/genres/{id}` | Delete a genre | Yes |
| GET | `/v1/genres/search` | Search genres by name or description | No |
| POST | `/v1/genres/bulk` | Create multiple genres at once | Yes |
| DELETE | `/v1/genres/bulk` | Delete multiple genres at once | Yes |

#### Search Genres
Search for genres by name or description with pagination support.

**Query Parameters:**
- `query` (required): Search term
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 10)

**Response:**
```json
{
  "data": [
    {
      "id": 1,
      "name": "Action",
      "created_at": "2024-02-20T19:27:00Z",
      "updated_at": "2024-02-20T19:27:00Z"
    }
  ],
  "pagination": {
    "total": 100,
    "page": 1,
    "limit": 10,
    "pages": 10
  }
}
```

#### Bulk Create Genres
Create multiple genres in a single request.

**Request Body:**
```json
{
  "genres": [
    {
      "name": "Action"
    },
    {
      "name": "Comedy"
    }
  ]
}
```

#### Bulk Delete Genres
Delete multiple genres in a single request.

**Request Body:**
```json
{
  "ids": [1, 2, 3]
}
```

### Tags

| Method | Endpoint | Description | Authentication Required |
|--------|----------|-------------|------------------------|
| GET | `/v1/tags` | Get all tags | No |
| GET | `/v1/tags/{id}` | Get tag by ID | No |
| POST | `/v1/tags` | Create a new tag | Yes |
| PUT | `/v1/tags/{id}` | Update a tag | Yes |
| DELETE | `/v1/tags/{id}` | Delete a tag | Yes |
| GET | `/v1/tags/search` | Search tags by name | No |
| POST | `/v1/tags/bulk` | Create multiple tags at once | Yes |
| DELETE | `/v1/tags/bulk` | Delete multiple tags at once | Yes |

### Reviews

- **GET** `/v1/reviews`: Get all reviews with pagination
- **GET** `/v1/reviews/{id}`: Get a specific review by ID
- **POST** `/v1/reviews`: Create a new review (Authenticated users only)
- **PUT** `/v1/reviews/{id}`: Update an existing review (Authenticated users only)
- **DELETE** `/v1/reviews/{id}`: Delete a review (Authenticated users only)
- **GET** `/v1/reviews/user/{userId}`: Get all reviews by a specific user
- **GET** `/v1/reviews/anime/{animeId}`: Get all reviews for a specific anime

### Favorites

- **GET** `/v1/favorites`: Get user's favorite anime (Authenticated users only)
- **POST** `/v1/favorites/{animeId}`: Add anime to favorites (Authenticated users only)
- **DELETE** `/v1/favorites/{animeId}`: Remove anime from favorites (Authenticated users only)
- **GET** `/v1/favorites/check/{animeId}`: Check if anime is in user's favorites (Authenticated users only)

## Documentation

### Swagger Documentation

Generate Swagger documentation:

```bash
swag init --dir ./cmd,./api/handlers,./api/models,./internal/errors --output ./cmd/docs
```

Access the Swagger UI at `http://localhost:8080/swagger/index.html`

### GoDoc Documentation

Start the godoc server:

```bash
godoc -http=:6060
```

Access the documentation at `http://localhost:6060/pkg/myanimeapi/`

## Testing

### Unit Tests
```bash
go test -v ./api/handlers
```

### Integration Tests
```bash
go test -v ./tests/integration
```

### End-to-End Tests
```bash
go test -v ./tests/e2e
```

### Smoke Tests
```bash
go test -v ./tests/smoke
```

### Test Coverage
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

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


## Contributing

1. Fork the repository
2. Create a new branch for your feature or bugfix
3. Commit your changes with clear and descriptive messages
4. Submit a pull request

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
[![DeepWiki Documentation](https://img.shields.io/badge/docs-DeepWiki-blue)](https://deepwiki.com/italo-gouveia/myanimeapi)
[![GitDiagram](https://img.shields.io/badge/architecture-GitDiagram-blue)](https://gitdiagram.com/italo-gouveia/myanimeapi)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments
- [Gorilla Mux](https://github.com/gorilla/mux) for routing
- [GORM](https://gorm.io/) for database interactions
- [JWT](https://jwt.io/) for authentication
- [Swagger](https://swagger.io/) for API documentation
- [Docker](https://www.docker.com/) for Containerization
- [Validator](https://github.com/go-playground/validator) for input validation

## Contact
For questions or feedback, please reach out to italogouveiadev@outlook.com.

