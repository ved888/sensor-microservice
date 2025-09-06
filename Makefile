# Makefile for Sensor Microservices
# Run protobuf code generation

PROTO_DIR = shared-proto
PROTO_FILE = $(PROTO_DIR)/sensor.proto

MICROSERVICES = microservice-a microservice-b

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
	@echo "🧹 Clean completed."


# -----------------------------------
# Run Microservices with proper order
# -----------------------------------
.PHONY: all run-a run-b run clean

MICROSERVICES = microservice-a microservice-b

all: run

run: run-b run-a

# Start Microservice B (gRPC + REST)
run-b:
	@echo "Starting Microservice B..."
	cd microservice-b && nohup go run cmd/main.go > microservice-b.log 2>&1 & echo $$! > microservice-b.pid
	@sleep 2 # optional short delay


# Start Microservice A (REST + gRPC client)
run-a:
	@echo "Starting Microservice A..."
	cd microservice-a && nohup go run cmd/main.go > microservice-a.log 2>&1 & echo $$! > microservice-a.pid

# Stop both microservices
clean:
	@echo "Stopping Microservice A..."
	@if [ -f microservice-a.pid ]; then kill $$(cat microservice-a.pid) 2>/dev/null || true; rm microservice-a.pid; fi
	@fuser -k 8080/tcp 2>/dev/null || true
	@echo "Stopping Microservice B..."
	@if [ -f microservice-b.pid ]; then kill $$(cat microservice-b.pid) 2>/dev/null || true; rm microservice-b.pid; fi
	@fuser -k 50051/tcp 2>/dev/null || true
	@echo "✅ All microservices stopped"
