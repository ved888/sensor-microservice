# Sensor Microservices

This project demonstrates a microservices setup with gRPC communication between services.

---

## 📂 Project Structure

```bash
.
├── microservice-a
│   ├── cmd
│   │   └── main.go
│   ├── internal
│   │   ├── api
│   │   │   ├── grpcclient
│   │   │   └── http
│   │   ├── repository
│   │   └── usecase
│   ├── model
│   └── pb
├── microservice-b
│   ├── cmd
│   │   └── main.go
│   ├── database
│   │   ├── db.go
│   │   └── migrations
│   ├── internal
│   │   ├── api
│   │   │   ├── grpcclient
│   │   │   └── http
│   │   ├── repository
│   │   └── usecase
│   ├── middleware
│   ├── model
│   └── pb
├── shared-proto
│   └── sensor.proto
└── Makefile
````
# 🚀 Protocol Buffers & gRPC Code Generation

This project uses **Protocol Buffers (protobuf)** and **gRPC** for communication between microservices.  
The `.proto` definitions are stored in the [`shared-proto/`](./shared-proto) folder.

---

## 📂 Proto file location


---

## ⚙️ Installation (first time only)

Make sure you have the required tools installed:

```bash
# Install protobuf compiler
sudo apt-get install -y protobuf-compiler

# Check version (must be >= 3.20)
protoc --version

# Install Go plugins for protobuf & gRPC
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Ensure Go bin path is available
export PATH="$PATH:$(go env GOPATH)/bin"
```
🔨 Generate Go code

Run the following commands from the project root:
```# Generate gRPC & protobuf code for Microservice A
protoc --go_out=microservice-a/pb --go-grpc_out=microservice-a/pb \
  --proto_path=shared-proto shared-proto/sensor.proto

# Generate gRPC & protobuf code for Microservice B
protoc --go_out=microservice-b/pb --go-grpc_out=microservice-b/pb \
  --proto_path=shared-proto shared-proto/sensor.proto
```
or
```
make proto-gen
```

### Clean Generated Files
If you want to remove generated files:
```
make clean-proto
```
After running, you will see generated files in each microservice:

```
microservice-a/pb/
 ├── sensor.pb.go
 └── sensor_grpc.pb.go

microservice-b/pb/
 ├── sensor.pb.go
 └── sensor_grpc.pb.go
```