# Selector 选择器模式

使用 Selector 实现多路复用等待。

## 基本 Selector

```go
selector := workflow.NewSelector(ctx)

timer := workflow.NewTimer(ctx, 30*time.Second)
signalChan := workflow.GetSignalChannel(ctx, "continue")

selector.AddFuture(timer, func(f workflow.Future) {
    // 定时器完成
})
selector.AddReceive(signalChan, func(c workflow.ReceiveChannel, more bool) {
    var payload string
    c.Receive(ctx, &payload)
    // 处理信号
})

// 阻塞直到其中一个条件满足
selector.Select(ctx)
```

## 多个 Future

```go
f1 := workflow.ExecuteActivity(ctx, Activity1)
f2 := workflow.ExecuteActivity(ctx, Activity2)

selector.AddFuture(f1, func(f workflow.Future) {
    var result1 string
    f.Get(ctx, &result1)
})
selector.AddFuture(f2, func(f workflow.Future) {
    var result2 string
    f.Get(ctx, &result2)
})

// 等待第一个完成
selector.Select(ctx)
```

## 循环选择

```go
for completed < total {
    selector.Select(ctx)
    completed++
}
```

## 最佳实践

1. **避免阻塞主线程**: 使用 selector 而非阻塞等待
2. **超时处理**: 组合使用定时器和信号
3. **优先级**: 先添加的先检查
