---
name: golemporal-developer
description: golemporal 框架开发专家。精通基于 protobuf 代码生成的 Temporal Go 开发模式，能够指导用户完成从 proto 定义、代码生成到 workflow/activity 实现和 worker 注册的完整流程。当用户提及 golemporal、protoc-gen-golemporal、proto workflow、proto activity、temporal 代码生成等关键词时必须使用此 skill。
compatibility: Go 1.25+, go.temporal.io/sdk, google.golang.org/protobuf
---

# golemporal Developer Skill

本 skill 提供 golemporal 框架的完整开发指导。golemporal 是一个基于 protobuf 代码生成的 Temporal SDK 框架，通过 protoc 插件从 proto service 定义自动生成类型安全的客户端、服务器和注册代码。

## 核心概念快速导航

| 主题 | 参考文档 | 说明 |
|------|----------|------|
| Proto 定义规范 | [proto-definition](./references/proto-definition.md) | Service 命名约定、message 定义 |
| 代码生成 | [code-generation](./references/code-generation.md) | protoc 命令、插件安装 |
| Workflow 实现 | [workflow-implementation](./references/workflow-implementation.md) | Server 接口、Activity 调用 |
| Activity 实现 | [workflow-implementation](./references/workflow-implementation.md) | Server 接口、context.Context |
| Worker 注册 | [worker-registration](./references/worker-registration.md) | 注册函数、多接口实现 |
| 启动选项 | [starter-options](./references/starter-options.md) | 所有 starter.Option 详解 |

## 目录

