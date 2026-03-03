# Await Signals 等待信号

处理乱序信号。

## 基本模式

```go
func AwaitSignalsWorkflow(ctx workflow.Context) ([]string, error) {
    signals := []string{}
    signalChan := workflow.GetSignalChannel(ctx, "input")

    // 持续等待信号直到满足条件
    workflow.Await(ctx, func() bool {
        return len(signals) >= 3
    })

    return signals, nil
}
```

## 带超时的 Await

```go
func AwaitWithTimeoutWorkflow(ctx workflow.Context) error {
    ok, _ := workflow.AwaitWithTimeout(ctx, 30*time.Second, func() bool {
        return false // 条件
    })

    if !ok {
        // 超时
    }
    return nil
}
```

## 处理乱序信号

```go
func OutOfOrderSignalsWorkflow(ctx workflow.Context) error {
    received := make(map[string]bool)
    signalChan := workflow.GetSignalChannel(ctx, "event")

    for len(received) < 3 {
        selector := workflow.NewSelector(ctx)
        selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
            var event string
            c.Receive(ctx, &event)
            received[event] = true
        })
        selector.Select(ctx)
    }

    return nil
}
```
