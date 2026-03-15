# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install build dependencies for CGO/SQLite (if needed, though modernc.org/sqlite shouldn't need it)
RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the application
RUN go build -o server .

# Run stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/server .

# Expose the port
EXPOSE 8080

# Command to run the application
CMD ["./server"]
