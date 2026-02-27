# golemporal

A Temporal SDK framework for Go that uses protobuf-based code generation to define workflows and activities.

## Overview

golemporal provides a code generation tool (`protoc-gen-golemporal`) that generates type-safe Temporal workflow clients, activity clients, and worker registration code from proto service definitions.

## Features

- **Proto-based workflow definition**: Define workflows and activities in proto files
- **Type-safe code generation**: Generated code provides compile-time type checking
- **Functional options**: Workflow start options via functional configuration
- **Automatic registration**: Generated worker registration functions

## Quick Start

### 1. Define Proto Services

```protobuf
// example.proto
syntax = "proto3";

package example;

option go_package = "example;example";

// Workflow service (ends with "Workflow")
service GreeterWorkflow {
  rpc Hello(HelloRequest) returns (HelloResponse);
}

// Activity service (ends with "Activity")
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

```bash
# Install the code generator
go install github.com/soyacen/golemporal/cmd/protoc-gen-golemporal@latest

# Generate code from proto files
protoc --proto_path=. --go_out=. --go_opt=paths=source_relative \
  --golemporal_out=. --golemporal_opt=paths=source_relative \
  example.proto
```

### 3. Implement Workflow and Activities

```go
// worker.go
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

    // Register workflow and activities
    example.RegisterGreeterWorkflowWorker(w,
        &GreeterWorkflowServer{
            activity: example.NewGreeterActivityClient(),
        },
        &GreeterActivityServer{},
    )

    w.Start()
    <-make(chan struct{})
}

type GreeterWorkflowServer struct {
    activity example.GreeterActivityClient
}

func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *example.HelloRequest) (*example.HelloResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("Starting HelloWorkflow")

    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)

    result, err := s.activity.Greet(ctx, &example.GreetRequest{Name: input.Name})
    if err != nil {
        return nil, err
    }

    return &example.HelloResponse{Message: result.Message}, nil
}

type GreeterActivityServer struct{}

func (s *GreeterActivityServer) Greet(ctx context.Context, input *example.GreetRequest) (*example.GreetResponse, error) {
    return &example.GreetResponse{Message: "Hello, " + input.Name + "!"}, nil
}
```

### 4. Start Workflows

```go
// starter.go
package main

import (
    "context"
    "log"

    "example"
    "go.temporal.io/sdk/client"
)

func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    gc := example.NewGreeterWorkflowClient(c, "my-task-queue")

    // Execute workflow
    result, err := gc.Hello(context.Background(), &example.HelloRequest{Name: "World"})
    if err != nil {
        log.Fatalln("Workflow failed", err)
    }

    log.Println("Result:", result.Message)
}
```

## Proto Service Naming

The code generator uses naming conventions:

- Services ending with `Workflow` are workflow services
- Services ending with `Activity` are activity services
- Only one workflow service per proto file is supported

## Workflow Options

Use functional options when starting workflows:

```go
result, err := gc.Hello(ctx, &example.HelloRequest{
    Name: "World",
}, starter.ID("my-workflow-id"))
```

Available options in `starter` package:
- `ID(string)` - Set workflow ID
- `TaskQueue(string)` - Override task queue
- `WorkflowExecutionTimeout(time.Duration)` - Execution timeout
- `WorkflowRunTimeout(time.Duration)` - Run timeout
- `WorkflowTaskTimeout(time.Duration)` - Task timeout
- `RetryPolicy(*temporal.RetryPolicy)` - Retry policy
- `CronSchedule(string)` - Cron schedule
- `Memo(map[string]any)` - Workflow memo
- And more...

## Example

See the `example/` directory for a complete working example:

```bash
# Terminal 1: Start worker
cd example/worker && go run main.go

# Terminal 2: Start workflow
cd example/starter && go run main.go
```

Requires a Temporal server running on `localhost:7233`.

## Dependencies

- Go 1.25+
- Temporal Go SDK
- Protocol Buffers
- protoc
