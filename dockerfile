# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod file (go.sum will be generated if needed)
COPY go.mod ./

# Download dependencies and generate go.sum
RUN go mod download && go mod tidy

# Copy the source code
COPY . .

# Build the application for linux/amd64 with optimizations
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/main ./cmd/main.go

# Stage 2: Create the final image
FROM alpine:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy migrations ONLY if needed by the application at runtime
COPY migrations ./migrations

# DO NOT COPY .env file! Configuration comes from Cloud Run environment.

# Expose the port Cloud Run expects (and the app should listen on via $PORT)
EXPOSE 8080

# Use ENTRYPOINT to run the application
ENTRYPOINT ["./main"]