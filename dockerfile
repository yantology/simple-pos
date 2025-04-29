# Stage 1: Build the application
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
# Build for linux/amd64 architecture, common for containers
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/main ./cmd/main.go

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy migrations (assuming they are needed at runtime, adjust if not)
COPY migrations ./migrations

# Copy .env file (ensure it exists or handle configuration differently)
# COPY .env .

# Expose the port the application runs on (adjust if different)
# The main.go swagger comment mentions 8000, config/app.go defaults to 3000 but reads APP_PORT
# Assuming APP_PORT will be set to 8000 in docker-compose
EXPOSE 8000

# Command to run the application
CMD ["./main"]
