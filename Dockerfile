FROM golang:1.23.2-alpine3.20
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/main/
CMD ["./main"]
