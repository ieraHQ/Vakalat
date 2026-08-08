# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY backend/api/go.mod backend/api/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/api/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /vakalat-backend

# Final stage
FROM alpine:latest

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /vakalat-backend ./

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Expose port
EXPOSE 8080

# Run the application
CMD ["./vakalat-backend"]