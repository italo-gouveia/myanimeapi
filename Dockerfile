# Stage 1: Build the Go binary
FROM golang:1.22 AS build

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -o main ./cmd/main.go

# Stage 2: Create a smaller image
FROM alpine:latest

WORKDIR /root/

# Install certificates
RUN apk --no-cache add ca-certificates

# Copy the binary from the build stage
COPY --from=build /app/main .

# Expose the port
EXPOSE 8080

# Command to run the binary
CMD ["./main"]
