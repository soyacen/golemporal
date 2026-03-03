# Context Propagation 上下文传播

在 Workflow 和 Activity 之间传播上下文信息。

```go
// 带上下文的 Activity
func ActivityWithContext(ctx context.Context, request Request) error {
    // 从 context 中获取值
    traceID := ctx.Value("trace_id")
    userID := ctx.Value("user_id")

    // 使用上下文信息
    log.Printf("Processing request: trace=%v, user=%v", traceID, userID)
    return nil
}

// 在 Workflow 中设置上下文
func WorkflowWithContext(ctx workflow.Context, input Input) error {
    AO := workflow.ActivityOptions{
        StartToCloseTimeout: time.Minute,
    }
    ctx = workflow.WithActivityOptions(ctx, AO)

    // 使用 context.Context 传递
    actCtx := context.WithValue(context.Background(), "trace_id", input.TraceID)
    actCtx = context.WithValue(actCtx, "user_id", input.UserID)

    return workflow.ExecuteActivity(ctx, ActivityWithContext, actCtx, request).Get(ctx, nil)
}
```
