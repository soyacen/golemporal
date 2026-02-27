# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

golemporal is a Temporal SDK framework for Go that uses protobuf-based code generation to define workflows and activities. It includes a protoc plugin (`protoc-gen-golemporal`) that generates type-safe Temporal workflow clients, activity clients, and worker registration code from proto service definitions.

## Common Commands

```bash
# Install dependencies
go mod tidy

# Build all packages
go build ./...

# Run tests
go test ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Install the protoc code generator
go install ./cmd/protoc-gen-golemporal

# Generate code from proto files (requires protoc and go install)
# Run from the example directory:
cd example && ./protoc.sh
```

## Architecture

### Proto Service Naming Convention

The code generator recognizes services by their naming convention:
- Services ending with `Workflow` are treated as workflow services (e.g., `GreeterWorkflow`)
- Services ending with `Activity` are treated as activity services (e.g., `AddActivity`)
- Only one workflow service per proto file is supported

### Generated Code Structure

For each proto file with Temporal services, the generator produces a `*_temporal.pb.go` file containing:

1. **Activity Client** - Interface for invoking activities from workflows, with constructor and implementation
2. **Activity Server** - Interface that activity implementations must satisfy
3. **Workflow Client** - Interface for starting workflow executions from starters, accepting functional options
4. **Workflow Server** - Interface that workflow implementations must satisfy
5. **Register Function** - Function to register workflow and activity implementations with a Temporal worker

### Key Components

| Path | Purpose |
|------|---------|
| `cmd/protoc-gen-golemporal/main.go` | Protoc plugin implementation |
| `starter/option.go` | Functional options for workflow start options (ID, timeouts, retry policy, etc.) |
| `example/starter/main.go` | Example workflow client (starts workflows) |
| `example/worker/main.go` | Example worker implementation (runs workflows/activities) |
| `example/api/example.proto` | Proto definitions with workflow and activity services |

### Workflow Implementation Pattern

1. Define proto services: one `Workflow` service, one or more `Activity` services
2. Run code generation: `./protoc.sh`
3. Implement workflow server interface (receives `workflow.Context`)
4. Implement activity server interfaces (receives `context.Context`)
5. Register both with a Temporal worker using the generated `Register*Worker` function
6. Use the generated workflow client to start workflows from your application

### Dependencies

- Go 1.25+
- Temporal Go SDK (`go.temporal.io/sdk`)
- Protocol Buffers (`google.golang.org/protobuf`)
- protoc (for code generation)

### Running Examples

Requires a Temporal server running on `localhost:7233` (default). Start the worker first, then the starter.