1. [快速开始](#快速开始)
2. [Proto 定义](#proto-定义)
3. [代码生成](#代码生成)
4. [Workflow 实现](#workflow-实现)
5. [Activity 实现](#activity-实现)
6. [Worker 注册](#worker-注册)
7. [启动工作流](#启动工作流)
8. [常见模式](#常见模式)
9. [注意事项](#注意事项)

---

## 快速开始

### 完整流程概览

```
定义 .proto → 运行 protoc 生成代码 → 实现 Workflow/Activity Server → 注册 Worker → 启动工作流
```

### 1. 定义 Proto 文件

```protobuf
syntax = "proto3";
package myapp.api;
option go_package = "github.com/myorg/myapp/api;api";

message HelloRequest {
  string name = 1;
}

message HelloResponse {
  string message = 1;
}

service HelloWorkflow {
  rpc Hello(HelloRequest) returns (HelloResponse);
}

service GreetActivity {
  rpc Greet(HelloRequest) returns (HelloResponse);
}
```

### 2. 生成代码

```bash
cd example && ./protoc.sh
```

生成 `api/hello_temporal.pb.go`，包含 `HelloWorkflowClient`、`HelloWorkflowServer`、`GreetActivityClient`、`GreetActivityServer` 和注册函数。

### 3. 实现 Workflow

```go
type HelloWorkflowServer struct {
    greetActivity api.GreetActivityClient
}

func (s *HelloWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)
    return s.greetActivity.Greet(ctx, input)
}
```

### 4. 实现 Activity

```go
type GreetActivityServer struct{}

func (s *GreetActivityServer) Greet(ctx context.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    return &api.HelloResponse{Message: "Hello, " + input.Name + "!"}, nil
}
```

### 5. 注册并运行 Worker

```go
func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln(err)
    }
    defer c.Close()

    w := worker.New(c, "my-task-queue", worker.Options{})

    wf := &HelloWorkflowServer{greetActivity: api.NewGreetActivityClient()}
    api.RegisterHelloWorkflow(w, wf)
    api.RegisterGreetActivity(w, &GreetActivityServer{})

    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln(err)
    }
}
```

### 6. 启动工作流

```go
hc := api.NewHelloWorkflowClient(c, "my-task-queue")
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.WaitResult(true))
if err != nil {
    log.Fatal(err)
}
log.Printf("result: %s, workflow_id: %s", result.Message, md.GetWorkflowId())
```

---

## Proto 定义

详细参考：[proto-definition](./references/proto-definition.md)

### 命名约定（必须遵守）

- **Workflow Service**：必须以 `Workflow` 结尾，如 `HelloWorkflow`、`OrderWorkflow`
- **Activity Service**：必须以 `Activity` 结尾，如 `GreetActivity`、`PaymentActivity`
- 区分大小写，必须是精确后缀匹配
- 一个 proto 文件可以包含 **多个** workflow 和 activity 服务

### 示例：多服务 proto

```protobuf
service HelloWorkflow {
  rpc Hello(HelloRequest) returns (HelloResponse);
}

service GoodbyeWorkflow {
  rpc Goodbye(GoodbyeRequest) returns (GoodbyeResponse);
}

service AddActivity {
  rpc Add(AddRequest) returns (AddResponse);
}

service MultiActivity {
  rpc Multi(MultiRequest) returns (MultiResponse);
}
```

---

## 代码生成

详细参考：[code-generation](./references/code-generation.md)

### 依赖安装

```bash
# 安装标准 Go protobuf 插件
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest

# 安装 golemporal 插件
go install github.com/soyacen/golemporal/cmd/protoc-gen-golemporal@latest
```

### protoc 命令

```bash
protoc \
  --proto_path=. \
  --proto_path=../../ \
  --go_out=. \
  --go_opt=paths=source_relative \
  --golemporal_out=. \
  --golemporal_opt=paths=source_relative \
  *.proto
```

### 生成产物

对于包含 `HelloWorkflow` 和 `GreetActivity` 的 `hello.proto`：

| 生成项 | 名称 | 用途 |
|--------|------|------|
| Activity Client | `GreetActivityClient` | workflow 内部调用 activity |
| Activity Server | `GreetActivityServer` | 实现 activity 逻辑 |
| Workflow Client | `HelloWorkflowClient` | 应用代码启动工作流 |
| Workflow Server | `HelloWorkflowServer` | 实现 workflow 逻辑 |
| 注册函数 | `RegisterGreetActivity` | 注册 activity 到 worker |
| 注册函数 | `RegisterHelloWorkflow` | 注册 workflow 到 worker |

---

## Workflow 实现

详细参考：[workflow-implementation](./references/workflow-implementation.md)

### Server 接口

```go
type HelloWorkflowServer interface {
    Hello(workflow.Context, *HelloRequest) (*HelloResponse, error)
}
```

实现要点：
- 方法接收 `workflow.Context`（不是 `context.Context`）
- 通过生成的 **Activity Client** 调用 activity，而非直接使用 `workflow.ExecuteActivity`
- 需要设置 `workflow.ActivityOptions`

### 在 Workflow 中调用 Activity

```go
func (s *HelloWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 使用生成的 Activity Client
    result, err := s.greetActivity.Greet(ctx, input)
    if err != nil {
        workflow.GetLogger(ctx).Error("activity failed", "error", err)
        return nil, err
    }
    return result, nil
}
```

---

## Activity 实现

详细参考：[workflow-implementation](./references/workflow-implementation.md)

### Server 接口

```go
type GreetActivityServer interface {
    Greet(context.Context, *HelloRequest) (*HelloResponse, error)
}
```

实现要点：
- 方法接收 `context.Context`（不是 `workflow.Context`）
- 可以直接调用外部服务、数据库等
- 实现应为纯函数，便于测试

```go
type GreetActivityServer struct{}

func (s *GreetActivityServer) Greet(ctx context.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    return &api.HelloResponse{Message: "Hello, " + input.Name + "!"}, nil
}
```

---

## Worker 注册

详细参考：[worker-registration](./references/worker-registration.md)

### 基本注册

```go
api.RegisterHelloWorkflow(w, &HelloWorkflowServer{
    greetActivity: api.NewGreetActivityClient(),
})
api.RegisterGreetActivity(w, &GreetActivityServer{})
```

### 一个 struct 实现多个接口

```go
type GreeterWorkflowServer struct {
    addActivity   api.AddActivityClient
    multiActivity api.MultiActivityClient
}

// 实现 HelloWorkflowServer
func (s *GreeterWorkflowServer) Hello(...) (...) { ... }

// 实现 GoodbyeWorkflowServer
func (s *GreeterWorkflowServer) Goodbye(...) (...) { ... }

// 注册时传入同一个实例
wf := &GreeterWorkflowServer{...}
api.RegisterHelloWorkflow(w, wf)
api.RegisterGoodbyeWorkflow(w, wf)
```

---

## 启动工作流

详细参考：[starter-options](./references/starter-options.md)

### 基础用法

```go
hc := api.NewHelloWorkflowClient(c, "my-task-queue")

// 异步启动（不等待结果）
md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"})
// 返回值: result=nil, md=执行元数据, err=错误

// 同步等待结果
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.WaitResult(true))
// 返回值: result=工作流输出, md=执行元数据, err=错误
```

### 返回值说明

Workflow Client 方法返回 `(*Output, *protobuf.Metadata, error)`：

- `Output`：工作流的返回结果（proto message 指针），仅当 `WaitResult=true` 时有值
- `Metadata`：执行元数据，包含 `WorkflowId`、`RunId`、`WorkflowType`、`TaskQueue`
- `error`：启动或执行过程中的错误

### 常用启动选项

```go
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.ID("my-workflow-id"),                           // 指定工作流 ID
    starter.WorkflowExecutionTimeout(30 * time.Minute),     // 执行超时
    starter.RetryPolicy(&temporal.RetryPolicy{MaximumAttempts: 3}),
    starter.CronSchedule("0 9 * * *"),                      // 定时调度
    starter.WaitResult(true),                               // 等待结果
)
```

---

## 常见模式

### 模式 1：一个 Workflow 调用多个 Activity

```go
type OrderWorkflowServer struct {
    paymentActivity api.PaymentActivityClient
    shipActivity    api.ShipActivityClient
}

func (s *OrderWorkflowServer) Process(ctx workflow.Context, input *api.OrderRequest) (*api.OrderResponse, error) {
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 调用 payment activity
    paymentResult, err := s.paymentActivity.Charge(ctx, &api.ChargeRequest{Amount: input.Amount})
    if err != nil {
        return nil, err
    }

    // 调用 ship activity
    shipResult, err := s.shippingActivity.Ship(ctx, &api.ShipRequest{OrderId: input.OrderId})
    if err != nil {
        return nil, err
    }

    return &api.OrderResponse{
        PaymentId: paymentResult.Id,
        ShipId:    shipResult.Id,
    }, nil
}
```

### 模式 2：多个 Workflow 共享 Activity

```go
// 同一个 struct 实现多个 Workflow Server 接口
type GreeterWorkflowServer struct {
    addActivity api.AddActivityClient
}

func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    result, _ := s.addActivity.Add(ctx, &api.AddRequest{Count: input.Count})
    return &api.HelloResponse{Result: result.Result}, nil
}

func (s *GreeterWorkflowServer) Goodbye(ctx workflow.Context, input *api.GoodbyeRequest) (*api.GoodbyeResponse, error) {
    result, _ := s.addActivity.Add(ctx, &api.AddRequest{Count: input.Count})
    return &api.GoodbyeResponse{Result: result.Result}, nil
}

// Worker 注册
wf := &GreeterWorkflowServer{addActivity: api.NewAddActivityClient()}
api.RegisterHelloWorkflow(w, wf)
api.RegisterGoodbyeWorkflow(w, wf)
api.RegisterAddActivity(w, &AddActivityServer{})
```

---

## 注意事项

### 1. DisableAlreadyRegisteredCheck

生成的注册函数使用 `DisableAlreadyRegisteredCheck: true`。如果同一个 worker 上重复注册同名 workflow 或 activity，不会报错，但后注册的会覆盖先注册的。确保注册逻辑正确。

### 2. 版本管理

插件版本硬编码在 `cmd/protoc-gen-golemporal/main.go` 中（如 `var Version = "v0.3.0"`）。通过 GitHub Actions release workflow 更新版本、打 tag 并发布 release。

### 3. Context 类型区分

| 场景 | Context 类型 |
|------|-------------|
| Workflow Server 方法 | `workflow.Context` |
| Activity Server 方法 | `context.Context` |
| Workflow Client 调用（启动工作流）| `context.Context` |

### 4. 重新生成代码的时机

- 修改 `.proto` 文件（新增/修改 service、message、rpc）
- 更新 golemporal 插件版本（生成代码格式可能变化）

### 5. sdk-go/ 目录说明

代码库中的 `sdk-go/` 是 temporalio/sdk-go 的 vendored 副本（独立 Go 模块），一般情况下不需要修改。

---

## 相关资源

- [golemporal GitHub 仓库](https://github.com/soyacen/golemporal)
- [Temporal Go SDK 文档](https://docs.temporal.io/dev-guide/golang)
- [Temporal Go SDK Skill](../temporal-go-sdk/SKILL.md) — 通用 Temporal Go 开发知识
