# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/healthchecker ./cmd/healthchecker

# Final stage
FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/healthchecker .

# Copy templates
COPY templates/ /app/templates/

# Expose port 8080
EXPOSE 8080

# Run the binary
CMD ["./healthchecker"]
