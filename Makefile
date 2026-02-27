.PHONY: all build install proto proto-example clean test

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Binary names
PROTOC_GEN_TEMPORAL=protoc-gen-golemporal

# Proto files
PROTO_FILES=$(wildcard example/*.proto)
PROTO_GO_FILES=$(PROTO_FILES:.proto=.pb.go)

# Default target
all: install

# Build the protoc plugin
build: $(PROTOC_GEN_TEMPORAL)

$(PROTOC_GEN_TEMPORAL):
	cd cmd/protoc-gen-golemporal && $(GOBUILD) -o ../../bin/$@ .

# Install the protoc plugin
install: $(PROTOC_GEN_TEMPORAL)
	cd cmd/protoc-gen-golemporal && $(GOINSTALL) .

# Install dependencies
deps:
	$(GOMOD) download
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GOGET) go.temporal.io/sdk@latest
	$(GOGET) go.uber.org/zap@latest

# Generate code from proto files
proto: install
	protoc \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		--proto_path=. \
		example/example.proto

# Generate starter code only
proto-starter: install
	protoc \
		--golemporal_out=starter=true:generated \
		--proto_path=. \
		example/example.proto

# Generate worker code only
proto-worker: install
	protoc \
		--golemporal_out=worker=true:generated \
		--proto_path=. \
		example/example.proto

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf bin/
	rm -rf generated/
	rm -f example/*.pb.go
	rm -f example/*.grpc.pb.go

# Run tests
test:
	$(GOTEST) -v ./...

# Run example worker
run-worker:
	cd example && $(GOCMD) run worker/main.go

# Run example starter
run-starter:
	cd example && $(GOCMD) run starter/main.go
