# Stage 1: Build the Go binary
FROM golang:1.25 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Install swag CLI for generating Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation (source dir updated after hexagonal refactor)
RUN swag init --dir ./cmd,./api/adapters/http,./api/models,./internal/errors --output ./cmd/docs

# Build the Go binary (statically linked)
RUN CGO_ENABLED=0 go build -o main ./cmd

# Stage 2: Run the Go binary
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/main .

# Ensure the binary is executable
RUN chmod +x /app/main

# Copy the Swagger documentation
COPY --from=builder /app/cmd/docs ./cmd/docs

# Change ownership of the application files to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]