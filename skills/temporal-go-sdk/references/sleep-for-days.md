# Sleep for Days 长睡眠

长时间睡眠的工作流。

```go
// 长时间睡眠工作流
func LongSleepWorkflow(ctx workflow.Context, duration time.Duration) error {
    AO := workflow.ActivityOptions{StartToCloseTimeout: duration + time.Hour}
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 使用 Timer 实现长时间等待
    // 注意: Temporal 会持久化 Timer，不会因为 Worker 重启而丢失
    timerFut := workflow.NewTimer(ctx, duration)

    // 可以添加提前取消的逻辑
    cancelCh := workflow.GetSignalChannel(ctx, "cancel")

    selector := workflow.NewSelector(ctx)
    selector.AddFuture(timerFut, func(f workflow.Future) {
        // Timer 完成
    })
    selector.AddReceive(cancelCh, func(c workflow.ReceiveChannel, more bool) {
        // 收到取消信号
    })

    selector.Select(ctx)

    return nil
}

// 或者使用 Async 模式
func AsyncLongSleepWorkflow(ctx workflow.Context, targetTime time.Time) error {
    // 计算等待时间
    waitDuration := targetTime.Sub(time.Now())

    if waitDuration > 0 {
        // 等待直到指定时间
        workflow.Sleep(ctx, waitDuration)
    }

    // 执行后续逻辑
    return workflow.ExecuteActivity(ctx, ExecuteAtTimeActivity, nil).Get(ctx, nil)
}
```
