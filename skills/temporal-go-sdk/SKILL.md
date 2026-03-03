---
name: temporal-go-developer
description: Go Temporal SDK 开发专家。精通 Go Temporal SDK，能够帮助开发者实现工作流、Activity、并发控制、错误处理、观测性等所有特性。当用户提及 Temporal、workflow、activity、child workflow、saga、temporal cloud、worker、task queue 等关键词时必须使用此 skill。提供从基础到高级的完整开发指导，包括代码模板、最佳实践和常见问题解决方案。
compatibility: Go 1.21+, go.temporal.io/sdk, Temporal Server
---

# Temporal Go SDK Developer Skill

本 skill 提供全面的 Go Temporal SDK 开发指导，涵盖 [temporalio/samples-go](https://github.com/temporalio/samples-go) 中的所有特性。

## 核心概念快速导航

| 特性 | 参考文档 | 典型用例 |
|------|----------|----------|
| 基础工作流 | [helloworld](./references/helloworld.md) | 简单业务流程 |
| Activity 重试 | [retryactivity](./references/retryactivity.md) | 外部 API 调用 |
| 子工作流 | [child-workflow](./references/child-workflow.md) | 模块化业务流程 |
| 并发执行 | [goroutine](./references/goroutine.md) | 并行任务处理 |
| 定时器 | [timer](./references/timer.md) | 延迟任务 |
| Cron 调度 | [cron](./references/cron.md) | 定时任务 |
| 查询工作流 | [query](./references/query.md) | 获取工作流状态 |
| 信号处理 | [signals](./references/signals.md) | 异步通信 |
| 更新工作流 | [update](./references/update.md) | 请求-响应模式 |
| 数据加密 | [encryption](./references/encryption.md) | 敏感数据保护 |
| Search Attributes | [searchattributes](./references/searchattributes.md) | 工作流检索 |
| Session | [sessions](./references/sessions.md) | 资源池管理 |
| 观测性 | [observability](./references/observability.md) | 链路追踪 |
| 认证 | [authentication](./references/authentication.md) | Temporal Cloud 连接 |
| Worker 版本化 | [worker-versioning](./references/worker-versioning.md) | 无停机部署 |
| Nexus | [nexus](./references/nexus.md) | 跨命名空间调用 |
| 动态工作流 | [dynamic-workflows](./references/dynamic-workflows.md) | 运行时决策 |
| Saga 模式 | [saga](./references/saga.md) | 分布式事务 |
| 条件分支 | [branch](./references/branch.md) | 多分支逻辑 |
| 选择执行 | [pickfirst](./references/pickfirst.md) | 竞态处理 |
| 取消处理 | [cancellation](./references/cancellation.md) | 优雅关闭 |
| 异步完成 | [expense](./references/expense.md) | 长运行任务 |
| 请求响应模式 | [reqrespactivity](./references/reqrespactivity.md) | 回调模式 |
| 轮询模式 | [polling](./references/polling.md) | 外部资源等待 |
| 多条件分支 | [choice-multi](./references/choice-multi.md) | 复杂条件 |
| 上下文传播 | [ctxpropagation](./references/ctxpropagation.md) | 传递上下文 |
| DSL 工作流 | [dsl](./references/dsl.md) | 配置化流程 |
| 积极启动 | [eager-workflow-start](./references/eager-workflow-start.md) | 低延迟启动 |
| 提前返回 | [early-return](./references/early-return.md) | 快速路径 |
| 外部配置 | [external-env-conf](./references/external-env-conf.md) | 动态配置 |
| 本地问候 | [greetingslocal](./references/greetingslocal.md) | 本地 Activity |
| gRPC 代理 | [grpc-proxy](./references/grpc-proxy.md) | 代理转发 |
| API Key 认证 | [helloworld-apikey](./references/helloworld-apikey.md) | API 认证 |
| mTLS 认证 | [helloworldmtls](./references/helloworldmtls.md) | 双向认证 |
| 日志拦截器 | [logger-interceptor](./references/logger-interceptor.md) | 日志增强 |
| Memo | [memo](./references/memo.md) | 工作流元数据 |
| 多历史回放 | [multi-history-replay](./references/multi-history-replay.md) | 测试覆盖 |
| 互斥锁 | [mutex](./references/mutex.md) | 资源竞争 |
| Nexus 取消 | [nexus-cancelation](./references/nexus-cancelation.md) | 取消操作 |
| Nexus 上下文 | [nexus-context-propagation](./references/nexus-context-propagation.md) | 上下文传递 |
| Nexus 多参数 | [nexus-multiple-arguments](./references/nexus-multiple-arguments.md) | 复杂请求 |
| PSO 采样 | [pso](./references/pso.md) | 周期性采样 |
| 故障恢复 | [recovery](./references/recovery.md) | 错误恢复 |
| 请求响应查询 | [reqrespquery](./references/reqrespquery.md) | 查询模式 |
| 请求响应更新 | [reqrespupdate](./references/reqrespupdate.md) | 更新模式 |
| 安全消息 | [safe_message_handler](./references/safe_message_handler.md) | 消息验证 |
| 定时调度 | [schedule](./references/schedule.md) | 复杂调度 |
| JWT 认证 | [serverjwtauth](./references/serverjwtauth.md) | Token 认证 |
| Session 失败 | [session-failure](./references/session-failure.md) | 会话恢复 |
| 购物车 | [shoppingcart](./references/shoppingcart.md) | 状态管理 |
| 长睡眠 | [sleep-for-days](./references/sleep-for-days.md) | 长时间等待 |
| Slog 适配器 | [slogadapter](./references/slogadapter.md) | 结构化日志 |
| Snappy 压缩 | [snappycompress](./references/snappycompress.md) | 数据压缩 |
| 分离合并 Future | [splitmerge-future](./references/splitmerge-future.md) | 并行合并 |
| 分离合并选择器 | [splitmerge-selector](./references/splitmerge-selector.md) | 选择器合并 |
| 独立 Activity | [standalone-activity](./references/standalone-activity.md) | 独立执行 |
| 延迟启动 | [start-delay](./references/start-delay.md) | 延迟执行 |
| 同步代理 | [synchronous-proxy](./references/synchronous-proxy.md) | 代理转发 |
| 测试夹具 | [temporal-fixtures](./references/temporal-fixtures.md) | 单元测试 |
| 类型化搜索属性 | [typed-searchattributes](./references/typed-searchattributes.md) | 强类型检索 |
| 可更新计时器 | [updatabletimer](./references/updatabletimer.md) | 动态超时 |
| 特定任务队列 | [worker-specific-task-queues](./references/worker-specific-task-queues.md) | 任务分配 |
| 安全拦截器 | [workflow-security-interceptor](./references/workflow-security-interceptor.md) | 权限控制 |
| Zap 适配器 | [zapadapter](./references/zapadapter.md) | Zap 日志 |
| Datadog 集成 | [datadog](./references/datadog.md) | 监控集成 |
| 动态配置 | [dynamicconfig](./references/dynamicconfig.md) | 运行时配置 |
| 滑动窗口批处理 | [batch-sliding-window](./references/batch-sliding-window.md) | 批量处理 |

## 目录

1. [快速开始](#快速开始)
2. [工作流 (Workflow)](#工作流-workflow)
3. [Activity](#activity)
4. [并发控制](#并发控制)
5. [错误处理与重试](#错误处理与重试)
6. [观测性](#观测性)
7. [安全与认证](#安全与认证)
8. [高级模式](#高级模式)
9. [最佳实践](#最佳实践)

---

## 快速开始

### 安装依赖

```bash
go get go.temporal.io/sdk
go get go.temporal.io/api/serviceerror
```

### 最小示例

```go
package main

import (
    "context"
    "fmt"

    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/worker"
    "go.temporal.io/sdk/workflow"
)

// 定义 Activity
func SayHello(ctx context.Context, name string) (string, error) {
    return fmt.Sprintf("Hello, %s!", name), nil
}

// 定义 Workflow
func HelloWorldWorkflow(ctx workflow.Context, name string) (string, error) {
   AO := workflow.ActivityOptions{
        StartToCloseTimeout: 10 * time.Second,
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    var result string
    err := workflow.ExecuteActivity(ctx, SayHello, name).Get(ctx, &result)
    return result, err
}

func main() {
    c, err := client.Dial(client.Options{})
    if err != nil {
        log.Fatalln("Unable to create client", err)
    }
    defer c.Close()

    w := worker.New(c, "hello-task-queue", worker.Options{})
    w.RegisterWorkflow(HelloWorldWorkflow)
    w.RegisterActivity(SayHello)

    if err := w.Run(worker.InterruptCh()); err != nil {
        log.Fatalln("Worker stopped", err)
    }
}
```

### 启动工作流

```go
workflowOptions := client.StartWorkflowOptions{
    ID:        "my-workflow-id",
    TaskQueue: "hello-task-queue",
}

result, err := c.ExecuteWorkflow(context.Background(), workflowOptions, HelloWorldWorkflow, "World")
```

---

## 工作流 (Workflow)

### 工作流定义原则

1. **确定性**: 同样的输入必须产生同样的输出
2. **幂等性**: 可以安全地重放
3. **无副作用**: 不直接调用外部服务，使用 Activity

### 工作流类型

#### 1. 基础工作流

参考: [helloworld](./references/helloworld.md)

```go
func MyWorkflow(ctx workflow.Context, input MyInput) (MyOutput, error) {
    // 简单的顺序执行
    AO := workflow.ActivityOptions{StartToCloseTimeout: time.Minute}
    ctx = workflow.WithActivityOptions(ctx, AO)

    var result string
    err := workflow.ExecuteActivity(ctx, MyActivity, input).Get(ctx, &result)
    return result, err
}
```

#### 2. 动态工作流

参考: [dynamic-workflows](./references/dynamic-workflows.md)

```go
func DynamicWorkflow(ctx workflow.Context) (interface{}, error) {
    // 根据运行时条件动态选择执行哪个 Activity
    return nil, nil
}
```

#### 3. 可更新工作流

参考: [update](./references/update.md)

```go
func UpdateableWorkflow(ctx workflow.Context) (int, error) {
    // 注册更新处理程序
    workflow.SetUpdateHandler(ctx, "add", func(ctx workflow.Context, n int) (int, error) {
        // 处理更新请求
        return n, nil
    })
    // ... 等待信号或其他
}
```

### 工作流特性

| 特性 | 示例 | 参考 |
|------|------|------|
| 延迟 | `workflow.NewTimer(ctx, duration)` | timer |
| 条件分支 | `workflow.Await(ctx, func() bool {...})` | choice-exclusive |
| 并行执行 | `workflow.Go(ctx, func(ctx workflow.Context) {...})` | goroutine |
| 子工作流 | `workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, input)` | child-workflow |
| ContinueAsNew | `workflow.ContinueAsNew(ctx, input)` | child-workflow-continue-as-new |
| Cron 调度 | `workflow.NewCronSchedule("@daily")` | cron |

---

## Activity

### Activity 定义

Activity 是工作流中执行具体操作的函数，必须满足以下条件：

1. 函数签名: `func(ctx context.Context, input InputType) (OutputType, error)`
2. 必须是导出函数或方法
3. 必须是可序列化的参数和返回值

### Activity 选项

```go
AO := workflow.ActivityOptions{
    StartToCloseTimeout: 10 * time.Minute,  // Activity 执行超时
    ScheduleToStartTimeout: 5 * time.Minute, // 等待被调度超时
    ScheduleToCloseTimeout: 15 * time.Minute,// 调度到完成总超时
    HeartbeatTimeout: 30 * time.Second,     // 心跳超时（长运行 Activity）
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient: 2.0,
        MaximumInterval:    100 * time.Second,
        MaximumAttempts:    5,
    },
}
ctx = workflow.WithActivityOptions(ctx, AO)
```

### Activity 特性

| 特性 | 示例 | 参考 |
|------|------|------|
| 异步完成 | `activity.RecordHeartbeat(ctx, "progress")` | expense |
| 重试策略 | `RetryPolicy` 配置 | retryactivity |
| 心跳 | `activity.GetHeartbeatDetails(ctx)` | retryactivity |
| 本地 Activity | `workflow.WithLocalActivityOptions` | greetingslocal |
| 独立 Activity | 直接从客户端调用 | standalone-activity |

### Activity 最佳实践

1. **始终设置合理的超时**: 避免永久阻塞
2. **长运行 Activity 必须心跳**: 让 Temporal 知道 Activity 存活
3. **使用重试策略**: 处理临时性失败
4. **幂等性设计**: Activity 可能重试

---

## 并发控制

### Goroutine (并发执行)

参考: [goroutine](./references/goroutine.md)

```go
// 并行执行多个 Activity
futures := make([]workflow.Future, len(tasks))
for i, task := range tasks {
    futures[i] = workflow.ExecuteActivity(ctx, ProcessTask, task)
}

// 等待所有完成
for _, f := range futures {
    var result TaskResult
    f.Get(ctx, &result)
}
```

### Selector (选择执行)

参考: [pickfirst](./references/pickfirst.md)

```go
selector := workflow.NewSelector(ctx)
selector.AddFuture(future1, func(f workflow.Future) {
    // 处理 future1 完成
})
selector.AddFuture(future2, func(f workflow.Future) {
    // 处理 future2 完成
})
selector.Select(ctx) // 阻塞直到其中一个完成
```

### Mutex (互斥锁)

参考: [mutex](./references/mutex.md)

使用 `workflow.ExecuteActivity` 配合分布式锁实现资源互斥。

### Session (资源池)

参考: [fileprocessing](./references/sessions.md)

```go
sessionOptions := &workflow.SessionOptions{
    ExecutionTimeout: 10 * time.Minute,
    CreationTimeout:  10 * time.Minute,
}
session, err := workflow.CreateSession(ctx, sessionOptions)
defer session.Close()

// 所有 Activity 将在同一 Worker 上执行
```

---

## 错误处理与重试

### 重试策略

参考: [retryactivity](./references/retryactivity.md)

```go
AO := workflow.ActivityOptions{
    RetryPolicy: &temporal.RetryPolicy{
        InitialInterval:    time.Second,
        BackoffCoefficient:  2.0,
        MaximumInterval:     100 * time.Second,
        MaximumAttempts:     5,
        NonRetryableErrorTypes: []string{"InvalidArgument"},
    },
}
```

### 取消处理

参考: [cancellation](./references/cancellation.md)

```go
// 监听取消
select {
case <-ctx.Done():
    // 执行清理 Activity
    workflow.ExecuteActivity(ctx, Cleanup)
    return ctx.Err()
default:
    // 继续执行
}
```

### Saga 模式

参考: [saga](./references/saga.md)

```go
// 定义补偿 Activity
type Saga struct {
    Activities []Activity
}

func (s *Saga) Add(compensation Activity) {
    s.Activities = append(s.Activities, compensation)
}

func (s *Saga) Compensate(ctx context.Context) error {
    for i := len(s.Activities) - 1; i >= 0; i-- {
        if err := workflow.ExecuteActivity(ctx, s.Activities[i]).Get(ctx, nil); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 观测性

### OpenTelemetry

参考: [opentelemetry](./references/observability.md)

```go
// 启动带有 OpenTelemetry 的 Worker
tracer := oteltrace.NewTracer()
c, _ := client.Dial(client.Options{
    TracingInterceptorOptions: interceptor.TracingInterceptorOptions{
        Tracer: tracer,
    },
})
```

### Metrics

参考: [metrics](./references/metrics.md)

```go
// 配置 Prometheus Metrics
c, _ := client.Dial(client.Options{
    MetricsHandler: metrics.NewPrometheusHandler(),
})
```

### 日志

```go
logger := workflow.GetLogger(ctx)
logger.Info("workflow started", "input", input)
```

---

## 安全与认证

### mTLS 认证

参考: [helloworldmtls](./references/authentication.md)

```go
c, err := client.Dial(client.Options{
    HostPort:  "my-cluster.tmprl.cloud:7233",
    Namespace: "my-namespace",
    ConnectionOptions: client.ConnectionOptions{
        TLS: &tls.Config{
            CertFile: "cert.pem",
            KeyFile:  "key.pem",
            RootCAFile: "ca.pem",
        },
    },
})
```

### API Key 认证

参考: [helloworld-apiKey](./references/authentication.md)

```go
c, err := client.Dial(client.Options{
    APIKey: "my-api-key",
})
```

---

## 高级模式

### Child Workflow

参考: [child-workflow](./references/child-workflow.md)

```go
// 启动子工作流
childFuture := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, childInput)
var childResult string
childFuture.Get(ctx, &childResult)
```

### ContinueAsNew

参考: [child-workflow-continue-as-new.md](./references/child-workflow.md)

用于处理大数据量，将状态转移到新的工作流执行。

```go
// 在工作流结束时调用
workflow.ContinueAsNew(ctx, NewWorkflowInput{ /* ... */ })
```

### 更新 (Update)

参考: [update](./references/update.md)

```go
// Workflow 端
workflow.SetUpdateHandler(ctx, "update-name", handler)

// Client 端
handle, _ := c.UpdateWorkflow(ctx, "workflow-id", "update-name", updateInput)
```

### 查询 (Query)

参考: [query](./references/query.md)

```go
// Workflow 端
workflow.SetQueryHandler(ctx, "get-state", func() (string, error) {
    return currentState, nil
})

// Client 端
value, _ := c.QueryWorkflow(ctx, "workflow-id", "", "get-state")
```

### Cron Schedule

参考: [cron](./references/cron.md)

```go
options := client.StartWorkflowOptions{
    CronSchedule: "@daily", // 每小时: "@hourly", 每分钟: "@every 1m"
}
```

### 调度 (Schedule)

参考: [schedule](./references/schedule.md)

```go
schedule, _ := c.CreateSchedule(ctx, client.ScheduleOptions{
    Spec: client.ScheduleSpec{
        CronExpressions: []string{"@daily"},
    },
    Action: &client.ScheduleWorkflowAction{
        ID:        "scheduled-workflow",
        TaskQueue: "my-task-queue",
        Workflow:  MyWorkflow,
    },
})
```

---

## 最佳实践

### 1. 正确使用 Context

- Workflow Context: `workflow.Context` - 用于工作流内部
- Activity Context: `context.Context` - 用于 Activity 内部

### 2. 避免的常见错误

| 错误 | 解决方案 |
|------|----------|
| 使用 time.Sleep | 使用 `workflow.NewTimer` |
| 使用 goroutine | 使用 `workflow.Go` |
| 访问全局变量 | 通过参数传递 |
| 调用外部服务 | 通过 Activity |

### 3. 数据序列化

- 使用 PoD 类型或使用 `temporalio/dataConverter` 配置自定义编码器
- 参考 [encryption](./references/encryption.md) 保护敏感数据

### 4. Worker 配置

```go
w := worker.New(c, "task-queue", worker.Options{
    MaxConcurrentWorkflowTaskExecutionSize: 100,
    MaxConcurrentActivityExecutionSize:    50,
    MaxConcurrentLocalActivityExecutionSize: 20,
    StickyScheduleToStartTimeout:         10 * time.Second,
})
```

### 5. Search Attributes 使用

参考: [searchattributes](./references/searchattributes.md)

```go
// 启动时添加 Search Attributes
options := client.StartWorkflowOptions{
    SearchAttributes: map[string]interface{}{
        "CustomStatus": []string{"pending", "running"},
    },
}
```

---

## 常见问题

### Q: 如何调试工作流?

A: 使用 `workflow.GetLogger(ctx).Info()` 添加日志，或使用 Temporal UI 查看历史。

### Q: Activity 超时怎么办?

A: 检查网络、检查 Activity 是否正确注册、增加超时时间。

### Q: 工作流卡住怎么办?

A: 使用 `workflow.AwaitWithTimeout` 添加超时，检查是否有死循环。

### Q: 如何处理长运行工作流?

A: 使用 `ContinueAsNew` 定期创建新执行，或使用 Cron 调度。

---

## 相关资源

- [Temporal Go SDK 文档](https://docs.temporal.io/dev-guide/golang)
- [Samples Go 仓库](https://github.com/temporalio/samples-go)
- [Temporal CLI](https://docs.temporal.io/cli)
- [Temporal Cloud](https://temporal.io/cloud)
