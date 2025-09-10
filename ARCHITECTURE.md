# Sensor Microservices Architecture

## System Overview

This document describes the architecture of the sensor microservices system, which consists of multiple data generators (Microservice A instances) communicating with a centralized data receiver and API service (Microservice B).

## Architecture Diagram

```mermaid
graph TB
    subgraph "Client Layer"
        WEB[Web Client]
        API_CLIENT[API Client]
        POSTMAN[Postman/Testing Tools]
    end

    subgraph "Load Balancer (Optional)"
        LB[Load Balancer]
    end

    subgraph "Microservice A Instances (Data Generators)"
        A1[Microservice A-1<br/>Humidity Sensor<br/>ID1: A, ID2: 1<br/>Port: 8080]
        A2[Microservice A-2<br/>Light Sensor<br/>ID1: B, ID2: 2<br/>Port: 8081]
        A3[Microservice A-3<br/>Motion Sensor<br/>ID1: C, ID2: 3<br/>Port: 8082]
        A4[Microservice A-4<br/>Pressure Sensor<br/>ID1: D, ID2: 3<br/>Port: 8083]
        A5[Microservice A-5<br/>Temperature Sensor<br/>ID1: E, ID2: 3<br/>Port: 8084]
        AN[Microservice A-N<br/>Custom Sensor<br/>ID1: X, ID2: N<br/>Port: 808N]
    end

    subgraph "Microservice B (Data Receiver & API)"
        B_GRPC[gRPC Server<br/>Port: 50051]
        B_REST[REST API Server<br/>Port: 8000]
        B_AUTH[JWT Authentication]
        B_SWAGGER[Swagger Documentation]
    end

    subgraph "Database Layer"
        MYSQL[(MySQL Database<br/>Port: 3306)]
        USERS_TABLE[(Users Table)]
        SENSORS_TABLE[(Sensor Readings Table)]
    end

    subgraph "Infrastructure"
        DOCKER[Docker Containers]
        COMPOSE[Docker Compose]
        NETWORK[Bridge Network]
    end

    %% Client connections
    WEB --> LB
    API_CLIENT --> LB
    POSTMAN --> LB
    LB --> B_REST

    %% Direct API connections
    API_CLIENT -.-> B_REST
    POSTMAN -.-> B_REST

    %% Microservice A to B communication
    A1 -->|gRPC Stream| B_GRPC
    A2 -->|gRPC Stream| B_GRPC
    A3 -->|gRPC Stream| B_GRPC
    A4 -->|gRPC Stream| B_GRPC
    A5 -->|gRPC Stream| B_GRPC
    AN -->|gRPC Stream| B_GRPC

    %% Microservice B internal connections
    B_GRPC --> B_REST
    B_REST --> B_AUTH
    B_REST --> B_SWAGGER

    %% Database connections
    B_REST --> MYSQL
    MYSQL --> USERS_TABLE
    MYSQL --> SENSORS_TABLE

    %% Infrastructure
    A1 -.-> DOCKER
    A2 -.-> DOCKER
    A3 -.-> DOCKER
    A4 -.-> DOCKER
    A5 -.-> DOCKER
    AN -.-> DOCKER
    B_GRPC -.-> DOCKER
    B_REST -.-> DOCKER
    MYSQL -.-> DOCKER
    DOCKER --> COMPOSE
    COMPOSE --> NETWORK

    %% Styling
    classDef microserviceA fill:#e1f5fe
    classDef microserviceB fill:#f3e5f5
    classDef database fill:#e8f5e8
    classDef client fill:#fff3e0
    classDef infrastructure fill:#fce4ec

    class A1,A2,A3,A4,A5,AN microserviceA
    class B_GRPC,B_REST,B_AUTH,B_SWAGGER microserviceB
    class MYSQL,USERS_TABLE,SENSORS_TABLE database
    class WEB,API_CLIENT,POSTMAN,LB client
    class DOCKER,COMPOSE,NETWORK infrastructure
```

## Data Flow

### 1. Data Generation Flow
```mermaid
sequenceDiagram
    participant A as Microservice A
    participant B as Microservice B
    participant DB as MySQL Database

    A->>A: Generate Sensor Data
    A->>B: gRPC Stream (SensorData)
    B->>DB: Store Sensor Reading
    DB-->>B: Confirmation
    B-->>A: ACK Response
```

### 2. API Request Flow
```mermaid
sequenceDiagram
    participant C as Client
    participant B as Microservice B
    participant DB as MySQL Database

    C->>B: REST API Request
    B->>B: JWT Validation
    B->>DB: Query Data
    DB-->>B: Return Data
    B-->>C: JSON Response
```

## Component Details

### Microservice A (Data Generators)
- **Purpose**: Generate sensor data streams
- **Technology**: Go, Echo Framework, gRPC Client
- **Features**:
    - Configurable sensor types (Temperature, Humidity, Pressure,Light,Motion etc.)
    - Adjustable data generation frequency via REST API
    - gRPC streaming to Microservice B
    - Swagger documentation

### Microservice B (Data Receiver & API)
- **Purpose**: Receive, store, and serve sensor data
- **Technology**: Go, Echo Framework, gRPC Server, MySQL
- **Features**:
    - gRPC server for receiving sensor data
    - REST API for data retrieval and manipulation
    - JWT-based authentication and authorization
    - Database operations with filtering and pagination
    - Swagger documentation

### Database Schema
- **Users Table**: User authentication and authorization
- **Sensor Readings Table**: Time-series sensor data storage
- **Indexes**: Optimized for queries by ID combinations and time ranges

### Infrastructure
- **Docker**: Containerization for all services
- **Docker Compose**: Orchestration and networking
- **MySQL**: Persistent data storage
- **Bridge Network**: Inter-service communication

## Scalability Features

1. **Horizontal Scaling**: Multiple Microservice A instances
2. **Load Distribution**: Each instance handles different sensor types
3. **Database Optimization**: Indexed queries for performance
4. **Container Orchestration**: Easy deployment and scaling
5. **Network Isolation**: Secure inter-service communication

## Security Features

1. **JWT Authentication**: Secure API access
2. **Role-based Authorization**: Admin and Analyst roles
3. **Network Isolation**: Docker bridge network
4. **Input Validation**: API payload validation
5. **Secure Communication**: gRPC and HTTPS support
