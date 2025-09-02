# # Builder stage
# FROM golang:1.24.5-alpine AS builder

# # Set working directory
# WORKDIR /app

# # Install git
# RUN apk add --no-cache git

# # Copy go.mod and go.sum first for caching
# COPY go.mod go.sum ./

# # Set Go build cache directory to persist via volume
# ENV GOCACHE=/root/.cache/go-build

# RUN go mod download

# # Copy the rest of the code
# COPY . .

# # Build the binary
# RUN go build -o server ./main.go

# # Final stage: lightweight runtime
# FROM alpine:latest

# WORKDIR /root/

# # Copy binary from builder
# COPY --from=builder /app/server .

# EXPOSE 8080

# # Command
# CMD ["./server"]

# -------------------
# Lightweight Dockerfile to bypass building inside Docker, only way I can 
# get around running out of storage within Docker.

# Use a minimal image
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Copy the pre-built binary
COPY server .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./server"]