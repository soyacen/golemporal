# Workflow 与 Activity 实现参考

本参考文档说明如何实现 golemporal 生成的 Workflow Server 和 Activity Server 接口。

## Workflow 实现

### 基本结构

```go
type HelloWorkflowServer struct {
    greetActivity api.GreetActivityClient
}

func NewHelloWorkflowServer() *HelloWorkflowServer {
    return &HelloWorkflowServer{
        greetActivity: api.NewGreetActivityClient(),
    }
}

func (s *HelloWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    // 1. 设置 Activity 选项
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 2. 调用 Activity（使用生成的客户端）
    result, err := s.greetActivity.Greet(ctx, input)
    if err != nil {
        return nil, err
    }

    // 3. 返回结果
    return result, nil
}
```

### 关键要点

| 要点 | 说明 |
|------|------|
| Context 类型 | 必须使用 `workflow.Context`，不能用 `context.Context` |
| Activity 调用 | 使用生成的 `*ActivityClient`，不要直接调用 `workflow.ExecuteActivity` |
| 选项设置 | 通过 `workflow.WithActivityOptions` 设置超时、重试等 |
| 返回值 | 返回 `(*Response, error)`，与 proto 定义的返回类型对应 |

### 在 Workflow 中调用多个 Activity

```go
func (s *OrderWorkflowServer) Process(ctx workflow.Context, input *api.OrderRequest) (*api.OrderResponse, error) {
    ao := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 顺序调用多个 Activity
    paymentResult, err := s.paymentActivity.Charge(ctx, &api.ChargeRequest{
        OrderId: input.OrderId,
        Amount:  input.Amount,
    })
    if err != nil {
        return nil, fmt.Errorf("payment failed: %w", err)
    }

    shipResult, err := s.shippingActivity.Ship(ctx, &api.ShipRequest{
        OrderId: input.OrderId,
        Address: input.Address,
    })
    if err != nil {
        return nil, fmt.Errorf("shipping failed: %w", err)
    }

    return &api.OrderResponse{
        PaymentId: paymentResult.TransactionId,
        TrackingId: shipResult.TrackingId,
    }, nil
}
```

### 并行调用 Activity

```go
func (s *HelloWorkflowServer) ParallelHello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    ao := workflow.ActivityOptions{StartToCloseTimeout: 10 * time.Second}
    ctx = workflow.WithActivityOptions(ctx, ao)

    // 并行启动两个 Activity
    future1 := workflow.ExecuteActivity(ctx, api.GreetActivity_Greet_ActitityType, &api.HelloRequest{Name: input.Name + "-1"})
    future2 := workflow.ExecuteActivity(ctx, api.GreetActivity_Greet_ActitityType, &api.HelloRequest{Name: input.Name + "-2"})

    var result1, result2 api.HelloResponse
    if err := future1.Get(ctx, &result1); err != nil {
        return nil, err
    }
    if err := future2.Get(ctx, &result2); err != nil {
        return nil, err
    }

    return &api.HelloResponse{
        Message: result1.Message + " | " + result2.Message,
    }, nil
}
```

### 使用 Logger

```go
func (s *HelloWorkflowServer) Hello(ctx workflow.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    logger := workflow.GetLogger(ctx)
    logger.Info("HelloWorkflow started", "name", input.Name)

    result, err := s.greetActivity.Greet(ctx, input)
    if err != nil {
        logger.Error("greet activity failed", "error", err)
        return nil, err
    }

    logger.Info("HelloWorkflow completed", "result", result.Message)
    return result, nil
}
```

### 设置不同的 Activity 选项

可以为不同的 Activity 设置不同的超时和重试策略：

```go
func (s *OrderWorkflowServer) Process(ctx workflow.Context, input *api.OrderRequest) (*api.OrderResponse, error) {
    // Payment Activity：需要重试
    paymentCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: 30 * time.Second,
        RetryPolicy: &temporal.RetryPolicy{
            InitialInterval:    time.Second,
            BackoffCoefficient: 2.0,
            MaximumAttempts:    5,
        },
    })

    // Shipping Activity：不需要重试
    shipCtx := workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Minute,
    })

    paymentResult, err := s.paymentActivity.Charge(paymentCtx, &api.ChargeRequest{...})
    if err != nil {
        return nil, err
    }

    shipResult, err := s.shippingActivity.Ship(shipCtx, &api.ShipRequest{...})
    if err != nil {
        return nil, err
    }

    return &api.OrderResponse{...}, nil
}
```

## Activity 实现

### 基本结构

```go
type GreetActivityServer struct{}

func (s *GreetActivityServer) Greet(ctx context.Context, input *api.HelloRequest) (*api.HelloResponse, error) {
    return &api.HelloResponse{Message: "Hello, " + input.Name + "!"}, nil
}
```

### 关键要点

| 要点 | 说明 |
|------|------|
| Context 类型 | 必须使用 `context.Context`，不能用 `workflow.Context` |
| 外部调用 | 可以调用 HTTP API、数据库、文件系统等外部服务 |
| 幂等性 | Activity 可能重试，实现必须是幂等的 |
| 返回值 | 返回 `(*Response, error)`，与 proto 定义的返回类型对应 |

### 访问 Activity 信息

```go
import "go.temporal.io/sdk/activity"

func (s *PaymentActivityServer) Charge(ctx context.Context, input *api.ChargeRequest) (*api.ChargeResponse, error) {
    // 获取 Activity 信息
    info := activity.GetInfo(ctx)
    logger := activity.GetLogger(ctx)

    logger.Info("Processing payment",
        "activity_id", info.ActivityID,
        "activity_type", info.ActivityType,
        "attempt", info.Attempt,
    )

    // 执行实际的支付操作
    transactionId, err := processPayment(input.Amount, input.Currency)
    if err != nil {
        return nil, err
    }

    return &api.ChargeResponse{TransactionId: transactionId, Status: "success"}, nil
}
```

### 长运行 Activity（心跳）

```go
func (s *DataProcessingActivityServer) Process(ctx context.Context, input *api.ProcessRequest) (*api.ProcessResponse, error) {
    total := len(input.Items)

    for i, item := range input.Items {
        // 处理每个项目
        if err := processItem(item); err != nil {
            return nil, err
        }

        // 发送心跳（长运行 Activity 必须心跳）
        activity.RecordHeartbeat(ctx, i+1)

        // 检查是否被取消
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }
    }

    return &api.ProcessResponse{ProcessedCount: total}, nil
}
```

### 依赖注入

Activity Server 可以包含依赖，便于测试和复用：

```go
type PaymentActivityServer struct {
    paymentGateway PaymentGateway
    db             *sql.DB
}

func NewPaymentActivityServer(gateway PaymentGateway, db *sql.DB) *PaymentActivityServer {
    return &PaymentActivityServer{paymentGateway: gateway, db: db}
}

func (s *PaymentActivityServer) Charge(ctx context.Context, input *api.ChargeRequest) (*api.ChargeResponse, error) {
    // 使用注入的依赖
    return s.paymentGateway.Charge(ctx, input)
}
```

## Context 类型速查

| 位置 | Context 类型 | 包 |
|------|-------------|-----|
| Workflow Server 方法 | `workflow.Context` | `go.temporal.io/sdk/workflow` |
| Activity Server 方法 | `context.Context` | 标准库 `context` |
| Workflow Client 启动方法 | `context.Context` | 标准库 `context` |
| Activity Client 调用（在 workflow 内）| `workflow.Context` | `go.temporal.io/sdk/workflow` |
