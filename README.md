# golemporal

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)

A Temporal SDK framework for Go that uses protobuf-based code generation to define workflows and activities.

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
  - [1. Define Proto Services](#1-define-proto-services)
  - [2. Generate Code](#2-generate-code)
  - [3. Implement Workflow and Activities](#3-implement-workflow-and-activities)
  - [4. Start Workflows](#4-start-workflows)
- [Generated Code Structure](#generated-code-structure)
- [Proto Service Naming](#proto-service-naming)
- [Workflow Options](#workflow-options)
- [Example](#example)
- [Project Structure](#project-structure)
- [Dependencies](#dependencies)
- [License](#license)

## Overview

golemporal provides a protoc plugin (`protoc-gen-golemporal`) that generates type-safe Temporal workflow clients, activity clients, and worker registration code from proto service definitions. It enables a **proto-first** approach to building Temporal workflows in Go.

## Features

- **📝 Proto-based workflow definition**: Define workflows and activities in proto files
- **🔒 Type-safe code generation**: Generated code provides compile-time type checking
- **⚙️ Functional options**: Workflow start options via functional configuration
- **🚀 Automatic registration**: Generated worker registration functions
- **🎯 Clean architecture**: Separation of workflow and activity concerns

## Installation

### Prerequisites

- Go 1.25 or later
- Protocol Buffers compiler (`protoc`)
- Go protoc plugins (`protoc-gen-go`)

### Install protoc-gen-golemporal

```bash
go install github.com/soyacen/golemporal/cmd/protoc-gen-golemporal@latest
```

### Install protoc plugins (if not already installed)

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
```

## Quick Start

### 1. Define Proto Services

Create a `.proto` file with your workflow and activity definitions:

```protobuf
// example.proto
syntax = "proto3";

package golemporal.example.api;

option go_package = "github.com/soyacen/golemporal/example/api;api";

// Workflow services (must end with "Workflow")
service HelloWorkflow {
  rpc Hello(HelloRequest) returns (HelloResponse);
}

service GoodbyeWorkflow {
  rpc Goodbye(GoodbyeRequest) returns (GoodbyeResponse);
}

// Activity services (must end with "Activity")
service AddActivity {
  rpc Add(AddRequest) returns (AddResponse);
}

service MultiActivity {
  rpc Multi(MultiRequest) returns (MultiResponse);
}

message HelloRequest {
  string name = 1;
  int32 count = 2;
}

message HelloResponse {
  string message = 1;
  int32 result = 2;
}

message GoodbyeRequest {
  string name = 1;
  int32 count = 2;
}

message GoodbyeResponse {
  string message = 1;
  int32 result = 2;
}

message AddRequest {
  int32 count = 1;
}

message AddResponse {
  int32 result = 1;
}

message MultiRequest {
  int32 count = 1;
}

message MultiResponse {
  int32 result = 1;
}
```

### 2. Generate Code

Run the protoc compiler to generate Go code:

```bash
protoc \
  --proto_path=. \
  --go_out=. \
  --go_opt=paths=source_relative \
  --golemporal_out=. \
  --golemporal_opt=paths=source_relative \
  example.proto
```

This generates:
- `example.pb.go` - Standard protobuf Go code
- `example_temporal.pb.go` - Temporal-specific generated code (clients, servers, registration)

### 3. Implement Workflow and Activities

Create a worker that implements the generated interfaces:

```go
// worker/main.go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/soyacen/golemporal/example/api"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "go.temporal.io/sdk/workflow"
)

func main() {
    // The client is a heavyweight object that should be created once per process.
    c, err := client.Dial(client.Options{
        HostPort: client.DefaultHostPort,
    })
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    // Create worker
    taskQueue := "golemporal-example"
    w := worker.New(c, taskQueue, worker.Options{})
    wf := &GreeterWorkflowServer{
        addActivity:   api.NewAddActivityClient(),
        multiActivity: api.NewMultiActivityClient(),
    }
    
    // Register workflows and activities using generated functions
    api.RegisterHelloWorkflow(w, wf)
    api.RegisterGoodbyeWorkflow(w, wf)
    api.RegisterAddActivity(w, &AddActivityServer{})
    api.RegisterMultiActivity(w, &MultiActivityServer{})

    if err := w.Start(); err != nil {
        log.Fatalln("Unable to start worker", err)
    }

    // Wait for interrupt
    <-make(chan struct{})
}

// GreeterWorkflowServer implements the workflow interfaces
type GreeterWorkflowServer struct {
    addActivity   api.AddActivityClient
    multiActivity api.MultiActivityClient
}

func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("HelloWorkflow starting")
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, ao)
    helloResult, err := s.addActivity.Add(ctx, &api.AddRequest{Count: input.GetCount()})
    if err != nil {
        logger.Error("activity failed", "error", err)
        return nil, err
    }
    logger.Info("HelloWorkflow completed")
    return &api.HelloResponse{
        Message: fmt.Sprintf("Hello, %s! (result: %d)", input.Name, helloResult.Result),
        Result:  helloResult.Result,
    }, nil
}

func (s *GreeterWorkflowServer) Goodbye(ctx workflow.Context, input *api.GoodbyeRequest) (*api.GoodbyeResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("GoodbyeWorkflow starting")
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, ao)
    result, err := s.multiActivity.Multi(ctx, &api.MultiRequest{Count: input.GetCount()})
    if err != nil {
        logger.Error("activity failed", "error", err)
        return nil, err
    }
    logger.Info("GoodbyeWorkflow completed")
    return &api.GoodbyeResponse{
        Message: fmt.Sprintf("Goodbye, %s! (result: %d)", input.Name, result.Result),
        Result:  result.Result,
    }, nil
}

// Activity implementations
type AddActivityServer struct{}

func (s *AddActivityServer) Add(ctx context.Context, input *api.AddRequest) (*api.AddResponse, error) {
    return &api.AddResponse{Result: input.Count + input.Count}, nil
}

type MultiActivityServer struct{}

func (s *MultiActivityServer) Multi(ctx context.Context, input *api.MultiRequest) (*api.MultiResponse, error) {
    return &api.MultiResponse{Result: input.Count * input.Count}, nil
}
```

### 4. Start Workflows

Use the generated client to start workflow executions:

```go
// starter/main.go
package main

import (
    "context"
    "log"

    "github.com/soyacen/golemporal/example/api"
    "github.com/soyacen/golemporal/starter"
    "go.temporal.io/sdk/client"
)

func main() {
    ctx := context.Background()
    c, err := client.Dial(client.Options{
        HostPort: client.DefaultHostPort,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer c.Close()

    taskQueue := "golemporal-example"

    // Create a workflow client and execute workflow
    hc := api.NewHelloWorkflowClient(c, taskQueue)
    var helloResult api.HelloResponse
    hlMd, err := hc.Hello(ctx, &api.HelloRequest{Name: "World", Count: 5}, starter.GetResult(&helloResult))
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("hello workflow completed, workflow_id: %s, workflow_type: %s, run_id: %s, message: %s, result: %d",
        hlMd.GetWorkflowId(), hlMd.GetWorkflowType(), hlMd.GetRunId(), helloResult.Message, helloResult.Result)

    gc := api.NewGoodbyeWorkflowClient(c, taskQueue)
    var goodbyeResult api.GoodbyeResponse
    bgMd, err := gc.Goodbye(ctx, &api.GoodbyeRequest{Name: "World", Count: 10}, starter.GetResult(&goodbyeResult))
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("goodbye workflow completed, workflow_id: %s, workflow_type: %s, run_id: %s, message: %s, result: %d",
        bgMd.GetWorkflowId(), bgMd.GetWorkflowType(), bgMd.GetRunId(), goodbyeResult.Message, goodbyeResult.Result)
}
```

## Generated Code Structure

The `*_temporal.pb.go` file contains the following generated code:

| Component | Description |
|-----------|-------------|
| `ActivityClient` | Interface for invoking activities from workflows |
| `ActivityServer` | Interface that activity implementations must satisfy |
| `WorkflowClient` | Interface for starting workflow executions |
| `WorkflowServer` | Interface that workflow implementations must satisfy |
| `Register*Activity` | Function to register activity with a Temporal worker |
| `Register*Workflow` | Function to register workflow with a Temporal worker |

### Workflow Execution Metadata

Workflow client methods return `(*protobuf.Metadata, error)` which contains execution information:

```go
type Metadata struct {
    WorkflowId   string // Unique workflow ID
    RunId        string // Run ID for this execution
    WorkflowType string // Workflow type name
}
```

Example usage:

```go
md, err := client.Hello(ctx, &api.HelloRequest{Name: "World"})
if err != nil {
    log.Fatal(err)
}
log.Printf("Workflow started: ID=%s, RunID=%s, Type=%s", 
    md.GetWorkflowId(), md.GetRunId(), md.GetWorkflowType())
```

## Proto Service Naming

The code generator recognizes services by their naming convention:

- **Workflow Services**: Must end with `Workflow` (e.g., `GreeterWorkflow`)
- **Activity Services**: Must end with `Activity` (e.g., `GreeterActivity`)

**Note:**
- Service names are case-sensitive

## Workflow Options

Use functional options from the `starter` package when starting workflows:

```go
import "github.com/soyacen/golemporal/starter"

md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.ID("my-workflow-id"),                           // Set workflow ID
    starter.TaskQueue("custom-queue"),                      // Override task queue
    starter.WorkflowExecutionTimeout(30 * time.Minute),     // Execution timeout
    starter.WorkflowRunTimeout(10 * time.Minute),           // Run timeout
    starter.WorkflowTaskTimeout(10 * time.Second),          // Task timeout
    starter.RetryPolicy(&temporal.RetryPolicy{               // Retry policy
        MaximumAttempts: 3,
    }),
    starter.CronSchedule("0 9 * * *"),                      // Cron schedule
    starter.Memo(map[string]any{"key": "value"}),           // Workflow memo
    starter.TypedSearchAttributes(temporal.NewSearchAttributes(...)), // Search attributes
)
```

### Available Options

| Option | Description |
|--------|-------------|
| `ID(string)` | Set a unique workflow ID |
| `TaskQueue(string)` | Override the default task queue |
| `WorkflowExecutionTimeout(time.Duration)` | Total workflow execution timeout |
| `WorkflowRunTimeout(time.Duration)` | Single workflow run timeout |
| `WorkflowTaskTimeout(time.Duration)` | Workflow task timeout |
| `WorkflowIDReusePolicy(enums.WorkflowIdReusePolicy)` | Workflow ID reuse policy |
| `WorkflowIDConflictPolicy(enums.WorkflowIdConflictPolicy)` | Workflow ID conflict policy |
| `WorkflowExecutionErrorWhenAlreadyStarted(bool)` | Return error when workflow already running |
| `RetryPolicy(*temporal.RetryPolicy)` | Retry policy for workflow |
| `CronSchedule(string)` | Cron schedule for periodic execution |
| `Memo(map[string]any)` | Workflow memo data |
| `TypedSearchAttributes(temporal.SearchAttributes)` | Typed search attributes for workflow |
| `EnableEagerStart(bool)` | Request eager execution if local worker available |
| `StartDelay(time.Duration)` | Delay before dispatching first workflow task |
| `StaticSummary(string)` | Static summary for the workflow |
| `StaticDetails(string)` | Static details for the workflow |
| `VersioningOverride(client.VersioningOverride)` | Versioning override for the workflow |
| `Priority(temporal.Priority)` | Priority for the workflow |
| `GetResult(proto.Message)` | Capture workflow result into a proto message pointer |

## Example

See the `example/` directory for a complete working example.

### Running the Example

1. **Start a Temporal server** (requires Temporal running on `localhost:7233`):
   ```bash
   # Using Temporal CLI
   temporal server start-dev
   
   # Or using Docker
   docker run --rm -p 7233:7233 temporalio/server:latest
   ```

2. **Start the worker** (Terminal 1):
   ```bash
   cd example/worker && go run main.go
   ```

3. **Start the workflow** (Terminal 2):
   ```bash
   cd example/starter && go run main.go
   ```

## Project Structure

```
.
├── cmd/protoc-gen-golemporal/   # Protoc plugin implementation
├── starter/                     # Workflow starter options (functional options)
├── protobuf/                    # Common protobuf definitions (Metadata)
├── example/                     # Complete working example
│   ├── api/                     # Proto definitions and generated code
│   ├── starter/                 # Workflow client example
│   └── worker/                  # Worker implementation example
├── skills/temporal-go-sdk/      # Temporal SDK reference documentation
├── go.mod                       # Go module definition
└── doc.go                       # Package documentation
```

## Dependencies

- **Go**: 1.25+
- [Temporal Go SDK](https://github.com/temporalio/sdk-go): Core Temporal SDK
- [Temporal API](https://github.com/temporalio/api-go): Temporal API definitions
- [Protocol Buffers](https://github.com/protocolbuffers/protobuf): Protocol Buffers compiler
- [Go Protobuf](https://google.golang.org/protobuf): Go protobuf runtime

## Contributing

Contributions are welcome! Please feel free to submit issues or pull requests.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

---

For more information, see the [Temporal Documentation](https://docs.temporal.io/) and [Temporal Go SDK](https://github.com/temporalio/sdk-go).
