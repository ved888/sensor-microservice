# Build Stage
FROM golang:1.24-alpine AS builder

# Set working directory inside the container
WORKDIR /app

# Install protobuf compiler, protoc, and bash shell
RUN apk add --no-cache protobuf-dev protoc bash

# Install Go plugins for protobuf code generation
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Install swag CLI tool to generate Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Add Go binaries to the PATH environment variable
ENV PATH="${PATH}:$(go env GOPATH)/bin"

# Copy go.mod and go.sum to leverage Docker cache and download dependencies
COPY microservice-b/go.mod microservice-b/go.sum ./
RUN go mod download

# Copy shared proto definitions to the container
COPY shared-proto/ ./shared-proto/

# Copy microservice-b source code into the container
COPY microservice-b/ ./

# Generate protobuf gRPC code from sensor.proto file
RUN protoc --go_out=pb --go-grpc_out=pb \
    --proto_path=shared-proto shared-proto/sensor.proto

# Generate Swagger documentation based on Go code annotations
RUN swag init -g cmd/main.go -o docs

# Build the Go application as a statically linked binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final Stage
FROM alpine:latest

# Install CA certificates to handle HTTPS requests in the app
RUN apk --no-cache add ca-certificates

# Set working directory inside the container
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Copy database migration files needed at runtime
COPY --from=builder /app/database ./database

# Expose the application HTTP and gRPC ports
EXPOSE 8000 50051

# Define the default command to run the application
CMD ["./main"]
