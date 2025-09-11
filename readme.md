# Sensor Microservices

This project demonstrates a microservices setup with gRPC communication between services.

---
## 🛠 Prerequisites

- Go >= 1.18 ([Install guide](https://golang.org/doc/install))
- Docker & Docker Compose ([Install guide](https://docs.docker.com/get-docker/))
- Protocol Buffers compiler (`protoc`) >= 3.20 ([Install guide](https://grpc.io/docs/protoc-installation/))
- MySQL 8+ running locally or via Docker

## 📥 Download / Clone Project
## Local Setup

1.Refer to [the official guide to install Go](https://golang.org/doc/install).
2.Install Go version _1.18_ or higher.

```sh
export GOPATH=$(go env GOPATH)
export PATH=$PATH:$GOPATH/bin
```

These should be applied in your
[`~/.bash_profile`](https://www.baeldung.com/linux/bashrc-vs-bash-profile-vs-profile).

For Mac users using Catalina or newer (or if you're using zsh), you should add this to
[`~/.zshenv`](https://carlosroso.com/the-right-way-to-migrate-your-bash-profile-to-zsh/) instead.

3. Test if you have Go properly installed by running `go version` in your terminal.

4. Clone the project by using below command

```sh
git clone git@github.com:ved888/sensor-microservice.git
cd sensor-microservice
```

## 📂 Project Structure

```bash
.
├── ARCHITECTURE.md
├── configs
│   ├── humidity.env
│   ├── light.env
│   ├── motion.env
│   ├── pressure.env
│   └── temperature.env
├── DATABASE_SCHEMA.md
├── docker-compose.yml
├── Makefile
├── microservice-a
│   ├── cmd
│   │   └── main.go
│   ├── Dockerfile
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   └── api
│   │       ├── grpcclient
│   │       │   ├── grpcclient.go
│   │       │   └── grpcclient_test.go
│   │       └── http
│   │           ├── http.go
│   │           └── http_test.go
│   ├── model
│   │   └── sensor_reading.go
│   └── pb
│       └── shared-proto
│           ├── sensor_grpc.pb.go
│           └── sensor.pb.go
├── microservice-b
│   ├── cmd
│   │   └── main.go
│   ├── database
│   │   ├── db.go
│   │   └── migrations
│   │       ├── 0001_create_users_table.down.sql
│   │       ├── 0001_create_users_table.up.sql
│   │       ├── 0002_create_sensor_readings_table.down.sql
│   │       ├── 0002_create_sensor_readings_table.up.sql
│   │       ├── 0003_alter_user_table.down.sql
│   │       └── 0003_alter_user_table.up.sql
│   ├── db.env
│   ├── Dockerfile
│   ├── docs
│   │   ├── docs.go
│   │   ├── swagger.json
│   │   └── swagger.yaml
│   ├── go.mod
│   ├── go.sum
│   ├── internal
│   │   ├── api
│   │   │   ├── grpc
│   │   │   │   └── grpc.go
│   │   │   └── http
│   │   │       ├── sensor.go
│   │   │       ├── sensor_test.go
│   │   │       ├── user.go
│   │   │       └── user_test.go
│   │   ├── repository
│   │   │   ├── sensor_repository.go
│   │   │   ├── sensor_repository_test.go
│   │   │   ├── user.go
│   │   │   └── user_repository_test.go
│   │   └── usecase
│   │       ├── user.go
│   │       └── user_usecase_test.go
│   ├── middleware
│   │   └── jwt.go
│   ├── model
│   │   ├── sensor_reading.go
│   │   └── user.go
│   ├── pb
│   │   └── shared-proto
│   │       ├── sensor_grpc.pb.go
│   │       └── sensor.pb.go
│   ├── README.md
│   └── utils
│       └── utils.go
├── postmanCollection
│   ├── sensor_microservice_api_Enviroment.postman_environment.json
│   └── sensor_microservice_api.postman_collection.json
├── readme.md
├── sensor_db.png
└── shared-proto
    └── sensor.proto
````

## 🏗️ Architecture & Database
This project includes dedicated documentation for system architecture and database schema:
### Architecture
* **File**: [architecture.md](ARCHITECTURE.md)
* **Contains**:
  * System overview and component details
  * Microservice A instances (data generators) and Microservice B (data receiver + API)
  * Data flow diagrams for sensor data generation and API requests
  * Deployment & scalability considerations
  * Security and JWT authentication overview
  
You can visualize the full architecture using the included Mermaid diagrams.

### Database Schema
* **File**: [database_schema.md](DATABASE_SCHEMA.md)
* **Contains**:
  * Entity Relationship Diagram (ERD)
  * Users table and sensor readings table with columns, data types, and constraints
  * Indexing strategy and query patterns
  * Sample SQL queries for filtering, aggregation, and pagination
  * Performance considerations and migration history

These documents provide a deeper understanding of how microservices interact, how data flows through the system, and how the database is structured for scalability and performance.

## 🚀 Protocol Buffers & gRPC Code Generation

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
```bash
# Generate gRPC & protobuf code for Microservice A
protoc --go_out=microservice-a/pb --go-grpc_out=microservice-a/pb \
  --proto_path=shared-proto shared-proto/sensor.proto

# Generate gRPC & protobuf code for Microservice B
protoc --go_out=microservice-b/pb --go-grpc_out=microservice-b/pb \
  --proto_path=shared-proto shared-proto/sensor.proto
```
or
```bash
make proto-gen
```

### Clean Generated Files
If you want to remove generated files:
```bash
make clean-proto
```
After running, you will see generated files in each microservice:

```bash
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
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### 🔨 Generate Swagger Docs
Generate Swagger documentation using swag:
```bash
# Using Makefile
make swagger-gen
```
Swagger docs are generated in `docs/` inside each microservice:
* Microservice A → `microservice-a/docs/swagger.json`
* Microservice B → `microservice-b/docs/swagger.json`

### 🌐 Access Swagger UI
* **Microservice B Swagger UI** → http://localhost:8000/swagger/index.html
* **Microservice A-1 (Humidity) Swagger UI** → http://localhost:8080/swagger/index.html
* **Microservice A-2 (Light) Swagger UI** → http://localhost:8081/swagger/index.html
* **Microservice A-3 (Motion) Swagger UI** → http://localhost:8082/swagger/index.html
* **Microservice A-4 (Pressure) Swagger UI** → http://localhost:8083/swagger/index.html
* **Microservice A-5 (Temperature) Swagger UI** → http://localhost:8084/swagger/index.html

## ▶️ Running the Microservices

This project uses to simplify running and stopping the services.  
Microservice **B** must start first (it provides gRPC + REST APIs), followed by Microservice **A** (data generator + gRPC client).

---

### 🏃 Run All Services

```bash
make docker-up
```
This will:

* Start **Microservice B** (gRPC on `:50051`, REST on `:8000`)
* Starts Microservice A instances (REST on :8080–:8084 depending on env)

### 🖥️ Start Individually
```bash
# Microservice B
make docker-up-b

# Microservice A (all instances)
make docker-up-a
```
### 🛑 Stop All Services
```bash
# Stop all
make docker-down

# Stop Microservice A only
make docker-down-a

# Stop Microservice B only
make docker-down-b
```

### 📜 Logs
```bash
# All logs
make logs

# Microservice A (all instances)
make logs-a

# Microservice B
make logs-b
```

### 🧪 Testing
Run tests using Makefile:
```bash
# All tests
make test

# Microservice A only
make test-a

# Microservice B only
make test-b

# Tests with coverage
make test-coverage
```
Coverage reports are generated in `coverage.html` for each service.

### 🌐 Access Services
* Microservice A REST API → http://localhost:8080
* Microservice B REST API → http://localhost:8000
* Microservice B gRPC → `localhost:50051`

## 🔐 Authentication (JWT-based)
This project uses JSON Web Tokens (JWT) to authenticate users and protect API endpoints. Authentication is implemented in Microservice B, which handles user signup, login, and validation for accessing protected routes.

### 📋 How It Works
Microservice B implements JWT authentication for API routes under /api.

#### 1. Signup (`POST /signup`)
Allows users to register by providing:
* `first_name` (optional)
* `last_name `(optional)
* `email` (required, must be unique)
* `password` (required, 6–60 characters)
* `role` (optional, defaults to "analyst", valid values: "admin", "analyst")

Example request:
```bash
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

#### 2. Login (`POST /login`)
Allows users to authenticate and obtain a JWT token by providing:
* `email` (required)
* `password` (required)

Example request:
```bash
curl -X POST http://localhost:8081/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "ved@example.com",
    "password": "securepassword123"
  }'
```
Example response:
```bash
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
#### 3. Accessing Protected Routes
Routes under `/api` are protected and require the JWT token.
Add the token in the request header as:
```bash
Authorization: Bearer <JWT_TOKEN>
```
Example:

```bash
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
