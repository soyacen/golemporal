# Nexus Context Propagation Nexus 上下文传播

在 Nexus 调用中传播上下文。

```go
// 带上下文的 Nexus 调用
func CallWithContext(ctx context.Context, client temporal.Client) error {
    // 创建带上下文的 options
    options := nexus.StartOperationOptions{
        Service:   "my-service",
        Operation: "process",
        Request:   myRequest{TraceID: "trace-123"},
        // 传播上下文
        Context: context.WithValue(ctx, "user_id", "user-456"),
    }

    handle, err := client.NexusStartOperation(ctx, options)
    if err != nil {
        return err
    }

    var result myResponse
    return handle.Get(ctx, &result)
}

// Nexus 服务处理上下文
type MyNexusService struct{}

func (s *MyNexusService) HandleOperation(ctx context.Context, request myRequest) (myResponse, error) {
    // 提取传播的上下文
    userID := ctx.Value("user_id")
    traceID := request.TraceID

    log.Printf("Processing: user=%v, trace=%v", userID, traceID)

    return processWithContext(userID, traceID)
}
```
