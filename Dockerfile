# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install required system packages
RUN apk add --no-cache git make protoc

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o /app/server ./cmd/server

# Final stage
FROM alpine:3.18

WORKDIR /app

# Install required runtime packages
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/server .

# Expose ports
EXPOSE 8080 50051

# Run the application
CMD ["./server"]
