# Worker 注册参考

本参考文档说明如何使用 golemporal 生成的注册函数将 Workflow 和 Activity 注册到 Temporal Worker。

## 基本注册

### 注册单个 Workflow

```go
func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln(err)
    }
    defer c.Close()

    w := worker.New(c, "my-task-queue", worker.Options{})

    // 创建 Workflow Server 实例
    wf := &HelloWorkflowServer{
        greetActivity: api.NewGreetActivityClient(),
    }

    // 注册 Workflow
    api.RegisterHelloWorkflow(w, wf)

    // 注册 Activity
    api.RegisterGreetActivity(w, &GreetActivityServer{})

    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln(err)
    }
}
```

### 注册函数签名

```go
// Workflow 注册
func RegisterHelloWorkflow(wk worker.Worker, server HelloWorkflowServer)

// Activity 注册
func RegisterGreetActivity(wk worker.Worker, server GreetActivityServer)
```

## 一个 Struct 实现多个接口

这是 golemporal 的常见模式：一个 struct 实现多个 Workflow 或 Activity Server 接口。

### 多个 Workflow 共享依赖

```go
type GreeterWorkflowServer struct {
    addActivity   api.AddActivityClient
    multiActivity api.MultiActivityClient
}

// 实现 HelloWorkflowServer
func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)

    result, err := s.addActivity.Add(ctx, &api.AddRequest{Count: input.Count})
    if err != nil {
        return nil, err
    }

    return &api.HelloResponse{
        Message: fmt.Sprintf("Hello, %s!", input.Name),
        Result:  result.Result,
    }, nil
}

// 实现 GoodbyeWorkflowServer
func (s *GreeterWorkflowServer) Goodbye(ctx workflow.Context, input *api.GoodbyeRequest) (*api.GoodbyeResponse, error) {
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)

    result, err := s.multiActivity.Multi(ctx, &api.MultiRequest{Count: input.Count})
    if err != nil {
        return nil, err
    }

    return &api.GoodbyeResponse{
        Message: fmt.Sprintf("Goodbye, %s!", input.Name),
        Result:  result.Result,
    }, nil
}

func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln(err)
    }
    defer c.Close()

    w := worker.New(c, "golemporal-example", worker.Options{})

    // 同一个实例注册到多个 Workflow
    wf := &GreeterWorkflowServer{
        addActivity:   api.NewAddActivityClient(),
        multiActivity: api.NewMultiActivityClient(),
    }
    api.RegisterHelloWorkflow(w, wf)
    api.RegisterGoodbyeWorkflow(w, wf)

    // 注册 Activity
    api.RegisterAddActivity(w, &AddActivityServer{})
    api.RegisterMultiActivity(w, &MultiActivityServer{})

    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln(err)
    }
}
```

### Activity 也可以共享

```go
// 一个 struct 实现多个 Activity Server 接口
type OrderActivityServer struct {
    paymentService  PaymentService
    shippingService ShippingService
}

// 实现 PaymentActivityServer
func (s *OrderActivityServer) Charge(ctx context.Context, input *api.ChargeRequest) (*api.ChargeResponse, error) {
    return s.paymentService.Charge(ctx, input)
}

// 实现 ShippingActivityServer
func (s *OrderActivityServer) Ship(ctx context.Context, input *api.ShipRequest) (*api.ShipResponse, error) {
    return s.shippingService.Ship(ctx, input)
}

// 注册
api.RegisterPaymentActivity(w, &OrderActivityServer{...})
api.RegisterShippingActivity(w, &OrderActivityServer{...})
```

## 注册顺序

注册顺序不影响运行时行为，但建议按以下顺序组织代码以提高可读性：

```go
func main() {
    // ... 创建 client 和 worker ...

    // 1. 创建所有 Server 实例
    workflowServer := NewWorkflowServer()
    activityServer := NewActivityServer()

    // 2. 注册所有 Workflow
    api.RegisterHelloWorkflow(w, workflowServer)
    api.RegisterGoodbyeWorkflow(w, workflowServer)

    // 3. 注册所有 Activity
    api.RegisterGreetActivity(w, activityServer)
    api.RegisterPaymentActivity(w, activityServer)

    // 4. 启动 Worker
    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln(err)
    }
}
```

## 禁用已注册检查

生成的注册函数内部使用 `DisableAlreadyRegisteredCheck: true`：

```go
func RegisterHelloWorkflow(wk worker.Worker, server HelloWorkflowServer) {
    wk.RegisterWorkflowWithOptions(server.Hello, workflow.RegisterOptions{
        Name:                          HelloWorkflow_Hello_WorkflowType,
        DisableAlreadyRegisteredCheck: true,
    })
}
```

这意味着：
- **不会报错**：如果同一个 workflow/activity 名称被重复注册，不会报错
- **后注册者覆盖**：后注册的实现会覆盖先注册的实现
- **注意事项**：确保不要意外重复注册，否则可能导致运行时行为不符合预期

## 完整 Worker 示例

```go
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

type GreeterWorkflowServer struct {
    addActivity   api.AddActivityClient
    multiActivity api.MultiActivityClient
}

func (s *GreeterWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("HelloWorkflow starting")
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)
    result, err := s.addActivity.Add(ctx, &api.AddRequest{Count: input.Count})
    if err != nil {
        logger.Error("activity failed", "error", err)
        return nil, err
    }
    logger.Info("HelloWorkflow completed")
    return &api.HelloResponse{
        Message: fmt.Sprintf("Hello, %s! (result: %d)", input.Name, result.Result),
        Result:  result.Result,
    }, nil
}

func (s *GreeterWorkflowServer) Goodbye(ctx workflow.Context, input *api.GoodbyeRequest) (*api.GoodbyeResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("GoodbyeWorkflow starting")
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)
    result, err := s.multiActivity.Multi(ctx, &api.MultiRequest{Count: input.Count})
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

type AddActivityServer struct{}

func (s *AddActivityServer) Add(ctx context.Context, input *api.AddRequest) (*api.AddResponse, error) {
    return &api.AddResponse{Result: input.Count + input.Count}, nil
}

type MultiActivityServer struct{}

func (s *MultiActivityServer) Multi(ctx context.Context, input *api.MultiRequest) (*api.MultiResponse, error) {
    return &api.MultiResponse{Result: input.Count * input.Count}, nil
}

func main() {
    c, err := client.Dial(client.Options{HostPort: client.DefaultHostPort})
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    taskQueue := "golemporal-example"
    w := worker.New(c, taskQueue, worker.Options{})

    wf := &GreeterWorkflowServer{
        addActivity:   api.NewAddActivityClient(),
        multiActivity: api.NewMultiActivityClient(),
    }
    api.RegisterHelloWorkflow(w, wf)
    api.RegisterGoodbyeWorkflow(w, wf)
    api.RegisterAddActivity(w, &AddActivityServer{})
    api.RegisterMultiActivity(w, &MultiActivityServer{})

    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln("Unable to start worker", err)
    }
}
```
