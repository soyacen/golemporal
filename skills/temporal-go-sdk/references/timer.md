# Timer 定时器

使用 Timer 实现延迟和超时。

## 创建定时器

```go
func TimerWorkflow(ctx workflow.Context, input Input) (string, error) {
    // 创建 5 秒定时器
    timer := workflow.NewTimer(ctx, 5*time.Second)

    // 等待定时器完成
    if err := timer.Get(ctx, nil); err != nil {
        return "", err
    }

    return "Timer completed after 5 seconds", nil
}
```

## 带条件的定时器

```go
func ConditionalTimerWorkflow(ctx workflow.Context) error {
    // 使用 Selector 等待多个条件之一
    selector := workflow.NewSelector(ctx)

    timer := workflow.NewTimer(ctx, 30*time.Second)
    signalChan := workflow.GetSignalChannel(ctx, "continue")

    var proceed bool
    selector.AddFuture(timer, func(f workflow.Future) {
        // 定时器完成
        proceed = false
    })
    selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
        // 收到信号
        c.Receive(ctx, nil)
        proceed = true
    })

    selector.Select(ctx)

    if proceed {
        // 继续执行
    } else {
        // 超时
    }
    return nil
}
```

## 可取消的定时器

```go
func CancellableTimerWorkflow(ctx workflow.Context) error {
    timer := workflow.NewTimer(ctx, 60*time.Second)

    select {
    case <-timer.GetChan():
        // 定时器完成
    case <-ctx.Done():
        // Context 取消（工作流取消）
        return ctx.Err()
    }
    return nil
}
```

## 更新定时器 (Updatable Timer)

参考: updatabletimer 示例

```go
func UpdatableTimerWorkflow(ctx workflow.Context) error {
    duration := 30 * time.Second
    timer := workflow.NewTimer(ctx, duration)

    selector := workflow.NewSelector(ctx)
    signalChan := workflow.GetSignalChannel(ctx, "update-duration")

    selector.AddFuture(timer, func(f workflow.Future) {
        // 定时器完成
    })
    selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
        var newDuration int
        c.Receive(ctx, &newDuration)
        duration = time.Duration(newDuration) * time.Second
        timer = workflow.NewTimer(ctx, duration)
    })

    selector.Select(ctx)
    return nil
}
```

## 最佳实践

1. **使用 Timer**: 永远使用 `workflow.NewTimer` 而非 `time.Sleep`
2. **Selector**: 使用 Selector 处理多个定时器或信号
3. **取消**: 监听 `ctx.Done()` 处理工作流取消
