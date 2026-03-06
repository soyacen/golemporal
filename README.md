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

package example;

option go_package = "example;example";

// Workflow service (must end with "Workflow")
service GreeterWorkflow {
  rpc Hello(HelloRequest) returns (HelloResponse);
}

// Activity service (must end with "Activity")
service GreeterActivity {
  rpc Greet(GreetRequest) returns (GreetResponse);
}

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}

message GreetRequest {
  string name = 1;
}

message GreetResponse {
  string message = 1;
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
    "log"
    "time"

    "example"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "go.temporal.io/sdk/workflow"
)

func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    w := worker.New(c, "my-task-queue", worker.Options{})

    // Register workflow and activities using generated functions
    example.RegisterGreeterWorkflow(w, &GreeterWorkflowServer{
        activity: example.NewGreeterActivityClient(),
    })
    example.RegisterGreeterActivity(w, &GreeterActivityServer{})

    if err := w.Start(); err != nil {
        log.Fatalln("Unable to start worker", err)
    }

    // Wait for interrupt
    <-make(chan struct{})
}

// GreeterWorkflowServer implements the generated GreeterWorkflowServer interface
type GreeterWorkflowServer struct {
    activity example.GreeterActivityClient
}

func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *example.HelloRequest) (*example.HelloResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("Starting HelloWorkflow", "name", input.Name)

    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)

    result, err := s.activity.Greet(ctx, &example.GreetRequest{Name: input.Name})
    if err != nil {
        logger.Error("Activity failed", "error", err)
        return nil, err
    }

    return &example.HelloResponse{Message: result.Message}, nil
}

// GreeterActivityServer implements the generated GreeterActivityServer interface
type GreeterActivityServer struct{}

func (s *GreeterActivityServer) Greet(ctx context.Context, input *example.GreetRequest) (*example.GreetResponse, error) {
    return &example.GreetResponse{Message: "Hello, " + input.Name + "!"}, nil
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

    "example"
    "github.com/soyacen/golemporal/starter"
    "go.temporal.io/sdk/client"
)

func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    // Create a workflow client
    gc := example.NewGreeterWorkflowClient(c, "my-task-queue")

    // Execute workflow with options
    var result example.HelloResponse
    err = gc.Hello(context.Background(), 
        &example.HelloRequest{Name: "World"},
        starter.ID("my-workflow-id"),
        starter.GetResult(&result),
    )
    if err != nil {
        log.Fatalln("Workflow failed", err)
    }

    log.Println("Result:", result.Message)
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

result, err := gc.Hello(ctx, &example.HelloRequest{Name: "World"},
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
    starter.SearchAttributes(map[string]any{"CustomKey": "value"}), // Search attributes
)
```

### Available Options

| Option | Description |
|--------|-------------|
| `ID(string)` | Set a unique workflow ID |
| `GetID(*string)` | Get the assigned workflow ID after execution |
| `GetRunID(*string)` | Get the run ID after execution |
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
| `GetResult(any)` | Capture workflow result into pointer |

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
