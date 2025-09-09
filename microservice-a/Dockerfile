# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory inside container
WORKDIR /app

# Install protobuf compiler, protoc, and bash shell
RUN apk add --no-cache protobuf-dev protoc bash

# Install Go protobuf plugins
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install swag CLI for Swagger docs generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Update PATH so installed binaries are available
ENV PATH="${PATH}:$(go env GOPATH)/bin"

# Copy go.mod and go.sum, then download dependencies
COPY microservice-a/go.mod microservice-a/go.sum ./
RUN go mod download

# Copy shared proto definitions to the container
COPY shared-proto/ ./shared-proto/

# Copy microservice-a source code into the container
COPY microservice-a/ ./

# Generate protobuf gRPC code from sensor.proto file
RUN protoc --go_out=pb --go-grpc_out=pb \
    --proto_path=shared-proto shared-proto/sensor.proto

# Generate Swagger documentation based on Go code annotations
RUN swag init -g cmd/main.go -o docs

# Build the Go application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./cmd/main.go

# Final stage
FROM alpine:latest

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the built binary from the builder
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./main"]
