FROM golang:1.23.2-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY ./cmd/main.go ./

COPY . .

# Build the application
RUN go build -o main .

# Expose port 3000
EXPOSE 8080

# Run the binary
CMD ["./main"]
