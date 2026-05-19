# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies for CGO (required by go-sqlite3)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary with static linking and optimizations
RUN CGO_ENABLED=1 GOOS=linux go build -ldflags="-s -w" -o crowd-vote-api main.go

# Run stage
FROM alpine:latest

# Install sqlite for debugging (optional) and ca-certificates
RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/crowd-vote-api .

# Create a directory for the database
RUN mkdir /app/data

# Environment variable for database path (optional, but good practice)
ENV DB_PATH=/app/data/crowd_vote.db

EXPOSE 8080

# Command to run the application
CMD ["./crowd-vote-api"]
