# Using alpine variant of golang image for smaller size
FROM golang:1.18-alpine AS builder

# Set working directory
WORKDIR /app

# Copy the Go application source code
COPY ./ /app

# Install necessary packages
# - git might be needed for fetching Go modules
# - ca-certificates is needed for HTTPS requests
RUN apk add --no-cache git ca-certificates

# Build the Go application
RUN go build -o /app/image-updater

# Create the final lightweight image
FROM alpine:3.15

# Copy only the binary from the builder stage
COPY --from=builder /app/image-updater /bin/

# Use the binary as the default entry point
ENTRYPOINT ["image-updater"]
