# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

golemporal is a protoc plugin (`protoc-gen-golemporal`) that generates type-safe Temporal workflow clients, activity clients, and worker registration code from protobuf service definitions. It enables a proto-first approach to building Temporal workflows in Go.

## Common Commands

```bash
# Install dependencies and clean up unused imports
go mod tidy

# Build all packages in the main module
go build ./...

# Run tests (note: the main module has no tests; sdk-go/ contains a vendored Temporal SDK with its own test suite)
go test ./...

# Run a single test
go test -v -run TestName ./path/to/package

# Install the protoc code generator
go install ./cmd/protoc-gen-golemporal

# Generate code from proto files (requires protoc, protoc-gen-go, and protoc-gen-golemporal installed)
cd example && ./protoc.sh
```

## Architecture

### Proto Service Naming Convention

The code generator recognizes services by name suffix (case-sensitive):
- Services ending with `Workflow` generate workflow clients, servers, and registration functions
- Services ending with `Activity` generate activity clients, servers, and registration functions
- Multiple workflow and activity services per proto file are supported

### Code Generation Pipeline

The plugin (`cmd/protoc-gen-golemporal/main.go`) uses the `google.golang.org/protobuf/compiler/protogen` API to generate `*_temporal.pb.go` files. For each proto file containing recognized services:

1. **Activity Client** (`New*ActivityClient`) - Interface with methods that call `workflow.ExecuteActivity` using generated type constants; used inside workflow implementations
2. **Activity Server** (`*ActivityServer`) - Interface that activity implementations must satisfy; methods receive `context.Context`
3. **Workflow Client** (`New*WorkflowClient`) - Interface for starting workflow executions from application code; methods accept `starter.Option` functional options and return `(*Output, *protobuf.Metadata, error)`
4. **Workflow Server** (`*WorkflowServer`) - Interface that workflow implementations must satisfy; methods receive `workflow.Context`
5. **Register Functions** (`Register*Activity`, `Register*Workflow`) - Register implementations with a Temporal worker using `DisableAlreadyRegisteredCheck: true`

### Key Components

| Path | Purpose |
|------|---------|
| `cmd/protoc-gen-golemporal/main.go` | Protoc plugin implementation; contains hardcoded `Version` variable |
| `starter/option.go` | Functional options for workflow start options (ID, timeouts, retry policy, WaitResult, etc.) |
| `protobuf/metadata.proto` | Common `Metadata` message for workflow execution info (workflow_id, run_id, workflow_type, task_queue) |
| `example/api/example.proto` | Proto definitions with workflow and activity services |
| `example/protoc.sh` | Code generation script (includes `--proto_path=../../` for resolving imports) |
| `example/worker/main.go` | Example worker implementation |
| `example/starter/main.go` | Example workflow client/starter |
| `sdk-go/` | Vendored copy of temporalio/sdk-go (separate Go module, not part of main module) |

### Workflow Client Result Handling

Workflow client methods return three values: `(*Output, *protobuf.Metadata, error)`. By default, `WaitResult` is false and the client returns immediately after starting the workflow (output is nil). Use `starter.WaitResult(true)` to block until the workflow completes and populate the output:

```go
hc := api.NewHelloWorkflowClient(c, taskQueue)
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.WaitResult(true))
```

### Worker Registration Pattern

A single struct can implement multiple workflow or activity server interfaces. Register each service separately:

```go
wf := &GreeterWorkflowServer{
    addActivity: api.NewAddActivityClient(),
}
api.RegisterHelloWorkflow(w, wf)
api.RegisterGoodbyeWorkflow(w, wf)
api.RegisterAddActivity(w, &AddActivityServer{})
```

### Version Management

The plugin version is hardcoded in `cmd/protoc-gen-golemporal/main.go` as `var Version = "v0.3.0"`. The GitHub Actions release workflow (`.github/workflows/release.yml`) updates this value, commits the change, creates a git tag, and publishes a GitHub release.

### Running Examples

Requires a Temporal server on `localhost:7233`. Start the worker first, then the starter:

```bash
# Terminal 1
cd example/worker && go run main.go

# Terminal 2
cd example/starter && go run main.go
```
