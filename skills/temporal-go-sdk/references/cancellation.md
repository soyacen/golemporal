# Cancellation 工作流取消

处理工作流取消和清理。

## 监听取消

```go
func CancellableWorkflow(ctx workflow.Context) error {
    // 创建可取消的 Activity
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: 5 * time.Minute,
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    err := workflow.ExecuteActivity(ctx, LongRunningActivity, input).Get(ctx, nil)
    if err != nil {
        // 处理错误
    }

    return nil
}

// 在 Activity 中监听取消
func LongRunningActivity(ctx context.Context, input Input) error {
    for {
        select {
        case <-ctx.Done():
            // 执行清理
            return ctx.Err()
        default:
            // 执行工作
        }
    }
}
```

## 取消处理程序

```go
func CancelableWorkflowWithCleanup(ctx workflow.Context) error {
    // 启动清理 Activity
    cleanupFuture := workflow.ExecuteActivity(ctx, CleanupActivity, cleanupData)

    selector := workflow.NewSelector(ctx)
    selector.AddFuture(cleanupFuture, func(f workflow.Future) {
        // 清理完成
    })
    selector.AddReceive(workflow.GetSignalChannel(ctx, "cancel"), func(c workflow.ReceiveChannel, more bool) {
        c.Receive(ctx, nil)
        // 收到取消信号
    })

    selector.Select(ctx)
    return nil
}
```

## 外部取消

```go
// 从客户端取消工作流
err := c.CancelWorkflow(context.Background(), "workflow-id", "")
```

## 最佳实践

1. **监听 ctx.Done()**: 始终检查 Context 取消
2. **清理 Activity**: 使用 defer 确保清理
3. **部分成功**: 允许部分完成的清理
