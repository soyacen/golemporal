# Updatable Timer 可更新计时器

可更新的计时器模式。

```go
// 可更新的计时器工作流
func UpdatableTimerWorkflow(ctx workflow.Context) (time.Time, error) {
    deadline := time.Now().Add(time.Hour)

    // 注册更新处理程序
    handler := workflow.UpdateWorkflow(ctx, "extendDeadline", func(newDuration time.Duration) time.Time {
        deadline = time.Now().Add(newDuration)
        return deadline
    })

    // 注册查询处理程序
    workflow.SetQueryHandler(ctx, "getDeadline", func() time.Time {
        return deadline
    })

    // 使用 Timer
    timer := workflow.NewTimer(ctx, time.Until(deadline))

    selector := workflow.NewSelector(ctx)
    selector.AddFuture(timer, func(f workflow.Future) {
        // Timer 完成
    })

    // 等待 Timer 或取消
    selector.Select(ctx)

    return deadline, nil
}

// 客户端更新 Timer
func ExtendTimer(client temporal.Client, workflowID string, newDuration time.Duration) error {
    handle := client.GetWorkflowHandle(workflowID)

    var newDeadline time.Time
    err := handle.Update(ctx, "extendDeadline", &newDeadline, newDuration)
    return err
}
```
