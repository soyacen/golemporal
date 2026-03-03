# Update 工作流更新

使用 Update API 实现请求-响应模式。

## 工作流端注册 Update

```go
func UpdateableWorkflow(ctx workflow.Context) (string, error) {
    currentState := "initial"

    // 注册更新处理程序
    workflow.SetUpdateHandler(ctx, "set-state", func(ctx workflow.Context, newState string) (string, error) {
        oldState := currentState
        currentState = newState
        return fmt.Sprintf("State changed from %s to %s", oldState, currentState), nil
    })

    // 带参数的更新
    workflow.SetUpdateHandler(ctx, "increment-counter", func(ctx workflow.Context, delta int) (int, error) {
        counter := 0
        counter += delta
        return counter, nil
    })

    // 等待
    workflow.Await(ctx, func() bool { return false })

    return currentState, nil
}
```

## 客户端发送 Update

```go
// 发送更新并等待结果
handle, err := c.UpdateWorkflow(context.Background(), "workflow-id", "", "set-state", "new-state")
if err != nil {
    return err
}

var result string
err = handle.Get(context.Background(), &result)
fmt.Println(result) // "State changed from initial to new-state"

// 非阻塞更新
handle, err = c.UpdateWorkflow(context.Background(), "workflow-id", "", "increment-counter", 5)
// 不等待结果
_ = handle
```

## 验证 Update

```go
// 更新处理程序可以返回验证错误
workflow.SetUpdateHandler(ctx, "update-name", func(ctx workflow.Context, name string) (string, error) {
    if len(name) < 3 {
        return "", fmt.Errorf("name too short: minimum 3 characters")
    }
    return "Name updated to: " + name, nil
})
```

## 最佳实践

1. **同步响应**: Update 适合需要同步响应的场景
2. **验证**: 可以在更新处理程序中返回验证错误
3. **超时设置**: 使用 Context 超时避免无限等待
4. **与 Signal 区分**: Update 有返回值，Signal 没有
