# Stage 1: Build
FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy seluruh proyek
COPY . .

# Build aplikasi dari cmd/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o customer-search ./cmd/main.go

# Stage 2: Run
FROM alpine:latest

WORKDIR /app

# Copy binary dari builder
COPY --from=builder /app/customer-search .

EXPOSE 8080
CMD ["./customer-search"]