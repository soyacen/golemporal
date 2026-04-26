# 启动选项参考

本参考文档说明 `starter` 包提供的所有 Workflow 启动选项及其用法。

## 基础用法

```go
import "github.com/soyacen/golemporal/starter"

hc := api.NewHelloWorkflowClient(c, "my-task-queue")

// 异步启动（不等待结果）
_, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"})

// 同步等待结果
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.WaitResult(true))
```

## 返回值说明

Workflow Client 方法返回三个值：`(*Output, *protobuf.Metadata, error)`

```go
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.WaitResult(true))

// result: *api.HelloResponse — 工作流输出（WaitResult=false 时为 nil）
// md: *protobuf.Metadata — 执行元数据
// err: error — 启动或执行错误
```

### Metadata 结构

```protobuf
message Metadata {
  string task_queue = 1;
  string workflow_id = 2;
  string run_id = 3;
  string workflow_type = 4;
}
```

```go
log.Printf("workflow_id: %s, run_id: %s, type: %s, queue: %s",
    md.GetWorkflowId(),
    md.GetRunId(),
    md.GetWorkflowType(),
    md.GetTaskQueue(),
)
```

## 所有选项

### WaitResult

控制是否阻塞等待工作流完成。

```go
// false（默认）：启动后立即返回，result 为 nil
_, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"})

// true：阻塞等待工作流完成，result 包含输出
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.WaitResult(true))
```

### ID

设置工作流 ID。如果不指定，Temporal 会自动生成 UUID。

```go
hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.ID("hello-world-001"))
```

### TaskQueue

覆盖默认的任务队列（在 `New*WorkflowClient` 中指定的队列）。

```go
hc.Hello(ctx, &api.HelloRequest{Name: "World"}, starter.TaskQueue("priority-queue"))
```

### 超时选项

```go
// 整个工作流执行的总超时（包括所有重试、Cron 触发等）
starter.WorkflowExecutionTimeout(30 * time.Minute)

// 单次工作流运行的超时
starter.WorkflowRunTimeout(10 * time.Minute)

// 工作流任务处理的超时
starter.WorkflowTaskTimeout(10 * time.Second)
```

### ID 复用策略

```go
import enums "go.temporal.io/api/enums/v1"

// 允许复用已完成的 workflow ID
starter.WorkflowIDReusePolicy(enums.WORKFLOW_ID_REUSE_POLICY_ALLOW_DUPLICATE)

// 拒绝复用
starter.WorkflowIDReusePolicy(enums.WORKFLOW_ID_REUSE_POLICY_REJECT_DUPLICATE)
```

### ID 冲突策略

```go
// 当 workflow ID 已存在时的处理策略
starter.WorkflowIDConflictPolicy(enums.WORKFLOW_ID_CONFLICT_POLICY_UNSPECIFIED)
```

### 已启动时返回错误

```go
// 当 workflow 已经在运行时返回错误
starter.WorkflowExecutionErrorWhenAlreadyStarted(true)
```

### 重试策略

```go
starter.RetryPolicy(&temporal.RetryPolicy{
    InitialInterval:        time.Second,
    BackoffCoefficient:     2.0,
    MaximumInterval:        100 * time.Second,
    MaximumAttempts:        5,
    NonRetryableErrorTypes: []string{"InvalidArgument"},
})
```

### Cron 调度

```go
// 标准 Cron 表达式
starter.CronSchedule("0 9 * * *")   // 每天 9:00
starter.CronSchedule("@hourly")     // 每小时
starter.CronSchedule("@daily")      // 每天

// 也可以使用 Temporal 的间隔语法
starter.CronSchedule("*/5 * * * *") // 每 5 分钟
```

### Memo

```go
starter.Memo(map[string]any{
    "requester": "user-123",
    "source":    "web-api",
})
```

### 搜索属性

```go
starter.TypedSearchAttributes(temporal.NewSearchAttributes(
    temporal.SearchAttributeKeyString("CustomStatus").ValueSet("pending"),
))
```

### 积极启动

```go
// 如果本地 worker 可用，请求立即执行（减少调度延迟）
starter.EnableEagerStart(true)
```

### 延迟启动

```go
// 延迟 5 分钟后才开始调度第一个工作流任务
starter.StartDelay(5 * time.Minute)
```

### 静态摘要和详情

```go
starter.StaticSummary("Order processing workflow")
starter.StaticDetails("Processing order #12345 for user-123")
```

### 版本化覆盖

```go
starter.VersioningOverride(client.VersioningOverride{
    Behavior: client.VersioningBehaviorPinned,
    PinnedVersion: "v1.2.3",
})
```

### 优先级

```go
starter.Priority(temporal.Priority{
    // 优先级配置
})
```

## 常见组合示例

### 1. 简单的同步调用

```go
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.WaitResult(true),
)
```

### 2. 带超时的同步调用

```go
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.WaitResult(true),
    starter.WorkflowExecutionTimeout(5 * time.Minute),
    starter.WorkflowRunTimeout(2 * time.Minute),
)
```

### 3. 定时调度

```go
_, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.CronSchedule("0 9 * * *"),
    starter.ID("daily-hello-job"),
)
```

### 4. 带重试的异步调用

```go
_, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.RetryPolicy(&temporal.RetryPolicy{
        MaximumAttempts: 3,
    }),
)
```

### 5. 带搜索属性和 Memo

```go
result, md, err := hc.Hello(ctx, &api.HelloRequest{Name: "World"},
    starter.WaitResult(true),
    starter.ID("hello-001"),
    starter.Memo(map[string]any{"source": "api"}),
    starter.TypedSearchAttributes(temporal.NewSearchAttributes(
        temporal.SearchAttributeKeyString("Status").ValueSet("running"),
    )),
)
```

## 选项组合规则

- 可以传入任意数量的选项，顺序无关
- 后传入的选项会覆盖先传入的同名选项
- 某些选项互斥或不常用，根据业务场景选择
