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

## 📘 Swagger API Documentation

Both microservices expose **REST APIs** documented using **Swagger (swaggo)**.

### 📦 Installation

Install `swag` CLI tool:
```
go install github.com/swaggo/swag/cmd/swag@latest
```

### 🔨 Generate Swagger Docs

From the root of the project, run:
```
# Generate Swagger docs for Microservice A
cd microservice-a
swag init -g cmd/main.go -o docs

# Generate Swagger docs for Microservice B
cd ../microservice-b
swag init -g cmd/main.go -o docs
```
This will create a `docs/` folder in each service containing `swagger.json` and `swagger.yaml`.

### 🚀 Run Services with Swagger

Start each service:
```
# Start Microservice A (REST on :8080)
cd microservice-a
go run cmd/main.go

# Start Microservice B (REST on :8081)
cd microservice-b
go run cmd/main.go
```

### 🌐 Access Swagger UI
* **Microservice A Swagger UI** → http://localhost:8080/swagger/index.html
* **Microservice B Swagger UI** → http://localhost:8081/swagger/index.html

# ▶️ Running the Microservices

This project uses a **Makefile** to simplify running and stopping the services.  
Microservice **B** must start first (it provides gRPC + REST APIs), followed by Microservice **A** (data generator + gRPC client).

---

## 🏃 Run All Services

```bash
make run
```
This will:

* Start **Microservice B** (gRPC on `:50051`, REST on `:8081`)
* Then start Microservice A (REST on `:8080`)

Logs will be written to:
* `microservice-a/microservice-a.log`
* `microservice-b/microservice-b.log`

Process IDs are stored in `.pid` files for easier shutdown.

### 🖥️ Start Individually

* Start **Microservice B** only:
```
make run-b
```

* Start **Microservice A** only:
```
make run-a
```
* start multiple Service A
``` 
#Temperature
 cd microservice-a
 ENV_FILE=../configs/temperature.env PORT=8084 go run cmd/main.go

# Humidity
cd microservice-a
ENV_FILE=../configs/humidity.env PORT=8080 go run cmd/main.go

# Pressure
cd microservice-a
ENV_FILE=../configs/pressure.env PORT=8083 go run cmd/main.go

# Light
cd microservice-a
ENV_FILE=../configs/light.env PORT=8081 go run cmd/main.go

# Motion
cd microservice-a
ENV_FILE=../configs/motion.env PORT=8082 go run cmd/main.go
```

### 🛑 Stop All Services

Gracefully stop both services:
```
make clean
```
This will:

* Kill processes saved in `.pid` files
* Free up ports `8080`, `8081`, and `50051`

### 🌐 Access Services
* Microservice A REST API → http://localhost:8080
* Microservice B REST API → http://localhost:8081
* Microservice B gRPC → `localhost:50051`

## 🔐 Authentication (JWT-based)
This project uses JSON Web Tokens (JWT) to authenticate users and protect API endpoints. Authentication is implemented in Microservice B, which handles user signup, login, and validation for accessing protected routes.

### 📋 How It Works
#### 1. Signup (POST /signup)
Allows users to register by providing:
* `first_name` (optional)
* `last_name `(optional)
* `email` (required, must be unique)
* `password` (required, 6–60 characters)
* `role` (optional, defaults to "analyst", valid values: "admin", "analyst")

Example request:
```
curl -X POST http://localhost:8081/signup \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "Ved",
    "last_name": "Verma",
    "email": "ved@example.com",
    "password": "securepassword123",
    "role": "admin"
  }'
```

#### 1. Login (POST /login)
Allows users to authenticate and obtain a JWT token by providing:
* `email` (required)
* `password` (required)

Example request:
```
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ved@example.com",
    "password": "securepassword123"
  }'
```
Example response:
```
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
#### 3. Accessing Protected Routes
Routes under `/api` are protected and require the JWT token.
Add the token in the request header as:
```
Authorization: Bearer <JWT_TOKEN>
```
Example:

```
curl -X GET "http://localhost:8081/api/sensors?id1=A&id2=1&page=1&limit=10" \
  -H "accept: application/json" \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InZlZDEyM0BnbWFpbC5jb20iLCJleHAiOjE3NTczMzMzMjUsInJvbGUiOiJhZG1pbiIsInVzZXJfaWQiOjF9.ND6vcnN0bbJbS6pVh2Cdx_DY6LONSfB_hjyWFyXhbTA"
```

### 📦 JWT Token Details
* **Algorithm:** `HS256` (HMAC with SHA-256)
* **Claims:**
  * user_id: User’s ID
  * email: User’s email address
  * role: User’s role (admin or analyst)
  * exp: Token expiration time (in UNIX timestamp)

## Postman Collection
You can import the Postman collection and environment to test the APIs:
- **Collection:** [sensor_microservice_api.postman_collection.json](./postmanCollection/sensor_microservice_api.postman_collection.json)
- **Environment:** [sensor_microservice_api_Environment.postman_environment.json](./postmanCollection/sensor_microservice_api_Environment.postman_environment.json)
