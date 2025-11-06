### -------- Build Stage --------
# Use an official lightweight Go image with Alpine Linux
FROM golang:1.24.5-alpine AS builder

# Install necessary build tools
RUN apk add --no-cache git build-base ca-certificates

# Set the working directory inside the container
WORKDIR /go/src/app

# Copy the entire application source code into the container
COPY . .

# Build the Go binary with verbose output
RUN go build -o /go/bin/main -v ./main.go


### -------- Runtime Stage --------
# Use a minimal Alpine Linux image for runtime
FROM alpine:latest

# Set the working directory
WORKDIR /go/src/app

# Add CA certificates for HTTPS calls (if needed by your app)
RUN apk --no-cache add ca-certificates tzdata curl

# Copy the compiled Go binary from the builder stage
COPY --from=builder /go/bin/main /main

# Copy the public folder to the container
COPY public ./public

# Define the entrypoint for the container
ENTRYPOINT ["/main"]

# Expose the application port
EXPOSE 30000

# Health check
HEALTHCHECK --interval=30s --timeout=30s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:30000/healthcheck || exit 1
