# Recovery 故障恢复

工作流故障恢复模式。

```go
// 可恢复的工作流
func RecoverableWorkflow(ctx workflow.Context, input Input) error {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
        RetryPolicy: &temporal.RetryPolicy{
            MaximumAttempts: 3,
            BackoffCoefficient: 2.0,
        },
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 使用 Selector 处理多种信号
    selector := workflow.NewSelector(ctx)

    var cancelSignal bool
    var progressSignal Progress

    selector.AddReceive(ctx, workflow.GetSignalChannel(ctx, "cancel"), func(c workflow.ReceiveChannel, more bool) {
        c.Receive(ctx, &cancelSignal)
    })

    selector.AddReceive(ctx, workflow.GetSignalChannel(ctx, "progress"), func(c workflow.ReceiveChannel, more bool) {
        c.Receive(ctx, &progressSignal)
    })

    // 主逻辑
    err := workflow.ExecuteActivity(ctx, RiskyActivity, input).Get(ctx, nil)
    if err != nil {
        // 恢复逻辑
        return workflow.ExecuteActivity(ctx, RecoverActivity, input).Get(ctx, nil)
    }

    return nil
}
```
