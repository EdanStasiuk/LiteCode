FROM golang:1.24.5-alpine AS builder
RUN apk add --no-cache git bash

WORKDIR /app

# Copy the entire project to preserve relative paths
COPY . .

# Move to worker module
WORKDIR /app/apps/worker

# Download dependencies for worker module
RUN go mod tidy

# Build the worker binary
RUN go build -o worker ./main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/apps/worker/worker .
RUN apk add --no-cache ca-certificates
CMD ["./worker"]
