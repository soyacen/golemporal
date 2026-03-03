# Signals 异步通信

使用 Signal 实现工作流与外部的异步通信。

## 信号处理程序

```go
func SignalWorkflow(ctx workflow.Context) (string, error) {
    // 注册信号处理程序
    signalChan := workflow.GetSignalChannel(ctx, "my-signal")

    var state string
    selector := workflow.NewSelector(ctx)
    selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
        var payload string
        c.Receive(ctx, &payload)
        state = payload
        workflow.GetLogger(ctx).Info("Received signal", "payload", payload)
    })

    // 等待信号
    for i := 0; i < 5; i++ {
        selector.Select(ctx)
    }

    return state, nil
}
```

## 发送信号

```go
// 从客户端发送信号
err := c.SignalWorkflow(context.Background(), "workflow-id", "", "my-signal", "new-state")
```

## 带参数的工作流信号

```go
type MySignal struct {
    Event string
    Data  string
}

func SignalWorkflowWithParams(ctx workflow.Context) error {
    signalChan := workflow.GetSignalChannel(ctx, "my-signal")

    var state string
    selector := workflow.NewSelector(ctx)
    selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
        var signal MySignal
        c.Receive(ctx, &signal)
        state = signal.Data
    })

    selector.Select(ctx)
    return nil
}

// 发送
err := c.SignalWorkflow(context.Background(), "workflow-id", "", "my-signal", MySignal{
    Event: "update",
    Data:  "value",
})
```

## 工作流内发送信号

```go
func WorkflowWithChildSignal(ctx workflow.Context) error {
    // 启动子工作流
    childFuture := workflow.ExecuteChildWorkflow(ctx, ChildWorkflow, input)

    // 发送信号给子工作流
    childHandle := workflow.GetChildWorkflowFuture(ctx, childFuture)
    // ... 可以通过 childHandle 获取结果
    return nil
}
```

## 最佳实践

1. **使用 channel**: 使用 `GetSignalChannel` 处理信号
2. **Selector**: 使用 Selector 监听多个信号
3. **幂等性**: 设计可重复处理的信号
4. **避免阻塞**: 信号处理应该快速返回
