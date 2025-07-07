# Use the official Go image as the base image
FROM golang:1.24-alpine AS builder

# Set the working directory
WORKDIR /app

# Install git and ca-certificates (needed for go mod download)
RUN apk add --no-cache git ca-certificates

COPY . .

RUN go build -o main .

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create a non-root user
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Create necessary directories
RUN mkdir -p /app/data /app/img && \
    chown -R appuser:appgroup /app

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Copy the example.env file (for reference)
COPY --from=builder /app/example.env .

# Change ownership of the binary
RUN chown appuser:appgroup main

# Switch to non-root user
USER appuser

# Expose any necessary ports (if needed for future features)
# EXPOSE 8080

# Create a volume for the database and other persistent data
VOLUME ["/app/data", "/app/img"]

# Set environment variables for the database path
ENV DB_PATH=/app/data/db.sqlite
ENV DEV_MODE=false

# The application will look for db.sqlite in the current directory
# We'll mount the external database to /app/data/db.sqlite
# The application needs to be modified to use the DB_PATH environment variable

# Run the application
CMD ["./main"] 