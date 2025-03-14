# MyAnimeAPI

MyAnimeAPI is a RESTful API built with Go (Golang) that allows users to manage anime, reviews, and user accounts. It provides endpoints for creating, reading, updating, and deleting (CRUD) anime and reviews, as well as user authentication and authorization using JWT (JSON Web Tokens).

## Features

- **User Management**: Register, authenticate, and manage user accounts.
- **Anime Management**: Create, read, update, and delete anime entries.
- **Review Management**: Create, read, update, and delete reviews for anime.
- **Authentication & Authorization**: Secure endpoints using JWT tokens.
- **Pagination**: Retrieve paginated lists of users and reviews.
- **Database Integration**: Uses PostgreSQL for data storage and GORM for ORM (Object-Relational Mapping).

## Technologies Used

- **Go (Golang)**: The primary programming language.
- **Gorilla Mux**: A powerful HTTP router and URL matcher for building Go web servers.
- **GORM**: An ORM library for Go that supports PostgreSQL, MySQL, SQLite, and more.
- **JWT (JSON Web Tokens)**: Used for user authentication and authorization.
- **PostgreSQL**: A powerful, open-source relational database system.
- **Docker**: Containerization for easy deployment and development.

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Docker (optional)

### Installation

Clone the repository:

```bash
git clone https://github.com/yourusername/myanimeapi.git
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
Generate the swagger documentation:

```bash
swag init --dir ./cmd,./pkg/handlers,./pkg/models --output ./cmd/docs
```

This command it will generate the swagger docs. And then, after you run the application locally, you will be abble to go to `http://localhost:8080/swagger/index.html`.

----------
Generate Mocks for DBInterface:
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
go test -v ./pkg/handlers
```

### Docker Setup

Build and run the Docker containers:

```bash
docker-compose up --build
```

This will start both the PostgreSQL database and the Go API server.

**Access the API:**  
The API will be available at `http://localhost:8080/v1`.
The Swagger it will be available at `http://localhost:8080/swagger/index.html`.

## API Endpoints

### Authentication

- **POST** `/v1/auth/register`: Register a new user.
- **POST** `/auth/authenticate`: Authenticate a user and receive a JWT token.

### Users

- **GET** `/v1/users`: Retrieve a paginated list of users (Admin only).
- **GET** `/v1/users/{id}`: Retrieve a specific user by ID.
- **POST** `/v1/users`: Create a new user.
- **PUT** `/v1/users/{id}`: Update an existing user.
- **DELETE** `/v1/users/{id}`: Delete a user.

### Anime

- **GET** `/v1/anime`: Retrieve all anime entries.
- **GET** `/v1/anime/{id}`: Retrieve a specific anime by ID.
- **POST** `/v1/anime`: Create a new anime entry (Authenticated users only).
- **PUT** `/v1/anime/{id}`: Update an existing anime entry (Authenticated users only).
- **DELETE** `/v1/anime/{id}`: Delete an anime entry (Authenticated users only).

### Reviews

- **GET** `/v1/reviews/{id}`: Retrieve a specific review by ID.
- **POST** `/v1/reviews`: Create a new review (Authenticated users only).
- **PUT** `/v1/reviews/{id}`: Update an existing review (Authenticated users only).
- **DELETE** `/v1/reviews/{id}`: Delete a review (Authenticated users only).

## Diagrams

### Architecture Diagram

```
+-------------------+       +-------------------+       +-------------------+
|   Client (HTTP)   | <---> |   Go API Server   | <---> |   PostgreSQL DB   |
+-------------------+       +-------------------+       +-------------------+
```

### Flow Diagram

![Client Request Flow](/myanimeapi/resources/client_request_flow.png)

### Model Entity Relationship

![Model Entity Relationship](/myanimeapi/resources/entity_model_relationship.png)

## Testing

To run the tests, use the following command:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
