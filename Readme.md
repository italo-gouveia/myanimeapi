# MyAnimeAPI

MyAnimeAPI is a RESTful API for managing anime, users, reviews, and authentication. It is built using Go, Gorilla Mux for routing, GORM for database interactions, and supports PostgreSQL as the database backend.

## Features

- **Anime Management**: Create, read, update, and delete anime entries.
- **User Management**: Register, authenticate, and manage users.
- **Review Management**: Add, update, and delete reviews for anime.
- **Authentication**: JWT-based authentication for secure access.
- **Pagination**: Paginated responses for large datasets.
- **Rate Limiting**: Protect endpoints from abuse with rate limiting.
- **Swagger Documentation**: Auto-generated API documentation.

## Technologies Used

- **Go**: Backend programming language.
- **Gorilla Mux**: HTTP router and dispatcher.
- **GORM**: ORM for database interactions.
- **PostgreSQL**: Relational database.
- **JWT**: JSON Web Tokens for authentication.
- **Swagger**: API documentation.
- **Docker**: Containerization for easy deployment and development.

## Getting Started

### Prerequisites

- Go 1.23 or higher
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
Generate the swagger documentation:

```bash
swag init --dir ./cmd,./pkg/handlers,./pkg/models --output ./cmd/docs
```

This command it will generate the swagger docs. And then, after you run the application locally, you will be abble to go to `http://localhost:8080/swagger/index.html`.

----------
Start the godoc server:

```bash
godoc -http=:6060
```
Open your browser and navigate to:

```bash
http://localhost:6060/pkg/myanimeapi/
```

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

- **GET** `/v1/anime`: Get all anime entries.
- **GET** `/v1/anime/{id}`: Retrieve a specific anime by ID.
- **POST** `/v1/anime`: Create a new anime entry (Authenticated users only).
- **PUT** `/v1/anime/{id}`: Update an existing anime entry (Authenticated users only).
- **DELETE** `/v1/anime/{id}`: Delete an anime entry (Authenticated users only).

### Reviews

- **GET** `/v1/reviews/{id}`: Get a specific review by ID.
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

Contributions are welcome! Please follow these steps:

1. Fork the repository.

2. Create a new branch for your feature or bugfix.

3. Commit your changes with clear and descriptive messages.

4. Submit a pull request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments
- [Gorilla Mux](https://github.com/gorilla/mux) for routing.

- [GORM](https://gorm.io/) for database interactions.

- [JWT](https://jwt.io/) for authentication.

- [Swagger](https://swagger.io/) for API documentation.

- [Docker](https://www.docker.com/) for Containerization

## Contact
For questions or feedback, please reach out to italogouveiadev@outlook.com.

