# Stage 1: Build
FROM golang:1.22 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /root/
RUN apk --no-cache add ca-certificates postgresql-client

COPY --from=build /app/main .

# Ensure main is executable
RUN chmod +x /root/main

EXPOSE 8080
CMD ["./main"]
