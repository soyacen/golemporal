# Nexus Cancellation Nexus 取消

取消 Nexus 操作。

```go
// 启动长时间运行的 Nexus 操作
func StartLongRunningOperation(ctx context.Context, client temporal.Client) (string, error) {
    handle, err := client.NexusStartOperation(ctx, nexus.StartOperationOptions{
        Service:   "my-service",
        Operation: "long-running",
        Request:   myRequest{},
    })

    if err != nil {
        return "", err
    }

    // 取消操作
    err = handle.Cancel(ctx)
    if err != nil {
        return "", err
    }

    // 等待取消完成
    var result myResponse
    err = handle.Get(ctx, &result)
    return result, err
}

// 处理 Nexus 取消请求
type LongRunningOperationHandler struct{}

func (h *LongRunningOperationHandler) HandleOperation(ctx context.Context, request myRequest) (myResponse, error) {
    // 检查取消信号
    select {
    case <-ctx.Done():
        return myResponse{}, ctx.Err()
    default:
        // 正常处理
    }

    // 执行长时间操作
    return processLongRunning(request)
}
```
