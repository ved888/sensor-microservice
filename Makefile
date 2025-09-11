# Makefile for Sensor Microservices (Docker + Test + Swagger + Per-Service)
# --------------------
# Configuration
# --------------------
PROTO_DIR = shared-proto
PROTO_FILE = $(PROTO_DIR)/sensor.proto
MICROSERVICES = microservice-a microservice-b

# --------------------
# Protobuf Commands
# --------------------
proto-gen:
	@for ms in $(MICROSERVICES); do \
		echo "Generating protobuf code for $$ms..."; \
		protoc --go_out=$$ms/pb --go-grpc_out=$$ms/pb \
		  --proto_path=$(PROTO_DIR) $(PROTO_FILE); \
	done
	@echo "✅ Protobuf code generation completed."

clean-proto:
	@for ms in $(MICROSERVICES); do \
		echo "Cleaning generated files in $$ms/pb..."; \
		rm -f $$ms/pb/*.pb.go; \
	done
	@echo "🧹 Protobuf clean completed."

# --------------------
# Docker Commands
# --------------------
.PHONY: all proto-gen clean-proto docker-up docker-down docker-clean test swagger-gen \
	docker-up-a docker-down-a docker-up-b docker-down-b

all: proto-gen docker-up

docker-up:
	@echo "Starting all services using Docker Compose..."
	docker-compose -f docker-compose.yml up -d
	@echo "✅ All services are up."

docker-down:
	@echo "Stopping all services using Docker Compose..."
	docker-compose -f docker-compose.yml down
	@echo "✅ All services are stopped."

docker-clean:
	@echo "Cleaning Docker containers, networks, volumes..."
	docker-compose -f docker-compose.yml down -v --rmi all --remove-orphans
	@echo "🧹 Docker clean completed."

# Run/Stop only Microservice A (all instances)
docker-up-a:
	@echo "Starting Microservice A (all instances) using Docker..."
	docker-compose -f docker-compose.yml up -d microservice-a-1 microservice-a-2 microservice-a-3 microservice-a-4 microservice-a-5
	@echo "✅ Microservice A instances are up."

docker-down-a:
	@echo "Stopping Microservice A (all instances) using Docker..."
	docker-compose -f docker-compose.yml stop microservice-a-1 microservice-a-2 microservice-a-3 microservice-a-4 microservice-a-5
	@echo "✅ Microservice A instances are stopped."

# Run/Stop only Microservice B
docker-up-b:
	@echo "Starting Microservice B using Docker..."
	docker-compose -f docker-compose.yml up -d microservice-b
	@echo "✅ Microservice B is up."

docker-down-b:
	@echo "Stopping Microservice B using Docker..."
	docker-compose -f docker-compose.yml stop microservice-b
	@echo "✅ Microservice B is stopped."

# --------------------
# Logs Commands
# --------------------
logs-a:
	@echo "Showing logs for Microservice A (all instances)..."
	docker-compose -f docker-compose.yml logs -f microservice-a-1 microservice-a-2 microservice-a-3 microservice-a-4 microservice-a-5

logs-b:
	@echo "Showing logs for Microservice B..."
	docker-compose -f docker-compose.yml logs -f microservice-b

logs:
	@echo "Showing logs for all services..."
	docker-compose -f docker-compose.yml logs -f

# --------------------
# Test Commands
# --------------------
test: proto-gen
	@echo "Running tests for Microservice A..."
	cd microservice-a && go test ./... -v
	@echo "Running tests for Microservice B..."
	cd microservice-b && go test ./... -v
	@echo "✅ All tests completed."

test-a: ## Run tests for Microservice A only
	@echo "$(BLUE)Testing Microservice A...$(NC)"
	@cd microservice-a && go test -v ./...

test-b: ## Run tests for Microservice B only
	@echo "$(BLUE)Testing Microservice B...$(NC)"
	@cd microservice-b && go test -v ./...

test-coverage: ## Run tests with coverage report
	@echo "$(BLUE)Running tests with coverage...$(NC)"
	@for ms in $(MICROSERVICES); do \
		if [ -d $$ms ]; then \
			echo "Coverage for $$ms:"; \
			cd $$ms && go test -coverprofile=coverage.out ./... || true; \
			if [ -f coverage.out ]; then \
				go tool cover -html=coverage.out -o coverage.html; \
			fi; \
			cd ..; \
		else \
			echo "⚠️ Directory $$ms not found, skipping..."; \
		fi; \
	done
	@echo "$(GREEN)✅ Coverage reports generated$(NC)"


# --------------------
# Swagger Commands
# --------------------
swagger-gen:
	@echo "Generating Swagger documentation for Microservice A..."
	cd microservice-a && swag init --output ./docs
	@echo "Generating Swagger documentation for Microservice B..."
	cd microservice-b && swag init --output ./docs
	@echo "✅ Swagger documentation generated."

# --------------------
# Clean All
# --------------------
clean: clean-proto docker-clean
	@echo "✅ All clean tasks completed."
