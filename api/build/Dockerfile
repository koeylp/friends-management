# Build Stage
FROM golang:1.23.2-alpine3.20 AS builder
WORKDIR /app

# Copy go.mod and go.sum from the api folder
COPY api/go.mod .
COPY api/go.sum .
RUN go mod download

# Copy the rest of the source code from api folder
COPY api/ .

# Build the Go application
RUN go build -o main ./cmd/main/

# Final Stage (Production)
FROM alpine:3.20
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy the .env file from the api/build folder
COPY api/.env .

CMD ["./main"]

# docker build -t friends-management -f api/build/Dockerfile .
