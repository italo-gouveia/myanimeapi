# Stage 1: Build the Go binary
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
RUN swag init --dir ./cmd,./pkg/handlers,./pkg/models --output ./cmd/docs

# Build the Go binary
RUN go build -o main ./cmd

# Stage 2: Run the Go binary
FROM alpine:latest

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy the Swagger documentation
COPY --from=builder /app/cmd/docs ./cmd/docs

# Copy the .env file
COPY .env .

# Change ownership of the application files to the non-root user
RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]